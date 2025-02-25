# Hashing Objects

Some precautions:

```go
// This will not work
type Hash struct {
  Pairs map[Object]Object
}
```

Supposed we do this with our programming language

```go
let hash = { "name": "Monkey" };
hash["name"]
```

We will end up with a `*object.String` with `.Value` being "name" mapped to a `*object.String` with `.Value` being "Monkey"

But the index operator will NOT work: The "name" string literal will be evaluated to a _brand new_ `*object.String`, and it will **point to a different memory location**, so the "Monkey" value will not be bound to it

Reason: There are two pointers pointing to the same value "name" but _they are at different locations_. Thus comparing these two pointers will tell us that they are not equal

```go
name1 := &object.String{Value: "name"}
monkey := &object.String{Value: "Monkey"}
pairs := map[object.Object]object.Object{}
pairs[name1] = monkey
fmt.Printf("pairs[name1]=%+v\n", pairs[name1])
// => pairs[name1]=&{Value:Monkey} - Pointer #1
name2 := &object.String{Value: "name"}
fmt.Printf("pairs[name2]=%+v\n", pairs[name2])
// => pairs[name2]=<nil> - Pointer #2
fmt.Printf("(name1 == name2)=%t\n", name1 == name2)
// => (name1 == name2)=false
```

What we need: We need to generate a hash key for a `*object.String` that is **equal** to the hash key of another `*object.String` with the same`.Value`. The same thing needs to be done for `*object.Integer` and `*object.Boolean`

> [!IMPORTANT] > **Between types, the hash keys always have to differ**

How do we solve this? By using the `HashKey` function hashes the string/integer/boolean key. The strings/integers/boolean values sharing the same content will also share the hash, and thus pointing to the same memory location.

```js
pairs[name1.HashKey()] = monkey
fmt.Printf("pairs[name1.HashKey()]=%+v\n", pairs[name1.HashKey()])
// => pairs[name1.HashKey()]=&{Value:Monkey}
name2 := &object.String{Value: "name"}
fmt.Printf("pairs[name2.HashKey()]=%+v\n", pairs[name2.HashKey()])
// => pairs[name2.HashKey()]=&{Value:Monkey}
fmt.Printf("(name1 == name2)=%t\n", name1 == name2)
// => (name1 == name2)=false
fmt.Printf("(name1.HashKey() == name2.HashKey())=%t\n",
name1.HashKey() == name2.HashKey())
// => (name1.HashKey() == name2.HashKey())=true
```
