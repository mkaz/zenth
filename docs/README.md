# Zenth Language Documentation

Zenth is a compiled programming language that blends Go's type discipline with Python's simplicity. It uses curly braces and semicolons for structure, and compiles to native binaries.

```zenth
fn main() {
    let name = "Zenth";
    println("Hello from {name}!");
}
```

## Documentation

- [Getting Started](getting-started.md) -- Installation, your first program, and building
- [Variables](variables.md) -- `let`, `var`, `const`, and type inference
- [Types](types.md) -- Built-in types and type system
- [Functions](functions.md) -- Declaring and calling functions
- [Strings](strings.md) -- String literals, interpolation, and raw strings
- [Control Flow](control-flow.md) -- `if`/`else`, `for` loops, and `match`
- [Objects](objects.md) -- Objects, methods, and `self`
- [Slices](slices.md) -- Dynamic arrays
- [Files](files.md) -- Reading files and file paths
- [Imports](imports.md) -- Using standard library modules

## Design Philosophy

- **No whitespace sensitivity** -- curly braces `{}` and semicolons `;` define structure
- **Immutable by default** -- `let` is immutable, `var` is mutable
- **Strongly typed with inference** -- function signatures require explicit types, locals are inferred
- **One loop keyword** -- `for` handles C-style, for-in, while-style, and infinite loops
- **Compiles to native binaries** -- via Go as a backend, inheriting its GC and cross-compilation
