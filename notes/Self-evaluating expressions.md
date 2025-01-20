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
