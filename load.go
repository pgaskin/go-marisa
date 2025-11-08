package marisa

import (
	"context"
	"errors"
	"io"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/pgaskin/go-marisa/internal/walloc"
	"github.com/pgaskin/go-marisa/internal/wexcept"
	"github.com/pgaskin/go-marisa/internal/wexport"
	"github.com/tetratelabs/wazero/api"
)

// Open opens a dictionary from a file.
func Open(name string) (*Trie, error) {
	var t Trie

	// only try mmap if it's likely to succeed and it's on a fully tested platform
	if (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && (runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64") {
		// attempt to get the size (and if it's not seekable, it's unlikely to be mappable either)
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if size, err := f.Seek(0, io.SeekEnd); err == nil {
			if err := t.MapFile(f, 0, size); err == nil {
				return &t, nil
			}
		}
	}

	// read the entire dictionary
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return New(b)
}

// New is shorthand for initializing a dictionary with [Trie.UnmarshalBinary].
// Using Load with a [bytes.Reader] may result in a more optimal in-memory
// layout.
func New(b []byte) (*Trie, error) {
	var t Trie
	if err := t.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return &t, nil
}

// Load is shorthand for initializing a dictionary with [Trie.ReadFrom].
func Load(r io.Reader) (*Trie, error) {
	var t Trie
	if _, err := t.ReadFrom(r); err != nil {
		return nil, err
	}
	return &t, nil
}

// MapFile mmaps a file and loads the dictionary from it. On error, the trie is
// left unchanged. If not supported by the current platform, an error matching
// [errors.ErrUnsupported] is returned.
func (t *Trie) MapFile(f *os.File, offset int64, length int64) error {
	if uint64(length) > maxAlloc {
		return errors.New("dictionary too large")
	}
	va := &walloc.VirtualAllocator{
		Fallback: &walloc.SliceAllocator{
			// if it falls back to this, we'll be returning an error anyways
			OverrideMax: scratchSpace,
		},
		// on 32-bit hosts, it's critical for this (unlike the SliceAllocator)
		// since it immediately reserves virtual address space, which we only
		// have 4 GiB of
		OverrideMax: uint64(length) + scratchSpace,
	}
	mod, err := instantiate(va)
	if err != nil {
		return err
	}
	if err := va.Err(); err != nil {
		return err
	}
	ptr, err := va.MapFile(context.Background(), mod.Module(), f, offset, length, false)
	if err != nil {
		return err
	}
	if _, err := mod.Call("marisa_new", uint64(ptr), uint64(length)); err != nil {
		var ex *wexcept.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, wexcept.StdException("runtime_error")) && strings.Contains(ex.What(), "size > avail_") {
				err = io.ErrUnexpectedEOF
			}
		}
		return err
	}
	return t.swap(mod)
}

// UnmarshalBinary copies b and maps the trie directly from it. This is faster
// than [Trie.ReadFrom], but may have a less optimal memory layout. On error,
// the trie is left unchanged.
func (t *Trie) UnmarshalBinary(b []byte) error {
	if uint64(len(b)) > min(math.MaxUint32, math.MaxInt) {
		return errors.New("dictionary too large")
	}
	sa := &walloc.SliceAllocator{
		OverrideMax: uint64(len(b)) + scratchSpace,
	}
	mod, err := instantiate(sa)
	if err != nil {
		return err
	}
	ptr, err := mod.Alloc(len(b))
	if err != nil {
		return err
	}
	if buf, ok := mod.Module().Memory().Read(ptr, uint32(len(b))); !ok {
		panic("bad allocation")
	} else {
		copy(buf, b)
	}
	if _, err := mod.Call("marisa_new", uint64(ptr), uint64(len(b))); err != nil {
		var ex *wexcept.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, wexcept.StdException("runtime_error")) && strings.Contains(ex.What(), "size > avail_") {
				err = io.ErrUnexpectedEOF
			}
		}
		return err
	}
	return t.swap(mod)
}

// ReadFrom reads a dictionary from r. On success, it will have read exactly the
// size of the dictionary. On error, the trie is left unchanged.
func (t *Trie) ReadFrom(r io.Reader) (int64, error) {
	// note: it won't actually read past in practice, since it reads exactly
	// what it wants with std::istream::read, and our stream impl is effectively
	// unbuffered
	sa := &walloc.SliceAllocator{
		OverrideMax: maxAlloc,
	}
	mod, err := instantiate(sa)
	if err != nil {
		return 0, err
	}
	c := &countReader{R: r}
	if _, err := mod.CallContext(withReader(context.Background(), c), "marisa_load"); err != nil {
		var ex *wexcept.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, wexcept.StdException("runtime_error")) && strings.Contains(ex.What(), "!stream_->read") {
				err = io.ErrUnexpectedEOF
			}
		}
		return c.N, err
	}
	return c.N, t.swap(mod)
}

type zeroReader struct{}

func (z zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

type readKey struct{}

func withReader(ctx context.Context, r io.Reader) context.Context {
	return context.WithValue(ctx, readKey{}, r)
}

var read = wexport.VII("read", func(ctx context.Context, m api.Module, p, n uint32) {
	r, ok := ctx.Value(readKey{}).(io.Reader)
	if !ok {
		panic("no active reader")
	}
	if n != 0 {
		if p != 0 {
			b, ok := m.Memory().Read(p, n)
			if !ok {
				panic("invalid pointer")
			}
			if _, err := io.ReadFull(r, b); err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				wexcept.Throw(err)
			}
		} else {
			if _, err := io.CopyN(io.Discard, r, int64(n)); err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				wexcept.Throw(err)
			}
		}
	}
}, "read", "buf", "n")

type countReader struct {
	N int64
	R io.Reader
}

func (c *countReader) Read(p []byte) (n int, err error) {
	n, err = c.R.Read(p)
	c.N += int64(n)
	return
}
