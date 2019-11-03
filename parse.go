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

func (p *Parser) ident() *Ident {
	tkn := p.expect(IDENT)
	return &Ident{Name: tkn.Val}
}

func (p *Parser) factor() Expr {
	if p.peek().Kind == NUMBER {
		tkn := p.expect(NUMBER)
		return &Lit{Kind: tkn.Kind, Val: tkn.Val}
	}

	if p.peek().Kind == IDENT {
		return p.ident()
	}

	panic(fmt.Sprintf("parse.go : invalid factor %s", p.peek()))
}

func (p *Parser) primary() Expr {
	if p.consume(LPAREN) {
		n := p.expr()
		p.expect(RPAREN)
		return n
	}

	factor := p.factor()
	if p.consume(LPAREN) {
		ident := factor.(*Ident)
		p.expect(RPAREN)
		return &CallFunc{FuncName: ident.Name}
	}
	return factor
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
		} else if p.consume(REM) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: REM, Left: left, Right: right}
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
		} else if p.consume(OR) {
			left := n.(Expr)
			right := p.mul().(Expr)
			n = &Binary{Kind: OR, Left: left, Right: right}
		} else if p.consume(AND) {
			left := n.(Expr)
			right := p.mul().(Expr)
			n = &Binary{Kind: AND, Left: left, Right: right}
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

func (p *Parser) binary() Expr {
	n := p.eq()
	for {
		if p.consume(LOR) {
			left := n.(Expr)
			right := p.eq().(Expr)
			n = &Binary{Kind: LOR, Left: left, Right: right}
		} else if p.consume(LAND) {
			left := n.(Expr)
			right := p.eq().(Expr)
			n = &Binary{Kind: LOR, Left: left, Right: right}
		} else {
			return n
		}
	}
}

func (p *Parser) expr() Expr {
	return p.binary()
}

// ex) x = 3
func (p *Parser) assign() Stmt {
	lhs := p.expr()

	var rhs Expr
	if p.consume(ASSIGN) {
		rhs = p.expr()
	} else if p.consume(SHL_ASSIGN) {
		rhs = &Binary{Kind: SHL, Left: lhs, Right: p.expr()}
	} else if p.consume(SHR_ASSIGN) {
		rhs = &Binary{Kind: SHR, Left: lhs, Right: p.expr()}
	} else if p.consume(ADD_ASSIGN) {
		rhs = &Binary{Kind: ADD, Left: lhs, Right: p.expr()}
	} else if p.consume(SUB_ASSIGN) {
		rhs = &Binary{Kind: SUB, Left: lhs, Right: p.expr()}
	} else if p.consume(MUL_ASSIGN) {
		rhs = &Binary{Kind: MUL, Left: lhs, Right: p.expr()}
	} else if p.consume(DIV_ASSIGN) {
		rhs = &Binary{Kind: DIV, Left: lhs, Right: p.expr()}
	} else if p.consume(REM_ASSIGN) {
		rhs = &Binary{Kind: REM, Left: lhs, Right: p.expr()}
	} else if p.consume(OR_ASSIGN) {
		rhs = &Binary{Kind: OR, Left: lhs, Right: p.expr()}
	} else if p.consume(AND_ASSIGN) {
		rhs = &Binary{Kind: AND, Left: lhs, Right: p.expr()}
	} else {
		return &ExprStmt{Exprs: []Expr{lhs}}
	}

	return &AssignStmt{
		Lhs: []Expr{lhs},
		Rhs: []Expr{rhs},
	}
}

// https://golang.org/ref/spec#Variable_declarations
func (p *Parser) varDecl() *DeclStmt {
	ident := p.expr().(*Ident)

	var specs []Spec
	if p.consume(ASSIGN) {
		// ex ) var x = 2
		initValues := p.expr()
		spec := &ValueSpec{
			Type:       NewInt(),
			Names:      []*Ident{ident},
			InitValues: []Expr{initValues},
		}
		specs = append(specs, spec)
	} else if p.peek().Kind == IDENT {
		// ex ) var x int
		typeName := p.expr().(*Ident)
		spec := &ValueSpec{
			TypeIdent: typeName,
			Names:     []*Ident{ident},
		}
		specs = append(specs, spec)
	}

	return &DeclStmt{
		Decl: &GenDecl{
			Kind:  VAR,
			Specs: specs,
		},
	}
}

func (p *Parser) ifStmt() *IfStmt {
	ifStmt := &IfStmt{}
	ifStmt.Cond = p.expr()
	ifStmt.Then = p.stmtBlock()
	return ifStmt
}

func (p *Parser) simpleStmt() Stmt {
	if p.consume(RETURN) {
		return &ReturnStmt{Exprs: []Expr{p.expr()}}
	}

	return p.assign()
}

func (p *Parser) stmt() Stmt {
	if p.consume(IF) {
		return p.ifStmt()
	}

	if p.consume(VAR) {
		return p.varDecl()
	}

	return p.simpleStmt()
}

func (p *Parser) stmtBlock() *BlockStmt {
	p.expect(LBRACE)
	b := &BlockStmt{}
	for !p.consume(RBRACE) {
		b.List = append(b.List, p.stmt())
	}
	return b
}

func (p *Parser) toplevel() *FuncDecl {
	funcDecl := FuncDecl{}

	if p.consume(FUNC) {
		funcDecl.FuncName = p.ident()
		p.expect(LPAREN)
		p.expect(RPAREN)

		if p.peek().Kind == IDENT {
			funcDecl.ReturnTypeIdent = p.expr().(*Ident)
		}

		funcDecl.Body = p.stmtBlock()
	}

	return &funcDecl
}

func (p *Parser) ParseFile(scope *Scope) *Ast {
	var nodes []Node

	for p.peek().Kind != EOF {
		nodes = append(nodes, p.toplevel())
	}

	return &Ast{
		Nodes: nodes,
	}
}
