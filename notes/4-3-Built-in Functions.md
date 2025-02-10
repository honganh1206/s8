# 4-3-Built-in Functions

Functions built into the interpreter, into the language itself

Such functions are **atomic/single-level** operation implemented directly in the host language

Built-in functions accept zero or more `object.Object` and returns an `object.Object`

We also create a separate environment for built-in functions

We do not need to unwrap the return values when it comes to built-in functions, since built-in functions do not return an `*object.ReturnValue`

```go
func lenBuiltin(args []Object) Object {
    // Direct return of an Object, never wrapped in ReturnValue
    // Since built-in functions do not need to handle nested return statemenets
    return &Integer{Value: int64(len(args[0].(*String).Value))}
}

```
