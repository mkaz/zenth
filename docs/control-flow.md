# Control Flow

## If / Else

```zenth
if x > 10 {
    Println("big");
} else if x > 5 {
    Println("medium");
} else {
    Println("small");
}
```

Conditions must be `Bool` -- there is no truthy/falsy coercion. Parentheses around the condition are not required.

## For Loops

`for` is the only loop keyword in Zenth. It supports five styles.

### C-Style For

The classic three-part loop with init, condition, and post:

```zenth
for var i = 0; i < 10; i++ {
    Println(i);
}
```

The init clause can declare a variable with `var` or `let`:

```zenth
for let start = 0; start < 5; start++ {
    Println(start);
}
```

### For-Range Count

Run a loop body a fixed number of times:

```zenth
for range 3 {
    Println("tick");
}

let n = 5;
for range n {
    Println("count");
}
```

The count expression must be `Int`.

### While-Style

A `for` with just a condition acts as a while loop:

```zenth
var running = true;
for running {
    process();
    running = check();
}
```

### Infinite Loop

A bare `for` loops forever (use `break` to exit):

```zenth
for {
    let input = read();
    if input == "quit" {
        break;
    }
    handle(input);
}
```

### For-In

Iterate over arrays and strings:

```zenth
let items = ["apple", "banana", "cherry"];
for item in items {
    Println(item);
}
```

With an index variable:

```zenth
for i, item in items {
    Println("{i}: {item}");
}
```

Using `Range()` and `Rangei()` for numeric iteration:

```zenth
// Range(start, end) — exclusive end
for i in Range(0, 5) {
    Println(i);       // 0 1 2 3 4
}

// Rangei(start, end) — inclusive end
for i in Rangei(1, 5) {
    Println(i);       // 1 2 3 4 5
}

// Optional step parameter
for i in Range(0, 10, 2) {
    Println(i);       // 0 2 4 6 8
}
```

`Range()` and `Rangei()` return a **range object**, not an array. Range objects are lightweight and iterable. They also support an O(1) `.contains()` method for bounds checking:

```zenth
let xrange = Rangei(0, 100);
let yrange = Rangei(0, 50);

if xrange.contains(x) && yrange.contains(y) {
    Println("point in bounds");
}
```

This replaces verbose compound comparisons:

```zenth
// before
if x >= 0 && x <= 100 && y >= 0 && y <= 50 {

// after
if xrange.contains(x) && yrange.contains(y) {
```

Ranges can also be stored in variables and reused across multiple checks. `Len()` on a range object returns the number of elements without materializing the sequence:

```zenth
let r = Range(0, 100, 3);
Println(Len(r));           // 34
Println(r.contains(99));   // false (99 is not a multiple of 3 from 0... but bounds check: 99 < 100 → true)
```

Note: `.contains()` performs a bounds check only. For `Range(0, 10, 3)`, `.contains(5)` returns `true` because 5 is within the bounds `[0, 10)`, not because 5 is in the sequence `[0, 3, 6, 9]`.

Use `_` to discard the loop variable when you only need the repetition:

```zenth
for _ in Range(0, 10) {
    Println("tick");
}
```

The `_` discard also works as the index or value in two-variable loops:

```zenth
for _, item in items {   // discard index
    Println(item);
}
for i, _ in items {      // discard value
    Println(i);
}
```

Iterating over a string yields individual characters as strings:

```zenth
for ch in "hello" {
    Println(ch);
}
```

## Break and Continue

`break` exits the innermost loop. `continue` skips to the next iteration:

```zenth
for var i = 0; i < 100; i++ {
    if i % 2 == 0 {
        continue;  // skip even numbers
    }
    if i > 10 {
        break;     // stop after 10
    }
    Println(i);
}
```

## If Expressions

`if` can be used as an expression to produce a value. Both branches must be present and return the same type:

```zenth
let label = if x > 5 { "big" } else { "small" };
```

Else-if chains work too:

```zenth
let grade = if x > 9 { "A" } else if x > 7 { "B" } else { "C" };
```

If expressions can be used inline:

```zenth
Println(if done { "yes" } else { "no" });
```

## Match

`match` selects a branch based on a value, similar to `switch` in other languages:

```zenth
let day = 3;
match day {
    1 => Println("Monday");
    2 => Println("Tuesday");
    3 => Println("Wednesday");
    4 => Println("Thursday");
    5 => Println("Friday");
    _ => Println("Weekend");
}
```

Use `_` as the wildcard/default case.

### Block Bodies

Match arms can use blocks for multiple statements:

```zenth
match status {
    200 => {
        Println("OK");
        handle_success();
    }
    404 => Println("Not Found");
    _ => Println("Other");
}
```

### Match Expressions

`match` can also be used as an expression that returns a value. Each arm uses `=>` followed by a single expression, with commas separating arms:

```zenth
let label = match x {
    1 => "one",
    2 => "two",
    _ => "other"
};
```

Match expressions can be used anywhere an expression is expected:

```zenth
Println(match status {
    200 => "OK",
    404 => "Not Found",
    _ => "Unknown"
});
```

All arms must return the same type.

### Practical Example: FizzBuzz

```zenth
fn fizzbuzz(n: Int) -> Str {
    if n % 15 == 0 {
        return "FizzBuzz";
    } else if n % 3 == 0 {
        return "Fizz";
    } else if n % 5 == 0 {
        return "Buzz";
    } else {
        return Str(n);
    }
}

fn main() {
    for var i = 1; i <= 20; i++ {
        Println(fizzbuzz(i));
    }
}
```
