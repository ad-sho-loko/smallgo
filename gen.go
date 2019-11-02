package main

import "fmt"

func emit(s string) {
	fmt.Println("  " + s)
}

func lgen(ast *Ast, n Node){
	switch v := n.(type) {
	case *Ident:
		emit("mov rax, rbp")
		fmt.Printf("  sub rax, %d\n", ast.Symbols[*v].Size)
		emit("push rax")
	}
}

func gen(ast *Ast, n Node) {
	switch v := n.(type) {
	case *ReturnStmt:
		for _, e := range v.Exprs {
			gen(ast, e)
		}
		emit("pop rax")
		emit("mov rsp, rbp")
		emit("pop rbp")
		emit("ret")

	case *AssignStmt:
		lgen(ast, v.Lhs)
		gen(ast, v.Rhs)
		emit("pop rdi")
		emit("pop rax")
		emit("mov [rax], rdi")

	case *DeclStmt:
		gen(ast, v.Decl)

	case *GenDecl:
		for _, spec := range v.Specs {
			gen(ast, spec)
		}

	case *ValueSpec:
		for i, expr := range v.InitValues {
			lgen(ast, v.Names[i])
			gen(ast, expr)
			emit("pop rdi")
			emit("pop rax")
			emit("mov [rax], rdi")
		}

	case *Ident:
		lgen(ast, v)
		emit("pop rax")
		emit("mov rax, [rax]")
		emit("push rax")

	case *Lit:
		fmt.Printf("  push %s\n", v.Val)

	case *Binary:
		gen(ast, v.Left)
		gen(ast, v.Right)
		emit("pop rdi")
		emit("pop rax")
		switch v.Kind {
		case ADD:
			emit("add rax, rdi")
		case SUB:
			emit("sub rax, rdi")
		case MUL:
			emit("imul rax, rdi")
		case DIV:
			emit("cqo")
			emit("idiv rdi")
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
		case LSS:
			emit("cmp rax, rdi")
			emit("setl al")
			emit("movzb rax, al")
		case LEQ:
			emit("cmp rax, rdi")
			emit("setle al")
			emit("movzb rax, al")
		}
		emit("push rax")
	default:
	}
}

func Gen(ast *Ast) {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println()
	fmt.Println("main:")
	fmt.Println("  push rbp")
	fmt.Println("  mov rbp, rsp")
	fmt.Printf("  sub rsp, %d\n", ast.FrameSize())

	for _, n := range ast.Nodes {
		gen(ast, n)
	}
}
