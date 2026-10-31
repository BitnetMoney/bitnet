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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// reorgJournalSize is how many reorganisations the node remembers.
const reorgJournalSize = 256

// ReorgRecord describes a chain reorganisation the node made, or refused under
// the reorg depth limit.
type ReorgRecord struct {
	Time       uint64         `json:"time"`       // Unix time the node saw it
	ForkNumber uint64         `json:"forkNumber"` // last block both chains share
	Dropped    uint64         `json:"dropped"`    // canonical blocks replaced (a lower bound when refused)
	Added      uint64         `json:"added"`      // blocks of the new chain after the fork
	OldHead    common.Hash    `json:"oldHead"`
	NewHead    common.Hash    `json:"newHead"`
	NewMiner   common.Address `json:"newMiner"` // miner of the new chain's head
	Refused    bool           `json:"refused"`
}

// reorgJournal keeps the most recent reorganisations in memory, oldest first.
type reorgJournal struct {
	mu      sync.Mutex
	records []ReorgRecord
}

func (j *reorgJournal) add(r ReorgRecord) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if r.Time == 0 {
		r.Time = uint64(time.Now().Unix())
	}
	j.records = append(j.records, r)
	if len(j.records) > reorgJournalSize {
		j.records = j.records[len(j.records)-reorgJournalSize:]
	}
}

// list returns the records newest first.
func (j *reorgJournal) list() []ReorgRecord {
	j.mu.Lock()
	defer j.mu.Unlock()
	out := make([]ReorgRecord, len(j.records))
	for i, r := range j.records {
		out[len(out)-1-i] = r
	}
	return out
}

// RecentReorgs returns the chain reorganisations this node made or refused since
// it started, newest first (at most reorgJournalSize).
func (bc *BlockChain) RecentReorgs() []ReorgRecord {
	return bc.reorgs.list()
}

// recordRefusedReorg is called by the fork choice when the depth limit refuses a
// heavier chain.
func (bc *BlockChain) recordRefusedReorg(current, extern *types.Header, depth uint64) {
	fork := current.Number.Uint64() - depth
	bc.reorgs.add(ReorgRecord{
		ForkNumber: fork,
		Dropped:    depth,
		Added:      extern.Number.Uint64() - fork,
		OldHead:    current.Hash(),
		NewHead:    extern.Hash(),
		NewMiner:   extern.Coinbase,
		Refused:    true,
	})
}
