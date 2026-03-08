# <img src="extras/branding/zenth-enso-trans.png" alt="Zenth logo" width="48" /> Zenth

A programming language built to match my personal tastes.

Zenth takes the parts I like from Python and Go and mashes them up into something hopefully consistent and simple. I like Python's readability and ease of use, and Go's compilation to a binary and toolset. For example, Zenth uses curly braces and semicolons instead of whitespace to define structure, requires types to be used (using inference where it can), and simplifies object creation and usage.

It compiles to native binaries using Go's toolchain, so gets the performance, garbage collection, and cross-compilation benefits of Go.

**Note:** This is a personal project. It's fun to hack on and pleasant to write small programs in, but probably shouldn't be used for anything too serious.

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

You need [Go](https://go.dev/dl/) 1.26+ and optionally [just](https://github.com/casey/just).

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
- **Objects with methods** -- define methods inside the object, `self` pre-declated
- **Useful object printing** -- objects print as constructor-like values by default
- **Typed hashmaps** -- create hashmaps with `hashmap(KeyType, ValueType)`
- **Tuples** -- fixed-size mixed values like `("a", 1)` with `.0`, `.1` access
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

### Vim 

Vim/Neovim syntax highlighting is included in `extras/vim/`. To install, copy or symlink it into your Vim config:

```sh
cp -r extras/vim/* ~/.vim/
# or for Neovim
cp -r extras/vim/* ~/.config/nvim/
```

This gives you syntax highlighting, filetype detection, and indentation for `.zn` files.

### LLMs

A skill is available at `extras/skill` to provide guidance to your LLM on how to
write Zenth programs. Copy the `skill` directory to your providers skills
folder, renaming to `zenth`. For Claude Code:

```sh
cp -r extras/skill ~/.claude/skills/zenth
```


## License

MIT
