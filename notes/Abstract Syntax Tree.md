# The Abstract Syntax Tree

- The AST is a **data structure** used for **internal representation** of the source code.
- The "Abstract" part in the AST means *certain details visible in the source code are omitted** in the AST like whitespaces, comments, braces, etc.
- Note that *there is NO one true, universal AST format*. AST implementations might share the concept but not the details, as the latter depend on the programming language(s) being parsed.
- Our AST will consits solely of nodes connected to each other
