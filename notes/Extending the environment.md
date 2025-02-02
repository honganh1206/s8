# Extending the environment

For [function calls](./3-10-Functions and Function Calls.md), we should _preserve bindings while also make new ones available_ with a new environment

We extend the environment by creating a new instance of `object.Environment` with a **pointer** to an environment it should extend to

## How the extended environment works

When the new environment invokes a `Get` method but it does not have a value associated with the given name, it _calls the `Get` method of the enclosing (outer) environment_ to ensure closure

If there is still no associated value, the enclosing environment will keep calling its enclosing environment until _there is no closing environment anymore_. In such cases, an error will be shown

> [!TIP]
> This concept of extending environments is a way to think about variable scopes: If something is not found in the inner scope, we look it up in the outer scope.
> The outer scope **encloses** the inner scope, and the inner scope **extends** the outer one

The language has [closures](Closures.md) and we use it by _passing the extending environment to `Eval` instead of the current environment of the function_
