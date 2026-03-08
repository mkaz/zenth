# Testing

Zenth has a built-in test runner that discovers and runs test functions. Test files live in a `tests/` directory and use the `_test.zn` suffix.

## Writing Tests

Create a file ending with `_test.zn` in a `tests/` directory. Define test functions that start with `test_`:

```zenth
// tests/math_test.zn

fn test_addition() {
    assert_eq(2 + 3, 5);
}

fn test_negative_numbers() {
    assert_eq(-3 + 3, 0);
    assert(-1 < 0);
}
```

Test files do not need a `fn main()` -- the test runner generates one automatically.

## Assertions

Two assertion functions are available in all Zenth programs, not just test files:

### assert(condition)

Panics if `condition` is `false`, reporting the file and line number:

```zenth
assert(len("hello") == 5);
assert(x > 0);
assert(items.exists("key"));
```

On failure:
```
assert failed at tests/math_test.zn:4
```

### assert_eq(got, expected)

Panics if `got != expected`, reporting the file, line, and both values:

```zenth
assert_eq(add(2, 3), 5);
assert_eq(name.upper(), "ALICE");
assert_eq(len(items), 3);
```

On failure:
```
assert_eq failed at tests/math_test.zn:8
  expected: 5
       got: 4
```

`assert_eq` uses the `==` operator, so both values must be the same comparable type.

## Running Tests

Run all tests in the default `tests/` directory:

```sh
zenth test
```

Run tests in a specific directory:

```sh
zenth test path/to/tests/
```

Run a single test file:

```sh
zenth test tests/math_test.zn
```

## Output Format

The test runner prints results in plain text:

```
tests/math_test.zn
  test_addition ... PASS
  test_negative_numbers ... PASS
tests/strings_test.zn
  test_upper ... PASS
  test_split ... FAIL
    assert_eq failed at tests/strings_test.zn:12
      expected: 3
           got: 2

3 passed, 1 failed
```

The exit code is `0` if all tests pass, `1` if any test fails.

## Conventions

- Test files: `*_test.zn` in a `tests/` directory
- Test functions: `fn test_*()` with no parameters and no return type
- One `tests/` directory at the project root is the default location
- Use `assert` for boolean conditions, `assert_eq` for value comparison
- Helper functions (without the `test_` prefix) can be defined in test files and called from test functions

## Example: Testing an Object

```zenth
// tests/rectangle_test.zn

obj Rectangle {
    width: f64;
    height: f64;

    fn area() -> f64 {
        return self.width * self.height;
    }
}

fn test_area() {
    let r = Rectangle{ width: 3.0, height: 4.0 };
    assert_eq(r.area(), 12.0);
}

fn test_zero_area() {
    let r = Rectangle{ width: 0.0, height: 100.0 };
    assert_eq(r.area(), 0.0);
}
```

## Example: Testing with Collections

```zenth
// tests/collections_test.zn

fn test_array_operations() {
    let nums = [3, 1, 4, 1, 5];
    assert_eq(nums.max(), 5);
    assert_eq(nums.min(), 1);
    assert_eq(nums.sum(), 14);

    let sorted = nums.sorted();
    assert_eq(sorted[0], 1);
    assert_eq(sorted[4], 5);
}

fn test_hashmap_lookup() {
    var scores = hashmap(str, int);
    scores["alice"] = 95;
    scores["bob"] = 80;

    assert(scores.exists("alice"));
    assert_eq(scores["alice"], 95);
    assert_eq(len(scores), 2);
}
```
