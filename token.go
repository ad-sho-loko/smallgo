package main

import (
	"fmt"
)

type Tokenizer struct {
	b   []byte
	pos int
}

type Token struct {
	Kind TokenKind
	Val  string
}

func (t Token) String() string {
	if t.Val == "" {
		return fmt.Sprintf("(%s)", tokenString[t.Kind])
	}
	return fmt.Sprintf("(%s `%s`)", tokenString[t.Kind], t.Val)
}

type TokenKind uint

const (
	NUMBER         TokenKind = iota + 1
	ADD                      // +
	SUB                      // -
	MUL                      // *
	DIV                      // /
	REM                      // %
	OR                       // |
	XOR                      // ^
	OR_ASSIGN                // |=
	AND                      // &
	AND_ASSIGN               // &=
	ADD_ASSIGN               // +=
	SUB_ASSIGN               // -=
	MUL_ASSIGN               // *=
	DIV_ASSIGN               // /=
	REM_ASSIGN               // %=
	XOR_ASSIGN               // ^=
	INC                      // ++
	DEC                      // --
	LPAREN                   // (
	RPAREN                   // )
	LBRACE                   // {
	RBRACE                   // }
	LBRACK                   // [
	RBRACK                   // ]
	COMMA                    // ,
	SEMICOLON                // ;
	ASSIGN                   // =
	EQL                      // ==
	SHL                      // <<
	SHL_ASSIGN               // <<=
	SHR                      // >>
	SHR_ASSIGN               // >>=
	NEQ                      // !=
	NOT                      // !
	LOR                      // ||
	LAND                     // &&
	AND_NOT                  // &^
	AND_NOT_ASSIGN           // &^=
	LSS                      // <
	LEQ                      // <=
	GTR                      // >
	GEQ                      // >=
	VAR                      // var
	RETURN                   // return
	FUNC                     // func
	IF                       // if
	ELSE                     // else
	FOR                      // for
	IDENT
	CHAR
	STRING
	EOF
)

var tokenString = map[TokenKind]string{
	NUMBER:         "NUMBER",
	ADD:            "ADD",
	SUB:            "SUB",
	MUL:            "MUL",
	DIV:            "DIV",
	REM:            "REM",
	XOR:            "XOR",
	OR:             "OR",
	OR_ASSIGN:      "OR_ASSIGN",
	AND:            "AND",
	AND_ASSIGN:     "AND_ASSIGN",
	ADD_ASSIGN:     "ADD_ASSIGN",
	SUB_ASSIGN:     "SUB_ASSIGN",
	MUL_ASSIGN:     "MUL_ASSIGN",
	DIV_ASSIGN:     "QUO_ASSIGN",
	REM_ASSIGN:     "REM_ASSIGN",
	XOR_ASSIGN:     "XOR_ASSIGN",
	INC:            "INC",
	DEC:            "DEC",
	LPAREN:         "LPAREN",
	RPAREN:         "RPAREN",
	LBRACE:         "LBRACE",
	RBRACE:         "RBRACE",
	LBRACK:         "LBRACK",
	RBRACK:         "RBRACK",
	COMMA:          "COMMA",
	SEMICOLON:      "SEMICOLON",
	ASSIGN:         "ASSIGN",
	EQL:            "EQL",
	SHL:            "SHL",
	SHL_ASSIGN:     "SHL_ASSIGN",
	SHR:            "SHR",
	SHR_ASSIGN:     "SHR_ASSIGN",
	NEQ:            "NEQ",
	NOT:            "NOT",
	LOR:            "LOR",
	LAND:           "LAND",
	AND_NOT:        "AND_NOT",
	AND_NOT_ASSIGN: "AND_NOT_ASSIGN",
	LSS:            "LSS",
	LEQ:            "LEQ",
	GTR:            "GTR",
	GEQ:            "GTR",
	VAR:            "VAR",
	RETURN:         "RETURN",
	FUNC:           "FUNC",
	IF:             "IF",
	ELSE:           "ELSE",
	FOR:            "FOR",
	IDENT:          "IDENT",
	CHAR:           "CHAR",
	STRING:         "STRING",
	EOF:            "EOF",
}

var keywords = map[string]TokenKind{
	"var":    VAR,
	"return": RETURN,
	"func":   FUNC,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
}

func (t TokenKind) String() string {
	v, ok := tokenString[t]
	if !ok {
		panic(fmt.Sprint("failed to stringfy the token"))
	}
	return v
}

func NewTokenizer(b []byte) *Tokenizer {
	return &Tokenizer{
		b:   b,
		pos: 0,
	}
}

func (t *Tokenizer) newToken(kind TokenKind, val string) *Token {
	return &Token{
		Kind: kind,
		Val:  val,
	}
}

func (t *Tokenizer) peek() byte {
	return t.b[t.pos]
}

func (t *Tokenizer) isEof() bool {
	return t.pos >= len(t.b)
}

func (t *Tokenizer) isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (t *Tokenizer) isKeyword(str string) bool {
	_, found := keywords[str]
	return found
}

func (t *Tokenizer) isIdentifer(str string) bool {
	for i, ch := range []byte(str) {
		if !t.isLetter(ch) && ch != '_' && (i == 0 && !t.isDigit(ch)) {
			return false
		}
	}

	return str != "" && !t.isKeyword(str)
}

func (t *Tokenizer) isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (t *Tokenizer) isSpace() bool {
	return t.peek() == ' ' || t.peek() == '\t' || t.peek() == '\r'
}

func (t *Tokenizer) skipSpace() {
	for ; !t.isEof() && t.isSpace(); t.pos++ {
	}
}

func (t *Tokenizer) switch2(kind1, kind2 TokenKind) TokenKind {
	if !t.isEof() && t.peek() == '=' {
		t.pos++
		return kind2
	}

	return kind1
}

func (t *Tokenizer) switch3(ch byte, kind1, kind2, kind3 TokenKind) TokenKind {
	if t.isEof() {
		return kind1
	}

	if t.peek() == '=' {
		t.pos++
		return kind2
	}

	if t.peek() == ch {
		t.pos++
		return kind3
	}

	return kind1
}

func (t *Tokenizer) switch4(ch byte, kind1, kind2, kind3, kind4 TokenKind) TokenKind {
	if t.isEof() {
		return kind1
	}

	if t.peek() == '=' {
		t.pos++
		return kind2
	}

	if t.peek() == ch {
		t.pos++
		if !t.isEof() && t.peek() == '=' {

			t.pos++
			return kind4
		}
		return kind3
	}

	return kind1
}

func (t *Tokenizer) readTypeOrIdent() *Token {
	s := ""

	for ; !t.isEof() && (t.isLetter(t.peek()) || t.isDigit(t.peek())); t.pos++ {
		s += string(t.peek())
	}

	if t.isKeyword(s) {
		keyword := keywords[s]
		return t.newToken(keyword, "")
	}

	if t.isIdentifer(s) {
		return t.newToken(IDENT, s)
	}

	panic("token.go : invalid string")
}

func (t *Tokenizer) readAsciiString() *Token{
	s := ""

	for !t.isEof() && t.peek() != '"'{
		s += string(t.peek())
		t.pos++
	}

	return &Token{Kind:STRING, Val:s}
}

func (t *Tokenizer) readNumeric() *Token {
	n := ""

	for ; !t.isEof() && t.isDigit(t.peek()); t.pos++ {
		n += string(t.peek())
	}

	return t.newToken(NUMBER, n)
}

func (t *Tokenizer) Tokenize() []*Token {
	var tokens []*Token

	for !t.isEof() {

		t.skipSpace()
		if t.isEof() {
			break
		}

		if t.isLetter(t.peek()) {
			tokens = append(tokens, t.readTypeOrIdent())
			continue
		}

		if t.isDigit(t.peek()) {
			tokens = append(tokens, t.readNumeric())
			continue
		}

		switch t.peek() {
		case '+':
			t.pos++
			kind := t.switch3('+', ADD, ADD_ASSIGN, INC)
			tokens = append(tokens, t.newToken(kind, ""))
		case '-':
			t.pos++
			kind := t.switch3('-', SUB, SUB_ASSIGN, DEC)
			tokens = append(tokens, t.newToken(kind, ""))
		case '*':
			t.pos++
			kind := t.switch2(MUL, MUL_ASSIGN)
			tokens = append(tokens, t.newToken(kind, ""))
		case '/':
			t.pos++
			kind := t.switch2(DIV, DIV_ASSIGN)
			tokens = append(tokens, t.newToken(kind, ""))
		case '%':
			t.pos++
			kind := t.switch2(REM, REM_ASSIGN)
			tokens = append(tokens, t.newToken(kind, ""))
		case '(':
			tokens = append(tokens, t.newToken(LPAREN, ""))
			t.pos++
		case ')':
			tokens = append(tokens, t.newToken(RPAREN, ""))
			t.pos++
		case '{':
			tokens = append(tokens, t.newToken(LBRACE, ""))
			t.pos++
		case '}':
			tokens = append(tokens, t.newToken(RBRACE, ""))
			t.pos++
		case '[':
			tokens = append(tokens, t.newToken(LBRACK, ""))
			t.pos++
		case ']':
			tokens = append(tokens, t.newToken(RBRACK, ""))
			t.pos++
		case '=':
			t.pos++
			kind := t.switch2(ASSIGN, EQL)
			tokens = append(tokens, t.newToken(kind, ""))
			t.pos++
		case '<':
			t.pos++
			kind := t.switch4('<', LSS, LEQ, SHL, SHL_ASSIGN)
			tokens = append(tokens, t.newToken(kind, ""))
		case '>':
			t.pos++
			kind := t.switch4('>', GTR, GEQ, SHR, SHR_ASSIGN)
			tokens = append(tokens, t.newToken(kind, ""))
		case '!':
			t.pos++
			kind := t.switch2(NOT, NEQ)
			tokens = append(tokens, t.newToken(kind, ""))
		case '|':
			t.pos++
			kind := t.switch3('|', OR, OR_ASSIGN, LOR)
			tokens = append(tokens, t.newToken(kind, ""))
		case '^':
			t.pos++
			kind := t.switch2(XOR, XOR_ASSIGN)
			tokens = append(tokens, t.newToken(kind, ""))
		case '&':
			t.pos++
			var kind TokenKind
			if !t.isEof() && t.peek() == '^' {
				t.pos++
				kind = t.switch2(AND_NOT, AND_NOT_ASSIGN)
			} else {
				kind = t.switch3('&', AND, AND_ASSIGN, LAND)
			}
			tokens = append(tokens, t.newToken(kind, ""))
		case ',':
			tokens = append(tokens, t.newToken(COMMA, ""))
			t.pos++
		case ';':
			tokens = append(tokens, t.newToken(SEMICOLON, ""))
			t.pos++
		case '\'':
			t.pos++
			tokens = append(tokens, t.newToken(CHAR, string(t.peek())))
			t.pos++
			if t.peek() != '\'' {
				panic("' not closed")
			}
			t.pos++
		case '"':
			t.pos++
			tokens = append(tokens, t.readAsciiString())
			t.pos++

		default:
			panic(fmt.Sprintf("token.go : invalid charactor %s(%#v)", string(t.peek()), t.peek()))
		}
	}

	tokens = append(tokens, t.newToken(EOF, ""))
	return tokens
}
