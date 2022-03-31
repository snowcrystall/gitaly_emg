package wiki

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/rubyserver"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
)

func testSuccessfulWikiUpdatePageRequest(t *testing.T, cfg config.Cfg, rubySrv *rubyserver.Server) {
	wikiRepoProto, wikiRepoPath := setupWikiRepo(t, cfg)
	wikiRepo := localrepo.NewTestRepo(t, cfg, wikiRepoProto)

	ctx, cancel := testhelper.Context()
	defer cancel()

	client := setupWikiService(t, cfg, rubySrv)

	writeWikiPage(t, client, wikiRepoProto, createWikiPageOpts{title: "Instálling Gitaly", content: []byte("foobar")})

	authorID := int32(1)
	authorUserName := []byte("ahmad")
	authorName := []byte("Ahmad Sherif")
	authorEmail := []byte("ahmad@gitlab.com")
	message := []byte("Add installation instructions")

	testCases := []struct {
		desc    string
		req     *gitalypb.WikiUpdatePageRequest
		content []byte
	}{
		{
			desc: "with user id and username",
			req: &gitalypb.WikiUpdatePageRequest{
				Repository: wikiRepoProto,
				PagePath:   []byte("//Instálling Gitaly"),
				Title:      []byte("Instálling Gitaly"),
				Format:     "markdown",
				CommitDetails: &gitalypb.WikiCommitDetails{
					Name:     authorName,
					Email:    authorEmail,
					Message:  message,
					UserId:   authorID,
					UserName: authorUserName,
				},
			},
			content: bytes.Repeat([]byte("Mock wiki page content"), 10000),
		},
		{
			desc: "without user id and username", // deprecate in gitlab 11.0 https://gitlab.com/gitlab-org/gitaly/issues/1154
			req: &gitalypb.WikiUpdatePageRequest{
				Repository: wikiRepoProto,
				PagePath:   []byte("//Instálling Gitaly"),
				Title:      []byte("Instálling Gitaly"),
				Format:     "markdown",
				CommitDetails: &gitalypb.WikiCommitDetails{
					Name:    authorName,
					Email:   authorEmail,
					Message: message,
				},
			},
			content: bytes.Repeat([]byte("Mock wiki page content 2"), 10000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			stream, err := client.WikiUpdatePage(ctx)
			require.NoError(t, err)

			request := tc.req
			nSends, err := sendBytes(tc.content, 1000, func(p []byte) error {
				request.Content = p

				if err := stream.Send(request); err != nil {
					return err
				}

				// Use a new response so we don't send other fields (Repository, ...) over and over
				request = &gitalypb.WikiUpdatePageRequest{}

				return nil
			})

			require.NoError(t, err)
			require.True(t, nSends > 1, "should have sent more than one message")

			_, err = stream.CloseAndRecv()
			require.NoError(t, err)

			headID := gittest.Exec(t, cfg, "-C", wikiRepoPath, "show", "--format=format:%H", "--no-patch", "HEAD")
			commit, err := wikiRepo.ReadCommit(ctx, git.Revision(headID))
			require.NoError(t, err, "look up git commit before merge is applied")

			require.Equal(t, authorName, commit.Author.Name, "author name mismatched")
			require.Equal(t, authorEmail, commit.Author.Email, "author email mismatched")
			require.Equal(t, message, commit.Subject, "message mismatched")

			pageContent := gittest.Exec(t, cfg, "-C", wikiRepoPath, "cat-file", "blob", "HEAD:Instálling-Gitaly.md")
			require.Equal(t, tc.content, pageContent, "mismatched content")
		})
	}
}

func testFailedWikiUpdatePageDueToValidations(t *testing.T, cfg config.Cfg, rubySrv *rubyserver.Server) {
	wikiRepo, _ := setupWikiRepo(t, cfg)

	client := setupWikiService(t, cfg, rubySrv)

	writeWikiPage(t, client, wikiRepo, createWikiPageOpts{title: "Installing Gitaly", content: []byte("foobar")})

	commitDetails := &gitalypb.WikiCommitDetails{
		Name:     []byte("Ahmad Sherif"),
		Email:    []byte("ahmad@gitlab.com"),
		Message:  []byte("Add installation instructions"),
		UserId:   int32(1),
		UserName: []byte("ahmad"),
	}

	testCases := []struct {
		desc    string
		request *gitalypb.WikiUpdatePageRequest
		code    codes.Code
	}{
		{
			desc: "empty page_path",
			request: &gitalypb.WikiUpdatePageRequest{
				Repository:    wikiRepo,
				Title:         []byte("Installing Gitaly"),
				Format:        "markdown",
				CommitDetails: commitDetails,
				Content:       []byte(""),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "page does not exist",
			request: &gitalypb.WikiUpdatePageRequest{
				Repository:    wikiRepo,
				PagePath:      []byte("//Installing Gibaly"),
				Title:         []byte("Installing Gitaly"),
				Format:        "markdown",
				CommitDetails: commitDetails,
				Content:       []byte(""),
			},
			code: codes.NotFound,
		},
		{
			desc: "empty title",
			request: &gitalypb.WikiUpdatePageRequest{
				Repository:    wikiRepo,
				PagePath:      []byte("//Installing Gitaly"),
				Format:        "markdown",
				CommitDetails: commitDetails,
				Content:       []byte(""),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty format",
			request: &gitalypb.WikiUpdatePageRequest{
				Repository:    wikiRepo,
				PagePath:      []byte("//Installing Gitaly.md"),
				Title:         []byte("Installing Gitaly"),
				CommitDetails: commitDetails,
				Content:       []byte(""),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty commit details",
			request: &gitalypb.WikiUpdatePageRequest{
				Repository: wikiRepo,
				PagePath:   []byte("//Installing Gitaly.md"),
				Title:      []byte("Installing Gitaly"),
				Format:     "markdown",
				Content:    []byte(""),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty commit details' name",
			request: &gitalypb.WikiUpdatePageRequest{
				Repository: wikiRepo,
				PagePath:   []byte("//Installing Gitaly.md"),
				Title:      []byte("Installing Gitaly"),
				Format:     "markdown",
				CommitDetails: &gitalypb.WikiCommitDetails{
					Email:    []byte("a@b.com"),
					Message:  []byte("A message"),
					UserId:   int32(1),
					UserName: []byte("username"),
				},
				Content: []byte(""),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty commit details' email",
			request: &gitalypb.WikiUpdatePageRequest{
				Repository: wikiRepo,
				PagePath:   []byte("//Installing Gitaly.md"),
				Title:      []byte("Installing Gitaly"),
				Format:     "markdown",
				CommitDetails: &gitalypb.WikiCommitDetails{
					Name:     []byte("A name"),
					Message:  []byte("A message"),
					UserId:   int32(1),
					UserName: []byte("username"),
				},
				Content: []byte(""),
			},
			code: codes.InvalidArgument,
		},
		{
			desc: "empty commit details' message",
			request: &gitalypb.WikiUpdatePageRequest{
				Repository: wikiRepo,
				PagePath:   []byte("//Installing Gitaly.md"),
				Title:      []byte("Installing Gitaly"),
				Format:     "markdown",
				CommitDetails: &gitalypb.WikiCommitDetails{
					Name:     []byte("A name"),
					Email:    []byte("a@b.com"),
					UserId:   int32(1),
					UserName: []byte("username"),
				},
				Content: []byte(""),
			},
			code: codes.InvalidArgument,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()

			stream, err := client.WikiUpdatePage(ctx)
			require.NoError(t, err)

			require.NoError(t, stream.Send(testCase.request))

			_, err = stream.CloseAndRecv()
			testhelper.RequireGrpcError(t, err, testCase.code)
		})
	}
}
