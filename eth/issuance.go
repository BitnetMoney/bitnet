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

package eth

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

var issuanceKeyPrefix = []byte("bitnet-issuance-") // + block number (uint64 big endian)

// issuanceCheckpoint is the supply after a block: the genesis allocation plus
// everything blocks and uncles created up to and including it.
type issuanceCheckpoint struct {
	Hash   common.Hash // the block it was taken at, to notice reorgs
	Supply *big.Int
	Uncles uint64 // uncles included up to and including the block
}

// issuanceIndex keeps a checkpoint of the supply every interval blocks, so the
// supply after any block is a checkpoint plus at most a few intervals of block
// rewards. Checkpoints are only taken finality blocks below the head, and each
// is checked against the canonical chain before it is used.
type issuanceIndex struct {
	chain    *core.BlockChain
	db       ethdb.KeyValueStore
	interval uint64
	finality uint64

	mu      sync.Mutex
	genesis *big.Int // genesis allocation
	indexed uint64   // highest checkpoint stored and verified (0: none)

	quit chan struct{}
	wg   sync.WaitGroup
}

func newIssuanceIndex(chain *core.BlockChain, db ethdb.KeyValueStore) *issuanceIndex {
	return &issuanceIndex{chain: chain, db: db, interval: 10000, finality: 1000, quit: make(chan struct{})}
}

func issuanceKey(number uint64) []byte {
	return binary.BigEndian.AppendUint64(append([]byte{}, issuanceKeyPrefix...), number)
}

// genesisSupply sums the genesis allocation stored with the chain.
func (idx *issuanceIndex) genesisSupply() *big.Int {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if idx.genesis != nil {
		return idx.genesis
	}
	total := new(big.Int)
	if spec := rawdb.ReadGenesisStateSpec(idx.db, idx.chain.Genesis().Hash()); len(spec) > 0 {
		var alloc core.GenesisAlloc
		if err := json.Unmarshal(spec, &alloc); err != nil {
			log.Warn("Could not read the genesis allocation for the supply index", "err", err)
		}
		for _, account := range alloc {
			if account.Balance != nil {
				total.Add(total, account.Balance)
			}
		}
	}
	idx.genesis = total
	return total
}

// checkpoint returns the stored checkpoint at number if it is still on the
// canonical chain.
func (idx *issuanceIndex) checkpoint(number uint64) (*issuanceCheckpoint, bool) {
	data, err := idx.db.Get(issuanceKey(number))
	if err != nil || len(data) == 0 {
		return nil, false
	}
	var cp issuanceCheckpoint
	if err := rlp.DecodeBytes(data, &cp); err != nil || idx.chain.GetCanonicalHash(number) != cp.Hash {
		return nil, false
	}
	return &cp, true
}

// accumulate adds the issuance and uncles of blocks from..to (inclusive) to cp.
func (idx *issuanceIndex) accumulate(cp *issuanceCheckpoint, from, to uint64) error {
	config := idx.chain.Config()
	for n := from; n <= to; n++ {
		header := idx.chain.GetHeaderByNumber(n)
		if header == nil {
			return fmt.Errorf("block %d missing from the local chain", n)
		}
		var uncles []*types.Header
		if header.UncleHash != types.EmptyUncleHash {
			body := idx.chain.GetBody(header.Hash())
			if body == nil {
				return fmt.Errorf("body of block %d missing", n)
			}
			uncles = body.Uncles
		}
		cp.Supply.Add(cp.Supply, ethash.Issuance(config, header, uncles))
		cp.Uncles += uint64(len(uncles))
		cp.Hash = header.Hash()
	}
	return nil
}

// base returns the highest usable checkpoint at or below number, or the genesis.
func (idx *issuanceIndex) base(number uint64) (uint64, *issuanceCheckpoint) {
	idx.mu.Lock()
	top := idx.indexed
	idx.mu.Unlock()
	if number < top {
		top = number
	}
	for c := top / idx.interval * idx.interval; c > 0; c -= idx.interval {
		if cp, ok := idx.checkpoint(c); ok {
			return c, cp
		}
	}
	return 0, &issuanceCheckpoint{Hash: idx.chain.Genesis().Hash(), Supply: new(big.Int).Set(idx.genesisSupply())}
}

// errIndexing is returned while the index is too far behind to answer quickly.
var errIndexing = errors.New("the supply index is still being built, try again later")

// supplyAt returns the supply and uncle count after block number. It reads at
// most maxScan blocks beyond the nearest checkpoint.
func (idx *issuanceIndex) supplyAt(number uint64, maxScan uint64) (*issuanceCheckpoint, error) {
	start, cp := idx.base(number)
	if number-start > maxScan {
		return nil, errIndexing
	}
	if err := idx.accumulate(cp, start+1, number); err != nil {
		return nil, err
	}
	return cp, nil
}

// lastCheckpoint returns the highest checkpoint at or below limit that is still
// on the canonical chain, or 0 if there is none.
func (idx *issuanceIndex) lastCheckpoint(limit uint64) uint64 {
	for c := limit / idx.interval * idx.interval; c > 0; c -= idx.interval {
		if _, ok := idx.checkpoint(c); ok {
			return c
		}
	}
	return 0
}

// catchUp stores every missing checkpoint that is at least finality blocks below
// the head, starting after the highest one still on the canonical chain (so a
// deep reorg rebuilds what it replaced). It returns early when the index stops.
func (idx *issuanceIndex) catchUp() error {
	head := idx.chain.CurrentBlock().Number.Uint64()
	if head < idx.finality {
		return nil
	}
	limit := head - idx.finality

	start := idx.lastCheckpoint(limit)
	cp := &issuanceCheckpoint{Hash: idx.chain.Genesis().Hash(), Supply: new(big.Int).Set(idx.genesisSupply())}
	if start > 0 {
		cp, _ = idx.checkpoint(start)
	}
	idx.mu.Lock()
	idx.indexed = start
	idx.mu.Unlock()

	for next := start + idx.interval; next <= limit; next += idx.interval {
		select {
		case <-idx.quit:
			return nil
		default:
		}
		if err := idx.accumulate(cp, next-idx.interval+1, next); err != nil {
			return err
		}
		data, err := rlp.EncodeToBytes(cp)
		if err != nil {
			return err
		}
		if err := idx.db.Put(issuanceKey(next), data); err != nil {
			return err
		}
		idx.mu.Lock()
		idx.indexed = next
		idx.mu.Unlock()
	}
	return nil
}

func (idx *issuanceIndex) start() {
	idx.wg.Add(1)
	go func() {
		defer idx.wg.Done()
		started := time.Now()
		logged := false
		for {
			if err := idx.catchUp(); err != nil {
				log.Warn("Supply index stalled", "err", err)
			} else if !logged {
				idx.mu.Lock()
				indexed := idx.indexed
				idx.mu.Unlock()
				if indexed > 0 {
					log.Info("Supply index ready", "checkpoint", indexed, "elapsed", common.PrettyDuration(time.Since(started)))
					logged = true
				}
			}
			select {
			case <-idx.quit:
				return
			case <-time.After(time.Minute):
			}
		}
	}()
}

func (idx *issuanceIndex) stop() {
	close(idx.quit)
	idx.wg.Wait()
}
