---
name: zenth
description: Zenth is a programming language that combines Go and Python syntax, it transpiles Zenth to Go and then uses Go build tools to compile down to a binary. Use this skill for writing Zenth programs.
---

# Zenth Programming Language - LLM Agent Guide

Zenth is a compiled programming language that blends Go's type discipline with Python's simplicity. It uses curly braces and semicolons for structure and compiles to native binaries.

- **File extension:** `.zn`
- **Documentation:** <https://github.com/mkaz/zenth/tree/trunk/docs>

## Quick Start

```zenth
// hello.zn
fn main() {
    println("Hello, World!");
}
```

To run a zenth program use:
```sh
zenth run hello.zn
```

## Build and Compile

To build a binary from a zenth program:
```sh
zenth build -o hello hello.zn
```

To output the intermediate Go code:
```sh
zenth build --emit-go hello.zn
```

## Language Reference

### Variables

```zenth
let x = 5;              // immutable, type inferred
let x: int = 5;         // immutable, explicit type
var counter = 0;         // mutable, type inferred
var counter: int = 0;    // mutable, explicit type
const PI = 3.14159;      // constant
INT_MAX                  // built-in: max int value (9223372036854775807)
INT_MIN                  // built-in: min int value (-9223372036854775808)
```

- `let` is immutable (cannot reassign)
- `var` is mutable (can reassign)
- `const` is a compile-time constant
- No `:=` operator; `let`/`var`/`const` is the declaration signal, `=` is always assignment
- Compound assignment: `+=`, `-=`, `*=`, `/=`
- Increment/decrement: `++`, `--` (statements, not expressions)

### Types

**Primitive:** `int`, `i8`, `i16`, `i32`, `i64`, `u8`, `u16`, `u32`, `u64`, `f32`, `f64`, `bool`, `str`, `byte`

**Composite:** `array(int)`, `tuple(str, int)`, `hashmap(str, int)`, `set(str)`, `set(tuple(int, int))`, obj types

**Type aliases:**
```zenth
type Grid = hashmap(Point, int);
type Pair = tuple(str, int);
```

**Type conversions:**
```zenth
let n = int("42");        // str to int
let f = f64(42);          // int to f64
let s = str(42);          // any to str
let b = int("ff", 16);   // str to int with base (255)
let bin = int("1010", 2); // binary string to int (10)
```

### Functions

```zenth
fn add(a: int, b: int) -> int {
    return a + b;
}

fn greet(name: str) {
    println("Hello, " + name);
}

// Default parameters
fn connect(host: str, port: int = 8080) {
    println("{host}:{port}");
}

// Named arguments at call site
connect(host="localhost", port=3000);
connect(host="localhost");  // uses default port

// Recursion works
fn factorial(n: int) -> int {
    if n <= 1 { return 1; }
    return n * factorial(n - 1);
}
```

Function signatures require explicit types. Return type is omitted for void functions.

**Multi-return functions** use parenthesized return types and `return a, b;` shorthand:

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

`-> (int, int)` is shorthand for `-> tuple(int, int)`. `return a, b;` is shorthand for `return tuple(a, b);`.

### Strings

Double-quoted strings support interpolation; single-quoted strings are raw:

```zenth
let name = "World";
println("Hello {name}!");          // interpolation
println("{a} + {b} = {a + b}");    // expressions in braces
println('Raw {name} string');      // no interpolation
println("Use \{braces\} literally"); // escaped braces
```

**String methods:**
```zenth
"hello".upper()              // "HELLO"
"HeLLo".lower()              // "hello"
"zenth".starts_with("zen")   // true
"zenth".ends_with("th")      // true
"fold along x=5".strip_prefix("fold along ") // "x=5"
"photo.png".strip_suffix(".png")             // "photo"
"  hi  ".strip()             // "hi"
"..hi..".strip(".")          // "hi"
"banana".find("na")          // 2 (index)
"banana".count("na")         // 2
"a-b-a".replace("a", "x")   // "x-b-x"
"a-b-a".replace("a", "x", 1) // "x-b-a" (limit replacements)
"a,b,c".split(",")           // ["a", "b", "c"]
"hello world".split()        // ["hello", "world"] (whitespace)
"hi".contains("h")           // true
".".repeat(5)                // "....."
len("hello")                 // 5
"hello".length()             // 5
"hello".to_int()             // error; use on numeric strings
"42".to_int()                // 42
"ff".to_int(16)              // 255 (with base)
"5".is_digit()               // true (single char is a decimal digit 0-9)
"a".is_digit()               // false
```

**String indexing and slicing:**
```zenth
let s = "hello";
let ch = s[0];      // "h" (returns str, not byte)
let sub = s[1:3];   // "el"
let rest = s[2:];   // "llo"
let head = s[:3];   // "hel"
```

### Control Flow

```zenth
// if/else
if x > 10 {
    println("big");
} else if x > 5 {
    println("medium");
} else {
    println("small");
}

// if as expression
let label = if x > 10 { "big" } else { "small" };

// C-style for loop
for var i = 0; i < 10; i++ {
    println(i);
}

// for-in over arrays
for item in items {
    println(item);
}
for i, item in items {
    println("{i}: {item}");
}

// for-in over strings (yields each character as str)
for ch in "hello" {
    print(ch);
}
for i, ch in "hello" {
    println("{i}: {ch}");
}

// for-range with count
for step in range(0, 100) {
    println(step);
}

// _ discard for unused loop variables
for _ in range(0, 10) {
    println("tick");
}
for _, item in items {   // discard index
    println(item);
}

// while-style
for running {
    process();
}

// infinite loop
for {
    listen();
    if done { break; }
}

// match statement
match day {
    1 => println("Monday");
    2 => println("Tuesday");
    _ => println("Other");
}

// match as expression (arms use commas, not semicolons)
let label = match day {
    1 => "Monday",
    2 => "Tuesday",
    _ => "Other"
};

// break and continue work in all loops
for var i = 0; i < 10; i++ {
    if i == 5 { break; }
    if i % 2 == 0 { continue; }
    println(i);
}
```

### Range Objects

`range()` and `rangei()` return a **range object** — lightweight, iterable, with O(1) `.contains()` check:

```zenth
let r  = range(0, 10);    // exclusive end: [0, 9]
let ri = rangei(0, 10);   // inclusive end: [0, 10]

// .contains() for bounds check (O(1))
println(r.contains(5));   // true
println(r.contains(10));  // false (exclusive)
println(ri.contains(10)); // true (inclusive)

// len() on range is O(1)
println(len(range(0, 100, 3)));  // 34

// Range objects are iterable in for-in
for i in range(0, 5) {
    println(i);  // 0 1 2 3 4
}

// Named ranges for reusable bounds checking
let xrange = rangei(0, 100);
let yrange = rangei(0, 50);
if xrange.contains(x) && yrange.contains(y) {
    println("in bounds");
}
```

Note: `.contains()` is a bounds check only, not sequence membership.

### Arrays

```zenth
let numbers = [1, 2, 3, 4, 5];
let names: array(str) = [];       // empty with type annotation

// Access and length
println(numbers[0]);
println(len(numbers));
println(numbers.length());

// Mutating (requires var)
var items: array(int) = [];
items[0] = 99;         // direct index assignment
items[0] += 1;         // compound assignment by index
items.add(10);         // append
items.push(20);        // same as add
let last = items.pop();     // remove and return last
let at = items.pop(0);     // remove and return at index
items.insert(1, 99);   // insert 99 at index 1, shifting later elements right
items.remove(2);       // remove element at index 2 (discards it)
items.extend([30, 40]);  // append all elements from another array

// Repeat — create arrays filled with repeated values
let zeros = [0].repeat(10);      // [0, 0, 0, ..., 0] (10 zeros)
let pattern = [1, 2].repeat(6);  // [1, 2, 1, 2, 1, 2]

// Containment
if numbers.exists(3) {
    println("found");
}

// Closures: map and filter
let doubled = numbers.map(fn(x) x * 2);
let evens = numbers.filter(fn(x) x % 2 == 0);

// Block closure with explicit types
let processed = numbers.map(fn(x: int) -> int {
    let y = x + 1;
    return y;
});

// Chaining
let result = [1, 2, 3, 4, 5].filter(fn(x) x > 2).map(fn(x) x * 10);

// Slicing
let sub = numbers[1:3];   // [2, 3]

// Concatenation
let combined = [1, 2] + [3, 4];

// Max, min, sum on numeric arrays
let biggest = numbers.max();    // 5
let smallest = numbers.min();   // 1
let total = numbers.sum();      // 15

// Reduce — combine elements with a closure
let product = numbers.reduce(fn(a, b) a * b, 1);  // 120
let sum = numbers.reduce(fn(a, b) a + b);          // 15

// Sorted (returns new sorted copy)
let s = [3, 1, 4].sorted();          // [1, 3, 4]
let sw = ["b", "a"].sorted();        // ["a", "b"]
let desc = [3, 1, 4].sorted("desc"); // [4, 3, 1] — largest to smallest
let asc = [3, 1, 4].sorted("asc");   // [1, 3, 4] — smallest to largest (default)

// Join string arrays into a single string
let words = ["hello", "world"];
println(words.join(" "));       // "hello world"
println(["a", "b"].join(","));  // "a,b"

// Zip — combine parallel arrays into array of tuples
let names = ["Alice", "Bob"];
let scores = [95, 87];
for pair in zip(names, scores) {
    println(pair.0 + "=" + str(pair.1));
}
// Zip 3+ arrays: zip(a, b, c) → array(tuple(T1, T2, T3))

// Type conversions for string arrays
let strs = ["1", "2", "3"];
let ints = strs.to_int();     // [1, 2, 3]
let floats = strs.to_f64();   // [1.0, 2.0, 3.0]

// Array destructuring — bind elements directly to variables
let [a, b, c] = [10, 20, 30];
println(a);  // 10

// Works with var (mutable) and any array expression
var [x, y] = [1, 2];
x = 99;

// Use _ to discard positions
let [first, _, third] = [1, 2, 3];

// Destructure from a method chain
let [l, w, h] = "12x4x8".split("x").to_int().sorted();
```

### Tuples

```zenth
// Positional tuple
let pair = tuple("hello", 42);
println(pair.0);           // "hello"
println(pair.1);           // "42"

// Named tuple
let person = tuple(name="Alice", age=30);
println(person.name);      // "Alice"
println(person.age);       // "30"
println(person.0);         // also works by index

// Type annotation
let t: tuple(key: str, value: int) = tuple(key="x", value=42);

// Destructuring
let (k, v) = t;

// Array of named tuples
let pairs: array(tuple(name: str, score: int)) = [];
pairs.add(tuple(name="Bob", score=95));
for p in pairs {
    println("{p.name}: {str(p.score)}");
}
```

### Hashmaps

```zenth
// Create with type arguments
var scores = hashmap(str, int);
scores["Alice"] = 95;
scores["Bob"] = 80;

// Access
println(scores["Alice"]);
println(len(scores));

// Default values (like Python's defaultdict)
var counts = hashmap(str, int, default=0);
counts["x"] += 1;  // no KeyError, starts from 0

// Iteration
for key, val in scores {
    println("{key}: {str(val)}");
}
for key in scores {       // keys only
    println(key);
}

// Methods
let keys = scores.keys();       // array(str)
let vals = scores.values();     // array(int)
if scores.exists("Alice") {     // check key existence
    println("found");
}

// Object keys (value-based lookup)
obj Point {
    x: int;
    y: int;
}
var grid = hashmap(Point, str);
grid[Point(x=1, y=2)] = "#";
println(grid[Point(x=1, y=2)]);  // "#" (same-value lookup works)

// Tuple keys (composite keys for multi-dimensional lookups)
var cells = hashmap(tuple(int, int), str);
cells[tuple(0, 0)] = "origin";
cells[tuple(1, 2)] = "point";
println(cells[tuple(0, 0)]);        // "origin"
println(cells.exists(tuple(1, 2))); // true
```

### Sets

```zenth
// Create a set
var visited = set(str);
visited.add("start");
visited.add("middle");
visited.exists("start");   // true
visited.remove("middle");
println(len(visited));       // 1
println(visited.length());   // 1

// Integer set
var nums = set(int);
nums.add(1);
nums.add(2);
nums.add(2);  // duplicate, ignored
println(len(nums));       // 2

// Iterate over set
for n in nums {
    println(n);
}

// Set of tuples (composite keys without string round-tripping)
var dots = set(tuple(int, int));
dots.add(tuple(6, 10));
dots.add(tuple(0, 14));
dots.add(tuple(6, 10));  // duplicate, ignored
if dots.exists(tuple(6, 10)) {
    println("found");
}
dots.remove(tuple(0, 14));
for dot in dots {
    println("{dot.0},{dot.1}");
}
```

### Objects and Methods

```zenth
obj Rectangle {
    width: f64;
    height: f64;

    fn area() -> f64 {
        return self.width * self.height;
    }

    fn scale(factor: f64) -> Rectangle {
        return Rectangle(
            width=self.width * factor,
            height=self.height * factor
        );
    }

    // Override default string representation
    fn string() -> str {
        return "Rect({self.width}x{self.height})";
    }
}

// Constructor uses named arguments
let r = Rectangle(width=10.0, height=5.0);
println(r.area());         // "50"
println(r);                // "Rect(10x5)"

// Default field values
obj Account {
    balance: int = 0;
    interest: f64 = 2.5;
}
let a = Account();                    // all defaults
let b = Account(balance=100);         // partial
let c = Account(interest=5.0, balance=1000);  // any order
```

- Methods use `self` to access fields and call other methods
- Fields are accessed with dot notation: `r.width`
- Constructors use named arguments: `Point(x=1, y=2)`

### Enums

Enums are string-backed. Without explicit values, the variant name is the value.

```zenth
enum Color {
    Red;        // value is "Red"
    Green;      // value is "Green"
    Blue;       // value is "Blue"
}

enum HexColor {
    Red = "#FF0000";
    Green = "#00FF00";
    Blue = "#0000FF";
}
```

- Access variants with `EnumName.Variant`: `Color.Red`, `HexColor.Red`
- Backed by strings: default is variant name, or explicit string value
- Compare with `==` and `!=`
- Prints as its string value: `println(Color.Red)` prints `"Red"`, `println(HexColor.Red)` prints `"#FF0000"`
- Use in `match` statements and expressions
- Use as type annotations to enforce valid values: `let c: Color = Color.Green;`

### Imports

```zenth
import "fmt";
import "math" as m;
```

Available modules map to Go stdlib: `fmt`, `math`, `os`, `strings`

### External Go Package Imports

Use `import_go` to import any Go package, including third-party libraries. The toolchain runs `go get` automatically before compiling.

```zenth
import_go "github.com/some/library";
import_go "github.com/some/library" as lib;
```

Calls on `import_go` packages bypass type checking and pass through to Go verbatim. Return values are untyped — you can assign them to variables and chain further calls.

```zenth
import_go "database/sql";
import_go "github.com/mattn/go-sqlite3" as _;

fn main() {
    let db = sql.Open("sqlite3", "./data.db");
    let rows = db.Query("SELECT id, name FROM users");
    for rows.Next() {
        var id: int = 0;
        var name = "";
        rows.Scan(&id, &name);
        println("{id}: {name}");
    }
}
```

### Local Modules

Import your own `.zn` files using a path starting with `./` or `../`. The module name is the filename stem:

```zenth
import "./utils";
import "./math/geometry" as geo;
import "../shared/helpers" as h;
```

A module file is a regular `.zn` file without `fn main`. It can contain functions, objects, constants, and type aliases:

```zenth
// utils.zn
fn double(x: int) -> int {
    return x * 2;
}
```

Call module definitions via `module.name()` syntax:

```zenth
import "./utils";

fn main() {
    println(str(utils.double(21)));  // "42"
}
```

**Directory modules:** If the import path is a directory, all `.zn` files in that directory are merged into one module:

```zenth
import "./shapes";  // loads shapes/rect.zn, shapes/circle.zn, etc.

fn main() {
    let r = shapes.Rectangle(width=3.0, height=4.0);
    println(str(r.area()));
}
```

Module objects are constructed with `module.ObjName(field=value)` syntax and their methods work normally after construction.

### Built-in Functions

| Function | Description |
|---|---|
| `print(...)` | Print without newline |
| `println(...)` | Print with newline |
| `len(x)` | Length of string, array, or hashmap |
| `str(x)` | Convert any value to string |
| `int(x)` | Convert string or float to int |
| `int(str, base)` | Convert string to int with base (2, 8, 16, etc.) |
| `f64(x)` | Convert to float64 |
| `abs(x)` | Absolute value (int or float) |
| `min(a, b, ...)` | Minimum of two or more values (all must be same type: int or f64) |
| `max(a, b, ...)` | Maximum of two or more values (all must be same type: int or f64) |
| `clamp(x, lo, hi)` | Clamp value to range |
| `round(x)` | Round float to nearest int |
| `floor(x)` | Floor of float |
| `ceil(x)` | Ceiling of float |
| `pow(base, exp)` | Exponentiation |
| `sqrt(x)` | Square root |
| `range(start, end)` | Exclusive range object |
| `range(start, end, step)` | Exclusive range object with step |
| `rangei(start, end)` | Inclusive range object |
| `rangei(start, end, step)` | Inclusive range object with step |
| `hashmap(K, V)` | Create empty hashmap |
| `hashmap(K, V, default=val)` | Create hashmap with default |
| `set(T)` | Create empty set |
| `tuple(...)` | Create tuple (positional or named) |
| `file(path)` | Create file handle |
| `flag(default=val)` | Command-line flag (name inferred from variable) |
| `exit(code)` | Exit program with status code |
| `assert(cond)` | Panic if `cond` is false (reports file:line) |
| `assert_eq(got, expected)` | Panic if `got != expected` (reports file:line and both values) |

### Built-in Constants

| Constant | Description |
|---|---|
| `INT_MAX` | Maximum `int` value (2^63-1 = 9223372036854775807) |
| `INT_MIN` | Minimum `int` value (-2^63 = -9223372036854775808) |

These are true constants and cannot be reassigned. Use them instead of magic numbers for sentinel values:

```zenth
var best = INT_MAX;       // use as "infinity"
if cost < best {
    best = cost;
}
```

### File I/O

```zenth
let f = file("data.txt");

if f.exists() {
    let content = f.read();       // entire file as str
    let lines = f.lines();        // array(str)
    let parts = f.sections();     // split on blank lines -> array(str)
    println("Name: " + f.name()); // filename
    println("Ext: " + f.ext());   // extension
}
```

### Command-Line Flags

Flag names are inferred from the variable name:

```zenth
let debug = flag(default=false);  // --debug flag (bool)
let count = flag(default=5);      // --count flag (int)
let msg = flag(default="hi");     // --msg flag (str)
```

Run with: `./program --debug --count=10 --msg="hello"`

### Multi-Variable Assignment

```zenth
var a = 1;
var b = 2;

// Swap without temp variable
a, b = b, a;
```

## Common Patterns

### Reading a File Line by Line
```zenth
fn main() {
    let data = file("input.txt");
    for i, line in data.lines() {
        println("{i}: {line}");
    }
}
```

### Hashmap Word Counter
```zenth
fn main() {
    var counts = hashmap(str, int, default=0);
    let words = ["apple", "banana", "apple", "cherry", "banana", "apple"];
    for word in words {
        counts[word] += 1;
    }
    for word, count in counts {
        println("{word}: {str(count)}");
    }
}
```

### Object with Methods
```zenth
obj Vec2 {
    x: f64;
    y: f64;

    fn add(other: Vec2) -> Vec2 {
        return Vec2(x=self.x + other.x, y=self.y + other.y);
    }

    fn magnitude() -> f64 {
        return sqrt(self.x * self.x + self.y * self.y);
    }

    fn string() -> str {
        return "({self.x}, {self.y})";
    }
}

fn main() {
    let a = Vec2(x=3.0, y=4.0);
    let b = Vec2(x=1.0, y=2.0);
    let c = a.add(b);
    println("Sum: {c}");
    println("Magnitude: {a.magnitude()}");
}
```

### FizzBuzz
```zenth
fn main() {
    for var i = 1; i <= 100; i++ {
        let result = match i % 15 {
            0 => "FizzBuzz",
            _ => match i % 3 {
                0 => "Fizz",
                _ => match i % 5 {
                    0 => "Buzz",
                    _ => str(i)
                }
            }
        };
        println(result);
    }
}
```

### Processing with Map/Filter
```zenth
fn main() {
    let data = file("numbers.txt");
    let values = data.lines().to_int();
    let big_doubled = values.filter(fn(x) x > 10).map(fn(x) x * 2);
    for v in big_doubled {
        println(v);
    }
}
```

## Testing

Zenth has a built-in test runner. Test files use `_test.zn` suffix and live in a `tests/` directory.

### Running Tests

```sh
zenth test               # run all tests in tests/
zenth test path/to/dir   # run tests in specific directory
zenth test file_test.zn  # run a single test file
```

### Writing Tests

Test functions start with `test_` and need no `fn main()`:

```zenth
// tests/math_test.zn

fn test_addition() {
    assert_eq(2 + 3, 5);
    assert(10 > 0);
}

fn test_strings() {
    assert_eq("hello".upper(), "HELLO");
    assert("hello".contains("ell"));
}
```

### Assertions

- `assert(cond)` -- panics with `file:line` if `cond` is false
- `assert_eq(got, expected)` -- panics with `file:line` and both values if `got != expected`

Both assertions are available in all Zenth programs, not just test files.

### Conventions

- Test files: `*_test.zn` in a `tests/` directory
- Test functions: `fn test_*()` with no parameters and no return type
- Helper functions (without `test_` prefix) can be defined in test files
- Exit code is `0` if all tests pass, `1` if any fail

## Important Notes

- Every program needs a `fn main() { ... }` entry point (except test files)
- Statements end with semicolons `;`
- Blocks use curly braces `{ }`
- `let` variables cannot be reassigned; use `var` for mutable state
- Function parameters require type annotations; local variables use inference
- String interpolation only works in double-quoted strings: `"hello {name}"`
- Single-quoted strings are raw (no interpolation, no escape sequences): `'raw'`
- The `for` keyword is used for all loops (no `while` keyword)
- `match` replaces `switch`; use `_` for the default case
- Objects use `self` (not `this`) for method access
- Object constructors always use named arguments: `Point(x=1, y=2)`
- **match statements** use semicolons: `1 => println("one");`
- **match expressions** (used as values) use commas: `1 => "one",`
- Trailing commas are allowed in array literals and function call arguments
- Empty array `[]` can be passed where the type is known from context (field, variable annotation)

