package operations

import (
	"context"
	"errors"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/updateref"
	"gitlab.com/gitlab-org/gitaly/v14/internal/metadata/featureflag"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UserCreateBranch(ctx context.Context, req *gitalypb.UserCreateBranchRequest) (*gitalypb.UserCreateBranchResponse, error) {
	if len(req.BranchName) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Bad Request (empty branch name)")
	}

	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty user")
	}

	if len(req.StartPoint) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty start point")
	}

	quarantineDir, quarantineRepo, err := s.quarantinedRepo(ctx, req.GetRepository(), featureflag.Quarantine)
	if err != nil {
		return nil, err
	}

	// BEGIN TODO: Uncomment if StartPoint started behaving sensibly
	// like BranchName. See
	// https://gitlab.com/gitlab-org/gitaly/-/issues/3331
	//
	// startPointReference, err := s.localrepo(req.GetRepository()).GetReference(ctx, "refs/heads/"+string(req.StartPoint))
	// startPointCommit, err := log.GetCommit(ctx, req.Repository, startPointReference.Target)
	startPointCommit, err := quarantineRepo.ReadCommit(ctx, git.Revision(req.StartPoint))
	// END TODO
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "revspec '%s' not found", req.StartPoint)
	}

	startPointOID, err := git.NewObjectIDFromHex(startPointCommit.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not parse start point commit ID: %v", err)
	}

	referenceName := git.NewReferenceNameFromBranchName(string(req.BranchName))
	_, err = quarantineRepo.GetReference(ctx, referenceName)
	if err == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Could not update %s. Please refresh and try again.", req.BranchName)
	} else if !errors.Is(err, git.ErrReferenceNotFound) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := s.updateReferenceWithHooks(ctx, req.GetRepository(), req.User, quarantineDir, referenceName, startPointOID, git.ZeroOID); err != nil {
		var preReceiveError updateref.PreReceiveError
		if errors.As(err, &preReceiveError) {
			return &gitalypb.UserCreateBranchResponse{
				PreReceiveError: preReceiveError.Message,
			}, nil
		}

		var updateRefError updateref.Error
		if errors.As(err, &updateRefError) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return nil, err
	}

	return &gitalypb.UserCreateBranchResponse{
		Branch: &gitalypb.Branch{
			Name:         req.BranchName,
			TargetCommit: startPointCommit,
		},
	}, nil
}

func validateUserUpdateBranchGo(req *gitalypb.UserUpdateBranchRequest) error {
	if req.User == nil {
		return status.Errorf(codes.InvalidArgument, "empty user")
	}

	if len(req.BranchName) == 0 {
		return status.Errorf(codes.InvalidArgument, "empty branch name")
	}

	if len(req.Oldrev) == 0 {
		return status.Errorf(codes.InvalidArgument, "empty oldrev")
	}

	if len(req.Newrev) == 0 {
		return status.Errorf(codes.InvalidArgument, "empty newrev")
	}

	return nil
}

func (s *Server) UserUpdateBranch(ctx context.Context, req *gitalypb.UserUpdateBranchRequest) (*gitalypb.UserUpdateBranchResponse, error) {
	// Validate the request
	if err := validateUserUpdateBranchGo(req); err != nil {
		return nil, err
	}

	newOID, err := git.NewObjectIDFromHex(string(req.Newrev))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not parse newrev: %v", err)
	}

	oldOID, err := git.NewObjectIDFromHex(string(req.Oldrev))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not parse oldrev: %v", err)
	}

	referenceName := git.NewReferenceNameFromBranchName(string(req.BranchName))

	quarantineDir, _, err := s.quarantinedRepo(ctx, req.GetRepository(), featureflag.Quarantine)
	if err != nil {
		return nil, err
	}

	if err := s.updateReferenceWithHooks(ctx, req.GetRepository(), req.User, quarantineDir, referenceName, newOID, oldOID); err != nil {
		var preReceiveError updateref.PreReceiveError
		if errors.As(err, &preReceiveError) {
			return &gitalypb.UserUpdateBranchResponse{
				PreReceiveError: preReceiveError.Message,
			}, nil
		}

		// An oddball response for compatibility with the old
		// Ruby code. The "Could not update..."  message is
		// exactly like the default updateRefError, except we
		// say "branch-name", not
		// "refs/heads/branch-name". See the
		// "Gitlab::Git::CommitError" case in the Ruby code.
		return nil, status.Errorf(codes.FailedPrecondition, "Could not update %s. Please refresh and try again.", req.BranchName)
	}

	return &gitalypb.UserUpdateBranchResponse{}, nil
}

func (s *Server) UserDeleteBranch(ctx context.Context, req *gitalypb.UserDeleteBranchRequest) (*gitalypb.UserDeleteBranchResponse, error) {
	// That we do the branch name & user check here first only in
	// UserDelete but not UserCreate is "intentional", i.e. it's
	// always been that way.
	if len(req.BranchName) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Bad Request (empty branch name)")
	}

	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Bad Request (empty user)")
	}

	referenceName := git.NewReferenceNameFromBranchName(string(req.BranchName))

	referenceValue, err := s.localrepo(req.GetRepository()).ResolveRevision(ctx, referenceName.Revision())
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "branch not found: %s", req.BranchName)
	}

	if err := s.updateReferenceWithHooks(ctx, req.Repository, req.User, nil, referenceName, git.ZeroOID, referenceValue); err != nil {
		var preReceiveError updateref.PreReceiveError
		if errors.As(err, &preReceiveError) {
			return &gitalypb.UserDeleteBranchResponse{
				PreReceiveError: preReceiveError.Message,
			}, nil
		}

		var updateRefError updateref.Error
		if errors.As(err, &updateRefError) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return nil, err
	}

	return &gitalypb.UserDeleteBranchResponse{}, nil
}
