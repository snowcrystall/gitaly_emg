package ref

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func TestPackRefsSuccessfulRequest(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	cfg, repoProto, repoPath, client := setupRefService(t)

	packedRefs := linesInPackfile(t, repoPath)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	// creates some new heads
	newBranches := 10
	for i := 0; i < newBranches; i++ {
		require.NoError(t, repo.UpdateRef(ctx, git.ReferenceName(fmt.Sprintf("refs/heads/new-ref-%d", i)), "refs/heads/master", git.ZeroOID))
	}

	// pack all refs
	_, err := client.PackRefs(ctx, &gitalypb.PackRefsRequest{Repository: repoProto})
	require.NoError(t, err)

	files, err := ioutil.ReadDir(filepath.Join(repoPath, "refs/heads"))
	require.NoError(t, err)
	assert.Len(t, files, 0, "git pack-refs --all should have packed all refs in refs/heads")
	assert.Equal(t, packedRefs+newBranches, linesInPackfile(t, repoPath), fmt.Sprintf("should have added %d new lines to the packfile", newBranches))

	// ensure all refs are reachable
	for i := 0; i < newBranches; i++ {
		gittest.Exec(t, cfg, "-C", repoPath, "show-ref", fmt.Sprintf("refs/heads/new-ref-%d", i))
	}
}

func linesInPackfile(t *testing.T, repoPath string) int {
	packFile, err := os.Open(filepath.Join(repoPath, "packed-refs"))
	require.NoError(t, err)
	defer packFile.Close()
	scanner := bufio.NewScanner(packFile)
	var refs int
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		refs++
	}
	return refs
}
