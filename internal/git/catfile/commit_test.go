package catfile

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testassert"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseCommit(t *testing.T) {
	info := &ObjectInfo{
		Oid:  "a984dfa4dee018c6d5f5f57ffec0d0e22763df16",
		Type: "commit",
	}

	// Valid-but-interesting commits should be test at the FindCommit level.
	// Invalid objects (that Git would complain about during fsck) can be
	// tested here.
	//
	// Once a repository contains a pathological object it can be hard to get
	// rid of it. Because of this I think it's nicer to ignore such objects
	// than to throw hard errors.
	testCases := []struct {
		desc string
		in   []byte
		out  *gitalypb.GitCommit
	}{
		{
			desc: "empty commit object",
			in:   []byte{},
			out:  &gitalypb.GitCommit{Id: info.Oid.String()},
		},
		{
			desc: "no email",
			in:   []byte("author Jane Doe"),
			out: &gitalypb.GitCommit{
				Id:     info.Oid.String(),
				Author: &gitalypb.CommitAuthor{Name: []byte("Jane Doe")},
			},
		},
		{
			desc: "unmatched <",
			in:   []byte("author Jane Doe <janedoe@example.com"),
			out: &gitalypb.GitCommit{
				Id:     info.Oid.String(),
				Author: &gitalypb.CommitAuthor{Name: []byte("Jane Doe")},
			},
		},
		{
			desc: "unmatched >",
			in:   []byte("author Jane Doe janedoe@example.com>"),
			out: &gitalypb.GitCommit{
				Id:     info.Oid.String(),
				Author: &gitalypb.CommitAuthor{Name: []byte("Jane Doe janedoe@example.com>")},
			},
		},
		{
			desc: "missing date",
			in:   []byte("author Jane Doe <janedoe@example.com> "),
			out: &gitalypb.GitCommit{
				Id:     info.Oid.String(),
				Author: &gitalypb.CommitAuthor{Name: []byte("Jane Doe"), Email: []byte("janedoe@example.com")},
			},
		},
		{
			desc: "date too high",
			in:   []byte("author Jane Doe <janedoe@example.com> 9007199254740993 +0200"),
			out: &gitalypb.GitCommit{
				Id: info.Oid.String(),
				Author: &gitalypb.CommitAuthor{
					Name:     []byte("Jane Doe"),
					Email:    []byte("janedoe@example.com"),
					Date:     &timestamppb.Timestamp{Seconds: 9223371974719179007},
					Timezone: []byte("+0200"),
				},
			},
		},
		{
			desc: "date negative",
			in:   []byte("author Jane Doe <janedoe@example.com> -1 +0200"),
			out: &gitalypb.GitCommit{
				Id: info.Oid.String(),
				Author: &gitalypb.CommitAuthor{
					Name:     []byte("Jane Doe"),
					Email:    []byte("janedoe@example.com"),
					Date:     &timestamppb.Timestamp{Seconds: 9223371974719179007},
					Timezone: []byte("+0200"),
				},
			},
		},
		{
			desc: "huge",
			in:   append([]byte("author "), bytes.Repeat([]byte("A"), 100000)...),
			out: &gitalypb.GitCommit{
				Id: info.Oid.String(),
				Author: &gitalypb.CommitAuthor{
					Name: bytes.Repeat([]byte("A"), 100000),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			info.Size = int64(len(tc.in))
			out, err := ParseCommit(bytes.NewBuffer(tc.in), info.Oid)
			require.NoError(t, err, "parse error")
			require.Equal(t, tc.out, out)
		})
	}
}

func TestGetCommit(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	_, c, _ := setupBatch(t, ctx)

	ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

	const commitSha = "2d1db523e11e777e49377cfb22d368deec3f0793"
	const commitMsg = "Correct test_env.rb path for adding branch\n"
	const blobSha = "c60514b6d3d6bf4bec1030f70026e34dfbd69ad5"

	testCases := []struct {
		desc     string
		revision string
		errStr   string
	}{
		{
			desc:     "commit",
			revision: commitSha,
		},
		{
			desc:     "not existing commit",
			revision: "not existing revision",
			errStr:   "object not found",
		},
		{
			desc:     "blob sha",
			revision: blobSha,
			errStr:   "object not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			c, err := GetCommit(ctx, c, git.Revision(tc.revision))

			if tc.errStr == "" {
				require.NoError(t, err)
				require.Equal(t, commitMsg, string(c.Body))
			} else {
				require.EqualError(t, err, tc.errStr)
			}
		})
	}
}

func TestGetCommitWithTrailers(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	cfg, c, testRepo := setupBatch(t, ctx)

	ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

	commit, err := GetCommitWithTrailers(ctx, git.NewExecCommandFactory(cfg), testRepo, c, "5937ac0a7beb003549fc5fd26fc247adbce4a52e")

	require.NoError(t, err)

	require.Equal(t, commit.Trailers, []*gitalypb.CommitTrailer{
		{
			Key:   []byte("Signed-off-by"),
			Value: []byte("Dmitriy Zaporozhets <dmitriy.zaporozhets@gmail.com>"),
		},
	})
}

func TestParseCommitAuthor(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		author   string
		expected *gitalypb.CommitAuthor
	}{
		{
			desc:     "empty author",
			author:   "",
			expected: &gitalypb.CommitAuthor{},
		},
		{
			desc:   "normal author",
			author: "Au Thor <au.thor@example.com> 1625121079 +0000",
			expected: &gitalypb.CommitAuthor{
				Name:     []byte("Au Thor"),
				Email:    []byte("au.thor@example.com"),
				Date:     timestamppb.New(time.Unix(1625121079, 0)),
				Timezone: []byte("+0000"),
			},
		},
		{
			desc:   "author with missing mail",
			author: "Au Thor <> 1625121079 +0000",
			expected: &gitalypb.CommitAuthor{
				Name:     []byte("Au Thor"),
				Date:     timestamppb.New(time.Unix(1625121079, 0)),
				Timezone: []byte("+0000"),
			},
		},
		{
			desc:   "author with missing date",
			author: "Au Thor <au.thor@example.com>",
			expected: &gitalypb.CommitAuthor{
				Name:  []byte("Au Thor"),
				Email: []byte("au.thor@example.com"),
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			testassert.ProtoEqual(t, tc.expected, parseCommitAuthor(tc.author))
		})
	}
}
