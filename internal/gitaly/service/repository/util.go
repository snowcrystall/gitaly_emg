package repository

import (
	"context"
	"fmt"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func (s *server) removeOriginInRepo(ctx context.Context, repository *gitalypb.Repository) error {
	cmd, err := s.gitCmdFactory.New(ctx, repository, git.SubCmd{Name: "remote", Args: []string{"remove", "origin"}}, git.WithRefTxHook(ctx, repository, s.cfg))

	if err != nil {
		return fmt.Errorf("remote cmd start: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("remote cmd wait: %v", err)
	}

	return nil
}
