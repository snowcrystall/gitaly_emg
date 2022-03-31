package commit

import (
	"context"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/log"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func (s *server) LastCommitForPath(ctx context.Context, in *gitalypb.LastCommitForPathRequest) (*gitalypb.LastCommitForPathResponse, error) {
	if err := validateLastCommitForPathRequest(in); err != nil {
		return nil, helper.ErrInvalidArgument(err)
	}

	resp, err := s.lastCommitForPath(ctx, in)
	if err != nil {
		return nil, helper.ErrInternal(err)
	}

	return resp, nil
}

func (s *server) lastCommitForPath(ctx context.Context, in *gitalypb.LastCommitForPathRequest) (*gitalypb.LastCommitForPathResponse, error) {
	path := string(in.GetPath())
	if len(path) == 0 || path == "/" {
		path = "."
	}

	repo := s.localrepo(in.GetRepository())
	c, err := s.catfileCache.BatchProcess(ctx, repo)
	if err != nil {
		return nil, err
	}

	options := in.GetGlobalOptions()

	// Preserve backwards compatibility with legacy LiteralPathspec
	// flag: These protobuf changes were not shipped in 13.1, so can be
	// remove before 13.2.
	if in.GetLiteralPathspec() {
		if options == nil {
			options = &gitalypb.GlobalOptions{}
		}

		options.LiteralPathspecs = true
	}

	commit, err := log.LastCommitForPath(ctx, s.gitCmdFactory, c, repo, git.Revision(in.GetRevision()), path, options)
	if log.IsNotFound(err) {
		return &gitalypb.LastCommitForPathResponse{}, nil
	}

	return &gitalypb.LastCommitForPathResponse{Commit: commit}, err
}

func validateLastCommitForPathRequest(in *gitalypb.LastCommitForPathRequest) error {
	if err := git.ValidateRevision(in.Revision); err != nil {
		return helper.ErrInvalidArgument(err)
	}
	return nil
}
