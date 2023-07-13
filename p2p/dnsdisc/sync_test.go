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

package dnsdisc

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestLinkCache(t *testing.T) {
	var lc linkCache

	// Check adding links.
	lc.addLink("1", "2")
	if !lc.changed {
		t.Error("changed flag not set")
	}
	lc.changed = false
	lc.addLink("1", "2")
	if lc.changed {
		t.Error("changed flag set after adding link that's already present")
	}
	lc.addLink("2", "3")
	lc.addLink("3", "1")
	lc.addLink("2", "4")
	lc.changed = false

	if !lc.isReferenced("3") {
		t.Error("3 not referenced")
	}
	if lc.isReferenced("6") {
		t.Error("6 is referenced")
	}

	lc.resetLinks("1", nil)
	if !lc.changed {
		t.Error("changed flag not set")
	}
	if len(lc.backrefs) != 0 {
		t.Logf("%+v", lc)
		t.Error("reference maps should be empty")
	}
}

func TestLinkCacheRandom(t *testing.T) {
	tags := make([]string, 1000)
	for i := range tags {
		tags[i] = strconv.Itoa(i)
	}

	// Create random links.
	var lc linkCache
	var remove []string
	for i := 0; i < 100; i++ {
		a, b := tags[rand.Intn(len(tags))], tags[rand.Intn(len(tags))]
		lc.addLink(a, b)
		remove = append(remove, a)
	}

	// Remove all the links.
	for _, s := range remove {
		lc.resetLinks(s, nil)
	}
	if len(lc.backrefs) != 0 {
		t.Logf("%+v", lc)
		t.Error("reference maps should be empty")
	}
}
