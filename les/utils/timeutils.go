// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package utils

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

type UpdateTimer struct {
	clock     mclock.Clock
	lock      sync.Mutex
	last      mclock.AbsTime
	threshold time.Duration
}

func NewUpdateTimer(clock mclock.Clock, threshold time.Duration) *UpdateTimer {
	// We don't accept the update threshold less than 0.
	if threshold < 0 {
		return nil
	}
	// Don't panic for lazy users
	if clock == nil {
		clock = mclock.System{}
	}
	return &UpdateTimer{
		clock:     clock,
		last:      clock.Now(),
		threshold: threshold,
	}
}

func (t *UpdateTimer) Update(callback func(diff time.Duration) bool) bool {
	return t.UpdateAt(t.clock.Now(), callback)
}

func (t *UpdateTimer) UpdateAt(at mclock.AbsTime, callback func(diff time.Duration) bool) bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	diff := time.Duration(at - t.last)
	if diff < 0 {
		diff = 0
	}
	if diff < t.threshold {
		return false
	}
	if callback(diff) {
		t.last = at
		return true
	}
	return false
}
