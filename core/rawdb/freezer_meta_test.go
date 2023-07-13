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

import (
	"os"
	"testing"
)

func TestReadWriteFreezerTableMeta(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create file %v", err)
	}
	err = writeMetadata(f, newMetadata(100))
	if err != nil {
		t.Fatalf("Failed to write metadata %v", err)
	}
	meta, err := readMetadata(f)
	if err != nil {
		t.Fatalf("Failed to read metadata %v", err)
	}
	if meta.Version != freezerVersion {
		t.Fatalf("Unexpected version field")
	}
	if meta.VirtualTail != uint64(100) {
		t.Fatalf("Unexpected virtual tail field")
	}
}

func TestInitializeFreezerTableMeta(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create file %v", err)
	}
	meta, err := loadMetadata(f, uint64(100))
	if err != nil {
		t.Fatalf("Failed to read metadata %v", err)
	}
	if meta.Version != freezerVersion {
		t.Fatalf("Unexpected version field")
	}
	if meta.VirtualTail != uint64(100) {
		t.Fatalf("Unexpected virtual tail field")
	}
}
