// Package walloc implements allocators for wazero.
package walloc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
)

var _ experimental.MemoryAllocator = (*SliceAllocator)(nil)
var _ experimental.MemoryAllocator = (*VirtualAllocator)(nil)

// SliceAllocator allocates instances of [SliceMemory].
type SliceAllocator struct {
	// OverrideMax reduces the reserved memory space to this value, or cap,
	// whichever is largest. This is mostly intended for use on 32-bit binaries,
	// which have limited virtual address space.
	OverrideMax uint64
}

func (a *SliceAllocator) Allocate(cap, max uint64) experimental.LinearMemory {
	if a.OverrideMax != 0 && a.OverrideMax >= cap {
		max = a.OverrideMax
	}
	return SliceMemory(cap, max)
}

// VirtualAllocator allocates instances of [VirtualMemory]. It must not be used
// concurrently. In general, you should make new one each time you instantiate a
// set of modules.
type VirtualAllocator struct {
	// OverrideMax reduces the reserved memory space to this value, or cap,
	// whichever is largest. This is mostly intended for use on 32-bit binaries,
	// which have limited virtual address space.
	OverrideMax uint64

	// Fallback is the fallback allocator to use. If nil, [SliceAllocator] is
	// used.
	Fallback experimental.MemoryAllocator

	vm []*virtualMemory
	ve []error
}

func (a *VirtualAllocator) String() string {
	var b strings.Builder
	b.WriteString(reflect.TypeOf(a).String())
	b.WriteString("{")
	var s bool
	for _, m := range a.vm {
		if m.backing() != nil {
			if s {
				b.WriteString(", ")
			}
			s = true
			b.WriteString(m.String())
		}
	}
	for _, m := range a.ve {
		if s {
			b.WriteString(", err:")
		}
		s = true
		b.WriteString(strconv.Quote(m.Error()))
	}
	b.WriteString("}")
	return b.String()
}

// Allocate implements [experimental.MemoryAllocator].
func (a *VirtualAllocator) Allocate(cap, max uint64) experimental.LinearMemory {
	if a.OverrideMax != 0 && a.OverrideMax >= cap {
		max = a.OverrideMax
	}
	if m, err := VirtualMemory(cap, max); err == nil {
		m := m.(*virtualMemory)
		if i := slices.IndexFunc(a.vm, func(m *virtualMemory) bool {
			return m.backing() == nil
		}); i != -1 {
			a.vm[i] = m
		} else {
			a.vm = append(a.vm, m)
		}
		return m
	} else {
		a.ve = append(a.ve, err)
	}
	if a.Fallback != nil {
		return a.Fallback.Allocate(cap, max)
	}
	return new(SliceAllocator).Allocate(cap, max)
}

// Err returns any errors encountered while allocating [VirtualMemory] instances.
func (a *VirtualAllocator) Err() error {
	return errors.Join(a.ve...)
}

// VirtualMemory gets the virtual memory backing mem, if any.
func (a *VirtualAllocator) virtualMemory(mem api.Memory) *virtualMemory {
	if len(a.vm) == 0 {
		return nil
	}
	if mem.Size() == 0 {
		if _, ok := mem.Grow(1); !ok {
			return nil
		}
	}
	p, ok := mem.Read(0, 1)
	if !ok {
		return nil
	}
	for _, m := range a.vm {
		if b := m.backing(); b != nil {
			if &b[0] == &p[0] {
				return m
			}
		}
	}
	return nil
}

// MapFile mmaps a file, returning the offset. If mod's [api.Module.Memory] is
// not a [VirtualMemory] successfully allocated by a, or libc doesn't export
// aligned_alloc, an error matching [errors.ErrUnsupported] will be returned. If
// free isn't also exported, errors will cause memory to leak. The offset will
// be aligned as required.
func (a *VirtualAllocator) MapFile(ctx context.Context, mod api.Module, f *os.File, offset, length int64, write bool) (uint32, error) {
	m := a.virtualMemory(mod.Memory())
	if m == nil {
		return 0, fmt.Errorf("%w: module memory does not support mmap", errors.ErrUnsupported)
	}

	fn := mod.ExportedFunction("aligned_alloc")
	if fn == nil {
		return 0, fmt.Errorf("%w: no libc support for aligned_alloc", errors.ErrUnsupported)
	}
	if t := fn.Definition().ParamTypes(); len(t) != 2 || t[0] != api.ValueTypeI32 || t[1] != api.ValueTypeI32 {
		return 0, fmt.Errorf("%w: wrong function signature for aligned_alloc", errors.ErrUnsupported)
	}
	if t := fn.Definition().ResultTypes(); len(t) != 1 || t[0] != api.ValueTypeI32 {
		return 0, fmt.Errorf("%w: wrong function signature for aligned_alloc", errors.ErrUnsupported)
	}

	if offset < 0 || length <= 0 {
		return 0, errors.New("invalid mapping")
	}

	psz := uint64(m.pageSize())  // page size
	rnd := uint64(psz - 1)       // page mask
	fof := uint64(offset) &^ rnd // file offset: offset rounded down to the nearest page
	mof := uint64(offset) & rnd  // virt offset: remainder of rounded offset
	fsz := mof + uint64(length)  // file size: remainder+length
	msz := (fsz + rnd) &^ rnd    // virt size: remainder+length rounded up to the nearest page

	if mof > math.MaxUint32 || msz > math.MaxUint32 {
		return 0, errors.New("offset or length out of range")
	}

	var ptr uint32
	if res, err := fn.Call(ctx, psz, msz); err != nil {
		return 0, err
	} else if res[0] == 0 {
		return 0, errors.New("memory allocation failed")
	} else {
		ptr = uint32(res[0])
	}
	if _, ok := mod.Memory().Read(ptr, uint32(msz)); !ok {
		return 0, errors.New("aligned_alloc returned bad pointer")
	}

	if err := m.mapFile(f, int(ptr), int64(fof), int(fsz), write); err != nil {
		var freeErr error
		if fn = mod.ExportedFunction("free"); fn != nil {
			_, freeErr = fn.Call(ctx, uint64(ptr))
		} else {
			freeErr = errors.ErrUnsupported
		}
		if freeErr != nil {
			err = fmt.Errorf("%w (cleanup failed: free: %v)", err, freeErr)
		}
		return 0, err
	}
	return ptr + uint32(mof), nil
}
