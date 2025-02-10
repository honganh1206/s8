# Right-binding Power

## Definition

- The higher the power is, the more tokens/operators/operands to the right of the current expressions (the future peek tokens) can we "bind" to it.

- An example would be this expression `a = b = c = 5` and this will be parsed as `a = (b = (c = 5))` because the `=` can "bind" (aka wrap their operands in _implicit_ parentheses) all the expressions to its right. The rightmost `=` will bind first and we go from right to left.

- Binding power can affect whether we continue parsing more tokens into our current expression or not.

## Examples

Consider the two expression statements: `1 + 2 + 3` and `!5` and `!5 + 3`

Case 1 - Regular Left-Associative/Infix Expression: `1 + 2 + 3`

```javascript
// When parsing "1 + 2 + 3"
// The precedence of + is PLUS (4)

parseExpression(LOWEST) {
    leftExp = 1
    // peekToken is +, precedence(+) > LOWEST, so we enter loop
    leftExp = (1 + 2)  // First +
    // peekToken is +, precedence(+) > LOWEST, so we enter loop again
    leftExp = ((1 + 2) + 3)  // Second +
}
```

Case 2 - Prefix Expression: `!5`

```javascript
// When parsing "!5"
// The precedence of ! is PREFIX (6)

parseExpression(LOWEST) {
    leftExp = !5  // ! has high precedence PREFIX
    // Even if there's a next operator like +
    // We won't enter the loop because PREFIX > most other precedences
    // So !5 stays as is and doesn't become part of a larger expression
}
```

Case 3 - Combining Prefix and Infix Expressions: `!5 + 3`

The parsing process:

1. Parse `!5` as a prefix expression
2. Peek to the operator `+`
3. As `!` > `+` in terms of precedence, the `!5` does not enter the loop and thus stays as a unit like `(!5)` and not a left child node
4. The final AST becomes `(!5) + 3` and not `!(5 + 3)`

[[Left-binding power]]
