package main

import (
	"fmt"
	"monkey/repl"
	"monkey/source"
	"os"
)

func main() {
	switch len(os.Args) {
	case 1:
		err := repl.Run()
		if err != nil {
			panic(err)
		}
	case 2:
		err := source.Run()
		if err != nil {
			panic(err)
		}
	default:
		fmt.Fprintf(os.Stderr, "usage: %s [source file]\n", os.Args[0])
		os.Exit(1)
	}
}
