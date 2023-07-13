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

package event

import "fmt"

func ExampleTypeMux() {
	type someEvent struct{ I int }
	type otherEvent struct{ S string }
	type yetAnotherEvent struct{ X, Y int }

	var mux TypeMux

	// Start a subscriber.
	done := make(chan struct{})
	sub := mux.Subscribe(someEvent{}, otherEvent{})
	go func() {
		for event := range sub.Chan() {
			fmt.Printf("Received: %#v\n", event.Data)
		}
		fmt.Println("done")
		close(done)
	}()

	// Post some events.
	mux.Post(someEvent{5})
	mux.Post(yetAnotherEvent{X: 3, Y: 4})
	mux.Post(someEvent{6})
	mux.Post(otherEvent{"whoa"})

	// Stop closes all subscription channels.
	// The subscriber goroutine will print "done"
	// and exit.
	mux.Stop()

	// Wait for subscriber to return.
	<-done

	// Output:
	// Received: event.someEvent{I:5}
	// Received: event.someEvent{I:6}
	// Received: event.otherEvent{S:"whoa"}
	// done
}
