package codegen

import (
	"strings"
	"testing"

	"github.com/mkaz/zenth/pkg/checker"
	"github.com/mkaz/zenth/pkg/lexer"
	"github.com/mkaz/zenth/pkg/parser"
)

// helper: lex → parse → check → generate, fail on error
func generate(t *testing.T, src string) string {
	t.Helper()
	l := lexer.New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}
	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	ch := checker.New()
	if err := ch.Check(prog); err != nil {
		t.Fatalf("check error: %v", err)
	}
	g := New()
	return g.Generate(prog)
}

// helper: assert that generated Go code contains a substring
func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("expected output to contain %q, got:\n%s", want, got)
	}
}

// helper: assert that generated Go code does NOT contain a substring
func assertNotContains(t *testing.T, got, notWant string) {
	t.Helper()
	if strings.Contains(got, notWant) {
		t.Errorf("expected output to NOT contain %q, got:\n%s", notWant, got)
	}
}

// ---------- Package and Imports ----------

func TestPackageMain(t *testing.T) {
	out := generate(t, `fn main() { Println("hello"); }`)
	assertContains(t, out, "package main")
}

func TestImportFmt(t *testing.T) {
	out := generate(t, `fn main() { Println("hello"); }`)
	assertContains(t, out, `"fmt"`)
}

func TestImportModule(t *testing.T) {
	out := generate(t, `import "strings"; fn main() { let x = strings.contains("hello", "ell"); }`)
	assertContains(t, out, `"strings"`)
}

func TestImportAlias(t *testing.T) {
	out := generate(t, `import "strings" as s; fn main() { let x = s.contains("hello", "ell"); }`)
	assertContains(t, out, `"strings"`)
}

// ---------- Function Declarations ----------

func TestFnMain(t *testing.T) {
	out := generate(t, `fn main() { Println("hello"); }`)
	assertContains(t, out, "func main()")
	assertContains(t, out, `fmt.Println("hello")`)
}

func TestFnWithReturn(t *testing.T) {
	out := generate(t, `fn add(a: int, b: int) -> int { return a + b; } fn main() { Println(Str(add(1, 2))); }`)
	assertContains(t, out, "func add(a int, b int) int")
	assertContains(t, out, "return a + b")
}

func TestFnNoParams(t *testing.T) {
	out := generate(t, `fn greet() { Println("hi"); } fn main() { greet(); }`)
	assertContains(t, out, "func greet()")
}

func TestFnMultipleParams(t *testing.T) {
	out := generate(t, `fn calc(a: int, b: int, c: int) -> int { return a + b + c; } fn main() { Println(Str(calc(1, 2, 3))); }`)
	assertContains(t, out, "func calc(a int, b int, c int) int")
}

// ---------- Variable Declarations ----------

func TestLetDecl(t *testing.T) {
	out := generate(t, `fn main() { let x = 5; Println(Str(x)); }`)
	assertContains(t, out, "x := 5")
}

func TestLetDeclWithType(t *testing.T) {
	out := generate(t, `fn main() { let x: int = 5; Println(Str(x)); }`)
	assertContains(t, out, "x")
	assertContains(t, out, "5")
}

func TestVarDecl(t *testing.T) {
	out := generate(t, `fn main() { var x = 5; x = 10; Println(Str(x)); }`)
	assertContains(t, out, "x := 5")
	assertContains(t, out, "x = 10")
}

func TestConstDecl(t *testing.T) {
	out := generate(t, `fn main() { const pi = 3.14; Println(Str(pi)); }`)
	assertContains(t, out, "pi := 3.14")
}

// ---------- Object Declarations ----------

func TestObjDecl(t *testing.T) {
	src := `
		obj Point {
			x: f64;
			y: f64;
		}
		fn main() {
			let p = Point(x=1.0, y=2.0);
			Println(Str(p.x));
		}
	`
	out := generate(t, src)
	assertContains(t, out, "type Point struct")
	assertContains(t, out, "X float64")
	assertContains(t, out, "Y float64")
}

func TestObjMethod(t *testing.T) {
	src := `
		obj Rect {
			w: f64;
			h: f64;

			fn area() -> f64 {
				return self.w * self.h;
			}
		}
		fn main() {
			let r = Rect(w=3.0, h=4.0);
			Println(Str(r.area()));
		}
	`
	out := generate(t, src)
	assertContains(t, out, "func (self *Rect) Area() float64")
	assertContains(t, out, "self.W * self.H")
}

// ---------- Enum Declarations ----------

func TestEnumDecl(t *testing.T) {
	src := `
		enum Color {
			Red;
			Green;
			Blue;
		}
		fn main() {
			let c = Color.Red;
			Println(Str(c));
		}
	`
	out := generate(t, src)
	assertContains(t, out, "type Color string")
	assertContains(t, out, `Color_Red Color = "Red"`)
	assertContains(t, out, `Color_Green Color = "Green"`)
	assertContains(t, out, `Color_Blue Color = "Blue"`)
}

// ---------- Control Flow ----------

func TestIfStmt(t *testing.T) {
	src := `fn main() { let x = 5; if x > 0 { Println("pos"); } else { Println("neg"); } }`
	out := generate(t, src)
	assertContains(t, out, "if x > 0")
	assertContains(t, out, "} else {")
}

func TestIfElseIfStmt(t *testing.T) {
	src := `fn main() { let x = 5; if x > 10 { Println("big"); } else if x > 3 { Println("med"); } else { Println("sm"); } }`
	out := generate(t, src)
	assertContains(t, out, "if x > 10")
	assertContains(t, out, "} else if x > 3")
	assertContains(t, out, "} else {")
}

func TestForCStyle(t *testing.T) {
	src := `fn main() { for var i = 0; i < 5; i++ { Println(Str(i)); } }`
	out := generate(t, src)
	assertContains(t, out, "for i := 0; i < 5; i++")
}

func TestForWhileStyle(t *testing.T) {
	src := `fn main() { var running = true; for running { running = false; } }`
	out := generate(t, src)
	assertContains(t, out, "for running")
}

func TestForInfinite(t *testing.T) {
	src := `fn main() { for { break; } }`
	out := generate(t, src)
	assertContains(t, out, "for {")
	assertContains(t, out, "break")
}

func TestForIn(t *testing.T) {
	src := `fn main() { let items = [1, 2, 3]; for item in items { Println(Str(item)); } }`
	out := generate(t, src)
	assertContains(t, out, "for _, item := range items")
}

func TestForInWithIndex(t *testing.T) {
	src := `fn main() { let items = [1, 2, 3]; for i, item in items { Println(Str(i)); } }`
	out := generate(t, src)
	assertContains(t, out, "for i, item := range items")
}

func TestMatchStmt(t *testing.T) {
	src := `fn main() { let x = 1; match x { 1 => Println("one"); 2 => Println("two"); _ => Println("other"); } }`
	out := generate(t, src)
	assertContains(t, out, "switch x")
	assertContains(t, out, "case 1:")
	assertContains(t, out, "case 2:")
	assertContains(t, out, "default:")
}

func TestLoopStmt(t *testing.T) {
	src := `fn main() { loop 3 { Println("hi"); } }`
	out := generate(t, src)
	assertContains(t, out, "< 3")
}

// ---------- Expressions ----------

func TestBinaryExpr(t *testing.T) {
	src := `fn main() { let x = 2 + 3; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "2 + 3")
}

func TestBinaryMul(t *testing.T) {
	src := `fn main() { let x = 4 * 5; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "4 * 5")
}

func TestUnaryNeg(t *testing.T) {
	src := `fn main() { let x = -5; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "-5")
}

func TestBooleanLogic(t *testing.T) {
	src := `fn main() { let a = true; let b = false; let c = a && b; let d = a || b; Println(Str(c)); Println(Str(d)); }`
	out := generate(t, src)
	assertContains(t, out, "a && b")
	assertContains(t, out, "a || b")
}

func TestStringConcat(t *testing.T) {
	src := `fn main() { let x = "hello" + " " + "world"; Println(x); }`
	out := generate(t, src)
	assertContains(t, out, `"hello" + " " + "world"`)
}

func TestComparisonOps(t *testing.T) {
	src := `fn main() { let a = 1 < 2; let b = 3 >= 2; let c = 1 == 1; let d = 1 != 2; Println(Str(a)); Println(Str(b)); Println(Str(c)); Println(Str(d)); }`
	out := generate(t, src)
	assertContains(t, out, "1 < 2")
	assertContains(t, out, "3 >= 2")
	assertContains(t, out, "1 == 1")
	assertContains(t, out, "1 != 2")
}

// ---------- Built-in Functions ----------

func TestPrintln(t *testing.T) {
	out := generate(t, `fn main() { Println("hello"); }`)
	assertContains(t, out, `fmt.Println("hello")`)
}

func TestPrint(t *testing.T) {
	out := generate(t, `fn main() { Print("hello"); }`)
	assertContains(t, out, `fmt.Print("hello")`)
}

func TestLen(t *testing.T) {
	src := `fn main() { let xs = [1, 2, 3]; let n = Len(xs); Println(Str(n)); }`
	out := generate(t, src)
	assertContains(t, out, "len(xs)")
}

func TestStrConversion(t *testing.T) {
	src := `fn main() { let x = 42; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "fmt.Sprint(x)")
}

// ---------- Array Literal ----------

func TestArrayLiteral(t *testing.T) {
	src := `fn main() { let xs = [1, 2, 3]; Println(Str(Len(xs))); }`
	out := generate(t, src)
	assertContains(t, out, "[]int{1, 2, 3}")
}

func TestEmptyArrayWithType(t *testing.T) {
	src := `fn main() { let xs: array(int) = []; Println(Str(Len(xs))); }`
	out := generate(t, src)
	assertContains(t, out, "[]int{}")
}

// ---------- String Interpolation ----------

func TestInterpString(t *testing.T) {
	src := `fn main() { let name = "world"; Println("hello {name}"); }`
	out := generate(t, src)
	assertContains(t, out, "fmt.Sprintf")
	assertContains(t, out, "name")
}

// ---------- Return Statement ----------

func TestReturnValue(t *testing.T) {
	src := `fn double(x: int) -> int { return x * 2; } fn main() { Println(Str(double(5))); }`
	out := generate(t, src)
	assertContains(t, out, "return x * 2")
}

func TestReturnBare(t *testing.T) {
	src := `fn noop() { return; } fn main() { noop(); }`
	out := generate(t, src)
	assertContains(t, out, "return")
}

// ---------- Assignment ----------

func TestAssignment(t *testing.T) {
	src := `fn main() { var x = 0; x = 42; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "x = 42")
}

func TestCompoundAssign(t *testing.T) {
	src := `fn main() { var x = 10; x += 5; x -= 2; x *= 3; x /= 2; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "x += 5")
	assertContains(t, out, "x -= 2")
	assertContains(t, out, "x *= 3")
	assertContains(t, out, "x /= 2")
}

func TestIncDec(t *testing.T) {
	src := `fn main() { var x = 0; x++; x--; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "x++")
	assertContains(t, out, "x--")
}

// ---------- Break and Continue ----------

func TestBreakContinue(t *testing.T) {
	src := `fn main() { for var i = 0; i < 10; i++ { if i == 5 { break; } if i == 3 { continue; } Println(Str(i)); } }`
	out := generate(t, src)
	assertContains(t, out, "break")
	assertContains(t, out, "continue")
}

// ---------- Object Constructor ----------

func TestObjConstructor(t *testing.T) {
	src := `
		obj Point {
			x: f64;
			y: f64;
		}
		fn main() {
			let p = Point(x=1.0, y=2.0);
			Println(Str(p.x));
		}
	`
	out := generate(t, src)
	assertContains(t, out, "&Point{")
	assertContains(t, out, "X: 1.0")
	assertContains(t, out, "Y: 2.0")
}

// ---------- Field Access ----------

func TestFieldAccess(t *testing.T) {
	src := `
		obj Point {
			x: f64;
			y: f64;
		}
		fn main() {
			let p = Point(x=1.0, y=2.0);
			Println(Str(p.x));
		}
	`
	out := generate(t, src)
	assertContains(t, out, "p.X")
}

// ---------- Index Access ----------

func TestIndexAccess(t *testing.T) {
	src := `fn main() { let xs = [10, 20, 30]; let v = xs[1]; Println(Str(v)); }`
	out := generate(t, src)
	assertContains(t, out, "xs[1]")
}

// ---------- Slice Expression ----------

func TestSliceExpr(t *testing.T) {
	src := `fn main() { let xs = [1, 2, 3, 4, 5]; let ys = xs[1:3]; Println(Str(Len(ys))); }`
	out := generate(t, src)
	assertContains(t, out, "xs[1:3]")
}

// ---------- Bool Literals ----------

func TestBoolLiterals(t *testing.T) {
	src := `fn main() { let a = true; let b = false; Println(Str(a)); Println(Str(b)); }`
	out := generate(t, src)
	assertContains(t, out, "true")
	assertContains(t, out, "false")
}

// ---------- Obj String Method ----------

func TestObjAutoString(t *testing.T) {
	src := `
		obj Point {
			x: f64;
			y: f64;
		}
		fn main() {
			let p = Point(x=1.0, y=2.0);
			Println(Str(p));
		}
	`
	out := generate(t, src)
	// Auto-generated String() method for the object
	assertContains(t, out, "func (self *Point) String() string")
}

// ---------- Interface Declarations ----------

func TestInterfaceDecl(t *testing.T) {
	src := `
		interface Greeter {
			fn greet() -> str;
		}
		fn main() { Println("hi"); }
	`
	out := generate(t, src)
	assertContains(t, out, "type Greeter interface")
	assertContains(t, out, "Greet() string")
}

// ---------- Type Alias ----------

func TestTypeAlias(t *testing.T) {
	src := `
		type ID = int;
		fn main() {
			let x: ID = 42;
			Println(Str(x));
		}
	`
	out := generate(t, src)
	// Type alias should resolve during codegen; the variable should use the underlying type
	assertContains(t, out, "42")
}

// ---------- Multiple Functions ----------

func TestMultipleFunctions(t *testing.T) {
	src := `
		fn square(x: int) -> int { return x * x; }
		fn main() { Println(Str(square(4))); }
	`
	out := generate(t, src)
	assertContains(t, out, "func square(x int) int")
	assertContains(t, out, "func main()")
}

// ---------- Nested Control Flow ----------

func TestNestedIfInFor(t *testing.T) {
	src := `
		fn main() {
			for var i = 0; i < 10; i++ {
				if i > 5 {
					Println("big");
				} else {
					Println("small");
				}
			}
		}
	`
	out := generate(t, src)
	assertContains(t, out, "for i := 0; i < 10; i++")
	assertContains(t, out, "if i > 5")
}

// ---------- Default Parameter Values ----------

func TestDefaultParams(t *testing.T) {
	src := `
		fn greet(name: str = "world") -> str {
			return "hello " + name;
		}
		fn main() {
			Println(greet());
		}
	`
	out := generate(t, src)
	assertContains(t, out, "func greet(name string) string")
}

// ---------- Empty Main ----------

func TestEmptyMain(t *testing.T) {
	out := generate(t, `fn main() {}`)
	assertContains(t, out, "func main()")
	assertContains(t, out, "package main")
}

// ---------- Modulo Operator ----------

func TestModuloOp(t *testing.T) {
	src := `fn main() { let x = 10 % 3; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "10 % 3")
}

// ---------- Enum with Explicit Values ----------

func TestEnumExplicitValues(t *testing.T) {
	src := `
		enum Status {
			Active = "active";
			Inactive = "inactive";
		}
		fn main() {
			let s = Status.Active;
			Println(Str(s));
		}
	`
	out := generate(t, src)
	assertContains(t, out, `Status_Active Status = "active"`)
	assertContains(t, out, `Status_Inactive Status = "inactive"`)
}

// ---------- Obj with Default Fields ----------

func TestObjDefaultField(t *testing.T) {
	src := `
		obj Config {
			debug: bool = false;
			name: str = "default";
		}
		fn main() {
			let c = Config();
			Println(Str(c.debug));
		}
	`
	out := generate(t, src)
	assertContains(t, out, "type Config struct")
}

// ---------- Unused Variable Suppression ----------

func TestLetUnusedSuppression(t *testing.T) {
	// let bindings generate _ = name to suppress Go unused-variable errors
	src := `fn main() { let x = 5; Println(Str(x)); }`
	out := generate(t, src)
	assertContains(t, out, "_ = x")
}

// ---------- Enum String Method ----------

func TestEnumStringMethod(t *testing.T) {
	src := `
		enum Dir { Up; Down; }
		fn main() {
			let d = Dir.Up;
			Println(Str(d));
		}
	`
	out := generate(t, src)
	assertContains(t, out, "func (e Dir) String() string")
	assertContains(t, out, "return string(e)")
}
