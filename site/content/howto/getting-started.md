---
title: "Getting Started"
weight: 1
aliases:
  - /tutorials/getting-started/
---

# Getting Started

## Prerequisites

Zenth compiles to native binaries through Go, so you need:

- [Go](https://go.dev/dl/) 1.26 or later
- [just](https://github.com/casey/just) command runner (optional, for convenience)

## Installation

Clone the repository and build the compiler:

```sh
git clone https://github.com/mkaz/zenth.git
cd zenth
just build
```

This produces a `zenth` binary in the current directory. To install, add it to your `$PATH`:

```sh
just install
```

## Your First Program

Create a file called `hello.zn`:

```zenth
fn main() {
    Println("Hello, World!");
}
```

Every Zenth program needs a `main` function as its entry point.

## Building and Running

**Compile and run in one step:**

```sh
zenth run hello.zn
```

**Compile to a binary:**

```sh
zenth build hello.zn
./hello
```

**Specify an output name:**

```sh
zenth build -o greet hello.zn
./greet
```

**View the generated Go source** (useful for debugging):

```sh
zenth build --emit-go hello.zn
```

## File Extension

Zenth source files use the `.zn` extension.

## A Slightly Bigger Example

Here's a Fibonacci program that shows functions, loops, and type annotations:

```zenth
fn fibonacci(n: Int) -> Int {
    if n <= 1 {
        return n;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
}

fn main() {
    for var i = 0; i < 15; i++ {
        Println(fibonacci(i));
    }
}
```

```sh
zenth run fibonacci.zn
```

## Next Steps

- [Variables](../reference/variables.md) — learn about `let`, `var`, and `const`
- [Functions](../reference/functions.md) — parameters, return types, and calling
- [Strings](../reference/strings.md) — interpolation and raw strings
