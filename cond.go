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
several ComparisonOperator (op) instances made available through
this package:

• Eq, or "equal to" (=)

• Ne, or "not equal to" (!=)	// USE WITH CAUTION!!

• Lt, or "less than" (<)

• Le, or "less than or equal" (<=)

• Gt, or "greater than" (>)

• Ge, or "greater than or equal" (>=)

... OR through a user-defined operator that conforms to the package
defined Operator interface.

By default, permitted expression (ex) values must honor these guidelines:

• Must be a non-zero string, OR ...

• Must be a valid instance of Stack (or an *alias* of Stack that is convertible back
to the Stack type), OR ...

• Must be a valid instance of any type that exports a stringer method (String())
intended to produce the Condition's final expression string representation

However when a PushPolicy function or method is added to an instance of this type,
greater control is afforded to the user in terms of what values will be accepted,
as well as the quality or state of such values.
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
newCondition initializes, (optionally sets) and returns a new instance of
*condition in one shot.
*/
func newCondition(kw any, op Operator, ex any) (r *condition) {
	r = initCondition()

	r.setKeyword(kw)    // keyword
	r.setOperator(op)   // operator
	r.setExpression(ex) // expr. value(s)

	return
}

func initCondition() (r *condition) {
	r = new(condition)
	r.cfg = new(nodeConfig)
	r.cfg.typ = cond

	return
}

/*
SetKeyword sets the receiver's keyword using the specified kw
input argument.
*/
func (r *Condition) SetKeyword(kw any) *Condition {
	if r.condition == nil {
		r.condition = initCondition()
	}

	r.condition.setKeyword(kw)
	return r
}

func (r *condition) setKeyword(kw any) {
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
func (r *Condition) SetOperator(op Operator) *Condition {
	if r.condition == nil {
		r.condition = initCondition()
	}

	r.condition.setOperator(op)
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
func (r *Condition) SetExpression(ex any) *Condition {
	if r.condition == nil {
		r.condition = initCondition()
	}

	r.condition.setExpression(ex)
	return r
}

func (r *condition) setExpression(ex any) {
	if v, ok := r.assertConditionExpressionValue(ex); ok {
		r.ex = v
	}
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
	if r.cfg.ppf != nil {
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
setError assigns (or clears) the underlying configuration error value.
*/
func (r *condition) setError(err error) {
	r.cfg.setError(err)
}

/*
isError returns a Boolean value indicative of whether the underlying
configuration contains a non-nil error instance.
*/
func (r condition) isError() bool {
	return r.error() != nil
}

/*
error returns the underlying error instance, whether nil or not.
*/
func (r condition) error() error {
	return r.cfg.err
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
func (r *condition) defaultAssertionExpressionHandler(x any) (X any) {
	// Try to find a push policy first and, IF
	// FOUND, run it and break out of the case
	// statement either way.
	if err, found := r.tryPushPolicy(x); err == nil {
		X = x
		return
	} else if err != nil {
		r.setError(err)
		return
	} else if found {
		return
	}

	// no push policy, so we'll see if the basic
	// guidelines were satisfied, at least ...
	if v, ok := stackTypeAliasConverter(x); ok {
		// a user-created type alias of Stack
		// was converted back to Stack without
		// any issues ...
		X = v
	} else if meth := getStringer(x); meth != nil {
		// whatever it is, it seems to have
		// a stringer method, at least ...
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
	if r.condition == nil {
		r.condition = initCondition()
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
Cond returns an instance of *Condition bearing the provided component values.
*/
func Cond(kw any, op Operator, ex any) Condition {
	return Condition{newCondition(kw, op, ex)}
}

/*
SetID assigns the provided string value (or lack thereof) to the receiver.
This is optional, and is usually only needed in complex Condition structures
in which "labeling" certain components may be advantageous. It has no effect
on an evaluation, nor should a name ever cause a validity check to fail.
*/
func (r Condition) SetID(n string) Condition {
	if r.condition == nil {
		r.condition = initCondition()
	}

	r.condition.cfg.setID(n)
	return r
}

/*
Name returns the name of the receiver instance, if set, else a zero string
will be returned. The presence or lack of a name has no effect on any of
the receiver's mechanics, and is strictly for convenience.
*/
func (r Condition) ID() string {
	return r.condition.cfg.id
}

/*
IsZero returns a Boolean value indicative of whether the receiver is nil,
or unset.
*/
func (r Condition) IsZero() bool {
	return r.condition.isZero()
}

func (r *condition) isZero() bool {
	if r == nil {
		return true
	}

	return len(r.kw) == 0 &&
		r.op == nil &&
		r.ex == nil &&
		r.cfg == nil
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

	if r.condition.isError() {
		err = r.condition.error()
		return
	}

	// if a validitypolicy was provided, use it
	if r.condition.cfg.vpf != nil {
		err = r.condition.cfg.vpf()
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
	if r.condition == nil {
		r.condition = initCondition()
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
	if r.condition == nil {
		r.condition = initCondition()
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
	if r.condition == nil {
		r.condition = initCondition()
	}

	r.condition.cfg.rpf = x
	return r
}

/*
Encap accepts input characters for use in controlled condition value
encapsulation. Acceptable input types are:

• string - a single string value will be used for both L and R
encapsulation.

• string slices - An instance of []string with two (2) values will
be used for L and R encapsulation using the first and second
slice values respectively. An instance of []string with only one (1)
value is identical to providing a single string value, in that both
L and R will use one value.
*/
func (r Condition) Encap(x ...any) Condition {
	if r.condition == nil {
		r.condition = initCondition()
	}

	r.condition.cfg.setEncap(x...)
	return r
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

	return r.cfg.positive(parens)
}

func (r *condition) toggleOpt(x cfgFlag) {
	r.cfg.toggleOpt(x)
}

func (r *condition) setOpt(x cfgFlag) {
	r.cfg.setOpt(x)
}

func (r *condition) unsetOpt(x cfgFlag) {
	r.cfg.unsetOpt(x)
}

/*
Padding sets the no-space-padding bit within the receiver.
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
			r.condition.cfg.setOpt(nspad)
		} else {
			r.condition.cfg.unsetOpt(nspad)
		}
	} else {
		r.condition.cfg.toggleOpt(nspad)
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

	return !r.cfg.positive(nspad)
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
	if r.condition == nil {
		r.condition = initCondition()
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
		return r.cfg.rpf()
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

func (r condition) msg(x any) (m Message) {
	m.Tag = `UNKNOWN`

	switch tv := x.(type) {
	case error:
		if tv != nil {
			m.Tag = `ERROR`
			m.Msg = tv.Error()
		}
	case string:
		if len(tv) > 0 {
			m.Tag = `DEBUG`
			m.Msg = tv
		}
	}

	m.Type = `C`
	if cat := r.getCat(); len(cat) > 0 {
		m.Type += sprintf("_%s", cat)
	}

	m.Time = now()
	m.ID = r.cfg.id
	m.Len = -1 // N/A
	m.Cap = -1 // N/A
	m.Addr = sprintf("%p", r)

	return
}

const badCond = `<invalid_condition>`
