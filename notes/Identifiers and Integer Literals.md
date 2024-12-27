# Identifiers and Integer Literals

 ## More on identifiers

- *Identifiers are expressions too*, so they can be used in all contexts such as `foo + bar`

- We add an *optional* check for semicolon when parsing expression statements - so we can do things like `5 + 5` in our REPL just like JS!


- We will have `prefixParseFn` function type to handle cases when we *encounter the associated token type in prefix position*, and `infixParseFn` to handle *token type in infix position case*

- The parsing functions of both types `prefixParsingFn` and `infixParsingFn` will follow this **same protocol**: We return the AST node with the current token and the literal value of the token **without** advancing the token. This is an important principle in [[Recursive-descent Parsing]]


## More on Integer Literals

- *Integer literals are expressions*. The value they produce is the integer itself for example `5;`

```js
// Places where integer literals can occur
let x = 5;
add(5, 10);
5 + 5 + 5;
```

- Note that the literal value of the AST node before parsing will be of `string` type and we need to convert that to `int` type when constructing our node.
