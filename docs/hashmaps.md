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
println(str(aliceScore));
```

## Hashmap Length

Use `len()` to get the number of key-value pairs in the hashmap:

```zenth
println(str(len(scores)));
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
