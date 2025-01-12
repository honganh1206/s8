# Call Expressions

Structure of the expression: `<expression>(<comma separated expressions>)`

Example: `add(2,3)`. Note that _`add` is an identifier and identifiers are expressions too._

The function `add()` is **bounded** to the identifier `add`, thus the identifier `add` returns the function `add()` when it is evaluated.

For this reason, we can just _go straight to the source -> replace the identifier with the function literal_

```js
funk(x, y) { x + y; }(2, 3) // This is valid
callsFunction(2, 3, fn(x, y) { x + y;}) // This is also valid

```

At this point, we might encounter any errors telling that we have yet to register `prefixParseFn` for call expressions because _there are no new token types in call expressions_. Thus, we only need to register prefix parsing functions for the parentheses

The left parentheses comes between the identifier `add` and the list of arguments, so we need to _register an `infixParseFn` for `tolen.LPAREN`_.

Flow: Parse the expression that is the function (either an identifier or a function literal `funk`) -> Check for an `infixParseFn` associated with `token.LPAREN` -> Call that parsing function with the already parsed expression as argument
