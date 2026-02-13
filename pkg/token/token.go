package token

import "fmt"

// Type represents a token type.
type Type int

const (
	// Special
	Illegal Type = iota
	EOF

	// Literals
	Ident     // variable names, function names
	IntLit    // 42
	FloatLit  // 3.14
	StringLit       // "hello"
	InterpStringLit // "hello {name}" (contains interpolation)

	// Operators
	Plus     // +
	Minus    // -
	Star     // *
	Slash    // /
	Percent  // %
	Assign   // =
	Walrus   // :=
	Arrow    // ->
	FatArrow // =>
	PlusPlus // ++
	MinusMinus // --

	// Comparison
	Eq  // ==
	Neq // !=
	Lt  // <
	Gt  // >
	Lte // <=
	Gte // >=

	// Logical
	And // &&
	Or  // ||
	Not // !

	// Delimiters
	LParen   // (
	RParen   // )
	LBrace   // {
	RBrace   // }
	LBracket // [
	RBracket // ]

	// Punctuation
	Comma     // ,
	Dot       // .
	Colon     // :
	Semicolon // ;

	// Keywords
	Fn
	Let
	Var
	Const
	Return
	If
	Else
	For
	In
	Match
	Struct
	Enum
	Interface
	Import
	From
	As
	Nil
	True
	False
	Break
	Continue
	PlusAssign  // +=
	MinusAssign // -=
	StarAssign  // *=
	SlashAssign // /=
)

var typeNames = map[Type]string{
	Illegal:    "Illegal",
	EOF:        "EOF",
	Ident:      "Ident",
	IntLit:     "IntLit",
	FloatLit:   "FloatLit",
	StringLit:       "StringLit",
	InterpStringLit: "InterpStringLit",
	Plus:       "+",
	Minus:      "-",
	Star:       "*",
	Slash:      "/",
	Percent:    "%",
	Assign:     "=",
	Walrus:     ":=",
	Arrow:      "->",
	FatArrow:   "=>",
	PlusPlus:   "++",
	MinusMinus: "--",
	Eq:         "==",
	Neq:        "!=",
	Lt:         "<",
	Gt:         ">",
	Lte:        "<=",
	Gte:        ">=",
	And:        "&&",
	Or:         "||",
	Not:        "!",
	LParen:     "(",
	RParen:     ")",
	LBrace:     "{",
	RBrace:     "}",
	LBracket:   "[",
	RBracket:   "]",
	Comma:      ",",
	Dot:        ".",
	Colon:      ":",
	Semicolon:  ";",
	Fn:         "fn",
	Let:        "let",
	Var:        "var",
	Const:      "const",
	Return:     "return",
	If:         "if",
	Else:       "else",
	For:        "for",
	In:         "in",
	Match:      "match",
	Struct:     "struct",
	Enum:       "enum",
	Interface:  "interface",
	Import:     "import",
	From:       "from",
	As:         "as",
	Nil:        "nil",
	True:       "true",
	False:      "false",
	Break:      "break",
	Continue:   "continue",
	PlusAssign:  "+=",
	MinusAssign: "-=",
	StarAssign:  "*=",
	SlashAssign: "/=",
}

func (t Type) String() string {
	if name, ok := typeNames[t]; ok {
		return name
	}
	return fmt.Sprintf("Token(%d)", int(t))
}

var keywords = map[string]Type{
	"fn":        Fn,
	"let":       Let,
	"var":       Var,
	"const":     Const,
	"return":    Return,
	"if":        If,
	"else":      Else,
	"for":       For,
	"in":        In,
	"match":     Match,
	"struct":    Struct,
	"enum":      Enum,
	"interface": Interface,
	"import":    Import,
	"from":      From,
	"as":        As,
	"nil":       Nil,
	"true":      True,
	"false":     False,
	"break":     Break,
	"continue":  Continue,
}

// LookupIdent returns the keyword token type for ident, or Ident if not a keyword.
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}

// Pos represents a source position.
type Pos struct {
	File   string
	Line   int
	Column int
}

func (p Pos) String() string {
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Column)
}

// Token represents a single token with its type, literal value, and position.
type Token struct {
	Type    Type
	Literal string
	Pos     Pos
}
