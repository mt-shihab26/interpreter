package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

// UnaryExpression implements the Expression interface.
type UnaryExpression struct {
	Token    token.Token
	Operator string
	RightExpression    Expression
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

// IfExpression implements the Expression interface.
type IfExpression struct {
	Token       token.Token
	ConditionExpression   Expression
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

// FunctionExpression implements the Expression interface.
type FunctionExpression struct {
	Token      token.Token
	ParameterExpressions []*IdentifierExpression
	BodyStatement       *BlockStatement
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
