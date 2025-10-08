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
