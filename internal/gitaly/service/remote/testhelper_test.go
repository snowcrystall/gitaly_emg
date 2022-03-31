package remote

import (
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/rubyserver"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testserver"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	defer testhelper.MustHaveNoChildProcess()

	cleanup := testhelper.Configure()
	defer cleanup()

	return m.Run()
}

func TestWithRubySidecar(t *testing.T) {
	t.Parallel()
	cfg := testcfg.Build(t)

	rubySrv := rubyserver.New(cfg)
	require.NoError(t, rubySrv.Start())
	t.Cleanup(rubySrv.Stop)

	fs := []func(t *testing.T, cfg config.Cfg, rubySrv *rubyserver.Server){
		testSuccessfulAddRemote,
		testAddRemoteTransactional,
	}

	for _, f := range fs {
		t.Run(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), func(t *testing.T) {
			f(t, cfg, rubySrv)
		})
	}
}

func setupRemoteServiceWithRuby(t *testing.T, cfg config.Cfg, rubySrv *rubyserver.Server, opts ...testserver.GitalyServerOpt) (config.Cfg, *gitalypb.Repository, string, gitalypb.RemoteServiceClient) {
	t.Helper()

	repo, repoPath := gittest.CloneRepo(t, cfg, cfg.Storages[0])

	addr := testserver.RunGitalyServer(t, cfg, rubySrv, func(srv *grpc.Server, deps *service.Dependencies) {
		gitalypb.RegisterRemoteServiceServer(srv, NewServer(
			deps.GetCfg(),
			deps.GetRubyServer(),
			deps.GetLocator(),
			deps.GetGitCmdFactory(),
			deps.GetCatfileCache(),
			deps.GetTxManager(),
		))
	}, opts...)
	cfg.SocketPath = addr

	client, conn := newRemoteClient(t, addr)
	t.Cleanup(func() { conn.Close() })

	return cfg, repo, repoPath, client
}

func setupRemoteService(t *testing.T, opts ...testserver.GitalyServerOpt) (config.Cfg, *gitalypb.Repository, string, gitalypb.RemoteServiceClient) {
	t.Helper()

	cfg := testcfg.Build(t)
	return setupRemoteServiceWithRuby(t, cfg, nil, opts...)
}

func newRemoteClient(t *testing.T, serverSocketPath string) (gitalypb.RemoteServiceClient, *grpc.ClientConn) {
	t.Helper()

	connOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(serverSocketPath, connOpts...)
	if err != nil {
		t.Fatal(err)
	}

	return gitalypb.NewRemoteServiceClient(conn), conn
}
