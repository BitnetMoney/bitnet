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

// Package abi implements the Ethereum ABI (Application Binary
// Interface).
//
// The Ethereum ABI is strongly typed, known at compile time
// and static. This ABI will handle basic type casting; unsigned
// to signed and visa versa. It does not handle slice casting such
// as unsigned slice to signed slice. Bit size type casting is also
// handled. ints with a bit size of 32 will be properly cast to int256,
// etc.
package abi
