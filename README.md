# go-stackage

[![Go Report Card](https://goreportcard.com/badge/github.com/JesseCoretta/go-stackage)](https://goreportcard.com/report/github.com/JesseCoretta/go-stackage) [![Go Reference](https://pkg.go.dev/badge/github.com/JesseCoretta/go-stackage.svg)](https://pkg.go.dev/github.com/JesseCoretta/go-stackage) [![CodeQL](https://github.com/JesseCoretta/go-stackage/workflows/CodeQL/badge.svg)](https://github.com/JesseCoretta/go-stackage/actions/workflows/github-code-scanning/codeql) [![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-stackage/blob/main/LICENSE) [![codecov](https://codecov.io/gh/JesseCoretta/go-stackage/graph/badge.svg?token=RLW4DHLKQP)](https://codecov.io/gh/JesseCoretta/go-stackage) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-stackage/issues) [![Experimental](https://img.shields.io/badge/experimental-blue?logoColor=blue&label=%F0%9F%A7%AA%20%F0%9F%94%AC&labelColor=blue&color=gray)](https://github.com/JesseCoretta/JesseCoretta/blob/main/EXPERIMENTAL.md) [![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/jessecoretta/go-stackage/go.yml?event=push)](https://github.com/JesseCoretta/go-stackage/actions/workflows/go.yml) [![Author](https://img.shields.io/badge/author-Jesse_Coretta-darkred?label=%F0%9F%94%BA&labelColor=indigo&color=maroon)](mailto:jesse.coretta@icloud.com) [![GitHub release (with filter)](https://img.shields.io/github/v/release/JesseCoretta/go-stackage)](https://github.com/JesseCoretta/go-stackage/releases) [![Help Animals](https://img.shields.io/badge/help_animals-gray?label=%F0%9F%90%BE%20%F0%9F%98%BC%20%F0%9F%90%B6&labelColor=yellow)](https://github.com/JesseCoretta/JesseCoretta/blob/main/DONATIONS.md)

## Summary

The stackage package implements flexible Stack and Condition types with many useful features. It can be used to create object-based Boolean statements, abstract mathematical constructs, simple lists and much more: the possibilities are endless!

## Features

  - Flexible Stack configuration controls, allowing custom presentation, push and validity policies to be executed (instead of default behavior) through the use of closure signature functions
  - Recursive design - Stacks can reside in Stacks. Conditions can reside in Stacks. Conditions can contain other Stacks. Whatever!
    - Conversely, recursion capabilities can also be easily disabled per instance!
  - Interrogable - Stacks and Conditions extend many interrogation features, allowing the state of an instance to be easily qualified or scrutinized
  - Resilient - Stack writeability can be toggled easily, allowing safe read-only operation without the need for mutexing
  - Traversible - Recursive values are easily navigated using the Stack.Traverse method
  - Fluent-style - Types which offer methods for cumulative configuration are written in fluent-form, allowing certain commands to be optionally "chained" together
  - Extensible logical operator framework, allowing custom operators to be added for specialized expressions instead of the package-provided ComparisonOperator constants
  - MuTeX capable for each Stack instance independent of its parent or child (no recursive locking mechanisms)
  - Adopters may wish to create a type alias of the Condition and/or Stack types; this is particularly easy, and will not impact normal operations when dealing with nested instances of derivative types
  - Assertion capabilities; go beyond merely crafting the string representation of a formula -- add an Evaluator function to conduct an interrogation of a value, matching procedures, and more!
  - Fast, reliable, useful, albeit very niche

## Status

Although fairly well-tested, This package is in its early stages and is undergoing active development. It should only be used in production environments while under heavy scrutiny and with great care.

## License

The stackage package (from [`go-stackage`](https://github.com/JesseCoretta/go-stackage)) is released under the terms of the MIT license. See the [`LICENSE`](https://github.com/JesseCoretta/go-stackage/blob/main/LICENSE) file in the repository root, or click the "MIT" badge above, for complete details.

