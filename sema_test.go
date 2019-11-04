package main

import (
	"testing"
)

func makeAst() *Ast {
	universe := NewScope("__universe", nil)
	universe.DeclType = builtinTypes
	return &Ast{
		CurrentScope: universe,
		TopScope:     universe,
	}
}

func TestSema_IdentIsTypeName(t *testing.T) {
	t.Skip()
	ast := makeAst()
	ast.Nodes = []Node{
		&ValueSpec{
			Type:      nil,
			TypeIdent: &Ident{Name: "int"},
			Names:     []*Ident{{Name: "int"}},
		},
	}

	ast.walkNode(ast.Nodes[0])
	if len(ast.semanticErrors) == 0 {
		t.Fail()
	}
}
