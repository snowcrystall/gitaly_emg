package stats

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
)

func TestRepositoryProfile(t *testing.T) {
	cfg := testcfg.Build(t)

	testRepo, testRepoPath := gittest.InitRepo(t, cfg, cfg.Storages[0])

	ctx, cancel := testhelper.Context()
	defer cancel()

	hasBitmap, err := HasBitmap(testRepoPath)
	require.NoError(t, err)
	require.False(t, hasBitmap, "repository should not have a bitmap initially")
	unpackedObjects, err := UnpackedObjects(testRepoPath)
	require.NoError(t, err)
	require.Zero(t, unpackedObjects)
	packfiles, err := GetPackfiles(testRepoPath)
	require.NoError(t, err)
	require.Empty(t, packfiles)
	packfilesCount, err := PackfilesCount(testRepoPath)
	require.NoError(t, err)
	require.Zero(t, packfilesCount)

	blobs := 10
	blobIDs := gittest.WriteBlobs(t, cfg, testRepoPath, blobs)

	unpackedObjects, err = UnpackedObjects(testRepoPath)
	require.NoError(t, err)
	require.Equal(t, int64(blobs), unpackedObjects)

	gitCmdFactory := git.NewExecCommandFactory(cfg)
	looseObjects, err := LooseObjects(ctx, gitCmdFactory, testRepo)
	require.NoError(t, err)
	require.Equal(t, int64(blobs), looseObjects)

	for _, blobID := range blobIDs {
		commitID := gittest.WriteCommit(t, cfg, testRepoPath,
			gittest.WithTreeEntries(gittest.TreeEntry{
				Mode: "100644", Path: "blob", OID: git.ObjectID(blobID),
			}),
			gittest.WithParents(),
		)
		gittest.Exec(t, cfg, "-C", testRepoPath, "update-ref", "refs/heads/"+blobID, commitID.String())
	}

	// write a loose object
	gittest.WriteBlobs(t, cfg, testRepoPath, 1)

	gittest.Exec(t, cfg, "-C", testRepoPath, "repack", "-A", "-b", "-d")

	unpackedObjects, err = UnpackedObjects(testRepoPath)
	require.NoError(t, err)
	require.Zero(t, unpackedObjects)
	looseObjects, err = LooseObjects(ctx, gitCmdFactory, testRepo)
	require.NoError(t, err)
	require.Equal(t, int64(1), looseObjects)

	// let a ms elapse for the OS to recognize the blobs have been written after the packfile
	time.Sleep(1 * time.Millisecond)

	// write another loose object
	blobID := gittest.WriteBlobs(t, cfg, testRepoPath, 1)[0]

	// due to OS semantics, ensure that the blob has a timestamp that is after the packfile
	theFuture := time.Now().Add(10 * time.Minute)
	require.NoError(t, os.Chtimes(filepath.Join(testRepoPath, "objects", blobID[0:2], blobID[2:]), theFuture, theFuture))

	unpackedObjects, err = UnpackedObjects(testRepoPath)
	require.NoError(t, err)
	require.Equal(t, int64(1), unpackedObjects)

	looseObjects, err = LooseObjects(ctx, gitCmdFactory, testRepo)
	require.NoError(t, err)
	require.Equal(t, int64(2), looseObjects)
}
