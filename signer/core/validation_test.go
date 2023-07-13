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

import "testing"

func TestPasswordValidation(t *testing.T) {
	testcases := []struct {
		pw         string
		shouldFail bool
	}{
		{"test", true},
		{"testtest\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98", true},
		{"placeOfInterest⌘", true},
		{"password\nwith\nlinebreak", true},
		{"password\twith\vtabs", true},
		// Ok passwords
		{"password WhichIsOk", false},
		{"passwordOk!@#$%^&*()", false},
		{"12301203123012301230123012", false},
	}
	for _, test := range testcases {
		err := ValidatePasswordFormat(test.pw)
		if err == nil && test.shouldFail {
			t.Errorf("password '%v' should fail validation", test.pw)
		} else if err != nil && !test.shouldFail {
			t.Errorf("password '%v' shound not fail validation, but did: %v", test.pw, err)
		}
	}
}
