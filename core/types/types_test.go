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

package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type devnull struct{ len int }

func (d *devnull) Write(p []byte) (int, error) {
	d.len += len(p)
	return len(p), nil
}

func BenchmarkEncodeRLP(b *testing.B) {
	benchRLP(b, true)
}

func BenchmarkDecodeRLP(b *testing.B) {
	benchRLP(b, false)
}

func benchRLP(b *testing.B, encode bool) {
	key, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	to := common.HexToAddress("0x00000000000000000000000000000000deadbeef")
	signer := NewLondonSigner(big.NewInt(1337))
	for _, tc := range []struct {
		name string
		obj  interface{}
	}{
		{
			"legacy-header",
			&Header{
				Difficulty: big.NewInt(10000000000),
				Number:     big.NewInt(1000),
				GasLimit:   8_000_000,
				GasUsed:    8_000_000,
				Time:       555,
				Extra:      make([]byte, 32),
			},
		},
		{
			"london-header",
			&Header{
				Difficulty: big.NewInt(10000000000),
				Number:     big.NewInt(1000),
				GasLimit:   8_000_000,
				GasUsed:    8_000_000,
				Time:       555,
				Extra:      make([]byte, 32),
				BaseFee:    big.NewInt(10000000000),
			},
		},
		{
			"receipt-for-storage",
			&ReceiptForStorage{
				Status:            ReceiptStatusSuccessful,
				CumulativeGasUsed: 0x888888888,
				Logs:              make([]*Log, 0),
			},
		},
		{
			"receipt-full",
			&Receipt{
				Status:            ReceiptStatusSuccessful,
				CumulativeGasUsed: 0x888888888,
				Logs:              make([]*Log, 0),
			},
		},
		{
			"legacy-transaction",
			MustSignNewTx(key, signer,
				&LegacyTx{
					Nonce:    1,
					GasPrice: big.NewInt(500),
					Gas:      1000000,
					To:       &to,
					Value:    big.NewInt(1),
				}),
		},
		{
			"access-transaction",
			MustSignNewTx(key, signer,
				&AccessListTx{
					Nonce:    1,
					GasPrice: big.NewInt(500),
					Gas:      1000000,
					To:       &to,
					Value:    big.NewInt(1),
				}),
		},
		{
			"1559-transaction",
			MustSignNewTx(key, signer,
				&DynamicFeeTx{
					Nonce:     1,
					Gas:       1000000,
					To:        &to,
					Value:     big.NewInt(1),
					GasTipCap: big.NewInt(500),
					GasFeeCap: big.NewInt(500),
				}),
		},
	} {
		if encode {
			b.Run(tc.name, func(b *testing.B) {
				b.ReportAllocs()
				var null = &devnull{}
				for i := 0; i < b.N; i++ {
					rlp.Encode(null, tc.obj)
				}
				b.SetBytes(int64(null.len / b.N))
			})
		} else {
			data, _ := rlp.EncodeToBytes(tc.obj)
			// Test decoding
			b.Run(tc.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					if err := rlp.DecodeBytes(data, tc.obj); err != nil {
						b.Fatal(err)
					}
				}
				b.SetBytes(int64(len(data)))
			})
		}
	}
}
