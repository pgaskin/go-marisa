// Package cxxerr wraps C++ exceptions.
package cxxerr

import (
	"cmp"
	"strings"
)

// Exception represents a C++ exception.
type Exception struct {
	typ  string
	what string
	std  Std
}

// Wrap wraps a C++ exception with the specified mangled type, unmangled std
// name (empty if not a std exception, std:: prefix is required), and optional
// message. If typ isn't provided but std is, it will be set automatically.
func Wrap(typ, std, what string) error {
	exc := new(Exception)
	exc.typ = typ
	if typ := simpleDemangleClass(typ); typ != "" {
		exc.typ = typ
	}
	if std, ok := strings.CutPrefix(std, stdExceptionPrefix); ok {
		exc.std = Std(std)
	}
	if exc.typ == "" && exc.std != "" {
		exc.typ = exc.std.Error()
	}
	exc.what = what
	return exc
}

// Error creates a new C++ std exception with the specified method.
func Error(std Std, what string) error {
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
// standard library type, use [Unwrap] or [errors.As] with [Std] as the
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
		if std, ok = std.Unwrap().(Std); !ok {
			return nil
		}
	}
	return Error(std, e.what)
}

// As converts an Exception to a StdException.
func (e *Exception) As(target any) bool {
	if e.std != "" {
		if target, ok := target.(*Std); ok {
			*target = e.std
			return true
		}
	}
	return false
}

// Is allows [errors.Is] to work with [Std] values.
func (e *Exception) Is(target error) bool {
	if e.std != "" {
		if target, ok := target.(Std); ok {
			return e.std.Is(target) // we don't need to check the parents ourselves since we implement Unwrap for each one
		}
	}
	return false
}

// Std represents a standard library exception type.
//
// https://en.cppreference.com/w/cpp/error/exception.html
type Std string

const (
	stdExceptionPrefix     = "std::"
	stdExceptionBase   Std = "exception"
)

const (
	LogicError       Std = "logic_error"
	RuntimeError     Std = "runtime_error"
	BadTypeid        Std = "bad_typeid"
	BadCast          Std = "bad_cast"
	BadAlloc         Std = "bad_alloc"
	BadException     Std = "bad_exception"
	BadVariantAccess Std = "bad_variant_access"

	InvalidArgument Std = "invalid_argument"
	DomainError     Std = "domain_error"
	LengthError     Std = "length_error"
	OutOfRange      Std = "out_of_range"
	FutureError     Std = "future_error"

	RangeError           Std = "range_error"
	OverflowError        Std = "overflow_error"
	UnderflowError       Std = "underflow_error"
	RegexError           Std = "regex_error"
	SystemError          Std = "system_error"
	NonexistentLocalTime Std = "nonexistent_local_time"
	AmbiguousLocalTime   Std = "ambiguous_local_time"
	FormatError          Std = "format_error"

	IostreamFailure Std = "ios_base::failure"
	FilesystemError Std = "filesystem::filesystem_error"

	BadAnyCast Std = "bad_any_cast"

	BadArrayNewLength Std = "bad_array_new_length"
)

func (std Std) Error() string {
	if std == "" {
		return stdExceptionBase.Error()
	}
	return stdExceptionPrefix + string(std)
}

// Parents returns a space-separated string of std and all parent types.
func (std Std) Parents() string {
	var b strings.Builder
	for x, ok := std, true; ok && x != "" && x != stdExceptionBase; x, ok = x.Unwrap().(Std) {
		if b.Len() != 0 {
			b.WriteString(" ")
		}
		b.WriteString(stdExceptionPrefix)
		b.WriteString(string(x))
	}
	return b.String()
}

// Unwrap gets the parent class of std, if any.
func (std Std) Unwrap() error {
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
func (std Std) Is(target error) bool {
	if target, ok := target.(Std); ok && std != "" {
		return std == target || target == "" || target == stdExceptionBase
	}
	return false
}

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
