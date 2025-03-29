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
- Builtins & Libs
- [ ] `sleep`
- [ ] `map`
- [ ] Reflection

## Future plans

- [ ] A custom garbage collector in C
- [ ] Rewrite the language in Rust

## Repos for reference

- [OK?](https://github.com/jesseduffield/OK)

---

The project isn't mature enough yet but the spirit of the effort is this: https://justforfunnoreally.dev/
