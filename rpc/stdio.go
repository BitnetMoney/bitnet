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

package rpc

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"time"
)

// DialStdIO creates a client on stdin/stdout.
func DialStdIO(ctx context.Context) (*Client, error) {
	return DialIO(ctx, os.Stdin, os.Stdout)
}

// DialIO creates a client which uses the given IO channels
func DialIO(ctx context.Context, in io.Reader, out io.Writer) (*Client, error) {
	return newClient(ctx, newClientTransportIO(in, out))
}

func newClientTransportIO(in io.Reader, out io.Writer) reconnectFunc {
	return func(context.Context) (ServerCodec, error) {
		return NewCodec(stdioConn{
			in:  in,
			out: out,
		}), nil
	}
}

type stdioConn struct {
	in  io.Reader
	out io.Writer
}

func (io stdioConn) Read(b []byte) (n int, err error) {
	return io.in.Read(b)
}

func (io stdioConn) Write(b []byte) (n int, err error) {
	return io.out.Write(b)
}

func (io stdioConn) Close() error {
	return nil
}

func (io stdioConn) RemoteAddr() string {
	return "/dev/stdin"
}

func (io stdioConn) SetWriteDeadline(t time.Time) error {
	return &net.OpError{Op: "set", Net: "stdio", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}
