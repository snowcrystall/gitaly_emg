package repository

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
)

func TestWriteRefSuccessful(t *testing.T) {
	cfg, repo, repoPath, client := setupRepositoryService(t)

	testCases := []struct {
		desc string
		req  *gitalypb.WriteRefRequest
	}{
		{
			desc: "shell update HEAD to refs/heads/master",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Ref:        []byte("HEAD"),
				Revision:   []byte("refs/heads/master"),
			},
		},
		{
			desc: "shell update refs/heads/master",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Ref:        []byte("refs/heads/master"),
				Revision:   []byte("b83d6e391c22777fca1ed3012fce84f633d7fed0"),
			},
		},
		{
			desc: "shell update refs/heads/master w/ validation",
			req: &gitalypb.WriteRefRequest{
				Repository:  repo,
				Ref:         []byte("refs/heads/master"),
				Revision:    []byte("498214de67004b1da3d820901307bed2a68a8ef6"),
				OldRevision: []byte("b83d6e391c22777fca1ed3012fce84f633d7fed0"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()
			_, err := client.WriteRef(ctx, tc.req)

			require.NoError(t, err)

			if bytes.Equal(tc.req.Ref, []byte("HEAD")) {
				content := testhelper.MustReadFile(t, filepath.Join(repoPath, "HEAD"))

				refRevision := bytes.Join([][]byte{[]byte("ref: "), tc.req.Revision, []byte("\n")}, nil)

				require.EqualValues(t, content, refRevision)
				return
			}
			rev := gittest.Exec(t, cfg, "--git-dir", repoPath, "log", "--pretty=%H", "-1", string(tc.req.Ref))

			rev = bytes.Replace(rev, []byte("\n"), nil, 1)

			require.Equal(t, string(tc.req.Revision), string(rev))
		})
	}
}

func TestWriteRefValidationError(t *testing.T) {
	_, repo, _, client := setupRepositoryService(t)

	testCases := []struct {
		desc string
		req  *gitalypb.WriteRefRequest
	}{
		{
			desc: "empty revision",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Ref:        []byte("refs/heads/master"),
			},
		},
		{
			desc: "empty ref name",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Revision:   []byte("498214de67004b1da3d820901307bed2a68a8ef6"),
			},
		},
		{
			desc: "non-prefixed ref name for shell",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Ref:        []byte("master"),
				Revision:   []byte("498214de67004b1da3d820901307bed2a68a8ef6"),
			},
		},
		{
			desc: "revision contains \\x00",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Ref:        []byte("refs/heads/master"),
				Revision:   []byte("012301230123\x001243"),
			},
		},
		{
			desc: "ref contains \\x00",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Ref:        []byte("refs/head\x00s/master\x00"),
				Revision:   []byte("0123012301231243"),
			},
		},
		{
			desc: "ref contains whitespace",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Ref:        []byte("refs/heads /master"),
				Revision:   []byte("0123012301231243"),
			},
		},
		{
			desc: "invalid revision",
			req: &gitalypb.WriteRefRequest{
				Repository: repo,
				Ref:        []byte("refs/heads/master"),
				Revision:   []byte("--output=/meow"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()
			_, err := client.WriteRef(ctx, tc.req)

			testhelper.RequireGrpcError(t, err, codes.InvalidArgument)
		})
	}
}
