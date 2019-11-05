package main

type Visiter interface{
	VisitNode(node Node)
	VisitStmt(stmt Stmt)
	VisitExpr(expr Expr)
	LeaveNode(node Node)
	LeaveStmt(stmt Stmt)
}

func walkExpr(v Visiter, expr Expr) {
	v.VisitExpr(expr)

	switch e := expr.(type) {

	case *Ident:

	case *Binary:
		walkExpr(v, e.Left)
		walkExpr(v, e.Right)

	case *FuncType:

	case *CallFunc:
		for _, arg := range e.Args {
			walkExpr(v, arg)
		}

	case *StarExpr:
		walkExpr(v, e.X)

	case *UnaryExpr:
		walkExpr(v, e.X)

	case *Lit:
	}
}

func walkStmt(v Visiter, stmt Stmt) {
	v.VisitStmt(stmt)

	switch s := stmt.(type) {

	case *IfStmt:
		walkExpr(v, s.Cond)
		walkStmt(v, s.Then)
		walkStmt(v, s.Else)

	case *ForStmt:
		walkStmt(v, s.Init)
		walkExpr(v, s.Cond)
		walkStmt(v, s.Post)
		walkStmt(v, s.Body)

	case *BlockStmt:
		for _, stmt := range s.List {
			walkStmt(v, stmt)
		}

	case *DeclStmt:
		walkNode(v, s.Decl)

	case *ExprStmt:
		for _, e := range s.Exprs {
			walkExpr(v, e)
		}

	case *AssignStmt:
		for _, e := range s.Lhs {
			walkExpr(v, e)
		}

		for _, e := range s.Rhs {
			walkExpr(v, e)
		}

	case *ReturnStmt:
		for _, e := range s.Exprs {
			walkExpr(v, e)
		}
	}

	v.LeaveStmt(stmt)
}

func walkNode(v Visiter, n Node) {
	v.VisitNode(n)

	switch typ := n.(type) {
	case *FuncDecl:
		walkExpr(v, typ.FuncName)
		walkExpr(v, typ.FuncType)
		walkStmt(v, typ.Body)

	case *GenDecl:
		for _, spec := range typ.Specs {
			walkNode(v, spec)
		}

	case *ValueSpec:
		for _, ident := range typ.Names {
			walkExpr(v, ident)
		}

		for _, e := range typ.InitValues {
			walkExpr(v,e)
		}
	}

	v.LeaveNode(n)
}

func WalkAst(ast *Ast, scope *Scope) error {
	// Phase1. Register symbols & Resolve types.
	for _, n := range ast.Nodes {
		resolver := &Resolver{
			ast:ast,
			TopScope:scope,
			CurrentScope:scope,
		}

		walkNode(resolver, n)
		if len(ast.semanticErrors) > 0 {
			exitErrors(ast.semanticErrors)
		}
	}

	return nil
}
