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

func (p *Parser) peek() *Token {
	return p.tokens[p.pos]
}

func (p *Parser) isEof() bool {
	return p.peek().Kind == EOF
}

func (p *Parser) is(t TokenKind) bool {
	return !p.isEof() && p.peek().Kind == t
}

func (p *Parser) consume(t TokenKind) bool {
	if !p.is(t) {
		return false
	}
	p.pos++
	return true
}

func (p *Parser) expect(t TokenKind) *Token {
	if p.peek().Kind != t {
		panic(fmt.Sprintf("parse.go : expect %s, but %s", t, p.peek().Kind))
	}
	c := p.peek()
	p.pos++
	return c
}

func (p *Parser) num() Expr {
	tkn := p.expect(NUMBER)
	return &Lit{Kind: tkn.Kind, Val: tkn.Val}
}

func (p *Parser) add() Node {
	n := p.num()
	for ;;{
		if p.consume(ADD) {
			left := n.(Expr)
			right := p.num().(Expr)
			n = &Op{Kind: ADD, Left: left, Right: right}
		} else if p.consume(SUB) {
			left := n.(Expr)
			right := p.num().(Expr)
			n = &Op{Kind: SUB, Left: left, Right: right}
		} else{
			return n
		}
	}
}

func (p *Parser) mul() Node{
	n := p.add()
	for ;;{
		if p.consume(MUL) {
			left := n.(Expr)
			right := p.add().(Expr)
			n = &Op{Kind: MUL, Left: left, Right: right}
		} else if p.consume(DIV) {
			left := n.(Expr)
			right := p.add().(Expr)
			n = &Op{Kind: DIV, Left: left, Right: right}
		} else{
			return n
		}
	}
}

func (p *Parser) Parse() Node {
	return p.mul()
}
