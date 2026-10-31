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

package ethash

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// Issuance must equal what the block reward code actually credits, for every
// uncle depth and on both the Bitnet and the test chain rules.
func TestIssuanceMatchesRewards(t *testing.T) {
	for _, config := range []*params.ChainConfig{params.MainnetChainConfig, params.TestChainConfig} {
		for depth := 1; depth <= 6; depth++ {
			number := big.NewInt(100)
			header := &types.Header{Number: number, Coinbase: common.Address{0xa0}}
			var uncles []*types.Header
			for i := 0; i < 2; i++ { // two uncles, the most a block may include
				uncles = append(uncles, &types.Header{
					Number:   new(big.Int).Sub(number, big.NewInt(int64(depth))),
					Coinbase: common.Address{0xb0, byte(i)},
				})
			}
			db, _ := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
			accumulateRewards(config, db, header, uncles)

			credited := new(big.Int).Set(db.GetBalance(header.Coinbase))
			for _, u := range uncles {
				credited.Add(credited, db.GetBalance(u.Coinbase))
			}
			if got := Issuance(config, header, uncles); got.Cmp(credited) != 0 {
				t.Errorf("chain %v, depth %d: Issuance %v, credited %v", config.ChainID, depth, got, credited)
			}
		}
	}
	// On Bitnet a block with no uncles creates exactly 1 BTN.
	header := &types.Header{Number: big.NewInt(1)}
	if got := Issuance(params.MainnetChainConfig, header, nil); got.Cmp(big.NewInt(1e18)) != 0 {
		t.Errorf("plain block issuance %v, want 1 BTN", got)
	}
}
