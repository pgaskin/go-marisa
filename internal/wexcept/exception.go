// Package wexcept throws and catches Go and C++ exceptions.
package wexcept

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/pgaskin/go-marisa/internal/cxxerr"
	"github.com/pgaskin/go-marisa/internal/wexport"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func Instantiate(ctx context.Context, runtime wazero.Runtime) (api.Module, error) {
	return wexport.Instantiate(ctx, runtime, "wexcept", cxxThrow)
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

// Catch returns the error from [Throw] in the returned error from a
// [api.Function] call, if any.
func Catch(err error) (error, bool) {
	if err != nil {
		var throw *thrownError
		if errors.As(err, &throw) {
			return throw.err, true
		}
	}
	return err, false
}

var cxxThrow = wexport.VIII("cxx_throw", func(ctx context.Context, mod api.Module, typ, std, what uint32) {
	type nestedThrowKey struct{}
	if ctx.Value(nestedThrowKey{}) == true {
		panic(fmt.Errorf("wexcept: post-throw callback threw"))
	}
	typStr, _ := cString(mod.Memory(), typ, 256)
	stdStr, _ := cString(mod.Memory(), std, 256)
	whatStr, _ := cString(mod.Memory(), what, 8192)
	exc := cxxerr.Wrap(typStr, stdStr, whatStr)
	if _, err := mod.ExportedFunction("wexcept_cxx_throw_destroy").Call(context.WithValue(ctx, nestedThrowKey{}, true)); err != nil {
		panic(fmt.Errorf("wexcept: failed to call post-throw callback: %w", err))
	}
	Throw(exc)
}, "cxx_throw", "typ", "std", "what")

func cString(memory api.Memory, ptr, maxLen uint32) (string, bool) {
	if ptr != 0 {
		if buf, ok := memory.Read(ptr, min(maxLen, memory.Size()-ptr)); ok {
			if i := bytes.IndexByte(buf, 0); i != -1 {
				return string(buf[:i]), true
			}
		}
	}
	return "", false
}
