package ast

import (
	"bytes"
	"fmt"
)

// Program implements the Node interface.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the literal of the first statement's token, or "" if the program is empty.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String reconstructs the whole program as Monkey source code by concatenating each statement's String().
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Tree renders the program as a tree with one indexed child per top-level statement.
func (p *Program) Tree() string {
	children := make([]treeChild, 0, len(p.Statements))
	for i, statement := range p.Statements {
		if statement == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("Statement[%d]", i), statement})
	}
	return renderTree("Program", children...)
}
