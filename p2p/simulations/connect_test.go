// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package simulations

import (
	"testing"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

func newTestNetwork(t *testing.T, nodeCount int) (*Network, []enode.ID) {
	t.Helper()
	adapter := adapters.NewSimAdapter(adapters.LifecycleConstructors{
		"noopwoop": func(ctx *adapters.ServiceContext, stack *node.Node) (node.Lifecycle, error) {
			return NewNoopService(nil), nil
		},
	})

	// create network
	network := NewNetwork(adapter, &NetworkConfig{
		DefaultService: "noopwoop",
	})

	// create and start nodes
	ids := make([]enode.ID, nodeCount)
	for i := range ids {
		conf := adapters.RandomNodeConfig()
		node, err := network.NewNodeWithConfig(conf)
		if err != nil {
			t.Fatalf("error creating node: %s", err)
		}
		if err := network.Start(node.ID()); err != nil {
			t.Fatalf("error starting node: %s", err)
		}
		ids[i] = node.ID()
	}

	if len(network.Conns) > 0 {
		t.Fatal("no connections should exist after just adding nodes")
	}

	return network, ids
}

func TestConnectToLastNode(t *testing.T) {
	net, ids := newTestNetwork(t, 10)
	defer net.Shutdown()

	first := ids[0]
	if err := net.ConnectToLastNode(first); err != nil {
		t.Fatal(err)
	}

	last := ids[len(ids)-1]
	for i, id := range ids {
		if id == first || id == last {
			continue
		}

		if net.GetConn(first, id) != nil {
			t.Errorf("connection must not exist with node(ind: %v, id: %v)", i, id)
		}
	}

	if net.GetConn(first, last) == nil {
		t.Error("first and last node must be connected")
	}
}

func TestConnectToRandomNode(t *testing.T) {
	net, ids := newTestNetwork(t, 10)
	defer net.Shutdown()

	err := net.ConnectToRandomNode(ids[0])
	if err != nil {
		t.Fatal(err)
	}

	var cc int
	for i, a := range ids {
		for _, b := range ids[i:] {
			if net.GetConn(a, b) != nil {
				cc++
			}
		}
	}

	if cc != 1 {
		t.Errorf("expected one connection, got %v", cc)
	}
}

func TestConnectNodesFull(t *testing.T) {
	tests := []struct {
		name      string
		nodeCount int
	}{
		{name: "no node", nodeCount: 0},
		{name: "single node", nodeCount: 1},
		{name: "2 nodes", nodeCount: 2},
		{name: "3 nodes", nodeCount: 3},
		{name: "even number of nodes", nodeCount: 12},
		{name: "odd number of nodes", nodeCount: 13},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			net, ids := newTestNetwork(t, test.nodeCount)
			defer net.Shutdown()

			err := net.ConnectNodesFull(ids)
			if err != nil {
				t.Fatal(err)
			}

			VerifyFull(t, net, ids)
		})
	}
}

func TestConnectNodesChain(t *testing.T) {
	net, ids := newTestNetwork(t, 10)
	defer net.Shutdown()

	err := net.ConnectNodesChain(ids)
	if err != nil {
		t.Fatal(err)
	}

	VerifyChain(t, net, ids)
}

func TestConnectNodesRing(t *testing.T) {
	net, ids := newTestNetwork(t, 10)
	defer net.Shutdown()

	err := net.ConnectNodesRing(ids)
	if err != nil {
		t.Fatal(err)
	}

	VerifyRing(t, net, ids)
}

func TestConnectNodesStar(t *testing.T) {
	net, ids := newTestNetwork(t, 10)
	defer net.Shutdown()

	pivotIndex := 2

	err := net.ConnectNodesStar(ids, ids[pivotIndex])
	if err != nil {
		t.Fatal(err)
	}

	VerifyStar(t, net, ids, pivotIndex)
}
