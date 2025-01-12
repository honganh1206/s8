# What our parser will do

What will our parser do? It **repeatedly** advances the tokens + checks the current token to decide whether to _call another parsing function_ or _throw an error_.

The flow of the `parseProgram()` function the key function of our parser:

```
 parseProgram()
|
|---> program = newProgramASTNode()
|
|---> advanceTokens()
|
|---> Loop: currentToken() != EOF_TOKEN
|       |
|       +--> IF currentToken() == LET_TOKEN
|       |       |
|       |       +--> statement = parseLetStatement()
|       |
|       +--> ELSE IF currentToken() == RETURN_TOKEN
|       |       |
|       |       +--> statement = parseReturnStatement()
|       |
|       +--> ELSE IF currentToken() == IF_TOKEN
|               |
|               +--> statement = parseIfStatement()
|
|---> IF statement != null
|       |
|       +--> program.Statements.push(statement)
|
|---> advanceTokens()
|
+---> return program
```

The basic idea behind `parseProgram()`: We first build the root node of the AST -> Build the child nodes + statements by _calling functions that know which AST node to construct_ based on the current token.

When creating a parser, we need to both set the current token + next token for cases like `5;`, in which the current token is `token.INT` and the next token indicates whether we are dealing with EOF or starting an arithmetic expression.

## Variable binding in s8

An expression produces a value, a statement does not. However, what exactly an expression or a statement is or which one produces the value _depends on the language_. A function literal could be an expression in some languages while in others it is not (it will be in s8 though!)

```js
// An example of a function literal
const add = function (a, b) {
  return a + b;
};
```

Using the `Identifier` structfor both variable names and as a general expression is a good practice: Sometimes we use something like `myVar + 5` as a standalone expression/variable reference, and if so we have to do this:

```go
// For variable declarations
type VariableName struct {
   Token token.Token
   Value string
}

// For variable references
type VariableReference struct {
   Token token.Token
   Value string
}
```

[[Recursive-descent Parsing]]

## Recursive `parseExpression()`

- Note that our `parseExpression()` function and those alike are quite **recursive**

```
 parseExpression()
|
|---> IF currentToken() == INTEGER_TOKEN
|       |
|       +--> IF nextToken() == PLUS_TOKEN
|       |       |
|       |       +--> return parseOperatorExpression()
|       |
|       +--> ELSE IF nextToken() == SEMICOLON_TOKEN
|               |
|               +--> return parseIntegerLiteral()
|
|---> ELSE IF currentToken() == LEFT_PAREN
|       |
|       +--> return parseGroupedExpression()
|
|---> [Additional conditions if any, represented as ...]
|
+---> return

---


parseOperatorExpression()
|
|---> operatorExpression = newOperatorExpression()
|       |
|       +--> operatorExpression.left = parseIntegerLiteral()
|       |
|       +--> advanceTokens()
|       |       |
|       |       +--> operatorExpression.operator = currentToken()
|       |
|       +--> advanceTokens()
|               |
|               +--> operatorExpression.right = parseExpression() // Recursion
|
+---> return operatorExpression
```

- When parsing an expression like `5 + 5`, we first parse `5 +` and then call the `parseExpression()` recursively because _there might be another operator expression_ like `5 + 5 * 10`.
