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
				if s := assert.String(); s != badStack && len(s) > 0 {
					S = assert
					converted = true
					return
				}
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
				if s := assert.String(); s != badCond && len(s) > 0 {
					C = assert
					converted = true
					return
				}
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

	return `<invalid_number>`
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
following types:

• Stack

• Condition

This interface type offers users an alternative to the tedium of
repeated type assertion for every Stack and Condition instance they
encounter. This may be particularly useful in situations where the
act of traversal is conducted upon a Stack instance that contains
myriad hierarchies and nested contexts of varying types.

This is not a complete "replacement" for the explicit use of package
defined types nor their aliased counterparts. The Interface interface
only extends methods that are read-only in nature AND common to both
of the above types (and any aliased counterparts).

To access the entire breadth of available methods for the underlying
type instance, manual type assertion shall be necessary.

Users SHOULD adopt this interface signature for use in their solutions
as needed, though it is not strictly required.
*/
type Interface interface {
	Len() int

	IsInit() bool
	IsZero() bool
	CanNest() bool
	IsParen() bool
	IsEncap() bool
	IsPadded() bool
	IsNesting() bool

	ID() string
	Addr() string
	String() string
	Category() string

	Err() error
	Valid() error

	Logger() *log.Logger
}
