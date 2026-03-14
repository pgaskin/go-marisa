package marisa

import (
	"errors"
	"io"
	"math"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/pgaskin/go-marisa/internal/cxxerr"
	"github.com/pgaskin/go-marisa/internal/wexcept"
	"github.com/pgaskin/go-marisa/internal/wmem"
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
	va, err := wmem.VirtualMemory(uint64(length), uint64(length)+scratchSpace)
	if err != nil {
		return err
	}
	mod, err := instantiate(va)
	if err != nil {
		return err
	}
	ptr, err := wmem.MapFile(mod.marisa, mod.mem, f, offset, length, false)
	if err != nil {
		return err
	}
	if err := func() (err error) {
		defer wexcept.Catch(&err)
		mod.marisa.XNew(int32(ptr), int32(length))
		return
	}(); err != nil {
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
	sa := wmem.SliceMemory(uint64(len(b)), uint64(len(b))+scratchSpace)
	mod, err := instantiate(sa)
	if err != nil {
		return err
	}
	ptr, err := mod.Alloc(len(b))
	if err != nil {
		return err
	}
	if buf, ok := wmem.Bytes(mod.mem, ptr, uint32(len(b))); !ok {
		panic("bad allocation")
	} else {
		copy(buf, b)
	}
	if err := func() (err error) {
		defer wexcept.Catch(&err)
		mod.marisa.XNew(int32(ptr), int32(uint32(len(b))))
		return
	}(); err != nil {
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
	sa := wmem.SliceMemory(8192, maxAlloc)
	mod, err := instantiate(sa)
	if err != nil {
		return 0, err
	}
	c := &countReader{R: r}
	if err := func() (err error) {
		defer wexcept.Catch(&err)
		mod.io.Reader = c
		defer func() { mod.io.Reader = nil }()
		mod.marisa.XLoad()
		return
	}(); err != nil {
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
	err := func() (err error) {
		defer wexcept.Catch(&err)
		t.mod.io.WriteBuffer = &b
		defer func() { t.mod.io.WriteBuffer = nil }()
		t.mod.marisa.XSave()
		return
	}()
	return b, err
}

// WriteTo serializes the dictionary to w.
func (t *Trie) WriteTo(w io.Writer) (int64, error) {
	if t.mod == nil {
		return 0, errors.New("dictionary not initialized")
	}
	c := &countWriter{W: w}
	err := func() (err error) {
		defer wexcept.Catch(&err)
		t.mod.io.Writer = c
		defer func() { t.mod.io.Writer = nil }()
		t.mod.marisa.XSave()
		return
	}()
	return c.N, err
}

type countWriter struct {
	N int64
	W io.Writer
}

func (c *countWriter) Write(p []byte) (n int, err error) {
	n, err = c.W.Write(p)
	c.N += int64(n)
	return
}

func (m *module) Xread(p, n int32) {
	if m.io.Reader == nil {
		panic("no active reader")
	}
	if n != 0 {
		if p != 0 {
			b, ok := wmem.Bytes(m.mem, uint32(p), uint32(n))
			if !ok {
				panic("invalid pointer")
			}
			if _, err := io.ReadFull(m.io.Reader, b); err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				wexcept.Throw(err)
			}
		} else {
			if _, err := io.CopyN(io.Discard, m.io.Reader, int64(n)); err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				wexcept.Throw(err)
			}
		}
	}
}

func (m *module) Xwrite(p, n int32) {
	if w := m.io.Writer; w != nil {
		if n != 0 {
			if p != 0 {
				b, ok := wmem.Bytes(m.mem, uint32(p), uint32(n))
				if !ok {
					panic("invalid pointer")
				}
				n, err := w.Write(b)
				if err != nil {
					wexcept.Throw(err)
				}
				if n != len(b) {
					wexcept.Throw(io.ErrShortWrite)
				}
			} else {
				if _, err := io.CopyN(w, zeroReader{}, int64(n)); err != nil {
					wexcept.Throw(err)
				}
			}
		}
		return
	}
	if b := m.io.WriteBuffer; b != nil {
		if n != 0 {
			if p != 0 {
				x, ok := wmem.Bytes(m.mem, uint32(p), uint32(n))
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
}
