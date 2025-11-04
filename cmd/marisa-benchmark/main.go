// Command marisa-benchmark is a Go re-implementation of the same command.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pgaskin/go-marisa"
	"github.com/pgaskin/go-marisa/internal"
	"github.com/spf13/pflag"
)

var (
	MinNumTries  = pflag.IntP("min-num-tries", "N", 1, "limit the number of tries ["+strconv.Itoa(marisa.MinNumTries)+", "+strconv.Itoa(marisa.MaxNumTries)+"]")
	MaxNumTries  = pflag.IntP("max-num-tries", "n", 5, "limit the number of tries ["+strconv.Itoa(marisa.MinNumTries)+", "+strconv.Itoa(marisa.MaxNumTries)+"]")
	TextTail     = pflag.BoolP("text-tail", "t", false, "build a dictionary with text TAIL (default)")
	BinaryTail   = pflag.BoolP("binary-tail", "b", false, "build a dictionary with binary TAIL")
	WeightOrder  = pflag.BoolP("weight-order", "w", false, "arrange siblings in weight order (default)")
	LabelOrder   = pflag.BoolP("label-order", "l", false, "arrange siblings in label order")
	CacheLevel   = pflag.IntP("cache-level", "c", 3, "specify the cache size [1, 5]")
	PredictOn    = pflag.BoolP("predict-on", "P", false, "include predictive search (default)")
	PredictOff   = pflag.BoolP("predict-off", "p", false, "skip predictive search")
	ReuseOn      = pflag.BoolP("reuse-on", "R", false, "reuse agents (default, but not supported in this version)")
	ReuseOff     = pflag.BoolP("reuse-off", "r", false, "don't reuse agents")
	PrintSpeed   = pflag.BoolP("print-speed", "S", false, "print speed [1000 keys/s] (default)")
	PrintTime    = pflag.BoolP("print-time", "s", false, "print time [ns/key]")
	DisableJIT   = pflag.Bool("disable-jit", false, "disable wazero jit")      // not in the original version
	DisableChunk = pflag.Bool("disable-chunk", false, "disable chunked-build") // not in the original version
	Help         = pflag.BoolP("help", "h", false, "print this help")
)

func main() {
	pflag.Parse()

	if *Help {
		fmt.Printf("usage: %s [options] file...\n%s", os.Args[0], pflag.CommandLine.FlagUsages())
		os.Exit(0)
	}

	if *ReuseOn || !*ReuseOff {
		internal.NoCacheQuery = false
	} else {
		internal.NoCacheQuery = true
	}
	if *DisableJIT {
		internal.NoJIT = true
	}
	if *DisableChunk {
		internal.NoChunkBuild = true
	}

	cfg := marisa.Config{
		CacheLevel: marisa.DefaultCache,
		TailMode:   marisa.DefaultTail,
		NodeOrder:  marisa.DefaultOrder,
	}
	if *MinNumTries < marisa.MinNumTries {
		fmt.Fprintf(os.Stderr, "error: option '-n' with an invalid argument: %d\n", *MinNumTries)
		os.Exit(1)
	}
	if *MaxNumTries > marisa.MaxNumTries {
		fmt.Fprintf(os.Stderr, "error: option '-N' with an invalid argument: %d\n", *MaxNumTries)
		os.Exit(2)
	}
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
		os.Exit(3)
	}

	fmt.Printf("Number of tries: %d - %d\n", *MinNumTries, *MaxNumTries)
	fmt.Print("TAIL mode: ")
	switch cfg.TailMode {
	case marisa.TextTail:
		fmt.Println("Text mode")
	case marisa.BinaryTail:
		fmt.Println("Binary mode")
	default:
		fmt.Println()
	}
	fmt.Print("Node order: ")
	switch cfg.NodeOrder {
	case marisa.LabelOrder:
		fmt.Println("Ascending label order")
	case marisa.WeightOrder:
		fmt.Println("Descending weight order")
	default:
		fmt.Println()
	}
	fmt.Print("Cache level: ")
	switch cfg.CacheLevel {
	case marisa.HugeCache:
		fmt.Println("Huge cache")
	case marisa.LargeCache:
		fmt.Println("Large cache")
	case marisa.NormalCache:
		fmt.Println("Normal cache")
	case marisa.SmallCache:
		fmt.Println("Small cache")
	case marisa.TinyCache:
		fmt.Println("Tiny cache")
	default:
		fmt.Println()
	}

	if err := marisa.Initialize(); err != nil {
		panic(err)
	}

	var (
		keyset  []string
		weights []float32
		total   int
	)
	for key, weight := range func(yield func(string, float32) bool) {
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
					os.Exit(10)
				}
				if err := readKeys(yield, f); err == errBreak {
					f.Close()
					return
				} else if err != nil {
					f.Close()
					fmt.Fprintf(os.Stderr, "error: failed read keys from %q: %v\n", name, err)
					os.Exit(10)
				}
				f.Close()
			}
		}
	} {
		keyset = append(keyset, key)
		weights = append(weights, weight)
		total += len(key)
	}
	fmt.Printf("Number of keys: %d\n", len(keyset))
	fmt.Printf("Total length: %d\n", total)

	fmt.Printf("------+----------+--------+--------+--------+--------+--------\n")
	fmt.Printf("%6s %10s %8s %8s %8s %8s %8s\n", "#tries", "size", "build", "lookup", "reverse", "prefix", "predict")
	fmt.Printf("%6s %10s %8s %8s %8s %8s %8s\n", "", "", "", "", "lookup", "search", "search")
	if *PrintSpeed || !*PrintTime {
		fmt.Printf("%6s %10s %8s %8s %8s %8s %8s\n", "", "[bytes]", "[K/s]", "[K/s]", "[K/s]", "[K/s]", "[K/s]")
	} else {
		fmt.Printf("%6s %10s %8s %8s %8s %8s %8s\n", "", "[bytes]", "[ns]", "[ns]", "[ns]", "[ns]", "[ns]")
	}
	fmt.Printf("------+----------+--------+--------+--------+--------+--------\n")
	for numTries := *MinNumTries; numTries <= *MaxNumTries; numTries++ {
		cfg.NumTries = numTries
		fmt.Printf("%6d", numTries)
		var trie marisa.Trie
		benchmarkBuild(keyset, weights, cfg, &trie)
		if trie.Size() != 0 {
			benchmarkLookup(&trie, keyset)
			benchmarkReverseLookup(&trie, keyset)
			benchmarkCommonPrefixSearch(&trie, keyset)
			benchmarkPredictiveSearch(&trie, keyset)
		}
		fmt.Println()
	}
	fmt.Printf("------+----------+--------+--------+--------+--------+--------\n")
}

func benchmarkBuild(keyset []string, weights []float32, cfg marisa.Config, trie *marisa.Trie) {
	start := time.Now()
	if err := trie.BuildWeights(func(yield func(string, float32) bool) {
		for i, key := range keyset {
			if !yield(key, float32(weights[i])) {
				return
			}
		}
	}, cfg); err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf(" %10d", trie.DiskSize())
	printTimeInfo(len(keyset), elapsed)
}

func benchmarkLookup(trie *marisa.Trie, keyset []string) {
	start := time.Now()
	for _, key := range keyset {
		_, ok, err := trie.Lookup(key)
		if err != nil {
			panic(err)
		}
		if !ok {
			panic("lookup failed: " + strconv.Quote(key))
		}
		// note: we don't do additional validation like the original one
	}
	printTimeInfo(len(keyset), time.Since(start))
}

func benchmarkReverseLookup(trie *marisa.Trie, keyset []string) {
	start := time.Now()
	for i := range trie.Size() {
		_, ok, err := trie.ReverseLookup(i)
		if err != nil {
			panic(err)
		}
		if !ok {
			panic("reverse lookup failed")
		}
		// note: we don't do additional validation like the original one
	}
	printTimeInfo(len(keyset), time.Since(start))
}

func benchmarkCommonPrefixSearch(trie *marisa.Trie, keyset []string) {
	start := time.Now()
	for _, key := range keyset {
		var err error
		for id, key := range trie.CommonPrefixSearch(key)(&err) {
			_ = id
			_ = key
		}
		if err != nil {
			panic(err)
		}
		// note: we don't do additional validation like the original one
	}
	printTimeInfo(len(keyset), time.Since(start))
}

func benchmarkPredictiveSearch(trie *marisa.Trie, keyset []string) {
	start := time.Now()
	for _, key := range keyset {
		var err error
		for id, key := range trie.PredictiveSearch(key)(&err) {
			_ = id
			_ = key
		}
		if err != nil {
			panic(err)
		}
		// note: we don't do additional validation like the original one
	}
	printTimeInfo(len(keyset), time.Since(start))
}

func printTimeInfo(n int, elapsed time.Duration) {
	if *PrintSpeed || !*PrintTime {
		if elapsed == 0 {
			fmt.Printf(" %8s ", "-")
		} else {
			fmt.Printf(" %8.2f", float64(n)/elapsed.Seconds()/1000.0)
		}
	} else {
		if elapsed == 0 || n == 0 {
			fmt.Printf(" %8s ", "-")
		} else {
			fmt.Printf(" %8.1f", 1000000000.0*elapsed.Seconds()/float64(n))
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
