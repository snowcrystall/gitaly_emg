package objectpool

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/objectpool"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func TestGetObjectPoolSuccess(t *testing.T) {
	cfg, repo, _, locator, client := setup(t)

	relativePoolPath := gittest.NewObjectPoolName(t)

	pool, err := objectpool.NewObjectPool(cfg, locator, git.NewExecCommandFactory(cfg), nil, repo.GetStorageName(), relativePoolPath)
	require.NoError(t, err)

	poolCtx, cancel := testhelper.Context()
	defer cancel()
	require.NoError(t, pool.Create(poolCtx, repo))
	require.NoError(t, pool.Link(poolCtx, repo))

	ctx, cancel := testhelper.Context()
	defer func() {
		require.NoError(t, pool.Remove(ctx))
	}()
	defer cancel()

	resp, err := client.GetObjectPool(ctx, &gitalypb.GetObjectPoolRequest{
		Repository: repo,
	})

	require.NoError(t, err)
	require.Equal(t, relativePoolPath, resp.GetObjectPool().GetRepository().GetRelativePath())
}

func TestGetObjectPoolNoFile(t *testing.T) {
	_, repoo, _, _, client := setup(t)

	ctx, cancel := testhelper.Context()
	defer cancel()

	resp, err := client.GetObjectPool(ctx, &gitalypb.GetObjectPoolRequest{
		Repository: repoo,
	})

	require.NoError(t, err)
	require.Nil(t, resp.GetObjectPool())
}

func TestGetObjectPoolBadFile(t *testing.T) {
	_, repo, repoPath, _, client := setup(t)

	alternatesFilePath := filepath.Join(repoPath, "objects", "info", "alternates")
	require.NoError(t, os.MkdirAll(filepath.Dir(alternatesFilePath), 0755))
	require.NoError(t, ioutil.WriteFile(alternatesFilePath, []byte("not-a-directory"), 0644))

	ctx, cancel := testhelper.Context()
	defer cancel()

	resp, err := client.GetObjectPool(ctx, &gitalypb.GetObjectPoolRequest{
		Repository: repo,
	})

	require.NoError(t, err)
	require.Nil(t, resp.GetObjectPool())
}
