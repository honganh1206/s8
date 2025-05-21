A toy programming language written in Go.

This language is based on the books "Writing An Interpreter In Go" and "Writing A Compiler In Go" by Thorsten Ball, but I intend to extend it further.

Why the name? The language is named after the football player Dominik Szoboszlai with jersey number 8 aka my fiance's favorite player :)

## Use

Just `go run ./src/main.go` for now and go with the flow from there

## Sample

Showcasing some features:

```js
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

... and many more!

## TODOs

Compilers

- [ ] Compile on RV64
- [ ] Compile to WebAssembly?

Operators

- [x] `?` as ternary operator
- [x] `~` as bitwise NOT operator
- [x] `^` as bitwise XOR operator
- [x] `|` as bitwise OR operator
- [x] `&` as bitwise AND operator
- [x] `++` for incrementing and `--` for decrementing
- [ ] `go`
- [ ] `select`
- [ ] `match`
- [ ] `.` to access fields

Object types

- [x] Float
- [ ] Double
- [ ] Lambda functions (a subset of anonymous functions)
- [ ] LazyObject
- [ ] Comments
- [ ] Struct
- [ ] Tuple
- [ ] Generics
- [ ] Channels
- [ ] Interface
- [ ] Range
- [ ] Procedure

Statements

- [ ] `switch`
- [ ] `foreach`
- [ ] `for`

Builtins & Libs

- [ ] `sleep`
- [ ] `map`
- [ ] `left` and `right` to return child nodes of an AST node
- [ ] `operator` to return the operator of an infix expression
- [ ] `arguments` to return an array of nodes in a `*ast.CallExpression`
- [ ] `children` function to return child nodes

Macros

- [ ] Make `quote` and `unquote` separate keywords
- [ ] Passing block statements to `quote`/`unquote`

Other nasty stuff

- [ ] Upgrade macro error handling system
- [x] Use `rune` instead of `byte` for chars
- [ ] Error handling by values

## Future plans

- [ ] Rewrite the interpreter + compiler in Zig

Refs:

## References

- [OK?](https://github.com/jesseduffield/OK)
- [straw](https://github.com/yjp20/turtle/tree/master/straw)
- [knox](https://github.com/AZHenley/knox)

---

The project isn't mature enough yet but the spirit of the effort is this: https://justforfunnoreally.dev/
