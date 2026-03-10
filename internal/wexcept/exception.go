// Package wexcept throws and catches Go and C++ exceptions.
package wexcept

import (
	"github.com/pgaskin/go-marisa/internal/cxxerr"
	"github.com/pgaskin/go-marisa/internal/wmem"
)

type Imports interface {
	Xwexcept_cxx_throw_destroy()
}

type Module struct {
	Memory  wmem.Memory
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
		if x := recover(); x != nil {
			if x, ok := x.(*thrownError); ok {
				*err = x.err
				return true
			}
			panic(x)
		}
	}
	return false
}

func (m *Module) Xcxx_throw(typ int32, std int32, what int32) {
	typStr, _ := wmem.CString(m.Memory, typ)
	stdStr, _ := wmem.CString(m.Memory, std)
	whatStr, _ := wmem.CString(m.Memory, what)
	exc := cxxerr.Wrap(typStr, stdStr, whatStr)
	m.Imports.Xwexcept_cxx_throw_destroy()
	Throw(exc)
}
