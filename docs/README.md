# Zenth Language Documentation

Zenth is a compiled programming language that blends Go's type discipline with Python's simplicity. It uses curly braces and semicolons for structure, and compiles to native binaries.

```zenth
fn main() {
    let name = "Zenth";
    Println("Hello from {name}!");
}
```

## Documentation

- [Getting Started](getting-started.md) -- Installation, your first program, and building
- [Variables](variables.md) -- `let`, `var`, `const`, and type inference
- [Types](types.md) -- Built-in types and type system
- [Functions](functions.md) -- Declaring and calling functions
- [Strings](strings.md) -- String literals, interpolation, and raw strings
- [Control Flow](control-flow.md) -- `if`/`else`, `for` loops, and `match`
- [Enums](enums.md) -- Named value types with fixed variants
- [Objects](objects.md) -- Objects, methods, and `self`
- [Arrays](arrays.md) -- Dynamic arrays
- [Tuples](tuples.md) -- Fixed-size mixed-type values
- [Hashmaps](hashmaps.md) -- Key-value collections
- [Sets](sets.md) -- Unordered collections of unique values
- [Dates](dates.md) -- Calendar dates, parsing, and formatting
- [Files](files.md) -- Reading files and file paths
- [Imports](imports.md) -- Using standard library modules
- [Modules](modules.md) -- Splitting programs across multiple files
- [Testing](testing.md) -- Writing and running tests

## Design Philosophy

- **Simplicity first** -- every feature should have one clear way to use it, with minimal ceremony
- **Borrow the best ideas** -- `if`/`else` as expressions, `match` for both statements and values, a single `for` keyword for every loop style -- proven patterns from Rust, Python, and Go, unified under a consistent syntax
- **Immutable by default** -- `let` is immutable, `var` opts into mutation; safe defaults, easy overrides
- **Strongly typed with inference** -- function signatures require explicit types, locals are inferred; clarity at boundaries, convenience everywhere else
- **No whitespace sensitivity** -- curly braces `{}` and semicolons `;` define structure unambiguously
- **Compiles to native binaries** -- Go as a backend gives you a GC, a runtime, and cross-compilation for free
