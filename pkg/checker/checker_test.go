package checker

import (
	"strings"
	"testing"

	"github.com/mkaz/zenth/pkg/lexer"
	"github.com/mkaz/zenth/pkg/parser"
)

// check is a helper that lexes, parses, and type-checks a Zenth source string.
func check(src string) error {
	l := lexer.New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		return err
	}
	p := parser.New(tokens)
	prog, err := p.Parse()
	if err != nil {
		return err
	}
	c := New()
	return c.Check(prog)
}

// ---------------------------------------------------------------------------
// 1. TestScope — test scope.go directly
// ---------------------------------------------------------------------------

func TestScopeDefineAndLookup(t *testing.T) {
	s := NewScope(nil)
	sym := &Symbol{Name: "x", Type: TypeInt}
	if ok := s.Define(sym); !ok {
		t.Fatal("Define should return true for first definition")
	}
	got := s.Lookup("x")
	if got == nil {
		t.Fatal("Lookup should find defined symbol")
	}
	if got.Name != "x" || !got.Type.Equals(TypeInt) {
		t.Fatalf("Lookup returned wrong symbol: %+v", got)
	}
}

func TestScopeLookupUndefined(t *testing.T) {
	s := NewScope(nil)
	if got := s.Lookup("missing"); got != nil {
		t.Fatalf("Lookup should return nil for undefined name, got %+v", got)
	}
}

func TestScopeDuplicateDefine(t *testing.T) {
	s := NewScope(nil)
	s.Define(&Symbol{Name: "x", Type: TypeInt})
	if ok := s.Define(&Symbol{Name: "x", Type: TypeStr}); ok {
		t.Fatal("Define should return false for duplicate in same scope")
	}
}

func TestScopeNestedLookup(t *testing.T) {
	parent := NewScope(nil)
	parent.Define(&Symbol{Name: "x", Type: TypeInt})

	child := NewScope(parent)
	got := child.Lookup("x")
	if got == nil {
		t.Fatal("child scope should find parent symbol")
	}
	if !got.Type.Equals(TypeInt) {
		t.Fatalf("expected int, got %s", got.Type)
	}
}

func TestScopeShadowing(t *testing.T) {
	parent := NewScope(nil)
	parent.Define(&Symbol{Name: "x", Type: TypeInt})

	child := NewScope(parent)
	child.Define(&Symbol{Name: "x", Type: TypeStr})

	got := child.Lookup("x")
	if got == nil {
		t.Fatal("child scope should find shadowed symbol")
	}
	if !got.Type.Equals(TypeStr) {
		t.Fatalf("expected str (shadowed), got %s", got.Type)
	}

	// Parent should still see int
	parentGot := parent.Lookup("x")
	if !parentGot.Type.Equals(TypeInt) {
		t.Fatalf("parent should still have int, got %s", parentGot.Type)
	}
}

// ---------------------------------------------------------------------------
// 2. TestTypes — test types.go
// ---------------------------------------------------------------------------

func TestBuiltinTypeString(t *testing.T) {
	tests := []struct {
		typ  ZType
		want string
	}{
		{TypeInt, "int"},
		{TypeI8, "i8"},
		{TypeI16, "i16"},
		{TypeI32, "i32"},
		{TypeI64, "i64"},
		{TypeU8, "u8"},
		{TypeU16, "u16"},
		{TypeU32, "u32"},
		{TypeU64, "u64"},
		{TypeF32, "f32"},
		{TypeF64, "f64"},
		{TypeBool, "bool"},
		{TypeStr, "str"},
		{TypeByte, "byte"},
		{TypeVoid, "void"},
		{TypeError, "error"},
		{TypeNil, "nil"},
		{TypeFile, "file"},
	}
	for _, tt := range tests {
		if got := tt.typ.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestObjTypeString(t *testing.T) {
	o := &ObjType{Name: "Point", Fields: map[string]ZType{"x": TypeF64}}
	if got := o.String(); got != "Point" {
		t.Errorf("ObjType.String() = %q, want %q", got, "Point")
	}
}

func TestSliceTypeString(t *testing.T) {
	s := &SliceType{Elem: TypeInt}
	if got := s.String(); got != "[]int" {
		t.Errorf("SliceType.String() = %q, want %q", got, "[]int")
	}
}

func TestHashmapTypeString(t *testing.T) {
	h := &HashmapType{Key: TypeStr, Value: TypeInt}
	if got := h.String(); got != "hashmap[str]int" {
		t.Errorf("HashmapType.String() = %q, want %q", got, "hashmap[str]int")
	}
}

func TestSetTypeString(t *testing.T) {
	s := &SetType{Elem: TypeStr}
	if got := s.String(); got != "set(str)" {
		t.Errorf("SetType.String() = %q, want %q", got, "set(str)")
	}
}

func TestTupleTypeString(t *testing.T) {
	tests := []struct {
		name string
		typ  *TupleType
		want string
	}{
		{"empty", &TupleType{}, "tuple[]"},
		{"positional", &TupleType{Elems: []ZType{TypeInt, TypeStr}}, "tuple[int, str]"},
		{"named", &TupleType{Elems: []ZType{TypeInt, TypeStr}, Names: []string{"a", "b"}}, "tuple[a: int, b: str]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.typ.String(); got != tt.want {
				t.Errorf("TupleType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEnumTypeString(t *testing.T) {
	e := &EnumType{Name: "Color", Variants: map[string]string{"Red": "Red"}}
	if got := e.String(); got != "Color" {
		t.Errorf("EnumType.String() = %q, want %q", got, "Color")
	}
}

func TestFuncTypeString(t *testing.T) {
	f := &FuncType{Params: []ZType{TypeInt}, Returns: TypeStr}
	if got := f.String(); got != "fn(...) -> str" {
		t.Errorf("FuncType.String() = %q, want %q", got, "fn(...) -> str")
	}
}

// Equals tests

func TestBuiltinTypeEquals(t *testing.T) {
	if !TypeInt.Equals(TypeInt) {
		t.Error("int should equal int")
	}
	if TypeInt.Equals(TypeStr) {
		t.Error("int should not equal str")
	}
	if TypeInt.Equals(&ObjType{Name: "int"}) {
		t.Error("BuiltinType should not equal ObjType even with same name")
	}
}

func TestObjTypeEquals(t *testing.T) {
	a := &ObjType{Name: "Point"}
	b := &ObjType{Name: "Point"}
	c := &ObjType{Name: "Rect"}
	if !a.Equals(b) {
		t.Error("same-name ObjType should be equal")
	}
	if a.Equals(c) {
		t.Error("different-name ObjType should not be equal")
	}
	if a.Equals(TypeInt) {
		t.Error("ObjType should not equal BuiltinType")
	}
}

func TestSliceTypeEquals(t *testing.T) {
	a := &SliceType{Elem: TypeInt}
	b := &SliceType{Elem: TypeInt}
	c := &SliceType{Elem: TypeStr}
	if !a.Equals(b) {
		t.Error("[]int should equal []int")
	}
	if a.Equals(c) {
		t.Error("[]int should not equal []str")
	}
	if a.Equals(TypeInt) {
		t.Error("SliceType should not equal BuiltinType")
	}
}

func TestHashmapTypeEquals(t *testing.T) {
	a := &HashmapType{Key: TypeStr, Value: TypeInt}
	b := &HashmapType{Key: TypeStr, Value: TypeInt}
	c := &HashmapType{Key: TypeStr, Value: TypeStr}
	d := &HashmapType{Key: TypeInt, Value: TypeInt}
	if !a.Equals(b) {
		t.Error("same hashmap types should be equal")
	}
	if a.Equals(c) {
		t.Error("different value types should not be equal")
	}
	if a.Equals(d) {
		t.Error("different key types should not be equal")
	}
	if a.Equals(TypeInt) {
		t.Error("HashmapType should not equal BuiltinType")
	}
}

func TestSetTypeEquals(t *testing.T) {
	a := &SetType{Elem: TypeStr}
	b := &SetType{Elem: TypeStr}
	c := &SetType{Elem: TypeInt}
	if !a.Equals(b) {
		t.Error("set(str) should equal set(str)")
	}
	if a.Equals(c) {
		t.Error("set(str) should not equal set(int)")
	}
	if a.Equals(TypeStr) {
		t.Error("SetType should not equal BuiltinType")
	}
}

func TestTupleTypeEquals(t *testing.T) {
	positional := &TupleType{Elems: []ZType{TypeInt, TypeStr}}
	positional2 := &TupleType{Elems: []ZType{TypeInt, TypeStr}}
	named := &TupleType{Elems: []ZType{TypeInt, TypeStr}, Names: []string{"a", "b"}}
	namedDiffNames := &TupleType{Elems: []ZType{TypeInt, TypeStr}, Names: []string{"x", "y"}}
	diffElems := &TupleType{Elems: []ZType{TypeStr, TypeStr}}
	diffLen := &TupleType{Elems: []ZType{TypeInt}}

	if !positional.Equals(positional2) {
		t.Error("identical positional tuples should be equal")
	}
	if positional.Equals(named) {
		t.Error("positional vs named should not be equal")
	}
	if named.Equals(namedDiffNames) {
		t.Error("different names should not be equal")
	}
	if positional.Equals(diffElems) {
		t.Error("different element types should not be equal")
	}
	if positional.Equals(diffLen) {
		t.Error("different lengths should not be equal")
	}
	if positional.Equals(TypeInt) {
		t.Error("TupleType should not equal BuiltinType")
	}
}

func TestEnumTypeEquals(t *testing.T) {
	a := &EnumType{Name: "Color"}
	b := &EnumType{Name: "Color"}
	c := &EnumType{Name: "Size"}
	if !a.Equals(b) {
		t.Error("same-name EnumType should be equal")
	}
	if a.Equals(c) {
		t.Error("different-name EnumType should not be equal")
	}
	if a.Equals(TypeInt) {
		t.Error("EnumType should not equal BuiltinType")
	}
}

func TestFuncTypeEquals(t *testing.T) {
	a := &FuncType{Params: []ZType{TypeInt}, Returns: TypeStr}
	b := &FuncType{Params: []ZType{TypeInt}, Returns: TypeStr}
	diffParams := &FuncType{Params: []ZType{TypeStr}, Returns: TypeStr}
	diffParamCount := &FuncType{Params: []ZType{TypeInt, TypeInt}, Returns: TypeStr}
	diffReturn := &FuncType{Params: []ZType{TypeInt}, Returns: TypeInt}

	if !a.Equals(b) {
		t.Error("identical FuncTypes should be equal")
	}
	if a.Equals(diffParams) {
		t.Error("different param types should not be equal")
	}
	if a.Equals(diffParamCount) {
		t.Error("different param counts should not be equal")
	}
	if a.Equals(diffReturn) {
		t.Error("different return types should not be equal")
	}
	if a.Equals(TypeInt) {
		t.Error("FuncType should not equal BuiltinType")
	}
}

// LookupBuiltinType

func TestLookupBuiltinType(t *testing.T) {
	known := map[string]ZType{
		"int":   TypeInt,
		"i8":    TypeI8,
		"i16":   TypeI16,
		"i32":   TypeI32,
		"i64":   TypeI64,
		"u8":    TypeU8,
		"u16":   TypeU16,
		"u32":   TypeU32,
		"u64":   TypeU64,
		"f32":   TypeF32,
		"f64":   TypeF64,
		"bool":  TypeBool,
		"str":   TypeStr,
		"byte":  TypeByte,
		"error": TypeError,
	}
	for name, want := range known {
		got := LookupBuiltinType(name)
		if got != want {
			t.Errorf("LookupBuiltinType(%q) = %v, want %v", name, got, want)
		}
	}
	if got := LookupBuiltinType("nonexistent"); got != nil {
		t.Errorf("LookupBuiltinType(\"nonexistent\") = %v, want nil", got)
	}
}

// IsNumeric, IsInteger, IsFloat

func TestIsNumeric(t *testing.T) {
	numeric := []ZType{TypeInt, TypeI8, TypeI16, TypeI32, TypeI64, TypeU8, TypeU16, TypeU32, TypeU64, TypeF32, TypeF64}
	for _, typ := range numeric {
		if !IsNumeric(typ) {
			t.Errorf("IsNumeric(%s) should be true", typ)
		}
	}
	nonNumeric := []ZType{TypeBool, TypeStr, TypeByte, TypeVoid, &ObjType{Name: "Foo"}}
	for _, typ := range nonNumeric {
		if IsNumeric(typ) {
			t.Errorf("IsNumeric(%s) should be false", typ)
		}
	}
}

func TestIsInteger(t *testing.T) {
	ints := []ZType{TypeInt, TypeI8, TypeI16, TypeI32, TypeI64, TypeU8, TypeU16, TypeU32, TypeU64}
	for _, typ := range ints {
		if !IsInteger(typ) {
			t.Errorf("IsInteger(%s) should be true", typ)
		}
	}
	nonInts := []ZType{TypeF32, TypeF64, TypeBool, TypeStr}
	for _, typ := range nonInts {
		if IsInteger(typ) {
			t.Errorf("IsInteger(%s) should be false", typ)
		}
	}
}

func TestIsFloat(t *testing.T) {
	if !IsFloat(TypeF32) {
		t.Error("IsFloat(f32) should be true")
	}
	if !IsFloat(TypeF64) {
		t.Error("IsFloat(f64) should be true")
	}
	nonFloats := []ZType{TypeInt, TypeI8, TypeBool, TypeStr, &ObjType{Name: "Foo"}}
	for _, typ := range nonFloats {
		if IsFloat(typ) {
			t.Errorf("IsFloat(%s) should be false", typ)
		}
	}
}

// PromoteNumeric

func TestPromoteNumeric(t *testing.T) {
	tests := []struct {
		name string
		a, b ZType
		want ZType
	}{
		{"same int", TypeInt, TypeInt, TypeInt},
		{"same f64", TypeF64, TypeF64, TypeF64},
		{"f32+int -> f32", TypeF32, TypeInt, TypeF32},
		{"int+f64 -> f64", TypeInt, TypeF64, TypeF64},
		{"f32+f64 -> f64", TypeF32, TypeF64, TypeF64},
		{"diff ints -> nil", TypeI8, TypeI16, nil},
		{"non-numeric -> nil", TypeStr, TypeInt, nil},
		{"both non-numeric -> nil", TypeStr, TypeBool, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PromoteNumeric(tt.a, tt.b)
			if tt.want == nil {
				if got != nil {
					t.Errorf("PromoteNumeric(%s, %s) = %s, want nil", tt.a, tt.b, got)
				}
			} else {
				if got == nil || !got.Equals(tt.want) {
					t.Errorf("PromoteNumeric(%s, %s) = %v, want %s", tt.a, tt.b, got, tt.want)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 3. TestCheckerSuccess — programs that should pass type checking
// ---------------------------------------------------------------------------

func TestCheckerSuccess(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{
			"let declaration",
			`let x = 5;`,
		},
		{
			"let with type annotation",
			`let x: int = 5;`,
		},
		{
			"var declaration",
			`var x = 10;`,
		},
		{
			"var with type annotation",
			`var x: int = 10;`,
		},
		{
			"const declaration",
			`const PI = 3;`,
		},
		{
			"var reassignment",
			`fn main() { var x = 1; x = 2; }`,
		},
		{
			"function with typed params and return",
			`fn add(a: int, b: int) -> int { return a + b; }`,
		},
		{
			"function void return",
			`fn greet(name: str) { println(name); }`,
		},
		{
			"if/else with bool condition",
			`fn main() { let x = true; if x { println("yes"); } else { println("no"); } }`,
		},
		{
			"C-style for loop",
			`fn main() { for var i = 0; i < 10; i++ { println(str(i)); } }`,
		},
		{
			"for-in over array",
			`fn main() { let items = [1, 2, 3]; for item in items { println(str(item)); } }`,
		},
		{
			"for-in with index",
			`fn main() { let items = [1, 2, 3]; for i, item in items { println(str(i)); } }`,
		},
		{
			"infinite for loop with break",
			`fn main() { for var i = 0; true; i++ { break; } }`,
		},
		{
			"match statement",
			`fn main() { let x = 1; match x { 1 => println("one"); 2 => println("two"); _ => println("other"); } }`,
		},
		{
			"obj with fields and method",
			"obj Rect {\n  width: f64;\n  height: f64;\n  fn area() -> f64 {\n    return self.width * self.height;\n  }\n}\nfn main() {\n  let r = Rect(width=3.0, height=4.0);\n  println(str(r.area()));\n}\n",
		},
		{
			"enum declaration and variant access",
			`enum Color { Red; Green; Blue; } fn main() { let c = Color.Red; }`,
		},
		{
			"enum comparison",
			`enum Color { Red; Green; Blue; } fn main() { let c = Color.Red; let same = c == Color.Green; }`,
		},
		{
			"array literal and indexing",
			`fn main() { let arr = [1, 2, 3]; let first = arr[0]; }`,
		},
		{
			"string operations",
			`let s = "hello" + " " + "world";`,
		},
		{
			"binary expressions with numeric types",
			`let a = 1 + 2; let b = 3 * 4; let c = 10 - 5; let d = 8 / 2;`,
		},
		{
			"type promotion int + f64",
			`let a = 1 + 2.0;`,
		},
		{
			"comparison operators",
			`fn main() { let a = 1 < 2; let b = 3 >= 3; let c = 4 == 4; let d = 5 != 6; }`,
		},
		{
			"logical operators",
			`fn main() { let a = true && false; let b = true || false; }`,
		},
		{
			"unary operators",
			`fn main() { let a = !true; let b = -5; }`,
		},
		{
			"nested scopes in function",
			`fn foo() { let x = 1; let y = 2; let z = x + y; }`,
		},
		{
			"function calling function",
			`fn double(n: int) -> int { return n * 2; } fn quad(n: int) -> int { return double(double(n)); }`,
		},
		{
			"built-in len",
			`fn main() { let s = "hello"; let n = len(s); }`,
		},
		{
			"built-in str conversion",
			`fn main() { let n = 42; let s = str(n); }`,
		},
		{
			"string for-in",
			`fn main() { for ch in "hello" { println(ch); } }`,
		},
		{
			"import statement",
			`import "fmt";`,
		},
		{
			"import with alias",
			`import "math" as m;`,
		},
		{
			"INT_MAX and INT_MIN constants",
			`let a = INT_MAX; let b = INT_MIN;`,
		},
		{
			"var compound assignment",
			`fn main() { var x = 10; x += 5; }`,
		},
		{
			"string compound assignment",
			`fn main() { var s = "hello"; s += " world"; }`,
		},
		{
			"discard with underscore",
			`let _ = 5;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := check(tt.src); err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 4. TestCheckerErrors — programs that should fail type checking
// ---------------------------------------------------------------------------

func TestCheckerErrors(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr string
	}{
		{
			"undefined identifier",
			`let x = y;`,
			"undefined identifier: y",
		},
		{
			"assign to immutable let",
			`fn main() { let x = 1; x = 2; }`,
			"cannot assign to immutable variable",
		},
		{
			"assign to constant",
			`fn main() { const X = 1; X = 2; }`,
			"cannot assign to immutable variable",
		},
		{
			"type mismatch in let annotation",
			`let x: int = "hello";`,
			"type mismatch",
		},
		{
			"type mismatch in var annotation",
			`var x: int = "hello";`,
			"type mismatch",
		},
		{
			"return type mismatch",
			`fn foo() -> int { return "hello"; }`,
			"return type mismatch",
		},
		{
			"unknown type in annotation",
			`let x: foobar = 1;`,
			"unknown type: foobar",
		},
		{
			"binary op type mismatch: str + int",
			`let x = "hello" + 5;`,
			"cannot add",
		},
		{
			"if condition not bool",
			`fn main() { if 42 { println("oops"); } }`,
			"if condition must be bool",
		},
		{
			"for condition not bool",
			`fn main() { for var i = 0; 42; i++ { break; } }`,
			"for condition must be bool",
		},
		{
			"logical op non-bool",
			`let x = 1 && 2;`,
			"logical operators require bool operands",
		},
		{
			"unary ! on non-bool",
			`let x = !5;`,
			"! requires bool",
		},
		{
			"unary - on non-numeric",
			`let x = -"hello";`,
			"- requires numeric type",
		},
		{
			"comparison type mismatch",
			`let x = "hello" < "world";`,
			"cannot compare",
		},
		{
			"function wrong arg count",
			`fn foo(a: int, b: int) -> int { return a + b; } let x = foo(1);`,
			"expects 2 arguments, got 1",
		},
		{
			"function arg type mismatch",
			`fn foo(a: int) -> int { return a; } let x = foo("bad");`,
			"has type str, expected int",
		},
		{
			"enum bad variant",
			`enum Color { Red; Green; Blue; } let c = Color.Yellow;`,
			"has no variant 'Yellow'",
		},
		{
			"assign type mismatch to var",
			`fn main() { var x = 1; x = "hello"; }`,
			"type mismatch",
		},
		{
			"compound assign on non-numeric",
			`fn main() { var x = "hello"; x -= "world"; }`,
			"compound assignment requires numeric type",
		},
		{
			"array element type mismatch",
			`let arr = [1, "two", 3];`,
			"array element type mismatch",
		},
		{
			"inc/dec on non-numeric",
			`fn main() { var s = "hello"; s++; }`,
			"++/-- requires numeric type",
		},
		{
			"cannot iterate non-iterable",
			`fn main() { let x = 5; for item in x { println(str(item)); } }`,
			"cannot iterate over",
		},
		{
			"function missing return value",
			`fn foo() -> int { return; }`,
			"function expects return type",
		},
		{
			"subtract on strings",
			`let x = "hello" - "world";`,
			"cannot use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := check(tt.src)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 5. TestNewChecker — verify built-in functions and constants are registered
// ---------------------------------------------------------------------------

func TestNewCheckerBuiltinConstants(t *testing.T) {
	c := New()
	intMax := c.scope.Lookup("INT_MAX")
	if intMax == nil {
		t.Fatal("INT_MAX should be defined in global scope")
	}
	if !intMax.Type.Equals(TypeInt) {
		t.Errorf("INT_MAX should be int, got %s", intMax.Type)
	}
	if !intMax.IsConst {
		t.Error("INT_MAX should be marked as const")
	}

	intMin := c.scope.Lookup("INT_MIN")
	if intMin == nil {
		t.Fatal("INT_MIN should be defined in global scope")
	}
	if !intMin.Type.Equals(TypeInt) {
		t.Errorf("INT_MIN should be int, got %s", intMin.Type)
	}
	if !intMin.IsConst {
		t.Error("INT_MIN should be marked as const")
	}
}

func TestNewCheckerBuiltinFunctions(t *testing.T) {
	c := New()
	builtins := []struct {
		name       string
		wantReturn ZType
	}{
		{"print", TypeVoid},
		{"println", TypeVoid},
		{"len", TypeInt},
		{"str", TypeStr},
		{"exit", TypeVoid},
		{"range", &SliceType{Elem: TypeInt}},
		{"rangei", &SliceType{Elem: TypeInt}},
		{"file", TypeFile},
		{"int", TypeInt},
		{"f64", TypeF64},
		{"flag", TypeVoid},
	}
	for _, tt := range builtins {
		t.Run(tt.name, func(t *testing.T) {
			info, ok := c.funcs[tt.name]
			if !ok {
				t.Fatalf("built-in function %q not registered", tt.name)
			}
			if !info.Return.Equals(tt.wantReturn) {
				t.Errorf("built-in %q return type = %s, want %s", tt.name, info.Return, tt.wantReturn)
			}
		})
	}
}

func TestNewCheckerKnownModules(t *testing.T) {
	c := New()
	modules := []string{"fmt", "math", "os", "strings", "str", "io"}
	for _, mod := range modules {
		if !c.modules[mod] {
			t.Errorf("module %q should be registered", mod)
		}
	}
}
