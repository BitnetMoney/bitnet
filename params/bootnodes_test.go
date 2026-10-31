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

import "testing"

// The Bitnet genesis must not pick up Ethereum's DNS node list.
func TestKnownDNSNetworkBitnet(t *testing.T) {
	for _, protocol := range []string{"all", "les", "snap"} {
		if url := KnownDNSNetwork(MainnetGenesisHash, protocol); url != "" {
			t.Errorf("protocol %q: got DNS list %q for the Bitnet genesis, want none", protocol, url)
		}
	}
}
