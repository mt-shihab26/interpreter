package source

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

func Run() error {
	fileName := os.Args[1]
	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return errors.New("parser errors:\n\t" + strings.Join(p.Errors(), "\n\t"))
	}
	env := object.NewEnvironment()
	evaluated := eval.Eval(program, env)
	fmt.Println(evaluated.Inspect())
	return nil
}
