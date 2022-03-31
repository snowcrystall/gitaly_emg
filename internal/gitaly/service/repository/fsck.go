package repository

import (
	"bytes"
	"context"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func (s *server) Fsck(ctx context.Context, req *gitalypb.FsckRequest) (*gitalypb.FsckResponse, error) {
	var stdout, stderr bytes.Buffer

	repo := req.GetRepository()

	cmd, err := s.gitCmdFactory.New(ctx, repo,
		git.SubCmd{Name: "fsck"},
		git.WithStdout(&stdout),
		git.WithStderr(&stderr),
	)
	if err != nil {
		return nil, err
	}

	if err = cmd.Wait(); err != nil {
		return &gitalypb.FsckResponse{Error: append(stdout.Bytes(), stderr.Bytes()...)}, nil
	}

	return &gitalypb.FsckResponse{}, nil
}
