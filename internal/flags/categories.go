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

package flags

import "github.com/urfave/cli/v2"

const (
	EthCategory        = "ETHEREUM"
	LightCategory      = "LIGHT CLIENT"
	DevCategory        = "DEVELOPER CHAIN"
	EthashCategory     = "ETHASH"
	TxPoolCategory     = "TRANSACTION POOL"
	PerfCategory       = "PERFORMANCE TUNING"
	AccountCategory    = "ACCOUNT"
	APICategory        = "API AND CONSOLE"
	NetworkingCategory = "NETWORKING"
	MinerCategory      = "MINER"
	GasPriceCategory   = "GAS PRICE ORACLE"
	VMCategory         = "VIRTUAL MACHINE"
	LoggingCategory    = "LOGGING AND DEBUGGING"
	MetricsCategory    = "METRICS AND STATS"
	MiscCategory       = "MISC"
	DeprecatedCategory = "ALIASED (deprecated)"
)

func init() {
	cli.HelpFlag.(*cli.BoolFlag).Category = MiscCategory
	cli.VersionFlag.(*cli.BoolFlag).Category = MiscCategory
}
