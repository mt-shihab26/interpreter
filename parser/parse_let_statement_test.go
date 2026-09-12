package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

// TestLetStatements checks that "let <identifier> = <value>;" statements parse with the right identifier and value.
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
		value := letStatement.(*ast.LetStatement).ValueExpression
		if !testLiteralExpression(t, value, test.value) {
			return
		}
	}
}

// TestLetStatementWithoutTrailingSemicolon checks that the trailing ";" is optional in a let statement.
func TestLetStatementWithoutTrailingSemicolon(t *testing.T) {
	program := testParseProgram(t, "let x = 5", 1)
	if actual := program.String(); actual != "let x = 5;" {
		t.Errorf("expected=%v, got=%v\n", "let x = 5;", actual)
	}
}

// TestParseMalformedLetStatementDoesNotPanic checks that a malformed "let" statement records a parser error and never leaves behind a nil-panicking statement.
func TestParseMalformedLetStatementDoesNotPanic(t *testing.T) {
	inputs := []string{
		"let = 5;",  // missing identifier
		"let x 5;",  // missing "="
		"let x = ;", // missing value expression
	}
	for _, input := range inputs {
		parser := New(lexer.New(input))
		program := parser.ParseProgram()
		if len(parser.Errors()) == 0 {
			t.Errorf("input=%q: expected at least 1 error, got 0\n", input)
		}
		for _, statement := range program.Statements {
			if statement == nil {
				t.Errorf("input=%q: program.Statements contains a nil statement\n", input)
				continue
			}
			_ = statement.String() // must not panic
		}
	}
}
