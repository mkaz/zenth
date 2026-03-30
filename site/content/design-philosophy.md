---
title: "Design Philosophy"
weight: 4
aliases:
  - /explanation/design-philosophy/
---

# Design Philosophy

Zenth is an attempt to combine a few traits that are often split across different languages. Python is easy to read and write, but it trades away static types and native binaries. Go has both, but its syntax can feel heavier than I want for day-to-day code. Zenth takes ideas from each and aims for a smaller, more regular surface area.

## Simplicity first

Each feature should have one primary way to use it. There is one loop keyword (`for`) that covers counting, iteration, and while-style loops. There is one way to define a type (`obj`). There are no alternate spellings for the same construct.

When considering new features, the bias is toward reusing what already exists instead of introducing parallel forms.

## Strongly typed with inference

Function signatures require explicit types. That keeps function boundaries easy to read and makes it clear what a piece of code accepts and returns.

Inside functions, local variables are inferred. You write `let x = 42;` instead of `let x: Int = 42;`. The compiler already knows the type, so the local annotation usually does not add much.

```zenth
fn add(a: Int, b: Int) -> Int {   // explicit at boundaries
    let result = a + b;           // inferred locally
    return result;
}
```

## Capitalized built-ins

Built-in functions are capitalized so they stand out from user-defined names in code. It also means common names remain available for local variables without colliding with the built-in.

```zenth
let max = Max(a, b);
```

That keeps the built-in visually distinct while still allowing `max` to be used as an ordinary variable name.

## No whitespace sensitivity

Zenth uses curly braces and semicolons. This is a deliberate departure from Python. Whitespace-sensitive syntax can make pasted code, indentation changes, and generated code more fragile. Braces and semicolons keep structure explicit regardless of formatting.

## Compiles to native binaries

Zenth transpiles to Go and then compiles with Go's toolchain. That gives you:

- **Native binaries** — ship a single file, no runtime required
- **Garbage collection** — no manual memory management
- **Cross-compilation** — build for any platform Go supports
- **Performance** — Go's compiler produces fast code with a low memory footprint

Transpilation also means Zenth can lean on Go's standard library for basic operations. File I/O, string handling, and math build on packages that are already mature and widely used.

## Borrow the best ideas

Zenth does not try to introduce entirely new concepts. Instead, it takes patterns that have worked well elsewhere and puts them under a more consistent syntax:

- `if`/`else` as expressions (from Rust)
- `match` for both statements and values (from Rust and Kotlin)
- A single `for` keyword for every loop style (from Go)
- String interpolation with `{variable}` (from Python f-strings, simplified)
- Named and default function arguments (from Python)

The goal is a language that feels familiar enough to start quickly, while still being consistent enough to stay readable as programs grow.
