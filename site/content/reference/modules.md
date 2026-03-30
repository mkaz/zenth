---
title: "Modules"
weight: 15
---

Zenth supports splitting programs across multiple files using local modules. A module is simply a `.zn` File (or a directory of `.zn` files) that you import into your program.

## Importing a Local Module

Use a path starting with `./` or `../` to import a local module:

```zenth
import "./utils";
import "./math/geometry";
import "../shared/helpers";
```

The module name is derived from the filename stem (the last path component without `.zn`). Use `as` to give it a custom alias:

```zenth
import "./mathutils" as math;
import "../geometry/shapes" as shapes;
```

## Using Module Functions

Call functions and access types from a module through its name (or alias):

```zenth
import "./mathutils";

fn main() {
    let result = mathutils.add(3, 4);
    Println(Str(result));
}
```

With an alias:

```zenth
import "./mathutils" as math;

fn main() {
    let result = math.add(3, 4);
    Println(Str(result));
}
```

## Module Files

A module file is a regular `.zn` file that **does not define `fn main`**. It can contain functions, objects, constants, and type aliases, all of which become available to any file that imports it.

```zenth
// mathutils.zn
const pi = 3.14159;

fn add(a: Int, b: Int) -> Int {
    return a + b;
}

fn multiply(a: Int, b: Int) -> Int {
    return a * b;
}
```

**Rules for module files:**
- Must not define `fn main`
- Can import other local modules
- Can define functions, objects, constants, and type aliases
- All top-level definitions are accessible to importers

## Using Module Objects

Objects defined in a module are available to importers. Construct them using `module.ObjName(field=value)` syntax:

```zenth
// shapes/rect.zn
obj Rectangle {
    width: Float;
    height: Float;

    fn area() -> Float {
        return self.width * self.height;
    }
}
```

```zenth
// main.zn
import "./shapes";

fn main() {
    let r = shapes.Rectangle(width=10.0, height=5.0);
    Println(Str(r.area()));  // "50"
}
```

## Directory Modules

If the import path points to a directory, Zenth loads **all `.zn` files** in that directory as a single module. The module name is the directory name:

```
project/
  main.zn
  geometry/
    rect.zn      ← contains Rectangle obj
    circle.zn    ← contains circle_area fn
```

```zenth
// main.zn
import "./geometry";

fn main() {
    let r = geometry.Rectangle(width=3.0, height=4.0);
    Println(Str(r.area()));

    let a = geometry.circle_area(5.0);
    Println(Str(a));
}
```

All definitions across the files in the directory are merged into one module namespace.

## Module Search

Zenth resolves module paths **relative to the file containing the import**:

- `import "./utils"` → looks for `utils.zn` in the same directory, or a `utils/` directory
- `import "../shared"` → looks one directory up

## Example: Multi-file Project

```
project/
  main.zn
  testmods/
    mathutils.zn
    shapes/
      rect.zn
      circle.zn
```

```zenth
// testmods/mathutils.zn
fn add(a: Int, b: Int) -> Int {
    return a + b;
}

fn multiply(a: Int, b: Int) -> Int {
    return a * b;
}
```

```zenth
// main.zn
import "./testmods/mathutils" as math;
import "./testmods/shapes";

fn main() {
    Println(Str(math.add(3, 4)));         // "7"
    Println(Str(math.multiply(6, 7)));    // "42"

    let r = shapes.Rectangle(width=3.0, height=4.0);
    Println(Str(r.area()));               // "12"
    Println(Str(shapes.circle_area(1.0))); // "3.14159"
}
```

## Compared to Standard Library Imports

| | Standard library | Local module |
|---|---|---|
| Path | `"fmt"`, `"math"`, etc. | `"./path/to/module"` |
| Source | Go stdlib | Your `.zn` files |
| Alias | `import "math" as m;` | `import "./utils" as u;` |

Both use the same `module.function()` call syntax.

## Module-Internal Function Calls

Functions within a module can call other functions defined in the same module. This allows you to build composable modules where functions delegate to each other:

```zenth
// mathutils.zn
fn add(a: Int, b: Int) -> Int {
    return a + b;
}

fn multiply(a: Int, b: Int) -> Int {
    return a * b;
}

// sum_and_double calls add and multiply within the same module
fn sum_and_double(a: Int, b: Int) -> Int {
    let s = add(a, b);
    return multiply(s, 2);
}
```

```zenth
// main.zn
import "./mathutils";

fn main() {
    Println(Str(mathutils.sum_and_double(5, 3)));  // "16"
}
```

## Qualified Types

When you need to refer to a type defined in a module, use `module.Type` syntax. This works in type annotations, function parameters, and generic type positions like `Array()` and `Tuple()`:

```zenth
// items.zn
obj Item {
    name: Str;
    value: Int;
}

fn make_item(name: Str, value: Int) -> Item {
    return Item(name=name, value=value);
}
```

```zenth
// main.zn
import "./items";

fn print_items(data: Array(items.Item)) {
    for item in data {
        Println("{item.name}: {Str(item.value)}");
    }
}

fn main() {
    let a = items.make_item("apple", 3);
    let b = items.make_item("banana", 5);
    let list: Array(items.Item) = [a, b];
    print_items(list);
}
```

Qualified types work anywhere a type is expected:
- Function parameters: `fn process(item: task.Task) { ... }`
- Variable annotations: `let items: Array(task.Task) = [];`
- Generic positions: `Array(items.Item)`, `Tuple(Int, task.Task)`, `Hashmap(Str, task.Task)`
