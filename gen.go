package main

import "fmt"

var argregs = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var argregs8 = []string{"dil", "sil", "dl", "cl", "r8b", "r9b"}

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

// for debug
func d(s string) {
	fmt.Println("# " + s)
}

func reg(i, size int) string {
	if size == 1 {
		return argregs8[i]
	}
	return argregs[i]
}

func ptr(size int) string {
	if size == 1 {
		return "BYTE PTR "
	}
	return ""
}

func lgenExpr(ast *Ast, e Expr) {
	switch v := e.(type) {
	case *Ident:
		emit("mov rax, rbp")
		emitf("sub rax, %d", v._Offset)
		emit("push rax")
	case *IndexExpr:
		lgenExpr(ast, v.X)

	default:
		panic("gen.go : invalid lgenExpr")
	}
}

func genExpr(ast *Ast, expr Expr) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *CallFunc:
		for _, arg := range e.Args {
			genExpr(ast, arg)
		}

		for i := range e.Args {
			emit("pop rax")
			argNum := len(e.Args) - i - 1
			emitf("mov %s, rax", argregs[argNum])
		}

		// We need to align to 16 byte boundary
		l := ast.L()
		emit("mov rax, rsp")
		emit("and rax, 15")
		emitf("jnz .Lcall.align%s", l)
		emit("mov rax, 0")
		emitf("call %s", e.FuncName)
		emitf("jmp .Lcall.end%s", l)
		emitfNoIndent(".Lcall.align%s:", l)
		emit("sub rsp, 8")
		emit("mov rax, 0")
		emitf("call %s", e.FuncName)
		emit("add rsp, 8")
		emitfNoIndent(".Lcall.end%s:", l)
		emit("push rax")

	case *Ident:
		d("# Ident")
		lgenExpr(ast, e)
		emit("pop rax")
		emit("mov rax, [rax]")
		emit("push rax")

	case *Lit:
		if e.Kind == CHAR {
			emitf("push %d", e.Val[0])
		} else if e.Kind == STRING{
			emitf("lea rax, %s[rip]", e.Val)
			emit("push rax")
		} else {
			emitf("push %s", e.Val)
		}

	case *StarExpr:
		d("StarExpr")
		genExpr(ast, e.X)
		emit("pop rax")
		emit("mov rax, [rax]")
		emit("push rax")

	case *UnaryExpr:
		d("UnaryExpr")
		lgenExpr(ast, e.X)

	case *IndexExpr:
		genExpr(ast, e.X)

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
		case XOR:
			emit("xor rax, rdi")
		case AND_NOT:
			emit("not rdi")
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
		emitf("sub rsp, %d", roundup(v._FrameSize, 16))

		argNum := 0
		for _, arg := range v.FuncType.Args {
			for _, name := range arg.Names {
				emitf("mov %s [rbp-%d], %s", ptr(name._Size), name._Offset, reg(argNum, name._Size))
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
		for _, stmt := range v.List {
			gen(ast, stmt)
		}

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
			lgenExpr(ast, v.Lhs[i])
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
			lgenExpr(ast, v.Names[i])
			genExpr(ast, expr)
			emit("pop rdi")
			emit("pop rax")
			emit("mov [rax], rdi")
		}
	}
}

func genStrings(ast *Ast){
	for i, s := range ast.strings{
		emitfNoIndent(ast.stringLabels[i])
		emitf(".string \"%s\"", s)
	}
}

func Gen(ast *Ast) {
	emitfNoIndent(".intel_syntax noprefix")
	emitfNoIndent(".global main")
	genStrings(ast)
	fmt.Println()

	for _, n := range ast.Nodes {
		gen(ast, n)
	}
}
