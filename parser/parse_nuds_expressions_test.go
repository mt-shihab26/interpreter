package parser

import (
	"testing"

	"monkey/ast"
)

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	testIdentifierExpression(t, expressionStatement.Expression, "foobar")
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	testIntegerExpression(t, expressionStatement.Expression, 5)
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		expressionStatement := testExpressionStatement(t, program.Statements[0])
		testBooleanExpression(t, expressionStatement.Expression, test.value)
	}

}

func TestParsingUnaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		right    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		expressionStatement := testExpressionStatement(t, program.Statements[0])
		testUnaryExpression(t, expressionStatement.Expression, test.operator, test.right)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	ifExpression, ok := expressionStatement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expressionStatement.Expression is not *ast.IfExpression. got=%T\n", expressionStatement.Expression)
	}
	if !testBinaryExpression(t, ifExpression.Condition, "x", "<", "y") {
		return
	}
	if len(ifExpression.Consequence.Statements) != 1 {
		t.Errorf("ifExpression.Consequence.Statements does not contain %v statements. got=%v\n", 1, len(ifExpression.Consequence.Statements))
	}
	consequence, ok := ifExpression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpression.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T\n", ifExpression.Consequence.Statements[0])
	}
	if !testIdentifierExpression(t, consequence.Expression, "x") {
		return
	}
	if ifExpression.Alternative != nil {
		t.Errorf("ifExpression.Alternative was not nil. got=%v\n", ifExpression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	ifExpression, ok := expressionStatement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expressionStatement.Expression is not *ast.IfExpression. got=%T\n", expressionStatement.Expression)
	}
	if !testBinaryExpression(t, ifExpression.Condition, "x", "<", "y") {
		return
	}
	if len(ifExpression.Consequence.Statements) != 1 {
		t.Errorf("ifExpression.Consequence.Statements does not contain %v statements. got=%v\n", 1, len(ifExpression.Consequence.Statements))
	}
	consequence, ok := ifExpression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpression.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T\n", ifExpression.Consequence.Statements[0])
	}
	if !testIdentifierExpression(t, consequence.Expression, "x") {
		return
	}
	alternative, ok := ifExpression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpression.Alternative.Statements[0] is not *ast.ExpressionStatement. got=%T\n", ifExpression.Alternative.Statements[0])
	}
	if !testIdentifierExpression(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	functionLiteralExpression, ok := expressionStatement.Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("expressionStatement.Expression is not *ast.FunctionExpression. got=%T\n", expressionStatement.Expression)
	}
	if len(functionLiteralExpression.Parameters) != 2 {
		t.Fatalf("functionLiteralExpression.Parameters does not contain %v parameters. got=%v\n", 2, len(functionLiteralExpression.Parameters))
	}
	testLiteralExpression(t, functionLiteralExpression.Parameters[0], "x")
	testLiteralExpression(t, functionLiteralExpression.Parameters[1], "y")
	if len(functionLiteralExpression.Body.Statements) != 1 {
		t.Fatalf("functionLiteralExpression.Body.Statements does not contain %v statements. got=%v\n", 1, len(functionLiteralExpression.Body.Statements))
	}
	bodyStatement := testExpressionStatement(t, functionLiteralExpression.Body.Statements[0])
	testBinaryExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		expressionStatement := testExpressionStatement(t, program.Statements[0])
		functionLiteralExpression, ok := expressionStatement.Expression.(*ast.FunctionExpression)
		if !ok {
			t.Fatalf("expressionStatement.Expression is not *ast.FunctionExpression. got=%T\n", expressionStatement.Expression)
		}
		if len(functionLiteralExpression.Parameters) != len(test.expectedParams) {
			t.Fatalf("functionLiteralExpression.Parameters does not contain %v parameters. got=%v\n", len(test.expectedParams), len(functionLiteralExpression.Parameters))
		}
		for i, identifier := range test.expectedParams {
			testLiteralExpression(t, functionLiteralExpression.Parameters[i], identifier)

		}
	}
}
