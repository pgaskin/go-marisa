package marisa_test

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	_ "embed"
	"encoding/hex"
	"io"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/pgaskin/go-marisa"
)

var (
	//go:embed testdata/words.gz
	wordsGz []byte
	words   []string
)

func init() {
	zr, err := gzip.NewReader(bytes.NewReader(wordsGz))
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
	if ss := sha1.Sum(buf); hex.EncodeToString(ss[:]) != "4a53051e1939ced3e07f069a5f58e4ff2dfa9b5b" {
		panic("word list changed")
	}
	words = strings.FieldsFunc(string(buf), func(r rune) bool { return r == '\n' })
}

// mustWordsTrie returns a new copy of the words trie.
func mustWordsTrie() *marisa.Trie {
	var trie marisa.Trie
	if err := trie.Build(slices.Values(words), marisa.Config{}); err != nil {
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
	// gzip -cd testdata/words.gz | marisa-build | sha1sum -
	const exp = "99604746ae19ad387a778e662a8b9014d43283e2"

	buf, err := mustWordsTrie().MarshalBinary()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if sha := sha1.Sum(buf); hex.EncodeToString(sha[:]) != exp {
		t.Errorf("error: does not match native marisa-build v0.3.1 output (sha1:%x)", sha)
	} else {
		t.Logf("words = sha1:%x", sha)
	}
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
		buf, err := trie.MarshalBinary()
		if err != nil {
			t.Fatalf("error: %v", err)
		}
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
