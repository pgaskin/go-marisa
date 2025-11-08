package marisa

import (
	"fmt"
	"iter"

	"github.com/pgaskin/go-marisa/internal"
	"github.com/pgaskin/go-marisa/internal/wwrap"
)

// query is a MARISA agent.
type query struct {
	noCopy   noCopy
	mod      *wwrap.Module
	ptr      uint32
	shortStr uint32 // pre-allocated shortQueryLen
	longStr  uint32
	res      [3]uint64
}

const shortQueryLen = 128

func (t *Trie) query() (*query, error) {
	if t.mod == nil {
		return nil, nil
	}

	var q *query
	if !internal.NoCacheQuery && t.qry != nil {
		q, t.qry = t.qry, nil
	} else {
		res, err := t.mod.Call("marisa_query_new")
		if err != nil {
			return nil, err
		}
		q = &query{
			mod: t.mod,
			ptr: uint32(res[0]),
		}
	}
	return q, nil
}

// queryString starts a new query for s. If t is not loaded, it returns nil.
func (t *Trie) queryString(s string) (*query, error) {
	if t.mod == nil {
		return nil, nil
	}

	q, err := t.query()
	if err != nil {
		return nil, err
	}

	var str uint32
	if !internal.NoCacheQuery && len(s) < shortQueryLen {
		if q.shortStr == 0 {
			q.shortStr, err = t.mod.Alloc(shortQueryLen)
			if err != nil {
				t.queryDone(q)
				return nil, err
			}
		}
		str = q.shortStr
	} else {
		str, err = t.mod.Alloc(len(s))
		if err != nil {
			t.queryDone(q)
			return nil, err
		}
		q.longStr = str
	}
	if buf, ok := t.mod.Module().Memory().Read(str, uint32(len(s))); !ok {
		panic("bad allocation")
	} else {
		copy(buf, s)
	}

	if _, err := t.mod.Call("marisa_query_set_str", uint64(q.ptr), uint64(str), uint64(len(s))); err != nil {
		t.queryDone(q)
		return nil, err
	}
	return q, nil
}

// queryID starts a new query for id. If t is not loaded, it returns nil.
func (t *Trie) queryID(id uint32) (*query, error) {
	if t.mod == nil {
		return nil, nil
	}

	q, err := t.query()
	if err != nil {
		return nil, err
	}

	if _, err := t.mod.Call("marisa_query_set_id", uint64(q.ptr), uint64(id)); err != nil {
		t.queryDone(q)
		return nil, err
	}
	return q, nil
}

func (t *Trie) queryDone(q *query) {
	if t.mod == nil || q == nil {
		return
	}
	if q.ptr == 0 {
		panic("double-free of query")
	}
	if _, err := q.mod.Call("marisa_query_clear", uint64(q.ptr)); err != nil {
		panic(fmt.Errorf("marisa: failed to free query: %w", err))
	}
	if q.longStr != 0 {
		q.mod.Free(q.longStr)
		q.longStr = 0
	}
	if !internal.NoCacheQuery && t.qry == nil {
		t.qry = q
		return
	}
	if q.shortStr != 0 {
		q.mod.Free(q.shortStr)
		q.shortStr = 0
	}
	if _, err := q.mod.Call("marisa_query_free", uint64(q.ptr)); err != nil {
		panic(fmt.Errorf("marisa: failed to free query: %w", err))
	}
	q.ptr = 0
}

// Next gets the next result for a query, returning true if a result is
// available. If q is nil, this always returns false.
func (q *query) Next(name string) (bool, error) {
	if q == nil {
		return false, nil
	}

	var ok bool
	if res, err := q.mod.Call(name, uint64(q.ptr)); err != nil {
		return false, err
	} else {
		ok = res[0] != 0
	}
	if ok {
		res, err := q.mod.Call("marisa_query_result", uint64(q.ptr))
		if err != nil {
			return false, err
		}
		q.res = [3]uint64(res)
	}
	return ok, nil
}

// ID returns the key ID. It must only be called after Next returns true.
func (q *query) ID() uint32 {
	return uint32(q.res[0])
}

// Key returns the key. It must only be called after Next returns true.
func (q *query) Key() string {
	b, ok := q.mod.Module().Memory().Read(uint32(q.res[1]), uint32(q.res[2]))
	if !ok {
		panic("bad pointer")
	}
	return string(b)
}

// Lookup checks whether a key is registered or not, returning its ID.
func (t *Trie) Lookup(key string) (uint32, bool, error) {
	q, err := t.queryString(key)
	if err != nil {
		return 0, false, err
	}
	defer t.queryDone(q)

	ok, err := q.Next("marisa_query_lookup")
	if err != nil {
		return 0, false, err
	}
	if !ok {
		return 0, false, nil
	}
	return q.ID(), true, nil
}

// ReverseLookup gets a key by its ID.
func (t *Trie) ReverseLookup(id uint32) (string, bool, error) {
	if id >= t.size {
		return "", false, nil // optimization
	}

	q, err := t.queryID(id)
	if err != nil {
		return "", false, err
	}
	defer t.queryDone(q)

	ok, err := q.Next("marisa_query_reverse_lookup")
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return q.Key(), true, nil
}

// Dump dumps all keys.
func (t *Trie) Dump(limit int) ([]Key, error) {
	return collectKeys(limit, t.DumpSeq())
}

// PredictiveSearch returns keys starting with a query string.
func (t *Trie) PredictiveSearch(query string, limit int) ([]Key, error) {
	return collectKeys(limit, t.search("marisa_query_predictive_search", query))
}

// CommonPrefixSearchSeq returns keys which equal any prefix of the query string.
func (t *Trie) CommonPrefixSearch(query string, limit int) ([]Key, error) {
	return collectKeys(limit, t.search("marisa_query_common_prefix_search", query))
}

type Key struct {
	ID  uint32
	Key string
}

func collectKeys(limit int, seq func(*error) iter.Seq2[uint32, string]) (res []Key, err error) {
	if limit > 0 {
		res = make([]Key, 0, limit)
	}
	for id, key := range seq(&err) {
		if limit == 0 {
			break
		}
		res = append(res, Key{id, key})
		if limit > 0 && len(res) >= limit {
			break
		}
	}
	return res, err
}

// DumpSeq dumps all keys.
func (t *Trie) DumpSeq() func(*error) iter.Seq2[uint32, string] {
	return t.search("marisa_query_predictive_search", "")
}

// PredictiveSearch returns keys starting with a query string.
func (t *Trie) PredictiveSearchSeq(query string) func(*error) iter.Seq2[uint32, string] {
	return t.search("marisa_query_predictive_search", query)
}

// CommonPrefixSearchSeq returns keys which equal any prefix of the query string.
func (t *Trie) CommonPrefixSearchSeq(query string) func(*error) iter.Seq2[uint32, string] {
	return t.search("marisa_query_common_prefix_search", query)
}

// search iterates over results for the specified query function.
func (t *Trie) search(name, query string) func(*error) iter.Seq2[uint32, string] {
	return func(err *error) iter.Seq2[uint32, string] {
		return func(yield func(uint32, string) bool) {
			*err = func() error {
				if t.mod == nil {
					return nil
				}

				q, err := t.queryString(query)
				if err != nil {
					return err
				}
				defer t.queryDone(q)

				for {
					ok, err := q.Next(name)
					if err != nil {
						return err
					}
					if !ok || !yield(q.ID(), q.Key()) {
						return nil
					}
				}
			}()
		}
	}
}
