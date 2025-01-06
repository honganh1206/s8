# Right-binding Power

## Definition

- The higher the power is, the more tokens/operators/operands to the right of the current expressions (the future peek tokens) can we "bind" to it.

- An example would be this expression `a = b = c = 5` and this will be parsed as `a = (b = (c = 5))` because the `=` can "bind" (aka wrap their operands in *implicit* parentheses) all the expressions to its right. The rightmost `=` will bind first and we go from right to left.

- Binding power can affect whether we continue parsing more tokens into our current expression or not.

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
    leftExp = !5  // ! has high precedence
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

## Counterpart: Left-binding Power

- Our call to `peekPrecedence()` method returns stands for the left-binding power of the next operator `p.peekToken`. Then the loop inside `parseExpression()` with the condition `precedence < p.peekPrecedence()` checks *if the left-binding power of the next operator/token is higher than our current right-binding power*

- **TLDR**: It is like a **tug-of-war** whether the next operator (peek token) has enough "power" to steal operands.

Some examples:

Ex1: `1 + 2 * 3`

```javascript
// Precedences: LOWEST = 0, PLUS = 4, PRODUCT = 5

parseExpression(LOWEST) {
    leftExp = 1
    // peek is +
    // LOWEST(0) < PLUS(4) -> true, enter loop
    // parse: 1 + ...

    parseExpression(PLUS) {
        leftExp = 2
        // peek is *
        // PLUS(4) < PRODUCT(5) -> true, enter loop
        // So 2 becomes left side of * aka the 2 gets "sucked" to the right
        // Results in: 1 + (2 * 3)
    }
}
```

Ex2: `5 * 2 + 3`

```javascript
// Precedences: LOWEST = 0, PLUS = 4, PRODUCT = 5

parseExpression(LOWEST) {
    leftExp = 5
    // peek is *
    // LOWEST(0) < PRODUCT(5) -> true, enter loop
    // parse: (5 * 2)

    // peek is +
    // PRODUCT(5) ≮ PLUS(4) -> false, don't enter loop
    // So (5 * 2) becomes left side of + aka the 2 gets "sucked" to the left
    // Results in: ((5 * 2) + 3)
}
```

Ex3: `-1 + 2`

```javascript
// Precedences: LOWEST = 0, PLUS = 4, PREFIX = 6

parseExpression(LOWEST) {
    leftExp = -1
    // peek is +
    // PREFIX(6) > PLUS(4) -> false, don't enter loop
    // No infixParseFn is going to get the number 1
    // Thus 1 is returned as the "right" arm of our prefix expression
    // parse: (-1) + 2
}
```

## Notable ideas

- *Every token has the same right- and left-binding power*, and the precedence of each token type inside the `precedences` table can be either used as right- or left-binding power depending on the context.

 1. When parsing `1 + 2 + 3`:
```go
parseExpression(LOWEST) {
    leftExp = 1
    // peek is +
    // Here PLUS acts as left-binding power (can we pull leftExp?)
    // LOWEST < PLUS -> true, so + can take 1 as left operand

    parseExpression(PLUS) {
        leftExp = 2
        // peek is second +
        // Here PLUS acts as right-binding power (can we be pulled right?)
        // PLUS ≮ PLUS -> false, so 2 stays with first +
    }
}
```

2. When parsing `5 * 2 + 3`:
```go
parseExpression(LOWEST) {
    leftExp = 5
    // peek is *
    // Here PRODUCT acts as left-binding power
    // LOWEST < PRODUCT -> true, so * can take 5 as left operand

    parseExpression(PRODUCT) {
        leftExp = 2
        // peek is +
        // Here PRODUCT acts as right-binding power
        // PLUS < PRODUCT -> false, so (5 * 2) stays as unit
    }
}
```

- One way to make our operators more right-associative is to *decrement its precendence*
