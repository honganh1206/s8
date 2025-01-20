# Expose values of the host language to the user of the interpreted language

TLDR: Which approach to choose depends on the **language design** and **performance requirements**[ExposeTake Java for example: It has both **primitive data types** (holding actual data values in memory) e.g., int, byte, short, etc., and **reference types** (store memory addresses pointing to objects) to the user

In Java, _the primitive data types do not have a huge representation inside the Java implementation_ but **closely map** to their native counterparts i.e., they **correspond directly** to how the computer hardware handles basic data types

Reference types on the other hand _require extra layers of abstraction and processing by the JVM_ before they can be used by the computer hardware

```java
// Primitives: creates a real copy
int a = 5;
int b = a;     // b gets its own 5
b = 10;        // only b changes to 10, a stays 5

// References: just copies the directions
StringBuilder s1 = new StringBuilder("hello");
StringBuilder s2 = s1;    // s2 points to the same data as s1
s2.append(" world");      // changes what both s1 and s2 point to
// Both s1 and s2 now show "hello world"

```

For Ruby, users _does not have access to primitive data types_: Everything is an object and thus wrapped inside an internal representation or TLDR: Everything is the same value type with each other, wrapping different values
