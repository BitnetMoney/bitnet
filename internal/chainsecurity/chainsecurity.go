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

// Package chainsecurity measures how concentrated mining is and how many
// confirmations a payment needs against a miner with a given share of the
// hashrate.
package chainsecurity

import (
	"bytes"
	"math"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

// Share is one miner's part of a window of blocks.
type Share struct {
	Address common.Address `json:"address"`
	Blocks  uint64         `json:"blocks"`
	Share   float64        `json:"share"`
}

// Shares counts the blocks of each miner, largest first (ties by address, so
// the order is stable).
func Shares(miners []common.Address) []Share {
	counts := make(map[common.Address]uint64)
	for _, m := range miners {
		counts[m]++
	}
	out := make([]Share, 0, len(counts))
	for addr, n := range counts {
		out = append(out, Share{Address: addr, Blocks: n, Share: float64(n) / float64(len(miners))})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Blocks != out[j].Blocks {
			return out[i].Blocks > out[j].Blocks
		}
		return bytes.Compare(out[i].Address[:], out[j].Address[:]) < 0
	})
	return out
}

// NakamotoCoefficient is the fewest miners that together produced more than
// half of the blocks, given shares sorted largest first. It is 0 when there are
// no blocks.
func NakamotoCoefficient(shares []Share) int {
	var total, sum uint64
	for _, s := range shares {
		total += s.Blocks
	}
	for i, s := range shares {
		sum += s.Blocks
		if 2*sum > total {
			return i + 1
		}
	}
	return 0
}

// AttackerSuccess is the probability that an attacker with share q of the
// hashrate ever overtakes the honest chain from z blocks behind: the formula in
// section 11 of the Bitcoin paper. It is 1 for q >= 0.5.
func AttackerSuccess(q float64, z int) float64 {
	if q >= 0.5 {
		return 1
	}
	if q <= 0 {
		return 0
	}
	p := 1 - q
	lambda := float64(z) * q / p
	poisson := math.Exp(-lambda)
	sum := 1.0
	for k := 0; k <= z; k++ {
		if k > 0 {
			poisson *= lambda / float64(k)
		}
		sum -= poisson * (1 - math.Pow(q/p, float64(z-k)))
	}
	return math.Max(sum, 0)
}

// MaxConfirmations bounds the search in Confirmations.
const MaxConfirmations = 10000

// Confirmations returns the fewest confirmations that bring an attacker with
// share q below the given probability of reversing a payment. ok is false when
// no count up to MaxConfirmations is enough, which is always the case for
// q >= 0.5: such a miner finds blocks faster than everyone else together.
func Confirmations(q, risk float64) (z int, ok bool) {
	if q >= 0.5 {
		return 0, false
	}
	for z = 1; z <= MaxConfirmations; z++ {
		if AttackerSuccess(q, z) < risk {
			return z, true
		}
	}
	return 0, false
}
