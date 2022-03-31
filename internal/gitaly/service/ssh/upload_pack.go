package ssh

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/command"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/pktline"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/stats"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"gitlab.com/gitlab-org/gitaly/v14/streamio"
)

func (s *server) SSHUploadPack(stream gitalypb.SSHService_SSHUploadPackServer) error {
	req, err := stream.Recv() // First request contains Repository only
	if err != nil {
		return helper.ErrInternal(err)
	}

	repository := ""
	if req.Repository != nil {
		repository = req.Repository.GlRepository
	}

	ctxlogrus.Extract(stream.Context()).WithFields(log.Fields{
		"GlRepository":     repository,
		"GitConfigOptions": req.GitConfigOptions,
		"GitProtocol":      req.GitProtocol,
	}).Debug("SSHUploadPack")

	if err = validateFirstUploadPackRequest(req); err != nil {
		return helper.ErrInvalidArgument(err)
	}

	if err = s.sshUploadPack(stream, req); err != nil {
		return helper.ErrInternal(err)
	}

	return nil
}

func (s *server) sshUploadPack(stream gitalypb.SSHService_SSHUploadPackServer, req *gitalypb.SSHUploadPackRequest) error {
	ctx, cancelCtx := context.WithCancel(stream.Context())
	defer cancelCtx()

	stdin := streamio.NewReader(func() ([]byte, error) {
		request, err := stream.Recv()
		return request.GetStdin(), err
	})

	// gRPC doesn't allow concurrent writes to a stream, so we need to
	// synchronize writing stdout and stderrr.
	var m sync.Mutex

	stdout := streamio.NewSyncWriter(&m, func(p []byte) error {
		return stream.Send(&gitalypb.SSHUploadPackResponse{Stdout: p})
	})

	stderr := streamio.NewSyncWriter(&m, func(p []byte) error {
		return stream.Send(&gitalypb.SSHUploadPackResponse{Stderr: p})
	})

	repoPath, err := s.locator.GetRepoPath(req.Repository)
	if err != nil {
		return err
	}

	git.WarnIfTooManyBitmaps(ctx, s.locator, req.GetRepository().StorageName, repoPath)

	config, err := git.ConvertConfigOptions(req.GitConfigOptions)
	if err != nil {
		return err
	}

	pr, pw := io.Pipe()
	defer pw.Close()
	stdin = io.TeeReader(stdin, pw)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			pr.Close()
		}()

		stats, err := stats.ParsePackfileNegotiation(pr)
		if err != nil {
			ctxlogrus.Extract(stream.Context()).WithError(err).Debug("failed parsing packfile negotiation")
			return
		}
		stats.UpdateMetrics(s.packfileNegotiationMetrics)
	}()

	commandOpts := []git.CmdOpt{
		git.WithGitProtocol(ctx, req),
		git.WithConfig(config...),
		git.WithPackObjectsHookEnv(ctx, req.Repository, s.cfg),
	}

	cmd, monitor, err := monitorStdinCommand(ctx, s.gitCmdFactory, stdin, stdout, stderr, git.SubCmd{
		Name: "upload-pack",
		Args: []string{repoPath},
	}, commandOpts...)

	if err != nil {
		return err
	}

	// upload-pack negotiation is terminated by either a flush, or the "done"
	// packet: https://github.com/git/git/blob/v2.20.0/Documentation/technical/pack-protocol.txt#L335
	//
	// "flush" tells the server it can terminate, while "done" tells it to start
	// generating a packfile. Add a timeout to the second case to mitigate
	// use-after-check attacks.
	go monitor.Monitor(pktline.PktDone(), s.uploadPackRequestTimeout, cancelCtx)

	if err := cmd.Wait(); err != nil {
		pw.Close()
		wg.Wait()

		if status, ok := command.ExitStatus(err); ok {
			return stream.Send(&gitalypb.SSHUploadPackResponse{
				ExitStatus: &gitalypb.ExitStatus{Value: int32(status)},
			})
		}
		return fmt.Errorf("cmd wait: %v", err)
	}

	pw.Close()
	wg.Wait()

	return nil
}

func validateFirstUploadPackRequest(req *gitalypb.SSHUploadPackRequest) error {
	if req.Stdin != nil {
		return fmt.Errorf("non-empty stdin in first request")
	}

	return nil
}
