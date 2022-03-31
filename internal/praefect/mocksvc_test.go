package praefect

import (
	"context"

	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	repoAccessorUnaryFunc func(context.Context, *mock.RepoRequest) (*emptypb.Empty, error)
	repoMutatorUnaryFunc  func(context.Context, *mock.RepoRequest) (*emptypb.Empty, error)
)

// mockSvc is an implementation of mock.SimpleServer for testing purposes. The
// gRPC stub can be updated by running `make proto`.
type mockSvc struct {
	mock.UnimplementedSimpleServiceServer
	repoAccessorUnary repoAccessorUnaryFunc
	repoMutatorUnary  repoMutatorUnaryFunc
}

// RepoAccessorUnary is implemented by a callback
func (m *mockSvc) RepoAccessorUnary(ctx context.Context, req *mock.RepoRequest) (*emptypb.Empty, error) {
	return m.repoAccessorUnary(ctx, req)
}

// RepoMutatorUnary is implemented by a callback
func (m *mockSvc) RepoMutatorUnary(ctx context.Context, req *mock.RepoRequest) (*emptypb.Empty, error) {
	return m.repoMutatorUnary(ctx, req)
}
