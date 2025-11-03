package wautil

import (
	"context"
	"fmt"
	"io"

	"github.com/tetratelabs/wazero/api"
)

// Handle represents a handle within the current call stack.
type Handle uint32

// NewHandle creates a handle for a Go object to be used in the current call stack.
func NewHandle(ctx context.Context, a any) (context.Context, Handle) {
	if a == nil {
		return ctx, 0
	}
	idx, _ := ctx.Value(Handle(0)).(Handle)
	idx++
	ctx = context.WithValue(ctx, Handle(0), idx)
	ctx = context.WithValue(ctx, Handle(idx), a)
	return ctx, idx
}

// GetHandle dereferences a handle, panicking if it's invalid (note that this
// will be caught be wazero and turned into an error for the [api.Function]
// call).
func GetHandle[T any](ctx context.Context, handle Handle) T {
	if handle == 0 {
		var z T
		return z
	}
	x, ok := ctx.Value(handle).(T)
	if !ok {
		panic(fmt.Errorf("gocpp: invalid %T handle %d", x, handle))
	}
	return x
}

// note: we're not using externrefs since they're unnecessarily complicated to
// store and have no real benefit other than typing

var _ = register(ExportFuncIIII("cxx_write", func(ctx context.Context, mod api.Module, handle Handle, ptr, size uint32) uint32 {
	w := GetHandle[io.Writer](ctx, handle)
	if w == nil {
		panic("gocpp: nil handle")
	}
	b, ok := mod.Memory().Read(ptr, size)
	if !ok {
		panic("gocpp: invalid pointer")
	}
	n, err := w.Write(b)
	if err != nil {
		Throw(err)
	}
	if n != len(b) {
		Throw(io.ErrShortWrite)
	}
	return uint32(n)
}, "io.Writer.Write", "handle", "ptr", "size", "n"))

var _ = register(ExportFuncIIII("cxx_read", func(ctx context.Context, mod api.Module, handle Handle, ptr, size uint32) uint32 {
	r := GetHandle[io.Reader](ctx, handle)
	if r == nil {
		panic("gocpp: nil handle")
	}
	b, ok := mod.Memory().Read(ptr, size)
	if !ok {
		panic("gocpp: invalid pointer")
	}
	for {
		n, err := r.Read(b)
		if err == io.EOF {
			return 0
		}
		if err != nil {
			Throw(err)
		}
		if n > 0 {
			return uint32(n)
		}
	}
}, "io.Reader.Read", "handle", "ptr", "size", "n"))
