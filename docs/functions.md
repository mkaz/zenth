# Functions

Functions are declared with the `fn` keyword. Every program needs a `main` function as its entry point.

## Declaring Functions

```zenth
fn greet(name: Str) {
    Println("Hello, {name}!");
}
```

Parameters require explicit type annotations using the `name: Type` syntax.

## Return Types

Specify a return type with `->`:

```zenth
fn add(a: Int, b: Int) -> Int {
    return a + b;
}
```

Functions without a `->` return type return nothing (void).

## Calling Functions

```zenth
fn main() {
    greet("World");
    let result = add(3, 4);
    Println(result);       // 7
}
```

## Multiple Parameters

```zenth
fn clamp(value: Int, low: Int, high: Int) -> Int {
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

The `return a, b;` syntax is shorthand for `return Tuple(a, b);`. The return type `(T1, T2)` is shorthand for `Tuple(T1, T2)`. Both forms are equivalent:

```zenth
fn divide(a: Int, b: Int) -> (Int,Int) {
    return a / b, a % b;  // quotient and remainder
}

fn swap(x: Str, y: Str) -> (Str,Str) {
    return y, x;
}
```

You can also use the full `Tuple(...)` syntax when you want to be explicit:

```zenth
fn get_pair() -> Tuple(Int,Int) {
    return Tuple(10, 20);
}
```

## Default Arguments

Function parameters can specify default values using `=`. If a caller omits a default parameter, the default value is used:

```zenth
fn greet(name: Str = "World") {
    Println("Hello, {name}!");
}

fn main() {
    greet();          // Hello, World!
    greet("Zenth");   // Hello, Zenth!
}
```

## Named Arguments

When calling a function, you can provide arguments by name, which allows you to pass them in any order. This is especially useful for functions with many parameters or default values:

```zenth
fn draw_rect(x: Int, y: Int, width: Int, height: Int = 10) {
    // ...
}

fn main() {
    draw_rect(width=20, x=5, y=5); // y=5, x=5, width=20, height=10
}
```

## Recursion

Functions can call themselves:

```zenth
fn factorial(n: Int) -> Int {
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
| `Print(x, [enabled])` | Print `x` without a newline; print only when `enabled` is `true` (default) |
| `Println(x, [enabled])` | Print `x` with a newline; print only when `enabled` is `true` (default) |
| `Len(x)` | Return the length of a string or array |
| `Str(x)` | Convert any value to its string representation |
| `Int(x)` / `Int(Str, base)` | Convert to `Int` (optional base for string input) |
| `Float(x)` | Convert to `Float` |
| `Abs(x)` | Absolute value for `Int` or `Float` |
| `Min(a, b, ...)` / `Max(a, b, ...)` | Minimum / maximum for two or more `Int` or `Float` values (all args must be the same type) |
| `Clamp(x, lo, hi)` | Clamp `x` between `lo` and `hi` (`Int` or `Float`) |
| `Round(x)` / `Floor(x)` / `Ceil(x)` | Floating-point rounding helpers (return `Float`) |
| `Pow(x, y)` / `Sqrt(x)` | Power and square root (return `Float`) |
| `Range(start, end[, step])` | Build an exclusive range object |
| `Rangei(start, end[, step])` | Build an inclusive range object |
| `Hashmap(K, V[, default=val])` | Build an empty hashmap (optional default value) |
| `Set(T)` | Build an empty set |
| `Tuple(a, b, ...)` | Build a tuple value |
| `Zip(a, b, ...)` | Combine two or more arrays into an array of tuples |
| `Env(name)` / `Env(name, default)` | Read an environment variable (returns `Str`; default used when unset) |
| `Args.flag(default=val)` | Declare a command-line flag with a default value |
| `Args.args()` | Get remaining positional arguments as `Array(Str)` |
| `Date.today()` | Get today's date as a `Date` object |
| `Date.from(str[, fmt])` | Parse date from string (default format: `%Y-%m-%d`) |
| `File(path)` | Build a file value for file methods |
| `Exit([code])` | Exit the program (`0` when omitted) |
| `Assert(cond)` / `AssertEq(got, expected)` | Built-in test/assertion helpers |

### Conversion Examples

```zenth
fn main() {
    let whole = Int("42");
    let hex = Int("ff", 16);
    let pi = Float("3.14159");
    let ratio = Float(7) / Float(2);

    Println("whole=" + Str(whole));
    Println("hex=" + Str(hex));
    Println("pi=" + Str(pi));
    Println("ratio=" + Str(ratio));
}
```

```zenth
fn main() {
    Println("hello");      // prints "hello\n"
    Print("no newline");   // prints without newline
    let debug = false;
    Println("debug line", debug); // prints only when debug is true
    Println(Len("abc"));       // prints "3"
    Println(Str(42));      // prints "42"
    Println(Str(Abs(-5)));  // prints "5"
}
```

## Zip

The `Zip()` function combines two or more arrays into an array of tuples, pairing elements at corresponding positions:

```zenth
let names = ["Alice", "Bob", "Charlie"];
let scores = [95, 87, 92];

for pair in Zip(names, scores) {
    Println(pair.0 + ": " + Str(pair.1));
}
// Alice: 95
// Bob: 87
// Charlie: 92
```

With three or more arrays:

```zenth
let x = [1, 2, 3];
let y = [4, 5, 6];
let z = [7, 8, 9];
for t in Zip(x, y, z) {
    Println(Str(t.0) + "," + Str(t.1) + "," + Str(t.2));
}
```

When arrays have different lengths, `Zip()` stops at the shortest:

```zenth
let a = [1, 2, 3, 4];
let b = [10, 20];
let zipped = Zip(a, b);  // [Tuple(1, 10), Tuple(2, 20)]
```

## Command-Line Arguments

The `Args` object provides access to command-line flags and positional arguments.

### Flags

`Args.flag(default=val)` declares a command-line flag. The flag name is inferred from the variable name, and the type is inferred from the default value (`Bool`, `Int`, or `Str`).

```zenth
fn main() {
    let debug = Args.flag(default=false);   // --debug
    let times = Args.flag(default=5);       // --times 10
    let msg = Args.flag(default="Hello");   // --msg "world"

    for var i = 0; i < times; i++ {
        Println(msg, !debug);
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

### Positional Arguments

`Args.args()` returns an `Array(Str)` of all remaining positional arguments (those not consumed by flags):

```zenth
fn main() {
    let verbose = Args.flag(default=false);
    let args = Args.args();

    for file in args {
        Println("Processing: {file}");
    }
}
```

```sh
./program --verbose file1.txt file2.txt
# args = ["file1.txt", "file2.txt"]
```

Positional arguments are everything that appears after the flags (or after `--` to force the end of flag parsing).

## Environment Variables

The `Env()` built-in reads environment variables. It always returns a `Str`.

```zenth
fn main() {
    // Read an environment variable (empty string if not set)
    let home = Env("HOME");
    Println("Home: " + home);

    // Provide a default value for when the variable is unset or empty
    let editor = Env("EDITOR", "vim");
    Println("Editor: " + editor);

    // Use in conditionals
    let debug = Env("DEBUG", "false");
    if debug == "true" {
        Println("Debug mode enabled");
    }
}
```

With one argument, `Env(name)` returns the value of the environment variable, or an empty string `""` if it is not set. With two arguments, `Env(name, default)` returns the default value when the variable is unset or empty.

## Naming Rules

Zenth reserves leading uppercase for system names:

- Built-in functions are capitalized (`Println`, `Len`, `Range`, etc.).
- User-defined **top-level** functions must start with a lowercase letter.
- User-defined variables (`let`, `var`, `const`, including destructuring and loop bindings) must start with a lowercase letter.

Examples:

```zenth
fn build_report() {    // ok
    let debug = Args.flag(default=false);   // ok
    let value = 10;                    // ok
}

fn BuildReport() { }   // error
let Debug = true;      // error
const MaxSize = 1024;  // error
```

## Closures

Closures are anonymous functions declared with `fn` in expression position. They are used with array methods like `.map()` and `.filter()`:

```zenth
let doubled = [1, 2, 3].map(fn(x) x * 2);
let evens = [1, 2, 3, 4].filter(fn(x) x % 2 == 0);
```

Parameter types are inferred from context when used with `.map()` or `.filter()`. You can also provide explicit types:

```zenth
fn(x: Int) -> Int x * 2
```

Block body closures use explicit `return`:

```zenth
let processed = nums.map(fn(x: Int) -> Int {
    let y = x + 10;
    return y;
});
```

## Function Types

Functions can be passed as arguments to other functions using function type annotations. The syntax for a function type is `Fn(ParamTypes) -> ReturnType`:

```zenth
fn apply(x: Int, f: Fn(Int) -> Int) -> Int {
    return f(x);
}

fn double(x: Int) -> Int {
    return x * 2;
}

fn main() {
    // Pass a named function
    Println(apply(5, double));     // 10

    // Pass a closure
    Println(apply(5, fn(x) x + 10));  // 15
}
```

### Function Type Syntax

Function types use `Fn(ParamTypes) -> ReturnType`. Omit the `-> ReturnType` for void functions:

```zenth
Fn(Int) -> Int           // takes Int, returns Int
Fn(Str, Int) -> Bool     // takes Str and Int, returns Bool
Fn(Int)                  // takes Int, returns nothing (void)
Fn() -> Str              // takes nothing, returns Str
```

Note: `Fn` (capitalized) is the type name, while `fn` (lowercase) is the keyword for declaring functions and closures.

### Higher-Order Functions

You can write your own functions that accept function parameters:

```zenth
fn count_matching(items: Array(Int), pred: Fn(Int) -> Bool) -> Int {
    var count = 0;
    for item in items {
        if pred(item) {
            count++;
        }
    }
    return count;
}

fn is_even(x: Int) -> Bool {
    return x % 2 == 0;
}

fn main() {
    let nums = [1, 2, 3, 4, 5, 6];
    Println(count_matching(nums, is_even));       // 3
    Println(count_matching(nums, fn(x) x > 3));   // 3
}
```

### Function-Typed Variables

Variables can hold function values:

```zenth
let f: Fn(Int) -> Int = double;
Println(f(7));    // 14
```

### Passing Functions to map/filter/reduce

Named functions can be passed directly to `.map()`, `.filter()`, and `.reduce()`:

```zenth
fn double(x: Int) -> Int { return x * 2; }
fn is_even(x: Int) -> Bool { return x % 2 == 0; }
fn add(a: Int, b: Int) -> Int { return a + b; }

fn main() {
    let nums = [1, 2, 3, 4, 5];
    let doubled = nums.map(double);        // [2, 4, 6, 8, 10]
    let evens = nums.filter(is_even);      // [2, 4]
    let total = nums.reduce(add, 0);       // 15
}
```

## Methods

Functions can be defined inside objects to act as methods. See [Objects](objects.md) for details.
