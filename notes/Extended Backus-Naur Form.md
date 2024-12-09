# EBNF

- An **extension** of BNF with additional notation for *easier and more concise grammar specification*
- Adds these elements:
  - `[]` for optional items (0 or 1 occurrence)
  - `{}` for repetition (0 or more occurrences)
  - `()` for grouping
  - `*` for zero or more repetitions
  - `+` for one or more repetitions
  - `?` for optional elements

Example in EBNF:
```ebnf
expression = term { ("+" | "-") term } ;
term = factor { ("*" | "/") factor } ;
factor = number | "(" expression ")" ;
```
