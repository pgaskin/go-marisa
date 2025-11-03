package wautil

import (
	"errors"
	"testing"
)

func TestException(t *testing.T) {
	assertExceptionString := func(typ string, std StdException, what, str string) {
		err := &Exception{typ: typ, std: std, what: what}
		act := err.Error()
		if act != str {
			t.Errorf("expected error %#v to be %q, got %q", err, str, act)
		}
	}

	assertExceptionString("std::invalid_argument", "invalid_argument", "test",
		"std::invalid_argument (std::logic_error): test")

	assertExceptionString("std::system_error", "system_error", "test",
		"std::system_error (std::runtime_error): test")

	assertExceptionString("std::ios_base::failure", "ios_base::failure", "test",
		"std::ios_base::failure (std::system_error std::runtime_error): test")

	assertExceptionString("test::example_error", "exception", "test",
		"test::example_error: test")

	assertExceptionString("test::example_error", "invalid_argument", "test",
		"test::example_error (std::invalid_argument std::logic_error): test")

	if p := StdException(""); !errors.As(NewException("logic_error", "test"), &p) {
		t.Errorf("should be able to convert an Exception into the underlying StdException")
	} else if p != "logic_error" {
		t.Errorf("wrong converted value")
	}

	if p := StdException(""); !errors.As(NewException("", "test"), &p) {
		t.Errorf("should be able to convert a generic Exception into the underlying StdException")
	} else if p != "exception" {
		t.Errorf("wrong converted value")
	}

	if p := StdException(""); !errors.As(&Exception{"test::example_error", "test", "logic_error"}, &p) {
		t.Errorf("should be able to convert an Exception with a non-std type into the underlying StdException")
	} else if p != "logic_error" {
		t.Errorf("wrong converted value")
	}

	if p := StdException(""); !errors.As(&Exception{"test::example_error", "test", "exception"}, &p) {
		t.Errorf("should be able to convert an Exception with a non-std type into the underlying generic StdException")
	} else if p != "exception" {
		t.Errorf("wrong converted value")
	}

	if !errors.Is(&Exception{"std::system_error", "test", "system_error"}, StdException("")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"std::system_error", "test", "system_error"}, StdException("exception")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"std::system_error", "test", "system_error"}, StdException("runtime_error")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"std::system_error", "test", "system_error"}, StdException("system_error")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"test::example_error", "test", "system_error"}, StdException("")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"test::example_error", "test", "system_error"}, StdException("exception")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"test::example_error", "test", "system_error"}, StdException("runtime_error")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"test::example_error", "test", "system_error"}, StdException("system_error")) {
		t.Errorf("errors.Is should work on StdException values")
	}

	assertExceptionString("test::something_else", "", "test",
		"test::something_else: test")
	if p := StdException(""); errors.As(&Exception{"test::something_else", "test", ""}, &p) {
		t.Errorf("should not get StdException for a non-exception type")
	}
	if errors.Is(&Exception{"test::something_else", "test", ""}, StdException("")) {
		t.Errorf("should not get StdException for a non-exception type")
	}
}

func TestStdException(t *testing.T) {
	if !errors.Is(StdException(""), StdException("")) {
		t.Errorf("empty error type should match itself")
	}
	if errors.Is(StdException(""), StdException("exception")) {
		t.Errorf("empty error type should not match anything else")
	}
	if errors.Is(StdException(""), StdException("runtime_error")) {
		t.Errorf("empty error type should not match anything else")
	}
	if errors.Is(StdException(""), StdException("system_error")) {
		t.Errorf("empty error type should not match anything else")
	}
	if errors.Is(StdException(""), StdException("ios_base::failure")) {
		t.Errorf("empty error type should not match anything else")
	}
	if errors.Is(StdException(""), StdException("sdfsdf")) {
		t.Errorf("empty error type should not match anything else")
	}

	if !errors.Is(StdException("exception"), StdException("")) {
		t.Errorf("base error type should match the empty error type")
	}
	if !errors.Is(StdException("exception"), StdException("exception")) {
		t.Errorf("base error type should match the base error type")
	}
	if errors.Is(StdException("exception"), StdException("runtime_error")) {
		t.Errorf("base error type should not match anything more specific")
	}
	if errors.Is(StdException("exception"), StdException("system_error")) {
		t.Errorf("base error type should not match anything more specific")
	}
	if errors.Is(StdException("exception"), StdException("ios_base::failure")) {
		t.Errorf("base error type should not match anything more specific")
	}
	if errors.Is(StdException("exception"), StdException("sdfsdf")) {
		t.Errorf("base error type should not match an unknown error")
	}

	if !errors.Is(StdException("runtime_error"), StdException("")) {
		t.Errorf("error type should match the empty error type")
	}
	if !errors.Is(StdException("runtime_error"), StdException("exception")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(StdException("runtime_error"), StdException("runtime_error")) {
		t.Errorf("error type should match itself")
	}
	if errors.Is(StdException("runtime_error"), StdException("system_error")) {
		t.Errorf("error type should not match anything more specific")
	}
	if errors.Is(StdException("runtime_error"), StdException("ios_base::failure")) {
		t.Errorf("error type should not match anything more specific")
	}
	if errors.Is(StdException("runtime_error"), StdException("sdfsdf")) {
		t.Errorf("error type should not match an unknown error")
	}
	if errors.Is(StdException("runtime_error"), StdException("invalid_argument")) {
		t.Errorf("error type should not match another error")
	}

	if !errors.Is(StdException("system_error"), StdException("")) {
		t.Errorf("error type should match the empty error type")
	}
	if !errors.Is(StdException("system_error"), StdException("exception")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(StdException("system_error"), StdException("runtime_error")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(StdException("system_error"), StdException("system_error")) {
		t.Errorf("error type should match itself")
	}
	if errors.Is(StdException("system_error"), StdException("ios_base::failure")) {
		t.Errorf("error type should not match anything more specific")
	}
	if errors.Is(StdException("system_error"), StdException("sdfsdf")) {
		t.Errorf("error type should not match an unknown error")
	}
	if errors.Is(StdException("system_error"), StdException("invalid_argument")) {
		t.Errorf("error type should not match another error")
	}

	if !errors.Is(StdException("ios_base::failure"), StdException("")) {
		t.Errorf("error type should match the empty error type")
	}
	if !errors.Is(StdException("ios_base::failure"), StdException("exception")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(StdException("ios_base::failure"), StdException("runtime_error")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(StdException("ios_base::failure"), StdException("system_error")) {
		t.Errorf("error type should match itself")
	}
	if !errors.Is(StdException("ios_base::failure"), StdException("ios_base::failure")) {
		t.Errorf("error type should match itself")
	}
	if errors.Is(StdException("ios_base::failure"), StdException("sdfsdf")) {
		t.Errorf("error type should not match an unknown error")
	}
	if errors.Is(StdException("ios_base::failure"), StdException("invalid_argument")) {
		t.Errorf("error type should not match another error")
	}
}

func TestSimpleDemangleClass(t *testing.T) {
	for _, tc := range [][2]string{
		{"", ""},
		{"1", ""},
		{"_Z", ""},
		{"_ZSt9exception", "std::exception"},
		{"_ZSt9exceptio", ""},
		{"_ZNSt3__18ios_base7failureE", "std::__1::ios_base::failure"},
		{"_ZNSt3__18ios_base7failure", ""},
		{"_ZNSt3__18ios_base7failureEa", ""},
		{"N3tmp3tmp22asdasdasdsdfgsdfsdfsdf11CustomErrorE", "tmp::tmp::asdasdasdsdfgsdfsdfsdf::CustomError"},
		{"4test", "test"},
		{"4testE", ""},
		{"N0E", ""},
		{"N0aE", ""},
		{"N01aE", ""},
		{"N1aE", "a"},
		{"N4testE", "test"},
		{"N4test3abcE", "test::abc"},
		{"N4test3_23E", "test::_23"},
		{"N4test3_234abcdE", "test::_23::abcd"},
		{"N4test3_234abcdeE", ""},
		{"N4test3_234abcd1E", ""},
	} {
		if s := simpleDemangleClass(tc[0]); s != tc[1] {
			t.Errorf("demangle(%q) != %q, got %q", tc[0], tc[1], s)
		}
	}
}
