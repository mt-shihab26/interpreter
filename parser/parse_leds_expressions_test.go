package parser

import (
	"testing"

	"monkey/ast"
)

func TestParsingBinaryExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     any
		operator string
		right    any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		expressionStatement := testExpressionStatement(t, program.Statements[0])
		testBinaryExpression(t, expressionStatement.Expression, test.left, test.operator, test.right)
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	callExpression, ok := expressionStatement.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("callExpression.Expression is not *ast.CallExpression. got=%T\n", callExpression)
	}
	if !testIdentifierExpression(t, callExpression.Function, "add") {
		return
	}
	if len(callExpression.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%v\n", len(callExpression.Arguments))
	}
	testLiteralExpression(t, callExpression.Arguments[0], 1)
	testBinaryExpression(t, callExpression.Arguments[1], 2, "*", 3)
	testBinaryExpression(t, callExpression.Arguments[2], 4, "+", 5)
}
