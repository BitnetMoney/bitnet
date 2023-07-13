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

package common

import (
	"testing"
)

func TestStorageSizeString(t *testing.T) {
	tests := []struct {
		size StorageSize
		str  string
	}{
		{2839274474874, "2.58 TiB"},
		{2458492810, "2.29 GiB"},
		{2381273, "2.27 MiB"},
		{2192, "2.14 KiB"},
		{12, "12.00 B"},
	}

	for _, test := range tests {
		if test.size.String() != test.str {
			t.Errorf("%f: got %q, want %q", float64(test.size), test.size.String(), test.str)
		}
	}
}

func TestStorageSizeTerminalString(t *testing.T) {
	tests := []struct {
		size StorageSize
		str  string
	}{
		{2839274474874, "2.58TiB"},
		{2458492810, "2.29GiB"},
		{2381273, "2.27MiB"},
		{2192, "2.14KiB"},
		{12, "12.00B"},
	}

	for _, test := range tests {
		if test.size.TerminalString() != test.str {
			t.Errorf("%f: got %q, want %q", float64(test.size), test.size.TerminalString(), test.str)
		}
	}
}
