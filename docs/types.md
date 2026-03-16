# Types

Zenth is strongly typed with type inference for local variables. Function signatures always require explicit types.

## Integer Types

| Type  | Size    | Range |
|-------|---------|-------|
| `int` | platform-sized | Default integer type |
| `i8`  | 8-bit signed  | -128 to 127 |
| `i16` | 16-bit signed | -32,768 to 32,767 |
| `i32` | 32-bit signed | -2^31 to 2^31-1 |
| `i64` | 64-bit signed | -2^63 to 2^63-1 |
| `u8`  | 8-bit unsigned  | 0 to 255 |
| `u16` | 16-bit unsigned | 0 to 65,535 |
| `u32` | 32-bit unsigned | 0 to 2^32-1 |
| `u64` | 64-bit unsigned | 0 to 2^64-1 |

```zenth
let x = 42;          // int (inferred)
let big: i64 = 999;
```

Integer literals can use underscores for readability:

```zenth
let million = 1_000_000;
```

### Built-in Integer Constants

Zenth provides built-in constants for the integer range limits:

| Constant  | Value                | Description |
|-----------|----------------------|-------------|
| `INT_MAX` | 9223372036854775807  | Maximum `int` value (2^63-1) |
| `INT_MIN` | -9223372036854775808 | Minimum `int` value (-2^63) |

These are true constants and cannot be reassigned:

```zenth
let big = INT_MAX;
let small = INT_MIN;

// Use in expressions
if score < INT_MAX {
    println("not at the limit yet");
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
| `f32` | 32-bit | Single precision |
| `f64` | 64-bit | Double precision (default for float literals) |

```zenth
let pi = 3.14159;         // f64 (inferred)
let ratio: f32 = 0.75;
```

## Boolean

```zenth
let active = true;
let done = false;
```

Booleans are used in `if` conditions and logical expressions.

## String

The `str` type holds text:

```zenth
let greeting = "Hello, World!";
```

See [Strings](strings.md) for interpolation and raw string details.

## Byte

The `byte` type represents a single byte (alias for `u8`):

```zenth
let b: byte = 65;
```

Indexing into a string returns a one-character `str`.

## Nil

`nil` represents the absence of a value for reference types (hashmaps, arrays, objects). It must be used with an explicit type annotation:

```zenth
var m: hashmap(str, int) = nil;
var items: array(int) = nil;
```

Primitive types (`int`, `str`, `bool`, `f64`) cannot be nil.

## Composite Types

### Range Objects

`range(start, end)` and `rangei(start, end)` return a **range object** — a lightweight value type that is iterable and supports O(1) containment checks.

```zenth
let r  = range(0, 10);    // exclusive: [0, 9]
let ri = rangei(0, 10);   // inclusive: [0, 10]
let rs = range(0, 20, 3); // with step: [0, 3, 6, 9, 12, 15, 18]
```

**Methods:**

| Method | Returns | Description |
|--------|---------|-------------|
| `.contains(x)` | `bool` | O(1) bounds check: is `x` within the range? |

**Built-ins that work on range objects:**

| Built-in | Description |
|----------|-------------|
| `len(r)` | Number of elements in the range (O(1)) |

**Iteration:**

Range objects are iterable in `for-in` loops:

```zenth
for i in range(0, 5) {
    println(i);   // 0 1 2 3 4
}

for i, v in rangei(10, 13) {
    println("{i}: {v}");
}
// 0: 10
// 1: 11
// 2: 12
// 3: 13
```

**Containment checks:**

```zenth
let r = range(1, 10);
println(r.contains(5));   // true
println(r.contains(10));  // false (exclusive end)

let ri = rangei(1, 10);
println(ri.contains(10)); // true (inclusive end)
```

**Note:** `.contains()` is a bounds check, not a sequence membership test. For `range(0, 10, 3)` (sequence `[0, 3, 6, 9]`), `.contains(5)` returns `true` because 5 is within the bounds `[0, 10)`.

### Arrays

Dynamic arrays use `array(Type)` syntax:

```zenth
let numbers: array(int) = [1, 2, 3, 4, 5];
```

See [Arrays](arrays.md) for more details.

### Hashmaps

Hashmaps are key/value collections. Create an empty hashmap with `hashmap(KeyType, ValueType)`:

```zenth
obj Point {
    x: int;
    y: int;
}

var grid = hashmap(Point, str);
let pt = Point(x=1, y=2);
grid[pt] = "#";
```

For object keys, hashmap lookup is value-based: another `Point(x=1, y=2)` resolves the same entry.

### Sets

Sets are unordered collections of unique values. Create an empty set with `set(Type)`:

```zenth
var visited = set(str);
visited.add("start");
visited.add("middle");
visited.exists("start");  // true
visited.remove("middle");
println(len(visited));       // 1
```

Sets support `int`, `str`, `bool`, and other comparable types as elements.

#### Sets of Tuples

Sets can hold tuples, enabling composite keys without string round-tripping:

```zenth
var dots = set(tuple(int, int));
dots.add(tuple(6, 10));
dots.add(tuple(0, 14));
dots.add(tuple(6, 10));  // duplicate, ignored
println(len(dots));       // 2

if dots.exists(tuple(6, 10)) {
    println("found");
}

for dot in dots {
    println("{dot.0},{dot.1}");
}
```

Tuple elements must be comparable types (scalars, strings, booleans). Tuples containing arrays or maps cannot be used as set elements.

### Tuples

Tuples are fixed-size ordered values that can hold mixed types. They support both positional and named fields:

```zenth
let t = tuple("a", 1);
println(t.0);      // positional access

let person = tuple(name="Alice", age=30);
println(person.name);  // named access
```

See [Tuples](tuples.md) for named tuples, destructuring, and usage in functions.

### Objects

User-defined types with named fields:

```zenth
obj Point {
    x: f64;
    y: f64;
}

let p = Point(x=1.0, y=2.0);
```

See [Objects](objects.md) for methods and `self`.

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

Use the `str()` built-in to convert any value to a string:

```zenth
let n = 42;
let s = str(n);    // "42"
println(str(3.14)); // "3.14"
```

Use `int()` to convert strings or floats to integers. With two arguments, the second specifies the base:

```zenth
let a = int("42");         // 42
let b = int(3.14);         // 3
let bin = int("1010", 2);  // 10 (binary)
let hex = int("ff", 16);   // 255 (hexadecimal)
let oct = int("77", 8);    // 63 (octal)
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

Arithmetic operators require matching numeric types on both sides. The `+` operator also works for string concatenation when both sides are `str`.
