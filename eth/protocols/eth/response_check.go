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
	"errors"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/rlp"
)

// Decoding a list makes Go values that are far larger than their encoding: an
// empty block body is three bytes on the wire and over seventy in memory. A peer
// could therefore send ten-megabyte messages of tiny items and make the node
// allocate many times that (CVE-2026-26313, fixed in go-ethereum 1.17.0 by a
// larger rework). The checks below count the items of a list without decoding
// them, and only decode what is expected.

// maxTransactionsPerMessage bounds the transactions in one broadcast or pooled
// transactions message. Broadcasts are cut at 100 KB and pooled transaction
// requests ask for at most 256, so honest messages stay far below it.
const maxTransactionsPerMessage = 4096

var (
	errOversizedResponse = errors.New("response has more items than requested")
	errTooManyItems      = errors.New("message has too many items")
)

// responseEnvelope is the outer shape of every eth/66 response: the request id
// and the list of items, not yet decoded.
type responseEnvelope struct {
	RequestId uint64
	List      rlp.RawValue
}

// requestedItems returns how many items a request asked for.
func requestedItems(req *Request) (int, bool) {
	switch data := req.data.(type) {
	case *GetBlockHeadersPacket66:
		if data.GetBlockHeadersPacket == nil || data.Amount > math.MaxInt32 {
			return 0, false
		}
		return int(data.Amount), true
	case *GetBlockBodiesPacket66:
		return len(data.GetBlockBodiesPacket), true
	case *GetReceiptsPacket66:
		return len(data.GetReceiptsPacket), true
	case *GetNodeDataPacket66:
		return len(data.GetNodeDataPacket), true
	}
	return 0, false
}

// countItems returns the number of items in an encoded list without decoding them.
func countItems(list rlp.RawValue) (int, error) {
	content, _, err := rlp.SplitList(list)
	if err != nil {
		return 0, err
	}
	return rlp.CountValues(content)
}

// readResponse reads the request id of a response and checks it before the items
// are decoded. It returns the encoded item list when the response answers one of
// this peer's pending requests of the right kind with no more items than were
// asked for. A response that matches no pending request, such as one arriving
// after its request timed out, gives ok false: it is dropped without decoding its
// items, but the peer is kept. More items than asked for is a protocol violation
// and returns an error, which disconnects the peer.
func readResponse(msg Decoder, peer *Peer, code uint64) (id uint64, list rlp.RawValue, ok bool, err error) {
	var env responseEnvelope
	if err := msg.Decode(&env); err != nil {
		return 0, nil, false, fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	req := peer.pendingRequest(env.RequestId)
	if req == nil || req.want != code {
		return env.RequestId, nil, false, nil
	}
	items, err := countItems(env.List)
	if err != nil {
		return 0, nil, false, fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if want, known := requestedItems(req); known && items > want {
		return 0, nil, false, fmt.Errorf("%w: %d items for a request of %d", errOversizedResponse, items, want)
	}
	return env.RequestId, env.List, true, nil
}

// decodeLimitedList decodes an encoded list into val if it has at most max items.
func decodeLimitedList(list rlp.RawValue, max int, val interface{}) error {
	items, err := countItems(list)
	if err != nil {
		return err
	}
	if items > max {
		return fmt.Errorf("%w: %d, at most %d", errTooManyItems, items, max)
	}
	return rlp.DecodeBytes(list, val)
}
