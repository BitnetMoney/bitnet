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

package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBatchLimits(t *testing.T) {
	server := newTestServer()
	defer server.Stop()
	server.SetBatchLimits(3, 100)
	httpsrv := httptest.NewServer(server)
	defer httpsrv.Close()

	// batch builds n test_echo calls with ids 1..n; each result is about 60 bytes.
	batch := func(n int) string {
		calls := make([]string, n)
		for i := range calls {
			calls[i] = fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"method":"test_echo","params":["%s",1]}`, i+1, strings.Repeat("x", 30))
		}
		return "[" + strings.Join(calls, ",") + "]"
	}
	post := func(body string) []*jsonrpcMessage {
		t.Helper()
		resp, err := http.Post(httpsrv.URL, "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var msgs []*jsonrpcMessage
		if err := json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
			t.Fatal(err)
		}
		return msgs
	}

	// Within both limits: every call is answered.
	if msgs := post(batch(1)); len(msgs) != 1 || msgs[0].Error != nil {
		t.Errorf("small batch: %v", msgs)
	}

	// Over the call limit: one error for the whole batch, with the first call's id.
	msgs := post(batch(4))
	if len(msgs) != 1 || msgs[0].Error == nil || msgs[0].Error.Message != errMsgBatchTooLarge || string(msgs[0].ID) != "1" {
		t.Fatalf("batch of 4 with a limit of 3: %+v", msgs)
	}

	// Over the response size limit: the second response passes 100 bytes, so the
	// third call gets an error instead of running.
	msgs = post(batch(3))
	if len(msgs) != 3 {
		t.Fatalf("got %d responses, want 3", len(msgs))
	}
	if msgs[0].Error != nil || msgs[1].Error != nil {
		t.Errorf("first two calls should succeed: %v %v", msgs[0].Error, msgs[1].Error)
	}
	if msgs[2].Error == nil || msgs[2].Error.Code != errcodeResponseTooLarge || msgs[2].Error.Message != errMsgResponseTooLarge {
		t.Errorf("third call: %+v, want response too large", msgs[2].Error)
	}
}

// Without limits, as for an RPC client, large batches still work.
func TestBatchWithoutLimits(t *testing.T) {
	server := newTestServer()
	defer server.Stop()
	httpsrv := httptest.NewServer(server)
	defer httpsrv.Close()

	calls := make([]string, 2000)
	for i := range calls {
		calls[i] = fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"method":"test_echo","params":["x",1]}`, i+1)
	}
	resp, err := http.Post(httpsrv.URL, "application/json", strings.NewReader("["+strings.Join(calls, ",")+"]"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var msgs []*jsonrpcMessage
	if err := json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2000 {
		t.Errorf("got %d responses, want 2000", len(msgs))
	}
}
