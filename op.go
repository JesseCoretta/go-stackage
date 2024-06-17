package stackage

/*
op.go contains well-known comparison operators that will
be used in the expression of a given condition, and also
provides the framework for extending operator concepts as
a whole.
*/

/*
Operator is an interface type that allows user-defined operators to be used within instances
of [Condition]. In rare cases, users may wish to utilize operators that go beyond the package
provided [ComparisonOperator] definitions (or just represent the same operators in a different
way). Defining types that conform to the signature of this interface type allows just that.
*/
type Operator interface {
	// String should return the preferred string
	// representation of the Operator instance.
	// This is ultimately the value that shall be
	// used during the string representation of
	// the instance of Condition to which the
	// Operator is assigned. Generally, this will
	// be something short and succinct (e.g.: `~=`)
	// but can conceivably be anything you want.
	String() string

	// Context returns the string representation
	// of the context behind the operator. As an
	// example, the Context for an instance of
	// ComparisonOperator is `comparison`. Users
	// should choose intuitive, helpful context
	// names for custom types when defining them.
	Context() string
}

/*
ComparisonOperator constants are intended to be used in singular
form when evaluating two (2) particular values.
*/
const (
	nco ComparisonOperator = iota // 0 <invalid_operator>
	Eq                            // 1 (=)
	Ne                            // 2 (!=)
	Lt                            // 3 (<)
	Gt                            // 4 (>)
	Le                            // 5 (<=)
	Ge                            // 6 (>=)
)

const badOp = `<invalid_operator>`
const compOpCtx = `comparison`

/*
ComparisonOperator is a uint8 enumerated type used for abstract
representation of the following well-known and package-provided
operators:

  - Equal (=)
  - Not Equal (!=)
  - Less Than (<)
  - Greater Than (>)
  - Less Than Or Equal (<=)
  - Greater Than Or Equal (>=)

Instances of this type should be passed to [Cond] as the 'op' value.
*/
type ComparisonOperator uint8

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r ComparisonOperator) String() (op string) {
	op = badOp

	switch r {
	case Eq:
		op = `=`
	case Ne:
		op = `!=`
	case Lt:
		op = `<`
	case Gt:
		op = `>`
	case Le:
		op = `<=`
	case Ge:
		op = `>=`
	}

	return
}

/*
Context returns the contextual label associated with instances of
this type as a string value.
*/
func (r ComparisonOperator) Context() string {
	return compOpCtx
}
