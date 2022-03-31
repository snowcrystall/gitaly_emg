package operations

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/backchannel"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/hook"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testassert"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testserver"
	"gitlab.com/gitlab-org/gitaly/v14/internal/transaction/txinfo"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSuccessfulUserDeleteTagRequest(t *testing.T) {
	t.Parallel()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	tagNameInput := "to-be-deleted-soon-tag"

	gittest.Exec(t, cfg, "-C", repoPath, "tag", tagNameInput)

	request := &gitalypb.UserDeleteTagRequest{
		Repository: repo,
		TagName:    []byte(tagNameInput),
		User:       gittest.TestUser,
	}

	_, err := client.UserDeleteTag(ctx, request)
	require.NoError(t, err)

	tags := gittest.Exec(t, cfg, "-C", repoPath, "tag")
	require.NotContains(t, string(tags), tagNameInput, "tag name still exists in tags list")
}

func TestSuccessfulGitHooksForUserDeleteTagRequest(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	tagNameInput := "to-be-déleted-soon-tag"

	request := &gitalypb.UserDeleteTagRequest{
		Repository: repo,
		TagName:    []byte(tagNameInput),
		User:       gittest.TestUser,
	}

	for _, hookName := range GitlabHooks {
		t.Run(hookName, func(t *testing.T) {
			gittest.Exec(t, cfg, "-C", repoPath, "tag", tagNameInput)

			hookOutputTempPath := gittest.WriteEnvToCustomHook(t, repoPath, hookName)

			_, err := client.UserDeleteTag(ctx, request)
			require.NoError(t, err)

			output := testhelper.MustReadFile(t, hookOutputTempPath)
			require.Contains(t, string(output), "GL_USERNAME="+gittest.TestUser.GlUsername)
		})
	}
}

func writeAssertObjectTypePreReceiveHook(t *testing.T, cfg config.Cfg) string {
	t.Helper()

	hook := fmt.Sprintf(`#!/usr/bin/env ruby

# We match a non-ASCII ref_name below
Encoding.default_external = Encoding::UTF_8
Encoding.default_internal = Encoding::UTF_8

expected_object_type = ARGV.shift
commands = STDIN.each_line.map(&:chomp)
unless commands.size == 1
  abort "expected 1 ref update command, got #{commands.size}"
end

old_value, new_value, ref_name = commands[0].split(' ', 3)
abort 'missing new_value' unless new_value

out = IO.popen(%%W[%s cat-file -t #{new_value}], &:read)
abort 'cat-file failed' unless $?.success?

if ref_name =~ /^refs\/[^\/]+\/skip-type-check-/
  exit 0
end

unless out.chomp == expected_object_type
  abort "pre-receive hook error: expected '#{ref_name}' update of '#{old_value}' (a) -> '#{new_value}' (b) for 'b' to be a '#{expected_object_type}' object, got '#{out}'"
end`, cfg.Git.BinPath)

	dir := testhelper.TempDir(t)
	hookPath := filepath.Join(dir, "pre-receive")

	require.NoError(t, ioutil.WriteFile(hookPath, []byte(hook), 0755))

	return hookPath
}

func writeAssertObjectTypeUpdateHook(t *testing.T, cfg config.Cfg) string {
	t.Helper()

	hook := fmt.Sprintf(`#!/usr/bin/env ruby

# We match a non-ASCII ref_name below
Encoding.default_external = Encoding::UTF_8
Encoding.default_internal = Encoding::UTF_8

expected_object_type = ARGV.shift
ref_name, old_value, new_value = ARGV[0..2]

abort "missing new_value" unless new_value

out = IO.popen(%%W[%s cat-file -t #{new_value}], &:read)
abort 'cat-file failed' unless $?.success?

if ref_name =~ /^refs\/[^\/]+\/skip-type-check-/
  exit 0
end

unless out.chomp == expected_object_type
  abort "update hook error: expected '#{ref_name}' update of '#{old_value}' (a) -> '#{new_value}' (b) for 'b' to be a '#{expected_object_type}' object, got '#{out}'"
end`, cfg.Git.BinPath)

	dir := testhelper.TempDir(t)
	hookPath := filepath.Join(dir, "pre-receive")

	require.NoError(t, ioutil.WriteFile(hookPath, []byte(hook), 0755))

	return hookPath
}

func TestSuccessfulUserCreateTagRequest(t *testing.T) {
	t.Parallel()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	targetRevision := "c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd"
	targetRevisionCommit, err := repo.ReadCommit(ctx, git.Revision(targetRevision))
	require.NoError(t, err)

	inputTagName := "to-be-créated-soon"

	preReceiveHook := writeAssertObjectTypePreReceiveHook(t, cfg)

	updateHook := writeAssertObjectTypeUpdateHook(t, cfg)

	testCases := []struct {
		desc               string
		tagName            string
		message            string
		targetRevision     string
		expectedTag        *gitalypb.Tag
		expectedObjectType string
	}{
		{
			desc:           "lightweight tag to commit",
			tagName:        inputTagName,
			targetRevision: targetRevision,
			expectedTag: &gitalypb.Tag{
				Name:         []byte(inputTagName),
				Id:           targetRevision,
				TargetCommit: targetRevisionCommit,
			},
			expectedObjectType: "commit",
		},
		{
			desc:           "annotated tag to commit",
			tagName:        inputTagName,
			targetRevision: targetRevision,
			message:        "This is an annotated tag",
			expectedTag: &gitalypb.Tag{
				Name: []byte(inputTagName),
				//Id: is a new object, filled in below
				TargetCommit: targetRevisionCommit,
				Message:      []byte("This is an annotated tag"),
				MessageSize:  24,
			},
			expectedObjectType: "tag",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			for hook, content := range map[string]string{
				"pre-receive": fmt.Sprintf("#!/bin/sh\n%s %s \"$@\"", preReceiveHook, testCase.expectedObjectType),
				"update":      fmt.Sprintf("#!/bin/sh\n%s %s \"$@\"", updateHook, testCase.expectedObjectType),
			} {
				gittest.WriteCustomHook(t, repoPath, hook, []byte(content))
			}

			request := &gitalypb.UserCreateTagRequest{
				Repository:     repoProto,
				TagName:        []byte(testCase.tagName),
				TargetRevision: []byte(testCase.targetRevision),
				User:           gittest.TestUser,
				Message:        []byte(testCase.message),
			}

			response, err := client.UserCreateTag(ctx, request)
			require.NoError(t, err, "error from calling RPC")
			require.Empty(t, response.PreReceiveError, "PreReceiveError must be empty, signalling the push was accepted")

			defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", inputTagName)

			responseOk := &gitalypb.UserCreateTagResponse{
				Tag: testCase.expectedTag,
			}
			// Fake up *.Id for annotated tags
			if len(testCase.expectedTag.Id) == 0 {
				id := gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", inputTagName)
				responseOk.Tag.Id = text.ChompBytes(id)
			}

			testassert.ProtoEqual(t, responseOk, response)

			tag := gittest.Exec(t, cfg, "-C", repoPath, "tag")
			require.Contains(t, string(tag), inputTagName)
		})
	}
}

func TestUserCreateTagWithTransaction(t *testing.T) {
	t.Parallel()
	cfg, repoProto, repoPath := testcfg.BuildWithRepo(t)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	hooksOutputDir := testhelper.TempDir(t)
	hooksOutputPath := filepath.Join(hooksOutputDir, "output")

	// We're creating a set of custom hooks which simply
	// write to a file. The intention is that we want to
	// check that the hooks only run on the primary node.
	hooks := []string{"pre-receive", "update", "post-receive"}
	for _, hook := range hooks {
		gittest.WriteCustomHook(t, repoPath, hook,
			[]byte(fmt.Sprintf("#!/bin/sh\necho %s >>%s\n", hook, hooksOutputPath)),
		)
	}

	// We're creating a custom server with a fake transaction server which
	// simply returns success for every call, but tracks the number of
	// calls. The server is then injected into the client's context to make
	// it available for transactional voting. We cannot use
	// runOperationServiceServer as it puts a Praefect server in between if
	// running Praefect tests, which would break our test setup.
	transactionServer := &testTransactionServer{}
	testserver.RunGitalyServer(t, cfg, nil, func(srv *grpc.Server, deps *service.Dependencies) {
		gitalypb.RegisterOperationServiceServer(srv, NewServer(
			deps.GetCfg(),
			nil,
			deps.GetHookManager(),
			deps.GetLocator(),
			deps.GetConnsPool(),
			deps.GetGitCmdFactory(),
			deps.GetCatfileCache(),
		))
		gitalypb.RegisterHookServiceServer(srv, hook.NewServer(deps.GetCfg(), deps.GetHookManager(), deps.GetGitCmdFactory(), deps.GetPackObjectsCache()))
	})

	ctx, cancel := testhelper.Context()
	defer cancel()

	// We're using internal gitaly socket to connect to the server.
	// This is kind of a hack when running tests with Praefect:
	// if we directly connect to the server created above, then our call
	// would be intercepted by Praefect, which would in turn replace the
	// transaction information we inject further down below. So we instead
	// use internal socket so we can circumvent Praefect and just talk
	// to Gitaly directly.
	client := newMuxedOperationClient(t, ctx, "unix://"+cfg.GitalyInternalSocketPath(), cfg.Auth.Token,
		backchannel.NewClientHandshaker(
			testhelper.DiscardTestEntry(t),
			func() backchannel.Server {
				srv := grpc.NewServer()
				gitalypb.RegisterRefTransactionServer(srv, transactionServer)
				return srv
			},
		),
	)

	for i, testCase := range []struct {
		desc    string
		primary bool
		message string
	}{
		{
			desc:    "primary creates a lightweight tag",
			primary: true,
		},
		{
			desc:    "secondary creates a lightweight tag",
			primary: false,
		},
		{
			desc:    "primary creates an annotated tag",
			primary: true,
			message: "foobar",
		},
		{
			desc:    "secondary creates an annotated tag",
			primary: false,
			message: "foobar",
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			*transactionServer = testTransactionServer{}

			if err := os.Remove(hooksOutputPath); err != nil {
				require.True(t, os.IsNotExist(err), "error when cleaning up work area: %v", err)
			}

			tagName := fmt.Sprintf("tag-%d", i)
			targetRevision := "c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd"
			targetCommit, err := repo.ReadCommit(ctx, git.Revision(targetRevision))
			require.NoError(t, err)

			request := &gitalypb.UserCreateTagRequest{
				Repository:     repoProto,
				TagName:        []byte(tagName),
				Message:        []byte(testCase.message),
				TargetRevision: []byte(targetRevision),
				User:           gittest.TestUser,
			}

			// We need to convert to an incoming context first in
			// order to preserve the feature flag.
			ctx = helper.OutgoingToIncoming(ctx)
			ctx, err = txinfo.InjectTransaction(ctx, 1, "node", testCase.primary)
			require.NoError(t, err)
			ctx = helper.IncomingToOutgoing(ctx)

			response, err := client.UserCreateTag(ctx, request)
			require.NoError(t, err)

			targetOID := text.ChompBytes(gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", "refs/tags/"+tagName))
			peeledOID := text.ChompBytes(gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", targetOID+"^{commit}"))
			targetOIDOK := targetOID
			if len(testCase.message) > 0 {
				targetOIDOK = peeledOID
			}
			require.Equal(t, targetOIDOK, targetRevision)

			testassert.ProtoEqual(t, &gitalypb.UserCreateTagResponse{
				Tag: &gitalypb.Tag{
					Name:         []byte(tagName),
					Message:      []byte(testCase.message),
					MessageSize:  int64(len(testCase.message)),
					Id:           targetOID,
					TargetCommit: targetCommit,
				},
			}, response)

			// Only the primary node should've executed hooks.
			if testCase.primary {
				contents := testhelper.MustReadFile(t, hooksOutputPath)
				require.Equal(t, "pre-receive\nupdate\npost-receive\n", string(contents))
			} else {
				require.NoFileExists(t, hooksOutputPath)
			}

			require.Equal(t, 2, transactionServer.called)
			transactionServer.called = 0
		})
	}
}

func TestUserCreateTagQuarantine(t *testing.T) {
	t.Parallel()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)
	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	// We set up a custom "update" hook which simply prints the new tag to stdout and then exits
	// with an error. Like this, we can both assert that the hook can see the quarantined tag,
	// and it allows us to fail the RPC before we migrate quarantined objects. Furthermore, we
	// also try whether we can print the tag's tagged object to assert that we can see objects
	// which are not part of the object quarantine.
	script := fmt.Sprintf(`#!/bin/sh
	%s cat-file -p $3^{commit} >/dev/null &&
	%s cat-file -p $3^{tag} &&
	exit 1`, cfg.Git.BinPath, cfg.Git.BinPath)
	gittest.WriteCustomHook(t, repoPath, "update", []byte(script))

	response, err := client.UserCreateTag(ctx, &gitalypb.UserCreateTagRequest{
		Repository:     repoProto,
		TagName:        []byte("quarantined-tag"),
		TargetRevision: []byte("c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd"),
		User:           gittest.TestUser,
		Timestamp:      timestamppb.New(time.Unix(1600000000, 0)),
		Message:        []byte("message"),
	})
	require.NoError(t, err)

	// Conveniently, the pre-receive error will now contain output from our custom hook and thus
	// the tag's contents.
	testassert.ProtoEqual(t, &gitalypb.UserCreateTagResponse{
		PreReceiveError: `object c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd
type commit
tag quarantined-tag
tagger Jane Doe <janedoe@gitlab.com> 1600000000 +0800

message`,
	}, response)

	// In case we use an object quarantine directory, the tag should not exist in the target
	// repository because the RPC failed to update the revision.
	tagExists, err := repo.HasRevision(ctx, "85d279b2cc85df37992e08f84707987321e8ef47^{tag}")
	require.NoError(t, err)
	require.False(t, tagExists, "tag should not have been migrated")
}

func TestSuccessfulUserCreateTagRequestAnnotatedLightweightDisambiguation(t *testing.T) {
	t.Parallel()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	preReceiveHook := writeAssertObjectTypePreReceiveHook(t, cfg)

	updateHook := writeAssertObjectTypeUpdateHook(t, cfg)

	testCases := []struct {
		desc    string
		message string
		objType string
		err     error
	}{
		{
			desc:    "error: contains null byte",
			message: "\000",
			err:     status.Error(codes.Unknown, "ArgumentError: string contains null byte"),
		},
		{
			desc:    "annotated: some control characters",
			message: "\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008",
			objType: "tag",
		},
		{
			desc:    "lightweight: empty message",
			message: "",
			objType: "commit",
		},
		{
			desc:    "lightweight: simple whitespace",
			message: " \t\t",
			objType: "commit",
		},
		{
			desc:    "lightweight: whitespace with newlines",
			message: "\t\n\f\r ",
			objType: "commit",
		},
		{
			desc:    "lightweight: simple Unicode whitespace",
			message: "\u00a0",
			objType: "tag",
		},
		{
			desc:    "lightweight: lots of Unicode whitespace",
			message: "\u0020\u00a0\u1680\u180e\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u200b\u202f\u205f\u3000\ufeff",
			objType: "tag",
		},
		{
			desc:    "annotated: dot",
			message: ".",
			objType: "tag",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			for hook, content := range map[string]string{
				"pre-receive": fmt.Sprintf("#!/bin/sh\n%s %s \"$@\"", preReceiveHook, testCase.objType),
				"update":      fmt.Sprintf("#!/bin/sh\n%s %s \"$@\"", updateHook, testCase.objType),
			} {
				gittest.WriteCustomHook(t, repoPath, hook, []byte(content))
			}

			tagName := "what-will-it-be"
			request := &gitalypb.UserCreateTagRequest{
				Repository:     repo,
				TagName:        []byte(tagName),
				TargetRevision: []byte("c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd"),
				User:           gittest.TestUser,
				Message:        []byte(testCase.message),
			}

			response, err := client.UserCreateTag(ctx, request)

			if testCase.err != nil {
				testassert.GrpcEqualErr(t, testCase.err, err)
			} else {
				defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", tagName)
				require.NoError(t, err)
				require.Empty(t, response.PreReceiveError)
			}
		})
	}
}

func TestSuccessfulUserCreateTagRequestWithParsedTargetRevision(t *testing.T) {
	t.Parallel()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	gittest.Exec(t, cfg, "-C", repoPath, "branch", "heads/master", "master~1")
	defer gittest.Exec(t, cfg, "-C", repoPath, "branch", "-d", "heads/master")
	gittest.Exec(t, cfg, "-C", repoPath, "branch", "refs/heads/master", "master~2")
	defer gittest.Exec(t, cfg, "-C", repoPath, "branch", "-d", "refs/heads/master")

	testCases := []struct {
		desc             string
		targetRevision   string
		expectedRevision string
	}{
		{
			desc:             "tag",
			targetRevision:   "v1.0.0",
			expectedRevision: "refs/tags/v1.0.0",
		},
		{
			desc:             "tag~",
			targetRevision:   "v1.0.0~",
			expectedRevision: "refs/tags/v1.0.0~",
		},
		{
			desc:             "tags/tag~",
			targetRevision:   "tags/v1.0.0~",
			expectedRevision: "refs/tags/v1.0.0~",
		},
		{
			desc:             "refs/tag~",
			targetRevision:   "refs/tags/v1.0.0~",
			expectedRevision: "refs/tags/v1.0.0~",
		},
		{
			desc:             "master",
			targetRevision:   "master",
			expectedRevision: "master",
		},
		{
			desc:             "heads/master",
			targetRevision:   "heads/master",
			expectedRevision: "refs/heads/heads/master",
		},
		{
			desc:             "refs/heads/master",
			targetRevision:   "refs/heads/master",
			expectedRevision: "refs/heads/master",
		},
		{
			desc:             "heads/refs/heads/master",
			targetRevision:   "heads/refs/heads/master",
			expectedRevision: "refs/heads/refs/heads/master",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			tagName := "what-will-it-be"
			request := &gitalypb.UserCreateTagRequest{
				Repository:     repo,
				TagName:        []byte(tagName),
				TargetRevision: []byte(testCase.targetRevision),
				User:           gittest.TestUser,
			}

			response, err := client.UserCreateTag(ctx, request)
			defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", tagName)
			require.NoError(t, err)
			require.Empty(t, response.PreReceiveError)

			parsedID := gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", tagName)
			require.Equal(t, text.ChompBytes(parsedID), response.Tag.TargetCommit.Id)
		})
	}
}

func TestSuccessfulUserCreateTagRequestToNonCommit(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	inputTagName := "to-be-créated-soon"

	preReceiveHook := writeAssertObjectTypePreReceiveHook(t, cfg)

	updateHook := writeAssertObjectTypeUpdateHook(t, cfg)

	testCases := []struct {
		desc               string
		tagName            string
		message            string
		targetRevision     string
		expectedTag        *gitalypb.Tag
		expectedObjectType string
	}{
		{
			desc:           "lightweight tag to tree",
			tagName:        inputTagName,
			targetRevision: "612036fac47c5d31c212b17268e2f3ba807bce1e",
			expectedTag: &gitalypb.Tag{
				Name: []byte(inputTagName),
				Id:   "612036fac47c5d31c212b17268e2f3ba807bce1e",
			},
			expectedObjectType: "tree",
		},
		{
			desc:           "lightweight tag to blob",
			tagName:        inputTagName,
			targetRevision: "dfaa3f97ca337e20154a98ac9d0be76ddd1fcc82",
			expectedTag: &gitalypb.Tag{
				Name: []byte(inputTagName),
				Id:   "dfaa3f97ca337e20154a98ac9d0be76ddd1fcc82",
			},
			expectedObjectType: "blob",
		},
		{
			desc:           "annotated tag to tree",
			tagName:        inputTagName,
			targetRevision: "612036fac47c5d31c212b17268e2f3ba807bce1e",
			message:        "This is an annotated tag",
			expectedTag: &gitalypb.Tag{
				Name: []byte(inputTagName),
				//Id: is a new object, filled in below
				TargetCommit: nil,
				Message:      []byte("This is an annotated tag"),
				MessageSize:  24,
			},
			expectedObjectType: "tag",
		},
		{
			desc:           "annotated tag to blob",
			tagName:        inputTagName,
			targetRevision: "dfaa3f97ca337e20154a98ac9d0be76ddd1fcc82",
			message:        "This is an annotated tag",
			expectedTag: &gitalypb.Tag{
				Name: []byte(inputTagName),
				//Id: is a new object, filled in below
				TargetCommit: nil,
				Message:      []byte("This is an annotated tag"),
				MessageSize:  24,
			},
			expectedObjectType: "tag",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			for hook, content := range map[string]string{
				"pre-receive": fmt.Sprintf("#!/bin/sh\n%s %s \"$@\"", preReceiveHook, testCase.expectedObjectType),
				"update":      fmt.Sprintf("#!/bin/sh\n%s %s \"$@\"", updateHook, testCase.expectedObjectType),
			} {
				gittest.WriteCustomHook(t, repoPath, hook, []byte(content))
			}

			request := &gitalypb.UserCreateTagRequest{
				Repository:     repo,
				TagName:        []byte(testCase.tagName),
				TargetRevision: []byte(testCase.targetRevision),
				User:           gittest.TestUser,
				Message:        []byte(testCase.message),
			}

			responseOk := &gitalypb.UserCreateTagResponse{
				Tag: testCase.expectedTag,
			}
			response, err := client.UserCreateTag(ctx, request)
			require.NoError(t, err)
			require.Empty(t, response.PreReceiveError)
			defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", inputTagName)

			// Fake up *.Id for annotated tags
			if len(testCase.expectedTag.Id) == 0 {
				tagID := gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", inputTagName)
				responseOk.Tag.Id = text.ChompBytes(tagID)
			}
			testassert.ProtoEqual(t, responseOk, response)

			peeledID := gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", inputTagName+"^{}")
			require.Equal(t, testCase.targetRevision, text.ChompBytes(peeledID))

			objectType := gittest.Exec(t, cfg, "-C", repoPath, "cat-file", "-t", inputTagName)
			require.Equal(t, testCase.expectedObjectType, text.ChompBytes(objectType))
		})
	}
}

func TestSuccessfulUserCreateTagNestedTags(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	preReceiveHook := writeAssertObjectTypePreReceiveHook(t, cfg)

	updateHook := writeAssertObjectTypeUpdateHook(t, cfg)

	testCases := []struct {
		desc             string
		targetObject     string
		targetObjectType string
		expectedTag      *gitalypb.Tag
	}{
		{
			desc:             "nested tags to commit",
			targetObject:     "c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd",
			targetObjectType: "commit",
		},
		{
			desc:             "nested tags to tree",
			targetObjectType: "tree",
			targetObject:     "612036fac47c5d31c212b17268e2f3ba807bce1e",
		},
		{
			desc:             "nested tags to blob",
			targetObject:     "dfaa3f97ca337e20154a98ac9d0be76ddd1fcc82",
			targetObjectType: "blob",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			// We resolve down to commit/tree/blob, but
			// we'll only ever push a "tag" here.
			hookObjectType := "tag"
			for hook, content := range map[string]string{
				"pre-receive": fmt.Sprintf("#!/bin/sh\n%s %s \"$@\"", preReceiveHook, hookObjectType),
				"update":      fmt.Sprintf("#!/bin/sh\n%s %s \"$@\"", updateHook, hookObjectType),
			} {
				gittest.WriteCustomHook(t, repoPath, hook, []byte(content))
			}

			targetObject := testCase.targetObject
			nestLevel := 2
			for i := 0; i <= nestLevel; i++ {
				tagName := fmt.Sprintf("nested-tag-%v", i)
				tagMessage := fmt.Sprintf("This is level %v of a nested annotated tag to %v", i, testCase.targetObject)
				request := &gitalypb.UserCreateTagRequest{
					Repository:     repoProto,
					TagName:        []byte(tagName),
					TargetRevision: []byte(targetObject),
					User:           gittest.TestUser,
					Message:        []byte(tagMessage),
				}
				response, err := client.UserCreateTag(ctx, request)
				require.NoError(t, err)
				require.Empty(t, response.PreReceiveError)
				defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", tagName)

				createdID := gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", tagName)
				createdIDStr := text.ChompBytes(createdID)
				responseOk := &gitalypb.UserCreateTagResponse{
					Tag: &gitalypb.Tag{
						Name: request.TagName,
						Id:   createdIDStr,
						//TargetCommit: is dymamically determined, filled in below
						Message:     request.Message,
						MessageSize: int64(len(request.Message)),
					},
				}
				// Fake it up for all levels, except for ^{} == "commit"
				responseOk.Tag.TargetCommit = response.Tag.TargetCommit
				if testCase.targetObjectType == "commit" {
					responseOk.Tag.TargetCommit, err = repo.ReadCommit(ctx, git.Revision(testCase.targetObject))
					require.NoError(t, err)
				}
				testassert.ProtoEqual(t, responseOk, response)

				peeledID := gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", tagName+"^{}")
				peeledIDStr := text.ChompBytes(peeledID)
				require.Equal(t, testCase.targetObject, peeledIDStr)

				// Set up the next level of nesting...
				targetObject = response.Tag.Id

				// Create a *lightweight* tag pointing
				// to our N-level
				// tag->[commit|tree|blob]. The "tag"
				// field name will not match the tag
				// name.
				tagNameLight := fmt.Sprintf("skip-type-check-light-%s", tagName)
				request = &gitalypb.UserCreateTagRequest{
					Repository:     repoProto,
					TagName:        []byte(tagNameLight),
					TargetRevision: []byte(createdIDStr),
					User:           gittest.TestUser,
				}
				response, err = client.UserCreateTag(ctx, request)
				defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", tagNameLight)
				require.NoError(t, err)
				require.Empty(t, response.PreReceiveError)

				responseOk = &gitalypb.UserCreateTagResponse{
					Tag: &gitalypb.Tag{
						Name:         request.TagName,
						Id:           testCase.targetObject,
						TargetCommit: responseOk.Tag.TargetCommit,
						Message:      nil,
						MessageSize:  0,
					},
				}
				testassert.ProtoEqual(t, responseOk, response)

				createdIDLight := gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", tagNameLight)
				createdIDLightStr := text.ChompBytes(createdIDLight)
				require.Equal(t, testCase.targetObject, createdIDLightStr)
			}
		})
	}
}

func TestUserCreateTagStableTagIDs(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, _, repo, _, client := setupOperationsService(t, ctx)

	response, err := client.UserCreateTag(ctx, &gitalypb.UserCreateTagRequest{
		Repository:     repo,
		TagName:        []byte("happy-tag"),
		TargetRevision: []byte("dfaa3f97ca337e20154a98ac9d0be76ddd1fcc82"),
		Message:        []byte("my message"),
		User:           gittest.TestUser,
		Timestamp:      &timestamppb.Timestamp{Seconds: 12345},
	})
	require.NoError(t, err)

	require.Equal(t, &gitalypb.Tag{
		Id:          "123b02f05cc249a7da87aae583babb8e4871cd65",
		Name:        []byte("happy-tag"),
		Message:     []byte("my message"),
		MessageSize: 10,
	}, response.Tag)
}

func TestUserDeleteTagSuccessfulDeletionOfPrefixedTag(t *testing.T) {
	t.Parallel()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	testCases := []struct {
		desc         string
		tagNameInput string
		tagCommit    string
		user         *gitalypb.User
		response     *gitalypb.UserDeleteTagResponse
		err          error
	}{
		{
			desc:         "possible to delete a tag called refs/tags/something",
			tagNameInput: "refs/tags/can-find-this",
			tagCommit:    "c642fe9b8b9f28f9225d7ea953fe14e74748d53b",
			user:         gittest.TestUser,
			response:     &gitalypb.UserDeleteTagResponse{},
			err:          nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			gittest.Exec(t, cfg, "-C", repoPath, "tag", testCase.tagNameInput, testCase.tagCommit)

			request := &gitalypb.UserDeleteTagRequest{
				Repository: repo,
				TagName:    []byte(testCase.tagNameInput),
				User:       testCase.user,
			}

			response, err := client.UserDeleteTag(ctx, request)
			testassert.GrpcEqualErr(t, testCase.err, err)
			testassert.ProtoEqual(t, testCase.response, response)

			refs := gittest.Exec(t, cfg, "-C", repoPath, "for-each-ref", "--", "refs/tags/"+testCase.tagNameInput)
			require.NotContains(t, string(refs), testCase.tagCommit, "tag kept because we stripped off refs/tags/*")
		})
	}
}

func TestUserCreateTagsuccessfulCreationOfPrefixedTag(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	testCases := []struct {
		desc                   string
		tagNameInput           string
		tagTargetRevisionInput string
		user                   *gitalypb.User
		err                    error
	}{
		{
			desc:                   "possible to create a tag called refs/tags/something",
			tagNameInput:           "refs/tags/can-create-this",
			tagTargetRevisionInput: "1a0b36b3cdad1d2ee32457c102a8c0b7056fa863",
			user:                   gittest.TestUser,
			err:                    nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", testCase.tagNameInput)

			request := &gitalypb.UserCreateTagRequest{
				Repository:     repoProto,
				TagName:        []byte(testCase.tagNameInput),
				TargetRevision: []byte(testCase.tagTargetRevisionInput),
				User:           testCase.user,
			}

			response, err := client.UserCreateTag(ctx, request)
			testassert.GrpcEqualErr(t, testCase.err, err)
			commitOk, err := repo.ReadCommit(ctx, git.Revision(testCase.tagTargetRevisionInput))
			require.NoError(t, err)

			responseOk := &gitalypb.UserCreateTagResponse{
				Tag: &gitalypb.Tag{
					Name:         []byte(testCase.tagNameInput),
					Id:           testCase.tagTargetRevisionInput,
					TargetCommit: commitOk,
				},
			}

			testassert.ProtoEqual(t, responseOk, response)

			refs := gittest.Exec(t, cfg, "-C", repoPath, "for-each-ref", "--", "refs/tags/"+testCase.tagNameInput)
			require.Contains(t, string(refs), testCase.tagTargetRevisionInput, "tag created, we did not strip off refs/tags/*")
		})
	}
}

func TestSuccessfulGitHooksForUserCreateTagRequest(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	projectPath := "project/path"
	repo.GlProjectPath = projectPath

	tagName := "new-tag"

	request := &gitalypb.UserCreateTagRequest{
		Repository:     repo,
		TagName:        []byte(tagName),
		TargetRevision: []byte("c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd"),
		User:           gittest.TestUser,
	}

	for _, hookName := range GitlabHooks {
		t.Run(hookName, func(t *testing.T) {
			defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", tagName)

			hookOutputTempPath := gittest.WriteEnvToCustomHook(t, repoPath, hookName)

			response, err := client.UserCreateTag(ctx, request)
			require.NoError(t, err)
			require.Empty(t, response.PreReceiveError)

			output := string(testhelper.MustReadFile(t, hookOutputTempPath))
			require.Contains(t, output, "GL_USERNAME="+gittest.TestUser.GlUsername)
			require.Contains(t, output, "GL_PROJECT_PATH="+projectPath)
		})
	}
}

func TestFailedUserDeleteTagRequestDueToValidation(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, _, repo, _, client := setupOperationsService(t, ctx)

	testCases := []struct {
		desc     string
		request  *gitalypb.UserDeleteTagRequest
		response *gitalypb.UserDeleteTagResponse
		err      error
	}{
		{
			desc: "empty user",
			request: &gitalypb.UserDeleteTagRequest{
				Repository: repo,
				TagName:    []byte("does-matter-the-name-if-user-is-empty"),
			},
			response: nil,
			err:      status.Error(codes.InvalidArgument, "empty user"),
		},
		{
			desc: "empty tag name",
			request: &gitalypb.UserDeleteTagRequest{
				Repository: repo,
				User:       gittest.TestUser,
			},
			response: nil,
			err:      status.Error(codes.InvalidArgument, "empty tag name"),
		},
		{
			desc: "non-existent tag name",
			request: &gitalypb.UserDeleteTagRequest{
				Repository: repo,
				User:       gittest.TestUser,
				TagName:    []byte("i-do-not-exist"),
			},
			response: nil,
			err:      status.Errorf(codes.FailedPrecondition, "tag not found: %s", "i-do-not-exist"),
		},
		{
			desc: "space in tag name",
			request: &gitalypb.UserDeleteTagRequest{
				Repository: repo,
				User:       gittest.TestUser,
				TagName:    []byte("a tag"),
			},
			response: nil,
			err:      status.Errorf(codes.FailedPrecondition, "tag not found: %s", "a tag"),
		},
		{
			desc: "newline in tag name",
			request: &gitalypb.UserDeleteTagRequest{
				Repository: repo,
				User:       gittest.TestUser,
				TagName:    []byte("a\ntag"),
			},
			response: nil,
			err:      status.Errorf(codes.FailedPrecondition, "tag not found: %s", "a\ntag"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			response, err := client.UserDeleteTag(ctx, testCase.request)
			testassert.GrpcEqualErr(t, testCase.err, err)
			testassert.ProtoEqual(t, testCase.response, response)
		})
	}
}

func TestFailedUserDeleteTagDueToHooks(t *testing.T) {
	t.Parallel()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	tagNameInput := "to-be-deleted-soon-tag"
	gittest.Exec(t, cfg, "-C", repoPath, "tag", tagNameInput)
	defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", tagNameInput)

	request := &gitalypb.UserDeleteTagRequest{
		Repository: repo,
		TagName:    []byte(tagNameInput),
		User:       gittest.TestUser,
	}

	hookContent := []byte("#!/bin/sh\necho GL_ID=$GL_ID\nexit 1")

	for _, hookName := range gitlabPreHooks {
		t.Run(hookName, func(t *testing.T) {
			gittest.WriteCustomHook(t, repoPath, hookName, hookContent)

			response, err := client.UserDeleteTag(ctx, request)
			require.NoError(t, err)
			require.Contains(t, response.PreReceiveError, "GL_ID="+gittest.TestUser.GlId)

			tags := gittest.Exec(t, cfg, "-C", repoPath, "tag")
			require.Contains(t, string(tags), tagNameInput, "tag name does not exist in tags list")
		})
	}
}

func TestFailedUserCreateTagDueToHooks(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, _, repo, repoPath, client := setupOperationsService(t, ctx)

	request := &gitalypb.UserCreateTagRequest{
		Repository:     repo,
		TagName:        []byte("new-tag"),
		TargetRevision: []byte("c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd"),
		User:           gittest.TestUser,
	}

	hookContent := []byte("#!/bin/sh\necho GL_ID=$GL_ID\nexit 1")

	for _, hookName := range gitlabPreHooks {
		gittest.WriteCustomHook(t, repoPath, hookName, hookContent)

		response, err := client.UserCreateTag(ctx, request)
		require.NoError(t, err)
		require.Contains(t, response.PreReceiveError, "GL_ID="+gittest.TestUser.GlId)
	}
}

func TestFailedUserCreateTagRequestDueToTagExistence(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, _, repo, _, client := setupOperationsService(t, ctx)

	testCases := []struct {
		desc           string
		tagName        string
		targetRevision string
		user           *gitalypb.User
		response       *gitalypb.UserCreateTagResponse
		err            error
	}{
		{
			desc:           "simple existing tag",
			tagName:        "v1.1.0",
			targetRevision: "master",
			user:           gittest.TestUser,
			response: &gitalypb.UserCreateTagResponse{
				Tag:    nil,
				Exists: true,
			},
			err: nil,
		},
		{
			desc:           "existing tag nonexisting target revision",
			tagName:        "v1.1.0",
			targetRevision: "does-not-exist",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Errorf(codes.FailedPrecondition, "revspec '%s' not found", "does-not-exist"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			request := &gitalypb.UserCreateTagRequest{
				Repository:     repo,
				TagName:        []byte(testCase.tagName),
				TargetRevision: []byte(testCase.targetRevision),
				User:           testCase.user,
			}

			response, err := client.UserCreateTag(ctx, request)
			testassert.GrpcEqualErr(t, testCase.err, err)
			testassert.ProtoEqual(t, testCase.response, response)
		})
	}
}

func TestFailedUserCreateTagRequestDueToValidation(t *testing.T) {
	t.Parallel()
	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, _, repo, _, client := setupOperationsService(t, ctx)

	injectedTag := "inject-tag\ntagger . <> 0 +0000\n\nInjected subject\n\n"
	testCases := []struct {
		desc           string
		tagName        string
		targetRevision string
		message        string
		user           *gitalypb.User
		response       *gitalypb.UserCreateTagResponse
		err            error
	}{
		{
			desc:           "empty target revision",
			tagName:        "shiny-new-tag",
			targetRevision: "",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Error(codes.InvalidArgument, "empty target revision"),
		},
		{
			desc:           "empty user",
			tagName:        "shiny-new-tag",
			targetRevision: "master",
			user:           nil,
			response:       nil,
			err:            status.Error(codes.InvalidArgument, "empty user"),
		},
		{
			desc:           "empty starting point",
			tagName:        "new-tag",
			targetRevision: "",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Error(codes.InvalidArgument, "empty target revision"),
		},
		{
			desc:           "non-existing starting point",
			tagName:        "new-tag",
			targetRevision: "i-dont-exist",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Errorf(codes.FailedPrecondition, "revspec '%s' not found", "i-dont-exist"),
		},
		{
			desc:           "space in lightweight tag name",
			tagName:        "a tag",
			targetRevision: "master",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Errorf(codes.Unknown, "Gitlab::Git::CommitError: Could not update refs/tags/%s. Please refresh and try again.", "a tag"),
		},
		{
			desc:           "space in annotated tag name",
			tagName:        "a tag",
			targetRevision: "master",
			message:        "a message",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Errorf(codes.Unknown, "Gitlab::Git::CommitError: Could not update refs/tags/%s. Please refresh and try again.", "a tag"),
		},
		{
			desc:           "newline in lightweight tag name",
			tagName:        "a\ntag",
			targetRevision: "master",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Errorf(codes.Unknown, "Gitlab::Git::CommitError: Could not update refs/tags/%s. Please refresh and try again.", "a\ntag"),
		},
		{
			desc:           "newline in annotated tag name",
			tagName:        "a\ntag",
			targetRevision: "master",
			message:        "a message",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Error(codes.Unknown, "Rugged::InvalidError: failed to parse signature - expected prefix doesn't match actual"),
		},
		{
			desc:           "injection in lightweight tag name",
			tagName:        injectedTag,
			targetRevision: "master",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Errorf(codes.Unknown, "Gitlab::Git::CommitError: Could not update refs/tags/%s. Please refresh and try again.", injectedTag),
		},
		{
			desc:           "injection in annotated tag name",
			tagName:        injectedTag,
			targetRevision: "master",
			message:        "a message",
			user:           gittest.TestUser,
			response:       nil,
			err:            status.Errorf(codes.Unknown, "Gitlab::Git::CommitError: Could not update refs/tags/%s. Please refresh and try again.", injectedTag),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			request := &gitalypb.UserCreateTagRequest{
				Repository:     repo,
				TagName:        []byte(testCase.tagName),
				TargetRevision: []byte(testCase.targetRevision),
				User:           testCase.user,
				Message:        []byte(testCase.message),
			}

			response, err := client.UserCreateTag(ctx, request)
			testassert.GrpcEqualErr(t, testCase.err, err)
			testassert.ProtoEqual(t, testCase.response, response)
		})
	}
}

func TestTagHookOutput(t *testing.T) {
	t.Parallel()

	ctx, cancel := testhelper.Context()
	defer cancel()

	ctx, cfg, repo, repoPath, client := setupOperationsService(t, ctx)

	testCases := []struct {
		desc        string
		hookContent string
		output      string
	}{
		{
			desc:        "empty stdout and empty stderr",
			hookContent: "#!/bin/sh\nexit 1",
			output:      "",
		},
		{
			desc:        "empty stdout and some stderr",
			hookContent: "#!/bin/sh\necho stderr >&2\nexit 1",
			output:      "stderr\n",
		},
		{
			desc:        "some stdout and empty stderr",
			hookContent: "#!/bin/sh\necho stdout\nexit 1",
			output:      "stdout\n",
		},
		{
			desc:        "some stdout and some stderr",
			hookContent: "#!/bin/sh\necho stdout\necho stderr >&2\nexit 1",
			output:      "stderr\n",
		},
		{
			desc:        "whitespace stdout and some stderr",
			hookContent: "#!/bin/sh\necho '   '\necho stderr >&2\nexit 1",
			output:      "stderr\n",
		},
		{
			desc:        "some stdout and whitespace stderr",
			hookContent: "#!/bin/sh\necho stdout\necho '   ' >&2\nexit 1",
			output:      "stdout\n",
		},
	}

	for _, hookName := range gitlabPreHooks {
		for _, testCase := range testCases {
			t.Run(hookName+"/"+testCase.desc, func(t *testing.T) {
				tagNameInput := "some-tag"
				createRequest := &gitalypb.UserCreateTagRequest{
					Repository:     repo,
					TagName:        []byte(tagNameInput),
					TargetRevision: []byte("master"),
					User:           gittest.TestUser,
				}
				deleteRequest := &gitalypb.UserDeleteTagRequest{
					Repository: repo,
					TagName:    []byte(tagNameInput),
					User:       gittest.TestUser,
				}

				gittest.WriteCustomHook(t, repoPath, hookName, []byte(testCase.hookContent))

				createResponse, err := client.UserCreateTag(ctx, createRequest)
				require.NoError(t, err)

				createResponseOk := &gitalypb.UserCreateTagResponse{
					Tag:             createResponse.Tag,
					Exists:          false,
					PreReceiveError: testCase.output,
				}
				testassert.ProtoEqual(t, createResponseOk, createResponse)

				defer gittest.Exec(t, cfg, "-C", repoPath, "tag", "-d", tagNameInput)
				gittest.Exec(t, cfg, "-C", repoPath, "tag", tagNameInput)

				deleteResponse, err := client.UserDeleteTag(ctx, deleteRequest)
				require.NoError(t, err)
				deleteResponseOk := &gitalypb.UserDeleteTagResponse{
					PreReceiveError: testCase.output,
				}
				testassert.ProtoEqual(t, deleteResponseOk, deleteResponse)
			})
		}
	}
}
