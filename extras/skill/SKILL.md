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
    Println("Hello, World!");
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
let x: Int = 5;         // immutable, explicit type
var counter = 0;         // mutable, type inferred
var counter: Int = 0;    // mutable, explicit type
const pi = 3.14159;      // constant
INT_MAX                  // built-in: max Int value (9223372036854775807)
INT_MIN                  // built-in: min Int value (-9223372036854775808)
```

- `let` is immutable (cannot reassign)
- `var` is mutable (can reassign)
- `const` is a compile-time constant
- User-defined variable names must start with a lowercase letter
- No `:=` operator; `let`/`var`/`const` is the declaration signal, `=` is always assignment
- Compound assignment: `+=`, `-=`, `*=`, `/=`
- Increment/decrement: `++`, `--` (statements, not expressions)

### Types

**Primitive:** `Int`, `Float`, `Bool`, `Str`, `Byte`, `Date`

**Composite:** `Array(Int)`, `Tuple(Str, Int)`, `Hashmap(Str, Int)`, `Set(Str)`, `Set(Tuple(Int, Int))`, obj types

**Function types:** `Fn(Int) -> Int`, `Fn(Str, Int) -> Bool`, `Fn(Int)` (void)

**Type aliases:**
```zenth
type Grid = Hashmap(Point, Int);
type Pair = Tuple(Str, Int);
```

**Type conversions:**
```zenth
let n = Int("42");        // Str to Int
let f = Float(42);        // Int to Float
let s = Str(42);          // any to Str
let b = Int("ff", 16);   // Str to Int with base (255)
let bin = Int("1010", 2); // binary string to Int (10)
```

### Functions

```zenth
fn add(a: Int, b: Int) -> Int {
    return a + b;
}

fn greet(name: Str) {
    Println("Hello, " + name);
}

// Default parameters
fn connect(host: Str, port: Int = 8080) {
    Println("{host}:{port}");
}

// Named arguments at call site
connect(host="localhost", port=3000);
connect(host="localhost");  // uses default port

// Recursion works
fn factorial(n: Int) -> Int {
    if n <= 1 { return 1; }
    return n * factorial(n - 1);
}
```

Function signatures require explicit types. Return type is omitted for void functions.

Top-level user-defined function names must start with a lowercase letter.

**Multi-return functions** use parenthesized return types and `return a, b;` shorthand:

```zenth
fn min_max(nums: Array(Int)) (Int,Int) {
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

`(Int,Int)` is shorthand for `-> Tuple(Int,Int)`. `return a, b;` is shorthand for `return Tuple(a, b);`.

**Function types** — functions can be passed as arguments using `Fn(Types) -> ReturnType` syntax:

```zenth
fn apply(x: Int, f: Fn(Int) -> Int) -> Int {
    return f(x);
}

fn double(x: Int) -> Int { return x * 2; }

fn main() {
    Println(apply(5, double));          // 10 (named function)
    Println(apply(5, fn(x) x + 10));   // 15 (closure)

    // Function-typed variables
    let f: Fn(Int) -> Int = double;
    Println(f(7));  // 14

    // Named functions work with map/filter/reduce
    let nums = [1, 2, 3];
    let doubled = nums.map(double);       // [2, 4, 6]
}
```

Function type syntax: `Fn(Int) -> Int`, `Fn(Str, Int) -> Bool`, `Fn(Int)` (void return).

### Strings

Double-quoted strings support interpolation; single-quoted strings are raw:

```zenth
let name = "World";
Println("Hello {name}!");          // interpolation
Println("{a} + {b} = {a + b}");    // expressions in braces
Println('Raw {name} string');      // no interpolation
Println("Use \{braces\} literally"); // escaped braces
```

**Multi-line strings** use triple double-quotes:

```zenth
let text = """
Hello
World
""";
// text is "Hello\nWorld"

// Supports interpolation
let name = "Zenth";
let msg = """
Welcome to {name}!
Enjoy coding.
""";
```

The first newline after opening `"""` and the last newline before closing `"""` are stripped.

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
"hello".pad_left(10)         // "     hello"
"hello".pad_right(10)        // "hello     "
"42".pad_left(5, "0")        // "00042"
Len("hello")                 // 5
"hello".length()             // 5
"hello".to_int()             // error; use on numeric strings
"42".to_int()                // 42
"ff".to_int(16)              // 255 (with base)
"5".is_digit()               // true (single char is a decimal digit 0-9)
"a".is_digit()               // false
(-9).sign()                  // -1
(0).sign()                   // 0
(7).sign()                   // 1
```

**String indexing and slicing:**
```zenth
let s = "hello";
let ch = s[0];      // "h" (returns Str, not Byte)
let sub = s[1:3];   // "el"
let rest = s[2:];   // "llo"
let head = s[:3];   // "hel"
```

### Control Flow

```zenth
// if/else
if x > 10 {
    Println("big");
} else if x > 5 {
    Println("medium");
} else {
    Println("small");
}

// if as expression
let label = if x > 10 { "big" } else { "small" };

// C-style for loop
for var i = 0; i < 10; i++ {
    Println(i);
}

// for-in over arrays
for item in items {
    Println(item);
}
for i, item in items {
    Println("{i}: {item}");
}

// for-in over strings (yields each character as Str)
for ch in "hello" {
    Print(ch);
}
for i, ch in "hello" {
    Println("{i}: {ch}");
}

// for-range with count
for step in Range(0, 100) {
    Println(step);
}

// _ discard for unused loop variables
for _ in Range(0, 10) {
    Println("tick");
}
for _, item in items {   // discard index
    Println(item);
}

// tuple destructuring in for-in
for (name, score) in Zip(names, scores) {
    Println("{name}: {score}");
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
    1 => Println("Monday");
    2 => Println("Tuesday");
    _ => Println("Other");
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
    Println(i);
}
```

### Range Objects

`Range()` and `Rangei()` return a **range object** — lightweight, iterable, with O(1) `.contains()` check:

```zenth
let r  = Range(0, 10);    // exclusive end: [0, 9]
let ri = Rangei(0, 10);   // inclusive end: [0, 10]

// .contains() for bounds check (O(1))
Println(r.contains(5));   // true
Println(r.contains(10));  // false (exclusive)
Println(ri.contains(10)); // true (inclusive)

// Len() on range is O(1)
Println(Len(Range(0, 100, 3)));  // 34

// Range objects are iterable in for-in
for i in Range(0, 5) {
    Println(i);  // 0 1 2 3 4
}

// Negative step for reverse iteration
for i in Range(5, -1, -1) {
    Println(i);  // 5 4 3 2 1 0
}

// Named ranges for reusable bounds checking
let xrange = Rangei(0, 100);
let yrange = Rangei(0, 50);
if xrange.contains(x) && yrange.contains(y) {
    Println("in bounds");
}
```

Note: `.contains()` is a bounds check only, not sequence membership.

### Arrays

```zenth
let numbers = [1, 2, 3, 4, 5];
let names: Array(Str) = [];       // empty with type annotation

// Access and length
Println(numbers[0]);
Println(Len(numbers));
Println(numbers.length());

// Mutating (requires var)
var items: Array(Int) = [];
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

// Safe indexing with default
let val = numbers.get(0, 0);     // returns numbers[0] or 0 if out of bounds
let x = numbers.get(99, -1);    // returns -1 (index out of bounds)
let last = numbers.last();       // final element (zero value if empty)

// Containment
if numbers.exists(3) {
    Println("found");
}

// Closures: map and filter
let doubled = numbers.map(fn(x) x * 2);
let evens = numbers.filter(fn(x) x % 2 == 0);

// Block closure with explicit types
let processed = numbers.map(fn(x: Int) -> Int {
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
Println(words.join(" "));       // "hello world"
Println(["a", "b"].join(","));  // "a,b"

// Enumerate — pair each element with its index
let fruits = ["apple", "banana", "cherry"];
let indexed = fruits.enumerate();  // Array(Tuple(Int, Str))
// Multi-param closures destructure tuples automatically:
let labeled = fruits.enumerate().map(fn(i, s) "{i}: {s}");
let evens = fruits.enumerate().filter(fn(i, s) i % 2 == 0);

// Zip — combine parallel arrays into array of tuples
let names = ["Alice", "Bob"];
let scores = [95, 87];
for pair in Zip(names, scores) {
    Println(pair.0 + "=" + Str(pair.1));
}
// Zip 3+ arrays: Zip(a, b, c) → Array(Tuple(T1, T2, T3))

// Type conversions for string arrays
let strs = ["1", "2", "3"];
let ints = strs.to_int();     // [1, 2, 3]
let floats = strs.to_float();   // [1.0, 2.0, 3.0]

// Array destructuring — bind elements directly to variables
let [a, b, c] = [10, 20, 30];
Println(a);  // 10

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
let pair = Tuple("hello", 42);
Println(pair.0);           // "hello"
Println(pair.1);           // "42"

// Named tuple
let person = Tuple(name="Alice", age=30);
Println(person.name);      // "Alice"
Println(person.age);       // "30"
Println(person.0);         // also works by index

// Type annotation
let t: Tuple(key: Str, value: Int) = Tuple(key="x", value=42);

// Destructuring
let (k, v) = t;

// Destructuring in for-in loops
for (name, score) in Zip(["Alice", "Bob"], [95, 87]) {
    Println("{name}: {score}");
}
for (i, item) in items.enumerate() {
    Println("{i}: {item}");
}

// Array of named tuples
let pairs: Array(Tuple(name: Str, score: Int)) = [];
pairs.add(Tuple(name="Bob", score=95));
for p in pairs {
    Println("{p.name}: {Str(p.score)}");
}
```

### Hashmaps

```zenth
// Create with type arguments
var scores = Hashmap(Str,Int);
scores["Alice"] = 95;
scores["Bob"] = 80;

// Access
Println(scores["Alice"]);
Println(Len(scores));

// Default values (like Python's defaultdict)
var counts = Hashmap(Str,Int, default=0);
counts["x"] += 1;  // no KeyError, starts from 0

// Defaults are generic, not Int-only
var names = Hashmap(Str,Str, default="");
Println(names["missing"]);  // ""

// Nested generics are supported
var recipe_in = Hashmap(Str, Array(Str));
recipe_in["cake"] = ["flour", "eggs"];
recipe_in["cake"].add("sugar");  // direct mutation through hashmap index

// Iteration
for key, val in scores {
    Println("{key}: {Str(val)}");
}
for key in scores {       // keys only
    Println(key);
}

// Methods
let keys = scores.keys();       // Array(Str)
let vals = scores.values();     // Array(Int)
if scores.exists("Alice") {     // check key existence
    Println("found");
}

// Object keys (value-based lookup)
obj Point {
    x: Int;
    y: Int;
}
var grid = Hashmap(Point,Str);
grid[Point(x=1, y=2)] = "#";
Println(grid[Point(x=1, y=2)]);  // "#" (same-value lookup works)

// Tuple keys (composite keys for multi-dimensional lookups)
var cells = Hashmap(Tuple(Int,Int),Str);
cells[Tuple(0, 0)] = "origin";
cells[Tuple(1, 2)] = "point";
Println(cells[Tuple(0, 0)]);        // "origin"
Println(cells.exists(Tuple(1, 2))); // true
```

### Sets

```zenth
// Create a set
var visited = Set(Str);
visited.add("start");
visited.add("middle");
visited.exists("start");   // true
visited.remove("middle");
Println(Len(visited));       // 1
Println(visited.length());   // 1

// Integer set
var nums = Set(Int);
nums.add(1);
nums.add(2);
nums.add(2);  // duplicate, ignored
Println(Len(nums));       // 2

// Iterate over set
for n in nums {
    Println(n);
}

// Set of tuples (composite keys without string round-tripping)
var dots = Set(Tuple(Int,Int));
dots.add(Tuple(6, 10));
dots.add(Tuple(0, 14));
dots.add(Tuple(6, 10));  // duplicate, ignored
if dots.exists(Tuple(6, 10)) {
    Println("found");
}
dots.remove(Tuple(0, 14));
for dot in dots {
    Println("{dot.0},{dot.1}");
}
```

### Objects and Methods

```zenth
obj Rectangle {
    width: Float;
    height: Float;

    fn area() -> Float {
        return self.width * self.height;
    }

    fn scale(factor: Float) -> Rectangle {
        return Rectangle(
            width=self.width * factor,
            height=self.height * factor
        );
    }

    // Override default string representation
    fn string() -> Str {
        return "Rect({self.width}x{self.height})";
    }
}

// Constructor uses named arguments
let r = Rectangle(width=10.0, height=5.0);
Println(r.area());         // "50"
Println(r);                // "Rect(10x5)"

// Default field values
obj Account {
    balance: Int = 0;
    interest: Float = 2.5;
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
- Prints as its string value: `Println(Color.Red)` prints `"Red"`, `Println(HexColor.Red)` prints `"#FF0000"`
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
        var id: Int = 0;
        var name = "";
        rows.Scan(&id, &name);
        Println("{id}: {name}");
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
fn double(x: Int) -> Int {
    return x * 2;
}
```

Call module definitions via `module.name()` syntax:

```zenth
import "./utils";

fn main() {
    Println(Str(utils.double(21)));  // "42"
}
```

**Directory modules:** If the import path is a directory, all `.zn` files in that directory are merged into one module:

```zenth
import "./shapes";  // loads shapes/rect.zn, shapes/circle.zn, etc.

fn main() {
    let r = shapes.Rectangle(width=3.0, height=4.0);
    Println(Str(r.area()));
}
```

Module objects are constructed with `module.ObjName(field=value)` syntax and their methods work normally after construction.

**Module-internal function calls:** Functions within a module can call other functions in the same module:

```zenth
// mathutils.zn
fn add(a: Int, b: Int) -> Int { return a + b; }

fn sum_and_double(a: Int, b: Int) -> Int {
    return add(a, b) * 2;  // calls add within the same module
}
```

**Qualified types:** Use `module.Type` syntax to refer to types from imported modules in type annotations and generics:

```zenth
import "./items";

fn process(data: Array(items.Item)) {
    for item in data {
        Println(item.name);
    }
}
```

Qualified types work in all type positions: function parameters, variable annotations, and generic type arguments (`Array(mod.Type)`, `Tuple(Int, mod.Type)`, etc.).

### Built-in Functions

| Function | Description |
|---|---|
| `Print(...)` | Print without newline |
| `Println(...)` | Print with newline |
| `Len(x)` | Length of string, array, or hashmap |
| `Str(x)` | Convert any value to String |
| `Int(x)` | Convert string or float to Int |
| `Int(Str, base)` | Convert string to Int with base (2, 8, 16, etc.) |
| `Float(x)` | Convert to Float |
| `Ord(Str)` | Unicode codepoint of first character |
| `Chr(Int)` | String containing the given codepoint |
| `Abs(x)` | Absolute value (Int or float) |
| `Min(a, b, ...)` | Minimum of two or more values (all must be same type: Int or Float) |
| `Max(a, b, ...)` | Maximum of two or more values (all must be same type: Int or Float) |
| `Clamp(x, lo, hi)` | Clamp value to range |
| `Round(x)` | Round float to nearest integer |
| `Floor(x)` | Floor of float |
| `Ceil(x)` | Ceiling of float |
| `Trunc(x)` | Truncate toward zero |
| `Pow(base, exp)` | Exponentiation |
| `Pow10(n)` | 10^n (`Int` argument, returns `Float`) |
| `Sqrt(x)` | Square root |
| `Cbrt(x)` | Cube root |
| `Hypot(x, y)` | sqrt(x^2 + y^2) |
| `Sin(x)` / `Cos(x)` / `Tan(x)` | Trig functions (radians) |
| `Asin(x)` / `Acos(x)` / `Atan(x)` | Inverse trig |
| `Atan2(y, x)` | Two-argument arctangent |
| `Sinh(x)` / `Cosh(x)` / `Tanh(x)` | Hyperbolic trig |
| `Asinh(x)` / `Acosh(x)` / `Atanh(x)` | Inverse hyperbolic |
| `Log(x)` / `Log2(x)` / `Log10(x)` / `Log1p(x)` | Logarithms |
| `Exp(x)` / `Exp2(x)` / `Expm1(x)` | Exponentials |
| `Mod(x, y)` / `Remainder(x, y)` | Float modulo / IEEE remainder |
| `Dim(x, y)` | max(x-y, 0) |
| `Copysign(x, y)` | x with sign of y |
| `NaN()` | IEEE 754 not-a-number |
| `Inf(sign)` | Infinity; sign is `Int` (1 or -1) |
| `IsNaN(x)` | True if x is NaN (`Float` argument, returns `Bool`) |
| `IsInf(x, sign)` | True if x is infinity in given direction (sign `Int`: 1, -1, or 0) |
| `Signbit(x)` | True if x is negative (`Float` argument, returns `Bool`) |
| `Erf(x)` / `Erfc(x)` | Error function / complement |
| `Gamma(x)` / `Lgamma(x)` | Gamma and log-Gamma |
| Array `.total()` | Sum of all elements (same type as array) |
| Array `.mean()` | Arithmetic mean (returns `Float`) |
| Array `.median()` | Median value (returns `Float`) |
| Array `.stdev()` | Population standard deviation (returns `Float`) |
| `Range(start, end)` | Exclusive range object |
| `Range(start, end, step)` | Exclusive range object with step |
| `Rangei(start, end)` | Inclusive range object |
| `Rangei(start, end, step)` | Inclusive range object with step |
| `Hashmap(K, V)` | Create empty hashmap |
| `Hashmap(K, V, default=val)` | Create hashmap with default |
| `Set(T)` | Create empty set |
| `Tuple(...)` | Create Tuple (positional or named) |
| `File(path)` | Create file handle |
| `Env(name)` / `Env(name, default)` | Read environment variable (returns Str; default when unset) |
| `Input(text)` | Display editable pre-filled text, return result after Enter (returns Str) |
| `Args.flag(default=val)` | Command-line flag (name inferred from variable) |
| `Args.args()` | Get remaining positional arguments as `Array(Str)` |
| `Date.today()` | Get today's date as a Date object |
| `Date.from(str)` | Parse date from string (default format: `%Y-%m-%d`) |
| `Date.from(str, fmt)` | Parse date from string with custom format |
| `Exit(code)` | Exit program with status code |
| `Assert(cond)` | Panic if `cond` is false (reports file:line) |
| `AssertEq(got, expected)` | Panic if `got != expected` (reports file:line and both values) |

### Built-in Constants

| Constant | Description |
|---|---|
| `INT_MAX` | Maximum `Int` value (2^63-1 = 9223372036854775807) |
| `INT_MIN` | Minimum `Int` value (-2^63 = -9223372036854775808) |

These are true constants and cannot be reassigned. Use them instead of magic numbers for sentinel values:

```zenth
var best = INT_MAX;       // use as "infinity"
if cost < best {
    best = cost;
}
```

### File I/O

```zenth
let f = File("data.txt");

if f.exists() {
    let content = f.read();       // entire file as Str
    let lines = f.lines();        // Array(Str)
    let parts = f.sections();     // split on blank lines -> Array(Str)
    Println("Name: " + f.name()); // filename
    Println("Ext: " + f.ext());   // extension
}

// Writing files
f.write("Hello world\n");        // create/overwrite file
f.append("More content\n");      // append to file (creates if missing)
```

### Dates

```zenth
// Today's date
let d = Date.today();
Println(d.format("%Y-%m-%d"));

// Parse from string (default format: %Y-%m-%d)
let epoch = Date.from("1970-01-01");

// Parse with custom format
let jan13 = Date.from("01/13/2007", "%m/%d/%Y");

// Add time (default unit: "days", also: "months", "years")
let tomorrow = d.add(1);
let next_month = d.add(1, "months");
let next_year = d.add(1, "years");

// Subtract time
let yesterday = d.sub(1);
let last_year = d.sub(1, "years");

// Chaining
let future = Date.today().add(1, "years").add(3, "months").add(10);
Println(future.format("%B %d, %Y"));
```

**Date methods:**
- `.format(fmt)` — format as string using Python-style specifiers (`%Y`, `%m`, `%d`, `%H`, `%M`, `%S`, `%B`, `%b`, `%A`, `%a`, `%y`)
- `.add(val)` / `.add(val, unit)` — add days/months/years, returns new Date
- `.sub(val)` / `.sub(val, unit)` — subtract days/months/years, returns new Date

### Command-Line Arguments

Flag names are inferred from the variable name:

```zenth
let debug = Args.flag(default=false);  // --debug (Bool)
let count = Args.flag(default=5);      // --count (Int)
let msg = Args.flag(default="hi");     // --msg (Str)
let args = Args.args();               // remaining positional args
```

Run with: `./program --debug --count=10 --msg="hello" file1.txt file2.txt`

### Environment Variables

```zenth
let home = Env("HOME");           // returns Str, empty if unset
let editor = Env("EDITOR", "vim"); // returns "vim" if EDITOR is unset
```

### User Input

```zenth
let name = Input("World");       // shows "World" as editable text, user can modify
Println("Hello, {name}!");

let age = Int(Input("25"));      // pre-fill "25", user can edit, result converted to Int
```

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
    let data = File("input.txt");
    for i, line in data.lines() {
        Println("{i}: {line}");
    }
}
```

### Hashmap Word Counter
```zenth
fn main() {
    var counts = Hashmap(Str,Int, default=0);
    let words = ["apple", "banana", "apple", "cherry", "banana", "apple"];
    for word in words {
        counts[word] += 1;
    }
    for word, count in counts {
        Println("{word}: {Str(count)}");
    }
}
```

### Object with Methods
```zenth
obj Vec2 {
    x: Float;
    y: Float;

    fn add(other: Vec2) -> Vec2 {
        return Vec2(x=self.x + other.x, y=self.y + other.y);
    }

    fn magnitude() -> Float {
        return Sqrt(self.x * self.x + self.y * self.y);
    }

    fn string() -> Str {
        return "({self.x}, {self.y})";
    }
}

fn main() {
    let a = Vec2(x=3.0, y=4.0);
    let b = Vec2(x=1.0, y=2.0);
    let c = a.add(b);
    Println("Sum: {c}");
    Println("Magnitude: {a.magnitude()}");
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
                    _ => Str(i)
                }
            }
        };
        Println(result);
    }
}
```

### Processing with Map/Filter
```zenth
fn main() {
    let data = File("numbers.txt");
    let values = data.lines().to_int();
    let big_doubled = values.filter(fn(x) x > 10).map(fn(x) x * 2);
    for v in big_doubled {
        Println(v);
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
    AssertEq(2 + 3, 5);
    Assert(10 > 0);
}

fn test_strings() {
    AssertEq("hello".upper(), "HELLO");
    Assert("hello".contains("ell"));
}
```

### Assertions

- `Assert(cond)` -- panics with `file:line` if `cond` is false
- `AssertEq(got, expected)` -- panics with `file:line` and both values if `got != expected`

Both assertions are available in all Zenth programs, not just test files.

### Testing with Module Imports

Test files can import local modules using `../` paths:

```zenth
// tests/task_test.zn
import "../src/task";

fn test_parse_task() {
    let t = task.parse_task("buy milk");
    AssertEq(t.name, "buy milk");
}

fn test_construct() {
    let t = task.Task(name="clean", state="done");
    AssertEq(t.state, "done");
}
```

All module features work in tests: function calls, object constructors with named arguments, method calls, and qualified type annotations. Import paths must start with `./` or `../`.

### Conventions

- Test files: `*_test.zn` in a `tests/` directory
- Test functions: `fn test_*()` with no parameters and no return type
- Helper functions (without `test_` prefix) can be defined in test files
- Exit code is `0` if all tests pass, `1` if any fail

### Bitwise Operators

Bitwise operators work on `Int` operands only:

| Operator | Description |
|----------|-------------|
| `&` | Bitwise AND |
| `\|` | Bitwise OR |

```zenth
let flags = 0b1100;
let mask  = 0b1010;
Println(Str(flags & mask));  // 8  (0b1000)
Println(Str(flags | mask));  // 14 (0b1110)
```

## Important Notes

- Every program needs a `fn main() { ... }` entry point (except test files)
- Statements end with semicolons `;`
- Blocks use curly braces `{ }`
- `let` variables cannot be reassigned; use `var` for mutable state
- Function parameters require type annotations; local variables use inference
- String interpolation only works in double-quoted strings: `"hello {name}"`
- Single-quoted strings are raw (no interpolation, no escape sequences): `'raw'`
- Triple double-quoted strings (`"""..."""`) are multi-line with interpolation support
- The `for` keyword is used for all loops (no `while` keyword)
- `match` replaces `switch`; use `_` for the default case
- Objects use `self` (not `this`) for method access
- Object constructors always use named arguments: `Point(x=1, y=2)`
- **match statements** use semicolons: `1 => Println("one");`
- **match expressions** (used as values) use commas: `1 => "one",`
- Trailing commas are allowed in array literals and function call arguments
- Empty array `[]` can be passed where the type is known from context (field, variable annotation)
