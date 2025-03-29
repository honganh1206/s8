# Syntactic macro systems

Treat code as **data** like how we use AST to turn s8 source code into Go structs and modify it with Go

We can do the syntactic macros **on the language itself**, not just in the host language: If a language X has syntactic macros, we can use language X to work with source code written in X

Aka the language _becomes self-aware_

Example: Elixir's `quote` function - Stop code from being evaluated

```elixir
# Instead of evaluating, quote returns a data structure that represents this expression
iex(1)> quote do: 10 + 5
# Return: a TUPLE containing the operator, metadata and a list of operands
{:+, [context: Elixir, import: Kernel], [10, 5]}
```

We can access it like any other tuple:

```elixir
iex(2)> exp = quote do: 10 + 5
{:+, [context: Elixir, import: Kernel], [10, 5]}
iex(3)> elem(exp, 0)
:+
iex(4)> elem(exp, 2)
[10, 5]
```

When dealing with dynamically injected numbers into the AST, we need the `unquote` function - jumping out of quote context and evaluate the code

```elixir
iex(8)> quote do: 10 + 5 + unquote(my_number)
{:+, [context: Elixir, import: Kernel],
[{:+, [context: Elixir, import: Kernel], [10, 5]}, 99]}
```

A simple macro in Elixir

```elixir
defmodule MacroExample do
  defmacro plus_to_minus(expression) do
    args = elem(expression, 2)
    quote do
      unquote(Enum.at(args, 0)) - unquote(Enum.at(args, 1))
    end
  end
end

# Usage
iex(1)> MacroExample.plus_to_minus 10 + 5
5
```
