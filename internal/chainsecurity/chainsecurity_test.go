// Copyright 2026 The Bitnet community authors
// This file is part of the Bitnet library.
//
// The Bitnet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Bitnet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

package chainsecurity

import (
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// The table in section 11 of the Bitcoin paper: confirmations for P < 0.001.
func TestConfirmationsMatchBitcoinPaper(t *testing.T) {
	table := []struct {
		q float64
		z int
	}{
		{0.10, 5}, {0.15, 8}, {0.20, 11}, {0.25, 15},
		{0.30, 24}, {0.35, 41}, {0.40, 89}, {0.45, 340},
	}
	for _, row := range table {
		z, ok := Confirmations(row.q, 0.001)
		if !ok || z != row.z {
			t.Errorf("q=%.2f: got %d (ok=%v), want %d", row.q, z, ok, row.z)
		}
	}
}

// The paper also lists the probabilities themselves for q=0.1.
func TestAttackerSuccessValues(t *testing.T) {
	want := map[int]float64{0: 1.0000000, 1: 0.2045873, 5: 0.0009137, 10: 0.0000012}
	for z, p := range want {
		if got := AttackerSuccess(0.1, z); math.Abs(got-p) > 1e-7 {
			t.Errorf("z=%d: got %.7f, want %.7f", z, got, p)
		}
	}
	if got := AttackerSuccess(0.5, 100); got != 1 {
		t.Errorf("q=0.5: got %v, want 1", got)
	}
	if _, ok := Confirmations(0.69, 0.001); ok {
		t.Error("a 69% miner must have no safe confirmation count")
	}
}

func TestSharesAndNakamoto(t *testing.T) {
	a, b, c := common.Address{1}, common.Address{2}, common.Address{3}
	var miners []common.Address
	for i := 0; i < 69; i++ {
		miners = append(miners, a)
	}
	for i := 0; i < 30; i++ {
		miners = append(miners, b)
	}
	miners = append(miners, c)

	shares := Shares(miners)
	if len(shares) != 3 || shares[0].Address != a || shares[0].Blocks != 69 || shares[1].Address != b || shares[2].Address != c {
		t.Fatalf("unexpected shares: %+v", shares)
	}
	if math.Abs(shares[0].Share-0.69) > 1e-9 {
		t.Errorf("share = %v, want 0.69", shares[0].Share)
	}
	if n := NakamotoCoefficient(shares); n != 1 {
		t.Errorf("Nakamoto coefficient = %d, want 1", n)
	}
	// Three equal miners: two are needed for a majority.
	even := Shares([]common.Address{a, b, c})
	if n := NakamotoCoefficient(even); n != 2 {
		t.Errorf("Nakamoto coefficient of three equal miners = %d, want 2", n)
	}
	if n := NakamotoCoefficient(nil); n != 0 {
		t.Errorf("Nakamoto coefficient of no blocks = %d, want 0", n)
	}
}
