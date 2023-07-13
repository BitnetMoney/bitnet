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

// Package syncx contains exotic synchronization primitives.
package syncx

// ClosableMutex is a mutex that can also be closed.
// Once closed, it can never be taken again.
type ClosableMutex struct {
	ch chan struct{}
}

func NewClosableMutex() *ClosableMutex {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &ClosableMutex{ch}
}

// TryLock attempts to lock cm.
// If the mutex is closed, TryLock returns false.
func (cm *ClosableMutex) TryLock() bool {
	_, ok := <-cm.ch
	return ok
}

// MustLock locks cm.
// If the mutex is closed, MustLock panics.
func (cm *ClosableMutex) MustLock() {
	_, ok := <-cm.ch
	if !ok {
		panic("mutex closed")
	}
}

// Unlock unlocks cm.
func (cm *ClosableMutex) Unlock() {
	select {
	case cm.ch <- struct{}{}:
	default:
		panic("Unlock of already-unlocked ClosableMutex")
	}
}

// Close locks the mutex, then closes it.
func (cm *ClosableMutex) Close() {
	_, ok := <-cm.ch
	if !ok {
		panic("Close of already-closed ClosableMutex")
	}
	close(cm.ch)
}
