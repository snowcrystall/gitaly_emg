package commit

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testassert"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
)

func TestSuccessfulGetCommitSignaturesRequest(t *testing.T) {
	t.Parallel()
	cfg, repo, repoPath, client := setupCommitServiceWithRepo(t, true)

	commitData := testhelper.MustReadFile(t, "testdata/dc00eb001f41dfac08192ead79c2377c588b82ee.commit")
	commit := text.ChompBytes(gittest.ExecStream(t, cfg, bytes.NewReader(commitData), "-C", repoPath, "hash-object", "-w", "-t", "commit", "--stdin", "--literally"))
	require.Equal(t, "dc00eb001f41dfac08192ead79c2377c588b82ee", commit)

	ctx, cancel := testhelper.Context()
	defer cancel()

	request := &gitalypb.GetCommitSignaturesRequest{
		Repository: repo,
		CommitIds: []string{
			"5937ac0a7beb003549fc5fd26fc247adbce4a52e", // has signature
			"e63f41fe459e62e1228fcef60d7189127aeba95a", // has no signature
			git.ZeroOID.String(),                       // does not exist
			"a17a9f66543673edf0a3d1c6b93bdda3fe600f32", // has signature
			"8cf8e80a5a0546e391823c250f2b26b9cf15ce88", // has signature and commit message > 4MB
			"dc00eb001f41dfac08192ead79c2377c588b82ee", // has signature and commit message without newline at the end
		},
	}

	expectedSignatures := []*gitalypb.GetCommitSignaturesResponse{
		{
			CommitId:   "5937ac0a7beb003549fc5fd26fc247adbce4a52e",
			Signature:  testhelper.MustReadFile(t, "testdata/commit-5937ac0a7beb003549fc5fd26fc247adbce4a52e-signature"),
			SignedText: testhelper.MustReadFile(t, "testdata/commit-5937ac0a7beb003549fc5fd26fc247adbce4a52e-signed-text"),
		},
		{
			CommitId:   "a17a9f66543673edf0a3d1c6b93bdda3fe600f32",
			Signature:  testhelper.MustReadFile(t, "testdata/gitlab-test-commit-a17a9f66543673edf0a3d1c6b93bdda3fe600f32-signature"),
			SignedText: testhelper.MustReadFile(t, "testdata/gitlab-test-commit-a17a9f66543673edf0a3d1c6b93bdda3fe600f32-signed-text"),
		},
		{
			CommitId:   "8cf8e80a5a0546e391823c250f2b26b9cf15ce88",
			Signature:  testhelper.MustReadFile(t, "testdata/gitaly-test-commit-8cf8e80a5a0546e391823c250f2b26b9cf15ce88-signature"),
			SignedText: testhelper.MustReadFile(t, "testdata/gitaly-test-commit-8cf8e80a5a0546e391823c250f2b26b9cf15ce88-signed-text"),
		},
		{
			CommitId:   "dc00eb001f41dfac08192ead79c2377c588b82ee",
			Signature:  testhelper.MustReadFile(t, "testdata/dc00eb001f41dfac08192ead79c2377c588b82ee-signed-no-newline-signature.txt"),
			SignedText: testhelper.MustReadFile(t, "testdata/dc00eb001f41dfac08192ead79c2377c588b82ee-signed-no-newline-signed-text.txt"),
		},
	}

	c, err := client.GetCommitSignatures(ctx, request)
	require.NoError(t, err)

	fetchedSignatures := readAllSignaturesFromClient(t, c)

	require.Len(t, fetchedSignatures, len(expectedSignatures))
	for i, expected := range expectedSignatures {
		testassert.ProtoEqual(t, expected, fetchedSignatures[i])
	}
}

func TestFailedGetCommitSignaturesRequest(t *testing.T) {
	t.Parallel()
	_, repo, _, client := setupCommitServiceWithRepo(t, true)

	testCases := []struct {
		desc    string
		request *gitalypb.GetCommitSignaturesRequest
		code    codes.Code
	}{
		{
			desc: "empty Repository",
			request: &gitalypb.GetCommitSignaturesRequest{
				Repository: nil,
				CommitIds:  []string{"5937ac0a7beb003549fc5fd26fc247adbce4a52e"},
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty CommitIds",
			request: &gitalypb.GetCommitSignaturesRequest{
				Repository: repo,
				CommitIds:  []string{},
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "commitIDS with shorthand sha",
			request: &gitalypb.GetCommitSignaturesRequest{
				Repository: repo,
				CommitIds:  []string{"5937ac0a7beb003549fc5fd26fc247adbce4a52e", "a17a9f6"},
			},
			code: codes.InvalidArgument,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()

			c, err := client.GetCommitSignatures(ctx, testCase.request)
			require.NoError(t, err)

			for {
				_, err = c.Recv()
				if err != nil {
					break
				}
			}

			testhelper.RequireGrpcError(t, err, testCase.code)
		})
	}
}

func readAllSignaturesFromClient(t *testing.T, c gitalypb.CommitService_GetCommitSignaturesClient) (signatures []*gitalypb.GetCommitSignaturesResponse) {
	t.Helper()

	for {
		resp, err := c.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		if resp.CommitId != "" {
			signatures = append(signatures, resp)
			// first message contains either Signature or SignedText so no need to append anything
			continue
		}

		currentSignature := signatures[len(signatures)-1]

		if len(resp.Signature) != 0 {
			currentSignature.Signature = append(currentSignature.Signature, resp.Signature...)
		} else if len(resp.SignedText) != 0 {
			currentSignature.SignedText = append(currentSignature.SignedText, resp.SignedText...)
		}
	}

	return
}
