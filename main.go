package main

import (
	"fmt"
	"os"
	"strconv"
)

func usage() {
	title := "Smallgo is a tiny go compiler aimed at self-compiled"
	help := "Usage : smallgo <prog>"
	fmt.Println(title)
	fmt.Println(help)
}

func atoi(v string) int {
	n, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return n
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n\n")
	fmt.Printf("main:\n")
	fmt.Printf("  mov rax, %d\n", atoi(os.Args[1]))
	fmt.Printf("  ret\n")
}
