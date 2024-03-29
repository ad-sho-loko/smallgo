package main

import (
	"fmt"
	"os"
)

func roundup(x, align int) int {
	return (x + align - 1) & ^(align - 1)
}

func _assert(b bool, msg string) {
	if !b {
		_, _ = fmt.Fprintf(os.Stderr, "[Assertion Failed]\n")
		_, _ = fmt.Fprintf(os.Stderr, msg+"\n")
		os.Exit(1)
	}
}

func exitErrors(errs []error) {
	for _, err := range errs {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
