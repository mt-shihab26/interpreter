package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
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

// TestLetStatementWithoutTrailingSemicolon checks that the trailing ";" is
// optional, since parseLetStatement only consumes it when present.
func TestLetStatementWithoutTrailingSemicolon(t *testing.T) {
	program := testParseProgram(t, "let x = 5", 1)
	if actual := program.String(); actual != "let x = 5;" {
		t.Errorf("expected=%v, got=%v\n", "let x = 5;", actual)
	}
}

// TestParseMalformedLetStatementDoesNotPanic exercises every way a "let"
// statement can be malformed (missing identifier, missing "=", missing
// value). Each case must record a parser error and must not panic: a
// dropped/broken let statement is a concrete *ast.LetStatement returned as
// nil from parseLetStatement, and parseStatement must convert that into a
// true nil ast.Statement -- not a non-nil interface wrapping a typed nil --
// or later calls like .String() on it will nil-dereference and panic.
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
