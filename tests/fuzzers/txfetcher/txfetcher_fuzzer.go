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

package txfetcher

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/fetcher"
)

var (
	peers []string
	txs   []*types.Transaction
)

func init() {
	// Random is nice, but we need it deterministic
	rand := rand.New(rand.NewSource(0x3a29))

	peers = make([]string, 10)
	for i := 0; i < len(peers); i++ {
		peers[i] = fmt.Sprintf("Peer #%d", i)
	}
	txs = make([]*types.Transaction, 65536) // We need to bump enough to hit all the limits
	for i := 0; i < len(txs); i++ {
		txs[i] = types.NewTransaction(rand.Uint64(), common.Address{byte(rand.Intn(256))}, new(big.Int), 0, new(big.Int), nil)
	}
}

func Fuzz(input []byte) int {
	// Don't generate insanely large test cases, not much value in them
	if len(input) > 16*1024 {
		return 0
	}
	verbose := false
	r := bytes.NewReader(input)

	// Reduce the problem space for certain fuzz runs. Small tx space is better
	// for testing clashes and in general the fetcher, but we should still run
	// some tests with large spaces to hit potential issues on limits.
	limit, err := r.ReadByte()
	if err != nil {
		return 0
	}
	switch limit % 4 {
	case 0:
		txs = txs[:4]
	case 1:
		txs = txs[:256]
	case 2:
		txs = txs[:4096]
	case 3:
		// Full run
	}
	// Create a fetcher and hook into it's simulated fields
	clock := new(mclock.Simulated)
	rand := rand.New(rand.NewSource(0x3a29)) // Same used in package tests!!!

	f := fetcher.NewTxFetcherForTests(
		func(common.Hash) bool { return false },
		func(txs []*types.Transaction) []error {
			return make([]error, len(txs))
		},
		func(string, []common.Hash) error { return nil },
		clock, rand,
	)
	f.Start()
	defer f.Stop()

	// Try to throw random junk at the fetcher
	for {
		// Read the next command and abort if we're done
		cmd, err := r.ReadByte()
		if err != nil {
			return 0
		}
		switch cmd % 4 {
		case 0:
			// Notify a new set of transactions:
			//   Byte 1:             Peer index to announce with
			//   Byte 2:             Number of hashes to announce
			//   Byte 3-4, 5-6, etc: Transaction indices (2 byte) to announce
			peerIdx, err := r.ReadByte()
			if err != nil {
				return 0
			}
			peer := peers[int(peerIdx)%len(peers)]

			announceCnt, err := r.ReadByte()
			if err != nil {
				return 0
			}
			announce := int(announceCnt) % (2 * len(txs)) // No point in generating too many duplicates

			var (
				announceIdxs = make([]int, announce)
				announces    = make([]common.Hash, announce)
			)
			for i := 0; i < len(announces); i++ {
				annBuf := make([]byte, 2)
				if n, err := r.Read(annBuf); err != nil || n != 2 {
					return 0
				}
				announceIdxs[i] = (int(annBuf[0])*256 + int(annBuf[1])) % len(txs)
				announces[i] = txs[announceIdxs[i]].Hash()
			}
			if verbose {
				fmt.Println("Notify", peer, announceIdxs)
			}
			if err := f.Notify(peer, announces); err != nil {
				panic(err)
			}

		case 1:
			// Deliver a new set of transactions:
			//   Byte 1:             Peer index to announce with
			//   Byte 2:             Number of hashes to announce
			//   Byte 3-4, 5-6, etc: Transaction indices (2 byte) to announce
			peerIdx, err := r.ReadByte()
			if err != nil {
				return 0
			}
			peer := peers[int(peerIdx)%len(peers)]

			deliverCnt, err := r.ReadByte()
			if err != nil {
				return 0
			}
			deliver := int(deliverCnt) % (2 * len(txs)) // No point in generating too many duplicates

			var (
				deliverIdxs = make([]int, deliver)
				deliveries  = make([]*types.Transaction, deliver)
			)
			for i := 0; i < len(deliveries); i++ {
				deliverBuf := make([]byte, 2)
				if n, err := r.Read(deliverBuf); err != nil || n != 2 {
					return 0
				}
				deliverIdxs[i] = (int(deliverBuf[0])*256 + int(deliverBuf[1])) % len(txs)
				deliveries[i] = txs[deliverIdxs[i]]
			}
			directFlag, err := r.ReadByte()
			if err != nil {
				return 0
			}
			direct := (directFlag % 2) == 0
			if verbose {
				fmt.Println("Enqueue", peer, deliverIdxs, direct)
			}
			if err := f.Enqueue(peer, deliveries, direct); err != nil {
				panic(err)
			}

		case 2:
			// Drop a peer:
			//   Byte 1: Peer index to drop
			peerIdx, err := r.ReadByte()
			if err != nil {
				return 0
			}
			peer := peers[int(peerIdx)%len(peers)]
			if verbose {
				fmt.Println("Drop", peer)
			}
			if err := f.Drop(peer); err != nil {
				panic(err)
			}

		case 3:
			// Move the simulated clock forward
			//   Byte 1: 100ms increment to move forward
			tickCnt, err := r.ReadByte()
			if err != nil {
				return 0
			}
			tick := time.Duration(tickCnt) * 100 * time.Millisecond
			if verbose {
				fmt.Println("Sleep", tick)
			}
			clock.Run(tick)
		}
	}
}
