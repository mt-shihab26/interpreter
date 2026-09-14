package main

import (
	"monkey/repl"
	"monkey/source"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		err := source.Run()
		if err != nil {
			panic(err)
		}
		return
	}
	err := repl.Run()
	if err != nil {
		panic(err)
	}
}
