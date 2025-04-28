# The stack

Where data is managed in an LIFO manner

The CPU needs the call stack because it needs to know which instruction to execute next, once the current function is fully executed. Without the call stack, the CPU will just execute the instruction _at the next higher address in memory_

The **return address** tells us which instruction to fetch next after the current function

> [!WARNING]
> Instructions are not laid out in memory in a linear order of execution

Using the stack is good for a compiler, since **function calls are often nested**, so the stack is a great data structure to handle order of function calls

Why _The Stack_ and not just a stack? Because the region used to implement a call stack is a **convention** - so common that it's been cast into hardware

GOAL: We are going to implement _a virtual call stack_

We have a memory address _pointing to the top of the stack_ and that address is stored in a register. We call that address a **stack pointer**, and it has a designated register in most CPUs
