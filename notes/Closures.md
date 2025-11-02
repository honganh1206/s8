# Closures

A **closure** is a function that "remembers" the variables and context from the environment where it was defined, even if it's later executed in a different scope.

> TLDR: Closure = Function + Its remembered environment

Closures **carry their own environments**. Closures always have _access to the bindings of the environment in which they were defined_, even much later and in any place.

```js
function makeCounter() {
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
// since the closure retains a reference (maybe a pointer) to the environment where count was declared
puts(counter()); // 1
puts(counter()); // 2
```

## ELI5

Imagine you have a **backpack** that can hold things (variables).

When you create a closure, you are making a _function_ that carries a backpack _filled with the variables from where it was born_. Even after you leave that place, the function still keeps and remembers what's inside the bag.

## Revising

We have a `Env` field in the `object.Function` to hold a `*object.Environment`, which we use to store global and local bindings.

When we evaluate an `*ast.FunctionLiteral` to an `*object.Function`, we put a pointer to the _current environment_ a.k.a the `Env` of this function. By doing this, _the function we created always has access to the environment in which it was created_, even when it finishes executing.

## Compiling closures

The conversion from `*ast.FunctionLiteral` to `*object.Function` and the setting of `Env` field of `*object.Function` will happen _at different times and in differen packages_ (They happen at the same time and in the same place when building the interpreter). This meeans we compile the functions in `compiler.go` and build up an environment in `vm.go`.

But the challenge: After we compile functions with their arguments then load them onto the stack, we need to give the compiled functions the ability to _hold bindings created only during runtime_ (That's tricky, isn't it? Compiled things hold runtime things, what are you talking about?). Also, the functions' instructions must _already reference such bindings_.

What are we going to do? Turn every function into a closure. When compiling a function's body, we will inspect each symbol we resolve to check whether it's a reference to a free variable or not. And we will transfer those free variables to our compiled function.

Or TLDR: Wrap `*object.CompiledFunction` inside `*object.Closure` when we execute `OpClosure`

What we need to make sure when compiling closures:

1. The compiler can detect references to free variables and load them onto the stack, even when they are already out of scope.
2. The compiled functions must be able to carry free variables with them.
3. The VM must not only resolve references to free variables correctly, but also store them on compiled functions' environments.

## Free variables

We have a new term: **Free variables** i.e., variables that are _neither defined in the current local scope_ nor _are they parameters of the current functions_.

> TLDR: We treat every non-local, non-global and non-built-in as a free variable.

"Free variables" is a relative term. A free variable in the current scope could be a local binding in the enclosing scope. See the example below:

```js
let a = 1;
let b = 2;
let firstLocal = funk() {
    // c and d are both free variables and local variables here
    let c = 3;
    let d = 4;
    a + b + c + d;
        let secondLocal = funk() {
            let e = 5;
            let f = 6;
            a + b + c + d + e + f;
        };
};
```

[Bonus: The point of closures -> TLDR: A generalization of a class](https://stackoverflow.com/questions/1305570/closures-why-are-they-so-useful)
[Bonus: When to use closure - Anytime you like, since it's an alternative to classes](https://stackoverflow.com/questions/256625/when-to-use-closure)

A few checks before we resolve a free variable:

- Has the name of the variable defined in the current scope/symbol table?
- Is it a global binding or a built-in function?

If all answers are no, it means the variable is defined as a local variable in the enclosing scope. Thus, it should be resolved as a free variable.

## Recursive closures

There might be an edge case where `OpGetLocal` is set BEFORE `OpSetLocal`. Have a look at the following example:

```js
let wrapper = function () {
  let countDown = function (x) {
    if (x == 0) {
      return 0;
    } else {
      countDown(x - 1);
    }
  };
  countDown(1);
};
wrapper();
```

Recall: The compiler iterates the symbols marked as free and emits necessary load instructions to get the free variables on to the stack.

While compiling the body of `countDown`, the compiler came across the reference to `countDown` and asked the symbol table to resolve it.

The symbol table realized that there is no symbol with the name `countDown` in the current scope, so it marked `countDown` as a free variable (Refer to `Resolve()` logic in `symbol_table.go`), ready to be loaded on to the stack.

After compiling the body of `countDown` and right before emitting the `OpClosure` instruction to turn `countDown` to a closure, the compiler iterates over the symbols marked as free and loads them on to the stack.

When the VM executes the `OpClosure` instruction that comes after the load instructions, it should be able to access the free variables and transfer them onto the `*object.Closure` (the backpack) the compiler creates.

However, when the VM executes the `OpClosure` instruction that comes after the variable-loading instructions, _the local index 0 has not yet been saved_. Instead, we have a Go `nil` on the stack, but that's _where the closure itself should end up_. (TLDR: We are executing something we have yet to create)

`countDown` is itself a free variable, and it refers to itself (recursion, duh!). But we have yet to turn it into a closure before we create a reference to it.

The solution?

1. We emit a new opcode instead of marking the recursive function as a free variable.
2. We make the compiler know the name of the function we are compiling.
3. We let the symbol table know when there is a self-reference by defining a new scope called `FunctionScope`. Only one symbol with that scope per symbol table. When we resolve a name to get a symbol with that scope, we know that it's the name of the current function we're compiling.
