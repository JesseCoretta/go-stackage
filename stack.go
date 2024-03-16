package stackage

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
	return Stack{newStack(list, false, capacity...)}
}

/*
And initializes and returns a new instance of Stack
configured as a Boolean ANDed stack.
*/
func And(capacity ...int) Stack {
	return Stack{newStack(and, false, capacity...)}
}

/*
Or initializes and returns a new instance of Stack
configured as a Boolean ORed stack.
*/
func Or(capacity ...int) Stack {
	return Stack{newStack(or, false, capacity...)}
}

/*
Not initializes and returns a new instance of Stack
configured as a Boolean NOTed stack.
*/
func Not(capacity ...int) Stack {
	return Stack{newStack(not, false, capacity...)}
}

/*
Basic initializes and returns a new instance of Stack, set
for basic operation only.

Please note that instances of this design are not eligible
for string representation, value encaps, delimitation, and
other presentation-related string methods. As such, a zero
string (“) shall be returned should String() be executed.

PresentationPolicy instances cannot be assigned to Stack
instances of this design.
*/
func Basic(capacity ...int) Stack {
	return Stack{newStack(basic, false, capacity...)}
}

/*
newStack initializes a new instance of *stack, configured
with the kind (t) requested by the user. This function
should only be executed when creating new instances of
Stack (or a Stack type alias), which embeds the *stack
type instance.
*/
func newStack(t stackType, fifo bool, c ...int) *stack {
	var (
		data map[string]string = make(map[string]string, 0)
		cfg  *nodeConfig       = new(nodeConfig)
		st   stack
	)

	cfg.typ = t
	cfg.ord = fifo
	data[`kind`] = t.String()
	data[`fifo`] = sprintf("%t", fifo)
	data[`fifo`] = sprintf("%t", fifo)

	if len(c) > 0 {
		if c[0] > 0 {
			cfg.cap = c[0] + 1 // 1 for cfg slice offset
			st = make(stack, 0, cfg.cap)
			data[`capacity`] = sprintf("%d", cfg.cap)
		}
	} else {
		st = make(stack, 0)
	}

	st = append(st, cfg)
	instance := &st

	data[`addr`] = ptrString(instance)

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
	return ptrString(r.stack)
}

func ptrString(x any) (addr string) {
	addr = sprintf(`uninitialized_%T`, x)
	if x != nil {
		addr = sprintf("%p", x)
	}
	return
}

/*
SetAuxiliary assigns aux, as initialized and optionally populated as
needed by the user, to the receiver instance. The aux input value may
be nil.

If no variadic input is provided, the default Auxiliary allocation
shall occur.

Note that this method shall obliterate any instance that may already
be present, regardless of the state of the input value aux.
*/
func (r Stack) SetAuxiliary(aux ...Auxiliary) Stack {
	if r.IsInit() {
		r.stack.setAuxiliary(aux...)
	}
	return r
}

/*
setAuxiliary is a private method called by Stack.SetAuxiliary.
*/
func (r *stack) setAuxiliary(aux ...Auxiliary) {
	cfg, _ := r.config()

	var _aux Auxiliary
	if len(aux) == 0 {
		_aux = make(Auxiliary, 0)
	} else {
		if aux[0] == nil {
			_aux = make(Auxiliary, 0)
		} else {
			_aux = aux[0]
		}
	}

	cfg.aux = _aux
}

/*
Auxiliary returns the instance of Auxiliary from within the receiver.
*/
func (r Stack) Auxiliary() (aux Auxiliary) {
	if r.IsInit() {
		aux = r.stack.auxiliary()
	}
	return
}

/*
auxiliary is a private method called by Stack.Auxiliary.
*/
func (r stack) auxiliary() (aux Auxiliary) {
	sc, _ := r.config()
	aux = sc.aux

	return
}

/*
IsFIFO returns a Boolean value indicative of whether the underlying
receiver instance exhibits First-In-First-Out behavior as it pertains
to the appending and truncation order of the receiver instance.

A value of false implies Last-In-Last-Out behavior, which is the
default ordering scheme imposed upon instances of this type.
*/
func (r Stack) IsFIFO() (is bool) {
	if r.IsInit() {
		is = r.stack.isFIFO()
	}
	return
}

/*
isFIFO is a private method called by the Stack.IsFIFO method, et al.
*/
func (r stack) isFIFO() bool {
	sc, _ := r.config()
	return sc.ord
}

/*
SetFIFO shall assign the bool instance to the underlying receiver
configuration, declaring the nature of the append/truncate scheme
to be honored.

  - A value of true shall impose First-In-First-Out behavior
  - A value of false (the default) shall impose Last-In-First-Out behavior

This setting shall impose no influence on any methods other than
the Pop method. In other words, Push, Defrag, Remove, Replace, et
al., will all operate in the same manner regardless.

Once set to the non-default value of true, this setting cannot be
changed nor toggled ever again for this instance and shall not be
subject to any override controls.

In short, once you go FIFO, you cannot go back.
*/
func (r Stack) SetFIFO(fifo bool) Stack {
	if r.IsInit() {
		r.stack.setFIFO(fifo)
	}
	return r
}

func (r *stack) setFIFO(fifo bool) {
	if sc, _ := r.config(); !sc.ord {
		sc.ord = fifo
	}
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
func (r Stack) Err() (err error) {
	if r.IsInit() {
		err = r.getErr()
	}
	return
}

/*
SetErr sets the underlying error value within the receiver
to the assigned input value err, whether nil or not.

This method may be most valuable to users who have chosen
to extend this type by aliasing, and wish to control the
handling of error conditions in another manner.
*/
func (r Stack) SetErr(err error) Stack {
	if r.IsInit() {
		r.stack.setErr(err)
	}
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
func (r stack) getErr() (err error) {
	sc, _ := r.config()
	err = sc.getErr()

	return
}

/*
kind returns the string representation of the kind value
set within the receiver's configuration value.
*/
func (r stack) kind() string {
	sc, _ := r.config()
	kind := sc.kind()

	return kind
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

	err = r.stack.valid()
	return
}

/*
valid is a private method called by Stack.Valid.
*/
func (r *stack) valid() (err error) {
	err = errorf("stack instance is not initialized")
	if r.isInit() {

		// try to see if the user provided a
		// validity function
		stk := Stack{r}
		err = nil
		if meth := stk.getValidityPolicy(); meth != nil {
			err = meth(r)
		}
	}

	return
}

/*
IsInit returns a Boolean value indicative of whether the
receiver has been initialized using any of the following
package-level functions:

  - And
  - Or
  - Not
  - List
  - Basic

This method does not take into account the presence (or
absence) of any user-provided values (e.g.: a length of
zero (0)) can still return true.
*/
func (r Stack) IsInit() (is bool) {
	if !r.IsZero() {
		is = r.stack.isInit()
	}
	return
}

/*
isInit is a private method called by Stack.IsInit.
*/
func (r *stack) isInit() bool {
	//var err error
	_, err := r.config()
	return err == nil
}

/*
Transfer will iterate the receiver (r) and add all slices
contained therein to the destination instance (dest).

The following circumstances will result in a false return:

  - Capacity constraints are in-force within the destination instance, and the transfer request (if larger than the sum number of available slices) cannot proceed as a result
  - The destination instance is nil, or has not been properly initialized
  - The receiver instance (r) contains no slices to transfer

The receiver instance (r) is not modified in any way as a
result of calling this method. If the receiver (source) should
undergo a call to its Reset() method following a call to the
Transfer method, only the source will be emptied, and of the
slices that have since been transferred instance shall remain
in the destination instance.
*/
func (r Stack) Transfer(dest Stack) (ok bool) {
	if r.Len() > 0 && dest.IsInit() {
		ok = r.transfer(dest.stack)
	}
	return
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
func (r *stack) transfer(dest *stack) (ok bool) {

	// if a capacity was set, make sure
	// the destination can handle it...
	if dest.cap() > 0 {
		if r.ulen() > dest.cap()-r.ulen() {
			return
		}
	}

	// xfer slices, without any regard for
	// nilness. Slice type is not subject
	// to discrimination.
	for i := 0; i < r.ulen(); i++ {
		sl, _, _ := r.index(i) // cfg offset handled by index method
		dest.push(sl)
	}

	// return result
	ok = dest.ulen() >= r.ulen()

	return
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
	if r != nil {
		if !r.positive(ronly) {
			if ok = i+1 <= r.ulen(); ok {
				(*r)[i+1] = x
			}
		}
	}

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
func (r Stack) Insert(x any, left int) (ok bool) {
	if r.IsInit() && x != nil {
		if !r.getState(ronly) {
			ok = r.stack.insert(x, left)
		}
	}
	return
}

/*
insert is a private method called by Stack.Insert.
*/
func (r *stack) insert(x any, left int) (ok bool) {

	// note the len before we start
	var u1 int = r.ulen()

	// bail out if a capacity has been set and
	// would be breached by this insertion.
	if u1+1 > r.cap()-1 && r.cap() != 0 {
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
	if r.IsInit() {
		r.stack.reset()
	}
}

/*
reset is a private method called by Stack.Reset.
*/
func (r *stack) reset() {
	if !r.positive(ronly) {

		var ct int = 0
		for i := r.ulen(); i > 0; i-- {
			ct++
			r.remove(i - 1)
		}
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
	if r.IsInit() {
		if !r.getState(ronly) {
			slice, ok = r.stack.remove(idx)
		}
	}
	return
}

/*
remove is a private method called by Stack.Remove.
*/
func (r *stack) remove(idx int) (slice any, ok bool) {

	r.lock()
	defer r.unlock()

	var found bool
	var index int
	if slice, index, found = r.index(idx); found {
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

		var R stack = make(stack, 0)
		R = append(R, cfg)
		R = append(R, contents...)

		*r = R

		// make sure we succeeded both in non-nilness
		// and in the expected integer length change.
		ok = slice != nil && u1-1 == r.ulen()
	}

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
	r.setState(parens, state...)
	return r
}

/*
IsNesting returns a Boolean value indicative of whether
at least one (1) slice member is either a Stack or Stack
type alias. If true, this indicates the relevant slice
descends into another hierarchical (nested) context.
*/
func (r Stack) IsNesting() (is bool) {
	if r.IsInit() {
		is = r.stack.isNesting()
	}
	return

}

/*
isNesting is a private method called by Stack.IsNesting.

When called, this method returns a Boolean value indicative
of whether the receiver contains one (1) or more slice elements
that match either of the following conditions:

  - Slice type is a stackage.Stack native type instance, OR ...
  - Slice type is a stackage.Stack type-aliased instance

A return value of true is thrown at the first of either
occurrence. Length of matched candidates is not significant
during the matching process.
*/
func (r stack) isNesting() (is bool) {


	// start iterating at index #1, thereby
	// skipping the configuration slice.
	for i := 1; i < r.len(); i++ {

		// perform a type switch on the
		// current index, thereby allowing
		// evaluation of slice types.
		switch tv := r[i].(type) {

		// native Stack instance
		case Stack:
			is = true

		// type alias stack instnaces, since
		// we have no knowledge of them here,
		// will be matched in default using
		// the stackTypeAliasConverter func.
		default:

			// If convertible is true, we know the
			// instance (tv) is a stack alias.
			_, is = stackTypeAliasConverter(tv)
		}

		if is {
			break
		}
	}

	return
}

/*
IsParen returns a Boolean value indicative of whether the
receiver is parenthetical.
*/
func (r Stack) IsParen() bool {
	return r.getState(parens)
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
	r.setState(cfold, state...)
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
	r.setState(negidx, state...)
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
	r.setState(fwdidx, state...)
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
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setListDelimiter(x)
		}
	}

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
	if x != nil {
		switch tv := x.(type) {
		case string:
			v = tv
		case rune:
			if tv != rune(0) {
				v = string(tv)
			}
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
encapsulation.

A single string value will be used for both L and R encapsulation.

An instance of []string with two (2) values will be used for L and R
encapsulation using the first and second slice values respectively.

An instance of []string with only one (1) value is identical to the
act of providing a single string value, in that both L and R will use
one (1) value.
*/
func (r Stack) Encap(x ...any) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setEncap(x...)
		}
	}

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
func (r Stack) IsEncap() (is bool) {
	if r.IsInit() {
		is = len(r.stack.getEncap()) > 0
	}
	return
}

/*
getEncap returns the current value encapsulation character pattern
set within the receiver instance.
*/
func (r *stack) getEncap() (encs [][]string) {
	sc, _ := r.config()
	encs = sc.enc

	return
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
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setID(id)
		}
	}

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
		id = ptrString(r)
	}

	r.lock()
	defer r.unlock()
	sc.setID(id)
}

/*
SetCategory assigns the provided string to the stack's internal category
value. This allows for a means of identifying a particular kind of stack
in the midst of many.
*/
func (r Stack) SetCategory(cat string) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setCat(cat)
		}
	}
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
func (r Stack) Category() (cat string) {
	if r.IsInit() {
		cat = r.stack.getCat()
	}
	return
}

/*
getCat is a private method called by Stack.Category.
*/
func (r *stack) getCat() (cat string) {
	sc, _ := r.config()
	cat = sc.cat

	return
}

/*
ID returns the assigned identifier string, if set, from within the underlying
stack configuration.
*/
func (r Stack) ID() (id string) {
	id = sprintf("uninitialized_%T", r)
	if r.IsInit() {
		id = r.stack.getID()
	}
	return
}

/*
getID is a private method called by Stack.ID.
*/
func (r *stack) getID() string {
	sc, _ := r.config()
	return sc.id
}

/*
LeadOnce sets the lead-once bit within the receiver. This
causes two (2) things to happen:

  - Only use the configured operator once in a stack, and ...
  - Only use said operator at the very beginning of the stack string value

Execution without a Boolean input value will *TOGGLE* the
current state of the lead-once bit (i.e.: true->false and
false->true)
*/
func (r Stack) LeadOnce(state ...bool) Stack {
	r.setState(lonce, state...)
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
	r.setState(nspad, state...)
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
	r.setState(nnest, state...)
	return r
}

/*
CanNest returns a Boolean value indicative of whether
the no-nesting bit is unset, thereby allowing the Push
of Stack and/or Stack type alias instances.

See also the IsNesting method.
*/
func (r Stack) CanNest() bool {
	return r.getState(nnest)
}

/*
IsPadded returns a Boolean value indicative of whether the
receiver pads its contents with a SPACE char (ASCII #32).
*/
func (r Stack) IsPadded() bool {
	return !r.getState(nspad)
}

/*
ReadOnly sets the receiver bit 'ronly' to a positive state.
This will prevent any writes to the receiver or its underlying
configuration.
*/
func (r Stack) ReadOnly(state ...bool) Stack {
	r.setState(ronly, state...)
	return r
}

/*
IsReadOnly returns a Boolean value indicative of whether the
receiver is set as read-only.
*/
func (r Stack) IsReadOnly() bool {
	return r.getState(ronly)
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
	if r.IsInit() {
		if !r.getState(ronly) {
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
		}
	}

	return r
}

/*
setSymbol is a private method called by Stack.Symbol.
*/
func (r *stack) setSymbol(sym string) {
	sc, _ := r.config()
	if sc.typ != list {
		r.lock()
		defer r.unlock()
		sc.setSymbol(sym)
	}
}

/*
getSymbol returns the symbol stored within the underlying *nodeConfig
instance slice.
*/
func (r stack) getSymbol() (sym string) {
	sc, _ := r.config()
	sym = sc.sym
	return
}

func (r Stack) getState(cf cfgFlag) (state bool) {
	if r.IsInit() {
		state = r.stack.positive(cf)
	}
	return
}

func (r Stack) setState(cf cfgFlag, state ...bool) {
	if r.IsInit() {
		if !r.getState(ronly) || cf == ronly {
			if len(state) > 0 {
				if state[0] {
					r.stack.setOpt(cf)
				} else {
					r.stack.unsetOpt(cf)
				}
			} else {
				r.stack.toggleOpt(cf)
			}
		}
	}
}

/*
toggleOpt returns the receiver value in fluent-style after a
locked modification of the underlying *nodeConfig instance
to TOGGLE the state of a particular option.
*/
func (r *stack) toggleOpt(x cfgFlag) *stack {
	cfg, _ := r.config()

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
	cfg, _ := r.config()

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
	cfg, _ := r.config()

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
	cfg, _ := r.config()
	result := cfg.positive(x)
	return result
}

/*
Mutex enables the receiver's mutual exclusion locking capabilities.

Subsequent calls of write-related methods, such as Push, Pop, Remove
and others, shall invoke MuTeX locking at the latest possible state
of processing, thereby minimizing the duration of a lock as much as
possible.
*/
func (r Stack) Mutex() Stack {
	if r.IsInit() {
		r.stack.setMutex()
	}
	return r
}

/*
setMutex is a private method called by Stack.Mutex.
*/
func (r *stack) setMutex() {
	sc, _ := r.config()
	sc.setMutex()
}

/*
CanMutex returns a Boolean value indicating whether the receiver
instance has been equipped with mutual exclusion locking features.

This does NOT indicate whether the receiver is actually locked.
*/
func (r Stack) CanMutex() (can bool) {
	if r.IsInit() {
		can = r.stack.canMutex()
	}
	return
}

/*
canMutex is a private method called by Stack.CanMutex.
*/
func (r stack) canMutex() bool {
	sc, _ := r.config()
	result := sc.mtx != nil
	return result
}

/*
lock will attempt to lock the receiver using sync.Mutex. If already
locked, the operation will block. If sync.Mutex was not enabled for
the receiver, nothing happens.
*/
func (r *stack) lock() {
	if r.canMutex() {
		if mutex, found := r.mutex(); found {
			sc, _ := r.config()
			_now := now()
			sc.ldr = &_now
			mutex.Lock()
		}
	}
}

/*
unlock will attempt to unlock the receiver using sync.Mutex. If not
already locked, the operation will block. If sync.Mutex was not enabled for
the receiver, nothing happens.
*/
func (r *stack) unlock() {
	if r.canMutex() {
		if mutex, found := r.mutex(); found {
			mutex.Unlock()
			sc, _ := r.config()
			sc.ldr = nil
		}
	}
}

/*
config returns the *nodeConfig instance found within the receiver
alongside an instance of error. If either the *stackInstance (sc)
is nil OR if the error (err) is non-nil, the receiver is deemed
totally invalid and unusable.
*/
func (r *stack) config() (sc *nodeConfig, err error) {
	if r != nil {
		var ok bool
		// verify slice #0 is a *nodeConfig
		// instance, or bail out.
		err = errorf("%T does not contain an expected %T instance; aborting", r, sc)
		if sc, ok = (*r)[0].(*nodeConfig); ok {
			err = nil
		}
	}

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
func (r Stack) Len() (i int) {
	if r.IsInit() {
		i = r.ulen()
	}
	return
}

/*
cap is a private method called during a variety of processes,
returning the integer representation of the capacity limit
imposed upon the receiver instance.

Note that this method returns the "true" capacity, which includes
the cfg slice offset value (+1) summed with the user-input capacity
limit figure. The exported Cap method, however, does NOT report this
figure, and only recognizes the user-input value. Don't confuse the
user :)
*/
func (r stack) cap() int {
	sc, _ := r.config()
	return sc.cap
}

/*
Cap returns the integer representation of a capacity limit imposed upon
the receiver. The return values shall be interpreted as follows:

  - A zero (0) value indicates that the receiver has NO capacity, as the instance is not properly initialized
  - A positive non-zero value (e.g.: >=1) reflects the capacity limit imposed upon the receiver instance
  - A minus one (-1) value indicates infinite capacity is available; no limit is imposed
*/
func (r Stack) Cap() (c int) {
	if r.IsInit() {
		offset := -1
		c = offset
		switch _c := r.cap(); _c {
		case 0:
			// interpret zero as minus 1
		default:
			// handle the cfg slice offset here, as
			// the value is +non-zero
			c = _c + offset // cfg.cap minus 1
		}
	}

	return
}

/*
Avail returns the available number of slices as an integer value by
subtracting the current length from a non-zero capacity.

If no capacity is set, this method returns minus one (-1), meaning
infinite capacity is available.

If the receiver (r) is uninitialized, zero (0) is returned.
*/
func (r Stack) Avail() (avail int) {
	if r.IsInit() {
		if avail = -1; r.cap() > 0 {
			avail = r.cap() - r.len()
		}
	}

	return
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
func (r Stack) Kind() (k string) {
	k = badStack
	if r.IsInit() {
		switch t, c := r.stack.typ(); c {
		case and, or, not, list, basic:
			k = t
		}
	}

	return
}

/*
String is a stringer method that returns the string representation
of the receiver.

Note that invalid Stack instances, as well as basic Stacks, are not
eligible for string representation.
*/
func (r Stack) String() (s string) {
	if r.IsInit() {
		s = r.stack.string()
	}
	return
}

func (r *stack) canString() (can bool, ot string, oc stackType) {
	if r != nil {
		if err := r.valid(); err == nil {
			ot, oc = r.typ()
			can = oc != 0x0 && oc != basic
		}
	}
	return
}

/*
string is a private method called by Stack.String.
*/
func (r *stack) string() (assembled string) {
	if can, ot, oc := r.canString(); can {

		var str []string
		// execute the user-authoried presentation
		// policy, if defined, instead of going any
		// further.
		if ppol := r.getPresentationPolicy(); ppol != nil {
			assembled = ppol(r)
			return
		}

		// Scan each slice and attempt stringification
		for i := 1; i < r.len(); i++ {
			// Handle slice value types through assertion
			if val := r.defaultAssertionHandler((*r)[i]); len(val) > 0 {
				str = append(str, val)
			}
		}

		// hand off our string slices, along with the outermost
		// type/code values, to the assembleStringStack worker.
		doPad := !r.positive(nspad) && r.getSymbol() == ``
		assembled = r.assembleStringStack(str, padValue(doPad, ot), oc)
	}

	return
}

/*
defaultAssertionHandler is a private method called by stack.string.
*/
func (r stack) defaultAssertionHandler(x any) (str string) {

	// str is assigned with 
	str = `UNKNOWN`

	// Whatever it is, it had better be one (1) of the following
	// possibilities in the following evaluated order:
	// • An initialized Stack, or type alias of Stack, or ...
	// • An initialized Condition, or type alias of Condition, or ...
	// • Something that has its own "stringer" (String()) method.
	//
	// BUG FIX: shadow the "ok" variable for conversion checks,
	// and use (Stack|Condition).IsInit to decide whether the
	// asserted value is usable. This resolves a panic bug that
	// was reported within go-aci, an application that imports
	// stackage.
	if Xs, _ := stackTypeAliasConverter(x); Xs.IsInit() {
		ik, ic := Xs.stack.typ() // make note of inner stack type
		if ic == not && len(Xs.getSymbol()) == 0 {
			// Handle NOTs a little differently
			// when nested and when not using
			// symbol operators ...
			ik = foldValue(Xs.positive(cfold), ik)
			str = ik + " " + Xs.String()
		} else {
			str = Xs.String()
		}

	} else if Xc, _ := conditionTypeAliasConverter(x); Xc.IsInit() {
		str = Xc.String()

	} else if meth := getStringer(x); meth != nil {
		// whatever it is, it seems to have
		// a stringer method, at least. If the
		// user is submitting a non-primitive
		// like a struct or a map, and NOT a
		// type alias of Stack/Condition, this
		// will be the condition that matches.
		str = padValue(!r.positive(nspad), r.encapv(meth()))
	} else if isKnownPrimitive(x) {
		// If its a Go primitive, string it (see misc.go).
		str = padValue(!r.positive(nspad), r.encapv(primitiveStringer(x)))
	}

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

	var fstr []string
	if r.positive(lonce) {
		// We're here because Lead-Once was requested,
		// so just place the type symbol/word at the
		// beginning and don't use it as a join value.
		if oc != list {
			fstr = append(fstr, ot)
		}

		// Append previous content as-is.
		fstr = append(fstr, str...)
	} else {
		// We're here because the user wants a symbol
		// or word to appear between every stack val
		// OR because the user is stringing a List.
		if oc == list {
			// Since we're dealing with a simple
			// list-style stack, use pad char as
			// the join value.
			if ljc := r.getListDelimiter(); len(ljc) > 0 {
				fstr = append(fstr, join(str, ljc))
			} else {
				fstr = append(fstr, join(str, pad))
			}
		} else {
			// Use the outerType as the join
			// char (e.g.: '&&', '||', 'AND',
			// et al).
			var tjn string
			var char string
			if len(r.getSymbol()) > 0 {
				if !r.positive(nspad) {
					char = string(rune(32))
				}

				sympad := padValue(!r.positive(nspad), char)
				tjn = join(str, sympad + ot + sympad)
			} else {
				char = string(rune(32)) // by default, use WHSP padding for symbol ops
				sympad := padValue(true, char)
				tjn = join(str, sympad + ot + sympad)
			}
			fstr = append(fstr, tjn)
		}
	}

	// Finally, join the completed slices using the
	// pad char, enclose in parenthesis (maybe), and
	// condense any consecutive WHSP/HTAB chars down
	// to one (1) WHSP char as needed.
	fpad := pad + join(fstr,pad) + pad
	result := condenseWHSP(r.paren(fpad))
	return result
}

/*
Traverse will "walk" a structure of stack elements using the path indices
provided. It returns the slice found at the final index, or nil, along with
a success-indicative Boolean value.

The semantics of "traversability" are as follows:

  - Any "nesting" instance must be a Stack or Stack type alias
  - Condition instances must either be the final requested element, OR must contain a Stack or Stack type alias instance through which the traversal process may continue
  - All other value types are returned as-is

If the traversal ended at any given value, it will be returned along with a
positive ok value letting the user know they arrived at the coordinates they
defined and that "something" was found.

If, however, any path elements remained and further traversal was NOT possible,
the last slice is returned as nil.

As the return type is any, the slice value must be manually type asserted.
*/
func (r Stack) Traverse(indices ...int) (slice any, ok bool) {
	if r.IsInit() {
		slice, ok, _ = r.stack.traverse(indices...)
	}
	return
}

/*
traverse is a private method called by Stack.Traverse.
*/
func (r stack) traverse(indices ...int) (slice any, ok, done bool) {

	if err := r.valid(); err == nil {
		if len(indices) == 0 {
			return
		}

		// begin "walking" path of int breadcrumbs ...
		for i := 0; i < len(indices); i++ {

			current := indices[i] // user-facing index number w/ offset

			if instance, _, found := r.index(current); found {

				// Begin assertion of possible traversable and non-traversable
				// values. We'll go as deep as possible, provided each nesting
				// instance is a Stack/Stack alias, or Condition/Condition alias
				// containing a Stack/Stack alias value.
				if slice, ok, done = r.traverseAssertionHandler(instance, i, indices...); !done {
					continue
				}
			}
			break
		}
	}

	return
}

/*
traverseAssertionHandler handles the type assertion processes during traversal of one (1) or more
nested Stack/Stack alias instances that may or may not reside in Condition/Condition alias instances.
*/
func (r stack) traverseAssertionHandler(x any, idx int, indices ...int) (slice any, ok, done bool) {
	if slice, ok, done = r.traverseStack(x, idx, indices...); ok {
		// The value was a Stack or Stack type alias.
	} else if slice, ok, done = r.traverseStackInCondition(x, idx, indices...); ok {
		// The value was a Condition, and there MAY be a Stack
		// or Stack type alias nested within said Condition ...
	} else if len(indices) <= 1 {
		// If we're at the end of the line, just return
		// whatever is there.
		slice = x
		ok = true
		done = true
	} else {
		// If we arrived here with more path elements left,
		// it would appear the path was invalid, or ill-suited
		// for this particular structure in the traversable
		// sense. Don't return the last slice, don't declare
		// done or ok since it (probably) isn't what they
		// wanted ...
		slice = nil
	}

	return
}

/*
traverseStackInCondition is the private recursive helper method for the stack.traverse method. This
method will traverse either a Stack *OR* Stack alias type fashioned by the user that is the Condition
instance's own value (i.e.: recurse stacks that reside in conditions, etc ...).
*/
func (r stack) traverseStackInCondition(u any, idx int, indices ...int) (slice any, ok, done bool) {
	if c, cOK := conditionTypeAliasConverter(u); cOK {
		// End of the line :)
		if len(indices) <= 1 {
			slice = c
			ok = true
			done = true
		} else {
			// We have leftovers. If the Condition's value is a
			// Stack *OR* a Stack alias, traverse it ...
			expr := c.Expression()
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

	if s, sOK := stackTypeAliasConverter(u); sOK {
		// End of the line :)
		if len(indices) <= 1 {
			slice = u
			ok = sOK
			done = true
		} else {
			// begin new Stack (tv/x) recursion beginning at the NEXT index ...
			return s.stack.traverse(indices[1:]...)
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

Negatives: When negative index support is enabled, a negative index
will not panic, rather the index number will be increased such that it
becomes positive and within the bounds of the stack length, and (perhaps
most importantly) aligned with the relative (intended) slice. To offer an
example, -2 would return the second-to-last slice. When negative index
support is NOT enabled, nil is returned for any index out of bounds along
with a Boolean value of false, although no panic will occur.

Positives: When forward index support is enabled, an index greater than
the length of the stack shall be reduced to the highest valid slice index.
For example, if an index of twenty (20) were used on a stack instance of
a length of ten (10), the index would transform to nine (9). When forward
index support is NOT enabled, nil is returned for any index out of bounds
along with a Boolean value of false, although no panic will occur.

In any scenario, a valid index within the bounds of the stack's length
returns the intended slice along with Boolean value of true.
*/
func (r Stack) Index(idx int) (slice any, ok bool) {
	if r.IsInit() {
		slice, _, ok = r.stack.index(idx)
	}
	return
}

/*
index is a private method called by Stack.Index.
*/
func (r stack) index(i int) (slice any, idx int, ok bool) {

	if L := r.ulen(); L > 0 {
		if i < 0 {
			if r.positive(negidx) && i-(i*2) <= L {
				i = factorNegIndex(i, L)
				ok = true
			}
		} else if i > L-1 {
			if r.positive(fwdidx) {
				i = L
				ok = true
			}
		} else {
			i++
			ok = true
		}

		// We're about to find out whether the index valid
		if ok {
			slice = r[i]
			idx = i
			ok = slice != nil
		}
	}

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
	if i += (l * 2); i > l-1 {
		i = (i - l) + 1
	} else {
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
		return `(` + pad + v + pad + `)`
	}
	return v
}

/*
encapv may (or may not) apply character encapsulation to the
input/return value, depending on the underlying receiver cfg.
*/
func (r stack) encapv(v string) (e string) {
	if r.stackType() != basic {
		e = encapValue(r.getEncap(), v)
	}
	return
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
		//} else if !(typ == list || typ == basic) {
		//kind = padValue(true, kind) // TODO: make this better
	}

	return
}

/*
Pop removes and returns the requisite slice value from the receiver
instance. A Boolean value is returned alongside, indicative of whether
an actual slice value was found.

The requisite slice shall be one (1) of the following, depending on
ordering mode in effect:

  - In the default mode -- LIFO -- this shall be the final slice (index "Stack.Len() - 1", or the "far right" element)
  - In the alternative mode -- FIFO -- this shall be the first slice (index 0, or the "far left" element)

Note that if the receiver is in an invalid state, or has a zero length,
nothing will be removed, and a meaningless value of true will be returned
alongside a nil slice value.
*/
func (r Stack) Pop() (popped any, ok bool) {
	if r.IsInit() {
		if !r.getState(ronly) {
			popped, ok = r.stack.pop()
		}
	}
	return
}

/*
pop is a private method called by Stack.Pop.
*/
func (r *stack) pop() (slice any, ok bool) {
	if r.ulen() < 1 {
		ok = true
		return
	}

	r.lock()
	defer r.unlock()

	var idx int

	if r.isFIFO() {
		idx = 1
		slice = (*r)[idx]
		pres := (*r)[idx+1:]
		(*r) = (*r)[:idx]
		*r = append(*r, pres...)
	} else {
		idx = len(*r) - 1
		slice = (*r)[idx]
		*r = (*r)[:idx]
	}

	ok = slice != nil

	return
}

/*
Push appends the provided value(s) to the receiver,
and returns the receiver in fluent form.

Note that if the receiver is in an invalid state, or
if maximum capacity has been set and reached, each of
the values intended for append shall be ignored.
*/
func (r Stack) Push(y ...any) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.push(y...)
		}
	}
	return r
}

/*
push is a private method called by Stack.Push.
*/
func (r *stack) push(x ...any) {

	r.lock()
	defer r.unlock()

	// try to see if the user provided a
	// push verification function
	if meth := r.getPushPolicy(); meth != nil {
		// use the user-provided function to scan
		// each pushed item for verification.
		r.methodAppend(meth, x...)
		return
	}

	// no push policy was found, just do it.
	r.genericAppend(x...)

	return
}

/*
Reverse shall re-order the receiver's current slices in a sequence that is the polar opposite
of the original.
*/
func (r Stack) Reverse() Stack {
	r.stack.reverse()
	return r
}

/*
reverse is a private niladic and void method called exclusively by Stack.Reverse.
*/
func (r *stack) reverse() {
	if r.ulen() > 0 {
		r.lock()
		defer r.unlock()
		for i,j := 1, r.len()-1; i < j; i, j = i+1,j-1 {
			(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
		}
	}
}

func (r *stack) canPushNester(x any) (can bool) {
	_, can = stackTypeAliasConverter(x)
	if !r.positive(nnest) {
		can = true
	}
	return
}

/*
methodAppend is a private method called by stack.push.
*/
func (r *stack) methodAppend(meth PushPolicy, x ...any) *stack {

	// use the user-provided function to scan
	// each pushed item for verification.
	var pct int
	for i := 0; i < len(x); i++ {
		var err error
		if !r.capReached() {
			if err = meth(x[i]); err != nil {
				r.setErr(err)
				break
			}

			*r = append(*r, x[i])
			pct++
		}
	}

	return r
}

/*
CapReached returns a Boolean value indicative of whether the receiver
has reached the maximum configured capacity.
*/
func (r Stack) CapReached() (cr bool) {
	if r.IsInit() {
		cr = r.capReached()
	}
	return

}

/*
capReached is a private method called by Stack.CapReached.
*/
func (r stack) capReached() (rc bool) {
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
func (r *stack) genericAppend(x ...any) {
	var pct int

	for i := 0; i < len(x); i++ {
		if r.canPushNester(x[i]) {
			if !r.capReached() {
				*r = append(*r, x[i])
				pct++
			}
		}
	}
}

/*
SetPushPolicy assigns the provided PushPolicy closure function
to the receiver, thereby enabling protection against undesired
appends to the Stack. The provided function shall be executed
by the Push method for each individual item being added.
*/
func (r Stack) SetPushPolicy(ppol PushPolicy) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setPushPolicy(ppol)
		}
	}
	return r
}

/*
setPushPolicy is a private method called by Stack.SetPushPolicy.
*/
func (r *stack) setPushPolicy(ppol PushPolicy) *stack {
	sc, _ := r.config()
	sc.ppf = ppol
	return r
}

/*
getPushPolicy is a private method called by Stack.Push.
*/
func (r *stack) getPushPolicy() PushPolicy {
	sc, _ := r.config()
	return sc.ppf
}

/*
SetPresentationPolicy assigns the provided PresentationPolicy
closure function to the receiver, thereby enabling full control
over the stringification of the receiver. Execution of this type's
String() method will execute the provided policy instead of the
package-provided routine.
*/
func (r Stack) SetPresentationPolicy(ppol PresentationPolicy) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setPresentationPolicy(ppol)
		}
	}
	return r
}

/*
setPresentationPolicy is a private method called by Stack.SetPresentationPolicy.
*/
func (r *stack) setPresentationPolicy(ppol PresentationPolicy) *stack {

	if r.stackType() == basic {
		return r
	}

	sc, _ := r.config()
	sc.rpf = ppol
	return r
}

/*
getPresentationPolicy is a private method called by stack.string.
*/
func (r *stack) getPresentationPolicy() PresentationPolicy {
	sc, _ := r.config()
	return sc.rpf
}

/*
SetValidityPolicy assigns the provided ValidityPolicy closure
function instance to the receiver, thereby allowing users to
introduce inline verification checks of a Stack to better
gauge its validity. The provided function shall be executed
by the Valid method.
*/
func (r Stack) SetValidityPolicy(vpol ValidityPolicy) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setValidityPolicy(vpol)
		}
	}
	return r
}

/*
setValidityPolicy is a private method called by Stack.SetValidityPolicy.
*/
func (r *stack) setValidityPolicy(vpol ValidityPolicy) *stack {
	sc, _ := r.config()
	sc.vpf = vpol
	return r
}

/*
getValidityPolicy is a private method called by Stack.Valid.
*/
func (r *stack) getValidityPolicy() ValidityPolicy {
	sc, _ := r.config()
	return sc.vpf
}

const badStack = `<invalid_stack>`
