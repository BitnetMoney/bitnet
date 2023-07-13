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

package client

import (
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/p2p/nodestate"
)

var (
	testSetup     = &nodestate.Setup{}
	sfTest1       = testSetup.NewFlag("test1")
	sfTest2       = testSetup.NewFlag("test2")
	sfTest3       = testSetup.NewFlag("test3")
	sfTest4       = testSetup.NewFlag("test4")
	sfiTestWeight = testSetup.NewField("nodeWeight", reflect.TypeOf(uint64(0)))
)

const iterTestNodeCount = 6

func TestWrsIterator(t *testing.T) {
	ns := nodestate.NewNodeStateMachine(nil, nil, &mclock.Simulated{}, testSetup)
	w := NewWrsIterator(ns, sfTest2, sfTest3.Or(sfTest4), sfiTestWeight)
	ns.Start()
	for i := 1; i <= iterTestNodeCount; i++ {
		ns.SetState(testNode(i), sfTest1, nodestate.Flags{}, 0)
		ns.SetField(testNode(i), sfiTestWeight, uint64(1))
	}
	next := func() int {
		ch := make(chan struct{})
		go func() {
			w.Next()
			close(ch)
		}()
		select {
		case <-ch:
		case <-time.After(time.Second * 5):
			t.Fatalf("Iterator.Next() timeout")
		}
		node := w.Node()
		ns.SetState(node, sfTest4, nodestate.Flags{}, 0)
		return testNodeIndex(node.ID())
	}
	set := make(map[int]bool)
	expset := func() {
		for len(set) > 0 {
			n := next()
			if !set[n] {
				t.Errorf("Item returned by iterator not in the expected set (got %d)", n)
			}
			delete(set, n)
		}
	}

	ns.SetState(testNode(1), sfTest2, nodestate.Flags{}, 0)
	ns.SetState(testNode(2), sfTest2, nodestate.Flags{}, 0)
	ns.SetState(testNode(3), sfTest2, nodestate.Flags{}, 0)
	set[1] = true
	set[2] = true
	set[3] = true
	expset()
	ns.SetState(testNode(4), sfTest2, nodestate.Flags{}, 0)
	ns.SetState(testNode(5), sfTest2.Or(sfTest3), nodestate.Flags{}, 0)
	ns.SetState(testNode(6), sfTest2, nodestate.Flags{}, 0)
	set[4] = true
	set[6] = true
	expset()
	ns.SetField(testNode(2), sfiTestWeight, uint64(0))
	ns.SetState(testNode(1), nodestate.Flags{}, sfTest4, 0)
	ns.SetState(testNode(2), nodestate.Flags{}, sfTest4, 0)
	ns.SetState(testNode(3), nodestate.Flags{}, sfTest4, 0)
	set[1] = true
	set[3] = true
	expset()
	ns.SetField(testNode(2), sfiTestWeight, uint64(1))
	ns.SetState(testNode(2), nodestate.Flags{}, sfTest2, 0)
	ns.SetState(testNode(1), nodestate.Flags{}, sfTest4, 0)
	ns.SetState(testNode(2), sfTest2, sfTest4, 0)
	ns.SetState(testNode(3), nodestate.Flags{}, sfTest4, 0)
	set[1] = true
	set[2] = true
	set[3] = true
	expset()
	ns.Stop()
}
