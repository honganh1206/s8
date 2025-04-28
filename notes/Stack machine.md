# Stack machine

Goal: Translate `1 + 2` to bytecode instructions that make use of a stack

How? Get the operands `1` and `2` from the topmost of the stack and tell the stack to add them :)

```js

 BEFORE           AFTER
+-----+          +-----+
|  2  |          |  3  |
+-----+          +-----+
|  1  |          |     |
+-----+          +-----+
```

Two instructions: One to push the operands to the stack, and another to add things on the stack

The opcode will not have "push" in its name, since it will not be solely about pushing things

Why? It is easy for the VM to take integer operands (fixed-size) and push it to the stack, but what if we want to work with string literals (variable-length dependent)?

[[The idea of constants]]

When we come across an integer literal (A constant expression), we keep track of the resulting `*object.Integer` by storing it in memory and assign to it a number

When we are done compiling and pass the instructions to the VM, we will hand over to the VM _all the constants we have found_ by putting them all in a data structure (the constant pool). The VM will then use the numbers assigned to the constants as indexes to retrieve values
