# Interpreter Index

## Introduction

- A common trait between interpreters: They the source code and evaluate it without producing visible, intermediate result for execution. This goes against compilers, which take the source code and produce output for machines to understand.
- More advanced are the interpreters that compile the source code into an internal representation of bytecode for later evaluation e.g., CPython, Lua, JVM, etc., and even more advanced are Just-In-Time (JIT) interpretersi that converting bytecode into native machine code during runtime!
- Some interpreters just parse the source code -> build an abstract syntax tree (AST) -> evaluate this tree like walking through the AST and interpreting it.

## Features

- C-like syntax
- Variable bindings
- Intergers and booleans
- Arithmetic expressions
- Built-in functions
- First-class and higher-order functions
- Closures
- A string data structure
- An array data structure
- A hash data structure

We are going tokenize and parse the s8 source code in a Read-Eval-Print Loop (REPL) as an interactive programming environment

## Index

- [Lexing](Lexing.md)
- [Parsing](Parsing.md)
- [Evaluation](Evaluation.md)
- [Extending the interpreter](Extending%20the%20interpreter.md)
