package stackage

/*
fcf.go contains first-class (closure) function signature and
interface definitions.
*/

/*
Evaluator is a first-class function signature type which may
be leveraged by users in order to compose matching functions
by which a given Condition shall be measured/compared, etc.

In other words, this type takes this package one step forward:
it no longer merely creates various expressions in the abstract
sense -- now it will actually *apply* them to real values, or
to gauge their verisimilitude in some other manner. The catch
is that the user will need to author the needed functions in
order to make such determinations practical.

Use of this feature is totally optional, and may be overkill
for most users.

The signature allows for variadic input of any type(s), and
shall require the return of a single any instance alongside
an error. The contents of the return 'any' instance (if not
nil) is entirely up to the user.
*/
type Evaluator func(...any) (any, error)

/*
PushPolicy is a first-class (closure) function signature
that may be leveraged by users in order to control what
types instances may be pushed into a Stack instance when
using its 'Push' method.

When authoring functions or methods that conform to this
signature, the idea is to return true for any value that
should be pushed, and false for all others.  This allows
for an opportunity to interdict potentially undesirable
Stack additions, unsupported types, etc.

A PushPolicy function or method is executed for each
element being added to a Stack via its Push method.
*/
type PushPolicy func(...any) error

/*
ValidityPolicy is a first-class (closure) function signature
that may be leveraged by users in order to better gauge the
validity of a stack based on its configuration and/or values.

A ValidityPolicy function or method is executed via the Stack
method Valid.
*/
type ValidityPolicy func(...any) error

/*
PresentationPolicy is a first-class (closure) function signature
that may be leveraged by users in order to better control value
presentation during the string representation process. Essentially,
one may write their own "stringer" (String()) function or method
and use it to override the default behavior of the package based
String method(s).

Note that basic Stack instances are ineligible for the process of
string representation, thus no PresentationPolicy may be set.
*/
type PresentationPolicy func(...any) string
