package main

import "fmt"

func emit(s string) {
	fmt.Println("  " + s)
}

func lgen(ast *Ast, e Expr) {
	switch v := e.(type) {
	case *Ident:
		sym, found := ast.CurrentScope.LookUpSymbol(v.Name)
		_assert(found, fmt.Sprintf("lookup failed : %s (scope=%s)", v.Name, ast.CurrentScope.Name))
		emit("mov rax, rbp")
		fmt.Printf("  sub rax, %d\n", sym.Offset)
		emit("push rax")
	default:
		panic("gen.go : invalid lgen")
	}
}

func genExpr(ast *Ast, expr Expr) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *CallFunc:
		for _, e := range e.Args{
			genExpr(ast, e)
			emit("pop rax")
			emit("mov rdi, rax")
		}
		fmt.Printf("  call %s\n", e.FuncName)
		emit("push rax")

	case *Ident:
		lgen(ast, e)
		emit("pop rax")
		emit("mov rax, [rax]")
		emit("push rax")

	case *Lit:
		fmt.Printf("  push %s\n", e.Val)

	case *Binary:
		genExpr(ast, e.Left)
		genExpr(ast, e.Right)
		emit("pop rdi")
		emit("pop rax")
		switch e.Kind {
		case ADD:
			emit("add rax, rdi")
		case SUB:
			emit("sub rax, rdi")
		case MUL:
			emit("imul rax, rdi")
		case DIV:
			emit("cqo")
			emit("idiv rdi")
		case REM:
			emit("cqo")
			emit("idiv rdi")
			emit("mov rax, rdx")
		case EQL:
			emit("cmp rax, rdi")
			emit("sete al")
			emit("movzb rax, al")
		case NEQ:
			emit("cmp rax, rdi")
			emit("setne al")
			emit("movzb rax, al")
		case GTR:
			emit("cmp rax, rdi")
			emit("setg al")
			emit("movzb rax, al")
		case GEQ:
			emit("cmp rax, rdi")
			emit("setge al")
			emit("movzb rax, al")
		case SHR:
			emit("mov cl, dil")
			emit("shr rax, cl")
		case LSS:
			emit("cmp rax, rdi")
			emit("setl al")
			emit("movzb rax, al")
		case LEQ:
			emit("cmp rax, rdi")
			emit("setle al")
			emit("movzb rax, al")
		case SHL:
			emit("mov cl, dil")
			emit("shl rax, cl")
		case OR:
			emit("or rax, rdi")
		case LOR:
			emit("or rax, rdi")
		case AND:
			emit("and rax, rdi")
		case LAND:
			emit("and rax, rdi")
		}
		emit("push rax")
	default:
	}
}

func gen(ast *Ast, n Node) {
	if n == nil {
		return
	}

	_, isExpr := n.(Expr)
	if isExpr {
		panic("gen() must be called in case of n is Expr")
	}

	switch v := n.(type) {
	case *FuncDecl:
		fmt.Printf("%s:\n", v.FuncName.Name)
		emit("push rbp")
		emit("mov rbp, rsp")
		fmt.Printf("  sub rsp, %d\n", ast.TopScope.frameSize() + ast.TopScope.Children[0].frameSize())

		for _, arg := range v.FuncType.Args{
			for _, name := range arg.Names{
				_, found := ast.CurrentScope.LookUpSymbol(name.Name)
				_assert(found, fmt.Sprintf("lookup failed : %s (scope=%s)", name.Name, ast.CurrentScope.Name))
				emit("mov [rsp], rdi")
			}
		}

		gen(ast, v.Body)
		emit("mov rax, 0")
		emit("mov rsp, rbp")
		emit("pop rbp")
		emit("ret")

	case *ReturnStmt:
		for _, e := range v.Exprs {
			genExpr(ast, e)
		}

		emit("pop rax")
		emit("mov rsp, rbp")
		emit("pop rbp")
		emit("ret")

	case *BlockStmt:
		ast.scopeDown()
		for _, stmt := range v.List {
			gen(ast, stmt)
		}
		ast.scopeUp()

	case *IfStmt:
		genExpr(ast, v.Cond)
		emit("pop rax")
		emit("cmp rax, 0")

		l1 := ast.L()
		l2 := ast.L()

		if v.Else == nil {
			emit("je .LEND" + l1)
		} else {
			emit("je .LELSE" + l2)
		}

		gen(ast, v.Then)

		if v.Else != nil {
			emit("jmp .LEND" + l1)
			fmt.Println(".LELSE" + l2 + ":")
			gen(ast, v.Else)
		}

		fmt.Println(".LEND" + l1 + ":")

	case *ForStmt:
		gen(ast, v.Init)
		l1 := ast.L()
		l2 := ast.L()

		emit(".LINIT" + l1 + ":")
		if v.Cond != nil {
			genExpr(ast, v.Cond)
			emit("pop rax")
			emit("cmp rax, 0")
			emit("je .LEND" + l2)
		}
		gen(ast, v.Body)
		gen(ast, v.Post)

		emit("jmp .LINIT" + l1)
		if v.Cond != nil {
			emit(".LEND" + l2 + ":")
		}

	case *AssignStmt:
		for i := range v.Lhs {
			lgen(ast, v.Lhs[i])
			genExpr(ast, v.Rhs[i])
		}

		emit("pop rdi")
		emit("pop rax")
		emit("mov [rax], rdi")

	case *DeclStmt:
		gen(ast, v.Decl)

	case *ExprStmt:
		for _, e := range v.Exprs {
			genExpr(ast, e)
		}

	case *GenDecl:
		for _, spec := range v.Specs {
			gen(ast, spec)
		}

	case *ValueSpec:
		for i, expr := range v.InitValues {
			lgen(ast, v.Names[i])
			genExpr(ast, expr)
			emit("pop rdi")
			emit("pop rax")
			emit("mov [rax], rdi")
		}
	}
}

func Gen(ast *Ast) {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println()

	for _, n := range ast.Nodes {
		gen(ast, n)
	}
}
