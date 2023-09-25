package stackage

/*
cond.go contains Condition-related methods and functions.
*/

import (
	"log"
)

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
several ComparisonOperator (op) instances made available through
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
  - Must be a valid instance of Stack (or an *alias* of Stack that is convertible back to the Stack type), OR ...
  - Must be a valid instance of any type that exports a stringer method (String()) intended to produce the Condition's final expression string representation

However when a PushPolicy function or method is added to an instance of this type,
greater control is afforded to the user in terms of what values will be accepted,
as well as the quality or state of such values.

Instances of this type -- similar to Stack instances -- MUST be initialized before use.
Initialization can occur as a result of executing the Cond package-level function, or
using the Init method extended through instances of this type. Initialization state may
be checked using the IsInit method.
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
	var data map[string]string = make(map[string]string, 0)

	r = new(condition)
	r.cfg = new(nodeConfig)
	r.cfg.log = newLogSystem(cLogDefault)
	r.cfg.log.lvl = logLevels(NoLogLevels)

	r.cfg.typ = cond
	data[`type`] = cond.String()

	if !logDiscard(r.cfg.log.log) {
		data[`laddr`] = sprintf("%p", r.cfg.log.log)
		data[`lpfx`] = r.cfg.log.log.Prefix()
		data[`lflags`] = sprintf("%d", r.cfg.log.log.Flags())
		data[`lwriter`] = sprintf("%T", r.cfg.log.log.Writer())
	}

	r.debug(`ALLOC`, data)

	return
}

/*
SetKeyword sets the receiver's keyword using the specified kw
input argument.
*/
func (r Condition) SetKeyword(kw any) Condition {
	if r.IsInit() {
		r.condition.setKeyword(kw)
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
SetOperator sets the receiver's comparison operator using the
specified Operator-qualifying input argument (op).
*/
func (r Condition) SetOperator(op Operator) Condition {
	if r.IsInit() {
		r.condition.setOperator(op)
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
specified ex input argument.
*/
func (r Condition) SetExpression(ex any) Condition {
	if r.IsInit() {
		r.condition.setExpression(ex)
	}
	return r
}

func (r *condition) setExpression(ex any) {
	if v, ok := r.assertConditionExpressionValue(ex); ok {
		r.ex = v
	}
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
func (r Condition) SetLogLevel(l ...any) Condition {
	if r.IsInit() {
		r.condition.setLogLevel(l...)
	}
	return r
}

func (r *condition) setLogLevel(l ...any) {
	r.calls(sprintf("%s: in:%T(%v;len:%d)",
		fmname(), l, l, len(l)))

	r.cfg.log.shift(l...)

	r.calls(sprintf("%s: out:%T(self)",
		fmname(), r))
}

/*
LogLevels returns the string representation of a comma-delimited list
of all active LogLevel values within the receiver.
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
UnsetLogLevel disables the specified LogLevel instance(s), thereby
instructing the logging subsystem to discard events submitted for
transcription to the underlying logger.
*/
func (r Condition) UnsetLogLevel(l ...any) Condition {
	if r.IsInit() {
		r.condition.unsetLogLevel(l...)
	}
	return r
}

func (r *condition) unsetLogLevel(l ...any) {
	fname := fmname()
	r.calls(sprintf("%s: in:%T(%v;len:%d)",
		fname, l, l, len(l)))

	r.cfg.log.unshift(l...)

	r.calls(sprintf("%s: out:%T(self)",
		fname, r))
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
func (r Condition) SetLogger(logger any) Condition {
	if r.IsInit() {
		r.condition.setLogger(logger)
	}

	return r
}

/*
Logger returns the *log.Logger instance. This can be used for quick
access to the log.Logger type's methods in a manner such as:

	r.Logger().Fatalf("We died")

It is not recommended to modify the return instance for the purpose
of disabling logging outright (see Stack.SetLogger method as well
as the SetDefaultConditionLogger package-level function for ways of
doing this easily).
*/
func (r Condition) Logger() (l *log.Logger) {
	if r.IsInit() {
		l = r.condition.logger()
	}
	return
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
SetCategory assigns the provided string to the receiver internal category value.
This allows for a means of identifying a particular kind of Condition in the midst
of many.
*/
func (r Condition) SetCategory(cat string) Condition {
	if r.IsInit() {
		r.condition.setCategory(cat)
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
Cond returns an instance of Condition bearing the provided component values.
This is intended to be used in situations where a Condition instance can be
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

This method may be useful in situations where a Condition will be assembled in a
"piecemeal" fashion (i.e.: incrementally), or if a Condition instance is slated to
be repurposed for use elsewhere (possibly in a repetative manner).
*/
func (r *Condition) Init() Condition {
	*r = Condition{condition: initCondition()}
	return *r
}

/*
SetID assigns the provided string value (or lack thereof) to the receiver.
This is optional, and is usually only needed in complex Condition structures
in which "labeling" certain components may be advantageous. It has no effect
on an evaluation, nor should a name ever cause a validity check to fail.

If the string `_random` is provided, a 24-character alphanumeric string is
randomly generated using math/rand and assigned as the ID.
*/
func (r Condition) SetID(id string) Condition {
	if r.IsInit() {
		if lc(id) == `_random` {
			id = randomID(randIDSize)
		} else if lc(id) == `_addr` {
			id = sprintf("%s", r.Addr())
		}

		r.condition.cfg.setID(id)
	}
	return r
}

/*
Len returns a "perceived" abstract length relating to the content (or lack
thereof) assigned to the receiver instance:

  - An uninitialized or zero instance returns zero (0)
  - An initialized instance with no Expression assigned (nil) returns zero (0)
  - A Stack or Stack type alias assigned as the Expression shall impose its own stack length as the return value (even if zero (0))

All other type instances assigned as an Expression shall result in a
return of one (1); this includes slice types, maps, arrays and any other
type that supports multiple values.

This capability was added to this type to mirror that of the Stack type in
order to allow additional functionality to be added to the Interface interface.
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
the IsZero method.
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

If a ValidityPolicy was set within the receiver, it shall be executed here.
If no ValidityPolicy was specified, only elements pertaining to basic viability
are checked.
*/
func (r Condition) Valid() (err error) {
	if !r.IsInit() {
		err = errorf("%T instance is nil", r)
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
		err = errorf("%T keyword value is zero", r)
		return
	}

	// verify comparison operator
	if cop := r.Operator(); cop != nil {
		if assert, ok := cop.(ComparisonOperator); ok {
			if !(1 <= int(assert) && int(assert) <= 6) {
				err = errorf("%T operator value is bogus", r)
				return
			}
		}
	}

	// verify expression value
	if r.Expression() == nil {
		err = errorf("%T expression value is nil", r)
	}

	return
}

/*
Evaluate uses the Evaluator closure function to apply the value (x)
to the receiver in order to conduct a matching/assertion test or
analysis for some reason. This is entirely up to the user.

A Boolean value returned indicative of the result. Note that if an
instance of Evaluator was not assigned to the Condition prior to
execution of this method, the return value shall always be false.
*/
func (r Condition) Evaluate(x ...any) (ev any, err error) {
	if r.IsInit() {
		if err = errorf("No %T.%T func/meth found", r, r.cfg.evl); r.cfg.evl != nil {
			ev, err = r.cfg.evl(x...)
		}
	}

	return
}

/*
SetEvaluator assigns the instance of Evaluator to the receiver. This
will allow the Evaluate method to return a more meaningful result.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetEvaluator(x Evaluator) Condition {
	if r.IsInit() {
		r.condition.cfg.evl = x
	}
	return r
}

/*
SetValidityPolicy assigns the instance of ValidityPolicy to the receiver.
This will allow the Valid method to return a more meaningful result.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetValidityPolicy(x ValidityPolicy) Condition {
	if r.IsInit() {
		r.condition.cfg.vpf = x
	}

	return r
}

/*
SetPresentationPolicy assigns the instance of PresentationPolicy to the receiver.
This will allow the user to leverage their own "stringer" method for automatic
use when this type's String method is called.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetPresentationPolicy(x PresentationPolicy) Condition {
	if r.IsInit() {
		r.condition.cfg.rpf = x
	}

	return r
}

/*
Encap accepts input characters for use in controlled condition value
encapsulation. Acceptable input types are:

A single string value will be used for both L and R encapsulation.

An instance of []string with two (2) values will be used for L and R
encapsulation using the first and second slice values respectively.

An instance of []string with only one (1) value is identical to the act of
providing a single string value, in that both L and R will use one value.
*/
func (r Condition) Encap(x ...any) Condition {
	if r.IsInit() {
		r.condition.cfg.setEncap(x...)
	}
	return r
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
NoNesting sets the no-nesting bit within the receiver. If
set to true, the receiver shall ignore any Stack or Stack
type alias instance when assigned using the SetExpression
method. In such a case, only primitives, etc., shall be
honored during the SetExpression operation.

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the nesting bit (i.e.: true->false and
false->true)
*/
func (r Condition) NoNesting(state ...bool) Condition {
	r.setState(nnest, state...)
	return r
}

/*
CanNest returns a Boolean value indicative of whether
the no-nesting bit is unset, thereby allowing a Stack
or Stack type alias instance to be set as the value.

See also the IsNesting method.
*/
func (r Condition) CanNest() (can bool) {
	if r.IsInit() {
		can = !r.condition.positive(nnest)
	}
	return
}

/*
IsNesting returns a Boolean value indicative of whether the
underlying expression value is either a Stack or Stack type
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
receiver instance's Expression value represents a Stack (or Stack
type alias) instance which exhibits First-In-First-Out behavior as
it pertains to the act of appending and truncating the receiver's
slices.

A value of false implies that no such Stack instance is set as the
expression, OR that the Stack exhibits Last-In-Last-Out behavior,
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

	fname := fmname()
	r.calls(sprintf("%s: in:niladic", fname))

	if stk, ok := stackTypeAliasConverter(r.ex); ok {
		result = stk.IsFIFO()
	}

	r.calls(sprintf("%s: out:%T(%v)",
		fmname(), result, result))

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
func (r Condition) Paren(state ...bool) Condition {
	r.setState(parens, state...)
	return r
}

/*
IsParen returns a Boolean value indicative of whether the
receiver is parenthetical.
*/
func (r Condition) IsParen() bool {
	return r.getState(parens)
}

func (r Condition) getState(cf cfgFlag) (state bool) {
	if r.IsInit() {
		state = r.condition.positive(cf)
	}
	return
}

func (r Condition) setState(cf cfgFlag, state ...bool) {
	if r.IsInit() {
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

func (r *condition) setLogger(logger any) {
	r.cfg.log.setLogger(logger)
}

func (r condition) logger() *log.Logger {
	return r.cfg.log.logger()
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
NoPadding sets the no-space-padding bit within the receiver.
String values within the receiver shall not be padded using
a single space character (ASCII #32).

A Boolean input value explicitly sets the bit as intended.
Execution without a Boolean input value will *TOGGLE* the
current state of the quotation bit (i.e.: true->false and
false->true)
*/
func (r Condition) NoPadding(state ...bool) Condition {
	r.setState(nspad, state...)
	return r
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
Operator returns the Operator interface type instance found within the
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
such as a Stack or a Go primitive, this method may be uncertain
as to what it should do. A bogus string may be returned.

In such a case, it may be necessary to subvert the default string
representation behavior demonstrated by instances of this type in
favor of a custom instance of the PresentationPolicy closure type
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

	s := sprintf("%s%s%s%s%s", r.kw, pad, r.op, pad, val)
	if r.cfg.positive(parens) {
		s = sprintf("(%s%s%s)", pad, s, pad)
	}

	return s
}

/*
mkmsg is the private method called by eventDispatch for the
purpose of Message assembly prior to submission to a logger.
*/
func (r *condition) mkmsg(typ string) (m Message, ok bool) {
	if r.isInit() {
		if len(typ) > 0 {
			m = Message{
				ID:   getLogID(r.cfg.id),
				Tag:  typ,
				Addr: ptrString(r),
				Type: sprintf("%T", *r),
				Time: timestamp(),
				Len:  -1, // initially N/A
				Cap:  -1, // initially N/A
			}
			ok = true
		}
	}

	return
}

/*
error conditions that are fatal and always serious
*/
func (r *condition) fatal(x any, data ...map[string]string) {
	if r != nil {
		r.eventDispatch(x, LogLevel5, `FATAL`, data...)
	}
}

/*
error conditions that are not fatal but potentially serious
*/
func (r *condition) error(x any, data ...map[string]string) {
	if r != nil {
		r.eventDispatch(x, LogLevel5, `ERROR`, data...)
	}
}

/*
extreme depth operational details
*/
func (r *condition) trace(x any, data ...map[string]string) {
	if r != nil {
		r.eventDispatch(x, LogLevel6, `TRACE`, data...)
	}
}

/*
relatively deep operational details
*/
func (r *condition) debug(x any, data ...map[string]string) {
	if r != nil {
		r.eventDispatch(x, LogLevel4, `DEBUG`, data...)
	}
}

/*
policy method operational details, as well as caps, r/o, etc.
*/
func (r *condition) policy(x any, data ...map[string]string) {
	if r != nil {
		r.eventDispatch(x, LogLevel2, `POLICY`, data...)
	}
}

/*
calls records in/out signatures and realtime meta-data regarding
individual method runtimes.
*/
func (r *condition) calls(x any, data ...map[string]string) {
	if r != nil {
		r.eventDispatch(x, LogLevel1, `CALL`, data...)
	}
}

/*
state records interrogations of, and changes to, the underlying
configuration value.
*/
func (r *condition) state(x any, data ...map[string]string) {
	if r != nil {
		r.eventDispatch(x, LogLevel3, `STATE`, data...)
	}
}

/*
eventDispatch is the main dispatcher of events of any severity.
A severity of FATAL (in any case) will result in a logger-driven
call of os.Exit.
*/
func (r condition) eventDispatch(x any, ll LogLevel, severity string, data ...map[string]string) {
	if !(r.cfg.log.positive(ll) ||
		eq(severity, `FATAL`) ||
		r.cfg.log.lvl == logLevels(AllLogLevels)) {
		return
	}

	printers := map[string]func(...any){
		`FATAL`:  r.logger().Fatalln,
		`STATE`:  r.logger().Println,
		`ERROR`:  r.logger().Println,
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

const badCond = `<invalid_condition>`
