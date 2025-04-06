# Modifying functions for unquote

Examples of modifying functions:

1. **Turn all numbers into their negative values**:

```go
modifier := func(node Node) Node {
    if number, ok := node.(*IntegerLiteral); ok {
        number.Value = -number.Value
        return number
    }
    return node
}

// Input AST: 5 + 3
// Output AST: -5 + -3
```

2. **Turn all identifiers to uppercase**:

```go
modifier := func(node Node) Node {
    if ident, ok := node.(*Identifier); ok {
        ident.Value = strings.ToUpper(ident.Value)
        return ident
    }
    return node
}

// Input AST: x + y
// Output AST: X + Y
```

3. **Replace all additions with subtractions**:

```go
modifier := func(node Node) Node {
    if infix, ok := node.(*InfixExpression); ok {
        if infix.Operator == "+" {
            infix.Operator = "-"
        }
        return infix
    }
    return node
}

// Input AST: 5 + 3 + 2
// Output AST: 5 - 3 - 2
```

4. **Complete node replacement example**:

```go
modifier := func(node Node) Node {
    if _, ok := node.(*IntegerLiteral); ok {
        // Replace any integer with boolean true
        return &Boolean{Value: true}
    }
    return node
}

// Input AST: 5 + 3
// Output AST: true + true
```
