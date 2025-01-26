# 3-9-Binding and the Environment

We add evaluation for `let` statements and bind the values to the identifiers

For example, we must ensure that in `let x = 5 * 5` the identifier `x` evaluates to 25 after interpreting that line

We evaluate `let` statement by _evaluating their value-producing expression_ and _keep track of the produced value under the specified name_

We evaluate the identifiers by _checking if we already have a value bound to the name_

## The environment

The environment is what we use to _keep track of value by associating them with a name_

It is **a hash map that associates strings with objects**

The environment will **persist** between calls to `Eval()`, as _we do not want the bindings to be in a different environment_ (thus being a pointer). For testing, however, we should create a new environment for each call to `testEval()` helper function
