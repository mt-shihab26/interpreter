package parser

import (
	"fmt"
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	program := testParseProgram(t, input, 3)
	tests := []struct {
		identifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, test := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, test.identifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return add(15);
`
	program := testParseProgram(t, input, 3)
	for _, statement := range program.Statements {
		if !testReturnStatement(t, statement) {
			return
		}
	}
}

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

func TestParsingPrefixExpressions(t *testing.T) {
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
		testPrefixExpression(t, expressionStatement.Expression, test.operator, test.right)
	}
}

func TestParsingInfixExpression(t *testing.T) {
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
		testInfixExpression(t, expressionStatement.Expression, test.left, test.operator, test.right)
	}
}

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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	ifExpression, ok := expressionStatement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expressionStatement.Expression is not *ast.IfExpression. got=%T\n", expressionStatement.Expression)
	}
	if !testInfixExpression(t, ifExpression.Condition, "x", "<", "y") {
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
	if !testInfixExpression(t, ifExpression.Condition, "x", "<", "y") {
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
	functionLiteralExpression, ok := expressionStatement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expressionStatement.Expression is not *ast.FunctionLiteral. got=%T\n", expressionStatement.Expression)
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
	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
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
		functionLiteralExpression, ok := expressionStatement.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("expressionStatement.Expression is not *ast.FunctionLiteral. got=%T\n", expressionStatement.Expression)
		}
		if len(functionLiteralExpression.Parameters) != len(test.expectedParams) {
			t.Fatalf("functionLiteralExpression.Parameters does not contain %v parameters. got=%v\n", len(test.expectedParams), len(functionLiteralExpression.Parameters))
		}
		for i, identifier := range test.expectedParams {
			testLiteralExpression(t, functionLiteralExpression.Parameters[i], identifier)

		}
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
	testInfixExpression(t, callExpression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpression.Arguments[2], 4, "+", 5)
}

func testParseProgram(t *testing.T, input string, statementsCount int) *ast.Program {
	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("parser.ParseProgram() returned nil\n")
	}
	if len(program.Statements) != statementsCount {
		t.Fatalf("program.Statements does not contain %v statements. got=%v\n", statementsCount, len(program.Statements))
	}
	return program
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %v errors\n", len(errors))
	for _, message := range errors {
		t.Errorf("parser error: %v\n", message)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral not 'let'. got=%v\n", statement.TokenLiteral())
		return false
	}
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement is not *ast.LetStatement. got=%T\n", statement)
		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not %v. got=%v\n", name, letStatement.Name.Value)
		return false
	}
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral() not %v. got=%v\n", name, letStatement.Name.TokenLiteral())
		return false
	}
	return true
}

func testReturnStatement(t *testing.T, statement ast.Statement) bool {
	returnStatement, ok := statement.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("statement is not *ast.ReturnStatement. got=%T\n", statement)
	}
	if returnStatement.TokenLiteral() != "return" {
		t.Errorf("returnStatement.TokenLiteral not 'return'. got=%v\n", returnStatement.TokenLiteral())
	}
	return true
}

func testExpressionStatement(t *testing.T, statement ast.Statement) *ast.ExpressionStatement {
	expressionStatement, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not *ast.ExpressionStatement. got=%T\n", statement)
	}
	return expressionStatement
}

func testIdentifierExpression(t *testing.T, expression ast.Expression, value string) bool {
	identifierExpression, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression is not *ast.Identifier. got=%T\n", expression)
		return false
	}
	if identifierExpression.Value != value {
		t.Errorf("identifierExpression.Value not %v. got=%v\n", value, identifierExpression.Value)
		return false
	}
	if identifierExpression.TokenLiteral() != value {
		t.Errorf("identifierExpression.TokenLiteral not %v. got=%v\n", value, identifierExpression.TokenLiteral())
		return false
	}
	return true
}

func testIntegerExpression(t *testing.T, expression ast.Expression, value int64) bool {
	integerExpression, ok := expression.(*ast.Integer)
	if !ok {
		t.Errorf("expression is not *ast.Integer. got=%T\n", expression)
		return false
	}
	if integerExpression.Value != value {
		t.Errorf("integerExpression.Value not %v. got=%v\n", value, integerExpression.Value)
		return false
	}
	if integerExpression.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("integerExpression.TokenLiteral not %v. got=%v\n", value, integerExpression.TokenLiteral())
		return false
	}
	return true
}

func testBooleanExpression(t *testing.T, expression ast.Expression, value bool) bool {
	booleanExpression, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expression is not *ast.Boolean. got=%T\n", expression)
		return false
	}
	if booleanExpression.Value != value {
		t.Errorf("booleanExpression.Value not %v. got=%v\n", value, booleanExpression.Value)
		return false
	}
	if booleanExpression.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("booleanExpression.TokenLiteral not %v. got=%v\n", value, booleanExpression.TokenLiteral())
		return false
	}
	return true
}

func testPrefixExpression(t *testing.T, expression ast.Expression, operator string, right any) bool {
	prefixExpression, ok := expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("expression is not *ast.PrefixExpression. got=%T\n", expression)
		return false
	}
	if prefixExpression.Operator != operator {
		t.Fatalf("prefixExpression.Operator is not %v. got=%v\n", operator, prefixExpression.Operator)
		return false
	}
	if !testLiteralExpression(t, prefixExpression.Right, right) {
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, expression ast.Expression, left any, operator string, right any) bool {
	infixExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("expression is not *ast.InfixExpression. got=%T\n", expression)
		return false
	}
	if !testLiteralExpression(t, infixExpression.Left, left) {
		return false
	}
	if infixExpression.Operator != operator {
		t.Fatalf("infixExpression.Operator is not %v. got=%v\n", operator, infixExpression.Operator)
		return false
	}
	if !testLiteralExpression(t, infixExpression.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerExpression(t, expression, int64(v))
	case int64:
		return testIntegerExpression(t, expression, v)
	case string:
		return testIdentifierExpression(t, expression, v)
	case bool:
		return testBooleanExpression(t, expression, v)
	}
	t.Errorf("type of expression not handled. got=%T\n", expression)
	return false
}
