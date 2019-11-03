package main

import "fmt"

func (ast *Ast) walkExpr(expr Expr) {
	switch e := expr.(type) {

	case *Ident:
		_, found := ast.CurrentScope.LookUpSymbol(e.Name)
		if !found {
			ast.semanticErrors = append(ast.semanticErrors,
				fmt.Errorf("undefined variable : %s", e.Name))
		}

	case *Lit, *Binary, *CallFunc:
	}
}

func (ast *Ast) walkStmt(stmt Stmt) {
	switch s := stmt.(type) {
	case *IfStmt:
		ast.walkExpr(s.Cond)
		ast.walkStmt(s.Then)
		ast.walkStmt(s.Else)
	case *BlockStmt:
		ast.createScope("__blockStmt")
		for _, stmt := range s.List {
			ast.walkStmt(stmt)
		}
		ast.exitScope()

	case *DeclStmt:
		ast.walkNode(s.Decl)

	case *ExprStmt:
		for _, e := range s.Exprs {
			ast.walkExpr(e)
		}

	case *AssignStmt:
		for _, e := range s.Lhs {
			ast.walkExpr(e)
		}

		for _, e := range s.Rhs {
			ast.walkExpr(e)
		}

	case *ReturnStmt:
		for _, e := range s.Exprs {
			ast.walkExpr(e)
		}
	}
}

func (ast *Ast) walkNode(n Node) {
	switch typ := n.(type) {
	case *FuncDecl:
		err := ast.CurrentScope.RegisterSymbol(typ.FuncName.Name, NewFunc())
		if err != nil {
			ast.semanticErrors = append(ast.semanticErrors, err)
		}

		if typ.ReturnTypeIdent != nil && typ.ReturnType == nil {
			t, err := ast.CurrentScope.ResolveType(typ.ReturnTypeIdent.Name)
			if err != nil {
				ast.semanticErrors = append(ast.semanticErrors, err)
			}

			typ.ReturnType = t
		}

		ast.walkExpr(typ.FuncName)
		// ast.walkNode(typ.ReturnType)
		ast.walkStmt(typ.Body)

	case *GenDecl:
		for _, spec := range typ.Specs {
			ast.walkNode(spec)
		}

	case *ValueSpec:
		if typ.Type == nil {
			t, err := ast.CurrentScope.ResolveType(typ.TypeIdent.Name)
			if err != nil {
				ast.semanticErrors = append(ast.semanticErrors, err)
			}
			typ.Type = t
		}

		for _, ident := range typ.Names {
			err := ast.CurrentScope.RegisterSymbol(ident.Name, typ.Type)
			if err != nil {
				ast.semanticErrors = append(ast.semanticErrors, err)
			}
			ast.walkExpr(ident)
		}

		for _, e := range typ.InitValues {
			ast.walkExpr(e)
		}
	}
}

func WalkAst(ast *Ast) error {
	for _, n := range ast.Nodes {
		ast.walkNode(n)
		if len(ast.semanticErrors) > 0 {
			exitErrors(ast.semanticErrors)
		}
	}
	return nil
}
