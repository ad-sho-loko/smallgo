package main

import (
	"testing"
)

func makeAst() *Ast {
	return &Ast{
		Symbols: make(map[Ident]*Symbol),
	}
}

func TestSema_IdentIsTypeName(t *testing.T) {
	ast := makeAst()
	ast.Nodes = []Node{
		&ValueSpec{
			Type:      nil,
			TypeIdent: &Ident{Name: "int"},
			Names:     []*Ident{{Name: "int"}},
		},
	}

	err := WalkAst(ast)
	if err == nil {
		t.Fail()
	}
}