//go:build unix

package walloc

import (
	"errors"
	"fmt"
	"math"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

// memory implements [Memory] for unix-like platforms.
type unixVirtualMemory struct {
	buf []byte // [:committed:reserved]
}

func init() {
	newVirtualMemoryImpl = newUnixVirtualMemory
}

func newUnixVirtualMemory(cap, max uint64) (virtualMemoryImpl, error) {
	var (
		rnd  = uint64(unix.Getpagesize() - 1)
		res  = (max + rnd) &^ rnd // round up to the page size
		com  = uint64(0)
		prot = unix.PROT_NONE
	)

	// commit memory only if cap=max.
	if cap == max {
		com = res
		prot = unix.PROT_READ | unix.PROT_WRITE
	}

	// ensure res fits in an int
	if res > math.MaxInt {
		return nil, unix.EOVERFLOW
	}

	// reserve the full address space (note: protected, private, anon mappings
	// should not commit memory)
	b, err := unix.Mmap(-1, 0, int(res), prot, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		switch err {
		case unix.ENOTSUP, unix.ENOSYS, unix.ENODEV:
			err = fmt.Errorf("%w (%s)", errors.ErrUnsupported, err)
		}
		return nil, err
	}
	return &unixVirtualMemory{buf: b[:com:len(b)]}, nil
}

func (m *unixVirtualMemory) Reallocate(size uint64) []byte {
	var (
		com = uint64(len(m.buf)) // committed memory
		res = uint64(cap(m.buf)) // address space
	)

	// grow the memory if required
	if size > res {
		return nil // failure
	}
	if com < size && size <= res {
		// geometrically grow the memory, rounded up to the page size
		rnd := uint64(unix.Getpagesize() - 1)
		new := com + com>>3
		new = min(max(size, new), res)
		new = (new + rnd) &^ rnd

		// commit the memory
		err := unix.Mprotect(m.buf[com:new], unix.PROT_READ|unix.PROT_WRITE)
		if err != nil {
			return nil // failure
		}
		m.buf = m.buf[:new]
	}

	// return a slice of memory limited to the committed amount
	return m.buf[:size:len(m.buf)]
}

func (m *unixVirtualMemory) Free() {
	// note: on unix (unlike windows), it's safe to munmap pages including ones
	// mapped from a file, so we don't need to keep track of the individual file
	// mappings and free those first

	err := unix.Munmap(m.buf[:cap(m.buf)])
	if err != nil {
		panic(fmt.Errorf("walloc: failed to unmap memory: %w", err))
	}
	m.buf = nil
}

func (m *unixVirtualMemory) backing() []byte {
	return m.buf
}

func (m *unixVirtualMemory) pageSize() int {
	return unix.Getpagesize()
}

func (m *unixVirtualMemory) mapFile(f *os.File, addr int, offset int64, length int, write bool) error {
	var (
		fd   = f.Fd()
		rnd  = uint64(unix.Getpagesize() - 1)
		prot = unix.PROT_READ
	)
	if addr < 0 || uint64(addr)&rnd != 0 || length <= 0 {
		return unix.EINVAL
	}
	if offset < 0 || uint64(offset)&rnd != 0 {
		return unix.EINVAL
	}
	if length <= 0 {
		return unix.EINVAL
	}
	if len(m.buf) < addr || len(m.buf)-addr < length {
		return unix.EINVAL
	}
	if write {
		prot |= unix.PROT_WRITE
	}
	_, err := unix.MmapPtr(int(fd), offset,
		unsafe.Pointer(&m.buf[addr]), uintptr(length),
		prot, unix.MAP_SHARED|unix.MAP_FIXED)
	return err
}
