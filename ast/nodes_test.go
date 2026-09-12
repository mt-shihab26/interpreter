package ast

import (
	"monkey/token"
	"testing"
)

// TestString checks that Program.String() reconstructs Monkey source code from the AST.
func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				IdentifierExpression: &IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				ValueExpression: &IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%v", program.String())
	}
}

// TestTree checks that Program.Tree() renders a nested LetStatement's Name/Value children as an indented, connector-based tree.
func TestTree(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				IdentifierExpression: &IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				ValueExpression: &IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	expected := "Program\n" +
		"└─ Statement[0]: LetStatement\n" +
		"   ├─ IdentifierExpression: IdentifierExpression \"myVar\"\n" +
		"   └─ ValueExpression: IdentifierExpression \"anotherVar\""
	if actual := program.Tree(); actual != expected {
		t.Errorf("program.Tree() wrong.\nexpected=\n%v\ngot=\n%v", expected, actual)
	}
}

// TestTreeEmptyProgram checks that a Program with no statements renders as just its own header line.
func TestTreeEmptyProgram(t *testing.T) {
	program := &Program{Statements: []Statement{}}
	if actual := program.Tree(); actual != "Program" {
		t.Errorf("program.Tree() wrong, got=%q", actual)
	}
}
