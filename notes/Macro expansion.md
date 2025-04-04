# Macro expansion

The 4th phase: The macro expansion phase - Between the 2nd phase (parsing) and 3rd phase (evaluation)

What we do? We evaluate all calls to macros and replace them with the return value of this evaluation

Macros take the source code as input (code as data) and return the source code. That is why we "expand" the source code, because each call might result in more code

Why before the evaluation phase? Because if not then it will be too late to evaluate :)

The macro expansion takes the AST and modifies it before it is evaluated.

## Steps

1st step: Extract the macro definition out from the AST and save it for later modification.

> [!WARNING]
> The macro removal is necessary to avoid tripping over on the macros in the evaluation phase

2nd step: Find the calls to macros and evaluate them

> [!IMPORTANT]
> In this phase we do not evaluate the arguments of the call before evaluate the body. More on [[Differences between functions and macros]]

## Where do we find the macros in the AST?
