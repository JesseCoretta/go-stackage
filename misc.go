package stackage

import (
	"errors"
	"fmt"
	"log"
	"math/rand" // not for crypto, don't worry :)
	"reflect"
	"runtime"
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
	split   func(string, string) []string       = strings.Split
	trimS   func(string) string                 = strings.TrimSpace
	join    func([]string, string) string       = strings.Join
	now     func() time.Time                    = time.Now
)

const (
	randChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randIDSize = 24
)

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

func isETravEligible(ok bool, x any) bool {
	return (isIntKeyedMap(x) || isSliceType(x)) && ok
}

func isIntKeyedMap(x any) bool {
	typ := typOf(x)
	if isPtr(x) {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Map {
		return false
	}
	return typ.Key().Kind() == reflect.Int
}

func isSliceType(x any) bool {
	typ := reflect.TypeOf(x)
	if isPtr(x) {
		typ = typ.Elem()
	}

	return typ.Kind() == reflect.Slice
}

/*
errorf wraps errors.New and returns a non-nil instance of error
based upon a non-nil/non-zero msg input value with optional args.
*/
func errorf(msg any, x ...any) error {
	switch tv := msg.(type) {
	case string:
		if len(tv) > 0 {
			return errors.New(sprintf(tv, x...))
		}
	case error:
		if tv != nil {
			return errors.New(sprintf(tv.Error(), x...))
		}
	}

	return nil
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
cleanly. In particular, this will be necessary during
the parsing (marshaling) of text rules into proper ACI
type instances.
*/
func condenseWHSP(b string) (a string) {
	// remove leading and trailing
	// WHSP/HTAB characters ...
	b = trimS(b)

	var last bool // previous char was WHSP or HTAB.

	for i := 0; i < len(b); i++ {
		c := rune(b[i])
		switch c {

		// match either whsp OR horizontal tab
		case rune(9), rune(32):
			if !last {
				last = !last
				a += string(rune(32)) // only add whsp (not htab) for consistency
			}

		// match all other chars ...
		default:
			if last {
				last = !last
			}
			a += string(c)
		}
	}

	return
}

/*
stackTypeAliasConverter attempts to convert any (u) back to a bonafide instance
of Stack. This will only work if input value u is a type alias of Stack. An
instance of Stack is returned along with a success-indicative Boolean value.
*/
func stackTypeAliasConverter(u any) (S Stack, converted bool) {
	if u == nil {
		return
	}

	// If it isn't a Stack alias, but is a
	// genuine Stack, just pass it back
	// with a thumbs-up ...
	if st, isStack := u.(Stack); isStack {
		S = st
		converted = isStack
		return
	}

	a := typOf(u) // current (src) type
	v := valOf(u) // current (src) value

	// unwrap any pointers for
	// maximum compatibility
	if isPtr(u) {
		a = a.Elem()
		v = v.Elem()
	}

	b := typOf(Stack{}) // target (dest) type
	if a.ConvertibleTo(b) {
		X := v.Convert(b).Interface()
		if assert, ok := X.(Stack); ok {
			if !assert.IsZero() {
				S = assert
				converted = true
				return
			}
		}
	}

	return
}

/*
conditionTypeAliasConverter attempts to convert any (u) back to a bonafide instance
of Condition. This will only work if input value u is a type alias of Condition. An
instance of Condition is returned along with a success-indicative Boolean value.
*/
func conditionTypeAliasConverter(u any) (C Condition, converted bool) {
	if u == nil {
		return
	}

	// If it isn't a Condition alias, but is a
	// genuine Condition, just pass it back
	// with a thumbs-up ...
	if co, isCond := u.(Condition); isCond {
		C = co
		converted = isCond
		return
	}

	a := typOf(u) // current (src) type
	v := valOf(u) // current (src) value

	// unwrap any pointers for
	// maximum compatibility
	if isPtr(u) {
		a = a.Elem()
		v = v.Elem()
	}

	b := typOf(Condition{}) // target (dest) type
	if a.ConvertibleTo(b) {
		X := v.Convert(b).Interface()
		if assert, ok := X.(Condition); ok {
			if !assert.IsZero() {
				C = assert
				converted = true
				return
			}
		}
	}

	return
}

/*
getStringer uses reflect to obtain and return a given
type instance's String ("stringer") method, if present.
If not, nil is returned.
*/
func getStringer(x any) func() string {
	if x == nil {
		return nil
	}

	v := valOf(x)
	if v.IsZero() {
		return nil
	}
	method := v.MethodByName(`String`)
	if method.Kind() == reflect.Invalid {
		return nil
	}

	if meth, ok := method.Interface().(func() string); ok {
		return meth
	}

	return nil
}

/*
getIDFunc uses reflect to obtain and return a given
type instance's ID method, if present. If not, a zero
string is returned.
*/
func getIDFromAny(x any) func() string {
	if x == nil {
		return nil
	}

	v := valOf(x)
	if v.IsZero() {
		return nil
	}

	method := v.MethodByName(`ID`)
	if method.Kind() == reflect.Invalid {
		return nil
	}

	if meth, ok := method.Interface().(func() string); ok {
		return meth
	}

	return nil
}

/*
a quick means of getting the caller name for logging purposes.
*/
func fmname() string {
	x, _, _, _ := runtime.Caller(1)
	name := split(runtime.FuncForPC(x).Name(), string(rune(46)))
	return uc(name[len(name)-1])
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
			v = sprintf("%s%s%s", sl[0], v, sl[0])
		case 2:
			// char 0 = L, char 1 = R
			v = sprintf("%s%s%s", sl[0], v, sl[1])
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
	return sprintf("%s%s%s", pad, value, pad)
}

/*
foldValue will apply lc (Strings.ToLower) and uc (Strings.ToUpper)
to the value based on the "do" disposition (do, or do not).
*/
func foldValue(do bool, value string) string {
	if len(value) == 0 {
		return value // ???
	}

	if do {
		if iuc(rune(value[0])) {
			return lc(value) // fold to lower
		}
		return uc(value) // fold to upper
	}

	return value // do not.
}

func isPtr(x any) bool {
	if x == nil {
		return false
	}

	return typOf(x).Kind() == reflect.Ptr
}

func isNumberPrimitive(x any) bool {
	switch x.(type) {
	case int, int8, int16, int32, int64,
		float32, float64, complex64, complex128,
		uint, uint8, uint16, uint32, uint64, uintptr:
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

func isKnownPrimitive(x any) bool {
	if isStringPrimitive(x) {
		return true
	} else if isNumberPrimitive(x) {
		return true
	} else if isBoolPrimitive(x) {
		return true
	}

	return false
}

func numberStringer(x any) string {
	switch tv := x.(type) {
	case float32, float64:
		return floatStringer(tv)
	case complex64, complex128:
		return complexStringer(tv)
	case int, int8, int16, int32, int64:
		return intStringer(tv)
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return uintStringer(tv)
	}

	return `NaN`
}

func primitiveStringer(x any) string {
	if !isKnownPrimitive(x) {
		return ``
	}
	if isBoolPrimitive(x) {
		return boolStringer(x)
	}
	if isNumberPrimitive(x) {
		return numberStringer(x)
	}
	if isStringPrimitive(x) {
		return x.(string)
	}

	return `unsupported_stringer_type`
}

func boolStringer(x any) string {
	return sprintf("%t", x.(bool))
}

func floatStringer(x any) string {
	switch tv := x.(type) {
	case float64:
		return sprintf("%.02f", tv)
	}

	return sprintf("%.02f", x.(float32))
}

func complexStringer(x any) string {
	switch tv := x.(type) {
	case complex128:
		return sprintf("%v", tv)
	}

	return sprintf("%v", x.(complex64))
}

func uintStringer(x any) string {
	switch tv := x.(type) {
	case uint8:
		return sprintf("%d", tv)
	case uint16:
		return sprintf("%d", tv)
	case uint32:
		return sprintf("%d", tv)
	case uint64:
		return sprintf("%d", tv)
	case uintptr:
		return sprintf("%d", tv)
	}

	return sprintf("%d", x.(uint))
}

func intStringer(x any) string {
	switch tv := x.(type) {
	case int8:
		return sprintf("%d", tv)
	case int16:
		return sprintf("%d", tv)
	case int32:
		return sprintf("%d", tv)
	case int64:
		return sprintf("%d", tv)
	}

	return sprintf("%d", x.(int))
}

/*
Interface is an interface type qualified through instances of the
following types (whether native, or type-aliased):

  - Stack
  - Condition

This interface type offers users an alternative to the tedium of
repeated type assertion for every Stack and Condition instance they
encounter. This may be particularly useful in situations where the
act of traversal is conducted upon a Stack instance that contains
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
