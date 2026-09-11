package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
)

// LetStatement implements the Statement interface.
type LetStatement struct {
	Token                token.Token
	IdentifierExpression *IdentifierExpression
	ValueExpression      Expression
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
	out.WriteString(ls.IdentifierExpression.String())
	out.WriteString(" = ")
	if ls.ValueExpression != nil {
		out.WriteString(ls.ValueExpression.String())
	}
	out.WriteString(";")
	return out.String()
}

func (ls *LetStatement) Tree() string {
	children := []treeChild{}
	if ls.IdentifierExpression != nil {
		children = append(children, treeChild{"Name", ls.IdentifierExpression})
	}
	if ls.ValueExpression != nil {
		children = append(children, treeChild{"Value", ls.ValueExpression})
	}
	return renderTree("LetStatement", children...)
}

// ReturnStatement implements the Statement interface.
type ReturnStatement struct {
	Token           token.Token
	ValueExpression Expression
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
	if rs.ValueExpression != nil {
		out.WriteString(rs.ValueExpression.String())
	}
	out.WriteString(";")
	return out.String()
}

func (rs *ReturnStatement) Tree() string {
	children := []treeChild{}
	if rs.ValueExpression != nil {
		children = append(children, treeChild{"Value", rs.ValueExpression})
	}
	return renderTree("ReturnStatement", children...)
}

// ExpressionStatement implements the Statement interface.
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

func (es *ExpressionStatement) Tree() string {
	children := []treeChild{}
	if es.Expression != nil {
		children = append(children, treeChild{"Expression", es.Expression})
	}
	return renderTree("ExpressionStatement", children...)
}

// BlockStatement implements the Statement interface.
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

func (bs *BlockStatement) Tree() string {
	children := make([]treeChild, 0, len(bs.Statements))
	for i, statement := range bs.Statements {
		if statement == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("[%d]", i), statement})
	}
	return renderTree("BlockStatement", children...)
}
