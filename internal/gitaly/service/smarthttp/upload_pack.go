package smarthttp

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/command"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/stats"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"gitlab.com/gitlab-org/gitaly/v14/streamio"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) PostUploadPack(stream gitalypb.SmartHTTPService_PostUploadPackServer) error {
	ctx := stream.Context()

	req, err := stream.Recv() // First request contains Repository only
	if err != nil {
		return err
	}

	if err := validateUploadPackRequest(req); err != nil {
		return err
	}

	h := sha1.New()

	stdinReader := io.TeeReader(streamio.NewReader(func() ([]byte, error) {
		resp, err := stream.Recv()

		return resp.GetData(), err
	}), h)

	stdin, collector := s.runStatsCollector(stream.Context(), stdinReader)
	defer collector.finish()

	var respBytes int64

	stdout := streamio.NewWriter(func(p []byte) error {
		respBytes += int64(len(p))
		return stream.Send(&gitalypb.PostUploadPackResponse{Data: p})
	})

	repoPath, err := s.locator.GetRepoPath(req.Repository)
	if err != nil {
		return err
	}

	git.WarnIfTooManyBitmaps(ctx, s.locator, req.GetRepository().GetStorageName(), repoPath)

	config, err := git.ConvertConfigOptions(req.GitConfigOptions)
	if err != nil {
		return err
	}

	commandOpts := []git.CmdOpt{
		git.WithStdin(stdin),
		git.WithStdout(stdout),
		git.WithGitProtocol(ctx, req),
		git.WithConfig(config...),
		git.WithPackObjectsHookEnv(ctx, req.Repository, s.cfg),
	}

	cmd, err := s.gitCmdFactory.NewWithoutRepo(ctx, git.SubCmd{
		Name:  "upload-pack",
		Flags: []git.Option{git.Flag{Name: "--stateless-rpc"}},
		Args:  []string{repoPath},
	}, commandOpts...)

	if err != nil {
		return status.Errorf(codes.Unavailable, "PostUploadPack: cmd: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		stats := collector.finish()

		if _, ok := command.ExitStatus(err); ok && stats.Deepen != "" {
			// We have seen a 'deepen' message in the request. It is expected that
			// git-upload-pack has a non-zero exit status: don't treat this as an
			// error.
			return nil
		}

		return status.Errorf(codes.Unavailable, "PostUploadPack: %v", err)
	}

	ctxlogrus.Extract(ctx).WithField("request_sha", fmt.Sprintf("%x", h.Sum(nil))).WithField("response_bytes", respBytes).Info("request details")

	return nil
}

func validateUploadPackRequest(req *gitalypb.PostUploadPackRequest) error {
	if req.Data != nil {
		return status.Errorf(codes.InvalidArgument, "PostUploadPack: non-empty Data")
	}

	return nil
}

type statsCollector struct {
	c       io.Closer
	statsCh chan stats.PackfileNegotiation
}

func (sc *statsCollector) finish() stats.PackfileNegotiation {
	sc.c.Close()
	return <-sc.statsCh
}

func (s *server) runStatsCollector(ctx context.Context, r io.Reader) (io.Reader, *statsCollector) {
	pr, pw := io.Pipe()
	sc := &statsCollector{
		c:       pw,
		statsCh: make(chan stats.PackfileNegotiation, 1),
	}

	go func() {
		defer close(sc.statsCh)

		stats, err := stats.ParsePackfileNegotiation(pr)
		if err != nil {
			ctxlogrus.Extract(ctx).WithError(err).Debug("failed parsing packfile negotiation")
			return
		}
		stats.UpdateMetrics(s.packfileNegotiationMetrics)

		sc.statsCh <- stats
	}()

	return io.TeeReader(r, pw), sc
}
