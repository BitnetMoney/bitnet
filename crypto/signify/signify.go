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
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	errInvalidKeyHeader = errors.New("incorrect key header")
	errInvalidKeyLength = errors.New("invalid, key length != 104")
)

func parsePrivateKey(key string) (k ed25519.PrivateKey, header []byte, keyNum []byte, err error) {
	keydata, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(keydata) != 104 {
		return nil, nil, nil, errInvalidKeyLength
	}
	if string(keydata[:2]) != "Ed" {
		return nil, nil, nil, errInvalidKeyHeader
	}
	return keydata[40:], keydata[:2], keydata[32:40], nil
}

// SignFile creates a signature of the input file.
//
// This accepts base64 keys in the format created by the 'signify' tool.
// The signature is written to the 'output' file.
func SignFile(input string, output string, key string, untrustedComment string, trustedComment string) error {
	// Pre-check comments and ensure they're set to something.
	if strings.IndexByte(untrustedComment, '\n') >= 0 {
		return errors.New("untrusted comment must not contain newline")
	}
	if strings.IndexByte(trustedComment, '\n') >= 0 {
		return errors.New("trusted comment must not contain newline")
	}
	if untrustedComment == "" {
		untrustedComment = "verify with " + input + ".pub"
	}
	if trustedComment == "" {
		trustedComment = fmt.Sprintf("timestamp:%d", time.Now().Unix())
	}

	filedata, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	skey, header, keyNum, err := parsePrivateKey(key)
	if err != nil {
		return err
	}

	// Create the main data signature.
	rawSig := ed25519.Sign(skey, filedata)
	var dataSig []byte
	dataSig = append(dataSig, header...)
	dataSig = append(dataSig, keyNum...)
	dataSig = append(dataSig, rawSig...)

	// Create the comment signature.
	var commentSigInput []byte
	commentSigInput = append(commentSigInput, rawSig...)
	commentSigInput = append(commentSigInput, []byte(trustedComment)...)
	commentSig := ed25519.Sign(skey, commentSigInput)

	// Create the output file.
	var out = new(bytes.Buffer)
	fmt.Fprintln(out, "untrusted comment:", untrustedComment)
	fmt.Fprintln(out, base64.StdEncoding.EncodeToString(dataSig))
	fmt.Fprintln(out, "trusted comment:", trustedComment)
	fmt.Fprintln(out, base64.StdEncoding.EncodeToString(commentSig))
	return os.WriteFile(output, out.Bytes(), 0644)
}
