
# Zenth Programming Language

A compiled programming language that blends Go's type discipline with Python's simplicity, using curly braces and semicolons for structure. Compiles to native binaries.

## Instructions

For every change ALWAYS follow these directives:

- For every language change, update the documentation thoroughly on how to use
  the new feature, include plenty of examples.

- For every language change, update the vim syntax in `extras/vim/` when relevant
    - For example: if a new keyword, operator, or other language construct is added

- For every change, increment the patch version number in `cmd/zenth/main.go`
    - For example: a change will increment 0.1.4 to 0.1.5

- For every change, look at the `extras/skill/SKILL.md` skill file and when
  relevant update so LLMs working with Zenth will have any updated information
  on how to use the language.

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

## How to Add a Built-in Function

Built-in functions touch four compiler packages plus docs/tests. Here is the checklist:

### 1. Register in the checker (`pkg/checker/checker.go`)

In `New()`, add to the `c.funcs` map (around line 89–147):

```go
c.funcs["MyFunc"] = &FuncInfo{
    Name:        "MyFunc",
    Params:      []ZType{TypeStr, TypeInt},  // parameter types
    Return:      TypeStr,                     // return type
    NumRequired: 1,                           // for optional params (omit if all required)
}
```

This handles standard type-checking automatically via `checkArgs()`. For built-ins
that need special validation (like `Flag`, `Zip`, `Hashmap`), add a special-case
block in `checkCallExpr()` (around line 1830–1850) that calls a dedicated
`checkMyFuncCall()` method instead.

### 2. Generate Go code (`pkg/codegen/codegen.go`)

In `genCallExpr()`, the big switch on `ident.Name` (starts around line 1828),
add a case:

```go
case "MyFunc":
    g.imports["os"] = ""  // if you need an import
    g.write("os.Something(")
    g.genArgList(c.Args)
    g.write(")")
    return
```

If you need a Go helper function, add a `needsMyFunc bool` field to the
`Generator` struct (around line 30–74), set it in the case above, and emit
the helper in the helpers block (around line 550–700) following the existing
`if g.needsXxx { ... }` pattern.

### 3. AST annotations (only for complex built-ins)

If the checker needs to pass extra metadata to codegen (like `Flag` does with
`FlagName`/`FlagGoType`), add fields to `CallExpr` in `pkg/ast/ast.go`.
Simple built-ins that just map to a Go call don't need this.

### 4. Tests (`pkg/driver/driver_test.go`, `testdata/`)

- Add a success fixture: `testdata/myfunc.zn` with `fn main() { ... }`
- Add a test entry in `driver_test.go` in the success test table with `contains` strings
- Add an error fixture: `testdata/errors/myfunc_bad_arg.zn`
- Add an error test entry with `errMsg` substring

### 5. Docs and extras

- `docs/functions.md` — add to the built-in table and add a detailed section with examples
- `extras/vim/syntax/zenth.vim` — add to the `zenthBuiltin` keyword list
- `extras/skill/SKILL.md` — add to the built-in table and add usage examples
- `cmd/zenth/main.go` — bump patch version
- `cmd/zenth/main_test.go` — update version string in `TestVersion`

### Key patterns

- Parser has no special cases for built-ins — they parse as regular `CallExpr` nodes
- The checker's `c.funcs` map handles arg count/type validation for standard built-ins
- Codegen's `g.imports` map manages Go import statements; set `g.imports["pkg"] = ""` as needed
- Helper functions use `g.needsXxx` bools to emit only when used (avoids unused-code errors in generated Go)

## Coding Style & Naming Conventions

Use idiomatic Go and keep files `gofmt`-formatted (tabs, standard imports).
Package names are short lowercase nouns (`lexer`, `checker`). Exported identifiers use `CamelCase`; internals use `camelCase`.
Tests live next to source files as `*_test.go` and use `TestXxx` names.
For `.zn` fixture files, use descriptive lowercase names (for example `multiassign.zn`, `type_mismatch.zn`).

## Language Syntax Reference

### Variables

```zenth
let x = 5;              // immutable, type inferred
let x: Int = 5;         // immutable, explicit type
var counter = 0;         // mutable, type inferred
var counter: Int = 0;    // mutable, explicit type
const PI = 3.14159;      // constant
```

No `:=` operator. The `let`/`var`/`const` keyword signals a declaration, `=` is always the assignment operator.

### Functions

```zenth
fn add(a: Int, b: Int) -> Int {
    return a + b;
}

fn greet(name: Str) {
    Println("Hello, " + name);
}
```

### Objects and Methods (OOP-style)

Methods are defined inside the obj body. Use `self` to access fields and call other methods.

```zenth
obj Rectangle {
    width: Float;
    height: Float;

    fn area() -> Float {
        return self.width * self.height;
    }
}

let r = Rectangle(width=3.0, height=4.0);
Println(Str(r.area()));
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

// C-style for loop
for var i = 0; i < 10; i++ {
    Println(Str(i));
}

// for-in
for item in items {
    Println(Str(item));
}
for i, item in items {
    Println(Str(i));
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

Built-in: `Int`, `Float`, `Bool`, `Str`, `Byte`

Composite: `Array(Int)` (slice), obj literals `Point(x=1.0, y=2.0)`, `Tuple(Str, Int)`

### Built-in Functions

`Print(...)`, `Println(...)`, `Len(x)`, `Str(x)`, `Flag(default=val)`

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

- Use `justfile` as convenience tool to build
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

