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

func (p *Parser) primary() Expr{
	if p.consume(LPAREN){
		n := p.expr()
		p.expect(RPAREN)
		return n
	}

	return p.num()
}

func (p *Parser) unary() Expr {
	if p.consume(ADD) {
	} else if p.consume(SUB) {
		return &Binary{Kind:SUB, Left:&Lit{Kind:NUMBER, Val:"0"}, Right:p.primary()}
	}

	return p.primary()
}

func (p *Parser) mul() Expr{
	n := p.unary()

	for ;;{
		if p.consume(MUL) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: MUL, Left: left, Right: right}
		} else if p.consume(DIV) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: DIV, Left: left, Right: right}
		} else{
			return n
		}
	}
}

func (p *Parser) expr() Expr {
	n := p.mul()
	for ;;{
		if p.consume(ADD) {
			left := n.(Expr)
			right := p.mul().(Expr)
			n = &Binary{Kind: ADD, Left: left, Right: right}
		} else if p.consume(SUB) {
			left := n.(Expr)
			right := p.mul().(Expr)
			n = &Binary{Kind: SUB, Left: left, Right: right}
		} else{
			return n
		}
	}
}

func (p *Parser) Parse() Node {
	return p.expr()
}
