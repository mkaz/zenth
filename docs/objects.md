# Objects

Objects are user-defined types with named fields. Methods are defined inside the object body using `self` to access fields.

## Declaring an Object

```zenth
obj Point {
    x: Float;
    y: Float;
}
```

Fields are declared with `name: Type;` syntax.

## Default Field Values

Fields can have default values. If a default is provided, the field becomes optional in the constructor:

```zenth
obj Account {
    balance: Int = 0;
    interest: Float = 2.5;
}
```

## Creating Instances

Use constructor syntax with named arguments `Name(field=value)`:

```zenth
let origin = Point(x=0.0, y=0.0);
let corner = Point(x=3.0, y=4.0);
```

Fields with defaults can be omitted or overridden:

```zenth
let a1 = Account();                           // all defaults
let a2 = Account(balance=100);                // override one
let a3 = Account(interest=5.0, balance=1000); // override multiple
```

## Accessing Fields

Use dot notation:

```zenth
Println(corner.x);       // 3
Println(corner.y);       // 4
```

## Methods

Methods are defined inside the object body. Use `self` to access the object's fields and call its other methods:

```zenth
obj Rectangle {
    width: Float;
    height: Float;

    fn area() -> Float {
        return self.width * self.height;
    }

    fn perimeter() -> Float {
        return 2.0 * (self.width + self.height);
    }

    fn describe() {
        Println("{self.width} x {self.height}, area = {self.area()}");
    }
}
```

## Calling Methods

```zenth
let r = Rectangle(width=10.0, height=5.0);
Println("Area: {r.area()}");         // Area: 50
Println("Perimeter: {r.perimeter()}"); // Perimeter: 30
r.describe();
```

## Methods with Parameters

Methods can take parameters in addition to `self`:

```zenth
obj Point {
    x: Float;
    y: Float;

    fn distance(other: Point) -> Float {
        let dx = self.x - other.x;
        let dy = self.y - other.y;
        return Sqrt(dx * dx + dy * dy);
    }
}
```

## Returning Strings

By default, printing an object uses constructor-like output:

```zenth
obj Point {
    x: Int;
    y: Int;
    val: Str;
}

fn main() {
    let p = Point(x=1, y=2, val="#");
    Println(p);  // Point(x=1, y=2, val="#")
}
```

You can override this by defining a `string()` method:

```zenth
obj Point {
    x: Float;
    y: Float;

    fn string() -> Str {
        return "({self.x}, {self.y})";
    }
}

fn main() {
    let p = Point(x=1.0, y=2.0);
    Println(p.string());  // (1, 2)
}
```

## Complete Example

```zenth
obj Circle {
    radius: Float;

    fn area() -> Float {
        return 3.14159 * self.radius * self.radius;
    }

    fn circumference() -> Float {
        return 2.0 * 3.14159 * self.radius;
    }
}

fn main() {
    let c = Circle(radius=5.0);
    Println("Circle with radius {c.radius}");
    Println("  Area: {c.area()}");
    Println("  Circumference: {c.circumference()}");
}
```
