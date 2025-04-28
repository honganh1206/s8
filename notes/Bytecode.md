# Bytecode

Bytecode is an **intermediate representation** of code that _sits between source code and machine code_. It is **NOT** native machine code, nor it is assembly language

Bytecode is a binary format and not really readable

Why the "Byte" prefix? Because each instruction is one byte in size

Bytecode is interpreted by [[Virtual Machines]] that is _part of the interpreter_.

Just like VMWare and VirtualBox emulate real machines and CPUs, such VMs _emulate a machine that particularly understands this specific bytecode format_.

The exact format of a bytecode and the **opcodes** (the instructions that make up the bytecode) a bytecode is composed of _varies and depends on the guest and host programming languages_

Consider this Python example:

```py
x = 5 + 3
```

After Python compiles it, it generates bytecode that might look like this

```bytecode
LOAD_CONST 0 (5)
LOAD_CONST 1 (3)
BINARY_ADD
STORE_NAME 0 (x)

```

Each instruction i.e., `LOAD_CONST, BINARY_ADD, STORE_NAME` is an **[[Opcode]]** - a basic instruction that the Python VM understands

> [!TIP]
> Python uses both an interpreter and a compiler! It generates `.pyc` bytecode files and stores them in `_pycache_` directory.

In general, the opcodes are _pretty similar_ to the mnemonics of most assembly languages, so most bytecode implementations share opcodes for `push` and `pop` for stack operations

[[Assembling and disassembling]]

Bytecode can be specialized and domain-specific. It is the machine language for the custom-built VM. Sometimes it contains domain-specific instructions for the domain-specific VM like the JVM with `invokeinterface` opcode
