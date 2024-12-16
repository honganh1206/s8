# Pratt Parsing/Top-down Operator Precedence

- The "Top Down Operator Precedence" by Vaughan Pratt is recently rediscovered and later popularized.
- This was invented as an *alternative* to parsers based on context-free grammars and the [[Backus-Naur Form]]

## The main difference

- Instead of associating parsing functions with grammar rules (like we did with `parseLetStatement()`, Pratt associated these functions with **single token types**. In detail, *each token type can have TWO associated parsing functions* depending on whether the token denotes as a prefix or an infix.
