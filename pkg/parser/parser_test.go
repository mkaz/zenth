package parser

import (
	"testing"

	"github.com/mkaz/zenth/pkg/ast"
	"github.com/mkaz/zenth/pkg/lexer"
	"github.com/mkaz/zenth/pkg/token"
)

// helper: lex and parse source, fail on error
func parse(t *testing.T, src string) *ast.Program {
	t.Helper()
	l := lexer.New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}
	p := New(tokens)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return prog
}

// helper: lex and parse source, expect a parse error
func parseExpectError(t *testing.T, src string) {
	t.Helper()
	l := lexer.New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		// lex error counts as expected failure
		return
	}
	p := New(tokens)
	_, err = p.Parse()
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

// ---------- Variable Declarations ----------

func TestLetDecl(t *testing.T) {
	prog := parse(t, `let x = 5;`)
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}
	stmt, ok := prog.Stmts[0].(*ast.LetStmt)
	if !ok {
		t.Fatalf("expected LetStmt, got %T", prog.Stmts[0])
	}
	if stmt.Name != "x" {
		t.Errorf("expected name 'x', got %q", stmt.Name)
	}
	if stmt.Type != nil {
		t.Errorf("expected no explicit type, got %v", stmt.Type)
	}
}

func TestLetDeclWithType(t *testing.T) {
	prog := parse(t, `let x: int = 5;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	if stmt.Name != "x" {
		t.Errorf("expected name 'x', got %q", stmt.Name)
	}
	if stmt.Type == nil || stmt.Type.Name != "int" {
		t.Errorf("expected type 'int', got %v", stmt.Type)
	}
}

func TestVarDecl(t *testing.T) {
	prog := parse(t, `var counter = 0;`)
	stmt, ok := prog.Stmts[0].(*ast.VarStmt)
	if !ok {
		t.Fatalf("expected VarStmt, got %T", prog.Stmts[0])
	}
	if stmt.Name != "counter" {
		t.Errorf("expected name 'counter', got %q", stmt.Name)
	}
}

func TestConstDecl(t *testing.T) {
	prog := parse(t, `const PI = 3.14;`)
	stmt, ok := prog.Stmts[0].(*ast.ConstStmt)
	if !ok {
		t.Fatalf("expected ConstStmt, got %T", prog.Stmts[0])
	}
	if stmt.Name != "PI" {
		t.Errorf("expected name 'PI', got %q", stmt.Name)
	}
}

// ---------- Function Declarations ----------

func TestFnDeclNoReturn(t *testing.T) {
	prog := parse(t, `fn greet(name: str) { Println(name); }`)
	fn, ok := prog.Stmts[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("expected FnDecl, got %T", prog.Stmts[0])
	}
	if fn.Name != "greet" {
		t.Errorf("expected name 'greet', got %q", fn.Name)
	}
	if len(fn.Params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(fn.Params))
	}
	if fn.Params[0].Name != "name" {
		t.Errorf("expected param name 'name', got %q", fn.Params[0].Name)
	}
	if fn.Params[0].Type.Name != "str" {
		t.Errorf("expected param type 'str', got %q", fn.Params[0].Type.Name)
	}
	if fn.ReturnType != nil {
		t.Errorf("expected no return type, got %v", fn.ReturnType)
	}
}

func TestFnDeclWithReturn(t *testing.T) {
	prog := parse(t, `fn add(a: int, b: int) -> int { return a + b; }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	if fn.Name != "add" {
		t.Errorf("expected name 'add', got %q", fn.Name)
	}
	if len(fn.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Params))
	}
	if fn.ReturnType == nil || fn.ReturnType.Name != "int" {
		t.Errorf("expected return type 'int', got %v", fn.ReturnType)
	}
}

func TestFnDeclNoParams(t *testing.T) {
	prog := parse(t, `fn main() { Println("hello"); }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	if fn.Name != "main" {
		t.Errorf("expected name 'main', got %q", fn.Name)
	}
	if len(fn.Params) != 0 {
		t.Errorf("expected 0 params, got %d", len(fn.Params))
	}
}

func TestFnDeclDefaultParam(t *testing.T) {
	prog := parse(t, `fn greet(name: str = "world") { Println(name); }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	if len(fn.Params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(fn.Params))
	}
	if fn.Params[0].Default == nil {
		t.Error("expected default value for param, got nil")
	}
}

// ---------- Object Declarations ----------

func TestObjDecl(t *testing.T) {
	src := `obj Point {
		x: f64;
		y: f64;
	}`
	prog := parse(t, src)
	obj, ok := prog.Stmts[0].(*ast.ObjDecl)
	if !ok {
		t.Fatalf("expected ObjDecl, got %T", prog.Stmts[0])
	}
	if obj.Name != "Point" {
		t.Errorf("expected name 'Point', got %q", obj.Name)
	}
	if len(obj.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(obj.Fields))
	}
	if obj.Fields[0].Name != "x" || obj.Fields[0].Type.Name != "f64" {
		t.Errorf("expected field 'x: f64', got '%s: %s'", obj.Fields[0].Name, obj.Fields[0].Type.Name)
	}
}

func TestObjDeclWithMethod(t *testing.T) {
	src := `obj Rect {
		w: f64;
		h: f64;

		fn area() -> f64 {
			return self.w * self.h;
		}
	}`
	prog := parse(t, src)
	obj := prog.Stmts[0].(*ast.ObjDecl)
	if len(obj.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(obj.Methods))
	}
	if obj.Methods[0].Name != "area" {
		t.Errorf("expected method 'area', got %q", obj.Methods[0].Name)
	}
	if obj.Methods[0].ReturnType == nil || obj.Methods[0].ReturnType.Name != "f64" {
		t.Errorf("expected return type 'f64'")
	}
}

// ---------- Enum Declaration ----------

func TestEnumDecl(t *testing.T) {
	src := `enum Color {
		Red;
		Green;
		Blue;
	}`
	prog := parse(t, src)
	e, ok := prog.Stmts[0].(*ast.EnumDecl)
	if !ok {
		t.Fatalf("expected EnumDecl, got %T", prog.Stmts[0])
	}
	if e.Name != "Color" {
		t.Errorf("expected name 'Color', got %q", e.Name)
	}
	if len(e.Variants) != 3 {
		t.Fatalf("expected 3 variants, got %d", len(e.Variants))
	}
	if e.Variants[0].Name != "Red" {
		t.Errorf("expected first variant 'Red', got %q", e.Variants[0].Name)
	}
}

// ---------- Import Declarations ----------

func TestImportDecl(t *testing.T) {
	prog := parse(t, `import "fmt";`)
	imp, ok := prog.Stmts[0].(*ast.ImportDecl)
	if !ok {
		t.Fatalf("expected ImportDecl, got %T", prog.Stmts[0])
	}
	if imp.Path != "fmt" {
		t.Errorf("expected path 'fmt', got %q", imp.Path)
	}
	if imp.Alias != "" {
		t.Errorf("expected no alias, got %q", imp.Alias)
	}
}

func TestImportDeclWithAlias(t *testing.T) {
	prog := parse(t, `import "math" as m;`)
	imp := prog.Stmts[0].(*ast.ImportDecl)
	if imp.Path != "math" {
		t.Errorf("expected path 'math', got %q", imp.Path)
	}
	if imp.Alias != "m" {
		t.Errorf("expected alias 'm', got %q", imp.Alias)
	}
}

// ---------- Control Flow ----------

func TestIfStmt(t *testing.T) {
	src := `fn main() { if x > 0 { Println("pos"); } else { Println("neg"); } }`
	prog := parse(t, src)
	fn := prog.Stmts[0].(*ast.FnDecl)
	ifStmt, ok := fn.Body.Stmts[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("expected IfStmt, got %T", fn.Body.Stmts[0])
	}
	if ifStmt.Condition == nil {
		t.Error("expected condition, got nil")
	}
	if ifStmt.Body == nil {
		t.Error("expected body, got nil")
	}
	if ifStmt.Else == nil {
		t.Error("expected else branch, got nil")
	}
}

func TestIfElseIfStmt(t *testing.T) {
	src := `fn main() { if x > 10 { Println("big"); } else if x > 5 { Println("med"); } else { Println("sm"); } }`
	prog := parse(t, src)
	fn := prog.Stmts[0].(*ast.FnDecl)
	ifStmt := fn.Body.Stmts[0].(*ast.IfStmt)
	// else branch should be another IfStmt
	elseIf, ok := ifStmt.Else.(*ast.IfStmt)
	if !ok {
		t.Fatalf("expected else-if to be IfStmt, got %T", ifStmt.Else)
	}
	if elseIf.Else == nil {
		t.Error("expected else branch on else-if, got nil")
	}
}

func TestForCStyle(t *testing.T) {
	src := `fn main() { for var i = 0; i < 10; i++ { Println(Str(i)); } }`
	prog := parse(t, src)
	fn := prog.Stmts[0].(*ast.FnDecl)
	forStmt, ok := fn.Body.Stmts[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("expected ForStmt, got %T", fn.Body.Stmts[0])
	}
	if forStmt.Init == nil {
		t.Error("expected init clause")
	}
	if forStmt.Condition == nil {
		t.Error("expected condition")
	}
	if forStmt.Post == nil {
		t.Error("expected post clause")
	}
}

func TestForInStmt(t *testing.T) {
	src := `fn main() { for item in items { Println(item); } }`
	prog := parse(t, src)
	fn := prog.Stmts[0].(*ast.FnDecl)
	forIn, ok := fn.Body.Stmts[0].(*ast.ForInStmt)
	if !ok {
		t.Fatalf("expected ForInStmt, got %T", fn.Body.Stmts[0])
	}
	if forIn.Value != "item" {
		t.Errorf("expected value 'item', got %q", forIn.Value)
	}
	if forIn.Index != "" {
		t.Errorf("expected no index, got %q", forIn.Index)
	}
}

func TestForInWithIndex(t *testing.T) {
	src := `fn main() { for i, item in items { Println(Str(i)); } }`
	prog := parse(t, src)
	fn := prog.Stmts[0].(*ast.FnDecl)
	forIn := fn.Body.Stmts[0].(*ast.ForInStmt)
	if forIn.Index != "i" {
		t.Errorf("expected index 'i', got %q", forIn.Index)
	}
	if forIn.Value != "item" {
		t.Errorf("expected value 'item', got %q", forIn.Value)
	}
}

func TestForWhileStyle(t *testing.T) {
	src := `fn main() { for running { process(); } }`
	prog := parse(t, src)
	fn := prog.Stmts[0].(*ast.FnDecl)
	forStmt, ok := fn.Body.Stmts[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("expected ForStmt, got %T", fn.Body.Stmts[0])
	}
	if forStmt.Condition == nil {
		t.Error("expected condition for while-style loop")
	}
	if forStmt.Init != nil {
		t.Error("expected no init for while-style loop")
	}
}

func TestForInfinite(t *testing.T) {
	src := `fn main() { for { break; } }`
	prog := parse(t, src)
	fn := prog.Stmts[0].(*ast.FnDecl)
	forStmt, ok := fn.Body.Stmts[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("expected ForStmt, got %T", fn.Body.Stmts[0])
	}
	if forStmt.Condition != nil {
		t.Error("expected no condition for infinite loop")
	}
}

func TestMatchStmt(t *testing.T) {
	src := `fn main() { match day { 1 => Println("Mon"); 2 => Println("Tue"); _ => Println("Other"); } }`
	prog := parse(t, src)
	fn := prog.Stmts[0].(*ast.FnDecl)
	matchStmt, ok := fn.Body.Stmts[0].(*ast.MatchStmt)
	if !ok {
		t.Fatalf("expected MatchStmt, got %T", fn.Body.Stmts[0])
	}
	if matchStmt.Subject == nil {
		t.Error("expected subject expression")
	}
	if len(matchStmt.Arms) != 3 {
		t.Fatalf("expected 3 arms, got %d", len(matchStmt.Arms))
	}
}

// ---------- Expressions ----------

func TestBinaryExpr(t *testing.T) {
	prog := parse(t, `let x = 1 + 2 * 3;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	// Should parse as 1 + (2 * 3) due to precedence
	bin, ok := stmt.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", stmt.Value)
	}
	if bin.Op != token.Plus {
		t.Errorf("expected Plus at top level, got %s", bin.Op)
	}
	// Right side should be 2 * 3
	right, ok := bin.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr on right, got %T", bin.Right)
	}
	if right.Op != token.Star {
		t.Errorf("expected Star on right, got %s", right.Op)
	}
}

func TestUnaryExpr(t *testing.T) {
	prog := parse(t, `let x = -5;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	unary, ok := stmt.Value.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", stmt.Value)
	}
	if unary.Op != token.Minus {
		t.Errorf("expected Minus, got %s", unary.Op)
	}
}

func TestBooleanNot(t *testing.T) {
	prog := parse(t, `let x = !flag;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	unary, ok := stmt.Value.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("expected UnaryExpr, got %T", stmt.Value)
	}
	if unary.Op != token.Not {
		t.Errorf("expected Not, got %s", unary.Op)
	}
}

func TestCallExpr(t *testing.T) {
	prog := parse(t, `fn main() { Println("hello", "world"); }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	exprStmt, ok := fn.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected ExprStmt, got %T", fn.Body.Stmts[0])
	}
	call, ok := exprStmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", exprStmt.Expr)
	}
	if len(call.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(call.Args))
	}
}

func TestFieldExpr(t *testing.T) {
	prog := parse(t, `let x = point.x;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	field, ok := stmt.Value.(*ast.FieldExpr)
	if !ok {
		t.Fatalf("expected FieldExpr, got %T", stmt.Value)
	}
	if field.Field != "x" {
		t.Errorf("expected field 'x', got %q", field.Field)
	}
}

func TestIndexExpr(t *testing.T) {
	prog := parse(t, `let x = items[0];`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	idx, ok := stmt.Value.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("expected IndexExpr, got %T", stmt.Value)
	}
	if idx.Object == nil || idx.Index == nil {
		t.Error("expected non-nil object and index")
	}
}

func TestSliceExpr(t *testing.T) {
	prog := parse(t, `let x = items[1:3];`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	slice, ok := stmt.Value.(*ast.SliceExpr)
	if !ok {
		t.Fatalf("expected SliceExpr, got %T", stmt.Value)
	}
	if slice.Low == nil || slice.High == nil {
		t.Error("expected both low and high")
	}
}

func TestArrayLiteral(t *testing.T) {
	prog := parse(t, `let xs = [1, 2, 3];`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	arr, ok := stmt.Value.(*ast.ArrayLitExpr)
	if !ok {
		t.Fatalf("expected ArrayLitExpr, got %T", stmt.Value)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr.Elements))
	}
}

func TestBoolLiterals(t *testing.T) {
	prog := parse(t, `let a = true; let b = false;`)
	a := prog.Stmts[0].(*ast.LetStmt)
	boolA, ok := a.Value.(*ast.BoolLitExpr)
	if !ok {
		t.Fatalf("expected BoolLitExpr, got %T", a.Value)
	}
	if !boolA.Value {
		t.Error("expected true")
	}
	b := prog.Stmts[1].(*ast.LetStmt)
	boolB := b.Value.(*ast.BoolLitExpr)
	if boolB.Value {
		t.Error("expected false")
	}
}

func TestNilLiteral(t *testing.T) {
	prog := parse(t, `let x = nil;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	_, ok := stmt.Value.(*ast.NilExpr)
	if !ok {
		t.Fatalf("expected NilExpr, got %T", stmt.Value)
	}
}

func TestIntLiteral(t *testing.T) {
	prog := parse(t, `let x = 42;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	intLit, ok := stmt.Value.(*ast.IntLitExpr)
	if !ok {
		t.Fatalf("expected IntLitExpr, got %T", stmt.Value)
	}
	if intLit.Value != 42 {
		t.Errorf("expected 42, got %d", intLit.Value)
	}
}

func TestFloatLiteral(t *testing.T) {
	prog := parse(t, `let x = 3.14;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	floatLit, ok := stmt.Value.(*ast.FloatLitExpr)
	if !ok {
		t.Fatalf("expected FloatLitExpr, got %T", stmt.Value)
	}
	if floatLit.Value != 3.14 {
		t.Errorf("expected 3.14, got %f", floatLit.Value)
	}
}

func TestStringLiteral(t *testing.T) {
	prog := parse(t, `let x = "hello";`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	strLit, ok := stmt.Value.(*ast.StringLitExpr)
	if !ok {
		t.Fatalf("expected StringLitExpr, got %T", stmt.Value)
	}
	if strLit.Value != "hello" {
		t.Errorf("expected 'hello', got %q", strLit.Value)
	}
}

// ---------- Statements ----------

func TestReturnStmt(t *testing.T) {
	prog := parse(t, `fn foo() -> int { return 42; }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	ret, ok := fn.Body.Stmts[0].(*ast.ReturnStmt)
	if !ok {
		t.Fatalf("expected ReturnStmt, got %T", fn.Body.Stmts[0])
	}
	if ret.Value == nil {
		t.Error("expected return value")
	}
}

func TestReturnBare(t *testing.T) {
	prog := parse(t, `fn foo() { return; }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	ret := fn.Body.Stmts[0].(*ast.ReturnStmt)
	if ret.Value != nil {
		t.Errorf("expected bare return (nil value), got %T", ret.Value)
	}
}

func TestAssignStmt(t *testing.T) {
	prog := parse(t, `fn main() { var x = 0; x = 5; }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	assign, ok := fn.Body.Stmts[1].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("expected AssignStmt, got %T", fn.Body.Stmts[1])
	}
	if assign.Op != token.Assign {
		t.Errorf("expected Assign, got %s", assign.Op)
	}
}

func TestCompoundAssign(t *testing.T) {
	tests := []struct {
		src  string
		op   token.Type
		name string
	}{
		{`fn main() { var x = 0; x += 1; }`, token.PlusAssign, "PlusAssign"},
		{`fn main() { var x = 0; x -= 1; }`, token.MinusAssign, "MinusAssign"},
		{`fn main() { var x = 1; x *= 2; }`, token.StarAssign, "StarAssign"},
		{`fn main() { var x = 4; x /= 2; }`, token.SlashAssign, "SlashAssign"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog := parse(t, tt.src)
			fn := prog.Stmts[0].(*ast.FnDecl)
			assign := fn.Body.Stmts[1].(*ast.AssignStmt)
			if assign.Op != tt.op {
				t.Errorf("expected %s, got %s", tt.op, assign.Op)
			}
		})
	}
}

func TestIncDecStmt(t *testing.T) {
	prog := parse(t, `fn main() { var x = 0; x++; x--; }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	inc, ok := fn.Body.Stmts[1].(*ast.IncDecStmt)
	if !ok {
		t.Fatalf("expected IncDecStmt, got %T", fn.Body.Stmts[1])
	}
	if inc.Op != token.PlusPlus {
		t.Errorf("expected PlusPlus, got %s", inc.Op)
	}
	dec := fn.Body.Stmts[2].(*ast.IncDecStmt)
	if dec.Op != token.MinusMinus {
		t.Errorf("expected MinusMinus, got %s", dec.Op)
	}
}

func TestBreakContinue(t *testing.T) {
	prog := parse(t, `fn main() { for { break; continue; } }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	forStmt := fn.Body.Stmts[0].(*ast.ForStmt)
	_, ok := forStmt.Body.Stmts[0].(*ast.BreakStmt)
	if !ok {
		t.Fatalf("expected BreakStmt, got %T", forStmt.Body.Stmts[0])
	}
	_, ok = forStmt.Body.Stmts[1].(*ast.ContinueStmt)
	if !ok {
		t.Fatalf("expected ContinueStmt, got %T", forStmt.Body.Stmts[1])
	}
}

// ---------- Type Expressions ----------

func TestArrayType(t *testing.T) {
	prog := parse(t, `let xs: array(int) = [1, 2];`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	if stmt.Type == nil {
		t.Fatal("expected type annotation")
	}
	if !stmt.Type.IsSlice {
		t.Error("expected IsSlice=true")
	}
}

func TestTupleDestruct(t *testing.T) {
	prog := parse(t, `let (a, b) = get_pair();`)
	td, ok := prog.Stmts[0].(*ast.TupleDestructStmt)
	if !ok {
		t.Fatalf("expected TupleDestructStmt, got %T", prog.Stmts[0])
	}
	if len(td.Names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(td.Names))
	}
	if td.Names[0] != "a" || td.Names[1] != "b" {
		t.Errorf("expected (a, b), got (%s, %s)", td.Names[0], td.Names[1])
	}
}

// ---------- Interface Declaration ----------

func TestInterfaceDecl(t *testing.T) {
	src := `interface Shape {
		fn area() -> f64;
		fn perimeter() -> f64;
	}`
	prog := parse(t, src)
	iface, ok := prog.Stmts[0].(*ast.InterfaceDecl)
	if !ok {
		t.Fatalf("expected InterfaceDecl, got %T", prog.Stmts[0])
	}
	if iface.Name != "Shape" {
		t.Errorf("expected name 'Shape', got %q", iface.Name)
	}
	if len(iface.Methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(iface.Methods))
	}
	if iface.Methods[0].Name != "area" {
		t.Errorf("expected method 'area', got %q", iface.Methods[0].Name)
	}
}

// ---------- Type Alias ----------

func TestTypeAlias(t *testing.T) {
	prog := parse(t, `type ID = int;`)
	ta, ok := prog.Stmts[0].(*ast.TypeAliasDecl)
	if !ok {
		t.Fatalf("expected TypeAliasDecl, got %T", prog.Stmts[0])
	}
	if ta.Name != "ID" {
		t.Errorf("expected name 'ID', got %q", ta.Name)
	}
	if ta.Type == nil || ta.Type.Name != "int" {
		t.Errorf("expected type 'int'")
	}
}

// ---------- Logical Operators ----------

func TestLogicalOperators(t *testing.T) {
	prog := parse(t, `let x = a && b || c;`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	// || has lower precedence than &&, so top should be ||
	bin, ok := stmt.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr, got %T", stmt.Value)
	}
	if bin.Op != token.Or {
		t.Errorf("expected Or at top level, got %s", bin.Op)
	}
	left, ok := bin.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected BinaryExpr on left, got %T", bin.Left)
	}
	if left.Op != token.And {
		t.Errorf("expected And on left, got %s", left.Op)
	}
}

// ---------- Comparison Operators ----------

func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		sym string
		op  token.Type
	}{
		{"==", token.Eq},
		{"!=", token.Neq},
		{"<", token.Lt},
		{">", token.Gt},
		{"<=", token.Lte},
		{">=", token.Gte},
	}
	for _, tt := range tests {
		t.Run(tt.sym, func(t *testing.T) {
			src := `let x = a ` + tt.sym + ` b;`
			prog := parse(t, src)
			stmt := prog.Stmts[0].(*ast.LetStmt)
			bin, ok := stmt.Value.(*ast.BinaryExpr)
			if !ok {
				t.Fatalf("expected BinaryExpr, got %T", stmt.Value)
			}
			if bin.Op != tt.op {
				t.Errorf("expected %s, got %s", tt.op, bin.Op)
			}
		})
	}
}

// ---------- Multiple Statements ----------

func TestMultipleTopLevel(t *testing.T) {
	src := `
		let x = 1;
		let y = 2;
		fn main() { Println(Str(x + y)); }
	`
	prog := parse(t, src)
	if len(prog.Stmts) != 3 {
		t.Fatalf("expected 3 top-level statements, got %d", len(prog.Stmts))
	}
	_, ok1 := prog.Stmts[0].(*ast.LetStmt)
	_, ok2 := prog.Stmts[1].(*ast.LetStmt)
	_, ok3 := prog.Stmts[2].(*ast.FnDecl)
	if !ok1 || !ok2 || !ok3 {
		t.Error("unexpected statement types")
	}
}

// ---------- Loop Statement ----------

func TestLoopStmt(t *testing.T) {
	prog := parse(t, `fn main() { loop 5 { Println("hi"); } }`)
	fn := prog.Stmts[0].(*ast.FnDecl)
	loop, ok := fn.Body.Stmts[0].(*ast.LoopStmt)
	if !ok {
		t.Fatalf("expected LoopStmt, got %T", fn.Body.Stmts[0])
	}
	if loop.Count == nil {
		t.Error("expected count expression")
	}
}

// ---------- Interpolated Strings ----------

func TestInterpStringExpr(t *testing.T) {
	prog := parse(t, `let x = "hello {name}";`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	interp, ok := stmt.Value.(*ast.InterpStringExpr)
	if !ok {
		t.Fatalf("expected InterpStringExpr, got %T", stmt.Value)
	}
	if len(interp.Parts) == 0 {
		t.Error("expected at least one interp part")
	}
}

// ---------- Method Chaining ----------

func TestMethodChain(t *testing.T) {
	prog := parse(t, `let x = items.length();`)
	stmt := prog.Stmts[0].(*ast.LetStmt)
	call, ok := stmt.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", stmt.Value)
	}
	field, ok := call.Callee.(*ast.FieldExpr)
	if !ok {
		t.Fatalf("expected FieldExpr as callee, got %T", call.Callee)
	}
	if field.Field != "length" {
		t.Errorf("expected field 'length', got %q", field.Field)
	}
}

// ---------- Error Cases ----------

func TestMissingSemicolon(t *testing.T) {
	parseExpectError(t, `let x = 5`)
}

func TestMissingCloseBrace(t *testing.T) {
	parseExpectError(t, `fn main() { Println("hello");`)
}

func TestMissingCloseParen(t *testing.T) {
	parseExpectError(t, `fn main( { }`)
}
