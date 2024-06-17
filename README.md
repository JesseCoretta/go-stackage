# go-stackage

[![Go Report Card](https://goreportcard.com/badge/github.com/JesseCoretta/go-stackage)](https://goreportcard.com/report/github.com/JesseCoretta/go-stackage) [![Go Reference](https://pkg.go.dev/badge/github.com/JesseCoretta/go-stackage.svg)](https://pkg.go.dev/github.com/JesseCoretta/go-stackage) [![CodeQL](https://github.com/JesseCoretta/go-stackage/workflows/CodeQL/badge.svg)](https://github.com/JesseCoretta/go-stackage/actions/workflows/github-code-scanning/codeql) [![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-stackage/blob/main/LICENSE) [![codecov](https://codecov.io/gh/JesseCoretta/go-stackage/graph/badge.svg?token=RLW4DHLKQP)](https://codecov.io/gh/JesseCoretta/go-stackage) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-stackage/issues) [![Experimental](https://img.shields.io/badge/experimental-blue?logoColor=blue&label=%F0%9F%A7%AA%20%F0%9F%94%AC&labelColor=blue&color=gray)](https://github.com/JesseCoretta/JesseCoretta/blob/main/EXPERIMENTAL.md) [![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/jessecoretta/go-stackage/go.yml?event=push)](https://github.com/JesseCoretta/go-stackage/actions/workflows/go.yml) [![Author](https://img.shields.io/badge/author-Jesse_Coretta-darkred?label=%F0%9F%94%BA&labelColor=indigo&color=maroon)](mailto:jesse.coretta@icloud.com) [![GitHub release (with filter)](https://img.shields.io/github/v/release/JesseCoretta/go-stackage)](https://github.com/JesseCoretta/go-stackage/releases) [![Help Animals](https://img.shields.io/badge/help_animals-gray?label=%F0%9F%90%BE%20%F0%9F%98%BC%20%F0%9F%90%B6&labelColor=yellow)](https://github.com/JesseCoretta/JesseCoretta/blob/main/DONATIONS.md)

![stacks03](https://github.com/JesseCoretta/go-stackage/assets/74126726/ec6429ac-d698-413b-8c1c-0ac8dc1e301c)

## Summary

The stackage package implements flexible Stack and Condition types with many useful features. It can be used to create object-based Boolean statements, abstract mathematical constructs, simple lists and much more: the possibilities are endless!

## Mission

The main goal of this package is provide an extremely reliable and accommodating stack/condition solution that is suitable for use in virtually any conceivable Go-based scenario in which objects of these types are needed. While extremely extensible and flexible, it should always be possible to use this package with no need for additional (non main) code while operating in extremely simple scenarios.

## Features

  - Stack instances are either LIFO (stack based, default) or FIFO (queue based)
    - FIFO is First-In/First-Out (like a line at your favorite deli: first come, first serve)
    - LIFO is Last-In/First-Out (like those plate-stacking apparatuses found in restaurant kitchens, in which the first plate inserted shall be the last plate removed)
  - Flexible Stack configuration controls, allowing custom stringer presentation, sorting closures, push controls and validity-checking policies to be imposed
  - Recursive design - Stacks can reside in Stacks. Conditions can reside in Stacks. Conditions can contain other Stacks. Whatever!
    - Eligible values are easily navigated using the Stack.Traverse method using an ordered sequence of indices, or slice index numbers
    - Conversely, recursion capabilities can also be easily disabled per instance!
  - Observable - flexible logging facilities, using the [log](https://pkg.go.dev/log) package are available globally, or on a per-instance basis
  - Interrogable - Stacks and Conditions extend many interrogation features, allowing the many facets and "states" of an instance to be queried simply
  - Resilient - Stack writeability can be toggled easily, allowing safe (albeit naïve) read-only operation without the need for mutexing
  - Fluent-style - Types which offer methods for cumulative configuration are written in fluent-form, allowing certain commands to be optionally "chained" together
  - Extensible
    - Logical operator framework allows custom operators to be added for specialized expressions, instead of the package-provided ComparisonOperator constants
    - Users can add their own Evaluator function to perform computational tasks, value interrogation, matching procedures ... pretty much anything you could imagine
  - Stack instances are (independently) MuTeX capable, thanks to the [sync](https://pkg.go.dev/sync) package
    - Recursive locking mechanisms are NOT supported due to my aversion to insanity
  - Adopters may create a type alias of the Condition and/or Stack types
    - See the [Type Aliasing](#type-aliasing) section below
  - Fast, reliable, useful, albeit very niche

## Status

Although fairly well-tested, this package is in its early stages and is undergoing active development. It should only be used in production environments while under heavy scrutiny and with great care.

## License

The stackage package, from [`go-stackage`](https://github.com/JesseCoretta/go-stackage), is released under the terms of the MIT license. See the [`LICENSE`](https://github.com/JesseCoretta/go-stackage/blob/main/LICENSE) file in the repository root, or click the License badge above, for complete details.

## Type Aliasing

When needed, users may opt to create their own derivative alias types of either the Stack or Condition types for more customized use in their application.

The caveat, naturally, is that users will be expected to wrap all of the package-provided methods (e.g.: `String`, `Push`, `Pop`, etc) they intend to use.

However, the upside is that the user may now write (extend) wholly _new_ methods that are unique to their own application, and _without_ having to resort to potentially awkward measures, such as embedding.

To create a derivative type based on the Stack type, simply do something similar to the following example in your code:

```
type MyStack stackage.Stack

// Here we extend a wholly new function. The input and output signatures
// are entirely defined at the discretion of the author and are shown in
// "pseudo code context" here.
func (r MyStack) NewMethodName([input signature]) [<output signature>] {
	// your custom code, do whatever!
}

// Here we wrap a pre-existing package-provided function, String, that
// one would probably intend to use.
// 
// To run the actual String method, we need to first CAST the custom
// type (r, MyStack) to a bonafide stackage.Stack instance as shown
// here. Unlike the above example, this is NOT "pseudo code" and will
// compile just fine.
//
// Repeat as needed for other methods that may be used.
func (r MyStack) String() string {
 	// return the result from a "TYPE CAST -> EXEC" call
	return stackage.Stack(r).String()
}

// For added convenience, adopters can write their own private "cast"
// method for quick transformation back to the derived Stack type.
// This allows easy access to base methods which the adopter has not
// explicitly wrapped.
func (r MyStack) cast() stackage.Stack {
	return stackage.Stack(r)
}
```

The procedure would be identical for a Condition alias -- just change the name and the derived stackage type from the first example line and modify as desired.

If you'd like to see a more complex working example of this concept in the wild, have a look at the [`go-aci`](https://github.com/JesseCoretta/go-aci) package, which makes **heavy use** of derivative stackage types.

