package marisa_test

import (
	"iter"
	"math/bits"
	"math/rand"
	"testing"
	"time"

	"github.com/pgaskin/go-marisa"
)

// note: the intent of these tests are to test the bindings, not to test
// marisa-trie's correctness

func TestQuery(t *testing.T) {
	t.Run("Zero", func(t *testing.T) {
		var trie marisa.Trie
		if x, ok, err := trie.Lookup(""); err != nil || ok || x != 0 {
			t.Errorf("query on uninitialized trie should return nothing")
		}
		if x, ok, err := trie.ReverseLookup(0); err != nil || ok || x != "" {
			t.Errorf("query on uninitialized trie should return nothing")
		}
		if a, err := trie.Dump(-1); err != nil || a != nil {
			t.Errorf("query on uninitialized trie should return nothing")
		}
		if a, err := trie.PredictiveSearch("", -1); err != nil || a != nil {
			t.Errorf("query on uninitialized trie should return nothing")
		}
		if a, err := trie.CommonPrefixSearch("", -1); err != nil || a != nil {
			t.Errorf("query on uninitialized trie should return nothing")
		}
		if empty, err := iterErrEmpty2(trie.DumpSeq()); err != nil || !empty {
			t.Errorf("query on uninitialized trie should return nothing")
		}
		if empty, err := iterErrEmpty2(trie.PredictiveSearchSeq("")); err != nil || !empty {
			t.Errorf("query on uninitialized trie should return nothing")
		}
		if empty, err := iterErrEmpty2(trie.CommonPrefixSearchSeq("")); err != nil || !empty {
			t.Errorf("query on uninitialized trie should return nothing")
		}
	})
	t.Run("Words", func(t *testing.T) {
		trie := mustWordsTrie()
		if x, ok, err := trie.Lookup("add"); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if ok || x != 0 {
			t.Errorf("incorrect result %d", x) // yeah, this word list doesn't have the word add lol
		}

		if x, ok, err := trie.Lookup("addend"); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !ok || x != 46435 {
			t.Errorf("incorrect result %d", x)
		}
		if x, ok, err := trie.ReverseLookup(46435); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !ok || x != "addend" {
			t.Errorf("incorrect result %q", x)
		}

		if a, err := trie.Dump(-1); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if len(a) != len(words) { // assuming they're all unique
			t.Errorf("incorrect result %d", len(a))
		}
		if a, err := trie.PredictiveSearch("addend", -1); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !iterValuesEqual2(marisaKeySeq(a), "addend", "addendum", "addendums", "addenda", "addends") {
			t.Errorf("incorrect result %#v", a)
		}
		if a, err := trie.CommonPrefixSearch("addend", -1); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !iterValuesEqual2(marisaKeySeq(a), "a", "ad", "addend") {
			t.Errorf("incorrect result %#v", a)
		}

		if a, err := trie.Dump(0); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if len(a) != 0 {
			t.Errorf("incorrect result %d", len(a))
		}
		if a, err := trie.PredictiveSearch("addend", 0); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !iterValuesEqual2(marisaKeySeq(a)) {
			t.Errorf("incorrect result %#v", a)
		}
		if a, err := trie.CommonPrefixSearch("addend", 0); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !iterValuesEqual2(marisaKeySeq(a)) {
			t.Errorf("incorrect result %#v", a)
		}

		if a, err := trie.Dump(2); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if len(a) != 2 {
			t.Errorf("incorrect result %d", len(a))
		}
		if a, err := trie.PredictiveSearch("addend", 2); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !iterValuesEqual2(marisaKeySeq(a), "addend", "addendum") {
			t.Errorf("incorrect result %#v", a)
		}
		if a, err := trie.CommonPrefixSearch("addend", 2); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !iterValuesEqual2(marisaKeySeq(a), "a", "ad") {
			t.Errorf("incorrect result %#v", a)
		}

		if n, err := iterErrCount2(trie.DumpSeq()); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if n != len(words) {
			t.Errorf("incorrect result")
		}
		if ok, err := iterErrValuesEqual2(trie.PredictiveSearchSeq("addend"), "addend", "addendum", "addendums", "addenda", "addends"); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !ok {
			t.Errorf("incorrect result")
		}
		if ok, err := iterErrValuesEqual2(trie.CommonPrefixSearchSeq("addend"), "a", "ad", "addend"); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !ok {
			t.Errorf("incorrect result")
		}

		// test allocation of very long queries
		if x, ok, err := trie.Lookup(string(filled(byte('a'), 2*1024*1024))); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if ok || x != 0 {
			t.Errorf("incorrect result %d", x)
		}

		// again to test caching
		if x, ok, err := trie.Lookup(string(filled(byte('a'), 2*1024*1024))); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if ok || x != 0 {
			t.Errorf("incorrect result %d", x)
		}

		// now test a shorter one
		if x, ok, err := trie.Lookup("pneumonoultramicroscopicsilicovolcanoconiosis"); err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !ok || x != 382453 {
			t.Errorf("incorrect result %d", x)
		}
	})
}

func TestNestedQuery(t *testing.T) {
	rnd := rand.New(rand.NewSource(1234567890))
	trie := mustWordsTrie()

	const (
		attempts = 2
		levels   = 6
	)
	t.Logf("trying %d nested queries for %d keys %d times", levels, len(words), attempts)

	for try := range attempts { // do it multiple times to test the cache logic
		var (
			now      = time.Now()
			err      = make([]error, levels)
			next     = make([]func() (uint32, string, bool), levels)
			stop     = make([]func(), levels)
			count    = make([]int, levels)
			limit    = make([]int, levels)
			expected int
		)
		for i := range levels {
			// note: we've already tested the different iterators, and they all work
			// the same way, so just use dump for testing since it's easy
			seq := trie.DumpSeq()(&err[i])
			next[i], stop[i] = iter.Pull2(seq)
			if bits.UintSize < 64 {
				// 32-bit is slow, use a much smaller limit
				limit[i] = rnd.Intn(5000) + 5
				expected += limit[i]
			} else {
				if i%2 == 0 {
					// stop some early
					limit[i] = rnd.Intn(len(words)-5) + 5
					expected += limit[i]
				} else {
					limit[i] = -1
					expected += len(words)
				}
			}
		}
		for {
			i, ok := rnd.Intn(levels), false
			for j := range levels {
				k := (i + j) % levels
				if next[k] != nil {
					i, ok = k, true
					break
				}
			}
			if !ok {
				break // no unfinished iterators left
			}
			if limit[i] < 0 || count[i] < limit[i] {
				if id, key, ok := next[i](); ok {
					if x, ok, err := trie.Lookup(key); err != nil {
						t.Fatalf("try %d: lookup failed: %v", try, err)
					} else if !ok {
						t.Fatalf("try %d: got bad key %q", try, key)
					} else if x != id {
						t.Fatalf("try %d: got wrong key id %d", try, x)
					}
					if x, ok, err := trie.ReverseLookup(id); err != nil {
						t.Fatalf("try %d: lookup failed: %v", try, err)
					} else if !ok {
						t.Fatalf("try %d: got bad id %q", try, id)
					} else if x != key {
						t.Fatalf("try %d: got wrong key %q", try, x)
					}
					count[i]++
					continue
				}
				t.Logf("try %d: iterator %d done %d", try, i, count[i])
			} else {
				t.Logf("try %d: iterator %d done %d (limit reached)", try, i, count[i])
			}
			stop[i]()
			next[i] = nil
			stop[i] = nil
		}
		for _, fn := range next {
			if fn != nil {
				panic("wtf")
			}
		}
		for i, err := range err {
			if err != nil {
				t.Errorf("try %d: iterator %d finished with error: %v", try, i, err)
			}
		}
		var (
			total   int
			elapsed = time.Since(now)
		)
		for _, count := range count {
			total += count
		}
		if total != expected {
			t.Errorf("try %d: expected %d items, got %d", try, expected, count)
		} else {
			t.Logf("try %d: got expected %d items in %s (%.0f ns/op)", try, expected, elapsed, float64(elapsed.Nanoseconds())/float64(expected))
		}
	}
}

func marisaKeySeq(k []marisa.Key) iter.Seq2[uint32, string] {
	return func(yield func(uint32, string) bool) {
		for _, k := range k {
			if !yield(k.ID, k.Key) {
				return
			}
		}
	}
}

func iterErrCount2[K, V comparable](seq func(*error) iter.Seq2[K, V]) (int, error) {
	var err error
	n := iterCount2(seq(&err))
	return n, err
}

func iterCount2[K, V comparable](seq iter.Seq2[K, V]) int {
	var n int
	for range seq {
		n++
	}
	return n
}

func iterErrValuesEqual2[K, V comparable](seq func(*error) iter.Seq2[K, V], values ...V) (bool, error) {
	var err error
	ok := iterValuesEqual2(seq(&err), values...)
	return ok, err
}

func iterValuesEqual2[K, V comparable](seq iter.Seq2[K, V], values ...V) bool {
	for _, v := range seq {
		if len(values) == 0 {
			return false
		}
		if values[0] != v {
			return false
		}
		values = values[1:]
	}
	return len(values) == 0
}

func iterErrEmpty2[K, V any](seq func(*error) iter.Seq2[K, V]) (bool, error) {
	var err error
	empty := iterEmpty2(seq(&err))
	return empty, err
}

func iterEmpty2[K, V any](seq iter.Seq2[K, V]) bool {
	for range seq {
		return false
	}
	return true
}
