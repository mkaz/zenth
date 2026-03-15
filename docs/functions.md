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
    println(result);       // 7
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

## Multi-Return Functions

Functions can return multiple values using a parenthesized return type. The caller unpacks the result with tuple destructuring:

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

The `return a, b;` syntax is shorthand for `return tuple(a, b);`. The return type `(T1, T2)` is shorthand for `tuple(T1, T2)`. Both forms are equivalent:

```zenth
fn divide(a: int, b: int) -> (int, int) {
    return a / b, a % b;  // quotient and remainder
}

fn swap(x: str, y: str) -> (str, str) {
    return y, x;
}
```

You can also use the full `tuple(...)` syntax when you want to be explicit:

```zenth
fn get_pair() -> tuple(int, int) {
    return tuple(10, 20);
}
```

## Default Arguments

Function parameters can specify default values using `=`. If a caller omits a default parameter, the default value is used:

```zenth
fn greet(name: str = "World") {
    println("Hello, {name}!");
}

fn main() {
    greet();          // Hello, World!
    greet("Zenth");   // Hello, Zenth!
}
```

## Named Arguments

When calling a function, you can provide arguments by name, which allows you to pass them in any order. This is especially useful for functions with many parameters or default values:

```zenth
fn draw_rect(x: int, y: int, width: int, height: int = 10) {
    // ...
}

fn main() {
    draw_rect(width=20, x=5, y=5); // y=5, x=5, width=20, height=10
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
| `print(x, [enabled])` | Print `x` without a newline; print only when `enabled` is `true` (default) |
| `println(x, [enabled])` | Print `x` with a newline; print only when `enabled` is `true` (default) |
| `len(x)` | Return the length of a string or slice |
| `str(x)` | Convert any value to its string representation |
| `abs(x)` | Absolute value for `int` or `f64` |
| `min(a, b, ...)` / `max(a, b, ...)` | Minimum / maximum for two or more `int` or `f64` values (all args must be the same type) |
| `clamp(x, lo, hi)` | Clamp `x` between `lo` and `hi` (`int` or `f64`) |
| `round(x)` / `floor(x)` / `ceil(x)` | Floating-point rounding helpers (return `f64`) |
| `pow(x, y)` / `sqrt(x)` | Power and square root (return `f64`) |
| `flag(default=val)` | Declare a command-line flag with a default value |

```zenth
fn main() {
    println("hello");      // prints "hello\n"
    print("no newline");   // prints without newline
    let debug = false;
    println("debug line", debug); // prints only when debug is true
    println(len("abc"));       // prints "3"
    println(str(42));      // prints "42"
    println(str(abs(-5)));  // prints "5"
}
```

## Command-Line Flags

The `flag()` built-in declares command-line flags. The flag name is inferred from the variable name, and the type is inferred from the default value (`bool`, `int`, or `str`).

```zenth
fn main() {
    let debug = flag(default=false);   // --debug
    let times = flag(default=5);       // --times 10
    let msg = flag(default="Hello");   // --msg "world"

    for var i = 0; i < times; i++ {
        println(msg, !debug);
    }
}
```

Build and run with flags:

```sh
zenth build -o greet greet.zn
./greet --debug --times 3 --msg "Hi"
```

Or pass flags directly with `zenth run`:

```sh
zenth run greet.zn --times 3 --msg "Hi"
```

When no flags are provided, the default values are used. Boolean flags are set to `true` by passing `--name` with no value.

## Closures

Closures are anonymous functions declared with `fn` in expression position. They are used with slice methods like `.map()` and `.filter()`:

```zenth
let doubled = [1, 2, 3].map(fn(x) x * 2);
let evens = [1, 2, 3, 4].filter(fn(x) x % 2 == 0);
```

Parameter types are inferred from context when used with `.map()` or `.filter()`. You can also provide explicit types:

```zenth
fn(x: int) -> int x * 2
```

Block body closures use explicit `return`:

```zenth
let processed = nums.map(fn(x: int) -> int {
    let y = x + 10;
    return y;
});
```

## Methods

Functions can be defined inside objects to act as methods. See [Objects](objects.md) for details.
