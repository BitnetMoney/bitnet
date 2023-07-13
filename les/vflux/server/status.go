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

package server

import (
	"reflect"

	"github.com/ethereum/go-ethereum/p2p/nodestate"
)

type peerWrapper struct{ clientPeer } // the NodeStateMachine type system needs this wrapper

// serverSetup is a wrapper of the node state machine setup, which contains
// all the created flags and fields used in the vflux server side.
type serverSetup struct {
	setup       *nodestate.Setup
	clientField nodestate.Field // Field contains the client peer handler

	// Flags and fields controlled by balance tracker. BalanceTracker
	// is responsible for setting/deleting these flags or fields.
	priorityFlag nodestate.Flags // Flag is set if the node has a positive balance
	updateFlag   nodestate.Flags // Flag is set whenever the node balance is changed(priority changed)
	balanceField nodestate.Field // Field contains the client balance for priority calculation

	// Flags and fields controlled by priority queue. Priority queue
	// is responsible for setting/deleting these flags or fields.
	activeFlag    nodestate.Flags // Flag is set if the node is active
	inactiveFlag  nodestate.Flags // Flag is set if the node is inactive
	capacityField nodestate.Field // Field contains the capacity of the node
	queueField    nodestate.Field // Field contains the information in the priority queue
}

// newServerSetup initializes the setup for state machine and returns the flags/fields group.
func newServerSetup() *serverSetup {
	setup := &serverSetup{setup: &nodestate.Setup{}}
	setup.clientField = setup.setup.NewField("client", reflect.TypeOf(peerWrapper{}))
	setup.priorityFlag = setup.setup.NewFlag("priority")
	setup.updateFlag = setup.setup.NewFlag("update")
	setup.balanceField = setup.setup.NewField("balance", reflect.TypeOf(&nodeBalance{}))
	setup.activeFlag = setup.setup.NewFlag("active")
	setup.inactiveFlag = setup.setup.NewFlag("inactive")
	setup.capacityField = setup.setup.NewField("capacity", reflect.TypeOf(uint64(0)))
	setup.queueField = setup.setup.NewField("queue", reflect.TypeOf(&ppNodeInfo{}))
	return setup
}
