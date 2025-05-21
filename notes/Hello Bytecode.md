# Hello Bytecode

What we will do:

Take the expression `1 + 2` -> Tokenize and parse it to a node in the AST -> Compile it to bytecode -> Execute it with the virtual machine

Data structures: String -> Tokens -> AST -> Bytecode -> Objects

[[Stack machine]]

[[The smallest compiler]]

Our VM should be able to fetch, decode and execute `OpConstant` instructions, and the results will be pushed on to the VM's stack
