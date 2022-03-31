package hook

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/backchannel"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/quarantine"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/transaction"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitlab"
	"gitlab.com/gitlab-org/gitaly/v14/internal/metadata/featureflag"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/internal/transaction/txinfo"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func TestUpdate_customHooks(t *testing.T) {
	cfg, repo, repoPath := testcfg.BuildWithRepo(t)

	hookManager := NewManager(config.NewLocator(cfg), transaction.NewManager(cfg, backchannel.NewRegistry()), gitlab.NewMockClient(), cfg)

	receiveHooksPayload := &git.ReceiveHooksPayload{
		UserID:   "1234",
		Username: "user",
		Protocol: "web",
	}

	ctx, cleanup := testhelper.Context()
	defer cleanup()

	payload, err := git.NewHooksPayload(cfg, repo, nil, receiveHooksPayload, git.UpdateHook, featureflag.RawFromContext(ctx)).Env()
	require.NoError(t, err)

	primaryPayload, err := git.NewHooksPayload(
		cfg,
		repo,
		&txinfo.Transaction{
			ID: 1234, Node: "primary", Primary: true,
		},
		receiveHooksPayload,
		git.UpdateHook,
		featureflag.RawFromContext(ctx),
	).Env()
	require.NoError(t, err)

	secondaryPayload, err := git.NewHooksPayload(
		cfg,
		repo,
		&txinfo.Transaction{
			ID: 1234, Node: "secondary", Primary: false,
		},
		receiveHooksPayload,
		git.UpdateHook,
		featureflag.RawFromContext(ctx),
	).Env()
	require.NoError(t, err)

	hash1 := strings.Repeat("1", 40)
	hash2 := strings.Repeat("2", 40)

	testCases := []struct {
		desc           string
		env            []string
		hook           string
		reference      string
		oldHash        string
		newHash        string
		expectedErr    string
		expectedStdout string
		expectedStderr string
	}{
		{
			desc:      "hook receives environment variables",
			env:       []string{payload},
			reference: "refs/heads/master",
			oldHash:   hash1,
			newHash:   hash2,
			hook:      "#!/bin/sh\nenv | grep -e '^GL_' -e '^GITALY_' | sort\n",
			expectedStdout: strings.Join([]string{
				"GL_ID=1234",
				fmt.Sprintf("GL_PROJECT_PATH=%s", repo.GetGlProjectPath()),
				"GL_PROTOCOL=web",
				fmt.Sprintf("GL_REPOSITORY=%s", repo.GetGlRepository()),
				"GL_USERNAME=user",
			}, "\n") + "\n",
		},
		{
			desc:           "hook receives arguments",
			env:            []string{payload},
			reference:      "refs/heads/master",
			oldHash:        hash1,
			newHash:        hash2,
			hook:           "#!/bin/sh\nprintf '%s\\n' \"$@\"\n",
			expectedStdout: fmt.Sprintf("refs/heads/master\n%s\n%s\n", hash1, hash2),
		},
		{
			desc:           "stdout and stderr are passed through",
			env:            []string{payload},
			reference:      "refs/heads/master",
			oldHash:        hash1,
			newHash:        hash2,
			hook:           "#!/bin/sh\necho foo >&1\necho bar >&2\n",
			expectedStdout: "foo\n",
			expectedStderr: "bar\n",
		},
		{
			desc:      "standard input is empty",
			env:       []string{payload},
			reference: "refs/heads/master",
			oldHash:   hash1,
			newHash:   hash2,
			hook:      "#!/bin/sh\ncat\n",
		},
		{
			desc:        "invalid script causes failure",
			env:         []string{payload},
			reference:   "refs/heads/master",
			oldHash:     hash1,
			newHash:     hash2,
			hook:        "",
			expectedErr: "exec format error",
		},
		{
			desc:        "errors are passed through",
			env:         []string{payload},
			reference:   "refs/heads/master",
			oldHash:     hash1,
			newHash:     hash2,
			hook:        "#!/bin/sh\nexit 123\n",
			expectedErr: "exit status 123",
		},
		{
			desc:           "errors are passed through with stderr and stdout",
			env:            []string{payload},
			reference:      "refs/heads/master",
			oldHash:        hash1,
			newHash:        hash2,
			hook:           "#!/bin/sh\necho foo >&1\necho bar >&2\nexit 123\n",
			expectedStdout: "foo\n",
			expectedStderr: "bar\n",
			expectedErr:    "exit status 123",
		},
		{
			desc:           "hook is executed on primary",
			env:            []string{primaryPayload},
			reference:      "refs/heads/master",
			oldHash:        hash1,
			newHash:        hash2,
			hook:           "#!/bin/sh\necho foo\n",
			expectedStdout: "foo\n",
		},
		{
			desc:      "hook is not executed on secondary",
			env:       []string{secondaryPayload},
			reference: "refs/heads/master",
			oldHash:   hash1,
			newHash:   hash2,
			hook:      "#!/bin/sh\necho foo\n",
		},
		{
			desc:        "hook fails with missing reference",
			env:         []string{payload},
			oldHash:     hash1,
			newHash:     hash2,
			expectedErr: "hook got no reference",
		},
		{
			desc:        "hook fails with missing old value",
			env:         []string{payload},
			reference:   "refs/heads/master",
			newHash:     hash2,
			expectedErr: "hook got invalid old value",
		},
		{
			desc:        "hook fails with missing new value",
			env:         []string{payload},
			reference:   "refs/heads/master",
			oldHash:     hash1,
			expectedErr: "hook got invalid new value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gittest.WriteCustomHook(t, repoPath, "update", []byte(tc.hook))

			var stdout, stderr bytes.Buffer
			err = hookManager.UpdateHook(ctx, repo, tc.reference, tc.oldHash, tc.newHash, tc.env, &stdout, &stderr)

			if tc.expectedErr != "" {
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedStdout, stdout.String())
			require.Equal(t, tc.expectedStderr, stderr.String())
		})
	}
}

func TestUpdate_quarantine(t *testing.T) {
	ctx, cleanup := testhelper.Context()
	defer cleanup()

	cfg, repoProto, repoPath := testcfg.BuildWithRepo(t)

	quarantine, err := quarantine.New(ctx, repoProto, config.NewLocator(cfg))
	require.NoError(t, err)

	quarantinedRepo := localrepo.NewTestRepo(t, cfg, quarantine.QuarantinedRepo())
	blobID, err := quarantinedRepo.WriteBlob(ctx, "", strings.NewReader("allyourbasearebelongtous"))
	require.NoError(t, err)

	hookManager := NewManager(config.NewLocator(cfg), nil, gitlab.NewMockClient(), cfg)

	script := fmt.Sprintf("#!/bin/sh\n%s cat-file -p '%s' || true\n",
		cfg.Git.BinPath, blobID.String())
	gittest.WriteCustomHook(t, repoPath, "update", []byte(script))

	for repo, isQuarantined := range map[*gitalypb.Repository]bool{
		quarantine.QuarantinedRepo(): true,
		repoProto:                    false,
	} {
		t.Run(fmt.Sprintf("quarantined: %v", isQuarantined), func(t *testing.T) {
			env, err := git.NewHooksPayload(cfg, repo, nil,
				&git.ReceiveHooksPayload{
					UserID:   "1234",
					Username: "user",
					Protocol: "web",
				},
				git.PreReceiveHook,
				featureflag.RawFromContext(ctx),
			).Env()
			require.NoError(t, err)

			var stdout, stderr bytes.Buffer
			require.NoError(t, hookManager.UpdateHook(ctx, repo, "refs/heads/master",
				git.ZeroOID.String(), git.ZeroOID.String(), []string{env}, &stdout, &stderr))

			if isQuarantined {
				require.Equal(t, "allyourbasearebelongtous", stdout.String())
				require.Empty(t, stderr.String())
			} else {
				require.Empty(t, stdout.String())
				require.Contains(t, stderr.String(), "Not a valid object name")
			}
		})
	}
}
