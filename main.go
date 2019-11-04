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

func parseOption(args []string) bool {
	if len(args) == 0 {
		return false
	}

	if args[0] == "--debug" {
		return true
	}

	return false
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	universe := NewScope("__universe", nil)
	universe.DeclType = builtinTypes

	isTrace := parseOption(os.Args[2:])

	tokens := NewTokenizer([]byte(os.Args[1])).Tokenize()
	ast := NewParser(tokens, isTrace).ParseFile(universe)
	ast.CurrentScope = universe
	ast.TopScope = universe

	err := WalkAst(ast)

	if err != nil {
		panic(err)
	}

	Gen(ast)
}
