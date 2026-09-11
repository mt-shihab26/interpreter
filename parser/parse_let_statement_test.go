package parser

import (
	"testing"

	"monkey/ast"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input      string
		identifier string
		value      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		letStatement := program.Statements[0]
		if !testLetStatement(t, letStatement, test.identifier) {
			return
		}
		value := letStatement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, value, test.value) {
			return
		}
	}
}
