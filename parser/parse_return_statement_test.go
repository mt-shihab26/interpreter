package parser

import (
	"testing"

	"monkey/ast"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input       string
		returnValue any
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		returnStatement := program.Statements[0]
		if !testReturnStatement(t, returnStatement) {
			return
		}
		value := returnStatement.(*ast.ReturnStatement).ReturnValue
		if !testLiteralExpression(t, value, test.returnValue) {
			return
		}
	}
}
