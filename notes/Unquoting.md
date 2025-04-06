# Unquoting

What is unquoting? It tells the interpreter/compiler to immediately evaluate/compile whatever is quoted during macro expansion time

We cannot just add another condition to `Eval()` to unquote, since _we have never called `Eval()`_. We cannot rely on the recursive nature of `Eval()` to evaluate the quoted literals

What we will do instead: Traverse the argument passed into `quote`, find the calls to `unquote` and pass the argument of the call to `Eval`.

A twist: We will have the **Modify** step by replacing the whole `*ast.CallExpression` involving `unquote` with the **result** of this call (as another `ast.Node` type) to `Eval`

We turn the result of `unquote` into a new AST node and modify the existing call to `unquote` with this newly created AST node

> [!IMPORTANT]
> But why a `ast.Node` instead of a `object`? The result of `unquote` might be evaluated during macro expansion time, but its result still needs to be a part of the normal program execution. Plus, doing so ensures the result is properly integrated into the surrounding code structure

[[Modifying functions for unquote]]
