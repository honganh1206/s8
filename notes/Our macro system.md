# Our macro system

Our language will have the `quote` and `unquote` functions similar to those of Elixir (which is modelled after the `define-macro` system of Lisp and Scheme).

How it works:

```bash
>> quote(foobar);
QUOTE(foobar)
>> quote(10 + 5);
QUOTE((10 + 5))
>> quote(foobar + 10 + 5 + barfoo);
QUOTE((((foobar + 10) + 5) + barfoo))
```

We will also have **macro literals**

```go
// Note: These calls will be evaluated in a different way
>> let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
>> reverse(2 + 2, 10 - 5);
1
```

Think of it like this: `quote` tells us to "skip this part, do not evaluate it", while `unquote` tells us to "except this one, evaluate this"

Our Go function will return a `*object.Quote` containing an **unevaluated** `ast.Node`. Inside this unevaluated node we will [unquote](./Unquoting.md) to evaluate the expressions.

This will evaluate the argument of the `unquote` call and replace the whole call expression (an `ast.Node`) with the **result** of that evaluation

> [!TIP]
> When we pass the `*object.Environment` into our call to `unquote`, we can do environment-aware evaluation inside unquote calls
