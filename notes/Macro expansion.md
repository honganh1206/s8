# Macro expansion

TLDR: Text replacement - The macro name is replaced with its defined content

```c
#define SQUARE(x) ((x) * (x))

int main() {
    int result = SQUARE(5);  // During macro expansion, this becomes: int result = ((5) * (5));
}
```

Macros take the source code as input (code as data) and return the source code. That is why we "expand" the source code, because each call might result in more code

The point of macros is that **they receive the raw syntax/code rather than evaluated values** (which is what functions do). It is for **code transformation** - the new code that will be evaluated later during the program execution

```js
// Create a new syntax structure to transform into the existing code
let unless = macro(condition, consequence) {
    quote(if (!unquote(condition)) {
        unquote(consequence)
    });
};

unless(x > 5, puts("x is not greater than 5"));
```

The macro expansion takes a part of the AST and modifies it before it is evaluated.

## When do we do macros?

The macro expansion phase is between the 2nd phase (parsing) and 3rd phase (evaluation)

Why before the evaluation phase? Because if not then it will be too late to evaluate/compile :)

## Steps

1st step: Extract the macro definition out from the AST and save it for later modification.

> [!WARNING]
> The macro removal is necessary to avoid tripping over on the macros in the evaluation phase

2nd step: Find the calls to macros and evaluate them

> [!IMPORTANT]
> In this phase we do not evaluate the arguments of the call before evaluate the body. More on [[Differences between functions and macros]]

## How it differs from unquote

Where an `unquote` only causes its single argument to be evaluated, macro calls result in _the body of the macro being evaluated_, and the arguments are available in the environment
