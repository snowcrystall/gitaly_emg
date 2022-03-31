package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/config"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type (
	virtualStorage string
	gitalyStorage  string
)

type nodePing struct {
	address string
	// set of storages this node hosts
	storages  map[gitalyStorage][]virtualStorage
	vStorages map[virtualStorage]struct{} // set of virtual storages node belongs to
	token     string                      // auth token
	err       error                       // any error during dial/ping
}

func flattenNodes(conf config.Config) map[string]*nodePing {
	nodeByAddress := map[string]*nodePing{} // key is address

	// flatten nodes between virtual storages
	for _, vs := range conf.VirtualStorages {
		vsName := virtualStorage(vs.Name)
		for _, node := range vs.Nodes {
			gsName := gitalyStorage(node.Storage)

			n, ok := nodeByAddress[node.Address]
			if !ok {
				n = &nodePing{
					storages:  map[gitalyStorage][]virtualStorage{},
					vStorages: map[virtualStorage]struct{}{},
				}
			}
			n.address = node.Address

			s := n.storages[gsName]
			n.storages[gsName] = append(s, vsName)

			n.vStorages[vsName] = struct{}{}
			n.token = node.Token
			nodeByAddress[node.Address] = n
		}
	}
	return nodeByAddress
}

type dialNodesSubcommand struct{}

func (s *dialNodesSubcommand) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("dial-nodes", flag.ExitOnError)
}

func (s *dialNodesSubcommand) Exec(flags *flag.FlagSet, conf config.Config) error {
	nodes := flattenNodes(conf)

	var wg sync.WaitGroup
	for _, n := range nodes {
		wg.Add(1)
		go func(n *nodePing) {
			defer wg.Done()
			n.checkNode()
		}(n)
	}
	wg.Wait()

	var err error
	for _, n := range nodes {
		if n.err != nil {
			err = n.err
		}
	}

	return err
}

func (npr *nodePing) dial() (*grpc.ClientConn, error) {
	return subCmdDial(npr.address, npr.token)
}

func (npr *nodePing) healthCheck(cc *grpc.ClientConn) (grpc_health_v1.HealthCheckResponse_ServingStatus, error) {
	hClient := grpc_health_v1.NewHealthClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := hClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return 0, err
	}

	return resp.GetStatus(), nil
}

func (npr *nodePing) isConsistent(cc *grpc.ClientConn) bool {
	praefect := gitalypb.NewServerServiceClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(npr.storages) == 0 {
		npr.log("ERROR: current configuration has no storages")
		return false
	}

	resp, err := praefect.ServerInfo(ctx, &gitalypb.ServerInfoRequest{})
	if err != nil {
		npr.log("ERROR: failed to receive state from the remote: %v", err)
		return false
	}

	if len(resp.StorageStatuses) == 0 {
		npr.log("ERROR: remote has no configured storages")
		return false
	}

	storagesSet := make(map[gitalyStorage]bool, len(resp.StorageStatuses))

	knownStoragesSet := make(map[gitalyStorage]bool, len(npr.storages))
	for k := range npr.storages {
		knownStoragesSet[k] = true
	}

	consistent := true
	for _, status := range resp.StorageStatuses {
		gStorage := gitalyStorage(status.StorageName)

		// only proceed if the gitaly storage belongs to a configured
		// virtual storage
		if len(npr.storages[gStorage]) == 0 {
			continue
		}

		if storagesSet[gStorage] {
			npr.log("ERROR: remote has duplicated storage: %q", status.StorageName)
			consistent = false
			continue
		}
		storagesSet[gStorage] = true

		if status.Readable && status.Writeable {
			npr.log(
				"SUCCESS: confirmed Gitaly storage %q in virtual storages %v is served",
				status.StorageName,
				npr.storages[gStorage],
			)
			delete(knownStoragesSet, gStorage) // storage found
		} else {
			npr.log("ERROR: storage %q is not readable or writable", status.StorageName)
			consistent = false
		}
	}

	for storage := range knownStoragesSet {
		npr.log("ERROR: configured storage was not reported by remote: %q", storage)
		consistent = false
	}

	return consistent
}

func (npr *nodePing) log(msg string, args ...interface{}) {
	log.Printf("[%s]: %s", npr.address, fmt.Sprintf(msg, args...))
}

func (npr *nodePing) checkNode() {
	npr.log("dialing...")
	cc, err := npr.dial()
	if err != nil {
		npr.log("ERROR: dialing failed: %v", err)
		npr.err = err
		return
	}
	defer cc.Close()
	npr.log("dialed successfully!")

	npr.log("checking health...")
	health, err := npr.healthCheck(cc)
	if err != nil {
		npr.log("ERROR: unable to request health check: %v", err)
		npr.err = err
		return
	}

	if health != grpc_health_v1.HealthCheckResponse_SERVING {
		npr.err = fmt.Errorf(
			"health check did not report serving, instead reported: %s",
			health.String())
		npr.log("ERROR: %v", npr.err)
		return
	}

	npr.log("SUCCESS: node is healthy!")

	npr.log("checking consistency...")
	if !npr.isConsistent(cc) {
		npr.err = errors.New("consistency check failed")
		npr.log("ERROR: %v", npr.err)
		return
	}
	npr.log("SUCCESS: node configuration is consistent!")
}
