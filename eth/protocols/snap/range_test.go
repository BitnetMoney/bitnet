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

package snap

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// Tests that given a starting hash and a density, the hash ranger can correctly
// split up the remaining hash space into a fixed number of chunks.
func TestHashRanges(t *testing.T) {
	tests := []struct {
		head   common.Hash
		chunks uint64
		starts []common.Hash
		ends   []common.Hash
	}{
		// Simple test case to split the entire hash range into 4 chunks
		{
			head:   common.Hash{},
			chunks: 4,
			starts: []common.Hash{
				{},
				common.HexToHash("0x4000000000000000000000000000000000000000000000000000000000000000"),
				common.HexToHash("0x8000000000000000000000000000000000000000000000000000000000000000"),
				common.HexToHash("0xc000000000000000000000000000000000000000000000000000000000000000"),
			},
			ends: []common.Hash{
				common.HexToHash("0x3fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
				common.HexToHash("0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
				common.HexToHash("0xbfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
				common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			},
		},
		// Split a divisible part of the hash range up into 2 chunks
		{
			head:   common.HexToHash("0x2000000000000000000000000000000000000000000000000000000000000000"),
			chunks: 2,
			starts: []common.Hash{
				{},
				common.HexToHash("0x9000000000000000000000000000000000000000000000000000000000000000"),
			},
			ends: []common.Hash{
				common.HexToHash("0x8fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
				common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			},
		},
		// Split the entire hash range into a non divisible 3 chunks
		{
			head:   common.Hash{},
			chunks: 3,
			starts: []common.Hash{
				{},
				common.HexToHash("0x5555555555555555555555555555555555555555555555555555555555555556"),
				common.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac"),
			},
			ends: []common.Hash{
				common.HexToHash("0x5555555555555555555555555555555555555555555555555555555555555555"),
				common.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab"),
				common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			},
		},
		// Split a part of hash range into a non divisible 3 chunks
		{
			head:   common.HexToHash("0x2000000000000000000000000000000000000000000000000000000000000000"),
			chunks: 3,
			starts: []common.Hash{
				{},
				common.HexToHash("0x6aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab"),
				common.HexToHash("0xb555555555555555555555555555555555555555555555555555555555555556"),
			},
			ends: []common.Hash{
				common.HexToHash("0x6aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				common.HexToHash("0xb555555555555555555555555555555555555555555555555555555555555555"),
				common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			},
		},
		// Split a part of hash range into a non divisible 3 chunks, but with a
		// meaningful space size for manual verification.
		//   - The head being 0xff...f0, we have 14 hashes left in the space
		//   - Chunking up 14 into 3 pieces is 4.(6), but we need the ceil of 5 to avoid a micro-last-chunk
		//   - Since the range is not divisible, the last interval will be shorter, capped at 0xff...f
		//   - The chunk ranges thus needs to be [..0, ..5], [..6, ..b], [..c, ..f]
		{
			head:   common.HexToHash("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0"),
			chunks: 3,
			starts: []common.Hash{
				{},
				common.HexToHash("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6"),
				common.HexToHash("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc"),
			},
			ends: []common.Hash{
				common.HexToHash("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5"),
				common.HexToHash("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffb"),
				common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			},
		},
	}
	for i, tt := range tests {
		r := newHashRange(tt.head, tt.chunks)

		var (
			starts = []common.Hash{{}}
			ends   = []common.Hash{r.End()}
		)
		for r.Next() {
			starts = append(starts, r.Start())
			ends = append(ends, r.End())
		}
		if len(starts) != len(tt.starts) {
			t.Errorf("test %d: starts count mismatch: have %d, want %d", i, len(starts), len(tt.starts))
		}
		for j := 0; j < len(starts) && j < len(tt.starts); j++ {
			if starts[j] != tt.starts[j] {
				t.Errorf("test %d, start %d: hash mismatch: have %x, want %x", i, j, starts[j], tt.starts[j])
			}
		}
		if len(ends) != len(tt.ends) {
			t.Errorf("test %d: ends count mismatch: have %d, want %d", i, len(ends), len(tt.ends))
		}
		for j := 0; j < len(ends) && j < len(tt.ends); j++ {
			if ends[j] != tt.ends[j] {
				t.Errorf("test %d, end %d: hash mismatch: have %x, want %x", i, j, ends[j], tt.ends[j])
			}
		}
	}
}
