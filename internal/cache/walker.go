// Package cache supplies background workers for periodically cleaning the
// cache folder on all storages listed in the config file. Upon configuration
// validation, one worker will be started for each storage. The worker will
// walk the cache directory tree and remove any files older than one hour. The
// worker will walk the cache directory every ten minutes.
package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/dontpanic"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/log"
)

func (c *DiskCache) logWalkErr(err error, path, msg string) {
	c.walkerErrorTotal.Inc()
	log.Default().
		WithField("path", path).
		WithError(err).
		Warn(msg)
}

func (c *DiskCache) cleanWalk(path string) error {
	defer time.Sleep(100 * time.Microsecond) // relieve pressure

	c.walkerCheckTotal.Inc()
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		c.logWalkErr(err, path, "unable to stat directory")
		return err
	}

	for _, e := range entries {
		ePath := filepath.Join(path, e.Name())

		if e.IsDir() {
			if err := c.cleanWalk(ePath); err != nil {
				return err
			}
			continue
		}

		c.walkerCheckTotal.Inc()
		if time.Since(e.ModTime()) < staleAge {
			continue // still fresh
		}

		// file is stale
		if err := os.Remove(ePath); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			c.logWalkErr(err, ePath, "unable to remove file")
			return err
		}
		c.walkerRemovalTotal.Inc()
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		c.logWalkErr(err, path, "unable to stat directory after walk")
		return err
	}

	if len(files) == 0 {
		c.walkerEmptyDirTotal.Inc()
		if err := os.Remove(path); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			c.logWalkErr(err, path, "unable to remove empty directory")
			return err
		}
		c.walkerEmptyDirRemovalTotal.Inc()
		c.walkerRemovalTotal.Inc()
	}

	return nil
}

const cleanWalkFrequency = 10 * time.Minute

func (c *DiskCache) walkLoop(walkPath string) {
	logger := logrus.WithField("path", walkPath)
	logger.Infof("Starting file walker for %s", walkPath)

	walkTick := time.NewTicker(cleanWalkFrequency)
	dontpanic.GoForever(time.Minute, func() {
		if err := c.cleanWalk(walkPath); err != nil {
			logger.Error(err)
		}
		<-walkTick.C
	})
}

func (c *DiskCache) startCleanWalker(cacheDir, stateDir string) {
	if c.cacheConfig.disableWalker {
		return
	}

	c.walkLoop(cacheDir)
	c.walkLoop(stateDir)
}

// moveAndClear will move the cache to the storage location's
// temporary folder, and then remove its contents asynchronously
func (c *DiskCache) moveAndClear(storage config.Storage) error {
	if c.cacheConfig.disableMoveAndClear {
		return nil
	}

	logger := logrus.WithField("storage", storage.Name)
	logger.Info("clearing disk cache object folder")

	tempPath, err := c.locator.TempDir(storage.Name)
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}

	if err := os.MkdirAll(tempPath, 0755); err != nil {
		return err
	}

	tmpDir, err := ioutil.TempDir(tempPath, "diskcache")
	if err != nil {
		return err
	}

	defer func() {
		dontpanic.Go(func() {
			start := time.Now()
			if err := os.RemoveAll(tmpDir); err != nil {
				logger.Errorf("unable to remove disk cache objects: %q", err)
				return
			}

			logger.Infof("cleared all cache object files in %s after %s", tmpDir, time.Since(start))
		})
	}()

	logger.Infof("moving disk cache object folder to %s", tmpDir)

	cachePath, err := c.locator.CacheDir(storage.Name)
	if err != nil {
		return fmt.Errorf("cache dir: %w", err)
	}

	if err := os.Rename(cachePath, filepath.Join(tmpDir, "moved")); err != nil {
		if os.IsNotExist(err) {
			logger.Info("disk cache object folder doesn't exist, no need to remove")
			return nil
		}

		return err
	}

	return nil
}

// StartWalkers starts the cache walker Goroutines. Initially, this function will try to clean up
// any preexisting cache directories.
func (c *DiskCache) StartWalkers() error {
	// Deduplicate storages by path.
	storageByPath := map[string]config.Storage{}
	for _, storage := range c.storages {
		storageByPath[storage.Path] = storage
	}

	for _, storage := range storageByPath {
		cacheDir, err := c.locator.CacheDir(storage.Name)
		if err != nil {
			return fmt.Errorf("cache dir: %w", err)
		}

		stateDir, err := c.locator.StateDir(storage.Name)
		if err != nil {
			return fmt.Errorf("state dir: %w", err)
		}

		if err := c.moveAndClear(storage); err != nil {
			return err
		}

		c.startCleanWalker(cacheDir, stateDir)
	}

	return nil
}
