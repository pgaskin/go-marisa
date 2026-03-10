// Package wexcept throws and catches Go and C++ exceptions.
package wexcept

import (
	"bytes"

	"github.com/pgaskin/go-marisa/internal/cxxerr"
)

type Imports interface {
	Xwexcept_cxx_throw_destroy()
}

type Exports interface {
	Xcxx_throw(typ, std, what int32)
}

type Module struct {
	Memory interface {
		Data() *[]byte
	}
	Imports Imports
}

// thrownError represents an opaque error thrown from within a host function.
type thrownError struct {
	err error
}

func (*thrownError) Error() string {
	return "wexcept: throw"
}

// Throw throws an error. It should only be called from within a host module
// function. Note that unlike regular C exceptions, this will not unwind the C++
// stack properly, which means destructors of local variables will not be
// executed.
func Throw(err error) {
	panic(&thrownError{err})
}

// Catch should be called in a defer statement to catch a thrown error, setting
// err to it if *err is nil. It re-throws other kinds of panics. It returns true
// if a thrown error was caught.
func Catch(err *error) bool {
	if *err == nil {
		x := recover()
		if x != nil {
			if x, ok := x.(*thrownError); ok {
				*err = x
				return true
			}
		}
		panic(x)
	}
	return false
}

func (m *Module) Xcxx_throw(typ int32, std int32, what int32) {
	mem := *m.Memory.Data()
	typStr, _ := cString(mem, typ)
	stdStr, _ := cString(mem, std)
	whatStr, _ := cString(mem, what)
	exc := cxxerr.Wrap(typStr, stdStr, whatStr)
	m.Imports.Xwexcept_cxx_throw_destroy()
	Throw(exc)
}

func cString(buf []byte, ptr int32) (string, bool) {
	if ptr != 0 && int(uint32(ptr)) < len(buf) {
		if buf, _, ok := bytes.Cut(buf[uint32(ptr):], []byte{0}); ok {
			return string(buf), true
		}
	}
	return "", false
}
