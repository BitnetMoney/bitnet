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
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

var (
	minDataGasPrice            = big.NewInt(params.BlobTxMinDataGasprice)
	dataGaspriceUpdateFraction = big.NewInt(params.BlobTxDataGaspriceUpdateFraction)
)

// CalcBlobFee calculates the blobfee from the header's excess data gas field.
func CalcBlobFee(excessDataGas *big.Int) *big.Int {
	// If this block does not yet have EIP-4844 enabled, return the starting fee
	if excessDataGas == nil {
		return big.NewInt(params.BlobTxMinDataGasprice)
	}
	return fakeExponential(minDataGasPrice, excessDataGas, dataGaspriceUpdateFraction)
}

// fakeExponential approximates factor * e ** (numerator / denominator) using
// Taylor expansion.
func fakeExponential(factor, numerator, denominator *big.Int) *big.Int {
	var (
		output = new(big.Int)
		accum  = new(big.Int).Mul(factor, denominator)
	)
	for i := 1; accum.Sign() > 0; i++ {
		output.Add(output, accum)

		accum.Mul(accum, numerator)
		accum.Div(accum, denominator)
		accum.Div(accum, big.NewInt(int64(i)))
	}
	return output.Div(output, denominator)
}
