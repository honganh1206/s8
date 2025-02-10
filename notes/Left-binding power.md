# Left-binding power

One way to make our operators more right-associative is to _decrement its precendence_

Our call to `peekPrecedence()` method returns stands for the left-binding power of the next operator `p.peekToken`. Then the loop inside `parseExpression()` with the condition `precedence < p.peekPrecedence()` checks _if the left-binding power of the next operator/token is higher than our current right-binding power_

**TLDR**: It is like a **tug-of-war** whether the next operator (peek token) has enough "power" to steal operands.

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

_Every token has the same right- and left-binding power_, and the precedence of each token type inside the `precedences` table can be either used as right- or left-binding power depending on the context.

1.  When parsing `1 + 2 + 3`:

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
