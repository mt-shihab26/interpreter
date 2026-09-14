package main

import (
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
	}
}
