package stackage

import (
	"log"
)

/*
Stack embeds slices of any ([]any) in pointer form
and extends methods allowing convenient interaction
with stack structures.
*/
type Stack struct {
	*stack
}

/*
stack represents the underlying ordered slice type, which
is embedded (in pointer form) within instances of Stack.
*/
type stack []any

/*
List initializes and returns a new instance of Stack
configured as a simple list. Stack instances of this
design can be delimited using the SetDelimiter method.
*/
func List(capacity ...int) Stack {
	return Stack{newStack(list, capacity...)}
}

/*
And initializes and returns a new instance of Stack
configured as a Boolean ANDed stack.
*/
func And(capacity ...int) Stack {
	return Stack{newStack(and, capacity...)}
}

/*
Or initializes and returns a new instance of Stack
configured as a Boolean ORed stack.
*/
func Or(capacity ...int) Stack {
	return Stack{newStack(or, capacity...)}
}

/*
Not initializes and returns a new instance of Stack
configured as a Boolean NOTed stack.
*/
func Not(capacity ...int) Stack {
	return Stack{newStack(not, capacity...)}
}

/*
Basic initializes and returns a new instance of Stack, set
for basic operation only.

Please note that instances of this design are not eligible
for string representation, value encaps, delimitation, and
other presentation-related string methods. As such, a zero
string (“) shall be returned should String() be executed.

PresentationPolicy instances cannot be assigned to Stacks
of this design.
*/
func Basic(capacity ...int) Stack {
	return Stack{newStack(basic, capacity...)}
}

/*
newStack initializes a new instance of *stack, configured
with the kind (t) requested by the user. This function
should only be executed when creating new instances of
Stack, which embeds the *stack type.
*/
func newStack(t stackType, c ...int) *stack {
	switch t {
	case and, or, not, list, basic:
		// ok
	default:
		return nil
	}

	var (
		data map[string]string = make(map[string]string, 0)
		cfg  *nodeConfig       = new(nodeConfig)
		st   stack
	)

	cfg.typ = t
	data[`kind`] = t.String()

	cfg.setLogger(sLogDefault)
	if !logDiscard(cfg.getLogger()) {
		data[`laddr`] = sprintf("%p", sLogDefault)
		data[`lpfx`] = sLogDefault.Prefix()
		data[`lflags`] = sprintf("%d", sLogDefault.Flags())
		data[`lwriter`] = sprintf("%T", sLogDefault.Writer())
	}

	if len(c) > 0 {
		if c[0] > 0 {
			st = make(stack, 0, c[0]+1) // 1 for cfg slice offset
			cfg.cap = c[0] + 1
			data[`capacity`] = sprintf("%d", cfg.cap)
		}
	} else {
		st = make(stack, 0)
	}

	st = append(st, cfg)
	instance := &st
	data[`addr`] = instance.addr()
	st.debug(`ALLOC`, data)

	return instance
}

/*
IsZero returns a Boolean value indicative of whether the
receiver is nil, or uninitialized.
*/
func (r Stack) IsZero() bool {
	return r.stack.isZero()
}

/*
isZero is a private method called by Stack.IsZero.
*/
func (r *stack) isZero() bool {
	return r == nil
}

/*
Addr returns the string representation of the pointer
address for the receiver. This may be useful for logging
or debugging operations.
*/
func (r Stack) Addr() string {
	if !r.IsInit() {
		return ``
	}

	return r.stack.addr()
}

/*
addr returns the string representation of the pointer
address for the receiver.
*/
func (r *stack) addr() string {
	if r.isZero() {
		return ``
	}

	return sprintf("%p", *r)
}

/*
Err returns the error residing within the receiver, or nil
if no error condition has been declared.

This method will be particularly useful for users who do
not care for fluent-style operations and would instead
prefer to operate in a more procedural fashion.

Note that a chained sequence of method calls of this type
shall potentially obscure error conditions along the way,
as each successive method may happily overwrite any error
instance already present.
*/
func (r Stack) Err() error {
	if !r.IsInit() {
		return nil
	}
	return r.getErr()
}

/*
SetErr sets the underlying error value within the receiver
to the assigned input value err, whether nil or not.

This method may be most valuable to users who have chosen
to extend this type by aliasing, and wish to control the
handling of error conditions in another manner.
*/
func (r Stack) SetErr(err error) Stack {
	if !r.IsInit() {
		return r
	}

	r.stack.setErr(err)
	return r
}

/*
setErr assigns an error instance, whether nil or not, to
the underlying receiver configuration.
*/
func (r *stack) setErr(err error) {
	sc, _ := r.config()
	sc.setErr(err)
}

/*
getErr returns the instance of error, whether nil or not, from
the underlying receiver configuration.
*/
func (r stack) getErr() error {
	sc, _ := r.config()
	return sc.getErr()
}

/*
kind returns the string representation of the kind value
set within the receiver's configuration value.
*/
func (r stack) kind() string {
	sc, _ := r.config()
	return sc.kind()
}

/*
Valid returns an error if the receiver lacks a configuration
value, or is unset as a whole. This method does not check to
see whether the receiver is in an error condition regarding
user content operations (see the E method).
*/
func (r Stack) Valid() (err error) {
	if r.stack == nil {
		err = errorf("embedded %T instance is nil", r.stack)
		return
	}

	if err = r.stack.valid(); err != nil {
		return
	}

	// try to see if the user provided a
	// validity function
	if meth := r.getValidityPolicy(); meth != nil {
		err = meth(r)
	}

	return
}

/*
valid is a private method called by Stack.Valid.
*/
func (r stack) valid() (err error) {
	if !r.isInit() {
		err = errorf("%T instance is not initialized", Stack{})
		return
	}

	_, err = r.config()
	return
}

/*
IsInit returns a Boolean value indicative of whether the
receiver has been initialized using any of the following
package-level functions:

• And

• Or

• Not

• List

• Basic

This method does not take into account the presence (or
absence) of any user-provided values (e.g.: a length of
zero (0)) can still return true.
*/
func (r Stack) IsInit() bool {
	if r.IsZero() {
		return false
	}

	return r.stack.isInit()
}

/*
isInit is a private method called by Stack.IsInit.
*/
func (r stack) isInit() bool {
	return r.stackType() != 0x0
}

/*
SetLogger assigns the specified logging facility to the receiver
instance.

Logging is available but is set to discard all events by default.

An active logging subsystem within the receiver supercedes the
default package logger.

The following types/values are permitted:

• string: `none`, `off`, `null`, `discard` will turn logging off *

• string: `stderr` or `errors` will set basic stderr logging *

• int: 0 will turn logging off

• int: 1 will set basic stderr logging

• *log.Logger: user-defined *log.Logger instance will be set; it
should not be nil

*Case is not significant in the string matching process.

Logging may also be set globally using the SetDefaultLogger
package level function. Similar semantics apply.
*/
func (r Stack) SetLogger(logger any) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	r.stack.setLogger(logger)
	return r
}

/*
Logger returns the *log.Logger instance. This can be used for quick
access to the log.Logger type's methods in a manner such as:

	r.Logger().Fatalf("We died")

It is not recommended to modify the return instance for the purpose
of disabling logging outright (see Stack.SetLogger method as well
as the SetDefaultLogger package-level function for ways of doing
this easily).
*/
func (r Stack) Logger() *log.Logger {
	if !r.IsInit() {
		return nil
	}

	return r.stack.logger()
}

/*
Transfer will iterate the receiver (r) and add all slices
contained therein to the destination instance.

The following circumstances will result in a false return:

- Capacity constraints are in-force within the destination
instance, and the transfer request (if larger than the sum
number of available slices) cannot proceed as a result

- The destination instance is nil, or has not been properly
initialized

- The receiver instance (r) contains no slices to transfer

The receiver instance (r) is not modified in any way as a
result of calling this method. If the receiver undergoes a
call to its Reset() method, only the receiver instance will
be emptied, and the transferred slices within the submitted
destination instance shall remain.
*/
func (r Stack) Transfer(dest Stack) bool {
	return r.transfer(dest.stack)
}

/*
transfer is a private method executed by the Stack.Transfer
method. It will return a *stack instance containing the same
slices as the receiver (r). Configuration is not copied, nor
is a destination *stack subject to initialization within this
method. Thus, the user must submit a *stack instance ready to
receive slices immediately.

Capacity enforcement is honored. If the source (r) contains
more slices than a non-zero destination capacity allows, the
operation is canceled outright and false is returned.
*/
func (r *stack) transfer(dest *stack) bool {

	if r.ulen() == 0 || dest == nil {
		// nothing to xfer
		return false
	}

	if !r.isInit() {
		// dest MUST be initialized in some way
		// or this happens ...
		return false
	}

	// if a capacity was set, make sure
	// the destination can handle it...
	if dest.cap() > 0 {
		if r.ulen() > dest.cap()-r.ulen() {
			// capacity is in-force, and
			// there are too many slices
			// to xfer.
			return false
		}
	}

	// xfer slices, without any regard for
	// nilness. Slice type is not subject
	// to discrimination.
	for i := 0; i < r.ulen(); i++ {
		sl, _, _ := r.index(i) // cfg offset handled by index method
		dest.push(sl)
	}

	// return newly populated stack
	return dest.ulen() >= r.ulen()
}

/*
Replace will overwrite slice idx using value x and returns a Boolean
value indicative of success.

If slice i does not exist (e.g.: idx > receiver len), then nothing is
altered and a false Boolean value is returned.

Use of the Replace method shall not result in fragmentation of the
receiver instance; this method does not honor any attempt to replace
any receiver slice value with nil.
*/
func (r Stack) Replace(x any, idx int) bool {
	return r.stack.replace(x, idx)
}

func (r *stack) replace(x any, i int) (ok bool) {
	// bail out if receiver or
	// input value is nil
	if r == nil || x == nil {
		return
	} else if !r.isInit() {
		return
	}

	if r.positive(ronly) {
		return
	}

	if i+1 > r.ulen() {
		return
	}

	(*r)[i+1] = x
	ok = true

	return
}

/*
Insert will insert value x to become the left index. For example,
using zero (0) as left shall result in value x becoming the first
slice within the receiver.

This method returns a Boolean value indicative of success. A value
of true indicates the receiver length became longer by one (1).

This method does not currently respond to forward/negative index
support. An integer value less than or equal to zero (0) shall
become zero (0). An integer value that exceeds the length of the
receiver shall become index len-1. A value that falls within the
bounds of the receiver's current length is inserted as intended.

Use of the Insert method shall not result in fragmentation of the
receiver instance, as any nil x value shall be discarded and not
considered for insertion into the stack.
*/
func (r Stack) Insert(x any, left int) bool {
	return r.stack.insert(x, left)
}

/*
insert is a private method called by Stack.Insert.
*/
func (r *stack) insert(x any, left int) (ok bool) {
	// bail out if receiver or
	// input value is nil
	if r == nil || x == nil {
		return
	}

	//fname := fmname()

	if r.positive(ronly) {
		return
	}

	// note the len before we start
	var u1 int = r.ulen()

	// bail out if a capacity has been set and
	// would be breached by this insertion.
	if u1+1 > r.cap() && r.cap() != 0 {
		return
	}

	r.lock()
	defer r.unlock()

	cfg, _ := r.config()

	// If left is greater-than-or-equal
	// to the user length, just push.
	if u1-1 < left {
		*r = append(*r, x)

		// Verify something was added
		ok = u1+1 == r.ulen()
		return
	}

	var R stack = make(stack, 0)
	left += 1

	// If left is less-than-or-equal to
	// zero (0), we'll use a new stack
	// alloc (R) and move items into it
	// in the appropriate order. The
	// new element will be the first
	// user slice.
	if left <= 1 {
		left = 1
		R = append(R, cfg)
		R = append(R, x)
		R = append(R, (*r)[left:]...)

		// If left falls within the user length,
		// append all elements up to and NOT
		// including the desired index, and also
		// append everything after desired index.
		// This leaves a slot into which we can
		// drop the new element (x)
	} else {
		R = append((*r)[:left+1], (*r)[left:]...)
		R[left] = x
	}

	// Verify something was added
	*r = R
	ok = u1+1 == r.ulen()

	return
}

/*
Reset will silently iterate and delete each slice found within
the receiver, leaving it unpopulated but still retaining its
active configuration. Nothing is returned.
*/
func (r Stack) Reset() {
	if !r.IsInit() {
		return
	}
	r.stack.reset()
}

/*
reset is a private method called by Stack.Reset.
*/
func (r *stack) reset() {
	if !r.isInit() {
		return
	}

	if r.positive(ronly) {
		return
	}

	r.lock()
	defer r.unlock()

	for i := r.ulen(); i > 0; i-- {
		r.remove(i - 1)
	}
}

/*
Remove will remove and return the Nth slice from the index,
along with a success-indicative Boolean value. A value of
true indicates the receiver length became shorter by one (1).

Use of the Remove method shall not result in fragmentation of
the stack: gaps resulting from the removal of slice instances
shall immediately be "collapsed" using the subsequent slices
available.
*/
func (r Stack) Remove(idx int) (slice any, ok bool) {
	if !r.IsInit() {
		return
	}

	return r.stack.remove(idx)
}

/*
remove is a private method called by Stack.Remove.
*/
func (r *stack) remove(idx int) (slice any, ok bool) {
	if !r.isInit() {
		return
	}

	fname := fmname()

	if r.positive(ronly) {
		return
	}

	r.lock()
	defer r.unlock()

	slice, index, found := r.index(idx)
	if !found {
		r.debug(sprintf("%s: idx %d not found", fname, idx))
		return
	} else if index == 0 {
		// I'm just paranoid
		return
	}

	// note the len before we start
	var u1 int = r.ulen()
	var contents []any
	var preserved int

	// Gather what we want to keep.
	for i := 1; i < r.len(); i++ {
		if index != i {
			preserved++
			contents = append(contents, (*r)[i])
		}
	}

	// zero out everything except the config slice
	cfg, _ := r.config()
	r.debug(sprintf("%s: ALLOC", fname))

	var R stack = make(stack, 0)
	R = append(R, cfg)

	r.debug(sprintf("%s: adding %d preserved elements", fname, preserved))

	R = append(R, contents...)

	r.debug(sprintf("%s: updating PTR contents [%s]", fname, r.addr()))

	*r = R

	// make sure we succeeded both in non-nilness
	// and in the expected integer length change.
	ok = slice != nil && u1-1 == r.ulen()
	r.debug(sprintf("%s: updated: %t", fname, ok))

	return
}

/*
Paren sets the string-encapsulation bit for parenthetical
expression within the receiver. The receiver shall undergo
parenthetical encapsulation ( (...) ) during the string
representation process. Individual string values shall not
be encapsulated in parenthesis, only the whole (current)
stack.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the encapsulation bit (i.e.: true->false
and false->true)
*/
func (r Stack) Paren(state ...bool) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.stack.setOpt(parens)
		} else {
			r.stack.unsetOpt(parens)
		}
	} else {
		r.stack.toggleOpt(parens)
	}

	return r
}

/*
IsNesting returns a Boolean value indicative of whether
at least one (1) slice member is either a Stack or Stack
type alias. If true, this indicates the relevant slice
descends into another hierarchical (nested) context.
*/
func (r Stack) IsNesting() bool {
	if !r.IsInit() {
		return false
	}

	return r.stack.isNesting()
}

/*
isNesting is a private method called by Stack.IsNesting.

When called, this method returns a Boolean value indicative
of whether the receiver contains one (1) or more slice elements
that match either of the following conditions:

• Slice type is a stackage.Stack native type instance, OR ...

• Slice type is a stackage.Stack type-aliased instance

A return value of true is thrown at the first of either
occurrence. Length of matched candidates is not significant
during the matching process.
*/
func (r stack) isNesting() bool {

	// start iterating at index #1, thereby
	// skipping the configuration slice.
	for i := 1; i < r.len(); i++ {

		// perform a type switch on the
		// current index, thereby allowing
		// evaluation of slice types.
		switch tv := r[i].(type) {

		// native Stack instance
		case Stack:
			return true

		// type alias stack instnaces, since
		// we have no knowledge of them here,
		// will be matched in default using
		// the stackTypeAliasConverter func.
		default:

			// If convertible is true, we know the
			// instance (tv) is a stack alias.
			if _, convertible := stackTypeAliasConverter(tv); convertible {
				return convertible
			}
		}
	}

	return false
}

/*
IsParen returns a Boolean value indicative of whether the
receiver is parenthetical.
*/
func (r Stack) IsParen() bool {
	if !r.IsInit() {
		return false
	}

	return r.stack.positive(parens)
}

/*
Stack will fold the case of logical Boolean operators which
are not represented through symbols. For example, `AND`
becomes `and`, or vice versa. This won't have any effect
on List-based receivers, or if symbols are used in place
of said Boolean words.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the case-folding bit (i.e.: true->false
and false->true)
*/
func (r Stack) Fold(state ...bool) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.stack.setOpt(cfold)
		} else {
			r.stack.unsetOpt(cfold)
		}
	} else {
		r.stack.toggleOpt(cfold)
	}

	return r
}

/*
NegativeIndices will enable negative index support when
using the Index method extended by this type. See the
method documentation for further details.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the negative indices bit (i.e.: true->false
and false->true)
*/
func (r Stack) NegativeIndices(state ...bool) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.stack.setOpt(negidx)
		} else {
			r.stack.unsetOpt(negidx)
		}
	} else {
		r.stack.toggleOpt(negidx)
	}

	return r
}

/*
ForwardIndices will enable forward index support when using
the Index method extended by this type. See the method
documentation for further details.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the forward indices bit (i.e.: true->false
and false->true)
*/
func (r Stack) ForwardIndices(state ...bool) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.stack.setOpt(fwdidx)
		} else {
			r.stack.unsetOpt(fwdidx)
		}
	} else {
		r.stack.toggleOpt(fwdidx)
	}

	return r
}

/*
SetDelimiter accepts input characters (as string, or a single rune) for use
in controlled stack value joining when the underlying stack type is a LIST.
In such a case, the input value shall be used for delimitation of all slice
values during the string representation process.

A zero string, the NTBS (NULL) character -- ASCII #0 -- or nil, shall unset
this value within the receiver.

If this method is executed using any other stack type, the operation has no
effect. If using Boolean AND, OR or NOT stacks and a character delimiter is
preferred over a Boolean WORD, see the Stack.Symbol method.
*/
func (r Stack) SetDelimiter(x any) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	r.stack.setListDelimiter(x)
	return r
}

/*
assertListDelimiter is a private function called (indirectly)
by the Stack.SetDelimiter method for the purpose of handing
type assertion for values that may express a particular value
intended to serve as delimiter character for a LIST Stack when
string representation is requested.
*/
func assertListDelimiter(x any) (v string) {
	if x == nil {
		return
	}

	switch tv := x.(type) {
	case string:
		v = tv
	case rune:
		if tv != rune(0) {
			v = string(tv)
		}
	}

	return
}

/*
Delimiter returns the delimiter string value currently set
within the receiver instance.
*/
func (r Stack) Delimiter() string {
	return r.stack.getListDelimiter()
}

/*
setListDelimiter is a private method called by Stack.SetDelimiter
*/
func (r *stack) setListDelimiter(x any) {
	sc, _ := r.config()
	sc.setListDelimiter(assertListDelimiter(x))
}

/*
getListDelimiter is a private method called by Stack.Delimiter.
*/
func (r *stack) getListDelimiter() string {
	sc, _ := r.config()
	return sc.getListDelimiter()
}

/*
Encap accepts input characters for use in controlled stack value
encapsulation. Acceptable input types are:

• string - a single string value will be used for both L and R
encapsulation.

• string slices - An instance of []string with two (2) values will
be used for L and R encapsulation using the first and second
slice values respectively. An instance of []string with only one (1)
value is identical to providing a single string value, in that both
L and R will use one value.
*/
func (r Stack) Encap(x ...any) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	r.stack.setEncap(x...)
	return r
}

/*
setEncap is a private method called by Stack.Encap.
*/
func (r *stack) setEncap(x ...any) {
	sc, _ := r.config()
	sc.setEncap(x...)
}

/*
IsEncap returns a Boolean value indicative of whether value encapsulation
characters have been set within the receiver.
*/
func (r Stack) IsEncap() bool {
	if !r.IsInit() {
		return false
	}

	return len(r.stack.getEncap()) > 0
}

/*
getEncap returns the current value encapsulation character pattern
set within the receiver instance.
*/
func (r *stack) getEncap() [][]string {
	sc, _ := r.config()
	return sc.enc
}

/*
SetID assigns the provided string value (or lack thereof) to the receiver.
This is optional, and is usually only needed in complex structures in which
"labeling" certain components may be advantageous. It has no effect on an
evaluation, nor should a name ever cause a validity check to fail.

If the string `_random` is provided, a 24-character alphanumeric string is
randomly generated using math/rand and assigned as the ID.
*/
func (r Stack) SetID(id string) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	r.stack.setID(id)
	return r
}

/*
setID is a private method called by Stack.SetID.
*/
func (r *stack) setID(id string) {
	sc, _ := r.config()

	if lc(id) == `_random` {
		id = randomID(randIDSize)
	} else if lc(id) == `_addr` {
		id = sprintf("%s", r.addr())
	}

	sc.setID(id)
}

/*
SetCategory assigns the provided string to the stack's internal category
value. This allows for a means of identifying a particular kind of stack
in the midst of many.
*/
func (r Stack) SetCategory(cat string) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	r.stack.setCat(cat)
	return r
}

/*
setCat is a private method called by Stack.SetCategory.
*/
func (r *stack) setCat(cat string) {
	sc, _ := r.config()
	sc.setCat(cat)
}

/*
Category returns the categorical label string value assigned to the receiver,
if set, else a zero string.
*/
func (r Stack) Category() string {
	if !r.IsInit() {
		return ``
	}

	return r.stack.getCat()
}

/*
getCat is a private method called by Stack.Category.
*/
func (r *stack) getCat() string {
	if !r.isInit() {
		return ``
	}

	sc, _ := r.config()
	return sc.cat
}

/*
ID returns the assigned identifier string, if set, from within the underlying
stack configuration.
*/
func (r Stack) ID() string {
	if !r.IsInit() {
		return ``
	}

	return r.stack.getID()
}

/*
getID is a private method called by Stack.ID.
*/
func (r *stack) getID() string {
	if !r.isInit() {
		return ``
	}

	sc, _ := r.config()
	return sc.id
}

/*
LeadOnce sets the lead-once bit within the receiver. This
causes two (2) things to happen:

• Only use the configured operator once in a stack, and ...

• Only use said operator at the very beginning of the
stack string value.

Execution without a Boolean input value will *TOGGLE* the
current state of the lead-once bit (i.e.: true->false and
false->true)
*/
func (r Stack) LeadOnce(state ...bool) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.stack.setOpt(lonce)
		} else {
			r.stack.unsetOpt(lonce)
		}
	} else {
		r.stack.toggleOpt(lonce)
	}

	return r
}

/*
NoPadding sets the no-space-padding bit within the receiver.
String values within the receiver shall not be padded using
a single space character (ASCII #32).

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the padding bit (i.e.: true->false and
false->true)
*/
func (r Stack) NoPadding(state ...bool) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.stack.setOpt(nspad)
		} else {
			r.stack.unsetOpt(nspad)
		}
	} else {
		r.stack.toggleOpt(nspad)
	}

	return r
}

/*
NoNesting sets the no-nesting bit within the receiver. If
set to true, the receiver shall ignore any Stack or Stack
type alias instance when pushed using the Push method. In
such a case, only primitives, Conditions, etc., shall be
honored during the Push operation.

Note this will only have an effect when not using a custom
PushPolicy. When using a custom PushPolicy, the user has
total control -- and full responsibility -- in deciding
what may or may not be pushed.

Also note that setting or unsetting this bit shall not, in
any way, have an impact on pre-existing Stack or Stack type
alias instances within the receiver. This bit only has an
influence on the Push method and only when set to true.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the nesting bit (i.e.: true->false and
false->true)
*/
func (r Stack) NoNesting(state ...bool) Stack {
	if !r.IsInit() {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.stack.setOpt(nnest)
		} else {
			r.stack.unsetOpt(nnest)
		}
	} else {
		r.stack.toggleOpt(nnest)
	}

	return r
}

/*
CanNest returns a Boolean value indicative of whether
the no-nesting bit is unset, thereby allowing the Push
of Stack and/or Stack type alias instances.

See also the IsNesting method.
*/
func (r Stack) CanNest() bool {
	if !r.IsInit() {
		return false
	}

	return !r.stack.positive(nnest)
}

/*
IsPadded returns a Boolean value indicative of whether the
receiver pads its contents with a SPACE char (ASCII #32).
*/
func (r Stack) IsPadded() bool {
	if !r.IsInit() {
		return false
	}

	return r.stack.positive(parens)
}

/*
ReadOnly sets the receiver bit 'ronly' to a positive state.
This will prevent any writes to the receiver or its underlying
configuration.
*/
func (r Stack) ReadOnly(state ...bool) Stack {
	if !r.IsInit() {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.stack.setOpt(ronly)
		} else {
			r.stack.unsetOpt(ronly)
		}
	} else {
		r.stack.toggleOpt(ronly)
	}

	return r
}

/*
IsReadOnly returns a Boolean value indicative of whether the
receiver is set as read-only.
*/
func (r Stack) IsReadOnly() bool {
	if !r.IsInit() {
		return false
	}

	return r.stack.positive(ronly)
}

/*
Symbol sets the provided symbol expression, which will be
a sequence of any characters desired, to represent various
Boolean operators without relying on words such as "AND".
If a non-zero sequence of characters is set, they will be
used to supplant the default word-based operators within
the given stack in which the symbol is configured.

Acceptable input types are string and rune.

Execution of this method with no arguments empty the symbol
store within the receiver, thereby returning to the default
word-based behavior.

This method has no effect on list-style stacks.
*/
func (r Stack) Symbol(c ...any) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	if len(c) == 0 {
		r.stack.setSymbol(``)
		return r
	}

	var str string
	for i := 0; i < len(c); i++ {
		switch tv := c[i].(type) {
		case string:
			str += tv
		case rune:
			char := string(tv)
			str += char
		}
	}

	r.stack.setSymbol(str)
	return r
}

/*
setSymbol is a private method called by Stack.Symbol.
*/
func (r *stack) setSymbol(sym string) {
	sc, err := r.config()
	if err != nil {
		return
	}
	if sc.typ != list {
		sc.setSymbol(sym)
	}
}

/*
getSymbol returns the symbol stored within the underlying *nodeConfig
instance slice.
*/
func (r stack) getSymbol() string {
	sc, err := r.config()
	if err != nil {
		return ``
	}

	return sc.sym
}

func (r *stack) setLogger(logger any) {
	cfg, err := r.config()
	if err != nil {
		return
	}

	r.lock()
	defer r.unlock()

	cfg.setLogger(resolveLogger(logger))
	r.debug(sprintf("%s: set %T", fmname(), logger))
}

func (r *stack) logger() *log.Logger {
	cfg, _ := r.config()
	return cfg.getLogger()
}

/*
toggleOpt returns the receiver value in fluent-style after a
locked modification of the underlying *nodeConfig instance
to TOGGLE the state of a particular option.
*/
func (r *stack) toggleOpt(x cfgFlag) *stack {
	cfg, err := r.config()
	if err != nil {
		return r
	}

	r.lock()
	defer r.unlock()

	cfg.toggleOpt(x)
	return r
}

/*
setOpt returns the receiver value in fluent-style after a
locked modification of the underlying *nodeConfig instance
to SET a particular option.
*/
func (r *stack) setOpt(x cfgFlag) *stack {
	cfg, err := r.config()
	if err != nil {
		return r
	}

	r.lock()
	defer r.unlock()

	cfg.setOpt(x)
	return r
}

/*
unsetOpt returns the receiver value in fluent-style after
a locked modification of the underlying *nodeConfig instance
to UNSET a particular option.
*/
func (r *stack) unsetOpt(x cfgFlag) *stack {
	cfg, err := r.config()
	if err != nil {
		return r
	}

	r.lock()
	defer r.unlock()

	cfg.unsetOpt(x)
	return r
}

/*
positive returns a Boolean value indicative of whether the specified
cfgFlag input value is "on" within the receiver's configuration
value.
*/
func (r stack) positive(x cfgFlag) bool {
	fname := fmname()

	cfg, err := r.config()
	if err != nil {
		return false
	}

	result := cfg.positive(x)
	r.debug(sprintf("%s: %s:%t", fname, x, result))

	return result
}

/*
Mutex enables the receiver's mutual exclusion locking capabilities.
*/
func (r Stack) Mutex() Stack {
	if !r.IsInit() {
		return r
	}

	r.stack.setMutex()
	return r
}

/*
setMutex is a private method called by Stack.Mutex.
*/
func (r *stack) setMutex() {
	fname := fmname()

	sc, err := r.config()
	if err != nil {
		return
	}

	r.debug(sprintf("%s: MuTeX implemented", fname))
	sc.setMutex()
}

/*
CanMutex returns a Boolean value indicating whether the receiver
instance has been equipped with mutual exclusion locking features.

This does NOT indicate whether the receiver is actually locked.
*/
func (r Stack) CanMutex() bool {
	if !r.IsInit() {
		return false
	}

	return r.stack.canMutex()
}

/*
canMutex is a private method called by Stack.CanMutex.
*/
func (r stack) canMutex() bool {
	sc, err := r.config()
	if err != nil {
		return false
	}

	result := sc.mtx != nil
	r.debug(sprintf("%s:%t", fmname(), result))

	return sc.mtx != nil
}

/*
lock will attempt to lock the receiver using sync.Mutex. If already
locked, the operation will block. If sync.Mutex was not enabled for
the receiver, nothing happens.
*/
func (r *stack) lock() {
	if !r.canMutex() {
		return
	}

	fname := fmname()

	if mutex, found := r.mutex(); found {
		sc, _ := r.config()
		_now := now()
		sc.ldr = &_now
		r.debug(sprintf("%s: LOCKING", fname))
		mutex.Lock()
	}
}

/*
unlock will attempt to unlock the receiver using sync.Mutex. If not
already locked, the operation will block. If sync.Mutex was not enabled for
the receiver, nothing happens.
*/
func (r *stack) unlock() {
	if !r.canMutex() {
		return
	}

	fname := fmname()

	if mutex, found := r.mutex(); found {
		mutex.Unlock()
		sc, _ := r.config()
		_then := *sc.ldr
		sc.ldr = nil
		r.debug(sprintf("%s: UNLOCKING; duration:%.09f sec",
			fname, now().Sub(_then).Seconds()))
	}
}

/*
config returns the *nodeConfig instance found within the receiver
alongside an instance of error. If either the *stackInstance (sc)
is nil OR if the error (err) is non-nil, the receiver is deemed
totally invalid and unusable.
*/
func (r stack) config() (sc *nodeConfig, err error) {
	if &r == nil || r.len() == 0 {
		err = errorf("%T instance is nil; aborting", r)
		r.error(err)
		return
	}

	var ok bool
	if sc, ok = r[0].(*nodeConfig); !ok {
		r.error(err)
		err = errorf("%T does not contain an expected %T instance; aborting", r, sc)
		return
	}

	err = sc.valid()
	return
}

/*
len returns the true integer length of the receiver instance.

See also stack.ulen().
*/
func (r stack) len() int {
	return len(r)
}

/*
Len returns the integer length of the receiver.
*/
func (r Stack) Len() int {
	if !r.IsInit() {
		return 0
	}

	return r.ulen()
}

/*
cat is a private method called during various string representation
processes.
*/
func (r stack) cap() int {
	sc, _ := r.config()
	return sc.cap
}

func (r Stack) Cap() int {
	if !r.IsInit() {
		return 0
	}

	return r.cap()
}

/*
Avail returns the available number of slices as an integer value by
subtracting the current length from a non-zero capacity.

If no capacity is set, this method returns minus one (-1), meaning
infinite capacity is available.

If the receiver (r) is uninitialized, zero (0) is returned.
*/
func (r Stack) Avail() int {
	if !r.IsInit() {
		return 0
	}

	return r.stack.avail()
}

func (r stack) avail() int {
	if r.cap() == 0 {
		return -1 // no cap set means "infinite capacity"
	}
	return r.cap() - r.len()
}

/*
ulen returns the practical integer length of the receiver. This does
*NOT* include slice zero (0), which is reserved for the configuration
instance, therefore this method always returns len()-1 for stacks with
a true length that is greater than or equal to two (<= 2), and returns
zero (0) for a true length of zero (0) OR one (1).
*/
func (r stack) ulen() (l int) {
	switch u := r.len(); u {
	case 0, 1:
		return 0
	default:
		l = u - 1
	}

	return
}

/*
Kind returns the string name of the type of receiver configuration.
*/
func (r Stack) Kind() string {
	if !r.IsInit() {
		return badStack
	}

	switch t, c := r.stack.typ(); c {
	case and, or, not, list, basic:
		return t
	}

	return badStack
}

/*
String is a stringer method that returns the string representation
of the receiver.

Note that invalid Stack instances, as well as basic Stacks, are not
eligible for string representation.
*/
func (r Stack) String() string {
	if !r.IsInit() {
		return ``
	}
	return r.stack.string()
}

/*
string is a private method called by Stack.String.
*/
func (r stack) string() string {
	// run validity checks, whether package default
	// or a user-authored validity policy ...
	if err := r.valid(); err != nil {
		return badStack
	}

	var str []string
	ot, oc := r.typ()
	if oc == 0x0 {
		return badStack
	} else if oc == basic {
		return ``
	}

	// execute the user-authoried presentation
	// policy, if defined, instead of going any
	// further.
	if ppol := r.getPresentationPolicy(); ppol != nil {
		return ppol(r)
	}

	// Scan each slice and attempt stringification
	for i := 1; i < r.len(); i++ {
		// Handle slice value types through assertion
		if val := r.stringAssertion(r[i]); len(val) > 0 {
			// Append the apparently valid
			// string value ...
			str = append(str, val)
		}
	}

	// hand off our string slices, along with the outermost
	// type/code values, to the assembleStringStack worker.
	return r.assembleStringStack(str, ot, oc)
}

/*
stringAssertion is the private method called during slice iteration
during stack.string runs.
*/
func (r stack) stringAssertion(x any) (value string) {
	switch tv := x.(type) {
	case string:
		// Slice is a raw string value (which
		// may be eligible for encapsulation)
		value = r.encapv(tv)
	default:
		// Catch-all; call defaultAssertionHandler
		// with the current interface slice as the
		// input argument (tv).
		value = r.defaultAssertionHandler(tv)
	}

	return
}

/*
defaultAssertionHandler is a private method called by stack.string.
*/
func (r stack) defaultAssertionHandler(x any) (str string) {
	// Whatever it is, it had better be one (1) of the following:
	// • A Stack, or type alias of Stack, or ...
	// • A Condition, or type alias of Condition, or ...
	// • Something that has its own "stringer" (String()) method.
	if Xs, ok := stackTypeAliasConverter(x); ok {
		ik, ic := Xs.stack.typ() // make note of inner stack type
		if ic == not && len(Xs.getSymbol()) == 0 {
			// Handle NOTs a little differently
			// when nested and when not using
			// symbol operators ...
			ik = foldValue(Xs.positive(cfold), ik)
			str = sprintf("%s %s", ik, Xs.String())
		} else {
			str = Xs.String()
		}
	} else if Xc, ok := conditionTypeAliasConverter(x); ok {
		str = Xc.String()
	} else if meth := getStringer(x); meth != nil {
		// whatever it is, it seems to have
		// a stringer method, at least ...
		str = padValue(!r.positive(nspad), r.encapv(meth()))
	}

	r.debug(sprintf("%s: %T produced: %s", fmname(), x, str))

	return
}

/*
assembleStringStack is a private method called by stack.string. This method
reduces the cyclomatic complexity of stack.string() by handling the end-stage
processing of a request for string representation of the receiver.
*/
func (r stack) assembleStringStack(str []string, ot string, oc stackType) string {

	// Padding char (or lack thereof)
	pad := padValue(!r.positive(nspad), ``)
	fname := fmname()

	var fstr []string
	if r.positive(lonce) {
		// We're here because Lead-Once was requested,
		// so just place the type symbol/word at the
		// beginning and don't use it as a join value.
		if oc != list {
			r.debug(sprintf("%s: [%s] += '%s'", fname, oc, ot))
			fstr = append(fstr, ot)
		}

		// Append previous content as-is.
		fstr = append(fstr, str...)
		r.debug(sprintf("%s: [%s] += '%v' [PREV]", fname, oc, str))
	} else {
		// We're here because the user wants a symbol
		// or word to appear between every stack val
		// OR because the user is stringing a List.
		if oc == list {
			// Since we're dealing with a simple
			// list-style stack, use pad char as
			// the join value.
			if ljc := r.getListDelimiter(); len(ljc) > 0 {
				r.debug(sprintf("%s: [%s] += '%s'", fname, oc, ljc))
				fstr = append(fstr, join(str, ljc))
			} else {
				r.debug(sprintf("%s: [%s] += '%s' [PADDING]", fname, oc, pad))
				fstr = append(fstr, join(str, pad))
			}
		} else {
			// Use the outerType as the join
			// char (e.g.: '&&', '||', 'AND',
			// et al).
			var tjn string
			if len(r.getSymbol()) > 0 {
				char := string(rune(32)) // by default, use WHSP padding for symbol ops
				if r.positive(nspad) {
					char = `` // ... unless user turns it off
				}

				sympad := padValue(!r.positive(nspad), char)
				tjn = join(str, sprintf("%s%s%s", sympad, ot, sympad))
			} else {
				tjn = join(str, ot)
			}
			r.debug(sprintf("%s: [%s] += '%s'", fname, oc, tjn))
			fstr = append(fstr, tjn)
		}
	}

	// Finally, join the completed slices using the
	// pad char, enclose in parenthesis (maybe), and
	// condense any consecutive WHSP/HTAB chars down
	// to one (1) WHSP char as needed.
	fpad := sprintf("%s%s%s", pad, join(fstr, pad), pad)
	result := condenseWHSP(r.paren(fpad))
	r.debug(sprintf("%s: result: '%s'", fname, result))
	return result
}

/*
Traverse will "walk" a structure of stack elements using the path indices
provided. It returns the slice found at the final index, or nil, along with
a success-indicative Boolean value.

The semantics of "traversability" are as follows:

• Any "nesting" instance must be a Stack or Stack type alias

• Condition instances must either be the final requested element, OR must
contain a Stack or Stack type alias instance through which to continue the
traversal process

• All other value types are returned as-is

If the traversal ended at any given value, it will be returned along with a
positive ok value letting the user know they arrived at the coordinates they
defined and that "something" was found.

If, however, any path elements remained and further traversal was NOT possible,
the last slice is returned but ok is not set positive, thereby letting the user
know they took a wrong turn somewhere.

As the return type is any, the slice value must be manually type asserted.
*/
func (r Stack) Traverse(indices ...int) (slice any, ok bool) {
	if !r.IsInit() {
		return nil, false
	}

	slice, ok, _ = r.stack.traverse(indices...)
	return
}

/*
traverse is a private method called by Stack.Traverse.
*/
func (r stack) traverse(indices ...int) (slice any, ok, done bool) {
	fname := fmname()
	if err := r.valid(); err != nil {
		r.error(sprintf("%s: %v", fname, err))
		return
	} else if len(indices) == 0 {
		r.error(sprintf("%s: non-traversable %T [NO_PATH]", fname, indices))
		return
	}

	// begin "walking" path of int breadcrumbs ...
	for i := 0; i < len(indices); i++ {

		current := indices[i] // user-facing index number w/ offset
		r.debug(sprintf("%s: idx:%d %v", fname, current, indices))

		instance, _, found := r.index(current) // don't expose the TRUE index
		if !found {
			r.debug(sprintf("%s: idx:%d not found", fname, current))
			break
		}

		id := getLogID(r.getID())

		// Begin assertion of possible traversable and non-traversable
		// values. We'll go as deep as possible, provided each nesting
		// instance is a Stack/Stack alias, or Condition/Condition alias
		// containing a Stack/Stack alias value.
		r.debug(sprintf("%s: descending into idx:%d %T::%s", fname, current, instance, id))

		if slice, ok, done = r.traverseAssertionHandler(instance, i, indices...); done {
			r.debug(sprintf("%s: returned from idx:%d %T::%s", fname, current, instance, id))
			break
		}
	}

	return
}

/*
Reveal processes the receiver instance and disenvelops needlessly
enveloped Stack slices.
*/
func (r Stack) Reveal() Stack {
	_ = r.stack.reveal()
	return r
}

/*
reveal is a private method called by Stack.Reveal.
*/
func (r *stack) reveal() (err error) {
	if !r.isInit() {
		// receiver not initialized
		return
	}

	if r.positive(ronly) {
		return // no error needed
	}

	r.lock()
	defer r.unlock()

	fname := fmname()

	// scan each slice (except the config
	// slice) and analyze its structure.
	for i := 0; i < r.len(); i++ {
		sl, _, _ := r.index(i) // cfg offset handled by index method, be honest
		if sl == nil {
			continue
		}

		// If the element is a stack, begin descent
		// through recursion.
		if outer, ook := stackTypeAliasConverter(sl); ook {
			id := getLogID(outer.getID())
			r.debug(sprintf("%s: descending into idx:%d %T::%s", fname, i, outer, id))
			if err = r.revealDescend(outer, i); err != nil {
				r.setErr(err)
				r.error(sprintf("%s: %T::%s %v", fname, outer, id, err))
				return
			}
		}
	}

	return
}

/*
revealDescent is a private method called by stack.reveal. It engages the second
level of processing of the provided inner Stack instance.

Its first order of business is to determine the length of the inner Stack instance
and take action based on the result:

- A length of zero (0) results in the deletion of the provided inner stack (at the
provided slice index (idx)) from the receiver instance.

- A length of one (1) results in the recursive call of stack.revealSingle at index
zero (0) of the provided inner stack.

- Any other length (>1) results in a new top-level recursive call of the top-level
stack.reveal private method, beginning the entire process anew at index zero(0) of
the provided inner stack.
*/
func (r *stack) revealDescend(inner Stack, idx int) (err error) {
	var updated any

	fname := fmname()

	id := getLogID(r.getID())

	// TODO: Still mulling the best way
	// to handle NOTs.
	if inner.stackType() != not {
		switch inner.Len() {
		case 0:
			r.remove(idx) // empty stack? remove it from outer (rcvr)
		case 1:
			// descend into inner slice #0
			child, _, _ := inner.index(0)
			cid := getLogID(``)
			assert, ok := child.(Interface)
			if ok {
				cid = getLogID(assert.ID())
			}

			r.debug(sprintf("%s: descending into single idx:%d %T::%s", fname, 0, child, cid))
			if err = r.revealSingle(0); err != nil {
				r.error(sprintf("%s: %T::%s %v", fname, child, cid, err))
				break
			}
			updated = child
		default:
			// begin new top-level reveal of inner
			// as a whole, scanning all +2 slices
			r.debug(sprintf("%s: descending into %T::%s", fname, 0, inner, id))
			if err = inner.reveal(); err != nil {
				r.error(sprintf("%s: %T::%s %v", fname, inner, id, err))
				break
			}
			updated = inner
		}
	}

	// If we have an updated reference
	// in-hand, replace whatever was
	// already present at index idx
	// within the receiver instance.
	if updated != nil {
		r.replace(updated, idx)
	}

	// Begin second pass-over before
	// return.
	err = inner.reveal()
	r.error(sprintf("%s: %T::%s %v", fname, inner, id, err))

	return
}

/*
revealSingle is a private method called by stack.revealDescend. It engages the third
level of processing of the receiver's slice instance found at the provided index (idx).
An error is returned at the conclusion of processing.

If a single Condition or Condition type alias instance (whether pointer-based or not)
is found at the slice index indicated, its Expression value is checked for the presence
of a Stack or Stack type alias. If one is present, a new top-level stack.reveal recursion
is launched.

If, on the other hand, a Stack or Stack type alias (again, pointer or not) is found at
the slice index indicates, a new top-level stack.reveal recursion is launched into it
directly.
*/
func (r *stack) revealSingle(idx int) (err error) {
	// get slice @ idx, bail if nil ...
	slice, _, _ := r.index(idx)
	if slice == nil {
		return
	}

	// If a condition ...
	if c, okc := conditionTypeAliasConverter(slice); okc {
		// ... If condition expression is a stack ...
		if inner, iok := stackTypeAliasConverter(c.Expression()); iok {
			// ... recurse into said stack expression
			if err = inner.reveal(); err == nil {
				// update the condition w/ new value
				c.SetExpression(inner)
			}
		}

	} else if inner, iok := stackTypeAliasConverter(slice); iok {
		// If a stack then recurse
		err = inner.reveal()
	}

	return
}

/*
traverseAssertionHandler handles the type assertion processes during traversal of one (1) or more
nested Stack/Stack alias instances that may or may not reside in Condition/Condition alias instances.
*/
func (r stack) traverseAssertionHandler(x any, idx int, indices ...int) (slice any, ok, done bool) {
	fname := fmname()
	if slice, ok, done = r.traverseStack(x, idx, indices...); ok {
		// The value was a Stack or Stack type alias.
		r.debug(sprintf("%s: idx %d: found traversable %T [SLICE]", fname, idx, x))
	} else if slice, ok, done = r.traverseStackInCondition(x, idx, indices...); ok {
		// The value was a Condition, and there MAY be a Stack
		// or Stack type alias nested within said Condition ...
		r.debug(sprintf("%s: idx %d: found traversable %T [SLICE]", fname, idx, x))
	} else if len(indices) <= 1 {
		// If we're at the end of the line, just return
		// whatever is there.
		slice = x
		ok = true
		done = true
		r.debug(sprintf("%s: idx %d: found %T [SLICE] at end of path", fname, idx, x))
	} else {
		// If we arrived here with more path elements left,
		// it would appear the path was invalid, or ill-suited
		// for this particular structure in the traversable
		// sense. Return the last slice, but don't declare
		// done or ok since it (probably) isn't what they
		// wanted ...
		slice = x
		r.debug(sprintf("%s: path leftovers in non-traversable %T:%d at %v", fname, x, idx, indices))
	}

	return
}

/*
traverseStackInCondition is the private recursive helper method for the stack.traverse method. This
method will traverse either a Stack *OR* Stack alias type fashioned by the user that is the Condition
instance's own value (i.e.: recurse stacks that reside in conditions, etc ...).
*/
func (r stack) traverseStackInCondition(u any, idx int, indices ...int) (slice any, ok, done bool) {
	fname := fmname()
	if c, cOK := conditionTypeAliasConverter(u); cOK {
		// End of the line :)
		if len(indices) <= 1 {
			slice = c
			ok = true
			done = true
			r.debug(sprintf("%s: found %T:%d [SLICE] at %v", fname, c, idx, indices))
		} else {
			// We have leftovers. If the Condition's value is a
			// Stack *OR* a Stack alias, traverse it ...
			expr := c.Expression()
			var id string = getLogID(``)
			if assert, aok := expr.(Interface); aok {
				id = getLogID(assert.ID())
			}
			r.debug(sprintf("%s: descending into %T:%d:%s at %v", fname, expr, idx, id, indices))
			return r.traverseStack(expr, idx, indices...)
		}
	}

	return
}

/*
traverseStack is the private recursive helper method for the stack.traverse method. This
method will traverse either a Stack *OR* Stack alias type fashioned by the user.
*/
func (r stack) traverseStack(u any, idx int, indices ...int) (slice any, ok, done bool) {
	fname := fmname()
	if s, sOK := stackTypeAliasConverter(u); sOK {
		// End of the line :)
		if len(indices) <= 1 {
			slice = s
			ok = slice != nil
			done = true
			r.debug(sprintf("%s: found %T:%d [SLICE] at %v", fname, s, idx, indices))
		} else {
			// begin new Stack (tv/x) recursion beginning at
			// the NEXT index ...
			var id string = getLogID(s.ID())
			r.debug(sprintf("%s: descending into %T:%d:%s at %v", fname, s, idx, id, indices))
			return s.stack.traverse(indices[idx+1:]...)
		}
	}

	return
}

/*
Index returns the Nth slice within the given receiver alongside the
true index number and a Boolean value indicative of a successful call
of a non-nil value.

This method supports the use of the following index values depending
on the configuration of the receiver.

• Negatives: When negative index support is enabled, a negative index
will not panic, rather the index number will be increased such that it
becomes positive and within the bounds of the stack length, and (perhaps
most importantly) aligned with the relative (intended) slice. To offer an
example, -2 would return the second-to-last slice. When negative index
support is NOT enabled, nil is returned for any index out of bounds along
with a Boolean value of false, although no panic will occur.

• Positives: When forward index support is enabled, an index greater than
the length of the stack shall be reduced to the highest valid slice index.
For example, if an index of twenty (20) were used on a stack instance of
a length of ten (10), the index would transform to nine (9). When forward
index support is NOT enabled, nil is returned for any index out of bounds
along with a Boolean value of false, although no panic will occur.

In any scenario, a valid index within the bounds of the stack's length
returns the intended slice along with Boolean value of true.
*/
func (r Stack) Index(idx int) (slice any, ok bool) {
	if !r.IsInit() {
		return nil, false
	}

	slice, _, ok = r.stack.index(idx)
	return
}

/*
index is a private method called by Stack.Index.
*/
func (r stack) index(i int) (slice any, idx int, ok bool) {
	fname := fmname()
	var id string = getLogID(r.getID())

	L := r.ulen()
	if L == 0 {
		r.debug(sprintf("%s: %T:%d:%s; no content", fname, r, i, id))
		return
	}

	if i < 0 {
		if !r.positive(negidx) {
			return
		}
		// We're negative, so let's increase
		// 'idx' to a positive number that
		// reflects the intended slice index
		// value.
		i = factorNegIndex(i, L)
	} else if i > L-1 {
		if !r.positive(fwdidx) {
			return
		}
		// If the user input an index
		// that was always greater than
		// the length of the stack, then
		// return the value for -1 (last
		// slice value).
		i = L
	} else {
		// The index was neither greater
		// than the length of the stack,
		// nor was it a negative index.
		// so just increment by one (1)
		// to account for *nodeConfig
		// offset index.
		i++
	}

	// We're about to find out whether
	// or not the index is really valid
	slice = r[i]
	idx = i
	ok = slice != nil
	r.debug(sprintf("%s: %T:%d:%s; found %T [nil:%t]", fname, r, i, id, slice, !ok))

	return
}

/*
factorNegIndex is run when negative index support
is enabled and a negative index is encountered.
*/
func factorNegIndex(i, l int) int {
	// We're negative, so let's increase
	// 'idx' to a positive number that
	// reflects the intended slice index
	// value.
	if i += (l * 2); i == 0 {
		// if i was increased and
		// landed on zero (0), add
		// one so we don't expose
		// the *nodeConfig slice.
		i++
	} else if i > l-1 {
		// if the index is higher than
		// the length of the stack, we
		// will reduce as needed.
		i = (i - l) + 1
	} else {
		// The index was neither zero (0)
		// nor a number higher than the
		// length of the stack, so just
		// increment by one (1) to account
		// for the *nodeConfig offset.
		i++
	}

	return i
}

/*
paren may (or may not) apply parenthetical encapsulation
to the input/return value, depending on the underlying
receiver configuration. This will only encapsulate the
outside of the stack, not its value(s).
*/
func (r stack) paren(v string) string {
	var pad string = string(rune(32))
	if sc, _ := r.config(); sc.positive(nspad) {
		pad = ``
	}

	if r.positive(parens) && r.stackType() != basic {
		return sprintf("(%s%s%s)", pad, v, pad)
	}
	return v
}

/*
encapv may (or may not) apply character encapsulation to the
input/return value, depending on the underlying receiver cfg.
*/
func (r stack) encapv(v string) string {
	if r.stackType() == basic {
		return v
	}
	return encapValue(r.getEncap(), v)
}

/*
typ returns the string representation of the "kind"
of receiver along with the appropriate uint8 value
representative of said kind.  If the receiver is in
an invalid state, a "null" kind is returned, along
with an implicit null uint8 value (0x0).
*/
func (r stack) typ() (kind string, typ stackType) {
	typ = r.stackType()
	kind = r.kind()
	if sym := r.getSymbol(); len(sym) > 0 {
		kind = sym
	} else if !(typ == list || typ == basic) {
		kind = padValue(true, kind) // TODO: make this better
	}

	return
}

/*
Pop removes and returns the final slice value from
the receiver instance. A Boolean value is returned
alongside, indicative of whether an actual slice
value was found. Note that if the receiver is in an
invalid state, or has a zero length, nothing will
be removed, and a meaningless value of true will be
returned alongside a nil slice value.
*/
func (r Stack) Pop() (any, bool) {
	if !r.IsInit() {
		return nil, false
	}

	return r.stack.pop()
}

/*
pop is a private method called by Stack.Pop.
*/
func (r *stack) pop() (any, bool) {
	fname := fmname()
	id := getLogID(r.getID())

	if r.ulen() < 1 {
		r.error(sprintf("%s: %T:%s; nothing poppable", fname, r, id))
		return nil, true
	}

	if r.positive(ronly) {
		return nil, false
	}

	r.lock()
	defer r.unlock()

	idx := len(*r) - 1
	slice := (*r)[idx]
	r.debug(sprintf("%s: %T:%d:%s; removing idx", fname, r, idx, id))
	*r = (*r)[:idx]
	ok := slice != nil
	r.debug(sprintf("%s: %T:%d:%s; returning [ok:%t]", fname, r, idx, id, ok))
	return slice, slice != nil
}

/*
Push appends the provided value(s) to the receiver,
and returns the receiver in fluent form.

Note that if the receiver is in an invalid state, or
if maximum capacity has been set and reached, each of
the values intended for append shall be ignored.
*/
func (r Stack) Push(y ...any) Stack {
	if !r.IsInit() {
		return r
	}

	r.stack.push(y...)
	return r
}

/*
push is a private method called by Stack.Push.
*/
func (r *stack) push(x ...any) *stack {
	fname := fmname()
	id := getLogID(r.getID())

	if !r.isInit() {
		return r
	}

	if err := r.valid(); err != nil {
		r.error(sprintf("%s: %T:%s [PAYLOAD:%d]: %v", fname, r, id, len(x), err))
		return r
	}

	if r.positive(ronly) {
		return r
	}

	r.lock()
	defer r.unlock()

	// try to see if the user provided a
	// push verification function
	if meth := r.getPushPolicy(); meth != nil {
		r.debug(sprintf("%s: %T::%s; delegating push to %T [PAYLOAD:%d]", fname, r, id, meth, len(x)))
		// use the user-provided function to scan
		// each pushed item for verification.
		r.methodAppend(meth, x...)
		return r
	}

	// no push policy was found, just do it.
	r.debug(sprintf("%s: %T::%s; generic push [PAYLOAD:%d]", fname, r, id, len(x)))
	r.genericAppend(x...)

	return r
}

/*
Defrag scans the receiver for breaks in the contiguity of slices and will collapse their formation
so that they become contiguous. The effective ordering of repositioned slices is preserved.

For example, this:

	"value1", nil, "value2", "value3", nil, nil, nil, "value4"

... would become ...

	"value1", "value2", "value3", "value4"

Such fragmentation is rare and may only occur if the receiver were initially assembled with explicit
nil instances specified as slices. This may also happen if the values were pointer references that
have since been annihilated through some external means. This method exists for these reasons, as well
as other corner-cases currently inconceivable.

The max integer, which defaults to fifty (50) when unset, shall result in the scan being terminated
when the number of nil slices encountered consecutively reaches the maximum. This is to prevent the
process from looping into eternity.
*/
func (r Stack) Defrag(max ...int) Stack {
	var m int = 50
	if !r.IsInit() {
		return r
	}

	if len(max) > 0 {
		if max[0] > 0 {
			m = max[0]
		}
	}

	if m <= 0 {
		return r
	}

	r.stack.defrag(m)
	return r
}

func (r *stack) defrag(max int) {

	fname := fmname()
	id := getLogID(r.getID())

	var start int = -1
	var spat []int = make([]int, r.len(), r.len())
	for i := 0; i < r.len(); i++ {
		if _, _, ok := r.index(i); !ok {
			if start == -1 {
				start = i // only set once
			}
			continue
		}
		spat[i] = 1
	}

	if start == -1 || max <= start {
		r.debug(sprintf("%s: %T::%s; no fragmented data, or max:%d<start:%d", fname, r, id, max, start))
		return
	}

	if r.positive(ronly) {
		return
	}

	tpat := r.implode(start, max, spat)
	last, err := r.verifyImplode(spat, tpat)
	if err != nil {
		r.setErr(err)
	} else if last >= 0 {
		// chop off the remaining consecutive nil slices
		(*r) = (*r)[:last+1]
	}

	return
}

func (r stack) verifyImplode(spat, tpat []int) (last int, err error) {
	last = -1
	fname := fmname()
	id := getLogID(r.getID())

	data := make(map[string]string, len(tpat))
	var fail bool
	for i := 1; i < len(spat); i++ {
		key := sprintf("S[%d]", i-1)
		result := spat[i] == tpat[i]
		if !result {
			fail = true
		}
		if tpat[i] != 0 {
			last = (len(data) + i) - len(tpat)
		}

		data[key] = sprintf("match:%t", result)
	}

	if fail {
		err = errorf("defragmentation failed; inconsistent slice results")
		r.error(sprintf("%s: %T::%s result:%t: %v", fname, r, id, !fail, err), data)
	} else {
		r.debug(sprintf("%s: %T::%s result:%t", fname, r, id, !fail), data)
	}

	last--

	return
}

func (r *stack) implode(start, max int, spat []int) (tpat []int) {
	var ct int
	tpat = make([]int, len(spat), len(spat))
	tpat[0] = 1 // cfg slice is exempt

	r.lock()
	defer r.unlock()

	fname := fmname()

	for {
		r.debug(sprintf("%s: %T implosion target scan at %d", fname, r, start))

		if ct >= max || start+ct >= r.ulen() {
			break
		}

		if (*r)[start+ct+1] == nil {
			ct++
			continue
		}

		r.debug(sprintf("%s: %T shift overwrite %d with %d", fname, r, start+1, start+ct+1))
		(*r)[start+1] = (*r)[start+ct+1]

		tpat[start+ct] = 1

		r.debug(sprintf("%s: %T nullify %d", fname, r, start+ct+1))
		(*r)[start+ct+1] = nil
		start = start + 1
		ct = 0
	}

	r.debug(sprintf("%s: %T target has imploded", fname, r))

	return
}

func (r *stack) canPushNester(x any) (ok bool) {
	if !r.positive(nnest) {
		return true
	}

	_, ok = stackTypeAliasConverter(x)
	return
}

/*
methodAppend is a private method called by stack.push.
*/
func (r *stack) methodAppend(meth PushPolicy, x ...any) *stack {
	fname := fmname()
	id := getLogID(r.getID())

	// use the user-provided function to scan
	// each pushed item for verification.
	capacity := r.cap()
	var pct int
	for i := 0; i < len(x); i++ {
		var err error
		if r.capReached() {
			// if capacity has been reached,
			// break out of loop.
			err = errorf("%T.Push failed: capacity violation (%d/%d slices added)",
				Stack{}, pct, len(x))
			r.setErr(err)

			r.policy(sprintf("%s: +X= %T:%d:%s capacity violation (total:%d > cap:%d)",
				fname, x[i], i, id, r.len()+len(x), capacity))
			break
		}

		if err = meth(x[i]); err == nil {
			r.policy(sprintf("%s: appending %T instance (idx:%d) per %T",
				fname, x[i], i, meth))
			*r = append(*r, x[i])
			pct++
		} else {
			r.setErr(err)
			r.policy(sprintf("%s: appending %T instance (idx:%d) failed per %T: %v",
				fname, x[i], i, meth, err))
		}
	}

	return r
}

/*
CapReached returns a Boolean value indicative of whether the receiver
has reached the maximum configured capacity.
*/
func (r Stack) CapReached() bool {
	if !r.IsInit() {
		return false
	}

	return r.capReached()
}

/*
capReached is a private method called by Stack.CapReached.
*/
func (r stack) capReached() bool {
	if r.cap() == 0 {
		return false
	}
	return r.len() == r.cap()
}

/*
genericAppend performs a normal append operation without the
involvement of a custom push policy. Each iteration shall verify
that maximum capacity --if one was specified-- is not exceeded.
*/
func (r *stack) genericAppend(x ...any) *stack {
	fname := fmname()
	var pct int
	for i := 0; i < len(x); i++ {
		var err error
		if !r.canPushNester(x[i]) {
			err = errorf("%T.Push failed: nesting violation (%d/%d slices added)",
				Stack{}, pct, len(x))
			r.setErr(err)

			r.policy(sprintf("%s: +X= %T [%d]: nesting violation",
				fname, x[i], i))
			continue
		}

		if r.capReached() {
			err = errorf("%T.Push failed: capacity violation (%d/%d slices added)",
				Stack{}, pct, len(x))
			r.setErr(err)

			r.policy(sprintf("%s: +X= %T [%d]: capacity violation (total:%d > cap:%d)",
				fname, x[i], i, r.len()+len(x), r.cap()))
			break
		}

		r.debug(sprintf("%s: += %T [%d]", fname, x[i], i))
		*r = append(*r, x[i])
		pct++
	}

	return r
}

/*
SetPushPolicy assigns the provided PushPolicy closure function
to the receiver, thereby enabling protection against undesired
appends to the Stack. The provided function shall be executed
by the Push method for each individual item being added.
*/
func (r Stack) SetPushPolicy(ppol PushPolicy) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	r.stack.setPushPolicy(ppol)
	return r
}

/*
setPushPolicy is a private method called by Stack.SetPushPolicy.
*/
func (r *stack) setPushPolicy(ppol PushPolicy) *stack {
	fname := fmname()
	rc, err := r.config()
	if err != nil {
		r.setErr(err)
		return r
	}
	r.policy(sprintf("%s: %T registered", fname, ppol))
	rc.ppf = ppol
	return r
}

/*
getPushPolicy is a private method called by Stack.Push.
*/
func (r *stack) getPushPolicy() PushPolicy {
	rc, _ := r.config()
	return rc.ppf
}

/*
SetPresentationPolicy assigns the provided PresentationPolicy
closure function to the receiver, thereby enabling full control
over the stringification of the receiver. Execution of this type's
String() method will execute the provided policy instead of the
package-provided routine.
*/
func (r Stack) SetPresentationPolicy(ppol PresentationPolicy) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	r.stack.setPresentationPolicy(ppol)
	return r
}

/*
setPresentationPolicy is a private method called by Stack.SetPresentationPolicy.
*/
func (r *stack) setPresentationPolicy(ppol PresentationPolicy) *stack {
	fname := fmname()

	if r.stackType() == basic {
		err := errorf("%s: %T incompatible with %T:%s",
			fname, ppol, r, r.stackType())
		r.setErr(err)
		r.policy(err)
		return r
	}

	rc, err := r.config()
	if err != nil {
		r.setErr(err)
		return r
	}

	r.policy(sprintf("%s: %T registered", fname, ppol))
	rc.rpf = ppol
	return r
}

/*
getPresentationPolicy is a private method called by stack.string.
*/
func (r *stack) getPresentationPolicy() PresentationPolicy {
	rc, _ := r.config()
	return rc.rpf
}

/*
SetValidityPolicy assigns the provided ValidityPolicy closure
function instance to the receiver, thereby allowing users to
introduce inline verification checks of a Stack to better
gauge its validity. The provided function shall be executed
by the Valid method.
*/
func (r Stack) SetValidityPolicy(vpol ValidityPolicy) Stack {
	if !r.IsInit() {
		return r
	}

	if r.stack.positive(ronly) {
		return r
	}

	r.stack.setValidityPolicy(vpol)
	return r
}

/*
setValidityPolicy is a private method called by Stack.SetValidityPolicy.
*/
func (r *stack) setValidityPolicy(vpol ValidityPolicy) *stack {
	rc, err := r.config()
	if err != nil {
		return r
	}

	fname := fmname()
	r.policy(sprintf("%s: %T registered", fname, vpol))
	rc.vpf = vpol
	return r
}

/*
getValidityPolicy is a private method called by Stack.Valid.
*/
func (r *stack) getValidityPolicy() ValidityPolicy {
	rc, _ := r.config()
	return rc.vpf
}

func (r *stack) mkmsg(typ string) (Message, bool) {
	if r.isZero() || len(typ) == 0 {
		return Message{}, false
	}

	return Message{
		ID:   getLogID(r.getID()),
		Tag:  typ,
		Addr: r.addr(),
		Type: sprintf("%T", *r),
		Time: timestamp(),
		Len:  r.ulen(),
		Cap:  r.cap(),
	}, true
}

func (r stack) fatal(x any, data ...map[string]string) {
	m, ok := r.mkmsg(`FATAL`)
	if !ok {
		return
	}

	if ok = m.setText(x); !ok {
		return
	}

	if len(data) > 0 {
		if data[0] != nil {
			m.Data = data[0]
		}
	}

	r.logger().Fatalln(m)
}

func (r stack) error(x any, data ...map[string]string) {
	m, ok := r.mkmsg(`ERROR`)
	if !ok {
		return
	}

	if ok = m.setText(x); !ok {
		return
	}

	if len(data) > 0 {
		if data[0] != nil {
			m.Data = data[0]
		}
	}

	r.logger().Println(m)
}

func (r stack) debug(x any, data ...map[string]string) {
	m, ok := r.mkmsg(`DEBUG`)
	if !ok {
		return
	}

	if ok = m.setText(x); !ok {
		return
	}

	if len(data) > 0 {
		if data[0] != nil {
			m.Data = data[0]
		}
	}

	r.logger().Println(m)
}

func (r stack) policy(x any, data ...map[string]string) {
	m, ok := r.mkmsg(`POLICY`)
	if !ok {
		return
	}

	if ok = m.setText(x); !ok {
		return
	}

	if len(data) > 0 {
		if data[0] != nil {
			m.Data = data[0]
		}
	}

	r.logger().Println(m)
}

const badStack = `<invalid_stack>`
