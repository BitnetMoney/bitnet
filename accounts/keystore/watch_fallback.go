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

//go:build (darwin && !cgo) || ios || (linux && arm64) || windows || (!darwin && !freebsd && !linux && !netbsd && !solaris)
// +build darwin,!cgo ios linux,arm64 windows !darwin,!freebsd,!linux,!netbsd,!solaris

// This is the fallback implementation of directory watching.
// It is used on unsupported platforms.

package keystore

type watcher struct {
	running  bool
	runEnded bool
}

func newWatcher(*accountCache) *watcher { return new(watcher) }
func (*watcher) start()                 {}
func (*watcher) close()                 {}

// enabled returns false on systems not supported.
func (*watcher) enabled() bool { return false }
