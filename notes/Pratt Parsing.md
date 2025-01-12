# Pratt Parsing/Top-down Operator Precedence

## Definition

- The "Top Down Operator Precedence" by Vaughan Pratt is recently rediscovered and later popularized.
- This was invented as an _alternative_ to parsers based on context-free grammars and the [[Backus-Naur Form]]

## The main difference

- Instead of associating parsing functions with grammar rules (like we did with `parseLetStatement()`, Pratt associated these functions with **single token type**.

- _Each token type can have TWO associated parsing functions_: A prefix function (when the token appears at the **start** of an expression like `-5`) and an infix function (when the token appears **between two expressions** like `10 + 2`)

## Implementation

- Core idea: **Associating** parsing functions with different token types. When we encounter a specific token type, we _invoke the specific parsing functions_ to get the AST node that represent the token.

```go
// For parsing expressions that start with a token
type prefixParseFn func() ast.Expression
// Examples:
// - Identifiers: "x"
// - Numbers: "5"
// - Prefix operators: "-x", "!true"
// - Parentheses: "(x + y)"

// For parsing operators between expressions
type infixParseFn func(ast.Expression) ast.Expression
// Examples:
// - Binary operators: "5 + 10", "x * y"
// - Function calls: "add(x, y)"
// - Array index: "array[0]"
```

### Examples

- A simple case: `5 + 10`

The process:

1. Parse `5` with prefix parsing function `parseIntegerLiteral()`
2. Enter the loop as we satisfy the conditions (not semicolon + LOWEST < SUM)
3. Call `parseInfixExpression()` and advance to `+` with the current left expression as `5`
4. (Inside `parseInfixExpression()`) Get the SUM precedence -> Advance to `10` -> Handle `10` with `parseIntegerLiteral()` since it does not satisfy the loop conditions

Call stack:

```js
parseExpressionStatement()
└── parseExpression(LOWEST)
   ├── parseIntegerLiteral() -> returns 5
   ├── enters infix loop (sees +)
   └── parseInfixExpression(left: 5)
       ├── creates InfixExpression(+)
       └── parseExpression(SUM)
           └── parseIntegerLiteral() -> returns 10
           // Doesn't enter infix loop (sees semicolon)
```

The final AST would look like:

```
    +
   / \
  5   10
```

- A more complex case: `5 + 10 * 2`

Call stack:

```js
parseExpressionStatement()
└── parseExpression(LOWEST)
    ├── parseIntegerLiteral() -> returns 5
    ├── enters infix loop (sees +)
    └── parseInfixExpression(left: 5)
        ├── creates InfixExpression(+)
        └── parseExpression(SUM)
            ├── parseIntegerLiteral() -> returns 10
            ├── enters infix loop (sees *)
            └── parseInfixExpression(left: 10)
                ├── creates InfixExpression(*)
                └── parseExpression(PRODUCT)
                    └── parseIntegerLiteral() -> returns 2
```

The result is an AST that looks like:

```
    +
   / \
  5   *
     / \
    10  2
```
