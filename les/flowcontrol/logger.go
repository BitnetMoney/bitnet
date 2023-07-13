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

package flowcontrol

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

// logger collects events in string format and discards events older than the
// "keep" parameter
type logger struct {
	events           map[uint64]logEvent
	writePtr, delPtr uint64
	keep             time.Duration
}

// logEvent describes a single event
type logEvent struct {
	time  mclock.AbsTime
	event string
}

// newLogger creates a new logger
func newLogger(keep time.Duration) *logger {
	return &logger{
		events: make(map[uint64]logEvent),
		keep:   keep,
	}
}

// add adds a new event and discards old events if possible
func (l *logger) add(now mclock.AbsTime, event string) {
	keepAfter := now - mclock.AbsTime(l.keep)
	for l.delPtr < l.writePtr && l.events[l.delPtr].time <= keepAfter {
		delete(l.events, l.delPtr)
		l.delPtr++
	}
	l.events[l.writePtr] = logEvent{now, event}
	l.writePtr++
}

// dump prints all stored events
func (l *logger) dump(now mclock.AbsTime) {
	for i := l.delPtr; i < l.writePtr; i++ {
		e := l.events[i]
		fmt.Println(time.Duration(e.time-now), e.event)
	}
}
