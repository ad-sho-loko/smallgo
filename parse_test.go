package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlus(t *testing.T) {
	tokens := []*Token{
		{Kind: NUMBER, Val: "2"},
		{Kind: PLUS, Val: ""},
		{Kind: NUMBER, Val: "3"},
	}

	p := NewParser(tokens)
	ast := p.Parse()

	op := ast.(*Op)
	assert.Equal(t, PLUS, op.Kind)

	left := op.Left.(*Lit)
	assert.Equal(t, NUMBER, left.Kind)
	assert.Equal(t, "2", left.Val)

	right := op.Right.(*Lit)
	assert.Equal(t, NUMBER, right.Kind)
	assert.Equal(t, "3", right.Val)
}
