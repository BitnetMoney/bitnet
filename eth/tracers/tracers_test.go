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

package tracers

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/tests"
)

func BenchmarkTransactionTrace(b *testing.B) {
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	from := crypto.PubkeyToAddress(key.PublicKey)
	gas := uint64(1000000) // 1M gas
	to := common.HexToAddress("0x00000000000000000000000000000000deadbeef")
	signer := types.LatestSignerForChainID(big.NewInt(1337))
	tx, err := types.SignNewTx(key, signer,
		&types.LegacyTx{
			Nonce:    1,
			GasPrice: big.NewInt(500),
			Gas:      gas,
			To:       &to,
		})
	if err != nil {
		b.Fatal(err)
	}
	txContext := vm.TxContext{
		Origin:   from,
		GasPrice: tx.GasPrice(),
	}
	context := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Coinbase:    common.Address{},
		BlockNumber: new(big.Int).SetUint64(uint64(5)),
		Time:        5,
		Difficulty:  big.NewInt(0xffffffff),
		GasLimit:    gas,
		BaseFee:     big.NewInt(8),
	}
	alloc := core.GenesisAlloc{}
	// The code pushes 'deadbeef' into memory, then the other params, and calls CREATE2, then returns
	// the address
	loop := []byte{
		byte(vm.JUMPDEST), //  [ count ]
		byte(vm.PUSH1), 0, // jumpdestination
		byte(vm.JUMP),
	}
	alloc[common.HexToAddress("0x00000000000000000000000000000000deadbeef")] = core.GenesisAccount{
		Nonce:   1,
		Code:    loop,
		Balance: big.NewInt(1),
	}
	alloc[from] = core.GenesisAccount{
		Nonce:   1,
		Code:    []byte{},
		Balance: big.NewInt(500000000000000),
	}
	_, statedb := tests.MakePreState(rawdb.NewMemoryDatabase(), alloc, false)
	// Create the tracer, the EVM environment and run it
	tracer := logger.NewStructLogger(&logger.Config{
		Debug: false,
		//DisableStorage: true,
		//EnableMemory: false,
		//EnableReturnData: false,
	})
	evm := vm.NewEVM(context, txContext, statedb, params.AllEthashProtocolChanges, vm.Config{Tracer: tracer})
	msg, err := core.TransactionToMessage(tx, signer, nil)
	if err != nil {
		b.Fatalf("failed to prepare transaction for tracing: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		snap := statedb.Snapshot()
		st := core.NewStateTransition(evm, msg, new(core.GasPool).AddGas(tx.Gas()))
		_, err = st.TransitionDb()
		if err != nil {
			b.Fatal(err)
		}
		statedb.RevertToSnapshot(snap)
		if have, want := len(tracer.StructLogs()), 244752; have != want {
			b.Fatalf("trace wrong, want %d steps, have %d", want, have)
		}
		tracer.Reset()
	}
}

func TestMemCopying(t *testing.T) {
	for i, tc := range []struct {
		memsize  int64
		offset   int64
		size     int64
		wantErr  string
		wantSize int
	}{
		{0, 0, 100, "", 100},    // Should pad up to 100
		{0, 100, 0, "", 0},      // No need to pad (0 size)
		{100, 50, 100, "", 100}, // Should pad 100-150
		{100, 50, 5, "", 5},     // Wanted range fully within memory
		{100, -50, 0, "offset or size must not be negative", 0},                        // Errror
		{0, 1, 1024*1024 + 1, "reached limit for padding memory slice: 1048578", 0},    // Errror
		{10, 0, 1024*1024 + 100, "reached limit for padding memory slice: 1048666", 0}, // Errror

	} {
		mem := vm.NewMemory()
		mem.Resize(uint64(tc.memsize))
		cpy, err := GetMemoryCopyPadded(mem, tc.offset, tc.size)
		if want := tc.wantErr; want != "" {
			if err == nil {
				t.Fatalf("test %d: want '%v' have no error", i, want)
			}
			if have := err.Error(); want != have {
				t.Fatalf("test %d: want '%v' have '%v'", i, want, have)
			}
			continue
		}
		if err != nil {
			t.Fatalf("test %d: unexpected error: %v", i, err)
		}
		if want, have := tc.wantSize, len(cpy); have != want {
			t.Fatalf("test %d: want %v have %v", i, want, have)
		}
	}
}
