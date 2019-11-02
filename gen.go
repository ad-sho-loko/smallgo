package main

import "fmt"

func emit(s string) {
	fmt.Println("  " + s)
}

func top(n Node) {
	switch v := n.(type) {
	case *Lit:
		fmt.Printf("  push %s\n", v.Val)
	case *Op:
		top(v.Left)
		top(v.Right)
		emit("pop rax")
		emit("pop rdx")
		emit("add rax, rdx")
		emit("push rax")
	default:
	}
}

func Gen(n Node) {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global main")
	fmt.Println()
	fmt.Println("main:")
	top(n)
	emit("pop rax")
	emit("ret")
}
