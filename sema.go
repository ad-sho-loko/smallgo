package main

//noinspection ALL
func walkNode(ast *Ast, n Node) error {
	switch typ := n.(type) {
	case *FuncDecl:
		if typ.ReturnTypeIdent != nil && typ.ReturnType == nil {
			t, err := ast.Scope.ResolveType(typ.ReturnTypeIdent.Name)
			if err != nil {
				return err
			}
			typ.ReturnType = t
		}
		walkNode(ast, typ.FuncName)
		walkNode(ast, typ.ReturnType)

		for _, b := range typ.Body.List {
			walkNode(ast, b)
		}

	case *DeclStmt:
		walkNode(ast, typ.Decl)

	case *GenDecl:
		for _, spec := range typ.Specs {
			walkNode(ast, spec)
		}

	case *ValueSpec:
		if typ.Type == nil {
			t, err := ast.Scope.ResolveType(typ.TypeIdent.Name)
			if err != nil {
				return err
			}
			typ.Type = t
		}

		for _, ident := range typ.Names {
			err := ast.Scope.RegisterSymbol(ident.Name, typ.Type)
			if err != nil {
				return err
			}
			walkNode(ast, ident)
		}

		walkNode(ast, typ.InitValues)
	case *Ident:
		// nop
	default:
	}

	return nil
}

func WalkAst(ast *Ast) error {
	for _, n := range ast.Nodes {
		err := walkNode(ast, n)
		if err != nil {
			return err
		}
	}

	return nil
}
