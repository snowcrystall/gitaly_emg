package repository

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
)

func TestSuccessfulCreateRepositoryFromURLRequest(t *testing.T) {
	t.Parallel()
	cfg, _, repoPath, client := setupRepositoryService(t)

	ctx, cancel := testhelper.Context()
	defer cancel()

	importedRepo := &gitalypb.Repository{
		RelativePath: "imports/test-repo-imported.git",
		StorageName:  cfg.Storages[0].Name,
	}
	importedRepoPath := filepath.Join(cfg.Storages[0].Path, importedRepo.GetRelativePath())

	user := "username123"
	password := "password321localhost"
	port, stopGitServer := gitServerWithBasicAuth(t, cfg, user, password, repoPath)
	defer func() {
		require.NoError(t, stopGitServer())
	}()

	url := fmt.Sprintf("http://%s:%s@localhost:%d/%s", user, password, port, filepath.Base(repoPath))

	req := &gitalypb.CreateRepositoryFromURLRequest{
		Repository: importedRepo,
		Url:        url,
	}

	_, err := client.CreateRepositoryFromURL(ctx, req)
	require.NoError(t, err)

	gittest.Exec(t, cfg, "-C", importedRepoPath, "fsck")

	remotes := gittest.Exec(t, cfg, "-C", importedRepoPath, "remote")
	require.NotContains(t, string(remotes), "origin")

	info, err := os.Lstat(filepath.Join(importedRepoPath, "hooks"))
	require.NoError(t, err)
	require.NotEqual(t, 0, info.Mode()&os.ModeSymlink)
}

func TestCloneRepositoryFromUrlCommand(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	userInfo := "user:pass%21%3F%40"
	repositoryFullPath := "full/path/to/repository"
	url := fmt.Sprintf("https://%s@www.example.com/secretrepo.git", userInfo)

	cfg := testcfg.Build(t)
	s := server{cfg: cfg, gitCmdFactory: git.NewExecCommandFactory(cfg)}
	cmd, err := s.cloneFromURLCommand(ctx, &gitalypb.Repository{}, url, repositoryFullPath, nil)
	require.NoError(t, err)

	expectedScrubbedURL := "https://www.example.com/secretrepo.git"
	expectedBasicAuthHeader := fmt.Sprintf("Authorization: Basic %s", base64.StdEncoding.EncodeToString([]byte("user:pass!?@")))
	expectedHeader := fmt.Sprintf("http.extraHeader=%s", expectedBasicAuthHeader)

	var args = cmd.Args()
	require.Contains(t, args, expectedScrubbedURL)
	require.Contains(t, args, expectedHeader)
	require.NotContains(t, args, userInfo)
}

func TestFailedCreateRepositoryFromURLRequestDueToExistingTarget(t *testing.T) {
	t.Parallel()
	cfg, client := setupRepositoryServiceWithoutRepo(t)

	ctx, cancel := testhelper.Context()
	defer cancel()

	testCases := []struct {
		desc     string
		repoPath string
		isDir    bool
	}{
		{
			desc:     "target is a directory",
			repoPath: "imports/test-repo-import-dir.git",
			isDir:    true,
		},
		{
			desc:     "target is a file",
			repoPath: "imports/test-repo-import-file.git",
			isDir:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			importedRepo := &gitalypb.Repository{
				RelativePath: "imports/test-repo-imported.git",
				StorageName:  cfg.Storages[0].Name,
			}
			importedRepoPath := filepath.Join(cfg.Storages[0].Path, importedRepo.GetRelativePath())

			if testCase.isDir {
				require.NoError(t, os.MkdirAll(importedRepoPath, 0770))
			} else {
				require.NoError(t, ioutil.WriteFile(importedRepoPath, nil, 0644))
			}
			t.Cleanup(func() { require.NoError(t, os.RemoveAll(importedRepoPath)) })

			req := &gitalypb.CreateRepositoryFromURLRequest{
				Repository: importedRepo,
				Url:        "https://gitlab.com/gitlab-org/gitlab-test.git",
			}

			_, err := client.CreateRepositoryFromURL(ctx, req)
			testhelper.RequireGrpcError(t, err, codes.InvalidArgument)
		})
	}
}

func TestPreventingRedirect(t *testing.T) {
	t.Parallel()
	cfg, client := setupRepositoryServiceWithoutRepo(t)

	ctx, cancel := testhelper.Context()
	defer cancel()

	importedRepo := &gitalypb.Repository{
		RelativePath: "imports/test-repo-imported.git",
		StorageName:  cfg.Storages[0].Name,
	}

	httpServerState, redirectingServer := StartRedirectingTestServer()
	defer redirectingServer.Close()

	req := &gitalypb.CreateRepositoryFromURLRequest{
		Repository: importedRepo,
		Url:        redirectingServer.URL,
	}

	_, err := client.CreateRepositoryFromURL(ctx, req)

	require.True(t, httpServerState.serverVisited, "git command should make the initial HTTP request")
	require.False(t, httpServerState.serverVisitedAfterRedirect, "git command should not follow HTTP redirection")

	require.Error(t, err)
	require.Contains(t, err.Error(), "The requested URL returned error: 301")
}

func gitServerWithBasicAuth(t testing.TB, cfg config.Cfg, user, pass, repoPath string) (int, func() error) {
	return gittest.GitServer(t, cfg, repoPath, basicAuthMiddleware(t, user, pass))
}

func basicAuthMiddleware(t testing.TB, user, pass string) func(http.ResponseWriter, *http.Request, http.Handler) {
	return func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		authUser, authPass, ok := r.BasicAuth()
		require.True(t, ok, "should contain basic auth")
		require.Equal(t, user, authUser, "username should match")
		require.Equal(t, pass, authPass, "password should match")
		next.ServeHTTP(w, r)
	}
}
