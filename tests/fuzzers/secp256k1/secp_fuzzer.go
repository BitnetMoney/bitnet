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

// build +gofuzz

package secp256k1

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	fuzz "github.com/google/gofuzz"
)

func Fuzz(input []byte) int {
	var (
		fuzzer = fuzz.NewFromGoFuzz(input)
		curveA = secp256k1.S256()
		curveB = btcec.S256()
		dataP1 []byte
		dataP2 []byte
	)
	// first point
	fuzzer.Fuzz(&dataP1)
	x1, y1 := curveB.ScalarBaseMult(dataP1)
	// second point
	fuzzer.Fuzz(&dataP2)
	x2, y2 := curveB.ScalarBaseMult(dataP2)
	resAX, resAY := curveA.Add(x1, y1, x2, y2)
	resBX, resBY := curveB.Add(x1, y1, x2, y2)
	if resAX.Cmp(resBX) != 0 || resAY.Cmp(resBY) != 0 {
		fmt.Printf("%s %s %s %s\n", x1, y1, x2, y2)
		panic(fmt.Sprintf("Addition failed: geth: %s %s btcd: %s %s", resAX, resAY, resBX, resBY))
	}
	return 0
}
