# Recursive-descent Parsing

## Definition

- A top-down parsing technique to analyze and process the _structure of a language based on its grammar rules_.
- This technique involes _a set of mutually recursive functions_ where each function is responsible for a grammar part of the language being parsed.

## How it works

- Each non-terminal (Expr/Term/Factor) in the grammar is represented by a function.
- These function _recursively_ calls each other to match the input against the grammar rule.
- The parser consumes tokens from the input stream (lexer) and match them to grammar constructs.

## Important ideas

1. Our parsing functions should **NEVER** advance the tokens too far. Supposed we are parsing a prefix expression `-5` with this method:

```go
// Should only handle `-5` and no more
func (p *Parser) parsePrefixExpression() ast.Expression {
    // Starts with currentToken being "-"
    expression := &ast.PrefixExpression{
        Token:    p.currentToken, // The "-" token
        Operator: p.currentToken.Literal,
    }

    // Advance to the next token to parse the right side
    p.nextToken()

    // Parse the right side (the "5")
    expression.Right = p.parseExpression(PREFIX)

    // Ends with currentToken being "5"
    return expression
}
```

1. When `parsePrefixExpression()` starts, `currentToken` is the token the function is parsing
2. The function should **only** consume tokens that are _part of its expression_
3. When `parsePrefixExpression()` is done executing, `currentToken` should be the **last** token of its expression

- What are the benefits?

1. Parsing functions become composable - each function knows exactly where it starts and stops.
2. We prevent tokens from accidentally skipped or processed twice.
3. We ensure each parsing function can parse only tokens that it is responsible for.

4. Rercursive-descent parsing is a **general** technique, while [[Pratt Parsing]] specifically uses recursive-descent parsing but adds more concepts like:

- Precedence (binding power)
- Prefix/infix parsing function types
- A core `parseExpression(precedence)` function that **recursively** calls itself
