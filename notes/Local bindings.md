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
