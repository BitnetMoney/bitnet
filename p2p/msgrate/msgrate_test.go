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

package msgrate

import "testing"

func TestCapacityOverflow(t *testing.T) {
	tracker := NewTracker(nil, 1)
	tracker.Update(1, 1, 100000)
	cap := tracker.Capacity(1, 10000000)
	if int32(cap) < 0 {
		t.Fatalf("Negative: %v", int32(cap))
	}
}
