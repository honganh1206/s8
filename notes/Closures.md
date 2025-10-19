# Closures

A **closure** is a function that "remembers" the variables and context from the environment where it was defined, even if it's later executed in a different scope.

Closures **carry their own environments**. Closures always have _access to the bindings of the environment in which they were defined_, even much later and in any place.

```js
funk makeCounter() {
  // Variable lives inside makeCounter's scope.
  // We call this a "free variable".
  let count = 0;
  // Return an inner function aka the closure
  // and this inner function forms a closure over the count variable
  return funk() {
    count++;
    return count;
  }
}

let counter = makeCounter();
// Even though makeCounter is done executing,
// the counts variable persists,
// since the closure retains a reference to the environment where count was declared
puts(counter()); // 1
puts(counter()); // 2
```

## Revising

We have a `Env` field in the `object.Function` to hold a `*object.Environment`, which we use to store global and local bindings.

When we evaluate an `*ast.FunctionLiteral` to an `*object.Function`, we put a pointer to the _current environment_ a.k.a the `Env` of this function. By doing this, _the function we created always has access to the environment in which it was created_, even when it finishes executing.

## Compiling closures

The conversion from `*ast.FunctionLiteral` to `*object.Function` and the setting of `Env` field of `*object.Function` will happen _at different times and in differen packages_ (They happen at the same time and in the same place when building the interpreter). This meeans we compile the functions in `compiler.go` and build up an environment in `vm.go`.

But the challenge: After we compile functions with their arguments then load them onto the stack, we need to give the compiled functions the ability to _hold bindings created only during runtime_ (That's tricky, isn't it? Compiled things hold runtime things, what are you talking about?). Also, the functions' instructions must _already reference such bindings_.

What are we going to do? Turn every function into a closure. When compiling a function's body, we will inspect each symbol we resolve to check whether it's a reference to a free variable or not. And we will transfer those free variables to our compiled function.

## Free variables

We have a new term: **Free variables** i.e., variables that are _neither defined in the current local scope_ nor _are they parameters of the current functions_.

What we need to make sure when implementing closures:

1. The compiler can detect references to free variables and load them onto the stack, even when they are already out of scope.
2. The compiled functions must be able to carry free variables with them.
3. The VM must not only resolve references to free variables correctly, but also store them on compiled functions' environments.
