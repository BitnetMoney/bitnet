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

package rawdb

import "github.com/ethereum/go-ethereum/ethdb"

// KeyLengthIterator is a wrapper for a database iterator that ensures only key-value pairs
// with a specific key length will be returned.
type KeyLengthIterator struct {
	requiredKeyLength int
	ethdb.Iterator
}

// NewKeyLengthIterator returns a wrapped version of the iterator that will only return key-value
// pairs where keys with a specific key length will be returned.
func NewKeyLengthIterator(it ethdb.Iterator, keyLen int) ethdb.Iterator {
	return &KeyLengthIterator{
		Iterator:          it,
		requiredKeyLength: keyLen,
	}
}

func (it *KeyLengthIterator) Next() bool {
	// Return true as soon as a key with the required key length is discovered
	for it.Iterator.Next() {
		if len(it.Iterator.Key()) == it.requiredKeyLength {
			return true
		}
	}

	// Return false when we exhaust the keys in the underlying iterator.
	return false
}
