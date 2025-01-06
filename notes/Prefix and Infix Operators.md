# Prefix and Infix Operators

## Prefix Operators aka Unary Expressions

- Structure: `<prefix-operator><expression>`

- Any expression can follow a prefix opearator as operand. This means that an AST node for a prefix operator must be *flexible* enough to *point to any expressions as its opearand*

```js
// These are all valid
!isGreaterThanZero(2);
5 + -add(10 + 2);
```

- When parsing prefix expression like `"-5"` we need to handle *multiple tokens* unlike when we parse a single number. In detail, `parseExpression()` is invoked when we process the token **after** the prefix operator.

- We use the **recursive approach** extensively here and it is very powerful! Supposed we have to parse this expression `!-5` - note that the right side of the prefix operator `!` is ANOTHER prefix opeartor `-`

```js
parsePrefixExpression() sees "!"
└── creates PrefixExpression node for "!"
    └── calls parseExpression(PREFIX) for "-"
        └── calls parsePrefixExpression() (because "-" is a prefix operator)
            └── creates PrefixExpression node for "-"
                └── calls parseExpression(PREFIX) for "5"
                    └── calls parseIntegerLiteral() for "5"
```


- After calling `parsePrefixExpression()` as a `prefixParseFn`, we are back to the outer call to `parseExpression()` with `precedence = LOWEST` (it is never permanently changed to PREFIX) as it was the call made by `parseExpressionStatement()` deep in the call stack of the recursion.

```go
parseExpressionStatement() {
    // Initial call
    expression := parseExpression(LOWEST) {
        // See '-', call parsePrefixExpression
        leftExp := parsePrefixExpression() {
            // Parse the operand with PREFIX precedence
            operand := parseExpression(PREFIX) {
                // Parse '1' as integer literal
                return &ast.IntegerLiteral{Value: 1}
            }
            // Return the prefix expression (-1)
            return &ast.PrefixExpression{"-", operand}
        }

        // Now back in original parseExpression(LOWEST)
        // leftExp is (-1)
        // peek is '+', LOWEST < PLUS is true
        // Enter loop to parse rest of expression: + 2
        // Result: (-1) + 2
    }
}
```

The resulting AST for `!-5` would look like:

```js
PrefixExpression {
    Token: BANG
    Operator: "!"
    Right: PrefixExpression {
        Token: MINUS
        Operator: "-"
        Right: IntegerLiteral {
            Token: INT
            Value: 5
        }
    }
}
```

## Infix Operators aka Binary Expressions (Why the "Binary"? Because we have **two operands**)

- Structure of an infix operator expression: `<expression> <infix-operator> <expression>`.

- We construct a **precedence table** with our operators mapped with corresponding precedences.
