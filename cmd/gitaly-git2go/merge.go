// +build static,system_libgit2

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	git "github.com/libgit2/git2go/v31"
	"gitlab.com/gitlab-org/gitaly/v14/cmd/gitaly-git2go/conflicts"
	"gitlab.com/gitlab-org/gitaly/v14/cmd/gitaly-git2go/git2goutil"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git2go"
)

type mergeSubcommand struct {
	request string
}

func (cmd *mergeSubcommand) Flags() *flag.FlagSet {
	flags := flag.NewFlagSet("merge", flag.ExitOnError)
	flags.StringVar(&cmd.request, "request", "", "git2go.MergeCommand")
	return flags
}

func (cmd *mergeSubcommand) Run(context.Context, io.Reader, io.Writer) error {
	request, err := git2go.MergeCommandFromSerialized(cmd.request)
	if err != nil {
		return err
	}

	if request.AuthorDate.IsZero() {
		request.AuthorDate = time.Now()
	}

	repo, err := git2goutil.OpenRepository(request.Repository)
	if err != nil {
		return fmt.Errorf("could not open repository: %w", err)
	}
	defer repo.Free()

	ours, err := lookupCommit(repo, request.Ours)
	if err != nil {
		return fmt.Errorf("ours commit lookup: %w", err)
	}

	theirs, err := lookupCommit(repo, request.Theirs)
	if err != nil {
		return fmt.Errorf("theirs commit lookup: %w", err)
	}

	mergeOpts, err := git.DefaultMergeOptions()
	if err != nil {
		return fmt.Errorf("could not create merge options: %w", err)
	}
	mergeOpts.RecursionLimit = git2go.MergeRecursionLimit

	index, err := repo.MergeCommits(ours, theirs, &mergeOpts)
	if err != nil {
		return fmt.Errorf("could not merge commits: %w", err)
	}
	defer index.Free()

	if index.HasConflicts() {
		if !request.AllowConflicts {
			return errors.New("could not auto-merge due to conflicts")
		}

		if err := resolveConflicts(repo, index); err != nil {
			return fmt.Errorf("could not resolve conflicts: %w", err)
		}
	}

	tree, err := index.WriteTreeTo(repo)
	if err != nil {
		return fmt.Errorf("could not write tree: %w", err)
	}

	committer := git.Signature(git2go.NewSignature(request.AuthorName, request.AuthorMail, request.AuthorDate))
	commit, err := repo.CreateCommitFromIds("", &committer, &committer, request.Message, tree, ours.Id(), theirs.Id())
	if err != nil {
		return fmt.Errorf("could not create merge commit: %w", err)
	}

	response := git2go.MergeResult{
		CommitID: commit.String(),
	}

	if err := response.SerializeTo(os.Stdout); err != nil {
		return err
	}

	return nil
}

func resolveConflicts(repo *git.Repository, index *git.Index) error {
	// We need to get all conflicts up front as resolving conflicts as we
	// iterate breaks the iterator.
	indexConflicts, err := getConflicts(index)
	if err != nil {
		return err
	}

	for _, conflict := range indexConflicts {
		if isConflictMergeable(conflict) {
			merge, err := conflicts.Merge(repo, conflict)
			if err != nil {
				return err
			}

			mergedBlob, err := repo.CreateBlobFromBuffer(merge.Contents)
			if err != nil {
				return err
			}

			mergedIndexEntry := git.IndexEntry{
				Path: merge.Path,
				Mode: git.Filemode(merge.Mode),
				Id:   mergedBlob,
			}

			if err := index.Add(&mergedIndexEntry); err != nil {
				return err
			}

			if err := index.RemoveConflict(merge.Path); err != nil {
				return err
			}
		} else {
			if conflict.Their != nil {
				// If a conflict has `Their` present, we add it back to the index
				// as we want those changes to be part of the merge.
				if err := index.Add(conflict.Their); err != nil {
					return err
				}

				if err := index.RemoveConflict(conflict.Their.Path); err != nil {
					return err
				}
			} else if conflict.Our != nil {
				// If a conflict has `Our` present, remove its conflict as we
				// don't want to include those changes.
				if err := index.RemoveConflict(conflict.Our.Path); err != nil {
					return err
				}
			} else {
				// If conflict has no `Their` and `Our`, remove the conflict to
				// mark it as resolved.
				if err := index.RemoveConflict(conflict.Ancestor.Path); err != nil {
					return err
				}
			}
		}
	}

	if index.HasConflicts() {
		return errors.New("index still has conflicts")
	}

	return nil
}

func isConflictMergeable(conflict git.IndexConflict) bool {
	conflictIndexEntriesCount := 0

	if conflict.Their != nil {
		conflictIndexEntriesCount++
	}

	if conflict.Our != nil {
		conflictIndexEntriesCount++
	}

	if conflict.Ancestor != nil {
		conflictIndexEntriesCount++
	}

	return conflictIndexEntriesCount >= 2
}

func getConflicts(index *git.Index) ([]git.IndexConflict, error) {
	var conflicts []git.IndexConflict

	iterator, err := index.ConflictIterator()
	if err != nil {
		return nil, err
	}
	defer iterator.Free()

	for {
		conflict, err := iterator.Next()
		if err != nil {
			if git.IsErrorCode(err, git.ErrIterOver) {
				break
			}
			return nil, err
		}

		conflicts = append(conflicts, conflict)
	}

	return conflicts, nil
}
