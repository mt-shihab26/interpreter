package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				IdentifierExpression: &IdentifierExpression{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				ValueExpression: &IdentifierExpression{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%v", program.String())
	}
}

// TestTree checks that Program.Tree() renders an indented, connector-based
// tree of the node and its descendants, recursing through a nested
// LetStatement into its Name/Value children.
func TestTree(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				IdentifierExpression: &IdentifierExpression{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				ValueExpression: &IdentifierExpression{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	expected := "Program\n" +
		"└─ [0]: LetStatement\n" +
		"   ├─ Name: IdentifierExpression \"myVar\"\n" +
		"   └─ Value: IdentifierExpression \"anotherVar\""
	if actual := program.Tree(); actual != expected {
		t.Errorf("program.Tree() wrong.\nexpected=\n%v\ngot=\n%v", expected, actual)
	}
}

// TestTreeEmptyProgram checks that a Program with no statements renders as
// just its own header line, with no dangling connectors.
func TestTreeEmptyProgram(t *testing.T) {
	program := &Program{Statements: []Statement{}}
	if actual := program.Tree(); actual != "Program" {
		t.Errorf("program.Tree() wrong, got=%q", actual)
	}
}
