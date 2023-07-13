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

package asm

import (
	"reflect"
	"testing"
)

func lexAll(src string) []token {
	ch := Lex([]byte(src), false)

	var tokens []token
	for i := range ch {
		tokens = append(tokens, i)
	}
	return tokens
}

func TestLexer(t *testing.T) {
	tests := []struct {
		input  string
		tokens []token
	}{
		{
			input:  ";; this is a comment",
			tokens: []token{{typ: lineStart}, {typ: eof}},
		},
		{
			input:  "0x12345678",
			tokens: []token{{typ: lineStart}, {typ: number, text: "0x12345678"}, {typ: eof}},
		},
		{
			input:  "0x123ggg",
			tokens: []token{{typ: lineStart}, {typ: number, text: "0x123"}, {typ: element, text: "ggg"}, {typ: eof}},
		},
		{
			input:  "12345678",
			tokens: []token{{typ: lineStart}, {typ: number, text: "12345678"}, {typ: eof}},
		},
		{
			input:  "123abc",
			tokens: []token{{typ: lineStart}, {typ: number, text: "123"}, {typ: element, text: "abc"}, {typ: eof}},
		},
		{
			input:  "0123abc",
			tokens: []token{{typ: lineStart}, {typ: number, text: "0123"}, {typ: element, text: "abc"}, {typ: eof}},
		},
		{
			input:  "00123abc",
			tokens: []token{{typ: lineStart}, {typ: number, text: "00123"}, {typ: element, text: "abc"}, {typ: eof}},
		},
		{
			input:  "@foo",
			tokens: []token{{typ: lineStart}, {typ: label, text: "foo"}, {typ: eof}},
		},
		{
			input:  "@label123",
			tokens: []token{{typ: lineStart}, {typ: label, text: "label123"}, {typ: eof}},
		},
	}

	for _, test := range tests {
		tokens := lexAll(test.input)
		if !reflect.DeepEqual(tokens, test.tokens) {
			t.Errorf("input %q\ngot:  %+v\nwant: %+v", test.input, tokens, test.tokens)
		}
	}
}
