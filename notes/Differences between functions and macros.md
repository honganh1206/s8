# Differences between functions and macros

While macros operate on the code itself (unevaluated expressions), functions operate on evaluated values

1. With Regular Function:

```scheme
// If double was a function:
let double = fn(x) { x + x };
double(2 + 3)

// Evaluation order:
1. First evaluates (2 + 3) = 5
2. Then calls double(5)
3. Result: 5 + 5 = 10
```

2. With Macro:

```scheme
let double = macro(x) { quote(+ unquote(x) unquote(x)) };
double(2 + 3)

// Expansion order:
1. Macro receives unevaluated expression (2 + 3)
2. Expands to: (+ (2 + 3) (2 + 3))
3. Then evaluates: (2 + 3) + (2 + 3) = 5 + 5 = 10
```

Key points

1. Macros receive AST nodes, not evaluated values
2. The expansion happens before evaluation
3. You can generate new code structures using quote/unquote
