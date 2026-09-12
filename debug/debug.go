package debug

import (
	"io"

	"monkey/ast"
	"monkey/object"
)

// PrintProgram writes program's source and AST to out, followed by evaluated's inspected value when evaluated is not nil.
func PrintProgram(out io.Writer, program *ast.Program, evaluated object.Object) {
	io.WriteString(out, "---CODE---\n")
	io.WriteString(out, program.String())
	io.WriteString(out, "\n---AST---\n")
	io.WriteString(out, program.Tree())
	io.WriteString(out, "\n")
	if evaluated != nil {
		io.WriteString(out, "\n---OUT---\n")
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}
