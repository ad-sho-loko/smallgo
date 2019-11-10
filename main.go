package main

import (
	"fmt"
	"io/ioutil"
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

func readFile(path string) ([]byte, error){
	f, err := ioutil.ReadFile(path)
	if err != nil{
		return nil, err
	}

	return f, err
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	universe := NewScope("__universe", nil)
	universe.DeclType = builtinTypes

	f := os.Args[1]
	isTrace := parseOption(os.Args[2:])

	b, err := readFile(f)
	if err != nil{
		panic(err)
	}

	tokens := NewTokenizer(b).Tokenize()
	ast := NewParser(tokens, isTrace).ParseFile()

	err = WalkAst(ast, universe)

	if err != nil {
		panic(err)
	}

	Gen(ast)
}
