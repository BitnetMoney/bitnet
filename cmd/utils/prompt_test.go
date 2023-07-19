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

// Package utils contains internal helper functions for go-ethereum commands.
package utils

import (
	"testing"
)

func TestGetPassPhraseWithList(t *testing.T) {
	type args struct {
		text         string
		confirmation bool
		index        int
		passwords    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test1",
			args{
				"text1",
				false,
				0,
				[]string{"zero", "one", "two"},
			},
			"zero",
		},
		{
			"test2",
			args{
				"text2",
				false,
				5,
				[]string{"zero", "one", "two"},
			},
			"two",
		},
		{
			"test3",
			args{
				"text3",
				true,
				1,
				[]string{"zero", "one", "two"},
			},
			"one",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPassPhraseWithList(tt.args.text, tt.args.confirmation, tt.args.index, tt.args.passwords); got != tt.want {
				t.Errorf("GetPassPhraseWithList() = %v, want %v", got, tt.want)
			}
		})
	}
}
