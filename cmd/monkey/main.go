package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Print(repl.MONKEY_FACE_HAPPY)
	fmt.Printf("Hello %v! This is the Monkey Programming Language!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
