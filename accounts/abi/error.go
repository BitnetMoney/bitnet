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

package abi

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Error struct {
	Name   string
	Inputs Arguments
	str    string

	// Sig contains the string signature according to the ABI spec.
	// e.g.	 error foo(uint32 a, int b) = "foo(uint32,int256)"
	// Please note that "int" is substitute for its canonical representation "int256"
	Sig string

	// ID returns the canonical representation of the error's signature used by the
	// abi definition to identify event names and types.
	ID common.Hash
}

func NewError(name string, inputs Arguments) Error {
	// sanitize inputs to remove inputs without names
	// and precompute string and sig representation.
	names := make([]string, len(inputs))
	types := make([]string, len(inputs))
	for i, input := range inputs {
		if input.Name == "" {
			inputs[i] = Argument{
				Name:    fmt.Sprintf("arg%d", i),
				Indexed: input.Indexed,
				Type:    input.Type,
			}
		} else {
			inputs[i] = input
		}
		// string representation
		names[i] = fmt.Sprintf("%v %v", input.Type, inputs[i].Name)
		if input.Indexed {
			names[i] = fmt.Sprintf("%v indexed %v", input.Type, inputs[i].Name)
		}
		// sig representation
		types[i] = input.Type.String()
	}

	str := fmt.Sprintf("error %v(%v)", name, strings.Join(names, ", "))
	sig := fmt.Sprintf("%v(%v)", name, strings.Join(types, ","))
	id := common.BytesToHash(crypto.Keccak256([]byte(sig)))

	return Error{
		Name:   name,
		Inputs: inputs,
		str:    str,
		Sig:    sig,
		ID:     id,
	}
}

func (e *Error) String() string {
	return e.str
}

func (e *Error) Unpack(data []byte) (interface{}, error) {
	if len(data) < 4 {
		return "", errors.New("invalid data for unpacking")
	}
	if !bytes.Equal(data[:4], e.ID[:4]) {
		return "", errors.New("invalid data for unpacking")
	}
	return e.Inputs.Unpack(data[4:])
}
