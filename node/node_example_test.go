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

package node_test

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/node"
)

// SampleLifecycle is a trivial network service that can be attached to a node for
// life cycle management.
//
// The following methods are needed to implement a node.Lifecycle:
//   - Start() error              - method invoked when the node is ready to start the service
//   - Stop() error               - method invoked when the node terminates the service
type SampleLifecycle struct{}

func (s *SampleLifecycle) Start() error { fmt.Println("Service starting..."); return nil }
func (s *SampleLifecycle) Stop() error  { fmt.Println("Service stopping..."); return nil }

func ExampleLifecycle() {
	// Create a network node to run protocols with the default values.
	stack, err := node.New(&node.Config{})
	if err != nil {
		log.Fatalf("Failed to create network node: %v", err)
	}
	defer stack.Close()

	// Create and register a simple network Lifecycle.
	service := new(SampleLifecycle)
	stack.RegisterLifecycle(service)

	// Boot up the entire protocol stack, do a restart and terminate
	if err := stack.Start(); err != nil {
		log.Fatalf("Failed to start the protocol stack: %v", err)
	}
	if err := stack.Close(); err != nil {
		log.Fatalf("Failed to stop the protocol stack: %v", err)
	}
	// Output:
	// Service starting...
	// Service stopping...
}
