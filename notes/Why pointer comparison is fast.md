# Why pointer comparison is fast

When comparing large structs, we use **pointer comparison** to compare only the memory addresses (typically 4-8 bytes)

```go
type LargeStruct struct {
    field1 string
    field2 []int
    field3 map[string]int
    // ... many more fields
}

// Pointer comparison (fast)
obj1, obj2 := &LargeStruct{...}, &LargeStruct{...}
isEqual := obj1 == obj2  // Just compares addresses (single operation)

// Value comparison (slower)
value1, value2 := LargeStruct{...}, LargeStruct{...}
isEqual = value1 == value2  // Must compare all fields recursively
```

Pointer comparison offers:

1. Fixed size comparison (just addresses)
2. Single CPU operation
3. No need to traverse data structures
4. No recursive comparisons needed

## Performance benefits

Pointer comparison and HashMap have some similar performance benefits

1. Direct access to memory address value

```
Pointer Comparison:
address1 == address2  // Direct number comparison

HashMap:
hash("key") -> direct memory location
```

2. Constant Time Operations of `O(1)`
