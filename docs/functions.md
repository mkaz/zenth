# Functions

Functions are declared with the `fn` keyword. Every program needs a `main` function as its entry point.

## Declaring Functions

```zenth
fn greet(name: str) {
    println("Hello, {name}!");
}
```

Parameters require explicit type annotations using the `name: Type` syntax.

## Return Types

Specify a return type with `->`:

```zenth
fn add(a: int, b: int) -> int {
    return a + b;
}
```

Functions without a `->` return type return nothing (void).

## Calling Functions

```zenth
fn main() {
    greet("World");
    let result = add(3, 4);
    println(str(result));  // 7
}
```

## Multiple Parameters

```zenth
fn clamp(value: int, low: int, high: int) -> int {
    if value < low {
        return low;
    }
    if value > high {
        return high;
    }
    return value;
}
```

## Recursion

Functions can call themselves:

```zenth
fn factorial(n: int) -> int {
    if n <= 1 {
        return 1;
    }
    return n * factorial(n - 1);
}
```

## Built-in Functions

Zenth provides several built-in functions that are always available:

| Function | Description |
|----------|-------------|
| `print(...)` | Print values without a newline |
| `println(...)` | Print values followed by a newline |
| `len(x)` | Return the length of a string or slice |
| `str(x)` | Convert any value to its string representation |

```zenth
fn main() {
    println("hello");      // prints "hello\n"
    print("no newline");   // prints without newline
    println(str(len("abc")));  // prints "3"
    println(str(42));      // prints "42"
}
```

## Methods

Functions can be defined inside structs to act as methods. See [Structs](structs.md) for details.
