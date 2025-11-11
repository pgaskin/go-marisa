//go:build ignore

// Command download fetches, patches, and amalgamates marisa-trie.
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"
)

const (
	Version = "v0.3.1"
	Repo    = "s-yata/marisa-trie"
	Commit  = "3e87d53b78e15f2f43783d5e376561a8c9722051"
)

var (
	Search = []string{
		"include",
		"lib",
	}
	Replace = map[string]string{
		"include/marisa/iostream.h":    "",
		"include/marisa/stdio.h":       "",
		"lib/marisa/grimoire/intrin.h": "../src/intrin.h",
		"lib/marisa/grimoire/io.h":     "../src/io.h",
	}
	Headers = []string{
		"include/marisa.h",
	}
	Sources = []string{
		"lib/marisa/grimoire/trie/tail.cc",
		"lib/marisa/grimoire/trie/louds-trie.cc",
		"lib/marisa/grimoire/vector/bit-vector.cc",
		"lib/marisa/agent.cc",
		"lib/marisa/keyset.cc",
		"lib/marisa/trie.cc",
	}
	Warnings = []string{
		"-Wall",
		"-Wextra",
		"-Wconversion",
		"-Wno-implicit-fallthrough",
		"-Wno-unused-function",
		"-Wno-unused-const-variable",
	}
	Source = "../lib/marisa.cc"
	Header = "../lib/marisa.h"
)

func patch(files map[string][]byte) error {
	{
		slog.Info("ensuring there aren't catch statements since we don't support stack unwinding")
		for name, src := range files {
			if !strings.HasPrefix(name, "lib/") && !strings.HasPrefix(name, "include/") {
				continue
			}
			if !strings.HasSuffix(name, ".h") && !strings.HasSuffix(name, ".cc") {
				continue
			}
			for line := range bytes.Lines(src) {
				if bytes.HasPrefix(bytes.TrimLeftFunc(line, isSpaceASCII), []byte("//")) {
					continue
				}
				if bytes.Contains(line, []byte("catch")) { // this is very rudimentary and can be improved if necessary
					return fmt.Errorf("found possible unsupported catch statement in %q", name)
				}
			}
		}
	}
	{
		slog.Info("patching lib/marisa/grimoire/vector/vector.h to ensure unused space in Vector is zeroed")
		name := "lib/marisa/grimoire/vector/vector.h"
		src, n := files[name], 0
		src, n = bytesTryReplaceAll(src,
			[]byte(`new_buf(new char[sizeof(T) * new_capacity]);`),
			[]byte(`new_buf(new char[sizeof(T) * new_capacity]());`))
		if n != 1 {
			return fmt.Errorf("failed to apply patch")
		}
		src, n = bytesTryReplaceAll(src,
			[]byte(`new_buf(new char[sizeof(T) * capacity]);`),
			[]byte(`new_buf(new char[sizeof(T) * capacity]());`))
		if n != 1 {
			return fmt.Errorf("failed to apply patch")
		}
		files[name] = src
	}
	{
		slog.Info("patching lib/marisa/grimoire/**/*.h to ensure i/o methods aren't inlined for better stack traces")
		for name, src := range files {
			if !strings.HasPrefix(name, "lib/marisa/grimoire/") || !strings.HasSuffix(name, ".h") {
				continue
			}
			src = bytes.ReplaceAll(src, []byte(`void map(Mapper`), []byte(`[[clang::noinline]] void map(Mapper`))
			src = bytes.ReplaceAll(src, []byte(`void read(Reader`), []byte(`[[clang::noinline]] void read(Reader`))
			src = bytes.ReplaceAll(src, []byte(`void write(Writer`), []byte(`[[clang::noinline]] void write(Writer`))
			files[name] = src
		}
	}
	return nil
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if err := run(); err != nil {
		slog.Error("failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	slog.Info("fetching code", "repo", Repo, "version", Version, "commit", Commit)
	files, err := fetch(Repo, Commit)
	if err != nil {
		return fmt.Errorf("fetch code: %w", err)
	}
	if err := patch(files); err != nil {
		return fmt.Errorf("apply patches: %w", err)
	}

	var pre bytes.Buffer
	pre.WriteString("// AUTOMATICALLY GENERATED, DO NOT EDIT!\n")
	pre.WriteString("// marisa-trie ")
	pre.WriteString(Version)
	pre.WriteString(" (")
	pre.WriteString(Commit)
	pre.WriteString(") (patched)\n")

	for _, old := range slices.Sorted(maps.Keys(Replace)) {
		if new := Replace[old]; new != "" {
			slog.Info("replacing include", "old", old, "new", new)
		} else {
			slog.Info("removing include", "old", old)
		}
	}

	copying, ok := files["COPYING.md"]
	if !ok {
		return fmt.Errorf("missing copying text")
	}
	pre.WriteString("//\n")
	for line := range bytes.Lines(bytes.Trim(copying, "\n")) {
		pre.WriteString("// ")
		pre.Write(bytes.TrimRightFunc(line, isSpaceASCII))
		pre.WriteString("\n")
	}
	pre.WriteString("//\n")
	pre.WriteString("// clang-format off\n")
	pre.WriteString("\n")

	repl := maps.Clone(Replace)

	slog.Info("amalgamating headers", "src", strings.Join(Sources, ","))
	h, hinc, err := amalgamate(files, Search, repl, Headers...)
	if err != nil {
		return fmt.Errorf("generate headers: %w", err)
	}
	for _, name := range hinc {
		repl[name] = "./" + Header
	}
	h = slices.Concat(pre.Bytes(), h)

	pre.WriteString("#ifdef __clang__\n")
	warnings(&pre, "clang", Warnings...)
	pre.WriteString("#elif __GNUC__\n")
	warnings(&pre, "GCC", Warnings...)
	pre.WriteString("#endif\n")
	pre.WriteString("\n")

	slog.Info("amalgamating sources", "src", strings.Join(Sources, ","))
	cc, _, err := amalgamate(files, Search, repl, Sources...)
	if err != nil {
		return fmt.Errorf("generate sources: %w", err)
	}
	for _, name := range hinc {
		repl[name] = "./" + Source
	}
	cc = slices.Concat(pre.Bytes(), cc)

	slog.Info("writing output")
	return errors.Join(
		os.WriteFile(Header, h, 0666),
		os.WriteFile(Source, cc, 0666),
	)
}

// fetch downloads a git repository tarball.
func fetch(repo, rev string) (map[string][]byte, error) {
	resp, err := http.Get("https://github.com/" + repo + "/archive/" + url.PathEscape(rev) + ".tar.gz")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %d", resp.StatusCode)
	}

	zr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var (
		tr    = tar.NewReader(zr)
		files = map[string][]byte{}
		pfx   string
	)
	for {
		f, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if f.Name == "pax_global_header" || f.FileInfo().IsDir() {
			continue
		}

		if pfx == "" {
			rest := f.Name
			rest, _ = strings.CutPrefix(rest, "./")
			_, rest, _ = strings.Cut(rest, "/")
			if rest == "" {
				continue
			}
			pfx = f.Name[:len(f.Name)-len(rest)]
		}

		name, ok := strings.CutPrefix(f.Name, pfx)
		if !ok {
			return nil, fmt.Errorf("file %q is not in detected top-level directory %q", f.Name, pfx)
		}

		buf, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		files[name] = buf
	}
	return files, nil
}

// warnings writes pragma warning lines for the specified warning flags.
func warnings(buf *bytes.Buffer, compiler string, flags ...string) {
	for _, flag := range flags {
		flag, ign := strings.CutPrefix(flag, "-Wno-")
		if ign {
			flag = "-W" + flag
		}
		buf.WriteString("#pragma ")
		buf.WriteString(compiler)
		buf.WriteString(" diagnostic ")
		if ign {
			buf.WriteString("ignored")
		} else {
			buf.WriteString("warning")
		}
		buf.WriteByte(' ')
		buf.WriteByte('"')
		buf.WriteString(flag)
		buf.WriteByte('"')
		buf.WriteString("\n")
	}
}

// amalgamate amalgamates the specified sources. Includes are searched in order.
// If the resolved import is in the replace map, it is replaced with another
// include or removed accordingly.
func amalgamate(files map[string][]byte, search []string, replace map[string]string, sources ...string) ([]byte, []string, error) {
	var (
		why   = map[string][]string{}
		needs []string
		err   error
	)
	for src, inc := range crawlIncludes(files, func(name, inc string) (string, error) {
		res, err := resolveInclude(files, name, inc, search)
		if err != nil {
			return "", nil
		}
		if _, ok := replace[res]; ok {
			return "", nil
		}
		return res, nil
	}, sources...)(&err) {
		by, ok := why[inc]
		if !ok {
			needs = append(needs, inc)
		}
		if !slices.Contains(by, src) {
			by = append(by, src)
		}
		why[inc] = by
	}
	if err != nil {
		return nil, nil, fmt.Errorf("crawl: %w", err)
	}

	for _, src := range sources {
		if !slices.Contains(needs, src) {
			needs = append(needs, src)
		}
	}

	var buf bytes.Buffer
	for _, name := range needs {
		if by, ok := why[name]; ok {
			slog.Debug("inlining", "src", name, "needed_by", strings.Join(by, ","))
		} else {
			slog.Debug("inlining", "src", name)
		}
		if _, ok := replace[name]; ok {
			panic("wtf") // these should have already been filtered out
		}
		src, ok := files[name]
		if !ok {
			return nil, nil, fmt.Errorf("inline %q: file not found", name)
		}
		src, err := replaceIncludes(src, func(inc string, sys bool) (string, bool, string, error) {
			if sys {
				return inc, sys, "", nil
			}
			res, err := resolveInclude(files, name, inc, search)
			if err != nil {
				return "", false, "", err
			}
			r, ok := replace[res]
			if !ok {
				return "", false, "amalgamate " + inc, nil
			}
			if r == "" {
				return r, false, "amalgamate " + inc, nil
			}
			return r, false, "amalgamate " + inc, nil
		})
		if err != nil {
			return nil, nil, fmt.Errorf("inline %q: replace includes: %w", name, err)
		}
		buf.WriteString(`#line 1 "`)
		buf.WriteString(name)
		buf.WriteString(`"`)
		buf.WriteByte('\n')
		buf.Write(src)
	}
	return buf.Bytes(), needs, nil
}

// crawlIncludes recursively crawls non-system includes referenced by the
// specified files, depth-first postorder.
func crawlIncludes(files map[string][]byte, resolve func(name, inc string) (string, error), sources ...string) func(*error) iter.Seq2[string, string] {
	return func(err *error) iter.Seq2[string, string] {
		return func(yield func(string, string) bool) {
			*err = func() error {
				for _, name := range sources {
					src, ok := files[name]
					if !ok {
						return fmt.Errorf("file %q not found", name)
					}
					for line := range bytes.Lines(src) {
						_, include, _, system, ok := cutInclude(line)
						if !ok {
							continue
						}
						if system {
							continue
						}
						res, err := resolve(name, string(include))
						if err != nil {
							err = fmt.Errorf("resolve include %q from %q: %w", include, name, err)
						}
						if res == "" {
							continue // replaced or removed, so don't go deeper
						}
						if crawlIncludes(files, resolve, res)(&err)(yield); err != nil {
							return fmt.Errorf("crawl %q: %w", name, err)
						}
						if !yield(name, res) {
							return nil
						}
					}
				}
				return nil
			}()
		}
	}
}

// replaceIncludes calls fn for each include in src, removing it or replacing it
// with newpath, optionally with a comment appended.
func replaceIncludes(src []byte, fn func(path string, sys bool) (newpath string, newsys bool, comment string, err error)) ([]byte, error) {
	var dst bytes.Buffer
	dst.Grow(len(src))
	for line := range bytes.Lines(src) {
		pre, inc, post, sys, ok := cutInclude(bytes.TrimSuffix(line, []byte{'\n'}))
		if !ok {
			dst.Write(line)
			continue
		}

		newinc, newsys, comment, err := fn(string(inc), sys)
		if err != nil {
			return nil, fmt.Errorf("replace %q: %w", line, err)
		}

		var repl []byte
		if len(newinc) != 0 {
			repl, ok = formatInclude(pre, []byte(newinc), post, newsys)
			if !ok {
				return nil, fmt.Errorf("replace %q: bad repalcement %q", line, repl)
			}
		}
		if comment != "" {
			if len(repl) != 0 {
				repl = append(repl, ' ')
			}
			repl = append(repl, "//"...)
			repl = append(repl, comment...)
		}
		dst.Write(repl)
		dst.WriteByte('\n')
	}
	if old, new := bytes.Count(src, []byte{'\n'}), bytes.Count(dst.Bytes(), []byte{'\n'}); old != new {
		panic("wtf") // line count should be identical
	}
	return dst.Bytes(), nil
}

// cutInclude splits an #include line.
func cutInclude(line []byte) (prefix, path, suffix []byte, sys, ok bool) {
	path = line
	line = bytes.TrimLeftFunc(path, isSpaceASCII)
	path, ok = bytes.CutPrefix(path, []byte("#"))
	if !ok {
		return line, nil, nil, false, false
	}
	path = bytes.TrimLeftFunc(path, isSpaceASCII)
	path, ok = bytes.CutPrefix(path, []byte("include"))
	if !ok {
		return line, nil, nil, false, false
	}
	path = bytes.TrimLeftFunc(path, isSpaceASCII)
	prefix = line[:len(line)-len(path)]
	switch path[0] {
	case '<':
		sys = true
		path, suffix, _ = bytes.Cut(path[1:], []byte(">"))
	case '"':
		path, suffix, _ = bytes.Cut(path[1:], []byte(`"`))
	default:
		return line, nil, nil, false, false
	}
	return prefix, path, suffix, sys, true
}

// formatInclude is the inverse of cutInclude.
func formatInclude(prefix, path, suffix []byte, sys bool) ([]byte, bool) {
	b := make([]byte, 0, len(prefix)+len(path)+len(suffix)+2)
	b = append(b, prefix...)
	if sys {
		b = append(b, '<')
	} else {
		b = append(b, '"')
	}
	b = append(b, path...)
	if sys {
		b = append(b, '>')
	} else {
		b = append(b, '"')
	}
	b = append(b, suffix...)
	if a, b, c, d, ok := cutInclude(b); !ok || !bytes.Equal(a, prefix) || !bytes.Equal(b, path) || !bytes.Equal(c, suffix) || d != sys {
		return b, false
	}
	return b, true
}

// resolveInclude attempts to resolve an include against the current directory
// and specified search path.
func resolveInclude(files map[string][]byte, name, inc string, search []string) (string, error) {
	name = path.Clean(name)
	if res := path.Clean(path.Join(path.Dir(name), inc)); res != name {
		if _, ok := files[res]; ok {
			return res, nil
		}
	}
	for _, dir := range search {
		if res := path.Clean(path.Join(dir, inc)); res != name {
			if _, ok := files[res]; ok {
				return res, nil
			}
		}
	}
	return "", fmt.Errorf("include %q does not exist", inc)
}

// isSpaceASCII returns true if c is an ASCII whitespace character.
func isSpaceASCII(c rune) bool {
	return c == '\n' || c == '\t' || c == '\r' || c == '\v' || c == '\f' || c == ' '
}

// bytesTryReplaceAll is like [bytes.ReplaceAll], but returns the number of
// replacements made.
func bytesTryReplaceAll(s, old, new []byte) ([]byte, int) {
	n := bytes.Count(s, old)
	s = bytes.ReplaceAll(s, old, new)
	return s, n
}
