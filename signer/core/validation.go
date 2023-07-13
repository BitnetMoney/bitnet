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

package core

import (
	"errors"
	"regexp"
)

var printable7BitAscii = regexp.MustCompile("^[A-Za-z0-9!\"#$%&'()*+,\\-./:;<=>?@[\\]^_`{|}~ ]+$")

// ValidatePasswordFormat returns an error if the password is too short, or consists of characters
// outside the range of the printable 7bit ascii set
func ValidatePasswordFormat(password string) error {
	if len(password) < 10 {
		return errors.New("password too short (<10 characters)")
	}
	if !printable7BitAscii.MatchString(password) {
		return errors.New("password contains invalid characters - only 7bit printable ascii allowed")
	}
	return nil
}
