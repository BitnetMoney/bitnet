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

package utils

import "testing"

func TestExecQueue(t *testing.T) {
	var (
		N        = 10000
		q        = NewExecQueue(N)
		counter  int
		execd    = make(chan int)
		testexit = make(chan struct{})
	)
	defer q.Quit()
	defer close(testexit)

	check := func(state string, wantOK bool) {
		c := counter
		counter++
		qf := func() {
			select {
			case execd <- c:
			case <-testexit:
			}
		}
		if q.CanQueue() != wantOK {
			t.Fatalf("CanQueue() == %t for %s", !wantOK, state)
		}
		if q.Queue(qf) != wantOK {
			t.Fatalf("Queue() == %t for %s", !wantOK, state)
		}
	}

	for i := 0; i < N; i++ {
		check("queue below cap", true)
	}
	check("full queue", false)
	for i := 0; i < N; i++ {
		if c := <-execd; c != i {
			t.Fatal("execution out of order")
		}
	}
	q.Quit()
	check("closed queue", false)
}
