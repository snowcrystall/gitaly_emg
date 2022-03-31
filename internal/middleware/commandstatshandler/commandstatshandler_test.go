package commandstatshandler

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"

	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcmwlogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/backchannel"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/ref"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/transaction"
	"gitlab.com/gitlab-org/gitaly/v14/internal/log"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createNewServer(t *testing.T, cfg config.Cfg) *grpc.Server {
	logger := testhelper.NewTestLogger(t)
	logrusEntry := logrus.NewEntry(logger).WithField("test", t.Name())

	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpcmw.ChainStreamServer(
			StreamInterceptor,
			grpcmwlogrus.StreamServerInterceptor(logrusEntry,
				grpcmwlogrus.WithTimestampFormat(log.LogTimestampFormat),
				grpcmwlogrus.WithMessageProducer(CommandStatsMessageProducer)),
		)),
		grpc.UnaryInterceptor(grpcmw.ChainUnaryServer(
			UnaryInterceptor,
			grpcmwlogrus.UnaryServerInterceptor(logrusEntry,
				grpcmwlogrus.WithTimestampFormat(log.LogTimestampFormat),
				grpcmwlogrus.WithMessageProducer(CommandStatsMessageProducer)),
		)),
	}

	server := grpc.NewServer(opts...)

	gitCommandFactory := git.NewExecCommandFactory(cfg)

	gitalypb.RegisterRefServiceServer(server, ref.NewServer(
		cfg,
		config.NewLocator(cfg),
		gitCommandFactory,
		transaction.NewManager(cfg, backchannel.NewRegistry()),
		catfile.NewCache(cfg),
	))

	return server
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestInterceptor(t *testing.T) {
	cleanup := testhelper.Configure()
	defer cleanup()

	cfg, repo, _ := testcfg.BuildWithRepo(t)

	logBuffer := &bytes.Buffer{}
	testhelper.NewTestLogger = func(tb testing.TB) *logrus.Logger {
		return &logrus.Logger{Out: logBuffer, Formatter: &logrus.JSONFormatter{}, Level: logrus.InfoLevel}
	}

	s := createNewServer(t, cfg)
	defer s.Stop()

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	go func() {
		err := s.Serve(listener)
		require.NoError(t, err)
	}()

	tests := []struct {
		name        string
		performRPC  func(ctx context.Context, client gitalypb.RefServiceClient)
		expectedLog string
	}{
		{
			name: "Unary",
			performRPC: func(ctx context.Context, client gitalypb.RefServiceClient) {
				req := &gitalypb.RefExistsRequest{Repository: repo, Ref: []byte("refs/foo")}

				_, err := client.RefExists(ctx, req)
				require.NoError(t, err)
			},
			expectedLog: "\"command.count\":1",
		},
		{
			name: "Stream",
			performRPC: func(ctx context.Context, client gitalypb.RefServiceClient) {
				req := &gitalypb.FindAllBranchNamesRequest{Repository: repo}

				stream, err := client.FindAllBranchNames(ctx, req)
				require.NoError(t, err)

				for {
					_, err := stream.Recv()
					if err == io.EOF {
						break
					}
					require.NoError(t, err)
				}
			},
			expectedLog: "\"command.count\":1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logBuffer.Reset()

			ctx, cancel := testhelper.Context()
			defer cancel()

			conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
			require.NoError(t, err)
			defer conn.Close()

			client := gitalypb.NewRefServiceClient(conn)

			tt.performRPC(ctx, client)
			require.Contains(t, logBuffer.String(), tt.expectedLog)
		})
	}
}
