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

package core

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

func TestMaxReorgDepth(t *testing.T) {
	genesis := &Genesis{Config: params.TestChainConfig, BaseFee: big.NewInt(params.InitialBaseFee)}
	genDb, canon, _ := GenerateChainWithGenesis(genesis, ethash.NewFaker(), 20, nil)
	head := canon[len(canon)-1]

	// fork returns a heavier chain that leaves canon after block forkAt.
	fork := func(forkAt, n int) []*types.Block {
		blocks, _ := GenerateChain(genesis.Config, canon[forkAt-1], ethash.NewFaker(), genDb, n, func(i int, b *BlockGen) {
			b.SetCoinbase(common.Address{0x01})
		})
		return blocks
	}
	deep := fork(5, 30)     // drops blocks 6..20: 15 blocks
	shallow := fork(15, 10) // drops blocks 16..20: 5 blocks

	inSync := time.Unix(int64(head.Time())+60, 0)
	behind := time.Unix(int64(head.Time())+int64(2*recentHead/time.Second), 0)

	tests := []struct {
		name  string
		limit uint64
		now   time.Time
		side  []*types.Block
		want  *types.Block
	}{
		{"no limit", 0, inSync, deep, deep[len(deep)-1]},
		{"deeper than the limit", 10, inSync, deep, head},
		{"within the limit", 10, inSync, shallow, shallow[len(shallow)-1]},
		{"head not recent", 10, behind, deep, deep[len(deep)-1]},
	}
	for _, tt := range tests {
		chain, err := NewBlockChain(rawdb.NewMemoryDatabase(), nil, genesis, nil, ethash.NewFaker(), vm.Config{}, nil, nil)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		chain.SetMaxReorgDepth(tt.limit)
		chain.forker.now = func() time.Time { return tt.now }

		if _, err := chain.InsertChain(canon); err != nil {
			t.Fatalf("%s: inserting the canonical chain: %v", tt.name, err)
		}
		if _, err := chain.InsertChain(tt.side); err != nil {
			t.Fatalf("%s: inserting the side chain: %v", tt.name, err)
		}
		if got := chain.CurrentBlock(); got.Hash() != tt.want.Hash() {
			t.Errorf("%s: head = #%d %x, want #%d %x", tt.name, got.Number, got.Hash(), tt.want.Number(), tt.want.Hash())
		}
		chain.Stop()
	}
}

func TestReorgJournal(t *testing.T) {
	genesis := &Genesis{Config: params.TestChainConfig, BaseFee: big.NewInt(params.InitialBaseFee)}
	genDb, canon, _ := GenerateChainWithGenesis(genesis, ethash.NewFaker(), 20, nil)
	fork := func(forkAt, n int) []*types.Block {
		blocks, _ := GenerateChain(genesis.Config, canon[forkAt-1], ethash.NewFaker(), genDb, n, func(i int, b *BlockGen) {
			b.SetCoinbase(common.Address{0x01})
		})
		return blocks
	}
	inSync := time.Unix(int64(canon[len(canon)-1].Time())+60, 0)

	// A reorg the node makes. The side chain takes over at height 21, the first
	// height above the canonical head (Bitnet keeps the first-seen block at equal
	// height), so it replaces blocks 16-20 with 16-21 and then simply grows.
	chain, _ := NewBlockChain(rawdb.NewMemoryDatabase(), nil, genesis, nil, ethash.NewFaker(), vm.Config{}, nil, nil)
	defer chain.Stop()
	chain.InsertChain(canon)
	shallow := fork(15, 10)
	chain.InsertChain(shallow)
	got := chain.RecentReorgs()
	if len(got) != 1 {
		t.Fatalf("recorded %d reorgs, want 1", len(got))
	}
	if r := got[0]; r.Refused || r.ForkNumber != 15 || r.Dropped != 5 || r.Added != 6 ||
		r.OldHead != canon[19].Hash() || r.NewHead != shallow[5].Hash() || r.NewMiner != (common.Address{0x01}) {
		t.Errorf("unexpected record: %+v", r)
	}

	// A reorg the depth limit refuses.
	limited, _ := NewBlockChain(rawdb.NewMemoryDatabase(), nil, genesis, nil, ethash.NewFaker(), vm.Config{}, nil, nil)
	defer limited.Stop()
	limited.SetMaxReorgDepth(10)
	limited.forker.now = func() time.Time { return inSync }
	limited.InsertChain(canon)
	deep := fork(5, 30)
	limited.InsertChain(deep)
	got = limited.RecentReorgs()
	if len(got) == 0 || !got[0].Refused || got[0].Dropped <= 10 || got[0].OldHead != canon[19].Hash() {
		t.Errorf("refused reorg not recorded: %+v", got)
	}
}
