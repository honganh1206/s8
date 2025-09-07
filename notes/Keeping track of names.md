# Keeping track of names

Main task: Have the identifiers correctly resolve to the values they were bound to.

We will use numbers to represent the identifiers.

Two new opcodes: `OpGetGlobal` and `OpSetGlobal`, each having a 16-bit wide operand the identifier's index as a number.

We use a slice to store global bindings, and the new opcodes will pop/get the value from the stack from the slice.

We use a **symbol table** to associate identifiers with information like scope, location, type, etc.

We first "define" an identifier and associate some information with it. Later on we "resolve" the identifier to this information.
