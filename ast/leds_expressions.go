package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

// BinaryExpression implements the Expression interface.
type BinaryExpression struct {
	Token           token.Token
	LeftExpression  Expression
	Operator        string
	RightExpression Expression
}

func (be *BinaryExpression) expressionNode() {

}
func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Literal
}

func (be *BinaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(be.LeftExpression.String())
	out.WriteString(" ")
	out.WriteString(be.Operator)
	out.WriteString(" ")
	out.WriteString(be.RightExpression.String())
	out.WriteString(")")
	return out.String()
}

func (be *BinaryExpression) Tree() string {
	children := []treeChild{}
	if be.LeftExpression != nil {
		children = append(children, treeChild{"Left", be.LeftExpression})
	}
	if be.RightExpression != nil {
		children = append(children, treeChild{"Right", be.RightExpression})
	}
	return renderTree(fmt.Sprintf("BinaryExpression %q", be.Operator), children...)
}

// CallExpression implements the Expression interface.
type CallExpression struct {
	Token               token.Token // The '(' Token
	FunctionExpression  Expression  // IdentifierExpression or FunctionExpression
	ArgumentExpressions []Expression
}

func (ce *CallExpression) expressionNode() {

}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, parameter := range ce.ArgumentExpressions {
		params = append(params, parameter.String())
	}
	out.WriteString(ce.FunctionExpression.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}

func (ce *CallExpression) Tree() string {
	children := make([]treeChild, 0, len(ce.ArgumentExpressions)+1)
	if ce.FunctionExpression != nil {
		children = append(children, treeChild{"Function", ce.FunctionExpression})
	}
	for i, argument := range ce.ArgumentExpressions {
		if argument == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("Argument[%d]", i), argument})
	}
	return renderTree("CallExpression", children...)
}
