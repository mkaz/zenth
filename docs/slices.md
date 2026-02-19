# Slices

Slices are dynamically-sized sequences of elements, all of the same type.

## Creating Slices

Use bracket syntax for slice literals:

```zenth
let numbers = [1, 2, 3, 4, 5];
let names = ["Alice", "Bob", "Charlie"];
```

With an explicit type annotation:

```zenth
let scores: []int = [100, 95, 87];
```

An empty slice can use a type annotation:

```zenth
var points: []Point = [];
points.add(Point(x=1, y=2));
```

## Accessing Elements

Use bracket indexing (zero-based):

```zenth
let first = numbers[0];   // 1
let third = numbers[2];   // 3
```

## Length

Use `len()` to get the number of elements:

```zenth
let items = [10, 20, 30];
println(str(len(items)));  // 3
```

## Iterating

Use `for-in` to loop over elements:

```zenth
let fruits = ["apple", "banana", "cherry"];
for fruit in fruits {
    println(fruit);
}
```

With an index:

```zenth
for i, fruit in fruits {
    println("{i}: {fruit}");
}
```

## Containment

Use `.exists()` to check if a slice contains a value:

```zenth
let nums = [10, 20, 30];
if nums.exists(20) {
    println("found 20");
}
```

## Example

```zenth
fn sum(numbers: []int) -> int {
    var total = 0;
    for n in numbers {
        total += n;
    }
    return total;
}

fn main() {
    let data = [1, 2, 3, 4, 5];
    println("Sum: {sum(data)}");
}
```
