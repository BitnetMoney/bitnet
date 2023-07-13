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

// signFile reads the contents of an input file and signs it (in armored format)
// with the key provided, placing the signature into the output file.

package signify

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/jedisct1/go-minisign"
)

var (
	testSecKey = "RWRCSwAAAABVN5lr2JViGBN8DhX3/Qb/0g0wBdsNAR/APRW2qy9Fjsfr12sK2cd3URUFis1jgzQzaoayK8x4syT4G3Gvlt9RwGIwUYIQW/0mTeI+ECHu1lv5U4Wa2YHEPIesVPyRm5M="
	testPubKey = "RWTAPRW2qy9FjsBiMFGCEFv9Jk3iPhAh7tZb+VOFmtmBxDyHrFT8kZuT"
)

func TestSignify(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	data := make([]byte, 1024)
	rand.Read(data)
	tmpFile.Write(data)

	if err = tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	err = SignFile(tmpFile.Name(), tmpFile.Name()+".sig", testSecKey, "clé", "croissants")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name() + ".sig")

	// Verify the signature using a golang library
	sig, err := minisign.NewSignatureFromFile(tmpFile.Name() + ".sig")
	if err != nil {
		t.Fatal(err)
	}

	pKey, err := minisign.NewPublicKey(testPubKey)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := pKey.VerifyFromFile(tmpFile.Name(), sig)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("invalid signature")
	}
}

func TestSignifyTrustedCommentTooManyLines(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	data := make([]byte, 1024)
	rand.Read(data)
	tmpFile.Write(data)

	if err = tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	err = SignFile(tmpFile.Name(), tmpFile.Name()+".sig", testSecKey, "", "crois\nsants")
	if err == nil || err.Error() == "" {
		t.Fatalf("should have errored on a multi-line trusted comment, got %v", err)
	}
	defer os.Remove(tmpFile.Name() + ".sig")
}

func TestSignifyTrustedCommentTooManyLinesLF(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	data := make([]byte, 1024)
	rand.Read(data)
	tmpFile.Write(data)

	if err = tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	err = SignFile(tmpFile.Name(), tmpFile.Name()+".sig", testSecKey, "crois\rsants", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name() + ".sig")
}

func TestSignifyTrustedCommentEmpty(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	data := make([]byte, 1024)
	rand.Read(data)
	tmpFile.Write(data)

	if err = tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	err = SignFile(tmpFile.Name(), tmpFile.Name()+".sig", testSecKey, "", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name() + ".sig")
}
