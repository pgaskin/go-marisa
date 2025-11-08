package marisa_test

import (
	"iter"
	"math/rand"
	"testing"
	"time"
)

// note: the intent of these tests are to test the bindings, not to test
// marisa-trie's correctness

func TestQuery(t *testing.T) {
	trie := mustWordsTrie()
	_ = trie
	t.Skipf("TODO") // TODO
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
			if i%2 == 0 {
				// stop some early
				limit[i] = rnd.Intn(len(words)-5) + 5
				expected += limit[i]
			} else {
				limit[i] = -1
				expected += len(words)
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
