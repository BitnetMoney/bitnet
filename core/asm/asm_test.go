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

package asm

import (
	"testing"

	"encoding/hex"
)

// Tests disassembling the instructions for valid evm code
func TestInstructionIteratorValid(t *testing.T) {
	cnt := 0
	script, _ := hex.DecodeString("61000000")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if err := it.Error(); err != nil {
		t.Errorf("Expected 2, but encountered error %v instead.", err)
	}
	if cnt != 2 {
		t.Errorf("Expected 2, but got %v instead.", cnt)
	}
}

// Tests disassembling the instructions for invalid evm code
func TestInstructionIteratorInvalid(t *testing.T) {
	cnt := 0
	script, _ := hex.DecodeString("6100")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if it.Error() == nil {
		t.Errorf("Expected an error, but got %v instead.", cnt)
	}
}

// Tests disassembling the instructions for empty evm code
func TestInstructionIteratorEmpty(t *testing.T) {
	cnt := 0
	script, _ := hex.DecodeString("")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if err := it.Error(); err != nil {
		t.Errorf("Expected 0, but encountered error %v instead.", err)
	}
	if cnt != 0 {
		t.Errorf("Expected 0, but got %v instead.", cnt)
	}
}
