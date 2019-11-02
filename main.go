package main

import (
	"fmt"
	"os"
)

func usage() {
	title := "Smallgo is a tiny go compiler aimed at self-compiled"
	help := "Usage : smallgo <prog>"
	fmt.Println(title)
	fmt.Println(help)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	t := NewTokenizer([]byte(os.Args[1]))
	tokens := t.Tokenize()
	p := NewParser(tokens)
	n := p.Parse()
	Gen(n)
}
