package repository

import (
	"context"
	"fmt"
	"math"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/repository"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/stats"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	repackCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitaly_repack_total",
			Help: "Counter of Git repack operations",
		},
		[]string{"bitmap"},
	)
)

func init() {
	prometheus.MustRegister(repackCounter)
}

// log2Threads returns the log-2 number of threads based on the number of
// provided CPUs. This prevents repacking operations from exhausting all
// available CPUs and increasing request latency
func log2Threads(numCPUs int) git.ValueFlag {
	n := math.Max(1, math.Floor(math.Log2(float64(numCPUs))))
	return git.ValueFlag{Name: "--threads", Value: fmt.Sprint(n)}
}

func (s *server) RepackFull(ctx context.Context, in *gitalypb.RepackFullRequest) (*gitalypb.RepackFullResponse, error) {
	options := []git.Option{
		git.Flag{Name: "-A"},
		git.Flag{Name: "--pack-kept-objects"},
		git.Flag{Name: "-l"},
		log2Threads(runtime.NumCPU()),
	}
	if err := s.repackCommand(ctx, in.GetRepository(), in.GetCreateBitmap(), options...); err != nil {
		return nil, err
	}
	return &gitalypb.RepackFullResponse{}, nil
}

func (s *server) RepackIncremental(ctx context.Context, in *gitalypb.RepackIncrementalRequest) (*gitalypb.RepackIncrementalResponse, error) {
	if err := s.repackCommand(ctx, in.GetRepository(), false); err != nil {
		return nil, err
	}
	return &gitalypb.RepackIncrementalResponse{}, nil
}

func (s *server) repackCommand(ctx context.Context, repo repository.GitRepo, bitmap bool, args ...git.Option) error {
	cmd, err := s.gitCmdFactory.New(ctx, repo,
		git.SubCmd{
			Name:  "repack",
			Flags: append([]git.Option{git.Flag{Name: "-d"}}, args...),
		},
		git.WithConfig(repackConfig(ctx, bitmap)...),
	)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return err
		}
		return status.Errorf(codes.Internal, err.Error())
	}

	if err := cmd.Wait(); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	if err = s.writeCommitGraph(ctx, repo, gitalypb.WriteCommitGraphRequest_SizeMultiple); err != nil {
		return err
	}

	stats.LogObjectsInfo(ctx, s.gitCmdFactory, repo)

	return nil
}

func repackConfig(ctx context.Context, bitmap bool) []git.ConfigPair {
	config := []git.ConfigPair{
		git.ConfigPair{Key: "pack.island", Value: "r(e)fs/heads"},
		git.ConfigPair{Key: "pack.island", Value: "r(e)fs/tags"},
		git.ConfigPair{Key: "pack.islandCore", Value: "e"},
		git.ConfigPair{Key: "repack.useDeltaIslands", Value: "true"},
	}

	if bitmap {
		config = append(config, git.ConfigPair{Key: "repack.writeBitmaps", Value: "true"})
		config = append(config, git.ConfigPair{Key: "pack.writeBitmapHashCache", Value: "true"})
	} else {
		config = append(config, git.ConfigPair{Key: "repack.writeBitmaps", Value: "false"})
	}

	repackCounter.WithLabelValues(fmt.Sprint(bitmap)).Inc()

	return config
}
