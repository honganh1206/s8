# Local bindings

Local bindings are local to functions, meaning they are only visible and accessible within a scope of a function.

```js
let globalSeed = 50;

let minusOne = fn() {
  let num = 1;
  globalSeed - num;
}
let minusTwo = fn() {
  let num = 2;
  globalSeed - num;
}

minusOne() + minusTwo()
```

For local bindings we need a new store for the VM, a separate one from the global bindings.

We need to extend the symbol table to handle to tell different scopes apart and determine what symbol belongs to which scope by _keeping track of the scope the symbols are in_.

> Defining symbols in one local table must not interfere with the definitions in another local table, and resolving global symbols in a nested local table still resolves to the correct symbols.

For the VM, we need to ensure the storage for the bindings are different: We create bindings with `OpSetLocal` and resolve such bindings with `OpGetLocal` - all of such operations must happen within their scope.

We are going to store locals on the stack, and the effort will be worth it.

We are going to reuse the stack field of our VM. What we would do: When we come across `OpCall`, we make a **hole** a.k.a a region of the stack, and the hole should be large enough to _store the locals of our function call_. Below the hole would be _all the values previously pushed on the stack_.

How do we create the hole? We increase the stack pointer _by the number of locals_ used by the function we're about to execute (remember to store the current stack pointer somewhere else).

By re-using the same stack, when we are done executing a function, we can just restore the stack to the state before the execution.

The `Frame` struct will be where we store the state of the stack before the function call.

> Remember to clean up the stack after we are done executing the function.
