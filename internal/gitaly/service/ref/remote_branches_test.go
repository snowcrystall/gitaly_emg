package ref

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
)

func TestSuccessfulFindAllRemoteBranchesRequest(t *testing.T) {
	ctx, cancel := testhelper.Context()
	defer cancel()

	cfg, repoProto, repoPath, client := setupRefService(t)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	remoteName := "my-remote"
	expectedBranches := map[string]string{
		"foo": "c7fbe50c7c7419d9701eebe64b1fdacc3df5b9dd",
		"bar": "60ecb67744cb56576c30214ff52294f8ce2def98",
	}
	excludedRemote := "my-remote-2"
	excludedBranches := map[string]string{
		"from-another-remote": "5937ac0a7beb003549fc5fd26fc247adbce4a52e",
	}

	for branchName, commitID := range expectedBranches {
		gittest.CreateRemoteBranch(t, cfg, repoPath, remoteName, branchName, commitID)
	}

	for branchName, commitID := range excludedBranches {
		gittest.CreateRemoteBranch(t, cfg, repoPath, excludedRemote, branchName, commitID)
	}

	request := &gitalypb.FindAllRemoteBranchesRequest{Repository: repoProto, RemoteName: remoteName}

	c, err := client.FindAllRemoteBranches(ctx, request)
	require.NoError(t, err)

	branches := readFindAllRemoteBranchesResponsesFromClient(t, c)
	require.Len(t, branches, len(expectedBranches))

	for branchName, commitID := range expectedBranches {
		targetCommit, err := repo.ReadCommit(ctx, git.Revision(commitID))
		require.NoError(t, err)

		expectedBranch := &gitalypb.Branch{
			Name:         []byte("refs/remotes/" + remoteName + "/" + branchName),
			TargetCommit: targetCommit,
		}

		require.Contains(t, branches, expectedBranch)
	}

	for branchName, commitID := range excludedBranches {
		targetCommit, err := repo.ReadCommit(ctx, git.Revision(commitID))
		require.NoError(t, err)

		excludedBranch := &gitalypb.Branch{
			Name:         []byte("refs/remotes/" + excludedRemote + "/" + branchName),
			TargetCommit: targetCommit,
		}

		require.NotContains(t, branches, excludedBranch)
	}
}

func TestInvalidFindAllRemoteBranchesRequest(t *testing.T) {
	_, repo, _, client := setupRefService(t)

	testCases := []struct {
		description string
		request     *gitalypb.FindAllRemoteBranchesRequest
	}{
		{
			description: "Invalid repo",
			request: &gitalypb.FindAllRemoteBranchesRequest{
				Repository: &gitalypb.Repository{
					StorageName:  "fake",
					RelativePath: "repo",
				},
			},
		},
		{
			description: "Empty repo",
			request:     &gitalypb.FindAllRemoteBranchesRequest{RemoteName: "myRemote"},
		},
		{
			description: "Empty remote name",
			request:     &gitalypb.FindAllRemoteBranchesRequest{Repository: repo},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()
			c, err := client.FindAllRemoteBranches(ctx, tc.request)
			require.NoError(t, err)

			var recvError error
			for recvError == nil {
				_, recvError = c.Recv()
			}

			testhelper.RequireGrpcError(t, recvError, codes.InvalidArgument)
		})
	}
}

func readFindAllRemoteBranchesResponsesFromClient(t *testing.T, c gitalypb.RefService_FindAllRemoteBranchesClient) []*gitalypb.Branch {
	var branches []*gitalypb.Branch

	for {
		r, err := c.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		branches = append(branches, r.GetBranches()...)
	}

	return branches
}
