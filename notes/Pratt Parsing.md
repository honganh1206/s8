# Pratt Parsing/Top-down Operator Precedence

## Definition

- The "Top Down Operator Precedence" by Vaughan Pratt is recently rediscovered and later popularized.
- This was invented as an *alternative* to parsers based on context-free grammars and the [[Backus-Naur Form]]

## The main difference

- Instead of associating parsing functions with grammar rules (like we did with `parseLetStatement()`, Pratt associated these functions with **single token type**.

- *Each token type can have TWO associated parsing functions*: A prefix function (when the token appears at the **start** of an expression like `-5`) and an infix function (when the token appears **between two expressions** like `10 + 2`)

## Implementation

- Core idea: **Associating** parsing functions with different token types. When we encounter a specific token type, we *invoke the specific parsing functions* to get the AST node that represent the token.

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

### How this works

```go
// Parsing "5 + 10 * 2"
 func (p *Parser) parseExpression(precedence int) ast.Expression {
    // First iteration: Parse "5"
    prefix := p.prefixParseFns[p.currentToken.Type]
    leftExp := prefix() // leftExo = 5

		// Loop starts and we see "+"
		// As there is no operator yet before "+", we compare "+" with LOWEST
    for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
        // Get the infix function for "+" or "*"
        // The "*" has a higher precedence so we group the "10 * 2" first
        infix := p.infixParseFns[p.peekToken.Type]
        if infix == nil {
            return leftExp
        }
        p.nextToken()
        // The key part:
        // This will RECURSIVELY parse "10 * 2" first
        leftExp = infix(leftExp)
    }

    return leftExp
}

// When parsing "+"
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression{
        Token:    p.currentToken,  // "+"
        Operator: p.currentToken.Literal,
        Left:     left,           // "5"
    }

    precedence := p.curPrecedence()  // precedence of "+"
    p.nextToken()  // move to "10"

    // This recursive call is crucial!
    // It will parse "10 * 2" as one unit because "*" has higher precedence
    expression.Right = p.parseExpression(precedence)

    return expression
}
```

The process:

1. Parse "5" as leftExp
2. See "+", enter loop
3. In parseInfixExpression for "+":
   - Set "5" as Left
   - **Recursively** call parseExpression for right side with "+" precedence
4. In recursive call for "10 * 2":
   - Parse "10" as leftExp
   - See "*", which has higher precedence than "+"
   - Create new infix expression with "10" as Left
   - Parse "2" as Right
   - Return "(10 * 2)" as one unit
5. Back in original "+", set "(10 * 2)" as Right

The result is an AST that looks like:
```
    +
   / \
  5   *
     / \
    10  2
```
