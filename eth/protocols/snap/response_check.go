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
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

// Snap responses are decoded before the syncer matches them to a request, and a
// list of one-byte items turns into a list of 24-byte slices (CVE-2026-26313).
// The handler therefore counts the items of each such list before decoding it,
// and refuses responses with more of them than this node ever asks for.

const (
	// maxStorageSetsPerResponse is the most accounts a storage request covers.
	maxStorageSetsPerResponse = maxRequestSize / 1024

	// maxProofNodes bounds a range proof, which is two paths through the trie.
	maxProofNodes = 1024
)

var errTooManyItems = errors.New("message has too many items")

// The responses with their lists left encoded.
type (
	accountRangeEnvelope struct {
		ID       uint64
		Accounts rlp.RawValue
		Proof    rlp.RawValue
	}
	storageRangesEnvelope struct {
		ID    uint64
		Slots rlp.RawValue
		Proof rlp.RawValue
	}
	byteCodesEnvelope struct {
		ID    uint64
		Codes rlp.RawValue
	}
	trieNodesEnvelope struct {
		ID    uint64
		Nodes rlp.RawValue
	}
)

// decodeLimitedList decodes an encoded list into val if it has at most max items.
func decodeLimitedList(list rlp.RawValue, max int, val interface{}) error {
	content, _, err := rlp.SplitList(list)
	if err != nil {
		return err
	}
	items, err := rlp.CountValues(content)
	if err != nil {
		return err
	}
	if items > max {
		return fmt.Errorf("%w: %d, at most %d", errTooManyItems, items, max)
	}
	return rlp.DecodeBytes(list, val)
}
