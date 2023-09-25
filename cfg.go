package stackage

import (
	"sync"
	"time"
)

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
	log *logSystem         // logging subsystem
	opt cfgFlag            // parens, cfold, lonce, etc...
	enc [][]string         // val encapsulators
	err error              // error pertaining to the outer type state (Condition/Stack)
	aux Auxiliary          // auxiliary admin-related object storage, user managed

	typ stackType   // stacks only: defines the typ/kind of stack
	sym string      // stacks only: user-controlled symbol char(s)
	ljc string      // [list] stacks only: joining delim
	mtx *sync.Mutex // stacks only: optional locking system
	ldr *time.Time  // for lock duration; ephemeral, nil if not locked / no locking capabilities
	ord bool        // true = FIFO, false = LIFO (default); applies to stacks only
}

/*
Auxiliary is a map[string]any type alias extended by this package. It
can be created within any Stack instance when [re]initialized using
the SetAuxiliary method extended through instances of the Stack type,
and can be accessed using the Auxiliary() method in similar fashion.

The Auxiliary type extends four (4) methods: Get, Set, Len and Unset.
These are purely for convenience. Given that instances of this type
can easily be cast to standard map[string]any by the user, the use of
these methods is entirely optional.

The Auxiliary map instance is available to be leveraged in virtually
any way deemed appropriate by the user. Its primary purpose is for
storage of any instance(s) pertaining to the *administration of the
stack*, as opposed to the storage of content normally submitted *into*
said stack.

Examples of suitable instance types for storage within the Auxiliary
map include, but are certainly not limited to:

  - HTTP server listener / mux
  - HTTP client, agent
  - Graphviz node data
  - Go Metrics meters, gauges, etc.
  - Redis cache
  - bytes.Buffer
  - ANTLR parser
  - text/template instances
  - channels

Which instances are considered suitable for storage within Auxiliary map
instances is entirely up to the user. This package shall not impose ANY
controls or restrictions regarding the content within this instances of
this type, nor its behavior.
*/
type Auxiliary map[string]any

/*
Len returns the integer length of the receiver, defining the number of
key/value pairs present.
*/
func (r Auxiliary) Len() int {
	return len(r)
}

/*
Get returns the value associated with key, alongside a presence-indicative
Boolean value (ok).

Even if found within the receiver instance, if value is nil, ok shall be
explicitly set to false prior to return.

Case is significant in the matching process.
*/
func (r Auxiliary) Get(key string) (value any, ok bool) {
	if r != nil {
		value, ok = r[key]
	}

	return
}

/*
Set associates key with value, and assigns to receiver instance. See
also the Unset method.

If the receiver is not initialized, a new allocation is made.
*/
func (r Auxiliary) Set(key string, value any) Auxiliary {
	if r != nil {
		r[key] = value
	}
	return r
}

/*
Unset removes the key/value pair, identified by key, from the receiver
instance, if found. See also the Set method.

This method internally calls the following builtin:

	delete(*rcvr,key)

Case is significant in the matching process.
*/
func (r Auxiliary) Unset(key string) Auxiliary {
	if _, found := r[key]; found {
		delete(r, key)
	}
	return r
}

/*
cfgFlag contains left-shifted bit values that can represent
one of several configuration "flag states".
*/
type cfgFlag uint16

var cfgFlagMap map[cfgFlag]string

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
	case cond:
		t = `CONDITION` // just for logging
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

func (r nodeConfig) getErr() error {
	return r.err
}

func (r *nodeConfig) setErr(err error) {
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
	nnest                      //   256 // stack/condition does not allow stack/stack alias instances as slice members or expression value
	etrav                      //   512 // enhanced traversal support (slices, int-keyed maps)
	_                          //  1024
	_                          //  2048
	_                          //  4096
	_                          //  8192
	_                          // 16384
	_                          // 32768
)

/*
isZero returns a Boolean value indicative of whether the
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
	if !r.isZero() {
		switch r.typ {
		case and, or, not, list, cond, basic:
			kind = foldValue(r.positive(cfold), r.typ.String())
		}
	}

	return
}

/*
valid returns an error if the receiver is considered to be
invalid or nil.
*/
func (r *nodeConfig) valid() (err error) {
	err = errorf("%T instance is nil; aborting", r)
	if !r.isZero() {
		err = errorf("%T instance defines no stack \"kind\", or %T is invalid", r, r)
		if r.typ != 0x0 {
			err = nil
		}
	}

	return
}

/*
positive returns a Boolean value indicative of whether the specified
cfgFlag input value is "on" within the receiver's opt field.
*/
func (r nodeConfig) positive(x cfgFlag) (is bool) {
	if err := r.valid(); err == nil {
		is = r.opt.positive(x)
	}
	return
}

/*
positive returns a Boolean value indicative of whether the specified
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
	if err = r.valid(); err == nil {
		r.opt.shift(x)
	}
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
setListDelimiter is a private method invoked by stack.setListDelimiter.
*/
func (r *nodeConfig) setListDelimiter(x string) {
	if r.typ == list {
		r.ljc = x
	}
}

/*
getListDelimiter is a private method invoked by stack.getListDelimiter.
*/
func (r nodeConfig) getListDelimiter() string {
	return r.ljc
}

/*
setStringSliceEncap is a private method called by nodeConfig.setEncap,
and determines which encapsulation method to call based on the encap
input length (x).
*/
func (r *nodeConfig) setStringSliceEncap(x []string) {
	switch len(x) {
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
	var found bool
	for i := 0; i < 2 && !found; i++ {
		for u := 0; u < len(r.enc) && !found; u++ {
			found = strInSlice(x[i], r.enc[u])
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

func (r cfgFlag) String() (f string) {
	for k, v := range cfgFlagMap {
		if k == r {
			f = v
			break
		}
	}

	return
}

/*
setMutex enables the receiver's mutual exclusion
locking capabilities.
*/
func (r *nodeConfig) setMutex() {
	if err := r.valid(); err == nil {
		if r.mtx == nil {
			r.mtx = &sync.Mutex{}
		}
	}
}

/*
unsetOpt sets the specified cfgFlag to "off" within the receiver's
opt field.
*/
func (r *nodeConfig) unsetOpt(x cfgFlag) (err error) {
	if err = r.valid(); err == nil {
		r.opt.unshift(x)
	}
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
	if err = r.valid(); err == nil {
		r.opt.toggle(x)
	}
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
Boolean value.
*/
func (r *stack) mutex() (*sync.Mutex, bool) {
	sc, _ := r.config()
	return sc.mtx, sc.mtx != nil
}

func init() {
	cfgFlagMap = map[cfgFlag]string{
		parens: `parenthetical`,
		negidx: `neg_index`,
		fwdidx: `fwd_index`,
		cfold:  `case_fold`,
		nspad:  `no_whsp_pad`,
		lonce:  `lead_once`,
		joinl:  `join_list`,
		ronly:  `read_only`,
		nnest:  `no_nest`,
		etrav:  `enhanced_traversal`,
	}
}
