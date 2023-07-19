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
	"net"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/urfave/cli/v2"
)

var (
	keyCommand = &cli.Command{
		Name:  "key",
		Usage: "Operations on node keys",
		Subcommands: []*cli.Command{
			keyGenerateCommand,
			keyToIDCommand,
			keyToNodeCommand,
			keyToRecordCommand,
		},
	}
	keyGenerateCommand = &cli.Command{
		Name:      "generate",
		Usage:     "Generates node key files",
		ArgsUsage: "keyfile",
		Action:    genkey,
	}
	keyToIDCommand = &cli.Command{
		Name:      "to-id",
		Usage:     "Creates a node ID from a node key file",
		ArgsUsage: "keyfile",
		Action:    keyToID,
		Flags:     []cli.Flag{},
	}
	keyToNodeCommand = &cli.Command{
		Name:      "to-enode",
		Usage:     "Creates an enode URL from a node key file",
		ArgsUsage: "keyfile",
		Action:    keyToURL,
		Flags:     []cli.Flag{hostFlag, tcpPortFlag, udpPortFlag},
	}
	keyToRecordCommand = &cli.Command{
		Name:      "to-enr",
		Usage:     "Creates an ENR from a node key file",
		ArgsUsage: "keyfile",
		Action:    keyToRecord,
		Flags:     []cli.Flag{hostFlag, tcpPortFlag, udpPortFlag},
	}
)

var (
	hostFlag = &cli.StringFlag{
		Name:  "ip",
		Usage: "IP address of the node",
		Value: "127.0.0.1",
	}
	tcpPortFlag = &cli.IntFlag{
		Name:  "tcp",
		Usage: "TCP port of the node",
		Value: 30303,
	}
	udpPortFlag = &cli.IntFlag{
		Name:  "udp",
		Usage: "UDP port of the node",
		Value: 30303,
	}
)

func genkey(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("need key file as argument")
	}
	file := ctx.Args().Get(0)

	key, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("could not generate key: %v", err)
	}
	return crypto.SaveECDSA(file, key)
}

func keyToID(ctx *cli.Context) error {
	n, err := makeRecord(ctx)
	if err != nil {
		return err
	}
	fmt.Println(n.ID())
	return nil
}

func keyToURL(ctx *cli.Context) error {
	n, err := makeRecord(ctx)
	if err != nil {
		return err
	}
	fmt.Println(n.URLv4())
	return nil
}

func keyToRecord(ctx *cli.Context) error {
	n, err := makeRecord(ctx)
	if err != nil {
		return err
	}
	fmt.Println(n.String())
	return nil
}

func makeRecord(ctx *cli.Context) (*enode.Node, error) {
	if ctx.NArg() != 1 {
		return nil, fmt.Errorf("need key file as argument")
	}

	var (
		file = ctx.Args().Get(0)
		host = ctx.String(hostFlag.Name)
		tcp  = ctx.Int(tcpPortFlag.Name)
		udp  = ctx.Int(udpPortFlag.Name)
	)
	key, err := crypto.LoadECDSA(file)
	if err != nil {
		return nil, err
	}

	var r enr.Record
	if host != "" {
		ip := net.ParseIP(host)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP address %q", host)
		}
		r.Set(enr.IP(ip))
	}
	if udp != 0 {
		r.Set(enr.UDP(udp))
	}
	if tcp != 0 {
		r.Set(enr.TCP(tcp))
	}

	if err := enode.SignV4(&r, key); err != nil {
		return nil, err
	}
	return enode.New(enode.ValidSchemes, &r)
}
