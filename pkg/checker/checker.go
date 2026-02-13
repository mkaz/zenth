package checker

import (
	"fmt"

	"github.com/mkaz/zenth/pkg/ast"
	"github.com/mkaz/zenth/pkg/token"
)

// FuncInfo stores function metadata for type checking.
type FuncInfo struct {
	Name        string
	Params      []ZType
	ParamNames  []string
	Return      ZType
	Receiver    string // empty for free functions
	NumRequired int    // number of params without defaults
}

// ObjInfo stores obj metadata.
type ObjInfo struct {
	Name   string
	Fields map[string]ZType
	Order  []string // field order for codegen
}

// Checker performs type checking and semantic analysis on a Zenth AST.
type Checker struct {
	scope   *Scope
	funcs   map[string]*FuncInfo   // "name" or "Type.name"
	objs    map[string]*ObjInfo
	errors  []string
	currentFunc *FuncInfo // for checking return types
}

// New creates a new Checker.
func New() *Checker {
	global := NewScope(nil)

	c := &Checker{
		scope:   global,
		funcs:   make(map[string]*FuncInfo),
		objs:    make(map[string]*ObjInfo),
	}

	// Register built-in functions
	c.funcs["print"] = &FuncInfo{
		Name:   "print",
		Params: []ZType{TypeStr},
		Return: TypeVoid,
	}
	c.funcs["println"] = &FuncInfo{
		Name:   "println",
		Params: []ZType{TypeStr},
		Return: TypeVoid,
	}
	c.funcs["len"] = &FuncInfo{
		Name:   "len",
		Params: []ZType{TypeStr}, // overloaded for slices too
		Return: TypeInt,
	}
	c.funcs["str"] = &FuncInfo{
		Name:   "str",
		Params: []ZType{TypeInt},
		Return: TypeStr,
	}
	c.funcs["range"] = &FuncInfo{
		Name:   "range",
		Params: []ZType{TypeInt, TypeInt},
		Return: &SliceType{Elem: TypeInt},
	}
	c.funcs["rangei"] = &FuncInfo{
		Name:   "rangei",
		Params: []ZType{TypeInt, TypeInt},
		Return: &SliceType{Elem: TypeInt},
	}

	return c
}

func (c *Checker) errorf(pos token.Pos, format string, args ...interface{}) {
	msg := fmt.Sprintf("%s: %s", pos, fmt.Sprintf(format, args...))
	c.errors = append(c.errors, msg)
}

func (c *Checker) pushScope() {
	c.scope = NewScope(c.scope)
}

func (c *Checker) popScope() {
	c.scope = c.scope.parent
}

// Check performs type checking on a Program.
func (c *Checker) Check(prog *ast.Program) error {
	// First pass: register all top-level declarations
	for _, stmt := range prog.Stmts {
		switch s := stmt.(type) {
		case *ast.ObjDecl:
			c.registerObj(s)
		case *ast.FnDecl:
			c.registerFunc(s)
		}
	}

	// Second pass: check all declarations
	for _, stmt := range prog.Stmts {
		c.checkNode(stmt)
	}

	if len(c.errors) > 0 {
		result := "type errors:\n"
		for _, e := range c.errors {
			result += "  " + e + "\n"
		}
		return fmt.Errorf("%s", result)
	}
	return nil
}

func (c *Checker) registerObj(s *ast.ObjDecl) {
	info := &ObjInfo{
		Name:   s.Name,
		Fields: make(map[string]ZType),
	}
	for _, f := range s.Fields {
		t := c.resolveTypeExpr(f.Type)
		info.Fields[f.Name] = t
		info.Order = append(info.Order, f.Name)
	}
	c.objs[s.Name] = info

	// Register methods defined inside the struct
	for _, m := range s.Methods {
		c.registerFunc(m)
	}
}

func (c *Checker) registerFunc(f *ast.FnDecl) {
	info := &FuncInfo{Name: f.Name}
	for _, p := range f.Params {
		pt := c.resolveTypeExpr(p.Type)
		info.Params = append(info.Params, pt)
		info.ParamNames = append(info.ParamNames, p.Name)
		if p.Default == nil {
			info.NumRequired++
		} else {
			// Type-check default value against declared param type
			defType := c.checkNode(p.Default)
			if !pt.Equals(defType) && defType != TypeNil {
				c.errorf(f.Pos(), "default value for '%s' has type %s, expected %s", p.Name, defType, pt)
			}
		}
	}
	if f.ReturnType != nil {
		info.Return = c.resolveTypeExpr(f.ReturnType)
	} else {
		info.Return = TypeVoid
	}
	if f.OwnerObj != "" {
		info.Receiver = f.OwnerObj
		key := f.OwnerObj + "." + f.Name
		c.funcs[key] = info
	} else {
		c.funcs[f.Name] = info
	}
}

func (c *Checker) resolveTypeExpr(t *ast.TypeExpr) ZType {
	if t == nil {
		return TypeVoid
	}
	if t.IsSlice && len(t.Params) > 0 {
		return &SliceType{Elem: c.resolveTypeExpr(t.Params[0])}
	}
	if bt := LookupBuiltinType(t.Name); bt != nil {
		return bt
	}
	if st, ok := c.objs[t.Name]; ok {
		return &ObjType{Name: st.Name, Fields: st.Fields}
	}
	c.errorf(t.Pos(), "unknown type: %s", t.Name)
	return TypeVoid
}

func (c *Checker) checkNode(node ast.Node) ZType {
	switch n := node.(type) {
	case *ast.FnDecl:
		return c.checkFnDecl(n)
	case *ast.ObjDecl:
		// Check methods defined inside the obj
		for _, m := range n.Methods {
			c.checkFnDecl(m)
		}
		return TypeVoid
	case *ast.InterfaceDecl:
		return TypeVoid // TODO: check interface
	case *ast.ImportDecl:
		return TypeVoid
	case *ast.Block:
		return c.checkBlock(n)
	case *ast.LetStmt:
		return c.checkLetStmt(n)
	case *ast.VarStmt:
		return c.checkVarStmt(n)
	case *ast.ConstStmt:
		return c.checkConstStmt(n)
	case *ast.AssignStmt:
		return c.checkAssignStmt(n)
	case *ast.MultiAssignStmt:
		return c.checkMultiAssignStmt(n)
	case *ast.ReturnStmt:
		return c.checkReturnStmt(n)
	case *ast.IfStmt:
		return c.checkIfStmt(n)
	case *ast.ForStmt:
		return c.checkForStmt(n)
	case *ast.ForInStmt:
		return c.checkForInStmt(n)
	case *ast.MatchStmt:
		return c.checkMatchStmt(n)
	case *ast.BreakStmt, *ast.ContinueStmt:
		return TypeVoid
	case *ast.IncDecStmt:
		t := c.checkNode(n.Operand)
		if !IsNumeric(t) {
			c.errorf(n.Pos(), "++/-- requires numeric type, got %s", t)
		}
		return TypeVoid
	case *ast.ExprStmt:
		return c.checkNode(n.Expr)
	case *ast.BinaryExpr:
		return c.checkBinaryExpr(n)
	case *ast.UnaryExpr:
		return c.checkUnaryExpr(n)
	case *ast.CallExpr:
		return c.checkCallExpr(n)
	case *ast.FieldExpr:
		return c.checkFieldExpr(n)
	case *ast.IndexExpr:
		return c.checkIndexExpr(n)
	case *ast.IdentExpr:
		return c.checkIdentExpr(n)
	case *ast.IntLitExpr:
		return TypeInt
	case *ast.FloatLitExpr:
		return TypeF64
	case *ast.StringLitExpr:
		return TypeStr
	case *ast.BoolLitExpr:
		return TypeBool
	case *ast.NilExpr:
		return TypeNil
	case *ast.ArrayLitExpr:
		return c.checkArrayLit(n)
	case *ast.ObjLitExpr:
		return c.checkObjLit(n)
	case *ast.IfExpr:
		return c.checkIfExpr(n)
	case *ast.InterpStringExpr:
		for _, part := range n.Parts {
			if part.IsExpr {
				c.checkNode(part.Expr)
			}
		}
		return TypeStr
	default:
		c.errorf(node.Pos(), "unhandled node type in checker: %T", node)
		return TypeVoid
	}
}

func (c *Checker) checkFnDecl(f *ast.FnDecl) ZType {
	prev := c.currentFunc
	info := c.funcs[f.Name]
	if f.OwnerObj != "" {
		info = c.funcs[f.OwnerObj+"."+f.Name]
	}
	c.currentFunc = info

	c.pushScope()

	// Bind self for methods
	if f.OwnerObj != "" {
		if st, ok := c.objs[f.OwnerObj]; ok {
			selfType := &ObjType{Name: st.Name, Fields: st.Fields}
			c.scope.Define(&Symbol{Name: "self", Type: selfType})
		}
	}

	// Bind parameters
	for i, p := range f.Params {
		c.scope.Define(&Symbol{Name: p.Name, Type: info.Params[i]})
	}

	c.checkNode(f.Body)
	c.popScope()
	c.currentFunc = prev
	return TypeVoid
}

func (c *Checker) checkBlock(b *ast.Block) ZType {
	c.pushScope()
	for _, stmt := range b.Stmts {
		c.checkNode(stmt)
	}
	c.popScope()
	return TypeVoid
}

func (c *Checker) checkLetStmt(s *ast.LetStmt) ZType {
	valType := c.checkNode(s.Value)
	if s.Type != nil {
		declared := c.resolveTypeExpr(s.Type)
		if !declared.Equals(valType) && valType != TypeNil {
			c.errorf(s.Pos(), "type mismatch: cannot assign %s to %s", valType, declared)
		}
		valType = declared
	}
	c.scope.Define(&Symbol{Name: s.Name, Type: valType, Mutable: false})
	return TypeVoid
}

func (c *Checker) checkVarStmt(s *ast.VarStmt) ZType {
	valType := c.checkNode(s.Value)
	if s.Type != nil {
		declared := c.resolveTypeExpr(s.Type)
		if !declared.Equals(valType) && valType != TypeNil {
			c.errorf(s.Pos(), "type mismatch: cannot assign %s to %s", valType, declared)
		}
		valType = declared
	}
	c.scope.Define(&Symbol{Name: s.Name, Type: valType, Mutable: true})
	return TypeVoid
}

func (c *Checker) checkConstStmt(s *ast.ConstStmt) ZType {
	valType := c.checkNode(s.Value)
	if s.Type != nil {
		declared := c.resolveTypeExpr(s.Type)
		if !declared.Equals(valType) && valType != TypeNil {
			c.errorf(s.Pos(), "type mismatch: cannot assign %s to %s", valType, declared)
		}
		valType = declared
	}
	c.scope.Define(&Symbol{Name: s.Name, Type: valType, IsConst: true})
	return TypeVoid
}

func (c *Checker) checkAssignStmt(s *ast.AssignStmt) ZType {
	targetType := c.checkNode(s.Target)
	valueType := c.checkNode(s.Value)

	// Check mutability
	if ident, ok := s.Target.(*ast.IdentExpr); ok {
		sym := c.scope.Lookup(ident.Name)
		if sym != nil && !sym.Mutable {
			c.errorf(s.Pos(), "cannot assign to immutable variable '%s' (use 'var' instead of 'let')", ident.Name)
		}
		if sym != nil && sym.IsConst {
			c.errorf(s.Pos(), "cannot assign to constant '%s'", ident.Name)
		}
	}

	if s.Op == token.Assign {
		if !targetType.Equals(valueType) && valueType != TypeNil {
			c.errorf(s.Pos(), "type mismatch: cannot assign %s to %s", valueType, targetType)
		}
	} else {
		// Compound assignment (+=, -=, etc.)
		if !IsNumeric(targetType) {
			// Allow string concatenation with +=
			if s.Op == token.PlusAssign && targetType.Equals(TypeStr) && valueType.Equals(TypeStr) {
				return TypeVoid
			}
			c.errorf(s.Pos(), "compound assignment requires numeric type, got %s", targetType)
		}
	}
	return TypeVoid
}

func (c *Checker) checkMultiAssignStmt(s *ast.MultiAssignStmt) ZType {
	if len(s.Targets) != len(s.Values) {
		c.errorf(s.Pos(), "multi-assign: %d targets but %d values", len(s.Targets), len(s.Values))
		return TypeVoid
	}
	for i, target := range s.Targets {
		targetType := c.checkNode(target)
		valueType := c.checkNode(s.Values[i])
		// Check mutability
		if ident, ok := target.(*ast.IdentExpr); ok {
			sym := c.scope.Lookup(ident.Name)
			if sym != nil && !sym.Mutable {
				c.errorf(s.Pos(), "cannot assign to immutable variable '%s' (use 'var' instead of 'let')", ident.Name)
			}
			if sym != nil && sym.IsConst {
				c.errorf(s.Pos(), "cannot assign to constant '%s'", ident.Name)
			}
		}
		if !targetType.Equals(valueType) && valueType != TypeNil {
			c.errorf(s.Pos(), "type mismatch: cannot assign %s to %s", valueType, targetType)
		}
	}
	return TypeVoid
}

func (c *Checker) checkReturnStmt(s *ast.ReturnStmt) ZType {
	if c.currentFunc == nil {
		c.errorf(s.Pos(), "return outside of function")
		return TypeVoid
	}
	if s.Value == nil {
		if c.currentFunc.Return != TypeVoid {
			c.errorf(s.Pos(), "function expects return type %s", c.currentFunc.Return)
		}
		return TypeVoid
	}
	valType := c.checkNode(s.Value)
	if !c.currentFunc.Return.Equals(valType) && valType != TypeNil && c.currentFunc.Return != TypeVoid {
		c.errorf(s.Pos(), "return type mismatch: expected %s, got %s", c.currentFunc.Return, valType)
	}
	return TypeVoid
}

func (c *Checker) checkIfStmt(s *ast.IfStmt) ZType {
	condType := c.checkNode(s.Condition)
	if !condType.Equals(TypeBool) {
		c.errorf(s.Condition.Pos(), "if condition must be bool, got %s", condType)
	}
	c.checkNode(s.Body)
	if s.Else != nil {
		c.checkNode(s.Else)
	}
	return TypeVoid
}

func (c *Checker) checkForStmt(s *ast.ForStmt) ZType {
	c.pushScope()
	if s.Init != nil {
		c.checkNode(s.Init)
	}
	if s.Condition != nil {
		condType := c.checkNode(s.Condition)
		if !condType.Equals(TypeBool) {
			c.errorf(s.Condition.Pos(), "for condition must be bool, got %s", condType)
		}
	}
	if s.Post != nil {
		c.checkNode(s.Post)
	}
	// Check body without pushing scope again (block will push its own)
	for _, stmt := range s.Body.Stmts {
		c.checkNode(stmt)
	}
	c.popScope()
	return TypeVoid
}

func (c *Checker) checkForInStmt(s *ast.ForInStmt) ZType {
	iterType := c.checkNode(s.Iterable)

	c.pushScope()
	switch t := iterType.(type) {
	case *SliceType:
		if s.Index != "" {
			c.scope.Define(&Symbol{Name: s.Index, Type: TypeInt})
		}
		c.scope.Define(&Symbol{Name: s.Value, Type: t.Elem})
	default:
		if iterType.Equals(TypeStr) {
			if s.Index != "" {
				c.scope.Define(&Symbol{Name: s.Index, Type: TypeInt})
			}
			c.scope.Define(&Symbol{Name: s.Value, Type: TypeStr})
		} else {
			c.errorf(s.Pos(), "cannot iterate over %s", iterType)
		}
	}

	for _, stmt := range s.Body.Stmts {
		c.checkNode(stmt)
	}
	c.popScope()
	return TypeVoid
}

func (c *Checker) checkMatchStmt(s *ast.MatchStmt) ZType {
	c.checkNode(s.Subject)
	for _, arm := range s.Arms {
		c.checkNode(arm.Pattern)
		c.checkNode(arm.Body)
	}
	return TypeVoid
}

func (c *Checker) checkBinaryExpr(e *ast.BinaryExpr) ZType {
	left := c.checkNode(e.Left)
	right := c.checkNode(e.Right)

	switch e.Op {
	case token.Plus:
		// Allow string concatenation
		if left.Equals(TypeStr) && right.Equals(TypeStr) {
			return TypeStr
		}
		if IsNumeric(left) && left.Equals(right) {
			return left
		}
		c.errorf(e.Pos(), "cannot add %s and %s", left, right)
		return left

	case token.Minus, token.Star, token.Slash, token.Percent:
		if IsNumeric(left) && left.Equals(right) {
			return left
		}
		c.errorf(e.Pos(), "cannot use %s with %s and %s", e.Op, left, right)
		return left

	case token.Eq, token.Neq:
		return TypeBool

	case token.Lt, token.Gt, token.Lte, token.Gte:
		if IsNumeric(left) && left.Equals(right) {
			return TypeBool
		}
		c.errorf(e.Pos(), "cannot compare %s and %s", left, right)
		return TypeBool

	case token.And, token.Or:
		if !left.Equals(TypeBool) || !right.Equals(TypeBool) {
			c.errorf(e.Pos(), "logical operators require bool operands, got %s and %s", left, right)
		}
		return TypeBool
	}

	return TypeVoid
}

func (c *Checker) checkUnaryExpr(e *ast.UnaryExpr) ZType {
	operand := c.checkNode(e.Operand)
	switch e.Op {
	case token.Not:
		if !operand.Equals(TypeBool) {
			c.errorf(e.Pos(), "! requires bool, got %s", operand)
		}
		return TypeBool
	case token.Minus:
		if !IsNumeric(operand) {
			c.errorf(e.Pos(), "- requires numeric type, got %s", operand)
		}
		return operand
	}
	return operand
}

func (c *Checker) checkCallExpr(e *ast.CallExpr) ZType {
	// Handle method calls: obj.method(args)
	if field, ok := e.Callee.(*ast.FieldExpr); ok {
		objType := c.checkNode(field.Object)
		methodKey := ""
		if st, ok := objType.(*ObjType); ok {
			methodKey = st.Name + "." + field.Field
		}
		if info, ok := c.funcs[methodKey]; ok {
			c.checkArgs(e, info)
			return info.Return
		}
		// Check for imported module function call (e.g., fmt.println)
		if ident, ok := field.Object.(*ast.IdentExpr); ok {
			qualName := ident.Name + "." + field.Field
			if info, ok := c.funcs[qualName]; ok {
				c.checkArgs(e, info)
				return info.Return
			}
		}
		// Allow any method call for now (stdlib calls we haven't registered)
		for _, arg := range e.Args {
			c.checkNode(arg)
		}
		return TypeVoid
	}

	// Handle free function calls
	if ident, ok := e.Callee.(*ast.IdentExpr); ok {
		if info, ok := c.funcs[ident.Name]; ok {
			c.checkArgs(e, info)
			return info.Return
		}
		// Built-in print/println accept any args
		if ident.Name == "print" || ident.Name == "println" {
			for _, arg := range e.Args {
				c.checkNode(arg)
			}
			return TypeVoid
		}
		c.errorf(e.Pos(), "undefined function: %s", ident.Name)
		return TypeVoid
	}

	// Type check args even if we can't resolve
	for _, arg := range e.Args {
		c.checkNode(arg)
	}
	return TypeVoid
}

func (c *Checker) checkArgs(e *ast.CallExpr, info *FuncInfo) {
	for _, arg := range e.Args {
		c.checkNode(arg)
	}
	// Flexible arg count for variadic builtins (print, println, etc.)
	if info.Name == "print" || info.Name == "println" {
		return
	}
	// range/rangei accept 2 or 3 int args
	if info.Name == "range" || info.Name == "rangei" {
		if len(e.Args) < 2 || len(e.Args) > 3 {
			c.errorf(e.Pos(), "%s expects 2 or 3 arguments, got %d", info.Name, len(e.Args))
		}
		return
	}
	if len(e.Args) < info.NumRequired || len(e.Args) > len(info.Params) {
		if info.NumRequired == len(info.Params) {
			c.errorf(e.Pos(), "%s expects %d arguments, got %d", info.Name, len(info.Params), len(e.Args))
		} else {
			c.errorf(e.Pos(), "%s expects %d to %d arguments, got %d", info.Name, info.NumRequired, len(info.Params), len(e.Args))
		}
	}
}

func (c *Checker) checkFieldExpr(e *ast.FieldExpr) ZType {
	objType := c.checkNode(e.Object)
	if st, ok := objType.(*ObjType); ok {
		if ft, ok := st.Fields[e.Field]; ok {
			return ft
		}
		c.errorf(e.Pos(), "obj %s has no field '%s'", st.Name, e.Field)
	}
	return TypeVoid
}

func (c *Checker) checkIndexExpr(e *ast.IndexExpr) ZType {
	objType := c.checkNode(e.Object)
	idxType := c.checkNode(e.Index)

	if !IsInteger(idxType) {
		c.errorf(e.Pos(), "index must be integer, got %s", idxType)
	}

	switch t := objType.(type) {
	case *SliceType:
		return t.Elem
	default:
		if objType.Equals(TypeStr) {
			return TypeByte
		}
		c.errorf(e.Pos(), "cannot index %s", objType)
		return TypeVoid
	}
}

func (c *Checker) checkIdentExpr(e *ast.IdentExpr) ZType {
	if e.Name == "_" {
		return TypeVoid
	}
	sym := c.scope.Lookup(e.Name)
	if sym != nil {
		return sym.Type
	}
	// Could be a struct name used as a type constructor
	if _, ok := c.objs[e.Name]; ok {
		return TypeVoid // struct names aren't values
	}
	// Allow unresolved idents for module names (fmt, math, etc.)
	return TypeVoid
}

func (c *Checker) checkArrayLit(e *ast.ArrayLitExpr) ZType {
	if len(e.Elements) == 0 {
		return &SliceType{Elem: TypeVoid}
	}
	firstType := c.checkNode(e.Elements[0])
	for i := 1; i < len(e.Elements); i++ {
		elemType := c.checkNode(e.Elements[i])
		if !firstType.Equals(elemType) {
			c.errorf(e.Elements[i].Pos(), "array element type mismatch: expected %s, got %s", firstType, elemType)
		}
	}
	return &SliceType{Elem: firstType}
}

func (c *Checker) checkObjLit(e *ast.ObjLitExpr) ZType {
	info, ok := c.objs[e.Name]
	if !ok {
		c.errorf(e.Pos(), "undefined obj: %s", e.Name)
		return TypeVoid
	}
	for _, f := range e.Fields {
		expected, ok := info.Fields[f.Name]
		if !ok {
			c.errorf(e.Pos(), "obj %s has no field '%s'", e.Name, f.Name)
			continue
		}
		actual := c.checkNode(f.Value)
		if !expected.Equals(actual) {
			c.errorf(e.Pos(), "field '%s': expected %s, got %s", f.Name, expected, actual)
		}
	}
	return &ObjType{Name: info.Name, Fields: info.Fields}
}

func (c *Checker) checkIfExpr(e *ast.IfExpr) ZType {
	condType := c.checkNode(e.Condition)
	if !condType.Equals(TypeBool) {
		c.errorf(e.Condition.Pos(), "if-expression condition must be bool, got %s", condType)
	}
	thenType := c.checkNode(e.Then)
	elseType := c.checkNode(e.Else)

	// For else-if chains, the elseType comes from the nested IfExpr
	if !thenType.Equals(elseType) {
		c.errorf(e.Pos(), "if-expression branches must have same type: then is %s, else is %s", thenType, elseType)
	}

	e.GoType = goTypeName(thenType)
	return thenType
}

func goTypeName(t ZType) string {
	switch ty := t.(type) {
	case *BuiltinType:
		switch ty.Name {
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
			return ty.Name
		}
	case *SliceType:
		return "[]" + goTypeName(ty.Elem)
	case *ObjType:
		return ty.Name
	default:
		return "interface{}"
	}
}
