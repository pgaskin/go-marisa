package marisa

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestFlags(t *testing.T) {
	testFlag(t, "NumTries", _MARISA_NUM_TRIES_MASK, _MARISA_DEFAULT_NUM_TRIES, numTriesFlag, flagNumTries)
	testFlag(t, "CacheLevel", _MARISA_CACHE_LEVEL_MASK, _MARISA_DEFAULT_CACHE, cacheLevelFlag, flagCacheLevel)
	testFlag(t, "TailMode", _MARISA_TAIL_MODE_MASK, _MARISA_DEFAULT_TAIL, tailModeFlag, flagTailMode)
	testFlag(t, "NodeOrder", _MARISA_NODE_ORDER_MASK, _MARISA_DEFAULT_ORDER, nodeOrderFlag, flagNodeOrder)
}

// testFlag ensures that a flag:
//   - round-trips from a go enum value to a marisa flag value and back
//   - has distinct go enum values and marisa flag value
//   - has a valid default value
//   - has marisa flag values which all are in the corresponding mask
//   - if an enum, has a String method with unique string values for valid flags, and unknown otherwise
//   - if an enum, values are numbered from 1 contiguously
func testFlag[T ~int](t *testing.T, name string,
	mask, def configFlag,
	toMarisa func(T) (configFlag, bool),
	toGo func(configFlag) T,
) {
	t.Run(name, func(t *testing.T) {
		var (
			zero   T
			isEnum = reflect.TypeFor[T]().PkgPath() != ""
		)
		if mask&^_MARISA_CONFIG_MASK != 0 {
			t.Fatalf("marisa flag mask is not within the config mask")
		}
		if def&^mask != 0 {
			t.Fatalf("marisa flag default is not within the mask")
		}
		if x, ok := toMarisa(zero); !ok {
			t.Errorf("zero go flag should be valid")
		} else if x != def {
			t.Errorf("zero go flag should convert to the default marisa flag, got %v", x)
		}
		if x := toGo(0); x != 0 {
			t.Errorf("zero marisa flag should convert to zero go flag")
		}
		for x := range _MARISA_CONFIG_MASK {
			if y, ok := toMarisa(toGo(x)); !ok {
				t.Errorf("marisa(go(%#x)) round-trip failed", x)
			} else if x == y {
				// valid value round-tripped correctly
			} else {
				// part of it got masked out, so it's an invalid value
			}
		}
		var upper T
		if isEnum {
			for x := range math.MaxInt {
				if _, ok := toMarisa(T(x)); !ok {
					break
				}
				upper++
			}
		}
		seen := make([]bool, upper)
		for x := range _MARISA_CONFIG_MASK {
			y, ok := toMarisa(toGo(x))
			if !ok {
				t.Errorf("marisa(go(%x)) round-trip failed", x)
			}
			if y&^mask != 0 {
				t.Errorf("marisa(go(%x)) is not within the mask", x)
			}
			if x != y {
				continue // part of it got masked out, so it's an invalid value
			}
			if isEnum {
				z := toGo(x)
				if z == 0 {
					t.Errorf("valid marisa flag %x should not be a zero go flag", x)
				}
				if z >= T(upper) {
					t.Fatalf("valid marisa flag %x should have a go flag in the range (0, %d)", x, upper)
				}
				seen[z] = true
			}
		}
		if isEnum {
			if _, ok := any(zero).(fmt.Stringer); !ok {
				t.Errorf("enum go flag does not implement fmt.Stringer")
			}
			if any(zero).(fmt.Stringer).String() != "unknown" {
				t.Errorf("zero go enum flag string should be unknown")
			}
			seenString := make(map[string]T, upper)
			for x := range upper {
				if x == 0 {
					continue
				}
				if !seen[x] {
					t.Errorf("no marisa flag resulted in go enum value %d", x)
				}
				if s := any(T(x)).(fmt.Stringer).String(); s == "unknown" || s == "" {
					t.Errorf("nonzero go enum value %d string should not be unknown or empty", x)
				} else {
					if y, ok := seenString[s]; ok {
						t.Errorf("duplicate string for go enum value %d and %d", x, y)
					}
					seenString[s] = x
				}
			}
			for x := range upper {
				y, _ := toMarisa(x)
				z := any(toGo(y)).(fmt.Stringer).String()
				t.Logf("%s(%d) %#x (%s)", reflect.TypeFor[T]().Name(), x, y, z)
			}
		}
	})
}
