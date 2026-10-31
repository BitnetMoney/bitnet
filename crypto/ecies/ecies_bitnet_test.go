// Copyright 2026 The Bitnet community authors
// This file is part of the Bitnet library.
//
// The Bitnet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Bitnet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

package ecies

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"
)

// A ciphertext with a valid MAC but less than one cipher block after the
// ephemeral key must be rejected, not sliced past its end. Anyone can build
// one, since anyone can encrypt to a node's public key.
func TestDecryptShortCiphertext(t *testing.T) {
	prv, err := GenerateKey(rand.Reader, DefaultCurve, nil)
	if err != nil {
		t.Fatal(err)
	}
	pub := &prv.PublicKey
	params, err := pubkeyParams(pub)
	if err != nil {
		t.Fatal(err)
	}
	R, err := GenerateKey(rand.Reader, pub.Curve, params)
	if err != nil {
		t.Fatal(err)
	}
	z, err := R.GenerateShared(pub, params.KeyLen, params.KeyLen)
	if err != nil {
		t.Fatal(err)
	}
	_, Km := deriveKeys(params.Hash(), z, nil, params.KeyLen)

	em := []byte{0x01} // shorter than the IV
	d := messageTag(params.Hash, Km, em, nil)
	Rb := elliptic.Marshal(pub.Curve, R.PublicKey.X, R.PublicKey.Y)
	ct := append(append(append([]byte{}, Rb...), em...), d...)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Decrypt panicked: %v", r)
		}
	}()
	if _, err := prv.Decrypt(ct, nil, nil); err != ErrInvalidMessage {
		t.Fatalf("Decrypt error = %v, want %v", err, ErrInvalidMessage)
	}
}

// Key agreement with a point that is not on the curve must fail: multiplying
// such points by the private key leaks bits of it.
func TestGenerateSharedOffCurve(t *testing.T) {
	prv, err := GenerateKey(rand.Reader, DefaultCurve, nil)
	if err != nil {
		t.Fatal(err)
	}
	p := prv.PublicKey.Curve.Params().P
	tests := []struct {
		name string
		x, y *big.Int
	}{
		{"not on the curve", big.NewInt(1), big.NewInt(1)},
		{"coordinate above the field", new(big.Int).Add(prv.PublicKey.X, p), prv.PublicKey.Y},
		{"missing coordinate", nil, nil},
	}
	for _, tt := range tests {
		pub := &PublicKey{X: tt.x, Y: tt.y, Curve: prv.PublicKey.Curve, Params: prv.PublicKey.Params}
		if _, err := prv.GenerateShared(pub, 16, 16); err != ErrInvalidPublicKey {
			t.Errorf("%s: GenerateShared error = %v, want %v", tt.name, err, ErrInvalidPublicKey)
		}
	}
}
