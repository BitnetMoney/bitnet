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

package netutil

import (
	crand "crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

const (
	opStatement = iota
	opContact
	opPredict
	opCheckFullCone
)

type iptrackTestEvent struct {
	op       int
	time     int // absolute, in milliseconds
	ip, from string
}

func TestIPTracker(t *testing.T) {
	tests := map[string][]iptrackTestEvent{
		"minStatements": {
			{opPredict, 0, "", ""},
			{opStatement, 0, "127.0.0.1", "127.0.0.2"},
			{opPredict, 1000, "", ""},
			{opStatement, 1000, "127.0.0.1", "127.0.0.3"},
			{opPredict, 1000, "", ""},
			{opStatement, 1000, "127.0.0.1", "127.0.0.4"},
			{opPredict, 1000, "127.0.0.1", ""},
		},
		"window": {
			{opStatement, 0, "127.0.0.1", "127.0.0.2"},
			{opStatement, 2000, "127.0.0.1", "127.0.0.3"},
			{opStatement, 3000, "127.0.0.1", "127.0.0.4"},
			{opPredict, 10000, "127.0.0.1", ""},
			{opPredict, 10001, "", ""}, // first statement expired
			{opStatement, 10100, "127.0.0.1", "127.0.0.2"},
			{opPredict, 10200, "127.0.0.1", ""},
		},
		"fullcone": {
			{opContact, 0, "", "127.0.0.2"},
			{opStatement, 10, "127.0.0.1", "127.0.0.2"},
			{opContact, 2000, "", "127.0.0.3"},
			{opStatement, 2010, "127.0.0.1", "127.0.0.3"},
			{opContact, 3000, "", "127.0.0.4"},
			{opStatement, 3010, "127.0.0.1", "127.0.0.4"},
			{opCheckFullCone, 3500, "false", ""},
		},
		"fullcone_2": {
			{opContact, 0, "", "127.0.0.2"},
			{opStatement, 10, "127.0.0.1", "127.0.0.2"},
			{opContact, 2000, "", "127.0.0.3"},
			{opStatement, 2010, "127.0.0.1", "127.0.0.3"},
			{opStatement, 3000, "127.0.0.1", "127.0.0.4"},
			{opContact, 3010, "", "127.0.0.4"},
			{opCheckFullCone, 3500, "true", ""},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) { runIPTrackerTest(t, test) })
	}
}

func runIPTrackerTest(t *testing.T, evs []iptrackTestEvent) {
	var (
		clock mclock.Simulated
		it    = NewIPTracker(10*time.Second, 10*time.Second, 3)
	)
	it.clock = &clock
	for i, ev := range evs {
		evtime := time.Duration(ev.time) * time.Millisecond
		clock.Run(evtime - time.Duration(clock.Now()))
		switch ev.op {
		case opStatement:
			it.AddStatement(ev.from, ev.ip)
		case opContact:
			it.AddContact(ev.from)
		case opPredict:
			if pred := it.PredictEndpoint(); pred != ev.ip {
				t.Errorf("op %d: wrong prediction %q, want %q", i, pred, ev.ip)
			}
		case opCheckFullCone:
			pred := fmt.Sprintf("%t", it.PredictFullConeNAT())
			if pred != ev.ip {
				t.Errorf("op %d: wrong prediction %s, want %s", i, pred, ev.ip)
			}
		}
	}
}

// This checks that old statements and contacts are GCed even if Predict* isn't called.
func TestIPTrackerForceGC(t *testing.T) {
	var (
		clock  mclock.Simulated
		window = 10 * time.Second
		rate   = 50 * time.Millisecond
		max    = int(window/rate) + 1
		it     = NewIPTracker(window, window, 3)
	)
	it.clock = &clock

	for i := 0; i < 5*max; i++ {
		e1 := make([]byte, 4)
		e2 := make([]byte, 4)
		crand.Read(e1)
		crand.Read(e2)
		it.AddStatement(string(e1), string(e2))
		it.AddContact(string(e1))
		clock.Run(rate)
	}
	if len(it.contact) > 2*max {
		t.Errorf("contacts not GCed, have %d", len(it.contact))
	}
	if len(it.statements) > 2*max {
		t.Errorf("statements not GCed, have %d", len(it.statements))
	}
}
