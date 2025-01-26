# 3-8-Error Handling

We do internal error handling at this point, such as wrong operators, unsupported operations or other errors arise during execution

> [!TIP] Extend error handling
> We keep things minimal, but for a production-ready interpreter we might want to add a stack trace, line, and column numbers of its origin and provide more than just a message

Errors created for unsupported operations must also **prevent further evaluation**. We do so in `evalBlockStatement` and `evalProgram`
