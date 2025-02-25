# 4-5-Hashes

The structure

```js
{<expression> : <expression>, <expression> : <expression>, ... }
```

How the hashes will look like:

```js
>> let myHash = {"name": "Jimmy", "age": 72, "band": "Led Zeppelin"};
>> myHash["name"]
Jimmy
>> myHash["age"]
72
>> myHash["band"]
Led Zeppelin
```

We use the index operator from [[4-4-Arrays]] to get the values out of the hash.

Note that the index values can be either strings/integers/boolean values

```js
>> let myHash = {true: "yes, a boolean", 99: "correct, an integer"};
>> myHash[true]
yes, a boolean
>> myHash[99]
correct, an integer
>> myHash[5 > 1]
yes, a boolean
>> myHash[100 - 1]
correct, an integer
```

One small problem: Our hash keys can be only either string/integer/boolean. How do we tell the parser that? Answer: _Validate hash key types during evaluation and generate possible errors there_

Why doing so? So we can prevent this:

```js
// key is valid as a hash key, even though it is an identifier
let key = "name";
let hash = { key: "Monkey" };
```

Since the `token.LBRACE` of a hash literal is in prefix position, so we make the parsing function for hash literals a `prefixParseFn`

[[Hashing Objects]]

## Why `Hash.Pairs` is of type `map[HashKey]HashPair` and not `map[HashKey]Object` ?

When printing a hash into the REPL, we want to _print the values contained in the hash as well as the keys_

We need to keep track of the _objects that generated the `HashKeys`_ by using `HashPairs` as values
