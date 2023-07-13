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

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/miner"
)

// maxTrackedPayloads is the maximum number of prepared payloads the execution
// engine tracks before evicting old ones. Ideally we should only ever track the
// latest one; but have a slight wiggle room for non-ideal conditions.
const maxTrackedPayloads = 10

// maxTrackedHeaders is the maximum number of executed payloads the execution
// engine tracks before evicting old ones. These are tracked outside the chain
// during initial sync to allow ForkchoiceUpdate to reference past blocks via
// hashes only. For the sync target it would be enough to track only the latest
// header, but snap sync also needs the latest finalized height for the ancient
// limit.
const maxTrackedHeaders = 96

// payloadQueueItem represents an id->payload tuple to store until it's retrieved
// or evicted.
type payloadQueueItem struct {
	id      engine.PayloadID
	payload *miner.Payload
}

// payloadQueue tracks the latest handful of constructed payloads to be retrieved
// by the beacon chain if block production is requested.
type payloadQueue struct {
	payloads []*payloadQueueItem
	lock     sync.RWMutex
}

// newPayloadQueue creates a pre-initialized queue with a fixed number of slots
// all containing empty items.
func newPayloadQueue() *payloadQueue {
	return &payloadQueue{
		payloads: make([]*payloadQueueItem, maxTrackedPayloads),
	}
}

// put inserts a new payload into the queue at the given id.
func (q *payloadQueue) put(id engine.PayloadID, payload *miner.Payload) {
	q.lock.Lock()
	defer q.lock.Unlock()

	copy(q.payloads[1:], q.payloads)
	q.payloads[0] = &payloadQueueItem{
		id:      id,
		payload: payload,
	}
}

// get retrieves a previously stored payload item or nil if it does not exist.
func (q *payloadQueue) get(id engine.PayloadID) *engine.ExecutionPayloadEnvelope {
	q.lock.RLock()
	defer q.lock.RUnlock()

	for _, item := range q.payloads {
		if item == nil {
			return nil // no more items
		}
		if item.id == id {
			return item.payload.Resolve()
		}
	}
	return nil
}

// has checks if a particular payload is already tracked.
func (q *payloadQueue) has(id engine.PayloadID) bool {
	q.lock.RLock()
	defer q.lock.RUnlock()

	for _, item := range q.payloads {
		if item == nil {
			return false
		}
		if item.id == id {
			return true
		}
	}
	return false
}

// headerQueueItem represents an hash->header tuple to store until it's retrieved
// or evicted.
type headerQueueItem struct {
	hash   common.Hash
	header *types.Header
}

// headerQueue tracks the latest handful of constructed headers to be retrieved
// by the beacon chain if block production is requested.
type headerQueue struct {
	headers []*headerQueueItem
	lock    sync.RWMutex
}

// newHeaderQueue creates a pre-initialized queue with a fixed number of slots
// all containing empty items.
func newHeaderQueue() *headerQueue {
	return &headerQueue{
		headers: make([]*headerQueueItem, maxTrackedHeaders),
	}
}

// put inserts a new header into the queue at the given hash.
func (q *headerQueue) put(hash common.Hash, data *types.Header) {
	q.lock.Lock()
	defer q.lock.Unlock()

	copy(q.headers[1:], q.headers)
	q.headers[0] = &headerQueueItem{
		hash:   hash,
		header: data,
	}
}

// get retrieves a previously stored header item or nil if it does not exist.
func (q *headerQueue) get(hash common.Hash) *types.Header {
	q.lock.RLock()
	defer q.lock.RUnlock()

	for _, item := range q.headers {
		if item == nil {
			return nil // no more items
		}
		if item.hash == hash {
			return item.header
		}
	}
	return nil
}
