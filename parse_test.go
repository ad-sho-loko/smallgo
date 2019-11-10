package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertFunc(t *testing.T, got, want *FuncDecl) {
	for i := range want.Body.List {
		assertNodeWalk(t, got.Body.List[i], want.Body.List[i])
	}
}

func assertNodeWalk(t *testing.T, got, want Node) {
	switch n := want.(type) {
	case *Binary:
		gotOp, ok := got.(*Binary)
		assert.True(t, ok)
		assert.Equal(t, gotOp.Kind, n.Kind)
		assertNodeWalk(t, gotOp.Left, n.Left)
		assertNodeWalk(t, gotOp.Right, n.Right)
	case *Lit:
		gotLit, ok := got.(*Lit)
		assert.True(t, ok)
		assert.Equal(t, gotLit.Kind, n.Kind)
		assert.Equal(t, gotLit.Val, n.Val)
	case *ReturnStmt:
		gotStmt, ok := got.(*ReturnStmt)
		assert.True(t, ok)
		for i, e := range n.Exprs {
			assertNodeWalk(t, gotStmt.Exprs[i], e)
		}
	case *ExprStmt:
		gotStmt, ok := got.(*ExprStmt)
		assert.True(t, ok)
		for i, e := range n.Exprs {
			assertNodeWalk(t, gotStmt.Exprs[i], e)
		}

	default:
		t.Fatal("you need to add the types")
	}
}

func TestParse_Add(t *testing.T) {
	test := struct {
		b    []*Token
		want []Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "3"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "2"},
			{Kind: EOF, Val: ""},
		},
		want: []Node{
			&Binary{
				Kind:  ADD,
				Left:  &Lit{Kind: NUMBER, Val: "3"},
				Right: &Lit{Kind: NUMBER, Val: "2"},
			},
		},
	}

	p := NewParser(test.b, false)
	ast := p.expr()
	assertNodeWalk(t, ast, test.want[0])
}

func TestParse_AddPolynomial(t *testing.T) {
	test := struct {
		b    []*Token
		want []Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "3"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "2"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "4"},
			{Kind: EOF, Val: ""},
		},
		want: []Node{
			&Binary{
				Kind: ADD,
				Left: &Binary{
					Kind:  ADD,
					Left:  &Lit{Kind: NUMBER, Val: "3"},
					Right: &Lit{Kind: NUMBER, Val: "2"},
				},
				Right: &Lit{
					Kind: NUMBER,
					Val:  "4",
				},
			},
		},
	}

	p := NewParser(test.b, false)
	ast := p.expr()
	assertNodeWalk(t, ast, test.want[0])
}

func TestParse_Mul(t *testing.T) {
	test := struct {
		b    []*Token
		want []Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "3"},
			{Kind: MUL, Val: ""},
			{Kind: NUMBER, Val: "2"},
			{Kind: EOF, Val: ""},
		},
		want: []Node{
			&Binary{
				Kind:  MUL,
				Left:  &Lit{Kind: NUMBER, Val: "3"},
				Right: &Lit{Kind: NUMBER, Val: "2"},
			},
		},
	}

	p := NewParser(test.b, false)
	ast := p.expr()
	assertNodeWalk(t, ast, test.want[0])
}

func TestParse_Precedence(t *testing.T) {
	test := struct {
		b    []*Token
		want []Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "2"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "3"},
			{Kind: MUL, Val: ""},
			{Kind: NUMBER, Val: "4"},
			{Kind: EOF, Val: ""},
		},
		want: []Node{
			&Binary{
				Kind: ADD,
				Left: &Lit{Kind: NUMBER, Val: "2"},
				Right: &Binary{
					Kind:  MUL,
					Left:  &Lit{Kind: NUMBER, Val: "3"},
					Right: &Lit{Kind: NUMBER, Val: "4"},
				},
			},
		},
	}

	p := NewParser(test.b, false)
	ast := p.expr()
	assertNodeWalk(t, ast, test.want[0])
}

func TestParse_Paren(t *testing.T) {
	test := struct {
		b    []*Token
		want []Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "2"},
			{Kind: MUL, Val: ""},
			{Kind: LPAREN, Val: ""},
			{Kind: NUMBER, Val: "3"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "4"},
			{Kind: RPAREN, Val: ""},
			{Kind: EOF, Val: ""},
		},
		want: []Node{
			&Binary{
				Kind: MUL,
				Left: &Lit{Kind: NUMBER, Val: "2"},
				Right: &Binary{
					Kind:  ADD,
					Left:  &Lit{Kind: NUMBER, Val: "3"},
					Right: &Lit{Kind: NUMBER, Val: "4"},
				},
			},
		},
	}

	p := NewParser(test.b, false)
	ast := p.expr()
	assertNodeWalk(t, ast, test.want[0])
}

func TestParse_ReturnStmt(t *testing.T) {
	test := struct {
		b    []*Token
		want []Node
	}{
		b: []*Token{
			{Kind: RETURN, Val: "2"},
			{Kind: NUMBER, Val: "5"},
			{Kind: EOF, Val: ""},
		},
		want: []Node{
			&ReturnStmt{
				Exprs: []Expr{
					&Lit{Kind: NUMBER, Val: "5"},
				},
			},
		},
	}

	p := NewParser(test.b, false)
	ast := p.stmt()
	assertNodeWalk(t, ast, test.want[0])
}

func TestParse_FuncDecl(t *testing.T) {
	test := struct {
		b    []*Token
		want *FuncDecl
	}{
		b: []*Token{
			{Kind: FUNC, Val: ""},
			{Kind: IDENT, Val: "main"},
			{Kind: LPAREN, Val: ""},
			{Kind: RPAREN, Val: ""},
			{Kind: IDENT, Val: "int"},
			{Kind: LBRACE, Val: ""},
			{Kind: RETURN, Val: ""},
			{Kind: NUMBER, Val: "5"},
			{Kind: RBRACE, Val: ""},
			{Kind: EOF, Val: ""},
		},
		want: &FuncDecl{
			FuncName: &Ident{Name: "main"},
			Body: &BlockStmt{
				List: []Stmt{&ReturnStmt{
					Exprs: []Expr{
						&Lit{Kind: NUMBER, Val: "5"},
					},
				},
				},
			},
		},
	}

	p := NewParser(test.b, false)
	ast := p.funcDecl()
	assertFunc(t, ast, test.want)
}

func TestParse_ExprStmt(t *testing.T) {
	test := struct {
		b    []*Token
		want *ExprStmt
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "5"},
			{Kind: EOF, Val: ""},
		},
		want: &ExprStmt{
			Exprs: []Expr{
				&Lit{Kind: NUMBER, Val: "5"},
			},
		},
	}

	p := NewParser(test.b, false)
	ast := p.stmt()
	assertNodeWalk(t, ast, test.want)
}

func TestParse_ReadField(t *testing.T) {
	// (x int)
	test := struct {
		b    []*Token
		want *Field
	}{
		b: []*Token{
			{Kind: IDENT, Val: "x"},
			{Kind: IDENT, Val: "int"},
			{Kind: EOF, Val: ""},
		},
		want: &Field{
			Names: []*Ident{{Name: "x"}},
			Type:  &TypeName{Name: "int"},
		},
	}
	f := NewParser(test.b, false).readField()
	assert.Equal(t, test.want, f)

	// (x, y int)
	test = struct {
		b    []*Token
		want *Field
	}{
		b: []*Token{
			{Kind: IDENT, Val: "x"},
			{Kind: COMMA, Val: ""},
			{Kind: IDENT, Val: "y"},
			{Kind: IDENT, Val: "int"},
			{Kind: EOF, Val: ""},
		},
		want: &Field{
			Names: []*Ident{{Name: "x"}, {Name: "y"}},
			Type:  &TypeName{Name: "int"},
		},
	}

	f = NewParser(test.b, false).readField()
	assert.Equal(t, test.want, f)

	// (int)
	test = struct {
		b    []*Token
		want *Field
	}{
		b: []*Token{
			{Kind: IDENT, Val: "int"},
			{Kind: EOF, Val: ""},
		},
		want: &Field{
			Names: nil,
			Type:  &TypeName{Name: "int"},
		},
	}

	f = NewParser(test.b, false).readField()
	assert.Equal(t, test.want, f)
}
