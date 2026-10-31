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

package txpool

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestReplacedTxEvents(t *testing.T) {
	t.Parallel()

	pool, key := setupPool()
	defer pool.Stop()
	testAddBalance(pool, crypto.PubkeyToAddress(key.PublicKey), big.NewInt(1_000_000_000))

	events := make(chan ReplacedTxEvent, 8)
	sub := pool.SubscribeReplacedTxsEvent(events)
	defer sub.Unsubscribe()

	expect := func(what string, old, new *types.Transaction) {
		t.Helper()
		select {
		case ev := <-events:
			if ev.Old.Hash() != old.Hash() || ev.New.Hash() != new.Hash() {
				t.Errorf("%s: got %x -> %x, want %x -> %x", what, ev.Old.Hash(), ev.New.Hash(), old.Hash(), new.Hash())
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: no replacement event", what)
		}
	}
	add := func(tx *types.Transaction) error { return pool.addRemoteSync(tx) }

	// A pending transaction replaced with a higher fee.
	pending := pricedTransaction(0, 100000, big.NewInt(10), key)
	if err := add(pending); err != nil {
		t.Fatal(err)
	}
	bumped := pricedTransaction(0, 100000, big.NewInt(20), key)
	if err := add(bumped); err != nil {
		t.Fatal(err)
	}
	expect("pending", pending, bumped)

	// A queued (future nonce) transaction replaced.
	queued := pricedTransaction(5, 100000, big.NewInt(10), key)
	if err := add(queued); err != nil {
		t.Fatal(err)
	}
	queuedBumped := pricedTransaction(5, 100000, big.NewInt(20), key)
	if err := add(queuedBumped); err != nil {
		t.Fatal(err)
	}
	expect("queued", queued, queuedBumped)

	// A replacement without enough fee bump is rejected and announces nothing.
	if err := add(pricedTransaction(0, 100000, big.NewInt(20), key)); err == nil {
		t.Fatal("underpriced replacement accepted")
	}
	select {
	case ev := <-events:
		t.Errorf("event for a rejected replacement: %x -> %x", ev.Old.Hash(), ev.New.Hash())
	case <-time.After(200 * time.Millisecond):
	}
}
