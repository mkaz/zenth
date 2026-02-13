package lexer

import (
	"testing"

	"github.com/mkaz/zenth/pkg/token"
)

func TestHelloWorld(t *testing.T) {
	src := `fn main() { println("Hello"); }`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	expected := []token.Type{
		token.Fn, token.Ident, token.LParen, token.RParen,
		token.LBrace, token.Ident, token.LParen, token.StringLit,
		token.RParen, token.Semicolon, token.RBrace, token.EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestOperators(t *testing.T) {
	src := `:= -> => ++ -- += -= == != <= >=`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	expected := []token.Type{
		token.Walrus, token.Arrow, token.FatArrow,
		token.PlusPlus, token.MinusMinus,
		token.PlusAssign, token.MinusAssign,
		token.Eq, token.Neq, token.Lte, token.Gte,
		token.EOF,
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("token %d: expected %s, got %s", i, expected[i], tok.Type)
		}
	}
}

func TestKeywords(t *testing.T) {
	src := `fn let var const return if else for in match obj`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	expected := []token.Type{
		token.Fn, token.Let, token.Var, token.Const, token.Return,
		token.If, token.Else, token.For, token.In, token.Match, token.Obj,
		token.EOF,
	}

	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("token %d: expected %s, got %s (%q)", i, expected[i], tok.Type, tok.Literal)
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	src := `42 3.14 1_000_000`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[0].Type != token.IntLit || tokens[0].Literal != "42" {
		t.Errorf("expected IntLit 42, got %s %q", tokens[0].Type, tokens[0].Literal)
	}
	if tokens[1].Type != token.FloatLit || tokens[1].Literal != "3.14" {
		t.Errorf("expected FloatLit 3.14, got %s %q", tokens[1].Type, tokens[1].Literal)
	}
	if tokens[2].Type != token.IntLit || tokens[2].Literal != "1000000" {
		t.Errorf("expected IntLit 1000000, got %s %q", tokens[2].Type, tokens[2].Literal)
	}
}

func TestComments(t *testing.T) {
	src := `// this is a comment
# this is also a comment
fn main() { /* block comment */ }`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	// All comments should be stripped
	if tokens[0].Type != token.Fn {
		t.Errorf("expected Fn after comments, got %s", tokens[0].Type)
	}
}

func TestStringEscapes(t *testing.T) {
	src := `"hello\nworld" "tab\there"`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[0].Literal != "hello\nworld" {
		t.Errorf("expected newline escape, got %q", tokens[0].Literal)
	}
	if tokens[1].Literal != "tab\there" {
		t.Errorf("expected tab escape, got %q", tokens[1].Literal)
	}
}

func TestInterpString(t *testing.T) {
	src := `"hello {name}"`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Type != token.InterpStringLit {
		t.Errorf("expected InterpStringLit, got %s", tokens[0].Type)
	}
	if tokens[0].Literal != "hello {name}" {
		t.Errorf("expected literal 'hello {name}', got %q", tokens[0].Literal)
	}
}

func TestPlainStringNoInterp(t *testing.T) {
	src := `"hello world"`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Type != token.StringLit {
		t.Errorf("expected StringLit, got %s", tokens[0].Type)
	}
}

func TestRawString(t *testing.T) {
	src := `'hello {name}'`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Type != token.StringLit {
		t.Errorf("expected StringLit for raw string, got %s", tokens[0].Type)
	}
	if tokens[0].Literal != "hello {name}" {
		t.Errorf("expected literal 'hello {name}', got %q", tokens[0].Literal)
	}
}

func TestEscapedBrace(t *testing.T) {
	src := `"hello \{world\}"`
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	// Escaped braces should NOT trigger interpolation
	if tokens[0].Type != token.StringLit {
		t.Errorf("expected StringLit for escaped braces, got %s", tokens[0].Type)
	}
	if tokens[0].Literal != "hello {world}" {
		t.Errorf("expected literal 'hello {world}', got %q", tokens[0].Literal)
	}
}

func TestPositionTracking(t *testing.T) {
	src := "fn main() {\n    println(\"hi\");\n}"
	l := New("test.zn", src)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	// fn should be at line 1, col 1
	if tokens[0].Pos.Line != 1 || tokens[0].Pos.Column != 1 {
		t.Errorf("expected fn at 1:1, got %d:%d", tokens[0].Pos.Line, tokens[0].Pos.Column)
	}

	// println should be at line 2
	for _, tok := range tokens {
		if tok.Literal == "println" {
			if tok.Pos.Line != 2 {
				t.Errorf("expected println at line 2, got line %d", tok.Pos.Line)
			}
			break
		}
	}
}
