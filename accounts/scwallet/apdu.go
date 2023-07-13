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

package scwallet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// commandAPDU represents an application data unit sent to a smartcard.
type commandAPDU struct {
	Cla, Ins, P1, P2 uint8  // Class, Instruction, Parameter 1, Parameter 2
	Data             []byte // Command data
	Le               uint8  // Command data length
}

// serialize serializes a command APDU.
func (ca commandAPDU) serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, ca.Cla); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, ca.Ins); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, ca.P1); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, ca.P2); err != nil {
		return nil, err
	}
	if len(ca.Data) > 0 {
		if err := binary.Write(buf, binary.BigEndian, uint8(len(ca.Data))); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, ca.Data); err != nil {
			return nil, err
		}
	}
	if err := binary.Write(buf, binary.BigEndian, ca.Le); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// responseAPDU represents an application data unit received from a smart card.
type responseAPDU struct {
	Data     []byte // response data
	Sw1, Sw2 uint8  // status words 1 and 2
}

// deserialize deserializes a response APDU.
func (ra *responseAPDU) deserialize(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("can not deserialize data: payload too short (%d < 2)", len(data))
	}

	ra.Data = make([]byte, len(data)-2)

	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &ra.Data); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &ra.Sw1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &ra.Sw2); err != nil {
		return err
	}
	return nil
}
