# 3-10-Functions and Function Calls

We will be able to do this:

```js
>> let add = funk(a, b, c, d) { return a + b + c + d};
>> add(1, 2, 3, 4);
10
```

And also passing around functions, higher-order functions and closures

```language
>> let callTwoTimes = funk(x, funk) { funk(funk(x)) };
>> callTwoTimes(3, addThree);
9
>> callTwoTimes(3, funk(x) { x + 1 });
5
```

To do so we need to:

1. Define an internal representation of functions on our `object` system
2. Add support for function calls to `Eval`

Note that [functions in our language are treated like any other values](./Function Literals.md): We can bind them to names/use them in expressions/pass them to other functions/return them from functions/etc.
