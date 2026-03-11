# Hashmaps

Hashmaps are unordered collections of key-value pairs.

## Creating Hashmaps

You can create a hashmap using the `hashmap(KeyType, ValueType)` built-in:

```zenth
var scores = hashmap(str, int);
```

## Assigning and Accessing

Use bracket syntax `[key]` to assign and access values:

```zenth
scores["Alice"] = 100;
scores["Bob"] = 95;

let aliceScore = scores["Alice"];
println(aliceScore);
```

## Hashmap Length

Use `len()` to get the number of key-value pairs in the hashmap:

```zenth
println(len(scores));
```

## Default Values

You can provide a default value for missing keys using the `default` named argument:

```zenth
var counts = hashmap(str, int, default=0);
counts["apples"] += 1;
counts["apples"] += 2;
println(counts["apples"]);       // prints "3"
println(counts["missing"]);      // prints "0"
```

When you access a missing key, the default is returned instead of Go's zero value. Compound assignments (`+=`, `-=`, `++`, etc.) on missing keys initialize from the default first.

## Iterating

Use `for` to iterate over a hashmap:

```zenth
var m = hashmap(str, int);
m["a"] = 1;
m["b"] = 2;

// iterate over key-value pairs
for k, v in m {
    println(k + "=" + str(v));
}

// iterate over keys only
for k in m {
    println(k);
}
```

### exists()

Use `.exists(key)` to check whether a key is present in the hashmap:

```zenth
var scores = hashmap(str, int);
scores["Alice"] = 95;

if scores.exists("Alice") {
    println("found Alice");
}
if !scores.exists("Bob") {
    println("Bob not found");
}
```

### keys() and values()

Use `.keys()` and `.values()` to get slices of keys or values:

```zenth
let keys = m.keys();
let vals = m.values();

for v in vals {
    println(v);
}
```

## Object Keys

You can use user-defined objects as hashmap keys. Hashmap lookup is value-based, so two instances with the same field values resolve to the same entry:

```zenth
obj Point {
    x: int;
    y: int;
}

var grid = hashmap(Point, str);
let pt1 = Point(x=1, y=2);
let pt2 = Point(x=1, y=2);

grid[pt1] = "#";
println(grid[pt2]); // Prints "#"
```
