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

package flags

import (
	"os"
	"os/user"
	"runtime"
	"testing"
)

func TestPathExpansion(t *testing.T) {
	user, _ := user.Current()
	var tests map[string]string

	if runtime.GOOS == "windows" {
		tests = map[string]string{
			`/home/someuser/tmp`:        `\home\someuser\tmp`,
			`~/tmp`:                     user.HomeDir + `\tmp`,
			`~thisOtherUser/b/`:         `~thisOtherUser\b`,
			`$DDDXXX/a/b`:               `\tmp\a\b`,
			`/a/b/`:                     `\a\b`,
			`C:\Documents\Newsletters\`: `C:\Documents\Newsletters`,
			`C:\`:                       `C:\`,
			`\\.\pipe\\pipe\geth621383`: `\\.\pipe\\pipe\geth621383`,
		}
	} else {
		tests = map[string]string{
			`/home/someuser/tmp`:        `/home/someuser/tmp`,
			`~/tmp`:                     user.HomeDir + `/tmp`,
			`~thisOtherUser/b/`:         `~thisOtherUser/b`,
			`$DDDXXX/a/b`:               `/tmp/a/b`,
			`/a/b/`:                     `/a/b`,
			`C:\Documents\Newsletters\`: `C:\Documents\Newsletters\`,
			`C:\`:                       `C:\`,
			`\\.\pipe\\pipe\geth621383`: `\\.\pipe\\pipe\geth621383`,
		}
	}

	os.Setenv(`DDDXXX`, `/tmp`)
	for test, expected := range tests {
		got := expandPath(test)
		if got != expected {
			t.Errorf(`test %s, got %s, expected %s\n`, test, got, expected)
		}
	}
}
