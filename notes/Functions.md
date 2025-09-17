# Functions

Implementing functions is challenging:

- How do we represent functions and even higher-order ones?
- How do we manage control flow? How do we get the VM to execute instructions in a function? Apparently we cannot get them all mixed up. And how to make return the control flow to the main execution, given that we have explicit `funk() { return 5; }` + implicit returns `funk() { 5; }`? And even a function that does not return at all `funk() {}`!
- How do we pass arguments into functions?

## Representing functions

Functions are essentially values. Functions can be bound to names, returned from other functions, passed to other functions as arguments, etc., and top of all, _they are produced by expressions_

The value the function literals produce _does not change_, just like other literals.

Since we compile functions to `*object.CompiledFunction` and treat them as constants, we will load them onto the stack using `OpConstant`

### Return from a function

Implicit and explicit returns will be compiled into the same bytecode.

Explicit return: We first compile the return statement so the return value ends up on top of the stack, then we omit a `OpReturnValue`.

Functions that return nothing are an edge case and for that, we will return `Null`

For a function declaration like `funk() { return 5 + 10; }`, we should have:

1. `OpConstant 0` to load 5 on the stack
2. `OpConstant 1` to load 10 on the stack
3. `OpAdd` to add them together
4. `OpReturnValue` to return value at the top of the stack

## Compiling functions

We cannot just call `Compile` to the `Body` field of `*ast.FunctionLiteral` - it would make a mess of instructions entangled with instructions of the `main` program. So what to do?

We use **scopes** - bundle our instructions in a compilation scope and use a stack of compilation scopes

> We need to carefully separate the instructions from a function (scoped) and the ones from the main flow.
