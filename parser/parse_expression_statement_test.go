package parser

import (
	"testing"

	"monkey/lexer"
)

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input          string
		expected       string
		statementCount int
	}{
		{"-a * b", "((-a) * b)", 1},
		{"!-a", "(!(-a))", 1},
		{"a + b + c", "((a + b) + c)", 1},
		{"a + b - c", "((a + b) - c)", 1},
		{"a * b * c", "((a * b) * c)", 1},
		{"a * b / c", "((a * b) / c)", 1},
		{"a + b / c", "(a + (b / c))", 1},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)", 1},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)", 2},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))", 1},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))", 1},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))", 1},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))", 1},
		{"true", "true", 1},
		{"false", "false", 1},
		{"3 > 5 == false", "((3 > 5) == false)", 1},
		{"3 < 5 == true", "((3 < 5) == true)", 1},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)", 1},
		{"(5 + 5) * 2", "((5 + 5) * 2)", 1},
		{"2 / (5 + 5)", "(2 / (5 + 5))", 1},
		{"-(5 + 5)", "(-(5 + 5))", 1},
		{"!(true == true)", "(!(true == true))", 1},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)", 1},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))", 1},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))", 1},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, test.statementCount)
		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected=%v, got=%v\n", test.expected, actual)
		}
	}
}

// TestExpressionStatementWithoutTrailingSemicolon checks that the trailing
// ";" is optional, since parseExpressionStatement only consumes it when
// present.
func TestExpressionStatementWithoutTrailingSemicolon(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5", "5"},
		{"x + y", "(x + y)"},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		if actual := program.String(); actual != test.expected {
			t.Errorf("input=%q expected=%v, got=%v\n", test.input, test.expected, actual)
		}
	}
}

// TestParseIllegalTokenRecordsError checks that a token the lexer can't
// classify records a parser error instead of panicking.
func TestParseIllegalTokenRecordsError(t *testing.T) {
	parser := New(lexer.New("@"))
	parser.ParseProgram()
	if len(parser.Errors()) == 0 {
		t.Fatalf("expected at least 1 error for an illegal token, got 0\n")
	}
}
