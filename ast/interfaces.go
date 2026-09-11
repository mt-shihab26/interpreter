package ast

type Node interface {
	TokenLiteral() string
	String() string
	// Tree returns a multi-line, indented representation of the node and
	// its descendants -- unlike String(), which reconstructs Monkey
	// source code, Tree() is meant for inspecting the shape of a parsed
	// AST while debugging.
	Tree() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}
