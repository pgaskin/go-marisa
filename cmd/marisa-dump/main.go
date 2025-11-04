// Command marisa-dump is a Go re-implementation of the same command.
//
// The options and output format are the same, but errors may differ.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/pgaskin/go-marisa"
	"github.com/spf13/pflag"
)

var (
	Delimiter      = pflag.StringP("delimiter", "d", "\n", "specify the delimiter")
	MmapDictionary = pflag.BoolP("mmap-dictionary", "m", false, "use memory-mapped i/o to load a dictionary (exclusive with -r)")
	ReadDictionary = pflag.BoolP("read-dictionary", "r", false, "read an entire dictionary into memory (exclusive with -m)")
	Help           = pflag.BoolP("help", "h", false, "print this help")
)

func main() {
	pflag.Parse()

	if *Help {
		fmt.Printf("usage: %s [options] file...\n%s", os.Args[0], pflag.CommandLine.FlagUsages())
		os.Exit(0)
	}

	names := pflag.Args()
	if len(names) == 0 {
		names = []string{"-"}
	}
	for _, name := range names {
		if code := dump(name); code != 0 {
			os.Exit(code)
		}
	}
}

func dump(name string) int {
	var trie marisa.Trie
	if name != "-" { // note: support for - for stdin is not in the original version
		fmt.Fprintf(os.Stderr, "input: %s\n", name)
		if *MmapDictionary || !*ReadDictionary { // notee: the automatic fallback if neither is explicitly stated is not in the original version
			if err := mmap(&trie, name); err != nil {
				if *MmapDictionary || !errors.Is(err, errors.ErrUnsupported) {
					fmt.Fprintf(os.Stderr, "error: failed to mmap dictionary %q: %v\n", name, err)
					os.Exit(10)
				} else if err := load(&trie, name); err != nil {
					fmt.Fprintf(os.Stderr, "error: failed to load dictionary %q: %v\n", name, err)
					os.Exit(11)
				}
			}
		} else {
			if err := load(&trie, name); err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to load dictionary %q: %v\n", name, err)
				os.Exit(11)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "input: <stdin>\n")
		if _, err := trie.ReadFrom(os.Stdin); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to load dictionary %q: %v\n", name, err)
			os.Exit(22)
		}
	}
	var err error
	var keys int
	for _, key := range trie.Dump()(&err) {
		if _, err := fmt.Printf("%s%s", key, *Delimiter); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to write stdout: %v\n", err)
			os.Exit(20)
		}
		keys++
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to dump dictionary %q: %v\n", name, err)
		os.Exit(21)
	}
	fmt.Fprintf(os.Stderr, "#keys: %d\n", keys)
	return 0
}

func mmap(trie *marisa.Trie, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if err := trie.MapFile(f, 0, fi.Size()); err != nil {
		return err
	}
	return nil
}

func load(trie *marisa.Trie, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := trie.ReadFrom(f); err != nil {
		return err
	}
	return nil
}
