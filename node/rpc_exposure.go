// Copyright 2026 The Bitnet community authors
// This file is part of the Bitnet library.
//
// The Bitnet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Bitnet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

package node

import (
	"fmt"
	"net"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

// unsafeModules control the node or reveal its internals: peers and data
// directory (admin), tracing and profiling (debug), accounts (personal) and
// mining (miner). Anyone who can reach them can use them.
var unsafeModules = []string{"admin", "debug", "personal", "miner"}

// isLoopback reports whether a listening host accepts only local connections.
func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// checkRPCExposure refuses to serve unsafe modules over HTTP or WebSocket on an
// address other machines can reach, unless AllowUnsafePublicRPC is set.
func checkRPCExposure(conf *Config) error {
	check := func(transport, host string, modules []string) error {
		if host == "" || isLoopback(host) {
			return nil
		}
		for _, m := range modules {
			for _, unsafe := range unsafeModules {
				if strings.EqualFold(strings.TrimSpace(m), unsafe) {
					if conf.AllowUnsafePublicRPC {
						log.Warn("Serving an unsafe API on a public address", "api", unsafe, "transport", transport, "host", host)
						continue
					}
					return fmt.Errorf("refusing to serve the %s API over %s on %s, an address other machines can reach: "+
						"it lets anyone who can connect control the node or read its internals. "+
						"Serve it on 127.0.0.1 instead, or pass --rpc.allow-unsafe-public if this is really intended", unsafe, transport, host)
				}
			}
		}
		return nil
	}
	if err := check("HTTP", conf.HTTPHost, conf.HTTPModules); err != nil {
		return err
	}
	return check("WebSocket", conf.WSHost, conf.WSModules)
}
