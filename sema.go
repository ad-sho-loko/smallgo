package main

import (
	"fmt"
	"strconv"
)

func (ast *Ast) walkExpr(expr Expr) {
	switch e := expr.(type) {

	case *Ident:
		sym, found := ast.CurrentScope.LookUpSymbol(e.Name)
		if !found {
			ast.semanticErrors = append(ast.semanticErrors,
				fmt.Errorf("undefined variable : %s", e.Name))
		}

		if sym.Type.Kind == Array {
			// e._Size = sym.Type.ArraySize
			e._Offset = sym.Offset
		} else {
			e._Size = sym.Type.Size
			e._Offset = sym.Offset
		}

	case *Binary:
		ast.walkExpr(e.Left)
		ast.walkExpr(e.Right)

	case *CallFunc:
		for _, arg := range e.Args {
			ast.walkExpr(arg)
		}

	case *StarExpr:
		ast.walkExpr(e.X)

	case *UnaryExpr:
		ast.walkExpr(e.X)

	case *Lit:
	}
}

func (ast *Ast) walkStmt(stmt Stmt) {
	switch s := stmt.(type) {
	case *IfStmt:
		ast.walkExpr(s.Cond)
		ast.walkStmt(s.Then)
		ast.walkStmt(s.Else)
	case *ForStmt:
		ast.walkStmt(s.Init)
		ast.walkExpr(s.Cond)
		ast.walkStmt(s.Post)
		ast.walkStmt(s.Body)
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

		ast.walkExpr(typ.FuncName)

		for _, arg := range typ.FuncType.Args {
			t, err := ast.CurrentScope.ResolveTop(arg.Type)
			if err != nil {
				ast.semanticErrors = append(ast.semanticErrors, err)
			}

			for _, ident := range arg.Names {
				err = ast.CurrentScope.RegisterSymbol(ident.Name, t)
				if err != nil {
					ast.semanticErrors = append(ast.semanticErrors, err)
				}

				ident._Offset = ast.CurrentScope.frameSize()
				ident._Size = t.Size
			}
		}

		for _, r := range typ.FuncType.Returns {
			r.Type, err = ast.CurrentScope.ResolveTop(r.Type)
			if err != nil {
				ast.semanticErrors = append(ast.semanticErrors, err)
			}
		}

		ast.walkStmt(typ.Body)
		typ._FrameSize = ast.TopScope.frameSize() + ast.TopScope.Children[0].frameSize()

	case *GenDecl:
		for _, spec := range typ.Specs {
			ast.walkNode(spec)
		}

	case *ValueSpec:
		var err error
		typ.Type, err = ast.CurrentScope.ResolveTop(typ.Type)
		if err != nil {
			ast.semanticErrors = append(ast.semanticErrors, err)
		}

		for _, ident := range typ.Names {
			_assert(typ.Type != nil, "type must not be nil here")

			t := typ.Type.(*Type)
			err := ast.CurrentScope.RegisterSymbol(ident.Name, t)
			if err != nil {
				ast.semanticErrors = append(ast.semanticErrors, err)
			}

			if t.Kind == Array {
				sizeLit := t.ArraySize.(*Lit)
				sizeInt, _ := strconv.Atoi(sizeLit.Val)
				ident._Size = sizeInt * t.PtrOf.(*Type).Size
				ident._Offset = ast.CurrentScope.frameSize()
			} else {
				ident._Size = t.Size
				ident._Offset = ast.CurrentScope.frameSize()
			}

			ast.walkExpr(ident) // ....
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
