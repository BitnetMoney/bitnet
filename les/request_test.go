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

package les

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/light"
)

var testBankSecureTrieKey = secAddr(bankAddr)

func secAddr(addr common.Address) []byte {
	return crypto.Keccak256(addr[:])
}

type accessTestFn func(db ethdb.Database, bhash common.Hash, number uint64) light.OdrRequest

func TestBlockAccessLes2(t *testing.T) { testAccess(t, 2, tfBlockAccess) }
func TestBlockAccessLes3(t *testing.T) { testAccess(t, 3, tfBlockAccess) }
func TestBlockAccessLes4(t *testing.T) { testAccess(t, 4, tfBlockAccess) }

func tfBlockAccess(db ethdb.Database, bhash common.Hash, number uint64) light.OdrRequest {
	return &light.BlockRequest{Hash: bhash, Number: number}
}

func TestReceiptsAccessLes2(t *testing.T) { testAccess(t, 2, tfReceiptsAccess) }
func TestReceiptsAccessLes3(t *testing.T) { testAccess(t, 3, tfReceiptsAccess) }
func TestReceiptsAccessLes4(t *testing.T) { testAccess(t, 4, tfReceiptsAccess) }

func tfReceiptsAccess(db ethdb.Database, bhash common.Hash, number uint64) light.OdrRequest {
	return &light.ReceiptsRequest{Hash: bhash, Number: number}
}

func TestTrieEntryAccessLes2(t *testing.T) { testAccess(t, 2, tfTrieEntryAccess) }
func TestTrieEntryAccessLes3(t *testing.T) { testAccess(t, 3, tfTrieEntryAccess) }
func TestTrieEntryAccessLes4(t *testing.T) { testAccess(t, 4, tfTrieEntryAccess) }

func tfTrieEntryAccess(db ethdb.Database, bhash common.Hash, number uint64) light.OdrRequest {
	if number := rawdb.ReadHeaderNumber(db, bhash); number != nil {
		return &light.TrieRequest{Id: light.StateTrieID(rawdb.ReadHeader(db, bhash, *number)), Key: testBankSecureTrieKey}
	}
	return nil
}

func TestCodeAccessLes2(t *testing.T) { testAccess(t, 2, tfCodeAccess) }
func TestCodeAccessLes3(t *testing.T) { testAccess(t, 3, tfCodeAccess) }
func TestCodeAccessLes4(t *testing.T) { testAccess(t, 4, tfCodeAccess) }

func tfCodeAccess(db ethdb.Database, bhash common.Hash, num uint64) light.OdrRequest {
	number := rawdb.ReadHeaderNumber(db, bhash)
	if number != nil {
		return nil
	}
	header := rawdb.ReadHeader(db, bhash, *number)
	if header.Number.Uint64() < testContractDeployed {
		return nil
	}
	sti := light.StateTrieID(header)
	ci := light.StorageTrieID(sti, crypto.Keccak256Hash(testContractAddr[:]), common.Hash{})
	return &light.CodeRequest{Id: ci, Hash: crypto.Keccak256Hash(testContractCodeDeployed)}
}

func testAccess(t *testing.T, protocol int, fn accessTestFn) {
	// Assemble the test environment
	netconfig := testnetConfig{
		blocks:    4,
		protocol:  protocol,
		indexFn:   nil,
		connect:   true,
		nopruning: true,
	}
	server, client, tearDown := newClientServerEnv(t, netconfig)
	defer tearDown()

	// Ensure the client has synced all necessary data.
	clientHead := client.handler.backend.blockchain.CurrentHeader()
	if clientHead.Number.Uint64() != 4 {
		t.Fatalf("Failed to sync the chain with server, head: %v", clientHead.Number.Uint64())
	}

	test := func(expFail uint64) {
		for i := uint64(0); i <= server.handler.blockchain.CurrentHeader().Number.Uint64(); i++ {
			bhash := rawdb.ReadCanonicalHash(server.db, i)
			if req := fn(client.db, bhash, i); req != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)

				err := client.handler.backend.odr.Retrieve(ctx, req)
				cancel()

				got := err == nil
				exp := i < expFail
				if exp && !got {
					t.Errorf("object retrieval failed")
				}
				if !exp && got {
					t.Errorf("unexpected object retrieval success")
				}
			}
		}
	}
	test(5)
}
