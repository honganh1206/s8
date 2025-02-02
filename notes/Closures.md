# Closures

Closures are **basically functions**, and they "close over" the environment they were defined in?

Closures **carry their own environments**

```js
// newAdder is a higher-order function that returns another function
let newAdder = funk(x) { funk(y) { x + y } };
// Here newAdder returns addTwo as a CLOSURE - it can access the value bound to parameter x
// The argument that addTwo receives will be bound to parameter y in the returning function
let addTwo = newAdder(2);
addTwo(3); // Return 5
```
