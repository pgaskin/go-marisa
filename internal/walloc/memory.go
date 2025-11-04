package walloc

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"unsafe"

	"github.com/tetratelabs/wazero/experimental"
)

// SliceMemory allocates a movable slice-backed memory. If allocation fails,
// it panics.
func SliceMemory(cap, max uint64) experimental.LinearMemory {
	return &sliceMemory{make([]byte, 0, cap), max}
}

type sliceMemory struct {
	buf []byte
	max uint64
}

func (m *sliceMemory) String() string {
	return fmt.Sprintf("(%T)(%p)[:%d:%d/%d]", m, &m.buf[0], len(m.buf), cap(m.buf), m.max)
}

func (m *sliceMemory) Free() {}

func (m *sliceMemory) Reallocate(size uint64) []byte {
	if size > m.max {
		return nil
	}
	if cap := uint64(cap(m.buf)); size > cap {
		m.buf = slices.Grow(m.buf[:cap], int(size-cap))
	}
	m.buf = m.buf[:size]
	return m.buf
}

// VirtualMemory reserves a non-movable region of virtual memory. An error
// matching [errors.ErrUnsupported] is returned if it is not supported for the
// current platform.
func VirtualMemory(cap, max uint64) (experimental.LinearMemory, error) {
	if newVirtualMemoryImpl == nil {
		return nil, fmt.Errorf("%w: virtual memory not supported for %s/%s", errors.ErrUnsupported, runtime.GOOS, runtime.GOARCH)
	}
	m, err := newVirtualMemoryImpl(cap, max)
	if err != nil {
		return nil, err
	}
	return &virtualMemory{m}, nil
}

type virtualMemory struct {
	virtualMemoryImpl
}

type virtualMemoryImpl interface {
	experimental.LinearMemory

	// backing returns the backing slice (&=addr len=committed, cap=reserved).
	backing() []byte

	// pageSize returns the page size, which must be a non-zero power of 2.
	pageSize() int

	// mapFile maps f into memory (after which it can be closed safely), where
	// addr and offset are a multiple of pageSize, and the length is greter than
	// zero. The mapping is released when the memory is freed.
	mapFile(f *os.File, addr int, offset int64, length int, write bool) error
}

var newVirtualMemoryImpl func(cap, max uint64) (virtualMemoryImpl, error)

func (m *virtualMemory) String() string {
	b := m.backing()
	return fmt.Sprintf("(%T)(%p)[:%d:%d]", m.virtualMemoryImpl, unsafe.SliceData(b), len(b), cap(b))
}
