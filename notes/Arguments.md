# Arguments

Arguments to function calls are a _special_ case of local bindings.

While local bindings are created explicitly by the user, arguments are _implicitly bound to names_.

We are going to store arguments on the stack. We push the arguments to the stack _right after the function has been pushed_.

But this approach begs the question(s): How do we know the exact number of arguments that are on top of the function object? And how do we reach that function for execution?

-> We give the `OpCall` opcode an operand that _holds_ the number of arguments of the call.

Compiling the arguments is as easy as iterating over each argument and recursively compile it. But the challenge _to use those arguments in the function's body_.

> At the time of a function call, the arguments will not sit on the stack. How do we access them while the function is executing?

We treat arguments in a way _no different than local bindings_ created in the same function. And we will treat them the same.

A coincidence: Arguements sit right above the function call, and that's also where local bindings are. So _we treat arguments as locals, and they would be exactly where they need to be_.

## A sinister bug behind the base pointer

We treat arguments as locals, so the VM must do some calculations with the base pointer and the number of locals. But there lies an issue: Our stack pointer is too high on the stack.

Expectation: The base pointer should remain constant while the stack pointer changes when new stuff is pushed into the stack. Right before we execute a function, we set basePointer to the current value of the stack pointer. Then we increase the stack pointer when we push new values onto the stack, thus creating a memory region to store local bindings and arguments.

```css
         vm.sp ──▶  +-----------+
                    |           |   ◀── basePointer+2
                    +-----------+
                    |   Arg 2   |   ◀── basePointer+1
                    +-----------+
                    |   Arg 1   |   ◀── basePointer
                    +-----------+
                    | Function  |
                    +-----------+
```

Hidden bug: After we pushed the arguments on the stack, **the base pointer and the stack pointer are set to the same value**. And when we do the `basePointer + local binding index` formula, the base pointer points to empty/undefined region of the stack.

```css
         vm.sp ──▶  +-----------+
                    |           |   ◀── basePointer+2
                    +-----------+
                    |           |   ◀── basePointer+1
                    +-----------+
                    |           |   ◀── basePointer
                    +-----------+
                    |   Arg 2   |
                    +-----------+
                    |   Arg 1   |
                    +-----------+
                    | Function  |
                    +-----------+
```

Solution: A new formula `basePointer = vm.sp - numArguments` to anchor the base pointer AFTER we have pushed arguments onto the stack.
