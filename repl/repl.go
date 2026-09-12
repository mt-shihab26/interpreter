package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "

// Start runs a read-eval-print loop over in, writing each line's result to out until in is exhausted.
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

// executeLine lexes and parses one line, writing its source/tree to out or its parser errors on failure.
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

// printParseErrors writes the sad monkey face followed by each parser error message to out.
func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE_SAD)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
