package marisa

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"math"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/pgaskin/go-marisa/internal/cxxerr"
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
		var ex *cxxerr.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, cxxerr.Std("runtime_error")) && strings.Contains(ex.What(), "size > avail_") {
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
		var ex *cxxerr.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, cxxerr.Std("runtime_error")) && strings.Contains(ex.What(), "size > avail_") {
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
		var ex *cxxerr.Exception
		if errors.As(err, &ex) {
			if errors.Is(ex, cxxerr.Std("runtime_error")) && strings.Contains(ex.What(), "!stream_->read") {
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

// MarshalBinary serializes the dictionary.
func (t *Trie) MarshalBinary() ([]byte, error) {
	return t.AppendBinary(nil)
}

func (t *Trie) AppendBinary(b []byte) ([]byte, error) {
	if t.mod == nil {
		return nil, errors.New("dictionary not initialized")
	}
	b = slices.Grow(b, int(t.ioSize))
	_, err := t.mod.CallContext(withWriteBuffer(context.Background(), &b), "marisa_save")
	return b, err
}

// WriteTo serializes the dictionary to w.
func (t *Trie) WriteTo(w io.Writer) (int64, error) {
	if t.mod == nil {
		return 0, errors.New("dictionary not initialized")
	}
	c := &countWriter{W: w}
	_, err := t.mod.CallContext(withWriter(context.Background(), c), "marisa_save")
	return c.N, err
}

type writeKey string

func withWriter(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, writeKey("writer"), w)
}

func withWriteBuffer(ctx context.Context, b *[]byte) context.Context {
	return context.WithValue(ctx, writeKey("buffer"), b)
}

var write = wexport.VII("write", func(ctx context.Context, m api.Module, p, n uint32) {
	if w, ok := ctx.Value(writeKey("writer")).(io.Writer); ok {
		if n != 0 {
			if p != 0 {
				b, ok := m.Memory().Read(p, n)
				if !ok {
					panic("invalid pointer")
				}
				maybePanicAtOffset(w, n, b)
				n, err := w.Write(b)
				if err != nil {
					wexcept.Throw(err)
				}
				if n != len(b) {
					wexcept.Throw(io.ErrShortWrite)
				}
			} else {
				maybePanicAtOffset(w, n, nil)
				if _, err := io.CopyN(w, zeroReader{}, int64(n)); err != nil {
					wexcept.Throw(err)
				}
			}
		}
		return
	}
	if b, ok := ctx.Value(writeKey("buffer")).(*[]byte); ok {
		if n != 0 {
			if p != 0 {
				x, ok := m.Memory().Read(p, n)
				if !ok {
					panic("invalid pointer")
				}
				*b = append(*b, x...)
			} else {
				*b = append(*b, make([]byte, n)...)
			}
		}
		return
	}
	panic("no active writer")
}, "write", "buf", "n")

type countWriter struct {
	N int64
	W io.Writer
}

func (c *countWriter) Write(p []byte) (n int, err error) {
	n, err = c.W.Write(p)
	c.N += int64(n)
	return
}

// debugPanicAtOffset causes a stack trace to be printed when a write overlaps
// the specified offsets. This is intended for debugging, or for figuring out
// what exactly a specific offset is for.
var debugPanicAtOffset uint32 = math.MaxUint32

func init() {
	if s := os.Getenv("MARISA_DEBUG_PANIC_AT_OFFSET"); s != "" {
		if n, err := strconv.ParseUint(s, 0, 32); err == nil {
			debugPanicAtOffset = uint32(n)
		}
	}
}

func maybePanicAtOffset(a any, n uint32, b []byte) {
	if debugPanicAtOffset == math.MaxUint32 {
		return
	}
	var o uint32
	switch a := a.(type) {
	case *countWriter:
		o = uint32(a.N)
	default:
		return
	}
	if o <= debugPanicAtOffset && debugPanicAtOffset-o < n {
		var s strings.Builder
		s.WriteString("write (MARISA_DEBUG_PANIC_AT_OFFSET)")
		s.WriteString("\nmatch: ")
		s.WriteString(strconv.FormatUint(uint64(o), 10))
		s.WriteString(" <= ")
		s.WriteString(strconv.FormatUint(uint64(debugPanicAtOffset), 10))
		s.WriteString(" < ")
		s.WriteString(strconv.FormatUint(uint64(o+n), 10))
		s.WriteString(" (")
		s.WriteString(strconv.FormatUint(uint64(n), 10))
		s.WriteString(")")
		s.WriteString("\ndata:")
		if b != nil {
			for _, l := range strings.Split(hex.Dump(b), "\n") {
				s.WriteString("\n\t")
				s.WriteString(l)
			}
		} else {
			s.WriteString("zeros")
		}
		s.WriteString("\n")
		panic(s.String())
	}
}
