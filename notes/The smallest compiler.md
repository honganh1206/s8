# The smallest compiler

What we do:

1. Traverse the AST to find the `*ast.IntegerLiteral` nodes
2. Evaluate them into `*object.Integer` objects
3. Add the objects to the constant pool
4. Emit the `OpConstant` instructions that reference such constants

What the compiler needs to do:

1. Walk the AST _recursively_ and find the `*ast.IntegerLiterals`
2. Evaluate them and turn them into `*object.Integers`
3. Add the objects to the constant pool
4. Add the `OpConstant` instructions to its internal `instructions` slice
