package ast

// Node is implemented by every AST node, statements and expressions alike.
type Node interface {
	// TokenLiteral returns the literal of the token the node was built from.
	TokenLiteral() string
	// String reconstructs the node (and its descendants) as Monkey source code.
	String() string
	// Tree returns a multi-line, indented representation of the node and
	// its descendants -- unlike String(), which reconstructs Monkey
	// source code, Tree() is meant for inspecting the shape of a parsed
	// AST while debugging.
	Tree() string
}

// Statement is implemented by AST nodes that represent a statement, e.g. LetStatement.
type Statement interface {
	Node
	// statementNode is a marker method with no purpose other than to
	// distinguish Statement from Expression at compile time.
	statementNode()
}

// Expression is implemented by AST nodes that represent an expression, e.g. BinaryExpression.
type Expression interface {
	Node
	// expressionNode is a marker method with no purpose other than to
	// distinguish Expression from Statement at compile time.
	expressionNode()
}
