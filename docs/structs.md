# Structs

Structs are user-defined types with named fields. Methods are defined inside the struct body using `self` to access fields.

## Declaring a Struct

```zenth
struct Point {
    x: f64;
    y: f64;
}
```

Fields are declared with `name: Type;` syntax.

## Creating Instances

Use struct literal syntax with `Name{ field: value }`:

```zenth
let origin = Point{ x: 0.0, y: 0.0 };
let corner = Point{ x: 3.0, y: 4.0 };
```

## Accessing Fields

Use dot notation:

```zenth
println(str(corner.x));  // 3
println(str(corner.y));  // 4
```

## Methods

Methods are defined inside the struct body. Use `self` to access the struct's fields and call its other methods:

```zenth
struct Rectangle {
    width: f64;
    height: f64;

    fn area() -> f64 {
        return self.width * self.height;
    }

    fn perimeter() -> f64 {
        return 2.0 * (self.width + self.height);
    }

    fn describe() {
        println("{self.width} x {self.height}, area = {self.area()}");
    }
}
```

## Calling Methods

```zenth
let r = Rectangle{ width: 10.0, height: 5.0 };
println("Area: {r.area()}");         // Area: 50
println("Perimeter: {r.perimeter()}"); // Perimeter: 30
r.describe();
```

## Methods with Parameters

Methods can take parameters in addition to `self`:

```zenth
struct Point {
    x: f64;
    y: f64;

    fn distance(other: Point) -> f64 {
        let dx = self.x - other.x;
        let dy = self.y - other.y;
        return dx + dy;
    }
}
```

## Returning Strings

A common pattern is to define a `string()` method for display:

```zenth
struct Point {
    x: f64;
    y: f64;

    fn string() -> str {
        return "({self.x}, {self.y})";
    }
}

fn main() {
    let p = Point{ x: 1.0, y: 2.0 };
    println(p.string());  // (1, 2)
}
```

## Complete Example

```zenth
struct Circle {
    radius: f64;

    fn area() -> f64 {
        return 3.14159 * self.radius * self.radius;
    }

    fn circumference() -> f64 {
        return 2.0 * 3.14159 * self.radius;
    }
}

fn main() {
    let c = Circle{ radius: 5.0 };
    println("Circle with radius {c.radius}");
    println("  Area: {c.area()}");
    println("  Circumference: {c.circumference()}");
}
```
