---
title: "Dates"
weight: 12
---

Zenth has a built-in `Date` type for working with calendar dates. Dates are created using static constructors and support formatting, addition, and subtraction.

## Creating Dates

### Today's date

```zenth
let d = Date.today();
```

Returns a `Date` object representing the current date (time portion is midnight UTC).

### From a string

```zenth
let epoch = Date.from("1970-01-01");
```

Parses a date string using the default format `%Y-%m-%d` (year-month-day).

### From a string with custom format

```zenth
let jan13 = Date.from("01/13/2007", "%m/%d/%Y");
```

The second argument is a Python-style format string. Available format specifiers:

| Specifier | Meaning | Example |
|-----------|---------|---------|
| `%Y` | 4-digit year | `2024` |
| `%y` | 2-digit year | `24` |
| `%m` | 2-digit month | `03` |
| `%d` | 2-digit day | `21` |
| `%H` | Hour (24h) | `15` |
| `%M` | Minute | `04` |
| `%S` | Second | `05` |
| `%B` | Full month name | `January` |
| `%b` | Abbreviated month | `Jan` |
| `%A` | Full weekday name | `Monday` |
| `%a` | Abbreviated weekday | `Mon` |

## Formatting Dates

Use `.format()` with a Python-style format string:

```zenth
let d = Date.from("2024-03-21");
Println(d.format("%Y-%m-%d"));    // 2024-03-21
Println(d.format("%d/%m/%Y"));    // 21/03/2024
Println(d.format("%B %d, %Y"));   // March 21, 2024
```

## Adding Time

Use `.add(value)` or `.add(value, unit)` to add time to a date. Returns a new `Date` (dates are immutable).

```zenth
let d = Date.from("2024-03-21");

// Default unit is "days"
let tomorrow = d.add(1);           // 2024-03-22
let next_week = d.add(7);          // 2024-03-28

// Explicit units
let next_month = d.add(1, "months");  // 2024-04-21
let next_year = d.add(1, "years");    // 2025-03-21
```

Supported units: `"days"` (default), `"months"`, `"years"`.

## Subtracting Time

Use `.sub(value)` or `.sub(value, unit)` to subtract time from a date. Returns a new `Date`.

```zenth
let d = Date.from("2024-03-21");

let yesterday = d.sub(1);             // 2024-03-20
let last_month = d.sub(1, "months");  // 2024-02-21
let last_year = d.sub(1, "years");    // 2023-03-21
```

## Chaining

All date methods return new `Date` objects, so calls can be chained:

```zenth
let d = Date.today().add(1, "years").add(3, "months").add(10);
Println(d.format("%Y-%m-%d"));
```

## Converting to String

Use `.format()` for controlled output. `Str()` works but returns Go's verbose time format:

```zenth
let d = Date.from("2024-03-21");
Println(d.format("%Y-%m-%d"));   // "2024-03-21" (preferred)
Println(Str(d));                  // "2024-03-21 00:00:00 +0000 UTC" (verbose)
```

## Complete Example

```zenth
fn main() {
    let today = Date.today();
    Println("Today: " + today.format("%B %d, %Y"));

    let deadline = today.add(30);
    Println("Deadline: " + deadline.format("%Y-%m-%d"));

    let epoch = Date.from("1970-01-01");
    let y2k = epoch.add(30, "years");
    Println("Y2K: " + y2k.format("%Y-%m-%d"));
}
```
