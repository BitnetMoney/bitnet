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

	"github.com/ethereum/go-ethereum/common/mclock"
)

func TestValueExpiration(t *testing.T) {
	var cases = []struct {
		input      ExpiredValue
		timeOffset Fixed64
		expect     uint64
	}{
		{ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(0), 128},
		{ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(1), 64},
		{ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(2), 32},
		{ExpiredValue{Base: 128, Exp: 2}, Uint64ToFixed64(2), 128},
		{ExpiredValue{Base: 128, Exp: 2}, Uint64ToFixed64(3), 64},
	}
	for _, c := range cases {
		if got := c.input.Value(c.timeOffset); got != c.expect {
			t.Fatalf("Value mismatch, want=%d, got=%d", c.expect, got)
		}
	}
}

func TestValueAddition(t *testing.T) {
	var cases = []struct {
		input      ExpiredValue
		addend     int64
		timeOffset Fixed64
		expect     uint64
		expectNet  int64
	}{
		// Addition
		{ExpiredValue{Base: 128, Exp: 0}, 128, Uint64ToFixed64(0), 256, 128},
		{ExpiredValue{Base: 128, Exp: 2}, 128, Uint64ToFixed64(0), 640, 128},

		// Addition with offset
		{ExpiredValue{Base: 128, Exp: 0}, 128, Uint64ToFixed64(1), 192, 128},
		{ExpiredValue{Base: 128, Exp: 2}, 128, Uint64ToFixed64(1), 384, 128},
		{ExpiredValue{Base: 128, Exp: 2}, 128, Uint64ToFixed64(3), 192, 128},

		// Subtraction
		{ExpiredValue{Base: 128, Exp: 0}, -64, Uint64ToFixed64(0), 64, -64},
		{ExpiredValue{Base: 128, Exp: 0}, -128, Uint64ToFixed64(0), 0, -128},
		{ExpiredValue{Base: 128, Exp: 0}, -192, Uint64ToFixed64(0), 0, -128},

		// Subtraction with offset
		{ExpiredValue{Base: 128, Exp: 0}, -64, Uint64ToFixed64(1), 0, -64},
		{ExpiredValue{Base: 128, Exp: 0}, -128, Uint64ToFixed64(1), 0, -64},
		{ExpiredValue{Base: 128, Exp: 2}, -128, Uint64ToFixed64(1), 128, -128},
		{ExpiredValue{Base: 128, Exp: 2}, -128, Uint64ToFixed64(2), 0, -128},
	}
	for _, c := range cases {
		if net := c.input.Add(c.addend, c.timeOffset); net != c.expectNet {
			t.Fatalf("Net amount mismatch, want=%d, got=%d", c.expectNet, net)
		}
		if got := c.input.Value(c.timeOffset); got != c.expect {
			t.Fatalf("Value mismatch, want=%d, got=%d", c.expect, got)
		}
	}
}

func TestExpiredValueAddition(t *testing.T) {
	var cases = []struct {
		input      ExpiredValue
		another    ExpiredValue
		timeOffset Fixed64
		expect     uint64
	}{
		{ExpiredValue{Base: 128, Exp: 0}, ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(0), 256},
		{ExpiredValue{Base: 128, Exp: 1}, ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(0), 384},
		{ExpiredValue{Base: 128, Exp: 0}, ExpiredValue{Base: 128, Exp: 1}, Uint64ToFixed64(0), 384},
		{ExpiredValue{Base: 128, Exp: 0}, ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(1), 128},
	}
	for _, c := range cases {
		c.input.AddExp(c.another)
		if got := c.input.Value(c.timeOffset); got != c.expect {
			t.Fatalf("Value mismatch, want=%d, got=%d", c.expect, got)
		}
	}
}

func TestExpiredValueSubtraction(t *testing.T) {
	var cases = []struct {
		input      ExpiredValue
		another    ExpiredValue
		timeOffset Fixed64
		expect     uint64
	}{
		{ExpiredValue{Base: 128, Exp: 0}, ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(0), 0},
		{ExpiredValue{Base: 128, Exp: 0}, ExpiredValue{Base: 128, Exp: 1}, Uint64ToFixed64(0), 0},
		{ExpiredValue{Base: 128, Exp: 1}, ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(0), 128},
		{ExpiredValue{Base: 128, Exp: 1}, ExpiredValue{Base: 128, Exp: 0}, Uint64ToFixed64(1), 64},
	}
	for _, c := range cases {
		c.input.SubExp(c.another)
		if got := c.input.Value(c.timeOffset); got != c.expect {
			t.Fatalf("Value mismatch, want=%d, got=%d", c.expect, got)
		}
	}
}

func TestLinearExpiredValue(t *testing.T) {
	var cases = []struct {
		value  LinearExpiredValue
		now    mclock.AbsTime
		expect uint64
	}{
		{LinearExpiredValue{
			Offset: 0,
			Val:    0,
			Rate:   mclock.AbsTime(1),
		}, 0, 0},

		{LinearExpiredValue{
			Offset: 1,
			Val:    1,
			Rate:   mclock.AbsTime(1),
		}, 0, 1},

		{LinearExpiredValue{
			Offset: 1,
			Val:    1,
			Rate:   mclock.AbsTime(1),
		}, mclock.AbsTime(2), 0},

		{LinearExpiredValue{
			Offset: 1,
			Val:    1,
			Rate:   mclock.AbsTime(1),
		}, mclock.AbsTime(3), 0},
	}
	for _, c := range cases {
		if value := c.value.Value(c.now); value != c.expect {
			t.Fatalf("Value mismatch, want=%d, got=%d", c.expect, value)
		}
	}
}

func TestLinearExpiredAddition(t *testing.T) {
	var cases = []struct {
		value  LinearExpiredValue
		amount int64
		now    mclock.AbsTime
		expect uint64
	}{
		{LinearExpiredValue{
			Offset: 0,
			Val:    0,
			Rate:   mclock.AbsTime(1),
		}, -1, 0, 0},

		{LinearExpiredValue{
			Offset: 1,
			Val:    1,
			Rate:   mclock.AbsTime(1),
		}, -1, 0, 0},

		{LinearExpiredValue{
			Offset: 1,
			Val:    2,
			Rate:   mclock.AbsTime(1),
		}, -1, mclock.AbsTime(2), 0},

		{LinearExpiredValue{
			Offset: 1,
			Val:    2,
			Rate:   mclock.AbsTime(1),
		}, -2, mclock.AbsTime(2), 0},
	}
	for _, c := range cases {
		if value := c.value.Add(c.amount, c.now); value != c.expect {
			t.Fatalf("Value mismatch, want=%d, got=%d", c.expect, value)
		}
	}
}
