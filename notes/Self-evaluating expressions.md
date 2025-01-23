# Self-evaluating expressions

Actually **literals in the land of `Eval`**

They are the easiest to evaluate because **they evaluate themselves**: If you type `5` into REPL the program should return `5`

## Integer literals

What to expect: `Eval` should return an `*object.Integer` whose `Value` field _contains the same integer_ as `*ast.IntegerLiteral.Value`

## Boolean literals

Note that we only have 2 values: `true` and `false`. Instead of creating a new `object.Boolean` every time we encounter a `true` or `false`, we **reference** them instead of creating new instances

```go
// Global variables
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)
```

The same thing works for `null`: We only need ONE `null` instance to reference throughout our evaluator

## Prefix Expressions

Prefix expressions in this section are just _operator expressions with one operator and one operand_

> [!IMPORTANT]
> We are defining the semantics of our language. A small change in the evaluation of operator expressions might cause something unintended in a part of a language that seems unrelated.

We first evaluate its operand and then use the result of this evaluation with the operator

## Infix Expressions

There are two groups of infix expressions: The group that produces boolean values e.g., `5 != 5` and the group that produces other types of values e.g., `5 + 5`

Our language will NOT support boolean arithmetic like adding/subtracting/dividing/multiplying booleans like C/C++ or Python

> [!IMPORTANT]
> We are doing **pointer comparison** here as our method arguments are _pointers instead of values_, and pointer comparison is much faster since _we are comparing the same instances_.

> [!WARNING]
> We can use pointer comparison for booleans as operands, but we cannot do so with integer operands since _we are always creating new pointers_. For this reason, `5 == 5` would be false since we are comparing pointers with `object.Object` as value wrapper.
