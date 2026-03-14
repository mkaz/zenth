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
let scores: array(int) = [100, 95, 87];
```

An empty slice can use a type annotation:

```zenth
var points: array(Point) = [];
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
println(len(items));       // 3
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

## Modifying Slices

Zenth provides methods to append, prepend, and remove elements from a slice. The slice must be assigned to a mutable `var` to use these methods.

### Adding Elements

- `.add(value)` appends an element to the end of the slice.
- `.push(value)` prepends an element to the beginning of the slice.

```zenth
var items = [2, 3];
items.add(4);   // [2, 3, 4]
items.push(1);  // [1, 2, 3, 4]
```

### Removing Elements

- `.pop()` removes and returns the last element of the slice.
- `.pop(index)` removes and returns the element at the specified index.

```zenth
var letters = ["a", "b", "c", "d"];
let last = letters.pop();    // "d", letters is now ["a", "b", "c"]
let first = letters.pop(0);  // "a", letters is now ["b", "c"]
```

## Max and Min

Use `.max()` and `.min()` to find the largest and smallest elements of a numeric array:

```zenth
let nums = [3, 1, 4, 1, 5, 9, 2, 6];
println(nums.max());       // 9
println(nums.min());       // 1

let floats = [3.14, 2.71, 1.41];
println(floats.max());     // 3.14
println(floats.min());     // 1.41
```

These methods work on any numeric array (`int`, `f64`, `i32`, etc.) and panic at runtime if called on an empty array.

## Sum

Use `.sum()` to compute the sum of a numeric array:

```zenth
let nums = [1, 2, 3, 4, 5];
println(nums.sum());       // 15

let floats = [1.5, 2.5, 3.0];
println(floats.sum());     // 7
```

Returns zero for empty arrays.

## Sorted

Use `.sorted()` to get a new sorted copy of an array. Works on numeric and string arrays:

```zenth
let unsorted = [3, 1, 4, 1, 5, 9];
let s = unsorted.sorted();  // [1, 1, 3, 4, 5, 9]
// unsorted is unchanged

let words = ["banana", "apple", "cherry"];
let sw = words.sorted();  // ["apple", "banana", "cherry"]
```

The original array is not modified. Strings are sorted lexicographically.

## Reduce

Use `.reduce()` to combine all elements into a single value using a closure that takes an accumulator and element:

```zenth
let nums = [1, 2, 3, 4, 5];

// Sum all elements
let total = nums.reduce(fn(a, b) a + b);
println(total);  // 15

// Product with initial value
let product = nums.reduce(fn(a, b) a * b, 1);
println(product);  // 120

// Find max using reduce
let biggest = nums.reduce(fn(a, b) if a > b { a } else { b });
println(biggest);  // 5
```

Without an initial value, the first element is used as the starting accumulator and reduction begins from the second element. Calling `reduce()` on an empty array without an initial value panics at runtime.

With an initial value, reduction starts from that value and processes all elements:

```zenth
let nums = [1, 2, 3];
let sum_from_100 = nums.reduce(fn(a, b) a + b, 100);
println(sum_from_100);  // 106
```

Works on any array type — numeric, string, etc. The closure must take two parameters of the element type and return the same type.

## Transforming Slices

Use `.map()` to transform each element and `.filter()` to select elements:

```zenth
let numbers = [1, 2, 3, 4, 5];

// Double each element
let doubled = numbers.map(fn(x) x * 2);
// [2, 4, 6, 8, 10]

// Keep only even numbers
let evens = numbers.filter(fn(x) x % 2 == 0);
// [2, 4]

// Chaining
let result = numbers.filter(fn(x) x > 2).map(fn(x) x * 10);
// [30, 40, 50]
```

Closures with a block body use explicit `return`:

```zenth
let processed = numbers.map(fn(x: int) -> int {
    let y = x + 10;
    return y;
});
```

## Example

```zenth
fn sum(numbers: array(int)) -> int {
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
