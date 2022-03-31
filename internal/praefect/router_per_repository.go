package praefect

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/datastore"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/nodes"
	"google.golang.org/grpc"
)

const (
	routeRepositoryAccessorPolicy            = "gitaly-route-repository-accessor-policy"
	routeRepositoryAccessorPolicyPrimaryOnly = "primary-only"
)

// errRepositoryNotFound is retuned when trying to operate on a non-existent repository.
var errRepositoryNotFound = errors.New("repository not found")

// errPrimaryUnassigned is returned when the primary node is not in the set of assigned nodes.
var errPrimaryUnassigned = errors.New("primary node is not assigned")

// AssignmentGetter is an interface for getting repository host node assignments.
type AssignmentGetter interface {
	// GetHostAssignments returns the names of the storages assigned to host the repository.
	// The primary node must always be assigned.
	GetHostAssignments(ctx context.Context, virtualStorage, relativePath string) ([]string, error)
}

// StaticStorageAssignments is a static assignment of the same storages in a virtual storage for every repository.
type StaticStorageAssignments map[string][]string

func (st StaticStorageAssignments) GetHostAssignments(ctx context.Context, virtualStorage, relativePath string) ([]string, error) {
	storages, ok := st[virtualStorage]
	if !ok {
		return nil, nodes.ErrVirtualStorageNotExist
	}

	return storages, nil
}

// ErrNoSuitableNode is returned when there is not suitable node to serve a request.
var ErrNoSuitableNode = errors.New("no suitable node to serve the request")

// ErrNoHealthyNodes is returned when there are no healthy nodes to serve a request.
var ErrNoHealthyNodes = errors.New("no healthy nodes")

// Connections is a set of connections to configured storage nodes by their virtual storages.
type Connections map[string]map[string]*grpc.ClientConn

// PrimaryGetter is an interface for getting a primary of a repository.
type PrimaryGetter interface {
	// GetPrimary returns the primary storage for a given repository.
	GetPrimary(ctx context.Context, virtualStorage string, relativePath string) (string, error)
}

// PerRepositoryRouter implements a router that routes requests respecting per repository primary nodes.
type PerRepositoryRouter struct {
	conns                     Connections
	ag                        AssignmentGetter
	pg                        PrimaryGetter
	rand                      Random
	hc                        HealthChecker
	csg                       datastore.ConsistentStoragesGetter
	defaultReplicationFactors map[string]int
}

// NewPerRepositoryRouter returns a new PerRepositoryRouter using the passed configuration.
func NewPerRepositoryRouter(conns Connections, pg PrimaryGetter, hc HealthChecker, rand Random, csg datastore.ConsistentStoragesGetter, ag AssignmentGetter, defaultReplicationFactors map[string]int) *PerRepositoryRouter {
	return &PerRepositoryRouter{
		conns:                     conns,
		pg:                        pg,
		rand:                      rand,
		hc:                        hc,
		csg:                       csg,
		ag:                        ag,
		defaultReplicationFactors: defaultReplicationFactors,
	}
}

func (r *PerRepositoryRouter) healthyNodes(virtualStorage string) ([]RouterNode, error) {
	conns, ok := r.conns[virtualStorage]
	if !ok {
		return nil, nodes.ErrVirtualStorageNotExist
	}

	healthyNodes := make([]RouterNode, 0, len(conns))
	for _, storage := range r.hc.HealthyNodes()[virtualStorage] {
		conn, ok := conns[storage]
		if !ok {
			return nil, fmt.Errorf("no connection to node %q/%q", virtualStorage, storage)
		}

		healthyNodes = append(healthyNodes, RouterNode{
			Storage:    storage,
			Connection: conn,
		})
	}

	if len(healthyNodes) == 0 {
		return nil, ErrNoHealthyNodes
	}

	return healthyNodes, nil
}

func (r *PerRepositoryRouter) pickRandom(nodes []RouterNode) (RouterNode, error) {
	if len(nodes) == 0 {
		return RouterNode{}, ErrNoSuitableNode
	}

	return nodes[r.rand.Intn(len(nodes))], nil
}

// RouteStorageAccessor routes requests for storage-scoped accessor RPCs. The
// only storage scoped accessor RPC is RemoteService/FindRemoteRepository,
// which in turn executes a command without a repository. This can be done by
// any Gitaly server as it doesn't depend on the state on the server.
func (r *PerRepositoryRouter) RouteStorageAccessor(ctx context.Context, virtualStorage string) (RouterNode, error) {
	healthyNodes, err := r.healthyNodes(virtualStorage)
	if err != nil {
		return RouterNode{}, err
	}

	return r.pickRandom(healthyNodes)
}

// RouteStorageMutator is not implemented here. The only storage scoped mutator RPC is related to namespace operations.
// These are not relevant anymore, given hashed storage is default everywhere, and should be eventually removed.
func (r *PerRepositoryRouter) RouteStorageMutator(ctx context.Context, virtualStorage string) (StorageMutatorRoute, error) {
	return StorageMutatorRoute{}, errors.New("RouteStorageMutator is not implemented on PerRepositoryRouter")
}

func (r *PerRepositoryRouter) RouteRepositoryAccessor(ctx context.Context, virtualStorage, relativePath string, forcePrimary bool) (RouterNode, error) {
	healthyNodes, err := r.healthyNodes(virtualStorage)
	if err != nil {
		return RouterNode{}, err
	}

	if forcePrimary {
		primary, err := r.pg.GetPrimary(ctx, virtualStorage, relativePath)
		if err != nil {
			return RouterNode{}, fmt.Errorf("get primary: %w", err)
		}

		for _, node := range healthyNodes {
			if node.Storage == primary {
				return node, nil
			}
		}

		return RouterNode{}, nodes.ErrPrimaryNotHealthy
	}

	consistentStorages, err := r.csg.GetConsistentStorages(ctx, virtualStorage, relativePath)
	if err != nil {
		return RouterNode{}, fmt.Errorf("consistent storages: %w", err)
	}

	healthyConsistentNodes := make([]RouterNode, 0, len(healthyNodes))
	for _, node := range healthyNodes {
		if _, ok := consistentStorages[node.Storage]; !ok {
			continue
		}

		healthyConsistentNodes = append(healthyConsistentNodes, node)
	}

	return r.pickRandom(healthyConsistentNodes)
}

func (r *PerRepositoryRouter) RouteRepositoryMutator(ctx context.Context, virtualStorage, relativePath string) (RepositoryMutatorRoute, error) {
	healthyNodes, err := r.healthyNodes(virtualStorage)
	if err != nil {
		return RepositoryMutatorRoute{}, err
	}

	primary, err := r.pg.GetPrimary(ctx, virtualStorage, relativePath)
	if err != nil {
		return RepositoryMutatorRoute{}, fmt.Errorf("get primary: %w", err)
	}

	healthySet := make(map[string]RouterNode)
	for _, node := range healthyNodes {
		healthySet[node.Storage] = node
	}

	if _, ok := healthySet[primary]; !ok {
		return RepositoryMutatorRoute{}, nodes.ErrPrimaryNotHealthy
	}

	consistentStorages, err := r.csg.GetConsistentStorages(ctx, virtualStorage, relativePath)
	if err != nil {
		return RepositoryMutatorRoute{}, fmt.Errorf("consistent storages: %w", err)
	}

	if _, ok := consistentStorages[primary]; !ok {
		return RepositoryMutatorRoute{}, ErrRepositoryReadOnly
	}

	assignedStorages, err := r.ag.GetHostAssignments(ctx, virtualStorage, relativePath)
	if err != nil {
		return RepositoryMutatorRoute{}, fmt.Errorf("get host assignments: %w", err)
	}

	var route RepositoryMutatorRoute
	for _, assigned := range assignedStorages {
		node, healthy := healthySet[assigned]
		if assigned == primary {
			route.Primary = node
			continue
		}

		if _, consistent := consistentStorages[node.Storage]; !consistent || !healthy {
			route.ReplicationTargets = append(route.ReplicationTargets, assigned)
			continue
		}

		route.Secondaries = append(route.Secondaries, node)
	}

	if (route.Primary == RouterNode{}) {
		// AssignmentGetter interface defines that the primary must always be assigned.
		// While this case should not commonly happen, we must handle it here for now.
		// This can be triggered if the primary is demoted and unassigned during the RPC call.
		// The three SQL queries above are done non-transactionally. Once the variable
		// replication factor and repository specific primaries are enabled by default, we should
		// combine the above queries in to a single call to remove this case and make the
		// whole operation more efficient.
		return RepositoryMutatorRoute{}, errPrimaryUnassigned
	}

	return route, nil
}

// RouteRepositoryCreation picks a random healthy node to act as the primary node and selects the secondary nodes
// if assignments are enabled. Healthy secondaries take part in the transaction, unhealthy secondaries are set as
// replication targets.
func (r *PerRepositoryRouter) RouteRepositoryCreation(ctx context.Context, virtualStorage string) (RepositoryMutatorRoute, error) {
	healthyNodes, err := r.healthyNodes(virtualStorage)
	if err != nil {
		return RepositoryMutatorRoute{}, err
	}

	primary, err := r.pickRandom(healthyNodes)
	if err != nil {
		return RepositoryMutatorRoute{}, err
	}

	replicationFactor := r.defaultReplicationFactors[virtualStorage]
	if replicationFactor == 1 {
		return RepositoryMutatorRoute{Primary: primary}, nil
	}

	var secondaryNodes []RouterNode
	for storage, conn := range r.conns[virtualStorage] {
		if storage == primary.Storage {
			continue
		}

		secondaryNodes = append(secondaryNodes, RouterNode{
			Storage:    storage,
			Connection: conn,
		})
	}

	// replicationFactor being zero indicates it has not been configured. If so, we fallback to the behavior
	// of no assignments, replicate everywhere and do not select assigned secondaries below.
	if replicationFactor > 0 {
		// Select random secondaries according to the default replication factor.
		r.rand.Shuffle(len(secondaryNodes), func(i, j int) {
			secondaryNodes[i], secondaryNodes[j] = secondaryNodes[j], secondaryNodes[i]
		})

		secondaryNodes = secondaryNodes[:replicationFactor-1]
	}

	var secondaries []RouterNode
	var replicationTargets []string
	for _, secondaryNode := range secondaryNodes {
		isHealthy := false
		for _, healthyNode := range healthyNodes {
			if healthyNode == secondaryNode {
				isHealthy = true
				break
			}
		}

		if isHealthy {
			secondaries = append(secondaries, secondaryNode)
			continue
		}

		replicationTargets = append(replicationTargets, secondaryNode.Storage)
	}

	return RepositoryMutatorRoute{
		Primary:            primary,
		Secondaries:        secondaries,
		ReplicationTargets: replicationTargets,
	}, nil
}
