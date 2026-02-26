# Variables

Zenth has three ways to declare variables: `let`, `var`, and `const`. All declarations use `=` for assignment -- there is no `:=` operator.

## Immutable Variables with `let`

`let` declares an immutable binding. Once assigned, its value cannot change:

```zenth
let name = "Zenth";
let year = 2026;
let pi = 3.14159;
let active = true;
```

Attempting to reassign a `let` variable is a compile-time error:

```zenth
let x = 5;
x = 10;  // error: cannot assign to immutable variable 'x'
```

## Mutable Variables with `var`

`var` declares a mutable binding that can be reassigned:

```zenth
var counter = 0;
counter = counter + 1;
counter += 1;
println(str(counter));  // 2
```

## Constants with `const`

`const` declares a named constant:

```zenth
const PI = 3.14159;
const MAX_SIZE = 1024;
```

## Type Inference

Types are inferred from the assigned value by default:

```zenth
let x = 42;        // int
let y = 3.14;      // f64
let s = "hello";   // str
let b = true;      // bool
```

## Explicit Type Annotations

You can annotate the type explicitly with `: Type`:

```zenth
let x: int = 5;
var name: str = "Zenth";
let ratio: f64 = 0.75;
```

Type annotations are required on function parameters and return types, but optional on local variables.

## Compound Assignment

Mutable variables support compound assignment operators:

```zenth
var n = 10;
n += 5;   // n is now 15
n -= 3;   // n is now 12
n *= 2;   // n is now 24
n /= 4;   // n is now 6
```

## Multiple Assignment

You can assign to multiple variables in a single statement. This is often used to swap values without needing a temporary variable:

```zenth
var x = 1;
var y = 2;
x, y = y, x;  // x becomes 2, y becomes 1
```

## Tuple Destructuring

You can destructure tuples directly in `let`, `var`, and `const` declarations:

```zenth
let (cmd, raw) = "forward 10".split_once(" ");
let value = int(raw);
```

Destructuring arity must match the tuple size:

```zenth
let (a, b) = tuple(1, 2);      // ok
let (x, y, z) = tuple(1, 2);   // error
```

## Increment and Decrement

```zenth
var i = 0;
i++;  // i is now 1
i--;  // i is now 0
```

## Summary

| Keyword | Mutable | Use case |
|---------|---------|----------|
| `let`   | No      | Most variables -- prefer immutable by default |
| `var`   | Yes     | Counters, accumulators, anything that changes |
| `const` | No      | Named constants |
