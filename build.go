package marisa

import (
	"errors"
	"iter"
	"math"

	"github.com/pgaskin/go-marisa/internal/walloc"
)

// Config specifies options for a dictionary. Any unspecified options will be
// set to their default.
type Config struct {
	NumTries   int
	CacheLevel CacheLevel
	TailMode   TailMode
	NodeOrder  NodeOrder
}

func configFlags(c Config) (flags configFlag, ok bool) {
	if f, ok := numTriesFlag(c.NumTries); ok {
		flags |= f
	} else {
		return 0, false
	}
	if f, ok := cacheLevelFlag(c.CacheLevel); ok {
		flags |= f
	} else {
		return 0, false
	}
	if f, ok := tailModeFlag(c.TailMode); ok {
		flags |= f
	} else {
		return 0, false
	}
	if f, ok := nodeOrderFlag(c.NodeOrder); ok {
		flags |= f
	} else {
		return 0, false
	}
	return flags, true
}

// Build builds a dictionary out of the specified set of keys, with a weight of
// 1 for each.
func (t *Trie) Build(keys iter.Seq[string], cfg Config) error {
	return t.BuildWeights(func(yield func(string, float32) bool) {
		for key := range keys {
			if !yield(key, 1.0) {
				return
			}
		}
	}, cfg)
}

// BuildWeights builds a dictionary out of the specified set of keys and
// weights. If a key is specified multiple times, the weights are accumulated.
func (t *Trie) BuildWeights(keys iter.Seq2[string, float32], cfg Config) error {
	flag, ok := configFlags(cfg)
	if !ok {
		return errors.New("invalid config")
	}

	sa := &walloc.SliceAllocator{
		OverrideMax: maxAlloc,
	}
	mod, err := instantiate(sa)
	if err != nil {
		return err
	}

	const step = 1024
	alloc := step
	ptr, err := mod.Alloc(step)
	if err != nil {
		return err
	}
	defer func() {
		if ptr != 0 {
			mod.Free(ptr)
		}
	}()
	for key, weight := range keys {
		n := len(key)
		if n > alloc {
			old := ptr
			ptr = 0
			mod.Free(old)
			alloc = (n + step - 1) / step * step
			ptr, err = mod.Alloc(alloc)
			if err != nil {
				return err
			}
		}
		buf, ok := mod.Module().Memory().Read(ptr, uint32(len(key)))
		if !ok {
			panic("bad allocation")
		}
		copy(buf, key)

		if _, err := mod.Call("marisa_build_push", uint64(ptr), uint64(n), uint64(math.Float32bits(weight))); err != nil {
			return err
		}
	}
	if _, err := mod.Call("marisa_build", uint64(flag)); err != nil {
		return err
	}
	return t.swap(mod)
}
