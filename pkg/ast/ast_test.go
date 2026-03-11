package ast

import (
	"testing"

	"github.com/mkaz/zenth/pkg/token"
)

func TestProgramPos(t *testing.T) {
	t.Run("empty program returns zero pos", func(t *testing.T) {
		p := &Program{}
		got := p.Pos()
		want := token.Pos{}
		if got != want {
			t.Errorf("empty Program.Pos() = %v, want %v", got, want)
		}
	})

	t.Run("non-empty program returns first statement pos", func(t *testing.T) {
		pos := token.Pos{File: "test.zn", Line: 3, Column: 7}
		p := &Program{
			Stmts: []Node{
				&LetStmt{TokenPos: pos},
				&VarStmt{TokenPos: token.Pos{File: "test.zn", Line: 5, Column: 1}},
			},
		}
		got := p.Pos()
		if got != pos {
			t.Errorf("Program.Pos() = %v, want %v", got, pos)
		}
	})
}

func TestNodePositions(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want token.Pos
	}{
		// Declarations
		{
			name: "FnDecl",
			node: &FnDecl{TokenPos: token.Pos{File: "test.zn", Line: 1, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 1, Column: 1},
		},
		{
			name: "ObjDecl",
			node: &ObjDecl{TokenPos: token.Pos{File: "test.zn", Line: 2, Column: 3}},
			want: token.Pos{File: "test.zn", Line: 2, Column: 3},
		},
		{
			name: "EnumDecl",
			node: &EnumDecl{TokenPos: token.Pos{File: "test.zn", Line: 3, Column: 5}},
			want: token.Pos{File: "test.zn", Line: 3, Column: 5},
		},
		{
			name: "InterfaceDecl",
			node: &InterfaceDecl{TokenPos: token.Pos{File: "test.zn", Line: 4, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 4, Column: 1},
		},
		{
			name: "TypeAliasDecl",
			node: &TypeAliasDecl{TokenPos: token.Pos{File: "test.zn", Line: 5, Column: 2}},
			want: token.Pos{File: "test.zn", Line: 5, Column: 2},
		},
		{
			name: "ImportDecl",
			node: &ImportDecl{TokenPos: token.Pos{File: "test.zn", Line: 6, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 6, Column: 1},
		},

		// Statements
		{
			name: "Block",
			node: &Block{TokenPos: token.Pos{File: "test.zn", Line: 7, Column: 10}},
			want: token.Pos{File: "test.zn", Line: 7, Column: 10},
		},
		{
			name: "LetStmt",
			node: &LetStmt{TokenPos: token.Pos{File: "test.zn", Line: 8, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 8, Column: 1},
		},
		{
			name: "VarStmt",
			node: &VarStmt{TokenPos: token.Pos{File: "test.zn", Line: 9, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 9, Column: 1},
		},
		{
			name: "ConstStmt",
			node: &ConstStmt{TokenPos: token.Pos{File: "test.zn", Line: 10, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 10, Column: 1},
		},
		{
			name: "TupleDestructStmt",
			node: &TupleDestructStmt{TokenPos: token.Pos{File: "test.zn", Line: 11, Column: 4}},
			want: token.Pos{File: "test.zn", Line: 11, Column: 4},
		},
		{
			name: "AssignStmt",
			node: &AssignStmt{TokenPos: token.Pos{File: "test.zn", Line: 12, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 12, Column: 1},
		},
		{
			name: "MultiAssignStmt",
			node: &MultiAssignStmt{TokenPos: token.Pos{File: "test.zn", Line: 13, Column: 2}},
			want: token.Pos{File: "test.zn", Line: 13, Column: 2},
		},
		{
			name: "ReturnStmt",
			node: &ReturnStmt{TokenPos: token.Pos{File: "test.zn", Line: 14, Column: 5}},
			want: token.Pos{File: "test.zn", Line: 14, Column: 5},
		},
		{
			name: "IfStmt",
			node: &IfStmt{TokenPos: token.Pos{File: "test.zn", Line: 15, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 15, Column: 1},
		},
		{
			name: "IfExpr",
			node: &IfExpr{TokenPos: token.Pos{File: "test.zn", Line: 16, Column: 9}},
			want: token.Pos{File: "test.zn", Line: 16, Column: 9},
		},
		{
			name: "ForStmt",
			node: &ForStmt{TokenPos: token.Pos{File: "test.zn", Line: 17, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 17, Column: 1},
		},
		{
			name: "ForInStmt",
			node: &ForInStmt{TokenPos: token.Pos{File: "test.zn", Line: 18, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 18, Column: 1},
		},
		{
			name: "LoopStmt",
			node: &LoopStmt{TokenPos: token.Pos{File: "test.zn", Line: 19, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 19, Column: 1},
		},
		{
			name: "MatchStmt",
			node: &MatchStmt{TokenPos: token.Pos{File: "test.zn", Line: 20, Column: 3}},
			want: token.Pos{File: "test.zn", Line: 20, Column: 3},
		},
		{
			name: "MatchExpr",
			node: &MatchExpr{TokenPos: token.Pos{File: "test.zn", Line: 21, Column: 7}},
			want: token.Pos{File: "test.zn", Line: 21, Column: 7},
		},
		{
			name: "BreakStmt",
			node: &BreakStmt{TokenPos: token.Pos{File: "test.zn", Line: 22, Column: 5}},
			want: token.Pos{File: "test.zn", Line: 22, Column: 5},
		},
		{
			name: "ContinueStmt",
			node: &ContinueStmt{TokenPos: token.Pos{File: "test.zn", Line: 23, Column: 5}},
			want: token.Pos{File: "test.zn", Line: 23, Column: 5},
		},
		{
			name: "IncDecStmt",
			node: &IncDecStmt{TokenPos: token.Pos{File: "test.zn", Line: 24, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 24, Column: 1},
		},

		// Expressions
		{
			name: "BinaryExpr",
			node: &BinaryExpr{TokenPos: token.Pos{File: "test.zn", Line: 25, Column: 3}},
			want: token.Pos{File: "test.zn", Line: 25, Column: 3},
		},
		{
			name: "UnaryExpr",
			node: &UnaryExpr{TokenPos: token.Pos{File: "test.zn", Line: 26, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 26, Column: 1},
		},
		{
			name: "CallExpr",
			node: &CallExpr{TokenPos: token.Pos{File: "test.zn", Line: 27, Column: 2}},
			want: token.Pos{File: "test.zn", Line: 27, Column: 2},
		},
		{
			name: "IndexExpr",
			node: &IndexExpr{TokenPos: token.Pos{File: "test.zn", Line: 28, Column: 4}},
			want: token.Pos{File: "test.zn", Line: 28, Column: 4},
		},
		{
			name: "SliceExpr",
			node: &SliceExpr{TokenPos: token.Pos{File: "test.zn", Line: 29, Column: 6}},
			want: token.Pos{File: "test.zn", Line: 29, Column: 6},
		},
		{
			name: "FieldExpr",
			node: &FieldExpr{TokenPos: token.Pos{File: "test.zn", Line: 30, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 30, Column: 1},
		},
		{
			name: "TupleLitExpr",
			node: &TupleLitExpr{TokenPos: token.Pos{File: "test.zn", Line: 31, Column: 5}},
			want: token.Pos{File: "test.zn", Line: 31, Column: 5},
		},
		{
			name: "IdentExpr",
			node: &IdentExpr{TokenPos: token.Pos{File: "test.zn", Line: 32, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 32, Column: 1},
		},
		{
			name: "IntLitExpr",
			node: &IntLitExpr{TokenPos: token.Pos{File: "test.zn", Line: 33, Column: 8}},
			want: token.Pos{File: "test.zn", Line: 33, Column: 8},
		},
		{
			name: "FloatLitExpr",
			node: &FloatLitExpr{TokenPos: token.Pos{File: "test.zn", Line: 34, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 34, Column: 1},
		},
		{
			name: "StringLitExpr",
			node: &StringLitExpr{TokenPos: token.Pos{File: "test.zn", Line: 35, Column: 12}},
			want: token.Pos{File: "test.zn", Line: 35, Column: 12},
		},
		{
			name: "BoolLitExpr",
			node: &BoolLitExpr{TokenPos: token.Pos{File: "test.zn", Line: 36, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 36, Column: 1},
		},
		{
			name: "NilExpr",
			node: &NilExpr{TokenPos: token.Pos{File: "test.zn", Line: 37, Column: 3}},
			want: token.Pos{File: "test.zn", Line: 37, Column: 3},
		},
		{
			name: "ArrayLitExpr",
			node: &ArrayLitExpr{TokenPos: token.Pos{File: "test.zn", Line: 38, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 38, Column: 1},
		},
		{
			name: "NamedArgExpr",
			node: &NamedArgExpr{TokenPos: token.Pos{File: "test.zn", Line: 39, Column: 5}},
			want: token.Pos{File: "test.zn", Line: 39, Column: 5},
		},
		{
			name: "InterpStringExpr",
			node: &InterpStringExpr{TokenPos: token.Pos{File: "test.zn", Line: 40, Column: 1}},
			want: token.Pos{File: "test.zn", Line: 40, Column: 1},
		},
		{
			name: "ClosureExpr",
			node: &ClosureExpr{TokenPos: token.Pos{File: "test.zn", Line: 41, Column: 9}},
			want: token.Pos{File: "test.zn", Line: 41, Column: 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.Pos()
			if got != tt.want {
				t.Errorf("%s.Pos() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExprStmtPos(t *testing.T) {
	innerPos := token.Pos{File: "test.zn", Line: 10, Column: 5}
	inner := &CallExpr{TokenPos: innerPos}
	stmt := &ExprStmt{Expr: inner}

	got := stmt.Pos()
	if got != innerPos {
		t.Errorf("ExprStmt.Pos() = %v, want %v (should delegate to inner expr)", got, innerPos)
	}

	// Also verify with a different inner expression type.
	identPos := token.Pos{File: "test.zn", Line: 20, Column: 3}
	stmt2 := &ExprStmt{Expr: &IdentExpr{TokenPos: identPos}}
	got2 := stmt2.Pos()
	if got2 != identPos {
		t.Errorf("ExprStmt.Pos() with IdentExpr = %v, want %v", got2, identPos)
	}
}
