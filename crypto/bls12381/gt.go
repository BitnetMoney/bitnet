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

package bls12381

import (
	"errors"
	"math/big"
)

// E is type for target group element
type E = fe12

// GT is type for target multiplicative group GT.
type GT struct {
	fp12 *fp12
}

func (e *E) Set(e2 *E) *E {
	return e.set(e2)
}

// One sets a new target group element to one
func (e *E) One() *E {
	e = new(fe12).one()
	return e
}

// IsOne returns true if given element equals to one
func (e *E) IsOne() bool {
	return e.isOne()
}

// Equal returns true if given two element is equal, otherwise returns false
func (g *E) Equal(g2 *E) bool {
	return g.equal(g2)
}

// NewGT constructs new target group instance.
func NewGT() *GT {
	fp12 := newFp12(nil)
	return &GT{fp12}
}

// Q returns group order in big.Int.
func (g *GT) Q() *big.Int {
	return new(big.Int).Set(q)
}

// FromBytes expects 576 byte input and returns target group element
// FromBytes returns error if given element is not on correct subgroup.
func (g *GT) FromBytes(in []byte) (*E, error) {
	e, err := g.fp12.fromBytes(in)
	if err != nil {
		return nil, err
	}
	if !g.IsValid(e) {
		return e, errors.New("invalid element")
	}
	return e, nil
}

// ToBytes serializes target group element.
func (g *GT) ToBytes(e *E) []byte {
	return g.fp12.toBytes(e)
}

// IsValid checks whether given target group element is in correct subgroup.
func (g *GT) IsValid(e *E) bool {
	r := g.New()
	g.fp12.exp(r, e, q)
	return r.isOne()
}

// New initializes a new target group element which is equal to one
func (g *GT) New() *E {
	return new(E).One()
}

// Add adds two field element `a` and `b` and assigns the result to the element in first argument.
func (g *GT) Add(c, a, b *E) {
	g.fp12.add(c, a, b)
}

// Sub subtracts two field element `a` and `b`, and assigns the result to the element in first argument.
func (g *GT) Sub(c, a, b *E) {
	g.fp12.sub(c, a, b)
}

// Mul multiplies two field element `a` and `b` and assigns the result to the element in first argument.
func (g *GT) Mul(c, a, b *E) {
	g.fp12.mul(c, a, b)
}

// Square squares an element `a` and assigns the result to the element in first argument.
func (g *GT) Square(c, a *E) {
	g.fp12.cyclotomicSquare(c, a)
}

// Exp exponents an element `a` by a scalar `s` and assigns the result to the element in first argument.
func (g *GT) Exp(c, a *E, s *big.Int) {
	g.fp12.cyclotomicExp(c, a, s)
}

// Inverse inverses an element `a` and assigns the result to the element in first argument.
func (g *GT) Inverse(c, a *E) {
	g.fp12.inverse(c, a)
}
