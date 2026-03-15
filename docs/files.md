# Files

The `file()` built-in creates a file object for reading files and inspecting file paths. File objects are value types -- there are no open handles and nothing to close.

## Creating a File Object

```zenth
let f = file("/path/to/data.txt");
```

## Methods

### exists

Returns `true` if the file exists on disk:

```zenth
let f = file("config.txt");
if f.exists() {
    println("found it");
}
```

### read

Returns the full file contents as a string. Exits with an error if the file cannot be read:

```zenth
let content = f.read();
println(content);
```

### lines

Returns the file contents split into lines as `array(str)`. Exits with an error if the file cannot be read:

```zenth
let lines = f.lines();
for i, line in lines {
    println("{i}: {line}");
}
```

### name

Returns the filename (base name) from the path:

```zenth
let f = file("/home/user/data.txt");
println(f.name());  // data.txt
```

### sections

Splits the file contents on blank lines, returning `array(str)`. This is useful for input that has sections separated by empty lines (common in Advent of Code problems):

```zenth
let parts = file("input.txt").sections();
for i, section in parts {
    println("Section {i}: {section}");
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
let f = file("report.csv");
println(f.ext());  // .csv
```

## Example

```zenth
fn main() {
    let f = file("input.txt");
    if !f.exists() {
        println("file not found: " + f.name());
        exit(1);
    }
    let lines = f.lines();
    println("Read {lines.length()} lines from {f.name()}");
}
```
