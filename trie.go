// Package marisa contains bindings for the marisa-trie library.
package marisa

import (
	_ "embed"
	"encoding"
	"fmt"
	"io"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/pgaskin/go-marisa/internal/cxxerr"
	"github.com/pgaskin/go-marisa/internal/marisa_wasm"
	"github.com/pgaskin/go-marisa/internal/wexcept"
	"github.com/pgaskin/go-marisa/internal/wmem"
)

//go:generate docker build --platform amd64 --progress plain --output . src

// Initialize was previously used to compile the wasm binary.
//
// Deprecated: This is no longer required.
func Initialize() error {
	return nil
}

// Trie is a read-only in-memory little-endian MARISA dictionary.
//
// At the moment, it must not be used concurrently. This restriction may be
// lifted in the future. It's okay to nest iterators or call methods from within
// one.
//
// On 64-bit systems, the maximum dictionary size is 4GiB. On 32-bit systems,
// the maximum dictionary size is 2 GiB. Note that if you build/load the same
// trie twice, it needs twice the amount of memory since it swaps it at the end.
type Trie struct {
	noCopy    noCopy // can't be copied since it's essentialy a handle
	mod       *module
	qry       *query // cache the last query (we'll usually only have one at a time unless someone is nesting iterators)
	size      uint32
	ioSize    uint32
	totalSize uint32
	numTries  uint32
	numNodes  uint32
	tailMode  TailMode
	nodeOrder NodeOrder
}

// binaryAppender is encoding.BinaryAppender (go1.24)
type binaryAppender interface {
	AppendBinary(b []byte) ([]byte, error)
}

var (
	_ fmt.Stringer               = (*Trie)(nil)
	_ encoding.BinaryMarshaler   = (*Trie)(nil)
	_ binaryAppender             = (*Trie)(nil)
	_ encoding.BinaryUnmarshaler = (*Trie)(nil)
	_ io.WriterTo                = (*Trie)(nil)
	_ io.ReaderFrom              = (*Trie)(nil)
)

const scratchSpace = 32 * 1024 * 1024             // scratch space to allocate when the size is known in advance
const maxAlloc = min(math.MaxUint32, math.MaxInt) // on 32-bit platforms, limit to 2GiB, on others, limit to 4GiB

type module struct {
	mem     wmem.Memory
	io      *marisaIOImpl
	wexcept *wexcept.Module
	marisa  *marisa_wasm.Module
}

// instantiate creates a new instance of the module.
func instantiate(mem wmem.Memory) (*module, error) {
	mod := &module{}
	mod.mem = mem
	mod.io = &marisaIOImpl{Memory: mod.mem}
	mod.wexcept = &wexcept.Module{Memory: mod.mem}
	mod.marisa = marisa_wasm.New(mod.mem, mod.io, mod.wexcept)
	mod.wexcept.Imports = mod.marisa
	runtime.SetFinalizer(mod, func(mod *module) {
		mod.mem.Free()
	})
	return mod, nil
}

func (m *module) Alloc(n int) (addr int32, err error) {
	if n != 0 {
		defer wexcept.Catch(&err)
		if n < 0 || int64(n) >= math.MaxInt32 {
			return 0, cxxerr.Error(cxxerr.BadAlloc, "size out of range")
		}
		addr = m.marisa.Xmalloc(int32(n))
	}
	return
}

func (m *module) Free(addr int32) {
	if addr != 0 {
		m.marisa.Xfree(addr)
	}
}

// swap sets the dictionary to use the specified mod containing an initialized
// dictionary, and updates the stats.
func (t *Trie) swap(mod *module) (err error) {
	defer wexcept.Catch(&err)
	size, ioSize, totalSize, numTries, numNodes, tailMode, nodeOrder := mod.marisa.Xmarisa_stat()
	*t = Trie{
		mod:       mod,
		qry:       nil,
		size:      uint32(size),
		ioSize:    uint32(ioSize),
		totalSize: uint32(totalSize),
		numTries:  uint32(numTries),
		numNodes:  uint32(numNodes),
		tailMode:  flagTailMode(configFlag(tailMode)),
		nodeOrder: flagNodeOrder(configFlag(nodeOrder)),
	}
	return nil
}

// String returns a human-readable description of the dictionary.
func (t *Trie) String() string {
	var b strings.Builder
	b.WriteString(reflect.TypeFor[Trie]().String())
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
		b.WriteString(" tail_mode=")
		b.WriteString(t.tailMode.String())
		b.WriteString(" node_order=")
		b.WriteString(t.nodeOrder.String())
	}
	b.WriteString(")")
	return b.String()
}

// Size returns the number of keys in the dictionary. Key are numbered from 0 to
// Size-1.
func (t *Trie) Size() uint32 {
	return t.size
}

// DiskSize returns the serialized size of the dictionary.
func (t *Trie) DiskSize() uint32 {
	return t.ioSize
}

// TotalSize returns the in-memory size of the dictionary.
func (t *Trie) TotalSize() uint32 {
	return t.totalSize
}

// NumTries returns the number of tries in the dictionary
func (t *Trie) NumTries() uint32 {
	return t.numTries
}

// NumNodes returns the number of nodes in the dictionary.
func (t *Trie) NumNodes() uint32 {
	return t.numNodes
}

// TailMode returns the tail mode of the dictionary. If unknown, it returns
// zero.
func (t *Trie) TailMode() TailMode {
	return t.tailMode
}

// NodeOrder returns the tail mode of the dictionary. If unknown, it returns
// zero.
func (t *Trie) NodeOrder() NodeOrder {
	return t.nodeOrder
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
