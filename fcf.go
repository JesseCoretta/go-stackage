package stackage

/*
fcf.go contains first-class (closure) function signature and
interface definitions.
*/

/*
Evaluator is a first-class function signature type which may
be leveraged by users in order to compose matching functions
by which a given [Condition] shall be measured/compared, etc.

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
types instances may be pushed into a [Stack] instance when
using its [Stack.Push] method.

When authoring functions or methods that conform to this
signature, the idea is to return true for any value that
should be pushed, and false for all others.  This allows
for an opportunity to interdict potentially undesirable
[Stack] additions, unsupported types, etc.

A PushPolicy function or method is executed for each element
being added to a [Stack] via its [Stack.Push] method.
*/
type PushPolicy func(...any) error

/*
ValidityPolicy is a first-class (closure) function signature
that may be leveraged by users in order to better gauge the
validity of a [Stack] based on its configuration and/or values.

A ValidityPolicy function or method is executed via the
[Stack.Valid] method.
*/
type ValidityPolicy func(...any) error

/*
PresentationPolicy is a first-class (closure) function signature
that may be leveraged by users in order to better control value
presentation during the string representation process. Essentially,
one may write their own "stringer" (String()) function or method
and use it to override the default behavior of the package based
String method(s).

Note that basic [Stack] instances are ineligible for the process of
string representation, thus no PresentationPolicy may be set.
*/
type PresentationPolicy func(...any) string

/*
EqualityPolicy is a first-class (closure) function signature that may be
leveraged by users in order to gain full control over the equality assertion
mechanism between two (2) [Stack] or [Stack]-alias instances, or two (2)
[Condition] or [Condition]-alias instances.

An EqualityPolicy may be set, or unset, using the [Stack.SetEqualityPolicy] and
[Condition.SetEqualityPolicy] methods as needed.  When the [Stack.IsEqual] or
[Condition.IsEqual] methods are called, the appropriate EqualityPolicy will be
invoked.

By default, when no custom EqualityPolicy is specified, the receiver will
compare itself to the input value (any) provided as described in the notes
of the [Stack.IsEqual] and [Condition.IsEqual] methods. It is noted that,
given sufficiently complex or large structures, the default mechanism is
fairly costly. A user-implemented alternative closure, when assigned to the
instance(s) in question, can mitigate this drawback through more selective
and fine-tuned assertion processes.
*/
type EqualityPolicy func(any, any) error

/*
Unmarshaler is a first-class (closure) function signature that may be
leveraged by users in order to gain full control over the unmarshaling
process of a [Stack] or [Stack]-alias instance and all of its contents
therein.

An Unmarshaler may be set, or unset, using the [Stack.SetUnmarshaler] and
[Condition.SetUnmarshaler] methods as needed. When the [Stack.Unmarshal]
method is called, it invokes the appropriate Unmarshaler.

By default, when no custom Unmarshaler is specified, the return type is
[]any, which is the type that is used to mimic any form of [Stack]. Any
[Condition] and [Condition] instances present are also mimicked, except
the copied instance -- if it contains a [Stack] or [Stack]-alias as the
[Condition.Expression] -- will have its [Stack] or [Stack]-alias swapped
with an instance of []any in recursive manner.  Note that none of the
original content is modified as as a result of this process, and when
completed, the user will be left with the original instance alongside a
nearly identical mimicry that may be freely dissected or modified as needed.

Users authoring their own Unmarshaler need only utilize the return []any
as an outer-envelope only -- the enveloped value(s) within may be of any
number and of any combination of types, assembled or [re]-structured in
any way desired.
*/
type Unmarshaler func(...any) ([]any, error)

/*
Marshaler is a first-class (closure) function signature that may be
leveraged by users in order to gain full control over the marshaling
process that deposits content into a [Stack] or [Stack]-alias instance.
*/
type Marshaler func(...any) (err error)

/*
LessFunc qualifies for [sort.Interface].
*/
type LessFunc func(int, int) bool
