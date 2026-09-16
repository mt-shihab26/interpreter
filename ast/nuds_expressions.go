package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"sort"
	"strings"
)

// UnaryExpression implements the Expression interface.
type UnaryExpression struct {
	Token           token.Token
	Operator        string
	RightExpression Expression
}

// expressionNode marks UnaryExpression as an ast.Expression.
func (ue *UnaryExpression) expressionNode() {

}

// TokenLiteral returns the operator token's literal, e.g. "-".
func (ue *UnaryExpression) TokenLiteral() string {
	return ue.Token.Literal
}

// Code reconstructs the expression as "(<operator><operand>)", e.g. "(-x)".
func (ue *UnaryExpression) Code() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ue.Operator)
	out.WriteString(ue.RightExpression.Code())
	out.WriteString(")")
	return out.String()
}

// Tree renders the expression with the operator in its header and its operand as a "Right" child.
func (ue *UnaryExpression) Tree() string {
	children := []treeChild{}
	if ue.RightExpression != nil {
		children = append(children, treeChild{"RightExpression", ue.RightExpression})
	}
	return renderTree(fmt.Sprintf("UnaryExpression %q", ue.Operator), children...)
}

// IdentifierExpression implements the Expression interface.
type IdentifierExpression struct {
	Token token.Token
	Value string
}

// expressionNode marks IdentifierExpression as an ast.Expression.
func (i *IdentifierExpression) expressionNode() {

}

// TokenLiteral returns the identifier token's literal, e.g. "foobar".
func (i *IdentifierExpression) TokenLiteral() string {
	return i.Token.Literal
}

// Code reconstructs the expression as the bare identifier name.
func (i *IdentifierExpression) Code() string {
	return i.Value
}

// Tree renders the expression as a single leaf line naming its value.
func (i *IdentifierExpression) Tree() string {
	return fmt.Sprintf("IdentifierExpression %q", i.Value)
}

// IntegerExpression implements the Expression interface.
type IntegerExpression struct {
	Token token.Token
	Value int64
}

// expressionNode marks IntegerExpression as an ast.Expression.
func (il *IntegerExpression) expressionNode() {

}

// TokenLiteral returns the integer token's literal, e.g. "5".
func (il *IntegerExpression) TokenLiteral() string {
	return il.Token.Literal
}

// Code reconstructs the expression as its original digit literal.
func (il *IntegerExpression) Code() string {
	return il.Token.Literal
}

// Tree renders the expression as a single leaf line naming its parsed value.
func (il *IntegerExpression) Tree() string {
	return fmt.Sprintf("IntegerExpression %d", il.Value)
}

// StringExpression implements the Expression interface.
type StringExpression struct {
	Token token.Token
	Value string
}

// expressionNode marks StringExpression as an ast.Expression.
func (il *StringExpression) expressionNode() {

}

// TokenLiteral returns the string token's literal, e.g. "foobar".
func (il *StringExpression) TokenLiteral() string {
	return il.Token.Literal
}

// Code reconstructs the expression as a double-quoted string literal.
func (il *StringExpression) Code() string {
	return fmt.Sprintf("%q", il.Value)
}

// Tree renders the expression as a single leaf line naming its parsed value.
func (il *StringExpression) Tree() string {
	return fmt.Sprintf("StringExpression %q", il.Value)
}

// BooleanExpression implements the Expression interface.
type BooleanExpression struct {
	Token token.Token
	Value bool
}

// expressionNode marks BooleanExpression as an ast.Expression.
func (b *BooleanExpression) expressionNode() {

}

// TokenLiteral returns the boolean token's literal, "true" or "false".
func (b *BooleanExpression) TokenLiteral() string {
	return b.Token.Literal
}

// Code reconstructs the expression as "true" or "false".
func (b *BooleanExpression) Code() string {
	return b.Token.Literal
}

// Tree renders the expression as a single leaf line naming its value.
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

// expressionNode marks IfExpression as an ast.Expression.
func (ie *IfExpression) expressionNode() {

}

// TokenLiteral returns the literal of the "if" token.
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// Code reconstructs the expression as "if<condition> <consequence>else <alternative>", omitting the else clause if there is none.
func (ie *IfExpression) Code() string {
	var out bytes.Buffer
	out.WriteString(ie.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(ie.ConditionExpression.Code())
	out.WriteString(") ")
	out.WriteString(ie.ConsequenceStatement.Code())
	if ie.AlternativeStatement != nil {
		out.WriteString(" else ")
		out.WriteString(ie.AlternativeStatement.Code())
	}
	return out.String()
}

// Tree renders the expression with its condition, consequence, and (if present) alternative as labeled children.
func (ie *IfExpression) Tree() string {
	children := []treeChild{}
	if ie.ConditionExpression != nil {
		children = append(children, treeChild{"ConditionExpression", ie.ConditionExpression})
	}
	if ie.ConsequenceStatement != nil {
		children = append(children, treeChild{"ConsequenceStatement", ie.ConsequenceStatement})
	}
	if ie.AlternativeStatement != nil {
		children = append(children, treeChild{"AlternativeStatement", ie.AlternativeStatement})
	}
	return renderTree("IfExpression", children...)
}

// FunctionExpression implements the Expression interface.
type FunctionExpression struct {
	Token                token.Token
	ParameterExpressions []*IdentifierExpression
	BodyStatement        *BlockStatement
}

// expressionNode marks FunctionExpression as an ast.Expression.
func (fe *FunctionExpression) expressionNode() {

}

// TokenLiteral returns the literal of the "fn" token.
func (fe *FunctionExpression) TokenLiteral() string {
	return fe.Token.Literal
}

// Code reconstructs the expression as "fn(<param>, <param>, ...)<body>".
func (fe *FunctionExpression) Code() string {
	var out bytes.Buffer
	params := []string{}
	for _, parameter := range fe.ParameterExpressions {
		params = append(params, parameter.Code())
	}
	out.WriteString(fe.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fe.BodyStatement.Code())
	return out.String()
}

// Tree renders the expression with each parameter as an indexed "Parameter[i]" child and the body as a "Body" child.
func (fe *FunctionExpression) Tree() string {
	children := make([]treeChild, 0, len(fe.ParameterExpressions)+1)
	for i, parameter := range fe.ParameterExpressions {
		if parameter == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("ParameterExpression[%d]", i), parameter})
	}
	if fe.BodyStatement != nil {
		children = append(children, treeChild{"BodyStatement", fe.BodyStatement})
	}
	return renderTree("FunctionExpression", children...)
}

// ArrayExpression implements the Expression interface.
type ArrayExpression struct {
	Token    token.Token
	Elements []Expression
}

// expressionNode marks FunctionExpression as an ast.Expression.
func (ae *ArrayExpression) expressionNode() {

}

// TokenLiteral returns the literal of the "fn" token.
func (ae *ArrayExpression) TokenLiteral() string {
	return ae.Token.Literal
}

// Code reconstructs the expression as "fn(<param>, <param>, ...)<body>".
func (ae *ArrayExpression) Code() string {
	var out bytes.Buffer
	params := []string{}
	for _, parameter := range ae.Elements {
		params = append(params, parameter.Code())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("]")
	return out.String()
}

// Tree renders the expression with each parameter as an indexed "Parameter[i]" child and the body as a "Body" child.
func (ae *ArrayExpression) Tree() string {
	children := make([]treeChild, 0, len(ae.Elements)+1)
	for i, parameter := range ae.Elements {
		if parameter == nil {
			continue
		}
		children = append(children, treeChild{fmt.Sprintf("Elements[%d]", i), parameter})
	}
	return renderTree("ArrayExpression", children...)
}

// ArrayExpression implements the Expression interface.
type HashExpression struct {
	Token token.Token
	Pairs map[Expression]Expression
}

// expressionNode marks HashExpression as an ast.Expression.
func (ae *HashExpression) expressionNode() {

}

// HashExpression returns the literal of the "{" token.
func (he *HashExpression) TokenLiteral() string {
	return he.Token.Literal
}

// Code reconstructs the expression as "fn(<param>, <param>, ...)<body>".
func (he *HashExpression) Code() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range he.Pairs {
		pairs = append(pairs, key.TokenLiteral()+":"+value.TokenLiteral())

	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// Tree renders each hash pair's key and value as indexed children.
func (he *HashExpression) Tree() string {
	type pair struct {
		key   Expression
		value Expression
	}
	pairs := make([]pair, 0, len(he.Pairs))
	for key, value := range he.Pairs {
		pairs = append(pairs, pair{key: key, value: value})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].key.Code() < pairs[j].key.Code()
	})

	children := make([]treeChild, 0, len(pairs)*2)
	for i, pair := range pairs {
		if pair.key != nil {
			children = append(children, treeChild{fmt.Sprintf("Pairs[%d].Key", i), pair.key})
		}
		if pair.value != nil {
			children = append(children, treeChild{fmt.Sprintf("Pairs[%d].Value", i), pair.value})
		}
	}
	return renderTree("HashExpression", children...)
}
