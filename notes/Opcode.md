# Opcode

The "operator" part of the instruction, sometimes also called "op". The opcode tells the VM what operation it will perform

Examples would be `0` to push a constant onto a stack, `1` to add top two numbers on the stack and `2` to jump to another instruction

An opcode is one byte wide, and the `PUSH` or `POP` are called **mnemonics**

The operands/arguments/parameters are also included in the bytecode. Unlike the opcodes, they do not have to be one-byte wide

In case an operand needs multiple bytes to be accurately represented, we need to consider how it is encoded - either little endian or big endian

```js
Instruction         Instruction Instruction
    |                   |           |
    |                   |           |
    v                   v           v
+-------------+     +-----------+   +-----+
| PUSH    505 |     | PUSH  205 |   | ADD |
+-------------+     +-----------+   +-----+
    ^     ^             ^   ^
    |     |             |   |
    |     |             |   |
 Opcode Two Byte     Opcode One Byte
       Operand             Operand
```

The mnemonics like `PUSH` do not show up in the actual bytecode, but replaced by the opcodes they refer to

Why `OpConstant` and not `OpPush`? We push **constants** from the constant pool, while `OpPush` is ambiguous - We do not know what we are pushing, it might be a constant, a variable a literal, etc.
