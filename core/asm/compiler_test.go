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
	"testing"
)

func TestCompiler(t *testing.T) {
	tests := []struct {
		input, output string
	}{
		{
			input: `
	GAS
	label:
	PUSH @label
`,
			output: "5a5b6300000001",
		},
		{
			input: `
	PUSH @label
	label:
`,
			output: "63000000055b",
		},
		{
			input: `
	PUSH @label
	JUMP
	label:
`,
			output: "6300000006565b",
		},
		{
			input: `
	JUMP @label
	label:
`,
			output: "6300000006565b",
		},
	}
	for _, test := range tests {
		ch := Lex([]byte(test.input), false)
		c := NewCompiler(false)
		c.Feed(ch)
		output, err := c.Compile()
		if len(err) != 0 {
			t.Errorf("compile error: %v\ninput: %s", err, test.input)
			continue
		}
		if output != test.output {
			t.Errorf("incorrect output\ninput: %sgot:  %s\nwant: %s\n", test.input, output, test.output)
		}
	}
}
