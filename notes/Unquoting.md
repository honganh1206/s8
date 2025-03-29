# Unquoting

No, we cannot just add another condition to `Eval()` to unquote, since _we have never called `Eval()`_. We cannot rely on the recursive nature of `Eval()` to evaluate the quoted literals

What we will do instead: Traverse the argument passed into `quote`, find the calls to `unquote` and pass the argument of the call to `Eval`.

A twist: We will have the **Modify** step by replacing the whole `*ast.CallExpression` involving `unquote` with the **result** of this call (as another `ast.Node` type) to `Eval`

We turn the result of `unquote` into a new AST node and modify the existing call to `unquote` with this newly created AST node

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
