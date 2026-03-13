package token

import (
	"testing"
)

func TestTypeString_Known(t *testing.T) {
	tests := []struct {
		typ  Type
		want string
	}{
		{Illegal, "Illegal"},
		{EOF, "EOF"},
		{Ident, "Ident"},
		{IntLit, "IntLit"},
		{FloatLit, "FloatLit"},
		{StringLit, "StringLit"},
		{InterpStringLit, "InterpStringLit"},
		{Plus, "+"},
		{Minus, "-"},
		{Star, "*"},
		{Slash, "/"},
		{Percent, "%"},
		{Assign, "="},
		{Walrus, ":="},
		{Arrow, "->"},
		{FatArrow, "=>"},
		{PlusPlus, "++"},
		{MinusMinus, "--"},
		{Eq, "=="},
		{Neq, "!="},
		{Lt, "<"},
		{Gt, ">"},
		{Lte, "<="},
		{Gte, ">="},
		{And, "&&"},
		{Or, "||"},
		{Not, "!"},
		{LParen, "("},
		{RParen, ")"},
		{LBrace, "{"},
		{RBrace, "}"},
		{LBracket, "["},
		{RBracket, "]"},
		{Comma, ","},
		{Dot, "."},
		{Colon, ":"},
		{Semicolon, ";"},
		{Fn, "fn"},
		{Let, "let"},
		{Var, "var"},
		{Const, "const"},
		{Return, "return"},
		{If, "if"},
		{Else, "else"},
		{For, "for"},
		{Loop, "loop"},
		{In, "in"},
		{Match, "match"},
		{Obj, "obj"},
		{Enum, "enum"},
		{Interface, "interface"},
		{Import, "import"},
		{From, "from"},
		{As, "as"},
		{Nil, "nil"},
		{True, "true"},
		{False, "false"},
		{Break, "break"},
		{Continue, "continue"},
		{TypeKw, "type"},
		{PlusAssign, "+="},
		{MinusAssign, "-="},
		{StarAssign, "*="},
		{SlashAssign, "/="},
	}

	for _, tt := range tests {
		got := tt.typ.String()
		if got != tt.want {
			t.Errorf("Type(%d).String() = %q, want %q", int(tt.typ), got, tt.want)
		}
	}
}

func TestTypeString_Unknown(t *testing.T) {
	unknown := Type(9999)
	got := unknown.String()
	want := "Token(9999)"
	if got != want {
		t.Errorf("Type(9999).String() = %q, want %q", got, want)
	}
}

func TestLookupIdent_Keywords(t *testing.T) {
	tests := []struct {
		ident string
		want  Type
	}{
		{"fn", Fn},
		{"let", Let},
		{"var", Var},
		{"const", Const},
		{"return", Return},
		{"if", If},
		{"else", Else},
		{"for", For},
		{"loop", Loop},
		{"in", In},
		{"match", Match},
		{"obj", Obj},
		{"enum", Enum},
		{"interface", Interface},
		{"import", Import},
		{"from", From},
		{"as", As},
		{"nil", Nil},
		{"true", True},
		{"false", False},
		{"break", Break},
		{"continue", Continue},
		{"type", TypeKw},
	}

	for _, tt := range tests {
		got := LookupIdent(tt.ident)
		if got != tt.want {
			t.Errorf("LookupIdent(%q) = %v, want %v", tt.ident, got, tt.want)
		}
	}
}

func TestLookupIdent_NonKeywords(t *testing.T) {
	nonKeywords := []string{
		"foo", "bar", "myVar", "x", "hello_world",
		"Fn", "LET", "Return", // case-sensitive: not keywords
		"fnn", "iff", "forr", // near-misses
		"", // empty string
	}

	for _, ident := range nonKeywords {
		got := LookupIdent(ident)
		if got != Ident {
			t.Errorf("LookupIdent(%q) = %v, want Ident", ident, got)
		}
	}
}

func TestPosString(t *testing.T) {
	tests := []struct {
		pos  Pos
		want string
	}{
		{Pos{"main.zn", 1, 1}, "main.zn:1:1"},
		{Pos{"src/foo.zn", 42, 10}, "src/foo.zn:42:10"},
		{Pos{"", 0, 0}, ":0:0"},
	}

	for _, tt := range tests {
		got := tt.pos.String()
		if got != tt.want {
			t.Errorf("Pos%+v.String() = %q, want %q", tt.pos, got, tt.want)
		}
	}
}

func TestTokenStruct(t *testing.T) {
	tok := Token{
		Type:    IntLit,
		Literal: "42",
		Pos:     Pos{File: "test.zn", Line: 5, Column: 3},
	}

	if tok.Type != IntLit {
		t.Errorf("tok.Type = %v, want IntLit", tok.Type)
	}
	if tok.Literal != "42" {
		t.Errorf("tok.Literal = %q, want %q", tok.Literal, "42")
	}
	if tok.Pos.File != "test.zn" {
		t.Errorf("tok.Pos.File = %q, want %q", tok.Pos.File, "test.zn")
	}
	if tok.Pos.Line != 5 {
		t.Errorf("tok.Pos.Line = %d, want 5", tok.Pos.Line)
	}
	if tok.Pos.Column != 3 {
		t.Errorf("tok.Pos.Column = %d, want 3", tok.Pos.Column)
	}
}

func TestTokenTypeConstants_Unique(t *testing.T) {
	// Verify all token types have distinct values
	seen := make(map[Type]string)
	allTypes := []struct {
		typ  Type
		name string
	}{
		{Illegal, "Illegal"}, {EOF, "EOF"},
		{Ident, "Ident"}, {IntLit, "IntLit"}, {FloatLit, "FloatLit"},
		{StringLit, "StringLit"}, {InterpStringLit, "InterpStringLit"},
		{Plus, "Plus"}, {Minus, "Minus"}, {Star, "Star"}, {Slash, "Slash"},
		{Percent, "Percent"}, {Assign, "Assign"}, {Walrus, "Walrus"},
		{Arrow, "Arrow"}, {FatArrow, "FatArrow"},
		{PlusPlus, "PlusPlus"}, {MinusMinus, "MinusMinus"},
		{Eq, "Eq"}, {Neq, "Neq"}, {Lt, "Lt"}, {Gt, "Gt"}, {Lte, "Lte"}, {Gte, "Gte"},
		{And, "And"}, {Or, "Or"}, {Not, "Not"},
		{LParen, "LParen"}, {RParen, "RParen"},
		{LBrace, "LBrace"}, {RBrace, "RBrace"},
		{LBracket, "LBracket"}, {RBracket, "RBracket"},
		{Comma, "Comma"}, {Dot, "Dot"}, {Colon, "Colon"}, {Semicolon, "Semicolon"},
		{Fn, "Fn"}, {Let, "Let"}, {Var, "Var"}, {Const, "Const"},
		{Return, "Return"}, {If, "If"}, {Else, "Else"}, {For, "For"},
		{Loop, "Loop"}, {In, "In"}, {Match, "Match"}, {Obj, "Obj"},
		{Enum, "Enum"}, {Interface, "Interface"}, {Import, "Import"},
		{From, "From"}, {As, "As"}, {Nil, "Nil"}, {True, "True"}, {False, "False"},
		{Break, "Break"}, {Continue, "Continue"}, {TypeKw, "TypeKw"},
		{PlusAssign, "PlusAssign"}, {MinusAssign, "MinusAssign"},
		{StarAssign, "StarAssign"}, {SlashAssign, "SlashAssign"},
	}

	for _, tt := range allTypes {
		if prev, ok := seen[tt.typ]; ok {
			t.Errorf("duplicate token type value %d: %s and %s", int(tt.typ), prev, tt.name)
		}
		seen[tt.typ] = tt.name
	}
}
