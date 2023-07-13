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

package rlpx

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestReadBufferReset(t *testing.T) {
	reader := bytes.NewReader(hexutil.MustDecode("0x010202030303040505"))
	var b readBuffer

	s1, _ := b.read(reader, 1)
	s2, _ := b.read(reader, 2)
	s3, _ := b.read(reader, 3)

	assert.Equal(t, []byte{1}, s1)
	assert.Equal(t, []byte{2, 2}, s2)
	assert.Equal(t, []byte{3, 3, 3}, s3)

	b.reset()

	s4, _ := b.read(reader, 1)
	s5, _ := b.read(reader, 2)

	assert.Equal(t, []byte{4}, s4)
	assert.Equal(t, []byte{5, 5}, s5)

	s6, err := b.read(reader, 2)

	assert.EqualError(t, err, "EOF")
	assert.Nil(t, s6)
}
