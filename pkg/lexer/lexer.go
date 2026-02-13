package lexer

import (
	"fmt"
	"unicode"

	"github.com/mkaz/zenth/pkg/token"
)

// Lexer tokenizes Zenth source code.
type Lexer struct {
	src      []rune
	file     string
	pos      int
	line     int
	col      int
	tokens   []token.Token
}

// New creates a new Lexer for the given source.
func New(filename, src string) *Lexer {
	return &Lexer{
		src:  []rune(src),
		file: filename,
		line: 1,
		col:  1,
	}
}

func (l *Lexer) curPos() token.Pos {
	return token.Pos{File: l.file, Line: l.line, Column: l.col}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	return l.src[l.pos]
}

func (l *Lexer) peekAt(offset int) rune {
	idx := l.pos + offset
	if idx >= len(l.src) {
		return 0
	}
	return l.src[idx]
}

func (l *Lexer) advance() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	ch := l.src[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) emit(typ token.Type, lit string, pos token.Pos) {
	l.tokens = append(l.tokens, token.Token{Type: typ, Literal: lit, Pos: pos})
}

func (l *Lexer) error(pos token.Pos, msg string) {
	l.emit(token.Illegal, msg, pos)
}

// Tokenize lexes the entire source and returns the token stream.
func (l *Lexer) Tokenize() ([]token.Token, error) {
	for l.pos < len(l.src) {
		l.skipWhitespaceAndComments()
		if l.pos >= len(l.src) {
			break
		}

		pos := l.curPos()
		ch := l.peek()

		switch {
		case ch == '"':
			l.readString(pos)
		case ch == '\'':
			l.readRawString(pos)
		case isDigit(ch):
			l.readNumber(pos)
		case isIdentStart(ch):
			l.readIdent(pos)
		default:
			l.readOperator(pos)
		}
	}

	l.emit(token.EOF, "", l.curPos())

	// Check for any illegal tokens
	for _, tok := range l.tokens {
		if tok.Type == token.Illegal {
			return l.tokens, fmt.Errorf("%s: %s", tok.Pos, tok.Literal)
		}
	}

	return l.tokens, nil
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.src) {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.advance()
			continue
		}
		// Line comment: // or #
		if (ch == '/' && l.peekAt(1) == '/') || ch == '#' {
			for l.pos < len(l.src) && l.peek() != '\n' {
				l.advance()
			}
			continue
		}
		// Block comment
		if ch == '/' && l.peekAt(1) == '*' {
			l.advance() // /
			l.advance() // *
			for l.pos < len(l.src) {
				if l.peek() == '*' && l.peekAt(1) == '/' {
					l.advance() // *
					l.advance() // /
					break
				}
				l.advance()
			}
			continue
		}
		break
	}
}

func (l *Lexer) readString(pos token.Pos) {
	l.advance() // opening "
	var lit []rune
	hasInterp := false
	for l.pos < len(l.src) {
		ch := l.peek()
		if ch == '"' {
			l.advance() // closing "
			if hasInterp {
				l.emit(token.InterpStringLit, string(lit), pos)
			} else {
				l.emit(token.StringLit, string(lit), pos)
			}
			return
		}
		if ch == '\\' {
			l.advance()
			esc := l.advance()
			switch esc {
			case 'n':
				lit = append(lit, '\n')
			case 't':
				lit = append(lit, '\t')
			case '\\':
				lit = append(lit, '\\')
			case '"':
				lit = append(lit, '"')
			case '0':
				lit = append(lit, 0)
			case '{':
				lit = append(lit, '{')
			case '}':
				lit = append(lit, '}')
			default:
				lit = append(lit, '\\', esc)
			}
			continue
		}
		if ch == '{' {
			hasInterp = true
		}
		if ch == '\n' {
			l.error(pos, "unterminated string literal")
			return
		}
		lit = append(lit, l.advance())
	}
	l.error(pos, "unterminated string literal")
}

func (l *Lexer) readRawString(pos token.Pos) {
	l.advance() // opening '
	var lit []rune
	for l.pos < len(l.src) {
		ch := l.peek()
		if ch == '\'' {
			l.advance() // closing '
			l.emit(token.StringLit, string(lit), pos)
			return
		}
		if ch == '\\' {
			l.advance()
			esc := l.advance()
			switch esc {
			case 'n':
				lit = append(lit, '\n')
			case 't':
				lit = append(lit, '\t')
			case '\\':
				lit = append(lit, '\\')
			case '\'':
				lit = append(lit, '\'')
			case '0':
				lit = append(lit, 0)
			default:
				lit = append(lit, '\\', esc)
			}
			continue
		}
		if ch == '\n' {
			l.error(pos, "unterminated string literal")
			return
		}
		lit = append(lit, l.advance())
	}
	l.error(pos, "unterminated string literal")
}

func (l *Lexer) readNumber(pos token.Pos) {
	var lit []rune
	isFloat := false

	for l.pos < len(l.src) {
		ch := l.peek()
		if ch == '.' && !isFloat && isDigit(l.peekAt(1)) {
			isFloat = true
			lit = append(lit, l.advance())
			continue
		}
		if ch == '_' {
			l.advance() // skip underscores in number literals (e.g., 1_000_000)
			continue
		}
		if !isDigit(ch) {
			break
		}
		lit = append(lit, l.advance())
	}

	if isFloat {
		l.emit(token.FloatLit, string(lit), pos)
	} else {
		l.emit(token.IntLit, string(lit), pos)
	}
}

func (l *Lexer) readIdent(pos token.Pos) {
	var lit []rune
	for l.pos < len(l.src) && isIdentPart(l.peek()) {
		lit = append(lit, l.advance())
	}
	word := string(lit)
	typ := token.LookupIdent(word)
	l.emit(typ, word, pos)
}

func (l *Lexer) readOperator(pos token.Pos) {
	ch := l.advance()
	next := l.peek()

	switch ch {
	case '+':
		switch next {
		case '+':
			l.advance()
			l.emit(token.PlusPlus, "++", pos)
		case '=':
			l.advance()
			l.emit(token.PlusAssign, "+=", pos)
		default:
			l.emit(token.Plus, "+", pos)
		}
	case '-':
		switch next {
		case '-':
			l.advance()
			l.emit(token.MinusMinus, "--", pos)
		case '>':
			l.advance()
			l.emit(token.Arrow, "->", pos)
		case '=':
			l.advance()
			l.emit(token.MinusAssign, "-=", pos)
		default:
			l.emit(token.Minus, "-", pos)
		}
	case '*':
		if next == '=' {
			l.advance()
			l.emit(token.StarAssign, "*=", pos)
		} else {
			l.emit(token.Star, "*", pos)
		}
	case '/':
		if next == '=' {
			l.advance()
			l.emit(token.SlashAssign, "/=", pos)
		} else {
			l.emit(token.Slash, "/", pos)
		}
	case '%':
		l.emit(token.Percent, "%", pos)
	case '=':
		if next == '=' {
			l.advance()
			l.emit(token.Eq, "==", pos)
		} else if next == '>' {
			l.advance()
			l.emit(token.FatArrow, "=>", pos)
		} else {
			l.emit(token.Assign, "=", pos)
		}
	case '!':
		if next == '=' {
			l.advance()
			l.emit(token.Neq, "!=", pos)
		} else {
			l.emit(token.Not, "!", pos)
		}
	case '<':
		if next == '=' {
			l.advance()
			l.emit(token.Lte, "<=", pos)
		} else {
			l.emit(token.Lt, "<", pos)
		}
	case '>':
		if next == '=' {
			l.advance()
			l.emit(token.Gte, ">=", pos)
		} else {
			l.emit(token.Gt, ">", pos)
		}
	case '&':
		if next == '&' {
			l.advance()
			l.emit(token.And, "&&", pos)
		} else {
			l.error(pos, fmt.Sprintf("unexpected character: &"))
		}
	case '|':
		if next == '|' {
			l.advance()
			l.emit(token.Or, "||", pos)
		} else {
			l.error(pos, fmt.Sprintf("unexpected character: |"))
		}
	case '(':
		l.emit(token.LParen, "(", pos)
	case ')':
		l.emit(token.RParen, ")", pos)
	case '{':
		l.emit(token.LBrace, "{", pos)
	case '}':
		l.emit(token.RBrace, "}", pos)
	case '[':
		l.emit(token.LBracket, "[", pos)
	case ']':
		l.emit(token.RBracket, "]", pos)
	case ',':
		l.emit(token.Comma, ",", pos)
	case '.':
		l.emit(token.Dot, ".", pos)
	case ':':
		if next == '=' {
			l.advance()
			l.emit(token.Walrus, ":=", pos)
		} else {
			l.emit(token.Colon, ":", pos)
		}
	case ';':
		l.emit(token.Semicolon, ";", pos)
	default:
		l.error(pos, fmt.Sprintf("unexpected character: %c", ch))
	}
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentStart(ch rune) bool {
	return ch == '_' || unicode.IsLetter(ch)
}

func isIdentPart(ch rune) bool {
	return ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch)
}
