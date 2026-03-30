---
title: "Files"
weight: 13
---

The `File()` built-in creates a file object for reading and writing files and inspecting file paths. File objects are value types -- there are no open handles and nothing to close.

## Creating a File Object

```zenth
let f = File("/path/to/data.txt");
```

## Methods

### exists

Returns `true` if the file exists on disk:

```zenth
let f = File("config.txt");
if f.exists() {
    Println("found it");
}
```

### read

Returns the full file contents as a string. Exits with an error if the file cannot be read:

```zenth
let content = f.read();
Println(content);
```

### lines

Returns the file contents split into lines as `Array(Str)`. Exits with an error if the file cannot be read:

```zenth
let lines = f.lines();
for i, line in lines {
    Println("{i}: {line}");
}
```

### name

Returns the filename (base name) from the path:

```zenth
let f = File("/home/user/data.txt");
Println(f.name());  // data.txt
```

### sections

Splits the file contents on blank lines, returning `Array(Str)`. This is useful for input that has sections separated by empty lines (common in Advent of Code problems):

```zenth
let parts = File("input.txt").sections();
for i, section in parts {
    Println("Section {i}: {section}");
}
```

For a file containing:
```
first block

second block

third block
```

`sections()` returns an array of three strings: `["first block", "second block", "third block"]`.

### ext

Returns the file extension, including the dot:

```zenth
let f = File("report.csv");
Println(f.ext());  // .csv
```

### write

Writes content to the file, creating it if it doesn't exist or overwriting it if it does. Exits with an error if the file cannot be written:

```zenth
let f = File("output.txt");
f.write("Hello world\n");
```

After this call, `output.txt` contains exactly `Hello world\n` regardless of any previous contents.

### append

Appends content to the end of the file. Creates the file if it doesn't exist. Exits with an error if the file cannot be written:

```zenth
let f = File("log.txt");
f.append("first line\n");
f.append("second line\n");
```

After these calls, `log.txt` contains:
```
first line
second line
```

`append()` is useful for building up files incrementally, such as log files or accumulated output:

```zenth
let out = File("results.txt");
out.write("");  // start fresh
for i in Range(0, 5) {
    out.append("result {i}\n");
}
```

## Examples

### Reading a file

```zenth
fn main() {
    let f = File("input.txt");
    if !f.exists() {
        Println("file not found: " + f.name());
        Exit(1);
    }
    let lines = f.lines();
    Println("Read {lines.length()} lines from {f.name()}");
}
```

### Writing and appending

```zenth
fn main() {
    let f = File("output.txt");

    // write creates or overwrites the file
    f.write("Hello world\n");

    // append adds to the end
    f.append("Goodbye world\n");

    // verify
    Println(f.read());
    // Hello world
    // Goodbye world
}
```
