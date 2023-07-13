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

package metrics

// CPUStats is the system and process CPU stats.
// All values are in seconds.
type CPUStats struct {
	GlobalTime float64 // Time spent by the CPU working on all processes
	GlobalWait float64 // Time spent by waiting on disk for all processes
	LocalTime  float64 // Time spent by the CPU working on this process
}
