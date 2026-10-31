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

package params

import (
	"testing"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

// A malformed entry would be skipped at startup with only a log line.
func TestMainnetBootnodesParse(t *testing.T) {
	if len(MainnetBootnodes) == 0 {
		t.Fatal("no mainnet bootnodes")
	}
	seen := make(map[enode.ID]bool)
	for _, url := range MainnetBootnodes {
		node, err := enode.Parse(enode.ValidSchemes, url)
		if err != nil {
			t.Errorf("%s: %v", url, err)
			continue
		}
		if node.IP() == nil || node.UDP() == 0 || node.TCP() == 0 {
			t.Errorf("%s: missing address", url)
		}
		if seen[node.ID()] {
			t.Errorf("%s: listed twice", url)
		}
		seen[node.ID()] = true
	}
}
