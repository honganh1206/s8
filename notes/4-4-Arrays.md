# 4-4-Arrays

Our language will have an array as **an ordered list of elements of DIFFERENT types**

```js
>> let myArr = ["Hello", "World", 100, funk(x){x*x}]
>> myArr[0]
Hello
>> myArr[3](2)
4
```

We introduce a new operator: the **index operator** `[]`

The basis of our language's array will be a Go slice of type `[]object.Object`

## Parsing array literals

We can pass comma-separated lists of expressions like [how function call arguments are
passed](./3-10-Functions and Function Calls.md) with `parseCallArguments()`

## Parsing index operator expression

Structure: `<expression>[<expression>]`

The index operator must also **have the highest precedence**

We pass the index operator with an `infixParseFn`, even though there is not a single between
operands. The point is that there is a `[` coming between the identifier/expression and the index
value

## Evaluating array literals

When evaluating array literals, we need to do a thing similar to when we evaluate a call expression:
_Stop the evaluation and return the error immediately_ if there is a problem

## Evaluating index operator expressions

For this we choose to **return NULL** when we use an index that is out of bound of an array

The left operand before `[` can be **any expression** and the index itself could be any expression

## Built-in functions for arrays

`rest` returns a new array containing the elements of the array passed as the argument **except the
1st one**. For this note that we are _returning a newly allocated array_ instead of modyfing the
array passed to `rest`
