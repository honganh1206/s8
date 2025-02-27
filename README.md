A toy programming language written in Go.

This is based on the books "Writing An Interpreter In Go" and "Writing A Compiler In Go" by Thorsten Ball, but I intend to extend it further.

## TODOs

- [x] Use `rune` instead of `byte` for chars
-  More token types from [bantam-go](https://github.com/obzva/bantam-go)
  - [x] `?` as ternary/conditional operator
  - [ ] `^` as exponent (It is called the caret)
  - [ ] `~` as bitwise NOT operator
  - [ ] `++` for incrementing and `--` for decrementing (Both postfix)
-  More object types
  - [ ] Float
- [ ] A custom garbage collector in C
- [ ] Rewrite the language in Rust

---

The project isn't mature enough yet but the spirit of the effort is this: https://justforfunnoreally.dev/
