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
	"testing"

	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
)

// signedNode returns a node whose record carries the given entries.
func signedNode(t *testing.T, entries ...enr.Entry) *enode.Node {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	var r enr.Record
	for _, e := range entries {
		r.Set(e)
	}
	if err := enode.SignV4(&r, key); err != nil {
		t.Fatal(err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func TestNodeFilter(t *testing.T) {
	backend := newTestBackend(8)
	defer backend.close()

	chain := backend.chain
	filter := NewNodeFilter(chain)
	head := chain.CurrentHeader()
	own := forkid.NewID(chain.Config(), chain.Genesis().Hash(), head.Number.Uint64(), head.Time)

	// Ethereum mainnet's fork ID at genesis: another chain sharing the discovery table.
	ethereum := forkid.ID{Hash: [4]byte{0xfc, 0x64, 0xec, 0x04}, Next: 1150000}

	tests := []struct {
		name string
		node *enode.Node
		want bool
	}{
		{"same chain", signedNode(t, &enrEntry{ForkID: own}), true},
		{"other chain", signedNode(t, &enrEntry{ForkID: ethereum}), false},
		{"no eth entry", signedNode(t), false},
	}
	for _, tt := range tests {
		if got := filter(tt.node); got != tt.want {
			t.Errorf("%s: filter = %v, want %v", tt.name, got, tt.want)
		}
	}
}
