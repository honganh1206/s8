# Hello Bytecode

What we will do:

Take the expression `1 + 2` -> Tokenize and parse it to a node in the AST -> Compile it to bytecode -> Execute it with the virtual machine

Data structures: String -> Tokens -> AST -> Bytecode -> Objects

[[Stack machine]]

[[The smallest compiler]]

Our VM should be able to fetch, decode and execute `OpConstant` instructions, and the results will be pushed on to the VM's stack

When we do `OpAdd` on two operands pushed to the stack, we need to _pop the operands off the stack_

> In some cases the order of the operands matters e.g., minus operator. We should NOT implicitly assume the `right` operand is the last one to be pushed on to the stack.
