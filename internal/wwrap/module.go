// Package wwrap wraps wazero modules.
package wwrap

import (
	"context"
	"fmt"
	"math"
	"math/bits"
	"unsafe"

	"github.com/pgaskin/go-marisa/internal/wexcept"
	"github.com/tetratelabs/wazero/api"
)

// Module wraps a module with various optimizations and enhancements.
type Module struct {
	ctx context.Context
	mod api.Module

	// cache technique inspired by github.com/ncruces/go-sqlite3@v0.29.1/sqlite.go
	funcs struct {
		fn   [32]api.Function
		fd   [32]fdef
		id   [32]*byte
		mask uint32
	}
	stack [32]uint64

	allocfn api.Function
	freefn  api.Function
}

type fdef struct {
	results int
}

func makeFdef(def api.FunctionDefinition) fdef {
	return fdef{
		results: len(def.ResultTypes()),
	}
}

func New(mod api.Module) *Module {
	return &Module{
		ctx: context.Background(),
		mod: mod,
	}
}

func (m *Module) Module() api.Module {
	return m.mod
}

// getfn gets an exported function, using the cache if possible.
func (m *Module) getfn(name string) (api.Function, fdef) {
	c := &m.funcs
	p := unsafe.StringData(name)
	for i := range c.id {
		if c.id[i] == p {
			c.id[i] = nil
			c.mask &^= uint32(1) << i
			return c.fn[i], c.fd[i]
		}
	}
	fn := m.mod.ExportedFunction(name)
	fd := makeFdef(fn.Definition())
	return fn, fd
}

// putfn adds en exported function to the cache, replacing the oldest one if the
// limit has been reached.
func (m *Module) putfn(name string, fn api.Function) {
	c := &m.funcs
	p := unsafe.StringData(name)
	i := bits.TrailingZeros32(^c.mask)
	if i < 32 {
		c.mask |= uint32(1) << i
	} else {
		c.mask = uint32(1)
		i = 0
	}
	c.id[i] = p
	c.fn[i] = fn
	c.fd[i] = makeFdef(fn.Definition())
}

// Call calls an exported function.
func (m *Module) Call(name string, params ...uint64) ([]uint64, error) {
	return m.CallContext(m.ctx, name, params...)
}

// CallContext calls an exported function, caching it based on the address of
// the string literal, and returning a slice of return values valid until the
// next call.
//
// It panics if the call signature is incorrect. If it returns an error, it will
// be one from [Catch].
func (m *Module) CallContext(ctx context.Context, name string, params ...uint64) ([]uint64, error) {
	copy(m.stack[:], params) // CallWithStack will return an error if there were too few params due to params being longer than the reused stack
	fn, fd := m.getfn(name)
	err := fn.CallWithStack(ctx, m.stack[:])
	if err != nil {
		if err, ok := wexcept.Catch(err); ok {
			return nil, err
		}
		panic(err)
	}
	m.putfn(name, fn)
	return m.stack[:fd.results], nil
}

// Alloc allocates memory, returning a slice pointing into it. If n is zero, nil
// is returned successfully.
func (m *Module) Alloc(n int) (addr uint32, buf []byte, err error) {
	if m.allocfn == nil {
		m.allocfn = m.mod.ExportedFunction("malloc")
	}
	if m.allocfn == nil {
		panic("wwrap: missing memory allocator helpers")
	}
	if n != 0 {
		if n < 0 || uint64(n) >= math.MaxUint32 {
			return 0, nil, wexcept.NewException(wexcept.BadAlloc, "size out of range")
		}
		m.stack[0] = uint64(n)
		if err := m.allocfn.CallWithStack(m.ctx, m.stack[:]); err != nil {
			return 0, nil, err
		}
		addr = uint32(m.stack[0])
		if addr == 0 {
			return 0, nil, wexcept.NewException(wexcept.BadAlloc, "failed to allocate memory")
		}
		var ok bool
		buf, ok = m.mod.Memory().Read(addr, uint32(n))
		if !ok {
			panic("wwrap: bad allocation")
		}
	}
	return addr, buf, nil
}

// Free frees memory. If addr is zero, no action is taken.
func (m *Module) Free(addr uint32) {
	if m.freefn == nil {
		m.freefn = m.mod.ExportedFunction("free")
	}
	if m.freefn == nil {
		panic("wwrap: missing memory allocator helpers")
	}
	if addr != 0 {
		m.stack[0] = uint64(addr)
		if err := m.freefn.CallWithStack(m.ctx, m.stack[:]); err != nil {
			panic(fmt.Errorf("wwrap: failed to free allocation: %w", err))
		}
	}
}
