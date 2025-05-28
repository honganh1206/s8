# Cleaning up the stack

After executing an operation like `ADD` we need to store the result somewhere outside of the stack

`ExpressionStatement` is not like `let` or `return` statement since the statement part works more like a _wrapper_ so the expression can occur on their own

What we need to do:

1. Define a new opcode to tell the VM to pop the topmost element.
2. Emit this opcode after _every_ expression statement
