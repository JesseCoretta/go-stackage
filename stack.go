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
is embedded (in pointer form) within instances of [Stack].
*/
type stack []any

/*
List initializes and returns a new instance of [Stack]
configured as a simple list. [Stack] instances of this
design can be delimited using the [Stack.SetDelimiter]
method.
*/
func List(capacity ...int) Stack {
	return Stack{newStack(list, false, capacity...)}
}

/*
And initializes and returns a new instance of [Stack]
configured as a Boolean [And] stack.
*/
func And(capacity ...int) Stack {
	return Stack{newStack(and, false, capacity...)}
}

/*
Or initializes and returns a new instance of [Stack]
configured as a Boolean [Or] stack.
*/
func Or(capacity ...int) Stack {
	return Stack{newStack(or, false, capacity...)}
}

/*
Not initializes and returns a new instance of [Stack]
configured as a Boolean [Not] stack.
*/
func Not(capacity ...int) Stack {
	return Stack{newStack(not, false, capacity...)}
}

/*
Basic initializes and returns a new instance of [Stack], set
for basic operation only.

Please note that instances of this design are not eligible
for string representation, value encaps, delimitation, and
other presentation-related string methods. As such, a zero
string (“) shall be returned should [Stack.String] be executed.

[PresentationPolicy] instances cannot be assigned to [Stack]
instances of this design.
*/
func Basic(capacity ...int) Stack {
	return Stack{newStack(basic, false, capacity...)}
}

/*
newStack initializes a new instance of *stack, configured
with the kind (t) requested by the user. This function
should only be executed when creating new instances of
[Stack](or a [Stack] type alias), which embeds the *stack
type instance.
*/
func newStack(t stackType, fifo bool, c ...int) *stack {
	var (
		cfg *nodeConfig = new(nodeConfig)
		st  stack
	)

	cfg.log = newLogSystem(sLogDefault)
	cfg.log.lvl = logLevels(sLogLevelDefault)

	cfg.typ = t
	cfg.ord = fifo

	if len(c) > 0 {
		if c[0] > 0 {
			cfg.cap = c[0] + 1 // 1 for cfg slice offset
			st = make(stack, 0, cfg.cap)
		}
	} else {
		st = make(stack, 0)
	}

	st = append(st, cfg)
	instance := &st

	return instance
}

/*
IsEmpty returns a Boolean value indicative of a receiver length of zero
(0).  This method wraps a call of [Stack.Len] == 0, and is only present
for compatibility and convenience reasons.
*/
func (r Stack) IsEmpty() bool {
	if r.IsInit() {
		return r.Len() == 0
	}

	return true
}

/*
IsZero returns a Boolean value indicative of whether the receiver is nil,
or uninitialized. Following use of [Stack.Free], this method will return
true.
*/
func (r Stack) IsZero() bool {
	return r.stack.isZero()
}

/*
isZero is a private method called by [Stack.IsZero].
*/
func (r *stack) isZero() bool {
	return r == nil
}

/*
SetLogLevel enables the specified [LogLevel] instance(s), thereby
instructing the logging subsystem to accept events for submission
and transcription to the underlying logger.

Users may also sum the desired bit values manually, and cast the
product as a LogLevel. For example, if STATE (4), DEBUG (8) and
TRACE (32) logging were desired, entering LogLevel(44) would be
the same as specifying LogLevel3, LogLevel4 and LogLevel6 in
variadic fashion.
*/
func (r Stack) SetLogLevel(l ...any) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			cfg, _ := r.config()
			cfg.log.shift(l...)
		}
	}

	return r
}

/*
LogLevels returns the string representation of a comma-delimited list
of all active [LogLevel] values within the receiver.
*/
func (r Stack) LogLevels() string {
	cfg, _ := r.config()
	return cfg.log.lvl.String()
}

/*
UnsetLogLevel disables the specified [LogLevel] instance(s), thereby
instructing the logging subsystem to discard events submitted for
transcription to the underlying logger.
*/
func (r Stack) UnsetLogLevel(l ...any) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			cfg, _ := r.config()
			cfg.log.unshift(l...)
		}
	}

	return r
}

/*
Addr returns the string representation of the pointer address for the
receiver. This may be useful for logging or debugging operations.

Note: this call uses [fmt.Sprintf].
*/
func (r Stack) Addr() string {
	return ptrString(r.stack)
}

func ptrString(x any) (addr string) {
	addr = `uninitialized`
	if x != nil {
		addr = sprintf("%p", x)
	}
	return
}

/*
Less returns a Boolean value indicative of whether the slice at index
i is deemed to be less than the slice at j in the context of ordering.

This method is intended to satisfy the func(int,int) bool signature
required by [sort.Interface].

See also [Stack.SetLessFunc] method for a means of specifying instances
of this function.

If no custom closure was assigned, the package-default closure is used,
which evaluates the string representation of values in order to conduct
alphabetical sorting. This means that both i and j must reference slice
values in one of the following states:

  - Type of slice instance must have its own String method
  - Value is a go primitive, such as a string or int

Equal values return false, as do zero string values.
*/
func (r Stack) Less(i, j int) (less bool) {
	if r.IsInit() {
		cfg, _ := r.stack.config()
		if meth := cfg.lss; meth != nil {
			less = meth(i, j)
		} else {
			less = r.stack.defaultLesser(i, j)
		}
	}

	return
}

func (r stack) defaultLesser(idx1, idx2 int) bool {
	slice1, _, _ := r.index(idx1)
	slice2, _, _ := r.index(idx2)

	var strs []string = make([]string, 2)
	for idx, slice := range []any{
		slice1,
		slice2,
	} {
		if slice != nil {
			var ok bool
			if strs[idx], ok = slice.(string); !ok {
				if meth := getStringer(slice); meth != nil {
					strs[idx] = meth()
				} else if isKnownPrimitive(slice) {
					strs[idx] = primitiveStringer(slice)
				}
			}

			if strs[idx] == `` {
				return false
			}
		}
	}

	switch scmp(strs[0], strs[1]) {
	case -1:
		return true
	}

	return false
}

/*
SetLessFunc assigns the provided closure instance to the receiver instance,
thereby allowing effective use of the [Stack.Less] method.

If a nil value, or no values, are submitted, the package-default sorting
mechanism will take precedence.
*/
func (r Stack) SetLessFunc(function ...LessFunc) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setLessFunc(function...)
		}
	}

	return r
}

func (r *stack) setLessFunc(function ...LessFunc) {
	var funk LessFunc
	if len(function) > 0 {
		if function[0] == nil {
			funk = r.defaultLesser
		} else {
			funk = function[0]
		}
	} else {
		funk = r.defaultLesser
	}

	cfg, _ := r.config()
	cfg.lss = funk

	return
}

/*
Swap implements the func(int,int) signature required by the [sort.Interface]
signature.
*/
func (r Stack) Swap(i, j int) {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.swap(i, j)
		}
	}
}

func (r *stack) swap(i, j int) {
	if ok := i <= r.ulen(); !ok {
		return
	} else if ok = j <= r.ulen(); !ok {
		return
	}

	i++
	j++

	r.lock()
	defer r.unlock()

	(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
}

/*
SetAuxiliary assigns aux, as initialized and optionally populated as
needed by the user, to the receiver instance. The aux input value may
be nil.

If no variadic input is provided, the default [Auxiliary] allocation
shall occur.

Note that this method shall obliterate any instance that may already
be present, regardless of the state of the input value aux.
*/
func (r Stack) SetAuxiliary(aux ...Auxiliary) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setAuxiliary(aux...)
		}
	}
	return r
}

/*
setAuxiliary is a private method called by [Stack.SetAuxiliary].
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
Auxiliary returns the instance of [Auxiliary] from within the receiver.
*/
func (r Stack) Auxiliary() (aux Auxiliary) {
	if r.IsInit() {
		aux = r.stack.auxiliary()
	}
	return
}

/*
auxiliary is a private method called by [Stack.Auxiliary].
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
isFIFO is a private method called by the [Stack.IsFIFO] method, et al.
*/
func (r stack) isFIFO() bool {
	sc, _ := r.config()
	result := sc.ord
	return result
}

/*
SetFIFO shall assign the bool instance to the underlying receiver
configuration, declaring the nature of the append/truncate scheme
to be honored.

  - A value of true shall impose First-In-First-Out behavior
  - A value of false (the default) shall impose Last-In-First-Out behavior

This setting shall impose no influence on any methods other than
the [Stack.Pop] method. In other words, [Stack.Push], [Stack.Defrag],
[Stack.Remove], [Stack.Replace], et al., will all operate in the same
manner regardless.

Once set to the non-default value of true, this setting cannot be
changed nor toggled ever again for this instance and shall not be
subject to any override controls.

In short, once you go FIFO, you cannot go back.
*/
func (r Stack) SetFIFO(fifo bool) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setFIFO(fifo)
		}
	}
	return r
}

func (r *stack) setFIFO(fifo bool) {
	if sc, _ := r.config(); !sc.ord {
		// can only change it once!
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

This may be used regardless of [Stack.IsReadOnly] status.
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
user content operations (see the [Stack.Err] method).
*/
func (r Stack) Valid() (err error) {
	if r.stack == nil {
		err = errorf("embedded instance is nil")
	} else if !r.stack.valid() {
		err = errorf("embedded instance is invalid")
	}

	return
}

/*
valid is a private method called by [Stack.Valid].
*/
func (r *stack) valid() (is bool) {
	if r.isInit() {
		// try to see if the user provided a
		// validity function
		stk := Stack{r}
		if meth := stk.getValidityPolicy(); meth != nil {
			if err := meth(r); err != nil {
				return
			}
		}
		is = true
	}

	return
}

/*
IsInit returns a Boolean value indicative of whether the
receiver has been initialized using any of the following
package-level functions:

  - [And]
  - [Or]
  - [Not]
  - [List]
  - [Basic]

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
isInit is a private method called by [Stack.IsInit].
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
  - *[log.Logger]: user-defined *[log.Logger] instance will be set; it should not be nil

Case is not significant in the string matching process.
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
Transfer will iterate the receiver (r) and add all slices contained
therein to the destination instance (dest), which must be a previously
initialized [Stack] or [Stack]-alias instance, else false is returned.

If capacity constraints are in-force within the destination instance,
and the transfer request cannot proceed due to it being larger than
the sum number of available slices, false is returned.

If [Stack.IsInit] returns false, this method returns false, as there
is nothing to transfer.

The receiver instance (r) is not modified in any way as a result of
calling this method. If the receiver (source) should undergo a call to
its [Stack.Reset] or [Stack.Free] methods following a call to the this
method, only the source will be emptied, and all of the slices that have
since been transferred instance shall remain in the destination instance.

A return value of true indicates a successful transfer.
*/
func (r Stack) Transfer(dest any) (ok bool) {
	if r.IsInit() {
		if s, sok := stackTypeAliasConverter(dest); sok {
			if !s.getState(ronly) {
				ok = r.transfer(s.stack)
			}
		}
	}

	return
}

/*
transfer is a private method executed by the [Stack.Transfer]
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
			// capacity is in-force, and
			// there are too many slices
			// to xfer.
			// err := errorf("failed: capacity violation"))
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
func (r Stack) Replace(x any, idx int) (ok bool) {
	if r.IsInit() && x != nil {
		if !r.getState(ronly) {
			ok = r.stack.replace(x, idx)
		}
	}

	return
}

func (r *stack) replace(x any, i int) (ok bool) {
	if r != nil {
		if ok = i+1 <= r.ulen(); ok {
			(*r)[i+1] = x
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
insert is a private method called by [Stack.Insert].
*/
func (r *stack) insert(x any, left int) (ok bool) {
	// note the len before we start
	var u1 int = r.ulen()

	// bail out if a capacity has been set and
	// would be breached by this insertion.
	if u1+1 > r.cap()-1 && r.cap() != 0 {
		//err := errorf("failed: capacity violation")
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
Free frees the receiver instance entirely, including the underlying
configuration. An error is returned if the instance is read-only or
uninitialized.

See also [Stack.Reset].
*/
func (r *Stack) Free() (err error) {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack = nil
			return
		}
		err = errorf("%T is read-only; cannot free", r)
	}

	return
}

/*
Reset will silently iterate and delete each slice found within
the receiver, leaving it unpopulated but still retaining its
active configuration. Nothing is returned.  No action is taken
if the receiver is empty.

See also [Stack.Free].
*/
func (r Stack) Reset() {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.reset()
		}
	}
}

/*
reset is a private method called by [Stack.Reset].
*/
func (r *stack) reset() {
	var ct int = 0
	for i := r.ulen(); i > 0; i-- {
		ct++
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
available.  No action is taken if the receiver is empty or
read-only.
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
remove is a private method called by [Stack.Remove].
*/
func (r *stack) remove(idx int) (slice any, ok bool) {

	var found bool
	var index int
	if slice, index, found = r.index(idx); found {
		// note the len before we start
		var u1 int = r.ulen()
		var contents []any
		var preserved int

		// zero out everything except the config slice
		cfg, _ := r.config()

		var R stack = make(stack, 0)
		R = append(R, cfg)

		r.lock()
		defer r.unlock()

		// Gather what we want to keep.
		for i := 1; i < r.len(); i++ {
			if index != i {
				preserved++
				contents = append(contents, (*r)[i])
			}
		}

		R = append(R, contents...)

		*r = R

		// make sure we succeeded both in non-nilness
		// and in the expected integer length change.
		ok = slice != nil && u1-1 == r.ulen()
	}

	return
}

/*
SetParen sets the string-encapsulation bit for parenthetical
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
func (r Stack) SetParen(state ...bool) Stack {
	r.setState(parens, state...)
	return r
}

/*
Deprecated: Use [Stack.SetParen].
*/
func (r Stack) Paren(state ...bool) Stack {
	return r.SetParen(state...)
}

/*
IsNesting returns a Boolean value indicative of whether
at least one (1) slice member is either a [Stack] or [Stack]
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
isNesting is a private method called by [Stack.IsNesting].

When called, this method returns a Boolean value indicative
of whether the receiver contains one (1) or more slice elements
that match either of the following conditions:

  - Slice type is a [Stack] native type instance, OR ...
  - Slice type is a [Stack] type-aliased instance

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
SetFold will fold the case of logical Boolean operators which
are not represented through symbols. For example, `AND` becomes
`and`, or vice versa. This won't have any effect on List-based
receivers, or if symbols are used in place of said Boolean words.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the case-folding bit (i.e.: true->false
and false->true)
*/
func (r Stack) SetFold(state ...bool) Stack {
	r.setState(cfold, state...)
	return r
}

/*
Deprecated: Use [Stack.SetFold].
*/
func (r Stack) Fold(state ...bool) Stack {
	return r.SetFold(state...)
}

/*
SetNegativeIndices will enable negative index support when using
the [Stack.Index] method extended by this type. See the method
documentation for further details.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the negative indices bit (i.e.: true->false
and false->true)
*/
func (r Stack) SetNegativeIndices(state ...bool) Stack {
	r.setState(negidx, state...)
	return r
}

/*
Deprecated: Use [Stack.SetForwardIndices].
*/
func (r Stack) NegativeIndices(state ...bool) Stack {
	return r.SetNegativeIndices(state...)
}

/*
SetForwardIndices will enable forward index support when using
the [Stack.Index] method extended by this type. See the method
documentation for further details.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the forward indices bit (i.e.: true->false
and false->true)
*/
func (r Stack) SetForwardIndices(state ...bool) Stack {
	r.setState(fwdidx, state...)
	return r
}

/*
Deprecated: Use [Stack.SetForwardIndices].
*/
func (r Stack) ForwardIndices(state ...bool) Stack {
	return r.SetForwardIndices(state...)
}

/*
SetDelimiter accepts input characters (as string, or a single rune) for use
in controlled [Stack] value joining when the underlying [Stack] type is a LIST.
In such a case, the input value shall be used for delimitation of all slice
values during the string representation process.

A zero string, the NTBS (NULL) character -- ASCII #0 -- or nil, shall unset
this value within the receiver.

If this method is executed using any other stack type, the operation has no
effect. If using Boolean AND, OR or NOT stacks and a character delimiter is
preferred over a Boolean WORD, see the [Stack.Symbol] method.
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
by the [Stack.SetDelimiter] method for the purpose of handing
type assertion for values that may express a particular value
intended to serve as delimiter character for a LIST [Stack] when
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
setListDelimiter is a private method called by [Stack.SetDelimiter]
*/
func (r *stack) setListDelimiter(x any) {
	sc, _ := r.config()
	sc.setListDelimiter(assertListDelimiter(x))
}

/*
getListDelimiter is a private method called by [Stack.Delimiter].
*/
func (r *stack) getListDelimiter() string {
	sc, _ := r.config()
	delim := sc.getListDelimiter()
	return delim
}

/*
SetEncap accepts input characters for use in controlled [Stack] value
encapsulation.

A single string value will be used for both L and R encapsulation.

An instance of []string with two (2) values will be used for L and R
encapsulation using the first and second slice values respectively.

An instance of []string with only one (1) value is identical to the
act of providing a single string value, in that both L and R will use
the one value.
*/
func (r Stack) SetEncap(x ...any) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setEncap(x...)
		}
	}

	return r
}

/*
Deprecated: Use [Stack.SetEncap].
*/
func (r Stack) Encap(x ...any) Stack {
	return r.SetEncap(x...)
}

/*
setEncap is a private method called by [Stack.Encap].
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
setID is a private method called by [Stack.SetID].
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
setCat is a private method called by [Stack.SetCategory].
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
getCat is a private method called by [Stack.Category].
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
	id = "unspecified"
	if r.IsInit() {
		id = r.stack.getID()
	}
	return
}

/*
getID is a private method called by [Stack.ID].
*/
func (r *stack) getID() string {
	sc, _ := r.config()
	return sc.id
}

/*
SetLeadOnce sets the lead-once bit within the receiver. This
causes two (2) things to happen:

  - Only use the configured operator once in a stack, and ...
  - Only use said operator at the very beginning of the stack string value

Execution without a Boolean input value will *TOGGLE* the
current state of the lead-once bit (i.e.: true->false and
false->true)
*/
func (r Stack) SetLeadOnce(state ...bool) Stack {
	r.setState(lonce, state...)
	return r
}

/*
Deprecated: Use [Stack.SetLeadOnce].
*/
func (r Stack) LeadOnce(state ...bool) Stack {
	return r.SetLeadOnce(state...)
}

/*
SetNoPadding sets the no-space-padding bit within the receiver.
String values within the receiver shall not be padded using
a single space character (ASCII #32).

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the padding bit (i.e.: true->false and
false->true)
*/
func (r Stack) SetNoPadding(state ...bool) Stack {
	r.setState(nspad, state...)
	return r
}

/*
Deprecated: Use [Stack.SetNoPadding].
*/
func (r Stack) NoPadding(state ...bool) Stack {
	return r.SetNoPadding(state...)
}

/*
SetNoNesting sets the no-nesting bit within the receiver. If
set to true, the receiver shall ignore any [Stack] or [Stack]
type alias instance when pushed using the [Stack.Push] method.
In such a case, only primitives, [Condition] instances, etc.,
shall be honored during the [Stack.Push] operation.

Note this will only have an effect when not using a custom
[PushPolicy]. When using a custom [PushPolicy], the user has
total control -- and full responsibility -- in deciding
what may or may not be pushed.

Also note that setting or unsetting this bit shall not, in
any way, have an impact on pre-existing [Stack] or [Stack] type
alias instances within the receiver. This bit only has an
influence on the [Stack.Push] method and only when set to true.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the nesting bit (i.e.: true->false and
false->true)
*/
func (r Stack) SetNoNesting(state ...bool) Stack {
	r.setState(nnest, state...)
	return r
}

/*
Deprecated: Use [Stack.SetNoNesting].
*/
func (r Stack) NoNesting(state ...bool) Stack {
	return r.SetNoNesting(state...)
}

/*
CanNest returns a Boolean value indicative of whether
the no-nesting bit is unset, thereby allowing the push
of [Stack] and/or [Stack] type alias instances.

See also the [Stack.IsNesting] method.
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
SetReadOnly sets the receiver bit 'ronly' to a positive state.
This will prevent any writes to the receiver or its underlying
configuration.
*/
func (r Stack) SetReadOnly(state ...bool) Stack {
	r.setState(ronly, state...)
	return r
}

/*
Deprecated: Use [Stack.SetReadOnly].
*/
func (r Stack) ReadOnly(state ...bool) Stack {
	return r.SetReadOnly(state...)
}

/*
IsReadOnly returns a Boolean value indicative of whether the
receiver is set as read-only.
*/
func (r Stack) IsReadOnly() bool {
	return r.getState(ronly)
}

/*
SetSymbol sets the provided symbol expression, which will be a sequence
of any characters desired, to represent various Boolean operators without
relying on words such as "AND". If a non-zero sequence of characters is
set, they will be used to supplant the default word-based operators within
the given stack in which the symbol is configured.

Acceptable input types are string and rune.

Execution of this method with no arguments empty the symbol store within
the receiver, thereby returning to the default word-based behavior.

This method has no effect on list-style [Stack] instances.
*/
func (r Stack) SetSymbol(c ...any) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setSymbol(c...)
		}
	}

	return r
}

/*
Deprecated: Use [Stack.SetSymbol].
*/
func (r Stack) Symbol(c ...any) Stack {
	return r.SetSymbol(c...)
}

/*
setSymbol is a private method called by [Stack.Symbol].
*/
func (r *stack) setSymbol(c ...any) {
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

	if sc, _ := r.config(); sc.typ != list {
		r.lock()
		defer r.unlock()
		sc.setSymbol(str)
	}
}

/*
getSymbol returns the symbol stored within the underlying *nodeConfig
instance slice.
*/
func (r stack) getSymbol() (sym string) {
	sc, _ := r.config()
	sym = sc.sym
	return sym
}

func (r *stack) setLogger(logger any) {
	cfg, _ := r.config()

	r.lock()
	defer r.unlock()

	cfg.log.setLogger(logger)
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
SetMutex enables the receiver's mutual exclusion locking capabilities.

Subsequent calls of write-related methods, such as [Stack.Push],
[Stack.Pop], [Stack.Remove] and others, shall invoke MuTeX locking at
the latest possible state of processing, thereby minimizing the duration
of a lock.
*/
func (r Stack) SetMutex() Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.stack.setMutex()
		}
	}
	return r
}

/*
Deprecated: Use [Stack.SetMutex].
*/
func (r Stack) Mutex() Stack {
	return r.SetMutex()
}

/*
setMutex is a private method called by [Stack.Mutex].
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
canMutex is a private method called by [Stack.CanMutex].
*/
func (r stack) canMutex() bool {
	sc, _ := r.config()
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
		err = unexpectedReceiverState
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
Len returns the integer length or "size" of the receiver.
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
		if _c := r.cap(); _c > 0 {
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
String is a stringer method that returns the string representation of
the receiver.

Note that invalid [Stack] instances, as well as basic [Stack] instances,
are not eligible for string representation.
*/
func (r Stack) String() (s string) {
	if r.IsInit() {
		s = r.stack.string()
	}
	return
}

func (r *stack) canString() (can bool, ot string, oc stackType) {
	if r != nil {
		if r.valid() {
			ot, oc = r.typ()
			can = oc != 0x0 && oc != basic
		}
	}
	return
}

/*
string is a private method called by [Stack.String].
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
				//if val := r.stringAssertion((*r)[i]); len(val) > 0 {
				// Append the apparently valid
				// string value ...
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
defaultAssertionHandler is a private method called by stack.string
and stack.isEqual.
*/
func (r stack) defaultAssertionHandler(x any) (str string) {

	// str is assigned with
	str = `UNKNOWN`

	// Whatever it is, it had better be one (1) of the following
	// possibilities in the following evaluated order:
	// • An initialized Stack, or type alias of Stack, or ...
	// • An initialized Condition, or type alias of Condition, or ...
	// • Something that has its own "stringer" (String()) method, or ...
	// • A Go primitive
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
			str = ik + ` ` + Xs.String()
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
	pad := padValue(!r.positive(nspad), "")

	builder := newStringBuilder()

	if r.positive(lonce) {
		if oc != list {
			builder.WriteString(ot)
		}
		for _, val := range str {
			builder.WriteString(val)
		}
	} else {
		if oc == list {
			var joinChar string
			if ljc := r.getListDelimiter(); len(ljc) > 0 {
				joinChar = ljc
			} else {
				joinChar = pad
			}
			builder.WriteString(join(str, joinChar))
		} else {
			var tjn string
			var char string
			if len(r.getSymbol()) > 0 {
				if !r.positive(nspad) {
					char = " "
				}
				sympad := padValue(!r.positive(nspad), char)
				j := sympad + ot + sympad
				tjn = join(str, j)
			} else {
				char = " "
				sympad := padValue(true, char)
				j := sympad + ot + sympad
				tjn = join(str, j)
			}
			builder.WriteString(tjn)
		}
	}

	fpad := pad + builder.String() + pad
	result := condenseWHSP(r.paren(fpad))

	return result
}

/*
Traverse will "walk" a structure of stack elements using the path indices
provided. It returns the slice found at the final index, or nil, along with
a success-indicative Boolean value.

The semantics of "traversability" are as follows:

  - Any "nesting" instance must be a [Stack] or [Stack] type alias
  - [Condition] instances must either be the final requested element, OR must contain a [Stack] or [Stack] type alias instance through which the traversal process may continue
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
traverse is a private method called by [Stack.Traverse].
*/
func (r stack) traverse(indices ...int) (slice any, ok, done bool) {
	if r.valid() {
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
Reveal processes the receiver instance and disenvelops needlessly
enveloped [Stack] slices.
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
reveal is a private method called by [Stack.Reveal].
*/
func (r *stack) reveal() (err error) {
	r.lock()
	defer r.unlock()

	// scan each slice (except the config
	// slice) and analyze its structure.
	for i := 0; i < r.len() && err == nil; i++ {
		if sl, _, _ := r.index(i); sl != nil { // cfg offset handled by index method, be honest
			// If the element is a stack, begin descent
			// through recursion.
			if outer, ook := stackTypeAliasConverter(sl); ook && outer.Len() > 0 {
				err = r.revealDescend(outer, i)
			}
		}
	}

	return
}

/*
revealDescend is a private method called by stack.reveal. It engages the second
level of processing of the provided inner [Stack] instance.

Its first order of business is to determine the length of the inner [Stack] instance
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

	// TODO: Still mulling the best way to handle NOTs.
	if inner.stackType() != not {
		switch inner.Len() {
		case 1:
			// descend into inner slice #0
			child, _, _ := inner.index(0)
			if assert, ok := child.(Interface); ok {
				if !assert.IsParen() && !inner.IsParen() {
					err = r.revealSingle(0)
					updated = child
				}
			}
		default:
			// begin new top-level reveal of inner
			// as a whole, scanning all +2 slices
			err = inner.reveal()
			updated = inner
		}
	}

	if err == nil {
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
	}

	return
}

/*
revealSingle is a private method called by stack.revealDescend. It engages the third
level of processing of the receiver's slice instance found at the provided index (idx).
An error is returned at the conclusion of processing.

If a single [Condition] or [Condition] type alias instance (whether pointer-based or not)
is found at the slice index indicated, its Expression value is checked for the presence
of a [Stack] or [Stack] type alias. If one is present, a new top-level stack.reveal recursion
is launched.

If, on the other hand, a [Stack] or [Stack] type alias (again, pointer or not) is found at
the slice index indicates, a new top-level stack.reveal recursion is launched into it
directly.
*/
func (r *stack) revealSingle(idx int) (err error) {

	// get slice @ idx, bail if nil ...
	if slice, _, _ := r.index(idx); slice != nil {
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
	}

	return
}

/*
traverseAssertionHandler handles the type assertion processes during traversal of one (1) or more
nested [Stack]/[Stack] alias instances that may or may not reside in [Condition]/[Condition] alias instances.
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
method will traverse either a [Stack] *OR* [Stack] alias type fashioned by the user that is the [Condition]
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
method will traverse either a [Stack] *OR* [Stack] alias type fashioned by the user.
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
Front returns the slice from the logical "front" of the receiver instance
alongside a Boolean value indicative of success.  The returned slice is
not removed from the receiver instance.

In LIFO mode (the default), this returns the right-most slice. In FIFO mode,
this returns the left-most slice, and is analogous to the concept of "top"
in other queue implementations.
*/
func (r Stack) Front() (slice any, ok bool) {
	if r.IsInit() {

		if r.IsFIFO() {
			for i := 0; i < r.Len(); i++ {
				if slice, ok = r.Index(i); ok {
					break
				}
			}
			return
		}

		for i := r.Len(); i > 0; i-- {
			if slice, ok = r.Index(i - 1); ok {
				break
			}
		}
	}

	return
}

/*
Back returns the slice from the logical "rear" of the receiver instance
alongside a Boolean value indicative of success.  The returned slice is
not removed from the receiver instance.

In LIFO mode (the default), this returns the left-most slice. In FIFO mode,
this returns the right-most slice.
*/
func (r Stack) Back() (slice any, ok bool) {
	if r.IsInit() {
		if !r.IsFIFO() {
			for i := 0; i < r.Len(); i++ {
				if slice, ok = r.Index(i); ok {
					break
				}
			}
			return
		}

		for i := r.Len(); i > 0; i-- {
			if slice, ok = r.Index(i - 1); ok {
				break
			}
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
becomes positive and within the bounds of the [Stack] length, and (perhaps
most importantly) aligned with the relative (intended) slice. To offer an
example, -2 would return the second-to-last slice. When negative index
support is NOT enabled, nil is returned for any index out of bounds along
with a Boolean value of false, although no panic will occur.

Positives: When forward index support is enabled, an index greater than
the length of the [Stack] shall be reduced to the highest valid slice index.
For example, if an index of twenty (20) were used on a [Stack] instance of
a length of ten (10), the index would transform to nine (9). When forward
index support is NOT enabled, nil is returned for any index out of bounds
along with a Boolean value of false, although no panic will occur.

In any scenario, a valid index within the bounds of the [Stack] instance's
length returns the intended slice along with Boolean value of true.
*/
func (r Stack) Index(idx int) (slice any, ok bool) {
	if r.IsInit() {
		slice, _, ok = r.stack.index(idx)
	}
	return
}

/*
index is a private method called by [Stack.Index].
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

  - In the default mode -- LIFO -- this shall be the final slice (index [Stack.Len] - 1", or the "far right" element)
  - In the alternative mode -- FIFO -- this shall be the first slice (index 0, or the "far left" element)

Note that if the receiver is in an invalid state, or has a zero length,
nothing will be removed.
*/
func (r Stack) Pop() (popped any, ok bool) {
	if !r.IsEmpty() {
		if !r.getState(ronly) {
			popped, ok = r.stack.pop()
		}
	}
	return
}

/*
pop is a private method called by [Stack.Pop].
*/
func (r *stack) pop() (slice any, ok bool) {

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
Push appends the provided value(s) to the receiver, and returns the
receiver in fluent form.

Note that if the receiver is in an invalid state, or if maximum capacity
has been set and reached, each of the values intended for append shall
be ignored.
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
push is a private method called by [Stack.Push].
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
	if !r.IsEmpty() {
		if !r.getState(ronly) {
			r.stack.reverse()
		}
	}
	return r
}

/*
reverse is a private niladic and void method called exclusively by [Stack.Reverse].
*/
func (r *stack) reverse() {

	r.lock()
	defer r.unlock()

	for i, j := 1, r.len()-1; i < j; i, j = i+1, j-1 {
		(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
	}
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

If run on a [Stack] or [Stack] type-alias that is currently in possession of one (1) or more nested [Stack]
or [Stack] type-alias instances, Defrag shall hierarchically traverse the structure and process it no
differently than the top-level instance. This applies to such Stack values nested with an instance of
[Condition] or [Condition] type-alias as well.

This is potentially a destructive method and is still very much considered EXPERIMENTAL. While all
tests yield expected results, those who use this method are advised to exercise extreme caution. The
most obvious note of caution pertains to the volatility of index numbers, which shall shift according
to the defragmentation's influence on the instance in question.  By necessity, [Stack.Len] return
values shall also change accordingly.
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
						if sub, ok := stackTypeAliasConverter(cub.Expression()); ok {
							// Condition expression contains a Stack/Stack alias
							sub.Defrag(m)
						}
					}
				}
			}
		}
	}

	return r
}

/*
IsEqual returns a Boolean value indicative of the outcome of a recursive
comparison of all values found within the receiver and input value o.

The parameters of this process are as follows:

  - All Go primitives (numbers, bool and string) are compared as-is
  - If both values are explicit nil, true is returned
  - Invalid instances return false in all cases
  - All pointer instances are dereferenced -- regardless of reference depth (e.g.: **string, *string, etc. become string) -- and then rechecked
  - Stack, Stack-alias, Condition and Condition-alias types utilize their respective `IsEqual` method
  - Structs are compared based on non-private field configuration, field order and underlying values; anonymous (embedded) fields are permitted
  - Slices and Arrays are compared based on matching capacity (if applicable), length, order and content only; the assertion process does not distinguish between the two types
  - Maps are compared based on matching length, keys and values
  - Functions and methods are compared based on their pointer addresses -- or, as a fallback, their respective I/O signatures; this allows distinct closures of like-signatures to qualify
  - Channels, UnsafePointers and Uintptrs are compared as-is
  - Interfaces are de-enveloped and then rechecked

This method is experimental, may change in future releases and can be
particularly costly wherever large or complex instances are concerned.
Please use sparingly.
*/
func (r Stack) IsEqual(o any) error {
	if !r.IsInit() {
		return errorf("Not initialized")
	}

	// handle stack/stack-alias assertion and
	// exit immediately if it fails due to a
	// bad type, or uninitialized input for o.
	if s, ok := stackTypeAliasConverter(o); ok {
		if sc, _ := r.config(); sc.eqf != nil {
			// use the user-authored closure assertion
			// with the original instance
			return sc.eqf(r, o)
		}

		// use default assertion with the converted
		// instance.
		return r.stack.isEqual(s.stack)
	}

	return errorf("Cannot perform equality assertion; bad input")
}

/*
isEqual is a private method called by the package-default IsEqual method.
It calls the top-level valuesEqual function, meant to compare two values,
which in turn calls any number of type-specific equality functions based
on the content encountered.
*/
func (r *stack) isEqual(o *stack) (err error) {
	// Before we bother to run functions,
	// lets see if the two instances are
	// actually the same pointer.
	if r == o {
		return nil
	}

	// compare len/cap of stacks
	if !capLenEqual(r.cap(), o.cap(), r.len(), o.len()) {
		err = errorf("Capacity or length mismatch")
		return
	}

	// Compare the kinds of stacks
	if r.kind() != o.kind() {
		err = errorf("Stack kind mismatch")
		return
	}

	// iterate each slice and compare using
	// the generic valuesEqual function ...
	for i := 0; i < r.ulen() && err == nil; i++ {
		isl, _, _ := r.index(i)
		jsl, _, _ := o.index(i)
		err = valuesEqual(isl, jsl)
	}

	return
}

/*
SetEqualityPolicy sets or unsets the [EqualityPolicy] within the receiver
instance.

When fed an [EqualityPolicy], it shall override the package default mechanism
beginning at the next call of [Stack.IsEqual].

When fed zero (0) [EqualityPolicy] instances, or a value of nil, the previously
specified instance will be removed, at which point the default behavior resumes.
*/
func (r Stack) SetEqualityPolicy(fn ...EqualityPolicy) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			sc, _ := r.config()

			r.lock()
			defer r.unlock()

			if len(fn) == 0 {
				sc.eqf = nil
			} else {
				sc.eqf = fn[0]
			}
		}
	}

	return r
}

/*
SetUnmarshaler sets or unsets the [Unmarshaler] within the receiver
instance.

When fed an [Unmarshaler], it shall override the package default mechanism
beginning at the next call of [Stack.Unmarshal].

When fed zero (0) [Unmarshaler] instances, or a value of nil, the previously
specified instance will be removed, at which point the default behavior
resumes.
*/
func (r Stack) SetUnmarshaler(fn ...Unmarshaler) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			sc, _ := r.config()

			r.lock()
			defer r.unlock()

			if len(fn) == 0 {
				sc.umf = nil
			} else {
				sc.umf = fn[0]
			}
		}
	}

	return r
}

/*
SetMarshaler sets or unsets the [Marshaler] within the receiver instance.

When fed an [Marshaler], it shall override the package default mechanism
beginning at the next call of [Stack.Marshal].

When fed zero (0) [Marshaler] instances, or a value of nil, the previously
specified instance will be removed, at which point the default behavior
resumes.
*/
func (r Stack) SetMarshaler(fn ...Marshaler) Stack {
	if r.IsInit() {
		if !r.getState(ronly) {
			sc, _ := r.config()

			r.lock()
			defer r.unlock()

			if len(fn) == 0 {
				sc.maf = nil
			} else {
				sc.maf = fn[0]
			}
		}
	}

	return r
}

/*
Unmarshal returns slices of any ([]interface{}) based on the raw contents
of the receiver instance. The output produced by this method is identical
to the desired [Stack.Marshal] method input.

This method is intended for generalized use, and may be overridden using
the [Stack.SetUnmarshaler] method.
*/
func (r Stack) Unmarshal() (slice []any, err error) {
	if r.IsInit() {
		if sc, _ := r.config(); sc.umf != nil {
			// use the user-authored closure unmarshaler
			slice, err = sc.umf()
		} else {
			// use default unmarshaler
			slice, err = r.stack.unmarshalDefault()
		}
	}

	return
}

/*
unmarshalDefault is a private method called by Stack.Unmarshal.
*/
func (r stack) unmarshalDefault() (slices []any, err error) {
	slices = append(slices, r.kind())
	for i := 0; i < r.ulen() && err == nil; i++ {
		slice, _, _ := r.index(i) // auto-skip config
		var subSlices []any
		if sub, ok := stackTypeAliasConverter(slice); ok {
			// Instance is Stack/Stack alias;
			// use native unmarshalDefault.
			if subSlices, err = sub.unmarshalDefault(); err == nil {
				slices = append(slices, subSlices)
			}
		} else if cub, ok := conditionTypeAliasConverter(slice); ok {
			// Instance is Condition/Condition alias;
			// use the Condition.Unmarshal method.
			if subSlices, err = cub.Unmarshal(); err == nil {
				slices = append(slices, subSlices)
			}
		} else {
			// Anything and everything -- even nil -- that is not a
			// Stack/Stack alias or Condition/Condition alias, will
			// be taken as-is.
			slices = append(slices, slice)
		}
	}

	return
}

/*
Marshal returns an error following an attempt to read the variadic 'in'
value(s) into the receiver instance. The appropriate input for this method
is the output produced by [Stack.Unmarshal].

This method is intended for generalized use, and may be overridden using
the [Stack.SetMarshaler] method.
*/
func (r *Stack) Marshal(in ...any) (err error) {
	if len(in) == 0 {
		err = errorf("Empty marshaler input")
	} else {
		var xs Stack
		var xc Condition

		if !r.IsInit() {
			// use default marshaler
			if xs, xc, err = marshalDefault(in); xs.IsInit() {
				r.stack = xs.stack
			} else if xc.IsInit() {
				err = errorf("Cannot Unmarshal Condition only; must envelope in Stack")
			}
		} else if sc, _ := r.config(); sc.maf != nil {
			// use the user-authored closure marshaler
			err = sc.maf(in...)
		} else {
			// use default marshaler
			if xs, xc, err = marshalDefault(in); xs.IsInit() {
				r.Push(xs)
			} else if xc.IsInit() {
				r.Push(xc)
			}
		}
	}

	return
}

func stackByWord(label string) Stack {
	switch uc(label) {
	case `LIST`:
		return List()
	case `AND`:
		return And()
	case `NOT`:
		return Not()
	case `OR`:
		return Or()
	}

	return Basic()
}

func extractConditionValues(in []any) (c Condition, ok bool) {
	if len(in) != 4 {
		return
	}
	var word string
	var op Operator

	if W, ok := in[1].(string); ok {
		word = W
	}
	if O, ok := in[2].(Operator); ok {
		op = O
	}
	if E, ok := in[3].([]any); ok {
		var xm Stack
		var xn Condition
		xm, xn, _ = marshalDefault(E)
		if xm.IsInit() {
			c = Cond(word, op, xm)
		} else if xn.IsInit() {
			c = Cond(word, op, xn)
		}
	} else {
		c = Cond(word, op, in[3])
	}

	return
}

func marshalDefault(in []any) (x Stack, c Condition, err error) {
	if len(in) == 0 {
		err = errorf("Empty input")
		return
	}

	// De-envelope needlessly enveloped value
	in = deenvelopeSingleStack(in)

	// The first string value in a stack indicates the
	// appropriate type of stack or condition
	lab, ok := in[0].(string)
	if !ok {
		err = errorf("Cannot unmarshal without stack label")
		return
	}

	switch uc(lab) {
	case `CONDITION`:
		// A condition is an instance of []any
		// with four values. The 1st is the
		// 'CONDITION' label. The 2nd is the
		// Keyword of a condition. The 3rd is
		// the Operator and the last is the
		// expression (value).  Convert this
		// to a proper instance of Condition.
		c, _ = extractConditionValues(in)
		return
	case `LIST`, `AND`, `OR`, `NOT`, `BASIC`:
		x = stackByWord(lab).Push(in[1:]...)
	default:
		// No idea what the value is, just
		// use a Basic
		x = Basic().Push(in...)
	}

	// Iterate through the new stack (x) and
	// Make sure there are no nested []any
	// instance. If found, try to convert
	// them to Stacks or Conditions by self
	// executing this same function and if
	// converted, replace the index with the
	// new value.
	for i := 0; i < x.Len(); i++ {
		slice, _ := x.Index(i)
		if tv, aok := slice.([]any); aok {
			var xz Stack
			var xc Condition
			if xz, xc, err = marshalDefault(tv); xz.IsInit() {
				// Was a stack; replace old slice
				x.Replace(xz, i)
			} else if xc.IsInit() {
				// Was a condition; replace old slice
				x.Replace(xc, i)
			}
		}
	}

	return
}

func deenvelopeSingleStack(in []any) []any {
	if len(in) == 1 {
		for {
			if inner, ok := in[0].([]any); ok {
				in = inner
			} else {
				break
			}
		}
	}

	return in
}

/*
calculateDefragMax is a private function executed exclusively by [Stack.Defrag], and
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

	return
}

func (r *stack) defrag(max int) {

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

	if !(start == -1 || max <= start) {
		tpat := r.implode(start, max, spat)
		last, err := r.verifyImplode(spat, tpat)
		r.setErr(err)
		if err == nil && last >= 0 {
			// chop off the remaining consecutive nil slices
			(*r) = (*r)[:last+1]
		}
	}

	return
}

func (r stack) verifyImplode(spat, tpat []int) (last int, err error) {
	last = -1

	err = errorf("defragmentation failed; inconsistent slice results")
	data := make(map[string]string, len(tpat))
	var fail bool
	for i := 1; i < len(spat); i++ {
		key := `S[` + itoa(i-1) + `]`
		result := spat[i] == tpat[i]
		fail = true
		if result {
			fail = false
		}
		if tpat[i] != 0 {
			last = (len(data) + i) - len(tpat)
		}

		data[key] = `match:` + bool2str(result)
	}

	if !fail {
		err = nil
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

	for {
		if ct >= max || start+ct >= r.ulen() {
			break
		}

		if (*r)[start+ct+1] == nil {
			ct++
			continue
		}

		(*r)[start+1] = (*r)[start+ct+1]

		tpat[start+ct] = 1

		(*r)[start+ct+1] = nil
		start = start + 1
		ct = 0
	}

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

	// use the user-provided function to scan
	// each pushed item for verification.
	var pct int
	for i := 0; i < len(x); i++ {
		var err error
		if !r.isFull() {
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
IsFull returns a Boolean value indicative of whether the receiver has reached
the maximum configured capacity. This method wraps [Stack.Len] == [Stack.Cap].
When no maximum capacity is specified, this method returns false.
*/
func (r Stack) IsFull() (full bool) {
	if r.IsInit() {
		full = r.isFull()
	}

	return
}

/*
Deprecated: Use [Stack.IsFull] instead.
*/
func (r Stack) CapReached() bool {
	return r.IsFull()
}

/*
isFull is a private method called by [Stack.IsFull].
*/
func (r stack) isFull() (rc bool) {
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
			if !r.isFull() {
				*r = append(*r, x[i])
				pct++
			}
		}
	}
}

/*
SetPushPolicy assigns the provided [PushPolicy] closure function
to the receiver, thereby enabling protection against undesired
appends to the [Stack]. The provided function shall be executed
by the [Stack.Push] method for each individual item being added.
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
setPushPolicy is a private method called by [Stack.SetPushPolicy].
*/
func (r *stack) setPushPolicy(ppol PushPolicy) *stack {
	sc, _ := r.config()
	sc.ppf = ppol
	return r
}

/*
getPushPolicy is a private method called by [Stack.Push].
*/
func (r *stack) getPushPolicy() PushPolicy {
	sc, _ := r.config()
	return sc.ppf
}

/*
SetPresentationPolicy assigns the provided [PresentationPolicy]
closure function to the receiver, thereby enabling full control
over the stringification of the receiver.  Execution of the
[Stack.String] method will execute the provided policy instead
of the package-provided routine.
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
setPresentationPolicy is a private method called by [Stack.SetPresentationPolicy].
*/
func (r *stack) setPresentationPolicy(ppol PresentationPolicy) *stack {

	if r.stackType() == basic {
		err := errorf("ppolicy incompatible with basic stack type")
		r.setErr(err)
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
SetValidityPolicy assigns the provided [ValidityPolicy] closure function
instance to the receiver, thereby allowing users to introduce inline
verification checks of a [Stack] to better gauge its validity. The provided
function shall be executed by the [Stack.Valid] method.
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
setValidityPolicy is a private method called by [Stack.SetValidityPolicy].
*/
func (r *stack) setValidityPolicy(vpol ValidityPolicy) *stack {
	sc, _ := r.config()
	sc.vpf = vpol
	return r
}

/*
getValidityPolicy is a private method called by [Stack.Valid].
*/
func (r *stack) getValidityPolicy() ValidityPolicy {
	sc, _ := r.config()
	return sc.vpf
}

/*
ConvertStack returns an instance of [Stack] alongside a Boolean value.

If the input value is a native [Stack], it is returned as-is alongside a
Boolean value of true.

If the input value is a [Stack]-alias, it is converted to a native [Stack]
instance and returned alongside a Boolean value of true.

Any other scenario returns a zero [Stack] alongside a Boolean value of
false.
*/
func ConvertStack(in any) (Stack, bool) {
	return stackTypeAliasConverter(in)
}

/*
stackTypeAliasConverter attempts to convert any (u) back to a bonafide
instance of Stack. This will only work if input value u is a type alias
of Stack.  An instance of Stack is returned along with a Boolean value
of true.
*/
func stackTypeAliasConverter(u any) (S Stack, converted bool) {
	if u != nil {
		// If it isn't a Stack alias, but is a
		// genuine Stack, just pass it back
		// with a thumbs-up ...
		if st, isStack := u.(Stack); isStack {
			S = st
			converted = isStack
			return
		}

		a, v, _ := derefPtr(typOf(u), valOf(u))
		b := typOf(Stack{}) // target (dest) type
		if a.ConvertibleTo(b) {
			X := v.Convert(b).Interface()
			if assert, ok := X.(Stack); ok {
				if !assert.IsZero() {
					S = assert
					converted = true
				}
			}
		}
	}

	return
}

const badStack = `<invalid_stack>`
