# Bytecode

Bytecode is an **intermediate representation** of code that _sits between source code and machine code_. It is **NOT** native machine code nor it is assembly language

Bytecode is interpreted by a **virtual machine** (VM) that is _part of the interpreter_.Just like VMWare and VirtualBox emulate real machines and CPUs, such VMs _emulate a machine that particularly understands this specific bytecode format_.

The exact format of a bytecode and the opcodes(the instructions that make up the bytecode) a bytecode is composed of _varies and depends on the guest and host programming languages_

Consider this Python example:

```py
x = 5 + 3
```

After Python compiles it, it generates bytecode that might look like this

```language
LOAD_CONST 0 (5)
LOAD_CONST 1 (3)
BINARY_ADD
STORE_NAME 0 (x)

```

Each instruction i.e., `LOAD_CONST, BINARY_ADD, STORE_NAME` is an **opcode** - a basic instruction that the Python VM understands

> [!TIP]
> Python uses both an interpreter and a compiler! It generates `.pyc` bytecode files and stores them in `_pycache_` directory.

In general, the opcodes are _pretty similar_ to the mnemonics of most assembly languages, so most bytecode implementations share opcodes for `push` and `pop` for stack operations
