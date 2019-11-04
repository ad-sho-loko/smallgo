package main

import "fmt"

var argregs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func emit(s string) {
	fmt.Println("  " + s)
}

func emitf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	fmt.Println("  " + s)
}

func emitfNoIndent(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	fmt.Println(s)
}

func lgen(ast *Ast, e Expr) {
	switch v := e.(type) {
	case *Ident:
		sym, found := ast.CurrentScope.LookUpSymbol(v.Name)
		_assert(found, fmt.Sprintf("lookup failed : %s (scope=%s)", v.Name, ast.CurrentScope.Name))
		emit("mov rax, rbp")
		emitf("sub rax, %d", sym.Offset)
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
		for i, e := range e.Args {
			genExpr(ast, e)
			emit("pop rax")
			emitf("mov %s, rax", argregs[i])
		}
		emitf("call %s", e.FuncName)
		emit("push rax")

	case *Ident:
		lgen(ast, e)
		emit("pop rax")
		emit("mov rax, [rax]")
		emit("push rax")

	case *Lit:
		emitf("push %s", e.Val)

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
		emitfNoIndent("%s:", v.FuncName.Name)
		emit("push rbp")
		emit("mov rbp, rsp")
		emitf("sub rsp, %d", ast.TopScope.frameSize()+ast.TopScope.Children[0].frameSize())
		argNum := 0
		for _, arg := range v.FuncType.Args {
			for _, name := range arg.Names {
				sym, found := ast.CurrentScope.LookUpSymbol(name.Name)
				_assert(found, fmt.Sprintf("lookup failed : %s (scope=%s)", name.Name, ast.CurrentScope.Name))
				emitf("mov [rbp-%d], %s", sym.Offset, argregs[argNum])
				argNum++
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
			emitf("je .LEND%s", l1)
		} else {
			emitf("je .LELSE%s", l2)
		}

		gen(ast, v.Then)

		if v.Else != nil {
			emit("jmp .LEND" + l1)
			emitfNoIndent(".LELSE%s:", l2)
			gen(ast, v.Else)
		}

		emitfNoIndent(".LEND%s:", l1)

	case *ForStmt:
		gen(ast, v.Init)
		l1 := ast.L()
		l2 := ast.L()

		emitfNoIndent(".LINIT%s:", l1)
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
			emitfNoIndent(".LEND%s:", l2)
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
	emitfNoIndent(".intel_syntax noprefix")
	emitfNoIndent(".global main")
	fmt.Println()

	for _, n := range ast.Nodes {
		gen(ast, n)
	}
}
