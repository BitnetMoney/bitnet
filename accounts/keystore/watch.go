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

//go:build (darwin && !ios && cgo) || freebsd || (linux && !arm64) || netbsd || solaris
// +build darwin,!ios,cgo freebsd linux,!arm64 netbsd solaris

package keystore

import (
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	ac       *accountCache
	running  bool // set to true when runloop begins
	runEnded bool // set to true when runloop ends
	starting bool // set to true prior to runloop starting
	quit     chan struct{}
}

func newWatcher(ac *accountCache) *watcher {
	return &watcher{
		ac:   ac,
		quit: make(chan struct{}),
	}
}

// enabled returns false on systems not supported.
func (*watcher) enabled() bool { return true }

// starts the watcher loop in the background.
// Start a watcher in the background if that's not already in progress.
// The caller must hold w.ac.mu.
func (w *watcher) start() {
	if w.starting || w.running {
		return
	}
	w.starting = true
	go w.loop()
}

func (w *watcher) close() {
	close(w.quit)
}

func (w *watcher) loop() {
	defer func() {
		w.ac.mu.Lock()
		w.running = false
		w.starting = false
		w.runEnded = true
		w.ac.mu.Unlock()
	}()
	logger := log.New("path", w.ac.keydir)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("Failed to start filesystem watcher", "err", err)
		return
	}
	defer watcher.Close()
	if err := watcher.Add(w.ac.keydir); err != nil {
		logger.Warn("Failed to watch keystore folder", "err", err)
		return
	}

	logger.Trace("Started watching keystore folder", "folder", w.ac.keydir)
	defer logger.Trace("Stopped watching keystore folder")

	w.ac.mu.Lock()
	w.running = true
	w.ac.mu.Unlock()

	// Wait for file system events and reload.
	// When an event occurs, the reload call is delayed a bit so that
	// multiple events arriving quickly only cause a single reload.
	var (
		debounceDuration = 500 * time.Millisecond
		rescanTriggered  = false
		debounce         = time.NewTimer(0)
	)
	// Ignore initial trigger
	if !debounce.Stop() {
		<-debounce.C
	}
	defer debounce.Stop()
	for {
		select {
		case <-w.quit:
			return
		case _, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Trigger the scan (with delay), if not already triggered
			if !rescanTriggered {
				debounce.Reset(debounceDuration)
				rescanTriggered = true
			}
			// The fsnotify library does provide more granular event-info, it
			// would be possible to refresh individual affected files instead
			// of scheduling a full rescan. For most cases though, the
			// full rescan is quick and obviously simplest.
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Info("Filsystem watcher error", "err", err)
		case <-debounce.C:
			w.ac.scanAccounts()
			rescanTriggered = false
		}
	}
}
