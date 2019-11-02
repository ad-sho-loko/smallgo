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
	EOF
)

var tokenString = map[TokenKind]string{
	NUMBER: "NUMBER",
	ADD:    "ADD",
	SUB:    "SUB",
	MUL:    "MUL",
	DIV:    "DIV",
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
	EOF:    "EOF",
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

func (t *Tokenizer) isNumeric() bool {
	return t.b[t.pos] >= '0' && t.b[t.pos] <= '9'
}

func (t *Tokenizer) isSpace() bool {
	return t.peek() == ' ' || t.peek() == '\t'
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

func (t *Tokenizer) readNumeric() string {
	n := ""

	for ; !t.isEof() && t.isNumeric(); t.pos++ {
		n += string(t.b[t.pos])
	}

	return n
}

func (t *Tokenizer) Tokenize() []*Token {
	var tokens []*Token

	for !t.isEof() {

		t.skipSpace()
		if t.isEof() {
			break
		}

		if t.isNumeric() {
			n := t.readNumeric()
			tokens = append(tokens, t.newToken(NUMBER, n))
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
