package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

// UnaryExpression implements the Expression interface.
type UnaryExpression struct {
	Token           token.Token
	Operator        string
	RightExpression Expression
}

func (ue *UnaryExpression) expressionNode() {

}
func (ue *UnaryExpression) TokenLiteral() string {
	return ue.Token.Literal
}

func (ue *UnaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ue.Operator)
	out.WriteString(ue.RightExpression.String())
	out.WriteString(")")
	return out.String()
}

func (ue *UnaryExpression) Tree() string {
	children := []treeChild{}
	if ue.RightExpression != nil {
		children = append(children, treeChild{"Right", ue.RightExpression})
	}
	return renderTree(fmt.Sprintf("UnaryExpression %q", ue.Operator), children...)
}

// IdentifierExpression implements the Expression interface.
type IdentifierExpression struct {
	Token token.Token
	Value string
}

func (i *IdentifierExpression) expressionNode() {

}
func (i *IdentifierExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IdentifierExpression) String() string {
	return i.Value
}

func (i *IdentifierExpression) Tree() string {
	return fmt.Sprintf("IdentifierExpression %q", i.Value)
}

// IntegerExpression implements the Expression interface.
type IntegerExpression struct {
	Token token.Token
	Value int64
}

func (il *IntegerExpression) expressionNode() {

}
func (il *IntegerExpression) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerExpression) String() string {
	return il.Token.Literal
}

func (il *IntegerExpression) Tree() string {
	return fmt.Sprintf("IntegerExpression %d", il.Value)
}

// BooleanExpression implements the Expression interface.
type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (b *BooleanExpression) expressionNode() {

}
func (b *BooleanExpression) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanExpression) String() string {
	return b.Token.Literal
}

func (b *BooleanExpression) Tree() string {
	return fmt.Sprintf("BooleanExpression %v", b.Value)
}

// IfExpression implements the Expression interface.
type IfExpression struct {
	Token                token.Token
	ConditionExpression  Expression
	ConsequenceStatement *BlockStatement
	AlternativeStatement *BlockStatement
}

func (ie *IfExpression) expressionNode() {

}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ie.TokenLiteral())
	out.WriteString(ie.ConditionExpression.String())
	out.WriteString(" ")
	out.WriteString(ie.ConsequenceStatement.String())
	if ie.AlternativeStatement != nil {
		out.WriteString("else ")
		out.WriteString(ie.AlternativeStatement.String())
	}
	return out.String()
}

func (ie *IfExpression) Tree() string {
	children := []treeChild{}
	if ie.ConditionExpression != nil {
		children = append(children, treeChild{"Condition", ie.ConditionExpression})
	}
	if ie.ConsequenceStatement != nil {
		children = append(children, treeChild{"Consequence", ie.ConsequenceStatement})
	}
	if ie.AlternativeStatement != nil {
		children = append(children, treeChild{"Alternative", ie.AlternativeStatement})
	}
	return renderTree("IfExpression", children...)
}

// FunctionExpression implements the Expression interface.
type FunctionExpression struct {
	Token                token.Token
	ParameterExpressions []*IdentifierExpression
	BodyStatement        *BlockStatement
}

func (fe *FunctionExpression) expressionNode() {

}
func (fe *FunctionExpression) TokenLiteral() string {
	return fe.Token.Literal
}

func (fe *FunctionExpression) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, parameter := range fe.ParameterExpressions {
		params = append(params, parameter.String())
	}
	out.WriteString(fe.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fe.BodyStatement.String())
	return out.String()
}

func (fe *FunctionExpression) Tree() string {
	children := make([]treeChild, 0, len(fe.ParameterExpressions)+1)
	for i, parameter := range fe.ParameterExpressions {
		if parameter == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("Parameter[%d]", i), parameter})
	}
	if fe.BodyStatement != nil {
		children = append(children, treeChild{"Body", fe.BodyStatement})
	}
	return renderTree("FunctionExpression", children...)
}
