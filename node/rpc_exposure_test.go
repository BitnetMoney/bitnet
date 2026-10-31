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
	"strings"
	"testing"
)

func TestCheckRPCExposure(t *testing.T) {
	tests := []struct {
		name    string
		conf    Config
		wantErr string // empty: allowed
	}{
		{"safe modules on a public address", Config{HTTPHost: "0.0.0.0", HTTPModules: []string{"eth", "net", "web3", "bitnet"}}, ""},
		{"admin on 0.0.0.0", Config{HTTPHost: "0.0.0.0", HTTPModules: []string{"eth", "admin"}}, "admin API over HTTP"},
		{"debug on a LAN address", Config{HTTPHost: "192.168.1.5", HTTPModules: []string{"debug"}}, "debug API over HTTP"},
		{"miner over WebSocket on ::", Config{WSHost: "::", WSModules: []string{"miner"}}, "miner API over WebSocket"},
		{"module name in other case", Config{HTTPHost: "0.0.0.0", HTTPModules: []string{"Personal"}}, "personal API"},
		{"admin on 127.0.0.1", Config{HTTPHost: "127.0.0.1", HTTPModules: []string{"admin", "debug"}}, ""},
		{"admin on localhost", Config{HTTPHost: "localhost", HTTPModules: []string{"admin"}}, ""},
		{"admin on ::1", Config{WSHost: "::1", WSModules: []string{"admin"}}, ""},
		{"HTTP disabled", Config{HTTPModules: []string{"admin"}}, ""},
		{"explicitly allowed", Config{HTTPHost: "0.0.0.0", HTTPModules: []string{"admin"}, AllowUnsafePublicRPC: true}, ""},
	}
	for _, tt := range tests {
		err := checkRPCExposure(&tt.conf)
		switch {
		case tt.wantErr == "" && err != nil:
			t.Errorf("%s: unexpected error: %v", tt.name, err)
		case tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)):
			t.Errorf("%s: got %v, want an error containing %q", tt.name, err, tt.wantErr)
		}
	}
}

// A node configured to serve admin on a public address does not start.
func TestNodeRefusesPublicAdmin(t *testing.T) {
	conf := testNodeConfig()
	conf.HTTPHost = "0.0.0.0"
	conf.HTTPPort = 0
	conf.HTTPModules = []string{"eth", "admin"}
	stack, err := New(conf)
	if err != nil {
		t.Fatal(err)
	}
	defer stack.Close()
	if err := stack.Start(); err == nil || !strings.Contains(err.Error(), "admin API") {
		t.Fatalf("Start: got %v, want the admin API refused", err)
	}
}
