---
title: "Imports"
weight: 14
---

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
| `strings` or `Str` | Go `strings` | `split`, `join`, `contains`, `replace`, `trim` |

## Calling Module Functions

```zenth
import "math";

fn main() {
    let x = math.sqrt(144.0);
    Println(x);  // 12

    let y = math.pow(2.0, 10.0);
    Println(y);  // 1024
}
```

## Example with String Operations

```zenth
import "strings";

fn main() {
    let csv = "apple,banana,cherry";
    let parts = strings.split(csv, ",");
    for part in parts {
        Println(part);
    }
}
```

## External Go Package Imports

Use `import_go` to import any Go package — including third-party libraries. The Zenth toolchain automatically fetches the package using `go get` before compiling.

```zenth
import_go "github.com/some/library";
```

With an alias:

```zenth
import_go "github.com/some/library" as lib;
```

Calls on `import_go` packages bypass Zenth's type checker and pass through directly to Go. This means you get full access to any Go library, with errors surfacing at compile time from Go rather than from Zenth.

```zenth
import_go "database/sql";
import_go "github.com/mattn/go-sqlite3" as _;

fn main() {
    let db = sql.Open("sqlite3", "./data.db");
    let rows = db.Query("SELECT id, name FROM users");
    for rows.Next() {
        var id: Int = 0;
        var name = "";
        rows.Scan(&id, &name);
        Println("{id}: {name}");
    }
}
```

> **Note:** `import_go "github.com/mattn/go-sqlite3"` requires CGO (`gcc` must be installed). For a pure-Go SQLite option, use `modernc.org/sqlite` instead.

Return values from `import_go` function calls are untyped — you can assign them to variables and chain further calls, but Zenth will not validate the types.

## Local Module Imports

To import your own `.zn` files, use a path starting with `./` or `../`:

```zenth
import "./utils";
import "./math/geometry" as geo;
```

See [Modules](modules.md) for full documentation on splitting programs across multiple files.

## Note on Built-ins

The functions `Print`, `Println`, `Len`, and `Str` are built-in and always available -- no import needed:

```zenth
fn main() {
    Println("no import required");
    Println(Len("hello"));
}
```
