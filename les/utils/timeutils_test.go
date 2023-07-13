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
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

func TestUpdateTimer(t *testing.T) {
	timer := NewUpdateTimer(mclock.System{}, -1)
	if timer != nil {
		t.Fatalf("Create update timer with negative threshold")
	}
	sim := &mclock.Simulated{}
	timer = NewUpdateTimer(sim, time.Second)
	if updated := timer.Update(func(diff time.Duration) bool { return true }); updated {
		t.Fatalf("Update the clock without reaching the threshold")
	}
	sim.Run(time.Second)
	if updated := timer.Update(func(diff time.Duration) bool { return true }); !updated {
		t.Fatalf("Doesn't update the clock when reaching the threshold")
	}
	if updated := timer.UpdateAt(sim.Now().Add(time.Second), func(diff time.Duration) bool { return true }); !updated {
		t.Fatalf("Doesn't update the clock when reaching the threshold")
	}
	timer = NewUpdateTimer(sim, 0)
	if updated := timer.Update(func(diff time.Duration) bool { return true }); !updated {
		t.Fatalf("Doesn't update the clock without threshold limitaion")
	}
}
