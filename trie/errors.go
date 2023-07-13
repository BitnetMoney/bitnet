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

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// MissingNodeError is returned by the trie functions (Get, Update, Delete)
// in the case where a trie node is not present in the local database. It contains
// information necessary for retrieving the missing node.
type MissingNodeError struct {
	Owner    common.Hash // owner of the trie if it's 2-layered trie
	NodeHash common.Hash // hash of the missing node
	Path     []byte      // hex-encoded path to the missing node
	err      error       // concrete error for missing trie node
}

// Unwrap returns the concrete error for missing trie node which
// allows us for further analysis outside.
func (err *MissingNodeError) Unwrap() error {
	return err.err
}

func (err *MissingNodeError) Error() string {
	if err.Owner == (common.Hash{}) {
		return fmt.Sprintf("missing trie node %x (path %x) %v", err.NodeHash, err.Path, err.err)
	}
	return fmt.Sprintf("missing trie node %x (owner %x) (path %x) %v", err.NodeHash, err.Owner, err.Path, err.err)
}
