# Imports

Zenth provides access to standard library functionality through imports. Under the hood, these map to Go standard library packages.

## Basic Import

```zenth
import "fmt";
```

## Aliased Import

Use `as` to give an import a shorter name:

```zenth
import "math" as m;
```

Then call functions through the alias:

```zenth
let root = m.sqrt(16.0);
```

## Available Modules

| Zenth module | Maps to | Common functions |
|-------------|---------|-----------------|
| `fmt` | Go `fmt` | `println`, `printf`, `sprintf` |
| `math` | Go `math` | `sqrt`, `abs`, `pow`, `min`, `max` |
| `os` | Go `os` | `exit`, `args` |
| `strings` or `str` | Go `strings` | `split`, `join`, `contains`, `replace`, `trim` |

## Calling Module Functions

```zenth
import "math";

fn main() {
    let x = math.sqrt(144.0);
    println(str(x));  // 12

    let y = math.pow(2.0, 10.0);
    println(str(y));  // 1024
}
```

## Example with String Operations

```zenth
import "strings";

fn main() {
    let csv = "apple,banana,cherry";
    let parts = strings.split(csv, ",");
    for part in parts {
        println(part);
    }
}
```

## Note on Built-ins

The functions `print`, `println`, `len`, and `str` are built-in and always available -- no import needed:

```zenth
fn main() {
    println("no import required");
    println(str(len("hello")));
}
```
