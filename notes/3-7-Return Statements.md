# 3-7-Return Statements

The return statements could be used _in the bodies of functions_ but also as _top-level statements_ in a Monkey program. They should also _stop the evaluation of a series of statements_ and _leave behind the value their expression has evaluated to_

```go
5 + 5 + 5;
return 10; // Evaluation stops and the called function should evaluate to 10
9 * 9 * 9;
```

To support return statements, we will be passing a "return value" through our evaluator? Whenever we encounter a `return`, we will wrap the value to be returned inside an `object.ReturnValue` to keep track of it later

## Nested block statements

The problem is that _we cannot unwrap the value of `object.ReturnValue`_ on first sight! We need to keep track of it so we can stop the execution in the outermost block statement. This approach is to _handle early return cases_ in nested blocks/scopes

The point is that our `object.ReturnValue` will remain wrapped until it bubbles up back to `evalProgram()`. At that point it can be unwrapped. To do so we need to call the `evalBlockStatement()` recursively to handle the nested block statements

```
evalProgram
    ↑
    unwrapped to just 10 at evalProgram()
    |
    └── evalBlockStatement (outer if)
        ↑
        | bubbles up unchanged
        ↑
    └── evalBlockStatement (inner if)
        ↑
        | bubbles up unchanged
        ↑
    └── ReturnValue{Value: 10}
```
