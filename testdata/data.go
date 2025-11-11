package testdata

import (
	"compress/gzip"
	"crypto/sha1"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"strings"
)

//go:embed *.gz
var testdata embed.FS

var (
	Words = load("words", "4a53051e1939ced3e07f069a5f58e4ff2dfa9b5b")
	Go125 = load("go125", "93902cc5140413de8eceee147c924dded686f4fc")
)

func load(name string, sha string) []string {
	r, err := testdata.Open(path.Join(name + ".gz"))
	if err != nil {
		panic(err)
	}
	zr, err := gzip.NewReader(r)
	if err != nil {
		panic(err)
	}
	buf, err := io.ReadAll(zr)
	if err != nil {
		panic(err)
	}
	if err := zr.Close(); err != nil {
		panic(err)
	}
	if ss := sha1.Sum(buf); hex.EncodeToString(ss[:]) != sha {
		panic(fmt.Errorf("testdata %q changed (%x)", name, ss))
	}
	return strings.FieldsFunc(string(buf), func(r rune) bool { return r == '\n' })
}
