package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mkaz/zenth/pkg/ast"
	"github.com/mkaz/zenth/pkg/lexer"
	"github.com/mkaz/zenth/pkg/token"
)

// Parser converts a token stream into an AST.
type Parser struct {
	tokens []token.Token
	pos    int
	errors []string
}

// New creates a new Parser.
func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) cur() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek() token.Type {
	return p.cur().Type
}

func (p *Parser) peekAt(offset int) token.Type {
	idx := p.pos + offset
	if idx >= len(p.tokens) {
		return token.EOF
	}
	return p.tokens[idx].Type
}

func (p *Parser) advance() token.Token {
	tok := p.cur()
	p.pos++
	return tok
}

func (p *Parser) expect(typ token.Type) token.Token {
	tok := p.cur()
	if tok.Type != typ {
		p.errorf(tok.Pos, "expected %s, got %s", typ, tok.Type)
		return tok
	}
	p.pos++
	return tok
}

func (p *Parser) errorf(pos token.Pos, format string, args ...any) {
	msg := fmt.Sprintf("%s: %s", pos, fmt.Sprintf(format, args...))
	p.errors = append(p.errors, msg)
}

// Parse parses the token stream into a Program AST.
func (p *Parser) Parse() (*ast.Program, error) {
	prog := &ast.Program{}

	for p.peek() != token.EOF {
		stmt := p.parseTopLevel()
		if stmt != nil {
			prog.Stmts = append(prog.Stmts, stmt)
		}
		if len(p.errors) > 10 {
			break
		}
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors:\n%s", joinErrors(p.errors))
	}
	return prog, nil
}

func (p *Parser) parseTopLevel() ast.Node {
	switch p.peek() {
	case token.Fn:
		return p.parseFnDecl()
	case token.Obj:
		return p.parseObjDecl()
	case token.Enum:
		return p.parseEnumDecl()
	case token.Interface:
		return p.parseInterfaceDecl()
	case token.Import:
		return p.parseImportDecl()
	case token.TypeKw:
		return p.parseTypeAliasDecl()
	case token.Let:
		return p.parseLetStmt()
	case token.Var:
		return p.parseVarStmt()
	case token.Const:
		return p.parseConstStmt()
	default:
		p.errorf(p.cur().Pos, "unexpected token at top level: %s", p.peek())
		p.advance()
		return nil
	}
}

func (p *Parser) parseFnDecl() *ast.FnDecl {
	fnTok := p.expect(token.Fn)
	decl := &ast.FnDecl{TokenPos: fnTok.Pos}

	decl.Name = p.expect(token.Ident).Literal

	// Parameters
	p.expect(token.LParen)
	decl.Params = p.parseParams()
	p.expect(token.RParen)

	// Return type
	if p.peek() == token.Arrow {
		p.advance()
		decl.ReturnType = p.parseTypeExpr()
	}

	// Body
	decl.Body = p.parseBlock()

	return decl
}

func (p *Parser) parseParams() []ast.Param {
	var params []ast.Param
	if p.peek() == token.RParen {
		return params
	}
	hasDefault := false
	for {
		nameTok := p.expect(token.Ident)
		name := nameTok.Literal
		p.expect(token.Colon)
		typ := p.parseTypeExpr()
		var def ast.Node
		if p.peek() == token.Assign {
			p.advance()
			def = p.parseExpr(0)
			hasDefault = true
		} else if hasDefault {
			p.errorf(nameTok.Pos, "required parameter '%s' cannot follow a parameter with a default value", name)
		}
		params = append(params, ast.Param{Name: name, Type: typ, Default: def})
		if p.peek() != token.Comma {
			break
		}
		p.advance() // consume comma
	}
	return params
}

func (p *Parser) parseTypeExpr() *ast.TypeExpr {
	pos := p.cur().Pos

	// Shorthand tuple type: (T1, T2, ...)
	if p.peek() == token.LParen {
		p.advance() // consume (
		var params []*ast.TypeExpr
		if p.peek() != token.RParen {
			for {
				params = append(params, p.parseTypeExpr())
				if p.peek() != token.Comma {
					break
				}
				p.advance() // consume ,
			}
		}
		p.expect(token.RParen)
		return &ast.TypeExpr{TokenPos: pos, Name: "tuple", IsTuple: true, Params: params}
	}

	// Array type: array(T)
	if p.peek() == token.Ident && p.cur().Literal == "array" && p.peekAt(1) == token.LParen {
		p.advance() // array
		p.expect(token.LParen)
		elem := p.parseTypeExpr()
		p.expect(token.RParen)
		return &ast.TypeExpr{TokenPos: pos, Name: "array", IsSlice: true, Params: []*ast.TypeExpr{elem}}
	}

	// Set type: set(T)
	if p.peek() == token.Ident && p.cur().Literal == "set" && p.peekAt(1) == token.LParen {
		p.advance() // set
		p.expect(token.LParen)
		elem := p.parseTypeExpr()
		p.expect(token.RParen)
		return &ast.TypeExpr{TokenPos: pos, Name: "set", IsSet: true, Params: []*ast.TypeExpr{elem}}
	}

	// Tuple type: tuple(T1, T2, ...) or tuple(name: T1, name: T2, ...)
	if p.peek() == token.Ident && p.cur().Literal == "tuple" && p.peekAt(1) == token.LParen {
		p.advance() // tuple
		p.expect(token.LParen)
		params := []*ast.TypeExpr{}
		var paramNames []string
		// Detect named tuple type: Ident followed by Colon
		named := p.peek() == token.Ident && p.peekAt(1) == token.Colon
		if p.peek() != token.RParen {
			if named {
				for {
					nameTok := p.expect(token.Ident)
					p.expect(token.Colon)
					paramNames = append(paramNames, nameTok.Literal)
					params = append(params, p.parseTypeExpr())
					if p.peek() != token.Comma {
						break
					}
					p.advance() // ,
				}
			} else {
				for {
					params = append(params, p.parseTypeExpr())
					if p.peek() != token.Comma {
						break
					}
					p.advance() // ,
				}
			}
		}
		p.expect(token.RParen)
		return &ast.TypeExpr{TokenPos: pos, Name: "tuple", IsTuple: true, Params: params, ParamNames: paramNames}
	}

	// Hashmap type: hashmap[K]V or hashmap(K, V)
	if p.peek() == token.Ident && p.cur().Literal == "hashmap" && p.peekAt(1) == token.LBracket {
		p.advance() // hashmap
		p.advance() // [
		key := p.parseTypeExpr()
		p.expect(token.RBracket)
		val := p.parseTypeExpr()
		return &ast.TypeExpr{TokenPos: pos, Name: "hashmap", IsHashmap: true, Params: []*ast.TypeExpr{key, val}}
	}
	if p.peek() == token.Ident && p.cur().Literal == "hashmap" && p.peekAt(1) == token.LParen {
		p.advance() // hashmap
		p.expect(token.LParen)
		key := p.parseTypeExpr()
		p.expect(token.Comma)
		val := p.parseTypeExpr()
		p.expect(token.RParen)
		return &ast.TypeExpr{TokenPos: pos, Name: "hashmap", IsHashmap: true, Params: []*ast.TypeExpr{key, val}}
	}

	name := p.expect(token.Ident).Literal
	return &ast.TypeExpr{TokenPos: pos, Name: name}
}

func (p *Parser) parseBlock() *ast.Block {
	tok := p.expect(token.LBrace)
	block := &ast.Block{TokenPos: tok.Pos}
	for p.peek() != token.RBrace && p.peek() != token.EOF {
		stmt := p.parseStmt()
		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		}
	}
	p.expect(token.RBrace)
	return block
}

func (p *Parser) parseStmt() ast.Node {
	switch p.peek() {
	case token.Let:
		return p.parseLetStmt()
	case token.Var:
		return p.parseVarStmt()
	case token.Const:
		return p.parseConstStmt()
	case token.Return:
		return p.parseReturnStmt()
	case token.If:
		return p.parseIfStmt()
	case token.For:
		return p.parseForStmt()
	case token.Loop:
		return p.parseLoopStmt()
	case token.Match:
		return p.parseMatchStmt()
	case token.Break:
		tok := p.advance()
		p.expect(token.Semicolon)
		return &ast.BreakStmt{TokenPos: tok.Pos}
	case token.Continue:
		tok := p.advance()
		p.expect(token.Semicolon)
		return &ast.ContinueStmt{TokenPos: tok.Pos}
	default:
		return p.parseExprOrAssignStmt()
	}
}

func (p *Parser) parseTupleBindingNames() []string {
	p.expect(token.LParen)
	var names []string
	for {
		names = append(names, p.expect(token.Ident).Literal)
		if p.peek() != token.Comma {
			break
		}
		p.advance()
	}
	p.expect(token.RParen)
	return names
}

func (p *Parser) parseArrayBindingNames() []string {
	p.expect(token.LBracket)
	var names []string
	for {
		names = append(names, p.expect(token.Ident).Literal)
		if p.peek() != token.Comma {
			break
		}
		p.advance()
	}
	p.expect(token.RBracket)
	return names
}

func (p *Parser) parseLetStmt() ast.Node {
	tok := p.expect(token.Let)
	if p.peek() == token.LParen {
		names := p.parseTupleBindingNames()
		p.expect(token.Assign)
		value := p.parseExpr(0)
		p.expect(token.Semicolon)
		return &ast.TupleDestructStmt{TokenPos: tok.Pos, Kind: token.Let, Names: names, Value: value}
	}
	if p.peek() == token.LBracket {
		names := p.parseArrayBindingNames()
		p.expect(token.Assign)
		value := p.parseExpr(0)
		p.expect(token.Semicolon)
		return &ast.ArrayDestructStmt{TokenPos: tok.Pos, Kind: token.Let, Names: names, Value: value}
	}
	stmt := &ast.LetStmt{TokenPos: tok.Pos}
	stmt.Name = p.expect(token.Ident).Literal

	if p.peek() == token.Colon {
		// let name: Type = expr
		p.advance()
		stmt.Type = p.parseTypeExpr()
		p.expect(token.Assign)
		stmt.Value = p.parseExpr(0)
	} else if p.peek() == token.Assign {
		// let name = expr  (infer type)
		p.advance()
		stmt.Infer = true
		stmt.Value = p.parseExpr(0)
	} else if p.peek() == token.Walrus {
		p.errorf(tok.Pos, "Zenth uses '=' for assignment, not ':='")
	} else {
		p.errorf(tok.Pos, "expected = or : after let variable name")
	}

	p.expect(token.Semicolon)
	return stmt
}

func (p *Parser) parseVarStmt() ast.Node {
	tok := p.expect(token.Var)
	if p.peek() == token.LParen {
		names := p.parseTupleBindingNames()
		p.expect(token.Assign)
		value := p.parseExpr(0)
		p.expect(token.Semicolon)
		return &ast.TupleDestructStmt{TokenPos: tok.Pos, Kind: token.Var, Names: names, Value: value}
	}
	if p.peek() == token.LBracket {
		names := p.parseArrayBindingNames()
		p.expect(token.Assign)
		value := p.parseExpr(0)
		p.expect(token.Semicolon)
		return &ast.ArrayDestructStmt{TokenPos: tok.Pos, Kind: token.Var, Names: names, Value: value}
	}
	stmt := &ast.VarStmt{TokenPos: tok.Pos}
	stmt.Name = p.expect(token.Ident).Literal

	if p.peek() == token.Colon {
		// var name: Type = expr
		p.advance()
		stmt.Type = p.parseTypeExpr()
		p.expect(token.Assign)
		stmt.Value = p.parseExpr(0)
	} else if p.peek() == token.Assign {
		// var name = expr  (infer type)
		p.advance()
		stmt.Infer = true
		stmt.Value = p.parseExpr(0)
	} else if p.peek() == token.Walrus {
		p.errorf(tok.Pos, "Zenth uses '=' for assignment, not ':='")
	} else {
		p.errorf(tok.Pos, "expected = or : after var variable name")
	}

	p.expect(token.Semicolon)
	return stmt
}

func (p *Parser) parseConstStmt() ast.Node {
	tok := p.expect(token.Const)
	if p.peek() == token.LParen {
		names := p.parseTupleBindingNames()
		p.expect(token.Assign)
		value := p.parseExpr(0)
		p.expect(token.Semicolon)
		return &ast.TupleDestructStmt{TokenPos: tok.Pos, Kind: token.Const, Names: names, Value: value}
	}
	if p.peek() == token.LBracket {
		names := p.parseArrayBindingNames()
		p.expect(token.Assign)
		value := p.parseExpr(0)
		p.expect(token.Semicolon)
		return &ast.ArrayDestructStmt{TokenPos: tok.Pos, Kind: token.Const, Names: names, Value: value}
	}
	stmt := &ast.ConstStmt{TokenPos: tok.Pos}
	stmt.Name = p.expect(token.Ident).Literal

	if p.peek() == token.Colon {
		// const name: Type = expr
		p.advance()
		stmt.Type = p.parseTypeExpr()
		p.expect(token.Assign)
		stmt.Value = p.parseExpr(0)
	} else if p.peek() == token.Assign {
		// const name = expr  (infer type)
		p.advance()
		stmt.Infer = true
		stmt.Value = p.parseExpr(0)
	} else if p.peek() == token.Walrus {
		p.errorf(tok.Pos, "Zenth uses '=' for assignment, not ':='")
	} else {
		p.errorf(tok.Pos, "expected = or : after const name")
	}

	p.expect(token.Semicolon)
	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	tok := p.expect(token.Return)
	stmt := &ast.ReturnStmt{TokenPos: tok.Pos}
	if p.peek() != token.Semicolon {
		first := p.parseExpr(0)
		// Multi-value return: return a, b, c; — desugars to return tuple(a, b, c);
		if p.peek() == token.Comma {
			elems := []ast.Node{first}
			for p.peek() == token.Comma {
				p.advance() // consume ,
				elems = append(elems, p.parseExpr(0))
			}
			stmt.Value = &ast.TupleLitExpr{TokenPos: tok.Pos, Elements: elems}
		} else {
			stmt.Value = first
		}
	}
	p.expect(token.Semicolon)
	return stmt
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	tok := p.expect(token.If)
	stmt := &ast.IfStmt{TokenPos: tok.Pos}
	stmt.Condition = p.parseExpr(0)
	stmt.Body = p.parseBlock()
	if p.peek() == token.Else {
		p.advance()
		if p.peek() == token.If {
			stmt.Else = p.parseIfStmt()
		} else {
			stmt.Else = p.parseBlock()
		}
	}
	return stmt
}

func (p *Parser) parseForStmt() ast.Node {
	tok := p.expect(token.For)

	// for { ... } -- infinite loop
	if p.peek() == token.LBrace {
		body := p.parseBlock()
		return &ast.ForStmt{TokenPos: tok.Pos, Body: body}
	}

	// for range count { ... } -- count-style
	if p.isForRangeCount() {
		return p.parseForRangeCountStmt(tok)
	}

	// Try to detect for...in: look for `ident in` or `ident, ident in`
	if p.isForIn() {
		return p.parseForInStmt(tok)
	}

	// for condition { ... } -- while-style
	// or for init; cond; post { ... } -- C-style
	return p.parseCStyleFor(tok)
}

func (p *Parser) isForRangeCount() bool {
	return p.peek() == token.Ident && p.cur().Literal == "range" && p.peekAt(1) != token.LParen
}

func (p *Parser) parseForRangeCountStmt(tok token.Token) *ast.LoopStmt {
	p.advance() // consume "range"
	count := p.parseExpr(0)
	body := p.parseBlock()
	return &ast.LoopStmt{TokenPos: tok.Pos, Count: count, Body: body}
}

func (p *Parser) isForIn() bool {
	// ident in ...
	if p.peekAt(0) == token.Ident && p.peekAt(1) == token.In {
		return true
	}
	// ident, ident in ...
	if p.peekAt(0) == token.Ident && p.peekAt(1) == token.Comma && p.peekAt(2) == token.Ident && p.peekAt(3) == token.In {
		return true
	}
	return false
}

func (p *Parser) parseForInStmt(tok token.Token) *ast.ForInStmt {
	stmt := &ast.ForInStmt{TokenPos: tok.Pos}

	first := p.expect(token.Ident).Literal
	if p.peek() == token.Comma {
		p.advance()
		stmt.Index = first
		stmt.Value = p.expect(token.Ident).Literal
	} else {
		stmt.Value = first
	}

	p.expect(token.In)
	stmt.Iterable = p.parseExpr(0)
	stmt.Body = p.parseBlock()
	return stmt
}

func (p *Parser) parseCStyleFor(tok token.Token) *ast.ForStmt {
	stmt := &ast.ForStmt{TokenPos: tok.Pos}

	// Check for while-style: expr { ... }
	// vs C-style: init; cond; post { ... }
	// We detect C-style by trying to parse the first part and seeing if a semicolon follows
	startPos := p.pos
	_ = startPos

	// Try: is the first semicolon before the first { ?
	if p.hasSemicolonBeforeBrace() {
		// C-style for
		stmt.Init = p.parseSimpleStmt()
		p.expect(token.Semicolon)
		if p.peek() != token.Semicolon {
			stmt.Condition = p.parseExpr(0)
		}
		p.expect(token.Semicolon)
		if p.peek() != token.LBrace {
			stmt.Post = p.parseSimpleStmtNoSemicolon()
		}
	} else {
		// While-style: just a condition
		stmt.Condition = p.parseExpr(0)
	}

	stmt.Body = p.parseBlock()
	return stmt
}

func (p *Parser) hasSemicolonBeforeBrace() bool {
	depth := 0
	for i := p.pos; i < len(p.tokens); i++ {
		switch p.tokens[i].Type {
		case token.LParen, token.LBracket:
			depth++
		case token.RParen, token.RBracket:
			depth--
		case token.Semicolon:
			if depth == 0 {
				return true
			}
		case token.LBrace:
			if depth == 0 {
				return false
			}
		case token.EOF:
			return false
		}
	}
	return false
}

// parseSimpleStmt parses a let/var or expression (for 'for' init clauses).
func (p *Parser) parseSimpleStmt() ast.Node {
	switch p.peek() {
	case token.Let:
		tok := p.expect(token.Let)
		if p.peek() == token.LParen {
			names := p.parseTupleBindingNames()
			p.expect(token.Assign)
			value := p.parseExpr(0)
			return &ast.TupleDestructStmt{TokenPos: tok.Pos, Kind: token.Let, Names: names, Value: value}
		}
		stmt := &ast.LetStmt{TokenPos: tok.Pos}
		stmt.Name = p.expect(token.Ident).Literal
		if p.peek() == token.Colon {
			p.advance()
			stmt.Type = p.parseTypeExpr()
			p.expect(token.Assign)
			stmt.Value = p.parseExpr(0)
		} else if p.peek() == token.Assign {
			p.advance()
			stmt.Infer = true
			stmt.Value = p.parseExpr(0)
		}
		return stmt
	case token.Var:
		tok := p.expect(token.Var)
		if p.peek() == token.LParen {
			names := p.parseTupleBindingNames()
			p.expect(token.Assign)
			value := p.parseExpr(0)
			return &ast.TupleDestructStmt{TokenPos: tok.Pos, Kind: token.Var, Names: names, Value: value}
		}
		stmt := &ast.VarStmt{TokenPos: tok.Pos}
		stmt.Name = p.expect(token.Ident).Literal
		if p.peek() == token.Colon {
			p.advance()
			stmt.Type = p.parseTypeExpr()
			p.expect(token.Assign)
			stmt.Value = p.parseExpr(0)
		} else if p.peek() == token.Assign {
			p.advance()
			stmt.Infer = true
			stmt.Value = p.parseExpr(0)
		}
		return stmt
	default:
		return p.parseSimpleStmtNoSemicolon()
	}
}

func (p *Parser) parseSimpleStmtNoSemicolon() ast.Node {
	expr := p.parseExpr(0)
	// Multi-assign: expr, expr, ... = expr, expr, ...
	if p.peek() == token.Comma {
		targets := []ast.Node{expr}
		for p.peek() == token.Comma {
			p.advance()
			targets = append(targets, p.parseExpr(0))
		}
		eqTok := p.expect(token.Assign)
		values := []ast.Node{p.parseExpr(0)}
		for p.peek() == token.Comma {
			p.advance()
			values = append(values, p.parseExpr(0))
		}
		return &ast.MultiAssignStmt{TokenPos: eqTok.Pos, Targets: targets, Values: values}
	}
	// Check for assignment
	if p.peek() == token.Assign || p.peek() == token.PlusAssign || p.peek() == token.MinusAssign || p.peek() == token.StarAssign || p.peek() == token.SlashAssign {
		op := p.advance()
		value := p.parseExpr(0)
		return &ast.AssignStmt{TokenPos: op.Pos, Target: expr, Op: op.Type, Value: value}
	}
	// Check for ++ or --
	if p.peek() == token.PlusPlus || p.peek() == token.MinusMinus {
		op := p.advance()
		return &ast.IncDecStmt{TokenPos: op.Pos, Operand: expr, Op: op.Type}
	}
	return &ast.ExprStmt{Expr: expr}
}

func (p *Parser) parseLoopStmt() *ast.LoopStmt {
	tok := p.expect(token.Loop)
	count := p.parseExpr(0)
	body := p.parseBlock()
	return &ast.LoopStmt{TokenPos: tok.Pos, Count: count, Body: body}
}

func (p *Parser) parseMatchStmt() *ast.MatchStmt {
	tok := p.expect(token.Match)
	stmt := &ast.MatchStmt{TokenPos: tok.Pos}
	stmt.Subject = p.parseExpr(0)
	p.expect(token.LBrace)
	for p.peek() != token.RBrace && p.peek() != token.EOF {
		arm := p.parseMatchArm()
		stmt.Arms = append(stmt.Arms, arm)
	}
	p.expect(token.RBrace)
	return stmt
}

func (p *Parser) parseMatchArm() ast.MatchArm {
	var arm ast.MatchArm
	arm.Pattern = p.parseExpr(0)
	p.expect(token.FatArrow)
	if p.peek() == token.LBrace {
		arm.Body = p.parseBlock()
	} else {
		arm.Body = p.parseStmt()
	}
	return arm
}

func (p *Parser) parseExprOrAssignStmt() ast.Node {
	expr := p.parseExpr(0)

	// Multi-assign: expr, expr, ... = expr, expr, ...;
	if p.peek() == token.Comma {
		targets := []ast.Node{expr}
		for p.peek() == token.Comma {
			p.advance()
			targets = append(targets, p.parseExpr(0))
		}
		eqTok := p.expect(token.Assign)
		values := []ast.Node{p.parseExpr(0)}
		for p.peek() == token.Comma {
			p.advance()
			values = append(values, p.parseExpr(0))
		}
		p.expect(token.Semicolon)
		return &ast.MultiAssignStmt{TokenPos: eqTok.Pos, Targets: targets, Values: values}
	}

	if p.peek() == token.Walrus {
		p.errorf(p.cur().Pos, "Zenth uses '=' for assignment, not ':='")
	}

	// Assignment
	if p.peek() == token.Assign || p.peek() == token.PlusAssign || p.peek() == token.MinusAssign || p.peek() == token.StarAssign || p.peek() == token.SlashAssign {
		op := p.advance()
		value := p.parseExpr(0)
		p.expect(token.Semicolon)
		return &ast.AssignStmt{TokenPos: op.Pos, Target: expr, Op: op.Type, Value: value}
	}

	// IncDec
	if p.peek() == token.PlusPlus || p.peek() == token.MinusMinus {
		op := p.advance()
		p.expect(token.Semicolon)
		return &ast.IncDecStmt{TokenPos: op.Pos, Operand: expr, Op: op.Type}
	}

	p.expect(token.Semicolon)
	return &ast.ExprStmt{Expr: expr}
}

// ---------- Expression Parsing (Pratt) ----------

func (p *Parser) parseExpr(minPrec int) ast.Node {
	left := p.parseUnary()

	for {
		prec := p.precedence(p.peek())
		if prec <= minPrec {
			break
		}
		op := p.advance()
		right := p.parseExpr(prec)
		left = &ast.BinaryExpr{
			TokenPos: op.Pos,
			Left:     left,
			Op:       op.Type,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) precedence(t token.Type) int {
	switch t {
	case token.Or:
		return 1
	case token.And:
		return 2
	case token.Eq, token.Neq:
		return 3
	case token.Lt, token.Gt, token.Lte, token.Gte:
		return 4
	case token.Plus, token.Minus:
		return 5
	case token.Star, token.Slash, token.Percent:
		return 6
	default:
		return 0
	}
}

func (p *Parser) parseUnary() ast.Node {
	if p.peek() == token.Not || p.peek() == token.Minus {
		op := p.advance()
		operand := p.parseUnary()
		return &ast.UnaryExpr{TokenPos: op.Pos, Op: op.Type, Operand: operand}
	}
	return p.parsePostfix()
}

func (p *Parser) parsePostfix() ast.Node {
	expr := p.parsePrimary()

	for {
		switch p.peek() {
		case token.LParen:
			expr = p.parseCallExpr(expr)
		case token.LBracket:
			tok := p.advance()
			// [:high] — slice with no low
			if p.peek() == token.Colon {
				p.advance()
				var high ast.Node
				if p.peek() != token.RBracket {
					high = p.parseExpr(0)
				}
				p.expect(token.RBracket)
				expr = &ast.SliceExpr{TokenPos: tok.Pos, Object: expr, High: high}
			} else {
				index := p.parseExpr(0)
				if p.peek() == token.Colon {
					// [low:] or [low:high] — slice
					p.advance()
					var high ast.Node
					if p.peek() != token.RBracket {
						high = p.parseExpr(0)
					}
					p.expect(token.RBracket)
					expr = &ast.SliceExpr{TokenPos: tok.Pos, Object: expr, Low: index, High: high}
				} else {
					// [index] — plain index
					p.expect(token.RBracket)
					expr = &ast.IndexExpr{TokenPos: tok.Pos, Object: expr, Index: index}
				}
			}
		case token.Dot:
			tok := p.advance()
			fieldTok := p.cur()
			var field string
			if fieldTok.Type == token.Ident || fieldTok.Type == token.IntLit {
				field = p.advance().Literal
			} else {
				p.errorf(fieldTok.Pos, "expected field name after '.', got %s", fieldTok.Type)
				field = "<error>"
			}
			expr = &ast.FieldExpr{TokenPos: tok.Pos, Object: expr, Field: field}
		default:
			return expr
		}
	}
}

func (p *Parser) parseCallExpr(callee ast.Node) *ast.CallExpr {
	tok := p.expect(token.LParen)
	call := &ast.CallExpr{TokenPos: tok.Pos, Callee: callee}
	if p.peek() != token.RParen {
		for {
			// Detect named arg: ident = value
			if p.peek() == token.Ident && p.peekAt(1) == token.Assign {
				nameTok := p.advance()
				eqTok := p.advance() // consume =
				value := p.parseExpr(0)
				call.Args = append(call.Args, &ast.NamedArgExpr{
					TokenPos: eqTok.Pos,
					Name:     nameTok.Literal,
					Value:    value,
				})
			} else {
				arg := p.parseExpr(0)
				call.Args = append(call.Args, arg)
			}
			if p.peek() != token.Comma {
				break
			}
			p.advance()
			// Allow trailing comma
			if p.peek() == token.RParen {
				break
			}
		}
	}
	p.expect(token.RParen)
	return call
}

func (p *Parser) parsePrimary() ast.Node {
	tok := p.cur()

	switch tok.Type {
	case token.IntLit:
		p.advance()
		val, _ := strconv.ParseInt(tok.Literal, 10, 64)
		return &ast.IntLitExpr{TokenPos: tok.Pos, Value: val}

	case token.FloatLit:
		p.advance()
		val, _ := strconv.ParseFloat(tok.Literal, 64)
		return &ast.FloatLitExpr{TokenPos: tok.Pos, Value: val}

	case token.StringLit:
		p.advance()
		return &ast.StringLitExpr{TokenPos: tok.Pos, Value: tok.Literal}

	case token.InterpStringLit:
		p.advance()
		return p.parseInterpString(tok)

	case token.True:
		p.advance()
		return &ast.BoolLitExpr{TokenPos: tok.Pos, Value: true}

	case token.False:
		p.advance()
		return &ast.BoolLitExpr{TokenPos: tok.Pos, Value: false}

	case token.Nil:
		p.advance()
		return &ast.NilExpr{TokenPos: tok.Pos}

	case token.Ident:
		p.advance()
		if tok.Literal == "tuple" && p.peek() == token.LParen {
			p.advance() // consume '('
			tuple := &ast.TupleLitExpr{TokenPos: tok.Pos}
			// Detect named tuple: first element is Ident followed by Assign
			named := p.peek() == token.Ident && p.peekAt(1) == token.Assign
			if named {
				for {
					nameTok := p.expect(token.Ident)
					p.expect(token.Assign)
					tuple.Names = append(tuple.Names, nameTok.Literal)
					tuple.Elements = append(tuple.Elements, p.parseExpr(0))
					if p.peek() != token.Comma {
						break
					}
					p.advance()
					if p.peek() == token.RParen {
						break
					}
					if p.peek() != token.Ident || p.peekAt(1) != token.Assign {
						p.errorf(p.cur().Pos, "cannot mix named and positional tuple fields")
						break
					}
				}
			} else {
				if p.peek() != token.RParen {
					tuple.Elements = append(tuple.Elements, p.parseExpr(0))
					for p.peek() == token.Comma {
						p.advance()
						if p.peek() == token.RParen {
							break
						}
						tuple.Elements = append(tuple.Elements, p.parseExpr(0))
					}
				}
			}
			p.expect(token.RParen)
			return tuple
		}
		return &ast.IdentExpr{TokenPos: tok.Pos, Name: tok.Literal}

	case token.LParen:
		pos := tok.Pos
		p.advance()
		expr := p.parseExpr(0)
		p.expect(token.RParen)
		return &ast.GroupedExpr{TokenPos: pos, Expr: expr}

	case token.LBracket:
		return p.parseArrayLit()

	case token.If:
		return p.parseIfExpr()

	case token.Match:
		return p.parseMatchExpr()

	case token.Fn:
		// fn( starts a closure expression
		if p.peekAt(1) == token.LParen {
			return p.parseClosureExpr()
		}
		p.errorf(tok.Pos, "unexpected fn in expression position")
		p.advance()
		return &ast.IdentExpr{TokenPos: tok.Pos, Name: "<error>"}

	default:
		p.errorf(tok.Pos, "unexpected token: %s (%q)", tok.Type, tok.Literal)
		p.advance()
		return &ast.IdentExpr{TokenPos: tok.Pos, Name: "<error>"}
	}
}

func (p *Parser) parseIfExpr() *ast.IfExpr {
	tok := p.expect(token.If)
	expr := &ast.IfExpr{TokenPos: tok.Pos}

	// Parse condition
	expr.Condition = p.parseExpr(0)

	// Parse then-branch: { expr }
	p.expect(token.LBrace)
	expr.Then = p.parseExpr(0)
	p.expect(token.RBrace)

	// else is mandatory for if-expressions
	p.expect(token.Else)

	// else-if chain or final else branch
	if p.peek() == token.If {
		expr.Else = p.parseIfExpr()
	} else {
		p.expect(token.LBrace)
		expr.Else = p.parseExpr(0)
		p.expect(token.RBrace)
	}

	return expr
}

func (p *Parser) parseMatchExpr() *ast.MatchExpr {
	tok := p.expect(token.Match)
	expr := &ast.MatchExpr{TokenPos: tok.Pos}
	expr.Subject = p.parseExpr(0)
	p.expect(token.LBrace)
	for p.peek() != token.RBrace && p.peek() != token.EOF {
		var arm ast.MatchExprArm
		arm.Pattern = p.parseExpr(0)
		p.expect(token.FatArrow)
		arm.Value = p.parseExpr(0)
		expr.Arms = append(expr.Arms, arm)
		// Require comma between arms, optional trailing comma
		if p.peek() == token.Comma {
			p.advance()
		} else if p.peek() != token.RBrace {
			p.expect(token.Comma)
		}
	}
	p.expect(token.RBrace)
	return expr
}

func (p *Parser) parseArrayLit() *ast.ArrayLitExpr {
	tok := p.expect(token.LBracket)
	lit := &ast.ArrayLitExpr{TokenPos: tok.Pos}
	if p.peek() != token.RBracket {
		for {
			elem := p.parseExpr(0)
			lit.Elements = append(lit.Elements, elem)
			if p.peek() != token.Comma {
				break
			}
			p.advance()
			// Allow trailing comma
			if p.peek() == token.RBracket {
				break
			}
		}
	}
	p.expect(token.RBracket)
	return lit
}

func (p *Parser) parseObjDecl() *ast.ObjDecl {
	tok := p.expect(token.Obj)
	decl := &ast.ObjDecl{TokenPos: tok.Pos}
	decl.Name = p.expect(token.Ident).Literal
	p.expect(token.LBrace)
	for p.peek() != token.RBrace && p.peek() != token.EOF {
		if p.peek() == token.Fn {
			// Method inside struct
			method := p.parseFnDecl()
			method.OwnerObj = decl.Name
			decl.Methods = append(decl.Methods, method)
		} else {
			// Field
			name := p.expect(token.Ident).Literal
			p.expect(token.Colon)
			typ := p.parseTypeExpr()
			var def ast.Node
			if p.peek() == token.Assign {
				p.advance()
				def = p.parseExpr(0)
			}
			p.expect(token.Semicolon)
			decl.Fields = append(decl.Fields, ast.Field{Name: name, Type: typ, Default: def})
		}
	}
	p.expect(token.RBrace)
	return decl
}

func (p *Parser) parseEnumDecl() *ast.EnumDecl {
	tok := p.expect(token.Enum)
	decl := &ast.EnumDecl{TokenPos: tok.Pos}
	decl.Name = p.expect(token.Ident).Literal
	p.expect(token.LBrace)
	for p.peek() != token.RBrace && p.peek() != token.EOF {
		name := p.expect(token.Ident).Literal
		variant := ast.EnumVariant{Name: name}
		if p.peek() == token.Assign {
			p.advance()
			valTok := p.cur()
			if valTok.Type != token.StringLit {
				p.errorf(valTok.Pos, "enum value must be a string literal")
				p.advance()
			} else {
				p.advance()
				variant.StrValue = &ast.StringLitExpr{TokenPos: valTok.Pos, Value: valTok.Literal}
			}
		}
		p.expect(token.Semicolon)
		decl.Variants = append(decl.Variants, variant)
	}
	p.expect(token.RBrace)
	return decl
}

func (p *Parser) parseInterfaceDecl() *ast.InterfaceDecl {
	tok := p.expect(token.Interface)
	decl := &ast.InterfaceDecl{TokenPos: tok.Pos}
	decl.Name = p.expect(token.Ident).Literal
	p.expect(token.LBrace)
	for p.peek() != token.RBrace && p.peek() != token.EOF {
		p.expect(token.Fn)
		name := p.expect(token.Ident).Literal
		p.expect(token.LParen)
		params := p.parseParams()
		p.expect(token.RParen)
		var retType *ast.TypeExpr
		if p.peek() == token.Arrow {
			p.advance()
			retType = p.parseTypeExpr()
		}
		p.expect(token.Semicolon)
		decl.Methods = append(decl.Methods, ast.MethodSig{Name: name, Params: params, ReturnType: retType})
	}
	p.expect(token.RBrace)
	return decl
}

func (p *Parser) parseImportDecl() *ast.ImportDecl {
	tok := p.expect(token.Import)
	decl := &ast.ImportDecl{TokenPos: tok.Pos}
	decl.Path = p.expect(token.StringLit).Literal
	decl.IsLocal = strings.HasPrefix(decl.Path, "./") || strings.HasPrefix(decl.Path, "../")
	if p.peek() == token.As {
		p.advance()
		decl.Alias = p.expect(token.Ident).Literal
	}
	p.expect(token.Semicolon)
	return decl
}

func (p *Parser) parseTypeAliasDecl() *ast.TypeAliasDecl {
	tok := p.expect(token.TypeKw)
	name := p.expect(token.Ident).Literal
	p.expect(token.Assign)
	typ := p.parseTypeExpr()
	p.expect(token.Semicolon)
	return &ast.TypeAliasDecl{TokenPos: tok.Pos, Name: name, Type: typ}
}

func (p *Parser) parseInterpString(tok token.Token) ast.Node {
	raw := tok.Literal
	var parts []ast.InterpPart
	var buf strings.Builder

	i := 0
	for i < len(raw) {
		if raw[i] == '{' {
			// Flush any accumulated text
			if buf.Len() > 0 {
				parts = append(parts, ast.InterpPart{Lit: buf.String()})
				buf.Reset()
			}
			// Find matching closing brace, accounting for nested braces
			i++ // skip opening {
			depth := 1
			start := i
			for i < len(raw) && depth > 0 {
				if raw[i] == '{' {
					depth++
				} else if raw[i] == '}' {
					depth--
				}
				if depth > 0 {
					i++
				}
			}
			exprStr := raw[start:i]
			if i < len(raw) {
				i++ // skip closing }
			}
			// Sub-lex and sub-parse the expression
			subLex := lexer.New(tok.Pos.File, exprStr)
			subTokens, err := subLex.Tokenize()
			if err != nil {
				p.errorf(tok.Pos, "error in interpolated expression: %s", err)
				continue
			}
			subParser := New(subTokens)
			expr := subParser.parseExpr(0)
			if len(subParser.errors) > 0 {
				for _, e := range subParser.errors {
					p.errors = append(p.errors, e)
				}
				continue
			}
			parts = append(parts, ast.InterpPart{IsExpr: true, Expr: expr})
		} else {
			buf.WriteByte(raw[i])
			i++
		}
	}
	// Flush remaining text
	if buf.Len() > 0 {
		parts = append(parts, ast.InterpPart{Lit: buf.String()})
	}

	return &ast.InterpStringExpr{TokenPos: tok.Pos, Parts: parts}
}

func (p *Parser) parseClosureExpr() *ast.ClosureExpr {
	tok := p.expect(token.Fn)
	expr := &ast.ClosureExpr{TokenPos: tok.Pos}

	// Parameters (with optional types)
	p.expect(token.LParen)
	expr.Params = p.parseClosureParams()
	p.expect(token.RParen)

	// Optional return type
	if p.peek() == token.Arrow {
		p.advance()
		expr.ReturnType = p.parseTypeExpr()
	}

	// Body: block or single expression
	if p.peek() == token.LBrace {
		expr.Body = p.parseBlock()
	} else {
		expr.Body = p.parseExpr(0)
	}

	return expr
}

func (p *Parser) parseClosureParams() []ast.Param {
	var params []ast.Param
	if p.peek() == token.RParen {
		return params
	}
	for {
		name := p.expect(token.Ident).Literal
		var typ *ast.TypeExpr
		if p.peek() == token.Colon {
			p.advance()
			typ = p.parseTypeExpr()
		}
		params = append(params, ast.Param{Name: name, Type: typ})
		if p.peek() != token.Comma {
			break
		}
		p.advance()
	}
	return params
}

func joinErrors(errs []string) string {
	var result strings.Builder
	for i, e := range errs {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString("  " + e)
	}
	return result.String()
}
