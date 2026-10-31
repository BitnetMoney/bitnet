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

package snap

import (
	"crypto/rand"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// recordingBackend keeps the last response handed to it.
type recordingBackend struct{ got Packet }

func (b *recordingBackend) Chain() *core.BlockChain                { panic("not used") }
func (b *recordingBackend) RunPeer(*Peer, Handler) error           { panic("not used") }
func (b *recordingBackend) PeerInfo(enode.ID) interface{}          { panic("not used") }
func (b *recordingBackend) Handle(peer *Peer, packet Packet) error { b.got = packet; return nil }

// handle sends one message to a fresh peer and runs the handler on it.
func handle(t *testing.T, code uint64, packet interface{}) (*recordingBackend, error) {
	t.Helper()
	app, net := p2p.MsgPipe()
	defer app.Close()
	var id enode.ID
	rand.Read(id[:])
	peer := NewPeer(SNAP1, p2p.NewPeer(id, "test", nil), net)
	go p2p.Send(app, code, packet)
	backend := new(recordingBackend)
	return backend, HandleMessage(backend, peer)
}

func blobs(n int) [][]byte {
	b := make([][]byte, n)
	for i := range b {
		b[i] = []byte{}
	}
	return b
}

func TestSnapResponseLimits(t *testing.T) {
	tests := []struct {
		name   string
		code   uint64
		packet interface{}
	}{
		{"bytecodes", ByteCodesMsg, &ByteCodesPacket{ID: 1, Codes: blobs(maxCodeRequestCount + 1)}},
		{"trie nodes", TrieNodesMsg, &TrieNodesPacket{ID: 1, Nodes: blobs(maxTrieRequestCount + 1)}},
		{"storage sets", StorageRangesMsg, &StorageRangesPacket{ID: 1, Slots: make([][]*StorageData, maxStorageSetsPerResponse+1)}},
		{"storage proof", StorageRangesMsg, &StorageRangesPacket{ID: 1, Proof: blobs(maxProofNodes + 1)}},
		{"account proof", AccountRangeMsg, &AccountRangePacket{ID: 1, Proof: blobs(maxProofNodes + 1)}},
	}
	for _, tt := range tests {
		backend, err := handle(t, tt.code, tt.packet)
		if err == nil || !strings.Contains(err.Error(), errTooManyItems.Error()) {
			t.Errorf("%s: got %v, want too many items", tt.name, err)
		}
		if backend.got != nil {
			t.Errorf("%s: an oversized response reached the syncer", tt.name)
		}
	}
	// Within the limits the response reaches the syncer as before.
	backend, err := handle(t, ByteCodesMsg, &ByteCodesPacket{ID: 7, Codes: [][]byte{{0x60}, {0x61}}})
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := backend.got.(*ByteCodesPacket); !ok || got.ID != 7 || len(got.Codes) != 2 {
		t.Errorf("valid response delivered as %#v", backend.got)
	}
}

// Three million empty bytecodes (3 MB on the wire) are refused without being
// decoded into 72 MB of slices.
func TestSnapOversizedResponseNotDecoded(t *testing.T) {
	app, net := p2p.MsgPipe()
	defer app.Close()
	var id enode.ID
	rand.Read(id[:])
	peer := NewPeer(SNAP1, p2p.NewPeer(id, "test", nil), net)
	go p2p.Send(app, ByteCodesMsg, &ByteCodesPacket{ID: 1, Codes: blobs(3_000_000)})

	// Wait for the message before measuring, so its encoding is not counted.
	msg, err := peer.rw.ReadMsg()
	if err != nil {
		t.Fatal(err)
	}
	relay, relayNet := p2p.MsgPipe()
	defer relay.Close()
	go relay.WriteMsg(msg)
	peer.rw = relayNet

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	err = HandleMessage(new(recordingBackend), peer)
	runtime.ReadMemStats(&after)

	if err == nil {
		t.Fatal("oversized response accepted")
	}
	if alloc := after.TotalAlloc - before.TotalAlloc; alloc > 16<<20 {
		t.Errorf("refusing the response allocated %d MB", alloc>>20)
	}
}
