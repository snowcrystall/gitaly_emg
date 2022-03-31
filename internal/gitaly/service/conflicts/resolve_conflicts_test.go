package conflicts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/hook"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testassert"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	user = &gitalypb.User{
		Name:  []byte("John Doe"),
		Email: []byte("johndoe@gitlab.com"),
		GlId:  "user-1",
	}
	conflictResolutionCommitMessage = "Solve conflicts"

	files = []map[string]interface{}{
		{
			"old_path": "files/ruby/popen.rb",
			"new_path": "files/ruby/popen.rb",
			"sections": map[string]string{
				"2f6fcd96b88b36ce98c38da085c795a27d92a3dd_14_14": "head",
			},
		},
		{
			"old_path": "files/ruby/regex.rb",
			"new_path": "files/ruby/regex.rb",
			"sections": map[string]string{
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_9_9":   "head",
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_21_21": "origin",
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_49_49": "origin",
			},
		},
	}
)

func TestSuccessfulResolveConflictsRequestHelper(t *testing.T) {
	cfg, repoProto, repoPath := SetupConfigAndRepo(t, true)
	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	ctx, cancel := testhelper.Context()
	defer cancel()

	missingAncestorPath := "files/missing_ancestor.txt"
	files := []map[string]interface{}{
		{
			"old_path": "files/ruby/popen.rb",
			"new_path": "files/ruby/popen.rb",
			"sections": map[string]string{
				"2f6fcd96b88b36ce98c38da085c795a27d92a3dd_14_14": "head",
			},
		},
		{
			"old_path": "files/ruby/regex.rb",
			"new_path": "files/ruby/regex.rb",
			"sections": map[string]string{
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_9_9":   "head",
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_21_21": "origin",
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_49_49": "origin",
			},
		},
		{
			"old_path": missingAncestorPath,
			"new_path": missingAncestorPath,
			"sections": map[string]string{
				"b760bfd3b1b1da380b4276eb30fb3b2b7e4f08e1_1_1": "origin",
			},
		},
	}

	filesJSON, err := json.Marshal(files)
	require.NoError(t, err)

	sourceBranch := "conflict-resolvable"
	targetBranch := "conflict-start"
	ourCommitOID := "1450cd639e0bc6721eb02800169e464f212cde06"   // part of branch conflict-resolvable
	theirCommitOID := "824be604a34828eb682305f0d963056cfac87b2d" // part of branch conflict-start
	ancestorCommitOID := "6907208d755b60ebeacb2e9dfea74c92c3449a1f"

	// introduce a conflict that exists on both branches, but not the
	// ancestor
	commitConflict := func(parentCommitID, branch, blob string) string {
		blobID, err := repo.WriteBlob(ctx, "", strings.NewReader(blob))
		require.NoError(t, err)
		gittest.Exec(t, cfg, "-C", repoPath, "read-tree", branch)
		gittest.Exec(t, cfg, "-C", repoPath,
			"update-index", "--add", "--cacheinfo", "100644", blobID.String(), missingAncestorPath,
		)
		treeID := bytes.TrimSpace(
			gittest.Exec(t, cfg, "-C", repoPath, "write-tree"),
		)
		commitID := bytes.TrimSpace(
			gittest.Exec(t, cfg, "-C", repoPath,
				"commit-tree", string(treeID), "-p", parentCommitID,
			),
		)
		gittest.Exec(t, cfg, "-C", repoPath, "update-ref", "refs/heads/"+branch, string(commitID))
		return string(commitID)
	}

	// sanity check: make sure the conflict file does not exist on the
	// common ancestor
	cmd := exec.CommandContext(ctx, "git", "cat-file", "-e", ancestorCommitOID+":"+missingAncestorPath)
	require.Error(t, cmd.Run())

	ourCommitOID = commitConflict(ourCommitOID, sourceBranch, "content-1")
	theirCommitOID = commitConflict(theirCommitOID, targetBranch, "content-2")
	hookCount := 0

	verifyFunc := func(t *testing.T, ctx context.Context, repo *gitalypb.Repository, pushOptions, env []string, stdin io.Reader, stdout, stderr io.Writer) error {
		changes, err := ioutil.ReadAll(stdin)
		require.NoError(t, err)
		pattern := fmt.Sprintf("%s .* refs/heads/%s\n", ourCommitOID, sourceBranch)
		require.Regexp(t, regexp.MustCompile(pattern), string(changes))
		require.Empty(t, pushOptions)
		hookCount++
		return nil
	}

	hookManager := hook.NewMockManager(t, verifyFunc, verifyFunc, hook.NopUpdate, hook.NopReferenceTransaction)
	client := SetupConflictsServiceWithConfig(t, &cfg, hookManager)

	mdGS := testhelper.GitalyServersMetadataFromCfg(t, cfg)
	mdFF, _ := metadata.FromOutgoingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Join(mdGS, mdFF))

	headerRequest := &gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
			Header: &gitalypb.ResolveConflictsRequestHeader{
				Repository:       repoProto,
				TargetRepository: repoProto,
				CommitMessage:    []byte(conflictResolutionCommitMessage),
				OurCommitOid:     ourCommitOID,
				TheirCommitOid:   theirCommitOID,
				SourceBranch:     []byte(sourceBranch),
				TargetBranch:     []byte(targetBranch),
				User:             user,
			},
		},
	}
	filesRequest1 := &gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
			FilesJson: filesJSON[:50],
		},
	}
	filesRequest2 := &gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
			FilesJson: filesJSON[50:],
		},
	}

	stream, err := client.ResolveConflicts(ctx)
	require.NoError(t, err)
	require.NoError(t, stream.Send(headerRequest))
	require.NoError(t, stream.Send(filesRequest1))
	require.NoError(t, stream.Send(filesRequest2))

	r, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.Empty(t, r.GetResolutionError())

	headCommit, err := repo.ReadCommit(ctx, git.Revision(sourceBranch))
	require.NoError(t, err)
	require.Contains(t, headCommit.ParentIds, ourCommitOID)
	require.Contains(t, headCommit.ParentIds, theirCommitOID)
	require.Equal(t, string(headCommit.Author.Email), "johndoe@gitlab.com")
	require.Equal(t, string(headCommit.Committer.Email), "johndoe@gitlab.com")
	require.Equal(t, string(headCommit.Subject), conflictResolutionCommitMessage)

	require.Equal(t, 2, hookCount)
}

func TestResolveConflictsWithRemoteRepo(t *testing.T) {
	hookManager := hook.NewMockManager(t, hook.NopPreReceive, hook.NopPostReceive, hook.NopUpdate, hook.NopReferenceTransaction)
	cfg, _, _, client := SetupConflictsService(t, true, hookManager)

	ctx, cancel := testhelper.Context()
	defer cancel()

	testhelper.BuildGitalySSH(t, cfg)
	testhelper.BuildGitalyHooks(t, cfg)

	sourceRepo, sourceRepoPath := gittest.CloneRepo(t, cfg, cfg.Storages[0])
	sourceBlobOID := gittest.WriteBlob(t, cfg, sourceRepoPath, []byte("contents-1\n"))
	sourceCommitOID := gittest.WriteCommit(t, cfg, sourceRepoPath,
		gittest.WithTreeEntries(gittest.TreeEntry{
			Path: "file.txt", OID: sourceBlobOID, Mode: "100644",
		}),
	)
	gittest.Exec(t, cfg, "-C", sourceRepoPath, "update-ref", "refs/heads/source", sourceCommitOID.String())

	targetRepo, targetRepoPath := gittest.CloneRepo(t, cfg, cfg.Storages[0])
	targetBlobOID := gittest.WriteBlob(t, cfg, targetRepoPath, []byte("contents-2\n"))
	targetCommitOID := gittest.WriteCommit(t, cfg, targetRepoPath,
		gittest.WithTreeEntries(gittest.TreeEntry{
			OID: targetBlobOID, Path: "file.txt", Mode: "100644",
		}),
	)
	gittest.Exec(t, cfg, "-C", targetRepoPath, "update-ref", "refs/heads/target", targetCommitOID.String())

	ctx = testhelper.MergeOutgoingMetadata(ctx, testhelper.GitalyServersMetadataFromCfg(t, cfg))

	stream, err := client.ResolveConflicts(ctx)
	require.NoError(t, err)

	filesJSON, err := json.Marshal([]map[string]interface{}{
		{
			"old_path": "file.txt",
			"new_path": "file.txt",
			"sections": map[string]string{
				"5436437fa01a7d3e41d46741da54b451446774ca_1_1": "origin",
			},
		},
	})
	require.NoError(t, err)

	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
			Header: &gitalypb.ResolveConflictsRequestHeader{
				Repository:       sourceRepo,
				TargetRepository: targetRepo,
				CommitMessage:    []byte(conflictResolutionCommitMessage),
				OurCommitOid:     sourceCommitOID.String(),
				TheirCommitOid:   targetCommitOID.String(),
				SourceBranch:     []byte("source"),
				TargetBranch:     []byte("target"),
				User:             user,
			},
		},
	}))
	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
			FilesJson: filesJSON,
		},
	}))

	response, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.Empty(t, response.GetResolutionError())

	require.Equal(t, []byte("contents-2\n"), gittest.Exec(t, cfg, "-C", sourceRepoPath, "cat-file", "-p", "refs/heads/source:file.txt"))
}

func TestResolveConflictsLineEndings(t *testing.T) {
	hookManager := hook.NewMockManager(t, hook.NopPreReceive, hook.NopPostReceive, hook.NopUpdate, hook.NopReferenceTransaction)
	cfg, repo, repoPath, client := SetupConflictsService(t, true, hookManager)

	ctx, cancel := testhelper.Context()
	defer cancel()
	ctx = testhelper.MergeOutgoingMetadata(ctx, testhelper.GitalyServersMetadataFromCfg(t, cfg))

	for _, tc := range []struct {
		desc             string
		ourContent       string
		theirContent     string
		resolutions      []map[string]interface{}
		expectedContents string
	}{
		{
			desc:             "only newline",
			ourContent:       "\n",
			theirContent:     "\n",
			resolutions:      []map[string]interface{}{},
			expectedContents: "\n",
		},
		{
			desc:         "conflicting newline with embedded character",
			ourContent:   "\nA\n",
			theirContent: "\nB\n",
			resolutions: []map[string]interface{}{
				{
					"old_path": "file.txt",
					"new_path": "file.txt",
					"sections": map[string]string{
						"5436437fa01a7d3e41d46741da54b451446774ca_2_2": "head",
					},
				},
			},
			expectedContents: "\nA\n",
		},
		{
			desc:         "conflicting carriage-return newlines",
			ourContent:   "A\r\nB\r\nC\r\nD\r\nE\r\n",
			theirContent: "A\r\nB\r\nX\r\nD\r\nE\r\n",
			resolutions: []map[string]interface{}{
				{
					"old_path": "file.txt",
					"new_path": "file.txt",
					"sections": map[string]string{
						"5436437fa01a7d3e41d46741da54b451446774ca_3_3": "origin",
					},
				},
			},
			expectedContents: "A\r\nB\r\nX\r\nD\r\nE\r\n",
		},
		{
			desc:         "conflict with no trailing newline",
			ourContent:   "A\nB",
			theirContent: "X\nB",
			resolutions: []map[string]interface{}{
				{
					"old_path": "file.txt",
					"new_path": "file.txt",
					"sections": map[string]string{
						"5436437fa01a7d3e41d46741da54b451446774ca_1_1": "head",
					},
				},
			},
			expectedContents: "A\nB",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ourOID := gittest.WriteBlob(t, cfg, repoPath, []byte(tc.ourContent))
			ourCommit := gittest.WriteCommit(t, cfg, repoPath,
				gittest.WithTreeEntries(gittest.TreeEntry{
					OID: ourOID, Path: "file.txt", Mode: "100644",
				}),
			)
			gittest.Exec(t, cfg, "-C", repoPath, "update-ref", "refs/heads/ours", ourCommit.String())

			theirOID := gittest.WriteBlob(t, cfg, repoPath, []byte(tc.theirContent))
			theirCommit := gittest.WriteCommit(t, cfg, repoPath,
				gittest.WithTreeEntries(gittest.TreeEntry{
					OID: theirOID, Path: "file.txt", Mode: "100644",
				}),
			)
			gittest.Exec(t, cfg, "-C", repoPath, "update-ref", "refs/heads/theirs", theirCommit.String())

			stream, err := client.ResolveConflicts(ctx)
			require.NoError(t, err)

			filesJSON, err := json.Marshal(tc.resolutions)
			require.NoError(t, err)

			require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
				ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
					Header: &gitalypb.ResolveConflictsRequestHeader{
						Repository:       repo,
						TargetRepository: repo,
						CommitMessage:    []byte(conflictResolutionCommitMessage),
						OurCommitOid:     ourCommit.String(),
						TheirCommitOid:   theirCommit.String(),
						SourceBranch:     []byte("ours"),
						TargetBranch:     []byte("theirs"),
						User:             user,
					},
				},
			}))
			require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
				ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
					FilesJson: filesJSON,
				},
			}))

			response, err := stream.CloseAndRecv()
			require.NoError(t, err)
			require.Empty(t, response.GetResolutionError())

			oursFile := gittest.Exec(t, cfg, "-C", repoPath, "cat-file", "-p", "refs/heads/ours:file.txt")
			require.Equal(t, []byte(tc.expectedContents), oursFile)
		})
	}
}

func TestResolveConflictsNonOIDRequests(t *testing.T) {
	hookManager := hook.NewMockManager(t, hook.NopPreReceive, hook.NopPostReceive, hook.NopUpdate, hook.NopReferenceTransaction)
	cfg, repoProto, _, client := SetupConflictsService(t, true, hookManager)

	ctx, cancel := testhelper.Context()
	defer cancel()
	ctx = testhelper.MergeOutgoingMetadata(ctx, testhelper.GitalyServersMetadataFromCfg(t, cfg))

	stream, err := client.ResolveConflicts(ctx)
	require.NoError(t, err)

	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
			Header: &gitalypb.ResolveConflictsRequestHeader{
				Repository:       repoProto,
				TargetRepository: repoProto,
				CommitMessage:    []byte(conflictResolutionCommitMessage),
				OurCommitOid:     "conflict-resolvable",
				TheirCommitOid:   "conflict-start",
				SourceBranch:     []byte("conflict-resolvable"),
				TargetBranch:     []byte("conflict-start"),
				User:             user,
			},
		},
	}))

	filesJSON, err := json.Marshal(files)
	require.NoError(t, err)
	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
			FilesJson: filesJSON,
		},
	}))

	_, err = stream.CloseAndRecv()
	testassert.GrpcEqualErr(t, status.Errorf(codes.Unknown, "Rugged::InvalidError: unable to parse OID - contains invalid characters"), err)
}

func TestResolveConflictsIdenticalContent(t *testing.T) {
	hookManager := hook.NewMockManager(t, hook.NopPreReceive, hook.NopPostReceive, hook.NopUpdate, hook.NopReferenceTransaction)
	cfg, repoProto, repoPath, client := SetupConflictsService(t, true, hookManager)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	ctx, cancel := testhelper.Context()
	defer cancel()

	sourceBranch := "conflict-resolvable"
	sourceOID, err := repo.ResolveRevision(ctx, git.Revision(sourceBranch))
	require.NoError(t, err)

	targetBranch := "conflict-start"
	targetOID, err := repo.ResolveRevision(ctx, git.Revision(targetBranch))
	require.NoError(t, err)

	tempDir := testhelper.TempDir(t)

	var conflictingPaths []string
	for _, rev := range []string{
		sourceOID.String(),
		"6907208d755b60ebeacb2e9dfea74c92c3449a1f",
		targetOID.String(),
	} {
		contents := gittest.Exec(t, cfg, "-C", repoPath, "cat-file", "-p", rev+":files/ruby/popen.rb")
		path := filepath.Join(tempDir, rev)
		require.NoError(t, ioutil.WriteFile(path, contents, 0644))
		conflictingPaths = append(conflictingPaths, path)
	}

	var conflictContents bytes.Buffer
	err = repo.ExecAndWait(ctx, git.SubCmd{
		Name: "merge-file",
		Flags: []git.Option{
			git.Flag{Name: "--quiet"},
			git.Flag{Name: "--stdout"},
			// We pass `-L` three times for each of the conflicting files.
			git.ValueFlag{Name: "-L", Value: "files/ruby/popen.rb"},
			git.ValueFlag{Name: "-L", Value: "files/ruby/popen.rb"},
			git.ValueFlag{Name: "-L", Value: "files/ruby/popen.rb"},
		},
		Args: conflictingPaths,
	}, git.WithStdout(&conflictContents))

	// The merge will result in a merge conflict and thus cause the command to fail.
	require.Error(t, err)
	require.Contains(t, conflictContents.String(), "<<<<<<")

	filesJSON, err := json.Marshal([]map[string]interface{}{
		{
			"old_path": "files/ruby/popen.rb",
			"new_path": "files/ruby/popen.rb",
			"content":  conflictContents.String(),
		},
		{
			"old_path": "files/ruby/regex.rb",
			"new_path": "files/ruby/regex.rb",
			"sections": map[string]string{
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_9_9":   "head",
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_21_21": "origin",
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_49_49": "origin",
			},
		},
	})
	require.NoError(t, err)

	ctx = testhelper.MergeOutgoingMetadata(ctx, testhelper.GitalyServersMetadataFromCfg(t, cfg))
	stream, err := client.ResolveConflicts(ctx)
	require.NoError(t, err)

	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
			Header: &gitalypb.ResolveConflictsRequestHeader{
				Repository:       repoProto,
				TargetRepository: repoProto,
				CommitMessage:    []byte(conflictResolutionCommitMessage),
				OurCommitOid:     sourceOID.String(),
				TheirCommitOid:   targetOID.String(),
				SourceBranch:     []byte(sourceBranch),
				TargetBranch:     []byte(targetBranch),
				User:             user,
			},
		},
	}))
	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
			FilesJson: filesJSON,
		},
	}))

	response, err := stream.CloseAndRecv()
	require.NoError(t, err)
	testassert.ProtoEqual(t, &gitalypb.ResolveConflictsResponse{
		ResolutionError: "Resolved content has no changes for file files/ruby/popen.rb",
	}, response)
}

func TestResolveConflictsStableID(t *testing.T) {
	hookManager := hook.NewMockManager(t, hook.NopPreReceive, hook.NopPostReceive, hook.NopUpdate, hook.NopReferenceTransaction)
	cfg, repoProto, _, client := SetupConflictsService(t, true, hookManager)

	repo := localrepo.NewTestRepo(t, cfg, repoProto)

	ctx, cancel := testhelper.Context()
	defer cancel()
	md := testhelper.GitalyServersMetadataFromCfg(t, cfg)
	ctx = testhelper.MergeOutgoingMetadata(ctx, md)

	stream, err := client.ResolveConflicts(ctx)
	require.NoError(t, err)

	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
			Header: &gitalypb.ResolveConflictsRequestHeader{
				Repository:       repoProto,
				TargetRepository: repoProto,
				CommitMessage:    []byte(conflictResolutionCommitMessage),
				OurCommitOid:     "1450cd639e0bc6721eb02800169e464f212cde06",
				TheirCommitOid:   "824be604a34828eb682305f0d963056cfac87b2d",
				SourceBranch:     []byte("conflict-resolvable"),
				TargetBranch:     []byte("conflict-start"),
				User:             user,
				Timestamp:        &timestamppb.Timestamp{Seconds: 12345},
			},
		},
	}))

	filesJSON, err := json.Marshal(files)
	require.NoError(t, err)
	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
			FilesJson: filesJSON,
		},
	}))

	response, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.Empty(t, response.GetResolutionError())

	resolvedCommit, err := repo.ReadCommit(ctx, git.Revision("conflict-resolvable"))
	require.NoError(t, err)
	require.Equal(t, &gitalypb.GitCommit{
		Id:     "a5ad028fd739d7a054b07c293e77c5b7aecc2435",
		TreeId: "febd97e4a09e71355a513d7e0b0b3808e2dabd28",
		ParentIds: []string{
			"1450cd639e0bc6721eb02800169e464f212cde06",
			"824be604a34828eb682305f0d963056cfac87b2d",
		},
		Subject:  []byte(conflictResolutionCommitMessage),
		Body:     []byte(conflictResolutionCommitMessage),
		BodySize: 15,
		Author: &gitalypb.CommitAuthor{
			Name:     user.Name,
			Email:    user.Email,
			Date:     &timestamppb.Timestamp{Seconds: 12345},
			Timezone: []byte("+0000"),
		},
		Committer: &gitalypb.CommitAuthor{
			Name:     user.Name,
			Email:    user.Email,
			Date:     &timestamppb.Timestamp{Seconds: 12345},
			Timezone: []byte("+0000"),
		},
	}, resolvedCommit)
}

func TestFailedResolveConflictsRequestDueToResolutionError(t *testing.T) {
	hookManager := hook.NewMockManager(t, hook.NopPreReceive, hook.NopPostReceive, hook.NopUpdate, hook.NopReferenceTransaction)
	cfg, repo, _, client := SetupConflictsService(t, true, hookManager)

	ctx, cancel := testhelper.Context()
	defer cancel()
	mdGS := testhelper.GitalyServersMetadataFromCfg(t, cfg)
	mdFF, _ := metadata.FromOutgoingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Join(mdGS, mdFF))

	files := []map[string]interface{}{
		{
			"old_path": "files/ruby/popen.rb",
			"new_path": "files/ruby/popen.rb",
			"content":  "",
		},
		{
			"old_path": "files/ruby/regex.rb",
			"new_path": "files/ruby/regex.rb",
			"sections": map[string]string{
				"6eb14e00385d2fb284765eb1cd8d420d33d63fc9_9_9": "head",
			},
		},
	}
	filesJSON, err := json.Marshal(files)
	require.NoError(t, err)

	headerRequest := &gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
			Header: &gitalypb.ResolveConflictsRequestHeader{
				Repository:       repo,
				TargetRepository: repo,
				CommitMessage:    []byte(conflictResolutionCommitMessage),
				OurCommitOid:     "1450cd639e0bc6721eb02800169e464f212cde06",
				TheirCommitOid:   "824be604a34828eb682305f0d963056cfac87b2d",
				SourceBranch:     []byte("conflict-resolvable"),
				TargetBranch:     []byte("conflict-start"),
				User:             user,
			},
		},
	}
	filesRequest := &gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
			FilesJson: filesJSON,
		},
	}

	stream, err := client.ResolveConflicts(ctx)
	require.NoError(t, err)
	require.NoError(t, stream.Send(headerRequest))
	require.NoError(t, stream.Send(filesRequest))

	r, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.Equal(t, r.GetResolutionError(), "Missing resolution for section ID: 6eb14e00385d2fb284765eb1cd8d420d33d63fc9_21_21")
}

func TestFailedResolveConflictsRequestDueToValidation(t *testing.T) {
	hookManager := hook.NewMockManager(t, hook.NopPreReceive, hook.NopPostReceive, hook.NopUpdate, hook.NopReferenceTransaction)
	cfg, repo, _, client := SetupConflictsService(t, true, hookManager)

	ctx, cancel := testhelper.Context()
	defer cancel()

	mdGS := testhelper.GitalyServersMetadataFromCfg(t, cfg)
	ourCommitOid := "1450cd639e0bc6721eb02800169e464f212cde06"
	theirCommitOid := "824be604a34828eb682305f0d963056cfac87b2d"
	commitMsg := []byte(conflictResolutionCommitMessage)
	sourceBranch := []byte("conflict-resolvable")
	targetBranch := []byte("conflict-start")

	testCases := []struct {
		desc         string
		header       *gitalypb.ResolveConflictsRequestHeader
		expectedCode codes.Code
		expectedErr  string
	}{
		{
			desc: "empty user",
			header: &gitalypb.ResolveConflictsRequestHeader{
				User:             nil,
				Repository:       repo,
				OurCommitOid:     ourCommitOid,
				TargetRepository: repo,
				TheirCommitOid:   theirCommitOid,
				CommitMessage:    commitMsg,
				SourceBranch:     sourceBranch,
				TargetBranch:     targetBranch,
			},
			expectedCode: codes.InvalidArgument,
			expectedErr:  "ResolveConflicts: empty User",
		},
		{
			desc: "empty repo",
			header: &gitalypb.ResolveConflictsRequestHeader{
				User:             user,
				Repository:       nil,
				OurCommitOid:     ourCommitOid,
				TargetRepository: repo,
				TheirCommitOid:   theirCommitOid,
				CommitMessage:    commitMsg,
				SourceBranch:     sourceBranch,
				TargetBranch:     targetBranch,
			},
			expectedCode: codes.InvalidArgument,
			// Praefect checks for an empty repository, too, but will raise a different
			// error message. Luckily, both Gitaly's and Praefect's error messages
			// contain "empty Repository".
			expectedErr: "empty Repository",
		},
		{
			desc: "empty target repo",
			header: &gitalypb.ResolveConflictsRequestHeader{
				User:             user,
				Repository:       repo,
				OurCommitOid:     ourCommitOid,
				TargetRepository: nil,
				TheirCommitOid:   theirCommitOid,
				CommitMessage:    commitMsg,
				SourceBranch:     sourceBranch,
				TargetBranch:     targetBranch,
			},
			expectedCode: codes.InvalidArgument,
			expectedErr:  "ResolveConflicts: empty TargetRepository",
		},
		{
			desc: "empty OurCommitId repo",
			header: &gitalypb.ResolveConflictsRequestHeader{
				User:             user,
				Repository:       repo,
				OurCommitOid:     "",
				TargetRepository: repo,
				TheirCommitOid:   theirCommitOid,
				CommitMessage:    commitMsg,
				SourceBranch:     sourceBranch,
				TargetBranch:     targetBranch,
			},
			expectedCode: codes.InvalidArgument,
			expectedErr:  "ResolveConflicts: empty OurCommitOid",
		},
		{
			desc: "empty TheirCommitId repo",
			header: &gitalypb.ResolveConflictsRequestHeader{
				User:             user,
				Repository:       repo,
				OurCommitOid:     ourCommitOid,
				TargetRepository: repo,
				TheirCommitOid:   "",
				CommitMessage:    commitMsg,
				SourceBranch:     sourceBranch,
				TargetBranch:     targetBranch,
			},
			expectedCode: codes.InvalidArgument,
			expectedErr:  "ResolveConflicts: empty TheirCommitOid",
		},
		{
			desc: "empty CommitMessage repo",
			header: &gitalypb.ResolveConflictsRequestHeader{
				User:             user,
				Repository:       repo,
				OurCommitOid:     ourCommitOid,
				TargetRepository: repo,
				TheirCommitOid:   theirCommitOid,
				CommitMessage:    nil,
				SourceBranch:     sourceBranch,
				TargetBranch:     targetBranch,
			},
			expectedCode: codes.InvalidArgument,
			expectedErr:  "ResolveConflicts: empty CommitMessage",
		},
		{
			desc: "empty SourceBranch repo",
			header: &gitalypb.ResolveConflictsRequestHeader{
				User:             user,
				Repository:       repo,
				OurCommitOid:     ourCommitOid,
				TargetRepository: repo,
				TheirCommitOid:   theirCommitOid,
				CommitMessage:    commitMsg,
				SourceBranch:     nil,
				TargetBranch:     targetBranch,
			},
			expectedCode: codes.InvalidArgument,
			expectedErr:  "ResolveConflicts: empty SourceBranch",
		},
		{
			desc: "empty TargetBranch repo",
			header: &gitalypb.ResolveConflictsRequestHeader{
				User:             user,
				Repository:       repo,
				OurCommitOid:     ourCommitOid,
				TargetRepository: repo,
				TheirCommitOid:   theirCommitOid,
				CommitMessage:    commitMsg,
				SourceBranch:     sourceBranch,
				TargetBranch:     nil,
			},
			expectedCode: codes.InvalidArgument,
			expectedErr:  "ResolveConflicts: empty TargetBranch",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			mdFF, _ := metadata.FromOutgoingContext(ctx)
			ctx = metadata.NewOutgoingContext(ctx, metadata.Join(mdGS, mdFF))

			stream, err := client.ResolveConflicts(ctx)
			require.NoError(t, err)

			headerRequest := &gitalypb.ResolveConflictsRequest{
				ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
					Header: testCase.header,
				},
			}
			require.NoError(t, stream.Send(headerRequest))

			_, err = stream.CloseAndRecv()
			testhelper.RequireGrpcError(t, err, testCase.expectedCode)
			require.Contains(t, err.Error(), testCase.expectedErr)
		})
	}
}

func TestResolveConflictsQuarantine(t *testing.T) {
	cfg, _, _, client := SetupConflictsService(t, true, nil)

	ctx, cancel := testhelper.Context()
	defer cancel()

	testhelper.BuildGitalySSH(t, cfg)
	testhelper.BuildGitalyHooks(t, cfg)

	sourceRepoProto, sourceRepoPath := gittest.CloneRepo(t, cfg, cfg.Storages[0])
	sourceBlobOID := gittest.WriteBlob(t, cfg, sourceRepoPath, []byte("contents-1\n"))
	sourceCommitOID := gittest.WriteCommit(t, cfg, sourceRepoPath,
		gittest.WithTreeEntries(gittest.TreeEntry{
			Path: "file.txt", OID: sourceBlobOID, Mode: "100644",
		}),
	)
	gittest.Exec(t, cfg, "-C", sourceRepoPath, "update-ref", "refs/heads/source", sourceCommitOID.String())

	// We set up a custom "update" hook which simply prints the new commit to stdout and then
	// exits with an error. Like this, we can both assert that the hook can see the quarantined
	// tag, and it allows us to fail the RPC before we migrate quarantined objects.
	script := fmt.Sprintf("#!/bin/sh\necho $3 && %s cat-file -p $3^{commit} && exit 1", cfg.Git.BinPath)
	gittest.WriteCustomHook(t, sourceRepoPath, "update", []byte(script))

	targetRepoProto, targetRepoPath := gittest.CloneRepo(t, cfg, cfg.Storages[0])
	targetBlobOID := gittest.WriteBlob(t, cfg, targetRepoPath, []byte("contents-2\n"))
	targetCommitOID := gittest.WriteCommit(t, cfg, targetRepoPath,
		gittest.WithTreeEntries(gittest.TreeEntry{
			Path: "file.txt", OID: targetBlobOID, Mode: "100644",
		}),
	)
	gittest.Exec(t, cfg, "-C", targetRepoPath, "update-ref", "refs/heads/target", targetCommitOID.String())

	ctx = testhelper.MergeOutgoingMetadata(ctx, testhelper.GitalyServersMetadataFromCfg(t, cfg))

	stream, err := client.ResolveConflicts(ctx)
	require.NoError(t, err)

	filesJSON, err := json.Marshal([]map[string]interface{}{
		{
			"old_path": "file.txt",
			"new_path": "file.txt",
			"sections": map[string]string{
				"5436437fa01a7d3e41d46741da54b451446774ca_1_1": "origin",
			},
		},
	})
	require.NoError(t, err)

	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_Header{
			Header: &gitalypb.ResolveConflictsRequestHeader{
				Repository:       sourceRepoProto,
				TargetRepository: targetRepoProto,
				CommitMessage:    []byte(conflictResolutionCommitMessage),
				OurCommitOid:     sourceCommitOID.String(),
				TheirCommitOid:   targetCommitOID.String(),
				SourceBranch:     []byte("source"),
				TargetBranch:     []byte("target"),
				User:             user,
				Timestamp:        &timestamppb.Timestamp{Seconds: 12345},
			},
		},
	}))
	require.NoError(t, stream.Send(&gitalypb.ResolveConflictsRequest{
		ResolveConflictsRequestPayload: &gitalypb.ResolveConflictsRequest_FilesJson{
			FilesJson: filesJSON,
		},
	}))

	response, err := stream.CloseAndRecv()
	require.EqualError(t, err, `rpc error: code = Unknown desc = af339cb882d1e3cf8d6751651e58bbaff0265d6e
tree 89fad81bbfa38070b90ca8f4c404625bf0999013
parent 29449b1d52cd77fd060a083a1de691bbaf12d8af
parent 26dac52be85c92742b2c0c19eb7303de9feccb63
author John Doe <johndoe@gitlab.com> 12345 +0000
committer John Doe <johndoe@gitlab.com> 12345 +0000

Solve conflicts`)
	require.Empty(t, response.GetResolutionError())

	// The file shouldn't have been updated and is thus expected to still have the same old
	// contents.
	require.Equal(t, []byte("contents-1\n"), gittest.Exec(t, cfg, "-C", sourceRepoPath, "cat-file", "-p", "refs/heads/source:file.txt"))

	// In case we use an object quarantine directory, the tag should not exist in the target
	// repository because the RPC failed to update the revision.
	exists, err := localrepo.NewTestRepo(t, cfg, sourceRepoProto).HasRevision(ctx, "af339cb882d1e3cf8d6751651e58bbaff0265d6e^{commit}")
	require.NoError(t, err)
	require.False(t, exists, "object should have not been migrated")
}
