# Enums

Enums define a type with a fixed set of named variants. Each variant has an integer value, starting from 0 by default.

## Declaring Enums

```zenth
enum Color {
    Red;
    Green;
    Blue;
}
```

Variants are automatically assigned values starting from 0: `Red` is 0, `Green` is 1, `Blue` is 2.

## Explicit Values

You can assign explicit integer values to variants. Subsequent variants auto-increment from the last explicit value:

```zenth
enum HttpStatus {
    Ok = 200;
    NotFound = 404;
    InternalError = 500;
}

enum Priority {
    Low = 1;
    Medium;    // 2
    High;      // 3
    Critical;  // 4
}
```

## Using Enums

Access variants with `EnumName.Variant` syntax:

```zenth
let c = Color.Red;
let status = HttpStatus.NotFound;
```

## Type Annotations

Enums can be used as type annotations:

```zenth
let c: Color = Color.Green;

fn describe(c: Color) {
    match c {
        Color.Red => println("red");
        Color.Green => println("green");
        Color.Blue => println("blue");
        _ => println("unknown");
    }
}
```

## Comparison

Enum values support equality comparison with `==` and `!=`:

```zenth
let c = Color.Red;

if c == Color.Red {
    println("it's red!");
}

if c != Color.Blue {
    println("not blue");
}
```

## Match

Enums work naturally with `match`:

```zenth
fn color_name(c: Color) -> str {
    return match c {
        Color.Red => "red",
        Color.Green => "green",
        Color.Blue => "blue",
        _ => "unknown",
    };
}
```

## String Conversion

Use `str()` to convert an enum value to its variant name:

```zenth
let c = Color.Green;
println(str(c));  // prints "Green"
```
