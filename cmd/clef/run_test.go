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

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/ethereum/go-ethereum/internal/cmdtest"
)

const registeredName = "clef-test"

type testproc struct {
	*cmdtest.TestCmd

	// template variables for expect
	Datadir   string
	Etherbase string
}

func init() {
	reexec.Register(registeredName, func() {
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}

func TestMain(m *testing.M) {
	// check if we have been reexec'd
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

// runClef spawns clef with the given command line args and adds keystore arg.
// This method creates a temporary  keystore folder which will be removed after
// the test exits.
func runClef(t *testing.T, args ...string) *testproc {
	ddir, err := os.MkdirTemp("", "cleftest-*")
	if err != nil {
		return nil
	}
	t.Cleanup(func() {
		os.RemoveAll(ddir)
	})
	return runWithKeystore(t, ddir, args...)
}

// runWithKeystore spawns clef with the given command line args and adds keystore arg.
// This method does _not_ create the keystore folder, but it _does_ add the arg
// to the args.
func runWithKeystore(t *testing.T, keystore string, args ...string) *testproc {
	args = append([]string{"--keystore", keystore}, args...)
	tt := &testproc{Datadir: keystore}
	tt.TestCmd = cmdtest.NewTestCmd(t, tt)
	// Boot "clef". This actually runs the test binary but the TestMain
	// function will prevent any tests from running.
	tt.Run(registeredName, args...)
	return tt
}

func (proc *testproc) input(text string) *testproc {
	proc.TestCmd.InputLine(text)
	return proc
}

/*
// waitForEndpoint waits for the rpc endpoint to appear, or
// aborts after 3 seconds.
func (proc *testproc) waitForEndpoint(t *testing.T) *testproc {
	t.Helper()
	timeout := 3 * time.Second
	ipc := filepath.Join(proc.Datadir, "clef.ipc")

	start := time.Now()
	for time.Since(start) < timeout {
		if _, err := os.Stat(ipc); !errors.Is(err, os.ErrNotExist) {
			t.Logf("endpoint %v opened", ipc)
			return proc
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Logf("stderr: \n%v", proc.StderrText())
	t.Logf("stdout: \n%v", proc.Output())
	t.Fatal("endpoint", ipc, "did not open within", timeout)
	return proc
}
*/
