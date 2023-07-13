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

package vm

import (
	"errors"

	"github.com/ethereum/go-ethereum/params"
)

// LookupInstructionSet returns the instructionset for the fork configured by
// the rules.
func LookupInstructionSet(rules params.Rules) (JumpTable, error) {
	switch {
	case rules.IsPrague:
		return newShanghaiInstructionSet(), errors.New("prague-fork not defined yet")
	case rules.IsCancun:
		return newShanghaiInstructionSet(), errors.New("cancun-fork not defined yet")
	case rules.IsShanghai:
		return newShanghaiInstructionSet(), nil
	case rules.IsMerge:
		return newMergeInstructionSet(), nil
	case rules.IsLondon:
		return newLondonInstructionSet(), nil
	case rules.IsBerlin:
		return newBerlinInstructionSet(), nil
	case rules.IsIstanbul:
		return newIstanbulInstructionSet(), nil
	case rules.IsConstantinople:
		return newConstantinopleInstructionSet(), nil
	case rules.IsByzantium:
		return newByzantiumInstructionSet(), nil
	case rules.IsEIP158:
		return newSpuriousDragonInstructionSet(), nil
	case rules.IsEIP150:
		return newTangerineWhistleInstructionSet(), nil
	case rules.IsHomestead:
		return newHomesteadInstructionSet(), nil
	}
	return newFrontierInstructionSet(), nil
}

// Stack returns the mininum and maximum stack requirements.
func (op *operation) Stack() (int, int) {
	return op.minStack, op.maxStack
}

// HasCost returns true if the opcode has a cost. Opcodes which do _not_ have
// a cost assigned are one of two things:
// - undefined, a.k.a invalid opcodes,
// - the STOP opcode.
// This method can thus be used to check if an opcode is "Invalid (or STOP)".
func (op *operation) HasCost() bool {
	// Ideally, we'd check this:
	//	return op.execute == opUndefined
	// However, go-lang does now allow that. So we'll just check some other
	// 'indicators' that this is an invalid op. Alas, STOP is impossible to
	// filter out
	return op.dynamicGas != nil || op.constantGas != 0
}
