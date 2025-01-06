# Interpreter Index

## Introduction

- A common trait between interpretersL They the source code and evaluate it  without producing visible, intermediate result for execution. This goes against compilers, which take the source code and produce output for machines to understand.
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

```
// Showcasing some features

// Variable binding
let age = 1;
let name = "Szoboszlai";
let result = 10 / 5;

// Data structures
let arr = [1, 2, 3, 4, 5];
arr[0]
let me = {"name": "Hong Anh", "age": 28"};
me["name"]

// Bind functions to names with implicit return
let explicitAdd = funk(a, b) { return a + b;}
let implicitADd = funk(a, b) { a + b};
let fib = funk(x) {
  if (x == 0) {
    0; // Implicit return
    } else {
      if (x == 1) {
       1;
      } else {
        fib(x - 1) + fib(x - 2);
      }
};

// Higher-order functions (functions that take other functions as arguments)
let twice = funk(f, x) {
  return f(f(x)); // Call the function passed as an argument two times
};

let addTwo = funk(x) {
  return x + 2;
};

twice(addTwo; 2); // Return the value of the first call

 ```
- We are going tokenize and parse the s8 source code in a Read-Eval-Print Loop (REPL) as an interactive programming environment

## Major components

- [x] Lexer

- [ ] Parser

- [ ] AST

- [ ] Internal object system

- [ ] Evaluator

## Index

- [[Lexing]]
- [[Parsing]]
- [[Evaluation]]
- [[Extending the Interpreter]]
- [[Going Further]]
