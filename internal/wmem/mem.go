package wmem

import (
	"bytes"
)

const (
	PageBits = 16
	PageSize = 1 << PageBits
)

// Memory is the imported memory interface for wasm2go, plus some additional
// methods.
type Memory interface {
	Data() *[]byte
	Grow(delta, max int32) int32
	Free()
}

func Pages(bytes uint32) int32 {
	return int32((int64(bytes) + PageSize - 1) >> PageBits)
}

func Bytes(m Memory, ptr, n int32) ([]byte, bool) {
	d := m.Data()
	if d == nil {
		return nil, false
	}
	b := *d
	if int(uint32(ptr)) >= len(b) {
		return nil, false
	}
	b = b[uint32(ptr):]
	if int(uint32(n)) >= len(b) {
		return b, false
	}
	return b[:uint32(n)], true
}

func CString(m Memory, ptr int32) (string, bool) {
	if ptr == 0 {
		return "", false
	}
	d := m.Data()
	if d == nil {
		return "", false
	}
	b := *d
	if int(uint32(ptr)) >= len(b) {
		return "", false
	}
	b, _, ok := bytes.Cut(b[uint32(ptr):], []byte{0})
	if !ok {
		return "", false
	}
	return string(b), true
}
