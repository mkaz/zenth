---
name: zenth
description: Zenth is a programming language that combines Go and Python syntax,
it transpiles Zenth to Go and then uses Go build tools to compile down to a
binary. Use this skill for writing Zenth programs.
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
```

- `let` is immutable (cannot reassign)
- `var` is mutable (can reassign)
- `const` is a compile-time constant
- No `:=` operator; `let`/`var`/`const` is the declaration signal, `=` is always assignment
- Compound assignment: `+=`, `-=`, `*=`, `/=`
- Increment/decrement: `++`, `--` (statements, not expressions)

### Types

**Primitive:** `int`, `i8`, `i16`, `i32`, `i64`, `u8`, `u16`, `u32`, `u64`, `f32`, `f64`, `bool`, `str`, `byte`

**Composite:** `array(int)`, `tuple(str, int)`, `hashmap(str, int)`, obj types

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
"  hi  ".strip()             // "hi"
"..hi..".strip(".")          // "hi"
"banana".find("na")          // 2 (index)
"banana".count("na")         // 2
"a-b-a".replace("a", "x")   // "x-b-x"
"a-b-a".replace("a", "x", 1) // "x-b-a" (limit replacements)
"a,b,c".split(",")           // ["a", "b", "c"]
"hello world".split()        // ["hello", "world"] (whitespace)
"hi".contains("h")           // true
len("hello")                 // 5
"hello".length()             // 5
"hello".to_int()             // error; use on numeric strings
"42".to_int()                // 42
"ff".to_int(16)              // 255 (with base)
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
    println(str(i));
}

// for-in over arrays/slices
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
    println(str(step));
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

// match as expression
let label = match day {
    1 => "Monday";
    2 => "Tuesday";
    _ => "Other";
};

// break and continue work in all loops
for var i = 0; i < 10; i++ {
    if i == 5 { break; }
    if i % 2 == 0 { continue; }
    println(str(i));
}
```

### Arrays (Slices)

```zenth
let numbers = [1, 2, 3, 4, 5];
let names: array(str) = [];       // empty with type annotation

// Access and length
println(str(numbers[0]));
println(str(len(numbers)));
println(str(numbers.length()));

// Mutating (requires var)
var items: array(int) = [];
items.add(10);         // append
items.push(20);        // same as add
let last = items.pop();     // remove and return last
let at = items.pop(0);     // remove and return at index

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

// Type conversions for string arrays
let strs = ["1", "2", "3"];
let ints = strs.to_int();     // [1, 2, 3]
let floats = strs.to_f64();   // [1.0, 2.0, 3.0]
```

### Tuples

```zenth
// Positional tuple
let pair = tuple("hello", 42);
println(pair.0);           // "hello"
println(str(pair.1));      // "42"

// Named tuple
let person = tuple(name="Alice", age=30);
println(person.name);      // "Alice"
println(str(person.age));  // "30"
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
println(str(scores["Alice"]));
println(str(len(scores)));

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

// Object keys (value-based lookup)
obj Point {
    x: int;
    y: int;
}
var grid = hashmap(Point, str);
grid[Point(x=1, y=2)] = "#";
println(grid[Point(x=1, y=2)]);  // "#" (same-value lookup works)
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
println(str(r.area()));    // "50"
println(str(r));           // "Rect(10x5)"

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

### Imports

```zenth
import "fmt";
import "math" as m;
```

Available modules map to Go stdlib: `fmt`, `math`, `os`, `strings`

### Built-in Functions

| Function | Description |
|---|---|
| `print(...)` | Print without newline |
| `println(...)` | Print with newline |
| `len(x)` | Length of string, array, or hashmap |
| `str(x)` | Convert any value to string |
| `int(x)` | Convert string or float to int |
| `f64(x)` | Convert to float64 |
| `abs(x)` | Absolute value (int or float) |
| `min(a, b)` | Minimum of two values |
| `max(a, b)` | Maximum of two values |
| `clamp(x, lo, hi)` | Clamp value to range |
| `round(x)` | Round float to nearest int |
| `floor(x)` | Floor of float |
| `ceil(x)` | Ceiling of float |
| `pow(base, exp)` | Exponentiation |
| `sqrt(x)` | Square root |
| `range(start, end)` | Exclusive range as array(int) |
| `range(start, end, step)` | Exclusive range with step |
| `rangei(start, end)` | Inclusive range |
| `rangei(start, end, step)` | Inclusive range with step |
| `hashmap(K, V)` | Create empty hashmap |
| `hashmap(K, V, default=val)` | Create hashmap with default |
| `tuple(...)` | Create tuple (positional or named) |
| `file(path)` | Create file handle |
| `flag(default=val)` | Command-line flag (name inferred from variable) |
| `exit(code)` | Exit program with status code |

### File I/O

```zenth
let f = file("data.txt");

if f.exists() {
    let content = f.read();       // entire file as str
    let lines = f.lines();        // array(str)
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
var a, b, c = 1, 2, 3;

// Swap
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
            0 => "FizzBuzz";
            _ => match i % 3 {
                0 => "Fizz";
                _ => match i % 5 {
                    0 => "Buzz";
                    _ => str(i);
                };
            };
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
        println(str(v));
    }
}
```

## Important Notes

- Every program needs a `fn main() { ... }` entry point
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

## What's Not Yet Implemented

- Enums / tagged unions
- Interfaces (parsed but not type-checked)
- Error handling (`try`)
- Multiple return values
- Concurrency / goroutines
- Generics
- Multi-file projects / user-defined packages
