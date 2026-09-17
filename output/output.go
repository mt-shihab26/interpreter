package output

import (
	"io"

	"monkey/ast"
	"monkey/object"
)

type Output struct {
	Writer    io.Writer
	Program   *ast.Program
	Evaluated object.Object
	Verbose   bool
}

func New(writer io.Writer, program *ast.Program, evaluated object.Object, verbose bool) *Output {
	return &Output{
		Writer:    writer,
		Program:   program,
		Evaluated: evaluated,
		Verbose:   verbose,
	}
}

func (o *Output) Print() {
	if o.Verbose {
		o.printVerbose()
	} else {
		o.printNormal()
	}
}

func (o *Output) printNormal() {
	if o.Evaluated != nil {
		io.WriteString(o.Writer, o.Evaluated.Inspect())
	}
}

func (o *Output) printVerbose() {
	io.WriteString(o.Writer, "---CODE---\n")
	io.WriteString(o.Writer, o.Program.Code())
	io.WriteString(o.Writer, "\n---AST---\n")
	io.WriteString(o.Writer, o.Program.Tree())
	io.WriteString(o.Writer, "\n")
	if o.Evaluated != nil {
		io.WriteString(o.Writer, "\n---OUT---\n")
		io.WriteString(o.Writer, o.Evaluated.Inspect())
		io.WriteString(o.Writer, "\n")
	}
}
