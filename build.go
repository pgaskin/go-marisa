package marisa

import (
	"encoding/binary"
	"errors"
	"iter"
	"math"

	"github.com/pgaskin/go-marisa/internal"
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

const chunkSize = 4 * 1024 * 1024

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

	var (
		free []uint32
		cptr uint32
		cbuf []byte
		cnum uint32
	)
	for key, weight := range keys {
		if csz := 8 + len(key); !internal.NoChunkBuild && csz < chunkSize {
			if cptr != 0 && csz > len(cbuf) {
				if _, err := mod.Call("marisa_build_push_chunk", uint64(cptr), uint64(cnum)); err != nil {
					return err
				}
				cptr, cbuf = 0, nil
			}
			if cptr == 0 {
				cptr, cbuf, err = mod.Alloc(chunkSize)
				if err != nil {
					return err
				}
				cnum = 0
				free = append(free, cptr)
			}
			var tmp []byte
			tmp, cbuf = cbuf[:csz], cbuf[csz:]
			binary.LittleEndian.PutUint32(tmp[0:4], uint32(len(key)))
			binary.LittleEndian.PutUint32(tmp[4:8], math.Float32bits(weight))
			copy(tmp[8:], key)
			cnum++
		} else {
			ptr, buf, err := mod.Alloc(len(key))
			if err != nil {
				return err
			}
			free = append(free, ptr)
			copy(buf, key)

			if _, err := mod.Call("marisa_build_push", uint64(ptr), uint64(len(key)), uint64(math.Float32bits(weight))); err != nil {
				return err
			}
		}
	}
	if !internal.NoChunkBuild && cptr != 0 {
		if _, err = mod.Call("marisa_build_push_chunk", uint64(cptr), uint64(cnum)); err != nil {
			return err
		}
		cptr, cbuf = 0, nil
	}
	if _, err := mod.Call("marisa_build", uint64(flag)); err != nil {
		return err
	}
	for _, ptr := range free {
		mod.Free(ptr)
	}
	return t.swap(mod)
}
