# Compiling Expressions

[[Cleaning up the stack]]

Booleans also exist as literal expressions e.g., `true;`

Comparison operators compare the two topmost elements on the stack, then tell the VM to pop them off and push the result back on

Why no `OpLessThan` and only `OpGreaterThan`? So we can do the **reordering of code**. Our compiler will reorder `3 < 5` to `5 > 3` to _keep the instruction set small_ and _keep the loop of our VM tighter_
