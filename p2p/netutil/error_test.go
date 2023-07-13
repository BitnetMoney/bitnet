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

package netutil

import (
	"net"
	"testing"
	"time"
)

// This test checks that isPacketTooBig correctly identifies
// errors that result from receiving a UDP packet larger
// than the supplied receive buffer.
func TestIsPacketTooBig(t *testing.T) {
	listener, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	sender, err := net.Dial("udp", listener.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer sender.Close()

	sendN := 1800
	recvN := 300
	for i := 0; i < 20; i++ {
		go func() {
			buf := make([]byte, sendN)
			for i := range buf {
				buf[i] = byte(i)
			}
			sender.Write(buf)
		}()

		buf := make([]byte, recvN)
		listener.SetDeadline(time.Now().Add(1 * time.Second))
		n, _, err := listener.ReadFrom(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			if !isPacketTooBig(err) {
				t.Fatalf("unexpected read error: %v", err)
			}
			continue
		}
		if n != recvN {
			t.Fatalf("short read: %d, want %d", n, recvN)
		}
		for i := range buf {
			if buf[i] != byte(i) {
				t.Fatalf("error in pattern")
			}
		}
	}
}
