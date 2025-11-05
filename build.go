package marisa

import (
	"cmp"
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
	// NumTries specifies the number of tries to use. Usually, more tries make a
	// dictionary space-efficient but time-inefficient.
	NumTries int

	// CacheLevel specifies the cache size. A larger cache enables faster search
	// but takes a more space.
	CacheLevel CacheLevel

	// TailMode specifies the kind of TAIL implementation.
	TailMode TailMode

	// NodeOrder specifies the arrangement of nodes, which affects the time cost
	// of matching and the order of predictive search.
	NodeOrder NodeOrder
}

const (
	MinNumTries = 1
	MaxNumTries = 127
)

type CacheLevel int

const (
	DefaultCache CacheLevel = iota
	HugeCache
	LargeCache
	NormalCache
	SmallCache
	TinyCache
)

type TailMode int

const (
	DefaultTail TailMode = iota

	// TextTail merges last labels as zero-terminated strings. So, it is
	// available if and only if the last labels do not contain a NULL character.
	// If TextTail is specified and a NULL character exists in the last labels,
	// the setting is automatically switched to MARISA_BINARY_TAIL.
	TextTail

	// BinaryTail also merges last labels but as byte sequences. It uses a bit
	// vector to detect the end of a sequence, instead of NULL characters. So,
	// BinaryTail requires a larger space if the average length of labels is
	// greater than 8.
	BinaryTail
)

type NodeOrder int

const (
	DefaultOrder NodeOrder = iota

	// LabelOrder arranges nodes in ascending label order. LabelOrder is useful
	// if an application needs to predict keys in label order.
	LabelOrder

	// WeightOrder arranges nodes in descending weight order. WeightOrder is
	// generally a better choice because it enables faster matching.
	WeightOrder
)

func numTriesFlag(v int) (uint32, bool) {
	if v = cmp.Or(v, 3); MinNumTries <= v && v <= MaxNumTries {
		return uint32(v), true
	}
	return 0, false
}

func cacheLevelFlag(v CacheLevel) (uint32, bool) {
	switch cmp.Or(v, NormalCache) {
	case HugeCache:
		return 0x00080, true
	case LargeCache:
		return 0x00100, true
	case NormalCache:
		return 0x00200, true
	case SmallCache:
		return 0x00400, true
	case TinyCache:
		return 0x00800, true
	default:
		return 0, false
	}
}

func tailModeFlag(v TailMode) (uint32, bool) {
	switch cmp.Or(v, TextTail) {
	case TextTail:
		return 0x01000, true
	case BinaryTail:
		return 0x02000, true
	default:
		return 0, false
	}
}

func nodeOrderFlag(v NodeOrder) (uint32, bool) {
	switch cmp.Or(v, WeightOrder) {
	case LabelOrder:
		return 0x10000, true
	case WeightOrder:
		return 0x20000, true
	default:
		return 0, false
	}
}

func configFlags(c Config) (flags uint32, ok bool) {
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
