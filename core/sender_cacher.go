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

package core

import (
	"runtime"

	"github.com/ethereum/go-ethereum/core/types"
)

// SenderCacher is a concurrent transaction sender recoverer and cacher.
var SenderCacher = newTxSenderCacher(runtime.NumCPU())

// txSenderCacherRequest is a request for recovering transaction senders with a
// specific signature scheme and caching it into the transactions themselves.
//
// The inc field defines the number of transactions to skip after each recovery,
// which is used to feed the same underlying input array to different threads but
// ensure they process the early transactions fast.
type txSenderCacherRequest struct {
	signer types.Signer
	txs    []*types.Transaction
	inc    int
}

// txSenderCacher is a helper structure to concurrently ecrecover transaction
// senders from digital signatures on background threads.
type txSenderCacher struct {
	threads int
	tasks   chan *txSenderCacherRequest
}

// newTxSenderCacher creates a new transaction sender background cacher and starts
// as many processing goroutines as allowed by the GOMAXPROCS on construction.
func newTxSenderCacher(threads int) *txSenderCacher {
	cacher := &txSenderCacher{
		tasks:   make(chan *txSenderCacherRequest, threads),
		threads: threads,
	}
	for i := 0; i < threads; i++ {
		go cacher.cache()
	}
	return cacher
}

// cache is an infinite loop, caching transaction senders from various forms of
// data structures.
func (cacher *txSenderCacher) cache() {
	for task := range cacher.tasks {
		for i := 0; i < len(task.txs); i += task.inc {
			types.Sender(task.signer, task.txs[i])
		}
	}
}

// Recover recovers the senders from a batch of transactions and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
func (cacher *txSenderCacher) Recover(signer types.Signer, txs []*types.Transaction) {
	// If there's nothing to recover, abort
	if len(txs) == 0 {
		return
	}
	// Ensure we have meaningful task sizes and schedule the recoveries
	tasks := cacher.threads
	if len(txs) < tasks*4 {
		tasks = (len(txs) + 3) / 4
	}
	for i := 0; i < tasks; i++ {
		cacher.tasks <- &txSenderCacherRequest{
			signer: signer,
			txs:    txs[i:],
			inc:    tasks,
		}
	}
}

// RecoverFromBlocks recovers the senders from a batch of blocks and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
func (cacher *txSenderCacher) RecoverFromBlocks(signer types.Signer, blocks []*types.Block) {
	count := 0
	for _, block := range blocks {
		count += len(block.Transactions())
	}
	txs := make([]*types.Transaction, 0, count)
	for _, block := range blocks {
		txs = append(txs, block.Transactions()...)
	}
	cacher.Recover(signer, txs)
}
