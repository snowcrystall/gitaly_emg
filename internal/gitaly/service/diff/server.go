package diff

import (
	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/catfile"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/localrepo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/git/repository"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/config"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/storage"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

const msgSizeThreshold = 5 * 1024

type server struct {
	gitalypb.UnimplementedDiffServiceServer
	MsgSizeThreshold int
	cfg              config.Cfg
	locator          storage.Locator
	gitCmdFactory    git.CommandFactory
	catfileCache     catfile.Cache
}

// NewServer creates a new instance of a gRPC DiffServer
func NewServer(cfg config.Cfg, locator storage.Locator, gitCmdFactory git.CommandFactory, catfileCache catfile.Cache) gitalypb.DiffServiceServer {
	return &server{
		MsgSizeThreshold: msgSizeThreshold,
		cfg:              cfg,
		locator:          locator,
		gitCmdFactory:    gitCmdFactory,
		catfileCache:     catfileCache,
	}
}

func (s *server) localrepo(repo repository.GitRepo) *localrepo.Repo {
	return localrepo.New(s.gitCmdFactory, s.catfileCache, repo, s.cfg)
}
