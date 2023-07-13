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

package state

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

// Tests that the node iterator indeed walks over the entire database contents.
func TestNodeIteratorCoverage(t *testing.T) {
	// Create some arbitrary test state to iterate
	db, sdb, root, _ := makeTestState()
	sdb.TrieDB().Commit(root, false)

	state, err := New(root, sdb, nil)
	if err != nil {
		t.Fatalf("failed to create state trie at %x: %v", root, err)
	}
	// Gather all the node hashes found by the iterator
	hashes := make(map[common.Hash]struct{})
	for it := NewNodeIterator(state); it.Next(); {
		if it.Hash != (common.Hash{}) {
			hashes[it.Hash] = struct{}{}
		}
	}
	// Check in-disk nodes
	var (
		seenNodes = make(map[common.Hash]struct{})
		seenCodes = make(map[common.Hash]struct{})
	)
	it := db.NewIterator(nil, nil)
	for it.Next() {
		ok, hash := isTrieNode(sdb.TrieDB().Scheme(), it.Key(), it.Value())
		if !ok {
			continue
		}
		seenNodes[hash] = struct{}{}
	}
	it.Release()

	// Check in-disk codes
	it = db.NewIterator(nil, nil)
	for it.Next() {
		ok, hash := rawdb.IsCodeKey(it.Key())
		if !ok {
			continue
		}
		if _, ok := hashes[common.BytesToHash(hash)]; !ok {
			t.Errorf("state entry not reported %x", it.Key())
		}
		seenCodes[common.BytesToHash(hash)] = struct{}{}
	}
	it.Release()

	// Cross check the iterated hashes and the database/nodepool content
	for hash := range hashes {
		_, ok := seenNodes[hash]
		if !ok {
			_, ok = seenCodes[hash]
		}
		if !ok {
			t.Errorf("failed to retrieve reported node %x", hash)
		}
	}
}

// isTrieNode is a helper function which reports if the provided
// database entry belongs to a trie node or not.
func isTrieNode(scheme string, key, val []byte) (bool, common.Hash) {
	if scheme == rawdb.HashScheme {
		if len(key) == common.HashLength {
			return true, common.BytesToHash(key)
		}
	}
	return false, common.Hash{}
}
