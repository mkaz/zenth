---
title: "Parse Puzzle Input"
weight: 2
---

# Parse Puzzle Input

This guide shows common parsing patterns for structured text input: one number per line, comma-separated values, blank-line sections, coordinates, key/value fields, and mixed text. These are the kinds of transformations that come up constantly in Advent of Code problems.

## Read a whole file and trim it

Use `.read()` when the input is one block of text, then `.strip()` to remove the trailing newline:

```zenth
fn main() {
    let raw = File("input.txt").read().strip();
    Println(raw);
}
```

This is useful for single-line inputs like `389125467` or `target area: x=20..30, y=-10..-5`.

## Parse one integer per line

When every line is a number, convert the full array in one step:

```zenth
fn main() {
    let nums = File("input.txt").lines().to_int();
    Println(nums.sum());
}
```

This is the shortest pattern for inputs like:

```text
199
200
208
210
```

## Parse comma-separated integers

For inputs stored on a single line, read the line, split on commas, then convert:

```zenth
fn main() {
    let line = File("input.txt").lines()[0];
    let nums = line.split(",").to_int();
    Println(nums);
}
```

If you already have the full file as a string:

```zenth
let nums = File("input.txt").read().strip().split(",").to_int();
```

## Parse blank-line sections

Use `.sections()` when the file is divided into groups separated by blank lines:

```zenth
fn main() {
    let groups = File("input.txt").sections();

    for group in groups {
        let nums = group.split("\n").to_int();
        Println(nums.sum());
    }
}
```

This pattern is useful for grouped calorie counts, multiple decks, batches of rules, or map tiles.

## Split each row into fields

Use `.split()` with no argument for whitespace-delimited data:

```zenth
fn main() {
    let rows = File("input.txt").lines();

    for row in rows {
        let fields = row.split();
        Println(fields);
    }
}
```

Convert the fields immediately when they should all be numbers:

```zenth
fn main() {
    let rows = File("input.txt").lines();

    for row in rows {
        let nums = row.split().to_int();
        Println(nums.max() - nums.min());
    }
}
```

## Use `split_once()` for two-part records

When each line has exactly one delimiter, `split_once()` is usually clearer than `split()`:

```zenth
let (cmd, value) = "forward 10".split_once(" ");
let n = Int(value);
```

For coordinate pairs:

```zenth
obj Point {
    x: Int;
    y: Int;
}

fn parse_point(raw: Str) -> Point {
    let (raw_x, raw_y) = raw.split_once(",");
    return Point(x=Int(raw_x), y=Int(raw_y));
}

fn parse_segment(line: Str) -> Tuple(Point, Point) {
    let (left, right) = line.split_once(" -> ");
    return Tuple(parse_point(left), parse_point(right));
}
```

This is a good fit for inputs like `498,4 -> 498,6` or `Player 1: 4`.

## Extract numbers from mixed text

For lines with fixed prefixes, peel the text away in stages:

```zenth
fn main() {
    let line = "target area: x=20..30, y=-10..-5";

    let parts = line.strip_prefix("target area: ").split(", ");
    let xparts = parts[0].strip_prefix("x=").split("..").to_int();
    let yparts = parts[1].strip_prefix("y=").split("..").to_int();

    let xmin = xparts[0];
    let xmax = xparts[1];
    let ymin = yparts[0];
    let ymax = yparts[1];

    Println("{xmin} {xmax} {ymin} {ymax}");
}
```

This approach works well when the format is regular and the separators are stable.

## Parse key/value records

Some inputs encode records as whitespace-separated `key:value` pairs:

```zenth
fn main() {
    let records = File("input.txt").sections();

    for record in records {
        var fields = Hashmap(Str, Str);

        for pair in record.split() {
            let (key, value) = pair.split_once(":");
            fields[key] = value;
        }

        Println(fields["byr"]);
    }
}
```

Using `.sections()` here is usually better than `.lines()` because one record can span multiple lines.

## Convert individual values to the right type

Use the built-in conversion functions when fields have different types:

```zenth
fn main() {
    let fields = "42 3.5".split();

    let count = Int(fields[0]);
    let ratio = Float(fields[1]);

    Println("{count} {ratio}");
}
```

For whole arrays, the array conversions are shorter:

```zenth
let ints = ["1", "2", "3"].to_int();
let floats = ["1.5", "2.5"].to_float();
```

## Build small helpers for repeated parsing

As soon as a format repeats, move the parsing into a function:

```zenth
fn parse_rule(line: Str) -> Tuple(Str, Int, Int) {
    let (name, rest) = line.split_once(": ");
    let ranges = rest.split(" or ");
    let left = ranges[0].split("-").to_int();
    return Tuple(name, left[0], left[1]);
}
```

Small parsing helpers keep the main solving code focused on the algorithm instead of string manipulation.

## A complete AoC-style example

This combines file reading, section splitting, and integer conversion:

```zenth
fn main() {
    let groups = File("input.txt").sections();

    var sums: Array(Int) = [];
    for group in groups {
        let nums = group.split("\n").to_int();
        sums.add(nums.sum());
    }

    sums = sums.sorted("desc");
    Println("Top group: {sums[0]}");
    Println("Top three: {sums[0] + sums[1] + sums[2]}");
}
```

## See also

- [Read and Process a File](read-a-file.md) — file input and output basics
- [Strings reference](../reference/strings.md) — `split`, `split_once`, `strip`, `strip_prefix`
- [Arrays reference](../reference/arrays.md) — array helpers like `.get()`, `.map()`, `.sorted()`
- [Types reference](../reference/types.md) — `Int()`, `Float()`, `.to_int()`, `.to_float()`
- [Files reference](../reference/files.md) — `read`, `lines`, `sections`
