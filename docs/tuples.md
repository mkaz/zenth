# Tuples

Tuples are fixed-size ordered values that can hold mixed types. They come in two forms: positional and named.

## Positional Tuples

Use `tuple(...)` to create a tuple with values accessed by numeric index:

```zenth
let t = tuple("hello", 42);
println(t.0);       // "hello"
println(t.1);       // "42"
```

With a type annotation:

```zenth
let pair: tuple(str, int) = tuple("age", 30);
```

## Named Tuples

Named tuples add field names so you can access elements by name. Use `=` for values and `:` for types:

```zenth
let person = tuple(name="Alice", age=30);
println(person.name);       // "Alice"
println(person.age);        // "30"
```

Numeric index access still works on named tuples:

```zenth
println(person.0);          // "Alice"
println(person.1);          // "30"
```

With a type annotation:

```zenth
let t: tuple(key: str, value: int) = tuple(key="x", value=42);
```

A tuple must be all-named or all-positional. Mixing is a compile error:

```zenth
let bad = tuple(name="foo", 42);  // error: cannot mix named and positional
```

Named and positional tuples are distinct types — you cannot assign one to the other:

```zenth
let t: tuple(name: str, age: int) = tuple("Alice", 30);  // error: type mismatch
```

## Function Parameters and Return Types

Tuples work in function signatures like any other type. For multi-return functions, use the `(T1, T2)` shorthand return type:

```zenth
fn min_max(nums: array(int)) -> (int, int) {
    var lo = nums[0];
    var hi = nums[0];
    for n in nums {
        if n < lo { lo = n; }
        if n > hi { hi = n; }
    }
    return lo, hi;
}

fn main() {
    let (lo, hi) = min_max([3, 1, 4, 1, 5, 9]);
    println("{lo} {hi}");  // "1 9"
}
```

The `(T1, T2)` return type is shorthand for `tuple(T1, T2)`. The `return a, b;` syntax is shorthand for `return tuple(a, b);`.

Named tuples also work in signatures:

```zenth
fn make_person(name: str, age: int) -> tuple(name: str, age: int) {
    return tuple(name=name, age=age);
}

fn greet(person: tuple(name: str, age: int)) {
    println(person.name + " is " + str(person.age));
}

fn main() {
    let p = make_person("Alice", 30);
    greet(p);
}
```

## Destructuring

Tuple destructuring binds each element to a variable by position:

```zenth
let t = tuple("hello", 42);
let (greeting, number) = t;
println(greeting);       // "hello"
println(number);         // "42"
```

This works with named tuples too — the names don't affect destructuring order:

```zenth
let person = tuple(name="Alice", age=30);
let (n, a) = person;
println(n);              // "Alice"
```

## Tuples as Hashmap and Set Keys

Tuples with comparable element types (built-in types like `int`, `str`, `f64`, `bool`) can be used as hashmap keys and set elements:

```zenth
var grid = hashmap(tuple(int, int), str);
grid[tuple(0, 0)] = "origin";
grid[tuple(1, 2)] = "point";
println(grid[tuple(0, 0)]);  // "origin"

var visited = set(tuple(int, int));
visited.add(tuple(3, 4));
println(visited.exists(tuple(3, 4)));  // true
```

See the [Hashmaps](hashmaps.md) and [Sets](sets.md) docs for more details.

## Arrays of Tuples

Tuples nest naturally inside arrays:

```zenth
let pairs: array(tuple(str, str)) = [];
pairs.add(tuple("left", "right"));
println(pairs[0].0 + ":" + pairs[0].1);
```

With named tuples:

```zenth
let people: array(tuple(name: str, score: int)) = [];
people.add(tuple(name="Bob", score=95));
people.add(tuple(name="Carol", score=88));

for p in people {
    println(p.name + ": " + str(p.score));
}
```
