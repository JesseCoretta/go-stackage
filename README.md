# go-stackage

[![Go Report Card](https://goreportcard.com/badge/github.com/JesseCoretta/go-stackage)](https://goreportcard.com/report/github.com/JesseCoretta/go-stackage) [![Go Reference](https://pkg.go.dev/badge/github.com/JesseCoretta/go-stackage.svg)](https://pkg.go.dev/github.com/JesseCoretta/go-stackage) [![CodeQL](https://github.com/JesseCoretta/go-stackage/workflows/CodeQL/badge.svg)](https://github.com/JesseCoretta/go-stackage/actions/workflows/github-code-scanning/codeql) [![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-stackage/blob/main/LICENSE) [![codecov](https://codecov.io/gh/JesseCoretta/go-stackage/graph/badge.svg?token=RLW4DHLKQP)](https://codecov.io/gh/JesseCoretta/go-stackage) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-stackage/issues) [![Experimental](https://img.shields.io/badge/experimental-blue?logoColor=blue&label=%F0%9F%A7%AA%20%F0%9F%94%AC&labelColor=blue&color=gray)](https://github.com/JesseCoretta/JesseCoretta/blob/main/EXPERIMENTAL.md) [![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/jessecoretta/go-stackage/go.yml?event=push)](https://github.com/JesseCoretta/go-stackage/actions/workflows/go.yml) [![Author](https://img.shields.io/badge/author-Jesse_Coretta-darkred?label=%F0%9F%94%BA&labelColor=indigo&color=maroon)](mailto:jesse.coretta@icloud.com) [![GitHub release (with filter)](https://img.shields.io/github/v/release/JesseCoretta/go-stackage)](https://github.com/JesseCoretta/go-stackage/releases) [![Help Animals](https://img.shields.io/badge/help_animals-gray?label=%F0%9F%90%BE%20%F0%9F%98%BC%20%F0%9F%90%B6&labelColor=yellow)](https://github.com/JesseCoretta/JesseCoretta/blob/main/DONATIONS.md)

## Summary

The stackage package implements flexible Stack and Condition types with many useful features. It can be used to create object-based Boolean statements, abstract mathematical constructs, simple lists and much more: the possibilities are endless!

## Mission

The main goal of this package is provide an extremely reliable and accommodating stack/condition solution that is suitable for use in virtually any conceivable Go-based scenario in which objects of these types are needed. While extremely extensible and flexible, it should always be possible to use this package with no need for additional (non main) code while operating in extremely simple scenarios.

## Features

  - Stack instances are either LIFO (stack based, default) or FIFO (queue based)
    - FIFO is First-In/First-Out (like a line at your favorite deli: first come, first serve)
    - LIFO is Last-In/First-Out (like those plate-stacking apparatuses found in restaurant kitchens, in which the first plate inserted shall be the last plate removed)
  - Flexible Stack configuration controls, allowing custom stringer presentation, push controls and validity-checking policies to be imposed
  - Recursive design - Stacks can reside in Stacks. Conditions can reside in Stacks. Conditions can contain other Stacks. Whatever!
    - Eligible values are easily navigated using the Stack.Traverse method using an ordered sequence of indices, or slice index numbers
    - Conversely, recursion capabilities can also be easily disabled per instance!
  - Observable - flexible logging facilities, using the [log](pkg.go.dev/log) package are available globally, or on a per-instance basis
  - Interrogable - Stacks and Conditions extend many interrogation features, allowing the many facets and "states" of an instance to be queried simply
  - Resilient - Stack writeability can be toggled easily, allowing safe (albeit naïve) read-only operation without the need for mutexing
  - Fluent-style - Types which offer methods for cumulative configuration are written in fluent-form, allowing certain commands to be optionally "chained" together
  - Extensible logical operator framework, allowing custom operators to be added for specialized expressions instead of the package-provided ComparisonOperator constants
  - MuTeX capable for each Stack instance independent of its parent or child (no recursive locking mechanisms)
  - Adopters may wish to create a type alias of the Condition and/or Stack types; this is particularly easy, and will not impact normal operations when dealing with nested instances of derivative types
    - This is as simple as doing `type NewTypeName stackage.Stack` or `type NewTypeName stackage.Condition` in your code, and then extending/wrapping methods as needed
  - Extended Evaluator capabilities made possible through closures
    - Users can add their own Evaluator function to perform computational tasks, value interrogation, matching procedures, pretty much anything you could imagine
  - Fast, reliable, useful, albeit very niche

## Status

Although fairly well-tested, this package is in its early stages and is undergoing active development. It should only be used in production environments while under heavy scrutiny and with great care.

## License

The stackage package (from [`go-stackage`](https://github.com/JesseCoretta/go-stackage)) is released under the terms of the MIT license. See the [`LICENSE`](https://github.com/JesseCoretta/go-stackage/blob/main/LICENSE) file in the repository root, or click the "MIT" badge above, for complete details.

