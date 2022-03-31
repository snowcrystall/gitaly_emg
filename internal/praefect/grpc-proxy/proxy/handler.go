// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// TODO: remove the following linter override when the deprecations are fixed
// in issue https://gitlab.com/gitlab-org/gitaly/issues/1663
//lint:file-ignore SA1019 Ignore all gRPC deprecations until issue #1663

package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"

	"gitlab.com/gitlab-org/gitaly/v14/internal/middleware/sentryhandler"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

// RegisterStreamHandlers sets up stream handlers for a set of gRPC methods for a given service.
// streamers is a map of method to grpc.StreamHandler eg:
//
// streamHandler := func(srv interface{}, stream ServerStream) error  {
//                       /** do some stuff **/
//                       return nil
//                  }
// RegisterStreamHandlers(grpcServer, "MyGrpcService", map[string]grpc.StreamHandler{"Method1": streamHandler})
// note: multiple calls with the same serviceName will result in a fatal
func RegisterStreamHandlers(server *grpc.Server, serviceName string, streamers map[string]grpc.StreamHandler) {
	desc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
	}

	for methodName, streamer := range streamers {
		streamDesc := grpc.StreamDesc{
			StreamName:    methodName,
			Handler:       streamer,
			ServerStreams: true,
			ClientStreams: true,
		}
		desc.Streams = append(desc.Streams, streamDesc)
	}

	server.RegisterService(desc, struct{}{})
}

// RegisterService sets up a proxy handler for a particular gRPC service and method.
// The behaviour is the same as if you were registering a handler method, e.g. from a codegenerated pb.go file.
//
// This can *only* be used if the `server` also uses grpcproxy.CodecForServer() ServerOption.
func RegisterService(server *grpc.Server, director StreamDirector, serviceName string, methodNames ...string) {
	streamer := &handler{director}
	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
	}
	for _, m := range methodNames {
		streamDesc := grpc.StreamDesc{
			StreamName:    m,
			Handler:       streamer.handler,
			ServerStreams: true,
			ClientStreams: true,
		}
		fakeDesc.Streams = append(fakeDesc.Streams, streamDesc)
	}
	server.RegisterService(fakeDesc, streamer)
}

// TransparentHandler returns a handler that attempts to proxy all requests that are not registered in the server.
// The indented use here is as a transparent proxy, where the server doesn't know about the services implemented by the
// backends. It should be used as a `grpc.UnknownServiceHandler`.
//
// This can *only* be used if the `server` also uses grpcproxy.CodecForServer() ServerOption.
func TransparentHandler(director StreamDirector) grpc.StreamHandler {
	streamer := &handler{director}
	return streamer.handler
}

type handler struct {
	director StreamDirector
}

type streamAndDestination struct {
	grpc.ClientStream
	destination Destination
	cancel      func()
}

// failDestinationWithErrors marks all of the destinations in the StreamParameters as
// having failed with the given error.
func failDestinationsWithError(params *StreamParameters, err error) {
	if params.Primary().ErrHandler != nil {
		_ = params.Primary().ErrHandler(err)
	}

	for _, secondary := range params.Secondaries() {
		if secondary.ErrHandler != nil {
			_ = secondary.ErrHandler(err)
		}
	}
}

// handler is where the real magic of proxying happens.
// It is invoked like any gRPC server stream and uses the gRPC server framing to get and receive bytes from the wire,
// forwarding it to a ClientStream established against the relevant ClientConn.
func (s *handler) handler(srv interface{}, serverStream grpc.ServerStream) (finalErr error) {
	// little bit of gRPC internals never hurt anyone
	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}

	peeker := newPeeker(serverStream)

	// We require that the director's returned context inherits from the serverStream.Context().
	params, err := s.director(serverStream.Context(), fullMethodName, peeker)
	if err != nil {
		return err
	}

	defer func() {
		err := params.RequestFinalizer()
		if finalErr == nil {
			finalErr = err
		}
	}()

	clientCtx, clientCancel := context.WithCancel(params.Primary().Ctx)
	defer clientCancel()
	// TODO(mwitkow): Add a `forwarded` header to metadata, https://en.wikipedia.org/wiki/X-Forwarded-For.

	primaryClientStream, err := grpc.NewClientStream(clientCtx, clientStreamDescForProxying, params.Primary().Conn, fullMethodName, params.CallOptions()...)
	if err != nil {
		failDestinationsWithError(params, fmt.Errorf("initiate primary stream: %w", err))
		return err
	}

	primaryStream := streamAndDestination{
		ClientStream: primaryClientStream,
		destination:  params.Primary(),
		cancel:       clientCancel,
	}

	var secondaryStreams []streamAndDestination
	for _, destination := range params.Secondaries() {
		clientCtx, clientCancel := context.WithCancel(destination.Ctx)
		defer clientCancel()

		secondaryClientStream, err := grpc.NewClientStream(clientCtx, clientStreamDescForProxying, destination.Conn, fullMethodName, params.CallOptions()...)
		if err != nil {
			failDestinationsWithError(params, fmt.Errorf("initiate secondary stream: %w", err))
			return err
		}
		secondaryStreams = append(secondaryStreams, streamAndDestination{
			ClientStream: secondaryClientStream,
			destination:  destination,
			cancel:       clientCancel,
		})
	}

	allStreams := append(secondaryStreams, primaryStream)

	// Explicitly *do not close* p2cErrChan and c2sErrChan, otherwise the select below will not terminate.
	// Channels do not have to be closed, it is just a control flow mechanism, see
	// https://groups.google.com/forum/#!msg/golang-nuts/pZwdYRGxCIk/qpbHxRRPJdUJ
	c2sErrChan := forwardClientToServers(serverStream, allStreams)
	p2cErrChan := forwardPrimaryToClient(primaryStream, serverStream)
	secondaryErrChan := receiveSecondaryStreams(secondaryStreams)

	// We need to wait for the streams from the primary and secondaries. However, we don't need
	// to wait for the p2c stream to finish because in some cases the client will close the
	// stream, signaling the end of a request/response cycle. If we wait for p2cErrChan, we are
	// effectively forcing the client to close the stream.
	var primaryProxied, secondariesProxied bool
	var clientErr, secondaryErr, primaryErr error

	for !(primaryProxied && secondariesProxied) {
		select {
		case clientErr = <-c2sErrChan:
			if clientErr != nil {
				// We may have gotten a receive error (stream disconnected, a read
				// error etc) in which case we need to cancel all streams. In order
				// to not exit before all Goroutines have concluded, we do not
				// return directly but instead return the error after the loop.
				cancelStreams(allStreams)
			}
		case secondaryErr = <-secondaryErrChan:
			secondariesProxied = true
		case primaryErr = <-p2cErrChan:
			// This happens when the clientStream has nothing else to offer (io.EOF) or
			// returned a gRPC error. In those two cases we may have received Trailers
			// as part of the call. In case of other errors (stream closed) the trailers
			// will be nil.
			trailer := primaryClientStream.Trailer()
			serverStream.SetTrailer(trailer)

			// clientErr will contain RPC error from client code. If not io.EOF, return
			// the RPC error as server stream error after the loop.
			if primaryErr != io.EOF {
				// If there is an error from the primary, we want to cancel all streams
				cancelStreams(allStreams)
			}

			primaryProxied = true
		}
	}

	if clientErr != nil {
		return status.Errorf(codes.Internal, "failed proxying c2s: %v", clientErr)
	}

	if primaryErr != nil && primaryErr != io.EOF {
		// we must not propagate Gitaly errors into Sentry
		sentryhandler.MarkToSkip(serverStream.Context())
		return primaryErr
	}

	if secondaryErr != nil {
		return status.Errorf(codes.Internal, "failed proxying to secondary: %v", secondaryErr)
	}

	return nil
}

func cancelStreams(streams []streamAndDestination) {
	for _, stream := range streams {
		stream.cancel()
	}
}

// receiveSecondaryStreams reads from the client streams of the secondaries and drops the message
// but returns an error to the channel if it encounters a non io.EOF error
func receiveSecondaryStreams(srcs []streamAndDestination) chan error {
	ret := make(chan error, 1)

	go func() {
		var g errgroup.Group
		for _, src := range srcs {
			src := src // rescoping for goroutine
			g.Go(func() error {
				for {
					if err := src.RecvMsg(&frame{}); err != nil {
						if errors.Is(err, io.EOF) {
							return nil
						}

						if src.destination.ErrHandler != nil {
							err = src.destination.ErrHandler(err)
						}

						src.cancel()
						return err
					}
				}
			})
		}

		ret <- g.Wait()
	}()
	return ret
}

func forwardPrimaryToClient(src streamAndDestination, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)

	go func() {
		var outerErr error
		f := &frame{}

		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				outerErr = err // this can be io.EOF which is happy case
				break
			}
			if i == 0 {
				// This is a bit of a hack, but client to server headers are only readable after first client msg is
				// received but must be written to server stream before the first msg is flushed.
				// This is the only place to do it nicely.
				md, err := src.Header()
				if err != nil {
					outerErr = err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					outerErr = err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				outerErr = err
				break
			}
		}

		if outerErr != nil && src.destination.ErrHandler != nil && !errors.Is(outerErr, io.EOF) {
			outerErr = src.destination.ErrHandler(outerErr)
		}

		ret <- outerErr
	}()

	return ret
}

func forwardFramesToServer(dst grpc.ClientStream, frameChan <-chan *frame) error {
	for f := range frameChan {
		if err := dst.SendMsg(f); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}

	// all messages redirected
	return dst.CloseSend()
}

func forwardClientToServers(src grpc.ServerStream, dsts []streamAndDestination) chan error {
	ret := make(chan error, 1)
	go func() {
		var g errgroup.Group

		frameChans := make([]chan<- *frame, 0, len(dsts))

		for _, dst := range dsts {
			dst := dst
			frameChan := make(chan *frame, 16)
			frameChan <- &frame{payload: dst.destination.Msg} // send re-written message
			frameChans = append(frameChans, frameChan)

			g.Go(func() error { return forwardFramesToServer(dst, frameChan) })
		}

		var outerErr error
		for {
			if err := receiveFromClientAndForward(src, frameChans); err != nil {
				if !errors.Is(err, io.EOF) {
					outerErr = err
				}

				break
			}
		}

		if err := g.Wait(); outerErr == nil {
			outerErr = err
		}

		ret <- outerErr
	}()
	return ret
}

func receiveFromClientAndForward(src grpc.ServerStream, frameChans []chan<- *frame) error {
	f := &frame{}

	if err := src.RecvMsg(f); err != nil {
		for _, frameChan := range frameChans {
			// signal no more data to redirect
			close(frameChan)
		}

		return err // this can be io.EOF which is happy case
	}

	for _, frameChan := range frameChans {
		frameChan <- f
	}

	return nil
}
