# Expressions in S8

 - Everything besides `let` and `return` statements is an expression. Consider some of them here:
  - Prefix operator `!true`
  - Infix operator `5 + 5`
  - Comparison operator `foo > bar`.
  - Identifiers as expressions `add(foo, bar)`
  - Function literals as first-class `let add = fn(x, y) { return x + y; }`
  - In-line If expression (not many have this) `let result = if (10 > 5) { true } else { false }`

## Brushing up some terminologies

- Expression: A question that always has an answer. Always produce a value
- Statement: A command or complete sentence. It does things.

```javascript
// Expressions (they produce values):
5 + 3        // produces 8
"hi" + "!"   // produces "hi!"
getName()    // produces whatever name it finds

// Statements (they do things):
let age = 5;         // declares a variable
if (hungry) {...}    // makes a decision
return cake;         // gives something back
```

- Prefix operator: An operator *in front of* its operand like `--5`
- Postfix operator: An operator *after* its operand like `foobar++`
- Infix opeartor: An operator that *sits between* its operands like `5 * 8`
- Binary expressions: Where the operator has two operands
- Operator precedence/Order of operations: The priority different operators have like `5 + 5 * 10`

## Preparing the Expression Statement

- This might sound weird: We need to add *expression statements* - a statement consists solely of **one expression** like a wrapper (TLDR: An expression with the `;` at the end)


```js
let x = 5;
x + 10; // Legal expression statement
```

- This feature exists mostly in scripting languages to *have one line consisting only of an expression*

- Tip: We add the `String()` method to our `Node` interface for debugging purposes!
