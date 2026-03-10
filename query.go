package marisa

import (
	"fmt"
	"iter"

	"github.com/pgaskin/go-marisa/internal"
	"github.com/pgaskin/go-marisa/internal/marisa_wasm"
	"github.com/pgaskin/go-marisa/internal/wexcept"
	"github.com/pgaskin/go-marisa/internal/wmem"
)

// query is a MARISA agent.
type query struct {
	noCopy   noCopy
	mod      *module
	ptr      uint32
	shortStr uint32 // pre-allocated shortQueryLen
	longStr  uint32
	res      [3]uint32
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
		ptr, err := func() (ptr uint32, err error) {
			defer wexcept.Catch(&err)
			ptr = uint32(t.mod.marisa.XQueryNew())
			return
		}()
		if err != nil {
			return nil, err
		}
		q = &query{
			mod: t.mod,
			ptr: ptr,
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
	if buf, ok := wmem.Bytes(t.mod.mem, str, uint32(len(s))); !ok {
		panic("bad allocation")
	} else {
		copy(buf, s)
	}

	if err := func() (err error) {
		defer wexcept.Catch(&err)
		t.mod.marisa.XQuerySetStr(int32(q.ptr), int32(str), int32(uint32(len(s))))
		return
	}(); err != nil {
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

	if err := func() (err error) {
		defer wexcept.Catch(&err)
		t.mod.marisa.XQuerySetID(int32(q.ptr), int32(id))
		return
	}(); err != nil {
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
	if err := func() (err error) {
		defer wexcept.Catch(&err)
		t.mod.marisa.XQueryClear(int32(q.ptr))
		return
	}(); err != nil {
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
	if err := func() (err error) {
		defer wexcept.Catch(&err)
		t.mod.marisa.XQueryFree(int32(q.ptr))
		return
	}(); err != nil {
		panic(fmt.Errorf("marisa: failed to free query: %w", err))
	}
	q.ptr = 0
}

// Next gets the next result for a query, returning true if a result is
// available. If q is nil, this always returns false.
func (q *query) Next(fn func(*marisa_wasm.Module, int32) int32) (bool, error) {
	if q == nil {
		return false, nil
	}

	var ok bool
	if res, err := func() (res int32, err error) {
		defer wexcept.Catch(&err)
		res = fn(q.mod.marisa, int32(q.ptr))
		return
	}(); err != nil {
		return false, err
	} else {
		ok = res != 0
	}
	if ok {
		res, err := func() (res [3]uint32, err error) {
			defer wexcept.Catch(&err)
			r0, r1, r2 := q.mod.marisa.XQueryResult(int32(q.ptr))
			res = [3]uint32{uint32(r0), uint32(r1), uint32(r2)}
			return
		}()
		if err != nil {
			return false, err
		}
		q.res = res
	}
	return ok, nil
}

// ID returns the key ID. It must only be called after Next returns true.
func (q *query) ID() uint32 {
	return uint32(q.res[0])
}

// Key returns the key. It must only be called after Next returns true.
func (q *query) Key() string {
	b, ok := wmem.Bytes(q.mod.mem, q.res[1], q.res[2])
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

	ok, err := q.Next((*marisa_wasm.Module).XQueryLookup)
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

	ok, err := q.Next((*marisa_wasm.Module).XQueryReverseLookup)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return q.Key(), true, nil
}

// Dump dumps all keys. If the limit is -1, all keys are returned.
func (t *Trie) Dump(limit int) ([]Key, error) {
	return collectKeys(limit, t.DumpSeq())
}

// PredictiveSearch returns keys starting with a query string. If the limit is
// -1, all keys are returned.
func (t *Trie) PredictiveSearch(query string, limit int) ([]Key, error) {
	return collectKeys(limit, t.search((*marisa_wasm.Module).XQueryPredictiveSearch, query))
}

// CommonPrefixSearchSeq returns keys which equal any prefix of the query
// string. If the limit is -1, all keys are returned.
func (t *Trie) CommonPrefixSearch(query string, limit int) ([]Key, error) {
	return collectKeys(limit, t.search((*marisa_wasm.Module).XQueryCommonPrefixSearch, query))
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
	return t.search((*marisa_wasm.Module).XQueryPredictiveSearch, "")
}

// PredictiveSearch returns keys starting with a query string.
func (t *Trie) PredictiveSearchSeq(query string) func(*error) iter.Seq2[uint32, string] {
	return t.search((*marisa_wasm.Module).XQueryPredictiveSearch, query)
}

// CommonPrefixSearchSeq returns keys which equal any prefix of the query string.
func (t *Trie) CommonPrefixSearchSeq(query string) func(*error) iter.Seq2[uint32, string] {
	return t.search((*marisa_wasm.Module).XQueryCommonPrefixSearch, query)
}

// search iterates over results for the specified query function.
func (t *Trie) search(fn func(*marisa_wasm.Module, int32) int32, query string) func(*error) iter.Seq2[uint32, string] {
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
					ok, err := q.Next(fn)
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
