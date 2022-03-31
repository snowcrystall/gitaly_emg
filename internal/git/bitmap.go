package git

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	grpcmwtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/packfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/storage"
)

var badBitmapRequestCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gitaly_bad_bitmap_request_total",
		Help: "RPC calls during which there was not exactly 1 packfile bitmap",
	},
	[]string{"method", "bitmaps"},
)

// WarnIfTooManyBitmaps checks for too many (more than one) bitmaps in
// repoPath, and if it finds any, it logs a warning. This is to help us
// investigate https://gitlab.com/gitlab-org/gitaly/issues/1728.
func WarnIfTooManyBitmaps(ctx context.Context, locator storage.Locator, storageName, repoPath string) {
	logEntry := ctxlogrus.Extract(ctx)

	storageRoot, err := locator.GetStorageByName(storageName)
	if err != nil {
		logEntry.WithError(err).Info("bitmap check failed")
		return
	}

	objdirs, err := ObjectDirectories(ctx, storageRoot, repoPath)
	if err != nil {
		logEntry.WithError(err).Info("bitmap check failed")
		return
	}

	var bitmapCount, packCount int
	seen := make(map[string]bool)
	for _, dir := range objdirs {
		if seen[dir] {
			continue
		}
		seen[dir] = true

		packs, err := packfile.List(dir)
		if err != nil {
			logEntry.WithError(err).Info("bitmap check failed")
			return
		}
		packCount += len(packs)

		for _, p := range packs {
			fi, err := os.Stat(strings.TrimSuffix(p, ".pack") + ".bitmap")
			if err == nil && !fi.IsDir() {
				bitmapCount++
			}
		}
	}

	if bitmapCount == 1 {
		// Exactly one bitmap: this is how things should be.
		return
	}

	if packCount == 0 {
		// If there are no packfiles we don't expect bitmaps nor do we care about
		// them.
		return
	}

	if bitmapCount > 1 {
		logEntry.WithField("bitmaps", bitmapCount).Warn("found more than one packfile bitmap in repository")
	}

	// The case where bitmapCount == 0 is likely to occur early in the life of a
	// repository. We don't want to spam our logs with that, so we count but
	// not log it.

	grpcMethod, ok := grpcmwtags.Extract(ctx).Values()["grpc.request.fullMethod"].(string)
	if !ok {
		return
	}

	badBitmapRequestCount.WithLabelValues(grpcMethod, strconv.Itoa(bitmapCount)).Inc()
}
