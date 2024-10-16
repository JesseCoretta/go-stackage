package stackage

import (
	"testing"
	"unsafe"
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
	if T, V, _ := derefPtr(typOf(ptr), valOf(ptr)); V.IsZero() || T.Kind() == 0x0 {
		t.Errorf("%s failed; pointer deref unsuccessful", t.Name())
		return
	}
}

func TestMiscCodecov(t *testing.T) {
	//for codecov
	sliceOrArrayKind()
	valueIsValid()
	getStringer(nil)
	isKnownPrimitive()

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

func TestPrimitiveEqualityAssertions(t *testing.T) {
	prims := []any{
		int(3),
		true,
		struct{}{},
		uint64(818),
		nil,
	}

	var prim string = `primitive`
	for i := 0; i < len(prims); i++ {
		same := valOf(prim)
		other := valOf(prims[i])
		if _, err := primitivesEqual(same, other); err == nil {
			t.Errorf("%s failed: expected error, got nothing", t.Name())
			return
		}
	}

}

func TestChannelEqualityAssertions(t *testing.T) {
	same := make(chan error, 5)
	if err := channelsEqual(same, same); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	var nochan chan error
	chans := []any{
		same,
		make(chan bool, 1),
		make(chan struct{}, 1),
		nil,
		[]string{`crap`},
		nochan,
	}

	for i := 0; i < len(chans); i++ {
		for j := 0; j < len(chans); j++ {
			if i == j {
				// don't compare anything to itself, as
				// we're trying to trigger errors with
				// bogus input.
				continue
			}
			f1 := chans[i]
			f2 := chans[j]
			if err := channelsEqual(f1, f2); err == nil {
				t.Errorf("%s failed: expected error, got nothing", t.Name())
			}
		}
	}
}

func TestSliceEqualityAssertions(t *testing.T) {

	a1 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	a2 := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if err := slicesEqual(a1, a2); err != nil {
		// should have worked; we treat
		// these two the same ...
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	multi := []any{
		[]string{`a`, `b`, `c`, `d`, `e`},
		[]any{struct{}{}, struct{ Name string }{}},
		[]any{struct{}{}, struct{}{}},
		struct{}{},
		&struct{}{},
		[10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]bool{},
		`this`,
		[]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'},
		[]any{`a`, 1, `b`, 2, `c`, 3, `d`, 4},
		[]any{struct{}{}, true, 3e14, 3.14159, `hello`, 'A'},
		true,
	}

	for i := 0; i < len(multi); i++ {
		for j := 0; j < len(multi); j++ {
			if i == j {
				// don't compare anything to itself, as
				// we're trying to trigger errors with
				// bogus input.
				continue
			}
			f1 := multi[i]
			f2 := multi[j]
			if err := slicesEqual(f1, f2); err == nil {
				t.Errorf("%s failed: expected error, got nothing", t.Name())
			}
		}
	}
}

func TestFunctionEqualityAssertions(t *testing.T) {
	same := TestFunctionEqualityAssertions // any func will do.
	if err := functionsEqual(same, same); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	funks := []any{
		func() string { return `` },
		struct{}{},
		func(...any) error { return nil },
		float64(1111.3),
		func() {},
		func(string, ...any) (int, error) { return 0, nil },
	}

	for i := 0; i < len(funks); i++ {
		for j := 0; j < len(funks); j++ {
			if i == j {
				// don't compare anything to itself, as
				// we're trying to trigger errors with
				// bogus input.
				continue
			}
			f1 := funks[i]
			f2 := funks[j]
			if err := functionsEqual(f1, f2); err == nil {
				t.Errorf("%s failed: expected error, got nothing", t.Name())
			}
		}
	}
}

func TestExtraEqualityAssertions(t *testing.T) {
	var us1 int = 42
	var us2 string = `forty two`
	var us3 string = ``
	var us4 any = struct{}{}
	var us5 rune = rune(44)
	var us6 byte = 0x1
	var us7 uint32 = uint32(444)
	var us8 map[string][]string = map[string][]string{`test`: {`1`, `2`}}
	var us9 any = &struct{}{}
	var us10 interface{}
	uns := []any{
		func() string { return `` },
		float64(1111.3),
		func() {},
		us8,
		uintptr(unsafe.Pointer(&us1)),
		unsafe.Pointer(&us2),
		func(...any) error { return nil },
		uintptr(unsafe.Pointer(&us3)),
		unsafe.Pointer(&us4),
		us4,
		func(string, ...any) (int, error) { return 0, nil },
		uintptr(unsafe.Pointer(&us5)),
		nil,
		unsafe.Pointer(&us6),
		uintptr(unsafe.Pointer(&us7)),
		struct{}{},
		unsafe.Pointer(&us8),
		uintptr(unsafe.Pointer(&us9)),
		us9,
		unsafe.Pointer(&us10),
		us10,
	}

	// non-matching instance checks (all should fail)
	for i := 0; i < len(uns); i++ {
		for j := 0; j < len(uns); j++ {
			if i == j {
				// don't compare anything to itself.
				continue
			}
			vi := valOf(uns[i])
			vj := valOf(uns[j])
			ki := vi.Kind()
			kj := vj.Kind()

			if err := matchExtra(ki, kj, vi, vj, uns[i], uns[j]); err == nil {
				t.Errorf("%s failed: expected error, got nothing", t.Name())
			}
		}
	}

	// same/same instance checks (should pass)
	for _, idx := range []int{0, 2, 5} {
		vi := valOf(uns[idx])
		vj := valOf(uns[idx])
		ki := vi.Kind()
		kj := vj.Kind()

		if err := matchExtra(ki, kj, vi, vj, uns[idx], uns[idx]); err != nil {
			t.Errorf("%s failed: %v", t.Name(), err)
			return
		}
	}
}

func TestStructEqualityAssertions(t *testing.T) {
	s := []any{
		struct {
			Handle string
			Age    uint16
			Height uint8
		}{Handle: "Courtney", Age: uint16(48), Height: uint8(168)},
		struct{ Name string }{Name: "Courtney"},
		rune(87),
		struct{ string }{`hello`},
		struct {
			Name string
			Age  uint8
		}{Name: "Courtney", Age: uint8(48)},
		struct{ Name string }{},
		[]byte(`HELP ME`),
		struct {
			Nombre string
			Age    uint16
		}{Nombre: "Courtney", Age: uint16(48)},
		struct {
			Name   string
			Age    uint16
			Height uint16
		}{Name: "Courtney", Height: uint16(168)},
		struct {
			Yo     string
			Config *struct{}
		}{Yo: "Courtney", Config: nil},
	}

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(s); j++ {
			if i == j {
				// don't compare anything to itself.
				continue
			}

			if err := structsEqual(s[i], s[j]); err == nil {
				t.Errorf("%s failed: expected error, got nothing", t.Name())
			}
		}
	}
}

func TestMapEqualityAssertions(t *testing.T) {
	same := map[bool]float32{true: float32(3.14159)}
	if err := mapsEqual(same, same); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	l := []any{
		map[bool]float32{true: float32(3.14159)},
		map[string]int{`key3`: 3, `key4`: 4},
		map[string]int{},
		map[string]int{`key3`: 3, `key4`: 4, `key87`: 87},
		map[any]any{`key3`: 3, uint64(4): 4, true: []int64{1, 2, 3, 4, 5, 6}},
		map[string]any{`key3`: 3, `key4`: int64(4), `key99`: []int{1, 2, 3, 4, 5, 6}},
		map[string]any{`key3`: 3, `key4`: 4, `key99`: []int{1, 2, 3, 4, 5}},
		[]any{1, 2, 3, 4},
		map[any]any{`key1`: true, float64(3.14159): map[string]bool{`this`: true}, 7: &struct{}{}, true: 1e45},
		nil,
		`string`,
		rune(77),
		map[string]string{`key5`: `5`, `key6`: `6`, `key7`: `7`},
		&struct{}{},
	}

	for i := 0; i < len(l); i++ {
		for j := 0; j < len(l); j++ {
			if i == j {
				// don't compare anything to itself.
				continue
			}

			if err := mapsEqual(l[i], l[j]); err == nil {
				t.Errorf("%s failed: expected error, got nothing", t.Name())
			}
		}
	}
}
