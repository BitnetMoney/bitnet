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

//go:build (amd64 && blsasm) || (amd64 && blsadx)
// +build amd64,blsasm amd64,blsadx

package bls12381

import (
	"golang.org/x/sys/cpu"
)

func init() {
	if !enableADX || !cpu.X86.HasADX || !cpu.X86.HasBMI2 {
		mul = mulNoADX
	}
}

// Use ADX backend for default
var mul func(c, a, b *fe) = mulADX

func square(c, a *fe) {
	mul(c, a, a)
}

func neg(c, a *fe) {
	if a.isZero() {
		c.set(a)
	} else {
		_neg(c, a)
	}
}

//go:noescape
func add(c, a, b *fe)

//go:noescape
func addAssign(a, b *fe)

//go:noescape
func ladd(c, a, b *fe)

//go:noescape
func laddAssign(a, b *fe)

//go:noescape
func double(c, a *fe)

//go:noescape
func doubleAssign(a *fe)

//go:noescape
func ldouble(c, a *fe)

//go:noescape
func sub(c, a, b *fe)

//go:noescape
func subAssign(a, b *fe)

//go:noescape
func lsubAssign(a, b *fe)

//go:noescape
func _neg(c, a *fe)

//go:noescape
func mulNoADX(c, a, b *fe)

//go:noescape
func mulADX(c, a, b *fe)
