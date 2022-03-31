package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestFailedUploadArchiveRequestDueToTimeout(t *testing.T) {
	cfg, repo, _ := testcfg.BuildWithRepo(t)

	serverSocketPath := runSSHServerWithOptions(t, cfg, []ServerOpt{WithArchiveRequestTimeout(100 * time.Microsecond)})

	client, conn := newSSHClient(t, serverSocketPath)
	defer conn.Close()

	ctx, cancel := testhelper.Context()
	defer cancel()

	stream, err := client.SSHUploadArchive(ctx)
	require.NoError(t, err)

	// The first request is not limited by timeout, but also not under attacker control
	require.NoError(t, stream.Send(&gitalypb.SSHUploadArchiveRequest{Repository: repo}))

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

func TestFailedUploadArchiveRequestDueToValidationError(t *testing.T) {
	cfg := testcfg.Build(t)

	serverSocketPath := runSSHServer(t, cfg)

	client, conn := newSSHClient(t, serverSocketPath)
	defer conn.Close()

	tests := []struct {
		Desc string
		Req  *gitalypb.SSHUploadArchiveRequest
		Code codes.Code
	}{
		{
			Desc: "Repository.RelativePath is empty",
			Req:  &gitalypb.SSHUploadArchiveRequest{Repository: &gitalypb.Repository{StorageName: cfg.Storages[0].Name, RelativePath: ""}},
			Code: codes.InvalidArgument,
		},
		{
			Desc: "Repository is nil",
			Req:  &gitalypb.SSHUploadArchiveRequest{Repository: nil},
			Code: codes.InvalidArgument,
		},
		{
			Desc: "Data exists on first request",
			Req:  &gitalypb.SSHUploadArchiveRequest{Repository: &gitalypb.Repository{StorageName: cfg.Storages[0].Name, RelativePath: "path/to/repo"}, Stdin: []byte("Fail")},
			Code: codes.InvalidArgument,
		},
	}

	for _, test := range tests {
		t.Run(test.Desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()
			stream, err := client.SSHUploadArchive(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if err = stream.Send(test.Req); err != nil {
				t.Fatal(err)
			}
			require.NoError(t, stream.CloseSend())

			err = testUploadArchiveFailedResponse(t, stream)
			testhelper.RequireGrpcError(t, err, test.Code)
		})
	}
}

func TestUploadArchiveSuccess(t *testing.T) {
	cfg, repo, _ := testcfg.BuildWithRepo(t)

	testhelper.BuildGitalySSH(t, cfg)

	serverSocketPath := runSSHServer(t, cfg)

	cmd := exec.Command(cfg.Git.BinPath, "archive", "master", "--remote=git@localhost:test/test.git")

	err := testArchive(t, cfg, serverSocketPath, repo, cmd)
	require.NoError(t, err)
}

func testArchive(t *testing.T, cfg config.Cfg, serverSocketPath string, testRepo *gitalypb.Repository, cmd *exec.Cmd) error {
	req := &gitalypb.SSHUploadArchiveRequest{Repository: testRepo}
	payload, err := protojson.Marshal(req)

	require.NoError(t, err)

	cmd.Env = []string{
		fmt.Sprintf("GITALY_ADDRESS=%s", serverSocketPath),
		fmt.Sprintf("GITALY_PAYLOAD=%s", payload),
		fmt.Sprintf("PATH=%s", ".:"+os.Getenv("PATH")),
		fmt.Sprintf(`GIT_SSH_COMMAND=%s upload-archive`, filepath.Join(cfg.BinDir, "gitaly-ssh")),
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %q", err, out)
	}
	if !cmd.ProcessState.Success() {
		return fmt.Errorf("Failed to run `git archive`: %q", out)
	}

	return nil
}

func testUploadArchiveFailedResponse(t *testing.T, stream gitalypb.SSHService_SSHUploadArchiveClient) error {
	var err error
	var res *gitalypb.SSHUploadArchiveResponse

	for err == nil {
		res, err = stream.Recv()
		require.Nil(t, res.GetStdout())
	}

	return err
}
