package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
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
	if !testIdentifierExpression(t, callExpression.FunctionExpression, "add") {
		return
	}
	if len(callExpression.ArgumentExpressions) != 3 {
		t.Fatalf("wrong length of arguments. got=%v\n", len(callExpression.ArgumentExpressions))
	}
	testLiteralExpression(t, callExpression.ArgumentExpressions[0], 1)
	testBinaryExpression(t, callExpression.ArgumentExpressions[1], 2, "*", 3)
	testBinaryExpression(t, callExpression.ArgumentExpressions[2], 4, "+", 5)
}

// TestParseEmptyCallArguments checks a call expression with no arguments.
func TestParseEmptyCallArguments(t *testing.T) {
	program := testParseProgram(t, "foo();", 1)
	callExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if !testIdentifierExpression(t, callExpression.FunctionExpression, "foo") {
		return
	}
	if len(callExpression.ArgumentExpressions) != 0 {
		t.Fatalf("callExpression.Arguments is not empty. got=%v\n", len(callExpression.ArgumentExpressions))
	}
}

// TestParseImmediatelyInvokedFunctionExpression checks that a call
// expression's callee can be a function literal, not just an identifier.
func TestParseImmediatelyInvokedFunctionExpression(t *testing.T) {
	program := testParseProgram(t, "fn(x) { x; }(5);", 1)
	callExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	functionExpression, ok := callExpression.FunctionExpression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("callExpression.Function is not *ast.FunctionExpression. got=%T\n", callExpression.FunctionExpression)
	}
	testLiteralExpression(t, functionExpression.ParameterExpressions[0], "x")
	if len(callExpression.ArgumentExpressions) != 1 {
		t.Fatalf("callExpression.Arguments does not contain 1 argument. got=%v\n", len(callExpression.ArgumentExpressions))
	}
	testIntegerExpression(t, callExpression.ArgumentExpressions[0], 5)
}

// TestParseChainedCallExpressions checks that a call's callee can itself be
// a call expression, e.g. curried invocation "add(1)(2)".
func TestParseChainedCallExpressions(t *testing.T) {
	program := testParseProgram(t, "add(1)(2);", 1)
	outerCall, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if len(outerCall.ArgumentExpressions) != 1 {
		t.Fatalf("outerCall.Arguments does not contain 1 argument. got=%v\n", len(outerCall.ArgumentExpressions))
	}
	testIntegerExpression(t, outerCall.ArgumentExpressions[0], 2)
	innerCall, ok := outerCall.FunctionExpression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("outerCall.Function is not *ast.CallExpression. got=%T\n", outerCall.FunctionExpression)
	}
	if !testIdentifierExpression(t, innerCall.FunctionExpression, "add") {
		return
	}
	if len(innerCall.ArgumentExpressions) != 1 {
		t.Fatalf("innerCall.Arguments does not contain 1 argument. got=%v\n", len(innerCall.ArgumentExpressions))
	}
	testIntegerExpression(t, innerCall.ArgumentExpressions[0], 1)
}

// TestParseTrailingCommaInCallArgumentsIsTolerated documents that the
// argument loop accepts (and silently ignores) a trailing comma before the
// closing delimiter, since it only requires a COMMA between two expressions
// rather than rejecting one right before RPAREN.
func TestParseTrailingCommaInCallArgumentsIsTolerated(t *testing.T) {
	program := testParseProgram(t, "foo(1, 2,);", 1)
	callExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if len(callExpression.ArgumentExpressions) != 2 {
		t.Fatalf("callExpression.Arguments does not contain 2 arguments. got=%v\n", len(callExpression.ArgumentExpressions))
	}
	testIntegerExpression(t, callExpression.ArgumentExpressions[0], 1)
	testIntegerExpression(t, callExpression.ArgumentExpressions[1], 2)
}

// TestParseUnterminatedCallArgumentsIsSilentlyDropped documents a known
// parser limitation: unlike grouped expressions and if conditions,
// parseCallExpression returns nil (as the ast.Expression interface, so no
// typed-nil trap here) without recording an error when its argument list
// runs into EOF instead of a closing ")". Parsing still completes without
// panicking, but the caller gets no diagnostic for genuinely malformed input.
func TestParseUnterminatedCallArgumentsIsSilentlyDropped(t *testing.T) {
	input := "foo(1, 2"
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
