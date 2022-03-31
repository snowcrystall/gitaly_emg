package praefect

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/backchannel"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	gconfig "gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/setup"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/v14/internal/listenmux"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/datastore"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/grpc-proxy/proxy"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/mock"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/nodes"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/nodes/tracker"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/protoregistry"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/service/transaction"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/transactions"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/promtest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testassert"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testserver"
	"gitlab.com/gitlab-org/gitaly/v14/internal/transaction/txinfo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/transaction/voting"
	"gitlab.com/gitlab-org/gitaly/v14/internal/version"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewBackchannelServerFactory(t *testing.T) {
	mgr := transactions.NewManager(config.Config{})

	logger := testhelper.DiscardTestEntry(t)
	registry := backchannel.NewRegistry()

	lm := listenmux.New(insecure.NewCredentials())
	lm.Register(backchannel.NewServerHandshaker(logger, registry, nil))

	server := grpc.NewServer(
		grpc.Creds(lm),
		grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
			id, err := backchannel.GetPeerID(stream.Context())
			if !assert.NoError(t, err) {
				return err
			}

			backchannelConn, err := registry.Backchannel(id)
			if !assert.NoError(t, err) {
				return err
			}

			resp, err := gitalypb.NewRefTransactionClient(backchannelConn).VoteTransaction(
				stream.Context(), &gitalypb.VoteTransactionRequest{
					ReferenceUpdatesHash: voting.VoteFromData([]byte{}).Bytes(),
				},
			)
			assert.Nil(t, resp)

			return err
		}),
	)
	defer server.Stop()

	ln, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	go server.Serve(ln)

	ctx, cancel := testhelper.Context()
	defer cancel()

	nodeSet, err := DialNodes(ctx, []*config.VirtualStorage{{
		Name:  "default",
		Nodes: []*config.Node{{Storage: "gitaly-1", Address: "tcp://" + ln.Addr().String()}},
	}}, nil, nil, backchannel.NewClientHandshaker(logger, NewBackchannelServerFactory(
		testhelper.DiscardTestEntry(t), transaction.NewServer(mgr),
	)))
	require.NoError(t, err)
	defer nodeSet.Close()

	clientConn := nodeSet["default"]["gitaly-1"].Connection
	testassert.GrpcEqualErr(t,
		status.Error(codes.NotFound, "transaction not found: 0"),
		clientConn.Invoke(ctx, "/Service/Method", &gitalypb.CreateBranchRequest{}, &gitalypb.CreateBranchResponse{}),
	)
}

func TestGitalyServerInfo(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	t.Run("gitaly responds with ok", func(t *testing.T) {
		firstCfg := testcfg.Build(t, testcfg.WithStorages("praefect-internal-1"))
		firstCfg.SocketPath = testserver.RunGitalyServer(t, firstCfg, nil, setup.RegisterAll, testserver.WithDisablePraefect())

		secondCfg := testcfg.Build(t, testcfg.WithStorages("praefect-internal-2"))
		secondCfg.SocketPath = testserver.RunGitalyServer(t, secondCfg, nil, setup.RegisterAll, testserver.WithDisablePraefect())

		conf := config.Config{
			VirtualStorages: []*config.VirtualStorage{
				{
					Name: "virtual-storage",
					Nodes: []*config.Node{
						{
							Storage: firstCfg.Storages[0].Name,
							Address: firstCfg.SocketPath,
						},
						{
							Storage: secondCfg.Storages[0].Name,
							Address: secondCfg.SocketPath,
						},
					},
				},
			},
		}

		nodeSet, err := DialNodes(ctx, conf.VirtualStorages, nil, nil, nil)
		require.NoError(t, err)
		t.Cleanup(nodeSet.Close)

		cc, _, cleanup := runPraefectServer(t, ctx, conf, buildOptions{
			withConnections: nodeSet.Connections(),
		})
		t.Cleanup(cleanup)

		gitVersion, err := git.CurrentVersion(ctx, git.NewExecCommandFactory(firstCfg))
		require.NoError(t, err)

		expected := &gitalypb.ServerInfoResponse{
			ServerVersion: version.GetVersion(),
			GitVersion:    gitVersion.String(),
			StorageStatuses: []*gitalypb.ServerInfoResponse_StorageStatus{
				{StorageName: conf.VirtualStorages[0].Name, Readable: true, Writeable: true, ReplicationFactor: 2},
			},
		}

		client := gitalypb.NewServerServiceClient(cc)
		actual, err := client.ServerInfo(ctx, &gitalypb.ServerInfoRequest{})
		require.NoError(t, err)
		for _, ss := range actual.StorageStatuses {
			ss.FsType = ""
			ss.FilesystemId = ""
		}
		require.True(t, proto.Equal(expected, actual), "expected: %v, got: %v", expected, actual)
	})

	t.Run("gitaly responds with error", func(t *testing.T) {
		cfg := testcfg.Build(t)

		conf := config.Config{
			VirtualStorages: []*config.VirtualStorage{
				{
					Name: "virtual-storage",
					Nodes: []*config.Node{
						{
							Storage: cfg.Storages[0].Name,
							Token:   cfg.Auth.Token,
							Address: "unix:///invalid.addr",
						},
					},
				},
			},
		}

		nodeSet, err := DialNodes(ctx, conf.VirtualStorages, nil, nil, nil)
		require.NoError(t, err)
		t.Cleanup(nodeSet.Close)

		cc, _, cleanup := runPraefectServer(t, ctx, conf, buildOptions{
			withConnections: nodeSet.Connections(),
		})
		t.Cleanup(cleanup)

		client := gitalypb.NewServerServiceClient(cc)
		actual, err := client.ServerInfo(ctx, &gitalypb.ServerInfoRequest{})
		require.NoError(t, err, "we expect praefect's server info to fail open even if the gitaly calls result in an error")
		require.Empty(t, actual.StorageStatuses, "got: %v", actual)
	})
}

func TestGitalyServerInfoBadNode(t *testing.T) {
	gitalySocket := testhelper.GetTemporaryGitalySocketFileName(t)
	healthSrv := testhelper.NewServerWithHealth(t, gitalySocket)
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_UNKNOWN)

	conf := config.Config{
		VirtualStorages: []*config.VirtualStorage{
			&config.VirtualStorage{
				Nodes: []*config.Node{
					&config.Node{
						Storage: "praefect-internal-1",
						Address: "unix://" + gitalySocket,
						Token:   "abc",
					},
				},
			},
		},
	}

	ctx, cancel := testhelper.Context()
	defer cancel()

	nodes, err := DialNodes(ctx, conf.VirtualStorages, nil, nil, nil)
	require.NoError(t, err)
	defer nodes.Close()

	cc, _, cleanup := runPraefectServer(t, ctx, conf, buildOptions{
		withConnections: nodes.Connections(),
	})
	defer cleanup()

	client := gitalypb.NewServerServiceClient(cc)

	metadata, err := client.ServerInfo(ctx, &gitalypb.ServerInfoRequest{})
	require.NoError(t, err)
	require.Len(t, metadata.GetStorageStatuses(), 0)
}

func TestDiskStatistics(t *testing.T) {
	praefectCfg := config.Config{VirtualStorages: []*config.VirtualStorage{{Name: "praefect"}}}
	for _, name := range []string{"gitaly-1", "gitaly-2"} {
		gitalyCfg := testcfg.Build(t)

		gitalyAddr := testserver.RunGitalyServer(t, gitalyCfg, nil, setup.RegisterAll, testserver.WithDisablePraefect())

		praefectCfg.VirtualStorages[0].Nodes = append(praefectCfg.VirtualStorages[0].Nodes, &config.Node{
			Storage: name,
			Address: gitalyAddr,
			Token:   gitalyCfg.Auth.Token,
		})
	}

	ctx, cancel := testhelper.Context()
	defer cancel()

	nodes, err := DialNodes(ctx, praefectCfg.VirtualStorages, nil, nil, nil)
	require.NoError(t, err)
	defer nodes.Close()

	cc, _, cleanup := runPraefectServer(t, ctx, praefectCfg, buildOptions{
		withConnections: nodes.Connections(),
	})
	defer cleanup()

	client := gitalypb.NewServerServiceClient(cc)

	diskStat, err := client.DiskStatistics(ctx, &gitalypb.DiskStatisticsRequest{})
	require.NoError(t, err)
	require.Len(t, diskStat.GetStorageStatuses(), 2)

	for _, storageStatus := range diskStat.GetStorageStatuses() {
		require.NotNil(t, storageStatus, "none of the storage statuses should be nil")
	}
}

func TestHealthCheck(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	cc, _, cleanup := runPraefectServer(t, ctx, config.Config{VirtualStorages: []*config.VirtualStorage{
		{
			Name:  "praefect",
			Nodes: []*config.Node{{Storage: "stub", Address: "unix:///stub-address", Token: ""}},
		},
	}}, buildOptions{})
	defer cleanup()

	// setup timeout only after praefect setup as db migration may require some time
	ctx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()

	client := grpc_health_v1.NewHealthClient(cc)
	_, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	require.NoError(t, err)
}

func TestRejectBadStorage(t *testing.T) {
	conf := config.Config{
		VirtualStorages: []*config.VirtualStorage{
			{
				Name: "praefect",
				Nodes: []*config.Node{
					{
						Storage: "praefect-internal-0",
						Address: "tcp://this-doesnt-matter",
					},
				},
			},
		},
	}

	ctx, cancel := testhelper.Context()
	defer cancel()

	cc, _, cleanup := runPraefectServer(t, ctx, conf, buildOptions{})
	defer cleanup()

	req := &gitalypb.GarbageCollectRequest{
		Repository: &gitalypb.Repository{
			StorageName:  "bad-name",
			RelativePath: "/path/to/hashed/storage",
		},
	}

	_, err := gitalypb.NewRepositoryServiceClient(cc).GarbageCollect(ctx, req)
	require.Error(t, err, status.New(codes.InvalidArgument, "repo scoped: invalid Repository"))
}

func TestWarnDuplicateAddrs(t *testing.T) {
	conf := config.Config{
		VirtualStorages: []*config.VirtualStorage{
			&config.VirtualStorage{
				Name: "default",
				Nodes: []*config.Node{
					&config.Node{
						Storage: "praefect-internal-0",
						Address: "tcp://abc",
					},
					&config.Node{
						Storage: "praefect-internal-1",
						Address: "tcp://xyz",
					},
				},
			},
			&config.VirtualStorage{
				Name: "praefect",
				Nodes: []*config.Node{
					&config.Node{
						Storage: "praefect-internal-0",
						Address: "tcp://abc",
					},
					&config.Node{
						Storage: "praefect-internal-1",
						Address: "tcp://xyz",
					},
				},
			},
		},
	}

	ctx, cancel := testhelper.Context()
	defer cancel()

	tLogger, hook := test.NewNullLogger()

	// instantiate a praefect server and trigger warning
	_, _, cleanup := runPraefectServer(t, ctx, conf, buildOptions{
		withLogger:  logrus.NewEntry(tLogger),
		withNodeMgr: nullNodeMgr{}, // to suppress node address issues
	})
	defer cleanup()

	for _, entry := range hook.AllEntries() {
		require.NotContains(t, entry.Message, "more than one backend node")
	}

	conf = config.Config{
		VirtualStorages: []*config.VirtualStorage{
			&config.VirtualStorage{
				Name: "praefect",
				Nodes: []*config.Node{
					&config.Node{
						Storage: "praefect-internal-0",
						Address: "tcp::/samesies",
					},
					&config.Node{
						Storage: "praefect-internal-1",
						Address: "tcp::/samesies",
					},
				},
			},
		},
	}

	tLogger, hook = test.NewNullLogger()

	// instantiate a praefect server and trigger warning
	_, _, cleanup = runPraefectServer(t, ctx, conf, buildOptions{
		withLogger:  logrus.NewEntry(tLogger),
		withNodeMgr: nullNodeMgr{}, // to suppress node address issues
	})
	defer cleanup()

	var found bool
	for _, entry := range hook.AllEntries() {
		if strings.Contains(entry.Message, "more than one backend node") {
			found = true
			break
		}
	}
	require.True(t, found, "expected to find error log")

	conf = config.Config{
		VirtualStorages: []*config.VirtualStorage{
			&config.VirtualStorage{
				Name: "default",
				Nodes: []*config.Node{
					&config.Node{
						Storage: "praefect-internal-0",
						Address: "tcp://abc",
					},
					&config.Node{
						Storage: "praefect-internal-1",
						Address: "tcp://xyz",
					},
				},
			},
			&config.VirtualStorage{
				Name: "praefect",
				Nodes: []*config.Node{
					&config.Node{
						Storage: "praefect-internal-0",
						Address: "tcp://abc",
					},
					&config.Node{
						Storage: "praefect-internal-2",
						Address: "tcp://xyz",
					},
				},
			},
		},
	}

	tLogger, hook = test.NewNullLogger()

	// instantiate a praefect server and trigger warning
	_, _, cleanup = runPraefectServer(t, ctx, conf, buildOptions{
		withLogger:  logrus.NewEntry(tLogger),
		withNodeMgr: nullNodeMgr{}, // to suppress node address issues
	})
	defer cleanup()

	for _, entry := range hook.AllEntries() {
		require.NotContains(t, entry.Message, "more than one backend node")
	}
}

func TestRemoveRepository(t *testing.T) {
	gitalyCfgs := make([]gconfig.Cfg, 3)
	repos := make([][]*gitalypb.Repository, 3)
	praefectCfg := config.Config{VirtualStorages: []*config.VirtualStorage{{Name: "praefect"}}}

	for i, name := range []string{"gitaly-1", "gitaly-2", "gitaly-3"} {
		cfgBuilder := testcfg.NewGitalyCfgBuilder(testcfg.WithStorages(name))
		gitalyCfgs[i], repos[i] = cfgBuilder.BuildWithRepoAt(t, "test-repository")

		gitalyAddr := testserver.RunGitalyServer(t, gitalyCfgs[i], nil, setup.RegisterAll, testserver.WithDisablePraefect())
		gitalyCfgs[i].SocketPath = gitalyAddr

		praefectCfg.VirtualStorages[0].Nodes = append(praefectCfg.VirtualStorages[0].Nodes, &config.Node{
			Storage: name,
			Address: gitalyAddr,
			Token:   gitalyCfgs[i].Auth.Token,
		})
	}

	verifyReposExistence := func(t *testing.T, code codes.Code) {
		for i, gitalyCfg := range gitalyCfgs {
			locator := gconfig.NewLocator(gitalyCfg)
			_, err := locator.GetRepoPath(repos[i][0])
			st, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, code, st.Code())
		}
	}

	verifyReposExistence(t, codes.OK)

	// TODO: once https://gitlab.com/gitlab-org/gitaly/-/issues/2703 is done and the replication manager supports
	// graceful shutdown, we can remove this code that waits for jobs to be complete
	queueInterceptor := datastore.NewReplicationEventQueueInterceptor(datastore.NewPostgresReplicationEventQueue(getDB(t)))
	jobsDoneCh := make(chan struct{}, 2)
	queueInterceptor.OnAcknowledge(func(ctx context.Context, state datastore.JobState, ids []uint64, queue datastore.ReplicationEventQueue) ([]uint64, error) {
		defer func() {
			if state == datastore.JobStateCompleted {
				jobsDoneCh <- struct{}{}
			}
		}()

		return queue.Acknowledge(ctx, state, ids)
	})

	repoStore := defaultRepoStore(praefectCfg)
	txMgr := defaultTxMgr(praefectCfg)
	nodeMgr, err := nodes.NewManager(testhelper.DiscardTestEntry(t), praefectCfg, nil,
		repoStore, promtest.NewMockHistogramVec(), protoregistry.GitalyProtoPreregistered,
		nil, backchannel.NewClientHandshaker(
			testhelper.DiscardTestEntry(t),
			NewBackchannelServerFactory(testhelper.DiscardTestEntry(t), transaction.NewServer(txMgr)),
		),
	)
	require.NoError(t, err)
	nodeMgr.Start(0, time.Hour)

	ctx, cancel := testhelper.Context()
	defer cancel()

	cc, _, cleanup := runPraefectServer(t, ctx, praefectCfg, buildOptions{
		withQueue:     queueInterceptor,
		withRepoStore: repoStore,
		withNodeMgr:   nodeMgr,
		withTxMgr:     txMgr,
	})
	defer cleanup()

	virtualRepo := proto.Clone(repos[0][0]).(*gitalypb.Repository)
	virtualRepo.StorageName = praefectCfg.VirtualStorages[0].Name

	_, err = gitalypb.NewRepositoryServiceClient(cc).RemoveRepository(ctx, &gitalypb.RemoveRepositoryRequest{
		Repository: virtualRepo,
	})
	require.NoError(t, err)

	for i := 0; i < cap(jobsDoneCh); i++ {
		<-jobsDoneCh
	}

	verifyReposExistence(t, codes.NotFound)
}

func pollUntilRemoved(t testing.TB, path string, deadline <-chan time.Time) {
	for {
		select {
		case <-deadline:
			require.Failf(t, "unable to detect path removal for %s", path)
		default:
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				return
			}
			require.NoError(t, err, "unexpected error while checking path %s", path)
		}
		time.Sleep(time.Millisecond)
	}
}

func TestRenameRepository(t *testing.T) {
	gitalyStorages := []string{"gitaly-1", "gitaly-2", "gitaly-3"}
	repoPaths := make([]string, len(gitalyStorages))
	praefectCfg := config.Config{
		VirtualStorages: []*config.VirtualStorage{{Name: "praefect"}},
		Failover:        config.Failover{Enabled: true},
	}

	var repo *gitalypb.Repository
	for i, storageName := range gitalyStorages {
		const relativePath = "test-repository"

		cfgBuilder := testcfg.NewGitalyCfgBuilder(testcfg.WithStorages(storageName))
		gitalyCfg, repos := cfgBuilder.BuildWithRepoAt(t, relativePath)
		if repo == nil {
			repo = repos[0]
		}

		gitalyAddr := testserver.RunGitalyServer(t, gitalyCfg, nil, setup.RegisterAll, testserver.WithDisablePraefect())

		praefectCfg.VirtualStorages[0].Nodes = append(praefectCfg.VirtualStorages[0].Nodes, &config.Node{
			Storage: storageName,
			Address: gitalyAddr,
			Token:   gitalyCfg.Auth.Token,
		})

		repoPaths[i] = filepath.Join(gitalyCfg.Storages[0].Path, relativePath)
	}

	var canCheckRepo sync.WaitGroup
	canCheckRepo.Add(2)

	evq := datastore.NewReplicationEventQueueInterceptor(datastore.NewPostgresReplicationEventQueue(getDB(t)))
	evq.OnAcknowledge(func(ctx context.Context, state datastore.JobState, ids []uint64, queue datastore.ReplicationEventQueue) ([]uint64, error) {
		defer canCheckRepo.Done()
		return queue.Acknowledge(ctx, state, ids)
	})

	ctx, cancel := testhelper.Context()
	defer cancel()

	cc, _, cleanup := runPraefectServer(t, ctx, praefectCfg, buildOptions{withQueue: evq})
	defer cleanup()

	// virtualRepo is a virtual repository all requests to it would be applied to the underline Gitaly nodes behind it
	virtualRepo := proto.Clone(repo).(*gitalypb.Repository)
	virtualRepo.StorageName = praefectCfg.VirtualStorages[0].Name

	repoServiceClient := gitalypb.NewRepositoryServiceClient(cc)

	newName, err := text.RandomHex(20)
	require.NoError(t, err)

	_, err = repoServiceClient.RenameRepository(ctx, &gitalypb.RenameRepositoryRequest{
		Repository:   virtualRepo,
		RelativePath: newName,
	})
	require.NoError(t, err)

	resp, err := repoServiceClient.RepositoryExists(ctx, &gitalypb.RepositoryExistsRequest{
		Repository: virtualRepo,
	})
	require.NoError(t, err)
	require.False(t, resp.GetExists(), "repo with old name must gone")

	// as we renamed the repo we need to update RelativePath before we could check if it exists
	renamedVirtualRepo := virtualRepo
	renamedVirtualRepo.RelativePath = newName

	// wait until replication jobs propagate changes to other storages
	// as we don't know which one will be used to check because of reads distribution
	canCheckRepo.Wait()

	for _, oldLocation := range repoPaths {
		pollUntilRemoved(t, oldLocation, time.After(10*time.Second))
		newLocation := filepath.Join(filepath.Dir(oldLocation), newName)
		require.DirExists(t, newLocation, "must be renamed on secondary from %q to %q", oldLocation, newLocation)
	}
}

type mockSmartHTTP struct {
	gitalypb.UnimplementedSmartHTTPServiceServer
	txMgr         *transactions.Manager
	m             sync.Mutex
	methodsCalled map[string]int
}

func (m *mockSmartHTTP) InfoRefsUploadPack(req *gitalypb.InfoRefsRequest, stream gitalypb.SmartHTTPService_InfoRefsUploadPackServer) error {
	m.m.Lock()
	defer m.m.Unlock()
	if m.methodsCalled == nil {
		m.methodsCalled = make(map[string]int)
	}

	m.methodsCalled["InfoRefsUploadPack"]++

	return stream.Send(&gitalypb.InfoRefsResponse{})
}

func (m *mockSmartHTTP) InfoRefsReceivePack(req *gitalypb.InfoRefsRequest, stream gitalypb.SmartHTTPService_InfoRefsReceivePackServer) error {
	m.m.Lock()
	defer m.m.Unlock()
	if m.methodsCalled == nil {
		m.methodsCalled = make(map[string]int)
	}

	m.methodsCalled["InfoRefsReceivePack"]++

	return stream.Send(&gitalypb.InfoRefsResponse{})
}

func (m *mockSmartHTTP) PostUploadPack(stream gitalypb.SmartHTTPService_PostUploadPackServer) error {
	m.m.Lock()
	defer m.m.Unlock()
	if m.methodsCalled == nil {
		m.methodsCalled = make(map[string]int)
	}

	m.methodsCalled["PostUploadPack"]++

	return stream.Send(&gitalypb.PostUploadPackResponse{})
}

func (m *mockSmartHTTP) PostReceivePack(stream gitalypb.SmartHTTPService_PostReceivePackServer) error {
	m.m.Lock()
	defer m.m.Unlock()
	if m.methodsCalled == nil {
		m.methodsCalled = make(map[string]int)
	}

	m.methodsCalled["PostReceivePack"]++

	var err error
	var req *gitalypb.PostReceivePackRequest
	for {
		req, err = stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return helper.ErrInternal(err)
		}

		if err := stream.Send(&gitalypb.PostReceivePackResponse{Data: req.GetData()}); err != nil {
			return helper.ErrInternal(err)
		}
	}

	ctx := stream.Context()

	tx, err := txinfo.TransactionFromContext(ctx)
	if err != nil {
		return helper.ErrInternal(err)
	}

	vote := voting.VoteFromData([]byte{})
	if err := m.txMgr.VoteTransaction(ctx, tx.ID, tx.Node, vote); err != nil {
		return helper.ErrInternal(err)
	}

	return nil
}

func (m *mockSmartHTTP) Called(method string) int {
	m.m.Lock()
	defer m.m.Unlock()

	return m.methodsCalled[method]
}

func newSmartHTTPGrpcServer(t *testing.T, cfg gconfig.Cfg, smartHTTPService gitalypb.SmartHTTPServiceServer) string {
	return testserver.RunGitalyServer(t, cfg, nil, func(srv *grpc.Server, deps *service.Dependencies) {
		gitalypb.RegisterSmartHTTPServiceServer(srv, smartHTTPService)
	}, testserver.WithDisablePraefect())
}

func TestProxyWrites(t *testing.T) {
	txMgr := transactions.NewManager(config.Config{})

	smartHTTP0, smartHTTP1, smartHTTP2 := &mockSmartHTTP{txMgr: txMgr}, &mockSmartHTTP{txMgr: txMgr}, &mockSmartHTTP{txMgr: txMgr}

	cfg0 := testcfg.Build(t, testcfg.WithStorages("praefect-internal-0"))
	addr0 := newSmartHTTPGrpcServer(t, cfg0, smartHTTP0)

	cfg1 := testcfg.Build(t, testcfg.WithStorages("praefect-internal-1"))
	addr1 := newSmartHTTPGrpcServer(t, cfg1, smartHTTP1)

	cfg2 := testcfg.Build(t, testcfg.WithStorages("praefect-internal-2"))
	addr2 := newSmartHTTPGrpcServer(t, cfg2, smartHTTP2)

	conf := config.Config{
		VirtualStorages: []*config.VirtualStorage{
			{
				Name: "default",
				Nodes: []*config.Node{
					{
						Storage: cfg0.Storages[0].Name,
						Address: addr0,
					},
					{
						Storage: cfg1.Storages[0].Name,
						Address: addr1,
					},
					{
						Storage: cfg2.Storages[0].Name,
						Address: addr2,
					},
				},
			},
		},
	}

	queue := datastore.NewPostgresReplicationEventQueue(getDB(t))
	entry := testhelper.DiscardTestEntry(t)

	nodeMgr, err := nodes.NewManager(entry, conf, nil, nil, promtest.NewMockHistogramVec(), protoregistry.GitalyProtoPreregistered, nil, nil)
	require.NoError(t, err)
	nodeMgr.Start(0, time.Hour)

	ctx, cancel := testhelper.Context()
	defer cancel()

	_, repo, _ := testcfg.BuildWithRepo(t)

	rs := datastore.MockRepositoryStore{
		GetConsistentStoragesFunc: func(ctx context.Context, virtualStorage, relativePath string) (map[string]struct{}, error) {
			return map[string]struct{}{cfg0.Storages[0].Name: {}, cfg1.Storages[0].Name: {}, cfg2.Storages[0].Name: {}}, nil
		},
	}

	coordinator := NewCoordinator(
		queue,
		rs,
		NewNodeManagerRouter(nodeMgr, rs),
		txMgr,
		conf,
		protoregistry.GitalyProtoPreregistered,
	)

	server := grpc.NewServer(
		grpc.ForceServerCodec(proxy.NewCodec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(coordinator.StreamDirector)),
	)

	socket := testhelper.GetTemporaryGitalySocketFileName(t)
	listener, err := net.Listen("unix", socket)
	require.NoError(t, err)

	go server.Serve(listener)
	defer server.Stop()

	client, _ := newSmartHTTPClient(t, "unix://"+socket)

	shard, err := nodeMgr.GetShard(ctx, conf.VirtualStorages[0].Name)
	require.NoError(t, err)

	for _, storage := range conf.VirtualStorages[0].Nodes {
		node, err := shard.GetNode(storage.Storage)
		require.NoError(t, err)
		waitNodeToChangeHealthStatus(ctx, t, node, true)
	}

	stream, err := client.PostReceivePack(ctx)
	require.NoError(t, err)

	payload := "some pack data"
	for i := 0; i < 10; i++ {
		require.NoError(t, stream.Send(&gitalypb.PostReceivePackRequest{
			Repository: repo,
			Data:       []byte(payload),
		}))
	}

	require.NoError(t, stream.CloseSend())

	var receivedData bytes.Buffer
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			require.FailNowf(t, "unexpected non io.EOF error: %v", err.Error())
		}

		_, err = receivedData.Write(resp.GetData())
		require.NoError(t, err)
	}

	assert.Equal(t, 1, smartHTTP0.Called("PostReceivePack"))
	assert.Equal(t, 1, smartHTTP1.Called("PostReceivePack"))
	assert.Equal(t, 1, smartHTTP2.Called("PostReceivePack"))
	assert.Equal(t, bytes.Repeat([]byte(payload), 10), receivedData.Bytes())
}

func TestErrorThreshold(t *testing.T) {
	backendToken := ""
	backend, cleanup := newMockDownstream(t, backendToken, &mockSvc{
		repoMutatorUnary: func(ctx context.Context, req *mock.RepoRequest) (*emptypb.Empty, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return &emptypb.Empty{}, errors.New("couldn't read metadata")
			}

			if md.Get("bad-header")[0] == "true" {
				return &emptypb.Empty{}, helper.ErrInternalf("something went wrong")
			}

			return &emptypb.Empty{}, nil
		},
		repoAccessorUnary: func(ctx context.Context, req *mock.RepoRequest) (*emptypb.Empty, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return &emptypb.Empty{}, errors.New("couldn't read metadata")
			}

			if md.Get("bad-header")[0] == "true" {
				return &emptypb.Empty{}, helper.ErrInternalf("something went wrong")
			}

			return &emptypb.Empty{}, nil
		},
	})
	defer cleanup()

	conf := config.Config{
		VirtualStorages: []*config.VirtualStorage{
			{
				Name: "default",
				Nodes: []*config.Node{
					{
						Storage: "praefect-internal-0",
						Address: backend,
					},
				},
			},
		},
		Failover: config.Failover{
			Enabled:          true,
			ElectionStrategy: "local",
		},
	}

	ctx, cancel := testhelper.Context()
	defer cancel()

	queue := datastore.NewPostgresReplicationEventQueue(getDB(t))
	entry := testhelper.DiscardTestEntry(t)

	testCases := []struct {
		desc     string
		accessor bool
		error    error
	}{
		{
			desc:     "read threshold reached",
			accessor: true,
			error:    errors.New(`read error threshold reached for storage "praefect-internal-0"`),
		},
		{
			desc:  "write threshold reached",
			error: errors.New(`write error threshold reached for storage "praefect-internal-0"`),
		},
	}

	registry, err := protoregistry.NewFromPaths("praefect/mock/mock.proto")
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			readThreshold := uint32(5000)
			writeThreshold := uint32(5)
			if tc.accessor {
				readThreshold = 5
				writeThreshold = 5000
			}

			errorTracker, err := tracker.NewErrors(ctx, 10*time.Hour, readThreshold, writeThreshold)
			require.NoError(t, err)

			rs := datastore.MockRepositoryStore{}
			nodeMgr, err := nodes.NewManager(entry, conf, nil, rs, promtest.NewMockHistogramVec(), registry, errorTracker, nil)
			require.NoError(t, err)

			coordinator := NewCoordinator(
				queue,
				rs,
				NewNodeManagerRouter(nodeMgr, rs),
				nil,
				conf,
				registry,
			)

			server := grpc.NewServer(
				grpc.ForceServerCodec(proxy.NewCodec()),
				grpc.UnknownServiceHandler(proxy.TransparentHandler(coordinator.StreamDirector)),
			)

			socket := testhelper.GetTemporaryGitalySocketFileName(t)
			listener, err := net.Listen("unix", socket)
			require.NoError(t, err)

			go server.Serve(listener)
			defer server.Stop()

			conn, err := dial("unix://"+socket, []grpc.DialOption{grpc.WithInsecure()})
			require.NoError(t, err)
			cli := mock.NewSimpleServiceClient(conn)

			_, repo, _ := testcfg.BuildWithRepo(t)

			node := nodeMgr.Nodes()["default"][0]
			require.Equal(t, "praefect-internal-0", node.GetStorage())

			// to get the node in a healthy starting state, three consecutive successful checks
			// are needed
			for i := 0; i < 3; i++ {
				healthy, err := node.CheckHealth(ctx)
				require.NoError(t, err)
				require.True(t, healthy)
			}

			for i := 0; i < 5; i++ {
				ctx := metadata.AppendToOutgoingContext(ctx, "bad-header", "true")

				handler := cli.RepoMutatorUnary
				if tc.accessor {
					handler = cli.RepoAccessorUnary
				}

				healthy, err := node.CheckHealth(ctx)
				require.NoError(t, err)
				require.True(t, healthy)

				_, err = handler(ctx, &mock.RepoRequest{Repo: repo})
				testassert.GrpcEqualErr(t, status.Error(codes.Internal, "something went wrong"), err)
			}

			healthy, err := node.CheckHealth(ctx)
			require.Equal(t, tc.error, err)
			require.False(t, healthy)
		})
	}
}

func newSmartHTTPClient(t *testing.T, serverSocketPath string) (gitalypb.SmartHTTPServiceClient, *grpc.ClientConn) {
	t.Helper()

	conn, err := grpc.Dial(serverSocketPath, grpc.WithInsecure())
	require.NoError(t, err)

	return gitalypb.NewSmartHTTPServiceClient(conn), conn
}
