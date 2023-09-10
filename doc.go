/*
Package stackage implements a flexible stack type optimized for use in creating and presenting conditional Boolean statements, abstract mathematical constructs, LDAP filter abstractions or simple lists.

# Features

  - Flexible Stack configuration controls, allowing custom presentation, push and validity policies to be executed (instead of default behavior) through the use of closure signature functions
  - Recursive design - Stacks can reside in Stacks. Conditions can reside in Stacks. Conditions can contain other Stacks. Whatever.
  - Traversible - Recursive values are easily navigated using the Stack.Traverse method
  - Fluent-style - types which offer methods for cumulative configuration are written in fluent-form, allowing certain commands to be optionally "chained" together
  - Extensible logical operator framework, allowing custom operators to be added for specialized expressions instead of the package-provided ComparisonOperator constants
  - MuTeX capable for each Stack instance independent of its parent or child (no recursive locking mechanisms)
  - Adopters may wish to create a type alias of the Condition and/or Stack types; this is particularly easy, and will not impact normal operations when dealing with nested instances of derivative types
  - Assertion capabilities; go beyond merely crafting the string representation of a formula -- add an Evaluator function to conduct an interrogation of a value, matching procedures, and more!
  - Fast, reliable, useful

# Status

This package is in its early stages, and is undergoing active development. It should NOT be used in any production capacity at this time.

# License

The stackage (go-stackage) package, from http://github.com/JesseCoretta/go-stackage, is available under the terms of the MIT license. For further details, see the LICENSE file within the aforementioned repository.
*/
package stackage
