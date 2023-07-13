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

package bitutil

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common/bitutil"
)

// Fuzz implements a go-fuzz fuzzer method to test various encoding method
// invocations.
func Fuzz(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	if data[0]%2 == 0 {
		return fuzzEncode(data[1:])
	}
	return fuzzDecode(data[1:])
}

// fuzzEncode implements a go-fuzz fuzzer method to test the bitset encoding and
// decoding algorithm.
func fuzzEncode(data []byte) int {
	proc, _ := bitutil.DecompressBytes(bitutil.CompressBytes(data), len(data))
	if !bytes.Equal(data, proc) {
		panic("content mismatch")
	}
	return 1
}

// fuzzDecode implements a go-fuzz fuzzer method to test the bit decoding and
// reencoding algorithm.
func fuzzDecode(data []byte) int {
	blob, err := bitutil.DecompressBytes(data, 1024)
	if err != nil {
		return 0
	}
	// re-compress it (it's OK if the re-compressed differs from the
	// original - the first input may not have been compressed at all)
	comp := bitutil.CompressBytes(blob)
	if len(comp) > len(blob) {
		// After compression, it must be smaller or equal
		panic("bad compression")
	}
	// But decompressing it once again should work
	decomp, err := bitutil.DecompressBytes(data, 1024)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(decomp, blob) {
		panic("content mismatch")
	}
	return 1
}
