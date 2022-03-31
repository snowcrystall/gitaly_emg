package commit

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"gitlab.com/gitlab-org/gitaly/v14/internal/command"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/trailerparser"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/chunk"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func (s *server) FindCommits(req *gitalypb.FindCommitsRequest, stream gitalypb.CommitService_FindCommitsServer) error {
	ctx := stream.Context()

	if err := git.ValidateRevisionAllowEmpty(req.Revision); err != nil {
		return helper.ErrInvalidArgument(err)
	}

	repo := s.localrepo(req.GetRepository())

	// Use Gitaly's default branch lookup function because that is already
	// migrated.
	if revision := req.Revision; len(revision) == 0 && !req.GetAll() {
		var err error
		req.Revision, err = defaultBranchName(ctx, repo)
		if err != nil {
			return helper.ErrInternal(fmt.Errorf("defaultBranchName: %v", err))
		}
	}
	// Clients might send empty paths. That is an error
	for _, path := range req.Paths {
		if len(path) == 0 {
			return helper.ErrInvalidArgument(errors.New("path is empty string"))
		}
	}

	if err := s.findCommits(ctx, req, stream); err != nil {
		return helper.ErrInternal(err)
	}

	return nil
}

func (s *server) findCommits(ctx context.Context, req *gitalypb.FindCommitsRequest, stream gitalypb.CommitService_FindCommitsServer) error {
	opts := git.ConvertGlobalOptions(req.GetGlobalOptions())
	repo := s.localrepo(req.GetRepository())

	logCmd, err := repo.Exec(ctx, getLogCommandSubCmd(req), opts...)
	if err != nil {
		return fmt.Errorf("error when creating git log command: %v", err)
	}

	batch, err := s.catfileCache.BatchProcess(ctx, repo)
	if err != nil {
		return fmt.Errorf("creating catfile: %v", err)
	}

	getCommits := NewGetCommits(logCmd, batch)

	if calculateOffsetManually(req) {
		if err := getCommits.Offset(int(req.GetOffset())); err != nil {
			// If we're at EOF, then it means that the offset has been greater than the
			// number of available commits. We do not treat this as an error, but
			// instead just return EOF ourselves.
			if errors.Is(err, io.EOF) {
				return nil
			}

			return fmt.Errorf("skipping to offset %d: %w", req.GetOffset(), err)
		}
	}

	if err := streamCommits(getCommits, stream, req.GetTrailers()); err != nil {
		return fmt.Errorf("error streaming commits: %v", err)
	}
	return nil
}

func calculateOffsetManually(req *gitalypb.FindCommitsRequest) bool {
	return req.GetFollow() && req.GetOffset() > 0
}

// GetCommits wraps a git log command that can be interated on to get individual commit objects
type GetCommits struct {
	scanner *bufio.Scanner
	batch   catfile.Batch
}

// NewGetCommits returns a new GetCommits object
func NewGetCommits(cmd *command.Command, batch catfile.Batch) *GetCommits {
	return &GetCommits{
		scanner: bufio.NewScanner(cmd),
		batch:   batch,
	}
}

// Scan indicates whether or not there are more commits to return
func (g *GetCommits) Scan() bool {
	return g.scanner.Scan()
}

// Err returns the first non EOF error
func (g *GetCommits) Err() error {
	return g.scanner.Err()
}

// Offset skips over a number of commits
func (g *GetCommits) Offset(offset int) error {
	for i := 0; i < offset; i++ {
		if !g.Scan() {
			err := g.Err()
			if err == nil {
				err = io.EOF
			}

			return fmt.Errorf("skipping commit: %w", err)
		}
	}
	return nil
}

// Commit returns the current commit
func (g *GetCommits) Commit(ctx context.Context, trailers bool) (*gitalypb.GitCommit, error) {
	logOutput := strings.TrimSpace(g.scanner.Text())
	var revAndTrailers []string
	var revision string

	if trailers {
		revAndTrailers = strings.SplitN(logOutput, "\000", 2)
		revision = revAndTrailers[0]
	} else {
		revision = logOutput
	}
	commit, err := catfile.GetCommit(ctx, g.batch, git.Revision(revision))
	if err != nil {
		return nil, fmt.Errorf("cat-file get commit %q: %v", revision, err)
	}

	if trailers && len(revAndTrailers) == 2 {
		commit.Trailers = trailerparser.Parse([]byte(revAndTrailers[1]))
	}

	return commit, nil
}

func streamCommits(getCommits *GetCommits, stream gitalypb.CommitService_FindCommitsServer, trailers bool) error {
	ctx := stream.Context()

	chunker := chunk.New(&commitsSender{
		send: func(commits []*gitalypb.GitCommit) error {
			return stream.Send(&gitalypb.FindCommitsResponse{
				Commits: commits,
			})
		},
	})

	for getCommits.Scan() {
		commit, err := getCommits.Commit(ctx, trailers)
		if err != nil {
			return err
		}

		if err := chunker.Send(commit); err != nil {
			return err
		}
	}
	if getCommits.Err() != nil {
		return fmt.Errorf("get commits: %v", getCommits.Err())
	}

	return chunker.Flush()
}

func getLogCommandSubCmd(req *gitalypb.FindCommitsRequest) git.SubCmd {
	logFormatOption := "--format=%H"
	if req.GetTrailers() {
		logFormatOption += "%x00%(trailers:unfold,separator=%x00)"
	}
	subCmd := git.SubCmd{Name: "log", Flags: []git.Option{git.Flag{Name: logFormatOption}}}

	//  We will perform the offset in Go because --follow doesn't play well with --skip.
	//  See: https://gitlab.com/gitlab-org/gitlab-ce/issues/3574#note_3040520
	if req.GetOffset() > 0 && !calculateOffsetManually(req) {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: fmt.Sprintf("--skip=%d", req.GetOffset())})
	}
	limit := req.GetLimit()
	if calculateOffsetManually(req) {
		limit += req.GetOffset()
	}
	subCmd.Flags = append(subCmd.Flags, git.Flag{Name: fmt.Sprintf("--max-count=%d", limit)})

	if req.GetFollow() && len(req.GetPaths()) > 0 {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: "--follow"})
	}
	if req.GetAuthor() != nil {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: fmt.Sprintf("--author=%s", string(req.GetAuthor()))})
	}
	if req.GetSkipMerges() {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: "--no-merges"})
	}
	if req.GetBefore() != nil {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: fmt.Sprintf("--before=%s", req.GetBefore().String())})
	}
	if req.GetAfter() != nil {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: fmt.Sprintf("--after=%s", req.GetAfter().String())})
	}
	if req.GetAll() {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: "--all"}, git.Flag{Name: "--reverse"})
	}
	if req.GetRevision() != nil {
		subCmd.Args = []string{string(req.GetRevision())}
	}
	if req.GetFirstParent() {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: "--first-parent"})
	}
	if len(req.GetPaths()) > 0 {
		for _, path := range req.GetPaths() {
			subCmd.PostSepArgs = append(subCmd.PostSepArgs, string(path))
		}
	}
	if req.GetOrder() == gitalypb.FindCommitsRequest_TOPO {
		subCmd.Flags = append(subCmd.Flags, git.Flag{Name: "--topo-order"})
	}

	return subCmd
}
