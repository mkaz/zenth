# Control Flow

## If / Else

```zenth
if x > 10 {
    println("big");
} else if x > 5 {
    println("medium");
} else {
    println("small");
}
```

Conditions must be `bool` -- there is no truthy/falsy coercion. Parentheses around the condition are not required.

## For Loops

`for` is the only loop keyword in Zenth. It supports five styles.

### C-Style For

The classic three-part loop with init, condition, and post:

```zenth
for var i = 0; i < 10; i++ {
    println(str(i));
}
```

The init clause can declare a variable with `var` or `let`:

```zenth
for let start = 0; start < 5; start++ {
    println(str(start));
}
```

### For-Range Count

Run a loop body a fixed number of times:

```zenth
for range 3 {
    println("tick");
}

let n = 5;
for range n {
    println("count");
}
```

The count expression must be `int`.

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

Iterate over slices and strings:

```zenth
let items = ["apple", "banana", "cherry"];
for item in items {
    println(item);
}
```

With an index variable:

```zenth
for i, item in items {
    println("{i}: {item}");
}
```

Using `range()` and `rangei()` for numeric iteration:

```zenth
// range(start, end) — exclusive end
for i in range(0, 5) {
    println(str(i));  // 0 1 2 3 4
}

// rangei(start, end) — inclusive end
for i in rangei(1, 5) {
    println(str(i));  // 1 2 3 4 5
}

// Optional step parameter
for i in range(0, 10, 2) {
    println(str(i));  // 0 2 4 6 8
}
```

Iterating over a string yields individual characters as strings:

```zenth
for ch in "hello" {
    println(ch);
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
    println(str(i));
}
```

## Match

`match` selects a branch based on a value, similar to `switch` in other languages:

```zenth
let day = 3;
match day {
    1 => println("Monday");
    2 => println("Tuesday");
    3 => println("Wednesday");
    4 => println("Thursday");
    5 => println("Friday");
    _ => println("Weekend");
}
```

Use `_` as the wildcard/default case.

### Block Bodies

Match arms can use blocks for multiple statements:

```zenth
match status {
    200 => {
        println("OK");
        handle_success();
    }
    404 => println("Not Found");
    _ => println("Other");
}
```

### Practical Example: FizzBuzz

```zenth
fn fizzbuzz(n: int) -> str {
    if n % 15 == 0 {
        return "FizzBuzz";
    } else if n % 3 == 0 {
        return "Fizz";
    } else if n % 5 == 0 {
        return "Buzz";
    } else {
        return str(n);
    }
}

fn main() {
    for var i = 1; i <= 20; i++ {
        println(fizzbuzz(i));
    }
}
```
