package main

func makeAst() *Ast {
	universe := NewScope("__universe", nil)
	universe.DeclType = builtinTypes
	return &Ast{
	}
}
