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
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

var (
	minerA = common.Address{0xaa}
	minerB = common.Address{0xbb}
	minerC = common.Address{0xcc}
)

// newBitnetTestChain builds ten blocks: 1-8 mined by A, 9-10 by B, and block 6
// includes an uncle mined by C.
func newBitnetTestChain(t *testing.T) *core.BlockChain {
	t.Helper()
	genesis := &core.Genesis{Config: params.TestChainConfig, BaseFee: big.NewInt(params.InitialBaseFee)}
	engine := ethash.NewFaker()
	_, blocks, _ := core.GenerateChainWithGenesis(genesis, engine, 10, func(i int, b *core.BlockGen) {
		if i < 8 {
			b.SetCoinbase(minerA)
		} else {
			b.SetCoinbase(minerB)
		}
		if i == 5 { // block 6: an uncle at height 5, a sibling of block 5
			b.AddUncle(&types.Header{ParentHash: b.PrevBlock(3).Hash(), Number: big.NewInt(5), Coinbase: minerC})
		}
	})
	chain, err := core.NewBlockChain(rawdb.NewMemoryDatabase(), nil, genesis, nil, engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := chain.InsertChain(blocks); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(chain.Stop)
	return chain
}

func TestBitnetMinerShares(t *testing.T) {
	api := NewBitnetAPI(newBitnetTestChain(t), nil, nil)

	res, err := api.MinerShares(nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.FromBlock != 1 || res.ToBlock != 10 || res.Blocks != 10 || res.Uncles != 1 {
		t.Errorf("window %d-%d, %d blocks, %d uncles; want 1-10, 10 blocks, 1 uncle", res.FromBlock, res.ToBlock, res.Blocks, res.Uncles)
	}
	if res.AvgBlockTime != 10 {
		t.Errorf("average block time %v, want 10", res.AvgBlockTime)
	}
	if len(res.Miners) != 2 || res.Miners[0].Address != minerA || res.Miners[0].Blocks != 8 || res.Miners[1].Address != minerB {
		t.Errorf("unexpected miners: %+v", res.Miners)
	}
	if res.NakamotoCoefficient != 1 {
		t.Errorf("Nakamoto coefficient %d, want 1", res.NakamotoCoefficient)
	}

	three := uint64(3)
	res, err = api.MinerShares(&three)
	if err != nil {
		t.Fatal(err)
	}
	if res.FromBlock != 8 || res.Miners[0].Address != minerB || res.Miners[0].Blocks != 2 {
		t.Errorf("last three blocks: from %d, miners %+v; want from 8 with B first", res.FromBlock, res.Miners)
	}
	for _, bad := range []uint64{0, maxShareWindow + 1} {
		if _, err := api.MinerShares(&bad); err == nil {
			t.Errorf("window %d: expected an error", bad)
		}
	}
}

func TestBitnetConfirmations(t *testing.T) {
	api := NewBitnetAPI(newBitnetTestChain(t), nil, nil)

	res, err := api.Confirmations(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// A mined 80% of the blocks: no confirmation count protects against it.
	if res.Largest.Safe || res.Largest.Address == nil || *res.Largest.Address != minerA {
		t.Errorf("largest miner advice: %+v", res.Largest)
	}
	// Against B's 20%, the Bitcoin paper's table gives 11 blocks; at 10 s each.
	if !res.Others.Safe || res.Others.Confirmations != 11 || math.Abs(res.Others.Minutes-11*10/60.0) > 1e-9 {
		t.Errorf("others advice: %+v", res.Others)
	}

	one := 0.1
	adv, err := api.ConfirmationsFor(one, nil)
	if err != nil || !adv.Safe || adv.Confirmations != 5 {
		t.Errorf("10%% attacker: %+v, %v; want 5 confirmations", adv, err)
	}
	bad := 1.5
	if _, err := api.Confirmations(&bad, nil); err == nil {
		t.Error("risk 1.5: expected an error")
	}
	if reorgs := api.Reorgs(); len(reorgs) != 0 {
		t.Errorf("reorgs on a straight chain: %+v", reorgs)
	}
}

func TestReplacedTransactionSamePayment(t *testing.T) {
	key, _ := crypto.GenerateKey()
	signer := types.LatestSigner(params.TestChainConfig)
	merchant, other := common.Address{0x11}, common.Address{0x22}
	tx := func(to common.Address, value, gasPrice int64) *types.Transaction {
		return types.MustSignNewTx(key, signer, &types.LegacyTx{Nonce: 7, To: &to, Value: big.NewInt(value), Gas: 21000, GasPrice: big.NewInt(gasPrice)})
	}
	pay := tx(merchant, 100, 10)

	speedUp := NewReplacedTransaction(txpool.ReplacedTxEvent{Old: pay, New: tx(merchant, 100, 20)}, signer)
	if !speedUp.SamePayment || speedUp.From != crypto.PubkeyToAddress(key.PublicKey) || speedUp.Nonce != 7 {
		t.Errorf("fee bump: %+v", speedUp)
	}
	redirect := NewReplacedTransaction(txpool.ReplacedTxEvent{Old: pay, New: tx(other, 100, 20)}, signer)
	if redirect.SamePayment || *redirect.Old.To != merchant || *redirect.New.To != other {
		t.Errorf("payment sent elsewhere: %+v", redirect)
	}
	if lower := NewReplacedTransaction(txpool.ReplacedTxEvent{Old: pay, New: tx(merchant, 1, 20)}, signer); lower.SamePayment {
		t.Error("a smaller payment to the same address must not count as the same payment")
	}
}
