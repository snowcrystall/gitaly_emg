package transaction

import (
	"context"
	"errors"

	"gitlab.com/gitlab-org/gitaly/v14/internal/transaction/txinfo"
	"gitlab.com/gitlab-org/gitaly/v14/internal/transaction/voting"
)

// MockManager is a mock Manager for use in tests.
type MockManager struct {
	VoteFn func(context.Context, txinfo.Transaction, voting.Vote) error
	StopFn func(context.Context, txinfo.Transaction) error
}

// Vote calls the MockManager's Vote function, if set. Otherwise, it returns an error.
func (m *MockManager) Vote(ctx context.Context, tx txinfo.Transaction, vote voting.Vote) error {
	if m.VoteFn == nil {
		return errors.New("mock does not implement Vote function")
	}
	return m.VoteFn(ctx, tx, vote)
}

// Stop calls the MockManager's Stop function, if set. Otherwise, it returns an error.
func (m *MockManager) Stop(ctx context.Context, tx txinfo.Transaction) error {
	if m.StopFn == nil {
		return errors.New("mock does not implement Stop function")
	}
	return m.StopFn(ctx, tx)
}
