package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/lexer"
	"monkey/parser"
)

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(out, "error reading input: %v\n", err)
			}
			return
		}
		line := scanner.Text()
		executeLine(out, line)
	}
}

func executeLine(out io.Writer, line string) {
	lex := lexer.New(line)
	parse := parser.New(lex)
	program := parse.ParseProgram()
	if len(parse.Errors()) != 0 {
		printParseErrors(out, parse.Errors())
		return
	}
	io.WriteString(out, "---\n")
	io.WriteString(out, program.String())
	io.WriteString(out, "\n")
	io.WriteString(out, "---\n")
	io.WriteString(out, program.Tree())
	io.WriteString(out, "\n")
	io.WriteString(out, "---\n")
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
