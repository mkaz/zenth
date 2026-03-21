# Strings

Zenth has two kinds of string literals: double-quoted strings with interpolation, and single-quoted raw strings.

## Double-Quoted Strings

Double-quoted strings (`"..."`) support interpolation with `{expression}`:

```zenth
let name = "World";
Println("Hello, {name}!");  // Hello, World!
```

Any expression can go inside the braces -- variables, arithmetic, method calls:

```zenth
let a = 3;
let b = 4;
Println("{a} + {b} = {a + b}");  // 3 + 4 = 7
```

Values are automatically converted to strings, so you don't need to call `Str()` inside interpolation:

```zenth
let count = 42;
Println("There are {count} items");  // There are 42 items
```

### Method Calls in Interpolation

You can call methods and access fields inside `{}`:

```zenth
let r = Rectangle{ width: 5.0, height: 10.0 };
Println("Area: {r.area()}");  // Area: 50
```

### Plain Strings

If a double-quoted string has no `{`, it behaves as a plain string with no special processing:

```zenth
Println("No interpolation here");
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
Println("Use \{braces\} literally");  // Use {braces} literally
```

## Single-Quoted Strings (Raw Strings)

Single-quoted strings (`'...'`) never interpolate. Braces are treated as plain characters:

```zenth
Println('Hello {name}');  // Hello {name}
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

## Multi-line Strings

Use triple double-quotes (`"""..."""`) for multi-line string literals. The content between the opening and closing `"""` is preserved as-is, including newlines and indentation:

```zenth
let text = """
Hello
World
""";
Println(text);  // prints "Hello\nWorld"
```

The first newline after the opening `"""` and the last newline before the closing `"""` are stripped, so the content is exactly what appears between them.

Multi-line strings support interpolation, just like regular double-quoted strings:

```zenth
let name = "Zenth";
let msg = """
Welcome to {name}!
Enjoy coding.
""";
```

They also support escape sequences (`\n`, `\t`, `\\`, `\"`, `\{`, `\}`, `\0`).

Multi-line strings are useful for help text, templates, and any string that spans multiple lines:

```zenth
Print("""
tasks — a pipe-delimited todo manager

Usage:
  tasks                    list active tasks
  tasks add 'title'        add a new task
  tasks done <id>          mark a task as done
""");
```

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

Use `Str()` to convert other types to strings:

```zenth
let n = 42;
let s = Str(n);    // "42"
let f = Str(3.14); // "3.14"
```

Inside interpolated strings, conversion is automatic -- `Str()` is only needed when you need a string value outside of interpolation.

### Parsing with Base

Use `.to_int(base)` to parse a string as an integer in a given base (2-36):

```zenth
let n = "ff".to_int(16);     // 255
let b = "1010".to_int(2);    // 10
let o = "77".to_int(8);      // 63
```

Without an argument, `.to_int()` parses as base 10.

## String Length

Use `Len()` or `.length()` to get the length of a string:

```zenth
let s = "hello";
Println(Len(s));           // 5
Println(s.length());       // 5
```

## Containment

Use `.contains()` to check if a string contains a substring:

```zenth
let s = "hello world";
if s.contains("world") {
    Println("found it");
}
```

## Case Conversion

Use `.upper()` and `.lower()` for case conversion:

```zenth
Println("hello".upper()); // HELLO
Println("HeLLo".lower()); // hello
```

## Prefix and Suffix

Use `.starts_with()` and `.ends_with()` to check for prefixes and suffixes:

```zenth
let name = "zenth";
Println(name.starts_with("zen"));      // true
Println(name.ends_with("th"));         // true
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
Println("--" + "  hi  ".strip() + "--");   // --hi--
Println("--" + "..hi..".strip(".") + "--"); // --hi--
```

## Find, Count, Replace

Use `.find()` and `.count()` for substring search and counting:

```zenth
Println("banana".find("na"));       // 2
Println("banana".count("na"));      // 2
```

Use `.replace(old, new)` to replace all matches, or `.replace(old, new, n)` to limit replacements:

```zenth
Println("a-b-a".replace("a", "x"));    // x-b-x
Println("a-b-a".replace("a", "x", 1)); // x-b-a
```

## Splitting Strings

Use the `.split()` method to divide a string into an array of substrings (`Array(Str)`). 
By default, `.split()` splits on whitespace. You can also pass a delimiter string:

```zenth
let words = "this is zenth".split();
Println(words[0]); // "this"

let parts = "a-b-c".split("-");
Println(parts[1]); // "b"
```

Use `.split_once(sep)` when you only need two pieces. It returns a tuple `(left, right)`:

```zenth
let (cmd, value) = "forward 10".split_once(" ");
Println(cmd);   // "forward"
Println(value); // "10"
```

## Character Classification

Use `.is_digit()` to check if a single-character string is a decimal digit (`0`–`9`):

```zenth
Println("5".is_digit());   // true
Println("0".is_digit());   // true
Println("a".is_digit());   // false
Println("[".is_digit());   // false
```

This is useful when iterating over a string character by character:

```zenth
var digits = 0;
for ch in "abc123" {
    if ch.is_digit() {
        digits += 1;
    }
}
Println(digits);  // 3
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

Indexing a string returns a one-character `Str`:

```zenth
let s = "hello";
let first = s[0];  // "h"
```

## Slicing

Use slicing syntax to extract a substring:

```zenth
let s = "hello";
let tail = s[1:];    // "ello"
let head = s[:3];    // "hel"
let mid = s[1:4];    // "ell"
```

String slicing is Unicode-aware (operates on runes, not bytes).
