package wmem

import (
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
)

// VirtualMemory reserves a non-movable region of virtual memory. An error
// matching [errors.ErrUnsupported] is returned if it is not supported for the
// current platform.
func VirtualMemory(cap, max uint64) (Memory, error) {
	if virtualMemoryImpl == nil {
		return nil, fmt.Errorf("%w: virtual memory not supported for %s/%s", errors.ErrUnsupported, runtime.GOOS, runtime.GOARCH)
	}
	return virtualMemoryImpl(cap, max)
}

var virtualMemoryImpl func(cap, max uint64) (Memory, error)

type mappableMemory interface {
	Memory
	pageSize() int
	mapFile(f *os.File, addr int, offset int64, length int, write bool) error
}

// MapFile mmaps a file, returning the offset. The module must expose the libc
// aligned_alloc and free functions. If mem does not support mmapping, an error
// matching [errors.ErrUnsupported] will be returned. The offset will be aligned
// as required.
func MapFile(mod interface {
	Xaligned_alloc(int32, int32) int32
	Xfree(int32)
}, mem Memory, f *os.File, offset, length int64, write bool) (int32, error) {
	m, ok := mem.(mappableMemory)
	if !ok {
		return 0, fmt.Errorf("%w: module memory does not support mmap", errors.ErrUnsupported)
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

	ptr := mod.Xaligned_alloc(int32(uint32(psz)), int32(uint32(msz)))
	if ptr == 0 {
		return 0, errors.New("memory allocation failed")
	}
	if _, ok := Bytes(mem, ptr, int32(msz)); !ok {
		return 0, errors.New("aligned_alloc returned bad pointer")
	}

	if err := m.mapFile(f, int(ptr), int64(fof), int(fsz), write); err != nil {
		mod.Xfree(ptr)
		return 0, err
	}
	return ptr + int32(mof), nil
}
