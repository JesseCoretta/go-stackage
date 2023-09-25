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

	cfg.log = newLogSystem(sLogDefault)
	cfg.log.lvl = logLevels(sLogLevelDefault)

	cfg.typ = t
	cfg.ord = fifo
	data[`kind`] = t.String()
	data[`fifo`] = sprintf("%t", fifo)

	if !logDiscard(cfg.log.log) {
		data[`laddr`] = sprintf("%p", cfg.log.log)
		data[`lpfx`] = cfg.log.log.Prefix()
		data[`lflags`] = sprintf("%d", cfg.log.log.Flags())
		data[`lwriter`] = sprintf("%T", cfg.log.log.Writer())
	}

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
	st.trace(`ALLOC`, data)

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
SetLogLevel enables the specified LogLevel instance(s), thereby
instructing the logging subsystem to accept events for submission
and transcription to the underlying logger.

Users may also sum the desired bit values manually, and cast the
product as a LogLevel. For example, if STATE (4), DEBUG (8) and
TRACE (32) logging were desired, entering LogLevel(44) would be
the same as specifying LogLevel3, LogLevel4 and LogLevel6 in
variadic fashion.
*/
func (r Stack) SetLogLevel(l ...any) Stack {
	cfg, _ := r.config()
	r.calls(sprintf("%s: in:%T(%v;len:%d)",
		fmname(), l, l, len(l)))

	cfg.log.shift(l...)

	r.calls(sprintf("%s: out:%T(self)",
		fmname(), r))

	return r
}

/*
LogLevels returns the string representation of a comma-delimited list
of all active LogLevel values within the receiver.
*/
func (r Stack) LogLevels() string {
	cfg, _ := r.config()
	return cfg.log.lvl.String()
}

/*
UnsetLogLevel disables the specified LogLevel instance(s), thereby
instructing the logging subsystem to discard events submitted for
transcription to the underlying logger.
*/
func (r Stack) UnsetLogLevel(l ...any) Stack {
	cfg, _ := r.config()

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%v;len:%d)",
		fname, l, l, len(l)))

	cfg.log.unshift(l...)

	r.calls(sprintf("%s: out:%T(self)",
		fname, r))

	return r
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

	fname := fmname()
	r.calls(sprintf("%s: in: variadic %T(len:%d)",
		fname, aux, len(aux)))

	var _aux Auxiliary
	if len(aux) == 0 {
		r.trace(sprintf("%s: ALLOC %T (no variadic input)",
			fname, _aux))
		_aux = make(Auxiliary, 0)
	} else {
		if aux[0] == nil {
			r.trace(sprintf("%s: ALLOC %T (nil variadic slice)",
				fname, _aux))
			_aux = make(Auxiliary, 0)
		} else {
			r.trace(sprintf("%s: assign user %T(len:%d)",
				fname, _aux, _aux.Len()))
			_aux = aux[0]
		}
	}

	cfg.aux = _aux
	r.debug(sprintf("%s: registered %T(len:%d)",
		fname, cfg.aux, cfg.aux.Len()))
	r.calls(sprintf("%s: out:void", fname))
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

	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))
	aux = sc.aux
	r.debug(sprintf("%s: get %T(len:%d)",
		fname, sc.aux, sc.aux.Len()))
	r.calls(sprintf("%s: out:%T(%d)",
		fname, aux, aux.Len()))

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

	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))
	result := sc.ord
	r.calls(sprintf("%s: out:%T(%v)",
		fmname(), result, result))

	return result
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
	sc, _ := r.config()

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%v)", fname, fifo, fifo))
	if !sc.ord {
		// can only change it once!
		r.state(sprintf("%s: %s state change (%t->%t; ok:%t)",
			fname, `ordering`, sc.ord, fifo, true))
		sc.ord = fifo
	}
	r.calls(sprintf("%s: out:void; result:%t", fname, sc.ord == fifo))

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
	fname := fmname()
	r.calls(sprintf("%s: in:%T(%v)", fname, err, err))
	r.state(sprintf("%s: %T state change (%v->%v; ok:%t)",
		fname, err, sc.err, err, true))
	sc.setErr(err)
	r.calls(sprintf("%s: out:void", fname))
}

/*
getErr returns the instance of error, whether nil or not, from
the underlying receiver configuration.
*/
func (r stack) getErr() (err error) {
	sc, _ := r.config()
	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))
	err = sc.getErr()
	r.calls(sprintf("%s: out:%T(nil:%t)",
		fname, err, err == nil))

	return
}

/*
kind returns the string representation of the kind value
set within the receiver's configuration value.
*/
func (r stack) kind() string {
	sc, _ := r.config()
	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))
	kind := sc.kind()
	r.calls(sprintf("%s: out:%T(%v;len:%d)",
		fname, kind, kind, len(kind)))

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
		fname := fmname()
		r.calls(sprintf("%s: in:niladic", fname))

		// try to see if the user provided a
		// validity function
		stk := Stack{r}
		err = nil
		if meth := stk.getValidityPolicy(); meth != nil {
			r.calls(sprintf("%s: executing validity_policy", fname))
			if err = meth(r); err != nil {
				r.policy(sprintf("%s: error: %v", fname, err))
			}
		}

		r.calls(sprintf("%s: out:error(nil:%t)", fname, err == nil))
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
SetLogger assigns the specified logging facility to the receiver
instance.

Logging is available but is set to discard all events by default.

An active logging subsystem within the receiver supercedes the
default package logger.

The following types/values are permitted:

  - string: `none`, `off`, `null`, `discard` will turn logging off
  - string: `stdout` will set basic STDOUT logging
  - string: `stderr` will set basic STDERR logging
  - int: 0 will turn logging off
  - int: 1 will set basic STDOUT logging
  - int: 2 will set basic STDERR logging
  - *log.Logger: user-defined *log.Logger instance will be set; it should not be nil

Case is not significant in the string matching process.

Logging may also be set globally using the SetDefaultLogger
package level function. Similar semantics apply.
*/
func (r Stack) SetLogger(logger any) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setLogger(logger)
		}
	}
	return r
}

/*
Logger returns the *log.Logger instance. This can be used for quick
access to the log.Logger type's methods in a manner such as:

	r.Logger().Fatalf("We died")

It is not recommended to modify the return instance for the purpose
of disabling logging outright (see Stack.SetLogger method as well
as the SetDefaultStackLogger package-level function for ways of
doing this easily).
*/
func (r Stack) Logger() (l *log.Logger) {
	if r.IsInit() {
		l = r.stack.logger()
	}
	return

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
	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))

	// if a capacity was set, make sure
	// the destination can handle it...
	if dest.cap() > 0 {
		if r.ulen() > dest.cap()-r.ulen() {
			// capacity is in-force, and
			// there are too many slices
			// to xfer.
			err := errorf("%s failed: capacity violation (%d/%d slices added)",
				fname, r.ulen(), dest.cap()-r.ulen())
			r.policy(err)
			return
		}
	}

	// xfer slices, without any regard for
	// nilness. Slice type is not subject
	// to discrimination.
	for i := 0; i < r.ulen(); i++ {
		sl, _, _ := r.index(i) // cfg offset handled by index method
		dest.push(sl)
		r.trace(sprintf("%s: %T:%d:%s :: transferring %T",
			fname, Stack{}, i, getLogID(r.getID()), sl))
	}

	// return result
	ok = dest.ulen() >= r.ulen()
	r.trace(sprintf("%s: %T::%s :: transfer result:%t",
		fname, Stack{}, getLogID(r.getID()), ok))
	r.calls(sprintf("%s: out:%T(%v)", fname, ok, ok))

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
			fname := fmname()
			r.calls(sprintf("%s: in:%T(%t),%T(%v:%d)",
				fname, x, x == nil, i, i, i))

			if ok = i+1 <= r.ulen(); ok {
				(*r)[i+1] = x
			}

			r.calls(sprintf("%s: out:%T(%t)",
				fname, ok, ok))
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
	fname := fmname()
	r.calls(sprintf("%s: in:%T(%t),%T(%v:%d)",
		fname, x, x == nil, left, left, left))

	// note the len before we start
	var u1 int = r.ulen()

	// bail out if a capacity has been set and
	// would be breached by this insertion.
	if u1+1 > r.cap() && r.cap() != 0 {
		err := errorf("%s failed: capacity violation (%d/%d slices added)",
			fname, 0, 1)
		r.policy(err)
		r.calls(sprintf("%s: out:%T(%v)",
			fname, ok, ok))
		return
	}

	r.lock()
	defer r.unlock()

	cfg, _ := r.config()

	// If left is greater-than-or-equal
	// to the user length, just push.
	if u1-1 < left {
		*r = append(*r, x)
		r.trace(sprintf("%s: append %T->%T",
			fname, x, *r))

		// Verify something was added
		ok = u1+1 == r.ulen()
		r.calls(sprintf("%s: out:%T(%v); ok:%t",
			fname, r, r, ok))
		return
	}

	var R stack = make(stack, 0)
	r.trace(sprintf("%s: ALLOC %T", fname, R))
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
		r.trace(sprintf("%s: append %T,%T,%T->%T",
			fname, cfg, x, (*r)[left:]))

		// If left falls within the user length,
		// append all elements up to and NOT
		// including the desired index, and also
		// append everything after desired index.
		// This leaves a slot into which we can
		// drop the new element (x)
	} else {
		R = append((*r)[:left+1], (*r)[left:]...)
		R[left] = x
		r.trace(sprintf("%s: append %T->%T",
			fname, (*r)[left:], R))
	}

	// Verify something was added
	*r = R
	r.trace(sprintf("%s: reset ptr %T(%p)->%T(%p)",
		fname, R, &R, *r, r))

	ok = u1+1 == r.ulen()
	r.calls(sprintf("%s: out:%T(%v); ok:%t",
		fname, r, r, ok))

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
		fname := fmname()

		r.calls(sprintf("%s: in:niladic", fname))

		var ct int = 0
		var plen int = r.ulen()
		for i := r.ulen(); i > 0; i-- {
			ct++
			r.trace(sprintf("%s: erase idx:%d (%d/%d)",
				fname, i-1, ct/r.ulen()))
			r.remove(i - 1)
		}

		r.debug(sprintf("%s: removed %d/%d slices", fname, ct, plen))
		r.calls(sprintf("%s: out:void", fname))
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

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%v)", fname, idx, idx))

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
				r.trace(sprintf("%s: preserving idx:%d (%T)", fname, idx, (*r)[i]))
				preserved++
				contents = append(contents, (*r)[i])
			}
		}

		// zero out everything except the config slice
		cfg, _ := r.config()
		r.trace(sprintf("%s: ALLOC", fname))

		var R stack = make(stack, 0)
		R = append(R, cfg)

		r.debug(sprintf("%s: adding %d preserved elements", fname, preserved))
		R = append(R, contents...)
		r.debug(sprintf("%s: updating PTR contents [%s]", fname, ptrString(r)))

		*r = R

		// make sure we succeeded both in non-nilness
		// and in the expected integer length change.
		ok = slice != nil && u1-1 == r.ulen()
		r.debug(sprintf("%s: updated: %t", fname, ok))
		r.calls(sprintf("%s: out:%T(nil:%t),%T(%t)", fname, slice, slice == nil, ok, ok))
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

	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))

	// start iterating at index #1, thereby
	// skipping the configuration slice.
	for i := 1; i < r.len(); i++ {
		r.trace(sprintf("%s: idx:%d nesting check", fname, i))

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

	r.calls(sprintf("%s: out:%T(%t)", fname, is, is))

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

	fname := fmname()
	r.calls(sprintf("%s: in:%T(nil:%t)", fname, x == nil))
	sc.setListDelimiter(assertListDelimiter(x))
	r.calls(sprintf("%s: out:void", fname))
}

/*
getListDelimiter is a private method called by Stack.Delimiter.
*/
func (r *stack) getListDelimiter() string {
	sc, _ := r.config()

	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))
	delim := sc.getListDelimiter()
	r.calls(sprintf("%s: out:%T(%v;len:%d)",
		fname, delim, delim, len(delim)))
	return delim
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

	fname := fmname()
	r.calls(sprintf("%s: in:variadic any (ct:%d)",
		fname, len(x)))
	sc.setEncap(x...)
	r.calls(sprintf("%s: out:void", fname))
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
	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))
	encs = sc.enc
	r.calls(sprintf("%s: out:%T(len:%d)", fname, encs, len(encs)))

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

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%v;len:%d)",
		fname, id, id, len(id)))

	if lc(id) == `_random` {
		id = randomID(randIDSize)
	} else if lc(id) == `_addr` {
		id = ptrString(r)
	}
	r.debug(sprintf("%s: setID: %v", fname, id))

	r.lock()
	defer r.unlock()

	sc.setID(id)
	r.calls(sprintf("%s: out:void", fname))
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

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%v;len:%d)",
		fname, cat, cat, len(cat)))
	sc.setCat(cat)
	r.calls(sprintf("%s: out:niladic", fname))
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
	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))
	cat = sc.cat
	r.calls(sprintf("%s: out:%T(%v;len:%d)",
		fname, cat, cat, len(cat)))

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

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%v;len:%d)",
		fname, sym, sym, len(sym)))

	if sc.typ != list {
		r.lock()
		defer r.unlock()
		sc.setSymbol(sym)
	}

	r.calls(sprintf("%s: out:niladic", fname))
}

/*
getSymbol returns the symbol stored within the underlying *nodeConfig
instance slice.
*/
func (r stack) getSymbol() (sym string) {
	sc, _ := r.config()
	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))
	sym = sc.sym
	r.calls(sprintf("%s: out:%T(%v;len:%d)",
		fname, sym, sym, len(sym)))
	return sym
}

func (r *stack) setLogger(logger any) {
	cfg, _ := r.config()

	fname := fmname()
	r.calls(sprintf("%s: in:%T(nil:%t)",
		fname, logger, logger == nil))

	r.lock()
	defer r.unlock()

	cfg.log.setLogger(logger)
	r.trace(sprintf("%s: set %T", fmname(), logger))
	r.calls(sprintf("%s: out:void", fname))
}

func (r *stack) logger() *log.Logger {
	cfg, _ := r.config()
	return cfg.log.logger()
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

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%s)",
		fname, x, x))

	r.lock()
	defer r.unlock()

	cur := cfg.positive(x)
	cfg.toggleOpt(x)
	after := cfg.positive(x)
	result := cur != after

	r.state(sprintf("%s: %s state change (%t->%t; ok:%t)",
		fmname(), x, cur, after, result))

	r.calls(sprintf("%s: out:%T(nil:%t)",
		fname, r, r == nil))

	return r
}

/*
setOpt returns the receiver value in fluent-style after a
locked modification of the underlying *nodeConfig instance
to SET a particular option.
*/
func (r *stack) setOpt(x cfgFlag) *stack {
	cfg, _ := r.config()

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%s)",
		fname, x, x))

	r.lock()
	defer r.unlock()

	cur := cfg.positive(x)
	cfg.setOpt(x)
	after := cfg.positive(x)
	result := cur != after

	r.state(sprintf("%s: %s state change (%t->%t; ok:%t)",
		fmname(), x, cur, after, result))

	r.calls(sprintf("%s: out:%T(nil:%t)",
		fname, r, r == nil))

	return r
}

/*
unsetOpt returns the receiver value in fluent-style after
a locked modification of the underlying *nodeConfig instance
to UNSET a particular option.
*/
func (r *stack) unsetOpt(x cfgFlag) *stack {
	cfg, _ := r.config()

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%s)",
		fname, x, x))

	r.lock()
	defer r.unlock()

	cur := cfg.positive(x)
	cfg.unsetOpt(x)
	after := cfg.positive(x)
	result := cur != after

	r.state(sprintf("%s: %s state change (%t->%t; ok:%t)",
		fmname(), x, cur, after, result))

	r.calls(sprintf("%s: out:%T(nil:%t)",
		fname, r, r == nil))

	return r
}

/*
positive returns a Boolean value indicative of whether the specified
cfgFlag input value is "on" within the receiver's configuration
value.
*/
func (r stack) positive(x cfgFlag) bool {
	cfg, _ := r.config()

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%s)",
		fname, x, x))
	result := cfg.positive(x)
	r.state(sprintf("%s: %s; positive:%t",
		fname, x, result))
	r.calls(sprintf("%s: out:%T(%t)",
		fname, result, result))

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
	fname := fmname()

	sc, _ := r.config()

	r.state(sprintf("%s: MuTeX implemented", fname))
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
	r.state(sprintf("%s: MuTeX implemented: %t", fmname(), result))

	return sc.mtx != nil
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
			r.state(sprintf("%s: LOCKING", fmname()))
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
			_then := *sc.ldr
			sc.ldr = nil
			r.state(sprintf("%s: UNLOCKING; duration:%.09f sec",
				fmname(), now().Sub(_then).Seconds()))
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

		r.error(err)
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
		switch _c := r.cap(); _c {
		case 0:
			c = offset // interpret zero as minus 1
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
		fname := fmname()
		r.calls(sprintf("%s: in:niladic", fname))

		var str []string
		// execute the user-authoried presentation
		// policy, if defined, instead of going any
		// further.
		if ppol := r.getPresentationPolicy(); ppol != nil {
			r.policy(sprintf("%s: call %T", fname, ppol))
			assembled = ppol(r)
			r.calls(sprintf("%s: out:%T(%v)",
				fname, assembled, assembled))
			return
		}

		// Scan each slice and attempt stringification
		for i := 1; i < r.len(); i++ {
			r.trace(sprintf("%s: iterate %T:%d:%s (slice %d/%d)",
				fname, (*r)[i], i-1, getLogID(``), i-1, r.len()-1))
			// Handle slice value types through assertion
			if val := r.stringAssertion((*r)[i]); len(val) > 0 {
				// Append the apparently valid
				// string value ...
				r.trace(sprintf("%s: %T:%d:%s += %T(%v;len:%d)",
					fname, (*r)[i], i-1, getLogID(``), val, val, len(val)))
				str = append(str, val)
			}
		}

		// hand off our string slices, along with the outermost
		// type/code values, to the assembleStringStack worker.
		doPad := !r.positive(nspad) && r.getSymbol() == ``
		assembled = r.assembleStringStack(str, padValue(doPad, ot), oc)

		r.calls(sprintf("%s: out:%T(%v;len:%d)",
			fname, assembled, assembled, len(assembled)))
	}

	return
}

/*
stringAssertion is the private method called during slice iteration
during stack.string runs.
*/
func (r stack) stringAssertion(x any) (value string) {
	fname := fmname()

	r.calls(sprintf("%s: in:%T(nil:%t)",
		fname, x, x == nil))

	switch tv := x.(type) {
	case string:
		// Slice is a raw string value (which
		// may be eligible for encapsulation)
		value = r.encapv(tv)
		r.trace(sprintf("%s: assert:%T->%T(%v;len:%d)",
			fname, tv, value, value, len(value)))
	default:
		// Catch-all; call defaultAssertionHandler
		// with the current interface slice as the
		// input argument (tv).
		value = r.defaultAssertionHandler(tv)
		r.trace(sprintf("%s: assert fallback:%T->%T(%v;len:%d)",
			fname, tv, value, value, len(value)))
	}

	r.calls(sprintf("%s: out:%T(%v;len:%d)",
		fname, value, value, len(value)))

	return
}

/*
defaultAssertionHandler is a private method called by stack.string.
*/
func (r stack) defaultAssertionHandler(x any) (str string) {
	fname := fmname()
	r.calls(sprintf("%s: in: %T(nil:%t)", fname, x, x == nil))

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

	r.debug(sprintf("%s: %T produced: %s",
		fmname(), x, str))

	r.calls(sprintf("%s: out: %T(%v;len:%d)",
		fname, str, str, len(str)))

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

	r.calls(sprintf("%s: in:%T(%v;len:%d),%T(%v;len:%d),%T(%s)",
		fmname(),
		str, str, len(str),
		ot, ot, len(ot),
		oc, oc,
		``))

	var fstr []string
	if r.positive(lonce) {
		// We're here because Lead-Once was requested,
		// so just place the type symbol/word at the
		// beginning and don't use it as a join value.
		if oc != list {
			r.trace(sprintf("%s: %s += '%s'", fname, oc, ot))
			fstr = append(fstr, ot)
		}

		// Append previous content as-is.
		fstr = append(fstr, str...)
		r.trace(sprintf("%s: %s += '%v' (PREV)", fname, oc, str))
	} else {
		// We're here because the user wants a symbol
		// or word to appear between every stack val
		// OR because the user is stringing a List.
		if oc == list {
			// Since we're dealing with a simple
			// list-style stack, use pad char as
			// the join value.
			if ljc := r.getListDelimiter(); len(ljc) > 0 {
				r.trace(sprintf("%s: %s += '%s'", fname, oc, ljc))
				fstr = append(fstr, join(str, ljc))
			} else {
				r.trace(sprintf("%s: %s += '%s' pad:%t", fname, oc, pad, !r.positive(nspad)))
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
				tjn = join(str, sprintf("%s%s%s", sympad, ot, sympad))
			} else {
				char = string(rune(32)) // by default, use WHSP padding for symbol ops
				sympad := padValue(true, char)
				tjn = join(str, sprintf("%s%s%s", sympad, ot, sympad))
			}
			r.trace(sprintf("%s: %s += '%s'", fname, oc, tjn))
			fstr = append(fstr, tjn)
		}
	}

	// Finally, join the completed slices using the
	// pad char, enclose in parenthesis (maybe), and
	// condense any consecutive WHSP/HTAB chars down
	// to one (1) WHSP char as needed.
	fpad := sprintf("%s%s%s", pad, join(fstr, pad), pad)
	result := condenseWHSP(r.paren(fpad))
	r.calls(sprintf("%s: out: %T(%v;len:%d)",
		fname, result, result, len(result)))
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
	fname := fmname()
	r.calls(sprintf("%s: in: variadic %T(%v;len:%d)",
		fname, indices, indices, indices))

	if err := r.valid(); err == nil {
		if len(indices) == 0 {
			r.error(sprintf("%s: non-traversable %T (NO_PATH)", fname, indices))
			r.calls(sprintf("%s: out:%T(nil:%t),%T(%t)",
				fname, slice, slice == nil, ok, ok))
			return
		}

		// begin "walking" path of int breadcrumbs ...
		for i := 0; i < len(indices); i++ {

			current := indices[i] // user-facing index number w/ offset
			r.trace(sprintf("%s: iterate idx:%d %v", fname, current, indices))

			if instance, _, found := r.index(current); found {
				id := getLogID(r.getID())

				// Begin assertion of possible traversable and non-traversable
				// values. We'll go as deep as possible, provided each nesting
				// instance is a Stack/Stack alias, or Condition/Condition alias
				// containing a Stack/Stack alias value.
				r.trace(sprintf("%s: attempting descent into idx:%d of %T::%s", fname, current, instance, id))
				if slice, ok, done = r.traverseAssertionHandler(instance, i, indices...); !done {
					continue
				}
			}
			break
		}
	}

	r.calls(sprintf("%s: out:%T(nil:%t),%T(%t)",
		fname, slice, slice == nil, ok, ok))

	return
}

/*
Reveal processes the receiver instance and disenvelops needlessly
enveloped Stack slices.
*/
func (r Stack) Reveal() Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			_ = r.stack.reveal()
		}
	}
	return r
}

/*
reveal is a private method called by Stack.Reveal.
*/
func (r *stack) reveal() (err error) {
	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))

	r.lock()
	defer r.unlock()

	// scan each slice (except the config
	// slice) and analyze its structure.
	for i := 0; i < r.len(); i++ {
		sl, _, _ := r.index(i) // cfg offset handled by index method, be honest
		r.trace(sprintf("%s: iterate idx:%d %T::", fname, i, sl))
		if sl == nil {
			continue
		}

		// If the element is a stack, begin descent
		// through recursion.
		if outer, ook := stackTypeAliasConverter(sl); ook && outer.Len() > 0 {
			id := getLogID(outer.getID())
			r.trace(sprintf("%s: descending into idx:%d %T::%s", fname, i, outer, id))
			if err = r.revealDescend(outer, i); err == nil {
				continue
			}
			r.setErr(err)
			r.error(sprintf("%s: %T::%s %v", fname, outer, id, err))
			break
		}
		r.calls(sprintf("%s: out:%T(nil:%t)", fname, err, err == nil))
	}

	return
}

/*
revealDescend is a private method called by stack.reveal. It engages the second
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
	r.calls(sprintf("%s: in: %T(len:%d),%T(%d)",
		fname, inner, inner.Len(), idx, idx))

	id := getLogID(r.getID())

	// TODO: Still mulling the best way to handle NOTs.
	if inner.stackType() != not {
		switch inner.Len() {
		case 1:
			// descend into inner slice #0
			child, _, _ := inner.index(0)
			cid := getLogID(``)
			assert, ok := child.(Interface)
			if ok {
				cid = getLogID(assert.ID())
			}

			if !assert.IsParen() && !inner.IsParen() {
				r.trace(sprintf("%s: descending into single idx:%d %T::%s",
					fname, 0, child, cid))

				err = r.revealSingle(0)

				r.error(sprintf("%s: %T::%s %v",
					fname, child, cid, err))

				updated = child
			}
		default:
			// begin new top-level reveal of inner
			// as a whole, scanning all +2 slices
			r.trace(sprintf("%s: descending into %T::%s",
				fname, 0, inner, id))

			err = inner.reveal()

			r.error(sprintf("%s: %T::%s %v",
				fname, inner, id, err))
			updated = inner
		}

		if err != nil {
			return
		}
	}

	// If we have an updated reference
	// in-hand, replace whatever was
	// already present at index idx
	// within the receiver instance.
	if updated != nil {
		r.trace(sprintf("%s: replace %T:%d:%s with %T",
			fname, r, idx, getLogID(r.getID()), updated))
		r.replace(updated, idx)
	}

	// Begin second pass-over before
	// return.
	err = inner.reveal()
	r.error(sprintf("%s: %T::%s %v",
		fname, inner, id, err))
	r.calls(sprintf("%s: out:error(nil:%t)",
		fname, err == nil))

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
	fname := fmname()
	r.trace(sprintf("%s: in:%T(%d)",
		fname, idx, idx))

	// get slice @ idx, bail if nil ...
	if slice, _, _ := r.index(idx); slice != nil {
		// If a condition ...
		if c, okc := conditionTypeAliasConverter(slice); okc {
			r.trace(sprintf("%s slice conversion to %T:ok", fname, c))
			// ... If condition expression is a stack ...
			if inner, iok := stackTypeAliasConverter(c.Expression()); iok {
				r.trace(sprintf("%s: slice assertion to %T:ok", fname, inner))
				// ... recurse into said stack expression
				if err = inner.reveal(); err == nil {
					r.trace(sprintf("%s: set expression to %T(len:%d)",
						fname, inner, inner.Len()))
					// update the condition w/ new value
					c.SetExpression(inner)
				}
			}
		} else if inner, iok := stackTypeAliasConverter(slice); iok {
			r.trace(sprintf("%s slice conversion to %T:ok", fname, inner))
			// If a stack then recurse
			err = inner.reveal()
		}

		r.calls(sprintf("%s: out:error(nil:%t)",
			fname, err == nil))
	}

	return
}

/*
traverseAssertionHandler handles the type assertion processes during traversal of one (1) or more
nested Stack/Stack alias instances that may or may not reside in Condition/Condition alias instances.
*/
func (r stack) traverseAssertionHandler(x any, idx int, indices ...int) (slice any, ok, done bool) {
	fname := fmname()
	r.calls(sprintf("%s: in:%T(nil:%t),%T(%d),%T(%v;len:%d)",
		fname, x, x == nil, idx, idx, indices, indices, len(indices)))

	if slice, ok, done = r.traverseStack(x, idx, indices...); ok {
		// The value was a Stack or Stack type alias.
		r.debug(sprintf("%s: idx %d: found traversable slice %T",
			fname, idx, slice))
	} else if slice, ok, done = r.traverseStackInCondition(x, idx, indices...); ok {
		// The value was a Condition, and there MAY be a Stack
		// or Stack type alias nested within said Condition ...
		r.debug(sprintf("%s: idx %d: found traversable slice %T",
			fname, idx, slice))
	} else if len(indices) <= 1 {
		// If we're at the end of the line, just return
		// whatever is there.
		slice = x
		ok = true
		done = true
		r.debug(sprintf("%s: idx %d: found target slice %T at end of path",
			fname, idx, x))
	} else {
		// If we arrived here with more path elements left,
		// it would appear the path was invalid, or ill-suited
		// for this particular structure in the traversable
		// sense. Don't return the last slice, don't declare
		// done or ok since it (probably) isn't what they
		// wanted ...
		slice = nil
		r.error(sprintf("%s: path leftovers in non-traversable %T:%d at %v",
			fname, x, idx, indices))
	}

	r.calls(sprintf("%s: out:%T(nil:%t),%T(ok:%t),%T(done:%t)",
		fname, slice, slice == nil, ok, ok, done, done))

	return
}

/*
traverseStackInCondition is the private recursive helper method for the stack.traverse method. This
method will traverse either a Stack *OR* Stack alias type fashioned by the user that is the Condition
instance's own value (i.e.: recurse stacks that reside in conditions, etc ...).
*/
func (r stack) traverseStackInCondition(u any, idx int, indices ...int) (slice any, ok, done bool) {
	fname := fmname()
	r.calls(sprintf("%s: in:%T(nil:%t),%T(%d),%T(%v;len:%d)",
		fname, u, u == nil, idx, idx, indices, indices, len(indices)))

	if c, cOK := conditionTypeAliasConverter(u); cOK {
		r.trace(sprintf("%s conversion to %T:ok", fname, c))
		// End of the line :)
		var id string = getLogID(``)
		if len(indices) <= 1 {
			slice = c
			ok = true
			done = true
			r.trace(sprintf("%s: found slice %T:%d at path %v",
				fname, c, idx, indices))
		} else {
			// We have leftovers. If the Condition's value is a
			// Stack *OR* a Stack alias, traverse it ...
			expr := c.Expression()
			if assert, aok := expr.(Interface); aok {
				id = getLogID(assert.ID())
			}
			r.trace(sprintf("%s: descending into %T:%d:%s of %T at %v", fname, expr, idx, id, r, indices))
			return r.traverseStack(expr, idx, indices...)
		}
	}

	r.calls(sprintf("%s: out:%T(nil:%t),%T(ok:%t),%T(done:%t)",
		fname, slice, slice == nil, ok, ok, done, done))

	return
}

/*
traverseStack is the private recursive helper method for the stack.traverse method. This
method will traverse either a Stack *OR* Stack alias type fashioned by the user.
*/
func (r stack) traverseStack(u any, idx int, indices ...int) (slice any, ok, done bool) {
	fname := fmname()
	r.calls(sprintf("%s: in:%T(nil:%t),%T(%d),%T(%v;len:%d)",
		fname, u, u == nil, idx, idx, indices, indices, len(indices)))

	if s, sOK := stackTypeAliasConverter(u); sOK {
		s.trace(sprintf("%s conversion to %T:ok", fname, s))
		// End of the line :)
		var id string = getLogID(s.ID())
		if len(indices) <= 1 {
			slice = u
			ok = sOK
			done = true
			s.trace(sprintf("%s: found slice %T:%d:%s of %T at %v", fname, s, idx, id, u, indices))
		} else {
			// begin new Stack (tv/x) recursion beginning at the NEXT index ...
			s.trace(sprintf("%s: descending into %T:%d:%s of %T at %v", fname, s, indices[idx], id, u, indices[idx:]))
			return s.stack.traverse(indices[1:]...)
		}
	}

	r.calls(sprintf("%s: out:%T(nil:%t),%T(ok:%t),%T(done:%t)",
		fname, slice, slice == nil, ok, ok, done, done))

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
	fname := fmname()
	var id string = getLogID(r.getID())
	r.calls(sprintf("%s: in:%T(%d)", fname, i, i))

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

	r.debug(sprintf("%s: %T:%d:%s; found %T [nil:%t]", fname, r, i, id, slice, !ok))
	r.calls(sprintf("%s: out:%T(nil:%t),%T(%d),%T(done:%t)",
		fname, slice, slice == nil, idx, idx, ok, ok))

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
		i++
	} else if i > l-1 {
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
		return sprintf("(%s%s%s)", pad, v, pad)
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
	fname := fmname()
	id := getLogID(r.getID())
	r.calls(sprintf("%s: in:niladic",
		fname))

	if r.ulen() < 1 {
		ok = true
		r.debug(sprintf("%s: %T:%s; nothing poppable", fname, r, id))
		return
	}

	r.lock()
	defer r.unlock()

	var idx, plen int
	var typ string
	plen = r.len() - 1

	if r.isFIFO() {
		typ = `FIFO`
		idx = 1
		slice = (*r)[idx]
		pres := (*r)[idx+1:]
		(*r) = (*r)[:idx]
		*r = append(*r, pres...)
	} else {
		typ = `LIFO`
		idx = len(*r) - 1
		slice = (*r)[idx]
		*r = (*r)[:idx]
	}

	r.debug(sprintf("%s (%s): %T:%d:%s; executing; len:%d",
		fname, typ, r, idx, id, plen))
	ok = slice != nil

	r.debug(sprintf("%s (%s): %T:%d:%s; ok:%t, len:%d",
		fname, typ, r, idx, id, ok, r.ulen()))

	r.calls(sprintf("%s: out: %T(nil:%t),%T(%t)",
		fname, slice, slice == nil, ok, ok))

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
	fname := fmname()
	id := getLogID(r.getID())

	r.calls(sprintf("%s: in: variadic any(len:%d)", fname, len(x)))

	r.lock()
	defer r.unlock()

	// try to see if the user provided a
	// push verification function
	if meth := r.getPushPolicy(); meth != nil {
		r.debug(sprintf("%s: %T::%s; delegating push to %T (payload:%d)", fname, r, id, meth, len(x)))
		// use the user-provided function to scan
		// each pushed item for verification.
		r.methodAppend(meth, x...)
		return
	}

	// no push policy was found, just do it.
	r.debug(sprintf("%s: %T::%s; generic push (payload:%d)", fname, r, id, len(x)))
	r.genericAppend(x...)

	r.calls(sprintf("%s: out %T(self)", fname, r))

	return
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

If run on a Stack or Stack type-alias that is currently in possession of one (1) or more nested Stack
or Stack type-alias instances, Defrag shall hierarchically traverse the structure and process it no
differently than the top-level instance. This applies to such Stack values nested with an instance of
Condition or Condition type-alias as well.

This is potentially a destructive method and is still very much considered EXPERIMENTAL. While all
tests yield expected results, those who use this method are advised to exercise extreme caution. The
most obvious note of caution pertains to the volatility of index numbers, which shall shift according
to the defragmentation's influence on the instance in question.  By necessity, Len return values shall
also change accordingly.
*/
func (r Stack) Defrag(max ...int) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			// to break defrag loop.
			m := calculateDefragMax(max...)

			r.stack.defrag(m) // defrag the stack itself

			// If the receiver instance is judged as nesting, we'll
			// recurse through stack, and defrag any other suitable
			// candidates for the operation. Targets are any Stack
			// or Condition instances, OR their aliased equivalents.
			if r.IsNesting() {
				for i := 0; i < r.Len(); i++ {
					slice, _ := r.Index(i)
					if sub, ok := stackTypeAliasConverter(slice); ok {
						// Instance is Stack/Stack alias
						sub.Defrag(m)
					} else if cub, ok := conditionTypeAliasConverter(slice); ok {
						// Instance is Condition/Condition alias
						if cub.IsNesting() {
							if sub, ok := stackTypeAliasConverter(cub.Expression()); ok {
								// Condition expression contains a Stack/Stack alias
								sub.Defrag(m)
							}
						}
					}
				}
			}
		}
	}

	return r
}

/*
calculateDefragMax is a private function executed exclusively by Stack.Defrag, and
exists simply to keep cyclomatics factors <=9 in the caller.
*/
func calculateDefragMax(max ...int) (m int) {
	var _m int = 50
	m = _m
	if len(max) > 0 {
		if max[0] > 0 {
			m = max[0]
		}
	}

	if m <= 0 {
		m = _m
	}

	return
}

func (r *stack) defrag(max int) {

	fname := fmname()
	r.calls(sprintf("%s: in: %T(%d)",
		fname, max, max))

	var start int = -1
	var spat []int = make([]int, r.len(), r.len())
	for i := 0; i < r.len(); i++ {
		if _, _, ok := r.index(i); !ok {
			if start == -1 {
				r.trace(sprintf("%s: starting nil identified - idx:%d", fname, i))
				start = i // only set once
			}
			continue
		}
		spat[i] = 1
	}

	if !(start == -1 || max <= start) {
		tpat := r.implode(start, max, spat)
		last, err := r.verifyImplode(spat, tpat)
		r.setErr(err)
		if err == nil && last >= 0 {
			// chop off the remaining consecutive nil slices
			r.trace(sprintf("%s: truncating %d/%d slices",
				fname, len((*r)[:last+1]), len(*r)))
			(*r) = (*r)[:last+1]
		}
	}

	r.calls(sprintf("%s: out:void", fname))

	return
}

func (r stack) verifyImplode(spat, tpat []int) (last int, err error) {
	last = -1
	fname := fmname()
	id := getLogID(r.getID())
	r.calls(sprintf("%s: in: %T,%T",
		fname, spat, tpat))

	err = errorf("defragmentation failed; inconsistent slice results")
	data := make(map[string]string, len(tpat))
	var fail bool
	for i := 1; i < len(spat); i++ {
		key := sprintf("S[%d]", i-1)
		result := spat[i] == tpat[i]
		fail = true
		if result {
			fail = false
		}
		if tpat[i] != 0 {
			last = (len(data) + i) - len(tpat)
		}

		data[key] = sprintf("match:%t", result)
	}

	if !fail {
		r.debug(sprintf("%s: %T::%s result:%t", fname, r, id, !fail), data)
		err = nil
	}

	last--

	r.trace(sprintf("%s: data", data))
	r.calls(sprintf("%s: out:%T(%d),error(nil:%t)",
		fname, last, last, err == nil))

	return
}

func (r *stack) implode(start, max int, spat []int) (tpat []int) {
	var ct int
	tpat = make([]int, len(spat), len(spat))
	tpat[0] = 1 // cfg slice is exempt

	fname := fmname()
	r.calls(sprintf("%s: in:%T(%d),%T(%d),%T(%v;len:%d)",
		fname, start, start, max, max, spat, spat, len(spat)))

	r.lock()
	defer r.unlock()

	r.debug(sprintf("%s: begin implosion target scan", fname))

	for {
		r.trace(sprintf("%s: %T implosion target scan at %d", fname, r, ct))

		if ct >= max || start+ct >= r.ulen() {
			break
		}

		if (*r)[start+ct+1] == nil {
			ct++
			continue
		}

		r.trace(sprintf("%s: %T shift overwrite %d with %d", fname, r, start+1, start+ct+1))
		(*r)[start+1] = (*r)[start+ct+1]

		tpat[start+ct] = 1

		r.trace(sprintf("%s: %T nullify %d", fname, r, start+ct+1))
		(*r)[start+ct+1] = nil
		start = start + 1
		ct = 0
	}

	r.debug(sprintf("%s: %T target has imploded", fname, r))
	r.calls(sprintf("%s: out:%T(%v;len:%d)",
		fname, tpat, tpat, len(tpat)))

	return
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
	fname := fmname()
	id := getLogID(r.getID())

	// use the user-provided function to scan
	// each pushed item for verification.
	var pct int
	for i := 0; i < len(x); i++ {
		var err error
		if !r.capReached() {
			if err = meth(x[i]); err != nil {
				r.setErr(err)
				r.policy(sprintf("%s: appending %T:%d:%s instance failed per %T: %v",
					fname, x[i], i, id, meth, err))
				break
			}

			r.policy(sprintf("%s: appending %T instance (idx:%d) per %T",
				fname, x[i], i, meth))
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
	fname := fmname()
	var pct int

	r.debug(sprintf("%s: begin += %d", fname, len(x)))

	for i := 0; i < len(x); i++ {
		r.trace(sprintf("%s: += %T (%d)", fname, x[i], i))
		if r.canPushNester(x[i]) {
			if !r.capReached() {
				r.trace(sprintf("%s: += %T (%d); ok:true", fname, x[i], i))
				*r = append(*r, x[i])
				pct++
			}
		}
	}

	r.debug(sprintf("%s: complete: %d/%d slices added", fname, pct, len(x)))
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
	fname := fmname()
	sc, _ := r.config()
	r.state(sprintf("%s: %T registered", fname, ppol))
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
	fname := fmname()

	if r.stackType() == basic {
		err := errorf("%s: %T incompatible with %T:%s",
			fname, ppol, r, r.stackType())
		r.setErr(err)
		r.policy(err)
		return r
	}

	sc, _ := r.config()
	r.state(sprintf("%s: %T registered", fname, ppol))
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
	fname := fmname()
	r.state(sprintf("%s: %T registered", fname, vpol))
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

/*
mkmsg is the private method called by eventDispatch for the
purpose of Message assembly prior to submission to a logger.
*/
func (r *stack) mkmsg(typ string) (m Message, ok bool) {
	if r.isInit() {
		if len(typ) > 0 {
			m = Message{
				ID:   getLogID(r.getID()),
				Tag:  typ,
				Addr: ptrString(r),
				Type: sprintf("%T", *r),
				Time: timestamp(),
				Len:  r.ulen(),
				Cap:  Stack{r}.Cap(), // user-facing cap is preferred
			}
			ok = true
		}
	}

	return
}

/*
error conditions that are fatal and always serious
*/
func (r *stack) fatal(x any, data ...map[string]string) {
	if r != nil && x != nil {
		r.eventDispatch(x, LogLevel5, `FATAL`, data...)
	}
}

/*
error conditions that are not fatal but potentially serious
*/
func (r *stack) error(x any, data ...map[string]string) {
	if r != nil && x != nil {
		r.eventDispatch(x, LogLevel5, `ERROR`, data...)
	}
}

/*
relatively deep operational details
*/
func (r *stack) debug(x any, data ...map[string]string) {
	if r != nil && x != nil {
		r.eventDispatch(x, LogLevel4, `DEBUG`, data...)
	}
}

/*
extreme depth operational details
*/
func (r *stack) trace(x any, data ...map[string]string) {
	if r != nil && x != nil {
		r.eventDispatch(x, LogLevel6, `TRACE`, data...)
	}
}

/*
policy method operational details, as well as caps, r/o, etc.
*/
func (r *stack) policy(x any, data ...map[string]string) {
	if r != nil && x != nil {
		r.eventDispatch(x, LogLevel2, `POLICY`, data...)
	}
}

/*
calls records in/out signatures and realtime meta-data regarding
individual method runtimes.
*/
func (r *stack) calls(x any, data ...map[string]string) {
	if r != nil && x != nil {
		r.eventDispatch(x, LogLevel1, `CALL`, data...)
	}
}

/*
state records interrogations of, and changes to, the underlying
configuration value.
*/
func (r *stack) state(x any, data ...map[string]string) {
	if r != nil && x != nil {
		r.eventDispatch(x, LogLevel3, `STATE`, data...)
	}
}

/*
eventDispatch is the main dispatcher of events of any severity.
A severity of FATAL (in any case) will result in a logger-driven
call of os.Exit.
*/
func (r stack) eventDispatch(x any, ll LogLevel, severity string, data ...map[string]string) {
	cfg, _ := r.config()
	if !(cfg.log.positive(ll) ||
		eq(severity, `FATAL`) ||
		cfg.log.lvl == logLevels(AllLogLevels)) {
		return
	}

	printers := map[string]func(...any){
		`FATAL`:  r.logger().Fatalln,
		`ERROR`:  r.logger().Println,
		`STATE`:  r.logger().Println,
		`CALL`:   r.logger().Println,
		`DEBUG`:  r.logger().Println,
		`TRACE`:  r.logger().Println,
		`POLICY`: r.logger().Println,
	}

	if m, ok := r.mkmsg(severity); ok {
		if ok = m.setText(x); ok {
			if len(data) > 0 {
				if data[0] != nil {
					m.Data = data[0]
					if _, ok := data[0][`FATAL`]; ok {
						severity = `ERROR`
					}
				}
			}
			printers[severity](m)
		}
	}
}

const badStack = `<invalid_stack>`
