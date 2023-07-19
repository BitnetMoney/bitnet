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
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/urfave/cli/v2"
)

var fileFlag = &cli.StringFlag{Name: "file"}

var enrdumpCommand = &cli.Command{
	Name:   "enrdump",
	Usage:  "Pretty-prints node records",
	Action: enrdump,
	Flags: []cli.Flag{
		fileFlag,
	},
}

func enrdump(ctx *cli.Context) error {
	var source string
	if file := ctx.String(fileFlag.Name); file != "" {
		if ctx.NArg() != 0 {
			return fmt.Errorf("can't dump record from command-line argument in -file mode")
		}
		var b []byte
		var err error
		if file == "-" {
			b, err = io.ReadAll(os.Stdin)
		} else {
			b, err = os.ReadFile(file)
		}
		if err != nil {
			return err
		}
		source = string(b)
	} else if ctx.NArg() == 1 {
		source = ctx.Args().First()
	} else {
		return fmt.Errorf("need record as argument")
	}

	r, err := parseRecord(source)
	if err != nil {
		return fmt.Errorf("INVALID: %v", err)
	}
	dumpRecord(os.Stdout, r)
	return nil
}

// dumpRecord creates a human-readable description of the given node record.
func dumpRecord(out io.Writer, r *enr.Record) {
	n, err := enode.New(enode.ValidSchemes, r)
	if err != nil {
		fmt.Fprintf(out, "INVALID: %v\n", err)
	} else {
		fmt.Fprintf(out, "Node ID: %v\n", n.ID())
		dumpNodeURL(out, n)
	}
	kv := r.AppendElements(nil)[1:]
	fmt.Fprintf(out, "Record has sequence number %d and %d key/value pairs.\n", r.Seq(), len(kv)/2)
	fmt.Fprint(out, dumpRecordKV(kv, 2))
}

func dumpNodeURL(out io.Writer, n *enode.Node) {
	var key enode.Secp256k1
	if n.Load(&key) != nil {
		return // no secp256k1 public key
	}
	fmt.Fprintf(out, "URLv4:   %s\n", n.URLv4())
}

func dumpRecordKV(kv []interface{}, indent int) string {
	// Determine the longest key name for alignment.
	var out string
	var longestKey = 0
	for i := 0; i < len(kv); i += 2 {
		key := kv[i].(string)
		if len(key) > longestKey {
			longestKey = len(key)
		}
	}
	// Print the keys, invoking formatters for known keys.
	for i := 0; i < len(kv); i += 2 {
		key := kv[i].(string)
		val := kv[i+1].(rlp.RawValue)
		pad := longestKey - len(key)
		out += strings.Repeat(" ", indent) + strconv.Quote(key) + strings.Repeat(" ", pad+1)
		formatter := attrFormatters[key]
		if formatter == nil {
			formatter = formatAttrRaw
		}
		fmtval, ok := formatter(val)
		if ok {
			out += fmtval + "\n"
		} else {
			out += hex.EncodeToString(val) + " (!)\n"
		}
	}
	return out
}

// parseNode parses a node record and verifies its signature.
func parseNode(source string) (*enode.Node, error) {
	if strings.HasPrefix(source, "enode://") {
		return enode.ParseV4(source)
	}
	r, err := parseRecord(source)
	if err != nil {
		return nil, err
	}
	return enode.New(enode.ValidSchemes, r)
}

// parseRecord parses a node record from hex, base64, or raw binary input.
func parseRecord(source string) (*enr.Record, error) {
	bin := []byte(source)
	if d, ok := decodeRecordHex(bytes.TrimSpace(bin)); ok {
		bin = d
	} else if d, ok := decodeRecordBase64(bytes.TrimSpace(bin)); ok {
		bin = d
	}
	var r enr.Record
	err := rlp.DecodeBytes(bin, &r)
	return &r, err
}

func decodeRecordHex(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("0x")) {
		b = b[2:]
	}
	dec := make([]byte, hex.DecodedLen(len(b)))
	_, err := hex.Decode(dec, b)
	return dec, err == nil
}

func decodeRecordBase64(b []byte) ([]byte, bool) {
	if bytes.HasPrefix(b, []byte("enr:")) {
		b = b[4:]
	}
	dec := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(dec, b)
	return dec[:n], err == nil
}

// attrFormatters contains formatting functions for well-known ENR keys.
var attrFormatters = map[string]func(rlp.RawValue) (string, bool){
	"id":   formatAttrString,
	"ip":   formatAttrIP,
	"ip6":  formatAttrIP,
	"tcp":  formatAttrUint,
	"tcp6": formatAttrUint,
	"udp":  formatAttrUint,
	"udp6": formatAttrUint,
}

func formatAttrRaw(v rlp.RawValue) (string, bool) {
	s := hex.EncodeToString(v)
	return s, true
}

func formatAttrString(v rlp.RawValue) (string, bool) {
	content, _, err := rlp.SplitString(v)
	return strconv.Quote(string(content)), err == nil
}

func formatAttrIP(v rlp.RawValue) (string, bool) {
	content, _, err := rlp.SplitString(v)
	if err != nil || len(content) != 4 && len(content) != 6 {
		return "", false
	}
	return net.IP(content).String(), true
}

func formatAttrUint(v rlp.RawValue) (string, bool) {
	var x uint64
	if err := rlp.DecodeBytes(v, &x); err != nil {
		return "", false
	}
	return strconv.FormatUint(x, 10), true
}
