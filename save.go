package marisa

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/pgaskin/go-marisa/internal/wexcept"
	"github.com/pgaskin/go-marisa/internal/wexport"
	"github.com/tetratelabs/wazero/api"
)

// MarshalBinary serializes the dictionary.
func (t *Trie) MarshalBinary() ([]byte, error) {
	if t.mod == nil {
		return nil, errors.New("dictionary not initialized")
	}
	return t.AppendBinary(nil)
}

func (t *Trie) AppendBinary(b []byte) ([]byte, error) {
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
				fmt.Println("sdfsdf", len(*b), len(x))
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
