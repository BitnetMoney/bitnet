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

package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/tests/fuzzers/snap"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: debug <file>\n")
		os.Exit(1)
	}
	crasher := os.Args[1]
	data, err := os.ReadFile(crasher)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading crasher %v: %v", crasher, err)
		os.Exit(1)
	}
	snap.FuzzTrieNodes(data)
}
