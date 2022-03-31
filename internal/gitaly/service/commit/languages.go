package commit

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/service/ref"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper/text"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

var errAmbigRef = errors.New("ambiguous reference")

func (s *server) CommitLanguages(ctx context.Context, req *gitalypb.CommitLanguagesRequest) (*gitalypb.CommitLanguagesResponse, error) {
	if err := git.ValidateRevisionAllowEmpty(req.Revision); err != nil {
		return nil, helper.ErrInvalidArgument(err)
	}

	repo := s.localrepo(req.GetRepository())

	revision := string(req.Revision)
	if revision == "" {
		defaultBranch, err := ref.DefaultBranchName(ctx, repo)
		if err != nil {
			return nil, err
		}
		revision = string(defaultBranch)
	}

	commitID, err := s.lookupRevision(ctx, repo, revision)
	if err != nil {
		return nil, helper.ErrInternalf("looking up revision: %w", err)
	}

	repoPath, err := repo.Path()
	if err != nil {
		return nil, helper.ErrInternalf("repository path: %w", err)
	}
	stats, err := s.linguist.Stats(ctx, s.cfg, repoPath, commitID)
	if err != nil {
		return nil, helper.ErrInternalf("language stats: %w", err)
	}

	resp := &gitalypb.CommitLanguagesResponse{}
	if len(stats) == 0 {
		return resp, nil
	}

	total := uint64(0)
	for _, count := range stats {
		total += count
	}

	if total == 0 {
		return nil, helper.ErrInternalf("linguist stats added up to zero: %v", stats)
	}

	for lang, count := range stats {
		l := &gitalypb.CommitLanguagesResponse_Language{
			Name:  lang,
			Share: float32(100*count) / float32(total),
			Color: s.linguist.Color(lang),
			Bytes: stats[lang],
		}
		resp.Languages = append(resp.Languages, l)
	}

	sort.Sort(languageSorter(resp.Languages))

	return resp, nil
}

type languageSorter []*gitalypb.CommitLanguagesResponse_Language

func (ls languageSorter) Len() int           { return len(ls) }
func (ls languageSorter) Swap(i, j int)      { ls[i], ls[j] = ls[j], ls[i] }
func (ls languageSorter) Less(i, j int) bool { return ls[i].Share > ls[j].Share }

func (s *server) lookupRevision(ctx context.Context, repo git.RepositoryExecutor, revision string) (string, error) {
	rev, err := s.checkRevision(ctx, repo, revision)
	if err != nil {
		switch err {
		case errAmbigRef:
			fullRev, err := s.disambiguateRevision(ctx, repo, revision)
			if err != nil {
				return "", err
			}

			rev, err = s.checkRevision(ctx, repo, fullRev)
			if err != nil {
				return "", err
			}
		default:
			return "", err
		}
	}

	return rev, nil
}

func (s *server) checkRevision(ctx context.Context, repo git.RepositoryExecutor, revision string) (string, error) {
	var stdout, stderr bytes.Buffer

	revParse, err := repo.Exec(ctx,
		git.SubCmd{Name: "rev-parse", Args: []string{revision}},
		git.WithStdout(&stdout),
		git.WithStderr(&stderr),
	)

	if err != nil {
		return "", err
	}

	if err = revParse.Wait(); err != nil {
		errMsg := strings.Split(stderr.String(), "\n")[0]
		return "", fmt.Errorf("%v: %v", err, errMsg)
	}

	if strings.HasSuffix(stderr.String(), "refname '"+revision+"' is ambiguous.\n") {
		return "", errAmbigRef
	}

	return text.ChompBytes(stdout.Bytes()), nil
}

func (s *server) disambiguateRevision(ctx context.Context, repo git.RepositoryExecutor, revision string) (string, error) {
	cmd, err := repo.Exec(ctx, git.SubCmd{
		Name:  "for-each-ref",
		Flags: []git.Option{git.Flag{Name: "--format=%(refname)"}},
		Args:  []string{"**/" + revision},
	})

	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(cmd)
	for scanner.Scan() {
		refName := scanner.Text()

		if strings.HasPrefix(refName, "refs/heads") {
			return refName, nil
		}
	}

	return "", fmt.Errorf("branch %v not found", revision)
}
