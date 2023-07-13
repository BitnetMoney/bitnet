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

package rlp

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

func decodeEncode(input []byte, val interface{}, i int) {
	if err := rlp.DecodeBytes(input, val); err == nil {
		output, err := rlp.EncodeToBytes(val)
		if err != nil {
			panic(err)
		}
		if !bytes.Equal(input, output) {
			panic(fmt.Sprintf("case %d: encode-decode is not equal, \ninput : %x\noutput: %x", i, input, output))
		}
	}
}

func Fuzz(input []byte) int {
	if len(input) == 0 {
		return 0
	}
	if len(input) > 500*1024 {
		return 0
	}

	var i int
	{
		rlp.Split(input)
	}
	{
		if elems, _, err := rlp.SplitList(input); err == nil {
			rlp.CountValues(elems)
		}
	}

	{
		rlp.NewStream(bytes.NewReader(input), 0).Decode(new(interface{}))
	}

	{
		decodeEncode(input, new(interface{}), i)
		i++
	}
	{
		var v struct {
			Int    uint
			String string
			Bytes  []byte
		}
		decodeEncode(input, &v, i)
		i++
	}

	{
		type Types struct {
			Bool  bool
			Raw   rlp.RawValue
			Slice []*Types
			Iface []interface{}
		}
		var v Types
		decodeEncode(input, &v, i)
		i++
	}
	{
		type AllTypes struct {
			Int    uint
			String string
			Bytes  []byte
			Bool   bool
			Raw    rlp.RawValue
			Slice  []*AllTypes
			Array  [3]*AllTypes
			Iface  []interface{}
		}
		var v AllTypes
		decodeEncode(input, &v, i)
		i++
	}
	{
		decodeEncode(input, [10]byte{}, i)
		i++
	}
	{
		var v struct {
			Byte [10]byte
			Rool [10]bool
		}
		decodeEncode(input, &v, i)
		i++
	}
	{
		var h types.Header
		decodeEncode(input, &h, i)
		i++
		var b types.Block
		decodeEncode(input, &b, i)
		i++
		var t types.Transaction
		decodeEncode(input, &t, i)
		i++
		var txs types.Transactions
		decodeEncode(input, &txs, i)
		i++
		var rs types.Receipts
		decodeEncode(input, &rs, i)
	}
	{
		i++
		var v struct {
			AnIntPtr  *big.Int
			AnInt     big.Int
			AnU256Ptr *uint256.Int
			AnU256    uint256.Int
			NotAnU256 [4]uint64
		}
		decodeEncode(input, &v, i)
	}
	return 1
}
