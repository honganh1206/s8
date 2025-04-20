# Compilers

Compilers come in all shapes in sizes, not jsut programming languages, including regexes, database queries, and even HTML templates

A compiler **transforms** computer code written in one programming language into another programming language

Compilers are **fundamentally** about translation

Interpreters and compilers are different when it comes to traversing the AST

Compilers generate source code in another language (target language), and the result would be executed by the computer

```js
+------------------+
|   Source Code    |
+--------+---------+
         |
      +--------+---------+
      | |  Lexer & Parser|
      +--------+---------+
        |
        v
+--------+---------+
|       AST        |
+--------+---------+
         |
      +--------+---------+
      |  | Optimizer     |
      +--------+---------+
        |
        v
+--------------------------+
| Internal Representation  |
+------------+-------------+
             |
      +------------+-------------+
      |      |  Code Generator  |
      +------------+-------------+
            |
            v
+------------+-------------+
|        Machine Code      |
+--------------------------+
```

The Internal Representation (IR) is better for optimization and translation into the target language than the AST

The optimization phase: Remove dead code, pre-calculate simple arithmetic, move not needed code away...

The code generator generates the code in the target language
