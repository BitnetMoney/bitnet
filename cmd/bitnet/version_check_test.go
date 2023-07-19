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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/jedisct1/go-minisign"
)

func TestVerification(t *testing.T) {
	// Signatures generated with `minisign`
	t.Run("minisig", func(t *testing.T) {
		// For this test, the pubkey is in testdata/minisign.pub
		// (the privkey is `minisign.sec`, if we want to expand this test. Password 'test' )
		pub := "RWQkliYstQBOKOdtClfgC3IypIPX6TAmoEi7beZ4gyR3wsaezvqOMWsp"
		testVerification(t, pub, "./testdata/vcheck/minisig-sigs/")
	})
	// Signatures generated with `signify-openbsd`
	t.Run("signify-openbsd", func(t *testing.T) {
		t.Skip("This currently fails, minisign expects 4 lines of data, signify provides only 2")
		// For this test, the pubkey is in testdata/signifykey.pub
		// (the privkey is `signifykey.sec`, if we want to expand this test. Password 'test' )
		pub := "RWSKLNhZb0KdATtRT7mZC/bybI3t3+Hv/O2i3ye04Dq9fnT9slpZ1a2/"
		testVerification(t, pub, "./testdata/vcheck/signify-sigs/")
	})
}

func testVerification(t *testing.T, pubkey, sigdir string) {
	// Data to verify
	data, err := os.ReadFile("./testdata/vcheck/data.json")
	if err != nil {
		t.Fatal(err)
	}
	// Signatures, with and without comments, both trusted and untrusted
	files, err := os.ReadDir(sigdir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		sig, err := os.ReadFile(filepath.Join(sigdir, f.Name()))
		if err != nil {
			t.Fatal(err)
		}
		err = verifySignature([]string{pubkey}, data, sig)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func versionUint(v string) int {
	mustInt := func(s string) int {
		a, err := strconv.Atoi(s)
		if err != nil {
			panic(v)
		}
		return a
	}
	components := strings.Split(strings.TrimPrefix(v, "v"), ".")
	a := mustInt(components[0])
	b := mustInt(components[1])
	c := mustInt(components[2])
	return a*100*100 + b*100 + c
}

// TestMatching can be used to check that the regexps are correct
func TestMatching(t *testing.T) {
	data, _ := os.ReadFile("./testdata/vcheck/vulnerabilities.json")
	var vulns []vulnJson
	if err := json.Unmarshal(data, &vulns); err != nil {
		t.Fatal(err)
	}
	check := func(version string) {
		vFull := fmt.Sprintf("Geth/%v-unstable-15339cf1-20201204/linux-amd64/go1.15.4", version)
		for _, vuln := range vulns {
			r, err := regexp.Compile(vuln.Check)
			vulnIntro := versionUint(vuln.Introduced)
			vulnFixed := versionUint(vuln.Fixed)
			current := versionUint(version)
			if err != nil {
				t.Fatal(err)
			}
			if vuln.Name == "Denial of service due to Go CVE-2020-28362" {
				// this one is not tied to geth-versions
				continue
			}
			if vulnIntro <= current && vulnFixed > current {
				// Should be vulnerable
				if !r.MatchString(vFull) {
					t.Errorf("Should be vulnerable, version %v, intro: %v, fixed: %v %v %v",
						version, vuln.Introduced, vuln.Fixed, vuln.Name, vuln.Check)
				}
			} else {
				if r.MatchString(vFull) {
					t.Errorf("Should not be flagged vulnerable, version %v, intro: %v, fixed: %v %v %d %d %d",
						version, vuln.Introduced, vuln.Fixed, vuln.Name, vulnIntro, current, vulnFixed)
				}
			}
		}
	}
	for major := 1; major < 2; major++ {
		for minor := 0; minor < 30; minor++ {
			for patch := 0; patch < 30; patch++ {
				vShort := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
				check(vShort)
			}
		}
	}
}

func TestGethPubKeysParseable(t *testing.T) {
	for _, pubkey := range gethPubKeys {
		_, err := minisign.NewPublicKey(pubkey)
		if err != nil {
			t.Errorf("Should be parseable")
		}
	}
}

func TestKeyID(t *testing.T) {
	type args struct {
		id [8]byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"@holiman key", args{id: extractKeyId(gethPubKeys[0])}, "FB1D084D39BAEC24"},
		{"second key", args{id: extractKeyId(gethPubKeys[1])}, "138B1CA303E51687"},
		{"third key", args{id: extractKeyId(gethPubKeys[2])}, "FD9813B2D2098484"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyID(tt.args.id); got != tt.want {
				t.Errorf("keyID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func extractKeyId(pubkey string) [8]byte {
	p, _ := minisign.NewPublicKey(pubkey)
	return p.KeyId
}
