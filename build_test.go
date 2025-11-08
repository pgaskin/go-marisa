package marisa_test

import (
	"crypto/rand"
	"io"
	"iter"
	"maps"
	"math/bits"
	"slices"
	"testing"

	"github.com/pgaskin/go-marisa"
)

func TestBuild(t *testing.T) {
	size, keys := uint32(0), []weightKey{
		{"a/b/c", 1},
		{"b/c/d", 2},
		{"b/a/d", 1.5},
		{"b/c/e", 1.5},
		{"b/x/a", 1},
		{"b/x/b", 1},
		{"b/x/c", 2},
		{"b/x/d", 1},
		{"a/b", 1},
		{"a", 1},
		{"b", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
		{"c", 1},
	}
	{
		s := slices.Collect(noWeightKeys(keys))
		slices.Sort(s)
		s = slices.Compact(s)
		size = uint32(len(s))
	}
	t.Run("Simple", func(t *testing.T) {
		var trie marisa.Trie
		if err := trie.Build(noWeightKeys(keys), marisa.Config{}); err != nil {
			t.Fatalf("error: %v", err)
		}
		if trie.Size() != size {
			t.Errorf("incorrect size")
		}
		exp := []string{"c", "b", "b/x/a", "b/x/b", "b/x/c", "b/x/d", "b/c/d", "b/c/e", "b/a/d", "a", "a/b", "a/b/c"}
		if act := mustTrieKeys(&trie); !slices.Equal(act, exp) {
			t.Errorf("incorrect keys:\nact: %#v\nexp: %#v", act, exp)
		}
	})
	t.Run("WeightOrder", func(t *testing.T) {
		var trie marisa.Trie
		if err := trie.BuildWeights(weightKeys(keys), marisa.Config{NodeOrder: marisa.WeightOrder}); err != nil {
			t.Fatalf("error: %v", err)
		}
		if trie.Size() != size {
			t.Errorf("incorrect size")
		}
		exp := []string{"b", "b/x/c", "b/x/a", "b/x/b", "b/x/d", "b/c/d", "b/c/e", "b/a/d", "c", "a", "a/b", "a/b/c"} // note: b is first since the weight includes the weight of all children
		if act := mustTrieKeys(&trie); !slices.Equal(act, exp) {
			t.Errorf("incorrect keys:\nact: %#v\nexp: %#v", act, exp)
		}
	})
	t.Run("LabelOrder", func(t *testing.T) {
		var trie marisa.Trie
		if err := trie.BuildWeights(weightKeys(keys), marisa.Config{NodeOrder: marisa.LabelOrder}); err != nil {
			t.Fatalf("error: %v", err)
		}
		if trie.Size() != size {
			t.Errorf("incorrect size")
		}
		exp := []string{"a", "a/b", "a/b/c", "b", "b/a/d", "b/c/d", "b/c/e", "b/x/a", "b/x/b", "b/x/c", "b/x/d", "c"}
		if act := mustTrieKeys(&trie); !slices.Equal(act, exp) {
			t.Errorf("incorrect keys:\nact: %#v\nexp: %#v", act, exp)
		}
	})
	t.Run("Config", func(t *testing.T) {
		cfg := marisa.Config{
			NumTries:   1,                 // default is 3
			CacheLevel: marisa.HugeCache,  // default is normal
			TailMode:   marisa.BinaryTail, // default is text
			NodeOrder:  marisa.LabelOrder, // default is weight
		}
		var trie marisa.Trie
		if err := trie.BuildWeights(weightKeys(keys), cfg); err != nil {
			t.Fatalf("error: %v", err)
		}
		if act := trie.NumTries(); act != uint32(cfg.NumTries) {
			t.Errorf("incorrect num tries %v", act)
		}
		if act := trie.TailMode(); act != cfg.TailMode {
			t.Errorf("incorrect tail mode %v", act)
		}
		if act := trie.NodeOrder(); act != cfg.NodeOrder {
			t.Errorf("incorrect node order %v", act)
		}
	})
	t.Run("ManyKeys", func(t *testing.T) {
		if bits.UintSize < 64 && testing.Short() {
			t.Skip("slow on 32-bit")
		}
		keys := map[string]bool{}
		for l := 1; l < 128; l++ {
			n := 1000
			if bits.UintSize < 64 {
				n = 5
			}
			for range n {
				key := make([]byte, l)
				if _, err := io.ReadFull(rand.Reader, key); err != nil {
					panic("wtf")
				}
				for i, b := range key {
					key[i] = 'a' + b%26
				}
				key[len(key)-1] = '$'
				keys[string(key)] = false
			}
		}
		var trie marisa.Trie
		if err := trie.Build(maps.Keys(keys), marisa.Config{}); err != nil {
			t.Fatalf("error: %v", err)
		}
		var err error
		for id, x := range trie.DumpSeq()(&err) {
			if x[len(x)-1] != '$' {
				t.Errorf("corrupt key (id=%d value=%q)", id, x)
				continue
			}
			if _, ok := keys[x]; !ok {
				t.Errorf("unexpected key (id=%d len=%d value=%q)", id, len(x), x)
				continue
			}
			keys[x] = true
		}
		if err != nil {
			panic(err)
		}
		var missing int
		for _, ok := range keys {
			if !ok {
				missing++
			}
		}
		if missing != 0 {
			t.Errorf("missing %d/%d keys", missing, len(keys))
		}
	})
	t.Run("LargeBinary", func(t *testing.T) {
		if bits.UintSize < 64 && testing.Short() {
			t.Skip("slow on 32-bit")
		}
		keys := map[string]bool{
			"test\x00test": false, // ensure we have at least one key with a null
		}
		for _, l := range []int{5, 15, 25, 123, 456, 512, 789, 900, 1000} {
			n := 1000
			if bits.UintSize < 64 {
				n = 5
			}
			for range n {
				key := make([]byte, l)
				if _, err := io.ReadFull(rand.Reader, key); err != nil {
					panic("wtf")
				}
				keys[string(key)] = false
			}
		}
		for _, l := range []int{300 * 1024, 500 * 1024, 500 * 1024, 800 * 1024, 800 * 1024} {
			key := make([]byte, l) // ensure we have at least a very large keys
			if _, err := io.ReadFull(rand.Reader, key); err != nil {
				panic("wtf")
			}
			keys[string(key)] = false
		}
		var trie marisa.Trie
		if err := trie.Build(maps.Keys(keys), marisa.Config{}); err != nil {
			t.Fatalf("error: %v", err)
		}
		if trie.TailMode() != marisa.BinaryTail {
			t.Errorf("expected binary tail mode")
		}
		var err error
		for id, x := range trie.DumpSeq()(&err) {
			if _, ok := keys[x]; !ok {
				t.Errorf("unexpected key (id=%d len=%d)", id, len(x))
				continue
			}
			keys[x] = true
		}
		if err != nil {
			panic(err)
		}
		var missing int
		for _, ok := range keys {
			if !ok {
				missing++
			}
		}
		if missing != 0 {
			t.Errorf("missing %d/%d keys", missing, len(keys))
		}
	})
	t.Run("Words", func(t *testing.T) {
		keys := map[string]bool{}
		for _, word := range words {
			keys[word] = false
		}
		var trie marisa.Trie
		if err := trie.Build(slices.Values(words), marisa.Config{}); err != nil {
			t.Fatalf("error: %v", err)
		}
		var err error
		for id, x := range trie.DumpSeq()(&err) {
			if _, ok := keys[x]; !ok {
				t.Errorf("unexpected key (id=%d len=%d value=%q)", id, len(x), x)
				continue
			}
			keys[x] = true
		}
		if err != nil {
			panic(err)
		}
		var missing int
		for _, ok := range keys {
			if !ok {
				missing++
			}
		}
		if missing != 0 {
			t.Errorf("missing %d/%d keys", missing, len(keys))
		}
	})
}

type weightKey struct {
	Key    string
	Weight float32
}

func weightKeys(k []weightKey) iter.Seq2[string, float32] {
	return func(yield func(string, float32) bool) {
		for _, k := range k {
			if !yield(k.Key, k.Weight) {
				return
			}
		}
	}
}

func noWeightKeys(k []weightKey) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, k := range k {
			if !yield(k.Key) {
				return
			}
		}
	}
}

func mustTrieKeys(t *marisa.Trie) []string {
	ks, err := t.Dump(-1)
	if err != nil {
		panic(err)
	}
	ss := make([]string, len(ks))
	for i, k := range ks {
		ss[i] = k.Key
	}
	return ss
}
