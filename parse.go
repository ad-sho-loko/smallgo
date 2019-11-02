package main

import "fmt"

type Parser struct {
	tokens []*Token
	pos    int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) now() *Token {
	return p.tokens[p.pos]
}

func (p *Parser) isEof() bool {
	return p.now().Kind == EOF
}

func (p *Parser) is(t TokenKind) bool {
	return !p.isEof() && p.now().Kind == t
}

func (p *Parser) consume(t TokenKind) bool {
	if !p.is(t) {
		return false
	}
	p.pos++
	return true
}

func (p *Parser) expect(t TokenKind) *Token {
	if p.now().Kind != t {
		panic(fmt.Sprintf("parse.go : expect %s, but %s", t, p.now().Kind))
	}
	c := p.now()
	p.pos++
	return c
}

func (p *Parser) num() Expr {
	tkn := p.expect(NUMBER)
	return &Lit{Kind: tkn.Kind, Val: tkn.Val}
}

func (p *Parser) add() Node {
	n := p.num()

	if p.consume(PLUS) {
		left := n.(*Lit)
		right := p.num().(*Lit)
		n = &Op{Kind: PLUS, Left: left, Right: right}
	}

	return n
}

func (p *Parser) Parse() Node {
	return p.add()
}
