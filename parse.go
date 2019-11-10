package main

import (
	"fmt"
	"os"
)

type Parser struct {
	tokens []*Token
	pos    int

	// for trace
	trace  bool
	indent int
}

func NewParser(tokens []*Token, isTrace bool) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		trace:  isTrace,
		indent: 0,
	}
}

func printTrace(p *Parser, msg string) {
	for i := p.indent; i > 0; i-- {
		fmt.Fprintf(os.Stderr, " ")
	}

	fmt.Fprintf(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "(%d;%v)", p.pos, p.tokens[p.pos].Kind)
	fmt.Fprintf(os.Stderr, "  %s", p.tokens[p.pos:])
	fmt.Fprintf(os.Stderr, "\n")
}

func trace(p *Parser, msg string) *Parser {
	printTrace(p, msg)
	p.indent++
	return p
}

func un(p *Parser) {
	p.indent--
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

func (p *Parser) readField() *Field {
	var names []*Ident
	var typ Expr

	p.consume(COMMA)
	names = append(names, p.readIdent().(*Ident))
	for p.consume(COMMA) {
		names = append(names, p.readIdent().(*Ident))
	}

	if p.peek().Kind == IDENT {
		typ = p.readType()
	} else {
		typ = &TypeName{Name: names[0].Name}
		names = nil
	}

	return &Field{
		Names: names,
		Type:  typ,
	}
}

func (p *Parser) readFields() []*Field {
	var fields []*Field

	for p.peek().Kind != RPAREN && p.peek().Kind != LBRACE {
		fields = append(fields, p.readField())
	}

	return fields
}

func (p *Parser) readArray() Expr {
	p.expect(LBRACK)
	arraySize := p.expr()
	p.expect(RBRACK)
	typ := p.readType()

	return &Type{
		Kind:      Array,
		PtrOf:     typ,
		ArraySize: arraySize,
	}
}

func (p *Parser) readType() Expr {
	// pointer type
	if p.consume(MUL) {
		return NewPointer(p.readType())
	}

	// array type
	if p.peek().Kind == LBRACK {
		return p.readArray()
	}

	tkn := p.expect(IDENT)
	return &TypeName{Name: tkn.Val}
}

func (p *Parser) readIdent() Expr {
	tkn := p.expect(IDENT)
	return &Ident{Name: tkn.Val}
}

func (p *Parser) factor() Expr {
	if p.trace {
		defer un(trace(p, "factor"))
	}

	if p.peek().Kind == NUMBER {
		tkn := p.expect(NUMBER)
		return &Lit{Kind: tkn.Kind, Val: tkn.Val}
	}

	if p.peek().Kind == CHAR {
		tkn := p.expect(CHAR)
		return &Lit{Kind: CHAR, Val: tkn.Val}
	}

	if p.peek().Kind == STRING{
		tkn := p.expect(STRING)
		return &Lit{Kind:STRING, Val:tkn.Val}
	}

	if p.peek().Kind == IDENT {
		return p.readIdent()
	}

	panic(fmt.Sprintf("parse.go : invalid factor %s", p.peek()))
}

func (p *Parser) primary() Expr {
	if p.trace {
		defer un(trace(p, "primary"))
	}

	if p.consume(LPAREN) {
		n := p.expr()
		p.expect(RPAREN)
		return n
	}

	factor := p.factor()

	// eg) call()
	if p.consume(LPAREN) {
		funcName := factor.(*Ident)
		callFunc := &CallFunc{FuncName: funcName.Name}

		for !p.consume(RPAREN) {
			callFunc.Args = append(callFunc.Args, p.expr())
			p.consume(COMMA)
		}

		return callFunc
	}

	// ex) x[0]
	if p.consume(LBRACK){
		index := p.expr()
		p.expect(RBRACK)

		ident := factor.(*Ident)
		return &IndexExpr{
			X:ident,
			Index:index,
		}
	}


	return factor
}

func (p *Parser) unary() Expr {
	if p.trace {
		defer un(trace(p, "unary"))
	}

	if p.consume(ADD) {
		// nop
	} else if p.consume(SUB) {
		return &Binary{Kind: SUB, Left: &Lit{Kind: NUMBER, Val: "0"}, Right: p.primary()}
	} else if p.consume(MUL) {
		return &StarExpr{X: p.expr()}
	} else if p.consume(AND) {
		return &UnaryExpr{X: p.expr()}
	}

	return p.primary()
}

func (p *Parser) mul() Expr {
	if p.trace {
		defer un(trace(p, "mul"))
	}

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
		} else if p.consume(AND) {
			left := n.(Expr)
			right := p.mul().(Expr)
			n = &Binary{Kind: AND, Left: left, Right: right}
		} else if p.consume(AND_NOT) {
			left := n.(Expr)
			right := p.unary().(Expr)
			n = &Binary{Kind: AND_NOT, Left: left, Right: right}
		} else {
			return n
		}
	}
}

func (p *Parser) add() Expr {
	if p.trace {
		defer un(trace(p, "add"))
	}

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
		} else if p.consume(XOR) {
			left := n.(Expr)
			right := p.mul().(Expr)
			n = &Binary{Kind: XOR, Left: left, Right: right}
		} else {
			return n
		}
	}
}

func (p *Parser) rel() Expr {
	if p.trace {
		defer un(trace(p, "rel"))
	}

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
	if p.trace {
		defer un(trace(p, "eq"))
	}

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
	if p.trace {
		defer un(trace(p, "binary"))
	}

	n := p.eq()
	for {
		if p.consume(LOR) {
			left := n.(Expr)
			right := p.eq().(Expr)
			n = &Binary{Kind: LOR, Left: left, Right: right}
		} else if p.consume(LAND) {
			left := n.(Expr)
			right := p.eq().(Expr)
			n = &Binary{Kind: LAND, Left: left, Right: right}
		} else {
			return n
		}
	}
}

func (p *Parser) expr() Expr {
	if p.trace {
		defer un(trace(p, "expr"))
	}

	return p.binary()
}

// ex) x = 3
func (p *Parser) assign() Stmt {
	if p.trace {
		defer un(trace(p, "assign"))
	}

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
	} else if p.consume(INC) {
		rhs = &Binary{Kind: ADD, Left: lhs, Right: &Lit{Val: "1"}}
	} else if p.consume(DEC) {
		rhs = &Binary{Kind: SUB, Left: lhs, Right: &Lit{Val: "1"}}
	} else if p.consume(OR_ASSIGN) {
		rhs = &Binary{Kind: OR, Left: lhs, Right: p.expr()}
	} else if p.consume(AND_ASSIGN) {
		rhs = &Binary{Kind: AND, Left: lhs, Right: p.expr()}
	} else if p.consume(XOR_ASSIGN) {
		rhs = &Binary{Kind: XOR, Left: lhs, Right: p.expr()}
	} else if p.consume(AND_NOT_ASSIGN) {
		rhs = &Binary{Kind: AND_NOT, Left: lhs, Right: p.expr()}
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
	if p.trace {
		defer un(trace(p, "varDecl"))
	}

	ident := p.readIdent().(*Ident)

	var specs []Spec
	if p.consume(ASSIGN) {
		// ex ) var x = 2
		initValues := p.expr()
		spec := &ValueSpec{
			Type:       &TypeName{Name: "int"},
			Names:      []*Ident{ident},
			InitValues: []Expr{initValues},
		}
		specs = append(specs, spec)
	} else if p.peek().Kind == IDENT || p.peek().Kind == MUL || p.peek().Kind == LBRACK {
		// ex ) var x int, var x *int, var x []int
		typeName := p.readType()

		spec := &ValueSpec{
			Type:  typeName,
			Names: []*Ident{ident},
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

func (p *Parser) forStmt() *ForStmt {
	if p.trace {
		defer un(trace(p, "forStmt"))
	}

	forStmt := &ForStmt{}

	// ex ) for { x = 5 }
	if p.peek().Kind == LBRACE {
		forStmt.Body = p.stmtBlock()
		return forStmt
	}

	if p.peek().Kind != SEMICOLON {
		forStmt.Init = p.simpleStmt()

		// ex ) for i < 10 {}
		e, ok := forStmt.Init.(*ExprStmt)
		if ok && len(e.Exprs) == 1 {
			forStmt.Init = nil
			forStmt.Cond = e.Exprs[0]
			forStmt.Body = p.stmtBlock()
			return forStmt
		}
	}

	p.expect(SEMICOLON)

	if !p.consume(SEMICOLON) {
		forStmt.Cond = p.expr()
		p.expect(SEMICOLON)
	}

	if p.peek().Kind != LBRACE {
		forStmt.Post = p.simpleStmt()
	}

	forStmt.Body = p.stmtBlock()
	return forStmt
}

func (p *Parser) ifStmt() *IfStmt {
	if p.trace {
		defer un(trace(p, "ifStmt"))
	}

	ifStmt := &IfStmt{}
	ifStmt.Cond = p.expr()
	ifStmt.Then = p.stmtBlock()
	if p.consume(ELSE) {
		if p.consume(IF) {
			ifStmt.Else = p.ifStmt()
		} else {
			ifStmt.Else = p.stmtBlock()
		}
	}
	return ifStmt
}

func (p *Parser) simpleStmt() Stmt {
	if p.trace {
		defer un(trace(p, "simpleStmt"))
	}

	if p.consume(RETURN) {
		return &ReturnStmt{Exprs: []Expr{p.expr()}}
	}

	return p.assign()
}

func (p *Parser) stmt() Stmt {
	if p.trace {
		defer un(trace(p, "stmt"))
	}

	if p.consume(IF) {
		return p.ifStmt()
	}

	if p.consume(FOR) {
		return p.forStmt()
	}

	if p.consume(VAR) {
		return p.varDecl()
	}

	return p.simpleStmt()
}

func (p *Parser) stmtBlock() *BlockStmt {
	if p.trace {
		defer un(trace(p, "stmtBlock"))
	}

	p.expect(LBRACE)
	b := &BlockStmt{}
	for !p.consume(RBRACE) {
		b.List = append(b.List, p.stmt())
	}
	return b
}

func (p *Parser) toplevel() *FuncDecl {
	if p.trace {
		defer un(trace(p, "toplevel"))
	}

	funcDecl := FuncDecl{}
	funcDecl.FuncType = &FuncType{}

	if p.consume(FUNC) {
		funcDecl.FuncName = p.readIdent().(*Ident)

		p.expect(LPAREN)
		if p.peek().Kind == IDENT {
			funcDecl.FuncType.Args = p.readFields()
		}
		p.expect(RPAREN)

		p.consume(LPAREN)
		if p.peek().Kind == IDENT {
			funcDecl.FuncType.Returns = p.readFields()
		}
		p.consume(RPAREN)

		funcDecl.Body = p.stmtBlock()
	}

	return &funcDecl
}

func (p *Parser) ParseFile() *Ast {
	if p.trace {
		un(trace(p, "PARSE START"))
	}

	var nodes []Node

	for p.peek().Kind != EOF {
		nodes = append(nodes, p.toplevel())
	}

	return &Ast{
		Nodes: nodes,
	}
}
