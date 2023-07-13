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
	"math/rand"
	"testing"
)

type testWrsItem struct {
	idx  int
	widx *int
}

func testWeight(i interface{}) uint64 {
	t := i.(*testWrsItem)
	w := *t.widx
	if w == -1 || w == t.idx {
		return uint64(t.idx + 1)
	}
	return 0
}

func TestWeightedRandomSelect(t *testing.T) {
	testFn := func(cnt int) {
		s := NewWeightedRandomSelect(testWeight)
		w := -1
		list := make([]testWrsItem, cnt)
		for i := range list {
			list[i] = testWrsItem{idx: i, widx: &w}
			s.Update(&list[i])
		}
		w = rand.Intn(cnt)
		c := s.Choose()
		if c == nil {
			t.Errorf("expected item, got nil")
		} else {
			if c.(*testWrsItem).idx != w {
				t.Errorf("expected another item")
			}
		}
		w = -2
		if s.Choose() != nil {
			t.Errorf("expected nil, got item")
		}
	}
	testFn(1)
	testFn(10)
	testFn(100)
	testFn(1000)
	testFn(10000)
	testFn(100000)
	testFn(1000000)
}
