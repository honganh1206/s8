A toy programming language written in Go.

This is based on the books "Writing An Interpreter In Go" and "Writing A Compiler In Go" by Thorsten Ball, but I intend to extend it further.

## TODOs

- [x] Use `rune` instead of `byte` for chars
- Token types from [bantam-go](https://github.com/obzva/bantam-go)
- [x] `?` as ternary/conditional operator
- [x] `~` as bitwise NOT operator
- [x] `^` as bitwise XOR operator
- [x] `|` as bitwise OR operator
- [x] `&` as bitwise AND operator
- [x] `++` for incrementing and `--` for decrementing (Both postfix/infix and prefix)
- Object types
- [x] Float
- [ ] Double
- [ ] Lambda functions (a subset of anonymous functions)
- [ ] LazyObject
- [ ] Comments
- [ ] Struct
- [ ] Tuple
- [ ] Generics
- Statements
- [ ] Iterators
- [ ] `switch` statement
- [ ] `foreach` statement
- [ ] `for` statement with index
- Builtins & Libs
- [ ] `sleep`
- [ ] `map`
- [ ] `left` and `right` to return child nodes of an AST node
- [ ] `operator` to return the operator of an infix expression
- [ ] `arguments` to return an array of nodes in a `*ast.CallExpression`
- [ ] `children` function to return child nodes
- Reflection
- Macros
- [ ] Make `quote` and `unquote` separate keywords
- [ ] Passing block statements to `quote`/`unquote`
- Other nasty stuff
- [ ] Upgrade macro error handling system

## Future plans

- [ ] Rewrite the interpreter + Compiler in Zig

Refs:

- [ ] Compile to WebAssembly

## References

- [OK?](https://github.com/jesseduffield/OK)
- [lox-zig](https://github.com/adrianchong518/lox-zig) (25 commits)
- [zigself](https://github.com/sin-ack/zigself) (>600 commits)
- [kiesel - JS engine](https://codeberg.org/kiesel-js/kiesel) (>2000 commits)

---

The project isn't mature enough yet but the spirit of the effort is this: https://justforfunnoreally.dev/
