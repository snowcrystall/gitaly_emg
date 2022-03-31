package git_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
)

func TestCheckRefFormat(t *testing.T) {
	cfg := testcfg.Build(t)

	ctx, cancel := testhelper.Context()
	defer cancel()

	gitCmdFactory := git.NewExecCommandFactory(cfg)

	for _, tc := range []struct {
		desc    string
		tagName string
		ok      bool
		err     error
	}{
		// Just trivial tests here, most of this is tested in
		// internal/gitaly/service/operations/tags_test.go
		{
			desc:    "unqualified name",
			tagName: "my-name",
			ok:      false,
			err:     git.CheckRefFormatError{},
		},
		{
			desc:    "fully-qualified name",
			tagName: "refs/heads/my-name",
			ok:      true,
			err:     nil,
		},
		{
			desc:    "basic tag",
			tagName: "refs/tags/my-tag",
			ok:      true,
			err:     nil,
		},
		{
			desc:    "invalid tag",
			tagName: "refs/tags/my tag",
			ok:      false,
			err:     git.CheckRefFormatError{},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ok, err := git.CheckRefFormat(ctx, gitCmdFactory, tc.tagName)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestReferenceName_NewReferenceNameFromBranchName(t *testing.T) {
	for _, tc := range []struct {
		desc      string
		reference string
		expected  string
	}{
		{
			desc:      "unqualified reference",
			reference: "master",
			expected:  "refs/heads/master",
		},
		{
			desc:      "partly qualified reference",
			reference: "heads/master",
			expected:  "refs/heads/heads/master",
		},
		{
			desc:      "fully qualified reference",
			reference: "refs/heads/master",
			expected:  "refs/heads/refs/heads/master",
		},
		{
			desc:      "weird branch name",
			reference: "refs/master",
			expected:  "refs/heads/refs/master",
		},
		{
			desc:      "tag is treated as a branch",
			reference: "refs/tags/master",
			expected:  "refs/heads/refs/tags/master",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ref := git.NewReferenceNameFromBranchName(tc.reference)
			require.Equal(t, ref.String(), tc.expected)
		})
	}
}

func TestReferenceName_Branch(t *testing.T) {
	for _, tc := range []struct {
		desc      string
		reference string
		expected  string
	}{
		{
			desc:      "fully qualified reference",
			reference: "refs/heads/master",
			expected:  "master",
		},
		{
			desc:      "nested branch",
			reference: "refs/heads/foo/master",
			expected:  "foo/master",
		},
		{
			desc:      "unqualified branch is not a branch",
			reference: "master",
			expected:  "",
		},
		{
			desc:      "tag is not a branch",
			reference: "refs/tags/master",
			expected:  "",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			branch, ok := git.ReferenceName(tc.reference).Branch()
			require.Equal(t, tc.expected, branch)
			require.Equal(t, tc.expected != "", ok)
		})
	}
}
