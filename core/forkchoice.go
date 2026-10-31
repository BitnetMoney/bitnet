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

package core

import (
	crand "crypto/rand"
	"errors"
	"math/big"
	mrand "math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header verification. It's implemented by both blockchain
// and lightchain.
type ChainReader interface {
	// Config retrieves the header chain's chain configuration.
	Config() *params.ChainConfig

	// GetTd returns the total difficulty of a local block.
	GetTd(common.Hash, uint64) *big.Int
}

// ForkChoice is the fork chooser based on the highest total difficulty of the
// chain(the fork choice used in the eth1) and the external fork choice (the fork
// choice used in the eth2). This main goal of this ForkChoice is not only for
// offering fork choice during the eth1/2 merge phase, but also keep the compatibility
// for all other proof-of-work networks.
type ForkChoice struct {
	chain ChainReader
	rand  *mrand.Rand

	// preserve is a helper function used in td fork choice.
	// Miners will prefer to choose the local mined block if the
	// local td is equal to the extern one. It can be nil for light
	// client
	preserve func(header *types.Header) bool

	// maxReorgDepth is how many canonical blocks a reorg may drop while the
	// local head is recent; 0 means no limit. See SetMaxReorgDepth.
	maxReorgDepth uint64
	now           func() time.Time
	onRefuse      func(current, extern *types.Header, depth uint64) // optional
}

// recentHead is how close to the wall clock the local head must be for the
// reorg depth limit to apply. A node that is still catching up may be on a
// stale branch and must be able to leave it.
const recentHead = time.Hour

// canonicalReader is implemented by chains that can look up their canonical
// chain, which the reorg depth limit needs to find the fork point.
type canonicalReader interface {
	GetHeader(common.Hash, uint64) *types.Header
	GetCanonicalHash(uint64) common.Hash
}

func NewForkChoice(chainReader ChainReader, preserve func(header *types.Header) bool) *ForkChoice {
	// Seed a fast but crypto originating random generator
	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		log.Crit("Failed to initialize random seed", "err", err)
	}
	return &ForkChoice{
		chain:    chainReader,
		rand:     mrand.New(mrand.NewSource(seed.Int64())),
		preserve: preserve,
		now:      time.Now,
	}
}

// SetMaxReorgDepth limits how many canonical blocks a reorg may drop; 0 means no
// limit. While the local head is recent, a heavier chain that forks off deeper
// than this is refused: it is kept as a side chain and an error is logged. This
// stops a miner with most of the hashrate from mining a long chain in private
// and replacing an hour of payments with it. The cost is that after a genuine
// network split that lasts longer than the limit, nodes on the lighter side do
// not follow the heavier side by themselves; an operator then lifts the limit
// or moves the head with debug_setHead.
func (f *ForkChoice) SetMaxReorgDepth(depth uint64) {
	f.maxReorgDepth = depth
}

// ReorgNeeded returns whether the reorg should be applied based on the given
// external header and local canonical chain, within the reorg depth limit.
func (f *ForkChoice) ReorgNeeded(current *types.Header, extern *types.Header) (bool, error) {
	reorg, err := f.reorgNeeded(current, extern)
	if err != nil || !reorg || f.maxReorgDepth == 0 {
		return reorg, err
	}
	if f.now().Sub(time.Unix(int64(current.Time), 0)) > recentHead {
		return true, nil
	}
	if depth, ok := f.forkDepth(current, extern); ok && depth > f.maxReorgDepth {
		log.Error("Refused deep chain reorganisation", "drop", depth, "limit", f.maxReorgDepth,
			"head", current.Number, "headhash", current.Hash(), "new", extern.Number, "newhash", extern.Hash())
		if f.onRefuse != nil {
			f.onRefuse(current, extern, depth)
		}
		return false, nil
	}
	return true, nil
}

// forkDepth returns how many blocks of the canonical chain up to current a switch
// to extern would drop. It stops walking once the depth exceeds the limit, and
// reports false if the fork point cannot be looked up.
func (f *ForkChoice) forkDepth(current, extern *types.Header) (uint64, bool) {
	chain, ok := f.chain.(canonicalReader)
	if !ok {
		return 0, false
	}
	head := current.Number.Uint64()
	h := extern
	for h.Number.Uint64() > head {
		if h = chain.GetHeader(h.ParentHash, h.Number.Uint64()-1); h == nil {
			return 0, false
		}
	}
	for chain.GetCanonicalHash(h.Number.Uint64()) != h.Hash() {
		if head-h.Number.Uint64() > f.maxReorgDepth || h.Number.Uint64() == 0 {
			return head - h.Number.Uint64(), true
		}
		if h = chain.GetHeader(h.ParentHash, h.Number.Uint64()-1); h == nil {
			return 0, false
		}
	}
	return head - h.Number.Uint64(), true
}

// reorgNeeded returns whether the reorg should be applied
// based on the given external header and local canonical chain.
// In the td mode, the new head is chosen if the corresponding
// total difficulty is higher. In the extern mode, the trusted
// header is always selected as the head.
func (f *ForkChoice) reorgNeeded(current *types.Header, extern *types.Header) (bool, error) {
	var (
		localTD  = f.chain.GetTd(current.Hash(), current.Number.Uint64())
		externTd = f.chain.GetTd(extern.Hash(), extern.Number.Uint64())
	)
	if localTD == nil || externTd == nil {
		return false, errors.New("missing td")
	}
	// Accept the new header as the chain head if the transition
	// is already triggered. We assume all the headers after the
	// transition come from the trusted consensus layer.
	if ttd := f.chain.Config().TerminalTotalDifficulty; ttd != nil && ttd.Cmp(externTd) <= 0 {
		return true, nil
	}

	// Reject the new header if it's the same as the current head.
	// Refer to https://eprint.iacr.org/2022/1020 for more info.
	// This is to prevent the uncle maker attack.
	if current.Number.Cmp(extern.Number) == 0 {
		return false, nil
	}

	// If the total difficulty is higher than our known, add it to the canonical chain
	if diff := externTd.Cmp(localTD); diff > 0 {
		return true, nil
	} else if diff < 0 {
		return false, nil
	}
	// Local and external difficulty is identical.
	// Second clause in the if statement reduces the vulnerability to selfish mining.
	// Please refer to http://www.cs.cornell.edu/~ie53/publications/btcProcFC.pdf
	reorg := false
	externNum, localNum := extern.Number.Uint64(), current.Number.Uint64()
	if externNum < localNum {
		reorg = true
	} else if externNum == localNum {
		var currentPreserve, externPreserve bool
		if f.preserve != nil {
			currentPreserve, externPreserve = f.preserve(current), f.preserve(extern)
		}
		reorg = !currentPreserve && (externPreserve || f.rand.Float64() < 0.5)
	}
	return reorg, nil
}
