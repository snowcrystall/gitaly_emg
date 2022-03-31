package objectpool

import (
	"context"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git/objectpool"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errInvalidPoolDir = helper.ErrInvalidArgument(objectpool.ErrInvalidPoolDir)

	// errMissingOriginRepository is returned when the request is missing the
	// origin repository.
	errMissingOriginRepository = helper.ErrInvalidArgumentf("no origin repository")

	// errMissingPool is returned when the request is missing the object pool.
	errMissingPool = helper.ErrInvalidArgumentf("no object pool repository")
)

func (s *server) CreateObjectPool(ctx context.Context, in *gitalypb.CreateObjectPoolRequest) (*gitalypb.CreateObjectPoolResponse, error) {
	if in.GetOrigin() == nil {
		return nil, errMissingOriginRepository
	}

	pool, err := s.poolForRequest(in)
	if err != nil {
		return nil, err
	}

	if pool.Exists() {
		return nil, status.Errorf(codes.FailedPrecondition, "pool already exists at: %v", pool.GetRelativePath())
	}

	if err := pool.Create(ctx, in.GetOrigin()); err != nil {
		return nil, err
	}

	return &gitalypb.CreateObjectPoolResponse{}, nil
}

func (s *server) DeleteObjectPool(ctx context.Context, in *gitalypb.DeleteObjectPoolRequest) (*gitalypb.DeleteObjectPoolResponse, error) {
	pool, err := s.poolForRequest(in)
	if err != nil {
		return nil, err
	}

	if err := pool.Remove(ctx); err != nil {
		return nil, err
	}

	return &gitalypb.DeleteObjectPoolResponse{}, nil
}

type poolRequest interface {
	GetObjectPool() *gitalypb.ObjectPool
}

func (s *server) poolForRequest(req poolRequest) (*objectpool.ObjectPool, error) {
	reqPool := req.GetObjectPool()

	poolRepo := reqPool.GetRepository()
	if poolRepo == nil {
		return nil, errMissingPool
	}

	pool, err := objectpool.NewObjectPool(s.cfg, s.locator, s.gitCmdFactory, s.catfileCache, poolRepo.GetStorageName(), poolRepo.GetRelativePath())
	if err != nil {
		if err == objectpool.ErrInvalidPoolDir {
			return nil, errInvalidPoolDir
		}

		return nil, helper.ErrInternal(err)
	}

	return pool, nil
}
