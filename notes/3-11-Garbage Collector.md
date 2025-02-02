# 3-11-Garbage Collector

Consider this snippet

```js
let counter = funk(x) {
  // If-else expression is evaluated again and again until the produced value is truthy
  // Each time this is about to be evaluated, a new object.Integers (100) is generated
  if (x > 100) {
    return true;
  } else {
  let foobar = 9999; // Value gets bound to the name foobar but never referenced again
    counter(x + 1);
  }
};

counter(0);
```

In each call, **a lot of objects are allocated**. Our objects are stored in **memory**, and the more objects we use, the more memory we need

The reason we are not running out of memory is that _Go's garbage collector (GC) is working under the hood_ since we are using Go as the host language.

Why don't we write our own GC? Then we have to disable Go's GC and allocate and free memory ourselves - in a language that by default prohibits doing so!
