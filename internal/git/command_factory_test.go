package git_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
)

func TestGitCommandProxy(t *testing.T) {
	cfg := testcfg.Build(t)

	requestReceived := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true
	}))
	defer ts.Close()

	oldHTTPProxy := os.Getenv("http_proxy")
	defer require.NoError(t, os.Setenv("http_proxy", oldHTTPProxy))
	require.NoError(t, os.Setenv("http_proxy", ts.URL))

	ctx, cancel := testhelper.Context()
	defer cancel()

	dir := testhelper.TempDir(t)

	gitCmdFactory := git.NewExecCommandFactory(cfg)
	cmd, err := gitCmdFactory.NewWithoutRepo(ctx, git.SubCmd{
		Name: "clone",
		Args: []string{"http://gitlab.com/bogus-repo", dir},
	}, git.WithDisabledHooks())
	require.NoError(t, err)

	err = cmd.Wait()
	require.NoError(t, err)
	require.True(t, requestReceived)
}

func TestExecCommandFactory_NewWithDir(t *testing.T) {
	cfg := testcfg.Build(t)

	gitCmdFactory := git.NewExecCommandFactory(cfg)

	t.Run("no dir specified", func(t *testing.T) {
		ctx, cancel := testhelper.Context()
		defer cancel()

		_, err := gitCmdFactory.NewWithDir(ctx, "", nil, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no 'dir' provided")
	})

	t.Run("runs in dir", func(t *testing.T) {
		repoPath := testhelper.TempDir(t)

		gittest.Exec(t, cfg, "init", repoPath)
		gittest.Exec(t, cfg, "-C", repoPath, "commit", "--allow-empty", "-m", "initial commit")

		ctx, cancel := testhelper.Context()
		defer cancel()

		var stderr bytes.Buffer
		cmd, err := gitCmdFactory.NewWithDir(ctx, repoPath, git.SubCmd{
			Name: "rev-parse",
			Args: []string{"master"},
		}, git.WithStderr(&stderr))
		require.NoError(t, err)

		revData, err := ioutil.ReadAll(cmd)
		require.NoError(t, err)

		require.NoError(t, cmd.Wait(), stderr.String())

		require.Equal(t, "99ed180822d96f70810847eba6d0d168c582258d", text.ChompBytes(revData))
	})

	t.Run("doesn't runs in non existing dir", func(t *testing.T) {
		ctx, cancel := testhelper.Context()
		defer cancel()

		var stderr bytes.Buffer
		_, err := gitCmdFactory.NewWithDir(ctx, "non-existing-dir", git.SubCmd{
			Name: "rev-parse",
			Args: []string{"master"},
		}, git.WithStderr(&stderr))
		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file or directory")
	})
}
