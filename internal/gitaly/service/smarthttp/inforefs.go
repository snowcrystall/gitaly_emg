package smarthttp

import (
	"context"
	"fmt"
	"io"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/pktline"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"gitlab.com/gitlab-org/gitaly/v14/streamio"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	uploadPackSvc  = "upload-pack"
	receivePackSvc = "receive-pack"
)

func (s *server) InfoRefsUploadPack(in *gitalypb.InfoRefsRequest, stream gitalypb.SmartHTTPService_InfoRefsUploadPackServer) error {
	repoPath, err := s.locator.GetRepoPath(in.GetRepository())
	if err != nil {
		return err
	}

	w := streamio.NewWriter(func(p []byte) error {
		return stream.Send(&gitalypb.InfoRefsResponse{Data: p})
	})

	return s.infoRefCache.tryCache(stream.Context(), in, w, func(w io.Writer) error {
		return s.handleInfoRefs(stream.Context(), uploadPackSvc, repoPath, in, w)
	})
}

func (s *server) InfoRefsReceivePack(in *gitalypb.InfoRefsRequest, stream gitalypb.SmartHTTPService_InfoRefsReceivePackServer) error {
	repoPath, err := s.locator.GetRepoPath(in.GetRepository())
	if err != nil {
		return err
	}
	w := streamio.NewWriter(func(p []byte) error {
		return stream.Send(&gitalypb.InfoRefsResponse{Data: p})
	})
	return s.handleInfoRefs(stream.Context(), receivePackSvc, repoPath, in, w)
}

func (s *server) handleInfoRefs(ctx context.Context, service, repoPath string, req *gitalypb.InfoRefsRequest, w io.Writer) error {
	ctxlogrus.Extract(ctx).WithFields(log.Fields{
		"service": service,
	}).Debug("handleInfoRefs")

	cmdOpts := []git.CmdOpt{git.WithGitProtocol(ctx, req)}
	if service == "receive-pack" {
		cmdOpts = append(cmdOpts, git.WithRefTxHook(ctx, req.Repository, s.cfg))
	}

	config, err := git.ConvertConfigOptions(req.GitConfigOptions)
	if err != nil {
		return err
	}
	cmdOpts = append(cmdOpts, git.WithConfig(config...))

	cmd, err := s.gitCmdFactory.NewWithoutRepo(ctx, git.SubCmd{
		Name:  service,
		Flags: []git.Option{git.Flag{Name: "--stateless-rpc"}, git.Flag{Name: "--advertise-refs"}},
		Args:  []string{repoPath},
	}, cmdOpts...)

	if err != nil {
		if _, ok := status.FromError(err); ok {
			return err
		}
		return status.Errorf(codes.Internal, "GetInfoRefs: cmd: %v", err)
	}

	if _, err := pktline.WriteString(w, fmt.Sprintf("# service=git-%s\n", service)); err != nil {
		return status.Errorf(codes.Internal, "GetInfoRefs: pktLine: %v", err)
	}

	if err := pktline.WriteFlush(w); err != nil {
		return status.Errorf(codes.Internal, "GetInfoRefs: pktFlush: %v", err)
	}

	if _, err := io.Copy(w, cmd); err != nil {
		return status.Errorf(codes.Internal, "GetInfoRefs: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return status.Errorf(codes.Internal, "GetInfoRefs: %v", err)
	}

	return nil
}
