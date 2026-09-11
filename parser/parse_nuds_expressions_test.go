package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
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
	if !testBinaryExpression(t, ifExpression.ConditionExpression, "x", "<", "y") {
		return
	}
	if len(ifExpression.ConsequenceStatement.Statements) != 1 {
		t.Errorf("ifExpression.Consequence.Statements does not contain %v statements. got=%v\n", 1, len(ifExpression.ConsequenceStatement.Statements))
	}
	consequence, ok := ifExpression.ConsequenceStatement.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpression.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T\n", ifExpression.ConsequenceStatement.Statements[0])
	}
	if !testIdentifierExpression(t, consequence.Expression, "x") {
		return
	}
	if ifExpression.AlternativeStatement != nil {
		t.Errorf("ifExpression.Alternative was not nil. got=%v\n", ifExpression.AlternativeStatement)
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
	if !testBinaryExpression(t, ifExpression.ConditionExpression, "x", "<", "y") {
		return
	}
	if len(ifExpression.ConsequenceStatement.Statements) != 1 {
		t.Errorf("ifExpression.Consequence.Statements does not contain %v statements. got=%v\n", 1, len(ifExpression.ConsequenceStatement.Statements))
	}
	consequence, ok := ifExpression.ConsequenceStatement.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpression.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T\n", ifExpression.ConsequenceStatement.Statements[0])
	}
	if !testIdentifierExpression(t, consequence.Expression, "x") {
		return
	}
	alternative, ok := ifExpression.AlternativeStatement.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpression.Alternative.Statements[0] is not *ast.ExpressionStatement. got=%T\n", ifExpression.AlternativeStatement.Statements[0])
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
	if len(functionLiteralExpression.ParameterExpressions) != 2 {
		t.Fatalf("functionLiteralExpression.Parameters does not contain %v parameters. got=%v\n", 2, len(functionLiteralExpression.ParameterExpressions))
	}
	testLiteralExpression(t, functionLiteralExpression.ParameterExpressions[0], "x")
	testLiteralExpression(t, functionLiteralExpression.ParameterExpressions[1], "y")
	if len(functionLiteralExpression.BodyStatement.Statements) != 1 {
		t.Fatalf("functionLiteralExpression.Body.Statements does not contain %v statements. got=%v\n", 1, len(functionLiteralExpression.BodyStatement.Statements))
	}
	bodyStatement := testExpressionStatement(t, functionLiteralExpression.BodyStatement.Statements[0])
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
		if len(functionLiteralExpression.ParameterExpressions) != len(test.expectedParams) {
			t.Fatalf("functionLiteralExpression.Parameters does not contain %v parameters. got=%v\n", len(test.expectedParams), len(functionLiteralExpression.ParameterExpressions))
		}
		for i, identifier := range test.expectedParams {
			testLiteralExpression(t, functionLiteralExpression.ParameterExpressions[i], identifier)

		}
	}
}

// TestParseEmptyFunctionLiteral checks a function literal with no parameters
// and an empty body.
func TestParseEmptyFunctionLiteral(t *testing.T) {
	program := testParseProgram(t, "fn() {};", 1)
	functionExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("expression is not *ast.FunctionExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if len(functionExpression.ParameterExpressions) != 0 {
		t.Fatalf("functionExpression.Parameters is not empty. got=%v\n", len(functionExpression.ParameterExpressions))
	}
	if len(functionExpression.BodyStatement.Statements) != 0 {
		t.Fatalf("functionExpression.Body.Statements is not empty. got=%v\n", len(functionExpression.BodyStatement.Statements))
	}
}

// TestParseEmptyIfBlock checks an if expression whose consequence has no
// statements and that has no else branch.
func TestParseEmptyIfBlock(t *testing.T) {
	program := testParseProgram(t, "if (x) {};", 1)
	ifExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expression is not *ast.IfExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if !testIdentifierExpression(t, ifExpression.ConditionExpression, "x") {
		return
	}
	if len(ifExpression.ConsequenceStatement.Statements) != 0 {
		t.Fatalf("ifExpression.Consequence.Statements is not empty. got=%v\n", len(ifExpression.ConsequenceStatement.Statements))
	}
	if ifExpression.AlternativeStatement != nil {
		t.Fatalf("ifExpression.Alternative was not nil. got=%v\n", ifExpression.AlternativeStatement)
	}
}

// TestParseDeeplyNestedGroupedExpression checks that repeated parentheses
// around a single literal still resolve to that literal.
func TestParseDeeplyNestedGroupedExpression(t *testing.T) {
	program := testParseProgram(t, "(((5)));", 1)
	testIntegerExpression(t, testExpressionStatement(t, program.Statements[0]).Expression, 5)
}

// TestParseDoubleUnaryMinus checks "--5", which is only unambiguous because
// MINUS is registered as both a nud (unary "-x") and a led (binary "a - b");
// here the second "-" is parsed as another unary nud, not a binary operator.
func TestParseDoubleUnaryMinus(t *testing.T) {
	program := testParseProgram(t, "--5;", 1)
	outer, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.UnaryExpression)
	if !ok {
		t.Fatalf("expression is not *ast.UnaryExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if outer.Operator != "-" {
		t.Fatalf("outer.Operator is not '-'. got=%v\n", outer.Operator)
	}
	testUnaryExpression(t, outer.RightExpression, "-", 5)
}

// TestParseTrailingCommaInFunctionParametersIsTolerated documents that the
// parameter loop accepts (and silently ignores) a trailing comma before the
// closing delimiter, since it only requires a COMMA between two identifiers
// rather than rejecting one right before RPAREN.
func TestParseTrailingCommaInFunctionParametersIsTolerated(t *testing.T) {
	program := testParseProgram(t, "fn(x, y,) {};", 1)
	functionExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("expression is not *ast.FunctionExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if len(functionExpression.ParameterExpressions) != 2 {
		t.Fatalf("functionExpression.Parameters does not contain 2 parameters. got=%v\n", len(functionExpression.ParameterExpressions))
	}
	testLiteralExpression(t, functionExpression.ParameterExpressions[0], "x")
	testLiteralExpression(t, functionExpression.ParameterExpressions[1], "y")
}

// TestParseUnterminatedGroupedExpressionRecordsError checks that a missing
// closing ")" is reported as a parser error rather than panicking or
// silently accepting the input.
func TestParseUnterminatedGroupedExpressionRecordsError(t *testing.T) {
	parser := New(lexer.New("(1 + 2"))
	parser.ParseProgram()
	if len(parser.Errors()) == 0 {
		t.Fatalf("expected at least 1 error for an unterminated grouped expression, got 0\n")
	}
}

// TestParseUnterminatedIfConditionRecordsError checks that a missing closing
// ")" on an if condition is reported as a parser error.
func TestParseUnterminatedIfConditionRecordsError(t *testing.T) {
	parser := New(lexer.New("if (x < y"))
	parser.ParseProgram()
	if len(parser.Errors()) == 0 {
		t.Fatalf("expected at least 1 error for an unterminated if condition, got 0\n")
	}
}

// TestParseUnterminatedFunctionParametersIsSilentlyDropped documents a known
// parser limitation: unlike grouped expressions and if conditions,
// parseFunctionExpression returns nil (as the ast.Expression interface, so
// no typed-nil trap here) without recording an error when its parameter list
// runs into EOF instead of a closing ")". Parsing still completes without
// panicking, but the caller gets no diagnostic for genuinely malformed input.
func TestParseUnterminatedFunctionParametersIsSilentlyDropped(t *testing.T) {
	input := "fn(x, y"
	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("input=%q: program.Statements does not contain 1 statement. got=%v\n", input, len(program.Statements))
	}
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	if expressionStatement.Expression != nil {
		t.Fatalf("input=%q: expected a nil expression for unterminated input, got=%T\n", input, expressionStatement.Expression)
	}
	if len(parser.Errors()) != 0 {
		t.Fatalf("input=%q: expected no recorded errors (documenting current behavior), got=%v\n", input, parser.Errors())
	}
}
