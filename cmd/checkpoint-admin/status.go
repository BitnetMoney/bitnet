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

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

var commandStatus = &cli.Command{
	Name:  "status",
	Usage: "Fetches the signers and checkpoint status of the oracle contract",
	Flags: []cli.Flag{
		nodeURLFlag,
	},
	Action: status,
}

// status fetches the admin list of specified registrar contract.
func status(ctx *cli.Context) error {
	// Create a wrapper around the checkpoint oracle contract
	addr, oracle := newContract(newRPCClient(ctx.String(nodeURLFlag.Name)))
	fmt.Printf("Oracle => %s\n", addr.Hex())
	fmt.Println()

	// Retrieve the list of authorized signers (admins)
	admins, err := oracle.Contract().GetAllAdmin(nil)
	if err != nil {
		return err
	}
	for i, admin := range admins {
		fmt.Printf("Admin %d => %s\n", i+1, admin.Hex())
	}
	fmt.Println()

	// Retrieve the latest checkpoint
	index, checkpoint, height, err := oracle.Contract().GetLatestCheckpoint(nil)
	if err != nil {
		return err
	}
	fmt.Printf("Checkpoint (published at #%d) %d => %s\n", height, index, common.Hash(checkpoint).Hex())

	return nil
}
