---
title: "Sets"
weight: 11
---

Sets are unordered collections of unique values.

## Creating Sets

Create an empty set with `Set(Type)`:

```zenth
var visited = Set(Str);
visited.add("start");
visited.add("middle");
```

You can also convert an array into a set with `Set(items)`:

```zenth
let letters = ["a", "b", "a"];
let unique = Set(letters);
Println(Len(unique));  // 2
```

## Core Operations

Use `.add(value)` to insert an item, `.exists(value)` to test membership, and `.remove(value)` to delete an item:

```zenth
var nums = Set(Int);
nums.add(1);
nums.add(2);
nums.add(2);  // duplicate, ignored

Println(nums.exists(1));  // true
nums.remove(1);
Println(Len(nums));       // 1
```

Sets also expose `.length()` when you prefer method syntax:

```zenth
Println(nums.length());
```

## Iteration

Iterate over a set with `for-in`:

```zenth
var colors = Set(Str);
colors.add("red");
colors.add("green");

for color in colors {
    Println(color);
}
```

Set iteration order is not guaranteed.

## Comparable Element Types

Set elements must be comparable types. Built-in scalar types work directly, and tuples are supported when all of their elements are comparable:

```zenth
var dots = Set(Tuple(Int, Int));
dots.add(Tuple(6, 10));
dots.add(Tuple(0, 14));
dots.add(Tuple(6, 10));  // duplicate, ignored

if dots.exists(Tuple(6, 10)) {
    Println("found");
}
```

Tuples containing arrays, hashmaps, or other non-comparable values cannot be used as set elements.
