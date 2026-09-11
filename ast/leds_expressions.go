package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

// BinaryExpression implements the Expression interface.
type BinaryExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {

}
func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Literal
}

func (be *BinaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(" ")
	out.WriteString(be.Operator)
	out.WriteString(" ")
	out.WriteString(be.Right.String())
	out.WriteString(")")
	return out.String()
}

func (be *BinaryExpression) Tree() string {
	children := []treeChild{}
	if be.Left != nil {
		children = append(children, treeChild{"Left", be.Left})
	}
	if be.Right != nil {
		children = append(children, treeChild{"Right", be.Right})
	}
	return renderTree(fmt.Sprintf("BinaryExpression %q", be.Operator), children...)
}

// CallExpression implements the Expression interface.
type CallExpression struct {
	Token     token.Token // The '(' Token
	Function  Expression  // IdentifierExpression or FunctionExpression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {

}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, parameter := range ce.Arguments {
		params = append(params, parameter.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}

func (ce *CallExpression) Tree() string {
	children := make([]treeChild, 0, len(ce.Arguments)+1)
	if ce.Function != nil {
		children = append(children, treeChild{"Function", ce.Function})
	}
	for i, argument := range ce.Arguments {
		if argument == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("Argument[%d]", i), argument})
	}
	return renderTree("CallExpression", children...)
}
