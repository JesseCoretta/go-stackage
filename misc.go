package stackage

import (
	"errors"
	"fmt"
	"log"
	"math/rand" // not for crypto, don't worry :)
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/*
frequently-used package import function aliases.
*/
var (
	typOf   func(any) reflect.Type              = reflect.TypeOf
	valOf   func(any) reflect.Value             = reflect.ValueOf
	printf  func(string, ...any) (int, error)   = fmt.Printf
	sprintf func(string, ...any) string         = fmt.Sprintf
	eq      func(string, string) bool           = strings.EqualFold
	lc      func(string) string                 = strings.ToLower
	ilc     func(rune) bool                     = unicode.IsLower
	uc      func(string) string                 = strings.ToUpper
	iuc     func(rune) bool                     = unicode.IsUpper
	rplc    func(string, string, string) string = strings.ReplaceAll
	qt      func(string) string                 = strconv.Quote
	uq      func(string) (string, error)        = strconv.Unquote
	itoa    func(int) string                    = strconv.Itoa
	split   func(string, string) []string       = strings.Split
	trimS   func(string) string                 = strings.TrimSpace
	join    func([]string, string) string       = strings.Join
	scmp    func(string, string) int            = strings.Compare
	now     func() time.Time                    = time.Now
)

const (
	randChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randIDSize = 24
)

func newStringBuilder() strings.Builder {
	return strings.Builder{}
}

func bool2str(b bool) string {
	if b {
		return `true`
	}
	return `false`
}

func randomID(n int) string {
	id := make([]byte, n)
	for i := range id {
		id[i] = randChars[rand.Int63()%int64(len(randChars))]
	}
	return string(id)
}

func timestamp() string {
	t := now()

	return sprintf("%04d%02d%02d%02d%02d%02d.%09d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond())
}

/*
errorf wraps errors.New and returns a non-nil instance of error
based upon a non-nil/non-zero msg input value with optional args.
*/
func errorf(msg any, x ...any) (err error) {
	switch tv := msg.(type) {
	case string:
		if len(tv) > 0 {
			err = errors.New(tv)
		}
	case error:
		if tv != nil {
			err = errors.New(tv.Error())
		}
	}

	return
}

func isPowerOfTwo(x int) bool {
	return x&(x-1) == 0
}

/*
strInSlice returns a Boolean value indicative of whether the
specified string (str) is present within slice. Please note
that case is a significant element in the matching process.
*/
func strInSlice(str string, slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if str == slice[i] {
			return true
		}
	}
	return false
}

/*
strInSliceFold returns a Boolean value indicative of whether
the specified string (str) is present within slice. Case is
not significant in the matching process.
*/
func strInSliceFold(str string, slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if eq(str, slice[i]) {
			return true
		}
	}
	return false
}

/*
isPtr returns a Boolean value indicative of whether kind
reflection revealed the presence of a pointer type.
*/
func isPtr(t reflect.Type) bool {
	if t == nil {
		return false
	}

	return t.Kind() == reflect.Ptr
}

func derefPtr(t reflect.Type, v reflect.Value) (reflect.Type, reflect.Value, reflect.Kind) {
	// loop to handle **type instances
	var k reflect.Kind
	for {
		if isPtr(t) {
			t = t.Elem()
			v = v.Elem()
			continue
		}
		break
	}
	k = v.Kind()

	return t, v, k
}

func assertReflect(x any) (at reflect.Type, av reflect.Value) {
	switch tv := x.(type) {
	case reflect.Value:
		at = tv.Type()
		av = tv
	default:
		at = typOf(tv)
		av = valOf(tv)
	}

	return
}

func sliceOrArrayKind(k ...reflect.Kind) bool {
	if len(k) == 0 {
		return false
	}

	for i := 0; i < len(k); i++ {
		if k[i] != reflect.Slice && k[i] != reflect.Array {
			return false
		}
	}

	return true
}

func valueIsValid(v ...reflect.Value) bool {
	if len(v) == 0 {
		return false
	}

	for i := 0; i < len(v); i++ {
		if !v[i].IsValid() {
			return false
		}
	}

	return true
}

func channelsEqual(x, y any) error {

	xft, xfv := assertReflect(x)
	yft, yfv := assertReflect(y)

	if !valueIsValid(xfv, yfv) {
		return errorf("Channel(s) invalid")
	}

	xfk := xft.Kind()
	yfk := yft.Kind()

	if xfk != reflect.Chan || yfk != reflect.Chan {
		return errorf("Channel kind mismatch")
	}

	if xft != yft {
		return errorf("Channel type mismatch")
	}

	if x != y {
		return errorf("Channel mismatch")
	}

	return nil
}

func functionsEqual(x, y any) error {

	if x == nil || y == nil {
		return errorf("Nil functions incomparable")
	}

	xft, _ := assertReflect(x)
	yft, _ := assertReflect(y)

	xfk := xft.Kind()
	yfk := yft.Kind()

	if xfk != reflect.Func || yfk != reflect.Func {
		return errorf("Function kind mismatch")
	}

	if xft != yft {
		return errorf("Function type mismatch")
	}

	// Try to match by signature elements...
	return nil
}

func primitivesEqual(x, y reflect.Value) (tried bool, err error) {
	if !valueIsValid(x, y) {
		err = errorf("Nil input")
		return
	}

	if isKnownPrimitive(x.Interface()) {
		tried = true
		if isKnownPrimitive(y.Interface()) {
			if !x.Equal(y) {
				err = errorf("primitive mismatch")
			}
			return
		}
		err = errorf("primitive incomparable to non-primitive")
	}

	return
}

func valuesEqual(x, y any) error {

	if x == nil && y == nil {
		return nil
	}

	_, xrv, xrk := derefPtr(assertReflect(x))
	_, yrv, yrk := derefPtr(assertReflect(y))

	if tried, err := primitivesEqual(xrv, yrv); tried {
		return err
	}

	switch xrk {
	case reflect.Struct:
		if tried, err := stackageStructsEqual(x, y); tried {
			return err
		}

		// whatever they are, handle manually
		return structsEqual(x, y)
	case reflect.Slice, reflect.Array:
		return slicesEqual(x, y)
	case reflect.Map:
		return mapsEqual(x, y)
	default:
		return matchExtra(xrk, yrk, xrv, yrv, x, y)
	}
}

func matchExtra(l, k reflect.Kind, a, b reflect.Value, x, y any) error {

	switch k {
	case reflect.Func:
		return functionsEqual(x, y)
	case reflect.Chan:
		return channelsEqual(x, y)
	case reflect.UnsafePointer, reflect.Uintptr:
		return uuptrsEqual(l, k, a, b)
	}

	return errorf("Unsupported type")
}

func uuptrsEqual(l, k reflect.Kind, x, y reflect.Value) error {
	if l == k {
		if x.Interface() != y.Interface() {
			return errorf("UnsafePointer mismatch")
		}
	} else {
		return errorf("Uintptr or unsafepointer kind mismatch")
	}

	return nil
}

func stackageStructsEqual(x, y any) (tried bool, err error) {
	// Are they both condition/condition alias?
	if icd, iokc := conditionTypeAliasConverter(x); iokc {
		tried = true
		if jcd, jokc := conditionTypeAliasConverter(y); jokc {
			err = icd.IsEqual(jcd)
			return
		}
	}

	// Are they both stack/stack alias?
	if ist, ioks := stackTypeAliasConverter(x); ioks {
		tried = true
		if jst, joks := stackTypeAliasConverter(y); joks {
			err = ist.IsEqual(jst)
			return
		}
	}

	err = errorf("Cannot compare stackage instances, cannot convert")

	return
}

func mapsEqual(x, y any) (err error) {

	xrt, xrv, xrk := derefPtr(assertReflect(x))
	yrt, yrv, yrk := derefPtr(assertReflect(y))

	if xrk != reflect.Map || xrk != yrk {
		err = errorf("Cannot compare non-map instances")
		return
	}

	if xrt != yrt {
		err = errorf("Map type mismatch")
		return
	}

	if xrv.Len() != yrv.Len() {
		err = errorf("Map length mismatch")
		return
	}

	for _, key := range xrv.MapKeys() {
		xval := xrv.MapIndex(key).Interface()
		yval := yrv.MapIndex(key).Interface()
		if err = valuesEqual(xval, yval); err != nil {
			return
		}
	}

	return
}

func structsEqual(x, y any) (err error) {
	xrt, xrv, xrk := derefPtr(assertReflect(x))
	yrt, yrv, yrk := derefPtr(assertReflect(y))

	if xrk != yrk {
		err = errorf("Struct type mismatch")
		return
	}

	if xrt.NumField() != yrt.NumField() {
		err = errorf("Struct field number mismatch")
		return
	}

	for i := 0; i < xrt.NumField() && err == nil; i++ {
		xtf := xrt.Field(i)
		xvf := xrv.Field(i)

		ytf := yrt.Field(i)
		yvf := yrv.Field(i)

		xn := xtf.Name
		yn := ytf.Name

		if xn != yn {
			xanon := xtf.Anonymous
			yanon := ytf.Anonymous

			if !(xanon && yanon) {
				err = errorf("Struct anonymous field mismatch failed")
				return
			}
		}

		err = valuesEqual(xvf.Interface(), yvf.Interface())
	}

	return
}

/*
capLenEqual will compare the input c1/c2 as capacity
indicators and l1/l2 as length indicators. If there
are no capacity constraints indicated, a comparison
is made based on length alone.
*/
func capLenEqual(c1, c2, l1, l2 int) bool {
	if c1 != 0 || c2 != 0 {
		return c1 == c2 && l1 == l2
	}
	return l1 == l2
}

func slicesEqual(x, y any) (err error) {

	_, xrv, xrk := derefPtr(assertReflect(x))
	_, yrv, yrk := derefPtr(assertReflect(y))

	if !sliceOrArrayKind(xrk, yrk) {
		err = errorf("Slice/array kind mismatch")
		return
	}

	if !capLenEqual(xrv.Cap(), yrv.Cap(), xrv.Len(), yrv.Len()) {
		err = errorf("Slice/array capacity or length mismatch")
		return
	}

	for i := 0; i < xrv.Len() && err == nil; i++ {
		_, xv, _ := derefPtr(xrv.Index(i).Type(), xrv.Index(i))
		_, yv, _ := derefPtr(yrv.Index(i).Type(), yrv.Index(i))

		// Get primitives out of the way
		var tried bool
		if tried, err = primitivesEqual(xv, yv); tried {
			return
		}

		err = valuesEqual(xv, yv)
	}

	return
}

/*
condenseWHSP returns input string b with all contiguous
WHSP or Horizontal TAB characters condensed into single
WHSP characters (ASCII #32). For example:

	" the   quick	      brown fox  "

... would become ...

	"the quick brown fox"

This function also removes any LEADING or TRAILING WHSP
characters using the strings.TrimSpace alias trimS.

This function, when combined with the act of replacing
all newline (ASCII #10, "\n") characters with a single
space, can help with the conversion of a multi-line or
indented "block value" into a single line value more
cleanly.
*/
func condenseWHSP(b string) string {
	b = trimS(b)

	var last bool // previous char was WHSP or HTAB.
	var builder strings.Builder

	for i := 0; i < len(b); i++ {
		c := rune(b[i])
		switch c {
		case rune(9), rune(32): // match either WHSP or horizontal tab
			if !last {
				last = true
				builder.WriteRune(rune(32)) // Add WHSP
			}
		default: // match all other characters
			if last {
				last = false
			}
			builder.WriteRune(c)
		}
	}

	return builder.String()
}

/*
getStringer uses reflect to obtain and return a given
type instance's String ("stringer") method, if present.
If not, nil is returned.
*/
func getStringer(x any) (meth func() string) {
	if x == nil {
		return nil
	}

	if v := valOf(x); !v.IsZero() {
		if method := v.MethodByName(`String`); method.Kind() != reflect.Invalid {
			if _meth, ok := method.Interface().(func() string); ok {
				meth = _meth
			}
		}
	}
	return
}

/*
encapValue will encapsulate value v using encapsulation scheme
enc, or the original string is returned if no scheme was set.
*/
func encapValue(enc [][]string, v string) string {
	if len(enc) == 0 {
		return v
	}

	for i := len(enc); i > 0; i-- {
		sl := enc[i-1]
		switch len(sl) {
		case 1:
			// use char 0 for both L and R
			v = sl[0] + v + sl[0]
		case 2:
			// char 0 = L, char 1 = R
			v = sl[0] + v + sl[1]
		}
	}

	return v
}

/*
padValue may, or may not, enclose the given string value input
within ASCII #32 (WHSP) "padding" characters. The Boolean "do"
input value will control whether or not padding is actually to
be invoked. Zero string values are ineligible for padding.
*/
func padValue(do bool, value string) string {
	var pad string = string(rune(32)) // whsp char
	if !do {
		// don't use an explicit NTBS value (i.e.:
		// rune(0)), else 0x0 will start showing up
		// in raw byte output. Just use an empty.
		pad = ``
	}

	if len(value) == 0 {
		return ``
	}
	return pad + value + pad
}

/*
foldValue will apply lc (Strings.ToLower) and uc (Strings.ToUpper)
to the value based on the "do" disposition (do, or do not).
*/
func foldValue(do bool, value string) (s string) {
	s = value
	if len(value) > 0 {
		if do {
			s = uc(value) // default to upper
			if iuc(rune(value[0])) {
				s = lc(value) // fold to lower
			}
		}
	}

	return
}

func isNumberPrimitive(x any) bool {
	switch x.(type) {
	case int, int8, int16, int32, int64,
		float32, float64, complex64, complex128,
		uint, uint8, uint16, uint32, uint64:
		return true
	}

	return false
}

func isStringPrimitive(x any) bool {
	switch x.(type) {
	case string:
		return true
	}

	return false
}

func isBoolPrimitive(x any) bool {
	switch x.(type) {
	case bool:
		return true
	}

	return false
}

func isKnownPrimitive(x ...any) (is bool) {
	if len(x) == 0 {
		return false
	}

	for i := 0; i < len(x); i++ {
		if isStringPrimitive(x[i]) {
			is = true
		} else if isNumberPrimitive(x[i]) {
			is = true
		} else if isBoolPrimitive(x[i]) {
			is = true
		} else {
			break
		}
	}

	return
}

func numberStringer(x any) (s string) {
	s = `NaN`
	switch tv := x.(type) {
	case float32, float64:
		s = floatStringer(tv)
	case complex64, complex128:
		s = complexStringer(tv)
	case int, int8, int16, int32, int64:
		s = intStringer(tv)
	case uint, uint8, uint16, uint32, uint64:
		s = uintStringer(tv)
	}

	return
}

func primitiveStringer(x any) (s string) {
	s = `unsupported_primitive_type`
	if isKnownPrimitive(x) {
		switch {
		case isBoolPrimitive(x):
			s = boolStringer(x)
		case isNumberPrimitive(x):
			s = numberStringer(x)
		case isStringPrimitive(x):
			s = x.(string)
		}
	}

	return
}

func boolStringer(x any) string {
	return bool2str(x.(bool))
}

func floatStringer(x any) (s string) {
	switch tv := x.(type) {
	case float32:
		s = strconv.FormatFloat(float64(tv), 'g', -1, 32)
	case float64:
		s = strconv.FormatFloat(tv, 'g', -1, 64)
	}

	return
}

func complexStringer(x any) (s string) {
	switch tv := x.(type) {
	case complex64:
		s = strconv.FormatComplex(complex128(tv), 'g', -1, 64)
	case complex128:
		s = strconv.FormatComplex(tv, 'g', -1, 128)
	}

	return
}

func uintStringer(x any) (s string) {
	switch tv := x.(type) {
	case uint:
		s = strconv.FormatUint(uint64(tv), 10)
	case uint8:
		s = strconv.FormatUint(uint64(tv), 10)
	case uint16:
		s = strconv.FormatUint(uint64(tv), 10)
	case uint32:
		s = strconv.FormatUint(uint64(tv), 10)
	case uint64:
		s = strconv.FormatUint(tv, 10)
	}

	return
}

func intStringer(x any) (s string) {
	switch tv := x.(type) {
	case int:
		s = strconv.FormatInt(int64(tv), 10)
	case int8:
		s = strconv.FormatInt(int64(tv), 10)
	case int16:
		s = strconv.FormatInt(int64(tv), 10)
	case int32:
		s = strconv.FormatInt(int64(tv), 10)
	case int64:
		s = strconv.FormatInt(tv, 10)
	}

	return
}

/*
Interface is an interface type qualified through instances of the
following types (whether native, or type-aliased):

  - [Stack]
  - [Condition]

This interface type offers users an alternative to the tedium of
repeated type assertion for every [Stack] and [Condition] instance they
encounter. This may be particularly useful in situations where the
act of traversal is conducted upon a [Stack] instance that contains
myriad hierarchies and nested contexts of varying types.

This is not a complete "replacement" for the explicit use of package
defined types nor their aliased counterparts. The Interface interface
only extends methods that are read-only in nature AND common to both
of the above categories.

To access the entire breadth of available methods for the underlying
type instance, manual type assertion shall be necessary.

Users SHOULD adopt this interface signature for use in their solutions
if and when needed, though it is not a requirement.
*/
type Interface interface {
	// Stack: Len returns an integer that represents the number of
	// slices present within.
	//
	// Condition: Len returns a one (1) if properly initialized
	// and set with a value that is not a Stack. Len returns a
	// zero (0) if invalid or otherwise not properly initialized.
	// Len returns the value of Stack.Len in deference, should the
	// expression itself be a Stack.
	Len() int

	// IsInit returns a Boolean value indicative of whether the
	// receiver instance is considered 'properly initialized'
	IsInit() bool

	// Stack: IsFIFO returns a Boolean value indicative of whether
	// the Stack instance exhibits First-In-First-Out behavior as it
	// pertains to the ingress and egress order of slices. See the
	// Stack.SetFIFO method for details on setting this behavior.
	//
	// Condition: IsFIFO returns a Boolean value indicative of whether
	// the receiver's Expression value contains a Stack instance AND
	// that Stack exhibits First-In-First-Out behavior as it pertains
	// to the ingress and egress order of slices. If no Stack value is
	// present within the Condition, false is always returned.
	IsFIFO() bool

	// IsZero returns a Boolean value indicative of whether the receiver
	// instance is considered nil, or unset.
	IsZero() bool

	// IsEqual performs an equality check against the receiver
	// instance and the input value.
	IsEqual(any) error

	// IsParen returns a Boolean value indicative of whether the receiver
	// instance, when represented as a string value, shall encapsulate in
	// parenthetical L [(] and R [)] characters.
	//
	// See also the Paren method for Condition and Stack instances.
	IsParen() bool

	// IsEncap returns a Boolean value indicative of whether the receiver
	// instance, when represented as a string value, shall encapsulate the
	// effective value within user-defined quotation characters.
	//
	// See also the Encap method for Condition and Stack instances.
	IsEncap() bool

	// Stack: IsPadded returns a Boolean value indicative of whether WHSP
	// (ASCII #32) padding shall be applied to the outermost ends of the
	// string representation of the receiver (but within parentheticals).
	//
	// Condition: IsPadded returns a Boolean value indicative of whether
	// WHSP (ASCII #32) padding shall be applied to the outermost ends of
	// the string representation of the receiver, as well as around it's
	// comparison operator, when applicable.
	//
	// See also the NoPadding method for either of these types.
	IsPadded() bool

	// Stack: IsNesting returns a Boolean value indicative of whether the
	// receiver contains one (1) or more slices that are Stack instances.
	//
	// Condition: IsNesting returns a Boolean value indicative of whether
	// the Expression value set within the receiver is a Stack instance.
	//
	// See also the NoNesting and CanNest method for either of these types.
	IsNesting() bool

	// Unmarshal will unmarshal the receiver instance.
	//
	// Stack or Stack-alias instances become []any.
	//
	// Condition or Condition-alias instances become an anonymous
	// struct as the sole slice member within []any.
	Unmarshal() ([]any, error)

	// Stack: CanNest returns a Boolean value indicative of whether the
	// receiver is allowed to append (Push) additional Stack instances
	// into its collection of slices.
	//
	// Condition: CanNest returns a Boolean value indicative of whether
	// the receiver is allowed to set a Stack instance as its Expression
	// value.
	//
	// See also the NoNesting and IsNesting methods for either of these
	// types.
	CanNest() bool

	// ID returns the identifier assigned to the receiver, if set. This
	// may be anything the user chooses to set, or it may be auto-assigned
	// using the `_addr` or `_random` input keywords.
	//
	// See also the SetID method for Condition and Stack instances.
	ID() string

	// Addr returns the string representation of the pointer address of
	// the embedded receiver instance. This is mainly useful in cases
	// where debugging/troubleshooting may be underway, and the ability
	// to distinguish unnamed instances would be beneficial. It may also
	// be used as an alternative to the tedium of manually naming objects.
	Addr() string

	// String returns the string representation of the receiver. Note that
	// if the receiver is a BASIC Stack, string representation shall not be
	// possible.
	String() string

	// Stack: Category returns the categorical string label assigned to the
	// receiver instance during initialization. This will be one of AND, OR,
	// NOT, LIST and BASIC (these values may be folded case-wise depending
	// on the state of the Fold bit, as set by the user).
	//
	// Condition: Category returns the categorical string label assigned to
	// the receiver instance by the user at any time. Condition categorical
	// labels are entirely user-controlled.
	//
	// See also the SetCategory method for either of these types.
	Category() string

	// Err returns the most recently set instance of error within the receiver.
	// This is particularly useful for users who dislike fluent-style method
	// execution and would prefer to operate in the traditional "if err != nil ..."
	// style.
	//
	// See also the SetErr method for Condition and Stack instances.
	Err() error

	// Valid returns an instance of error resulting from a cursory review of the
	// receiver without any special context. Nilness, quality-of-initialization
	// and other rudiments are checked.
	//
	// This method is most useful when users apply a ValidityPolicy method or
	// function to the receiver instance, allowing full control over what they
	// deem "valid". See the SetValidityPolicy method for Condition and Stack
	// instances for details.
	Valid() error

	// Logger returns the underlying instance of *log.Logger, which may be set by
	// the package by defaults, or supplied by the user in a piecemeal manner.
	//
	// See also the SetLogger method for Condition and Stack instances.
	Logger() *log.Logger
}

var (
	unexpectedReceiverState error = errorf("Receiver does not contain an expected instance; aborting")
)
