# Strings, Arrays and Hashes

## Compiling arrays

Arrays are the 1st composite data type i.e., they are composed out of other data types.

> We cannot treat array literals as constant expressions.

The values of an array can change during either compile time or runtime. Since rrays could be _any_ type of expression - integer literal, string concatenation, function literal, etc., _we can reliably determine what they evaluate to during runtime_.

This leads to a design change: Instead of building an array during compile time and pass it to the VM, we tell the VM to **build an array on its own**.

Flow: Compile all elements of an array to N values on the stack -> Emit `OpArray` instruction with operand as N.
