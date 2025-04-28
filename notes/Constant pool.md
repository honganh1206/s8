# Constant pool

At compile time, all values (as literals) are stored in one big array/vector

An instruction that needs a constant carries a small integer operand called **the index**

The VM/runtime keeps the pool in memory at some base address `BP`

To fetch entry `i` the VM computes something like address `address = BP + i * sizeof(slot)`

Example with a JVM-style VM

```text
ldc #7    ; “load constant pool entry 7”
```

At execution, the JVM does `value = constant_pool[7]` and pushes the `value` onto the operand stack
