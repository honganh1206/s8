# Central Processing Unit

The basic design of a computer according to the Von Neumann architecture

```js
              +-------------+
              |    Input    |
              +------+------+
                     |
                     v
 +------------------------------------------+
 |                                          |
 |       +------------------------+         |
 |       |        Memory          |         |
 |  +----+------------------------+----+    |
 |  |    |         |        |     |    |    |
 |  |    |         |        |     |    |    |
 |  |    |         |        |     |    |    |
 |  +----+---------+--------+-----+----+    |
 |       ^                   ^              |
 |       |                   |              |
 |       v                   v              |
 |  +----+---------------------------+      |
 |  |  Central Processing Unit (CPU) |      |
 |  +----+---------------------------+      |
 |       ^                   ^              |
 |       |                   |              |
 |       v                   v              |
 |  +----+---------------------------+      |
 |  |        Mass Storage            |      |
 |  +--------------------------------+      |
 |                                          |
 +------------------------------------------+
              |
              v
      +---------------+
      |    Output     |
      +---------------+
```

How it works:

1. Fetch an instruction from memory
2. Decode the instruction
3. Execute the instruction
4. Repeat from step 1

But how does the CPU address different parts of memory? The CPU uses numbers as addresses when accessing data in memory

The idea of telling the CPU where to store and retrieve data is good learning, but memory access today is abstracted away by layers of security and performance optimizations

A computer's memory also stores **programs** - the instructions fetched to the program counter - but it goes to a region in memory different from where the memory stores data

If the CPU attempts to decode data that is not a valid instruction, the CPU will respond in a way depending on how it is designed e.g., trigger an event, give the program a chance to recover, stop, etc.
