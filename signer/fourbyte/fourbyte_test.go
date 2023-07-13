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

package fourbyte

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Tests that all the selectors contained in the 4byte database are valid.
func TestEmbeddedDatabase(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for id, selector := range db.embedded {
		abistring, err := parseSelector(selector)
		if err != nil {
			t.Errorf("Failed to convert selector to ABI: %v", err)
			continue
		}
		abistruct, err := abi.JSON(strings.NewReader(string(abistring)))
		if err != nil {
			t.Errorf("Failed to parse ABI: %v", err)
			continue
		}
		m, err := abistruct.MethodById(common.Hex2Bytes(id))
		if err != nil {
			t.Errorf("Failed to get method by id (%s): %v", id, err)
			continue
		}
		if m.Sig != selector {
			t.Errorf("Selector mismatch: have %v, want %v", m.Sig, selector)
		}
	}
}

// Tests that custom 4byte datasets can be handled too.
func TestCustomDatabase(t *testing.T) {
	// Create a new custom 4byte database with no embedded component
	tmpdir := t.TempDir()
	filename := fmt.Sprintf("%s/4byte_custom.json", tmpdir)

	db, err := NewWithFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	db.embedded = make(map[string]string)

	// Ensure the database is empty, insert and verify
	calldata := common.Hex2Bytes("a52c101edeadbeef")
	if _, err = db.Selector(calldata); err == nil {
		t.Fatalf("Should not find a match on empty database")
	}
	if err = db.AddSelector("send(uint256)", calldata); err != nil {
		t.Fatalf("Failed to save file: %v", err)
	}
	if _, err = db.Selector(calldata); err != nil {
		t.Fatalf("Failed to find a match for abi signature: %v", err)
	}
	// Check that the file as persisted to disk by creating a new instance
	db2, err := NewFromFile(filename)
	if err != nil {
		t.Fatalf("Failed to create new abidb: %v", err)
	}
	if _, err = db2.Selector(calldata); err != nil {
		t.Fatalf("Failed to find a match for persisted abi signature: %v", err)
	}
}
