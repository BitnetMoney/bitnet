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

	"github.com/ethereum/go-ethereum/common"
)

func bigFromHex(hex string) *big.Int {
	return new(big.Int).SetBytes(common.FromHex(hex))
}

// decodeFieldElement expects 64 byte input with zero top 16 bytes,
// returns lower 48 bytes.
func decodeFieldElement(in []byte) ([]byte, error) {
	if len(in) != 64 {
		return nil, errors.New("invalid field element length")
	}
	// check top bytes
	for i := 0; i < 16; i++ {
		if in[i] != byte(0x00) {
			return nil, errors.New("invalid field element top bytes")
		}
	}
	out := make([]byte, 48)
	copy(out[:], in[16:])
	return out, nil
}
