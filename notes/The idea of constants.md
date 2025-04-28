# The idea of constants

In this context, **constant** is short for _constant expression_ and refers to expressions whose value does not change and can be determined at compile time

We do not need to run the program to know what the constant expressions evaluate to, since _the compiler can find them in the code and store the value they evaluate to_

```js

+--------+    +--------+    +---------+    +----------------+
| Lexer  | -> | Parser | -> | Compiler| -> | Virtual Machine|
+--------+    +--------+    +---------+    +----------------+

     Compile Time                               Run Time
<-------------------------------------->  <----------------->
```

After that, the compiler can reference the constants in the instructions it generates instead of embedding the values directly in them, since the generated instructions _already contain the constant_ (either as IR or in memory)

A plain integer can serve as an index into a data structure that holds all constants, often called a [[Constant pool]]
