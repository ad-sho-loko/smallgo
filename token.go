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

type TokenKind uint

const (
	NUMBER TokenKind = iota + 1
	ADD              // +
	SUB              // -
	MUL              // *
	DIV              // /
	MOD              // %
	LPAREN           // (
	RPAREN           // )
	ASSIGN           // =
	EQL              // ==
	NEQ              // !=
	NOT              // !
	LSS              // <
	LEQ              // <=
	GTR              // >
	GEQ              // >=
	VAR              // var
	RETURN           // return
	IDENT
	EOF
)

var tokenString = map[TokenKind]string{
	NUMBER: "NUMBER",
	ADD:    "ADD",
	SUB:    "SUB",
	MUL:    "MUL",
	DIV:    "DIV",
	MOD:    "MOD",
	LPAREN: "LPAREN",
	RPAREN: "RPAREN",
	ASSIGN: "ASSIGN",
	EQL:    "EQL",
	NEQ:    "NEQ",
	NOT:    "NOT",
	LSS:    "LSS",
	LEQ:    "LEQ",
	GTR:    "GTR",
	GEQ:    "GTR",
	VAR:    "VAR",
	RETURN: "RETURN",
	IDENT:  "IDENT",
	EOF:    "EOF",
}

var keywords = map[string]TokenKind{
	"var":    VAR,
	"return": RETURN,
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

func (t *Tokenizer) readString() *Token {
	s := ""

	for ; !t.isEof() && t.isLetter(t.peek()); t.pos++ {
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
			tokens = append(tokens, t.readString())
			continue
		}

		if t.isDigit(t.peek()) {
			tokens = append(tokens, t.readNumeric())
			continue
		}

		switch t.peek() {
		case '+':
			tokens = append(tokens, t.newToken(ADD, ""))
			t.pos++
		case '-':
			tokens = append(tokens, t.newToken(SUB, ""))
			t.pos++
		case '*':
			tokens = append(tokens, t.newToken(MUL, ""))
			t.pos++
		case '/':
			tokens = append(tokens, t.newToken(DIV, ""))
			t.pos++
		case '%':
			tokens = append(tokens, t.newToken(MOD, ""))
			t.pos++
		case '(':
			tokens = append(tokens, t.newToken(LPAREN, ""))
			t.pos++
		case ')':
			tokens = append(tokens, t.newToken(RPAREN, ""))
			t.pos++
		case '=':
			t.pos++
			kind := t.switch2(ASSIGN, EQL)
			tokens = append(tokens, t.newToken(kind, ""))
			t.pos++
		case '<':
			t.pos++
			kind := t.switch2(LSS, LEQ)
			tokens = append(tokens, t.newToken(kind, ""))
		case '>':
			t.pos++
			kind := t.switch2(GTR, GEQ)
			tokens = append(tokens, t.newToken(kind, ""))
		case '!':
			t.pos++
			kind := t.switch2(NOT, NEQ)
			tokens = append(tokens, t.newToken(kind, ""))
			t.pos++
		default:
			panic(fmt.Sprintf("token.go : invalid charactor %s(%#v)", string(t.peek()), t.peek()))
		}
	}

	tokens = append(tokens, t.newToken(EOF, ""))
	return tokens
}
