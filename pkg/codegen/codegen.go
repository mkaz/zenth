package codegen

import (
	"fmt"
	"strings"

	"github.com/mkaz/zenth/pkg/ast"
	"github.com/mkaz/zenth/pkg/token"
)

// Generator translates a Zenth AST to Go source code.
type Generator struct {
	buf         strings.Builder
	indent      int
	imports     map[string]string // Go import path -> alias (or empty)
	objs        map[string]*ast.ObjDecl
	funcs       map[string]*ast.FnDecl // "name" or "StructName.methodName"
	needsRange  bool
	needsRangei bool
}

// New creates a new code Generator.
func New() *Generator {
	return &Generator{
		imports: make(map[string]string),
		objs: make(map[string]*ast.ObjDecl),
		funcs:   make(map[string]*ast.FnDecl),
	}
}

// Generate produces Go source code from a Zenth AST.
func (g *Generator) Generate(prog *ast.Program) string {
	// First pass: collect objects, functions, and imports
	for _, stmt := range prog.Stmts {
		switch s := stmt.(type) {
		case *ast.ObjDecl:
			g.objs[s.Name] = s
			for _, m := range s.Methods {
				g.funcs[s.Name+"."+m.Name] = m
			}
		case *ast.FnDecl:
			g.funcs[s.Name] = s
		case *ast.ImportDecl:
			g.addImport(s)
		}
	}

	// Generate all top-level declarations into a buffer
	var body strings.Builder
	for _, stmt := range prog.Stmts {
		switch stmt.(type) {
		case *ast.ImportDecl:
			continue // handled separately
		default:
			old := g.buf
			g.buf = strings.Builder{}
			g.genNode(stmt)
			body.WriteString(g.buf.String())
			g.buf = old
		}
	}

	// Now build the final output
	g.buf.Reset()
	g.writeln("package main")
	g.writeln("")

	// Always import fmt for print/println
	g.imports["fmt"] = ""

	if len(g.imports) > 0 {
		g.writeln("import (")
		g.indent++
		for path, alias := range g.imports {
			if alias != "" {
				g.writef("%s %q\n", alias, path)
			} else {
				g.writef("%q\n", path)
			}
		}
		g.indent--
		g.writeln(")")
		g.writeln("")
	}

	if g.needsRange {
		g.writeln("func zenth_range(start, end, step int) []int {")
		g.writeln("\tvar r []int")
		g.writeln("\tif step > 0 {")
		g.writeln("\t\tfor i := start; i < end; i += step { r = append(r, i) }")
		g.writeln("\t} else if step < 0 {")
		g.writeln("\t\tfor i := start; i > end; i += step { r = append(r, i) }")
		g.writeln("\t}")
		g.writeln("\treturn r")
		g.writeln("}")
		g.writeln("")
	}
	if g.needsRangei {
		g.writeln("func zenth_rangei(start, end, step int) []int {")
		g.writeln("\tvar r []int")
		g.writeln("\tif step > 0 {")
		g.writeln("\t\tfor i := start; i <= end; i += step { r = append(r, i) }")
		g.writeln("\t} else if step < 0 {")
		g.writeln("\t\tfor i := start; i >= end; i += step { r = append(r, i) }")
		g.writeln("\t}")
		g.writeln("\treturn r")
		g.writeln("}")
		g.writeln("")
	}

	g.buf.WriteString(body.String())

	return g.buf.String()
}

func (g *Generator) addImport(imp *ast.ImportDecl) {
	goPath := mapImportPath(imp.Path)
	g.imports[goPath] = imp.Alias
}

func mapImportPath(zenthPath string) string {
	switch zenthPath {
	case "fmt":
		return "fmt"
	case "math":
		return "math"
	case "os":
		return "os"
	case "strings", "str":
		return "strings"
	case "io":
		return "os" // map to os for file I/O
	default:
		return zenthPath
	}
}

func (g *Generator) write(s string) {
	g.buf.WriteString(s)
}

func (g *Generator) writeln(s string) {
	g.writeIndent()
	g.buf.WriteString(s)
	g.buf.WriteString("\n")
}

func (g *Generator) writef(format string, args ...interface{}) {
	g.writeIndent()
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) writeIndent() {
	for i := 0; i < g.indent; i++ {
		g.buf.WriteString("\t")
	}
}

func (g *Generator) genNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.FnDecl:
		g.genFnDecl(n)
	case *ast.ObjDecl:
		g.genObjDecl(n)
	case *ast.InterfaceDecl:
		g.genInterfaceDecl(n)
	case *ast.Block:
		g.genBlock(n)
	case *ast.LetStmt:
		g.genLetStmt(n)
	case *ast.VarStmt:
		g.genVarStmt(n)
	case *ast.ConstStmt:
		g.genConstStmt(n)
	case *ast.AssignStmt:
		g.genAssignStmt(n)
	case *ast.MultiAssignStmt:
		g.genMultiAssignStmt(n)
	case *ast.ReturnStmt:
		g.genReturnStmt(n)
	case *ast.IfStmt:
		g.genIfStmt(n)
	case *ast.ForStmt:
		g.genForStmt(n)
	case *ast.ForInStmt:
		g.genForInStmt(n)
	case *ast.MatchStmt:
		g.genMatchStmt(n)
	case *ast.BreakStmt:
		g.writeln("break")
	case *ast.ContinueStmt:
		g.writeln("continue")
	case *ast.IncDecStmt:
		g.genIncDecStmt(n)
	case *ast.ExprStmt:
		g.writeIndent()
		g.genExpr(n.Expr)
		g.write("\n")
	default:
		g.writef("// unhandled: %T\n", node)
	}
}

func (g *Generator) genFnDecl(f *ast.FnDecl) {
	g.writeIndent()
	g.write("func ")
	if f.OwnerObj != "" {
		g.write("(self *")
		g.write(f.OwnerObj)
		g.write(") ")
		g.write(exportName(f.Name))
	} else {
		g.write(f.Name)
	}
	g.write("(")
	for i, p := range f.Params {
		if i > 0 {
			g.write(", ")
		}
		g.write(p.Name)
		g.write(" ")
		g.write(genTypeExpr(p.Type))
	}
	g.write(")")
	if f.ReturnType != nil {
		g.write(" ")
		g.write(genTypeExpr(f.ReturnType))
	}
	g.write(" {\n")
	g.indent++
	if f.Body != nil {
		for _, stmt := range f.Body.Stmts {
			g.genNode(stmt)
		}
	}
	g.indent--
	g.writeln("}")
	g.writeln("")
}

func (g *Generator) genObjDecl(s *ast.ObjDecl) {
	g.writef("type %s struct {\n", s.Name)
	g.indent++
	for _, f := range s.Fields {
		g.writef("%s %s\n", exportName(f.Name), genTypeExpr(f.Type))
	}
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Generate methods
	for _, m := range s.Methods {
		g.genFnDecl(m)
	}
}

func (g *Generator) genInterfaceDecl(iface *ast.InterfaceDecl) {
	g.writef("type %s interface {\n", iface.Name)
	g.indent++
	for _, m := range iface.Methods {
		g.writeIndent()
		g.write(exportName(m.Name))
		g.write("(")
		for i, p := range m.Params {
			if i > 0 {
				g.write(", ")
			}
			g.write(p.Name)
			g.write(" ")
			g.write(genTypeExpr(p.Type))
		}
		g.write(")")
		if m.ReturnType != nil {
			g.write(" ")
			g.write(genTypeExpr(m.ReturnType))
		}
		g.write("\n")
	}
	g.indent--
	g.writeln("}")
	g.writeln("")
}

func (g *Generator) genBlock(b *ast.Block) {
	for _, stmt := range b.Stmts {
		g.genNode(stmt)
	}
}

func (g *Generator) genLetStmt(s *ast.LetStmt) {
	g.writeIndent()
	if s.Infer || s.Type == nil {
		g.write(s.Name + " := ")
	} else {
		g.write("var " + s.Name + " " + genTypeExpr(s.Type) + " = ")
	}
	g.genExpr(s.Value)
	g.write("\n")
	// Suppress Go's "declared and not used" by referencing with _
	g.writef("_ = %s\n", s.Name)
}

func (g *Generator) genVarStmt(s *ast.VarStmt) {
	g.writeIndent()
	if s.Infer || s.Type == nil {
		g.write(s.Name + " := ")
	} else {
		g.write("var " + s.Name + " " + genTypeExpr(s.Type) + " = ")
	}
	g.genExpr(s.Value)
	g.write("\n")
}

func (g *Generator) genConstStmt(s *ast.ConstStmt) {
	g.writeIndent()
	// Go const requires compile-time constant expressions
	// Use var for safety since Zenth const may reference runtime values
	g.write(s.Name + " := ")
	g.genExpr(s.Value)
	g.write("\n")
	g.writef("_ = %s\n", s.Name)
}

func (g *Generator) genAssignStmt(s *ast.AssignStmt) {
	g.writeIndent()
	g.genExpr(s.Target)
	switch s.Op {
	case token.Assign:
		g.write(" = ")
	case token.PlusAssign:
		g.write(" += ")
	case token.MinusAssign:
		g.write(" -= ")
	case token.StarAssign:
		g.write(" *= ")
	case token.SlashAssign:
		g.write(" /= ")
	}
	g.genExpr(s.Value)
	g.write("\n")
}

func (g *Generator) genMultiAssignStmt(s *ast.MultiAssignStmt) {
	g.writeIndent()
	for i, target := range s.Targets {
		if i > 0 {
			g.write(", ")
		}
		g.genExpr(target)
	}
	g.write(" = ")
	for i, value := range s.Values {
		if i > 0 {
			g.write(", ")
		}
		g.genExpr(value)
	}
	g.write("\n")
}

func (g *Generator) genReturnStmt(s *ast.ReturnStmt) {
	g.writeIndent()
	g.write("return")
	if s.Value != nil {
		g.write(" ")
		g.genExpr(s.Value)
	}
	g.write("\n")
}

func (g *Generator) genIfStmt(s *ast.IfStmt) {
	g.writeIndent()
	g.write("if ")
	g.genExpr(s.Condition)
	g.write(" {\n")
	g.indent++
	for _, stmt := range s.Body.Stmts {
		g.genNode(stmt)
	}
	g.indent--
	if s.Else != nil {
		switch e := s.Else.(type) {
		case *ast.IfStmt:
			g.writeIndent()
			g.write("} else ")
			// Remove indent since genIfStmt will add its own
			old := g.indent
			g.indent = 0
			g.write("if ")
			g.genExpr(e.Condition)
			g.write(" {\n")
			g.indent = old + 1
			for _, stmt := range e.Body.Stmts {
				g.genNode(stmt)
			}
			g.indent = old
			if e.Else != nil {
				g.genElseChain(e.Else)
			} else {
				g.writeln("}")
			}
		case *ast.Block:
			g.writeln("} else {")
			g.indent++
			for _, stmt := range e.Stmts {
				g.genNode(stmt)
			}
			g.indent--
			g.writeln("}")
		}
	} else {
		g.writeln("}")
	}
}

func (g *Generator) genElseChain(node ast.Node) {
	switch e := node.(type) {
	case *ast.IfStmt:
		g.writeIndent()
		g.write("} else if ")
		g.genExpr(e.Condition)
		g.write(" {\n")
		g.indent++
		for _, stmt := range e.Body.Stmts {
			g.genNode(stmt)
		}
		g.indent--
		if e.Else != nil {
			g.genElseChain(e.Else)
		} else {
			g.writeln("}")
		}
	case *ast.Block:
		g.writeln("} else {")
		g.indent++
		for _, stmt := range e.Stmts {
			g.genNode(stmt)
		}
		g.indent--
		g.writeln("}")
	}
}

func (g *Generator) genForStmt(s *ast.ForStmt) {
	g.writeIndent()
	if s.Init == nil && s.Condition == nil && s.Post == nil {
		// Infinite loop
		g.write("for {\n")
	} else if s.Init == nil && s.Post == nil {
		// While-style
		g.write("for ")
		g.genExpr(s.Condition)
		g.write(" {\n")
	} else {
		// C-style
		g.write("for ")
		if s.Init != nil {
			g.genForClause(s.Init)
		}
		g.write("; ")
		if s.Condition != nil {
			g.genExpr(s.Condition)
		}
		g.write("; ")
		if s.Post != nil {
			g.genForClause(s.Post)
		}
		g.write(" {\n")
	}
	g.indent++
	for _, stmt := range s.Body.Stmts {
		g.genNode(stmt)
	}
	g.indent--
	g.writeln("}")
}

func (g *Generator) genForClause(node ast.Node) {
	switch n := node.(type) {
	case *ast.VarStmt:
		g.write(n.Name + " := ")
		g.genExpr(n.Value)
	case *ast.LetStmt:
		g.write(n.Name + " := ")
		g.genExpr(n.Value)
	case *ast.AssignStmt:
		g.genExpr(n.Target)
		switch n.Op {
		case token.Assign:
			g.write(" = ")
		case token.PlusAssign:
			g.write(" += ")
		}
		g.genExpr(n.Value)
	case *ast.MultiAssignStmt:
		for i, target := range n.Targets {
			if i > 0 {
				g.write(", ")
			}
			g.genExpr(target)
		}
		g.write(" = ")
		for i, value := range n.Values {
			if i > 0 {
				g.write(", ")
			}
			g.genExpr(value)
		}
	case *ast.IncDecStmt:
		g.genExpr(n.Operand)
		if n.Op == token.PlusPlus {
			g.write("++")
		} else {
			g.write("--")
		}
	case *ast.ExprStmt:
		g.genExpr(n.Expr)
	}
}

func (g *Generator) genForInStmt(s *ast.ForInStmt) {
	g.writeIndent()
	g.write("for ")
	if s.Index != "" {
		g.write(s.Index)
	} else {
		g.write("_")
	}
	g.write(", ")
	g.write(s.Value)
	g.write(" := range ")
	g.genExpr(s.Iterable)
	g.write(" {\n")
	g.indent++
	for _, stmt := range s.Body.Stmts {
		g.genNode(stmt)
	}
	g.indent--
	g.writeln("}")
}

func (g *Generator) genMatchStmt(s *ast.MatchStmt) {
	g.writeIndent()
	g.write("switch ")
	g.genExpr(s.Subject)
	g.write(" {\n")
	for _, arm := range s.Arms {
		g.writeIndent()
		// Check for wildcard _
		if ident, ok := arm.Pattern.(*ast.IdentExpr); ok && ident.Name == "_" {
			g.write("default:\n")
		} else {
			g.write("case ")
			g.genExpr(arm.Pattern)
			g.write(":\n")
		}
		g.indent++
		switch body := arm.Body.(type) {
		case *ast.Block:
			for _, stmt := range body.Stmts {
				g.genNode(stmt)
			}
		default:
			g.genNode(body)
		}
		g.indent--
	}
	g.writeln("}")
}

func (g *Generator) genIncDecStmt(s *ast.IncDecStmt) {
	g.writeIndent()
	g.genExpr(s.Operand)
	if s.Op == token.PlusPlus {
		g.write("++")
	} else {
		g.write("--")
	}
	g.write("\n")
}

// ---------- Expression generation ----------

func (g *Generator) genExpr(node ast.Node) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if n.PromoteLeft != "" {
			g.write(n.PromoteLeft + "(")
			g.genExpr(n.Left)
			g.write(")")
		} else {
			g.genExpr(n.Left)
		}
		g.write(" " + goOp(n.Op) + " ")
		if n.PromoteRight != "" {
			g.write(n.PromoteRight + "(")
			g.genExpr(n.Right)
			g.write(")")
		} else {
			g.genExpr(n.Right)
		}
	case *ast.UnaryExpr:
		g.write(goOp(n.Op))
		g.genExpr(n.Operand)
	case *ast.CallExpr:
		g.genCallExpr(n)
	case *ast.IndexExpr:
		g.genExpr(n.Object)
		g.write("[")
		g.genExpr(n.Index)
		g.write("]")
	case *ast.FieldExpr:
		g.genExpr(n.Object)
		g.write(".")
		g.write(goFieldName(n.Field, n.Object))
	case *ast.IdentExpr:
		g.write(n.Name)
	case *ast.IntLitExpr:
		g.write(fmt.Sprintf("%d", n.Value))
	case *ast.FloatLitExpr:
		s := fmt.Sprintf("%g", n.Value)
		if !strings.Contains(s, ".") && !strings.Contains(s, "e") {
			s += ".0"
		}
		g.write(s)
	case *ast.StringLitExpr:
		g.write(fmt.Sprintf("%q", n.Value))
	case *ast.BoolLitExpr:
		if n.Value {
			g.write("true")
		} else {
			g.write("false")
		}
	case *ast.NilExpr:
		g.write("nil")
	case *ast.ArrayLitExpr:
		g.genArrayLit(n)
	case *ast.NamedArgExpr:
		g.genExpr(n.Value)
	case *ast.IfExpr:
		g.genIfExpr(n)
	case *ast.InterpStringExpr:
		g.genInterpString(n)
	default:
		g.write(fmt.Sprintf("/* unhandled expr: %T */", node))
	}
}

func (g *Generator) genCallExpr(c *ast.CallExpr) {
	// Translate built-in functions
	if ident, ok := c.Callee.(*ast.IdentExpr); ok {
		switch ident.Name {
		case "print":
			g.write("fmt.Print(")
			g.genArgList(c.Args)
			g.write(")")
			return
		case "println":
			g.write("fmt.Println(")
			g.genArgList(c.Args)
			g.write(")")
			return
		case "len":
			g.write("len(")
			g.genArgList(c.Args)
			g.write(")")
			return
		case "str":
			g.write("fmt.Sprint(")
			g.genArgList(c.Args)
			g.write(")")
			return
		case "range":
			g.needsRange = true
			g.write("zenth_range(")
			g.genArgList(c.Args)
			if len(c.Args) == 2 {
				g.write(", 1")
			}
			g.write(")")
			return
		case "rangei":
			g.needsRangei = true
			g.write("zenth_rangei(")
			g.genArgList(c.Args)
			if len(c.Args) == 2 {
				g.write(", 1")
			}
			g.write(")")
			return
		}
	}

	// Handle obj constructor calls
	if ident, ok := c.Callee.(*ast.IdentExpr); ok {
		if decl, ok := g.objs[ident.Name]; ok {
			g.genObjConstructor(decl, c.Args)
			return
		}
	}

	// Translate module function calls (fmt.println -> fmt.Println)
	if field, ok := c.Callee.(*ast.FieldExpr); ok {
		if ident, ok := field.Object.(*ast.IdentExpr); ok {
			mapped := mapModuleCall(ident.Name, field.Field)
			if mapped != "" {
				g.write(mapped)
				g.write("(")
				g.genArgList(c.Args)
				g.write(")")
				return
			}
		}
	}

	// Fill in default args for user-defined functions
	args := c.Args
	if ident, ok := c.Callee.(*ast.IdentExpr); ok {
		if decl, ok := g.funcs[ident.Name]; ok && len(args) < len(decl.Params) {
			args = g.fillDefaults(args, decl)
		}
	} else if field, ok := c.Callee.(*ast.FieldExpr); ok {
		// Check for struct method calls
		if ident, ok := field.Object.(*ast.IdentExpr); ok {
			key := ident.Name + "." + field.Field
			if decl, ok := g.funcs[key]; ok && len(args) < len(decl.Params) {
				args = g.fillDefaults(args, decl)
			}
		}
	}

	g.genExpr(c.Callee)
	g.write("(")
	g.genArgList(args)
	g.write(")")
}

func (g *Generator) fillDefaults(args []ast.Node, decl *ast.FnDecl) []ast.Node {
	filled := make([]ast.Node, len(args))
	copy(filled, args)
	for i := len(args); i < len(decl.Params); i++ {
		if decl.Params[i].Default != nil {
			filled = append(filled, decl.Params[i].Default)
		}
	}
	return filled
}

func (g *Generator) genArgList(args []ast.Node) {
	for i, arg := range args {
		if i > 0 {
			g.write(", ")
		}
		g.genExpr(arg)
	}
}

func (g *Generator) genArrayLit(a *ast.ArrayLitExpr) {
	// Try to infer element type from first element
	g.write("[]interface{}{")
	for i, elem := range a.Elements {
		if i > 0 {
			g.write(", ")
		}
		g.genExpr(elem)
	}
	g.write("}")
}

func (g *Generator) genObjConstructor(decl *ast.ObjDecl, args []ast.Node) {
	// Build map of provided named args
	provided := make(map[string]ast.Node)
	for _, arg := range args {
		if named, ok := arg.(*ast.NamedArgExpr); ok {
			provided[named.Name] = named.Value
		}
	}

	g.write("&" + decl.Name + "{")
	first := true
	for _, f := range decl.Fields {
		val, ok := provided[f.Name]
		if !ok {
			val = f.Default
		}
		if val == nil {
			continue
		}
		if !first {
			g.write(", ")
		}
		first = false
		g.write(exportName(f.Name) + ": ")
		g.genExpr(val)
	}
	g.write("}")
}

func (g *Generator) genInterpString(s *ast.InterpStringExpr) {
	// Build fmt.Sprintf("...%v...", arg1, arg2, ...)
	var format strings.Builder
	var args []ast.Node
	for _, part := range s.Parts {
		if part.IsExpr {
			format.WriteString("%v")
			args = append(args, part.Expr)
		} else {
			// Escape % as %% and quote special chars for Go string literal
			for _, ch := range part.Lit {
				if ch == '%' {
					format.WriteString("%%")
				} else {
					format.WriteRune(ch)
				}
			}
		}
	}
	g.write("fmt.Sprintf(")
	g.write(fmt.Sprintf("%q", format.String()))
	for _, arg := range args {
		g.write(", ")
		g.genExpr(arg)
	}
	g.write(")")
}

func (g *Generator) genIfExpr(e *ast.IfExpr) {
	goType := e.GoType
	if goType == "" {
		goType = "interface{}"
	}
	g.write("func() " + goType + " { ")
	g.genIfExprBody(e)
	g.write(" }()")
}

func (g *Generator) genIfExprBody(e *ast.IfExpr) {
	g.write("if ")
	g.genExpr(e.Condition)
	g.write(" { return ")
	g.genExpr(e.Then)
	g.write(" }")
	// else branch
	if inner, ok := e.Else.(*ast.IfExpr); ok {
		g.write(" else ")
		g.genIfExprBody(inner)
	} else {
		g.write(" else { return ")
		g.genExpr(e.Else)
		g.write(" }")
	}
}

// ---------- Helpers ----------

func goOp(op token.Type) string {
	switch op {
	case token.Plus:
		return "+"
	case token.Minus:
		return "-"
	case token.Star:
		return "*"
	case token.Slash:
		return "/"
	case token.Percent:
		return "%"
	case token.Eq:
		return "=="
	case token.Neq:
		return "!="
	case token.Lt:
		return "<"
	case token.Gt:
		return ">"
	case token.Lte:
		return "<="
	case token.Gte:
		return ">="
	case token.And:
		return "&&"
	case token.Or:
		return "||"
	case token.Not:
		return "!"
	default:
		return "?"
	}
}

func genTypeExpr(t *ast.TypeExpr) string {
	if t == nil {
		return ""
	}
	if t.IsSlice && len(t.Params) > 0 {
		return "[]" + genTypeExpr(t.Params[0])
	}
	return mapTypeName(t.Name)
}

func mapTypeName(name string) string {
	switch name {
	case "int":
		return "int"
	case "i8":
		return "int8"
	case "i16":
		return "int16"
	case "i32":
		return "int32"
	case "i64":
		return "int64"
	case "u8":
		return "uint8"
	case "u16":
		return "uint16"
	case "u32":
		return "uint32"
	case "u64":
		return "uint64"
	case "f32":
		return "float32"
	case "f64":
		return "float64"
	case "bool":
		return "bool"
	case "str":
		return "string"
	case "byte":
		return "byte"
	default:
		// User-defined obj types are always pointers
		return "*" + name
	}
}

func exportName(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func goFieldName(field string, _ ast.Node) string {
	// All field/method access on objects gets capitalized for Go export
	return exportName(field)
}

func mapModuleCall(module, fn string) string {
	switch module {
	case "fmt":
		switch fn {
		case "println":
			return "fmt.Println"
		case "print":
			return "fmt.Print"
		case "printf":
			return "fmt.Printf"
		case "sprintf", "format":
			return "fmt.Sprintf"
		}
	case "math":
		switch fn {
		case "sqrt":
			return "math.Sqrt"
		case "abs":
			return "math.Abs"
		case "pow":
			return "math.Pow"
		case "min":
			return "math.Min"
		case "max":
			return "math.Max"
		case "pi":
			return "math.Pi"
		}
	case "os":
		switch fn {
		case "exit":
			return "os.Exit"
		case "args":
			return "os.Args"
		}
	case "strings", "str":
		switch fn {
		case "split":
			return "strings.Split"
		case "join":
			return "strings.Join"
		case "contains":
			return "strings.Contains"
		case "replace":
			return "strings.ReplaceAll"
		case "trim":
			return "strings.TrimSpace"
		}
	}
	return ""
}
