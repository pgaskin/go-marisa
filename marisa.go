// Package marisa contains bindings for the marisa-trie library.
package marisa

import (
	"bytes"
	"cmp"
	"context"
	_ "embed"
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"iter"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/pgaskin/go-marisa/internal"
	"github.com/pgaskin/go-marisa/internal/walloc"
	"github.com/pgaskin/go-marisa/internal/wexcept"
	"github.com/pgaskin/go-marisa/internal/wexport"
	"github.com/pgaskin/go-marisa/internal/wwrap"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
)

//go:generate docker build --platform amd64 --pull --no-cache --progress plain --output wasm wasm
//go:embed wasm/marisa.wasm
var wasm []byte

// Trie is a read-only in-memory little-endian MARISA dictionary.
//
// At the moment, it must not be used concurrently. This restriction may be
// lifted in the future.
//
// On 64-bit systems, the maximum dictionary size is 4GiB. On 32-bit
// systems, the maximum dictionary size is 2 GiB. Note that if you build/load
// the same trie twice, it needs twice the amount of memory since it swaps it at
// the end.
type Trie struct {
	mod       *wwrap.Module
	qry       *query // cache the last query (we'll usually only have one at a time unless someone is nesting iterators)
	size      uint32
	ioSize    uint32
	totalSize uint32
	numTries  uint32
	numNodes  uint32
}

const scratchSpace = 32 * 1024 * 1024             // scratch space to allocate when the size is known in advance
const maxAlloc = min(math.MaxUint32, math.MaxInt) // on 32-bit platforms, limit to 2GiB, on others, limit to 4GiB

// TODO: support cloning or shared/copy-on-write or wasm threads for concurrent use?

var (
	_ fmt.Stringer               = (*Trie)(nil)
	_ encoding.BinaryMarshaler   = (*Trie)(nil)
	_ encoding.BinaryUnmarshaler = (*Trie)(nil)
	_ io.WriterTo                = (*Trie)(nil)
	_ io.ReaderFrom              = (*Trie)(nil)
)

var instance struct {
	runtime  wazero.Runtime
	compiled wazero.CompiledModule
	err      error
	once     sync.Once
}

// Initialize compiles the wasm binary.
//
// This is called implicitly when [Trie] is used for the first time.
func Initialize() error {
	instance.once.Do(initialize)
	return instance.err
}

// initialize compiles the module and instantiates the host modules.
func initialize() {
	ctx := context.Background()
	cfg := wazero.NewRuntimeConfig()
	if internal.NoJIT {
		cfg = wazero.NewRuntimeConfigInterpreter()
	}

	cfg = cfg.WithCoreFeatures(
		api.CoreFeatureMutableGlobal |
			api.CoreFeatureBulkMemoryOperations |
			api.CoreFeatureMultiValue |
			api.CoreFeatureNonTrappingFloatToIntConversion |
			api.CoreFeatureSignExtensionOps |
			api.CoreFeatureSIMD)

	instance.runtime = wazero.NewRuntimeWithConfig(ctx, cfg)

	_, instance.err = wexcept.Instantiate(ctx, instance.runtime)
	if instance.err != nil {
		return
	}

	_, instance.err = wexport.Instantiate(ctx, instance.runtime, "marisa", read, write)
	if instance.err != nil {
		return
	}

	instance.compiled, instance.err = instance.runtime.CompileModule(ctx, wasm)
}

// instantiate creates a new instance of the module.
func instantiate(alloc experimental.MemoryAllocator) (*wwrap.Module, error) {
	if err := Initialize(); err != nil {
		return nil, err
	}
	mod, err := instance.runtime.InstantiateModule(experimental.WithMemoryAllocator(context.Background(), alloc), instance.compiled, wazero.NewModuleConfig().WithName(""))
	if err != nil {
		return nil, err
	}
	return wwrap.New(mod), nil
}

// Config specifies options for a dictionary. Any unspecified options will be
// set to their default.
type Config struct {
	// NumTries specifies the number of tries to use. Usually, more tries make a
	// dictionary space-efficient but time-inefficient.
	NumTries int

	// CacheLevel specifies the cache size. A larger cache enables faster search
	// but takes a more space.
	CacheLevel CacheLevel

	// TailMode specifies the kind of TAIL implementation.
	TailMode TailMode

	// NodeOrder specifies the arrangement of nodes, which affects the time cost
	// of matching and the order of predictive search.
	NodeOrder NodeOrder
}

func (c Config) build() (uint32, bool) {
	if c.NumTries = cmp.Or(c.NumTries, DefaultNumTries); c.NumTries&^numTriesMask != 0 {
		return 0, false
	}
	if c.CacheLevel = cmp.Or(c.CacheLevel, DefaultCache); c.CacheLevel&^cacheLevelMask != 0 {
		return 0, false
	}
	if c.TailMode = cmp.Or(c.TailMode, DefaultTail); c.TailMode&^tailModeMask != 0 {
		return 0, false
	}
	if c.NodeOrder = cmp.Or(c.NodeOrder, DefaultOrder); c.NodeOrder&^nodeOrderMask != 0 {
		return 0, false
	}
	return uint32(c.NumTries) | uint32(c.CacheLevel) | uint32(c.TailMode) | uint32(c.NodeOrder), true
}

const (
	MinNumTries     = 0x00001
	MaxNumTries     = 0x0007F
	DefaultNumTries = 0x00003
)

type CacheLevel uint32

const (
	HugeCache    CacheLevel = 0x00080
	LargeCache   CacheLevel = 0x00100
	NormalCache  CacheLevel = 0x00200
	SmallCache   CacheLevel = 0x00400
	TinyCache    CacheLevel = 0x00800
	DefaultCache CacheLevel = NormalCache
)

type TailMode uint32

const (
	// TextTail merges last labels as zero-terminated strings. So, it is
	// available if and only if the last labels do not contain a NULL character.
	// If TextTail is specified and a NULL character exists in the last labels,
	// the setting is automatically switched to MARISA_BINARY_TAIL.
	TextTail TailMode = 0x01000

	// BinaryTail also merges last labels but as byte sequences. It uses a bit
	// vector to detect the end of a sequence, instead of NULL characters. So,
	// BinaryTail requires a larger space if the average length of labels is
	// greater than 8.
	BinaryTail TailMode = 0x02000

	DefaultTail TailMode = TextTail
)

type NodeOrder uint32

const (
	// LabelOrder arranges nodes in ascending label order. LabelOrder is useful
	// if an application needs to predict keys in label order.
	LabelOrder NodeOrder = 0x10000

	// WeightOrder arranges nodes in descending weight order. WeightOrder is
	// generally a better choice because it enables faster matching.
	WeightOrder NodeOrder = 0x20000

	DefaultOrder NodeOrder = WeightOrder
)

const (
	numTriesMask   = 0x0007F
	cacheLevelMask = 0x00F80
	tailModeMask   = 0x0F000
	nodeOrderMask  = 0xF0000
)

// Build builds a trie out of the specified set of keys, with a weight of 1 for
// each.
func (t *Trie) Build(keys iter.Seq[string], cfg Config) error {
	return t.BuildWeights(func(yield func(string, float32) bool) {
		for key := range keys {
			if !yield(key, 1.0) {
				return
			}
		}
	}, cfg)
}

const chunkSize = 4 * 1024 * 1024

// BuildWeights builds a trie out of the specified set of keys and weights. If a key is
// specified multiple times, the weights are accumulated.
func (t *Trie) BuildWeights(keys iter.Seq2[string, float32], cfg Config) error {
	flag, ok := cfg.build()
	if !ok {
		return errors.New("invalid config")
	}

	sa := &walloc.SliceAllocator{
		OverrideMax: maxAlloc,
	}
	mod, err := instantiate(sa)
	if err != nil {
		return err
	}

	var (
		free []uint32
		cptr uint32
		cbuf []byte
		cnum uint32
	)
	for key, weight := range keys {
		if csz := 8 + len(key); !internal.NoChunkBuild && csz < chunkSize {
			if cptr != 0 && csz > len(cbuf) {
				if _, err := mod.Call("marisa_build_push_chunk", uint64(cptr), uint64(cnum)); err != nil {
					return err
				}
				cptr, cbuf = 0, nil
			}
			if cptr == 0 {
				cptr, cbuf, err = mod.Alloc(chunkSize)
				if err != nil {
					return err
				}
				cnum = 0
				free = append(free, cptr)
			}
			var tmp []byte
			tmp, cbuf = cbuf[:csz], cbuf[csz:]
			binary.LittleEndian.PutUint32(tmp[0:4], uint32(len(key)))
			binary.LittleEndian.PutUint32(tmp[4:8], math.Float32bits(weight))
			copy(tmp[8:], key)
			cnum++
		} else {
			ptr, buf, err := mod.Alloc(len(key))
			if err != nil {
				return err
			}
			free = append(free, ptr)
			copy(buf, key)

			if _, err := mod.Call("marisa_build_push", uint64(ptr), uint64(len(key)), uint64(math.Float32bits(weight))); err != nil {
				return err
			}
		}
	}
	if !internal.NoChunkBuild && cptr != 0 {
		if _, err = mod.Call("marisa_build_push_chunk", uint64(cptr), uint64(cnum)); err != nil {
			return err
		}
		cptr, cbuf = 0, nil
	}
	if _, err := mod.Call("marisa_build", uint64(flag)); err != nil {
		return err
	}
	for _, ptr := range free {
		mod.Free(ptr)
	}
	return t.swap(mod)
}

// swap sets the trie to use the specified mod containing an initialized
// dictionary, and updates the stats.
func (t *Trie) swap(mod *wwrap.Module) error {
	res, err := mod.Call("marisa_stat")
	if err != nil {
		return err
	}
	*t = Trie{
		mod:       mod,
		qry:       nil,
		size:      uint32(res[0]),
		ioSize:    uint32(res[1]),
		totalSize: uint32(res[2]),
		numTries:  uint32(res[3]),
		numNodes:  uint32(res[4]),
	}
	return nil
}

// MapFile mmaps a file and loads the dictionary from it. On error, the trie is
// left unchanged. If not supported by the current platform, an error matching
// [errors.ErrUnsupported] is returned.
func (t *Trie) MapFile(f *os.File, offset int64, length int64) error {
	if uint64(length) > maxAlloc {
		return errors.New("dictionary too large")
	}
	va := &walloc.VirtualAllocator{
		Fallback: &walloc.SliceAllocator{
			// if it falls back to this, we'll be returning an error anyways
			OverrideMax: scratchSpace,
		},
		// on 32-bit hosts, it's critical for this (unlike the SliceAllocator)
		// since it immediately reserves virtual address space, which we only
		// have 4 GiB of
		OverrideMax: uint64(length) + scratchSpace,
	}
	mod, err := instantiate(va)
	if err != nil {
		return err
	}
	if err := va.Err(); err != nil {
		return err
	}
	ptr, err := va.MapFile(context.Background(), mod.Module(), f, offset, length, false)
	if err != nil {
		return err
	}
	if _, err := mod.Call("marisa_new", uint64(ptr), uint64(length)); err != nil {
		var ex *wexcept.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, wexcept.StdException("runtime_error")) && strings.Contains(ex.What(), "size > avail_") {
				err = io.ErrUnexpectedEOF
			}
		}
		return err
	}
	return t.swap(mod)
}

// UnmarshalBinary copies b and maps the trie directly from it. This is faster
// than [Trie.ReadFrom], but may have a less optimal memory layout. On error,
// the trie is left unchanged.
func (t *Trie) UnmarshalBinary(b []byte) error {
	if uint64(len(b)) > min(math.MaxUint32, math.MaxInt) {
		return errors.New("dictionary too large")
	}
	sa := &walloc.SliceAllocator{
		OverrideMax: uint64(len(b)) + scratchSpace,
	}
	mod, err := instantiate(sa)
	if err != nil {
		return err
	}
	ptr, buf, err := mod.Alloc(len(b))
	if err != nil {
		return err
	}
	copy(buf, b)
	if _, err := mod.Call("marisa_new", uint64(ptr), uint64(len(b))); err != nil {
		var ex *wexcept.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, wexcept.StdException("runtime_error")) && strings.Contains(ex.What(), "size > avail_") {
				err = io.ErrUnexpectedEOF
			}
		}
		return err
	}
	return t.swap(mod)
}

// ReadFrom reads a dictionary from r. On success, it will have read exactly the
// size of the dictionary. On error, the trie is left unchanged.
func (t *Trie) ReadFrom(r io.Reader) (int64, error) {
	// note: it won't actually read past in practice, since it reads exactly
	// what it wants with std::istream::read, and our stream impl is effectively
	// unbuffered
	sa := &walloc.SliceAllocator{
		OverrideMax: maxAlloc,
	}
	mod, err := instantiate(sa)
	if err != nil {
		return 0, err
	}
	c := &countReader{R: r}
	if _, err := mod.CallContext(withReader(context.Background(), c), "marisa_load"); err != nil {
		var ex *wexcept.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, wexcept.StdException("runtime_error")) && strings.Contains(ex.What(), "!stream_->read") {
				err = io.ErrUnexpectedEOF
			}
		}
		return c.N, err
	}
	return c.N, t.swap(mod)
}

// MarshalBinary serializes the dictionary.
func (t *Trie) MarshalBinary() ([]byte, error) {
	if t.mod == nil {
		return nil, errors.New("dictionary not initialized")
	}
	buf := bytes.NewBuffer(make([]byte, 0, t.ioSize))
	if _, err := t.WriteTo(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteTo serializes the dictionary to w.
func (t *Trie) WriteTo(w io.Writer) (int64, error) {
	if t.mod == nil {
		return 0, errors.New("dictionary not initialized")
	}
	c := &countWriter{W: w}
	_, err := t.mod.CallContext(withWriter(context.Background(), c), "marisa_save")
	return c.N, err
}

// String returns a human-readable description of the trie.
func (t Trie) String() string {
	var b strings.Builder
	b.WriteString(reflect.TypeOf(t).String())
	b.WriteString("(")
	if t.mod == nil {
		b.WriteString("uninitialized")
	} else {
		b.WriteString("size=")
		b.WriteString(strconv.FormatUint(uint64(t.size), 10))
		b.WriteString(" io_size=")
		b.WriteString(strconv.FormatUint(uint64(t.ioSize), 10))
		b.WriteString(" total_size=")
		b.WriteString(strconv.FormatUint(uint64(t.totalSize), 10))
		b.WriteString(" num_tries=")
		b.WriteString(strconv.FormatUint(uint64(t.numTries), 10))
		b.WriteString(" num_nodes=")
		b.WriteString(strconv.FormatUint(uint64(t.numNodes), 10))
	}
	b.WriteString(")")
	return b.String()
}

// Size returns the number of keys in the dictionary. Key are numbered from 0 to
// Size-1.
func (t Trie) Size() uint32 {
	return t.size
}

// DiskSize returns the serialized size of the dictionary.
func (t Trie) DiskSize() uint32 {
	return t.ioSize
}

// TotalSize returns the in-memory size of the dictionary.
func (t Trie) TotalSize() uint32 {
	return t.ioSize
}

// NumTries returns the number of tries in the dictionary
func (t Trie) NumTries() uint32 {
	return t.numTries
}

// NumNodes returns the number of nodes in the dictionary.
func (t Trie) NumNodes() uint32 {
	return t.numNodes
}

// query is a MARISA agent.
type query struct {
	mod      *wwrap.Module
	ptr      uint32
	shortStr uint32 // pre-allocated shortQueryLen
	shortBuf []byte
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

	var (
		str uint32
		buf []byte
	)
	if !internal.NoCacheQuery && len(s) < shortQueryLen {
		if q.shortStr == 0 {
			q.shortStr, q.shortBuf, err = t.mod.Alloc(shortQueryLen)
			if err != nil {
				t.queryDone(q)
				return nil, err
			}
		}
		str = q.shortStr
		buf = q.shortBuf
	} else {
		str, buf, err = t.mod.Alloc(len(s))
		if err != nil {
			t.queryDone(q)
			return nil, err
		}
		q.longStr = str
	}
	copy(buf, s)

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

// ReverseLookup restores a key from its ID.
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

// CommonPrefixSearch searches keys from the possible prefixes of a query
// string.
func (t *Trie) CommonPrefixSearch(query string) func(*error) iter.Seq2[uint32, string] {
	return t.search("marisa_query_common_prefix_search", query)
}

// PredictiveSearch searches keys starting with a query string.
func (t *Trie) PredictiveSearch(query string) func(*error) iter.Seq2[uint32, string] {
	return t.search("marisa_query_predictive_search", query)
}

// Dump dumps all keys.
//
// Other functions on the trie MUST NOT be called while the iterator is being
// used.
func (t *Trie) Dump() func(*error) iter.Seq2[uint32, string] {
	return t.search("marisa_query_predictive_search", "")
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

type ioKey string

func withReader(ctx context.Context, r io.Reader) context.Context {
	return context.WithValue(ctx, ioKey("read"), r)
}

func withWriter(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, ioKey("write"), w)
}

var read = wexport.VII("read", func(ctx context.Context, m api.Module, p, n uint32) {
	r, ok := ctx.Value(ioKey("read")).(io.Reader)
	if !ok {
		panic("no active reader")
	}
	if n != 0 {
		if p != 0 {
			b, ok := m.Memory().Read(p, n)
			if !ok {
				panic("invalid pointer")
			}
			if _, err := io.ReadFull(r, b); err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				wexcept.Throw(err)
			}
		} else {
			if _, err := io.CopyN(io.Discard, r, int64(n)); err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				wexcept.Throw(err)
			}
		}
	}
}, "read", "buf", "n")

var write = wexport.VII("write", func(ctx context.Context, m api.Module, p, n uint32) {
	w, ok := ctx.Value(ioKey("write")).(io.Writer)
	if !ok {
		panic("no active writer")
	}
	if n != 0 {
		if p != 0 {
			b, ok := m.Memory().Read(p, n)
			if !ok {
				panic("invalid pointer")
			}
			maybePanicAtOffset(w, n, b)
			n, err := w.Write(b)
			if err != nil {
				wexcept.Throw(err)
			}
			if n != len(b) {
				wexcept.Throw(io.ErrShortWrite)
			}
		} else {
			maybePanicAtOffset(w, n, nil)
			if _, err := io.CopyN(w, zeroReader{}, int64(n)); err != nil {
				wexcept.Throw(err)
			}
		}
	}
}, "write", "buf", "n")

type zeroReader struct{}

func (z zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

type countWriter struct {
	N int64
	W io.Writer
}

func (c *countWriter) Write(p []byte) (n int, err error) {
	n, err = c.W.Write(p)
	c.N += int64(n)
	return
}

type countReader struct {
	N int64
	R io.Reader
}

func (c *countReader) Read(p []byte) (n int, err error) {
	n, err = c.R.Read(p)
	c.N += int64(n)
	return
}

// debugPanicAtOffset causes a stack trace to be printed when a write overlaps
// the specified offsets. This is intended for debugging, or for figuring out
// what exactly a specific offset is for.
var debugPanicAtOffset uint32 = math.MaxUint32

func init() {
	if s := os.Getenv("MARISA_DEBUG_PANIC_AT_OFFSET"); s != "" {
		if n, err := strconv.ParseUint(s, 0, 32); err == nil {
			debugPanicAtOffset = uint32(n)
		}
	}
}

func maybePanicAtOffset(a any, n uint32, b []byte) {
	if debugPanicAtOffset == math.MaxUint32 {
		return
	}
	var o uint32
	switch a := a.(type) {
	case *countWriter:
		o = uint32(a.N)
	default:
		return
	}
	if o <= debugPanicAtOffset && debugPanicAtOffset-o < n {
		var s strings.Builder
		s.WriteString("write (MARISA_DEBUG_PANIC_AT_OFFSET)")
		s.WriteString("\nmatch: ")
		s.WriteString(strconv.FormatUint(uint64(o), 10))
		s.WriteString(" <= ")
		s.WriteString(strconv.FormatUint(uint64(debugPanicAtOffset), 10))
		s.WriteString(" < ")
		s.WriteString(strconv.FormatUint(uint64(o+n), 10))
		s.WriteString(" (")
		s.WriteString(strconv.FormatUint(uint64(n), 10))
		s.WriteString(")")
		s.WriteString("\ndata:")
		if b != nil {
			for _, l := range strings.Split(hex.Dump(b), "\n") {
				s.WriteString("\n\t")
				s.WriteString(l)
			}
		} else {
			s.WriteString("zeros")
		}
		s.WriteString("\n")
		panic(s.String())
	}
}
