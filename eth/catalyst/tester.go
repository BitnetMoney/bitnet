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

package catalyst

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
)

// FullSyncTester is an auxiliary service that allows Geth to perform full sync
// alone without consensus-layer attached. Users must specify a valid block as
// the sync target. This tester can be applied to different networks, no matter
// it's pre-merge or post-merge, but only for full-sync.
type FullSyncTester struct {
	api    *ConsensusAPI
	block  *types.Block
	closed chan struct{}
	wg     sync.WaitGroup
}

// RegisterFullSyncTester registers the full-sync tester service into the node
// stack for launching and stopping the service controlled by node.
func RegisterFullSyncTester(stack *node.Node, backend *eth.Ethereum, block *types.Block) (*FullSyncTester, error) {
	cl := &FullSyncTester{
		api:    NewConsensusAPI(backend),
		block:  block,
		closed: make(chan struct{}),
	}
	stack.RegisterLifecycle(cl)
	return cl, nil
}

// Start launches the beacon sync with provided sync target.
func (tester *FullSyncTester) Start() error {
	tester.wg.Add(1)
	go func() {
		defer tester.wg.Done()

		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Don't bother downloader in case it's already syncing.
				if tester.api.eth.Downloader().Synchronising() {
					continue
				}
				// Short circuit in case the target block is already stored
				// locally. TODO(somehow terminate the node stack if target
				// is reached).
				if tester.api.eth.BlockChain().HasBlock(tester.block.Hash(), tester.block.NumberU64()) {
					log.Info("Full-sync target reached", "number", tester.block.NumberU64(), "hash", tester.block.Hash())
					return
				}
				// Trigger beacon sync with the provided block header as
				// trusted chain head.
				err := tester.api.eth.Downloader().BeaconSync(downloader.FullSync, tester.block.Header(), nil)
				if err != nil {
					log.Info("Failed to beacon sync", "err", err)
				}

			case <-tester.closed:
				return
			}
		}
	}()
	return nil
}

// Stop stops the full-sync tester to stop all background activities.
// This function can only be called for one time.
func (tester *FullSyncTester) Stop() error {
	close(tester.closed)
	tester.wg.Wait()
	return nil
}
