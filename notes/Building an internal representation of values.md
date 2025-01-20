# Building an internal representation of values

There are different choices when building an internal representation of values: Some use native types (integers, booleans, etc.) of the **host language** to represent the vaues of the interpreted languages, while some represent values/objects as pointers, and some have them mixed up

But the most important thing: **_How you represent a string of your interpreted language depends on how a string can be represented in the language the interpreter is written in_** i.e., An interpreter written in Ruby cannot represent values the same way an interpreter written in C can

Another important thing: Some interpreted languages only need representations of _primitive data types_ like integers or bytes, but some may need representations of lists, dictionaries, functions or compound data

Also: We must consider **the speed and the memory consumption while evaluating programs**. There cannot be a fast interpreter with a slow and bloated object systems

Adding a garbage collector means we have to _keep track of the values in the system_. But that is for another day if you do not really care about the performance

> [!IMPORTANT]
> Read the source code of some popular interpreters such as the [Wren source code](https://github.com/munificent/wren) to learn about different representations
