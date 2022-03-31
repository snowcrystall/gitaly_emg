package ssh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/command"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
)

type cloneCommand struct {
	command      *exec.Cmd
	repository   *gitalypb.Repository
	server       string
	featureFlags []string
	gitConfig    string
	gitProtocol  string
	cfg          config.Cfg
}

func runTestWithAndWithoutConfigOptions(t *testing.T, tf func(t *testing.T, opts ...testcfg.Option), opts ...testcfg.Option) {
	t.Run("no config options", func(t *testing.T) { tf(t) })

	if len(opts) > 0 {
		t.Run("with config options", func(t *testing.T) {
			tf(t, opts...)
		})
	}
}

func (cmd cloneCommand) execute(t *testing.T) error {
	req := &gitalypb.SSHUploadPackRequest{
		Repository:  cmd.repository,
		GitProtocol: cmd.gitProtocol,
	}
	if cmd.gitConfig != "" {
		req.GitConfigOptions = strings.Split(cmd.gitConfig, " ")
	}
	payload, err := protojson.Marshal(req)
	require.NoError(t, err)

	var flagPairs []string
	for _, flag := range cmd.featureFlags {
		flagPairs = append(flagPairs, fmt.Sprintf("%s:true", flag))
	}

	cmd.command.Env = []string{
		fmt.Sprintf("GITALY_ADDRESS=%s", cmd.server),
		fmt.Sprintf("GITALY_PAYLOAD=%s", payload),
		fmt.Sprintf("GITALY_FEATUREFLAGS=%s", strings.Join(flagPairs, ",")),
		fmt.Sprintf("PATH=.:%s", os.Getenv("PATH")),
		fmt.Sprintf(`GIT_SSH_COMMAND=%s upload-pack`, filepath.Join(cmd.cfg.BinDir, "gitaly-ssh")),
	}

	out, err := cmd.command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %q", err, out)
	}
	if !cmd.command.ProcessState.Success() {
		return fmt.Errorf("Failed to run `git clone`: %q", out)
	}

	return nil
}

func (cmd cloneCommand) test(t *testing.T, cfg config.Cfg, repoPath string, localRepoPath string) (string, string, string, string) {
	t.Helper()

	defer func() { require.NoError(t, os.RemoveAll(localRepoPath)) }()

	err := cmd.execute(t)
	require.NoError(t, err)

	remoteHead := text.ChompBytes(gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", "master"))
	localHead := text.ChompBytes(gittest.Exec(t, cfg, "-C", localRepoPath, "rev-parse", "master"))

	remoteTags := text.ChompBytes(gittest.Exec(t, cfg, "-C", repoPath, "tag"))
	localTags := text.ChompBytes(gittest.Exec(t, cfg, "-C", localRepoPath, "tag"))

	return localHead, remoteHead, localTags, remoteTags
}

func TestFailedUploadPackRequestDueToTimeout(t *testing.T) {
	runTestWithAndWithoutConfigOptions(t, testFailedUploadPackRequestDueToTimeout, testcfg.WithPackObjectsCacheEnabled())
}

func testFailedUploadPackRequestDueToTimeout(t *testing.T, opts ...testcfg.Option) {
	cfg, repo, _ := testcfg.BuildWithRepo(t, opts...)

	serverSocketPath := runSSHServerWithOptions(t, cfg, []ServerOpt{WithUploadPackRequestTimeout(10 * time.Microsecond)})

	client, conn := newSSHClient(t, serverSocketPath)
	defer conn.Close()

	ctx, cancel := testhelper.Context()
	defer cancel()

	stream, err := client.SSHUploadPack(ctx)
	require.NoError(t, err)

	// The first request is not limited by timeout, but also not under attacker control
	require.NoError(t, stream.Send(&gitalypb.SSHUploadPackRequest{Repository: repo}))

	// Because the client says nothing, the server would block. Because of
	// the timeout, it won't block forever, and return with a non-zero exit
	// code instead.
	requireFailedSSHStream(t, func() (int32, error) {
		resp, err := stream.Recv()
		if err != nil {
			return 0, err
		}

		var code int32
		if status := resp.GetExitStatus(); status != nil {
			code = status.Value
		}

		return code, nil
	})
}

func requireFailedSSHStream(t *testing.T, recv func() (int32, error)) {
	done := make(chan struct{})
	var code int32
	var err error

	go func() {
		for err == nil {
			code, err = recv()
		}
		close(done)
	}()

	select {
	case <-done:
		require.Equal(t, io.EOF, err)
		require.NotEqual(t, 0, code, "exit status")
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for SSH stream")
	}
}

func TestFailedUploadPackRequestDueToValidationError(t *testing.T) {
	cfg := testcfg.Build(t)

	serverSocketPath := runSSHServer(t, cfg)

	client, conn := newSSHClient(t, serverSocketPath)
	defer conn.Close()

	tests := []struct {
		Desc string
		Req  *gitalypb.SSHUploadPackRequest
		Code codes.Code
	}{
		{
			Desc: "Repository.RelativePath is empty",
			Req:  &gitalypb.SSHUploadPackRequest{Repository: &gitalypb.Repository{StorageName: cfg.Storages[0].Name, RelativePath: ""}},
			Code: codes.InvalidArgument,
		},
		{
			Desc: "Repository is nil",
			Req:  &gitalypb.SSHUploadPackRequest{Repository: nil},
			Code: codes.InvalidArgument,
		},
		{
			Desc: "Data exists on first request",
			Req:  &gitalypb.SSHUploadPackRequest{Repository: &gitalypb.Repository{StorageName: cfg.Storages[0].Name, RelativePath: "path/to/repo"}, Stdin: []byte("Fail")},
			Code: codes.InvalidArgument,
		},
	}

	for _, test := range tests {
		t.Run(test.Desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()
			stream, err := client.SSHUploadPack(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if err = stream.Send(test.Req); err != nil {
				t.Fatal(err)
			}
			require.NoError(t, stream.CloseSend())

			err = testPostUploadPackFailedResponse(t, stream)
			testhelper.RequireGrpcError(t, err, test.Code)
		})
	}
}

func TestUploadPackCloneSuccess(t *testing.T) {
	runTestWithAndWithoutConfigOptions(t, testUploadPackCloneSuccess, testcfg.WithPackObjectsCacheEnabled())
}

func testUploadPackCloneSuccess(t *testing.T, opts ...testcfg.Option) {
	cfg, repo, repoPath := testcfg.BuildWithRepo(t, opts...)

	testhelper.BuildGitalyHooks(t, cfg)
	testhelper.BuildGitalySSH(t, cfg)

	negotiationMetrics := prometheus.NewCounterVec(prometheus.CounterOpts{}, []string{"feature"})

	serverSocketPath := runSSHServerWithOptions(t, cfg, []ServerOpt{WithPackfileNegotiationMetrics(negotiationMetrics)})

	localRepoPath := testhelper.TempDir(t)

	tests := []struct {
		cmd    *exec.Cmd
		desc   string
		deepen float64
	}{
		{
			cmd:    exec.Command(cfg.Git.BinPath, "clone", "git@localhost:test/test.git", localRepoPath),
			desc:   "full clone",
			deepen: 0,
		},
		{
			cmd:    exec.Command(cfg.Git.BinPath, "clone", "--depth", "1", "git@localhost:test/test.git", localRepoPath),
			desc:   "shallow clone",
			deepen: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			negotiationMetrics.Reset()

			cmd := cloneCommand{
				repository: repo,
				command:    tc.cmd,
				server:     serverSocketPath,
				cfg:        cfg,
			}
			lHead, rHead, _, _ := cmd.test(t, cfg, repoPath, localRepoPath)
			require.Equal(t, lHead, rHead, "local and remote head not equal")

			metric, err := negotiationMetrics.GetMetricWithLabelValues("deepen")
			require.NoError(t, err)
			require.Equal(t, tc.deepen, promtest.ToFloat64(metric))
		})
	}
}

func TestUploadPackWithPackObjectsHook(t *testing.T) {
	cfg, repo, _ := testcfg.BuildWithRepo(t, testcfg.WithPackObjectsCacheEnabled())

	filterDir := testhelper.TempDir(t)
	outputPath := filepath.Join(filterDir, "output")
	cfg.BinDir = filterDir

	testhelper.BuildGitalyHooks(t, cfg)
	testhelper.BuildGitalySSH(t, cfg)

	hookScript := fmt.Sprintf(
		`#!/bin/sh
echo 'I was invoked' >'%s'
shift
exec '%s' "$@"
`,
		outputPath, cfg.Git.BinPath)

	// We're using a custom pack-objetcs hook for git-upload-pack. In order
	// to assure that it's getting executed as expected, we're writing a
	// custom script which replaces the hook binary. It doesn't do anything
	// special, but writes an error message and errors out and should thus
	// cause the clone to fail with this error message.
	testhelper.WriteExecutable(t, filepath.Join(filterDir, "gitaly-hooks"), []byte(hookScript))

	serverSocketPath := runSSHServer(t, cfg)

	localRepoPath := testhelper.TempDir(t)

	err := cloneCommand{
		repository: repo,
		command:    exec.Command(cfg.Git.BinPath, "clone", "git@localhost:test/test.git", localRepoPath),
		server:     serverSocketPath,
		cfg:        cfg,
	}.execute(t)
	require.NoError(t, err)

	require.Equal(t, []byte("I was invoked\n"), testhelper.MustReadFile(t, outputPath))
}

func TestUploadPackWithoutSideband(t *testing.T) {
	runTestWithAndWithoutConfigOptions(t, testUploadPackWithoutSideband, testcfg.WithPackObjectsCacheEnabled())
}

func testUploadPackWithoutSideband(t *testing.T, opts ...testcfg.Option) {
	cfg, repo, _ := testcfg.BuildWithRepo(t, opts...)

	testhelper.BuildGitalySSH(t, cfg)
	testhelper.BuildGitalyHooks(t, cfg)

	serverSocketPath := runSSHServer(t, cfg)

	// While Git knows the side-band-64 capability, some other clients don't. There is no way
	// though to have Git not use that capability, so we're instead manually crafting a packfile
	// negotiation without that capability and send it along.
	negotiation := bytes.NewBuffer([]byte{})
	gittest.WritePktlineString(t, negotiation, "want 1e292f8fedd741b75372e19097c76d327140c312 multi_ack_detailed thin-pack include-tag ofs-delta agent=git/2.29.1")
	gittest.WritePktlineString(t, negotiation, "want 1e292f8fedd741b75372e19097c76d327140c312")
	gittest.WritePktlineFlush(t, negotiation)
	gittest.WritePktlineString(t, negotiation, "done")

	request := &gitalypb.SSHUploadPackRequest{
		Repository: repo,
	}
	payload, err := protojson.Marshal(request)
	require.NoError(t, err)

	// As we're not using the sideband, the remote process will write both to stdout and stderr.
	// Those simultaneous writes to both stdout and stderr created a race as we could've invoked
	// two concurrent `SendMsg`s on the gRPC stream. And given that `SendMsg` is not thread-safe
	// a deadlock would result.
	uploadPack := exec.Command(filepath.Join(cfg.BinDir, "gitaly-ssh"), "upload-pack", "dontcare", "dontcare")
	uploadPack.Env = []string{
		fmt.Sprintf("GITALY_ADDRESS=%s", serverSocketPath),
		fmt.Sprintf("GITALY_PAYLOAD=%s", payload),
		fmt.Sprintf("PATH=.:%s", os.Getenv("PATH")),
	}
	uploadPack.Stdin = negotiation

	out, err := uploadPack.CombinedOutput()
	require.NoError(t, err)
	require.True(t, uploadPack.ProcessState.Success())
	require.Contains(t, string(out), "refs/heads/master")
	require.Contains(t, string(out), "Counting objects")
	require.Contains(t, string(out), "PACK")
}

func TestUploadPackCloneWithPartialCloneFilter(t *testing.T) {
	runTestWithAndWithoutConfigOptions(t, testUploadPackCloneWithPartialCloneFilter, testcfg.WithPackObjectsCacheEnabled())
}

func testUploadPackCloneWithPartialCloneFilter(t *testing.T, opts ...testcfg.Option) {
	cfg, repo, _ := testcfg.BuildWithRepo(t, opts...)

	testhelper.BuildGitalySSH(t, cfg)
	testhelper.BuildGitalyHooks(t, cfg)

	serverSocketPath := runSSHServer(t, cfg)

	// Ruby file which is ~1kB in size and not present in HEAD
	blobLessThanLimit := "6ee41e85cc9bf33c10b690df09ca735b22f3790f"
	// Image which is ~100kB in size and not present in HEAD
	blobGreaterThanLimit := "18079e308ff9b3a5e304941020747e5c39b46c88"

	tests := []struct {
		desc      string
		repoTest  func(t *testing.T, repoPath string)
		cloneArgs []string
	}{
		{
			desc: "full_clone",
			repoTest: func(t *testing.T, repoPath string) {
				gittest.GitObjectMustExist(t, cfg.Git.BinPath, repoPath, blobGreaterThanLimit)
			},
			cloneArgs: []string{"clone", "git@localhost:test/test.git"},
		},
		{
			desc: "partial_clone",
			repoTest: func(t *testing.T, repoPath string) {
				gittest.GitObjectMustNotExist(t, cfg.Git.BinPath, repoPath, blobGreaterThanLimit)
			},
			cloneArgs: []string{"clone", "--filter=blob:limit=2048", "git@localhost:test/test.git"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			// Run the clone with filtering enabled in both runs. The only
			// difference is that in the first run, we have the
			// UploadPackFilter flag disabled.
			localPath := testhelper.TempDir(t)

			cmd := cloneCommand{
				repository: repo,
				command:    exec.Command(cfg.Git.BinPath, append(tc.cloneArgs, localPath)...),
				server:     serverSocketPath,
				cfg:        cfg,
			}
			err := cmd.execute(t)
			defer func() { require.NoError(t, os.RemoveAll(localPath)) }()
			require.NoError(t, err, "clone failed")

			gittest.GitObjectMustExist(t, cfg.Git.BinPath, localPath, blobLessThanLimit)
			tc.repoTest(t, localPath)
		})
	}
}

func TestUploadPackCloneSuccessWithGitProtocol(t *testing.T) {
	runTestWithAndWithoutConfigOptions(t, testUploadPackCloneSuccessWithGitProtocol, testcfg.WithPackObjectsCacheEnabled())
}

func testUploadPackCloneSuccessWithGitProtocol(t *testing.T, opts ...testcfg.Option) {
	cfg, repo, repoPath := testcfg.BuildWithRepo(t, opts...)

	testhelper.BuildGitalySSH(t, cfg)
	testhelper.BuildGitalyHooks(t, cfg)

	localRepoPath := testhelper.TempDir(t)

	tests := []struct {
		cmd  *exec.Cmd
		desc string
	}{
		{
			cmd:  exec.Command(cfg.Git.BinPath, "clone", "git@localhost:test/test.git", localRepoPath),
			desc: "full clone",
		},
		{
			cmd:  exec.Command(cfg.Git.BinPath, "clone", "--depth", "1", "git@localhost:test/test.git", localRepoPath),
			desc: "shallow clone",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			readProto, cfg := gittest.EnableGitProtocolV2Support(t, cfg)

			serverSocketPath := runSSHServer(t, cfg)

			cmd := cloneCommand{
				repository:  repo,
				command:     tc.cmd,
				server:      serverSocketPath,
				gitProtocol: git.ProtocolV2,
				cfg:         cfg,
			}

			lHead, rHead, _, _ := cmd.test(t, cfg, repoPath, localRepoPath)
			require.Equal(t, lHead, rHead, "local and remote head not equal")

			envData := readProto()
			require.Contains(t, envData, fmt.Sprintf("GIT_PROTOCOL=%s\n", git.ProtocolV2))
		})
	}
}

func TestUploadPackCloneHideTags(t *testing.T) {
	cfg, repo, repoPath := testcfg.BuildWithRepo(t)

	testhelper.BuildGitalySSH(t, cfg)
	testhelper.BuildGitalyHooks(t, cfg)

	serverSocketPath := runSSHServer(t, cfg)

	localRepoPath := testhelper.TempDir(t)

	cmd := exec.Command(cfg.Git.BinPath, "clone", "--mirror", "git@localhost:test/test.git", localRepoPath)
	cmd.Env = os.Environ()
	cmd.Env = append(command.GitEnv, cmd.Env...)
	cloneCmd := cloneCommand{
		repository: repo,
		command:    cmd,
		server:     serverSocketPath,
		gitConfig:  "transfer.hideRefs=refs/tags",
		cfg:        cfg,
	}
	_, _, lTags, rTags := cloneCmd.test(t, cfg, repoPath, localRepoPath)

	if lTags == rTags {
		t.Fatalf("local and remote tags are equal. clone failed: %q != %q", lTags, rTags)
	}
	if tag := "v1.0.0"; !strings.Contains(rTags, tag) {
		t.Fatalf("sanity check failed, tag %q not found in %q", tag, rTags)
	}
}

func TestUploadPackCloneFailure(t *testing.T) {
	cfg, repo, _ := testcfg.BuildWithRepo(t)

	serverSocketPath := runSSHServer(t, cfg)

	localRepoPath := testhelper.TempDir(t)

	cmd := cloneCommand{
		repository: &gitalypb.Repository{
			StorageName:  "foobar",
			RelativePath: repo.GetRelativePath(),
		},
		command: exec.Command(cfg.Git.BinPath, "clone", "git@localhost:test/test.git", localRepoPath),
		server:  serverSocketPath,
		cfg:     cfg,
	}
	err := cmd.execute(t)
	require.Error(t, err, "clone didn't fail")
}

func testPostUploadPackFailedResponse(t *testing.T, stream gitalypb.SSHService_SSHUploadPackClient) error {
	var err error
	var res *gitalypb.SSHUploadPackResponse

	for err == nil {
		res, err = stream.Recv()
		require.Nil(t, res.GetStdout())
	}

	return err
}
