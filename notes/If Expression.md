# If Expressions

For this we will have to deal with _different token and expression types_

The If Expressions are similar to the Ifs in other languages. Note that the `else` is optional

```js
if (x > y) {
  return x;
} else {
  return y;
}

// Or just this
if (x > y) {
  return x;
}
```

Note that in s8 **if-else-conditionals are expressions**, so we do not need the return statement here

```js
let foobar = if (x > y) { x } else { y };
```

Structure of our if-else conditionals: `if (<condition>) <consequence> else <alternative>`

Our `<consequence>` and `<alternative>` will be of type `BlockStatement` as they are a _series of statements_
