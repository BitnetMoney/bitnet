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

package misc

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func TestCalcBlobFee(t *testing.T) {
	tests := []struct {
		excessDataGas int64
		blobfee       int64
	}{
		{0, 1},
		{1542706, 1},
		{1542707, 2},
		{10 * 1024 * 1024, 111},
	}
	have := CalcBlobFee(nil)
	if have.Int64() != params.BlobTxMinDataGasprice {
		t.Errorf("nil test: blobfee mismatch: have %v, want %v", have, params.BlobTxMinDataGasprice)
	}
	for i, tt := range tests {
		have := CalcBlobFee(big.NewInt(tt.excessDataGas))
		if have.Int64() != tt.blobfee {
			t.Errorf("test %d: blobfee mismatch: have %v want %v", i, have, tt.blobfee)
		}
	}
}

func TestFakeExponential(t *testing.T) {
	tests := []struct {
		factor      int64
		numerator   int64
		denominator int64
		want        int64
	}{
		// When numerator == 0 the return value should always equal the value of factor
		{1, 0, 1, 1},
		{38493, 0, 1000, 38493},
		{0, 1234, 2345, 0}, // should be 0
		{1, 2, 1, 6},       // approximate 7.389
		{1, 4, 2, 6},
		{1, 3, 1, 16}, // approximate 20.09
		{1, 6, 2, 18},
		{1, 4, 1, 49}, // approximate 54.60
		{1, 8, 2, 50},
		{10, 8, 2, 542}, // approximate 540.598
		{11, 8, 2, 596}, // approximate 600.58
		{1, 5, 1, 136},  // approximate 148.4
		{1, 5, 2, 11},   // approximate 12.18
		{2, 5, 2, 23},   // approximate 24.36
		{1, 50000000, 2225652, 5709098764},
	}
	for i, tt := range tests {
		f, n, d := big.NewInt(tt.factor), big.NewInt(tt.numerator), big.NewInt(tt.denominator)
		original := fmt.Sprintf("%d %d %d", f, n, d)
		have := fakeExponential(f, n, d)
		if have.Int64() != tt.want {
			t.Errorf("test %d: fake exponential mismatch: have %v want %v", i, have, tt.want)
		}
		later := fmt.Sprintf("%d %d %d", f, n, d)
		if original != later {
			t.Errorf("test %d: fake exponential modified arguments: have\n%v\nwant\n%v", i, later, original)
		}
	}
}
