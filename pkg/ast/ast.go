package ast

import "github.com/mkaz/zenth/pkg/token"

// Node is the interface all AST nodes implement.
type Node interface {
	Pos() token.Pos
	nodeMarker()
}

// Program is the root of the AST.
type Program struct {
	Stmts []Node
}

func (p *Program) Pos() token.Pos {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].Pos()
	}
	return token.Pos{}
}
func (p *Program) nodeMarker() {}

// ---------- Declarations ----------

// FnDecl represents a function declaration.
type FnDecl struct {
	TokenPos   token.Pos
	Name       string
	Params     []Param
	ReturnType *TypeExpr // nil means no return
	Body       *Block
	OwnerObj   string // non-empty for methods defined inside an obj
}

func (f *FnDecl) Pos() token.Pos { return f.TokenPos }
func (f *FnDecl) nodeMarker()    {}

// Param represents a function parameter.
type Param struct {
	Name    string
	Type    *TypeExpr
	Default Node // nil if no default value
}

// ObjDecl represents an obj declaration.
type ObjDecl struct {
	TokenPos token.Pos
	Name     string
	Fields   []Field
	Methods  []*FnDecl
}

func (s *ObjDecl) Pos() token.Pos { return s.TokenPos }
func (s *ObjDecl) nodeMarker()    {}

// Field represents an obj field.
type Field struct {
	Name    string
	Type    *TypeExpr
	Default Node // nil if no default
}

// InterfaceDecl represents an interface declaration.
type InterfaceDecl struct {
	TokenPos token.Pos
	Name     string
	Methods  []MethodSig
}

func (i *InterfaceDecl) Pos() token.Pos { return i.TokenPos }
func (i *InterfaceDecl) nodeMarker()    {}

// MethodSig represents a method signature in an interface.
type MethodSig struct {
	Name       string
	Params     []Param
	ReturnType *TypeExpr
}

// ImportDecl represents an import statement.
type ImportDecl struct {
	TokenPos token.Pos
	Path     string
	Alias    string // empty means use default name
}

func (i *ImportDecl) Pos() token.Pos { return i.TokenPos }
func (i *ImportDecl) nodeMarker()    {}

// ---------- Types ----------

// TypeExpr represents a type expression.
type TypeExpr struct {
	TokenPos token.Pos
	Name     string      // "int", "str", "bool", etc. or a user type name
	Params   []*TypeExpr // for generics like map[K]V or []T
	IsSlice  bool        // []T
	IsMap    bool        // map[K]V
	IsArray  bool        // [N]T
	ArrayLen int         // for fixed arrays
}

func (t *TypeExpr) Pos() token.Pos { return t.TokenPos }
func (t *TypeExpr) nodeMarker()    {}

// ---------- Statements ----------

// Block represents a { ... } block of statements.
type Block struct {
	TokenPos token.Pos
	Stmts    []Node
}

func (b *Block) Pos() token.Pos { return b.TokenPos }
func (b *Block) nodeMarker()    {}

// LetStmt represents: let name [: type] = expr;
type LetStmt struct {
	TokenPos token.Pos
	Name     string
	Type     *TypeExpr // nil if using :=
	Value    Node      // expression
	Infer    bool      // true if :=
}

func (l *LetStmt) Pos() token.Pos { return l.TokenPos }
func (l *LetStmt) nodeMarker()    {}

// VarStmt represents: var name [: type] = expr;
type VarStmt struct {
	TokenPos token.Pos
	Name     string
	Type     *TypeExpr
	Value    Node
	Infer    bool
}

func (v *VarStmt) Pos() token.Pos { return v.TokenPos }
func (v *VarStmt) nodeMarker()    {}

// ConstStmt represents: const name [: type] = expr;
type ConstStmt struct {
	TokenPos token.Pos
	Name     string
	Type     *TypeExpr
	Value    Node
	Infer    bool
}

func (c *ConstStmt) Pos() token.Pos { return c.TokenPos }
func (c *ConstStmt) nodeMarker()    {}

// AssignStmt represents: target = expr; or target += expr; etc.
type AssignStmt struct {
	TokenPos token.Pos
	Target   Node       // IdentExpr, IndexExpr, or FieldExpr
	Op       token.Type // Assign, PlusAssign, MinusAssign, etc.
	Value    Node
}

func (a *AssignStmt) Pos() token.Pos { return a.TokenPos }
func (a *AssignStmt) nodeMarker()    {}

// MultiAssignStmt represents: t1, t2 = v1, v2;
type MultiAssignStmt struct {
	TokenPos token.Pos
	Targets  []Node // lvalues (IdentExpr, IndexExpr, FieldExpr)
	Values   []Node // rvalue expressions
}

func (m *MultiAssignStmt) Pos() token.Pos { return m.TokenPos }
func (m *MultiAssignStmt) nodeMarker()    {}

// ReturnStmt represents: return [expr];
type ReturnStmt struct {
	TokenPos token.Pos
	Value    Node // nil for bare return
}

func (r *ReturnStmt) Pos() token.Pos { return r.TokenPos }
func (r *ReturnStmt) nodeMarker()    {}

// IfStmt represents: if cond { ... } [else { ... }]
type IfStmt struct {
	TokenPos  token.Pos
	Condition Node
	Body      *Block
	Else      Node // *Block or *IfStmt (else if) or nil
}

func (i *IfStmt) Pos() token.Pos { return i.TokenPos }
func (i *IfStmt) nodeMarker()    {}

// IfExpr represents an if-expression (used as a value).
// Each branch is a single expression: if cond { expr } else { expr }
type IfExpr struct {
	TokenPos  token.Pos
	Condition Node
	Then      Node   // single expression (then-branch value)
	Else      Node   // single expression, or another *IfExpr (else-if chain)
	GoType    string // Go type string, set by checker for codegen
}

func (i *IfExpr) Pos() token.Pos { return i.TokenPos }
func (i *IfExpr) nodeMarker()    {}

// ForStmt represents: for [init]; [cond]; [post] { ... }
type ForStmt struct {
	TokenPos  token.Pos
	Init      Node // may be nil
	Condition Node // may be nil (infinite loop)
	Post      Node // may be nil
	Body      *Block
}

func (f *ForStmt) Pos() token.Pos { return f.TokenPos }
func (f *ForStmt) nodeMarker()    {}

// ForInStmt represents: for [index,] value in iterable { ... }
type ForInStmt struct {
	TokenPos token.Pos
	Index    string // empty if no index variable
	Value    string
	Iterable Node
	Body     *Block
	IterStr  bool // set by checker when iterating over a string
}

func (f *ForInStmt) Pos() token.Pos { return f.TokenPos }
func (f *ForInStmt) nodeMarker()    {}

// LoopStmt represents: loop count { ... }
type LoopStmt struct {
	TokenPos token.Pos
	Count    Node
	Body     *Block
}

func (l *LoopStmt) Pos() token.Pos { return l.TokenPos }
func (l *LoopStmt) nodeMarker()    {}

// MatchStmt represents: match expr { arms }
type MatchStmt struct {
	TokenPos token.Pos
	Subject  Node
	Arms     []MatchArm
}

func (m *MatchStmt) Pos() token.Pos { return m.TokenPos }
func (m *MatchStmt) nodeMarker()    {}

// MatchArm represents: pattern => stmt or { block }
type MatchArm struct {
	Pattern Node // expression or _ (IdentExpr with name "_")
	Body    Node // single statement or Block
}

// BreakStmt represents: break;
type BreakStmt struct {
	TokenPos token.Pos
}

func (b *BreakStmt) Pos() token.Pos { return b.TokenPos }
func (b *BreakStmt) nodeMarker()    {}

// ContinueStmt represents: continue;
type ContinueStmt struct {
	TokenPos token.Pos
}

func (c *ContinueStmt) Pos() token.Pos { return c.TokenPos }
func (c *ContinueStmt) nodeMarker()    {}

// ExprStmt wraps an expression used as a statement.
type ExprStmt struct {
	Expr Node
}

func (e *ExprStmt) Pos() token.Pos { return e.Expr.Pos() }
func (e *ExprStmt) nodeMarker()    {}

// IncDecStmt represents i++ or i--
type IncDecStmt struct {
	TokenPos token.Pos
	Operand  Node
	Op       token.Type // PlusPlus or MinusMinus
}

func (i *IncDecStmt) Pos() token.Pos { return i.TokenPos }
func (i *IncDecStmt) nodeMarker()    {}

// ---------- Expressions ----------

// BinaryExpr represents: left op right
type BinaryExpr struct {
	TokenPos      token.Pos
	Left          Node
	Op            token.Type
	Right         Node
	PromoteLeft   string // Go type to cast left operand to (set by checker)
	PromoteRight  string // Go type to cast right operand to (set by checker)
	SliceConcat   bool   // true when + means slice concatenation (set by checker)
	SliceContains bool   // true when "in" means slice containment (set by checker)
}

func (b *BinaryExpr) Pos() token.Pos { return b.TokenPos }
func (b *BinaryExpr) nodeMarker()    {}

// UnaryExpr represents: op operand
type UnaryExpr struct {
	TokenPos token.Pos
	Op       token.Type
	Operand  Node
}

func (u *UnaryExpr) Pos() token.Pos { return u.TokenPos }
func (u *UnaryExpr) nodeMarker()    {}

// CallExpr represents: callee(args)
type CallExpr struct {
	TokenPos        token.Pos
	Callee          Node
	Args            []Node
	SliceMethod     bool   // set by checker for built-in slice methods
	SliceConvTarget string // set by checker for to_int/to_f64/to_str ("int", "float64", "string")
	StringMethod    string // set by checker for built-in string methods (e.g. "split")
	ResolvedFunc    string // set by checker: function key ("name" or "Type.method")
}

func (c *CallExpr) Pos() token.Pos { return c.TokenPos }
func (c *CallExpr) nodeMarker()    {}

// IndexExpr represents: object[index]
type IndexExpr struct {
	TokenPos token.Pos
	Object   Node
	Index    Node
	StrIndex bool // set by checker when indexing a string
}

func (i *IndexExpr) Pos() token.Pos { return i.TokenPos }
func (i *IndexExpr) nodeMarker()    {}

// FieldExpr represents: object.field
type FieldExpr struct {
	TokenPos token.Pos
	Object   Node
	Field    string
}

func (f *FieldExpr) Pos() token.Pos { return f.TokenPos }
func (f *FieldExpr) nodeMarker()    {}

// IdentExpr represents a variable reference.
type IdentExpr struct {
	TokenPos token.Pos
	Name     string
}

func (i *IdentExpr) Pos() token.Pos { return i.TokenPos }
func (i *IdentExpr) nodeMarker()    {}

// IntLitExpr represents an integer literal.
type IntLitExpr struct {
	TokenPos token.Pos
	Value    int64
}

func (i *IntLitExpr) Pos() token.Pos { return i.TokenPos }
func (i *IntLitExpr) nodeMarker()    {}

// FloatLitExpr represents a float literal.
type FloatLitExpr struct {
	TokenPos token.Pos
	Value    float64
}

func (f *FloatLitExpr) Pos() token.Pos { return f.TokenPos }
func (f *FloatLitExpr) nodeMarker()    {}

// StringLitExpr represents a string literal.
type StringLitExpr struct {
	TokenPos token.Pos
	Value    string
}

func (s *StringLitExpr) Pos() token.Pos { return s.TokenPos }
func (s *StringLitExpr) nodeMarker()    {}

// BoolLitExpr represents a boolean literal.
type BoolLitExpr struct {
	TokenPos token.Pos
	Value    bool
}

func (b *BoolLitExpr) Pos() token.Pos { return b.TokenPos }
func (b *BoolLitExpr) nodeMarker()    {}

// NilExpr represents the nil literal.
type NilExpr struct {
	TokenPos token.Pos
}

func (n *NilExpr) Pos() token.Pos { return n.TokenPos }
func (n *NilExpr) nodeMarker()    {}

// ArrayLitExpr represents: [expr, expr, ...]
type ArrayLitExpr struct {
	TokenPos token.Pos
	Elements []Node
}

func (a *ArrayLitExpr) Pos() token.Pos { return a.TokenPos }
func (a *ArrayLitExpr) nodeMarker()    {}

// NamedArgExpr represents a named argument: name=value
type NamedArgExpr struct {
	TokenPos token.Pos
	Name     string
	Value    Node
}

func (n *NamedArgExpr) Pos() token.Pos { return n.TokenPos }
func (n *NamedArgExpr) nodeMarker()    {}

// InterpStringExpr represents a string with interpolated expressions: "hello {name}"
type InterpStringExpr struct {
	TokenPos token.Pos
	Parts    []InterpPart
}

func (s *InterpStringExpr) Pos() token.Pos { return s.TokenPos }
func (s *InterpStringExpr) nodeMarker()    {}

// InterpPart is one segment of an interpolated string.
type InterpPart struct {
	IsExpr bool
	Lit    string // text content (when IsExpr is false)
	Expr   Node   // parsed expression (when IsExpr is true)
}
