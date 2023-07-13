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

package bloombits

import (
	"bytes"
	crand "crypto/rand"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

// Tests that batched bloom bits are correctly rotated from the input bloom
// filters.
func TestGenerator(t *testing.T) {
	// Generate the input and the rotated output
	var input, output [types.BloomBitLength][types.BloomByteLength]byte

	for i := 0; i < types.BloomBitLength; i++ {
		for j := 0; j < types.BloomBitLength; j++ {
			bit := byte(rand.Int() % 2)

			input[i][j/8] |= bit << byte(7-j%8)
			output[types.BloomBitLength-1-j][i/8] |= bit << byte(7-i%8)
		}
	}
	// Crunch the input through the generator and verify the result
	gen, err := NewGenerator(types.BloomBitLength)
	if err != nil {
		t.Fatalf("failed to create bloombit generator: %v", err)
	}
	for i, bloom := range input {
		if err := gen.AddBloom(uint(i), bloom); err != nil {
			t.Fatalf("bloom %d: failed to add: %v", i, err)
		}
	}
	for i, want := range output {
		have, err := gen.Bitset(uint(i))
		if err != nil {
			t.Fatalf("output %d: failed to retrieve bits: %v", i, err)
		}
		if !bytes.Equal(have, want[:]) {
			t.Errorf("output %d: bit vector mismatch have %x, want %x", i, have, want)
		}
	}
}

func BenchmarkGenerator(b *testing.B) {
	var input [types.BloomBitLength][types.BloomByteLength]byte
	b.Run("empty", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Crunch the input through the generator and verify the result
			gen, err := NewGenerator(types.BloomBitLength)
			if err != nil {
				b.Fatalf("failed to create bloombit generator: %v", err)
			}
			for j, bloom := range &input {
				if err := gen.AddBloom(uint(j), bloom); err != nil {
					b.Fatalf("bloom %d: failed to add: %v", i, err)
				}
			}
		}
	})
	for i := 0; i < types.BloomBitLength; i++ {
		crand.Read(input[i][:])
	}
	b.Run("random", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Crunch the input through the generator and verify the result
			gen, err := NewGenerator(types.BloomBitLength)
			if err != nil {
				b.Fatalf("failed to create bloombit generator: %v", err)
			}
			for j, bloom := range &input {
				if err := gen.AddBloom(uint(j), bloom); err != nil {
					b.Fatalf("bloom %d: failed to add: %v", i, err)
				}
			}
		}
	})
}
