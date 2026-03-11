# Strings

Zenth has two kinds of string literals: double-quoted strings with interpolation, and single-quoted raw strings.

## Double-Quoted Strings

Double-quoted strings (`"..."`) support interpolation with `{expression}`:

```zenth
let name = "World";
println("Hello, {name}!");  // Hello, World!
```

Any expression can go inside the braces -- variables, arithmetic, method calls:

```zenth
let a = 3;
let b = 4;
println("{a} + {b} = {a + b}");  // 3 + 4 = 7
```

Values are automatically converted to strings, so you don't need to call `str()` inside interpolation:

```zenth
let count = 42;
println("There are {count} items");  // There are 42 items
```

### Method Calls in Interpolation

You can call methods and access fields inside `{}`:

```zenth
let r = Rectangle{ width: 5.0, height: 10.0 };
println("Area: {r.area()}");  // Area: 50
```

### Plain Strings

If a double-quoted string has no `{`, it behaves as a plain string with no special processing:

```zenth
println("No interpolation here");
```

### Escape Sequences

| Escape | Character |
|--------|-----------|
| `\n` | Newline |
| `\t` | Tab |
| `\\` | Backslash |
| `\"` | Double quote |
| `\{` | Literal `{` (prevents interpolation) |
| `\}` | Literal `}` |
| `\0` | Null byte |

To include a literal brace in a double-quoted string, escape it:

```zenth
println("Use \{braces\} literally");  // Use {braces} literally
```

## Single-Quoted Strings (Raw Strings)

Single-quoted strings (`'...'`) never interpolate. Braces are treated as plain characters:

```zenth
println('Hello {name}');  // Hello {name}
```

This is useful for strings that contain braces as literal text, such as templates, format strings, or code snippets.

Single-quoted strings support the same escape sequences as double-quoted strings, except `\'` instead of `\"`:

| Escape | Character |
|--------|-----------|
| `\n` | Newline |
| `\t` | Tab |
| `\\` | Backslash |
| `\'` | Single quote |
| `\0` | Null byte |

## String Concatenation

Strings can be concatenated with `+`:

```zenth
let first = "Hello";
let last = "World";
let full = first + ", " + last;
```

With interpolation, concatenation is rarely needed:

```zenth
let full = "{first}, {last}";
```

## String Conversion

Use `str()` to convert other types to strings:

```zenth
let n = 42;
let s = str(n);    // "42"
let f = str(3.14); // "3.14"
```

Inside interpolated strings, conversion is automatic -- `str()` is only needed when you need a string value outside of interpolation.

### Parsing with Base

Use `.to_int(base)` to parse a string as an integer in a given base (2-36):

```zenth
let n = "ff".to_int(16);     // 255
let b = "1010".to_int(2);    // 10
let o = "77".to_int(8);      // 63
```

Without an argument, `.to_int()` parses as base 10.

## String Length

Use `len()` or `.length()` to get the length of a string:

```zenth
let s = "hello";
println(len(s));           // 5
println(s.length());       // 5
```

## Containment

Use `.contains()` to check if a string contains a substring:

```zenth
let s = "hello world";
if s.contains("world") {
    println("found it");
}
```

## Case Conversion

Use `.upper()` and `.lower()` for case conversion:

```zenth
println("hello".upper()); // HELLO
println("HeLLo".lower()); // hello
```

## Prefix and Suffix

Use `.starts_with()` and `.ends_with()` to check for prefixes and suffixes:

```zenth
let name = "zenth";
println(name.starts_with("zen"));      // true
println(name.ends_with("th"));         // true
```

Use `.strip_prefix()` and `.strip_suffix()` to remove a prefix or suffix from a string. If the string does not have the given prefix/suffix, it is returned unchanged:

```zenth
let line = "fold along x=5";
let spec = line.strip_prefix("fold along "); // "x=5"

let file = "photo.png";
let name = file.strip_suffix(".png"); // "photo"

// No match — returns unchanged
let same = "hello".strip_prefix("xyz"); // "hello"
```

## Trimming

Use `.strip()` to trim whitespace, or pass characters to trim from both ends:

```zenth
println("--" + "  hi  ".strip() + "--");   // --hi--
println("--" + "..hi..".strip(".") + "--"); // --hi--
```

## Find, Count, Replace

Use `.find()` and `.count()` for substring search and counting:

```zenth
println("banana".find("na"));       // 2
println("banana".count("na"));      // 2
```

Use `.replace(old, new)` to replace all matches, or `.replace(old, new, n)` to limit replacements:

```zenth
println("a-b-a".replace("a", "x"));    // x-b-x
println("a-b-a".replace("a", "x", 1)); // x-b-a
```

## Splitting Strings

Use the `.split()` method to divide a string into a slice of substrings (`array(str)`). 
By default, `.split()` splits on whitespace. You can also pass a delimiter string:

```zenth
let words = "this is zenth".split();
println(words[0]); // "this"

let parts = "a-b-c".split("-");
println(parts[1]); // "b"
```

Use `.split_once(sep)` when you only need two pieces. It returns a tuple `(left, right)`:

```zenth
let (cmd, value) = "forward 10".split_once(" ");
println(cmd);   // "forward"
println(value); // "10"
```

## Repeating

Use `.repeat(n)` to repeat a string `n` times:

```zenth
let dots = ".".repeat(5);   // "....."
let row = "ab".repeat(3);   // "ababab"
let empty = "x".repeat(0);  // ""
```

This is useful for building grid rows, padding, or repeated patterns.

## Indexing

Indexing a string returns a one-character `str`:

```zenth
let s = "hello";
let first = s[0];  // "h"
```

## Slicing

Use slice syntax to extract a substring:

```zenth
let s = "hello";
let tail = s[1:];    // "ello"
let head = s[:3];    // "hel"
let mid = s[1:4];    // "ell"
```

String slicing is Unicode-aware (operates on runes, not bytes).
