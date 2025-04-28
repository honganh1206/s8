# The smallest compiler

What we do:

1. Traverse the AST to find the `*ast.IntegerLiteral` nodes
2. Evaluate them into `*object.Integer` objects
3. Add the objects to the constant pool
4. Emit the `OpConstant` instructions that reference such constants
