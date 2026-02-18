# Zenth Programming Language

A compiled programming language that blends Go's type discipline with Python's simplicity, using curly braces and semicolons for structure. Compiles to native binaries.

## Instructions

For every change ALWAYS follow these directives:

- For every language change, update the vim syntax in `extras/vim/` when relevant
    - For example: if a new keyword, operator, or other language construct is added

- For every change, increment the patch version number in `cmd/zenth/main.go`
    - For example: a change will increment 0.1.4 to 0.1.5

## Build, Test, and Development Commands

- `just build`: build CLI binary at `./zenth` (`go build -o zenth ./cmd/zenth`).
- `just test`: run all Go tests (`go test ./...`), including fixture-driven driver tests.
- `just install`: install CLI to your Go bin path (`go install ./cmd/zenth`).
- `just clean`: remove local binary (`rm -f zenth`).
- `./zenth run testdata/hello.zn`: compile and execute a sample program.
- `./zenth build -o hello testdata/hello.zn`: compile fixture to a named binary.

## Architecture

Zenth is a transpiler: `.zn` source → Go source → `go build` → native binary. This gives us Go's GC, runtime, and cross-compilation for free.

```
source.zn → Lexer → Parser → Type Checker → Go Codegen → go build → native binary
```

The compiler is written in Go. Each stage lives in its own package under `pkg/`.

## Project Structure

```
cmd/zenth/          CLI entry point (build, run, version)
pkg/token/          Token types and position tracking
pkg/lexer/          Hand-written lexer
pkg/parser/         Recursive descent parser with Pratt expression parsing
pkg/ast/            AST node definitions
pkg/checker/        Type checker and semantic analysis (scope, mutability, types)
pkg/codegen/        Go source code generator
pkg/driver/         Build orchestration (temp dir, go build, cleanup)
testdata/           Example .zn programs (also used as integration tests)
testdata/errors/    Programs that must fail to compile (negative tests)
docs/               Language reference; keep aligned with code changes
extras/vim/         Vim/Neovim syntax, ftdetect, and indent support
```

## Coding Style & Naming Conventions

Use idiomatic Go and keep files `gofmt`-formatted (tabs, standard imports).
Package names are short lowercase nouns (`lexer`, `checker`). Exported identifiers use `CamelCase`; internals use `camelCase`.
Tests live next to source files as `*_test.go` and use `TestXxx` names.
For `.zn` fixture files, use descriptive lowercase names (for example `multiassign.zn`, `type_mismatch.zn`).

## Language Syntax Reference

### Variables

```zenth
let x = 5;              // immutable, type inferred
let x: int = 5;         // immutable, explicit type
var counter = 0;         // mutable, type inferred
var counter: int = 0;    // mutable, explicit type
const PI = 3.14159;      // constant
```

No `:=` operator. The `let`/`var`/`const` keyword signals a declaration, `=` is always the assignment operator.

### Functions

```zenth
fn add(a: int, b: int) -> int {
    return a + b;
}

fn greet(name: str) {
    println("Hello, " + name);
}
```

### Objects and Methods (OOP-style)

Methods are defined inside the obj body. Use `self` to access fields and call other methods.

```zenth
obj Rectangle {
    width: f64;
    height: f64;

    fn area() -> f64 {
        return self.width * self.height;
    }
}

let r = Rectangle{ width: 3.0, height: 4.0 };
println(str(r.area()));
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

// C-style for loop
for var i = 0; i < 10; i++ {
    println(str(i));
}

// for-in
for item in items {
    println(item);
}
for i, item in items {
    println(str(i));
}

// while-style
for running {
    process();
}

// infinite loop
for {
    listen();
}

// match
match day {
    1 => println("Monday");
    2 => println("Tuesday");
    _ => println("Other");
}
```

### Types

Built-in: `int`, `i8`, `i16`, `i32`, `i64`, `u8`, `u16`, `u32`, `u64`, `f32`, `f64`, `bool`, `str`, `byte`

Composite: `[]int` (slice), obj literals `Point{ x: 1.0, y: 2.0 }`

### Built-in Functions

`print(...)`, `println(...)`, `len(x)`, `str(x)`, `flag(default=val)`

### Imports

```zenth
import "fmt";
import "math" as m;
```

Standard library modules map to Go stdlib: `fmt`, `math`, `os`, `strings`/`str`

## Design Decisions

- **No whitespace sensitivity** — curly braces `{}` and semicolons `;` for structure
- **`let` is immutable, `var` is mutable** — catches accidental mutation at compile time
- **Simple `=` assignment** — no `:=`. `let`/`var`/`const` keyword is the declaration signal
- **Strongly typed with inference** — function signatures require explicit types, locals are inferred
- **OOP-style methods** — methods defined inside the obj body, use `self` to access fields
- **`for` is the only loop** — three styles: C-style, for-in, while-style
- **`match` instead of switch** — with `_` as default/wildcard
- **File extension** — `.zn`
- **Compile-time error checking** — type mismatches, mutability violations, undefined variables all caught before codegen
- **Transpile to Go** — inherits GC, goroutines, cross-compilation; `go build` produces the final native binary

## Conventions

- Use `justfile` (not Makefile)
- Obj fields are lowercase in Zenth, exported (capitalized) in generated Go
- Method names are lowercase in Zenth, exported in generated Go
- `self` is the implicit receiver name for methods

## Testing

- Unit tests: `pkg/lexer/lexer_test.go`
- Integration tests: `pkg/driver/driver_test.go` — builds and runs every `testdata/*.zn` file
- Negative tests: `testdata/errors/` — programs that must fail with specific error messages
- Run all: `just test`
- When adding language features, include at least one success fixture and one error case when applicable

## Commit & Pull Request Guidelines

Recent history uses short, imperative commit subjects (for example `Add file docs`, `Update vim instructions`).
Prefer one focused change per commit and keep subject lines under ~72 chars.
PRs should include: purpose, key implementation notes, tests run (`just test`), and docs updates for language/user-facing behavior changes.

## Current Status (v0.1)

**Working:** functions, let/var/const, type inference, objects (`obj`) with OOP-style methods and `self`, if/else/else-if, for loops (C-style, for-in, while-style, infinite), match statements, compile-time type errors, mutability checking, string/int/float/bool types, slices, obj literals, break/continue, compound assignment (`+=`, `-=`, etc.), `++`/`--`, imports, multi-variable assignment.

**Not yet implemented:** enums/tagged unions, interfaces (parsed but not checked), error handling (`try`), multi-return, concurrency, generics, multi-file projects, user-defined packages.
