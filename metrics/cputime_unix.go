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

//go:build !windows && !js
// +build !windows,!js

package metrics

import (
	syscall "golang.org/x/sys/unix"

	"github.com/ethereum/go-ethereum/log"
)

// getProcessCPUTime retrieves the process' CPU time since program startup.
func getProcessCPUTime() float64 {
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		log.Warn("Failed to retrieve CPU time", "err", err)
		return 0
	}
	return float64(usage.Utime.Sec+usage.Stime.Sec) + float64(usage.Utime.Usec+usage.Stime.Usec)/1000000 //nolint:unconvert
}
