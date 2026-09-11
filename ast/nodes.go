package ast

import (
	"bytes"
	"fmt"
)

// Program implements the Node interface.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) Tree() string {
	children := make([]treeChild, 0, len(p.Statements))
	for i, statement := range p.Statements {
		if statement == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("[%d]", i), statement})
	}
	return renderTree("Program", children...)
}
