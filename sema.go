package main

func walkIdentList(ast *Ast, idents []*Ident) {
	for _, ident := range idents {
		walkNode(ast, ident)
	}
}

func walkNode(ast *Ast, n Node) {
	switch typ := n.(type) {
	case *DeclStmt:
		walkNode(ast, typ.Decl)

	case *GenDecl:
		for _, spec := range typ.Specs {
			walkNode(ast, spec)
		}

	case *ValueSpec:
		for _, ident := range typ.Names {
			ast.Symbols[*ident] = typ.Type
		}
		walkIdentList(ast, typ.Names)
		walkNode(ast, typ.InitValues)

	case *Ident:
	default:
	}
}

func WalkAst(ast *Ast) {
	for _, n := range ast.Nodes {
		walkNode(ast, n)
	}
}
