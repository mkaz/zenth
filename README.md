# Zenth

A programming language built to match my personal tastes.

Zenth takes the parts I like from Python and Go -- Python's readability and ease of use, Go's type safety and compiled speed -- and combines them into something consistent and simple. Curly braces and semicolons for structure (no whitespace debates), type inference where it's obvious, explicit types where it matters.

It compiles to native binaries through Go, so you get real performance, garbage collection, and cross-compilation without thinking about it.

This is a personal project. It's fun to hack on and pleasant to write small programs in, but probably shouldn't be used for anything too serious.

## What it looks like

```zenth
fn main() {
    let name = "World";
    println("Hello, {name}!");

    let numbers = [1, 2, 3, 4, 5];
    for i, n in numbers {
        println("{i}: {n}");
    }
}
```

A few things to notice: `let` for immutable variables, string interpolation with `{}`, and `for-in` that just works.

## A slightly bigger example

```zenth
obj Circle {
    radius: f64;

    fn area() -> f64 {
        return 3.14159 * self.radius * self.radius;
    }
}

fn fizzbuzz(n: int) -> str {
    return if n % 15 == 0 { "FizzBuzz" }
      else if n % 3 == 0  { "Fizz" }
      else if n % 5 == 0  { "Buzz" }
      else                 { str(n) };
}

fn main() {
    let c = Circle(radius=5.0);
    println("Area: {c.area()}");

    for i in range(1, 21) {
        println(fizzbuzz(i));
    }
}
```

## Quick start

You need [Go](https://go.dev/dl/) 1.21+ and optionally [just](https://github.com/casey/just).

```sh
git clone https://github.com/mkaz/zenth.git
cd zenth
just build
```

Write a program:

```zenth
fn main() {
    println("Hello from Zenth!");
}
```

Run it:

```sh
zenth run hello.zn
```

Or compile to a binary:

```sh
zenth build hello.zn
./hello
```

## Key ideas

- **Immutable by default** -- `let` is immutable, `var` is mutable, caught at compile time
- **Type inference** -- locals are inferred, function signatures are explicit
- **One loop keyword** -- `for` does C-style, for-in, for-range-count, while-style, and infinite loops
- **Objects with methods** -- define methods inside the object, use `self`
- **`match` over switch** -- with `_` as the wildcard
- **Compiles to native binaries** -- via Go, so it's fast and portable

## Documentation

See the [docs/](docs/README.md) for the full language reference:

- [Getting Started](docs/getting-started.md)
- [Variables](docs/variables.md)
- [Types](docs/types.md)
- [Functions](docs/functions.md)
- [Strings](docs/strings.md)
- [Control Flow](docs/control-flow.md)
- [Objects](docs/objects.md)
- [Slices](docs/slices.md)
- [Files](docs/files.md)
- [Imports](docs/imports.md)

## Editor support

Vim/Neovim syntax highlighting is included in `extras/vim/`. To install, copy or symlink it into your Vim config:

```sh
cp -r extras/vim/* ~/.vim/
# or for Neovim
cp -r extras/vim/* ~/.config/nvim/
```

This gives you syntax highlighting, filetype detection, and indentation for `.zn` files.

## Status

Zenth is early and evolving. The basics work well -- functions, variables, objects, control flow, type checking, slices, file I/O -- but features like enums, interfaces, error handling, and multi-file projects aren't there yet. See the [docs](docs/README.md) for what's currently supported.

## License

MIT
