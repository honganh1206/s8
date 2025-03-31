# Macro system

Why? Enable code transformation and abstraction at compile-time. However, macros are (sometimes) _pure-text replacement_ (like in C), so they are **error-prone** due to the lack of type checking

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

[[Our macro system]]

[[Macro expansion]]
