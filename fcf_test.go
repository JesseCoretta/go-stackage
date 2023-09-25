package stackage

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

/*
This example demonstrates a custom user-authored stringer that
assumes full responsibility for the string representation of a
Stack instance.
*/
func ExampleStack_SetPresentationPolicy() {
	// slices we intend to push
	// into our stack.
	slices := []any{
		`safe_text`,
		`sensitive_text`,
		1.261,
		`safe_text`,
		`sensitive_text`,
		[]string{`health`, `your`, `for`, `bad`, `is`, `smoking`, `ziggurat`},
		`irrelevant_text`,
		`sensitive_text`,
		`safe_text`,
		[]string{`homes.`, `feelings,`, `your`, `about`, `talk`, `to`, `need`, `You`},
	}

	// create a List stack, push the above slices
	// in variadic form and create/set a ppolicy.
	myStack := List()
	myStack.Push(slices...).SetPresentationPolicy(func(x ...any) string {
		var retval []string
		for i := 0; i < myStack.Len(); i++ {
			slice, _ := myStack.Index(i)
			switch tv := slice.(type) {
			case string:
				if tv == `safe_text` {
					retval = append(retval, tv+` (vetted)`) // declare safe
				} else if tv == `sensitive_text` {
					retval = append(retval, `[REDACTED]`) // redact
				} else {
					retval = append(retval, tv) // as-is
				}
			case []string:
				var message []string
				for j := len(tv); j > 0; j-- {
					message = append(message, tv[j-1])
				}
				retval = append(retval, strings.Join(message, string(rune(32))))
			default:
				retval = append(retval, fmt.Sprintf("%v (%T)", tv, tv))
			}
		}

		// Since we're doing our own stringer, we can't rely on the
		// SetDelimiter method (nor any method used in the string
		// representation process).
		return strings.Join(retval, ` || `)
	})

	fmt.Printf("%s", myStack)
	// output: safe_text (vetted) || [REDACTED] || 1.261 (float64) || safe_text (vetted) || [REDACTED] || ziggurat smoking is bad for your health || irrelevant_text || [REDACTED] || safe_text (vetted) || You need to talk about your feelings, homes.
}

/*
This example demonstrates a custom user-authored stringer that
assumes full responsibility for the string representation of a
Condition instance.

Here, we import "encoding/base64" to decode a string into the
true underlying value prior to calling fmt.Sprintf to present
and return the preferred string value.
*/
func ExampleCondition_SetPresentationPolicy() {
	// create a Condition and set a custom
	// ppolicy to handle base64 encoded
	// values.
	myCondition := Cond(`keyword`, Eq, `MTM4NDk5`)

	myCondition.SetPresentationPolicy(func(x ...any) string {
		val := myCondition.Expression().(string)
		decoded, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			myCondition.SetErr(err)
			//fmt.Printf("%v", err) // optional
			return `invalid_condition`
		}

		// optional
		//myCondition.SetExpression(decoded)

		// return new string value
		return fmt.Sprintf("%v %v %v",
			myCondition.Keyword(),
			myCondition.Operator(),
			string(decoded),
		)
	})

	fmt.Printf("%s", myCondition)
	// Output: keyword = 138499
}

/*
This example demonstrates the use of a custom PushPolicy instance
assigned to the receiver instance to limit valid push candidates
to only integers, and only those that are positive and are powers
of two (n²).

Generally speaking, even though the input signature for instances
of the PushPolicy type is variadic, it is recommended that it be
used in a unary manner, as we will demonstrate by calling slice #0
explicitly and exclusively.
*/
func ExampleStack_SetPushPolicy() {
	myStack := List()
	myStack.SetPushPolicy(func(x ...any) (err error) {
		if len(x) == 0 {
			return // no error because nothing was pushed
		}

		switch tv := x[0].(type) {
		case int:
			// We allow int values, but only certain
			// ones will be accepted for push.
			if tv <= 0 {
				// value cannot be negative nor can it be zero
				err = fmt.Errorf("%T->%T push denied: unsupported integer value: %d <= 0", tv, myStack, tv)
			} else if !isPowerOfTwo(tv) {
				// value must be power of two
				err = fmt.Errorf("%T->%T push denied: unsupported integer value: %d != n²", tv, myStack, tv)
			}
		default:
			// If not an integer, its most definitely bogus.
			err = fmt.Errorf("%T->%T push denied: unsupported type %T", tv, myStack, tv)
		}

		return
	})

	// prepare some values to be pushed
	// into the myStack instance.
	values := []any{
		3.14159,
		2,
		`2`,
		14,
		nil,
		-81,
		uint8(3),
		128,
		[]int{1, 2, 3, 4},
		1023,
		1,
		`peter`,
		nil,
	}

	// Push each slice individually and allow for the
	// opportunity to check for errors at each loop
	// iteration. Though this particular technique is
	// not required, it can be useful when dealing with
	// unvetted or unscrubbed data.
	var errct int
	for i := 0; i < len(values); i++ {
		myStack.Push(values[i])

		// For the sake of a simple example,
		// let's just count the errors. Push
		// what we can, don't quit.
		if err := myStack.Err(); err != nil {
			errct++
		}

		/*
			// Other error handling options for
			// the sake of a helpful example. One
			// or more of the actions below may
			// be appropriate, given the context
			// of the importing app ...
			if err := myStack.Err(); err != nil {
				//fmt.Println(err)
				//err = fmt.Errorf("Error with additional info prepended: %v", err)
				//panic(err)
				//return
				//continue or break
				//some other action
			}
		*/
	}

	fmt.Printf("Successful pushes: %d/%d", myStack.Len(), len(values))
	// Output: Successful pushes: 3/13
}

/*
This example demonstrates the analytical capabilities made
possible through use of the ValidityPolicy type.

Setting a scanning policy, as demonstrates below, can allow
a more in-depth analysis of the receiver instance to occur.

The Stack.Valid method shall abandon its "vanilla checks" of
the receiver and will, instead, report solely on the result
of calling the specified user-authored ValidityPolicy.
*/
func ExampleStack_SetValidityPolicy() {
	// a custom type made by the user,
	// because why not?
	type demoStruct struct {
		Type  string
		Value any
	}

	// Initialize our stack, and
	// push into it immediately.
	myStack := Basic().Push(
		`healthy stuff`,
		`healthy stuff`,
		`healthy stuff`,
		`healthy stuff`,
		demoStruct{
			Type:  `text`,
			Value: `complete_garbage`,
		},
	)

	// Set a ValidityPolicy. Users can handle
	// errors in any manner they wish, such as
	// using one the following actions:
	//
	//   - wrap multiple errors in one
	//   - encode error as JSON or some other encoding
	//   - error filtering (e.g.: ignore all except X, Y and Z)
	//   - craft preliminary error (e.g.: Found N errors; see logs for details)
	//
	// For the purpose of this simple example, we
	// are only scanning for demoStruct instances
	// and are only throwing an error should the
	// Value field not contain the literal string
	// of 'healthy stuff'.
	myStack.SetValidityPolicy(func(_ ...any) (err error) {

		// Iterate stack slices per length
		for i := 0; i < myStack.Len(); i++ {
			// Grab each stack index (slice)
			slice, _ := myStack.Index(i)

			// Perform type switch on slice
			switch tv := slice.(type) {
			case demoStruct:
				if tv.Type == `text` {
					if assert, _ := tv.Value.(string); assert != `healthy stuff` {
						// What a piece of junk!
						err = fmt.Errorf("Undesirable string value '%s' detected in struct Value field [%T @ idx:%d]", assert, myStack, i)
					}
				}
			}

			// We choose not to go any
			// further when errors appear.
			if err != nil {
				break
			}
		}

		return
	})

	// Execute Stack.Valid, which returns an error instance,
	// directly into a fmt.Printf call.
	fmt.Printf("%v", myStack.Valid())
	// Output: Undesirable string value 'complete_garbage' detected in struct Value field [stackage.Stack @ idx:4]
}

func ExampleCondition_SetValidityPolicy() {
	myCondition := Cond(`keyword`, Eq, int(3))
	myCondition.SetValidityPolicy(func(_ ...any) (err error) {
		assert, _ := myCondition.Expression().(int)
		if assert%2 != 0 {
			err = fmt.Errorf("%T is not even (%d)", assert, assert)
		}
		return
	})
	fmt.Printf("%v", myCondition.Valid())
	// Output: int is not even (3)
}

/*
This example demonstrates the use of the SetEvaluator method to
assign an instance of the Evaluator closure signature to the
receiver instance, allowing the creation and return of a value
based upon some procedure defined by the user.

Specifically, the code below demonstrates the conversion a degree
value to a radian value using the formula shown in the example
below:
*/
func ExampleCondition_SetEvaluator() {
	degrees := float64(83.1)
	c := Cond(`degrees`, Eq, degrees)

	// we don't really need input values for
	// this one, since we read directly from
	// the Condition instance c. And we don't
	// need an error value since this is just
	// an example, therefore we use shadowing
	// (_) as needed.
	c.SetEvaluator(func(_ ...any) (R any, _ error) {
		expr := c.Expression()
		D, _ := expr.(float64) // Don't shadow 'ok' in real-life.

		// I could have imported "math",
		// but this is all we need.
		pi := 3.14159265358979323846264338327950288419716939937510582097494459
		R = float64(D*pi) / 180
		return
	})

	radians, _ := c.Evaluate()
	fmt.Printf("%.02f° x π/180° = %.02frad", degrees, radians.(float64))
	// Output: 83.10° x π/180° = 1.45rad
}

func TestFCF_codecov(t *testing.T) {
	stk := Basic().Push(
		float32(1.1),
		float32(1.2),
		float32(1.3),
	)
	stk.SetPresentationPolicy(func(_ ...any) (_ string) {
		return `doomed string`
	})

	fmt.Printf("%s", stk)
	// Output:
}
