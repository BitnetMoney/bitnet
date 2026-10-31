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
	"crypto/rand"
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
)

// newCheckPeer returns a peer whose outgoing messages are read and discarded,
// so that requests can be sent and handlers called directly.
func newCheckPeer(t *testing.T) (*testBackend, *Peer) {
	t.Helper()
	backend := newTestBackend(4)
	app, net := p2p.MsgPipe()
	var id enode.ID
	rand.Read(id[:])
	peer := NewPeer(ETH66, p2p.NewPeer(id, "test", nil), net, backend.TxPool())
	go func() {
		for {
			msg, err := app.ReadMsg()
			if err != nil {
				return
			}
			msg.Discard()
		}
	}()
	t.Cleanup(func() {
		peer.Close()
		app.Close()
		backend.close()
	})
	return backend, peer
}

func rlpMsg(t *testing.T, code uint64, val interface{}) p2p.Msg {
	t.Helper()
	size, r, err := rlp.EncodeToReader(val)
	if err != nil {
		t.Fatal(err)
	}
	return p2p.Msg{Code: code, Size: uint32(size), Payload: r}
}

func emptyBodies(n int) BlockBodiesPacket {
	bodies := make(BlockBodiesPacket, n)
	for i := range bodies {
		bodies[i] = new(BlockBody)
	}
	return bodies
}

// A response nobody asked for is dropped without decoding its items: a million
// empty bodies (3 MB on the wire) must not turn into a hundred megabytes of Go
// values.
func TestUnsolicitedResponseNotDecoded(t *testing.T) {
	backend, peer := newCheckPeer(t)
	msg := rlpMsg(t, BlockBodiesMsg, &BlockBodiesPacket66{RequestId: 42, BlockBodiesPacket: emptyBodies(1_000_000)})

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	err := handleBlockBodies66(backend, msg, peer)
	runtime.ReadMemStats(&after)

	if err != nil {
		t.Fatalf("unsolicited response: got error %v, want it dropped quietly", err)
	}
	if alloc := after.TotalAlloc - before.TotalAlloc; alloc > 16<<20 {
		t.Errorf("handling the unsolicited response allocated %d MB", alloc>>20)
	}
}

func TestResponseItemsAgainstRequest(t *testing.T) {
	backend, peer := newCheckPeer(t)
	sink := make(chan *Response, 1)
	hashes := []common.Hash{{0x01}, {0x02}}

	// More bodies than requested: a protocol violation that drops the peer.
	req, err := peer.RequestBodies(hashes, sink)
	if err != nil {
		t.Fatal(err)
	}
	msg := rlpMsg(t, BlockBodiesMsg, &BlockBodiesPacket66{RequestId: req.id, BlockBodiesPacket: emptyBodies(3)})
	if err := handleBlockBodies66(backend, msg, peer); !errors.Is(err, errOversizedResponse) {
		t.Fatalf("3 bodies for 2 hashes: got %v, want errOversizedResponse", err)
	}
	req.Close()

	// As many as requested: delivered as before.
	req, err = peer.RequestBodies(hashes, sink)
	if err != nil {
		t.Fatal(err)
	}
	errc := make(chan error, 1)
	msg = rlpMsg(t, BlockBodiesMsg, &BlockBodiesPacket66{RequestId: req.id, BlockBodiesPacket: emptyBodies(2)})
	go func() { errc <- handleBlockBodies66(backend, msg, peer) }()
	res := <-sink
	if got := len(*res.Res.(*BlockBodiesPacket)); got != 2 {
		t.Errorf("delivered %d bodies, want 2", got)
	}
	res.Done <- nil
	if err := <-errc; err != nil {
		t.Errorf("valid response: %v", err)
	}
}

// acceptingBackend accepts transactions, so that the transaction handlers run
// up to the point where the packet is checked.
type acceptingBackend struct{ *testBackend }

func (acceptingBackend) AcceptTxs() bool { return true }

// A transaction broadcast with absurdly many items is refused before decoding.
func TestTooManyBroadcastTransactions(t *testing.T) {
	testBackend, peer := newCheckPeer(t)
	backend := acceptingBackend{testBackend}
	txs := make(TransactionsPacket, maxTransactionsPerMessage+1)
	for i := range txs {
		txs[i] = types.NewTx(&types.LegacyTx{})
	}
	// The handler reports it as a decoding error, which disconnects the peer.
	err := handleTransactions(backend, rlpMsg(t, TransactionsMsg, txs), peer)
	if !errors.Is(err, errDecode) || !strings.Contains(err.Error(), errTooManyItems.Error()) {
		t.Fatalf("got %v, want a decoding error for too many items", err)
	}
}
