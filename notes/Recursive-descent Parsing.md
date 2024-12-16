# Recursive-descent Parsing

## Definition

- A top-down parsing technique to analyze and process the *structure of a language based on its grammar rules*.
- This technique involes *a set of mutually recursive functions* where each function is responsible for a grammar part of the language being parsed.

## How it works

- Each non-terminal (Expr/Term/Factor) in the grammar is represented by a function.
- These function *recursively* calls each other to match the input against the grammar rule.
- The parser consumes tokens from the input stream (lexer) and match them to grammar constructs.
