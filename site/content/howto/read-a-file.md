---
title: "Read and Process a File"
weight: 1
---

# Read and Process a File

This guide shows how to read a file, process its contents line by line, and write the results to another file.

## Read an entire file as a string

```zenth
fn main() {
    let content = File("input.txt").read();
    Println(content);
}
```

## Read a file line by line

Use `.lines()` to get an array of strings, one per line:

```zenth
fn main() {
    let lines = File("input.txt").lines();
    for i, line in lines {
        Println("{i}: {line}");
    }
}
```

## Check if a file exists first

```zenth
fn main() {
    let f = File("input.txt");
    if !f.exists() {
        Println("file not found: " + f.name());
        Exit(1);
    }
    let lines = f.lines();
    Println("Read {lines.length()} lines");
}
```

## Process and write results to a new file

Read lines, transform them, and write the output:

```zenth
fn main() {
    let lines = File("data.txt").lines();
    let out = File("output.txt");
    out.write("");  // start fresh

    for line in lines {
        let upper = line.upper();
        out.append(upper + "\n");
    }

    Println("Wrote {lines.length()} lines to output.txt");
}
```

## Read sections separated by blank lines

Some files use blank lines to separate blocks of data. Use `.sections()` to split on those boundaries:

```zenth
fn main() {
    let parts = File("input.txt").sections();
    for i, section in parts {
        Println("Section {i}: {section}");
    }
}
```

## See also

- [Parse Puzzle Input](parse-data.md) — common input parsing patterns for Advent of Code style problems
- [Files reference](../reference/files.md) — full list of File methods
- [Strings reference](../reference/strings.md) — string manipulation methods
