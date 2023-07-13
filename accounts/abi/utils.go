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

import "fmt"

// ResolveNameConflict returns the next available name for a given thing.
// This helper can be used for lots of purposes:
//
//   - In solidity function overloading is supported, this function can fix
//     the name conflicts of overloaded functions.
//   - In golang binding generation, the parameter(in function, event, error,
//     and struct definition) name will be converted to camelcase style which
//     may eventually lead to name conflicts.
//
// Name conflicts are mostly resolved by adding number suffix. e.g. if the abi contains
// Methods "send" and "send1", ResolveNameConflict would return "send2" for input "send".
func ResolveNameConflict(rawName string, used func(string) bool) string {
	name := rawName
	ok := used(name)
	for idx := 0; ok; idx++ {
		name = fmt.Sprintf("%s%d", rawName, idx)
		ok = used(name)
	}
	return name
}
