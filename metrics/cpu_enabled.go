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

//go:build !ios && !js
// +build !ios,!js

package metrics

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/shirou/gopsutil/cpu"
)

// ReadCPUStats retrieves the current CPU stats.
func ReadCPUStats(stats *CPUStats) {
	// passing false to request all cpu times
	timeStats, err := cpu.Times(false)
	if err != nil {
		log.Error("Could not read cpu stats", "err", err)
		return
	}
	if len(timeStats) == 0 {
		log.Error("Empty cpu stats")
		return
	}
	// requesting all cpu times will always return an array with only one time stats entry
	timeStat := timeStats[0]
	stats.GlobalTime = timeStat.User + timeStat.Nice + timeStat.System
	stats.GlobalWait = timeStat.Iowait
	stats.LocalTime = getProcessCPUTime()
}
