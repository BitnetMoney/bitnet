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

package dnsdisc

import (
	"errors"
	"fmt"
)

// Entry parse errors.
var (
	errUnknownEntry = errors.New("unknown entry type")
	errNoPubkey     = errors.New("missing public key")
	errBadPubkey    = errors.New("invalid public key")
	errInvalidENR   = errors.New("invalid node record")
	errInvalidChild = errors.New("invalid child hash")
	errInvalidSig   = errors.New("invalid base64 signature")
	errSyntax       = errors.New("invalid syntax")
)

// Resolver/sync errors
var (
	errNoRoot        = errors.New("no valid root found")
	errNoEntry       = errors.New("no valid tree entry found")
	errHashMismatch  = errors.New("hash mismatch")
	errENRInLinkTree = errors.New("enr entry in link tree")
	errLinkInENRTree = errors.New("link entry in ENR tree")
)

type nameError struct {
	name string
	err  error
}

func (err nameError) Error() string {
	if ee, ok := err.err.(entryError); ok {
		return fmt.Sprintf("invalid %s entry at %s: %v", ee.typ, err.name, ee.err)
	}
	return err.name + ": " + err.err.Error()
}

type entryError struct {
	typ string
	err error
}

func (err entryError) Error() string {
	return fmt.Sprintf("invalid %s entry: %v", err.typ, err.err)
}
