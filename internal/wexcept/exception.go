// Package wexcept throws and catches Go and C++ exceptions.
package wexcept

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"strings"

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

// Exception represents a C++ exception.
type Exception struct {
	typ  string
	what string
	std  StdException
}

// NewException creates a new C++ std exception.
func NewException(std StdException, what string) error {
	if strings.HasPrefix(string(std), stdExceptionPrefix) {
		panic("invalid StdException")
	}
	return &Exception{
		typ:  std.Error(),
		what: what,
		std:  cmp.Or(std, stdExceptionBase),
	}
}

// Type returns the C++ type name of the exception. To get the first known
// standard library type, use [Unwrap] or [errors.As] with [StdException] as the
// target.
func (e *Exception) Type() string {
	return e.typ
}

// What returns the error message.
func (e *Exception) What() string {
	return e.what
}

// Error returns a string describing the error message, including all parent
// error types.
func (e *Exception) Error() string {
	var b strings.Builder
	b.WriteString(cmp.Or(e.Type(), stdExceptionBase.Error()))
	if u, ok := e.Unwrap().(*Exception); ok && u.std != "" && u.std != stdExceptionBase {
		b.WriteString(" (")
		b.WriteString(u.std.Parents())
		b.WriteString(")")
	}
	if s := e.What(); s != "" {
		b.WriteString(": ")
		b.WriteString(e.What())
	}
	return b.String()
}

// Unwrap returns a copy of e with the type set to the parent stdlib class, if
// any.
func (e *Exception) Unwrap() error {
	std := e.std
	if std == "" {
		return nil
	}
	if e.typ == std.Error() {
		var ok bool
		if std, ok = std.Unwrap().(StdException); !ok {
			return nil
		}
	}
	return NewException(std, e.what)
}

// As converts an Exception to a StdException.
func (e *Exception) As(target any) bool {
	if e.std != "" {
		if target, ok := target.(*StdException); ok {
			*target = e.std
			return true
		}
	}
	return false
}

// Is allows [errors.Is] to work with [StdException] values.
func (e *Exception) Is(target error) bool {
	if e.std != "" {
		if target, ok := target.(StdException); ok {
			return e.std.Is(target) // we don't need to check the parents ourselves since we implement Unwrap for each one
		}
	}
	return false
}

// StdException represents a standard library exception type.
//
// https://en.cppreference.com/w/cpp/error/exception.html
type StdException string

const (
	stdExceptionPrefix              = "std::"
	stdExceptionBase   StdException = "exception"
)

const (
	LogicError       StdException = "logic_error"
	RuntimeError     StdException = "runtime_error"
	BadTypeid        StdException = "bad_typeid"
	BadCast          StdException = "bad_cast"
	BadAlloc         StdException = "bad_alloc"
	BadException     StdException = "bad_exception"
	BadVariantAccess StdException = "bad_variant_access"

	InvalidArgument StdException = "invalid_argument"
	DomainError     StdException = "domain_error"
	LengthError     StdException = "length_error"
	OutOfRange      StdException = "out_of_range"
	FutureError     StdException = "future_error"

	RangeError           StdException = "range_error"
	OverflowError        StdException = "overflow_error"
	UnderflowError       StdException = "underflow_error"
	RegexError           StdException = "regex_error"
	SystemError          StdException = "system_error"
	NonexistentLocalTime StdException = "nonexistent_local_time"
	AmbiguousLocalTime   StdException = "ambiguous_local_time"
	FormatError          StdException = "format_error"

	IostreamFailure StdException = "ios_base::failure"
	FilesystemError StdException = "filesystem::filesystem_error"

	BadAnyCast StdException = "bad_any_cast"

	BadArrayNewLength StdException = "bad_array_new_length"
)

func (std StdException) Error() string {
	if std == "" {
		return stdExceptionBase.Error()
	}
	return stdExceptionPrefix + string(std)
}

// Parents returns a space-separated string of std and all parent types.
func (std StdException) Parents() string {
	var b strings.Builder
	for x, ok := std, true; ok && x != "" && x != stdExceptionBase; x, ok = x.Unwrap().(StdException) {
		if b.Len() != 0 {
			b.WriteString(" ")
		}
		b.WriteString(stdExceptionPrefix)
		b.WriteString(string(x))
	}
	return b.String()
}

// Unwrap gets the parent class of std, if any.
func (std StdException) Unwrap() error {
	if std != "" && std != stdExceptionBase {
		switch std {
		case LogicError, RuntimeError, BadTypeid, BadCast, BadAlloc, BadException, BadVariantAccess:
			return stdExceptionBase
		case InvalidArgument, DomainError, LengthError, OutOfRange, FutureError:
			return LogicError
		case RangeError, OverflowError, UnderflowError, RegexError, SystemError, NonexistentLocalTime, AmbiguousLocalTime, FormatError:
			return RuntimeError
		case IostreamFailure, FilesystemError:
			return SystemError
		case BadAnyCast:
			return BadCast
		case BadArrayNewLength:
			return BadAlloc
		}
	}
	return nil
}

// Is returns true if std is the same error as target. Use [errors.Is] to check
// the entire hierachy.
func (std StdException) Is(target error) bool {
	if target, ok := target.(StdException); ok && std != "" {
		return std == target || target == "" || target == stdExceptionBase
	}
	return false
}

var cxxThrow = wexport.VIII("cxx_throw", func(ctx context.Context, mod api.Module, typ, std, what uint32) {
	type nestedThrowKey struct{}
	if ctx.Value(nestedThrowKey{}) == true {
		panic(fmt.Errorf("wexcept: post-throw callback threw"))
	}
	exc := new(Exception)
	if s, ok := cString(mod.Memory(), typ, 256); ok {
		exc.typ = cmp.Or(simpleDemangleClass(s), s)
	}
	if s, ok := cString(mod.Memory(), std, 256); ok {
		if s, ok := strings.CutPrefix(s, stdExceptionPrefix); ok {
			exc.std = StdException(strings.TrimPrefix(s, stdExceptionPrefix))
		}
	}
	if s, ok := cString(mod.Memory(), what, 8192); ok {
		exc.what = s
	}
	if exc.typ == "" && exc.std != "" {
		exc.typ = exc.std.Error()
	}
	if _, err := mod.ExportedFunction("wexcept_cxx_throw_destroy").Call(context.WithValue(ctx, nestedThrowKey{}, true)); err != nil {
		panic(fmt.Errorf("wexcept: failed to call post-throw callback: %w", err))
	}
	Throw(exc)
}, "cxx_throw", "typ", "std", "what")

// simpleDemangleClass demangles a small subset of C++ class names (for the
// Itanium C++ ABI). If invalid or unsupported, an empty string is returned.
// Notably, it does not support templates.
func simpleDemangleClass(s string) string {
	var b strings.Builder
	s, _ = strings.CutPrefix(s, "_Z")
	s, nested := strings.CutPrefix(s, "N")
	s, std := strings.CutPrefix(s, "St")
	if std {
		b.WriteString("std")
	}
	for s != "" {
		var n int
		for s != "" && !(n == 0 && s[0] == '0') && n <= len(s) && '0' <= s[0] && s[0] <= '9' {
			n *= 10
			n += int(s[0] - '0')
			s = s[1:]
		}
		if n == 0 || n > len(s) {
			return ""
		}
		if b.Len() != 0 {
			b.WriteString("::")
		}
		b.WriteString(s[:n])
		s = s[n:]
		if !nested && s != "" {
			break
		}
		if !nested || s == "E" {
			return b.String()
		}
	}
	return ""
}

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
