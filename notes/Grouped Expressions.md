# Grouped Expressions

- How our method `parseGroupedExpression()` works: It handles grouped expressions that are _enclosed in parentheses_ and _maintains proper precedence_ by treating the grouped expressions as **a single unit** e.g., In `5 * (2 + 3)` the `(2 + 3)` should be evaluated first.
- Function flow: When encountering '(', we parse everything until reaching ')'  and return the expression as a single unit
```go
// When this function is called, currentToken is at '('
p.nextToken()  // Move to the token after '('
expr := p.parseExpression(LOWEST)  // Parse everything until ')'
if !p.expectPeek(token.RPAREN) {  // Ensure the expression ends with ')'
    return nil // In case of a syntax error
}
return expr
```

