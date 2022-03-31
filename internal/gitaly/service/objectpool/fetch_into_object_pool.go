package objectpool

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git/objectpool"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/stats"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func (s *server) FetchIntoObjectPool(ctx context.Context, req *gitalypb.FetchIntoObjectPoolRequest) (*gitalypb.FetchIntoObjectPoolResponse, error) {
	if err := validateFetchIntoObjectPoolRequest(req); err != nil {
		return nil, helper.ErrInvalidArgument(err)
	}

	objectPool, err := objectpool.FromProto(s.cfg, s.locator, s.gitCmdFactory, s.catfileCache, req.GetObjectPool())
	if err != nil {
		return nil, helper.ErrInvalidArgument(fmt.Errorf("object pool invalid: %v", err))
	}

	if err := objectPool.FetchFromOrigin(ctx, req.GetOrigin()); err != nil {
		return nil, helper.ErrInternal(err)
	}

	stats.LogObjectsInfo(ctx, s.gitCmdFactory, req.ObjectPool.Repository)

	return &gitalypb.FetchIntoObjectPoolResponse{}, nil
}

func validateFetchIntoObjectPoolRequest(req *gitalypb.FetchIntoObjectPoolRequest) error {
	if req.GetOrigin() == nil {
		return errors.New("origin is empty")
	}

	if req.GetObjectPool() == nil {
		return errors.New("object pool is empty")
	}

	originRepository, poolRepository := req.GetOrigin(), req.GetObjectPool().GetRepository()

	if originRepository.GetStorageName() != poolRepository.GetStorageName() {
		return errors.New("origin has different storage than object pool")
	}

	return nil
}
