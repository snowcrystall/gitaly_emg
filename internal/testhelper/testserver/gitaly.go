package testserver

import (
	"context"
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	gitalyauth "gitlab.com/gitlab-org/gitaly/v14/auth"
	"gitlab.com/gitlab-org/gitaly/v14/client"
	"gitlab.com/gitlab-org/gitaly/v14/internal/backchannel"
	"gitlab.com/gitlab-org/gitaly/v14/internal/cache"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/auth"
	gitalylog "gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config/log"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/hook"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/linguist"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/rubyserver"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/server"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/storage"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/transaction"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitlab"
	praefectconfig "gitlab.com/gitlab-org/gitaly/v14/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/datastore/glsql"
	"gitlab.com/gitlab-org/gitaly/v14/internal/streamcache"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// RunGitalyServer starts gitaly server based on the provided cfg and returns a connection address.
// It accepts addition Registrar to register all required service instead of
// calling service.RegisterAll explicitly because it creates a circular dependency
// when the function is used in on of internal/gitaly/service/... packages.
func RunGitalyServer(t testing.TB, cfg config.Cfg, rubyServer *rubyserver.Server, registrar func(srv *grpc.Server, deps *service.Dependencies), opts ...GitalyServerOpt) string {
	_, gitalyAddr, disablePraefect := runGitaly(t, cfg, rubyServer, registrar, opts...)

	praefectBinPath, ok := os.LookupEnv("GITALY_TEST_PRAEFECT_BIN")
	if !ok || disablePraefect {
		return gitalyAddr
	}

	praefectAddr, _ := runPraefectProxy(t, cfg, gitalyAddr, praefectBinPath)

	// In case we're running with a Praefect proxy, it will use Gitaly's health information to
	// inform routing decisions. The Gitaly node thus must be healthy.
	waitHealthy(t, cfg, gitalyAddr, 3, time.Second)

	return praefectAddr
}

// createDatabase create a new database with randomly generated name and returns it back to the caller.
func createDatabase(t testing.TB) string {
	db := glsql.GetDB(t)
	return db.Name
}

func runPraefectProxy(t testing.TB, cfg config.Cfg, gitalyAddr, praefectBinPath string) (string, func()) {
	tempDir := testhelper.TempDir(t)

	praefectServerSocketPath := "unix://" + testhelper.GetTemporaryGitalySocketFileName(t)

	pgport := os.Getenv("PGPORT")
	port, err := strconv.Atoi(pgport)
	require.NoError(t, err)

	dbName := createDatabase(t)

	conf := praefectconfig.Config{
		AllowLegacyElectors: true,
		SocketPath:          praefectServerSocketPath,
		Auth: auth.Config{
			Token: cfg.Auth.Token,
		},
		DB: praefectconfig.DB{
			Host:    os.Getenv("PGHOST"),
			Port:    port,
			User:    os.Getenv("PGUSER"),
			DBName:  dbName,
			SSLMode: "disable",
		},
		Failover: praefectconfig.Failover{
			Enabled:           true,
			ElectionStrategy:  praefectconfig.ElectionStrategyLocal,
			BootstrapInterval: config.Duration(time.Microsecond),
			MonitorInterval:   config.Duration(time.Second),
		},
		Replication: praefectconfig.DefaultReplicationConfig(),
		Logging: gitalylog.Config{
			Format: "json",
			Level:  "panic",
		},
	}

	// Only single storage will be served by the praefect instance.
	// We can't include more as it is invalid to use same address for
	// different storages.
	conf.VirtualStorages = []*praefectconfig.VirtualStorage{{
		Name: cfg.Storages[0].Name,
		Nodes: []*praefectconfig.Node{{
			Storage: cfg.Storages[0].Name,
			Address: gitalyAddr,
			Token:   cfg.Auth.Token,
		}},
	}}

	configFilePath := filepath.Join(tempDir, "config.toml")
	configFile, err := os.Create(configFilePath)
	require.NoError(t, err)
	defer testhelper.MustClose(t, configFile)

	require.NoError(t, toml.NewEncoder(configFile).Encode(&conf))
	require.NoError(t, configFile.Sync())

	cmd := exec.Command(praefectBinPath, "-config", configFilePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	require.NoError(t, cmd.Start())

	waitHealthy(t, cfg, praefectServerSocketPath, 3, time.Second)

	t.Cleanup(func() { _ = cmd.Wait() })
	shutdown := func() { _ = cmd.Process.Kill() }
	t.Cleanup(shutdown)
	return praefectServerSocketPath, shutdown
}

// GitalyServer is a helper that carries additional info and
// functionality about gitaly (+praefect) server.
type GitalyServer struct {
	shutdown func()
	address  string
}

// Shutdown terminates running gitaly (+praefect) server.
func (gs GitalyServer) Shutdown() {
	gs.shutdown()
}

// Address returns address of the running gitaly (or praefect) service.
func (gs GitalyServer) Address() string {
	return gs.address
}

// StartGitalyServer creates and runs gitaly (and praefect as a proxy) server.
func StartGitalyServer(t testing.TB, cfg config.Cfg, rubyServer *rubyserver.Server, registrar func(srv *grpc.Server, deps *service.Dependencies), opts ...GitalyServerOpt) GitalyServer {
	gitalySrv, gitalyAddr, disablePraefect := runGitaly(t, cfg, rubyServer, registrar, opts...)

	praefectBinPath, ok := os.LookupEnv("GITALY_TEST_PRAEFECT_BIN")
	if !ok || disablePraefect {
		return GitalyServer{
			shutdown: gitalySrv.Stop,
			address:  gitalyAddr,
		}
	}

	praefectAddr, shutdownPraefect := runPraefectProxy(t, cfg, gitalyAddr, praefectBinPath)
	return GitalyServer{
		shutdown: func() {
			shutdownPraefect()
			gitalySrv.Stop()
		},
		address: praefectAddr,
	}
}

// waitHealthy executes health check request `retries` times and awaits each `timeout` period to respond.
// After `retries` unsuccessful attempts it returns an error.
// Returns immediately without an error once get a successful health check response.
func waitHealthy(t testing.TB, cfg config.Cfg, addr string, retries int, timeout time.Duration) {
	grpcOpts := []grpc.DialOption{grpc.WithInsecure()}
	if cfg.Auth.Token != "" {
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(gitalyauth.RPCCredentialsV2(cfg.Auth.Token)))
	}

	conn, err := grpc.Dial(addr, grpcOpts...)
	require.NoError(t, err)
	defer conn.Close()

	for i := 0; i < retries; i++ {
		if IsHealthy(conn, timeout) {
			return
		}
	}

	require.FailNow(t, "server not yet ready to serve")
}

// IsHealthy creates a health client to passed in connection and send `Check` request.
// It waits for `timeout` duration to get response back.
// It returns `true` only if remote responds with `SERVING` status.
func IsHealthy(conn *grpc.ClientConn, timeout time.Duration) bool {
	healthClient := healthpb.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
	if err != nil {
		return false
	}

	if resp.Status != healthpb.HealthCheckResponse_SERVING {
		return false
	}

	return true
}

func runGitaly(t testing.TB, cfg config.Cfg, rubyServer *rubyserver.Server, registrar func(srv *grpc.Server, deps *service.Dependencies), opts ...GitalyServerOpt) (*grpc.Server, string, bool) {
	t.Helper()

	var gsd gitalyServerDeps
	for _, opt := range opts {
		gsd = opt(gsd)
	}

	deps := gsd.createDependencies(t, cfg, rubyServer)
	t.Cleanup(func() { gsd.conns.Close() })

	srv, err := server.NewGitalyServerFactory(
		cfg,
		gsd.logger.WithField("test", t.Name()),
		deps.GetBackchannelRegistry(),
		deps.GetDiskCache(),
	).CreateExternal(cfg.TLS.CertPath != "" && cfg.TLS.KeyPath != "")
	require.NoError(t, err)
	t.Cleanup(srv.Stop)

	registrar(srv, deps)
	if _, found := srv.GetServiceInfo()["grpc.health.v1.Health"]; !found {
		// we should register health service as it is used for the health checks
		// praefect service executes periodically (and on the bootstrap step)
		healthpb.RegisterHealthServer(srv, health.NewServer())
	}

	// listen on internal socket
	if cfg.InternalSocketDir != "" {
		internalSocketDir := filepath.Dir(cfg.GitalyInternalSocketPath())
		sds, err := os.Stat(internalSocketDir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				require.NoError(t, os.MkdirAll(internalSocketDir, 0700))
				t.Cleanup(func() { require.NoError(t, os.RemoveAll(internalSocketDir)) })
			} else {
				require.FailNow(t, err.Error())
			}
		} else {
			require.True(t, sds.IsDir())
		}

		internalListener, err := net.Listen("unix", cfg.GitalyInternalSocketPath())
		require.NoError(t, err)
		go srv.Serve(internalListener)
	}

	var listener net.Listener
	var addr string
	switch {
	case cfg.TLSListenAddr != "":
		listener, err = net.Listen("tcp", cfg.TLSListenAddr)
		require.NoError(t, err)
		_, port, err := net.SplitHostPort(listener.Addr().String())
		require.NoError(t, err)
		addr = "tls://localhost:" + port
	case cfg.ListenAddr != "":
		listener, err = net.Listen("tcp", cfg.ListenAddr)
		require.NoError(t, err)
		addr = "tcp://" + listener.Addr().String()
	default:
		serverSocketPath := testhelper.GetTemporaryGitalySocketFileName(t)
		listener, err = net.Listen("unix", serverSocketPath)
		require.NoError(t, err)
		addr = "unix://" + serverSocketPath
	}

	go srv.Serve(listener)

	return srv, addr, gsd.disablePraefect
}

type gitalyServerDeps struct {
	disablePraefect  bool
	logger           *logrus.Logger
	conns            *client.Pool
	locator          storage.Locator
	txMgr            transaction.Manager
	hookMgr          hook.Manager
	gitlabClient     gitlab.Client
	gitCmdFactory    git.CommandFactory
	linguist         *linguist.Instance
	backchannelReg   *backchannel.Registry
	catfileCache     catfile.Cache
	diskCache        cache.Cache
	packObjectsCache streamcache.Cache
}

func (gsd *gitalyServerDeps) createDependencies(t testing.TB, cfg config.Cfg, rubyServer *rubyserver.Server) *service.Dependencies {
	if gsd.logger == nil {
		gsd.logger = testhelper.DiscardTestLogger(t)
	}

	if gsd.conns == nil {
		gsd.conns = client.NewPool()
	}

	if gsd.locator == nil {
		gsd.locator = config.NewLocator(cfg)
	}

	if gsd.gitlabClient == nil {
		gsd.gitlabClient = gitlab.NewMockClient()
	}

	if gsd.backchannelReg == nil {
		gsd.backchannelReg = backchannel.NewRegistry()
	}

	if gsd.txMgr == nil {
		gsd.txMgr = transaction.NewManager(cfg, gsd.backchannelReg)
	}

	if gsd.hookMgr == nil {
		gsd.hookMgr = hook.NewManager(gsd.locator, gsd.txMgr, gsd.gitlabClient, cfg)
	}

	if gsd.gitCmdFactory == nil {
		gsd.gitCmdFactory = git.NewExecCommandFactory(cfg)
	}

	if gsd.linguist == nil {
		var err error
		gsd.linguist, err = linguist.New(cfg)
		require.NoError(t, err)
	}

	if gsd.catfileCache == nil {
		cache := catfile.NewCache(cfg)
		gsd.catfileCache = cache
		t.Cleanup(cache.Stop)
	}

	if gsd.diskCache == nil {
		gsd.diskCache = cache.New(cfg, gsd.locator)
	}

	if gsd.packObjectsCache == nil {
		gsd.packObjectsCache = streamcache.New(cfg.PackObjectsCache, gsd.logger)
	}

	return &service.Dependencies{
		Cfg:                 cfg,
		RubyServer:          rubyServer,
		ClientPool:          gsd.conns,
		StorageLocator:      gsd.locator,
		TransactionManager:  gsd.txMgr,
		GitalyHookManager:   gsd.hookMgr,
		GitCmdFactory:       gsd.gitCmdFactory,
		Linguist:            gsd.linguist,
		BackchannelRegistry: gsd.backchannelReg,
		GitlabClient:        gsd.gitlabClient,
		CatfileCache:        gsd.catfileCache,
		DiskCache:           gsd.diskCache,
		PackObjectsCache:    gsd.packObjectsCache,
	}
}

// GitalyServerOpt is a helper type to shorten declarations.
type GitalyServerOpt func(gitalyServerDeps) gitalyServerDeps

// WithLogger sets a logrus.Logger instance that will be used for gitaly services initialisation.
func WithLogger(logger *logrus.Logger) GitalyServerOpt {
	return func(deps gitalyServerDeps) gitalyServerDeps {
		deps.logger = logger
		return deps
	}
}

// WithLocator sets a storage.Locator instance that will be used for gitaly services initialisation.
func WithLocator(locator storage.Locator) GitalyServerOpt {
	return func(deps gitalyServerDeps) gitalyServerDeps {
		deps.locator = locator
		return deps
	}
}

// WithGitLabClient sets gitlab.Client instance that will be used for gitaly services initialisation.
func WithGitLabClient(gitlabClient gitlab.Client) GitalyServerOpt {
	return func(deps gitalyServerDeps) gitalyServerDeps {
		deps.gitlabClient = gitlabClient
		return deps
	}
}

// WithHookManager sets hook.Manager instance that will be used for gitaly services initialisation.
func WithHookManager(hookMgr hook.Manager) GitalyServerOpt {
	return func(deps gitalyServerDeps) gitalyServerDeps {
		deps.hookMgr = hookMgr
		return deps
	}
}

// WithTransactionManager sets transaction.Manager instance that will be used for gitaly services initialisation.
func WithTransactionManager(txMgr transaction.Manager) GitalyServerOpt {
	return func(deps gitalyServerDeps) gitalyServerDeps {
		deps.txMgr = txMgr
		return deps
	}
}

// WithDisablePraefect disables setup and usage of the praefect as a proxy before the gitaly service.
func WithDisablePraefect() GitalyServerOpt {
	return func(deps gitalyServerDeps) gitalyServerDeps {
		deps.disablePraefect = true
		return deps
	}
}

// WithBackchannelRegistry sets backchannel.Registry instance that will be used for gitaly services initialisation.
func WithBackchannelRegistry(backchannelReg *backchannel.Registry) GitalyServerOpt {
	return func(deps gitalyServerDeps) gitalyServerDeps {
		deps.backchannelReg = backchannelReg
		return deps
	}
}

// WithDiskCache sets the cache.Cache instance that will be used for gitaly services initialisation.
func WithDiskCache(diskCache cache.Cache) GitalyServerOpt {
	return func(deps gitalyServerDeps) gitalyServerDeps {
		deps.diskCache = diskCache
		return deps
	}
}
