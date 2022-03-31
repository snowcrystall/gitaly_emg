package repository

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/remoterepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func (s *server) FetchSourceBranch(ctx context.Context, req *gitalypb.FetchSourceBranchRequest) (*gitalypb.FetchSourceBranchResponse, error) {
	if err := git.ValidateRevision(req.GetSourceBranch()); err != nil {
		return nil, helper.ErrInvalidArgument(err)
	}

	if err := git.ValidateRevision(req.GetTargetRef()); err != nil {
		return nil, helper.ErrInvalidArgument(err)
	}

	targetRepo := s.localrepo(req.GetRepository())

	sourceRepo, err := remoterepo.New(ctx, req.GetSourceRepository(), s.conns)
	if err != nil {
		return nil, helper.ErrInternal(err)
	}

	var sourceOid git.ObjectID
	var containsObject bool

	// If source and target repository are the same, then we know that both
	// are local. We can thus optimize and locally resolve the reference
	// instead of using an RPC call. We also know that we can always skip
	// the fetch as the object will be available.
	if helper.RepoPathEqual(req.GetRepository(), req.GetSourceRepository()) {
		var err error

		sourceOid, err = targetRepo.ResolveRevision(ctx, git.Revision(req.GetSourceBranch()))
		if err != nil {
			if errors.Is(err, git.ErrReferenceNotFound) {
				return &gitalypb.FetchSourceBranchResponse{Result: false}, nil
			}
			return nil, helper.ErrInternal(err)
		}

		containsObject = true
	} else {
		var err error

		sourceOid, err = sourceRepo.ResolveRevision(ctx, git.Revision(req.GetSourceBranch()))
		if err != nil {
			if errors.Is(err, git.ErrReferenceNotFound) {
				return &gitalypb.FetchSourceBranchResponse{Result: false}, nil
			}
			return nil, helper.ErrInternal(err)
		}

		// Otherwise, if the source is a remote repository, we check
		// whether the target repo already contains the desired object.
		// If so, we can skip the fetch.
		containsObject, err = targetRepo.HasRevision(ctx, sourceOid.Revision()+"^{commit}")
		if err != nil {
			return nil, helper.ErrInternal(err)
		}
	}

	// There's no need to perform the fetch if we already have the object
	// available.
	if !containsObject {
		if err := targetRepo.FetchInternal(
			ctx,
			req.GetSourceRepository(),
			[]string{sourceOid.String()},
			localrepo.FetchOpts{Tags: localrepo.FetchOptsTagsNone},
		); err != nil {
			// Design quirk: if the fetch failse, this RPC returns Result: false, but no error.
			if errors.As(err, &localrepo.ErrFetchFailed{}) {
				ctxlogrus.Extract(ctx).
					WithField("oid", sourceOid.String()).
					WithError(err).Warn("git fetch failed")
				return &gitalypb.FetchSourceBranchResponse{Result: false}, nil
			}

			return nil, err
		}
	}

	if err := targetRepo.UpdateRef(ctx, git.ReferenceName(req.GetTargetRef()), sourceOid, ""); err != nil {
		return nil, err
	}

	return &gitalypb.FetchSourceBranchResponse{Result: true}, nil
}
