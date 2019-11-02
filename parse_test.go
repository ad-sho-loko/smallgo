package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func walkAssert(t *testing.T, got, want Node) {
	switch n := want.(type) {
	case *Binary:
		gotOp, ok := got.(*Binary)
		assert.True(t, ok)
		assert.Equal(t, gotOp.Kind, n.Kind)
		walkAssert(t, gotOp.Left, n.Left)
		walkAssert(t, gotOp.Left, n.Left)
	case *Lit:
		gotLit, ok := got.(*Lit)
		assert.True(t, ok)
		assert.Equal(t, gotLit.Kind, n.Kind)
		assert.Equal(t, gotLit.Val, n.Val)
	default:
		t.Fail()
	}
}

func TestParse_Add(t *testing.T) {
	test := struct {
		b    []*Token
		want Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "3"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "2"},
			{Kind:EOF, Val:""},
		},
		want: &Binary{
			Kind:  ADD,
			Left:  &Lit{Kind: NUMBER, Val: "3"},
			Right: &Lit{Kind: NUMBER, Val: "2"},
		},
	}

	p := NewParser(test.b)
	ast := p.Parse()
	walkAssert(t, ast, test.want)
}

func TestParse_AddPolynomial(t *testing.T) {
	test := struct {
		b    []*Token
		want Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "3"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "2"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "4"},
			{Kind:EOF, Val:""},
		},
		want: &Binary{
			Kind:ADD,
			Left:&Binary{
				Kind:  ADD,
				Left:  &Lit{Kind: NUMBER, Val: "3"},
				Right: &Lit{Kind: NUMBER, Val: "2"},
			},
			Right:&Lit{
				Kind:NUMBER,
				Val:"4",
			},
		},
	}

	p := NewParser(test.b)
	ast := p.Parse()
	walkAssert(t, ast, test.want)
}


func TestParse_Mul(t *testing.T) {
	test := struct {
		b    []*Token
		want Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "3"},
			{Kind: MUL, Val: ""},
			{Kind: NUMBER, Val: "2"},
			{Kind:EOF, Val:""},
		},
		want: &Binary{
			Kind:  MUL,
			Left:  &Lit{Kind: NUMBER, Val: "3"},
			Right: &Lit{Kind: NUMBER, Val: "2"},
		},
	}

	p := NewParser(test.b)
	ast := p.Parse()
	walkAssert(t, ast, test.want)
}

func TestParse_Precedence(t *testing.T) {
	test := struct {
		b    []*Token
		want Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "2"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "3"},
			{Kind: MUL, Val: ""},
			{Kind: NUMBER, Val: "4"},
			{Kind:EOF, Val:""},
		},
		want: &Binary{
			Kind:  ADD,
			Left:  &Lit{Kind: NUMBER, Val: "2"},
			Right: &Binary{
				Kind:MUL,
				Left:  &Lit{Kind: NUMBER, Val: "3"},
				Right:  &Lit{Kind: NUMBER, Val: "4"},
			},
		},
	}

	p := NewParser(test.b)
	ast := p.Parse()
	walkAssert(t, ast, test.want)
}


func TestParse_Paren(t *testing.T) {
	test := struct {
		b    []*Token
		want Node
	}{
		b: []*Token{
			{Kind: NUMBER, Val: "2"},
			{Kind: MUL, Val: ""},
			{Kind: LPAREN, Val: ""},
			{Kind: NUMBER, Val: "3"},
			{Kind: ADD, Val: ""},
			{Kind: NUMBER, Val: "4"},
			{Kind: RPAREN, Val: ""},
			{Kind:EOF, Val:""},
		},
		want: &Binary{
			Kind:  MUL,
			Left:  &Lit{Kind: NUMBER, Val: "2"},
			Right: &Binary{
				Kind:ADD,
				Left:  &Lit{Kind: NUMBER, Val: "3"},
				Right:  &Lit{Kind: NUMBER, Val: "4"},
			},
		},
	}

	p := NewParser(test.b)
	ast := p.Parse()
	walkAssert(t, ast, test.want)
}