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
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

// NewStateSync create a new state trie download scheduler.
func NewStateSync(root common.Hash, database ethdb.KeyValueReader, onLeaf func(keys [][]byte, leaf []byte) error, scheme string) *trie.Sync {
	// Register the storage slot callback if the external callback is specified.
	var onSlot func(keys [][]byte, path []byte, leaf []byte, parent common.Hash, parentPath []byte) error
	if onLeaf != nil {
		onSlot = func(keys [][]byte, path []byte, leaf []byte, parent common.Hash, parentPath []byte) error {
			return onLeaf(keys, leaf)
		}
	}
	// Register the account callback to connect the state trie and the storage
	// trie belongs to the contract.
	var syncer *trie.Sync
	onAccount := func(keys [][]byte, path []byte, leaf []byte, parent common.Hash, parentPath []byte) error {
		if onLeaf != nil {
			if err := onLeaf(keys, leaf); err != nil {
				return err
			}
		}
		var obj types.StateAccount
		if err := rlp.Decode(bytes.NewReader(leaf), &obj); err != nil {
			return err
		}
		syncer.AddSubTrie(obj.Root, path, parent, parentPath, onSlot)
		syncer.AddCodeEntry(common.BytesToHash(obj.CodeHash), path, parent, parentPath)
		return nil
	}
	syncer = trie.NewSync(root, database, onAccount, scheme)
	return syncer
}
