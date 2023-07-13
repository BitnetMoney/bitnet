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

package server

import (
	"github.com/ethereum/go-ethereum/metrics"
)

var (
	totalActiveCapacityGauge = metrics.NewRegisteredGauge("vflux/server/active/capacity", nil)
	totalActiveCountGauge    = metrics.NewRegisteredGauge("vflux/server/active/count", nil)
	totalInactiveCountGauge  = metrics.NewRegisteredGauge("vflux/server/inactive/count", nil)

	clientConnectedMeter    = metrics.NewRegisteredMeter("vflux/server/clientEvent/connected", nil)
	clientActivatedMeter    = metrics.NewRegisteredMeter("vflux/server/clientEvent/activated", nil)
	clientDeactivatedMeter  = metrics.NewRegisteredMeter("vflux/server/clientEvent/deactivated", nil)
	clientDisconnectedMeter = metrics.NewRegisteredMeter("vflux/server/clientEvent/disconnected", nil)

	capacityQueryZeroMeter    = metrics.NewRegisteredMeter("vflux/server/capQueryZero", nil)
	capacityQueryNonZeroMeter = metrics.NewRegisteredMeter("vflux/server/capQueryNonZero", nil)
)
