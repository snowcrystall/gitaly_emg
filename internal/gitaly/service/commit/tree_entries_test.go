package commit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
)

func TestSuccessfulGetTreeEntriesWithCurlyBraces(t *testing.T) {
	t.Parallel()
	cfg, repo, repoPath, client := setupCommitServiceWithRepo(t, false)

	normalFolderName := "issue-46261/folder"
	curlyFolderName := "issue-46261/{{curly}}"
	normalFolder := filepath.Join(repoPath, normalFolderName)
	curlyFolder := filepath.Join(repoPath, curlyFolderName)

	require.NoError(t, os.MkdirAll(normalFolder, 0755))
	require.NoError(t, os.MkdirAll(curlyFolder, 0755))

	testhelper.MustRunCommand(t, nil, "touch", filepath.Join(normalFolder, "/test1.txt"))
	testhelper.MustRunCommand(t, nil, "touch", filepath.Join(curlyFolder, "/test2.txt"))

	gittest.Exec(t, cfg, "-C", repoPath, "add", "--all")
	gittest.Exec(t, cfg, "-C", repoPath, "commit", "-m", "Test commit")

	testCases := []struct {
		description string
		revision    []byte
		path        []byte
		recursive   bool
		filename    []byte
	}{
		{
			description: "with a normal folder",
			revision:    []byte("master"),
			path:        []byte(normalFolderName),
			filename:    []byte("issue-46261/folder/test1.txt"),
		},
		{
			description: "with a folder with curly braces",
			revision:    []byte("master"),
			path:        []byte(curlyFolderName),
			filename:    []byte("issue-46261/{{curly}}/test2.txt"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			request := &gitalypb.GetTreeEntriesRequest{
				Repository: repo,
				Revision:   []byte("HEAD"),
				Path:       testCase.path,
				Recursive:  testCase.recursive,
			}

			ctx, cancel := testhelper.Context()
			defer cancel()
			c, err := client.GetTreeEntries(ctx, request)
			require.NoError(t, err)

			fetchedEntries, _ := getTreeEntriesFromTreeEntryClient(t, c, nil)
			require.Equal(t, 1, len(fetchedEntries))
			require.Equal(t, testCase.filename, fetchedEntries[0].FlatPath)
		})
	}
}

func TestSuccessfulGetTreeEntries(t *testing.T) {
	t.Parallel()
	commitID := "d25b6d94034242f3930dfcfeb6d8d9aac3583992"
	rootOid := "21bdc8af908562ae485ed46d71dd5426c08b084a"

	_, repo, _, client := setupCommitServiceWithRepo(t, true)

	rootEntries := []*gitalypb.TreeEntry{
		{
			Oid:       "fd90a3d2d21d6b4f9bec2c33fb7f49780c55f0d2",
			RootOid:   rootOid,
			Path:      []byte(".DS_Store"),
			FlatPath:  []byte(".DS_Store"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "470ad2fcf1e33798f1afc5781d08e60c40f51e7a",
			RootOid:   rootOid,
			Path:      []byte(".gitignore"),
			FlatPath:  []byte(".gitignore"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "fdaada1754989978413d618ee1fb1c0469d6a664",
			RootOid:   rootOid,
			Path:      []byte(".gitmodules"),
			FlatPath:  []byte(".gitmodules"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "c74175afd117781cbc983664339a0f599b5bb34e",
			RootOid:   rootOid,
			Path:      []byte("CHANGELOG"),
			FlatPath:  []byte("CHANGELOG"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "c1788657b95998a2f177a4f86d68a60f2a80117f",
			RootOid:   rootOid,
			Path:      []byte("CONTRIBUTING.md"),
			FlatPath:  []byte("CONTRIBUTING.md"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "50b27c6518be44c42c4d87966ae2481ce895624c",
			RootOid:   rootOid,
			Path:      []byte("LICENSE"),
			FlatPath:  []byte("LICENSE"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "95d9f0a5e7bb054e9dd3975589b8dfc689e20e88",
			RootOid:   rootOid,
			Path:      []byte("MAINTENANCE.md"),
			FlatPath:  []byte("MAINTENANCE.md"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "bf757025c40c62e6ffa6f11d3819c769a76dbe09",
			RootOid:   rootOid,
			Path:      []byte("PROCESS.md"),
			FlatPath:  []byte("PROCESS.md"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "faaf198af3a36dbf41961466703cc1d47c61d051",
			RootOid:   rootOid,
			Path:      []byte("README.md"),
			FlatPath:  []byte("README.md"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "998707b421c89bd9a3063333f9f728ef3e43d101",
			RootOid:   rootOid,
			Path:      []byte("VERSION"),
			FlatPath:  []byte("VERSION"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "3c122d2b7830eca25235131070602575cf8b41a1",
			RootOid:   rootOid,
			Path:      []byte("encoding"),
			FlatPath:  []byte("encoding"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "b4a3321157f6e80c42b031ecc9ba79f784c8a557",
			RootOid:   rootOid,
			Path:      []byte("files"),
			FlatPath:  []byte("files"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "6fd00c6336d6385ef6efe553a29107b35d18d380",
			RootOid:   rootOid,
			Path:      []byte("level-0"),
			FlatPath:  []byte("level-0"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "409f37c4f05865e4fb208c771485f211a22c4c2d",
			RootOid:   rootOid,
			Path:      []byte("six"),
			FlatPath:  []byte("six"),
			Type:      gitalypb.TreeEntry_COMMIT,
			Mode:      0160000,
			CommitOid: commitID,
		},
	}

	// Order: Tree, Blob, Submodules
	sortedRootEntries := append(rootEntries[10:13], rootEntries[0:10]...)
	sortedRootEntries = append(sortedRootEntries, rootEntries[13])
	sortedAndPaginated := []*gitalypb.TreeEntry{rootEntries[10], rootEntries[11], rootEntries[12], rootEntries[0]}

	filesDirEntries := []*gitalypb.TreeEntry{
		{
			Oid:       "60d7a906c2fd9e4509aeb1187b98d0ea7ce827c9",
			RootOid:   rootOid,
			Path:      []byte("files/.DS_Store"),
			FlatPath:  []byte("files/.DS_Store"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "2132d150328bd9334cc4e62a16a5d998a7e399b9",
			RootOid:   rootOid,
			Path:      []byte("files/flat"),
			FlatPath:  []byte("files/flat/path/correct"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "a1e8f8d745cc87e3a9248358d9352bb7f9a0aeba",
			RootOid:   rootOid,
			Path:      []byte("files/html"),
			FlatPath:  []byte("files/html"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "5e147e3af6740ee83103ec2ecdf846cae696edd1",
			RootOid:   rootOid,
			Path:      []byte("files/images"),
			FlatPath:  []byte("files/images"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "7853101769f3421725ddc41439c2cd4610e37ad9",
			RootOid:   rootOid,
			Path:      []byte("files/js"),
			FlatPath:  []byte("files/js"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "fd581c619bf59cfdfa9c8282377bb09c2f897520",
			RootOid:   rootOid,
			Path:      []byte("files/markdown"),
			FlatPath:  []byte("files/markdown"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "b59dbe4a27371d53e61bf3cb8bef66be53572db0",
			RootOid:   rootOid,
			Path:      []byte("files/ruby"),
			FlatPath:  []byte("files/ruby"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
	}

	recursiveEntries := []*gitalypb.TreeEntry{
		{
			Oid:       "d564d0bc3dd917926892c55e3706cc116d5b165e",
			RootOid:   rootOid,
			Path:      []byte("level-0/level-1-1"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
			RootOid:   rootOid,
			Path:      []byte("level-0/level-1-1/.gitkeep"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
		{
			Oid:       "02366a40d0cde8191e43a8c5b821176c0668522c",
			RootOid:   rootOid,
			Path:      []byte("level-0/level-1-2"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "d564d0bc3dd917926892c55e3706cc116d5b165e",
			RootOid:   rootOid,
			Path:      []byte("level-0/level-1-2/level-2"),
			Type:      gitalypb.TreeEntry_TREE,
			Mode:      040000,
			CommitOid: commitID,
		},
		{
			Oid:       "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
			RootOid:   rootOid,
			Path:      []byte("level-0/level-1-2/level-2/.gitkeep"),
			Type:      gitalypb.TreeEntry_BLOB,
			Mode:      0100644,
			CommitOid: commitID,
		},
	}

	testCases := []struct {
		description string
		revision    []byte
		path        []byte
		recursive   bool
		sortBy      gitalypb.GetTreeEntriesRequest_SortBy
		entries     []*gitalypb.TreeEntry
		pageToken   string
		pageLimit   int32
		cursor      string
	}{
		{
			description: "with root path",
			revision:    []byte(commitID),
			path:        []byte("."),
			entries:     rootEntries,
		},
		{
			description: "with a folder",
			revision:    []byte(commitID),
			path:        []byte("files"),
			entries:     filesDirEntries,
		},
		{
			description: "with recursive",
			revision:    []byte(commitID),
			path:        []byte("level-0"),
			recursive:   true,
			entries:     recursiveEntries,
		},
		{
			description: "with a file",
			revision:    []byte(commitID),
			path:        []byte(".gitignore"),
			entries:     nil,
		},
		{
			description: "with a non-existing path",
			revision:    []byte(commitID),
			path:        []byte("i-dont/exist"),
			entries:     nil,
		},
		{
			description: "with root path and sorted by trees first",
			revision:    []byte(commitID),
			path:        []byte("."),
			entries:     sortedRootEntries,
			sortBy:      gitalypb.GetTreeEntriesRequest_TREES_FIRST,
		},
		{
			description: "with root path and sorted by trees with pagination",
			revision:    []byte(commitID),
			path:        []byte("."),
			entries:     sortedAndPaginated,
			pageLimit:   4,
			sortBy:      gitalypb.GetTreeEntriesRequest_TREES_FIRST,
			cursor:      "fd90a3d2d21d6b4f9bec2c33fb7f49780c55f0d2",
		},
		{
			description: "with pagination parameters",
			revision:    []byte(commitID),
			path:        []byte("."),
			entries:     rootEntries[3:6],
			pageToken:   "fdaada1754989978413d618ee1fb1c0469d6a664",
			pageLimit:   3,
			cursor:      rootEntries[5].Oid,
		},
		{
			description: "with pagination parameters larger than length",
			revision:    []byte(commitID),
			path:        []byte("."),
			entries:     rootEntries[12:],
			pageToken:   "b4a3321157f6e80c42b031ecc9ba79f784c8a557",
			pageLimit:   20,
		},
		{
			description: "with pagination limit of -1",
			revision:    []byte(commitID),
			path:        []byte("."),
			entries:     rootEntries[2:],
			pageToken:   "470ad2fcf1e33798f1afc5781d08e60c40f51e7a",
			pageLimit:   -1,
		},
		{
			description: "with pagination limit of 0",
			revision:    []byte(commitID),
			path:        []byte("."),
			pageToken:   "470ad2fcf1e33798f1afc5781d08e60c40f51e7a",
			pageLimit:   0,
		},
		{
			description: "with a blank pagination token",
			revision:    []byte(commitID),
			path:        []byte("."),
			pageToken:   "",
			entries:     rootEntries[0:2],
			pageLimit:   2,
			cursor:      rootEntries[1].Oid,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			request := &gitalypb.GetTreeEntriesRequest{
				Repository: repo,
				Revision:   testCase.revision,
				Path:       testCase.path,
				Recursive:  testCase.recursive,
				Sort:       testCase.sortBy,
			}

			if testCase.pageToken != "" || testCase.pageLimit > 0 {
				request.PaginationParams = &gitalypb.PaginationParameter{
					PageToken: testCase.pageToken,
					Limit:     testCase.pageLimit,
				}
			}

			ctx, cancel := testhelper.Context()
			defer cancel()
			c, err := client.GetTreeEntries(ctx, request)

			require.NoError(t, err)
			fetchedEntries, cursor := getTreeEntriesFromTreeEntryClient(t, c, nil)
			require.Equal(t, testCase.entries, fetchedEntries)

			if testCase.pageLimit > 0 && len(testCase.entries) < len(rootEntries) {
				require.NotNil(t, cursor)
				require.Equal(t, testCase.cursor, cursor.NextCursor)
			}
		})
	}
}

func TestUnsuccessfulGetTreeEntries(t *testing.T) {
	commitID := "d25b6d94034242f3930dfcfeb6d8d9aac3583992"

	_, repo, _, client := setupCommitServiceWithRepo(t, true)

	testCases := []struct {
		description   string
		revision      []byte
		path          []byte
		pageToken     string
		expectedError error
	}{
		{
			description:   "with non-existent token",
			revision:      []byte(commitID),
			path:          []byte("."),
			pageToken:     "non-existent",
			expectedError: fmt.Errorf("could not get find starting OID: non-existent"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			request := &gitalypb.GetTreeEntriesRequest{
				Repository: repo,
				Revision:   testCase.revision,
				Path:       testCase.path,
			}

			if testCase.pageToken != "" {
				request.PaginationParams = &gitalypb.PaginationParameter{
					PageToken: testCase.pageToken,
				}
			}

			ctx, cancel := testhelper.Context()
			defer cancel()
			c, err := client.GetTreeEntries(ctx, request)
			require.NoError(t, err)

			fetchedEntries, cursor := getTreeEntriesFromTreeEntryClient(t, c, testCase.expectedError)

			require.Empty(t, fetchedEntries)
			require.Nil(t, cursor)
		})
	}
}

func getTreeEntriesFromTreeEntryClient(t *testing.T, client gitalypb.CommitService_GetTreeEntriesClient, expectedError error) ([]*gitalypb.TreeEntry, *gitalypb.PaginationCursor) {
	t.Helper()

	var entries []*gitalypb.TreeEntry
	var cursor *gitalypb.PaginationCursor
	firstEntryReceived := false

	for {
		resp, err := client.Recv()
		if err == io.EOF {
			break
		}

		if expectedError == nil {
			require.NoError(t, err)
			entries = append(entries, resp.Entries...)

			if !firstEntryReceived {
				cursor = resp.PaginationCursor
				firstEntryReceived = true
			} else {
				require.Equal(t, nil, resp.PaginationCursor)
			}
		} else {
			require.Error(t, expectedError, err)
			break
		}
	}
	return entries, cursor
}

func TestSuccessfulGetTreeEntries_FlatPathMaxDeep_SingleFoldersStructure(t *testing.T) {
	t.Parallel()
	cfg, repo, repoPath, client := setupCommitServiceWithRepo(t, false)

	folderName := "1/2/3/4/5/6/7/8/9/10/11/12"
	require.GreaterOrEqual(t, strings.Count(strings.Trim(folderName, "/"), "/"), defaultFlatTreeRecursion, "sanity check: construct folder deeper than default recursion value")

	nestedFolder := filepath.Join(repoPath, folderName)
	require.NoError(t, os.MkdirAll(nestedFolder, 0755))
	// put single file into the deepest directory
	testhelper.MustRunCommand(t, nil, "touch", filepath.Join(nestedFolder, ".gitkeep"))
	gittest.Exec(t, cfg, "-C", repoPath, "add", "--all")
	gittest.Exec(t, cfg, "-C", repoPath, "commit", "-m", "Deep folder struct")

	commitID := text.ChompBytes(gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", "HEAD"))
	rootOid := text.ChompBytes(gittest.Exec(t, cfg, "-C", repoPath, "rev-parse", "HEAD^{tree}"))

	// make request to folder that contains nothing except one folder
	request := &gitalypb.GetTreeEntriesRequest{
		Repository: repo,
		Revision:   []byte(commitID),
		Path:       []byte("1"),
		Recursive:  false,
	}

	ctx, cancel := testhelper.Context()
	defer cancel()

	// request entries of the tree with single-folder structure on each level
	entriesClient, err := client.GetTreeEntries(ctx, request)
	require.NoError(t, err)

	fetchedEntries, _ := getTreeEntriesFromTreeEntryClient(t, entriesClient, nil)
	// We know that there is a directory "1/2/3/4/5/6/7/8/9/10/11/12"
	// but here we only get back "1/2/3/4/5/6/7/8/9/10/11".
	// This proves that FlatPath recursion is bounded, which is the point of this test.
	require.Equal(t, []*gitalypb.TreeEntry{{
		Oid:       "c836b95b37958e7179f5a42a32b7197b5dec7321",
		RootOid:   rootOid,
		Path:      []byte("1/2"),
		FlatPath:  []byte("1/2/3/4/5/6/7/8/9/10/11"),
		Type:      gitalypb.TreeEntry_TREE,
		Mode:      040000,
		CommitOid: commitID,
	}}, fetchedEntries)
}

func TestFailedGetTreeEntriesRequestDueToValidationError(t *testing.T) {
	t.Parallel()
	_, repo, _, client := setupCommitServiceWithRepo(t, true)

	revision := []byte("d42783470dc29fde2cf459eb3199ee1d7e3f3a72")
	path := []byte("a/b/c")

	rpcRequests := []*gitalypb.GetTreeEntriesRequest{
		{Repository: &gitalypb.Repository{StorageName: "fake", RelativePath: "path"}, Revision: revision, Path: path}, // Repository doesn't exist
		{Repository: nil, Revision: revision, Path: path},                                                             // Repository is nil
		{Repository: repo, Revision: nil, Path: path},                                                                 // Revision is empty
		{Repository: repo, Revision: revision},                                                                        // Path is empty
		{Repository: repo, Revision: []byte("--output=/meow"), Path: path},                                            // Revision is invalid
	}

	for _, rpcRequest := range rpcRequests {
		t.Run(fmt.Sprintf("%v", rpcRequest), func(t *testing.T) {
			ctx, cancel := testhelper.Context()
			defer cancel()
			c, err := client.GetTreeEntries(ctx, rpcRequest)
			require.NoError(t, err)

			err = drainTreeEntriesResponse(c)
			testhelper.RequireGrpcError(t, err, codes.InvalidArgument)
		})
	}
}

func drainTreeEntriesResponse(c gitalypb.CommitService_GetTreeEntriesClient) error {
	var err error
	for err == nil {
		_, err = c.Recv()
	}
	return err
}
