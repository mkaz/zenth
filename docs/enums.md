# Enums

Enums define a type with a fixed set of named variants. Each variant is backed by a string value — by default the variant name itself.

## Declaring Enums

```zenth
enum Color {
    Red;
    Green;
    Blue;
}
```

Without explicit values, each variant's string value is its name: `"Red"`, `"Green"`, `"Blue"`.

## Explicit String Values

You can assign explicit string values to variants:

```zenth
enum HexColor {
    Red = "#FF0000";
    Green = "#00FF00";
    Blue = "#0000FF";
}

enum Stoplight {
    Red = "red";
    Yellow = "yellow";
    Green = "green";
}
```

## Using Enums

Access variants with `EnumName.Variant` syntax:

```zenth
let c = Color.Red;
let hex = HexColor.Red;
```

## Type Safety

Enums enforce that a variable can only hold one of the declared variants. Use type annotations to constrain values:

```zenth
let c: Color = Color.Green;

fn describe(c: Color) {
    match c {
        Color.Red => Println("red");
        Color.Green => Println("green");
        Color.Blue => Println("blue");
        _ => Println("unknown");
    }
}
```

A `Stoplight` variable can only be `Stoplight.Red`, `Stoplight.Yellow`, or `Stoplight.Green` — no arbitrary strings allowed.

## Comparison

Enum values support equality comparison with `==` and `!=`:

```zenth
let c = Color.Red;

if c == Color.Red {
    Println("it's red!");
}

if c != Color.Blue {
    Println("not blue");
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

## Printing

Enum values print as their string value:

```zenth
let c = Color.Green;
Println(c);  // prints "Green"

let h = HexColor.Red;
Println(h);  // prints "#FF0000"
```
