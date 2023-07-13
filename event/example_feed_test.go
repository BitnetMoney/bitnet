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

package event_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/event"
)

func ExampleFeed_acknowledgedEvents() {
	// This example shows how the return value of Send can be used for request/reply
	// interaction between event consumers and producers.
	var feed event.Feed
	type ackedEvent struct {
		i   int
		ack chan<- struct{}
	}

	// Consumers wait for events on the feed and acknowledge processing.
	done := make(chan struct{})
	defer close(done)
	for i := 0; i < 3; i++ {
		ch := make(chan ackedEvent, 100)
		sub := feed.Subscribe(ch)
		go func() {
			defer sub.Unsubscribe()
			for {
				select {
				case ev := <-ch:
					fmt.Println(ev.i) // "process" the event
					ev.ack <- struct{}{}
				case <-done:
					return
				}
			}
		}()
	}

	// The producer sends values of type ackedEvent with increasing values of i.
	// It waits for all consumers to acknowledge before sending the next event.
	for i := 0; i < 3; i++ {
		acksignal := make(chan struct{})
		n := feed.Send(ackedEvent{i, acksignal})
		for ack := 0; ack < n; ack++ {
			<-acksignal
		}
	}
	// Output:
	// 0
	// 0
	// 0
	// 1
	// 1
	// 1
	// 2
	// 2
	// 2
}
