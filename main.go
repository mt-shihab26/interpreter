package main

import (
	"monkey/repl"
)

func main() {
	err := repl.Run()
	if err != nil {
		panic(err)
	}
}
