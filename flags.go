package marisa

import (
	"cmp"
)

type configFlag uint32

// These are the values from marisa/base.h. They are not exported since they are
// not part of the Go API.
const (
	_MARISA_MIN_NUM_TRIES     configFlag = 0x00001
	_MARISA_MAX_NUM_TRIES     configFlag = 0x0007F
	_MARISA_DEFAULT_NUM_TRIES configFlag = 0x00003

	_MARISA_HUGE_CACHE    configFlag = 0x00080
	_MARISA_LARGE_CACHE   configFlag = 0x00100
	_MARISA_NORMAL_CACHE  configFlag = 0x00200
	_MARISA_SMALL_CACHE   configFlag = 0x00400
	_MARISA_TINY_CACHE    configFlag = 0x00800
	_MARISA_DEFAULT_CACHE configFlag = _MARISA_NORMAL_CACHE

	_MARISA_TEXT_TAIL    configFlag = 0x01000
	_MARISA_BINARY_TAIL  configFlag = 0x02000
	_MARISA_DEFAULT_TAIL configFlag = _MARISA_TEXT_TAIL

	_MARISA_LABEL_ORDER   configFlag = 0x10000
	_MARISA_WEIGHT_ORDER  configFlag = 0x20000
	_MARISA_DEFAULT_ORDER configFlag = _MARISA_WEIGHT_ORDER

	_MARISA_NUM_TRIES_MASK   configFlag = 0x0007F
	_MARISA_CACHE_LEVEL_MASK configFlag = 0x00F80
	_MARISA_TAIL_MODE_MASK   configFlag = 0x0F000
	_MARISA_NODE_ORDER_MASK  configFlag = 0xF0000
	_MARISA_CONFIG_MASK      configFlag = 0xFFFFF
)

// NumTries specifies the number of tries to use. Usually, more tries make a
// dictionary space-efficient but time-inefficient.
//
// The limits are an implementation detail.
const (
	MinNumTries = int(_MARISA_MIN_NUM_TRIES)
	MaxNumTries = int(_MARISA_MAX_NUM_TRIES)
)

// CacheLevel specifies the cache size. A larger cache enables faster search but
// takes a more space.
type CacheLevel int

const (
	HugeCache CacheLevel = iota + 1
	LargeCache
	NormalCache
	SmallCache
	TinyCache
)

// TailMode specifies the kind of TAIL implementation.
type TailMode int

const (
	// TextTail merges last labels as zero-terminated strings. So, it is
	// available if and only if the last labels do not contain a NULL character.
	// If TextTail is specified and a NULL character exists in the last labels,
	// the setting is automatically switched to MARISA_BINARY_TAIL.
	TextTail TailMode = iota + 1

	// BinaryTail also merges last labels but as byte sequences. It uses a bit
	// vector to detect the end of a sequence, instead of NULL characters. So,
	// BinaryTail requires a larger space if the average length of labels is
	// greater than 8.
	BinaryTail
)

// NodeOrder specifies the arrangement of nodes, which affects the time cost of
// matching and the order of predictive search.
type NodeOrder int

const (
	// LabelOrder arranges nodes in ascending label order. LabelOrder is useful
	// if an application needs to predict keys in label order.
	LabelOrder NodeOrder = iota + 1

	// WeightOrder arranges nodes in descending weight order. WeightOrder is
	// generally a better choice because it enables faster matching.
	WeightOrder
)

func (c CacheLevel) String() string {
	switch c {
	case HugeCache:
		return "huge"
	case LargeCache:
		return "large"
	case NormalCache:
		return "normal"
	case SmallCache:
		return "small"
	case TinyCache:
		return "tiny"
	default:
		return "unknown"
	}
}

func (t TailMode) String() string {
	switch t {
	case TextTail:
		return "text"
	case BinaryTail:
		return "binary"
	default:
		return "unknown"
	}
}

func (t NodeOrder) String() string {
	switch t {
	case LabelOrder:
		return "label"
	case WeightOrder:
		return "weight"
	default:
		return "unknown"
	}
}

func numTriesFlag(v int) (configFlag, bool) {
	if v == 0 {
		return _MARISA_DEFAULT_NUM_TRIES, true
	}
	if MinNumTries <= v && v <= MaxNumTries {
		return configFlag(v), true
	}
	return 0, false
}

func cacheLevelFlag(v CacheLevel) (configFlag, bool) {
	if v == 0 {
		return _MARISA_DEFAULT_CACHE, true
	}
	switch v {
	case HugeCache:
		return _MARISA_HUGE_CACHE, true
	case LargeCache:
		return _MARISA_LARGE_CACHE, true
	case NormalCache:
		return _MARISA_NORMAL_CACHE, true
	case SmallCache:
		return _MARISA_SMALL_CACHE, true
	case TinyCache:
		return _MARISA_TINY_CACHE, true
	default:
		return 0, false
	}
}

func tailModeFlag(v TailMode) (configFlag, bool) {
	if v == 0 {
		return _MARISA_DEFAULT_TAIL, true
	}
	switch v {
	case TextTail:
		return _MARISA_TEXT_TAIL, true
	case BinaryTail:
		return _MARISA_BINARY_TAIL, true
	default:
		return 0, false
	}
}

func nodeOrderFlag(v NodeOrder) (configFlag, bool) {
	if v == 0 {
		return _MARISA_DEFAULT_ORDER, true
	}
	switch cmp.Or(v, WeightOrder) {
	case LabelOrder:
		return _MARISA_LABEL_ORDER, true
	case WeightOrder:
		return _MARISA_WEIGHT_ORDER, true
	default:
		return 0, false
	}
}

func flagNumTries(v configFlag) int {
	return int(v & _MARISA_NUM_TRIES_MASK)
}

func flagCacheLevel(f configFlag) CacheLevel {
	switch f & _MARISA_CACHE_LEVEL_MASK {
	case _MARISA_HUGE_CACHE:
		return HugeCache
	case _MARISA_LARGE_CACHE:
		return LargeCache
	case _MARISA_NORMAL_CACHE:
		return NormalCache
	case _MARISA_SMALL_CACHE:
		return SmallCache
	case _MARISA_TINY_CACHE:
		return TinyCache
	default:
		return 0
	}
}

func flagTailMode(f configFlag) TailMode {
	switch f & _MARISA_CONFIG_MASK {
	case _MARISA_TEXT_TAIL:
		return TextTail
	case _MARISA_BINARY_TAIL:
		return BinaryTail
	default:
		return 0
	}
}

func flagNodeOrder(f configFlag) NodeOrder {
	switch f & _MARISA_NODE_ORDER_MASK {
	case _MARISA_LABEL_ORDER:
		return LabelOrder
	case _MARISA_WEIGHT_ORDER:
		return WeightOrder
	default:
		return 0
	}
}
