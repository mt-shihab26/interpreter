package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

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

type UnaryExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
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
	out.WriteString(ue.Right.String())
	out.WriteString(")")
	return out.String()
}

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

type LetStatement struct {
	Token token.Token
	Name  *IdentifierExpression
	Value Expression
}

func (ls *LetStatement) statementNode() {

}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {

}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {

}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {

}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {

}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ie.TokenLiteral())
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type FunctionExpression struct {
	Token      token.Token
	Parameters []*IdentifierExpression
	Body       *BlockStatement
}

func (fe *FunctionExpression) expressionNode() {

}
func (fe *FunctionExpression) TokenLiteral() string {
	return fe.Token.Literal
}

func (fe *FunctionExpression) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, parameter := range fe.Parameters {
		params = append(params, parameter.String())
	}
	out.WriteString(fe.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fe.Body.String())
	return out.String()
}

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
