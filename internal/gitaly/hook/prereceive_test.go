package hook

import (
	"bytes"
	"context"
	"errors"
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
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/metadata/featureflag"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/internal/transaction/txinfo"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func TestPrereceive_customHooks(t *testing.T) {
	cfg, repo, repoPath := testcfg.BuildWithRepo(t)

	hookManager := NewManager(config.NewLocator(cfg), transaction.NewManager(cfg, backchannel.NewRegistry()), gitlab.NewMockClient(), cfg)

	receiveHooksPayload := &git.ReceiveHooksPayload{
		UserID:   "1234",
		Username: "user",
		Protocol: "web",
	}

	ctx, cleanup := testhelper.Context()
	defer cleanup()

	payload, err := git.NewHooksPayload(cfg, repo, nil, receiveHooksPayload, git.PreReceiveHook, featureflag.RawFromContext(ctx)).Env()
	require.NoError(t, err)

	primaryPayload, err := git.NewHooksPayload(
		cfg,
		repo,
		&txinfo.Transaction{
			ID: 1234, Node: "primary", Primary: true,
		},
		receiveHooksPayload,
		git.PreReceiveHook,
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
		git.PreReceiveHook,
		featureflag.RawFromContext(ctx),
	).Env()
	require.NoError(t, err)

	testCases := []struct {
		desc           string
		env            []string
		pushOptions    []string
		hook           string
		stdin          string
		expectedErr    string
		expectedStdout string
		expectedStderr string
	}{
		{
			desc:  "hook receives environment variables",
			env:   []string{payload},
			hook:  "#!/bin/sh\nenv | grep -e '^GL_' -e '^GITALY_' | sort\n",
			stdin: "change\n",
			expectedStdout: strings.Join([]string{
				"GL_ID=1234",
				fmt.Sprintf("GL_PROJECT_PATH=%s", repo.GetGlProjectPath()),
				"GL_PROTOCOL=web",
				fmt.Sprintf("GL_REPOSITORY=%s", repo.GetGlRepository()),
				"GL_USERNAME=user",
			}, "\n") + "\n",
		},
		{
			desc:        "hook receives push options",
			env:         []string{payload},
			pushOptions: []string{"mr.create", "mr.merge_when_pipeline_succeeds"},
			hook:        "#!/bin/sh\nenv | grep -e '^GL_' -e '^GITALY_' -e '^GIT_PUSH_' | sort\n",
			stdin:       "change\n",
			expectedStdout: strings.Join([]string{
				"GIT_PUSH_OPTION_0=mr.create",
				"GIT_PUSH_OPTION_1=mr.merge_when_pipeline_succeeds",
				"GIT_PUSH_OPTION_COUNT=2",
				"GL_ID=1234",
				fmt.Sprintf("GL_PROJECT_PATH=%s", repo.GetGlProjectPath()),
				"GL_PROTOCOL=web",
				fmt.Sprintf("GL_REPOSITORY=%s", repo.GetGlRepository()),
				"GL_USERNAME=user",
			}, "\n") + "\n",
		},
		{
			desc:           "hook can write to stderr and stdout",
			env:            []string{payload},
			hook:           "#!/bin/sh\necho foo >&1 && echo bar >&2\n",
			stdin:          "change\n",
			expectedStdout: "foo\n",
			expectedStderr: "bar\n",
		},
		{
			desc:           "hook receives standard input",
			env:            []string{payload},
			hook:           "#!/bin/sh\ncat\n",
			stdin:          "foo\n",
			expectedStdout: "foo\n",
		},
		{
			desc:           "hook succeeds without consuming stdin",
			env:            []string{payload},
			hook:           "#!/bin/sh\necho foo\n",
			stdin:          "ignore me\n",
			expectedStdout: "foo\n",
		},
		{
			desc:        "invalid hook results in error",
			env:         []string{payload},
			hook:        "",
			stdin:       "change\n",
			expectedErr: "exec format error",
		},
		{
			desc:        "failing hook results in error",
			env:         []string{payload},
			hook:        "#!/bin/sh\nexit 123",
			stdin:       "change\n",
			expectedErr: "exit status 123",
		},
		{
			desc:           "hook is executed on primary",
			env:            []string{primaryPayload},
			hook:           "#!/bin/sh\necho foo\n",
			stdin:          "change\n",
			expectedStdout: "foo\n",
		},
		{
			desc:  "hook is not executed on secondary",
			env:   []string{secondaryPayload},
			hook:  "#!/bin/sh\necho foo\n",
			stdin: "change\n",
		},
		{
			desc:        "missing changes cause error",
			env:         []string{payload},
			expectedErr: "hook got no reference updates",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gittest.WriteCustomHook(t, repoPath, "pre-receive", []byte(tc.hook))

			var stdout, stderr bytes.Buffer
			err = hookManager.PreReceiveHook(ctx, repo, tc.pushOptions, tc.env, strings.NewReader(tc.stdin), &stdout, &stderr)

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

func TestPrereceive_quarantine(t *testing.T) {
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
	gittest.WriteCustomHook(t, repoPath, "pre-receive", []byte(script))

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

			stdin := strings.NewReader(fmt.Sprintf("%s %s refs/heads/master",
				git.ZeroOID, git.ZeroOID))

			var stdout, stderr bytes.Buffer
			require.NoError(t, hookManager.PreReceiveHook(ctx, repo, nil,
				[]string{env}, stdin, &stdout, &stderr))

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

type prereceiveAPIMock struct {
	allowed    func(context.Context, gitlab.AllowedParams) (bool, string, error)
	prereceive func(context.Context, string) (bool, error)
}

func (m *prereceiveAPIMock) Allowed(ctx context.Context, params gitlab.AllowedParams) (bool, string, error) {
	return m.allowed(ctx, params)
}

func (m *prereceiveAPIMock) PreReceive(ctx context.Context, glRepository string) (bool, error) {
	return m.prereceive(ctx, glRepository)
}

func (m *prereceiveAPIMock) Check(ctx context.Context) (*gitlab.CheckInfo, error) {
	return nil, errors.New("unexpected call")
}

func (m *prereceiveAPIMock) PostReceive(context.Context, string, string, string, ...string) (bool, []gitlab.PostReceiveMessage, error) {
	return true, nil, errors.New("unexpected call")
}

func TestPrereceive_gitlab(t *testing.T) {
	cfg, repo, repoPath := testcfg.BuildWithRepo(t)

	payload, err := git.NewHooksPayload(cfg, repo, nil, &git.ReceiveHooksPayload{
		UserID:   "1234",
		Username: "user",
		Protocol: "web",
	}, git.PreReceiveHook, nil).Env()
	require.NoError(t, err)

	standardEnv := []string{payload}

	testCases := []struct {
		desc           string
		env            []string
		changes        string
		allowed        func(*testing.T, context.Context, gitlab.AllowedParams) (bool, string, error)
		prereceive     func(*testing.T, context.Context, string) (bool, error)
		expectHookCall bool
		expectedErr    error
	}{
		{
			desc:    "allowed change",
			env:     standardEnv,
			changes: "changes\n",
			allowed: func(t *testing.T, ctx context.Context, params gitlab.AllowedParams) (bool, string, error) {
				require.Equal(t, repoPath, params.RepoPath)
				require.Equal(t, repo.GlRepository, params.GLRepository)
				require.Equal(t, "1234", params.GLID)
				require.Equal(t, "web", params.GLProtocol)
				require.Equal(t, "changes\n", params.Changes)
				return true, "", nil
			},
			prereceive: func(t *testing.T, ctx context.Context, glRepo string) (bool, error) {
				require.Equal(t, repo.GlRepository, glRepo)
				return true, nil
			},
			expectHookCall: true,
		},
		{
			desc:    "disallowed change",
			env:     standardEnv,
			changes: "changes\n",
			allowed: func(t *testing.T, ctx context.Context, params gitlab.AllowedParams) (bool, string, error) {
				return false, "you shall not pass", nil
			},
			expectHookCall: false,
			expectedErr:    NotAllowedError{Message: "you shall not pass"},
		},
		{
			desc:    "allowed returns error",
			env:     standardEnv,
			changes: "changes\n",
			allowed: func(t *testing.T, ctx context.Context, params gitlab.AllowedParams) (bool, string, error) {
				return false, "", errors.New("oops")
			},
			expectHookCall: false,
			expectedErr:    NotAllowedError{Message: "GitLab: oops"},
		},
		{
			desc:    "prereceive rejects",
			env:     standardEnv,
			changes: "changes\n",
			allowed: func(t *testing.T, ctx context.Context, params gitlab.AllowedParams) (bool, string, error) {
				return true, "", nil
			},
			prereceive: func(t *testing.T, ctx context.Context, glRepo string) (bool, error) {
				return false, nil
			},
			expectHookCall: true,
			expectedErr:    errors.New(""),
		},
		{
			desc:    "prereceive errors",
			env:     standardEnv,
			changes: "changes\n",
			allowed: func(t *testing.T, ctx context.Context, params gitlab.AllowedParams) (bool, string, error) {
				return true, "", nil
			},
			prereceive: func(t *testing.T, ctx context.Context, glRepo string) (bool, error) {
				return false, errors.New("prereceive oops")
			},
			expectHookCall: true,
			expectedErr:    helper.ErrInternalf("calling pre_receive endpoint: %v", errors.New("prereceive oops")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cleanup := testhelper.Context()
			defer cleanup()

			gitlabAPI := prereceiveAPIMock{
				allowed: func(ctx context.Context, params gitlab.AllowedParams) (bool, string, error) {
					return tc.allowed(t, ctx, params)
				},
				prereceive: func(ctx context.Context, glRepo string) (bool, error) {
					return tc.prereceive(t, ctx, glRepo)
				},
			}

			hookManager := NewManager(config.NewLocator(cfg), transaction.NewManager(cfg, backchannel.NewRegistry()), &gitlabAPI, cfg)

			gittest.WriteCustomHook(t, repoPath, "pre-receive", []byte("#!/bin/sh\necho called\n"))

			var stdout, stderr bytes.Buffer
			err = hookManager.PreReceiveHook(ctx, repo, nil, tc.env, strings.NewReader(tc.changes), &stdout, &stderr)

			if tc.expectedErr != nil {
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			if tc.expectHookCall {
				require.Equal(t, "called\n", stdout.String())
			} else {
				require.Empty(t, stdout.String())
			}
			require.Empty(t, stderr.String())
		})
	}
}
