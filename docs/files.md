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
