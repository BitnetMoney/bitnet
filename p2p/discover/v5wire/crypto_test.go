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

package v5wire

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

func TestVector_ECDH(t *testing.T) {
	var (
		staticKey = hexPrivkey("0xfb757dc581730490a1d7a00deea65e9b1936924caaea8f44d476014856b68736")
		publicKey = hexPubkey(crypto.S256(), "0x039961e4c2356d61bedb83052c115d311acb3a96f5777296dcf297351130266231")
		want      = hexutil.MustDecode("0x033b11a2a1f214567e1537ce5e509ffd9b21373247f2a3ff6841f4976f53165e7e")
	)
	result := ecdh(staticKey, publicKey)
	check(t, "shared-secret", result, want)
}

func TestVector_KDF(t *testing.T) {
	var (
		ephKey = hexPrivkey("0xfb757dc581730490a1d7a00deea65e9b1936924caaea8f44d476014856b68736")
		cdata  = hexutil.MustDecode("0x000000000000000000000000000000006469736376350001010102030405060708090a0b0c00180102030405060708090a0b0c0d0e0f100000000000000000")
		net    = newHandshakeTest()
	)
	defer net.close()

	destKey := &testKeyB.PublicKey
	s := deriveKeys(sha256.New, ephKey, destKey, net.nodeA.id(), net.nodeB.id(), cdata)
	t.Logf("ephemeral-key = %#x", ephKey.D)
	t.Logf("dest-pubkey = %#x", EncodePubkey(destKey))
	t.Logf("node-id-a = %#x", net.nodeA.id().Bytes())
	t.Logf("node-id-b = %#x", net.nodeB.id().Bytes())
	t.Logf("challenge-data = %#x", cdata)
	check(t, "initiator-key", s.writeKey, hexutil.MustDecode("0xdccc82d81bd610f4f76d3ebe97a40571"))
	check(t, "recipient-key", s.readKey, hexutil.MustDecode("0xac74bb8773749920b0d3a8881c173ec5"))
}

func TestVector_IDSignature(t *testing.T) {
	var (
		key    = hexPrivkey("0xfb757dc581730490a1d7a00deea65e9b1936924caaea8f44d476014856b68736")
		destID = enode.HexID("0xbbbb9d047f0488c0b5a93c1c3f2d8bafc7c8ff337024a55434a0d0555de64db9")
		ephkey = hexutil.MustDecode("0x039961e4c2356d61bedb83052c115d311acb3a96f5777296dcf297351130266231")
		cdata  = hexutil.MustDecode("0x000000000000000000000000000000006469736376350001010102030405060708090a0b0c00180102030405060708090a0b0c0d0e0f100000000000000000")
	)

	sig, err := makeIDSignature(sha256.New(), key, cdata, ephkey, destID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("static-key = %#x", key.D)
	t.Logf("challenge-data = %#x", cdata)
	t.Logf("ephemeral-pubkey = %#x", ephkey)
	t.Logf("node-id-B = %#x", destID.Bytes())
	expected := "0x94852a1e2318c4e5e9d422c98eaf19d1d90d876b29cd06ca7cb7546d0fff7b484fe86c09a064fe72bdbef73ba8e9c34df0cd2b53e9d65528c2c7f336d5dfc6e6"
	check(t, "id-signature", sig, hexutil.MustDecode(expected))
}

func TestDeriveKeys(t *testing.T) {
	t.Parallel()

	var (
		n1    = enode.ID{1}
		n2    = enode.ID{2}
		cdata = []byte{1, 2, 3, 4}
	)
	sec1 := deriveKeys(sha256.New, testKeyA, &testKeyB.PublicKey, n1, n2, cdata)
	sec2 := deriveKeys(sha256.New, testKeyB, &testKeyA.PublicKey, n1, n2, cdata)
	if sec1 == nil || sec2 == nil {
		t.Fatal("key agreement failed")
	}
	if !reflect.DeepEqual(sec1, sec2) {
		t.Fatalf("keys not equal:\n  %+v\n  %+v", sec1, sec2)
	}
}

func check(t *testing.T, what string, x, y []byte) {
	t.Helper()

	if !bytes.Equal(x, y) {
		t.Errorf("wrong %s: %#x != %#x", what, x, y)
	} else {
		t.Logf("%s = %#x", what, x)
	}
}

func hexPrivkey(input string) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic(err)
	}
	return key
}

func hexPubkey(curve elliptic.Curve, input string) *ecdsa.PublicKey {
	key, err := DecodePubkey(curve, hexutil.MustDecode(input))
	if err != nil {
		panic(err)
	}
	return key
}
