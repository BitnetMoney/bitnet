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

package bind_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

var waitDeployedTests = map[string]struct {
	code        string
	gas         uint64
	wantAddress common.Address
	wantErr     error
}{
	"successful deploy": {
		code:        `6060604052600a8060106000396000f360606040526008565b00`,
		gas:         3000000,
		wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
	},
	"empty code": {
		code:        ``,
		gas:         300000,
		wantErr:     bind.ErrNoCodeAfterDeploy,
		wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
	},
}

func TestWaitDeployed(t *testing.T) {
	for name, test := range waitDeployedTests {
		backend := backends.NewSimulatedBackend(
			core.GenesisAlloc{
				crypto.PubkeyToAddress(testKey.PublicKey): {Balance: big.NewInt(10000000000000000)},
			},
			10000000,
		)
		defer backend.Close()

		// Create the transaction
		head, _ := backend.HeaderByNumber(context.Background(), nil) // Should be child's, good enough
		gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(1))

		tx := types.NewContractCreation(0, big.NewInt(0), test.gas, gasPrice, common.FromHex(test.code))
		tx, _ = types.SignTx(tx, types.HomesteadSigner{}, testKey)

		// Wait for it to get mined in the background.
		var (
			err     error
			address common.Address
			mined   = make(chan struct{})
			ctx     = context.Background()
		)
		go func() {
			address, err = bind.WaitDeployed(ctx, backend, tx)
			close(mined)
		}()

		// Send and mine the transaction.
		backend.SendTransaction(ctx, tx)
		backend.Commit()

		select {
		case <-mined:
			if err != test.wantErr {
				t.Errorf("test %q: error mismatch: want %q, got %q", name, test.wantErr, err)
			}
			if address != test.wantAddress {
				t.Errorf("test %q: unexpected contract address %s", name, address.Hex())
			}
		case <-time.After(2 * time.Second):
			t.Errorf("test %q: timeout", name)
		}
	}
}

func TestWaitDeployedCornerCases(t *testing.T) {
	backend := backends.NewSimulatedBackend(
		core.GenesisAlloc{
			crypto.PubkeyToAddress(testKey.PublicKey): {Balance: big.NewInt(10000000000000000)},
		},
		10000000,
	)
	defer backend.Close()

	head, _ := backend.HeaderByNumber(context.Background(), nil) // Should be child's, good enough
	gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(1))

	// Create a transaction to an account.
	code := "6060604052600a8060106000396000f360606040526008565b00"
	tx := types.NewTransaction(0, common.HexToAddress("0x01"), big.NewInt(0), 3000000, gasPrice, common.FromHex(code))
	tx, _ = types.SignTx(tx, types.HomesteadSigner{}, testKey)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	backend.SendTransaction(ctx, tx)
	backend.Commit()
	notContentCreation := errors.New("tx is not contract creation")
	if _, err := bind.WaitDeployed(ctx, backend, tx); err.Error() != notContentCreation.Error() {
		t.Errorf("error missmatch: want %q, got %q, ", notContentCreation, err)
	}

	// Create a transaction that is not mined.
	tx = types.NewContractCreation(1, big.NewInt(0), 3000000, gasPrice, common.FromHex(code))
	tx, _ = types.SignTx(tx, types.HomesteadSigner{}, testKey)

	go func() {
		contextCanceled := errors.New("context canceled")
		if _, err := bind.WaitDeployed(ctx, backend, tx); err.Error() != contextCanceled.Error() {
			t.Errorf("error missmatch: want %q, got %q, ", contextCanceled, err)
		}
	}()

	backend.SendTransaction(ctx, tx)
	cancel()
}
