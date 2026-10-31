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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// ReplacedTxEvent is sent when a transaction in the pool is replaced by another
// one from the same sender with the same nonce. Only one of them can be mined.
type ReplacedTxEvent struct {
	Old *types.Transaction
	New *types.Transaction
}

// SubscribeReplacedTxsEvent registers a subscription of ReplacedTxEvent. Events
// are sent after the pool has processed the replacement.
func (pool *TxPool) SubscribeReplacedTxsEvent(ch chan<- ReplacedTxEvent) event.Subscription {
	return pool.scope.Track(pool.replaceFeed.Subscribe(ch))
}
