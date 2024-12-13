# Recursive-descent Parsing

- The flow of the `parseProgram()` function - the key function of our parser:

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

-
