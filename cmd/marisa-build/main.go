// Command marisa-build is a Go re-implementation of the same command.
//
// The options and output format are the same, but errors may differ.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pgaskin/go-marisa"
	"github.com/spf13/pflag"
)

var (
	NumTries    = pflag.IntP("num-tries", "n", 3, "limit the number of tries ["+strconv.Itoa(marisa.MinNumTries)+", "+strconv.Itoa(marisa.MaxNumTries)+"]")
	TextTail    = pflag.BoolP("text-tail", "t", false, "build a dictionary with text TAIL (default)")
	BinaryTail  = pflag.BoolP("binary-tail", "b", false, "build a dictionary with binary TAIL (exclusive with -t)")
	WeightOrder = pflag.BoolP("weight-order", "w", false, "arrange siblings in weight order (default)")
	LabelOrder  = pflag.BoolP("label-order", "l", false, "arrange siblings in label order")
	CacheLevel  = pflag.IntP("cache-level", "c", 3, "specify the cache size [1, 5]")
	Output      = pflag.StringP("output", "o", "", "write tries to the file (default: stdout)")
	Help        = pflag.BoolP("help", "h", false, "print this help")
)

func main() {
	pflag.Parse()

	var cfg marisa.Config

	if *Help {
		fmt.Printf("usage: %s [options] file...\n%s", os.Args[0], pflag.CommandLine.FlagUsages())
		os.Exit(0)
	}

	if *NumTries < marisa.MinNumTries || *NumTries > marisa.MaxNumTries {
		fmt.Fprintf(os.Stderr, "error: option '-n' with an invalid argument: %d\n", *NumTries)
		os.Exit(1) // yes, this is what the original one returns
	}
	cfg.NumTries = *NumTries

	if *TextTail {
		cfg.TailMode = marisa.TextTail
	}

	if *BinaryTail {
		cfg.TailMode = marisa.BinaryTail
	}

	if *WeightOrder {
		cfg.NodeOrder = marisa.WeightOrder
	}

	if *LabelOrder {
		cfg.NodeOrder = marisa.LabelOrder
	}

	switch *CacheLevel {
	case 1:
		cfg.CacheLevel = marisa.TinyCache
	case 2:
		cfg.CacheLevel = marisa.SmallCache
	case 3:
		cfg.CacheLevel = marisa.NormalCache
	case 4:
		cfg.CacheLevel = marisa.LargeCache
	case 5:
		cfg.CacheLevel = marisa.HugeCache
	default:
		fmt.Fprintf(os.Stderr, "error: option '-c' with an invalid argument: %d\n", *CacheLevel)
		os.Exit(2)
	}

	var trie marisa.Trie

	if err := trie.BuildWeights(func(yield func(string, float32) bool) {
		names := pflag.Args()
		if len(names) == 0 {
			names = []string{"-"}
		}
		for _, name := range names {
			if name == "-" { // note: support for '-' for stdin is not in the original version
				if err := readKeys(yield, os.Stdin); err == errBreak {
					return
				} else if err != nil {
					if err == errBreak {
						return
					}
					fmt.Fprintf(os.Stderr, "error: failed read keys: %v\n", err)
					os.Exit(10)
				}
			} else {
				f, err := os.Open(name)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: failed to open %q: %v\n", name, err)
					os.Exit(11)
				}
				if err := readKeys(yield, f); err == errBreak {
					f.Close()
					return
				} else if err != nil {
					f.Close()
					fmt.Fprintf(os.Stderr, "error: failed read keys from %q: %v\n", name, err)
					os.Exit(12)
				}
				f.Close()
			}
		}
	}, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to build a dictionary: %v\n", err)
		os.Exit(20)
	}

	fmt.Fprintf(os.Stderr, "#keys: %d\n", trie.Size())
	fmt.Fprintf(os.Stderr, "#nodes: %d\n", trie.NumNodes())
	fmt.Fprintf(os.Stderr, "size: %d\n", trie.DiskSize())

	if *Output != "" && *Output != "-" { // note: support for '-' for stdout is not in the original version
		if err := save(&trie, *Output); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to write a dictionary to file: %v\n", err)
			os.Exit(30)
		}
	} else {
		if _, err := trie.WriteTo(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to write a dictionary to standard output: %v\n", err)
			os.Exit(33)
		}
	}
}

var errBreak = errors.New("break")

func readKeys(yield func(string, float32) bool, r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		key := sc.Text()
		weight := float32(1.0)
		if i := strings.LastIndexByte(key, '\t'); i != -1 {
			if v, err := strconv.ParseFloat(key[i+1:], 32); err == nil {
				key = key[:i]
				weight = float32(v)
			}
		}
		if !yield(key, weight) {
			return errBreak
		}
	}
	return sc.Err()
}

func save(t *marisa.Trie, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := t.WriteTo(f); err != nil {
		return err
	}
	return f.Close()
}
