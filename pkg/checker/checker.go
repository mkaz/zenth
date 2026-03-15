package checker

import (
	"fmt"
	"path"
	"strconv"
	"strings"

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
	Name     string
	Fields   map[string]ZType
	Order    []string        // field order for codegen
	Defaults map[string]bool // fields that have default values
}

// EnumInfo stores enum metadata.
type EnumInfo struct {
	Name     string
	Variants map[string]string // variant name -> string value
	Order    []string          // variant order
}

// Checker performs type checking and semantic analysis on a Zenth AST.
type Checker struct {
	scope           *Scope
	funcs           map[string]*FuncInfo // "name" or "Type.name"
	objs            map[string]*ObjInfo
	enums           map[string]*EnumInfo
	modules         map[string]bool
	typeAliases     map[string]*ast.TypeExpr
	errors          []string
	currentFunc     *FuncInfo // for checking return types
	pendingFlagName string    // set before checking a let/var value for flag()
}

// New creates a new Checker.
func New() *Checker {
	global := NewScope(nil)

	c := &Checker{
		scope:       global,
		funcs:       make(map[string]*FuncInfo),
		objs:        make(map[string]*ObjInfo),
		enums:       make(map[string]*EnumInfo),
		modules:     make(map[string]bool),
		typeAliases: make(map[string]*ast.TypeExpr),
	}

	// Known module names for qualified calls (fmt.println, math.sqrt, etc.).
	c.modules["fmt"] = true
	c.modules["math"] = true
	c.modules["os"] = true
	c.modules["strings"] = true
	c.modules["str"] = true
	c.modules["io"] = true

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
	c.funcs["exit"] = &FuncInfo{
		Name:        "exit",
		Params:      []ZType{TypeInt},
		Return:      TypeVoid,
		NumRequired: 0,
	}
	c.funcs["range"] = &FuncInfo{
		Name:   "range",
		Params: []ZType{TypeInt, TypeInt},
		Return: TypeRange,
	}
	c.funcs["rangei"] = &FuncInfo{
		Name:   "rangei",
		Params: []ZType{TypeInt, TypeInt},
		Return: TypeRangei,
	}
	c.funcs["file"] = &FuncInfo{
		Name:   "file",
		Params: []ZType{TypeStr},
		Return: TypeFile,
	}
	c.funcs["int"] = &FuncInfo{
		Name:   "int",
		Params: []ZType{TypeInt},
		Return: TypeInt,
	}
	c.funcs["f64"] = &FuncInfo{
		Name:   "f64",
		Params: []ZType{TypeF64},
		Return: TypeF64,
	}
	c.funcs["flag"] = &FuncInfo{
		Name:        "flag",
		Params:      []ZType{TypeVoid}, // placeholder; actual type inferred from default
		ParamNames:  []string{"default"},
		Return:      TypeVoid, // return type set dynamically
		NumRequired: 1,
	}

	// Register built-in constants
	global.Define(&Symbol{Name: "INT_MAX", Type: TypeInt, IsConst: true})
	global.Define(&Symbol{Name: "INT_MIN", Type: TypeInt, IsConst: true})

	return c
}

func (c *Checker) errorf(pos token.Pos, format string, args ...any) {
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
		case *ast.EnumDecl:
			c.registerEnum(s)
		case *ast.FnDecl:
			c.registerFunc(s)
		case *ast.ImportDecl:
			c.registerImport(s)
		case *ast.TypeAliasDecl:
			c.typeAliases[s.Name] = s.Type
		}
	}

	// Second pass: check all declarations
	for _, stmt := range prog.Stmts {
		c.checkNode(stmt)
	}

	if len(c.errors) > 0 {
		var result strings.Builder
		result.WriteString("type errors:\n")
		for _, e := range c.errors {
			result.WriteString("  " + e + "\n")
		}
		return fmt.Errorf("%s", result.String())
	}
	return nil
}

func (c *Checker) registerObj(s *ast.ObjDecl) {
	info := &ObjInfo{
		Name:     s.Name,
		Fields:   make(map[string]ZType),
		Defaults: make(map[string]bool),
	}
	for _, f := range s.Fields {
		t := c.resolveTypeExpr(f.Type)
		info.Fields[f.Name] = t
		info.Order = append(info.Order, f.Name)
		if f.Default != nil {
			defType := c.checkNode(f.Default)
			if !t.Equals(defType) && defType != TypeNil {
				c.errorf(s.Pos(), "default value for field '%s' has type %s, expected %s", f.Name, defType, t)
			}
			info.Defaults[f.Name] = true
		}
	}
	c.objs[s.Name] = info

	// Register methods defined inside the obj
	for _, m := range s.Methods {
		c.registerFunc(m)
	}
}

func (c *Checker) registerEnum(e *ast.EnumDecl) {
	info := &EnumInfo{
		Name:     e.Name,
		Variants: make(map[string]string),
	}
	for _, v := range e.Variants {
		val := v.Name // default: variant name is the string value
		if v.StrValue != nil {
			val = v.StrValue.Value
		}
		info.Variants[v.Name] = val
		info.Order = append(info.Order, v.Name)
	}
	c.enums[e.Name] = info
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

func (c *Checker) registerImport(imp *ast.ImportDecl) {
	if imp.Alias != "" {
		c.modules[imp.Alias] = true
		return
	}
	c.modules[path.Base(imp.Path)] = true
}

func (c *Checker) isModuleName(name string) bool {
	return c.modules[name]
}

func (c *Checker) resolveTypeExpr(t *ast.TypeExpr) ZType {
	if t == nil {
		return TypeVoid
	}
	if t.IsHashmap && len(t.Params) == 2 {
		return &HashmapType{
			Key:   c.resolveTypeExpr(t.Params[0]),
			Value: c.resolveTypeExpr(t.Params[1]),
		}
	}
	if t.IsSet && len(t.Params) > 0 {
		return &SetType{Elem: c.resolveTypeExpr(t.Params[0])}
	}
	if t.IsSlice && len(t.Params) > 0 {
		return &SliceType{Elem: c.resolveTypeExpr(t.Params[0])}
	}
	if t.IsTuple {
		elems := make([]ZType, 0, len(t.Params))
		for _, p := range t.Params {
			elems = append(elems, c.resolveTypeExpr(p))
		}
		var names []string
		if len(t.ParamNames) > 0 {
			names = t.ParamNames
		}
		return &TupleType{Elems: elems, Names: names}
	}
	if bt := LookupBuiltinType(t.Name); bt != nil {
		return bt
	}
	if st, ok := c.objs[t.Name]; ok {
		return &ObjType{Name: st.Name, Fields: st.Fields}
	}
	if ei, ok := c.enums[t.Name]; ok {
		return &EnumType{Name: ei.Name, Variants: ei.Variants}
	}
	if alias, ok := c.typeAliases[t.Name]; ok {
		return c.resolveTypeExpr(alias)
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
	case *ast.EnumDecl:
		return TypeVoid
	case *ast.InterfaceDecl:
		return TypeVoid // TODO: check interface
	case *ast.ImportDecl:
		return TypeVoid
	case *ast.TypeAliasDecl:
		return TypeVoid
	case *ast.Block:
		return c.checkBlock(n)
	case *ast.LetStmt:
		return c.checkLetStmt(n)
	case *ast.VarStmt:
		return c.checkVarStmt(n)
	case *ast.ConstStmt:
		return c.checkConstStmt(n)
	case *ast.TupleDestructStmt:
		return c.checkTupleDestructStmt(n)
	case *ast.ArrayDestructStmt:
		return c.checkArrayDestructStmt(n)
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
	case *ast.LoopStmt:
		return c.checkLoopStmt(n)
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
	case *ast.SliceExpr:
		return c.checkSliceExpr(n)
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
	case *ast.TupleLitExpr:
		return c.checkTupleLit(n)
	case *ast.NamedArgExpr:
		return c.checkNode(n.Value)
	case *ast.IfExpr:
		return c.checkIfExpr(n)
	case *ast.MatchExpr:
		return c.checkMatchExpr(n)
	case *ast.GroupedExpr:
		return c.checkNode(n.Expr)
	case *ast.ClosureExpr:
		return c.checkClosureExpr(n, nil)
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
	c.pendingFlagName = s.Name
	valType := c.checkNode(s.Value)
	c.pendingFlagName = ""
	if s.Type != nil {
		declared := c.resolveTypeExpr(s.Type)
		if c.isTypedEmptySliceAssignment(s.Value, declared, valType) {
			valType = declared
		} else if !declared.Equals(valType) && valType != TypeNil {
			c.errorf(s.Pos(), "type mismatch: cannot assign %s to %s", valType, declared)
		}
		valType = declared
	}
	sym := &Symbol{Name: s.Name, Type: valType, Mutable: false}
	if call, ok := s.Value.(*ast.CallExpr); ok && call.HashmapDefaultVal != nil {
		sym.DefaultExpr = call.HashmapDefaultVal
	}
	c.scope.Define(sym)
	return TypeVoid
}

func (c *Checker) checkVarStmt(s *ast.VarStmt) ZType {
	c.pendingFlagName = s.Name
	valType := c.checkNode(s.Value)
	c.pendingFlagName = ""
	if s.Type != nil {
		declared := c.resolveTypeExpr(s.Type)
		if c.isTypedEmptySliceAssignment(s.Value, declared, valType) {
			valType = declared
		} else if !declared.Equals(valType) && valType != TypeNil {
			c.errorf(s.Pos(), "type mismatch: cannot assign %s to %s", valType, declared)
		}
		valType = declared
	}
	sym := &Symbol{Name: s.Name, Type: valType, Mutable: true}
	if call, ok := s.Value.(*ast.CallExpr); ok && call.HashmapDefaultVal != nil {
		sym.DefaultExpr = call.HashmapDefaultVal
	}
	c.scope.Define(sym)
	return TypeVoid
}

func (c *Checker) checkConstStmt(s *ast.ConstStmt) ZType {
	valType := c.checkNode(s.Value)
	if s.Type != nil {
		declared := c.resolveTypeExpr(s.Type)
		if c.isTypedEmptySliceAssignment(s.Value, declared, valType) {
			valType = declared
		} else if !declared.Equals(valType) && valType != TypeNil {
			c.errorf(s.Pos(), "type mismatch: cannot assign %s to %s", valType, declared)
		}
		valType = declared
	}
	c.scope.Define(&Symbol{Name: s.Name, Type: valType, IsConst: true})
	return TypeVoid
}

func (c *Checker) checkTupleDestructStmt(s *ast.TupleDestructStmt) ZType {
	valueType := c.checkNode(s.Value)
	tt, ok := valueType.(*TupleType)
	if !ok {
		c.errorf(s.Pos(), "tuple destructuring requires tuple value, got %s", valueType)
		for _, name := range s.Names {
			if name == "_" {
				continue
			}
			sym := &Symbol{Name: name, Type: TypeVoid}
			switch s.Kind {
			case token.Var:
				sym.Mutable = true
			case token.Const:
				sym.IsConst = true
			}
			c.scope.Define(sym)
		}
		return TypeVoid
	}
	if len(tt.Elems) != len(s.Names) {
		c.errorf(s.Pos(), "tuple destructuring arity mismatch: %d names, %d values", len(s.Names), len(tt.Elems))
		return TypeVoid
	}
	s.ElemGoTypes = make([]string, len(tt.Elems))
	for i, name := range s.Names {
		s.ElemGoTypes[i] = goTypeName(tt.Elems[i])
		if name == "_" {
			continue
		}
		sym := &Symbol{Name: name, Type: tt.Elems[i]}
		switch s.Kind {
		case token.Var:
			sym.Mutable = true
		case token.Const:
			sym.IsConst = true
		}
		c.scope.Define(sym)
	}
	return TypeVoid
}

func (c *Checker) checkArrayDestructStmt(s *ast.ArrayDestructStmt) ZType {
	valueType := c.checkNode(s.Value)
	st, ok := valueType.(*SliceType)
	if !ok {
		c.errorf(s.Pos(), "array destructuring requires array value, got %s", valueType)
		for _, name := range s.Names {
			if name == "_" {
				continue
			}
			sym := &Symbol{Name: name, Type: TypeVoid}
			switch s.Kind {
			case token.Var:
				sym.Mutable = true
			case token.Const:
				sym.IsConst = true
			}
			c.scope.Define(sym)
		}
		return TypeVoid
	}
	s.ElemType = goTypeName(st.Elem)
	for _, name := range s.Names {
		if name == "_" {
			continue
		}
		sym := &Symbol{Name: name, Type: st.Elem}
		switch s.Kind {
		case token.Var:
			sym.Mutable = true
		case token.Const:
			sym.IsConst = true
		}
		c.scope.Define(sym)
	}
	return TypeVoid
}

func (c *Checker) checkAssignStmt(s *ast.AssignStmt) ZType {
	targetType := c.checkNode(s.Target)
	valueType := c.checkNode(s.Value)

	if field, ok := s.Target.(*ast.FieldExpr); ok && field.TupleAccess {
		c.errorf(s.Pos(), "cannot assign to tuple element .%d", field.TupleIndex)
	}

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
		if c.isTypedEmptySliceAssignment(s.Value, targetType, valueType) {
			return TypeVoid
		}
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
		if field, ok := target.(*ast.FieldExpr); ok && field.TupleAccess {
			c.errorf(s.Pos(), "cannot assign to tuple element .%d", field.TupleIndex)
		}
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
			if c.isTypedEmptySliceAssignment(s.Values[i], targetType, valueType) {
				continue
			}
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

	// Helper to define a loop variable in scope, skipping "_" discard
	defineLoopVar := func(name string, typ ZType) {
		if name != "_" {
			c.scope.Define(&Symbol{Name: name, Type: typ})
		}
	}

	c.pushScope()
	switch t := iterType.(type) {
	case *RangeType:
		s.IterRange = true
		s.IterRangeInclusive = t.Inclusive
		if s.Index != "" {
			defineLoopVar(s.Index, TypeInt)
		}
		defineLoopVar(s.Value, TypeInt)
	case *SliceType:
		if s.Index != "" {
			defineLoopVar(s.Index, TypeInt)
		}
		defineLoopVar(s.Value, t.Elem)
	case *HashmapType:
		s.IterHashmap = true
		if objKey, ok := t.Key.(*ObjType); ok {
			s.IterHashmapObjKey = true
			s.IterHashmapObjType = objKey.Name
		}
		if s.Index != "" {
			defineLoopVar(s.Index, t.Key)
			defineLoopVar(s.Value, t.Value)
		} else {
			// for v in m iterates over keys only
			defineLoopVar(s.Value, t.Key)
		}
	case *SetType:
		s.IterSet = true
		if s.Index != "" {
			c.errorf(s.Pos(), "set iteration does not support index variable")
		}
		if tt, ok := t.Elem.(*TupleType); ok {
			s.IterSetTupleStruct = tupleStructName(tt)
			s.IterSetTupleFieldTypes = tupleFieldGoTypes(tt)
		}
		defineLoopVar(s.Value, t.Elem)
	default:
		if iterType.Equals(TypeStr) {
			s.IterStr = true
			if s.Index != "" {
				defineLoopVar(s.Index, TypeInt)
			}
			defineLoopVar(s.Value, TypeStr)
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

func (c *Checker) checkLoopStmt(s *ast.LoopStmt) ZType {
	countType := c.checkNode(s.Count)
	if !countType.Equals(TypeInt) {
		c.errorf(s.Count.Pos(), "for-range count must be int, got %s", countType)
	}
	for _, stmt := range s.Body.Stmts {
		c.checkNode(stmt)
	}
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
		if promoted := PromoteNumeric(left, right); promoted != nil {
			c.annotatePromotion(e, left, right, promoted)
			return promoted
		}
		// Slice concatenation
		if ls, ok := left.(*SliceType); ok {
			if rs, ok := right.(*SliceType); ok {
				if ls.Elem.Equals(rs.Elem) {
					e.SliceConcat = true
					return left
				}
				c.errorf(e.Pos(), "cannot concatenate %s and %s (element types differ)", left, right)
				return left
			}
		}
		c.errorf(e.Pos(), "cannot add %s and %s", left, right)
		return left

	case token.Minus, token.Star, token.Slash, token.Percent:
		if IsNumeric(left) && left.Equals(right) {
			return left
		}
		if promoted := PromoteNumeric(left, right); promoted != nil {
			c.annotatePromotion(e, left, right, promoted)
			return promoted
		}
		c.errorf(e.Pos(), "cannot use %s with %s and %s", e.Op, left, right)
		return left

	case token.Eq, token.Neq:
		if IsNumeric(left) && IsNumeric(right) && !left.Equals(right) {
			if promoted := PromoteNumeric(left, right); promoted != nil {
				c.annotatePromotion(e, left, right, promoted)
			}
		}
		return TypeBool

	case token.Lt, token.Gt, token.Lte, token.Gte:
		if IsNumeric(left) && left.Equals(right) {
			return TypeBool
		}
		if promoted := PromoteNumeric(left, right); promoted != nil {
			c.annotatePromotion(e, left, right, promoted)
			return TypeBool
		}
		c.errorf(e.Pos(), "cannot compare %s and %s", left, right)
		return TypeBool

	case token.In:
		c.errorf(e.Pos(), "'in' is not a binary operator; use .exists() for slices and .contains() for strings")
		return TypeBool

	case token.And, token.Or:
		if !left.Equals(TypeBool) || !right.Equals(TypeBool) {
			c.errorf(e.Pos(), "logical operators require bool operands, got %s and %s", left, right)
		}
		return TypeBool
	}

	return TypeVoid
}

func (c *Checker) annotatePromotion(e *ast.BinaryExpr, left, right, promoted ZType) {
	if !left.Equals(promoted) {
		e.PromoteLeft = goTypeName(promoted)
	}
	if !right.Equals(promoted) {
		e.PromoteRight = goTypeName(promoted)
	}
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
		// Module function call (fmt.println, math.sqrt, etc.)
		if ident, ok := field.Object.(*ast.IdentExpr); ok && c.isModuleName(ident.Name) {
			for _, arg := range e.Args {
				if named, ok := arg.(*ast.NamedArgExpr); ok {
					c.errorf(named.Pos(), "named arguments are not supported for module call %s.%s", ident.Name, field.Field)
					c.checkNode(named.Value)
					continue
				}
				c.checkNode(arg)
			}
			return TypeVoid
		}

		objType := c.checkNode(field.Object)
		if objType.Equals(TypeVoid) {
			for _, arg := range e.Args {
				c.checkNode(arg)
			}
			return TypeVoid
		}
		methodKey := ""
		if st, ok := objType.(*ObjType); ok {
			methodKey = st.Name + "." + field.Field
		}
		if info, ok := c.funcs[methodKey]; ok {
			e.ResolvedFunc = methodKey
			c.checkArgs(e, info)
			return info.Return
		}
		// Check for built-in slice methods
		if sliceType, ok := objType.(*SliceType); ok {
			switch field.Field {
			case "pop":
				if len(e.Args) > 1 {
					c.errorf(e.Pos(), "pop() takes 0 or 1 arguments, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !IsInteger(argType) {
						c.errorf(e.Args[0].Pos(), "pop() index must be integer, got %s", argType)
					}
				}
				e.SliceMethod = true
				return sliceType.Elem
			case "add":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "add() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !sliceType.Elem.Equals(argType) {
						c.errorf(e.Args[0].Pos(), "add() argument type %s does not match slice element type %s", argType, sliceType.Elem)
					}
				}
				e.SliceMethod = true
				return TypeVoid
			case "push":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "push() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !sliceType.Elem.Equals(argType) {
						c.errorf(e.Args[0].Pos(), "push() argument type %s does not match slice element type %s", argType, sliceType.Elem)
					}
				}
				e.SliceMethod = true
				return TypeVoid
			case "length":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "length() takes no arguments, got %d", len(e.Args))
				}
				e.SliceMethod = true
				return TypeInt
			case "exists":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "exists() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !sliceType.Elem.Equals(argType) {
						c.errorf(e.Args[0].Pos(), "exists() argument type %s does not match slice element type %s", argType, sliceType.Elem)
					}
				}
				e.SliceMethod = true
				return TypeBool
			case "map":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "map() takes exactly 1 argument, got %d", len(e.Args))
					return &SliceType{Elem: TypeVoid}
				}
				closure, ok := e.Args[0].(*ast.ClosureExpr)
				if !ok {
					c.errorf(e.Args[0].Pos(), "map() argument must be a closure")
					c.checkNode(e.Args[0])
					return &SliceType{Elem: TypeVoid}
				}
				closureType := c.checkClosureExpr(closure, sliceType.Elem)
				ft, ok := closureType.(*FuncType)
				if !ok {
					return &SliceType{Elem: TypeVoid}
				}
				e.SliceMethod = true
				return &SliceType{Elem: ft.Returns}
			case "filter":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "filter() takes exactly 1 argument, got %d", len(e.Args))
					return sliceType
				}
				closure, ok := e.Args[0].(*ast.ClosureExpr)
				if !ok {
					c.errorf(e.Args[0].Pos(), "filter() argument must be a closure")
					c.checkNode(e.Args[0])
					return sliceType
				}
				closureType := c.checkClosureExpr(closure, sliceType.Elem)
				if ft, ok := closureType.(*FuncType); ok {
					if !ft.Returns.Equals(TypeBool) {
						c.errorf(e.Args[0].Pos(), "filter() closure must return bool, got %s", ft.Returns)
					}
				}
				e.SliceMethod = true
				return sliceType
			case "to_int":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "to_int() takes no arguments, got %d", len(e.Args))
				}
				if !sliceType.Elem.Equals(TypeStr) {
					c.errorf(e.Pos(), "to_int() requires []str, got %s", objType)
				}
				e.SliceMethod = true
				e.SliceConvTarget = "int"
				return &SliceType{Elem: TypeInt}
			case "to_f64":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "to_f64() takes no arguments, got %d", len(e.Args))
				}
				if !sliceType.Elem.Equals(TypeStr) {
					c.errorf(e.Pos(), "to_f64() requires []str, got %s", objType)
				}
				e.SliceMethod = true
				e.SliceConvTarget = "float64"
				return &SliceType{Elem: TypeF64}
			case "to_str":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "to_str() takes no arguments, got %d", len(e.Args))
				}
				if !sliceType.Elem.Equals(TypeInt) && !sliceType.Elem.Equals(TypeF64) {
					c.errorf(e.Pos(), "to_str() requires []int or []f64, got %s", objType)
				}
				e.SliceMethod = true
				e.SliceConvTarget = "string"
				return &SliceType{Elem: TypeStr}
			case "max":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "max() takes no arguments, got %d", len(e.Args))
				}
				if !IsNumeric(sliceType.Elem) {
					c.errorf(e.Pos(), "max() requires a numeric array, got %s", objType)
				}
				e.SliceMethod = true
				return sliceType.Elem
			case "min":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "min() takes no arguments, got %d", len(e.Args))
				}
				if !IsNumeric(sliceType.Elem) {
					c.errorf(e.Pos(), "min() requires a numeric array, got %s", objType)
				}
				e.SliceMethod = true
				return sliceType.Elem
			case "sum":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "sum() takes no arguments, got %d", len(e.Args))
				}
				if !IsNumeric(sliceType.Elem) {
					c.errorf(e.Pos(), "sum() requires a numeric array, got %s", objType)
				}
				e.SliceMethod = true
				return sliceType.Elem
			case "sorted":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "sorted() takes no arguments, got %d", len(e.Args))
				}
				if !IsNumeric(sliceType.Elem) && !sliceType.Elem.Equals(TypeStr) {
					c.errorf(e.Pos(), "sorted() requires a numeric or string array, got %s", objType)
				}
				e.SliceMethod = true
				return sliceType
			case "reduce":
				if len(e.Args) < 1 || len(e.Args) > 2 {
					c.errorf(e.Pos(), "reduce() takes 1 or 2 arguments (closure [, initial]), got %d", len(e.Args))
					return TypeVoid
				}
				closure, ok := e.Args[0].(*ast.ClosureExpr)
				if !ok {
					c.errorf(e.Args[0].Pos(), "reduce() first argument must be a closure")
					c.checkNode(e.Args[0])
					return TypeVoid
				}
				if len(closure.Params) != 2 {
					c.errorf(e.Args[0].Pos(), "reduce() closure must take exactly 2 parameters, got %d", len(closure.Params))
					return TypeVoid
				}
				closureType := c.checkClosureExpr(closure, sliceType.Elem)
				ft, ok := closureType.(*FuncType)
				if !ok {
					return TypeVoid
				}
				if !ft.Returns.Equals(sliceType.Elem) {
					c.errorf(e.Args[0].Pos(), "reduce() closure must return %s, got %s", sliceType.Elem, ft.Returns)
				}
				if len(e.Args) == 2 {
					initType := c.checkNode(e.Args[1])
					if !initType.Equals(sliceType.Elem) {
						c.errorf(e.Args[1].Pos(), "reduce() initial value type %s does not match element type %s", initType, sliceType.Elem)
					}
				}
				e.SliceMethod = true
				return sliceType.Elem
			case "insert":
				if len(e.Args) != 2 {
					c.errorf(e.Pos(), "insert() takes exactly 2 arguments (index, elem), got %d", len(e.Args))
				}
				if len(e.Args) >= 1 {
					argType := c.checkNode(e.Args[0])
					if !IsInteger(argType) {
						c.errorf(e.Args[0].Pos(), "insert() index must be integer, got %s", argType)
					}
				}
				if len(e.Args) == 2 {
					argType := c.checkNode(e.Args[1])
					if !sliceType.Elem.Equals(argType) {
						c.errorf(e.Args[1].Pos(), "insert() element type %s does not match slice element type %s", argType, sliceType.Elem)
					}
				}
				e.SliceMethod = true
				return TypeVoid
			case "remove":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "remove() takes exactly 1 argument (index), got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !IsInteger(argType) {
						c.errorf(e.Args[0].Pos(), "remove() index must be integer, got %s", argType)
					}
				}
				e.SliceMethod = true
				return TypeVoid
			}
		}
		// Check for built-in file methods
		if objType.Equals(TypeFile) {
			switch field.Field {
			case "exists":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "exists() takes no arguments, got %d", len(e.Args))
				}
				e.SliceMethod = true
				return TypeBool
			case "read":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "read() takes no arguments, got %d", len(e.Args))
				}
				e.SliceMethod = true
				return TypeStr
			case "lines":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "lines() takes no arguments, got %d", len(e.Args))
				}
				e.SliceMethod = true
				return &SliceType{Elem: TypeStr}
			case "name":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "name() takes no arguments, got %d", len(e.Args))
				}
				e.SliceMethod = true
				return TypeStr
			case "ext":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "ext() takes no arguments, got %d", len(e.Args))
				}
				e.SliceMethod = true
				return TypeStr
			case "sections":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "sections() takes no arguments, got %d", len(e.Args))
				}
				e.SliceMethod = true
				return &SliceType{Elem: TypeStr}
			}
		}
		// Check for built-in range methods
		if _, ok := objType.(*RangeType); ok {
			switch field.Field {
			case "contains":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "contains() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !IsInteger(argType) {
						c.errorf(e.Args[0].Pos(), "contains() argument must be int, got %s", argType)
					}
				}
				e.RangeMethod = true
				return TypeBool
			default:
				c.errorf(e.Pos(), "range has no method %q", field.Field)
				for _, arg := range e.Args {
					c.checkNode(arg)
				}
				return TypeVoid
			}
		}
		// Check for built-in string methods
		if objType.Equals(TypeStr) {
			switch field.Field {
			case "to_int":
				if len(e.Args) > 1 {
					c.errorf(e.Pos(), "to_int() takes 0 or 1 arguments, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !IsInteger(argType) {
						c.errorf(e.Args[0].Pos(), "to_int() base must be int, got %s", argType)
					}
				}
				e.SliceMethod = true
				return TypeInt
			case "to_f64":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "to_f64() takes no arguments, got %d", len(e.Args))
				}
				e.SliceMethod = true
				return TypeF64
			}
		}

		// Check for built-in string methods
		if objType.Equals(TypeStr) {
			switch field.Field {
			case "split":
				if len(e.Args) > 1 {
					c.errorf(e.Pos(), "split() takes 0 or 1 arguments, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "split() separator must be str, got %s", argType)
					}
				}
				e.StringMethod = "split"
				return &SliceType{Elem: TypeStr}
			case "split_once":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "split_once() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "split_once() separator must be str, got %s", argType)
					}
				}
				e.StringMethod = "split_once"
				return &TupleType{Elems: []ZType{TypeStr, TypeStr}}
			case "upper":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "upper() takes no arguments, got %d", len(e.Args))
				}
				e.StringMethod = "upper"
				return TypeStr
			case "lower":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "lower() takes no arguments, got %d", len(e.Args))
				}
				e.StringMethod = "lower"
				return TypeStr
			case "starts_with":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "starts_with() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "starts_with() argument must be str, got %s", argType)
					}
				}
				e.StringMethod = "starts_with"
				return TypeBool
			case "ends_with":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "ends_with() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "ends_with() argument must be str, got %s", argType)
					}
				}
				e.StringMethod = "ends_with"
				return TypeBool
			case "strip":
				if len(e.Args) > 1 {
					c.errorf(e.Pos(), "strip() takes 0 or 1 arguments, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "strip() argument must be str, got %s", argType)
					}
				}
				e.StringMethod = "strip"
				return TypeStr
			case "find":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "find() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "find() argument must be str, got %s", argType)
					}
				}
				e.StringMethod = "find"
				return TypeInt
			case "count":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "count() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "count() argument must be str, got %s", argType)
					}
				}
				e.StringMethod = "count"
				return TypeInt
			case "replace":
				if len(e.Args) < 2 || len(e.Args) > 3 {
					c.errorf(e.Pos(), "replace() takes 2 or 3 arguments, got %d", len(e.Args))
				}
				if len(e.Args) >= 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "replace() argument 1 must be str, got %s", argType)
					}
				}
				if len(e.Args) >= 2 {
					argType := c.checkNode(e.Args[1])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[1].Pos(), "replace() argument 2 must be str, got %s", argType)
					}
				}
				if len(e.Args) == 3 {
					argType := c.checkNode(e.Args[2])
					if !IsInteger(argType) {
						c.errorf(e.Args[2].Pos(), "replace() argument 3 must be int, got %s", argType)
					}
				}
				e.StringMethod = "replace"
				return TypeStr
			case "length":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "length() takes no arguments, got %d", len(e.Args))
				}
				e.StringMethod = "length"
				return TypeInt
			case "contains":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "contains() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "contains() argument must be str, got %s", argType)
					}
				}
				e.StringMethod = "contains"
				return TypeBool
			case "strip_prefix":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "strip_prefix() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "strip_prefix() argument must be str, got %s", argType)
					}
				}
				e.StringMethod = "strip_prefix"
				return TypeStr
			case "strip_suffix":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "strip_suffix() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !argType.Equals(TypeStr) {
						c.errorf(e.Args[0].Pos(), "strip_suffix() argument must be str, got %s", argType)
					}
				}
				e.StringMethod = "strip_suffix"
				return TypeStr
			case "repeat":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "repeat() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !IsInteger(argType) {
						c.errorf(e.Args[0].Pos(), "repeat() argument must be int, got %s", argType)
					}
				}
				e.StringMethod = "repeat"
				return TypeStr
			case "is_digit":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "is_digit() takes no arguments, got %d", len(e.Args))
				}
				e.StringMethod = "is_digit"
				return TypeBool
			}
		}

		if IsInteger(objType) {
			switch field.Field {
			case "to_base":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "to_base() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !IsInteger(argType) {
						c.errorf(e.Args[0].Pos(), "to_base() argument must be int, got %s", argType)
					}
				}
				e.StringMethod = "to_base"
				return TypeStr
			}
		}

		// Check for built-in hashmap methods
		if hmType, ok := objType.(*HashmapType); ok {
			switch field.Field {
			case "keys":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "keys() takes no arguments, got %d", len(e.Args))
				}
				e.HashmapMethod = "keys"
				return &SliceType{Elem: hmType.Key}
			case "values":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "values() takes no arguments, got %d", len(e.Args))
				}
				e.HashmapMethod = "values"
				return &SliceType{Elem: hmType.Value}
			case "exists":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "exists() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !hmType.Key.Equals(argType) {
						c.errorf(e.Args[0].Pos(), "exists() argument type %s does not match hashmap key type %s", argType, hmType.Key)
					}
				}
				e.HashmapMethod = "exists"
				return TypeBool
			}
		}

		// Check for built-in set methods
		if setType, ok := objType.(*SetType); ok {
			// Helper to annotate tuple struct info on the call expression
			annotateTupleStruct := func() {
				if tt, ok := setType.Elem.(*TupleType); ok {
					e.SetTupleStruct = tupleStructName(tt)
					e.SetTupleFieldTypes = tupleFieldGoTypes(tt)
				}
			}
			switch field.Field {
			case "add":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "add() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !setType.Elem.Equals(argType) {
						c.errorf(e.Args[0].Pos(), "add() argument type %s does not match set element type %s", argType, setType.Elem)
					}
				}
				e.SetMethod = "add"
				annotateTupleStruct()
				return TypeVoid
			case "exists":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "exists() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !setType.Elem.Equals(argType) {
						c.errorf(e.Args[0].Pos(), "exists() argument type %s does not match set element type %s", argType, setType.Elem)
					}
				}
				e.SetMethod = "exists"
				annotateTupleStruct()
				return TypeBool
			case "remove":
				if len(e.Args) != 1 {
					c.errorf(e.Pos(), "remove() takes exactly 1 argument, got %d", len(e.Args))
				}
				if len(e.Args) == 1 {
					argType := c.checkNode(e.Args[0])
					if !setType.Elem.Equals(argType) {
						c.errorf(e.Args[0].Pos(), "remove() argument type %s does not match set element type %s", argType, setType.Elem)
					}
				}
				e.SetMethod = "remove"
				annotateTupleStruct()
				return TypeVoid
			case "length":
				if len(e.Args) != 0 {
					c.errorf(e.Pos(), "length() takes no arguments, got %d", len(e.Args))
				}
				e.SetMethod = "length"
				return TypeInt
			}
		}

		// Unknown method call on a typed value should be a semantic error.
		for _, arg := range e.Args {
			c.checkNode(arg)
		}
		c.errorf(e.Pos(), "type %s has no method '%s'", objType, field.Field)
		return TypeVoid
	}

	// Handle free function calls
	if ident, ok := e.Callee.(*ast.IdentExpr); ok {
		if isNumericBuiltin(ident.Name) {
			return c.checkNumericBuiltinCall(e, ident.Name)
		}
		if ident.Name == "hashmap" {
			return c.checkHashmapConstructor(e)
		}
		if ident.Name == "set" {
			return c.checkSetConstructor(e)
		}
		if ident.Name == "flag" {
			return c.checkFlagCall(e)
		}
		if ident.Name == "assert" {
			return c.checkAssertCall(e)
		}
		if ident.Name == "assert_eq" {
			return c.checkAssertEqCall(e)
		}
		// Check if this is an obj constructor call
		if _, ok := c.objs[ident.Name]; ok {
			return c.checkObjConstructor(e, ident.Name)
		}
		if info, ok := c.funcs[ident.Name]; ok {
			e.ResolvedFunc = ident.Name
			c.checkArgs(e, info)
			// int(str, base) — two-argument form
			if info.Name == "int" && len(e.Args) == 2 {
				argType := c.checkNode(e.Args[0])
				baseType := c.checkNode(e.Args[1])
				if !argType.Equals(TypeStr) {
					c.errorf(e.Args[0].Pos(), "int() with base requires first argument to be str, got %s", argType)
				}
				if !IsInteger(baseType) {
					c.errorf(e.Args[1].Pos(), "int() base must be int, got %s", baseType)
				}
				e.IntBaseCall = true
				return TypeInt
			}
			// Conversion builtins on slices: int([]str) -> []int, etc.
			if (info.Name == "int" || info.Name == "f64" || info.Name == "str") && len(e.Args) == 1 {
				argType := c.checkNode(e.Args[0])
				if st, ok := argType.(*SliceType); ok {
					switch info.Name {
					case "int":
						if st.Elem.Equals(TypeStr) {
							e.SliceConvFunc = "int"
							return &SliceType{Elem: TypeInt}
						}
					case "f64":
						if st.Elem.Equals(TypeStr) {
							e.SliceConvFunc = "f64"
							return &SliceType{Elem: TypeF64}
						}
					case "str":
						e.SliceConvFunc = "str"
						return &SliceType{Elem: TypeStr}
					}
				}
			}
			return info.Return
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

func isNumericBuiltin(name string) bool {
	switch name {
	case "abs", "min", "max", "clamp", "round", "floor", "ceil", "pow", "sqrt":
		return true
	default:
		return false
	}
}

func (c *Checker) getPositionalArgsForBuiltin(e *ast.CallExpr, name string) []ast.Node {
	args := make([]ast.Node, 0, len(e.Args))
	for _, arg := range e.Args {
		if named, ok := arg.(*ast.NamedArgExpr); ok {
			c.errorf(named.Pos(), "named arguments are not supported for %s()", name)
			args = append(args, named.Value)
			continue
		}
		args = append(args, arg)
	}
	return args
}

func isBasicNumberType(t ZType) bool {
	return t.Equals(TypeInt) || t.Equals(TypeF64)
}

func (c *Checker) checkNumericBuiltinCall(e *ast.CallExpr, name string) ZType {
	args := c.getPositionalArgsForBuiltin(e, name)
	switch name {
	case "abs":
		if len(args) != 1 {
			c.errorf(e.Pos(), "abs() takes exactly 1 argument, got %d", len(args))
			for _, arg := range args {
				c.checkNode(arg)
			}
			return TypeVoid
		}
		t := c.checkNode(args[0])
		if t.Equals(TypeInt) {
			e.NumericMethod = "abs_int"
			return TypeInt
		}
		if t.Equals(TypeF64) {
			e.NumericMethod = "abs_f64"
			return TypeF64
		}
		c.errorf(args[0].Pos(), "abs() argument must be int or f64, got %s", t)
		return TypeVoid
	case "min", "max":
		if len(args) < 2 {
			c.errorf(e.Pos(), "%s() requires at least 2 arguments, got %d", name, len(args))
			for _, arg := range args {
				c.checkNode(arg)
			}
			return TypeVoid
		}
		allInt := true
		allF64 := true
		for _, arg := range args {
			t := c.checkNode(arg)
			if !t.Equals(TypeInt) {
				allInt = false
			}
			if !t.Equals(TypeF64) {
				allF64 = false
			}
		}
		if allInt {
			if len(args) == 2 {
				e.NumericMethod = name + "_int"
			} else {
				e.NumericMethod = name + "_int_variadic"
			}
			return TypeInt
		}
		if allF64 {
			if len(args) == 2 {
				e.NumericMethod = name + "_f64"
			} else {
				e.NumericMethod = name + "_f64_variadic"
			}
			return TypeF64
		}
		c.errorf(e.Pos(), "%s() arguments must all be int or all be f64", name)
		return TypeVoid
	case "clamp":
		if len(args) != 3 {
			c.errorf(e.Pos(), "clamp() takes exactly 3 arguments, got %d", len(args))
			for _, arg := range args {
				c.checkNode(arg)
			}
			return TypeVoid
		}
		t1 := c.checkNode(args[0])
		t2 := c.checkNode(args[1])
		t3 := c.checkNode(args[2])
		if t1.Equals(TypeInt) && t2.Equals(TypeInt) && t3.Equals(TypeInt) {
			e.NumericMethod = "clamp_int"
			return TypeInt
		}
		if t1.Equals(TypeF64) && t2.Equals(TypeF64) && t3.Equals(TypeF64) {
			e.NumericMethod = "clamp_f64"
			return TypeF64
		}
		c.errorf(e.Pos(), "clamp() arguments must all be int or all be f64, got %s, %s, %s", t1, t2, t3)
		return TypeVoid
	case "round", "floor", "ceil", "sqrt":
		if len(args) != 1 {
			c.errorf(e.Pos(), "%s() takes exactly 1 argument, got %d", name, len(args))
			for _, arg := range args {
				c.checkNode(arg)
			}
			return TypeVoid
		}
		t := c.checkNode(args[0])
		if !isBasicNumberType(t) {
			c.errorf(args[0].Pos(), "%s() argument must be int or f64, got %s", name, t)
			return TypeVoid
		}
		e.NumericMethod = name
		return TypeF64
	case "pow":
		if len(args) != 2 {
			c.errorf(e.Pos(), "pow() takes exactly 2 arguments, got %d", len(args))
			for _, arg := range args {
				c.checkNode(arg)
			}
			return TypeVoid
		}
		t1 := c.checkNode(args[0])
		t2 := c.checkNode(args[1])
		if !isBasicNumberType(t1) {
			c.errorf(args[0].Pos(), "pow() argument 1 must be int or f64, got %s", t1)
		}
		if !isBasicNumberType(t2) {
			c.errorf(args[1].Pos(), "pow() argument 2 must be int or f64, got %s", t2)
		}
		e.NumericMethod = "pow"
		return TypeF64
	default:
		for _, arg := range args {
			c.checkNode(arg)
		}
		return TypeVoid
	}
}

func (c *Checker) checkArgs(e *ast.CallExpr, info *FuncInfo) {
	// Flexible arg count for variadic builtins (print, println, etc.)
	if info.Name == "print" || info.Name == "println" {
		if len(e.Args) < 1 || len(e.Args) > 2 {
			c.errorf(e.Pos(), "%s() expects 1 or 2 arguments, got %d", info.Name, len(e.Args))
			for _, arg := range e.Args {
				if named, ok := arg.(*ast.NamedArgExpr); ok {
					c.errorf(named.Pos(), "named arguments are not supported for %s()", info.Name)
					c.checkNode(named.Value)
					continue
				}
				c.checkNode(arg)
			}
			return
		}
		if named, ok := e.Args[0].(*ast.NamedArgExpr); ok {
			c.errorf(named.Pos(), "named arguments are not supported for %s()", info.Name)
			c.checkNode(named.Value)
		} else {
			c.checkNode(e.Args[0])
		}
		if len(e.Args) == 2 {
			if named, ok := e.Args[1].(*ast.NamedArgExpr); ok {
				c.errorf(named.Pos(), "named arguments are not supported for %s()", info.Name)
				argType := c.checkNode(named.Value)
				if !argType.Equals(TypeBool) {
					c.errorf(named.Pos(), "second argument to %s must be bool, got %s", info.Name, argType)
				}
			} else {
				argType := c.checkNode(e.Args[1])
				if !argType.Equals(TypeBool) {
					c.errorf(e.Args[1].Pos(), "second argument to %s must be bool, got %s", info.Name, argType)
				}
			}
		}
		return
	}
	// Conversion builtins accept any single arg (int() also accepts 2 for base)
	if info.Name == "int" || info.Name == "f64" || info.Name == "str" {
		for _, arg := range e.Args {
			if named, ok := arg.(*ast.NamedArgExpr); ok {
				c.errorf(named.Pos(), "named arguments are not supported for %s()", info.Name)
				c.checkNode(named.Value)
				continue
			}
			c.checkNode(arg)
		}
		if info.Name == "int" && len(e.Args) == 2 {
			// int(str, base) form
			return
		}
		if len(e.Args) != 1 {
			c.errorf(e.Pos(), "%s() expects exactly 1 argument, got %d", info.Name, len(e.Args))
		}
		return
	}
	if info.Name == "len" {
		for _, arg := range e.Args {
			if named, ok := arg.(*ast.NamedArgExpr); ok {
				c.errorf(named.Pos(), "named arguments are not supported for len()")
				c.checkNode(named.Value)
				continue
			}
		}
		if len(e.Args) != 1 {
			c.errorf(e.Pos(), "len() expects exactly 1 argument, got %d", len(e.Args))
			return
		}
		argType := c.checkNode(e.Args[0])
		if _, ok := argType.(*RangeType); ok {
			e.LenArgIsRange = true
		} else if _, ok := argType.(*SliceType); !ok && !argType.Equals(TypeStr) {
			if _, ok := argType.(*HashmapType); !ok {
				if _, ok := argType.(*SetType); !ok {
					c.errorf(e.Args[0].Pos(), "argument 1 to len has type %s, expected str, slice, hashmap, set, or range", argType)
				}
			}
		}
		return
	}
	// range/rangei accept 2 or 3 int args
	if info.Name == "range" || info.Name == "rangei" {
		for i, arg := range e.Args {
			if named, ok := arg.(*ast.NamedArgExpr); ok {
				c.errorf(named.Pos(), "named arguments are not supported for %s()", info.Name)
				c.checkNode(named.Value)
				continue
			}
			argType := c.checkNode(arg)
			if !IsInteger(argType) {
				c.errorf(arg.Pos(), "argument %d to %s must be int, got %s", i+1, info.Name, argType)
			}
		}
		if len(e.Args) < 2 || len(e.Args) > 3 {
			c.errorf(e.Pos(), "%s expects 2 or 3 arguments, got %d", info.Name, len(e.Args))
		}
		return
	}

	resolved := make([]ast.Node, len(info.Params))
	assigned := make([]bool, len(info.Params))
	positionalIndex := 0
	seenNamed := false

	for _, arg := range e.Args {
		if named, ok := arg.(*ast.NamedArgExpr); ok {
			seenNamed = true
			idx := -1
			for i, name := range info.ParamNames {
				if name == named.Name {
					idx = i
					break
				}
			}
			if idx < 0 {
				c.errorf(named.Pos(), "%s() has no parameter named '%s'", info.Name, named.Name)
				c.checkNode(named.Value)
				continue
			}
			if assigned[idx] {
				c.errorf(named.Pos(), "duplicate argument for parameter '%s' in %s()", named.Name, info.Name)
				c.checkNode(named.Value)
				continue
			}
			valType := c.checkNode(named.Value)
			if !info.Params[idx].Equals(valType) && valType != TypeNil {
				c.errorf(named.Pos(), "argument '%s' to %s has type %s, expected %s", named.Name, info.Name, valType, info.Params[idx])
			}
			resolved[idx] = named.Value
			assigned[idx] = true
			continue
		}

		if seenNamed {
			c.errorf(arg.Pos(), "positional argument cannot follow named arguments in %s()", info.Name)
		}
		if positionalIndex >= len(info.Params) {
			c.checkNode(arg)
			positionalIndex++
			continue
		}
		valType := c.checkNode(arg)
		if !info.Params[positionalIndex].Equals(valType) && valType != TypeNil {
			c.errorf(arg.Pos(), "argument %d to %s has type %s, expected %s", positionalIndex+1, info.Name, valType, info.Params[positionalIndex])
		}
		resolved[positionalIndex] = arg
		assigned[positionalIndex] = true
		positionalIndex++
	}

	providedCount := 0
	for _, ok := range assigned {
		if ok {
			providedCount++
		}
	}
	if providedCount < info.NumRequired || providedCount > len(info.Params) {
		if info.NumRequired == len(info.Params) {
			c.errorf(e.Pos(), "%s expects %d arguments, got %d", info.Name, len(info.Params), providedCount)
		} else {
			c.errorf(e.Pos(), "%s expects %d to %d arguments, got %d", info.Name, info.NumRequired, len(info.Params), providedCount)
		}
		return
	}
	for i := 0; i < info.NumRequired; i++ {
		if !assigned[i] {
			c.errorf(e.Pos(), "%s missing required argument '%s'", info.Name, info.ParamNames[i])
		}
	}
}

func (c *Checker) checkFieldExpr(e *ast.FieldExpr) ZType {
	// Check for enum variant access: EnumName.Variant
	if ident, ok := e.Object.(*ast.IdentExpr); ok {
		if ei, ok := c.enums[ident.Name]; ok {
			if _, ok := ei.Variants[e.Field]; ok {
				return &EnumType{Name: ei.Name, Variants: ei.Variants}
			}
			c.errorf(e.Pos(), "enum %s has no variant '%s'", ei.Name, e.Field)
			return TypeVoid
		}
	}
	if ident, ok := e.Object.(*ast.IdentExpr); ok && c.isModuleName(ident.Name) {
		return TypeVoid
	}
	objType := c.checkNode(e.Object)
	if objType.Equals(TypeVoid) {
		return TypeVoid
	}
	if st, ok := objType.(*ObjType); ok {
		if ft, ok := st.Fields[e.Field]; ok {
			return ft
		}
		c.errorf(e.Pos(), "obj %s has no field '%s'", st.Name, e.Field)
		return TypeVoid
	}
	if tt, ok := objType.(*TupleType); ok {
		// Try numeric index first (works for both named and positional)
		idx, err := strconv.Atoi(e.Field)
		if err == nil {
			if idx < 0 || idx >= len(tt.Elems) {
				c.errorf(e.Pos(), "tuple index %d out of range (len=%d)", idx, len(tt.Elems))
				return TypeVoid
			}
			e.TupleAccess = true
			e.TupleIndex = idx
			e.TupleElemGoType = goTypeName(tt.Elems[idx])
			return tt.Elems[idx]
		}
		// Try named field access
		if tt.Names != nil {
			for i, name := range tt.Names {
				if name == e.Field {
					e.TupleAccess = true
					e.TupleIndex = i
					e.TupleElemGoType = goTypeName(tt.Elems[i])
					return tt.Elems[i]
				}
			}
			c.errorf(e.Pos(), "named tuple has no field '%s'", e.Field)
			return TypeVoid
		}
		c.errorf(e.Pos(), "tuple field must be numeric index, got '%s'", e.Field)
		return TypeVoid
	}
	c.errorf(e.Pos(), "cannot access field '%s' on %s", e.Field, objType)
	return TypeVoid
}

func (c *Checker) checkIndexExpr(e *ast.IndexExpr) ZType {
	objType := c.checkNode(e.Object)
	idxType := c.checkNode(e.Index)

	switch t := objType.(type) {
	case *SliceType:
		if !IsInteger(idxType) {
			c.errorf(e.Pos(), "index must be integer, got %s", idxType)
		}
		return t.Elem
	case *HashmapType:
		if !t.Key.Equals(idxType) {
			c.errorf(e.Pos(), "hashmap key type mismatch: expected %s, got %s", t.Key, idxType)
		}
		if _, ok := t.Key.(*ObjType); ok {
			e.HashmapObjKey = true
		}
		// Propagate default value from symbol to IndexExpr
		if ident, ok := e.Object.(*ast.IdentExpr); ok {
			if sym := c.scope.Lookup(ident.Name); sym != nil && sym.DefaultExpr != nil {
				e.HashmapDefaultVal = sym.DefaultExpr
			}
		}
		return t.Value
	default:
		if objType.Equals(TypeStr) {
			if !IsInteger(idxType) {
				c.errorf(e.Pos(), "index must be integer, got %s", idxType)
			}
			e.StrIndex = true
			return TypeStr
		}
		c.errorf(e.Pos(), "cannot index %s", objType)
		return TypeVoid
	}
}

func (c *Checker) checkSliceExpr(e *ast.SliceExpr) ZType {
	objType := c.checkNode(e.Object)
	if e.Low != nil {
		lowType := c.checkNode(e.Low)
		if !IsInteger(lowType) {
			c.errorf(e.Pos(), "slice index must be integer, got %s", lowType)
		}
	}
	if e.High != nil {
		highType := c.checkNode(e.High)
		if !IsInteger(highType) {
			c.errorf(e.Pos(), "slice index must be integer, got %s", highType)
		}
	}
	if objType.Equals(TypeStr) {
		e.StrSlice = true
		return TypeStr
	}
	if _, ok := objType.(*SliceType); ok {
		return objType
	}
	c.errorf(e.Pos(), "cannot slice %s", objType)
	return TypeVoid
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
	// Could be an enum name (used for EnumName.Variant access)
	if _, ok := c.enums[e.Name]; ok {
		return TypeVoid // enum names aren't values by themselves
	}
	if c.isModuleName(e.Name) {
		return TypeVoid
	}
	c.errorf(e.Pos(), "undefined identifier: %s", e.Name)
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
	sliceType := &SliceType{Elem: firstType}
	e.GoType = goTypeName(sliceType)
	return sliceType
}

func (c *Checker) checkTupleLit(e *ast.TupleLitExpr) ZType {
	elems := make([]ZType, 0, len(e.Elements))
	for _, elem := range e.Elements {
		elems = append(elems, c.checkNode(elem))
	}
	if e.Names != nil {
		// Check for duplicate field names
		seen := make(map[string]bool)
		for _, name := range e.Names {
			if seen[name] {
				c.errorf(e.Pos(), "duplicate field name '%s' in named tuple", name)
			}
			seen[name] = true
		}
		return &TupleType{Elems: elems, Names: e.Names}
	}
	return &TupleType{Elems: elems}
}

func (c *Checker) checkHashmapConstructor(e *ast.CallExpr) ZType {
	// Separate positional args from named args (default=expr)
	var posArgs []ast.Node
	var defaultExpr ast.Node
	for _, arg := range e.Args {
		if na, ok := arg.(*ast.NamedArgExpr); ok {
			if na.Name == "default" {
				defaultExpr = na.Value
			} else {
				c.errorf(na.Pos(), "hashmap() unknown named argument '%s'", na.Name)
			}
		} else {
			posArgs = append(posArgs, arg)
		}
	}

	if len(posArgs) != 2 {
		c.errorf(e.Pos(), "hashmap() expects exactly 2 type arguments, got %d", len(posArgs))
		return &HashmapType{Key: TypeVoid, Value: TypeVoid}
	}

	keyType := c.resolveTypeRefArg(posArgs[0])
	valType := c.resolveTypeRefArg(posArgs[1])

	if keyType == nil {
		keyType = TypeVoid
	}
	if valType == nil {
		valType = TypeVoid
	}
	if !isHashmapKeyType(keyType) && !keyType.Equals(TypeVoid) {
		c.errorf(posArgs[0].Pos(), "invalid hashmap key type: %s", keyType)
	}

	// Validate default expression type
	if defaultExpr != nil {
		defType := c.checkNode(defaultExpr)
		if !valType.Equals(defType) && !valType.Equals(TypeVoid) {
			c.errorf(defaultExpr.Pos(), "hashmap default type mismatch: expected %s, got %s", valType, defType)
		}
		e.HashmapDefaultVal = defaultExpr
	}

	e.HashmapCtor = true
	if objKey, ok := keyType.(*ObjType); ok {
		e.HashmapObjKey = true
		e.HashmapKeyGoType = objKey.Name
	} else {
		e.HashmapKeyGoType = goTypeName(keyType)
	}
	e.HashmapValGoType = goTypeName(valType)
	return &HashmapType{Key: keyType, Value: valType}
}

func (c *Checker) checkSetConstructor(e *ast.CallExpr) ZType {
	if len(e.Args) != 1 {
		c.errorf(e.Pos(), "set() expects exactly 1 type argument, got %d", len(e.Args))
		return &SetType{Elem: TypeVoid}
	}

	elemType := c.resolveTypeRefArg(e.Args[0])
	if elemType == nil {
		elemType = TypeVoid
	}

	e.SetCtor = true

	// Handle set(tuple(...)) with struct representation
	if tt, ok := elemType.(*TupleType); ok {
		if !isComparableType(tt) {
			c.errorf(e.Pos(), "set(tuple(...)) requires all tuple elements to be comparable types")
			e.SetElemGoType = "interface{}"
			return &SetType{Elem: elemType}
		}
		structName := tupleStructName(tt)
		e.SetElemGoType = structName
		e.SetTupleStruct = structName
		e.SetTupleFieldTypes = tupleFieldGoTypes(tt)
		return &SetType{Elem: elemType}
	}

	e.SetElemGoType = goTypeName(elemType)
	return &SetType{Elem: elemType}
}

func (c *Checker) resolveTypeRefArg(arg ast.Node) ZType {
	// Simple identifier: str, int, Point, etc.
	if ident, ok := arg.(*ast.IdentExpr); ok {
		if bt := LookupBuiltinType(ident.Name); bt != nil {
			return bt
		}
		if st, ok := c.objs[ident.Name]; ok {
			return &ObjType{Name: st.Name, Fields: st.Fields}
		}
		if alias, ok := c.typeAliases[ident.Name]; ok {
			return c.resolveTypeExpr(alias)
		}
		c.errorf(arg.Pos(), "unknown type: %s", ident.Name)
		return nil
	}
	// Nested generic: array(T), hashmap(K, V), etc. parsed as CallExpr
	if call, ok := arg.(*ast.CallExpr); ok {
		if callee, ok := call.Callee.(*ast.IdentExpr); ok {
			switch callee.Name {
			case "array":
				if len(call.Args) != 1 {
					c.errorf(arg.Pos(), "array() type expects exactly 1 argument, got %d", len(call.Args))
					return nil
				}
				elem := c.resolveTypeRefArg(call.Args[0])
				if elem == nil {
					return nil
				}
				return &SliceType{Elem: elem}
			case "hashmap":
				// Filter out named args (like default=)
				var posArgs []ast.Node
				for _, a := range call.Args {
					if _, ok := a.(*ast.NamedArgExpr); !ok {
						posArgs = append(posArgs, a)
					}
				}
				if len(posArgs) != 2 {
					c.errorf(arg.Pos(), "hashmap() type expects exactly 2 type arguments, got %d", len(posArgs))
					return nil
				}
				keyType := c.resolveTypeRefArg(posArgs[0])
				valType := c.resolveTypeRefArg(posArgs[1])
				if keyType == nil || valType == nil {
					return nil
				}
				return &HashmapType{Key: keyType, Value: valType}
			case "tuple":
				elems := make([]ZType, 0, len(call.Args))
				for _, a := range call.Args {
					et := c.resolveTypeRefArg(a)
					if et == nil {
						return nil
					}
					elems = append(elems, et)
				}
				return &TupleType{Elems: elems}
			case "set":
				if len(call.Args) != 1 {
					c.errorf(arg.Pos(), "set() type expects exactly 1 argument, got %d", len(call.Args))
					return nil
				}
				elem := c.resolveTypeRefArg(call.Args[0])
				if elem == nil {
					return nil
				}
				return &SetType{Elem: elem}
			}
		}
	}
	// TupleLitExpr: tuple(T1, T2) parsed as a tuple literal in expression context
	if tup, ok := arg.(*ast.TupleLitExpr); ok {
		elems := make([]ZType, 0, len(tup.Elements))
		for _, elem := range tup.Elements {
			et := c.resolveTypeRefArg(elem)
			if et == nil {
				return nil
			}
			elems = append(elems, et)
		}
		return &TupleType{Elems: elems, Names: nil}
	}
	c.errorf(arg.Pos(), "type arguments must be type names, got %T", arg)
	return nil
}

func isHashmapKeyType(t ZType) bool {
	switch t.(type) {
	case *BuiltinType, *ObjType:
		return true
	default:
		return false
	}
}

func (c *Checker) isTypedEmptySliceAssignment(value ast.Node, targetType ZType, valueType ZType) bool {
	lit, ok := value.(*ast.ArrayLitExpr)
	if !ok || len(lit.Elements) != 0 {
		return false
	}
	targetSlice, ok := targetType.(*SliceType)
	if !ok {
		return false
	}
	valueSlice, ok := valueType.(*SliceType)
	if !ok {
		return false
	}
	if !valueSlice.Elem.Equals(TypeVoid) || targetSlice.Elem.Equals(TypeVoid) {
		return false
	}
	lit.GoType = goTypeName(targetType)
	return true
}

func (c *Checker) checkObjConstructor(e *ast.CallExpr, name string) ZType {
	info := c.objs[name]

	// All args must be named
	seen := make(map[string]bool)
	for _, arg := range e.Args {
		named, ok := arg.(*ast.NamedArgExpr)
		if !ok {
			c.errorf(arg.Pos(), "%s constructor requires named arguments (field=value)", name)
			continue
		}
		if seen[named.Name] {
			c.errorf(named.Pos(), "duplicate field '%s' in %s constructor", named.Name, name)
			continue
		}
		seen[named.Name] = true
		expected, ok := info.Fields[named.Name]
		if !ok {
			c.errorf(named.Pos(), "obj %s has no field '%s'", name, named.Name)
			continue
		}
		actual := c.checkNode(named.Value)
		if !expected.Equals(actual) {
			if !c.isTypedEmptySliceAssignment(named.Value, expected, actual) {
				c.errorf(named.Pos(), "field '%s': expected %s, got %s", named.Name, expected, actual)
			}
		}
	}

	// Check that all fields without defaults are provided
	for _, fieldName := range info.Order {
		if !seen[fieldName] && !info.Defaults[fieldName] {
			c.errorf(e.Pos(), "%s constructor missing required field '%s'", name, fieldName)
		}
	}

	return &ObjType{Name: info.Name, Fields: info.Fields}
}

func (c *Checker) checkFlagCall(e *ast.CallExpr) ZType {
	if c.pendingFlagName == "" {
		c.errorf(e.Pos(), "flag() must be assigned to a variable (let x = flag(default=...))")
		return TypeVoid
	}

	// Find the "default" named arg
	var defaultNode ast.Node
	for _, arg := range e.Args {
		if named, ok := arg.(*ast.NamedArgExpr); ok {
			if named.Name == "default" {
				defaultNode = named.Value
			} else {
				c.errorf(named.Pos(), "flag() has no parameter named '%s'", named.Name)
			}
		} else {
			c.errorf(arg.Pos(), "flag() requires named argument: default=<value>")
		}
	}
	if defaultNode == nil {
		c.errorf(e.Pos(), "flag() requires a 'default' argument")
		return TypeVoid
	}

	valType := c.checkNode(defaultNode)
	var goType string
	switch {
	case valType.Equals(TypeBool):
		goType = "bool"
	case valType.Equals(TypeInt):
		goType = "int"
	case valType.Equals(TypeStr):
		goType = "string"
	default:
		c.errorf(e.Pos(), "flag() default must be bool, int, or str, got %s", valType)
		return TypeVoid
	}

	e.FlagName = c.pendingFlagName
	e.FlagGoType = goType
	return valType
}

func (c *Checker) checkClosureExpr(e *ast.ClosureExpr, expectedParamType ZType) ZType {
	c.pushScope()
	defer c.popScope()

	var paramTypes []ZType
	var goParams []string

	for _, p := range e.Params {
		var pType ZType
		if p.Type != nil {
			pType = c.resolveTypeExpr(p.Type)
		} else if expectedParamType != nil {
			pType = expectedParamType
		} else {
			c.errorf(e.Pos(), "closure parameter '%s' requires a type annotation (no context to infer from)", p.Name)
			pType = TypeVoid
		}
		paramTypes = append(paramTypes, pType)
		c.scope.Define(&Symbol{Name: p.Name, Type: pType})
		goParams = append(goParams, p.Name+" "+goTypeName(pType))
	}

	var retType ZType
	if _, isBlock := e.Body.(*ast.Block); isBlock {
		// Block body: use return type annotation or infer void
		if e.ReturnType != nil {
			retType = c.resolveTypeExpr(e.ReturnType)
		} else {
			retType = TypeVoid
		}
		// Set up a temporary FuncInfo so return statements can be checked
		prev := c.currentFunc
		c.currentFunc = &FuncInfo{Name: "<closure>", Return: retType}
		c.checkNode(e.Body)
		c.currentFunc = prev
	} else {
		// Single expression: implicit return
		retType = c.checkNode(e.Body)
		if e.ReturnType != nil {
			declared := c.resolveTypeExpr(e.ReturnType)
			if !declared.Equals(retType) {
				c.errorf(e.Pos(), "closure return type mismatch: declared %s, body returns %s", declared, retType)
			}
			retType = declared
		}
	}

	e.GoParams = strings.Join(goParams, ", ")
	e.GoReturn = goTypeName(retType)

	return &FuncType{Params: paramTypes, Returns: retType}
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

func (c *Checker) checkMatchExpr(e *ast.MatchExpr) ZType {
	c.checkNode(e.Subject)
	if len(e.Arms) == 0 {
		c.errorf(e.Pos(), "match expression must have at least one arm")
		return TypeVoid
	}
	var firstType ZType
	for i, arm := range e.Arms {
		c.checkNode(arm.Pattern)
		armType := c.checkNode(arm.Value)
		if i == 0 {
			firstType = armType
		} else if !armType.Equals(firstType) {
			c.errorf(arm.Value.Pos(), "match expression arms must have same type: first arm is %s, this arm is %s", firstType, armType)
		}
	}
	e.GoType = goTypeName(firstType)
	return firstType
}

func (c *Checker) checkAssertCall(e *ast.CallExpr) ZType {
	args := c.getPositionalArgsForBuiltin(e, "assert")
	if len(args) != 1 {
		c.errorf(e.Pos(), "assert() takes exactly 1 argument, got %d", len(args))
		return TypeVoid
	}
	argType := c.checkNode(args[0])
	if !argType.Equals(TypeBool) {
		c.errorf(args[0].Pos(), "assert() argument must be bool, got %s", argType)
	}
	e.ResolvedFunc = "assert"
	return TypeVoid
}

func (c *Checker) checkAssertEqCall(e *ast.CallExpr) ZType {
	args := c.getPositionalArgsForBuiltin(e, "assert_eq")
	if len(args) != 2 {
		c.errorf(e.Pos(), "assert_eq() takes exactly 2 arguments, got %d", len(args))
		return TypeVoid
	}
	gotType := c.checkNode(args[0])
	expectedType := c.checkNode(args[1])
	if !gotType.Equals(expectedType) {
		c.errorf(args[1].Pos(), "assert_eq() arguments must be the same type, got %s and %s", gotType, expectedType)
	}
	e.ResolvedFunc = "assert_eq"
	return TypeVoid
}

// isComparableType checks if a Zenth type maps to a comparable Go type (usable as map key).
func isComparableType(t ZType) bool {
	switch ty := t.(type) {
	case *BuiltinType:
		return true // all builtin types are comparable in Go
	case *ObjType:
		return true // obj types are structs, comparable if all fields are comparable
	case *TupleType:
		// Tuples become []interface{} normally (not comparable), but if all
		// elements are comparable builtins we can generate a struct instead.
		for _, elem := range ty.Elems {
			if !isComparableType(elem) {
				return false
			}
		}
		return true
	default:
		return false // slices, maps, sets, etc. are not comparable
	}
}

// tupleStructName generates a Go struct name for a comparable tuple used as set/map key.
func tupleStructName(tt *TupleType) string {
	parts := make([]string, len(tt.Elems))
	for i, elem := range tt.Elems {
		parts[i] = goTypeName(elem)
	}
	return "ZenthTuple_" + strings.Join(parts, "_")
}

// tupleFieldGoTypes returns the Go type names for each tuple element.
func tupleFieldGoTypes(tt *TupleType) []string {
	types := make([]string, len(tt.Elems))
	for i, elem := range tt.Elems {
		types[i] = goTypeName(elem)
	}
	return types
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
	case *HashmapType:
		keyName := goTypeName(ty.Key)
		// Obj keys use struct value types (not pointers) in Go maps
		if _, ok := ty.Key.(*ObjType); ok {
			keyName = strings.TrimPrefix(keyName, "*")
		}
		return "map[" + keyName + "]" + goTypeName(ty.Value)
	case *SetType:
		if tt, ok := ty.Elem.(*TupleType); ok {
			return "map[" + tupleStructName(tt) + "]struct{}"
		}
		return "map[" + goTypeName(ty.Elem) + "]struct{}"
	case *TupleType:
		return "[]interface{}"
	case *ObjType:
		return "*" + ty.Name
	case *EnumType:
		return ty.Name
	default:
		return "interface{}"
	}
}
