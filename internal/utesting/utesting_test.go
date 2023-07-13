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

package utesting

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestTest(t *testing.T) {
	tests := []Test{
		{
			Name: "successful test",
			Fn:   func(t *T) {},
		},
		{
			Name: "failing test",
			Fn: func(t *T) {
				t.Log("output")
				t.Error("failed")
			},
		},
		{
			Name: "panicking test",
			Fn: func(t *T) {
				panic("oh no")
			},
		},
	}
	results := RunTests(tests, nil)

	if results[0].Failed || results[0].Output != "" {
		t.Fatalf("wrong result for successful test: %#v", results[0])
	}
	if !results[1].Failed || results[1].Output != "output\nfailed\n" {
		t.Fatalf("wrong result for failing test: %#v", results[1])
	}
	if !results[2].Failed || !strings.HasPrefix(results[2].Output, "panic: oh no\n") {
		t.Fatalf("wrong result for panicking test: %#v", results[2])
	}
}

var outputTests = []Test{
	{
		Name: "TestWithLogs",
		Fn: func(t *T) {
			t.Log("output line 1")
			t.Log("output line 2\noutput line 3")
		},
	},
	{
		Name: "TestNoLogs",
		Fn:   func(t *T) {},
	},
	{
		Name: "FailWithLogs",
		Fn: func(t *T) {
			t.Log("output line 1")
			t.Error("failed 1")
		},
	},
	{
		Name: "FailMessage",
		Fn: func(t *T) {
			t.Error("failed 2")
		},
	},
	{
		Name: "FailNoOutput",
		Fn: func(t *T) {
			t.Fail()
		},
	},
}

func TestOutput(t *testing.T) {
	var buf bytes.Buffer
	RunTests(outputTests, &buf)

	want := regexp.MustCompile(`
^-- RUN TestWithLogs
 output line 1
 output line 2
 output line 3
-- OK TestWithLogs \([^)]+\)
-- OK TestNoLogs \([^)]+\)
-- RUN FailWithLogs
 output line 1
 failed 1
-- FAIL FailWithLogs \([^)]+\)
-- RUN FailMessage
 failed 2
-- FAIL FailMessage \([^)]+\)
-- FAIL FailNoOutput \([^)]+\)
2/5 tests passed.
$`[1:])
	if !want.MatchString(buf.String()) {
		t.Fatalf("output does not match: %q", buf.String())
	}
}

func TestOutputTAP(t *testing.T) {
	var buf bytes.Buffer
	RunTAP(outputTests, &buf)

	want := `
1..5
ok 1 TestWithLogs
# output line 1
# output line 2
# output line 3
ok 2 TestNoLogs
not ok 3 FailWithLogs
# output line 1
# failed 1
not ok 4 FailMessage
# failed 2
not ok 5 FailNoOutput
`
	if buf.String() != want[1:] {
		t.Fatalf("output does not match: %q", buf.String())
	}
}
