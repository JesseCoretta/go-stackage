package stackage

import "sync"

/*
nodeConfig contains configuration information that is
used by the circumscribing instance of Stack or Condition.
An instance of the *nodeConfig type shall always occupy
slice zero (0) of a *stack value, else the stack is totally
invalid. An instance of the *nodeConfig type shall always
occupy the cfg struct field of the *condition type, though
not all config fields will apply.
*/
type nodeConfig struct {
	id  string             // optional identifier
	cat string             // optional categorical identifier
	cap int                // optional stack capacity
	evl Evaluator          // closure evaluator
	ppf PushPolicy         // closure filterer
	vpf ValidityPolicy     // closure validator
	rpf PresentationPolicy // closure stringer
	msg chan Message       // optional debug/err chan interface
	opt cfgFlag            // parens, cfold, lonce, etc...
	enc [][]string         // val encapsulators
	err error              // error pertaining to the outer type state (Condition/Stack)

	typ stackType   // stacks only: defines the typ/kind of stack
	sym string      // stacks only: user-controlled symbol char(s)
	ljc string      // [list] stacks only: joining delim
	mtx *sync.Mutex // stacks only: optional locking system
}

/*
cfgFlag contains left-shifted bit values that can represent
one of several configuration "flag states".
*/
type cfgFlag uint16

/*
ONE of 'and', 'or', 'not' or 'list'
belongs in the nodeConfig.typ struct
field value.
*/
const (
	_ stackType = iota
	and
	or
	not
	list
	cond
	basic
)

type stackType uint8

/*
String is a stringer method that returns the string
representation of the receiver.
*/
func (r stackType) String() string {
	t := badStack

	switch r {
	case and:
		t = `AND`
	case or:
		t = `OR`
	case not:
		t = `NOT`
	case list:
		t = `LIST`
	case basic:
		t = `BASIC`
	}

	return t
}

func (r nodeConfig) stackType() stackType {
	return r.typ
}

func (r stack) stackType() stackType {
	sc, _ := r.config()
	return sc.stackType()
}

func (r nodeConfig) isError() bool {
	return r.err != nil
}

func (r nodeConfig) error() error {
	return r.err
}

func (r *nodeConfig) setError(err error) {
	r.err = err
}

/*
The left-shifted cfgFlag values belong
in the first (1st) slice of cfgOpts.
These will come into play during a stack's
string representation.
*/
const (
	parens cfgFlag = 1 << iota //     1 // current stack (not its values) should be encapsulated it in parenthesis in string representation
	cfold                      //     2 // fold case of 'AND', 'OR' and 'NOT' to 'and', 'or' and 'not' or vice versa
	nspad                      //     4 // don't pad slices with a single space character when represented as a string
	lonce                      //     8 // only use operator once per stack, and only as the leading element; mainly for LDAP filters
	negidx                     //    16 // enable negative index support
	fwdidx                     //    32 // enable forward index support
	joinl                      //    64 // list joining value
	ronly                      //   128 // stack is read-only
	_                          //   256
	_                          //   512
	_                          //  1024
	_                          //  2048
	_                          //  4096
	_                          //  8192
	_                          // 16384
	_                          // 32768
)

/*
isZero returns a boolean value indicative of whether the
receiver is nil, or uninitialized.
*/
func (r *nodeConfig) isZero() bool {
	return r == nil
}

func (r *nodeConfig) setID(id string) {
	r.id = id
}

func (r *nodeConfig) setCat(cat string) {
	r.cat = cat
}

/*
kind returns the string representation of the kind value
set within the receiver's typ field. Case will be folded
if the cfold bit is enabled.
*/
func (r *nodeConfig) kind() (kind string) {
	kind = `null`
	if r.isZero() {
		return
	}

	switch r.typ {
	case and, or, not, list:
		kind = r.typ.String()
	}

	if r.positive(cfold) {
		kind = lc(kind)
	}

	return
}

/*
valid returns an error if the receiver is considered to be
invalid or nil.
*/
func (r *nodeConfig) valid() (err error) {
	if r.isZero() {
		err = errorf("%T instance is nil; aborting", r)
		return
	}
	if r.typ == 0x0 {
		err = errorf("%T instance defines no stack \"kind\", or %T is invalid", r, r)
		return
	}

	return r.error()
}

/*
positive returns a boolean value indicative of whether the specified
cfgFlag input value is "on" within the receiver's opt field.
*/
func (r nodeConfig) positive(x cfgFlag) bool {
	if err := r.valid(); err != nil {
		return false
	}

	return r.opt.positive(x)
}

/*
positive returns a boolean value indicative of whether the specified
cfgFlag input value is "on" within the receiver's uint16 bit value.
*/
func (r cfgFlag) positive(x cfgFlag) bool {
	return r&x != 0
}

/*
setOpt sets the specified cfgFlag to "on" within the receiver's
opt field.
*/
func (r *nodeConfig) setOpt(x cfgFlag) (err error) {
	if err = r.valid(); err != nil {
		return
	}

	r.opt.shift(x)
	return
}

/*
setEncap accepts input characters for use in controlled stack value
encapsulation.

For example, providing []string{`(`, `)`} (i.e.: a string slice that
contains L and R parens) will cause the individual values within the
current stack to undergo parenthetical encapsulation.

If a string slice containing only one (1) character is provided, that
character shall be used for both L and R.

If multiple slices are provided, they will each be used incrementally.
For example, if the following slices are received in the order shown:

- []string{`(`,`)`}, []string{`"`}

... then each value within the stack shall be encapsulated as ("value").

No single character can appear in more than one []string value.
*/
func (r *nodeConfig) setEncap(x ...any) {
	if len(x) == 0 {
		r.enc = [][]string{}
		return
	}

	for sl := 0; sl < len(x); sl++ {
		switch tv := x[sl].(type) {
		case string:
			r.setStringSliceEncap([]string{tv})
		case []string:
			r.setStringSliceEncap(tv)
		}
	}
}

/*
setJoinDelim is a private method invoked by stack.setJoinDelim.
*/
func (r *nodeConfig) setJoinDelim(x string) {
	r.ljc = x
}

/*
setStringSliceEncap is a private method called by nodeConfig.setEncap,
and determines which encapsulation method to call based on the encap
input length (x).
*/
func (r *nodeConfig) setStringSliceEncap(x []string) {
	switch len(x) {
	case 0:
		return
	case 1:
		r.setStringSliceEncapOne(x)
	default:
		r.setStringSliceEncapTwo(x)
	}
}

/*
setStringSliceEncapOne is a private method called by setStringSliceEncap
for the purpose of assigning a single character for L and R encapsulation.
*/
func (r *nodeConfig) setStringSliceEncapOne(x []string) {
	var found bool
	for u := 0; u < len(r.enc); u++ {
		if found = strInSlice(x[0], r.enc[u]); found {
			break
		}
	}

	if !found {
		r.enc = append(r.enc, x)
	}
}

/*
setStringSliceEncapOne is a private method called by setStringSliceEncap
for the purpose of assigning a pair of characters for L and R encapsulation.
*/
func (r *nodeConfig) setStringSliceEncapTwo(x []string) {
	if len(x[0])|len(x[1]) == 0 {
		return
	}

	var found bool
	for i := 0; i < 2; i++ {
		for u := 0; u < len(r.enc); u++ {
			if found = strInSlice(x[i], r.enc[u]); found {
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		r.enc = append(r.enc, x)
	}
}

/*
setSymbol assigns the given string to the receiver, and will use
this value for as a particular operator symbol that is appropriate
for the stack type. This will only have an effect is the symbol
bit is positive.
*/
func (r *nodeConfig) setSymbol(sym string) {
	r.sym = sym
}

/*
shift sets the specified cfgFlag to "on" within the receiver's
uint16 bit value.
*/
func (r *cfgFlag) shift(x cfgFlag) {
	*r |= x
}

/*
unsetOpt sets the specified cfgFlag to "off" within the receiver's
opt field.
*/
func (r *nodeConfig) setMutex() {
	if err := r.valid(); err != nil {
		return
	}

	if r.mtx == nil {
		r.mtx = &sync.Mutex{}
	}
}

/*
unsetOpt sets the specified cfgFlag to "off" within the receiver's
opt field.
*/
func (r *nodeConfig) unsetOpt(x cfgFlag) (err error) {
	if err = r.valid(); err != nil {
		return
	}

	r.opt.unshift(x)
	return
}

/*
unsetOpt sets the specified cfgFlag to "off" within the receiver's
uint16 bit value.
*/
func (r *cfgFlag) unshift(x cfgFlag) {
	*r = *r &^ x
}

/*
toggleOpt will swap the receiver's opt field to reflect a state contrary
to its current value in terms of the input value. In other words, if the
state is "on" for a flag, it will be toggled to "off" and vice versa.
*/
func (r *nodeConfig) toggleOpt(x cfgFlag) (err error) {
	if err = r.valid(); err != nil {
		return
	}

	r.opt.toggle(x)
	return
}

/*
toggleOpt will shift the receiver's uint16 bits to reflect a state contrary
to its current value in terms of the input value. In other words, if the
state is "on" for a flag, it will be toggled to "off" and vice versa.
*/
func (r *cfgFlag) toggle(x cfgFlag) {
	if r.positive(x) {
		r.unshift(x)
		return
	}
	r.shift(x)
}

/*
mutex returns the *sync.Mutext instance, alongside a presence-indicative
boolean value.
*/
func (r *stack) mutex() (mutex *sync.Mutex, found bool) {
	sc, _ := r.config()
	return sc.mtx, sc.mtx != nil
}

func (r *nodeConfig) canWriteMessage() bool {
	return r.msg != nil
}
