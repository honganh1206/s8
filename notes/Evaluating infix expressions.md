# Evaluating infix expressions

There are two groups of infix expressions: The group that produces boolean values e.g., `5 != 5` and the group that produces other types of values e.g., `5 + 5`

Our language will NOT support boolean arithmetic like adding/subtracting/dividing/multiplying booleans like C/C++ or Python

> [!IMPORTANT]
> We are doing **pointer comparison** here as our method arguments are _pointers instead of values_, and pointer comparison is much faster since _we are comparing the same instances_.

> [!WARNING]
> We can use pointer comparison for booleans as operands, but we cannot do so with integer operands since _we are always creating new pointers_. For this reason, `5 == 5` would be false since we are comparing pointers with `object.Object` as value wrapper.
