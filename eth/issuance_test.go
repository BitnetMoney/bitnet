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
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	dead      = common.HexToAddress("0x000000000000000000000000000000000000dEaD")
	recipient = common.Address{0x55}
)

// newIssuanceTestChain builds 40 blocks on the Bitnet rules: a genesis
// allocation, two alternating miners, uncles of depth 1 and 2, transfers that
// pay fees, and one transfer to a burn address. It returns the chain, its
// database and every account that can hold coins.
func newIssuanceTestChain(t *testing.T) (*core.BlockChain, ethdb.Database, []common.Address) {
	t.Helper()
	key, _ := crypto.GenerateKey()
	bank := crypto.PubkeyToAddress(key.PublicKey)
	config := params.MainnetChainConfig // Bitnet
	genesis := &core.Genesis{Config: config, Alloc: core.GenesisAlloc{bank: {Balance: big.NewInt(5e18)}}}
	signer := types.LatestSigner(config)
	engine := ethash.NewFaker()

	_, blocks, _ := core.GenerateChainWithGenesis(genesis, engine, 40, func(i int, b *core.BlockGen) {
		b.SetCoinbase([]common.Address{minerA, minerB}[i%2])
		if i%3 == 2 { // an uncle of depth 1
			b.AddUncle(&types.Header{ParentHash: b.PrevBlock(i - 2).Hash(), Number: big.NewInt(int64(i)), Coinbase: minerC})
		}
		if i%5 == 4 { // an uncle of depth 2
			b.AddUncle(&types.Header{ParentHash: b.PrevBlock(i - 3).Hash(), Number: big.NewInt(int64(i - 1)), Coinbase: minerC})
		}
		if i%4 == 1 {
			to := recipient
			if i == 21 {
				to = dead
			}
			b.AddTx(types.MustSignNewTx(key, signer, &types.LegacyTx{
				Nonce: b.TxNonce(bank), To: &to, Value: big.NewInt(1e15), Gas: 21000, GasPrice: big.NewInt(2e9),
			}))
		}
	})
	db := rawdb.NewMemoryDatabase()
	chain, err := core.NewBlockChain(db, nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := chain.InsertChain(blocks); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(chain.Stop)
	return chain, db, []common.Address{bank, recipient, dead, minerA, minerB, minerC}
}

// balancesAt sums the balances of the given accounts in the state after block n:
// an independent count of the supply.
func balancesAt(t *testing.T, chain *core.BlockChain, n uint64, accounts []common.Address) *big.Int {
	t.Helper()
	state, err := chain.StateAt(chain.GetHeaderByNumber(n).Root)
	if err != nil {
		t.Fatalf("state of block %d: %v", n, err)
	}
	sum := new(big.Int)
	for _, a := range accounts {
		sum.Add(sum, state.GetBalance(a))
	}
	return sum
}

func TestIssuanceIndex(t *testing.T) {
	chain, db, accounts := newIssuanceTestChain(t)
	idx := newIssuanceIndex(chain, db)
	idx.interval, idx.finality = 4, 3
	head := chain.CurrentBlock().Number.Uint64()

	if _, err := idx.supplyAt(head, 5); err != errIndexing {
		t.Fatalf("before indexing: got %v, want errIndexing", err)
	}
	if err := idx.catchUp(); err != nil {
		t.Fatal(err)
	}
	if want := (head - 3) / 4 * 4; idx.indexed != want {
		t.Fatalf("indexed up to %d, want %d", idx.indexed, want)
	}
	// check compares every block with the balances; maxScan is how far past a
	// checkpoint supplyAt may read.
	check := func(stage string, maxScan uint64) {
		t.Helper()
		var uncles uint64
		for n := uint64(0); n <= head; n++ {
			if n > 0 {
				uncles += uint64(len(chain.GetBlockByNumber(n).Uncles()))
			}
			cp, err := idx.supplyAt(n, maxScan)
			if err != nil {
				t.Fatalf("%s: block %d: %v", stage, n, err)
			}
			if want := balancesAt(t, chain, n, accounts); cp.Supply.Cmp(want) != 0 {
				t.Errorf("%s: block %d: supply %v, balances sum to %v", stage, n, cp.Supply, want)
			}
			if cp.Hash != chain.GetCanonicalHash(n) || cp.Uncles != uncles {
				t.Errorf("%s: block %d: hash %x uncles %d, want %x and %d", stage, n, cp.Hash, cp.Uncles, chain.GetCanonicalHash(n), uncles)
			}
		}
		if uncles == 0 {
			t.Fatal("the test chain has no uncles")
		}
	}
	check("indexed", 10)

	// A reorg below the top checkpoints: every checkpoint from 8 up no longer
	// matches the canonical chain. They are skipped, and catchUp rebuilds them.
	for c := uint64(8); c <= idx.indexed; c += 4 {
		bogus, _ := rlp.EncodeToBytes(&issuanceCheckpoint{Hash: common.Hash{0x01}, Supply: big.NewInt(1)})
		db.Put(issuanceKey(c), bogus)
	}
	// Without them the nearest checkpoint is block 4, more than 10 blocks back.
	if _, err := idx.supplyAt(15, 10); err != errIndexing {
		t.Errorf("past the scan limit: got %v, want errIndexing", err)
	}
	check("replaced checkpoints skipped", 100)
	if err := idx.catchUp(); err != nil {
		t.Fatal(err)
	}
	if cp, ok := idx.checkpoint(8); !ok || cp.Supply.Cmp(balancesAt(t, chain, 8, accounts)) != 0 {
		t.Errorf("checkpoint 8 not rebuilt: %+v", cp)
	}
	check("rebuilt", 10)
}

func TestBitnetSupply(t *testing.T) {
	chain, db, accounts := newIssuanceTestChain(t)
	idx := newIssuanceIndex(chain, db)
	idx.interval, idx.finality = 4, 3
	if err := idx.catchUp(); err != nil {
		t.Fatal(err)
	}
	api := NewBitnetAPI(chain, nil, idx)

	res, err := api.Supply(nil)
	if err != nil {
		t.Fatal(err)
	}
	head := chain.CurrentBlock().Number.Uint64()
	total := balancesAt(t, chain, head, accounts)
	if res.Block != head || res.Supply != formatBTN(total) || res.SupplyWei.ToInt().Cmp(total) != 0 {
		t.Errorf("supply at head: %+v, want %s BTN", res, formatBTN(total))
	}
	// 0.001 BTN went to 0x...dEaD in block 22.
	if res.Burned != "0.001" || res.Circulating != formatBTN(new(big.Int).Sub(total, big.NewInt(1e15))) {
		t.Errorf("burned %q circulating %q", res.Burned, res.Circulating)
	}
	if _, err := NewBitnetAPI(chain, nil, nil).Supply(nil); err == nil {
		t.Error("supply without an index: expected an error")
	}
}

func TestFormatBTN(t *testing.T) {
	for wei, want := range map[string]string{
		"0":                         "0",
		"1000000000000000000":       "1",
		"1500000000000000000":       "1.5",
		"1":                         "0.000000000000000001",
		"8278381031250000000000000": "8278381.03125",
	} {
		v, _ := new(big.Int).SetString(wei, 10)
		if got := formatBTN(v); got != want {
			t.Errorf("formatBTN(%s) = %s, want %s", wei, got, want)
		}
	}
}

// The background loop indexes on start and returns promptly on stop.
func TestIssuanceIndexStartStop(t *testing.T) {
	chain, db, _ := newIssuanceTestChain(t)
	idx := newIssuanceIndex(chain, db)
	idx.interval, idx.finality = 4, 3
	idx.start()
	deadline := time.Now().Add(5 * time.Second)
	for {
		idx.mu.Lock()
		indexed := idx.indexed
		idx.mu.Unlock()
		if indexed > 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the background loop did not index")
		}
		time.Sleep(10 * time.Millisecond)
	}
	stopped := make(chan struct{})
	go func() { idx.stop(); close(stopped) }()
	select {
	case <-stopped:
	case <-time.After(5 * time.Second):
		t.Fatal("stop did not return")
	}
}
