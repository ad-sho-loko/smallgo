package main

import "fmt"

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
	PLUS
	EOF
)

var tokenString = map[TokenKind]string{
	NUMBER: "NUMBER",
	PLUS:   "PLUS",
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

func (t *Tokenizer) now() byte {
	return t.b[t.pos]
}

func (t *Tokenizer) isEof() bool {
	return t.pos >= len(t.b)
}

func (t *Tokenizer) isNumeric() bool {
	return t.b[t.pos] >= '0' && t.b[t.pos] <= '9'
}

func (t *Tokenizer) isSpace() bool {
	return t.now() == ' ' || t.now() == '\t'
}

func (t *Tokenizer) skipSpace() {
	for ; !t.isEof() && t.isSpace(); t.pos++ {
	}
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

		if t.now() == '+' {
			tokens = append(tokens, t.newToken(PLUS, ""))
			t.pos++
			continue
		}

		if t.isNumeric() {
			n := t.readNumeric()
			tokens = append(tokens, t.newToken(NUMBER, n))
			continue
		}

		panic(fmt.Sprintf("token.go : invalid charactor %s(%#v)", string(t.now()), t.now()))
	}

	tokens = append(tokens, t.newToken(EOF, ""))
	return tokens
}
