package stackage

import (
	"log"
)

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
	r.cfg.log.lvl = logLevels(cLogLevelDefault)

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
	if !r.IsInit() {
		return r
	}

	r.condition.setKeyword(kw)
	return r
}

func (r *condition) setKeyword(kw any) {
	if r == nil {
		r = initCondition()
	}

	switch tv := kw.(type) {
	case string:
		r.kw = tv
	default:
		r.kw = assertKeyword(tv)
	}
}

/*
SetOperator sets the receiver's comparison operator using the
specified Operator-qualifying input argument (op).
*/
func (r Condition) SetOperator(op Operator) Condition {
	if !r.IsInit() {
		return r
	}

	r.condition.setOperator(op)
	return r
}

func (r *condition) setOperator(op Operator) {
	if r == nil {
		r = initCondition()
	}

	if len(op.Context()) > 0 && len(op.String()) > 0 {
		r.op = op
	}
}

/*
SetExpression sets the receiver's expression value(s) using the
specified ex input argument.
*/
func (r Condition) SetExpression(ex any) Condition {
	if !r.IsInit() {
		return r
	}

	r.condition.setExpression(ex)
	return r
}

func (r *condition) setExpression(ex any) {
	if r == nil {
		r = initCondition()
	}

	if v, ok := r.assertConditionExpressionValue(ex); ok {
		r.ex = v
	}
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
	if r.IsZero() {
		return r
	}

	r.condition.setLogger(logger)
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
func (r Condition) Logger() *log.Logger {
	if r.IsZero() {
		return nil
	}
	return r.condition.logger()
}

/*
assertKeyword returns a valid keyword string value based on the input value.
If the value is a string, it is returned as-is. Else, if a custom type that
possesses its own stringer method, the returned value from that method is
returned.
*/
func assertKeyword(x any) (s string) {
	switch tv := x.(type) {
	case string:
		s = tv
	default:
		if meth := getStringer(x); meth != nil {
			s = meth()
		}
	}

	return
}

/*
tryPushPolicy will execute a push policy if one is set, and will return an error
and a Boolean value indicative of the presence of said policy.
*/
func (r *condition) tryPushPolicy(x any) (err error, found bool) {
	if found = r.cfg.ppf != nil; found {
		// if we have a policy, always set
		// found to true, regardless of the
		// outcome
		err = r.cfg.ppf(x)
	}

	// if no policy, err is always nil
	// and found is always false
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
func (r Condition) Err() error {
	if !r.IsInit() {
		return nil
	}
	return r.condition.getErr()
}

/*
SetErr sets the underlying error value within the receiver
to the assigned input value err, whether nil or not.

This method may be most valuable to users who have chosen
to extend this type by aliasing, and wish to control the
handling of error conditions in another manner.
*/
func (r Condition) SetErr(err error) Condition {
	if !r.IsInit() {
		return r
	}
	r.condition.setErr(err)
	return r
}

/*
setErr assigns an error instance, whether nil or not, to
the underlying receiver configuration.
*/
func (r *condition) setErr(err error) {
	if r.cfg == nil {
		return
	}
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

	// Try to find a push policy first and, IF
	// FOUND, run it and break out of the case
	// statement either way.
	if err, found := r.tryPushPolicy(x); found {
		X = x
		return
	} else if err != nil {
		r.setErr(err)
		return
	}

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

	} else if meth := getStringer(x); meth != nil {
		// whatever it is, it seems to have
		// a stringer method, at least ...
		X = x

	} else if isKnownPrimitive(x) {
		// value is one of go's builtin
		// numerical primitives, which
		// are string represented using
		// sprintf.
		X = primitiveStringer(x)
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
	if !r.IsInit() {
		return r
	}

	r.condition.setCategory(cat)
	return r
}

/*
Category returns the categorical label string value assigned to the receiver, if
set, else a zero string.
*/
func (r Condition) Category() string {
	if r.IsZero() {
		return ``
	}
	return r.condition.getCat()
}

/*
Cond returns an instance of Condition bearing the provided component values.
This is intended to be used in situations where a Condition instance can be
created in one shot.
*/
func Cond(kw any, op Operator, ex any) Condition {
	return Condition{newCondition(kw, op, ex)}
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
	if !r.IsInit() {
		return r
	}

	if lc(id) == `_random` {
		id = randomID(randIDSize)
	} else if lc(id) == `_addr` {
		id = sprintf("%s", r.condition.addr())
	}

	r.condition.cfg.setID(id)
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
func (r Condition) Addr() string {
	if !r.IsInit() {
		return ``
	}

	return r.condition.addr()
}

/*
addr returns the string representation of the pointer
address for the receiver.
*/
func (r condition) addr() string {
	if r.isZero() {
		return ``
	}

	return sprintf("%p", &r)
}

/*
Name returns the name of the receiver instance, if set, else a zero string
will be returned. The presence or lack of a name has no effect on any of
the receiver's mechanics, and is strictly for convenience.
*/
func (r Condition) ID() string {
	return r.condition.getID()
}

func (r condition) getID() string {
	if r.isZero() {
		return ``
	}
	return r.cfg.id
}

/*
IsInit will verify that the internal pointer instance of the receiver has
been properly initialized. This method executes a preemptive execution of
the IsZero method.
*/
func (r Condition) IsInit() bool {
	if r.IsZero() {
		return false
	}

	return r.condition.isInit()
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
	if r.condition == nil {
		return true
	}

	return r.condition.isZero()
}

func (r *condition) isZero() bool {
	if r == nil {
		return true
	} else if r.cfg == nil {
		return true
	}

	return false
}

/*
Valid returns an instance of error, identifying any issues perceived with
the state of the receiver.

If a ValidityPolicy was set within the receiver, it shall be executed here.
If no ValidityPolicy was specified, only a nilness is checked
*/
func (r Condition) Valid() (err error) {
	if r.condition.isZero() {
		err = errorf("%T instance is nil", r)
		return
	}

	if err = r.condition.getErr(); err != nil {
		return
	}

	// if a validitypolicy was provided, use it
	if r.condition.cfg.vpf != nil {
		err = r.condition.cfg.vpf(r)
		return
	}

	// no validity policy was provided, just check
	// what we can.
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
func (r Condition) Evaluate(x ...any) error {
	if r.cfg.evl != nil {
		return r.cfg.evl(r, x...)
	}

	return errorf("No %T function or method was set within %T")
}

/*
SetEvaluator assigns the instance of Evaluator to the receiver. This
will allow the Evaluate method to return a more meaningful result.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetEvaluator(x Evaluator) Condition {
	if !r.IsInit() {
		return r
	}

	r.condition.cfg.evl = x
	return r
}

/*
SetValidityPolicy assigns the instance of ValidityPolicy to the receiver.
This will allow the Valid method to return a more meaningful result.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetValidityPolicy(x ValidityPolicy) Condition {
	if !r.IsInit() {
		return r
	}

	r.condition.cfg.vpf = x
	return r
}

/*
SetPresentationPolicy assigns the instance of PresentationPolicy to the receiver.
This will allow the user to leverage their own "stringer" method for automatic
use when this type's String method is called.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetPresentationPolicy(x PresentationPolicy) Condition {
	if !r.IsInit() {
		return r
	}

	r.condition.cfg.rpf = x
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
	if r.condition == nil {
		r.condition = initCondition()
	}

	r.condition.cfg.setEncap(x...)
	return r
}

/*
IsEncap returns a Boolean value indicative of whether value encapsulation
characters have been set within the receiver.
*/
func (r Condition) IsEncap() bool {
	if r.IsZero() {
		return false
	}

	return len(r.condition.getEncap()) > 0
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
	if !r.IsInit() {
		return r
	}

	if len(state) > 0 {
		if state[0] {
			r.condition.setOpt(nnest)
		} else {
			r.condition.unsetOpt(nnest)
		}
	} else {
		r.condition.toggleOpt(nnest)
	}

	return r
}

/*
CanNest returns a Boolean value indicative of whether
the no-nesting bit is unset, thereby allowing a Stack
or Stack type alias instance to be set as the value.

See also the IsNesting method.
*/
func (r Condition) CanNest() bool {
	if r.IsZero() {
		return false
	}

	return !r.condition.positive(nnest)
}

/*
IsNesting returns a Boolean value indicative of whether the
underlying expression value is either a Stack or Stack type
alias. If true, this indicates the expression value descends
into another hierarchical (nested) context.
*/
func (r Condition) IsNesting() bool {
	if r.IsZero() {
		return false
	}

	return r.condition.isNesting()
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
func (r Condition) IsFIFO() bool {
	if r.condition == nil {
		r.condition = initCondition()
		return false
	}
	return r.condition.isFIFO()
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
	if r.condition == nil {
		r.condition = initCondition()
	}

	if len(state) > 0 {
		if state[0] {
			r.condition.setOpt(parens)
		} else {
			r.condition.unsetOpt(parens)
		}
	} else {
		r.condition.toggleOpt(parens)
	}

	return r
}

/*
IsParen returns a Boolean value indicative of whether the
receiver is parenthetical.
*/
func (r Condition) IsParen() bool {
	if r.IsZero() {
		return false
	}

	return r.condition.positive(parens)
}

func (r *condition) setLogger(logger any) {
	r.cfg.log.setLogger(logger)
}

func (r condition) logger() *log.Logger {
	return r.cfg.log.logger()
}

func (r *condition) toggleOpt(x cfgFlag) {
	if r == nil {
		r = initCondition()
	}

	r.cfg.toggleOpt(x)
}

func (r *condition) setOpt(x cfgFlag) {
	if r.isZero() || r == nil {
		r = initCondition()
	}

	r.cfg.setOpt(x)
}

func (r *condition) unsetOpt(x cfgFlag) {
	if r == nil {
		r = initCondition()
	}

	r.cfg.unsetOpt(x)
}

func (r *condition) positive(x cfgFlag) bool {
	if r == nil {
		r = initCondition()
	}
	return r.cfg.positive(x)
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
	if r.condition == nil {
		r.condition = initCondition()
	}

	if len(state) > 0 {
		if state[0] {
			r.condition.setOpt(nspad)
		} else {
			r.condition.unsetOpt(nspad)
		}
	} else {
		r.condition.toggleOpt(nspad)
	}

	return r
}

/*
IsPadded returns a Boolean value indicative of whether the
receiver pads its contents with a SPACE char (ASCII #32).
*/
func (r Condition) IsPadded() bool {
	if r.IsZero() {
		return false
	}

	return !r.condition.positive(nspad)
}

/*
SetPushPolicy assigns the instance of PushPolicy to the receiver. This
will allow the Set method to control what elements may (or may not) be
set as the expression value within the receiver.

See the documentation for the Set method for information on the default
behavior without the involvement of a PushPolicy instance.

Specifying nil shall disable this capability if enabled.
*/
func (r Condition) SetPushPolicy(x PushPolicy) Condition {
	if !r.IsInit() {
		return r
	}

	r.condition.cfg.ppf = x
	return r
}

/*
Expression returns the expression value(s) stored within the receiver, or
nil if unset. A valid receiver instance MUST always possess a non-nil
expression value.
*/
func (r Condition) Expression() any {
	if r.IsZero() {
		return nil
	}

	return r.condition.ex
}

/*
Operator returns the Operator interface type instance found within the
receiver.
*/
func (r Condition) Operator() Operator {
	if r.IsZero() {
		return nco
	}
	return r.condition.op
}

/*
Keyword returns the Keyword interface type instance found within the
receiver.
*/
func (r Condition) Keyword() string {
	if r.IsZero() {
		return ``
	}

	return r.condition.kw
}

/*
String is a stringer method that returns the string representation
of the receiver instance. It will only function if the receiver is
in good standing, and passes validity checks.
*/
func (r Condition) String() string {
	if err := r.Valid(); err != nil {
		return badCond
	}
	return r.condition.string()
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

	val := encapValue(r.cfg.enc, sprintf("%s", r.ex))
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
func (r *condition) mkmsg(typ string) (Message, bool) {
	if r.isZero() || len(typ) == 0 {
		return Message{}, false
	}

	return Message{
		ID:   getLogID(r.getID()),
		Tag:  typ,
		Addr: r.addr(),
		Type: sprintf("%T", *r),
		Time: timestamp(),
		Len:  -1, // initially N/A
		Cap:  -1, // initially N/A
	}, true
}

/*
error conditions that are fatal and always serious
*/
func (r condition) fatal(x any, data ...map[string]string) {
	r.eventDispatch(x, LogLevel5, `FATAL`, data...)
}

/*
error conditions that are not fatal but potentially serious
*/
func (r condition) error(x any, data ...map[string]string) {
	r.eventDispatch(x, LogLevel5, `ERROR`, data...)
}

/*
extreme depth operational details
*/
func (r condition) trace(x any, data ...map[string]string) {
	r.eventDispatch(x, LogLevel6, `TRACE`, data...)
}

/*
relatively deep operational details
*/
func (r condition) debug(x any, data ...map[string]string) {
	r.eventDispatch(x, LogLevel4, `DEBUG`, data...)
}

/*
policy method operational details, as well as caps, r/o, etc.
*/
func (r condition) policy(x any, data ...map[string]string) {
	r.eventDispatch(x, LogLevel2, `POLICY`, data...)
}

/*
calls records in/out signatures and realtime meta-data regarding
individual method runtimes.
*/
func (r condition) calls(x any, data ...map[string]string) {
	r.eventDispatch(x, LogLevel1, `CALL`, data...)
}

/*
state records interrogations of, and changes to, the underlying
configuration value.
*/
func (r condition) state(x any, data ...map[string]string) {
	r.eventDispatch(x, LogLevel3, `STATE`, data...)
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

	m, ok := r.mkmsg(severity)
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

	if eq(severity, `fatal`) {
		r.logger().Fatalln(m)
	}

	r.logger().Println(m)
}

const badCond = `<invalid_condition>`
