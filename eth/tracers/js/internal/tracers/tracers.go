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

// Package tracers contains the actual JavaScript tracer assets.
package tracers

import (
	"embed"
	"io/fs"
	"strings"
	"unicode"
)

//go:embed *.js
var files embed.FS

// Load reads the built-in JS tracer files embedded in the binary and
// returns a mapping of tracer name to source.
func Load() (map[string]string, error) {
	var assetTracers = make(map[string]string)
	err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		b, err := fs.ReadFile(files, path)
		if err != nil {
			return err
		}
		name := camel(strings.TrimSuffix(path, ".js"))
		assetTracers[name] = string(b)
		return nil
	})
	return assetTracers, err
}

// camel converts a snake cased input string into a camel cased output.
func camel(str string) string {
	pieces := strings.Split(str, "_")
	for i := 1; i < len(pieces); i++ {
		pieces[i] = string(unicode.ToUpper(rune(pieces[i][0]))) + pieces[i][1:]
	}
	return strings.Join(pieces, "")
}
