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

`nil` represents the absence of a value:

```zenth
let nothing = nil;
```

## Composite Types

### Slices

Dynamic arrays use `[]Type` syntax:

```zenth
let numbers: []int = [1, 2, 3, 4, 5];
```

See [Slices](slices.md) for more details.

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

### Structs

User-defined types with named fields:

```zenth
struct Point {
    x: f64;
    y: f64;
}

let p = Point{ x: 1.0, y: 2.0 };
```

See [Objects](objects.md) for methods and `self`.

## Type Conversions

Use the `str()` built-in to convert any value to a string:

```zenth
let n = 42;
let s = str(n);    // "42"
println(str(3.14)); // "3.14"
```

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
