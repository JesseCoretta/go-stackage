package stackage

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unicode"
)

/*
frequently-used package import function aliases.
*/
var (
	typOf   func(any) reflect.Type            = reflect.TypeOf
	valOf   func(any) reflect.Value           = reflect.ValueOf
	printf  func(string, ...any) (int, error) = fmt.Printf
	sprintf func(string, ...any) string       = fmt.Sprintf
	lc      func(string) string               = strings.ToLower
	ilc     func(rune) bool                   = unicode.IsLower
	uc      func(string) string               = strings.ToUpper
	iuc     func(rune) bool                   = unicode.IsUpper
	split   func(string, string) []string     = strings.Split
	trimS   func(string) string               = strings.TrimSpace
	join    func([]string, string) string     = strings.Join
	now     func() time.Time                  = time.Now
)

/*
Message is an optional type for use when a user-supplied Message channel has
been initialized and provided to one (1) or more Stack or Condition instances.

Instances of this type shall contain diagnostic, error and debug information
pertaining to current operations of the given Stack or Condition instance.
*/
type Message struct {
	ID   string    `json:"id"`
	Msg  string    `json:"message"`
	Tag  string    `json:"message_tag"`
	Type string    `json:"message_type"`
	Addr string    `json:"memory_address"`
	Len  int       `json:"current_length"`
	Cap  int       `json:"maximum_length,omitempty"` // omit if zero, meaning no cap was set.
	Time time.Time `json:"current_time"`
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r Message) String() string {
	if !r.Valid() {
		return ``
	}

	time := sprintf("[%04d%02d%02d%02d%02d%02d]",
		r.Time.Year(),
		r.Time.Month(),
		r.Time.Day(),
		r.Time.Hour(),
		r.Time.Minute(),
		r.Time.Second())

	var msgid string = sprintf("[%s][%s::%s]", time, r.Type, r.Addr)
	if len(r.ID) > 0 {
		msgid = sprintf("[%s::%s(%s)]", r.Type, r.ID, msgid)
	}
	var lencap string = sprintf("[%d/%d]", r.Len, r.Cap)

	return sprintf("%s%s [%s] %s", msgid, lencap, r.Tag, r.Msg)
}

/*
Valid returns a Boolean value indicative of whether the receiver
is perceived to be valid.
*/
func (r Message) Valid() bool {
	return (r.Type != `UNKNOWN` &&
		len(r.Addr) > 0 &&
		!r.Time.IsZero() &&
		len(r.Msg) > 0 &&
		len(r.Tag) > 0)
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
	// If it isn't a Stack alias, but is a
	// genuine Stack, just pass it back
	// with a thumbs-up ...
	if st, isStack := u.(Stack); isStack {
		S = st
		converted = isStack
		return
	}

	a := typOf(u)       // current (src) type
	b := typOf(Stack{}) // target (dest) type
	if a.ConvertibleTo(b) {
		X := valOf(u).Convert(b).Interface()
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
	// If it isn't a Condition alias, but is a
	// genuine Condition, just pass it back
	// with a thumbs-up ...
	if co, isCond := u.(Condition); isCond {
		C = co
		converted = isCond
		return
	}

	a := typOf(u)       // current (src) type
	b := typOf(Stack{}) // target (dest) type
	if a.ConvertibleTo(b) {
		X := valOf(u).Convert(b).Interface()
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
a quick means of getting the caller name for logging purposes.
*/
func fmname() string {
	x, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(x).Name()
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
