package namespace

import (
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitaly/storage"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

type server struct {
	gitalypb.UnimplementedNamespaceServiceServer
	locator storage.Locator
}

// NewServer creates a new instance of a gRPC namespace server
func NewServer(locator storage.Locator) gitalypb.NamespaceServiceServer {
	return &server{locator: locator}
}
