package source

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/output"
	"monkey/parser"
)

func Run() error {
	fileName := os.Args[1]
	extension := getExtension(fileName)
	if extension != "mx" {
		return fmt.Errorf("invalid file extension: .%q, expected .mx", extension)
	}
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
	output.PrintProgram(os.Stdout, program, evaluated)
	return nil
}

func getExtension(fileName string) string {
	split := strings.Split(fileName, "/")
	fileName = split[len(split)-1]
	return strings.Split(fileName, ".")[1]
}
