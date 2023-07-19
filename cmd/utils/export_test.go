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

package utils

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/rlp"
)

// TestExport does basic sanity checks on the export/import functionality
func TestExport(t *testing.T) {
	f := fmt.Sprintf("%v/tempdump", os.TempDir())
	defer func() {
		os.Remove(f)
	}()
	testExport(t, f)
}

func TestExportGzip(t *testing.T) {
	f := fmt.Sprintf("%v/tempdump.gz", os.TempDir())
	defer func() {
		os.Remove(f)
	}()
	testExport(t, f)
}

type testIterator struct {
	index int
}

func newTestIterator() *testIterator {
	return &testIterator{index: -1}
}

func (iter *testIterator) Next() (byte, []byte, []byte, bool) {
	if iter.index >= 999 {
		return 0, nil, nil, false
	}
	iter.index += 1
	if iter.index == 42 {
		iter.index += 1
	}
	return OpBatchAdd, []byte(fmt.Sprintf("key-%04d", iter.index)),
		[]byte(fmt.Sprintf("value %d", iter.index)), true
}

func (iter *testIterator) Release() {}

func testExport(t *testing.T, f string) {
	err := ExportChaindata(f, "testdata", newTestIterator(), make(chan struct{}))
	if err != nil {
		t.Fatal(err)
	}
	db := rawdb.NewMemoryDatabase()
	err = ImportLDBData(db, f, 5, make(chan struct{}))
	if err != nil {
		t.Fatal(err)
	}
	// verify
	for i := 0; i < 1000; i++ {
		v, err := db.Get([]byte(fmt.Sprintf("key-%04d", i)))
		if (i < 5 || i == 42) && err == nil {
			t.Fatalf("expected no element at idx %d, got '%v'", i, string(v))
		}
		if !(i < 5 || i == 42) {
			if err != nil {
				t.Fatalf("expected element idx %d: %v", i, err)
			}
			if have, want := string(v), fmt.Sprintf("value %d", i); have != want {
				t.Fatalf("have %v, want %v", have, want)
			}
		}
	}
	v, err := db.Get([]byte(fmt.Sprintf("key-%04d", 1000)))
	if err == nil {
		t.Fatalf("expected no element at idx %d, got '%v'", 1000, string(v))
	}
}

// testDeletion tests if the deletion markers can be exported/imported correctly
func TestDeletionExport(t *testing.T) {
	f := fmt.Sprintf("%v/tempdump", os.TempDir())
	defer func() {
		os.Remove(f)
	}()
	testDeletion(t, f)
}

// TestDeletionExportGzip tests if the deletion markers can be exported/imported
// correctly with gz compression.
func TestDeletionExportGzip(t *testing.T) {
	f := fmt.Sprintf("%v/tempdump.gz", os.TempDir())
	defer func() {
		os.Remove(f)
	}()
	testDeletion(t, f)
}

type deletionIterator struct {
	index int
}

func newDeletionIterator() *deletionIterator {
	return &deletionIterator{index: -1}
}

func (iter *deletionIterator) Next() (byte, []byte, []byte, bool) {
	if iter.index >= 999 {
		return 0, nil, nil, false
	}
	iter.index += 1
	if iter.index == 42 {
		iter.index += 1
	}
	return OpBatchDel, []byte(fmt.Sprintf("key-%04d", iter.index)), nil, true
}

func (iter *deletionIterator) Release() {}

func testDeletion(t *testing.T, f string) {
	err := ExportChaindata(f, "testdata", newDeletionIterator(), make(chan struct{}))
	if err != nil {
		t.Fatal(err)
	}
	db := rawdb.NewMemoryDatabase()
	for i := 0; i < 1000; i++ {
		db.Put([]byte(fmt.Sprintf("key-%04d", i)), []byte(fmt.Sprintf("value %d", i)))
	}
	err = ImportLDBData(db, f, 5, make(chan struct{}))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		v, err := db.Get([]byte(fmt.Sprintf("key-%04d", i)))
		if i < 5 || i == 42 {
			if err != nil {
				t.Fatalf("expected element at idx %d, got '%v'", i, err)
			}
			if have, want := string(v), fmt.Sprintf("value %d", i); have != want {
				t.Fatalf("have %v, want %v", have, want)
			}
		}
		if !(i < 5 || i == 42) {
			if err == nil {
				t.Fatalf("expected no element idx %d: %v", i, string(v))
			}
		}
	}
}

// TestImportFutureFormat tests that we reject unsupported future versions.
func TestImportFutureFormat(t *testing.T) {
	f := fmt.Sprintf("%v/tempdump-future", os.TempDir())
	defer func() {
		os.Remove(f)
	}()
	fh, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if err := rlp.Encode(fh, &exportHeader{
		Magic:    exportMagic,
		Version:  500,
		Kind:     "testdata",
		UnixTime: uint64(time.Now().Unix()),
	}); err != nil {
		t.Fatal(err)
	}
	db2 := rawdb.NewMemoryDatabase()
	err = ImportLDBData(db2, f, 0, make(chan struct{}))
	if err == nil {
		t.Fatal("Expected error, got none")
	}
	if !strings.HasPrefix(err.Error(), "incompatible version") {
		t.Fatalf("wrong error: %v", err)
	}
}
