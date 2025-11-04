// Command marisa-reverse-lookup is a Go re-implementation of the same command.
//
// The options and output format are the same, but errors may differ.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/pgaskin/go-marisa"
	"github.com/spf13/pflag"
)

var (
	MmapDictionary = pflag.BoolP("mmap-dictionary", "m", false, "use memory-mapped i/o to load a dictionary (exclusive with -r)")
	ReadDictionary = pflag.BoolP("read-dictionary", "r", false, "read an entire dictionary into memory (exclusive with -m)")
	Help           = pflag.BoolP("help", "h", false, "print this help")
)

func main() {
	pflag.Parse()

	if *Help {
		fmt.Printf("usage: %s [options] file\n%s", os.Args[0], pflag.CommandLine.FlagUsages())
		os.Exit(0)
	}
	if pflag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "error: dictionary is not specified\n")
		os.Exit(10)
	}
	if pflag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "error: more than one dictionaries are specified\n")
		os.Exit(11)
	}
	name := pflag.Arg(0)

	var trie marisa.Trie
	if *MmapDictionary || !*ReadDictionary { // notee: the automatic fallback if neither is explicitly stated is not in the original version
		if err := mmap(&trie, name); err != nil {
			if *MmapDictionary || !errors.Is(err, errors.ErrUnsupported) {
				fmt.Fprintf(os.Stderr, "error: failed to mmap dictionary %q: %v\n", name, err)
				os.Exit(20)
			} else if err := load(&trie, name); err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to load dictionary %q: %v\n", name, err)
				os.Exit(21)
			}
		}
	} else {
		if err := load(&trie, name); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to load dictionary %q: %v\n", name, err)
			os.Exit(21)
		}
	}

	var failed bool
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		id, err := strconv.ParseUint(sc.Text(), 10, 32)
		if err != nil {
			break // to match the behaviour of the original version
		}
		key, ok, err := trie.ReverseLookup(uint32(id))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: reverse lookup failed: %v\n", err)
			os.Exit(30)
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "error: reverse lookup failed: id out of range\n")
			os.Exit(30)
		}
		if _, err := fmt.Printf("%d\t%s\n", id, key); err != nil {
			failed = true
		}
	}
	if failed {
		fmt.Fprintf(os.Stderr, "error: failed to write results to standard output\n")
		os.Exit(30)
	}
	_ = sc.Err() // to match the behaviour of the original version
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
