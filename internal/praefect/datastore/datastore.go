// Package datastore provides data models and datastore persistence abstractions
// for tracking the state of repository replicas.
//
// See original design discussion:
// https://gitlab.com/gitlab-org/gitaly/issues/1495
package datastore

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JobState is an enum that indicates the state of a job
type JobState string

func (js JobState) String() string {
	return string(js)
}

const (
	// JobStateReady indicates the job is now ready to proceed.
	JobStateReady = JobState("ready")
	// JobStateInProgress indicates the job is being processed by a worker.
	JobStateInProgress = JobState("in_progress")
	// JobStateCompleted indicates the job is now complete.
	JobStateCompleted = JobState("completed")
	// JobStateCancelled indicates the job was cancelled. This can occur if the
	// job is no longer relevant (e.g. a node is moved out of a repository).
	JobStateCancelled = JobState("cancelled")
	// JobStateFailed indicates the job did not succeed. The Replicator will retry
	// failed jobs.
	JobStateFailed = JobState("failed")
	// JobStateDead indicates the job was retried up to the maximum retries.
	JobStateDead = JobState("dead")
)

// ChangeType indicates what kind of change the replication is propagating
type ChangeType string

const (
	// UpdateRepo is when a replication updates a repository in place
	UpdateRepo = ChangeType("update")
	// CreateRepo is when a replication creates a repo
	CreateRepo = ChangeType("create")
	// DeleteRepo is when a replication deletes a repo
	DeleteRepo = ChangeType("delete")
	// DeleteReplica change type indicates that the targeted replica is due for deletion.
	DeleteReplica = ChangeType("delete_replica")
	// RenameRepo is when a replication renames repo
	RenameRepo = ChangeType("rename")
	// GarbageCollect is when replication runs gc
	GarbageCollect = ChangeType("gc")
	// RepackFull is when replication runs a full repack
	RepackFull = ChangeType("repack_full")
	// RepackIncremental is when replication runs an incremental repack
	RepackIncremental = ChangeType("repack_incremental")
	// Cleanup is when replication runs a repo cleanup
	Cleanup = ChangeType("cleanup")
	// PackRefs is when replication optimizes references in a repo
	PackRefs = ChangeType("pack_refs")
	// WriteCommitGraph is when replication writes a commit graph
	WriteCommitGraph = ChangeType("write_commit_graph")
	// MidxRepack is when replication does a multi-pack-index repack
	MidxRepack = ChangeType("midx_repack")
	// OptimizeRepository is when replication optimizes a repository
	OptimizeRepository = ChangeType("optimize_repository")
)

func (ct ChangeType) String() string {
	return string(ct)
}

// Params represent additional information required to process event after fetching it from storage.
// It must be JSON encodable/decodable to persist it without problems.
type Params map[string]interface{}

// Scan assigns a value from a database driver.
func (p *Params) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	d, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unexpected type received: %T", value)
	}

	return json.Unmarshal(d, p)
}

// Value returns a driver Value.
func (p Params) Value() (driver.Value, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

// GetBool returns the boolean parameter associated with the given key. Returns an error if either
// the key does not exist, or if the value is not a bool.
func (p Params) GetBool(key string) (bool, error) {
	value, found := p[key]
	if !found {
		return false, errors.New("key does not exist")
	}

	booleanValue, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("value is of unexpected type %T", value)
	}

	return booleanValue, nil
}
