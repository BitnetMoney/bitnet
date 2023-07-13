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

func ExampleNewSubscription() {
	// Create a subscription that sends 10 integers on ch.
	ch := make(chan int)
	sub := event.NewSubscription(func(quit <-chan struct{}) error {
		for i := 0; i < 10; i++ {
			select {
			case ch <- i:
			case <-quit:
				fmt.Println("unsubscribed")
				return nil
			}
		}
		return nil
	})

	// This is the consumer. It reads 5 integers, then aborts the subscription.
	// Note that Unsubscribe waits until the producer has shut down.
	for i := range ch {
		fmt.Println(i)
		if i == 4 {
			sub.Unsubscribe()
			break
		}
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// unsubscribed
}
