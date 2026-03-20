# Arrays

Arrays are dynamically-sized sequences of elements, all of the same type.

## Creating Arrays

Use bracket syntax for array literals:

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

Use `.repeat(n)` to create an array by repeating elements:

```zenth
var zeros = [0].repeat(10);      // [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
let dashes = ["-"].repeat(5);    // ["-", "-", "-", "-", "-"]
let pattern = [1, 2].repeat(6);  // [1, 2, 1, 2, 1, 2]
```

When the source array has multiple elements, they tile to fill the requested count.

## Accessing Elements

Use bracket indexing (zero-based):

```zenth
let first = numbers[0];   // 1
let third = numbers[2];   // 3
```

## Length

Use `Len()` to get the number of elements:

```zenth
let items = [10, 20, 30];
Println(Len(items));       // 3
```

## Iterating

Use `for-in` to loop over elements:

```zenth
let fruits = ["apple", "banana", "cherry"];
for fruit in fruits {
    Println(fruit);
}
```

With an index:

```zenth
for i, fruit in fruits {
    Println("{i}: {fruit}");
}
```

## Containment

Use `.exists()` to check if an array contains a value:

```zenth
let nums = [10, 20, 30];
if nums.exists(20) {
    Println("found 20");
}
```

## Modifying Arrays

Zenth provides methods to append, prepend, and remove elements from an array. The array must be assigned to a mutable `var` to use these methods.

### Index Assignment

Assign directly to an index using bracket notation. The array must be a mutable `var`:

```zenth
var nums = [10, 20, 30];
nums[1] = 99;            // [10, 99, 30]
nums[0] = nums[0] + 5;   // [15, 99, 30]
```

Index assignment on an immutable `let` binding is a compile-time error.

### Adding Elements

- `.add(value)` appends an element to the end of the array.
- `.push(value)` prepends an element to the beginning of the array.
- `.extend(other)` appends all elements from another array of the same type.

```zenth
var items = [2, 3];
items.add(4);   // [2, 3, 4]
items.push(1);  // [1, 2, 3, 4]

var nums = [1, 2];
nums.extend([3, 4, 5]);  // [1, 2, 3, 4, 5]
```

### Inserting Elements

- `.insert(index, value)` inserts an element at the specified index, shifting later elements right.

```zenth
var nums = [1, 2, 4, 5];
nums.insert(2, 3);  // [1, 2, 3, 4, 5]
```

### Removing Elements

- `.pop()` removes and returns the last element of the array.
- `.pop(index)` removes and returns the element at the specified index.
- `.remove(index)` removes the element at the specified index (discards it).

```zenth
var letters = ["a", "b", "c", "d"];
let last = letters.pop();    // "d", letters is now ["a", "b", "c"]
let first = letters.pop(0);  // "a", letters is now ["b", "c"]

var items = [10, 20, 30];
items.remove(1);  // removes 20, items is now [10, 30]
```

Use `.remove()` instead of `.pop()` when you don't need the removed value.

## Max and Min

Use `.max()` and `.min()` to find the largest and smallest elements of a numeric array:

```zenth
let nums = [3, 1, 4, 1, 5, 9, 2, 6];
Println(nums.max());       // 9
Println(nums.min());       // 1

let floats = [3.14, 2.71, 1.41];
Println(floats.max());     // 3.14
Println(floats.min());     // 1.41
```

These methods work on any numeric array (`int`, `f64`, `i32`, etc.) and panic at runtime if called on an empty array.

## Sum

Use `.sum()` to compute the sum of a numeric array:

```zenth
let nums = [1, 2, 3, 4, 5];
Println(nums.sum());       // 15

let floats = [1.5, 2.5, 3.0];
Println(floats.sum());     // 7
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

Pass `"desc"` to sort in descending order (largest to smallest), or `"asc"` for ascending (smallest to largest, the default):

```zenth
let nums = [3, 1, 4, 1, 5];
let asc = nums.sorted("asc");    // [1, 1, 3, 4, 5]
let desc = nums.sorted("desc");  // [5, 4, 3, 1, 1]

let words = ["banana", "apple", "cherry"];
let rev = words.sorted("desc");  // ["cherry", "banana", "apple"]
```

The original array is not modified. Strings are sorted lexicographically.

## Reduce

Use `.reduce()` to combine all elements into a single value using a closure that takes an accumulator and element:

```zenth
let nums = [1, 2, 3, 4, 5];

// Sum all elements
let total = nums.reduce(fn(a, b) a + b);
Println(total);  // 15

// Product with initial value
let product = nums.reduce(fn(a, b) a * b, 1);
Println(product);  // 120

// Find max using reduce
let biggest = nums.reduce(fn(a, b) if a > b { a } else { b });
Println(biggest);  // 5
```

Without an initial value, the first element is used as the starting accumulator and reduction begins from the second element. Calling `reduce()` on an empty array without an initial value panics at runtime.

With an initial value, reduction starts from that value and processes all elements:

```zenth
let nums = [1, 2, 3];
let sum_from_100 = nums.reduce(fn(a, b) a + b, 100);
Println(sum_from_100);  // 106
```

Works on any array type — numeric, string, etc. The closure must take two parameters of the element type and return the same type.

## Join

Use `.join(separator)` to concatenate all elements of a string array into a single string, with the given separator between each element:

```zenth
let words = ["hello", "world"];
Println(words.join(" "));      // "hello world"

let csv = ["a", "b", "c"];
Println(csv.join(","));        // "a,b,c"

// Empty separator concatenates directly
let letters = ["x", "y", "z"];
Println(letters.join(""));     // "xyz"

// Multi-character separator
let parts = ["2026", "03", "15"];
Println(parts.join("-"));      // "2026-03-15"
```

Only works on `array(str)`. Calling `.join()` on a non-string array is a compile-time error.

## Transforming Arrays

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

## Destructuring

Use array destructuring to bind array elements directly to variables. The syntax mirrors slice literals: `let [a, b, c] = expr`:

```zenth
let nums = [10, 20, 30];
let [a, b, c] = nums;
Println(a);  // 10
Println(b);  // 20
Println(c);  // 30
```

Works with `let`, `var`, and `const`:

```zenth
var [x, y] = [1, 2];
x = 99;      // ok: x is mutable
```

Use `_` to discard positions you don't need:

```zenth
let [first, _, third] = [1, 2, 3];
Println(first);  // 1
Println(third);  // 3
```

Especially useful with method chains:

```zenth
let line = "12x4x8";
let [l, w, h] = line.split("x").to_int().sorted();
Println(l);  // 4
Println(w);  // 8
Println(h);  // 12
```

The number of binding names must not exceed the array length at runtime; a runtime panic occurs if the array is too short.

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
    Println("Sum: {sum(data)}");
}
```
