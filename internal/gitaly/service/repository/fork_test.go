package repository

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	gitalyauth "gitlab.com/gitlab-org/gitaly/v14/auth"
	gclient "gitlab.com/gitlab-org/gitaly/v14/client"
	"gitlab.com/gitlab-org/gitaly/v14/internal/backchannel"
	"gitlab.com/gitlab-org/gitaly/v14/internal/cache"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/hook"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/rubyserver"
	gserver "gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/server"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/commit"
	hookservice "gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/hook"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/ref"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/remote"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/ssh"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/transaction"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitlab"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testserver"
	gitalyx509 "gitlab.com/gitlab-org/gitaly/v14/internal/x509"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func TestSuccessfulCreateForkRequest(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name          string
		secure        bool
		beforeRequest func(repoPath string)
	}{
		{
			name:   "secure",
			secure: true,
		},
		{
			name: "insecure",
		},
		{
			name: "existing empty directory target",
			beforeRequest: func(repoPath string) {
				require.NoError(t, os.MkdirAll(repoPath, 0755))
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg, repo, _ := testcfg.BuildWithRepo(t)

			testhelper.BuildGitalyHooks(t, cfg)
			testhelper.BuildGitalySSH(t, cfg)

			var (
				client gitalypb.RepositoryServiceClient
				conn   *grpc.ClientConn
			)

			if tt.secure {
				testPool := injectCustomCATestCerts(t, &cfg)
				cfg.TLSListenAddr = runSecureServer(t, cfg, nil)
				client, conn = newSecureRepoClient(t, cfg.TLSListenAddr, cfg.Auth.Token, testPool)
				defer conn.Close()
			} else {
				client, cfg.SocketPath = runRepositoryService(t, cfg, nil)
			}

			ctxOuter, cancel := testhelper.Context()
			defer cancel()

			md := testhelper.GitalyServersMetadataFromCfg(t, cfg)
			ctx := metadata.NewOutgoingContext(ctxOuter, md)

			forkedRepo := &gitalypb.Repository{
				RelativePath: "forks/test-repo-fork.git",
				StorageName:  repo.GetStorageName(),
			}

			forkedRepoPath := filepath.Join(cfg.Storages[0].Path, forkedRepo.GetRelativePath())
			require.NoError(t, os.RemoveAll(forkedRepoPath))

			if tt.beforeRequest != nil {
				tt.beforeRequest(forkedRepoPath)
			}

			req := &gitalypb.CreateForkRequest{
				Repository:       forkedRepo,
				SourceRepository: repo,
			}

			_, err := client.CreateFork(ctx, req)
			require.NoError(t, err)
			defer func() { require.NoError(t, os.RemoveAll(forkedRepoPath)) }()

			gittest.Exec(t, cfg, "-C", forkedRepoPath, "fsck")

			remotes := gittest.Exec(t, cfg, "-C", forkedRepoPath, "remote")
			require.NotContains(t, string(remotes), "origin")

			info, err := os.Lstat(filepath.Join(forkedRepoPath, "hooks"))
			require.NoError(t, err)
			require.NotEqual(t, 0, info.Mode()&os.ModeSymlink)
		})
	}
}

func newSecureRepoClient(t testing.TB, addr, token string, pool *x509.CertPool) (gitalypb.RepositoryServiceClient, *grpc.ClientConn) {
	t.Helper()

	connOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			RootCAs:    pool,
			MinVersion: tls.VersionTLS12,
		})),
		grpc.WithPerRPCCredentials(gitalyauth.RPCCredentialsV2(token)),
	}

	conn, err := gclient.Dial(addr, connOpts)
	require.NoError(t, err)

	return gitalypb.NewRepositoryServiceClient(conn), conn
}

func TestFailedCreateForkRequestDueToExistingTarget(t *testing.T) {
	t.Parallel()
	cfg, repo, _, client := setupRepositoryService(t)

	ctxOuter, cancel := testhelper.Context()
	defer cancel()

	md := testhelper.GitalyServersMetadataFromCfg(t, cfg)
	ctx := metadata.NewOutgoingContext(ctxOuter, md)

	testCases := []struct {
		desc     string
		repoPath string
		isDir    bool
	}{
		{
			desc:     "target is a non-empty directory",
			repoPath: "forks/test-repo-fork-dir.git",
			isDir:    true,
		},
		{
			desc:     "target is a file",
			repoPath: "forks/test-repo-fork-file.git",
			isDir:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			forkedRepo := &gitalypb.Repository{
				RelativePath: testCase.repoPath,
				StorageName:  repo.StorageName,
			}

			forkedRepoPath := filepath.Join(cfg.Storages[0].Path, forkedRepo.GetRelativePath())

			if testCase.isDir {
				require.NoError(t, os.MkdirAll(forkedRepoPath, 0770))
				require.NoError(t, ioutil.WriteFile(
					filepath.Join(forkedRepoPath, "config"),
					nil,
					0644,
				))
			} else {
				require.NoError(t, ioutil.WriteFile(forkedRepoPath, nil, 0644))
			}
			defer func() { require.NoError(t, os.RemoveAll(forkedRepoPath)) }()

			req := &gitalypb.CreateForkRequest{
				Repository:       forkedRepo,
				SourceRepository: repo,
			}

			_, err := client.CreateFork(ctx, req)
			testhelper.RequireGrpcError(t, err, codes.InvalidArgument)
		})
	}
}

func injectCustomCATestCerts(t *testing.T, cfg *config.Cfg) *x509.CertPool {
	certFile, keyFile := testhelper.GenerateCerts(t)

	cfg.TLS.CertPath = certFile
	cfg.TLS.KeyPath = keyFile

	revertEnv := testhelper.ModifyEnvironment(t, gitalyx509.SSLCertFile, certFile)
	t.Cleanup(revertEnv)

	caPEMBytes := testhelper.MustReadFile(t, certFile)
	pool := x509.NewCertPool()
	require.True(t, pool.AppendCertsFromPEM(caPEMBytes))

	return pool
}

func runSecureServer(t *testing.T, cfg config.Cfg, rubySrv *rubyserver.Server) string {
	t.Helper()

	registry := backchannel.NewRegistry()
	locator := config.NewLocator(cfg)
	cache := cache.New(cfg, locator)
	server, err := gserver.New(true, cfg, testhelper.DiscardTestEntry(t), registry, cache)
	require.NoError(t, err)
	listener, addr := testhelper.GetLocalhostListener(t)

	txManager := transaction.NewManager(cfg, registry)
	hookManager := hook.NewManager(locator, txManager, gitlab.NewMockClient(), cfg)
	gitCmdFactory := git.NewExecCommandFactory(cfg)
	catfileCache := catfile.NewCache(cfg)

	gitalypb.RegisterRepositoryServiceServer(server, NewServer(cfg, rubySrv, locator, txManager, gitCmdFactory, catfileCache))
	gitalypb.RegisterHookServiceServer(server, hookservice.NewServer(cfg, hookManager, gitCmdFactory, nil))
	gitalypb.RegisterRemoteServiceServer(server, remote.NewServer(cfg, rubySrv, locator, gitCmdFactory, catfileCache, txManager))
	gitalypb.RegisterSSHServiceServer(server, ssh.NewServer(cfg, locator, gitCmdFactory, txManager))
	gitalypb.RegisterRefServiceServer(server, ref.NewServer(cfg, locator, gitCmdFactory, txManager, catfileCache))
	gitalypb.RegisterCommitServiceServer(server, commit.NewServer(cfg, locator, gitCmdFactory, nil, catfileCache))
	errQ := make(chan error, 1)

	// This creates a secondary GRPC server which isn't "secure". Reusing
	// the one created above won't work as its internal socket would be
	// protected by the same TLS certificate.

	cfg.TLS.KeyPath = ""
	testserver.RunGitalyServer(t, cfg, nil, func(srv *grpc.Server, deps *service.Dependencies) {
		gitalypb.RegisterHookServiceServer(srv, hookservice.NewServer(deps.GetCfg(), deps.GetHookManager(), deps.GetGitCmdFactory(), deps.GetPackObjectsCache()))
	})

	t.Cleanup(func() { require.NoError(t, <-errQ) })

	t.Cleanup(server.Stop)
	go func() { errQ <- server.Serve(listener) }()

	return "tls://" + addr
}
