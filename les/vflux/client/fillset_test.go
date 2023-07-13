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
	"crypto/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/p2p/nodestate"
)

type testIter struct {
	waitCh chan struct{}
	nodeCh chan *enode.Node
	node   *enode.Node
}

func (i *testIter) Next() bool {
	if _, ok := <-i.waitCh; !ok {
		return false
	}
	i.node = <-i.nodeCh
	return true
}

func (i *testIter) Node() *enode.Node {
	return i.node
}

func (i *testIter) Close() {
	close(i.waitCh)
}

func (i *testIter) push() {
	var id enode.ID
	rand.Read(id[:])
	i.nodeCh <- enode.SignNull(new(enr.Record), id)
}

func (i *testIter) waiting(timeout time.Duration) bool {
	select {
	case i.waitCh <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

func TestFillSet(t *testing.T) {
	ns := nodestate.NewNodeStateMachine(nil, nil, &mclock.Simulated{}, testSetup)
	iter := &testIter{
		waitCh: make(chan struct{}),
		nodeCh: make(chan *enode.Node),
	}
	fs := NewFillSet(ns, iter, sfTest1)
	ns.Start()

	expWaiting := func(i int, push bool) {
		for ; i > 0; i-- {
			if !iter.waiting(time.Second * 10) {
				t.Fatalf("FillSet not waiting for new nodes")
			}
			if push {
				iter.push()
			}
		}
	}

	expNotWaiting := func() {
		if iter.waiting(time.Millisecond * 100) {
			t.Fatalf("FillSet unexpectedly waiting for new nodes")
		}
	}

	expNotWaiting()
	fs.SetTarget(3)
	expWaiting(3, true)
	expNotWaiting()
	fs.SetTarget(100)
	expWaiting(2, true)
	expWaiting(1, false)
	// lower the target before the previous one has been filled up
	fs.SetTarget(0)
	iter.push()
	expNotWaiting()
	fs.SetTarget(10)
	expWaiting(4, true)
	expNotWaiting()
	// remove all previously set flags
	ns.ForEach(sfTest1, nodestate.Flags{}, func(node *enode.Node, state nodestate.Flags) {
		ns.SetState(node, nodestate.Flags{}, sfTest1, 0)
	})
	// now expect FillSet to fill the set up again with 10 new nodes
	expWaiting(10, true)
	expNotWaiting()

	fs.Close()
	ns.Stop()
}
