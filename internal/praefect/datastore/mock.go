package datastore

import "context"

// MockReplicationEventQueue is a helper for tests that implements ReplicationEventQueue
// and allows for parametrizing behavior.
type MockReplicationEventQueue struct {
	ReplicationEventQueue
	EnqueueFunc func(context.Context, ReplicationEvent) (ReplicationEvent, error)
}

func (m *MockReplicationEventQueue) Enqueue(ctx context.Context, event ReplicationEvent) (ReplicationEvent, error) {
	return m.EnqueueFunc(ctx, event)
}
