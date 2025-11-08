package marisa_test

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"testing"

	"github.com/pgaskin/go-marisa"
	"github.com/pgaskin/go-marisa/internal/wexcept"
)

func TestIO(t *testing.T) {
	trie := mustWordsTrie()

	expected, err := trie.MarshalBinary()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// note: we've already verified the hash in [TestReproducibility]

	filename := filepath.Join(t.TempDir(), "words.dat")
	if err := os.WriteFile(filename, expected, 0666); err != nil {
		panic(err)
	}

	t.Run("MarshalBinary", func(t *testing.T) {
		if buf, err := trie.MarshalBinary(); err != nil {
			t.Fatalf("error: %v", err)
		} else if !bytes.Equal(buf, expected) {
			t.Errorf("MarshalBinary differs")
		}
	})

	t.Run("AppendBinary", func(t *testing.T) {
		if buf, err := trie.AppendBinary(nil); err != nil {
			t.Fatalf("error: %v", err)
		} else if !bytes.Equal(buf, expected) {
			t.Errorf("AppendBinary(nil) differs")
		}

		if buf, err := trie.AppendBinary(filled(byte(0xFF), len(expected))[:0]); err != nil {
			t.Fatalf("error: %v", err)
		} else if !bytes.Equal(buf, expected) {
			t.Errorf("AppendBinary(make([]byte, 0, exact)) differs")
		}

		if buf, err := trie.AppendBinary(filled(byte(0xFF), len(expected)*2)[:0]); err != nil {
			t.Fatalf("error: %v", err)
		} else if !bytes.Equal(buf, expected) {
			t.Errorf("AppendBinary(make([]byte, 0, extra)) differs")
		}

		if buf, err := trie.AppendBinary(filled(byte(0xFF), 16)[:0]); err != nil {
			t.Fatalf("error: %v", err)
		} else if !bytes.Equal(buf, expected) {
			t.Errorf("AppendBinary(make([]byte, 0, short)) differs")
		}

		if buf, err := trie.AppendBinary(filled(byte(0xFF), 16)[:8]); err != nil {
			t.Fatalf("error: %v", err)
		} else if !bytes.Equal(buf[8:], expected) {
			t.Errorf("AppendBinary(make([]byte, 8, short)) differs")
		} else if !bytes.Equal(buf[:8], filled(byte(0xFF), 8)) {
			t.Errorf("AppendBinary(make([]byte, 8, short)) overwrote existing data")
		}
	})

	t.Run("WriteTo", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if n, err := trie.WriteTo(buf); err != nil {
			t.Fatalf("error: %v", err)
		} else if n != int64(len(expected)) {
			t.Errorf("wrong count %d", n)
		} else if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("WriteTo differs")
		}
		buf.Reset()

		if n, err := trie.WriteTo(&limitedWriter{buf, 1024}); !errors.Is(err, errWriteLimit) {
			if err == nil {
				t.Fatalf("expected write error")
			} else {
				t.Fatalf("wrong write error: %v", err)
			}
		} else if n != int64(buf.Len()) {
			t.Errorf("wrong count %d", n)
		}
		buf.Reset()
	})

	checkTrie := func(trie *marisa.Trie) bool {
		buf, err := trie.MarshalBinary()
		if err != nil {
			return false
		}
		return bytes.Equal(buf, expected)
	}

	t.Run("UnmarshalBinary", func(t *testing.T) {
		var trie marisa.Trie

		if err := trie.UnmarshalBinary(expected); err != nil {
			t.Fatalf("error: %v", err)
		} else if !checkTrie(&trie) {
			t.Errorf("round-trip failed")
		}

		if err := trie.UnmarshalBinary(nil); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("expected unexpected eof for nil (got %v)", err)
		} else if !checkTrie(&trie) {
			t.Errorf("old trie should still be valid on error")
		}

		if err := trie.UnmarshalBinary(expected[:500]); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("expected unexpected eof for truncated (got %v)", err)
		} else if !checkTrie(&trie) {
			t.Errorf("old trie should still be valid on error")
		}

		if err := trie.UnmarshalBinary(append(slices.Clone(expected[:500]), filled(byte(0xFF), 10*1024*1024)...)); !errors.Is(err, wexcept.RuntimeError) {
			t.Fatalf("expected runtime error for junk (got %v)", err)
		} else if !checkTrie(&trie) {
			t.Errorf("old trie should still be valid on error")
		}

		if err := trie.UnmarshalBinary(filled(byte(0xFF), 10*1024*1024)); !errors.Is(err, wexcept.RuntimeError) {
			t.Fatalf("expected runtime error for junk (got %v)", err)
		} else if !checkTrie(&trie) {
			t.Errorf("old trie should still be valid on error")
		}
	})

	t.Run("ReadFrom", func(t *testing.T) {
		var trie marisa.Trie

		if n, err := trie.ReadFrom(bytes.NewReader(expected)); err != nil {
			t.Fatalf("error: %v", err)
		} else if n != int64(len(expected)) {
			t.Errorf("incorrect n %d", n)
		} else if !checkTrie(&trie) {
			t.Errorf("round-trip failed")
		}

		if n, err := trie.ReadFrom(&shortReader{bytes.NewReader(expected), 5}); err != nil {
			t.Fatalf("error: %v", err)
		} else if n != int64(len(expected)) {
			t.Errorf("incorrect n %d", n)
		} else if !checkTrie(&trie) {
			t.Errorf("round-trip failed")
		}

		if n, err := trie.ReadFrom(bytes.NewReader(nil)); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("expected unexpected eof for nil (got %v)", err)
		} else if n != 0 {
			t.Errorf("incorrect n %d", n)
		} else if !checkTrie(&trie) {
			t.Errorf("old trie should still be valid on error")
		}

		if n, err := trie.ReadFrom(bytes.NewReader(expected[:500])); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("expected unexpected eof for truncated (got %v)", err)
		} else if n != 500 {
			t.Errorf("incorrect n %d", n)
		} else if !checkTrie(&trie) {
			t.Errorf("old trie should still be valid on error")
		}

		cr := &countReader{R: bytes.NewReader(append(slices.Clone(expected[:500]), filled(byte(0xFF), 10*1024*1024)...))}
		if n, err := trie.ReadFrom(cr); !errors.Is(err, wexcept.RuntimeError) {
			t.Fatalf("expected runtime error for junk (got %v)", err)
		} else if cr.N != n {
			t.Errorf("incorrect n %d", n)
		} else if n != 152136 { // this may need to be updated if marisa-trie is changed
			t.Errorf("incorrect n %d", n)
		} else if !checkTrie(&trie) {
			t.Errorf("old trie should still be valid on error")
		}

		cr = &countReader{R: bytes.NewReader(filled(byte(0xFF), 10*1024*1024))}
		if n, err := trie.ReadFrom(cr); !errors.Is(err, wexcept.RuntimeError) {
			t.Fatalf("expected runtime error for junk (got %v)", err)
		} else if cr.N != n {
			t.Errorf("incorrect n %d", n)
		} else if n != 16 { // this may need to be updated if marisa-trie is changed
			t.Errorf("incorrect n %d", n)
		} else if !checkTrie(&trie) {
			t.Errorf("old trie should still be valid on error")
		}
	})

	t.Run("MapFile", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			f, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			var trie marisa.Trie
			if err := trie.MapFile(f, 0, int64(len(expected))); err != nil {
				if errors.Is(err, errors.ErrUnsupported) {
					t.Skipf("unsupported platform: %v", err)
				}
				t.Errorf("error: %v", err)
			} else if !checkTrie(&trie) {
				t.Errorf("round-trip failed")
			}
		})

		t.Run("Aligned", func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "offset.dat")
			if err := os.WriteFile(filename, slices.Concat(filled(byte(0xFF), syscall.Getpagesize()*3), expected, filled(byte(0xFF), syscall.Getpagesize())), 0666); err != nil {
				panic(err)
			}

			f, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			var trie marisa.Trie
			if err := trie.MapFile(f, int64(syscall.Getpagesize()*3), int64(len(expected))); err != nil {
				if errors.Is(err, errors.ErrUnsupported) {
					t.Skipf("unsupported platform: %v", err)
				}
				t.Errorf("error: %v", err)
			} else if !checkTrie(&trie) {
				t.Errorf("round-trip failed")
			}
		})

		t.Run("Unaligned", func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "offset.dat")
			if err := os.WriteFile(filename, slices.Concat(filled(byte(0xFF), 32), expected, filled(byte(0xFF), 32)), 0666); err != nil {
				panic(err)
			}

			f, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			var trie marisa.Trie
			if err := trie.MapFile(f, 32, int64(len(expected))); err != nil {
				if errors.Is(err, errors.ErrUnsupported) {
					t.Skipf("unsupported platform: %v", err)
				}
				t.Errorf("error: %v", err)
			} else if !checkTrie(&trie) {
				t.Errorf("round-trip failed")
			}
		})
	})

	t.Run("Open", func(t *testing.T) {
		if trie, err := marisa.Open(filename); err != nil {
			if errors.Is(err, errors.ErrUnsupported) {
				t.Errorf("should have fallen back from mmap to regular open")
			}
			t.Errorf("error: %v", err)
		} else if !checkTrie(trie) {
			t.Errorf("round-trip failed")
		}

		if trie, err := marisa.Open(filepath.Join(t.TempDir(), "nonexistent")); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("expected not found, got: %v", err)
		} else if trie != nil {
			t.Errorf("expected trie to be nil if not found")
		}
	})

	t.Run("New", func(t *testing.T) {
		if trie, err := marisa.New(expected); err != nil {
			t.Errorf("error: %v", err)
		} else if !checkTrie(trie) {
			t.Errorf("round-trip failed")
		}
	})

	t.Run("Load", func(t *testing.T) {
		if trie, err := marisa.Load(bytes.NewReader(expected)); err != nil {
			t.Errorf("error: %v", err)
		} else if !checkTrie(trie) {
			t.Errorf("round-trip failed")
		}
	})
}

func filled[T any](v T, n int) []T {
	s := make([]T, n)
	for i := range n {
		s[i] = v
	}
	return s
}

var errWriteLimit = errors.New("limit reached")

type limitedWriter struct {
	W io.Writer
	N int64
}

func (w *limitedWriter) Write(p []byte) (n int, err error) {
	if int64(len(p)) > w.N {
		return 0, errWriteLimit
	}
	n, err = w.W.Write(p)
	if err == nil && n != len(p) {
		err = io.ErrShortWrite
	}
	w.N -= int64(n)
	return
}

type shortReader struct {
	R io.Reader
	N int
}

func (r *shortReader) Read(p []byte) (n int, err error) {
	if len(p) > r.N {
		p = p[:r.N]
	}
	n, err = r.R.Read(p)
	if err == nil && n != len(p) {
		err = io.ErrShortWrite
	}
	return
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
