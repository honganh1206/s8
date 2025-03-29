# Macro system

Why? Enable code transformation and abstraction at compile-time. However, macros sometimes are _pure-text replacement_ (like in C), so they are **error-prone** due to the lack of type checking

When? When functions and regular constructs are not enough e.g., generating complex compile-time code

```c
#define MIN(a, b) ((a) < (b) ? (a) : (b)) // Work with any type

#include <stdio.h>

int main() {
    printf("Min of 5 and 10: %d\n", MIN(5, 10));
    printf("Min of 3.2 and 4.5: %.1f\n", MIN(3.2, 4.5));
    return 0;
}
```

Two broad categories: **[[Text-substitution macro systems]]** (search and replace) and **[syntactic macro systems](./Syntactic macro systems.md)** (code as data)

## What we will do

Our language will have the `quote` and `unquote` functions similar to those of Elixir (which is modelled after the `define-macro` system of Lisp and Scheme).

How it works:

```bash
>> quote(foobar);
QUOTE(foobar)
>> quote(10 + 5);
QUOTE((10 + 5))
>> quote(foobar + 10 + 5 + barfoo);
QUOTE((((foobar + 10) + 5) + barfoo))
```

We will also have **macro literals**

```bash
# NOte: These calls will be evaluated in a different way
>> let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
>> reverse(2 + 2, 10 - 5);
1
```

Think of it like this: `quote` tells us to "skip this part", while `unquote` tells us to "except this one, evaluate this"

[[Unquoting]]

```go

```
