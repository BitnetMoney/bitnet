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
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/forkid"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/urfave/cli/v2"
)

var (
	nodesetCommand = &cli.Command{
		Name:  "nodeset",
		Usage: "Node set tools",
		Subcommands: []*cli.Command{
			nodesetInfoCommand,
			nodesetFilterCommand,
		},
	}
	nodesetInfoCommand = &cli.Command{
		Name:      "info",
		Usage:     "Shows statistics about a node set",
		Action:    nodesetInfo,
		ArgsUsage: "<nodes.json>",
	}
	nodesetFilterCommand = &cli.Command{
		Name:      "filter",
		Usage:     "Filters a node set",
		Action:    nodesetFilter,
		ArgsUsage: "<nodes.json> filters..",

		SkipFlagParsing: true,
	}
)

func nodesetInfo(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("need nodes file as argument")
	}

	ns := loadNodesJSON(ctx.Args().First())
	fmt.Printf("Set contains %d nodes.\n", len(ns))
	showAttributeCounts(ns)
	return nil
}

// showAttributeCounts prints the distribution of ENR attributes in a node set.
func showAttributeCounts(ns nodeSet) {
	attrcount := make(map[string]int)
	var attrlist []interface{}
	for _, n := range ns {
		r := n.N.Record()
		attrlist = r.AppendElements(attrlist[:0])[1:]
		for i := 0; i < len(attrlist); i += 2 {
			key := attrlist[i].(string)
			attrcount[key]++
		}
	}

	var keys []string
	var maxlength int
	for key := range attrcount {
		keys = append(keys, key)
		if len(key) > maxlength {
			maxlength = len(key)
		}
	}
	sort.Strings(keys)
	fmt.Println("ENR attribute counts:")
	for _, key := range keys {
		fmt.Printf("%s%s: %d\n", strings.Repeat(" ", maxlength-len(key)+1), key, attrcount[key])
	}
}

func nodesetFilter(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("need nodes file as argument")
	}
	// Parse -limit.
	limit, err := parseFilterLimit(ctx.Args().Tail())
	if err != nil {
		return err
	}
	// Parse the filters.
	filter, err := andFilter(ctx.Args().Tail())
	if err != nil {
		return err
	}

	// Load nodes and apply filters.
	ns := loadNodesJSON(ctx.Args().First())
	result := make(nodeSet)
	for id, n := range ns {
		if filter(n) {
			result[id] = n
		}
	}
	if limit >= 0 {
		result = result.topN(limit)
	}
	writeNodesJSON("-", result)
	return nil
}

type nodeFilter func(nodeJSON) bool

type nodeFilterC struct {
	narg int
	fn   func([]string) (nodeFilter, error)
}

var filterFlags = map[string]nodeFilterC{
	"-limit":       {1, trueFilter}, // needed to skip over -limit
	"-ip":          {1, ipFilter},
	"-min-age":     {1, minAgeFilter},
	"-eth-network": {1, ethFilter},
	"-les-server":  {0, lesFilter},
	"-snap":        {0, snapFilter},
}

// parseFilters parses nodeFilters from args.
func parseFilters(args []string) ([]nodeFilter, error) {
	var filters []nodeFilter
	for len(args) > 0 {
		fc, ok := filterFlags[args[0]]
		if !ok {
			return nil, fmt.Errorf("invalid filter %q", args[0])
		}
		if len(args)-1 < fc.narg {
			return nil, fmt.Errorf("filter %q wants %d arguments, have %d", args[0], fc.narg, len(args)-1)
		}
		filter, err := fc.fn(args[1 : 1+fc.narg])
		if err != nil {
			return nil, fmt.Errorf("%s: %v", args[0], err)
		}
		filters = append(filters, filter)
		args = args[1+fc.narg:]
	}
	return filters, nil
}

// parseFilterLimit parses the -limit option in args. It returns -1 if there is no limit.
func parseFilterLimit(args []string) (int, error) {
	limit := -1
	for i, arg := range args {
		if arg == "-limit" {
			if i == len(args)-1 {
				return -1, errors.New("-limit requires an argument")
			}
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				return -1, fmt.Errorf("invalid -limit %q", args[i+1])
			}
			limit = n
		}
	}
	return limit, nil
}

// andFilter parses node filters in args and returns a single filter that requires all
// of them to match.
func andFilter(args []string) (nodeFilter, error) {
	checks, err := parseFilters(args)
	if err != nil {
		return nil, err
	}
	f := func(n nodeJSON) bool {
		for _, filter := range checks {
			if !filter(n) {
				return false
			}
		}
		return true
	}
	return f, nil
}

func trueFilter(args []string) (nodeFilter, error) {
	return func(n nodeJSON) bool { return true }, nil
}

func ipFilter(args []string) (nodeFilter, error) {
	_, cidr, err := net.ParseCIDR(args[0])
	if err != nil {
		return nil, err
	}
	f := func(n nodeJSON) bool { return cidr.Contains(n.N.IP()) }
	return f, nil
}

func minAgeFilter(args []string) (nodeFilter, error) {
	minage, err := time.ParseDuration(args[0])
	if err != nil {
		return nil, err
	}
	f := func(n nodeJSON) bool {
		age := n.LastResponse.Sub(n.FirstResponse)
		return age >= minage
	}
	return f, nil
}

func ethFilter(args []string) (nodeFilter, error) {
	var filter forkid.Filter
	switch args[0] {
	case "mainnet":
		filter = forkid.NewStaticFilter(params.MainnetChainConfig, params.MainnetGenesisHash)
	case "rinkeby":
		filter = forkid.NewStaticFilter(params.RinkebyChainConfig, params.RinkebyGenesisHash)
	case "goerli":
		filter = forkid.NewStaticFilter(params.GoerliChainConfig, params.GoerliGenesisHash)
	case "sepolia":
		filter = forkid.NewStaticFilter(params.SepoliaChainConfig, params.SepoliaGenesisHash)
	default:
		return nil, fmt.Errorf("unknown network %q", args[0])
	}

	f := func(n nodeJSON) bool {
		var eth struct {
			ForkID forkid.ID
			Tail   []rlp.RawValue `rlp:"tail"`
		}
		if n.N.Load(enr.WithEntry("eth", &eth)) != nil {
			return false
		}
		return filter(eth.ForkID) == nil
	}
	return f, nil
}

func lesFilter(args []string) (nodeFilter, error) {
	f := func(n nodeJSON) bool {
		var les struct {
			Tail []rlp.RawValue `rlp:"tail"`
		}
		return n.N.Load(enr.WithEntry("les", &les)) == nil
	}
	return f, nil
}

func snapFilter(args []string) (nodeFilter, error) {
	f := func(n nodeJSON) bool {
		var snap struct {
			Tail []rlp.RawValue `rlp:"tail"`
		}
		return n.N.Load(enr.WithEntry("snap", &snap)) == nil
	}
	return f, nil
}
