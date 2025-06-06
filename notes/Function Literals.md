# Function Literals

This is our function literals--They will be of type **Expression**:

```js
funk(x, y) {
  return x + y;
}

let myFunction = funk(x, y) { return x + y; }
```

Structure: `funk <parameters> <block statement>` and parameters are just _a list of identifiers_ like this `(<parameter one>, <parameter two>, ...)`

Note that _the parameter list could be empty_

```js
funk() {
  return foobar + barfoo;
}
```

A function literal could work as _an expression in a return statement_, and this return statement is inside another function literal

```js
funk() {
  return funk(x, y) { return x > y; };
}
```

We can even use a function literal as **an argument** for another function literal

```js
myFunc(x, y, funk(x, y) { return x > y; });
```

We can also have **anonymous functions**: We create an unamed, inline function definiton, followed by an argument

```js
funk(x) { x; }(5) // Immediately invoking this anynymous function with the argument 5
// Or written more verbosely
let my_function = funk(x) { x; };
my_function(5);
```
