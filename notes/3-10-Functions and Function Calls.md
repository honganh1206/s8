# 3-10-Functions and Function Calls

We will be able to do this:

```js
>> let add = funk(a, b, c, d) { return a + b + c + d};
>> add(1, 2, 3, 4);
10
```

And also passing around functions, higher-order functions and closures

```js
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

> [!IMPORTANT]
> When evaluating an `*ast.CallExpression`, the expression could be either an identifier (named function e.g., `add(5 + 10)`) or a function literal (anonymous function e.g., `funk(x) { x; }(5)`)
> When evaluating arguments like `add(2 + 2, 5 * 5);`, we need to pass `4` and `25` as arguments, NOT the two expressions `2 + 2` and `5 * 5`

## Evaluating the body

We cannot just call `Eval` and pass the body of the function. The body can contain the **references** to the parameters of the function, and we need to _change the environment in which the function is evaluated_

```js
// This should print two lines 10 and 5 respectively
// But if we overwrite the current environment before evaluating the body of `printNum`
// The last line would result 10 being printed as no binding happen to the parameter happens
let i = 5;
let printNum = funk(i) {
  puts(i); // Assign the value to the parameter-unrelated to the global variable above
}

printNum(10); // Print the 10 here
puts(i); // Print the 5 here
```
