package stats

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/gittest"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper"
	"gitlab.com/gitlab-org/gitaly/v14/internal/testhelper/testcfg"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

func TestLogObjectInfo(t *testing.T) {
	cfg := testcfg.Build(t)

	repo1, repoPath1 := gittest.CloneRepo(t, cfg, cfg.Storages[0])
	repo2, repoPath2 := gittest.CloneRepo(t, cfg, cfg.Storages[0])

	ctx, cancel := testhelper.Context()
	defer cancel()

	logBuffer := &bytes.Buffer{}
	log := &logrus.Logger{Out: logBuffer, Formatter: &logrus.JSONFormatter{}, Level: logrus.InfoLevel}
	testCtx := ctxlogrus.ToContext(ctx, log.WithField("test", "logging"))
	gitCmdFactory := git.NewExecCommandFactory(cfg)

	requireLog := func(msg string) map[string]interface{} {
		var out map[string]interface{}
		require.NoError(t, json.NewDecoder(strings.NewReader(msg)).Decode(&out))
		const key = "count_objects"
		require.Contains(t, out, key, "there is no any information about statistics")
		countObjects := out[key].(map[string]interface{})
		require.Contains(t, countObjects, "count")
		require.Contains(t, countObjects, "size")
		require.Contains(t, countObjects, "in-pack")
		require.Contains(t, countObjects, "packs")
		require.Contains(t, countObjects, "size-pack")
		require.Contains(t, countObjects, "garbage")
		require.Contains(t, countObjects, "size-garbage")
		return countObjects
	}

	t.Run("shared repo with multiple alternates", func(t *testing.T) {
		locator := config.NewLocator(cfg)
		storagePath, err := locator.GetStorageByName(repo1.GetStorageName())
		require.NoError(t, err)

		tmpDir, err := ioutil.TempDir(storagePath, "")
		require.NoError(t, err)
		defer func() { require.NoError(t, os.RemoveAll(tmpDir)) }()

		// clone existing local repo with two alternates
		gittest.Exec(t, cfg, "clone", "--shared", repoPath1, "--reference", repoPath1, "--reference", repoPath2, tmpDir)

		logBuffer.Reset()
		LogObjectsInfo(testCtx, gitCmdFactory, &gitalypb.Repository{
			StorageName:  repo1.StorageName,
			RelativePath: filepath.Join(strings.TrimPrefix(tmpDir, storagePath), ".git"),
		})

		countObjects := requireLog(logBuffer.String())
		require.ElementsMatch(t, []string{repoPath1 + "/objects", repoPath2 + "/objects"}, countObjects["alternate"])
	})

	t.Run("repo without alternates", func(t *testing.T) {
		logBuffer.Reset()
		LogObjectsInfo(testCtx, gitCmdFactory, repo2)

		countObjects := requireLog(logBuffer.String())
		require.Contains(t, countObjects, "prune-packable")
	})
}
