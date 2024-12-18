package stackage

/*
cond.go contains Condition-related methods and functions.
*/

/*
Condition describes a single evaluative statement, i.e.:

	       op
	       |
	       v
	person = "Jesse"
	   ^         ^
	   |         |
	   kw        ex

The keyword (kw) shall always represent an abstract user-defined
string construct against which the expression value (ex) is to be
evaluated in some manner.

The disposition of the evaluation is expressed through one (1) of
several [ComparisonOperator] (op) instances made available through
this package:

  - Eq, or "equal to" (=)
  - Ne, or "not equal to" (!=)
  - Lt, or "less than" (<)
  - Le, or "less than or equal" (<=)
  - Gt, or "greater than" (>)
  - Ge, or "greater than or equal" (>=)

... OR through a user-defined operator that conforms to the package
defined Operator interface.

By default, permitted expression (ex) values must honor these guidelines:

  - Must be a non-zero string, OR ...
  - Must be a valid instance of [Stack] (or an *alias* of [Stack] that is convertible back to the [Stack] type), OR ...
  - Must be a valid instance of any type that exports a stringer method (String()) intended to produce the [Condition] instances final expression string representation

However when a [PushPolicy] function or method is added to an instance of this type,
greater control is afforded to the user in terms of what values will be accepted,
as well as the quality or state of such values.

Instances of this type -- similar to [Stack] instances -- MUST be initialized before use.
Initialization can occur as a result of executing the [Cond] package-level function, or
using the [Condition.Init] method extended through instances of this type. Initialization
state may be checked using the [Condition.IsInit] method.
*/
type Condition struct {
	*condition
}

/*
condition is the private embedded type to be circumscribed within instances
of Condition.
*/
type condition struct {
	cfg *nodeConfig
	kw  string
	op  Operator
	ex  any // expression value
}

/*
newCondition obtains, (optionally sets) and returns a new instance of
*condition in one shot.
*/
func newCondition(kw any, op Operator, ex any) (r *condition) {
	r = initCondition()

	r.setKeyword(kw)    // keyword
	r.setOperator(op)   // operator
	r.setExpression(ex) // expr. value(s)

	return
}

/*
initCondition is the central initializer function for an instance
of *condition, which is the embedded type instance found within
(properly initialized) Condition instances.
*/
func initCondition() (r *condition) {

	r = new(condition)
	r.cfg = new(nodeConfig)
	r.cfg.log = newLogSystem(cLogDefault)
	r.cfg.log.lvl = logLevels(NoLogLevels)

	r.cfg.typ = cond

	return
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
func (r Condition) SetAuxiliary(aux ...Auxiliary) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.setAuxiliary(aux...)
		}
	}
	return r
}

/*
setAuxiliary is a private method called by Condition.SetAuxiliary.
*/
func (r *condition) setAuxiliary(aux ...Auxiliary) {
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

	r.cfg.aux = _aux
}

/*
Auxiliary returns the instance of [Auxiliary] from within the receiver.
*/
func (r Condition) Auxiliary() (aux Auxiliary) {
	if r.IsInit() {
		aux = r.condition.auxiliary()
	}
	return
}

/*
auxiliary is a private method called by [Condition.Auxiliary].
*/
func (r condition) auxiliary() (aux Auxiliary) {
	aux = r.cfg.aux
	return
}

/*
SetKeyword sets the receiver's keyword using the specified kw
input argument.
*/
func (r Condition) SetKeyword(kw any) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.setKeyword(kw)
		}
	}

	return r
}

func (r *condition) setKeyword(kw any) {
	switch tv := kw.(type) {
	case string:
		r.kw = tv
	default:
		if meth := getStringer(tv); meth != nil {
			r.kw = meth()
		}
	}
}

/*
SetOperator sets the receiver's [ComparisonOperator] using the
specified [Operator]-qualifying input argument (op).
*/
func (r Condition) SetOperator(op Operator) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.setOperator(op)
		}
	}
	return r
}

func (r *condition) setOperator(op Operator) {
	if len(op.Context()) > 0 && len(op.String()) > 0 {
		r.op = op
	}
}

/*
SetExpression sets the receiver's expression value(s) using the
specified ex input argument.  See also the [Condition.Expression]
method.
*/
func (r Condition) SetExpression(ex any) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.setExpression(ex)
		}
	}
	return r
}

func (r *condition) setExpression(ex any) {
	if v, ok := r.assertConditionExpressionValue(ex); ok {
		r.ex = v
	}
}

/*
SetEqualityPolicy sets or unsets the [EqualityPolicy] within the receiver
instance.

When fed an [EqualityPolicy], it shall override the package default mechanism
beginning at the next call of [Condition.IsEqual].

When fed zero (0) [EqualityPolicy] instances, or a value of nil, the previously
specified instance will be removed, at which point the default behavior resumes.
*/
func (r Condition) SetEqualityPolicy(fn ...EqualityPolicy) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			if len(fn) == 0 {
				r.condition.cfg.eqf = nil
			} else {
				r.condition.cfg.eqf = fn[0]
			}
		}
	}

	return r
}

/*
SetUnmarshaler sets or unsets the [Unmarshaler] within the receiver
instance.

When fed an [Unmarshaler], it shall override the package default mechanism
beginning at the next call of [Condition.Unmarshal].

When fed zero (0) [Unmarshaler] instances, or a value of nil, the previously
specified instance will be removed, at which point the default behavior resumes.
*/
func (r Condition) SetUnmarshaler(fn ...Unmarshaler) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			if len(fn) == 0 {
				r.condition.cfg.umf = nil
			} else {
				r.condition.cfg.umf = fn[0]
			}
		}
	}

	return r
}

/*
Unmarshal returns an instance of []any containing the unmarshaled instance
of the receiver. This can be used for use in deep inspections of [Condition]
or [Condition]-alias instances.

Those instances that contain complex nested [Stack] or [Stack]-alias instances
specified as the [Condition.Expression] will have those instances unmarshaled
by way of the [Stack.Unmarshal] method.

All other type instances used for the [Condition.Expression] are added to the
return value as-is in all cases, even if nil.

Note that the underlying configuration within any [Condition] or [Condition]-alias
instance is lost during the transfer, thus the return value cannot easily
be used in the reverse context, meaning it cannot aid in marshaling a new
[Condition] or [Condition]-alias instance under ordinary circumstances.
*/
func (r Condition) Unmarshal() (slice []any, err error) {
	if r.IsInit() {
		if fn := r.condition.cfg.umf; fn != nil {
			// use the user-authored closure unmarshaler
			slice, err = fn()
		} else {
			// use default unmarshaler
			slice, err = r.condition.unmarshalDefault()
		}
	}

	return
}

/*
unmarshalDefault is a private method called by Condition.Unmarshal.
*/
func (r condition) unmarshalDefault() (slice []any, err error) {
	var nexpr any
	if s, ok := stackTypeAliasConverter(r.ex); ok {
		nexpr, err = s.Unmarshal() // unmarshaled stack/stack-alias
	} else {
		nexpr = r.ex // orig
	}

	slice = []any{
		`CONDITION`,
		r.kw,
		r.op,
		nexpr,
	}

	return
}

/*
IsEqual returns a Boolean value indicative of the outcome of a recursive
comparison of all values found within the receiver and input value o.

The parameters of this process are described within the [Stack.IsEqual]
notes.
*/
func (r Condition) IsEqual(o any) (err error) {
	if r.IsInit() {
		// handle condition/condition-alias assertion
		// and exit immediately if it fails due to a
		// bad type, or uninitialized input for o.
		if s, ok := conditionTypeAliasConverter(o); ok {
			if fn := r.condition.cfg.eqf; fn != nil {
				// use the user-authored closure assertion
				err = fn(r, o)
			} else {
				// use default assertion
				err = r.condition.isEqual(s.condition)
			}
		}
	}

	return
}

func (r *condition) isEqual(o *condition) error {
	if r.kw != o.kw {
		return errorf("Condition keyword mismatch")
	}

	if r.op.String() != o.op.String() {
		return errorf("Condition operator mismatch")
	}

	if r.op.Context() != o.op.Context() {
		return errorf("Condition operator (context) mismatch")
	}

	iexpr := r.ex
	jexpr := o.ex

	return valuesEqual(iexpr, jexpr)
}

/*
SetLogLevel enables the specified [LogLevel] instance(s), thereby
instructing the logging subsystem to accept events for submission
and transcription to the underlying logger.

Users may also sum the desired bit values manually, and cast the
product as a [LogLevel]. For example, if STATE (4), DEBUG (8) and
TRACE (32) logging were desired, entering LogLevel(44) would be
the same as specifying LogLevel3, LogLevel4 and LogLevel6 in
variadic fashion.
*/
func (r Condition) SetLogLevel(l ...any) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.setLogLevel(l...)
		}
	}
	return r
}

func (r *condition) setLogLevel(l ...any) {
	r.cfg.log.shift(l...)
}

/*
LogLevels returns the string representation of a comma-delimited list
of all active [LogLevel] values within the receiver.
*/
func (r Condition) LogLevels() (l string) {
	if r.IsInit() {
		l = r.condition.logLevels()
	}
	return
}

func (r condition) logLevels() (l string) {
	return r.cfg.log.lvl.String()
}

/*
UnsetLogLevel disables the specified [LogLevel] instance(s), thereby
instructing the logging subsystem to discard events submitted for
transcription to the underlying logger.
*/
func (r Condition) UnsetLogLevel(l ...any) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.unsetLogLevel(l...)
		}
	}
	return r
}

func (r *condition) unsetLogLevel(l ...any) {
	r.cfg.log.unshift(l...)
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
func (r Condition) SetLogger(logger any) Condition {
	if r.IsInit() {
		r.condition.setLogger(logger)
	}

	return r
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
func (r Condition) Err() (err error) {
	if r.IsInit() {
		err = r.condition.getErr()
	}
	return
}

/*
SetErr sets the underlying error value within the receiver
to the assigned input value err, whether nil or not.

This method may be most valuable to users who have chosen
to extend this type by aliasing, and wish to control the
handling of error conditions in another manner.

This may be used regardless of [Condition.IsReadOnly] status.
*/
func (r Condition) SetErr(err error) Condition {
	if r.IsInit() {
		r.condition.setErr(err)
	}
	return r
}

/*
setErr assigns an error instance, whether nil or not, to
the underlying receiver configuration.
*/
func (r *condition) setErr(err error) {
	r.cfg.setErr(err)
}

/*
error returns the underlying error instance, whether nil or not.
*/
func (r condition) getErr() error {
	return r.cfg.getErr()
}

/*
assertConditionExpressionValue returns a string value alongside a
success-indicative Boolean value. This method is used during the
set-execution of the receiver, and is designed to prevent unwanted
types from being assigned as the expression value (ex).
*/
func (r *condition) assertConditionExpressionValue(x any) (X any, ok bool) {

	switch tv := x.(type) {
	case string:
		if len(tv) > 0 {
			X = tv
		}
	default:
		X = r.defaultAssertionExpressionHandler(x)
	}

	ok = X != nil && !r.cfg.isError()
	return
}

/*
defaultAssertionExpressionHandler is the catch-all private method called
by condition.assertConditionExpressionValue.
*/
func (r condition) defaultAssertionExpressionHandler(x any) (X any) {
	// no push policy, so we'll see if the basic
	// guidelines were satisfied, at least ...
	if _, ok := stackTypeAliasConverter(x); ok {
		if r.positive(nnest) {
			return
		}

		// a user-created type alias of Stack
		// was converted back to Stack without
		// any issues ...
		X = x
	} else {
		X = x
	}

	return
}

/*
setCategory assigns the provided categorical label string value
to the underlying configuration.
*/
func (r *condition) setCategory(cat string) {
	r.cfg.cat = cat
}

/*
getCat returns the categorical label string value assigned to
the underlying configuration, or a zero string if unset.
*/
func (r condition) getCat() string {
	return r.cfg.cat
}

/*
SetCategory assigns the provided string to the receiver internal category
value. This allows for a means of identifying a particular kind of [Condition]
in the midst of many.
*/
func (r Condition) SetCategory(cat string) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.setCategory(cat)
		}
	}
	return r
}

/*
Category returns the categorical label string value assigned to the receiver, if
set, else a zero string.
*/
func (r Condition) Category() (cat string) {
	if r.IsInit() {
		cat = r.condition.getCat()
	}
	return
}

/*
Cond returns an instance of [Condition] bearing the provided component values.
This is intended to be used in situations where a [Condition] instance can be
created in one shot.
*/
func Cond(kw any, op Operator, ex any) (c Condition) {
	c = Condition{newCondition(kw, op, ex)}
	if err := c.Valid(); err != nil {
		c.SetErr(err)
	}
	return
}

/*
Init will [re-]initialize the receiver's contents and return them to an unset,
but assignable, state. This is a destructive method: the embedded pointer within
the receiver instance shall be totally annihilated.

This method is niladic and fluent in nature. No input is required, and the only
element returned is the receiver itself.

This method may be useful in situations where a [Condition] will be assembled in a
"piecemeal" fashion (i.e.: incrementally), or if a [Condition] instance is slated to
be repurposed for use elsewhere (possibly in a repetative manner).
*/
func (r *Condition) Init() Condition {
	*r = Condition{condition: initCondition()}
	return *r
}

/*
Free frees the receiver instance entirely, including the underlying
configuration. An error is returned if the instance is read-only and
cannot be freed.
*/
func (r *Condition) Free() (err error) {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition = nil
			return
		}
		err = errorf("%T is read-only; cannot free", r)
	}

	return
}

/*
SetID assigns the provided string value (or lack thereof) to the receiver.
This is optional, and is usually only needed in complex [Condition] structures
in which "labeling" certain components may be advantageous. It has no effect
on an evaluation, nor should a name ever cause a validity check to fail.

If the string `_random` is provided, a 24-character alphanumeric string is
randomly generated using math/rand and assigned as the ID.
*/
func (r Condition) SetID(id string) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			if lc(id) == `_random` {
				id = randomID(randIDSize)
			} else if lc(id) == `_addr` {
				id = r.Addr()
			}

			r.condition.cfg.setID(id)
		}
	}
	return r
}

/*
Len returns a "perceived" abstract length relating to the content (or lack
thereof) assigned to the receiver instance:

  - An uninitialized or zero instance returns zero (0)
  - An initialized instance with no [Condition.Expression] assigned (nil) returns zero (0)
  - A [Stack] or [Stack] type alias assigned as the [Condition.Expression] shall impose its own length as the return value (even if zero (0))

All other type instances assigned as an [Condition.Expression] shall result in a
return of one (1); this includes slice types, maps, arrays and any other
type that supports multiple values.

This capability was added to this type to mirror that of the [Stack] type in
order to allow additional functionality to be added to the [Interface] interface.
*/
func (r Condition) Len() int {
	if !r.IsInit() {
		return 0
	}

	if r.Expression() == nil {
		return 0
	}

	if stk, ok := stackTypeAliasConverter(r.Expression()); ok {
		return stk.Len()
	}

	return 1
}

/*
Addr returns the string representation of the pointer
address for the receiver. This may be useful for logging
or debugging operations.

Note: this method calls [fmt.Sprintf].
*/
func (r Condition) Addr() (addr string) {
	if r.IsInit() {
		p := r.condition
		addr = sprintf("%p", p)
	}
	return
}

/*
Name returns the name of the receiver instance, if set, else a zero string
will be returned. The presence or lack of a name has no effect on any of
the receiver's mechanics, and is strictly for convenience.
*/
func (r Condition) ID() (id string) {
	if r.IsInit() {
		id = r.condition.cfg.id
	}
	return
}

/*
IsInit will verify that the internal pointer instance of the receiver has
been properly initialized. This method executes a preemptive execution of
the [Condition.IsZero] method.
*/
func (r Condition) IsInit() (is bool) {
	if !r.IsZero() {
		is = r.condition.isInit()
	}
	return
}

/*
isInit is a private method called by Condition.IsInit.
*/
func (r *condition) isInit() bool {
	return r.cfg.typ == cond
}

/*
IsZero returns a Boolean value indicative of whether the receiver is nil,
or unset.
*/
func (r Condition) IsZero() bool {
	return r.condition == nil
}

/*
Valid returns an instance of error, identifying any serious issues perceived
with the state of the receiver.

Non-serious (interim) errors such as denied pushes, capacity violations, etc.,
are not shown by this method.

If a [ValidityPolicy] was set within the receiver, it shall be executed here.
If no [ValidityPolicy] was specified, only elements pertaining to basic viability
are checked.
*/
func (r Condition) Valid() (err error) {
	if !r.IsInit() {
		err = errorf("condition instance is nil")
		return
	}

	// if a validitypolicy was provided, use it
	if r.condition.cfg.vpf != nil {
		err = r.condition.cfg.vpf(r)
		return
	}

	// no validity policy was provided, just check
	// what we can.

	// verify keyword
	if kw := r.Keyword(); len(kw) == 0 {
		err = errorf("keyword value is zero")
		return
	}

	// verify comparison operator
	if cop := r.Operator(); cop != nil {
		if assert, ok := cop.(ComparisonOperator); ok {
			if !(1 <= int(assert) && int(assert) <= 6) {
				err = errorf("operator value is bogus")
				return
			}
		}
	}

	// verify expression value
	if r.Expression() == nil {
		err = errorf("expression value is nil")
	}

	return
}

/*
Evaluate uses the [Evaluator] closure function to apply the value (x)
to the receiver in order to conduct a matching/assertion test or
analysis for some reason. This is entirely up to the user.

An expression value is returned alongside an error. Note that if an
instance of [Evaluator] was not assigned to the [Condition] prior to
execution of this method, the return value shall always be false.
*/
func (r Condition) Evaluate(x ...any) (ev any, err error) {
	if r.IsInit() {
		if err = errorf("No func/meth found"); r.cfg.evl != nil {
			ev, err = r.cfg.evl(x...)
		}
	}

	return
}

/*
SetEvaluator assigns the instance of [Evaluator] to the receiver. This
will allow the [Condition.Evaluate] method to return a more meaningful
result.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetEvaluator(x Evaluator) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.cfg.evl = x
		}
	}
	return r
}

/*
SetValidityPolicy assigns the instance of [ValidityPolicy] to the receiver.
This will allow the [Condition.Valid] method to return a more meaningful
result.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetValidityPolicy(x ValidityPolicy) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.cfg.vpf = x
		}
	}

	return r
}

/*
SetPresentationPolicy assigns the instance of [PresentationPolicy] to the receiver.
This will allow the user to leverage their own "stringer" method for automatic
use when the [Condition.String] method is called.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetPresentationPolicy(x PresentationPolicy) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.cfg.rpf = x
		}
	}

	return r
}

/*
SetEncap accepts input characters for use in controlled condition value
encapsulation. Acceptable input types are:

A single string value will be used for both L and R encapsulation.

An instance of []string with two (2) values will be used for L and R
encapsulation using the first and second slice values respectively.

An instance of []string with only one (1) value is identical to the act of
providing a single string value, in that both L and R will use the one value.
*/
func (r Condition) SetEncap(x ...any) Condition {
	if r.IsInit() {
		if !r.getState(ronly) {
			r.condition.cfg.setEncap(x...)
		}
	}
	return r
}

/*
Deprecated: Use [Condition.SetEncap].
*/
func (r Condition) Encap(x ...any) Condition {
	return r.SetEncap(x...)
}

/*
IsEncap returns a Boolean value indicative of whether value encapsulation
characters have been set within the receiver.
*/
func (r Condition) IsEncap() (is bool) {
	if r.IsInit() {
		is = len(r.condition.getEncap()) > 0
	}
	return
}

func (r *condition) getEncap() [][]string {
	return r.cfg.enc
}

/*
SetNoNesting sets the no-nesting bit within the receiver. If
set to true, the receiver shall ignore any [Stack] or [Stack]
type alias instance when assigned using the [Condition.SetExpression]
method. In such a case, only primitives, etc., shall be honored during
the [Condition.SetExpression] operation.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the nesting bit (i.e.: true->false and
false->true)
*/
func (r Condition) SetNoNesting(state ...bool) Condition {
	r.setState(nnest, state...)
	return r
}

/*
Deprecated: Use [Condition.SetNoNesting].
*/
func (r Condition) NoNesting(state ...bool) Condition {
	return r.SetNoNesting(state...)
}

/*
CanNest returns a Boolean value indicative of whether
the no-nesting bit is unset, thereby allowing a [Stack]
or [Stack] type alias instance to be set as the value.

See also the [Condition.IsNesting] method.
*/
func (r Condition) CanNest() (can bool) {
	if r.IsInit() {
		can = !r.getState(nnest)
	}
	return
}

/*
IsNesting returns a Boolean value indicative of whether the
underlying expression value is either a [Stack] or [Stack] type
alias. If true, this indicates the expression value descends
into another hierarchical (nested) context.
*/
func (r Condition) IsNesting() (is bool) {
	if r.IsInit() {
		is = r.condition.isNesting()
	}
	return
}

/*
isNesting is a private method called by Condition.IsNesting.
*/
func (r condition) isNesting() bool {
	// If convertible is true, we know the
	// instance (tv) is a stack alias.
	_, convertible := stackTypeAliasConverter(r.ex)
	return convertible
}

/*
IsFIFO returns a Boolean value indicative of whether the underlying
receiver instance's [Condition.Expression] value represents a [Stack]
(or [Stack] type alias) instance which exhibits First-In-First-Out
behavior as it pertains to the act of appending and truncating the
receiver's slices.

A value of false implies that no such [Stack] instance is set as the
expression, OR that the [Stack] exhibits Last-In-Last-Out behavior,
which is the default ingress/egress scheme imposed upon instances
of this type.
*/
func (r Condition) IsFIFO() (is bool) {
	if r.IsInit() {
		is = r.condition.isFIFO()
	}
	return
}

/*
isFIFO is a private method called by the Condition.IsFIFO method.
*/
func (r condition) isFIFO() (result bool) {

	if stk, ok := stackTypeAliasConverter(r.ex); ok {
		result = stk.IsFIFO()
	}

	return
}

/*
SetParen sets the string-encapsulation bit for parenthetical
expression within the receiver. The receiver shall undergo
parenthetical encapsulation ( (...) ) during the string
representation process. Individual string values shall not
be encapsulated in parenthesis, only the whole (current)
[Stack] instance.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the encapsulation bit (i.e.: true->false
and false->true)
*/
func (r Condition) SetParen(state ...bool) Condition {
	r.setState(parens, state...)
	return r
}

/*
Deprecated: Use [Condition.SetParen].
*/
func (r Condition) Paren(state ...bool) Condition {
	return r.SetParen(state...)
}

/*
IsParen returns a Boolean value indicative of whether the
receiver is parenthetical.
*/
func (r Condition) IsParen() bool {
	return r.getState(parens)
}

/*
SetReadOnly sets the receiver bit 'ronly' to a positive state.
This will prevent any writes to the receiver or its underlying
configuration.
*/
func (r Condition) SetReadOnly(state ...bool) Condition {
	r.setState(ronly, state...)
	return r
}

/*
IsReadOnly returns a Boolean value indicative of whether the
receiver is set as read-only.
*/
func (r Condition) IsReadOnly() bool {
	return r.getState(ronly)
}

func (r Condition) getState(cf cfgFlag) (state bool) {
	if r.IsInit() {
		state = r.condition.positive(cf)
	}
	return
}

func (r Condition) setState(cf cfgFlag, state ...bool) {
	if r.IsInit() {
		if !r.getState(ronly) || cf == ronly {
			if len(state) > 0 {
				if state[0] {
					r.condition.setOpt(cf)
				} else {
					r.condition.unsetOpt(cf)
				}
			} else {
				r.condition.toggleOpt(cf)
			}
		}
	}
}

func (r *condition) setLogger(logger any) {
	r.cfg.log.setLogger(logger)
}

func (r *condition) toggleOpt(cf cfgFlag) {
	r.cfg.toggleOpt(cf)
}

func (r *condition) setOpt(cf cfgFlag) {
	r.cfg.setOpt(cf)
}

func (r *condition) unsetOpt(cf cfgFlag) {
	r.cfg.unsetOpt(cf)
}

func (r *condition) positive(cf cfgFlag) bool {
	return r.cfg.positive(cf)
}

/*
SetNoPadding sets the no-space-padding bit within the receiver.
String values within the receiver shall not be padded using
a single space character (ASCII #32).

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the quotation bit (i.e.: true->false and
false->true)
*/
func (r Condition) SetNoPadding(state ...bool) Condition {
	r.setState(nspad, state...)
	return r
}

/*
Deprecated: Use [Condition.SetNoPadding].
*/
func (r Condition) NoPadding(state ...bool) Condition {
	return r.SetNoPadding(state...)
}

/*
IsPadded returns a Boolean value indicative of whether the
receiver pads its contents with a SPACE char (ASCII #32).
*/
func (r Condition) IsPadded() (is bool) {
	return !r.getState(nspad)
}

/*
Expression returns the expression value(s) stored within the receiver, or
nil if unset. A valid receiver instance MUST always possess a non-nil
expression value.
*/
func (r Condition) Expression() (ex any) {
	if r.IsInit() {
		ex = r.condition.ex
	}
	return
}

/*
Operator returns the [Operator] interface type instance found within the
receiver.
*/
func (r Condition) Operator() (op Operator) {
	if r.IsInit() {
		op = r.condition.op
	}
	return
}

/*
Keyword returns the Keyword interface type instance found within the
receiver.
*/
func (r Condition) Keyword() (kw string) {
	if r.IsInit() {
		kw = r.condition.kw
	}
	return
}

/*
String is a stringer method that returns the string representation
of the receiver instance. It will only function if the receiver is
in good standing, and passes validity checks.

Note that if the underlying expression value is not a known type,
such as a [Stack] or a Go primitive, this method may be uncertain
as to what it should do. A bogus string may be returned.

In such a case, it may be necessary to subvert the default string
representation behavior demonstrated by instances of this type in
favor of a custom instance of the [PresentationPolicy] closure type
for maximum control.
*/
func (r Condition) String() (s string) {
	if err := r.Valid(); err == nil {
		s = r.condition.string()
	}
	return
}

/*
string is a stringer method that returns the string representation
of the receiver instance.
*/
func (r condition) string() string {
	if r.cfg.rpf != nil {
		return r.cfg.rpf(r)
	}

	// begin default presentation
	// handler ...
	var raw string
	if meth := getStringer(r.ex); meth != nil {
		raw = meth()
	} else {
		raw = primitiveStringer(r.ex)
	}

	val := encapValue(r.cfg.enc, raw)
	var pad string = string(rune(32))
	if r.cfg.positive(nspad) {
		pad = ``
	}

	s := r.kw + pad + r.op.String() + pad + val
	if r.cfg.positive(parens) {
		s = `(` + pad + s + pad + `)`
	}

	return s
}

/*
ConvertCondition returns an instance of [Condition] alongside a Boolean
value.

If the input value is a native [Condition], it is returned as-is alongside
a Boolean value of true.

If the input value is a [Condition]-alias, it is converted to a native
[Condition] instance and returned alongside a Boolean value of true.

Any other scenario returns a zero [Condition] alongside a Boolean value
of false.
*/
func ConvertCondition(in any) (Condition, bool) {
	return conditionTypeAliasConverter(in)
}

/*
conditionTypeAliasConverter attempts to convert any (u) back to a bonafide instance
of Condition. This will only work if input value u is a type alias of Condition. An
instance of Condition is returned along with a success-indicative Boolean value.
*/
func conditionTypeAliasConverter(u any) (C Condition, converted bool) {
	if u != nil {
		// If it isn't a Condition alias, but is a
		// genuine Condition, just pass it back
		// with a thumbs-up ...
		if co, isCond := u.(Condition); isCond {
			C = co
			converted = isCond
			return
		}

		a, v, _ := derefPtr(typOf(u), valOf(u))
		b := typOf(Condition{}) // target (dest) type
		if a.ConvertibleTo(b) {
			X := v.Convert(b).Interface()
			if assert, ok := X.(Condition); ok {
				if !assert.IsZero() {
					C = assert
					converted = true
				}
			}
		}
	}

	return
}

const badCond = `<invalid_condition>`
