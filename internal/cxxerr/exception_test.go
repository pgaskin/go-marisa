package cxxerr

import (
	"errors"
	"testing"
)

func TestException(t *testing.T) {
	assertExceptionString := func(typ string, std Std, what, str string) {
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

	if p := Std(""); !errors.As(Error("logic_error", "test"), &p) {
		t.Errorf("should be able to convert an Exception into the underlying StdException")
	} else if p != "logic_error" {
		t.Errorf("wrong converted value")
	}

	if p := Std(""); !errors.As(Error("", "test"), &p) {
		t.Errorf("should be able to convert a generic Exception into the underlying StdException")
	} else if p != "exception" {
		t.Errorf("wrong converted value")
	}

	if p := Std(""); !errors.As(&Exception{"test::example_error", "test", "logic_error"}, &p) {
		t.Errorf("should be able to convert an Exception with a non-std type into the underlying StdException")
	} else if p != "logic_error" {
		t.Errorf("wrong converted value")
	}

	if p := Std(""); !errors.As(&Exception{"test::example_error", "test", "exception"}, &p) {
		t.Errorf("should be able to convert an Exception with a non-std type into the underlying generic StdException")
	} else if p != "exception" {
		t.Errorf("wrong converted value")
	}

	if !errors.Is(&Exception{"std::system_error", "test", "system_error"}, Std("")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"std::system_error", "test", "system_error"}, Std("exception")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"std::system_error", "test", "system_error"}, Std("runtime_error")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"std::system_error", "test", "system_error"}, Std("system_error")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"test::example_error", "test", "system_error"}, Std("")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"test::example_error", "test", "system_error"}, Std("exception")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"test::example_error", "test", "system_error"}, Std("runtime_error")) {
		t.Errorf("errors.Is should work on StdException values")
	}
	if !errors.Is(&Exception{"test::example_error", "test", "system_error"}, Std("system_error")) {
		t.Errorf("errors.Is should work on StdException values")
	}

	assertExceptionString("test::something_else", "", "test",
		"test::something_else: test")
	if p := Std(""); errors.As(&Exception{"test::something_else", "test", ""}, &p) {
		t.Errorf("should not get StdException for a non-exception type")
	}
	if errors.Is(&Exception{"test::something_else", "test", ""}, Std("")) {
		t.Errorf("should not get StdException for a non-exception type")
	}
}

func TestStdException(t *testing.T) {
	if !errors.Is(Std(""), Std("")) {
		t.Errorf("empty error type should match itself")
	}
	if errors.Is(Std(""), Std("exception")) {
		t.Errorf("empty error type should not match anything else")
	}
	if errors.Is(Std(""), Std("runtime_error")) {
		t.Errorf("empty error type should not match anything else")
	}
	if errors.Is(Std(""), Std("system_error")) {
		t.Errorf("empty error type should not match anything else")
	}
	if errors.Is(Std(""), Std("ios_base::failure")) {
		t.Errorf("empty error type should not match anything else")
	}
	if errors.Is(Std(""), Std("sdfsdf")) {
		t.Errorf("empty error type should not match anything else")
	}

	if !errors.Is(Std("exception"), Std("")) {
		t.Errorf("base error type should match the empty error type")
	}
	if !errors.Is(Std("exception"), Std("exception")) {
		t.Errorf("base error type should match the base error type")
	}
	if errors.Is(Std("exception"), Std("runtime_error")) {
		t.Errorf("base error type should not match anything more specific")
	}
	if errors.Is(Std("exception"), Std("system_error")) {
		t.Errorf("base error type should not match anything more specific")
	}
	if errors.Is(Std("exception"), Std("ios_base::failure")) {
		t.Errorf("base error type should not match anything more specific")
	}
	if errors.Is(Std("exception"), Std("sdfsdf")) {
		t.Errorf("base error type should not match an unknown error")
	}

	if !errors.Is(Std("runtime_error"), Std("")) {
		t.Errorf("error type should match the empty error type")
	}
	if !errors.Is(Std("runtime_error"), Std("exception")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(Std("runtime_error"), Std("runtime_error")) {
		t.Errorf("error type should match itself")
	}
	if errors.Is(Std("runtime_error"), Std("system_error")) {
		t.Errorf("error type should not match anything more specific")
	}
	if errors.Is(Std("runtime_error"), Std("ios_base::failure")) {
		t.Errorf("error type should not match anything more specific")
	}
	if errors.Is(Std("runtime_error"), Std("sdfsdf")) {
		t.Errorf("error type should not match an unknown error")
	}
	if errors.Is(Std("runtime_error"), Std("invalid_argument")) {
		t.Errorf("error type should not match another error")
	}

	if !errors.Is(Std("system_error"), Std("")) {
		t.Errorf("error type should match the empty error type")
	}
	if !errors.Is(Std("system_error"), Std("exception")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(Std("system_error"), Std("runtime_error")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(Std("system_error"), Std("system_error")) {
		t.Errorf("error type should match itself")
	}
	if errors.Is(Std("system_error"), Std("ios_base::failure")) {
		t.Errorf("error type should not match anything more specific")
	}
	if errors.Is(Std("system_error"), Std("sdfsdf")) {
		t.Errorf("error type should not match an unknown error")
	}
	if errors.Is(Std("system_error"), Std("invalid_argument")) {
		t.Errorf("error type should not match another error")
	}

	if !errors.Is(Std("ios_base::failure"), Std("")) {
		t.Errorf("error type should match the empty error type")
	}
	if !errors.Is(Std("ios_base::failure"), Std("exception")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(Std("ios_base::failure"), Std("runtime_error")) {
		t.Errorf("error type should match the base error type")
	}
	if !errors.Is(Std("ios_base::failure"), Std("system_error")) {
		t.Errorf("error type should match itself")
	}
	if !errors.Is(Std("ios_base::failure"), Std("ios_base::failure")) {
		t.Errorf("error type should match itself")
	}
	if errors.Is(Std("ios_base::failure"), Std("sdfsdf")) {
		t.Errorf("error type should not match an unknown error")
	}
	if errors.Is(Std("ios_base::failure"), Std("invalid_argument")) {
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
