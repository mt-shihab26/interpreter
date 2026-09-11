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
		value := returnStatement.(*ast.ReturnStatement).ValueExpression
		if !testLiteralExpression(t, value, test.returnValue) {
			return
		}
	}
}

// TestReturnStatementWithoutTrailingSemicolon checks that the trailing ";" is
// optional, since parseReturnStatement only consumes it when present.
func TestReturnStatementWithoutTrailingSemicolon(t *testing.T) {
	program := testParseProgram(t, "return 5", 1)
	if actual := program.String(); actual != "return 5;" {
		t.Errorf("expected=%v, got=%v\n", "return 5;", actual)
	}
}
