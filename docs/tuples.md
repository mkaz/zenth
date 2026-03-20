# Tuples

Tuples are fixed-size ordered values that can hold mixed types. They come in two forms: positional and named.

## Positional Tuples

Use `Tuple(...)` to create a tuple with values accessed by numeric index:

```zenth
let t = Tuple("hello", 42);
Println(t.0);       // "hello"
Println(t.1);       // "42"
```

With a type annotation:

```zenth
let pair: Tuple(Str,Int) = Tuple("age", 30);
```

## Named Tuples

Named tuples add field names so you can access elements by name. Use `=` for values and `:` for types:

```zenth
let person = Tuple(name="Alice", age=30);
Println(person.name);       // "Alice"
Println(person.age);        // "30"
```

Numeric index access still works on named tuples:

```zenth
Println(person.0);          // "Alice"
Println(person.1);          // "30"
```

With a type annotation:

```zenth
let t: Tuple(key: Str, value: Int) = Tuple(key="x", value=42);
```

A tuple must be all-named or all-positional. Mixing is a compile error:

```zenth
let bad = Tuple(name="foo", 42);  // error: cannot mix named and positional
```

Named and positional tuples are distinct types — you cannot assign one to the other:

```zenth
let t: Tuple(name: Str, age: Int) = Tuple("Alice", 30);  // error: type mismatch
```

## Function Parameters and Return Types

Tuples work in function signatures like any other type. For multi-return functions, use the `(T1, T2)` shorthand return type:

```zenth
fn min_max(nums: Array(Int)) -> (Int,Int) {
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
    Println("{lo} {hi}");  // "1 9"
}
```

The `(T1, T2)` return type is shorthand for `Tuple(T1, T2)`. The `return a, b;` syntax is shorthand for `return Tuple(a, b);`.

Named tuples also work in signatures:

```zenth
fn make_person(name: Str, age: Int) -> Tuple(name: Str, age: Int) {
    return Tuple(name=name, age=age);
}

fn greet(person: Tuple(name: Str, age: Int)) {
    Println(person.name + " is " + Str(person.age));
}

fn main() {
    let p = make_person("Alice", 30);
    greet(p);
}
```

## Destructuring

Tuple destructuring binds each element to a variable by position:

```zenth
let t = Tuple("hello", 42);
let (greeting, number) = t;
Println(greeting);       // "hello"
Println(number);         // "42"
```

This works with named tuples too — the names don't affect destructuring order:

```zenth
let person = Tuple(name="Alice", age=30);
let (n, a) = person;
Println(n);              // "Alice"
```

## Tuples as Hashmap and Set Keys

Tuples with comparable element types (built-in types like `Int`, `Str`, `Float`, `Bool`) can be used as hashmap keys and set elements:

```zenth
var grid = Hashmap(Tuple(Int,Int),Str);
grid[Tuple(0, 0)] = "origin";
grid[Tuple(1, 2)] = "point";
Println(grid[Tuple(0, 0)]);  // "origin"

var visited = Set(Tuple(Int,Int));
visited.add(Tuple(3, 4));
Println(visited.exists(Tuple(3, 4)));  // true
```

See the [Hashmaps](hashmaps.md) and [Sets](sets.md) docs for more details.

## Arrays of Tuples

Tuples nest naturally inside arrays:

```zenth
let pairs: Array(Tuple(Str,Str)) = [];
pairs.add(Tuple("left", "right"));
Println(pairs[0].0 + ":" + pairs[0].1);
```

With named tuples:

```zenth
let people: Array(Tuple(name: Str, score: Int)) = [];
people.add(Tuple(name="Bob", score=95));
people.add(Tuple(name="Carol", score=88));

for p in people {
    Println(p.name + ": " + Str(p.score));
}
```
