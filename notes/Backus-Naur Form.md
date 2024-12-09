
# BNF
- A formal notation used to describe the syntax of programming languages or similar structured languages
- Uses a set of production rules with these elements:
  - `<symbol>` - Non-terminal symbols (variables) in angle brackets
  - Terminal symbols (literal values) without brackets
  - `::=` means "is defined as"
  - `|` means "or" (alternative)

```bnf
// Example of BNF
<expression> ::= <term> | <expression> + <term>
<term> ::= <factor> | <term> * <factor>
<factor> ::= <number> | ( <expression> )

// More concrete examples - EcmaScript syntax
PrimaryExpression ::= "this"
| ObjectLiteral
| ( "(" Expression ")" )
| Identifier
| ArrayLiteral
| Literal

Literal ::= ( <DECIMAL_LITERAL>
| <HEX_INTEGER_LITERAL>
| <STRING_LITERAL>
| <BOOLEAN_LITERAL>
| <NULL_LITERAL>
| <REGULAR_EXPRESSION_LITERAL> )

Identifier ::= <IDENTIFIER_NAME>
ArrayLiteral ::= "[" ( ( Elision )? "]"
| ElementList Elision "]"
| ( ElementList )? "]" )

ElementList ::= ( Elision )? AssignmentExpression
 ( Elision AssignmentExpression )*

Elision ::= ( "," )+

ObjectLiteral ::= "{" ( PropertyNameAndValueList )? "}"

PropertyNameAndValueList ::= PropertyNameAndValue ( "," PropertyNameAndValue
| "," )*

PropertyNameAndValue ::= PropertyName ":" AssignmentExpression

PropertyName ::= Identifier
| <STRING_LITERAL>
| <DECIMAL_LITERAL>
```
