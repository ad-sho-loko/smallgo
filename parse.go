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

func (p *Parser) factor() Expr {
	if p.peek().Kind == NUMBER {
		tkn := p.expect(NUMBER)
		return &Lit{Kind: tkn.Kind, Val: tkn.Val}
	}

	if p.peek().Kind == IDENT {
		tkn := p.expect(IDENT)
		return &Ident{Name: tkn.Val}
	}

	panic(fmt.Sprintf("parse.go : invalid factor %s", p.peek()))
}

func (p *Parser) primary() Expr {
	if p.consume(LPAREN) {
		n := p.expr()
		p.expect(RPAREN)
		return n
	}

	return p.factor()
}

func (p *Parser) unary() Expr {
	if p.consume(ADD) {
	} else if p.consume(SUB) {
		return &Binary{Kind: SUB, Left: &Lit{Kind: NUMBER, Val: "0"}, Right: p.primary()}
	}

	return p.primary()
}

func (p *Parser) mul() Expr {
	n := p.unary()

	for {
		if p.consume(MUL) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: MUL, Left: left, Right: right}
		} else if p.consume(DIV) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: DIV, Left: left, Right: right}
		} else if p.consume(MOD) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: MOD, Left: left, Right: right}
		} else if p.consume(SHL) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: SHL, Left: left, Right: right}
		} else if p.consume(SHR) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: SHR, Left: left, Right: right}
		} else {
			return n
		}
	}
}

func (p *Parser) add() Expr {
	n := p.mul()
	for {
		if p.consume(ADD) {
			left := n.(Expr)
			right := p.mul().(Expr)
			n = &Binary{Kind: ADD, Left: left, Right: right}
		} else if p.consume(SUB) {
			left := n.(Expr)
			right := p.mul().(Expr)
			n = &Binary{Kind: SUB, Left: left, Right: right}
		} else {
			return n
		}
	}
}

func (p *Parser) rel() Expr {
	n := p.add()
	for {
		if p.consume(LSS) {
			left := n.(Expr)
			right := p.add().(Expr)
			n = &Binary{Kind: LSS, Left: left, Right: right}
		} else if p.consume(GTR) {
			left := n.(Expr)
			right := p.add().(Expr)
			n = &Binary{Kind: GTR, Left: left, Right: right}
		} else if p.consume(GEQ) {
			left := n.(Expr)
			right := p.add().(Expr)
			n = &Binary{Kind: GEQ, Left: left, Right: right}
		} else if p.consume(LEQ) {
			left := n.(Expr)
			right := p.add().(Expr)
			n = &Binary{Kind: LEQ, Left: left, Right: right}
		} else {
			return n
		}
	}
}

func (p *Parser) eq() Expr {
	n := p.rel()
	for {
		if p.consume(EQL) {
			left := n.(Expr)
			right := p.rel().(Expr)
			n = &Binary{Kind: EQL, Left: left, Right: right}
		} else if p.consume(NEQ) {
			left := n.(Expr)
			right := p.rel().(Expr)
			n = &Binary{Kind: NEQ, Left: left, Right: right}
		} else {
			return n
		}
	}
}

func (p *Parser) expr() Expr {
	return p.eq()
}

// ex) x = 3
func (p *Parser) assign() *AssignStmt {
	lhs := p.expr()
	p.expect(ASSIGN)
	rhs := p.expr()

	return &AssignStmt{
		Lhs: []Expr{lhs},
		Rhs: []Expr{rhs},
	}
}

func (p *Parser) decl() *DeclStmt {
	ident := p.expr().(*Ident)
	p.expect(ASSIGN)
	return &DeclStmt{
		Decl: &GenDecl{
			Kind: VAR,
			Specs: []Spec{
				&ValueSpec{
					Type:       NewInt(),
					Names:      []*Ident{ident},
					InitValues: []Expr{p.expr()},
				},
			},
		},
	}
}

func (p *Parser) stmt() Stmt {
	if p.consume(RETURN) {
		return &ReturnStmt{Exprs: []Expr{p.expr()}}
	} else if p.consume(VAR) {
		return p.decl()
	}

	panic("parse.go : invalid statement.")
}

func (p *Parser) Parse() *Ast {
	var nodes []Node

	for p.peek().Kind != EOF {
		nodes = append(nodes, p.stmt())
	}

	return &Ast{
		Nodes:   nodes,
		Symbols: make(map[Ident]*Symbol),
	}
}
