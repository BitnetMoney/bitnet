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

// Package mclock is a wrapper for a monotonic clock source
package mclock

import (
	"time"

	_ "unsafe" // for go:linkname
)

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// AbsTime represents absolute monotonic time.
type AbsTime int64

// Now returns the current absolute monotonic time.
func Now() AbsTime {
	return AbsTime(nanotime())
}

// Add returns t + d as absolute time.
func (t AbsTime) Add(d time.Duration) AbsTime {
	return t + AbsTime(d)
}

// Sub returns t - t2 as a duration.
func (t AbsTime) Sub(t2 AbsTime) time.Duration {
	return time.Duration(t - t2)
}

// The Clock interface makes it possible to replace the monotonic system clock with
// a simulated clock.
type Clock interface {
	Now() AbsTime
	Sleep(time.Duration)
	NewTimer(time.Duration) ChanTimer
	After(time.Duration) <-chan AbsTime
	AfterFunc(d time.Duration, f func()) Timer
}

// Timer is a cancellable event created by AfterFunc.
type Timer interface {
	// Stop cancels the timer. It returns false if the timer has already
	// expired or been stopped.
	Stop() bool
}

// ChanTimer is a cancellable event created by NewTimer.
type ChanTimer interface {
	Timer

	// The channel returned by C receives a value when the timer expires.
	C() <-chan AbsTime
	// Reset reschedules the timer with a new timeout.
	// It should be invoked only on stopped or expired timers with drained channels.
	Reset(time.Duration)
}

// System implements Clock using the system clock.
type System struct{}

// Now returns the current monotonic time.
func (c System) Now() AbsTime {
	return Now()
}

// Sleep blocks for the given duration.
func (c System) Sleep(d time.Duration) {
	time.Sleep(d)
}

// NewTimer creates a timer which can be rescheduled.
func (c System) NewTimer(d time.Duration) ChanTimer {
	ch := make(chan AbsTime, 1)
	t := time.AfterFunc(d, func() {
		// This send is non-blocking because that's how time.Timer
		// behaves. It doesn't matter in the happy case, but does
		// when Reset is misused.
		select {
		case ch <- c.Now():
		default:
		}
	})
	return &systemTimer{t, ch}
}

// After returns a channel which receives the current time after d has elapsed.
func (c System) After(d time.Duration) <-chan AbsTime {
	ch := make(chan AbsTime, 1)
	time.AfterFunc(d, func() { ch <- c.Now() })
	return ch
}

// AfterFunc runs f on a new goroutine after the duration has elapsed.
func (c System) AfterFunc(d time.Duration, f func()) Timer {
	return time.AfterFunc(d, f)
}

type systemTimer struct {
	*time.Timer
	ch <-chan AbsTime
}

func (st *systemTimer) Reset(d time.Duration) {
	st.Timer.Reset(d)
}

func (st *systemTimer) C() <-chan AbsTime {
	return st.ch
}
