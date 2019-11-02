package main

import "fmt"

func emit(s string) {
	fmt.Println("  " + s)
}

func gen(n Node) {
	switch v := n.(type) {
	case *Lit:
		fmt.Printf("  push %s\n", v.Val)
	case *Binary:
		gen(v.Left)
		gen(v.Right)

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
			emit("cqo\n")
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

func Gen(n Node) {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println()
	fmt.Println("main:")
	gen(n)
	emit("pop rax")
	emit("ret")
}
