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

package miner

import (
	"testing"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/beacon"
)

// powEngine is a proof-of-work engine that reports a fixed hashrate.
type powEngine struct {
	consensus.Engine
	rate float64
}

func (e powEngine) Hashrate() float64 { return e.rate }

// The node runs ethash inside the beacon wrapper; the hashrate must still be seen.
func TestHashrateThroughBeaconWrapper(t *testing.T) {
	inner := powEngine{rate: 8200}
	for _, tt := range []struct {
		name   string
		engine consensus.Engine
	}{
		{"bare engine", inner},
		{"beacon wrapper", beacon.New(inner)},
	} {
		miner := &Miner{engine: tt.engine}
		if got := miner.Hashrate(); got != 8200 {
			t.Errorf("%s: hashrate %d, want 8200", tt.name, got)
		}
	}
}
