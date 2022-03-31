package repository

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	log "github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/housekeeping"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/stats"
	"gitlab.com/gitlab-org/gitaly/v14/internal/helper"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/status"
)

func (s *server) GarbageCollect(ctx context.Context, in *gitalypb.GarbageCollectRequest) (*gitalypb.GarbageCollectResponse, error) {
	ctxlogger := ctxlogrus.Extract(ctx)
	ctxlogger.WithFields(log.Fields{
		"WriteBitmaps": in.GetCreateBitmap(),
	}).Debug("GarbageCollect")

	repo := s.localrepo(in.GetRepository())

	if err := s.cleanupRepo(ctx, repo); err != nil {
		return nil, err
	}

	if err := s.cleanupKeepArounds(ctx, repo); err != nil {
		return nil, err
	}

	// Perform housekeeping to cleanup stale lockfiles that may block GC
	if err := housekeeping.Perform(ctx, repo); err != nil {
		ctxlogger.WithError(err).Warn("Pre gc housekeeping failed")
	}

	if err := s.gc(ctx, in); err != nil {
		return nil, err
	}

	if err := s.writeCommitGraph(ctx, repo, gitalypb.WriteCommitGraphRequest_SizeMultiple); err != nil {
		return nil, err
	}

	stats.LogObjectsInfo(ctx, s.gitCmdFactory, repo)

	return &gitalypb.GarbageCollectResponse{}, nil
}

func (s *server) gc(ctx context.Context, in *gitalypb.GarbageCollectRequest) error {
	config := append(repackConfig(ctx, in.CreateBitmap), git.ConfigPair{Key: "gc.writeCommitGraph", Value: "false"})

	var flags []git.Option
	if in.Prune {
		flags = append(flags, git.Flag{Name: "--prune=30.minutes.ago"})
	}

	cmd, err := s.gitCmdFactory.New(ctx, in.GetRepository(),
		git.SubCmd{Name: "gc", Flags: flags},
		git.WithConfig(config...),
	)

	if err != nil {
		if git.IsInvalidArgErr(err) {
			return helper.ErrInvalidArgumentf("GarbageCollect: gitCommand: %v", err)
		}

		if _, ok := status.FromError(err); ok {
			return err
		}
		return helper.ErrInternal(fmt.Errorf("GarbageCollect: gitCommand: %v", err))
	}

	if err := cmd.Wait(); err != nil {
		return helper.ErrInternal(fmt.Errorf("GarbageCollect: cmd wait: %v", err))
	}

	return nil
}

func (s *server) cleanupKeepArounds(ctx context.Context, repo *localrepo.Repo) error {
	repoPath, err := repo.Path()
	if err != nil {
		return nil
	}

	batch, err := s.catfileCache.BatchProcess(ctx, repo)
	if err != nil {
		return nil
	}

	keepAroundsPrefix := "refs/keep-around"
	keepAroundsPath := filepath.Join(repoPath, keepAroundsPrefix)

	refInfos, err := ioutil.ReadDir(keepAroundsPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, info := range refInfos {
		if info.IsDir() {
			continue
		}

		refName := fmt.Sprintf("%s/%s", keepAroundsPrefix, info.Name())
		path := filepath.Join(repoPath, keepAroundsPrefix, info.Name())

		if err = checkRef(ctx, batch, refName, info); err == nil {
			continue
		}

		if err := s.fixRef(ctx, repo, batch, path, refName, info.Name()); err != nil {
			return err
		}
	}

	return nil
}

func checkRef(ctx context.Context, batch catfile.Batch, refName string, info os.FileInfo) error {
	if info.Size() == 0 {
		return errors.New("checkRef: Ref file is empty")
	}

	_, err := batch.Info(ctx, git.Revision(refName))
	return err
}

func (s *server) fixRef(ctx context.Context, repo *localrepo.Repo, batch catfile.Batch, refPath string, name string, sha string) error {
	// So the ref is broken, let's get rid of it
	if err := os.RemoveAll(refPath); err != nil {
		return err
	}

	// If the sha is not in the the repository, we can't fix it
	if _, err := batch.Info(ctx, git.Revision(sha)); err != nil {
		return nil
	}

	// The name is a valid sha, recreate the ref
	return repo.ExecAndWait(ctx, git.SubCmd{
		Name: "update-ref",
		Args: []string{name, sha},
	}, git.WithRefTxHook(ctx, repo, s.cfg))
}
