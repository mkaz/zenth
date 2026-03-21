# Types

Zenth is strongly typed with type inference for local variables. Function signatures always require explicit types.

## Integer Types

Zenth has a single integer type: `Int` (64-bit signed).

```zenth
let x = 42;          // Int (inferred)
let big: Int = 999;
```

Integer literals can use underscores for readability:

```zenth
let million = 1_000_000;
```

### Built-in Integer Constants

Zenth provides built-in constants for the integer range limits:

| Constant  | Value                | Description |
|-----------|----------------------|-------------|
| `INT_MAX` | 9223372036854775807  | Maximum `Int` value (2^63-1) |
| `INT_MIN` | -9223372036854775808 | Minimum `Int` value (-2^63) |

These are true constants and cannot be reassigned:

```zenth
let big = INT_MAX;
let small = INT_MIN;

// Use in expressions
if score < INT_MAX {
    Println("not at the limit yet");
}
```

### Base Conversion

Use `.to_base(base)` to convert an integer to a string in a given base (2-36):

```zenth
let hex = 255.to_base(16);   // "ff"
let bin = 10.to_base(2);     // "1010"
let oct = 63.to_base(8);     // "77"
```

## Floating-Point Types

| Type  | Size | Description |
|-------|------|-------------|
| `Float` | 64-bit | Floating-point type (default for float literals) |

```zenth
let pi = 3.14159;         // Float (inferred)
let ratio: Float = 0.75;
```

Zenth has a single floating-point type: `Float`. Legacy names like `f32` and `f64` are not valid.

```zenth
let a: Float = 1.5;
let b = Float("2.25");
let parts = ["3.5", "4.5"];
let nums = parts.to_float();
```

## Boolean

```zenth
let active = true;
let done = false;
```

Booleans are used in `if` conditions and logical expressions.

## String

The `Str` type holds text:

```zenth
let greeting = "Hello, World!";
```

See [Strings](strings.md) for interpolation and raw string details.

## Byte

The `Byte` type represents a single byte:

```zenth
let b: Byte = 65;
```

Indexing into a string returns a one-character `Str`.

## Nil

`nil` represents the absence of a value for reference types (hashmaps, arrays, objects). It must be used with an explicit type annotation:

```zenth
var m: Hashmap(Str, Int) = nil;
var items: Array(Int) = nil;
```

Primitive types (`Int`, `Str`, `Bool`, `Float`) cannot be nil.

## Composite Types

### Range Objects

`Range(start, end)` and `Rangei(start, end)` return a **range object** — a lightweight value type that is iterable and supports O(1) containment checks.

```zenth
let r  = Range(0, 10);    // exclusive: [0, 9]
let ri = Rangei(0, 10);   // inclusive: [0, 10]
let rs = Range(0, 20, 3); // with step: [0, 3, 6, 9, 12, 15, 18]
```

**Methods:**

| Method | Returns | Description |
|--------|---------|-------------|
| `.contains(x)` | `Bool` | O(1) bounds check: is `x` within the range? |

**Built-ins that work on range objects:**

| Built-in | Description |
|----------|-------------|
| `Len(r)` | Number of elements in the Range (O(1)) |

**Iteration:**

Range objects are iterable in `for-in` loops:

```zenth
for i in Range(0, 5) {
    Println(i);   // 0 1 2 3 4
}

for i, v in Rangei(10, 13) {
    Println("{i}: {v}");
}
// 0: 10
// 1: 11
// 2: 12
// 3: 13
```

**Containment checks:**

```zenth
let r = Range(1, 10);
Println(r.contains(5));   // true
Println(r.contains(10));  // false (exclusive end)

let ri = Rangei(1, 10);
Println(ri.contains(10)); // true (inclusive end)
```

**Note:** `.contains()` is a bounds check, not a sequence membership test. For `Range(0, 10, 3)` (sequence `[0, 3, 6, 9]`), `.contains(5)` returns `true` because 5 is within the bounds `[0, 10)`.

### Arrays

Dynamic arrays use `Array(Type)` syntax:

```zenth
let numbers: Array(Int) = [1, 2, 3, 4, 5];
```

See [Arrays](arrays.md) for more details.

### Hashmaps

Hashmaps are key/value collections. Create an empty hashmap with `Hashmap(KeyType, ValueType)`:

```zenth
obj Point {
    x: Int;
    y: Int;
}

var grid = Hashmap(Point, Str);
let pt = Point(x=1, y=2);
grid[pt] = "#";
```

For object keys, hashmap lookup is value-based: another `Point(x=1, y=2)` resolves the same entry.

### Sets

Sets are unordered collections of unique values. Create an empty set with `Set(Type)`:

```zenth
var visited = Set(Str);
visited.add("start");
visited.add("middle");
visited.exists("start");  // true
visited.remove("middle");
Println(Len(visited));       // 1
```

Sets support `Int`, `Str`, `Bool`, and other comparable types as elements.

#### Sets of Tuples

Sets can hold tuples, enabling composite keys without string round-tripping:

```zenth
var dots = Set(Tuple(Int, Int));
dots.add(Tuple(6, 10));
dots.add(Tuple(0, 14));
dots.add(Tuple(6, 10));  // duplicate, ignored
Println(Len(dots));       // 2

if dots.exists(Tuple(6, 10)) {
    Println("found");
}

for dot in dots {
    Println("{dot.0},{dot.1}");
}
```

Tuple elements must be comparable types (scalars, strings, booleans). Tuples containing arrays or maps cannot be used as set elements.

### Tuples

Tuples are fixed-size ordered values that can hold mixed types. They support both positional and named fields:

```zenth
let t = Tuple("a", 1);
Println(t.0);      // positional access

let person = Tuple(name="Alice", age=30);
Println(person.name);  // named access
```

See [Tuples](tuples.md) for named tuples, destructuring, and usage in functions.

### Objects

User-defined types with named fields:

```zenth
obj Point {
    x: Float;
    y: Float;
}

let p = Point(x=1.0, y=2.0);
```

See [Objects](objects.md) for methods and `self`.

### Function Types

Function types describe the signature of a function value. Use them in parameter lists, variable declarations, and type aliases:

```zenth
Fn(Int) -> Int           // takes Int, returns Int
Fn(Str, Int) -> Bool     // takes Str and Int, returns Bool
Fn(Int)                  // takes Int, returns nothing
```

See [Functions](functions.md#function-types) for examples of higher-order functions and passing functions as values.

### Enums

Enums define a type with a fixed set of named variants:

```zenth
enum Direction {
    North;
    South;
    East;
    West;
}

let d = Direction.North;
```

See [Enums](enums.md) for explicit values, comparison, and match usage.

## Type Conversions

Use the `Str()` built-in to convert any value to a string:

```zenth
let n = 42;
let s = Str(n);    // "42"
Println(Str(3.14)); // "3.14"
```

Use `Int()` to convert strings or floats to integers. With two arguments, the second specifies the base:

```zenth
let a = Int("42");         // 42
let b = Int(3.14);         // 3
let bin = Int("1010", 2);  // 10 (binary)
let hex = Int("ff", 16);   // 255 (hexadecimal)
let oct = Int("77", 8);    // 63 (octal)
```

This is equivalent to `"ff".to_int(16)` but reads more naturally as a conversion function.

## Operators

### Arithmetic

| Operator | Description |
|----------|-------------|
| `+` | Addition (also string concatenation) |
| `-` | Subtraction |
| `*` | Multiplication |
| `/` | Division |
| `%` | Modulo |

### Comparison

| Operator | Description |
|----------|-------------|
| `==` | Equal |
| `!=` | Not equal |
| `<`  | Less than |
| `>`  | Greater than |
| `<=` | Less than or equal |
| `>=` | Greater than or equal |

### Logical

| Operator | Description |
|----------|-------------|
| `&&` | Logical AND |
| `\|\|` | Logical OR |
| `!`  | Logical NOT |

Arithmetic operators require matching numeric types on both sides. The `+` operator also works for string concatenation when both sides are `Str`.
