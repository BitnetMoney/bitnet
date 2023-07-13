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

package p2p

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

func TestExpHeap(t *testing.T) {
	var h expHeap

	var (
		basetime = mclock.AbsTime(10)
		exptimeA = basetime.Add(2 * time.Second)
		exptimeB = basetime.Add(3 * time.Second)
		exptimeC = basetime.Add(4 * time.Second)
	)
	h.add("b", exptimeB)
	h.add("a", exptimeA)
	h.add("c", exptimeC)

	if h.nextExpiry() != exptimeA {
		t.Fatal("wrong nextExpiry")
	}
	if !h.contains("a") || !h.contains("b") || !h.contains("c") {
		t.Fatal("heap doesn't contain all live items")
	}

	h.expire(exptimeA.Add(1), nil)
	if h.nextExpiry() != exptimeB {
		t.Fatal("wrong nextExpiry")
	}
	if h.contains("a") {
		t.Fatal("heap contains a even though it has already expired")
	}
	if !h.contains("b") || !h.contains("c") {
		t.Fatal("heap doesn't contain all live items")
	}
}
