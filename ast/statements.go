package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

// LetStatement implements the Statement interface.
type LetStatement struct {
	Token                token.Token
	IdentifierExpression *IdentifierExpression
	ValueExpression      Expression
}

// statementNode marks LetStatement as an ast.Statement.
func (ls *LetStatement) statementNode() {

}

// TokenLiteral returns the literal of the "let" token.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// Code reconstructs the statement as "let <identifier> = <value>;".
func (ls *LetStatement) Code() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.IdentifierExpression.Code())
	out.WriteString(" = ")
	if ls.ValueExpression != nil {
		out.WriteString(ls.ValueExpression.Code())
	}
	out.WriteString(";")
	return out.String()
}

// Tree renders the statement with the identifier as a "Name" child and the value as a "Value" child.
func (ls *LetStatement) Tree() string {
	children := []treeChild{}
	if ls.IdentifierExpression != nil {
		children = append(children, treeChild{"IdentifierExpression", ls.IdentifierExpression})
	}
	if ls.ValueExpression != nil {
		children = append(children, treeChild{"ValueExpression", ls.ValueExpression})
	}
	return renderTree("LetStatement", children...)
}

// ReturnStatement implements the Statement interface.
type ReturnStatement struct {
	Token           token.Token
	ValueExpression Expression
}

// statementNode marks ReturnStatement as an ast.Statement.
func (rs *ReturnStatement) statementNode() {

}

// TokenLiteral returns the literal of the "return" token.
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// Code reconstructs the statement as "return <value>;".
func (rs *ReturnStatement) Code() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")
	if rs.ValueExpression != nil {
		out.WriteString(rs.ValueExpression.Code())
	}
	out.WriteString(";")
	return out.String()
}

// Tree renders the statement with the value as a "Value" child.
func (rs *ReturnStatement) Tree() string {
	children := []treeChild{}
	if rs.ValueExpression != nil {
		children = append(children, treeChild{"ValueExpression", rs.ValueExpression})
	}
	return renderTree("ReturnStatement", children...)
}

// ExpressionStatement implements the Statement interface.
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// statementNode marks ExpressionStatement as an ast.Statement.
func (es *ExpressionStatement) statementNode() {

}

// TokenLiteral returns the literal of the expression's first token.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// Code reconstructs the statement as "<expression>;", or "" if it holds none.
func (es *ExpressionStatement) Code() string {
	if es.Expression != nil {
		return es.Expression.Code() + ";"
	}
	return ""
}

// Tree renders the statement with the wrapped expression as an "Expression" child.
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

// statementNode marks BlockStatement as an ast.Statement.
func (bs *BlockStatement) statementNode() {

}

// TokenLiteral returns the literal of the opening "{" token.
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// Code reconstructs the block as its statements' source, one per indented line.
func (bs *BlockStatement) Code() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString(indent(s.Code()))
		out.WriteString("\n")
	}
	out.WriteString("}")
	return out.String()
}

// indent prefixes every line of s with a tab, so nested blocks compound their indentation.
func indent(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = "\t" + line
	}
	return strings.Join(lines, "\n")
}

// Tree renders the block with one indexed child per statement.
func (bs *BlockStatement) Tree() string {
	children := make([]treeChild, 0, len(bs.Statements))
	for i, statement := range bs.Statements {
		if statement == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("Statement[%d]", i), statement})
	}
	return renderTree("BlockStatement", children...)
}
