// Package marisa contains bindings for the marisa-trie library.
package marisa

import (
	"context"
	_ "embed"
	"encoding"
	"fmt"
	"io"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/pgaskin/go-marisa/internal"
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
	mod       *wwrap.Module
	qry       *query // cache the last query (we'll usually only have one at a time unless someone is nesting iterators)
	size      uint32
	ioSize    uint32
	totalSize uint32
	numTries  uint32
	numNodes  uint32
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

// instantiate creates a new instance of the module.
func instantiate(alloc experimental.MemoryAllocator) (*wwrap.Module, error) {
	if err := Initialize(); err != nil {
		return nil, err
	}
	mod, err := instance.runtime.InstantiateModule(experimental.WithMemoryAllocator(context.Background(), alloc), instance.compiled, wazero.NewModuleConfig().WithName(""))
	if err != nil {
		return nil, err
	}
	w := wwrap.New(mod)
	runtime.SetFinalizer(w, func(w *wwrap.Module) {
		w.Module().Close(context.Background())
	})
	return w, nil
}

// swap sets the dictionary to use the specified mod containing an initialized
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

// String returns a human-readable description of the dictionary.
func (t *Trie) String() string {
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
func (t *Trie) Size() uint32 {
	return t.size
}

// DiskSize returns the serialized size of the dictionary.
func (t *Trie) DiskSize() uint32 {
	return t.ioSize
}

// TotalSize returns the in-memory size of the dictionary.
func (t *Trie) TotalSize() uint32 {
	return t.ioSize
}

// NumTries returns the number of tries in the dictionary
func (t *Trie) NumTries() uint32 {
	return t.numTries
}

// NumNodes returns the number of nodes in the dictionary.
func (t *Trie) NumNodes() uint32 {
	return t.numNodes
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
