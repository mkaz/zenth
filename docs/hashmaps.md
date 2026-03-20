# Hashmaps

Hashmaps are unordered collections of key-value pairs.

## Creating Hashmaps

You can create a hashmap using the `Hashmap(KeyType, ValueType)` built-in:

```zenth
var scores = Hashmap(str, int);
```

## Assigning and Accessing

Use bracket syntax `[key]` to assign and access values:

```zenth
scores["Alice"] = 100;
scores["Bob"] = 95;

let aliceScore = scores["Alice"];
Println(aliceScore);
```

## Hashmap Length

Use `Len()` to get the number of key-value pairs in the hashmap:

```zenth
Println(Len(scores));
```

## Default Values

You can provide a default value for missing keys using the `default` named argument:

```zenth
var counts = Hashmap(str, int, default=0);
counts["apples"] += 1;
counts["apples"] += 2;
Println(counts["apples"]);       // prints "3"
Println(counts["missing"]);      // prints "0"
```

When you access a missing key, the default is returned instead of Go's zero value. Compound assignments (`+=`, `-=`, `++`, etc.) on missing keys initialize from the default first.

## Iterating

Use `for` to iterate over a hashmap:

```zenth
var m = Hashmap(str, int);
m["a"] = 1;
m["b"] = 2;

// iterate over key-value pairs
for k, v in m {
    Println(k + "=" + Str(v));
}

// iterate over keys only
for k in m {
    Println(k);
}
```

### exists()

Use `.exists(key)` to check whether a key is present in the hashmap:

```zenth
var scores = Hashmap(str, int);
scores["Alice"] = 95;

if scores.exists("Alice") {
    Println("found Alice");
}
if !scores.exists("Bob") {
    Println("Bob not found");
}
```

### keys() and values()

Use `.keys()` and `.values()` to get arrays of keys or values:

```zenth
let keys = m.keys();
let vals = m.values();

for v in vals {
    Println(v);
}
```

## Object Keys

You can use user-defined objects as hashmap keys. Hashmap lookup is value-based, so two instances with the same field values resolve to the same entry:

```zenth
obj Point {
    x: int;
    y: int;
}

var grid = Hashmap(Point, str);
let pt1 = Point(x=1, y=2);
let pt2 = Point(x=1, y=2);

grid[pt1] = "#";
Println(grid[pt2]); // Prints "#"
```

## Tuple Keys

Tuples can also be used as hashmap keys, which is useful for multi-dimensional lookups or composite keys:

```zenth
var grid = Hashmap(Tuple(int, int), str);
grid[Tuple(0, 0)] = "origin";
grid[Tuple(1, 2)] = "point";

Println(grid[Tuple(0, 0)]);  // "origin"
Println(grid[Tuple(1, 2)]);  // "point"

// Check existence
if grid.exists(Tuple(0, 0)) {
    Println("found origin");
}

// Iterate
for key in grid {
    Println(Str(key.0) + "," + Str(key.1));
}
```

All tuple element types must be comparable (built-in types like `int`, `str`, `f64`, `bool`). Tuples containing arrays or other non-comparable types cannot be used as keys.
