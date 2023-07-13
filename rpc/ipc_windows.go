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

//go:build windows
// +build windows

package rpc

import (
	"context"
	"net"
	"time"

	"gopkg.in/natefinch/npipe.v2"
)

// This is used if the dialing context has no deadline. It is much smaller than the
// defaultDialTimeout because named pipes are local and there is no need to wait so long.
const defaultPipeDialTimeout = 2 * time.Second

// ipcListen will create a named pipe on the given endpoint.
func ipcListen(endpoint string) (net.Listener, error) {
	return npipe.Listen(endpoint)
}

// newIPCConnection will connect to a named pipe with the given endpoint as name.
func newIPCConnection(ctx context.Context, endpoint string) (net.Conn, error) {
	timeout := defaultPipeDialTimeout
	if deadline, ok := ctx.Deadline(); ok {
		timeout = deadline.Sub(time.Now())
		if timeout < 0 {
			timeout = 0
		}
	}
	return npipe.DialTimeout(endpoint, timeout)
}
