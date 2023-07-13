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

package trie

import "github.com/ethereum/go-ethereum/common"

// tracer tracks the changes of trie nodes. During the trie operations,
// some nodes can be deleted from the trie, while these deleted nodes
// won't be captured by trie.Hasher or trie.Committer. Thus, these deleted
// nodes won't be removed from the disk at all. Tracer is an auxiliary tool
// used to track all insert and delete operations of trie and capture all
// deleted nodes eventually.
//
// The changed nodes can be mainly divided into two categories: the leaf
// node and intermediate node. The former is inserted/deleted by callers
// while the latter is inserted/deleted in order to follow the rule of trie.
// This tool can track all of them no matter the node is embedded in its
// parent or not, but valueNode is never tracked.
//
// Besides, it's also used for recording the original value of the nodes
// when they are resolved from the disk. The pre-value of the nodes will
// be used to construct trie history in the future.
//
// Note tracer is not thread-safe, callers should be responsible for handling
// the concurrency issues by themselves.
type tracer struct {
	inserts    map[string]struct{}
	deletes    map[string]struct{}
	accessList map[string][]byte
}

// newTracer initializes the tracer for capturing trie changes.
func newTracer() *tracer {
	return &tracer{
		inserts:    make(map[string]struct{}),
		deletes:    make(map[string]struct{}),
		accessList: make(map[string][]byte),
	}
}

// onRead tracks the newly loaded trie node and caches the rlp-encoded
// blob internally. Don't change the value outside of function since
// it's not deep-copied.
func (t *tracer) onRead(path []byte, val []byte) {
	t.accessList[string(path)] = val
}

// onInsert tracks the newly inserted trie node. If it's already
// in the deletion set (resurrected node), then just wipe it from
// the deletion set as it's "untouched".
func (t *tracer) onInsert(path []byte) {
	if _, present := t.deletes[string(path)]; present {
		delete(t.deletes, string(path))
		return
	}
	t.inserts[string(path)] = struct{}{}
}

// onDelete tracks the newly deleted trie node. If it's already
// in the addition set, then just wipe it from the addition set
// as it's untouched.
func (t *tracer) onDelete(path []byte) {
	if _, present := t.inserts[string(path)]; present {
		delete(t.inserts, string(path))
		return
	}
	t.deletes[string(path)] = struct{}{}
}

// reset clears the content tracked by tracer.
func (t *tracer) reset() {
	t.inserts = make(map[string]struct{})
	t.deletes = make(map[string]struct{})
	t.accessList = make(map[string][]byte)
}

// copy returns a deep copied tracer instance.
func (t *tracer) copy() *tracer {
	var (
		inserts    = make(map[string]struct{})
		deletes    = make(map[string]struct{})
		accessList = make(map[string][]byte)
	)
	for path := range t.inserts {
		inserts[path] = struct{}{}
	}
	for path := range t.deletes {
		deletes[path] = struct{}{}
	}
	for path, blob := range t.accessList {
		accessList[path] = common.CopyBytes(blob)
	}
	return &tracer{
		inserts:    inserts,
		deletes:    deletes,
		accessList: accessList,
	}
}

// markDeletions puts all tracked deletions into the provided nodeset.
func (t *tracer) markDeletions(set *NodeSet) {
	for path := range t.deletes {
		// It's possible a few deleted nodes were embedded
		// in their parent before, the deletions can be no
		// effect by deleting nothing, filter them out.
		if _, ok := set.accessList[path]; !ok {
			continue
		}
		set.markDeleted([]byte(path))
	}
}
