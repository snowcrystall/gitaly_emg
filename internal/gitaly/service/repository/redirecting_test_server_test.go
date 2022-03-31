package repository

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/command"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
)

const redirectURL = "/redirect_url"

// RedirectingTestServerState holds information about whether the server was visited and redirect was happened
type RedirectingTestServerState struct {
	serverVisited              bool
	serverVisitedAfterRedirect bool
}

// StartRedirectingTestServer starts the test server with initial state
func StartRedirectingTestServer() (*RedirectingTestServerState, *httptest.Server) {
	state := &RedirectingTestServerState{}
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == redirectURL {
				state.serverVisitedAfterRedirect = true
			} else {
				state.serverVisited = true
				http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
			}
		}),
	)

	return state, server
}

func TestRedirectingServerRedirects(t *testing.T) {
	t.Parallel()
	cfg := testcfg.Build(t)
	dir := testhelper.TempDir(t)

	httpServerState, redirectingServer := StartRedirectingTestServer()

	// we only test for redirection, this command can fail after that
	cmd := exec.Command(cfg.Git.BinPath, "-c", "http.followRedirects=true", "clone", "--bare", redirectingServer.URL, dir)
	cmd.Env = append(command.GitEnv, cmd.Env...)
	cmd.Run()

	redirectingServer.Close()

	require.True(t, httpServerState.serverVisited, "git command should make the initial HTTP request")
	require.True(t, httpServerState.serverVisitedAfterRedirect, "git command should follow the HTTP redirect")
}
