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

## String Length

Use `len()` to get the length of a string:

```zenth
let s = "hello";
println(str(len(s)));  // 5
```

## Indexing

Indexing a string returns a `byte`:

```zenth
let s = "hello";
let first = s[0];  // byte value of 'h'
```
