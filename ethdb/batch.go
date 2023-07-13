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

package ethdb

// IdealBatchSize defines the size of the data batches should ideally add in one
// write.
const IdealBatchSize = 100 * 1024

// Batch is a write-only database that commits changes to its host database
// when Write is called. A batch cannot be used concurrently.
type Batch interface {
	KeyValueWriter

	// ValueSize retrieves the amount of data queued up for writing.
	ValueSize() int

	// Write flushes any accumulated data to disk.
	Write() error

	// Reset resets the batch for reuse.
	Reset()

	// Replay replays the batch contents.
	Replay(w KeyValueWriter) error
}

// Batcher wraps the NewBatch method of a backing data store.
type Batcher interface {
	// NewBatch creates a write-only database that buffers changes to its host db
	// until a final write is called.
	NewBatch() Batch

	// NewBatchWithSize creates a write-only database batch with pre-allocated buffer.
	NewBatchWithSize(size int) Batch
}

// HookedBatch wraps an arbitrary batch where each operation may be hooked into
// to monitor from black box code.
type HookedBatch struct {
	Batch

	OnPut    func(key []byte, value []byte) // Callback if a key is inserted
	OnDelete func(key []byte)               // Callback if a key is deleted
}

// Put inserts the given value into the key-value data store.
func (b HookedBatch) Put(key []byte, value []byte) error {
	if b.OnPut != nil {
		b.OnPut(key, value)
	}
	return b.Batch.Put(key, value)
}

// Delete removes the key from the key-value data store.
func (b HookedBatch) Delete(key []byte) error {
	if b.OnDelete != nil {
		b.OnDelete(key)
	}
	return b.Batch.Delete(key)
}
