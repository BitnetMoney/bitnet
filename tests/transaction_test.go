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

package tests

import (
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	txt := new(testMatcher)
	// These can't be parsed, invalid hex in RLP
	txt.skipLoad("^ttWrongRLP/.*")
	// We don't allow more than uint64 in gas amount
	// This is a pseudo-consensus vulnerability, but not in practice
	// because of the gas limit
	txt.skipLoad("^ttGasLimit/TransactionWithGasLimitxPriceOverflow.json")
	// We _do_ allow more than uint64 in gas price, as opposed to the tests
	// This is also not a concern, as long as tx.Cost() uses big.Int for
	// calculating the final cozt
	txt.skipLoad(".*TransactionWithGasPriceOverflow.*")

	// The nonce is too large for uint64. Not a concern, it means geth won't
	// accept transactions at a certain point in the distant future
	txt.skipLoad("^ttNonce/TransactionWithHighNonce256.json")

	// The value is larger than uint64, which according to the test is invalid.
	// Geth accepts it, which is not a consensus issue since we use big.Int's
	// internally to calculate the cost
	txt.skipLoad("^ttValue/TransactionWithHighValueOverflow.json")
	txt.walk(t, transactionTestDir, func(t *testing.T, name string, test *TransactionTest) {
		cfg := params.MainnetChainConfig
		if err := txt.checkFailure(t, test.Run(cfg)); err != nil {
			t.Error(err)
		}
	})
}
