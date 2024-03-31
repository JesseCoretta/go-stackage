package stackage

import (
	"testing"
)

var strInSliceMap map[int]map[int][]bool = map[int]map[int][]bool{
	// case match
	0: {
		0: {true, true, true, true, true},
		1: {true, true, true, true, true},
	},

	// case fold
	1: {
		0: {true, true, true, true, true},
		1: {true, true, true, true, true},
	},
}

func TestStrInSlice(t *testing.T) {
	for idx, fn := range []func(string, []string) bool{
		strInSlice,
		strInSliceFold,
	} {
		for i, values := range [][]string{
			{`cAndidate1`, `blarGetty`, `CANndidate7`, `squatcobbler`, `<censored>`},
			{`Ó-aîï4Åø´øH«w%);<wÃ¯`, `piles`, `4378295fmitty`, string(rune(0)), `broccolI`},
		} {
			for j, val := range values {
				result_expected := strInSliceMap[idx][i][j]

				// warp the candidate value such that
				// it no longer matches the slice from
				// whence it originates. j² is used as
				// its quicker and less stupid than
				// adding a rand generator.
				if isPowerOfTwo(j) {
					var R []rune = []rune(val)
					for g, h := 0, len(R)-1; g < h; g, h = g+1, h-1 {
						R[g], R[h] = R[h], R[g]
					}
					val = string(R)
					result_expected = !result_expected // invert
				}

				result_received := fn(val, values)
				if result_expected != result_received {
					t.Errorf("%s[%d->%d] failed; []byte(%v) in %v: %t (wanted %t)",
						t.Name(), i, j, []byte(val), values, result_received, result_expected)
					return
				}
			}
		}
	}
}

func TestDeref(t *testing.T) {
	c := Cond(`this`, Eq, `that`)
	ptr := &c
	if T, V := derefPtr(ptr); V.IsZero() || T.Kind() == 0x0 {
		t.Errorf("%s failed; pointer deref unsuccessful", t.Name())
		return
	}
}

func TestMiscCodecov(t *testing.T) {
	for idx, val := range []any{
		uint(1),
		int(1),
		uint8(1),
		int8(1),
		uint16(1),
		int16(1),
		uint32(1),
		int32(1),
		uint64(1),
		int64(1),
	} {
		if idx%2 == 0 {
			_ = uintStringer(val)
		} else {
			_ = intStringer(val)
		}
	}

	for _, val := range []any{
		complex64(1),
		complex128(1),
	} {
		_ = complexStringer(val)
	}

	var ls logSystem
	_ = ls.positive(1)

	ls = logSystem{
		log: cLogDefault,
	}
	_ = ls.positive(1)
	logDiscard(nil)
	logDiscard(sLogDefault)

	_ = timestamp()
}
