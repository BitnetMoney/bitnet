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

package secp256k1

import "testing"

func TestFuzzer(t *testing.T) {
	test := "00000000N0000000/R00000000000000000U0000S0000000mkhP000000000000000U"
	Fuzz([]byte(test))
}
