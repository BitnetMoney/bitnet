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
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

func filledStateDB() *StateDB {
	state, _ := New(common.Hash{}, NewDatabase(rawdb.NewMemoryDatabase()), nil)

	// Create an account and check if the retrieved balance is correct
	addr := common.HexToAddress("0xaffeaffeaffeaffeaffeaffeaffeaffeaffeaffe")
	skey := common.HexToHash("aaa")
	sval := common.HexToHash("bbb")

	state.SetBalance(addr, big.NewInt(42)) // Change the account trie
	state.SetCode(addr, []byte("hello"))   // Change an external metadata
	state.SetState(addr, skey, sval)       // Change the storage trie
	for i := 0; i < 100; i++ {
		sk := common.BigToHash(big.NewInt(int64(i)))
		state.SetState(addr, sk, sk) // Change the storage trie
	}
	return state
}

func TestCopyAndClose(t *testing.T) {
	db := filledStateDB()
	prefetcher := newTriePrefetcher(db.db, db.originalRoot, "")
	skey := common.HexToHash("aaa")
	prefetcher.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()})
	prefetcher.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()})
	time.Sleep(1 * time.Second)
	a := prefetcher.trie(common.Hash{}, db.originalRoot)
	prefetcher.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()})
	b := prefetcher.trie(common.Hash{}, db.originalRoot)
	cpy := prefetcher.copy()
	cpy.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()})
	cpy.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()})
	c := cpy.trie(common.Hash{}, db.originalRoot)
	prefetcher.close()
	cpy2 := cpy.copy()
	cpy2.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()})
	d := cpy2.trie(common.Hash{}, db.originalRoot)
	cpy.close()
	cpy2.close()
	if a.Hash() != b.Hash() || a.Hash() != c.Hash() || a.Hash() != d.Hash() {
		t.Fatalf("Invalid trie, hashes should be equal: %v %v %v %v", a.Hash(), b.Hash(), c.Hash(), d.Hash())
	}
}

func TestUseAfterClose(t *testing.T) {
	db := filledStateDB()
	prefetcher := newTriePrefetcher(db.db, db.originalRoot, "")
	skey := common.HexToHash("aaa")
	prefetcher.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()})
	a := prefetcher.trie(common.Hash{}, db.originalRoot)
	prefetcher.close()
	b := prefetcher.trie(common.Hash{}, db.originalRoot)
	if a == nil {
		t.Fatal("Prefetching before close should not return nil")
	}
	if b != nil {
		t.Fatal("Trie after close should return nil")
	}
}

func TestCopyClose(t *testing.T) {
	db := filledStateDB()
	prefetcher := newTriePrefetcher(db.db, db.originalRoot, "")
	skey := common.HexToHash("aaa")
	prefetcher.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()})
	cpy := prefetcher.copy()
	a := prefetcher.trie(common.Hash{}, db.originalRoot)
	b := cpy.trie(common.Hash{}, db.originalRoot)
	prefetcher.close()
	c := prefetcher.trie(common.Hash{}, db.originalRoot)
	d := cpy.trie(common.Hash{}, db.originalRoot)
	if a == nil {
		t.Fatal("Prefetching before close should not return nil")
	}
	if b == nil {
		t.Fatal("Copy trie should return nil")
	}
	if c != nil {
		t.Fatal("Trie after close should return nil")
	}
	if d == nil {
		t.Fatal("Copy trie should not return nil")
	}
}
