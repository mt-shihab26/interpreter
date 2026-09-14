package source

import (
	"fmt"
	"os"
)

func Run() error {
	filename := os.Args[1]
	fmt.Printf("%+v\n", filename)
	return nil
}
