package marisa_test

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"iter"
	"os"
	"path"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/pgaskin/go-marisa"
	"github.com/pgaskin/go-marisa/internal"
)

func init() {
	flag.BoolVar(&internal.NoJIT, "marisa.nojit", false, "disable jit")
	flag.BoolVar(&internal.NoCacheQuery, "marisa.nocachequery", false, "disable query agent caching")
}

func TestMain(m *testing.M) {
	flag.Parse()

	if internal.NoJIT {
		fmt.Println("marisa: jit disabled by flag")
	}
	if internal.NoCacheQuery {
		fmt.Println("marisa: query agent caching disabled by flag")
	}

	defer os.Exit(m.Run())
}

var (
	//go:embed testdata
	testdata     embed.FS
	EnglishWords = mustLoadTestdata("words.gz", "4a53051e1939ced3e07f069a5f58e4ff2dfa9b5b")
	Go125        = mustLoadTestdata("go125.gz", "93902cc5140413de8eceee147c924dded686f4fc")
)

func mustLoadTestdata(name string, sha string) []string {
	r, err := testdata.Open(path.Join("testdata", name))
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

// mustWordsTrieData returns the serialized words trie.
var mustWordsTrieData = sync.OnceValue(func() []byte {
	var trie marisa.Trie
	if err := trie.Build(slices.Values(EnglishWords), marisa.Config{}); err != nil {
		panic(err)
	}
	buf, err := trie.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return buf
})

// mustWordsTrie returns a new copy of the words trie.
func mustWordsTrie() *marisa.Trie {
	var trie marisa.Trie
	if err := trie.UnmarshalBinary(mustWordsTrieData()); err != nil {
		panic(err)
	}
	return &trie
}

// TestHammerInstantiate ensures we don't crash during parallel instantiation.
func TestHammerInstantiate(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	t.Logf("GOMAXPROCS = %d", runtime.GOMAXPROCS(0))
	var wg sync.WaitGroup
	for range runtime.GOMAXPROCS(0) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 5000 {
				new(marisa.Trie).UnmarshalBinary([]byte{0})
			}
		}()
	}
	wg.Wait()
}

func TestReproducibility(t *testing.T) {
	// printf | marisa-build | sha1sum -
	testReproducibility(t, "Empty", "1aa6c451104c2c1b24ecb66ecb84bde2403c49b1", slices.Values([]string{}))

	// echo | marisa-build | sha1sum -
	testReproducibility(t, "Blank", "db55aeb8613305b910d42cc00b56edb53e8a3ff0", slices.Values([]string{""}))

	// printf '%s\n' {a..z}{a..z}{a..z} | marisa-build | sha1sum -
	testReproducibility(t, "Letters", "bd9586bf7f6984ea693980058de34331f4e47eae", func(yield func(string) bool) {
		for a := 'a'; a <= 'z'; a++ {
			for b := 'a'; b <= 'z'; b++ {
				for c := 'a'; c <= 'z'; c++ {
					if !yield(string(a) + string(b) + string(c)) {
						return
					}
				}
			}
		}
	})

	// gzip -cd testdata/words.gz | marisa-build | sha1sum -
	testReproducibility(t, "Words", "99604746ae19ad387a778e662a8b9014d43283e2", slices.Values(EnglishWords))

	// gzip -cd testdata/go125.gz | marisa-build | sha1sum -
	testReproducibility(t, "Go125", "e8d5188d58eabc2928dcf20cb25e374137b13674", slices.Values(Go125))
}

func testReproducibility(t *testing.T, name, sha string, seq iter.Seq[string]) {
	t.Run(name, func(t *testing.T) {
		var trie marisa.Trie
		if err := trie.Build(seq, marisa.Config{}); err != nil {
			t.Errorf("error: %v", err)
		}
		buf, err := trie.MarshalBinary()
		if err != nil {
			t.Errorf("error: %v", err)
		}
		if sum := sha1.Sum(buf); hex.EncodeToString(sum[:]) != sha {
			t.Errorf("error: does not match native marisa-build v0.3.1 output (sha1:%x)", sum)
		} else {
			t.Logf("sha1:%x", sum)
		}
	})
}

func TestString(t *testing.T) {
	t.Run("Zero", func(t *testing.T) {
		if s := new(marisa.Trie).String(); !strings.HasSuffix(s, ".Trie(uninitialized)") {
			t.Errorf("incorrect String() value %q for zero trie", s)
		} else {
			t.Logf("zero = %s", s)
		}
	})
	t.Run("Words", func(t *testing.T) {
		if s := mustWordsTrie().String(); !strings.HasSuffix(s, ".Trie(size=466550 io_size=1413352 total_size=1412654 num_tries=3 num_nodes=608368 tail_mode=text node_order=weight)") {
			t.Errorf("incorrect String() value %q for words trie", s)
		} else {
			t.Logf("words = %s", s)
		}
	})
}

func TestStats(t *testing.T) {
	t.Run("Zero", func(t *testing.T) {
		var trie marisa.Trie
		if act := trie.Size(); act != 0 {
			t.Errorf("expected size to be zero, got %d", act)
		}
		if act := trie.DiskSize(); act != 0 {
			t.Errorf("expected disk size to be zero, got %d", act)
		}
		if act := trie.TotalSize(); act != 0 {
			t.Errorf("expected total size to be zero, got %d", act)
		}
		if act := trie.NumTries(); act != 0 {
			t.Errorf("expected num tries to be zero, got %d", act)
		}
		if act := trie.NumNodes(); act != 0 {
			t.Errorf("expected num nodes to be zero, got %d", act)
		}
		if act := trie.TailMode(); act != 0 {
			t.Errorf("expected tail mode to be zero, got %d", act)
		}
		if act := trie.NodeOrder(); act != 0 {
			t.Errorf("expected node order to be zero, got %d", act)
		}
	})
	t.Run("Words", func(t *testing.T) {
		trie := mustWordsTrie()
		buf := mustWordsTrieData()
		if exp, act := uint32(466550), trie.Size(); act != exp {
			t.Errorf("expected size to be %d, got %d", exp, act)
		}
		if exp, act := uint32(len(buf)), trie.DiskSize(); act != exp {
			t.Errorf("expected disk size to be %d, got %d", exp, act)
		}
		if exp, act := uint32(1412654), trie.TotalSize(); act != exp {
			t.Errorf("expected total size to be %d, got %d", exp, act)
		}
		if exp, act := uint32(3), trie.NumTries(); act != exp {
			t.Errorf("expected num tries to be %d, got %d", exp, act)
		}
		if exp, act := uint32(608368), trie.NumNodes(); act != exp {
			t.Errorf("expected num nodes to be %d, got %d", exp, act)
		}
		if exp, act := marisa.TextTail, trie.TailMode(); act != exp {
			t.Errorf("expected tail mode to be %d, got %d", exp, act)
		}
		if exp, act := marisa.WeightOrder, trie.NodeOrder(); act != exp {
			t.Errorf("expected node order to be %d, got %d", exp, act)
		}
	})
}

func BenchmarkTrie(b *testing.B) {
	benchmarkTrie(b, "Words",
		slices.Values(EnglishWords),
		slices.Values([]string{"nonexistent---", "testing", "forethoughtfulness"}),
		slices.Values([]uint32{0, 1234}),
		slices.Values([]string{"inter", "nondeter", "un", "testing"}),
		slices.Values([]string{"forethoughtfulness", "unthinkingly"}),
	)
	benchmarkTrie(b, "Go125",
		slices.Values(Go125),
		slices.Values([]string{"nonexistent---", "go/api/go1.25.txt", "go/src/cmd/vendor/golang.org/x/tools/internal/analysisinternal/typeindex/typeindex.go"}),
		slices.Values([]uint32{0, 1234}),
		slices.Values([]string{"go/src/cmd/vendor/golang.org", "go/src/go/ast/"}),
		slices.Values([]string{"go/api/go1.25.txt", "go/src/cmd/vendor/golang.org/x/tools/internal/analysisinternal/typeindex/typeindex.go"}),
	)
}

func benchmarkTrie(b *testing.B, name string,
	keys iter.Seq[string],
	lookup iter.Seq[string],
	reverseLookup iter.Seq[uint32],
	predictiveSearch iter.Seq[string],
	commonPrefixSearch iter.Seq[string],
) {
	b.Run(name, func(b *testing.B) {
		marisa.Initialize()
		var (
			numKeys    int
			keyBytes   int64
			trieBytes  []byte
			trieConfig marisa.Config
			newTrie    func() *marisa.Trie
		)
		{
			for key := range keys {
				numKeys++
				keyBytes += int64(len(key))
			}
			var trie marisa.Trie
			if err := trie.Build(keys, trieConfig); err != nil {
				b.Fatalf("build trie: %v", err)
			}
			if buf, err := trie.MarshalBinary(); err != nil {
				b.Fatalf("marshal trie: %v", err)
			} else {
				trieBytes = buf
			}
			newTrie = func() *marisa.Trie {
				var trie marisa.Trie
				if err := trie.UnmarshalBinary(trieBytes); err != nil {
					panic(err)
				}
				return &trie
			}
		}
		b.Run("Build", func(b *testing.B) {
			b.SetBytes(keyBytes)
			b.ResetTimer()
			for range b.N {
				if err := new(marisa.Trie).Build(keys, trieConfig); err != nil {
					panic(err)
				}
			}
			b.ReportMetric(float64(numKeys), "keys/op")
			b.ReportMetric(float64(numKeys)*float64(b.N)/float64(b.Elapsed().Seconds()), "keys/s")
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(numKeys)/float64(b.N), "ns/key")
			b.ReportAllocs()
		})
		b.Run("ReadFrom", func(b *testing.B) {
			b.SetBytes(int64(len(trieBytes)))
			b.ResetTimer()
			c := &readCounter{R: bytes.NewReader(trieBytes)}
			for range b.N {
				c.R.(*bytes.Reader).Reset(trieBytes)
				if _, err := new(marisa.Trie).ReadFrom(c); err != nil {
					panic(err)
				}
			}
			b.ReportMetric(float64(c.N)/float64(b.N), "reads/op")
			b.ReportAllocs()
		})
		b.Run("WriteTo", func(b *testing.B) {
			trie := newTrie()
			b.SetBytes(int64(len(trieBytes)))
			b.ResetTimer()
			c := &writeCounter{W: io.Discard}
			for range b.N {
				if _, err := trie.WriteTo(c); err != nil {
					panic(err)
				}
			}
			b.ReportMetric(float64(c.N)/float64(b.N), "writes/op")
			b.ReportAllocs()
		})
		b.Run("UnmarshalBinary", func(b *testing.B) {
			b.SetBytes(int64(len(trieBytes)))
			b.ResetTimer()
			for range b.N {
				if err := new(marisa.Trie).UnmarshalBinary(trieBytes); err != nil {
					panic(err)
				}
			}
			b.ReportAllocs()
		})
		b.Run("MarshalBinary", func(b *testing.B) {
			trie := newTrie()
			b.SetBytes(int64(len(trieBytes)))
			b.ResetTimer()
			for range b.N {
				if _, err := trie.MarshalBinary(); err != nil {
					panic(err)
				}
			}
			b.ReportAllocs()
		})
		b.Run("DumpSeq", func(b *testing.B) {
			trie := newTrie()
			trie.Dump(0) // ensure we have a cached agent
			b.ResetTimer()
			var results int64
			for range b.N {
				var err error
				for range trie.DumpSeq()(&err) {
					results++
				}
				if err != nil {
					panic(err)
				}
			}
			b.ReportAllocs()
			b.ReportMetric(float64(results)/float64(b.Elapsed().Seconds()), "keys/s")
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(results), "ns/key")
		})
		for query := range lookup {
			b.Run("Lookup", func(b *testing.B) {
				trie := newTrie()
				trie.Lookup("") // ensure we have a cached agent
				b.ResetTimer()
				for range b.N {
					if _, _, err := trie.Lookup(query); err != nil {
						panic(err)
					}
				}
				b.ReportAllocs()
				b.ReportMetric(float64(b.N)/float64(b.Elapsed().Seconds()), "keys/s")
			})
		}
		for query := range reverseLookup {
			b.Run("ReverseLookup", func(b *testing.B) {
				trie := newTrie()
				trie.ReverseLookup(0) // ensure we have a cached agent
				b.ResetTimer()
				for range b.N {
					if _, _, err := trie.ReverseLookup(query); err != nil {
						panic(err)
					}
				}
				b.ReportAllocs()
				b.ReportMetric(float64(b.N)/float64(b.Elapsed().Seconds()), "keys/s")
			})
		}
		for query := range predictiveSearch {
			b.Run("PredictiveSearchSeq", func(b *testing.B) {
				trie := newTrie()
				trie.PredictiveSearch("", 0) // ensure we have a cached agent
				b.ResetTimer()
				var results int64
				for range b.N {
					var err error
					for range trie.PredictiveSearchSeq(query)(&err) {
						results++
					}
					if err != nil {
						panic(err)
					}
				}
				b.ReportAllocs()
				b.ReportMetric(float64(results)/float64(b.N), "keys/op")
				b.ReportMetric(float64(results)/float64(b.Elapsed().Seconds()), "keys/s")
				b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(results), "ns/key")
			})
		}
		for query := range commonPrefixSearch {
			b.Run("CommonPrefixSearchSeq", func(b *testing.B) {
				trie := newTrie()
				trie.CommonPrefixSearch("", 0) // ensure we have a cached agent
				b.ResetTimer()
				var results int64
				for range b.N {
					var err error
					for range trie.CommonPrefixSearchSeq(query)(&err) {
						results++
					}
					if err != nil {
						panic(err)
					}
				}
				b.ReportAllocs()
				b.ReportMetric(float64(results)/float64(b.N), "keys/op")
				b.ReportMetric(float64(results)/float64(b.Elapsed().Seconds()), "keys/s")
				b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(results), "ns/key")
			})
		}
		b.Run("LookupAvg", func(b *testing.B) {
			trie := newTrie()
			trie.Lookup("") // ensure we have a cached agent
			b.ResetTimer()
			for range b.N {
				for query := range keys {
					if _, _, err := trie.Lookup(query); err != nil {
						panic(err)
					}
				}
			}
			b.ReportMetric(float64(b.N*numKeys)/float64(b.Elapsed().Seconds()), "keys/s")
		})
		b.Run("ReverseLookupAvg", func(b *testing.B) {
			trie := newTrie()
			trie.ReverseLookup(0) // ensure we have a cached agent
			b.ResetTimer()
			for range b.N {
				for query := range trie.Size() {
					if _, _, err := trie.ReverseLookup(query); err != nil {
						panic(err)
					}
				}
			}
			b.ReportMetric(float64(b.N*numKeys)/float64(b.Elapsed().Seconds()), "keys/s")
		})
		b.Run("PredictiveSearchSeqAvg", func(b *testing.B) {
			trie := newTrie()
			trie.PredictiveSearch("", 0) // ensure we have a cached agent
			b.ResetTimer()
			var results int64
			for range b.N {
				for query := range keys {
					var err error
					for range trie.PredictiveSearchSeq(query)(&err) {
						results++
					}
					if err != nil {
						panic(err)
					}
				}
			}
			b.ReportMetric(float64(results)/float64(b.Elapsed().Seconds()), "result/s")
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(results), "ns/result")
			b.ReportMetric(float64(numKeys), "queries/op")
		})
		b.Run("CommonPrefixSearchSeqAvg", func(b *testing.B) {
			trie := newTrie()
			trie.CommonPrefixSearch("", 0) // ensure we have a cached agent
			b.ResetTimer()
			var results int64
			for range b.N {
				for query := range keys {
					var err error
					for range trie.CommonPrefixSearchSeq(query)(&err) {
						results++
					}
					if err != nil {
						panic(err)
					}
				}
			}
			b.ReportMetric(float64(results)/float64(b.Elapsed().Seconds()), "result/s")
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(results), "ns/result")
			b.ReportMetric(float64(numKeys), "queries/op")
		})
	})
}

type writeCounter struct {
	W io.Writer
	N uint
}

func (c *writeCounter) Write(p []byte) (n int, err error) {
	c.N++
	return c.W.Write(p)
}

type readCounter struct {
	R io.Reader
	N uint
}

func (c *readCounter) Read(p []byte) (n int, err error) {
	c.N++
	return c.R.Read(p)
}
