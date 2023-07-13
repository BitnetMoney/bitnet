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

package rlp_test

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

func ExampleEncoderBuffer() {
	var w bytes.Buffer

	// Encode [4, [5, 6]] to w.
	buf := rlp.NewEncoderBuffer(&w)
	l1 := buf.List()
	buf.WriteUint64(4)
	l2 := buf.List()
	buf.WriteUint64(5)
	buf.WriteUint64(6)
	buf.ListEnd(l2)
	buf.ListEnd(l1)

	if err := buf.Flush(); err != nil {
		panic(err)
	}
	fmt.Printf("%X\n", w.Bytes())
	// Output:
	// C404C20506
}
