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
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/chainsecurity"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	defaultShareWindow = 1000  // blocks counted when no window is given
	maxShareWindow     = 10000 // largest window a call may ask for
	defaultRisk        = 0.001 // chance of reversal the confirmation advice aims below
)

// BitnetAPI offers figures about the chain's security that are otherwise read
// from explorers: who mines the blocks, how many confirmations a payment needs,
// the chain reorganisations this node saw, and transactions replaced in its pool.
type BitnetAPI struct {
	chain    *core.BlockChain
	txpool   *txpool.TxPool
	issuance *issuanceIndex
}

// NewBitnetAPI creates the API for the bitnet namespace. The pool and the
// supply index may be nil; the methods that need them then return an error.
func NewBitnetAPI(chain *core.BlockChain, pool *txpool.TxPool, issuance *issuanceIndex) *BitnetAPI {
	return &BitnetAPI{chain: chain, txpool: pool, issuance: issuance}
}

// MinerShares is the result of bitnet_minerShares.
type MinerShares struct {
	FromBlock           uint64                `json:"fromBlock"`
	ToBlock             uint64                `json:"toBlock"`
	Blocks              uint64                `json:"blocks"`
	Uncles              uint64                `json:"uncles"`       // uncles included by these blocks
	AvgBlockTime        float64               `json:"avgBlockTime"` // seconds between blocks
	Miners              []chainsecurity.Share `json:"miners"`       // largest first
	NakamotoCoefficient int                   `json:"nakamotoCoefficient"`
}

// MinerShares counts who mined the newest blocks (1000 by default, at most
// 10000) and how many miners together produce more than half of them.
func (api *BitnetAPI) MinerShares(blocks *uint64) (*MinerShares, error) {
	n := uint64(defaultShareWindow)
	if blocks != nil {
		n = *blocks
	}
	if n == 0 || n > maxShareWindow {
		return nil, errors.New("blocks must be between 1 and 10000")
	}
	head := api.chain.CurrentBlock().Number.Uint64()
	if head == 0 {
		return nil, errors.New("no blocks after genesis yet")
	}
	if n > head {
		n = head
	}
	res := &MinerShares{FromBlock: head - n + 1, ToBlock: head, Blocks: n}
	miners := make([]common.Address, 0, n)
	var first, last uint64
	for num := res.FromBlock; num <= head; num++ {
		block := api.chain.GetBlockByNumber(num)
		if block == nil {
			return nil, errors.New("block missing from the local chain")
		}
		miners = append(miners, block.Coinbase())
		res.Uncles += uint64(len(block.Uncles()))
		if num == res.FromBlock {
			first = block.Time()
		}
		last = block.Time()
	}
	if n > 1 {
		res.AvgBlockTime = float64(last-first) / float64(n-1)
	}
	res.Miners = chainsecurity.Shares(miners)
	res.NakamotoCoefficient = chainsecurity.NakamotoCoefficient(res.Miners)
	return res, nil
}

// ConfirmationAdvice is how many confirmations stop one attacker.
type ConfirmationAdvice struct {
	Address       *common.Address `json:"address,omitempty"`
	Share         float64         `json:"share"`
	Safe          bool            `json:"safe"` // false: no confirmation count is enough
	Confirmations int             `json:"confirmations,omitempty"`
	Minutes       float64         `json:"minutes,omitempty"` // at the window's average block time
}

// Confirmations is the result of bitnet_confirmations.
type Confirmations struct {
	FromBlock uint64             `json:"fromBlock"`
	ToBlock   uint64             `json:"toBlock"`
	Risk      float64            `json:"risk"`
	Largest   ConfirmationAdvice `json:"largest"` // against the largest miner
	Others    ConfirmationAdvice `json:"others"`  // against every other miner (the second largest)
}

func advice(share, risk, blockTime float64) ConfirmationAdvice {
	a := ConfirmationAdvice{Share: share}
	if z, ok := chainsecurity.Confirmations(share, risk); ok {
		a.Safe, a.Confirmations, a.Minutes = true, z, float64(z)*blockTime/60
	}
	return a
}

func validRisk(risk *float64) (float64, error) {
	if risk == nil {
		return defaultRisk, nil
	}
	if *risk <= 0 || *risk >= 1 {
		return 0, errors.New("risk must be between 0 and 1")
	}
	return *risk, nil
}

// Confirmations says how many blocks to wait so that the largest miner of the
// newest blocks, and separately any other miner, has less than the given chance
// (0.1% by default) of reversing a payment. Against a miner with half of the
// hashrate or more no count is enough, and Safe is false.
func (api *BitnetAPI) Confirmations(risk *float64, blocks *uint64) (*Confirmations, error) {
	r, err := validRisk(risk)
	if err != nil {
		return nil, err
	}
	shares, err := api.MinerShares(blocks)
	if err != nil {
		return nil, err
	}
	res := &Confirmations{FromBlock: shares.FromBlock, ToBlock: shares.ToBlock, Risk: r}
	if len(shares.Miners) > 0 {
		top := shares.Miners[0]
		res.Largest = advice(top.Share, r, shares.AvgBlockTime)
		res.Largest.Address = &top.Address
	}
	if len(shares.Miners) > 1 {
		second := shares.Miners[1]
		res.Others = advice(second.Share, r, shares.AvgBlockTime)
		res.Others.Address = &second.Address
	} else {
		res.Others = advice(0, r, shares.AvgBlockTime)
	}
	return res, nil
}

// ConfirmationsFor says how many blocks to wait against an attacker with the
// given share of the hashrate (0 to 1).
func (api *BitnetAPI) ConfirmationsFor(share float64, risk *float64) (*ConfirmationAdvice, error) {
	if share < 0 || share > 1 {
		return nil, errors.New("share must be between 0 and 1")
	}
	r, err := validRisk(risk)
	if err != nil {
		return nil, err
	}
	var blockTime float64
	if shares, err := api.MinerShares(nil); err == nil {
		blockTime = shares.AvgBlockTime
	}
	a := advice(share, r, blockTime)
	return &a, nil
}

// Reorgs lists the chain reorganisations this node made, or refused under
// --reorg.maxdepth, since it started; newest first.
func (api *BitnetAPI) Reorgs() []core.ReorgRecord {
	return api.chain.RecentReorgs()
}

// ReplacedTxInfo is one side of a replacement.
type ReplacedTxInfo struct {
	Hash     common.Hash     `json:"hash"`
	To       *common.Address `json:"to"`
	Value    *hexutil.Big    `json:"value"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
}

// ReplacedTransaction is what bitnet_subscribe("replacedTransactions") sends.
type ReplacedTransaction struct {
	From  common.Address `json:"from"`
	Nonce hexutil.Uint64 `json:"nonce"`
	Old   ReplacedTxInfo `json:"old"`
	New   ReplacedTxInfo `json:"new"`
	// SamePayment is true when the new transaction has the same recipient, value
	// and data, as when a sender only raises the fee. When it is false, the old
	// payment will not happen: the sender cancelled it or sent the coins elsewhere.
	SamePayment bool `json:"samePayment"`
}

func replacedTxInfo(tx *types.Transaction) ReplacedTxInfo {
	return ReplacedTxInfo{Hash: tx.Hash(), To: tx.To(), Value: (*hexutil.Big)(tx.Value()), GasPrice: (*hexutil.Big)(tx.GasPrice())}
}

func samePayment(a, b *types.Transaction) bool {
	if (a.To() == nil) != (b.To() == nil) || (a.To() != nil && *a.To() != *b.To()) {
		return false
	}
	return a.Value().Cmp(b.Value()) == 0 && bytes.Equal(a.Data(), b.Data())
}

// NewReplacedTransaction describes a replacement for subscribers.
func NewReplacedTransaction(ev txpool.ReplacedTxEvent, signer types.Signer) *ReplacedTransaction {
	from, _ := types.Sender(signer, ev.New) // already validated by the pool
	return &ReplacedTransaction{
		From:        from,
		Nonce:       hexutil.Uint64(ev.New.Nonce()),
		Old:         replacedTxInfo(ev.Old),
		New:         replacedTxInfo(ev.New),
		SamePayment: samePayment(ev.Old, ev.New),
	}
}

// ReplacedTransactions notifies when a transaction in this node's pool is
// replaced by another from the same sender with the same nonce. A merchant who
// accepts payments before they are mined can watch for its own address as the
// old recipient with samePayment false: that payment is being taken back.
// Subscribe with bitnet_subscribe("replacedTransactions") over WebSocket or IPC.
func (api *BitnetAPI) ReplacedTransactions(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}
	if api.txpool == nil {
		return &rpc.Subscription{}, errors.New("this node has no transaction pool")
	}
	rpcSub := notifier.CreateSubscription()
	signer := types.LatestSigner(api.chain.Config())

	go func() {
		events := make(chan txpool.ReplacedTxEvent, 128)
		sub := api.txpool.SubscribeReplacedTxsEvent(events)
		defer sub.Unsubscribe()
		for {
			select {
			case ev := <-events:
				notifier.Notify(rpcSub.ID, NewReplacedTransaction(ev, signer))
			case <-rpcSub.Err():
				return
			case <-notifier.Closed():
				return
			}
		}
	}()
	return rpcSub, nil
}

// burnAddresses hold coins nobody can spend.
var burnAddresses = []common.Address{{}, common.HexToAddress("0x000000000000000000000000000000000000dEaD")}

// maxSupplyScan is how many blocks past the nearest checkpoint bitnet_supply
// reads; further than that, it asks the caller to wait for the index.
const maxSupplyScan = 30000

// Supply is the result of bitnet_supply. Amounts are exact decimals in BTN.
type Supply struct {
	Block       uint64       `json:"block"`
	Hash        common.Hash  `json:"hash"`
	Supply      string       `json:"supply"` // genesis allocation plus every block and uncle reward
	SupplyWei   *hexutil.Big `json:"supplyWei"`
	Uncles      uint64       `json:"uncles"`                // uncles included up to the block
	Burned      string       `json:"burned,omitempty"`      // at 0x0 and 0x...dEaD, if the block's state is available
	Circulating string       `json:"circulating,omitempty"` // supply minus burned
}

// formatBTN writes an amount in wei as an exact decimal in BTN.
func formatBTN(wei *big.Int) string {
	unit := big.NewInt(1e18)
	whole, frac := new(big.Int).QuoRem(wei, unit, new(big.Int))
	if frac.Sign() == 0 {
		return whole.String()
	}
	digits := strings.TrimRight(fmt.Sprintf("%018s", frac.String()), "0")
	return whole.String() + "." + digits
}

// Supply returns how many BTN exist after a block (the latest by default),
// counted from the block rewards and uncle rewards of every block. Fees only
// move coins, and genesis allocations are included. The node builds an index
// in the background; until it is ready, blocks far from the start of the chain
// return an error asking to retry.
func (api *BitnetAPI) Supply(number *rpc.BlockNumber) (*Supply, error) {
	if api.issuance == nil {
		return nil, errors.New("this node keeps no supply index")
	}
	head := api.chain.CurrentBlock()
	n := head.Number.Uint64()
	if number != nil && *number >= 0 {
		if uint64(*number) > n {
			return nil, errors.New("block not found")
		}
		n = uint64(*number)
	}
	cp, err := api.issuance.supplyAt(n, maxSupplyScan)
	if err != nil {
		return nil, err
	}
	res := &Supply{
		Block:     n,
		Hash:      cp.Hash,
		Supply:    formatBTN(cp.Supply),
		SupplyWei: (*hexutil.Big)(cp.Supply),
		Uncles:    cp.Uncles,
	}
	if header := api.chain.GetHeaderByNumber(n); header != nil {
		if state, err := api.chain.StateAt(header.Root); err == nil {
			burned := new(big.Int)
			for _, addr := range burnAddresses {
				burned.Add(burned, state.GetBalance(addr))
			}
			res.Burned = formatBTN(burned)
			res.Circulating = formatBTN(new(big.Int).Sub(cp.Supply, burned))
		}
	}
	return res, nil
}
