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

// expressionNode marks BinaryExpression as an ast.Expression.
func (be *BinaryExpression) expressionNode() {

}

// TokenLiteral returns the operator token's literal, e.g. "+".
func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Literal
}

// Code reconstructs the expression as "(<left> <operator> <right>)", fully parenthesized.
func (be *BinaryExpression) Code() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(be.LeftExpression.Code())
	out.WriteString(" ")
	out.WriteString(be.Operator)
	out.WriteString(" ")
	out.WriteString(be.RightExpression.Code())
	out.WriteString(")")
	return out.String()
}

// Tree renders the expression with the operator in its header and its two operands as "Left"/"Right" children.
func (be *BinaryExpression) Tree() string {
	children := []treeChild{}
	if be.LeftExpression != nil {
		children = append(children, treeChild{"LeftExpression", be.LeftExpression})
	}
	if be.RightExpression != nil {
		children = append(children, treeChild{"RightExpression", be.RightExpression})
	}
	return renderTree(fmt.Sprintf("BinaryExpression %q", be.Operator), children...)
}

// CallExpression implements the Expression interface.
type CallExpression struct {
	Token               token.Token // The '(' Token
	FunctionExpression  Expression  // IdentifierExpression or FunctionExpression
	ArgumentExpressions []Expression
}

// expressionNode marks CallExpression as an ast.Expression.
func (ce *CallExpression) expressionNode() {

}

// TokenLiteral returns the literal of the opening "(" token.
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// Code reconstructs the expression as "<callee>(<arg>, <arg>, ...)".
func (ce *CallExpression) Code() string {
	var out bytes.Buffer
	params := []string{}
	for _, parameter := range ce.ArgumentExpressions {
		params = append(params, parameter.Code())
	}
	out.WriteString(ce.FunctionExpression.Code())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}

// Tree renders the expression with the callee as a "Function" child and each argument as an indexed "Argument[i]" child.
func (ce *CallExpression) Tree() string {
	children := make([]treeChild, 0, len(ce.ArgumentExpressions)+1)
	if ce.FunctionExpression != nil {
		children = append(children, treeChild{"FunctionExpression", ce.FunctionExpression})
	}
	for i, argument := range ce.ArgumentExpressions {
		if argument == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("ArgumentExpression[%d]", i), argument})
	}
	return renderTree("CallExpression", children...)
}

// IndexExpression implements the Expression interface.
type IndexExpression struct {
	Token            token.Token
	LeftExpression   Expression
	NumberExpression Expression
}

// expressionNode marks IndexExpression as an ast.Expression.
func (ce *IndexExpression) expressionNode() {

}

// TokenLiteral returns the literal of the opening "(" token.
func (ce *IndexExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// Code reconstructs the expression as "<callee>(<arg>, <arg>, ...)".
func (ce *IndexExpression) Code() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ce.LeftExpression.Code())
	out.WriteString("[")
	out.WriteString(ce.NumberExpression.Code())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}

// Tree renders the index expression with the left expression and index expression as children.
func (ie *IndexExpression) Tree() string {
	children := make([]treeChild, 0, 2)
	if ie.LeftExpression != nil {
		children = append(children,
			treeChild{"LeftExpression", ie.LeftExpression},
		)
	}
	if ie.NumberExpression != nil {
		children = append(children,
			treeChild{"NumberExpression", ie.NumberExpression},
		)
	}
	return renderTree("IndexExpression", children...)
}
