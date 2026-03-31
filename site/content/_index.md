---
title: "Zenth"
type: docs
---

# Zenth

A compiled programming language that blends Go's type discipline with Python's simplicity.

Zenth transpiles to Go and compiles to native binaries, giving you Go's performance, garbage collection, and cross-compilation while writing in a syntax designed for readability.

```zenth
fn main() {
    let name = "World";
    Println("Hello, {name}!");

    let numbers = [1, 2, 3, 4, 5];
    for i, n in numbers {
        Println("{i}: {n}");
    }
}
```

```zenth
obj Circle {
    radius: Float;

    fn area() -> Float {
        return 3.14159 * self.radius * self.radius;
    }
}

fn fizzbuzz(n: Int) -> Str {
    return if n % 15 == 0 { "FizzBuzz" }
      else if n % 3 == 0  { "Fizz" }
      else if n % 5 == 0  { "Buzz" }
      else                 { Str(n) };
}

fn main() {
    let c = Circle(radius=5.0);
    Println("Area: {c.area()}");

    for i in Range(1, 21) {
        Println(fizzbuzz(i));
    }
}
```

## Documentation

This documentation is organized into practical guides, reference material, and a top-level page for the language's design philosophy:

- **[How-to Guides](howto/_index.md)** — Setup and task-focused guides, including getting started
- **[Reference](reference/_index.md)** — Complete descriptions of Zenth's syntax and features
- **[Design Philosophy](design-philosophy.md)** — Background and design decisions behind the language
