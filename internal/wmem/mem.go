package wmem

import (
	"bytes"
)

const (
	PageBits = 16
	PageSize = 1 << PageBits
)

// Memory is the imported memory interface for wasm2go.
type Memory interface {
	Data() *[]byte
	Grow(delta, max int32) int32
}

// FreeableMemory is [Memory] with a method to free resources from it.
type FreeableMemory interface {
	Memory
	Free()
}

func Bytes(m Memory, ptr, n uint32) ([]byte, bool) {
	d := m.Data()
	if d == nil {
		return nil, false
	}
	b := *d
	if int(ptr) >= len(b) {
		return nil, false
	}
	b = b[ptr:]
	if int(n) >= len(b) {
		return b, false
	}
	return b[:n], true
}

func CString(m Memory, ptr uint32) (string, bool) {
	if ptr == 0 {
		return "", false
	}
	d := m.Data()
	if d == nil {
		return "", false
	}
	b := *d
	if int(ptr) >= len(b) {
		return "", false
	}
	b, _, ok := bytes.Cut(b[ptr:], []byte{0})
	if !ok {
		return "", false
	}
	return string(b), true
}
