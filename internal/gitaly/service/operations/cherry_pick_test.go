package operations

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/metadata/featureflag"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testassert"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServer_UserCherryPick_successful(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickSuccessful)
}

func testServerUserCherryPickSuccessful(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	destinationBranch := "cherry-picking-dst"
	gittest.Exec(t, cfg, "-C", repoPath, "branch", destinationBranch, "master")

	masterHeadCommit, err := repo.ReadCommit(ctx, "master")
	require.NoError(t, err)

	cherryPickedCommit, err := repo.ReadCommit(ctx, "8a0f2ee90d940bfb0ba1e14e8214b0649056e4ab")
	require.NoError(t, err)

	testRepoCopy, testRepoCopyPath := gittest.CloneRepo(t, cfg, cfg.Storages[0]) // read-only repo

	gittest.Exec(t, cfg, "-C", testRepoCopyPath, "branch", destinationBranch, "master")

	testCases := []struct {
		desc         string
		request      *gitalypb.UserCherryPickRequest
		branchUpdate *gitalypb.OperationBranchUpdate
	}{
		{
			desc: "branch exists",
			request: &gitalypb.UserCherryPickRequest{
				Repository: repoProto,
				User:       gittest.TestUser,
				Commit:     cherryPickedCommit,
				BranchName: []byte(destinationBranch),
				Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
			},
			branchUpdate: &gitalypb.OperationBranchUpdate{},
		},
		{
			desc: "nonexistent branch + start_repository == repository",
			request: &gitalypb.UserCherryPickRequest{
				Repository:      repoProto,
				User:            gittest.TestUser,
				Commit:          cherryPickedCommit,
				BranchName:      []byte("to-be-cherry-picked-into-1"),
				Message:         []byte("Cherry-picking " + cherryPickedCommit.Id),
				StartBranchName: []byte("master"),
			},
			branchUpdate: &gitalypb.OperationBranchUpdate{BranchCreated: true},
		},
		{
			desc: "nonexistent branch + start_repository != repository",
			request: &gitalypb.UserCherryPickRequest{
				Repository:      repoProto,
				User:            gittest.TestUser,
				Commit:          cherryPickedCommit,
				BranchName:      []byte("to-be-cherry-picked-into-2"),
				Message:         []byte("Cherry-picking " + cherryPickedCommit.Id),
				StartRepository: testRepoCopy,
				StartBranchName: []byte("master"),
			},
			branchUpdate: &gitalypb.OperationBranchUpdate{BranchCreated: true},
		},
		{
			desc: "nonexistent branch + empty start_repository",
			request: &gitalypb.UserCherryPickRequest{
				Repository:      repoProto,
				User:            gittest.TestUser,
				Commit:          cherryPickedCommit,
				BranchName:      []byte("to-be-cherry-picked-into-3"),
				Message:         []byte("Cherry-picking " + cherryPickedCommit.Id),
				StartBranchName: []byte("master"),
			},
			branchUpdate: &gitalypb.OperationBranchUpdate{BranchCreated: true},
		},
		{
			desc: "branch exists with dry run",
			request: &gitalypb.UserCherryPickRequest{
				Repository: testRepoCopy,
				User:       gittest.TestUser,
				Commit:     cherryPickedCommit,
				BranchName: []byte(destinationBranch),
				Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
				DryRun:     true,
			},
			branchUpdate: &gitalypb.OperationBranchUpdate{},
		},
		{
			desc: "nonexistent branch + start_repository == repository with dry run",
			request: &gitalypb.UserCherryPickRequest{
				Repository:      testRepoCopy,
				User:            gittest.TestUser,
				Commit:          cherryPickedCommit,
				BranchName:      []byte("to-be-cherry-picked-into-1"),
				Message:         []byte("Cherry-picking " + cherryPickedCommit.Id),
				StartBranchName: []byte("master"),
				DryRun:          true,
			},
			branchUpdate: &gitalypb.OperationBranchUpdate{BranchCreated: true},
		},
		{
			desc: "nonexistent branch + start_repository != repository with dry run",
			request: &gitalypb.UserCherryPickRequest{
				Repository:      testRepoCopy,
				User:            gittest.TestUser,
				Commit:          cherryPickedCommit,
				BranchName:      []byte("to-be-cherry-picked-into-2"),
				Message:         []byte("Cherry-picking " + cherryPickedCommit.Id),
				StartRepository: testRepoCopy,
				StartBranchName: []byte("master"),
				DryRun:          true,
			},
			branchUpdate: &gitalypb.OperationBranchUpdate{BranchCreated: true},
		},
		{
			desc: "nonexistent branch + empty start_repository with dry run",
			request: &gitalypb.UserCherryPickRequest{
				Repository:      testRepoCopy,
				User:            gittest.TestUser,
				Commit:          cherryPickedCommit,
				BranchName:      []byte("to-be-cherry-picked-into-3"),
				Message:         []byte("Cherry-picking " + cherryPickedCommit.Id),
				StartBranchName: []byte("master"),
				DryRun:          true,
			},
			branchUpdate: &gitalypb.OperationBranchUpdate{BranchCreated: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			response, err := client.UserCherryPick(ctx, testCase.request)
			require.NoError(t, err)

			testRepo := localrepo.NewTestRepo(t, cfg, testCase.request.Repository)
			headCommit, err := testRepo.ReadCommit(ctx, git.Revision(testCase.request.BranchName))
			require.NoError(t, err)

			expectedBranchUpdate := testCase.branchUpdate
			expectedBranchUpdate.CommitId = headCommit.Id

			require.Equal(t, expectedBranchUpdate, response.BranchUpdate)
			require.Empty(t, response.CreateTreeError)
			require.Empty(t, response.CreateTreeErrorCode)

			if testCase.request.DryRun {
				testassert.ProtoEqual(t, masterHeadCommit, headCommit)
			} else {
				require.Equal(t, testCase.request.Message, headCommit.Subject)
				require.Equal(t, masterHeadCommit.Id, headCommit.ParentIds[0])
			}
		})
	}
}

func TestServer_UserCherryPick_successfulGitHooks(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickSuccessfulGitHooks)
}

func testServerUserCherryPickSuccessfulGitHooks(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	destinationBranch := "cherry-picking-dst"
	gittest.Exec(t, cfg, "-C", repoPath, "branch", destinationBranch, "master")

	cherryPickedCommit, err := repo.ReadCommit(ctx, "8a0f2ee90d940bfb0ba1e14e8214b0649056e4ab")
	require.NoError(t, err)

	request := &gitalypb.UserCherryPickRequest{
		Repository: repoProto,
		User:       gittest.TestUser,
		Commit:     cherryPickedCommit,
		BranchName: []byte(destinationBranch),
		Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
	}

	var hookOutputFiles []string
	for _, hookName := range GitlabHooks {
		hookOutputTempPath := gittest.WriteEnvToCustomHook(t, repoPath, hookName)
		hookOutputFiles = append(hookOutputFiles, hookOutputTempPath)
	}

	response, err := client.UserCherryPick(ctx, request)
	require.NoError(t, err)
	require.Empty(t, response.PreReceiveError)

	for _, file := range hookOutputFiles {
		output := string(testhelper.MustReadFile(t, file))
		require.Contains(t, output, "GL_USERNAME="+gittest.TestUser.GlUsername)
	}
}

func TestServer_UserCherryPick_stableID(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickStableID)
}

func testServerUserCherryPickStableID(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	destinationBranch := "cherry-picking-dst"
	gittest.Exec(t, cfg, "-C", repoPath, "branch", destinationBranch, "master")

	commitToPick, err := repo.ReadCommit(ctx, "8a0f2ee90d940bfb0ba1e14e8214b0649056e4ab")
	require.NoError(t, err)

	request := &gitalypb.UserCherryPickRequest{
		Repository: repoProto,
		User:       gittest.TestUser,
		Commit:     commitToPick,
		BranchName: []byte(destinationBranch),
		Message:    []byte("Cherry-picking " + commitToPick.Id),
		Timestamp:  &timestamppb.Timestamp{Seconds: 12345},
	}

	response, err := client.UserCherryPick(ctx, request)
	require.NoError(t, err)
	require.Empty(t, response.PreReceiveError)
	require.Equal(t, "b17aeac93194cf2385b32623494ebce66efbacad", response.BranchUpdate.CommitId)

	pickedCommit, err := repo.ReadCommit(ctx, git.Revision(response.BranchUpdate.CommitId))
	require.NoError(t, err)
	require.Equal(t, &gitalypb.GitCommit{
		Id:        "b17aeac93194cf2385b32623494ebce66efbacad",
		Subject:   []byte("Cherry-picking " + commitToPick.Id),
		Body:      []byte("Cherry-picking " + commitToPick.Id),
		BodySize:  55,
		ParentIds: []string{"1e292f8fedd741b75372e19097c76d327140c312"},
		TreeId:    "5f1b6bcadf0abc482a19454aeaa219a5998db083",
		Author: &gitalypb.CommitAuthor{
			Name:  []byte("Ahmad Sherif"),
			Email: []byte("me@ahmadsherif.com"),
			Date: &timestamppb.Timestamp{
				Seconds: 1487337076,
			},
			Timezone: []byte("+0200"),
		},
		Committer: &gitalypb.CommitAuthor{
			Name:  gittest.TestUser.Name,
			Email: gittest.TestUser.Email,
			Date: &timestamppb.Timestamp{
				Seconds: 12345,
			},
			Timezone: []byte("+0000"),
		},
	}, pickedCommit)
}

func TestServer_UserCherryPick_failedValidations(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickFailedValidations)
}

func testServerUserCherryPickFailedValidations(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, _, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	cherryPickedCommit, err := repo.ReadCommit(ctx, "8a0f2ee90d940bfb0ba1e14e8214b0649056e4ab")
	require.NoError(t, err)

	destinationBranch := "cherry-picking-dst"

	testCases := []struct {
		desc    string
		request *gitalypb.UserCherryPickRequest
		code    codes.Code
	}{
		{
			desc: "empty user",
			request: &gitalypb.UserCherryPickRequest{
				Repository: repoProto,
				User:       nil,
				Commit:     cherryPickedCommit,
				BranchName: []byte(destinationBranch),
				Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty commit",
			request: &gitalypb.UserCherryPickRequest{
				Repository: repoProto,
				User:       gittest.TestUser,
				Commit:     nil,
				BranchName: []byte(destinationBranch),
				Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty branch name",
			request: &gitalypb.UserCherryPickRequest{
				Repository: repoProto,
				User:       gittest.TestUser,
				Commit:     cherryPickedCommit,
				BranchName: nil,
				Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty message",
			request: &gitalypb.UserCherryPickRequest{
				Repository: repoProto,
				User:       gittest.TestUser,
				Commit:     cherryPickedCommit,
				BranchName: []byte(destinationBranch),
				Message:    nil,
			},
			code: codes.InvalidArgument,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := client.UserCherryPick(ctx, testCase.request)
			testhelper.RequireGrpcError(t, err, testCase.code)
		})
	}
}

func TestServer_UserCherryPick_failedWithPreReceiveError(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickFailedWithPreReceiveError)
}

func testServerUserCherryPickFailedWithPreReceiveError(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	destinationBranch := "cherry-picking-dst"
	gittest.Exec(t, cfg, "-C", repoPath, "branch", destinationBranch, "master")

	cherryPickedCommit, err := repo.ReadCommit(ctx, "8a0f2ee90d940bfb0ba1e14e8214b0649056e4ab")
	require.NoError(t, err)

	request := &gitalypb.UserCherryPickRequest{
		Repository: repoProto,
		User:       gittest.TestUser,
		Commit:     cherryPickedCommit,
		BranchName: []byte(destinationBranch),
		Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
	}

	hookContent := []byte("#!/bin/sh\necho GL_ID=$GL_ID\nexit 1")

	for _, hookName := range GitlabPreHooks {
		t.Run(hookName, func(t *testing.T) {
			gittest.WriteCustomHook(t, repoPath, hookName, hookContent)

			response, err := client.UserCherryPick(ctx, request)
			require.NoError(t, err)
			require.Contains(t, response.PreReceiveError, "GL_ID="+gittest.TestUser.GlId)
		})
	}
}

func TestServer_UserCherryPick_failedWithCreateTreeError(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickFailedWithCreateTreeError)
}

func testServerUserCherryPickFailedWithCreateTreeError(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	destinationBranch := "cherry-picking-dst"
	gittest.Exec(t, cfg, "-C", repoPath, "branch", destinationBranch, "master")

	// This commit already exists in master
	cherryPickedCommit, err := repo.ReadCommit(ctx, "4a24d82dbca5c11c61556f3b35ca472b7463187e")
	require.NoError(t, err)

	request := &gitalypb.UserCherryPickRequest{
		Repository: repoProto,
		User:       gittest.TestUser,
		Commit:     cherryPickedCommit,
		BranchName: []byte(destinationBranch),
		Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
	}

	response, err := client.UserCherryPick(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, response.CreateTreeError)
	require.Equal(t, gitalypb.UserCherryPickResponse_EMPTY, response.CreateTreeErrorCode)
}

func TestServer_UserCherryPick_failedWithCommitError(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickFailedWithCommitError)
}

func testServerUserCherryPickFailedWithCommitError(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	sourceBranch := "cherry-pick-src"
	destinationBranch := "cherry-picking-dst"
	gittest.Exec(t, cfg, "-C", repoPath, "branch", destinationBranch, "master")
	gittest.Exec(t, cfg, "-C", repoPath, "branch", sourceBranch, "8a0f2ee90d940bfb0ba1e14e8214b0649056e4ab")

	cherryPickedCommit, err := repo.ReadCommit(ctx, git.Revision(sourceBranch))
	require.NoError(t, err)

	request := &gitalypb.UserCherryPickRequest{
		Repository:      repoProto,
		User:            gittest.TestUser,
		Commit:          cherryPickedCommit,
		BranchName:      []byte(sourceBranch),
		Message:         []byte("Cherry-picking " + cherryPickedCommit.Id),
		StartBranchName: []byte(destinationBranch),
	}

	response, err := client.UserCherryPick(ctx, request)
	require.NoError(t, err)
	require.Equal(t, "Branch diverged", response.CommitError)
}

func TestServer_UserCherryPick_failedWithConflict(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickFailedWithConflict)
}

func testServerUserCherryPickFailedWithConflict(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	destinationBranch := "cherry-picking-dst"
	gittest.Exec(t, cfg, "-C", repoPath, "branch", destinationBranch, "conflict_branch_a")

	// This commit cannot be applied to the destinationBranch above
	cherryPickedCommit, err := repo.ReadCommit(ctx, git.Revision("f0f390655872bb2772c85a0128b2fbc2d88670cb"))
	require.NoError(t, err)

	request := &gitalypb.UserCherryPickRequest{
		Repository: repoProto,
		User:       gittest.TestUser,
		Commit:     cherryPickedCommit,
		BranchName: []byte(destinationBranch),
		Message:    []byte("Cherry-picking " + cherryPickedCommit.Id),
	}

	response, err := client.UserCherryPick(ctx, request)
	require.NoError(t, err)
	require.NotEmpty(t, response.CreateTreeError)
	require.Equal(t, gitalypb.UserCherryPickResponse_CONFLICT, response.CreateTreeErrorCode)
}

func TestServer_UserCherryPick_successfulWithGivenCommits(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickSuccessfulWithGivenCommits)
}

func testServerUserCherryPickSuccessfulWithGivenCommits(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	testCases := []struct {
		desc           string
		startRevision  git.Revision
		cherryRevision git.Revision
	}{
		{
			desc:           "merge commit",
			startRevision:  "281d3a76f31c812dbf48abce82ccf6860adedd81",
			cherryRevision: "6907208d755b60ebeacb2e9dfea74c92c3449a1f",
		},
	}

	for i, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			destinationBranch := fmt.Sprintf("cherry-picking-%d", i)

			gittest.Exec(t, cfg, "-C", repoPath, "branch", destinationBranch, testCase.startRevision.String())

			commit, err := repo.ReadCommit(ctx, testCase.cherryRevision)
			require.NoError(t, err)

			request := &gitalypb.UserCherryPickRequest{
				Repository: repoProto,
				User:       gittest.TestUser,
				Commit:     commit,
				BranchName: []byte(destinationBranch),
				Message:    []byte("Cherry-picking " + testCase.cherryRevision.String()),
			}

			response, err := client.UserCherryPick(ctx, request)
			require.NoError(t, err)

			newHead, err := repo.ReadCommit(ctx, git.Revision(destinationBranch))
			require.NoError(t, err)

			expectedResponse := &gitalypb.UserCherryPickResponse{
				BranchUpdate: &gitalypb.OperationBranchUpdate{CommitId: newHead.Id},
			}

			testassert.ProtoEqual(t, expectedResponse, response)

			require.Equal(t, request.Message, newHead.Subject)
			require.Equal(t, testCase.startRevision.String(), newHead.ParentIds[0])
		})
	}
}

func TestServer_UserCherryPick_quarantine(t *testing.T) {
	t.Parallel()

	testhelper.NewFeatureSets([]featureflag.FeatureFlag{
		featureflag.Quarantine,
	}).Run(t, testServerUserCherryPickQuarantine)
}

func testServerUserCherryPickQuarantine(t *testing.T, ctx context.Context) {
	ctx, cfg, repoProto, repoPath, client := setupOperationsService(t, ctx)
	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	// Set up a hook that parses the new object and then aborts the update. Like this, we can
	// assert that the object does not end up in the main repository.
	hookScript := fmt.Sprintf("#!/bin/sh\n%s rev-parse $3^{commit} && exit 1", cfg.Git.BinPath)
	gittest.WriteCustomHook(t, repoPath, "update", []byte(hookScript))

	commit, err := repo.ReadCommit(ctx, "8a0f2ee90d940bfb0ba1e14e8214b0649056e4ab")
	require.NoError(t, err)

	request := &gitalypb.UserCherryPickRequest{
		Repository: repoProto,
		User:       gittest.TestUser,
		Commit:     commit,
		BranchName: []byte("refs/heads/master"),
		Message:    []byte("Message"),
	}

	response, err := client.UserCherryPick(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, response)

	oid, err := git.NewObjectIDFromHex(strings.TrimSpace(response.PreReceiveError))
	require.NoError(t, err)
	exists, err := repo.HasRevision(ctx, oid.Revision()+"^{commit}")
	require.NoError(t, err)

	// The new commit will be in the target repository in case quarantines are disabled.
	// Otherwise, it should've been discarded.
	require.Equal(t, !featureflag.Quarantine.IsEnabled(ctx), exists)
}
