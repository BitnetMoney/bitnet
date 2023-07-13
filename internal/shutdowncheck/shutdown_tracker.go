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

package shutdowncheck

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

// ShutdownTracker is a service that reports previous unclean shutdowns
// upon start. It needs to be started after a successful start-up and stopped
// after a successful shutdown, just before the db is closed.
type ShutdownTracker struct {
	db     ethdb.Database
	stopCh chan struct{}
}

// NewShutdownTracker creates a new ShutdownTracker instance and has
// no other side-effect.
func NewShutdownTracker(db ethdb.Database) *ShutdownTracker {
	return &ShutdownTracker{
		db:     db,
		stopCh: make(chan struct{}),
	}
}

// MarkStartup is to be called in the beginning when the node starts. It will:
// - Push a new startup marker to the db
// - Report previous unclean shutdowns
func (t *ShutdownTracker) MarkStartup() {
	if uncleanShutdowns, discards, err := rawdb.PushUncleanShutdownMarker(t.db); err != nil {
		log.Error("Could not update unclean-shutdown-marker list", "error", err)
	} else {
		if discards > 0 {
			log.Warn("Old unclean shutdowns found", "count", discards)
		}
		for _, tstamp := range uncleanShutdowns {
			t := time.Unix(int64(tstamp), 0)
			log.Warn("Unclean shutdown detected", "booted", t,
				"age", common.PrettyAge(t))
		}
	}
}

// Start runs an event loop that updates the current marker's timestamp every 5 minutes.
func (t *ShutdownTracker) Start() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rawdb.UpdateUncleanShutdownMarker(t.db)
			case <-t.stopCh:
				return
			}
		}
	}()
}

// Stop will stop the update loop and clear the current marker.
func (t *ShutdownTracker) Stop() {
	// Stop update loop.
	t.stopCh <- struct{}{}
	// Clear last marker.
	rawdb.PopUncleanShutdownMarker(t.db)
}
