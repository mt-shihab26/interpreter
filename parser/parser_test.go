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
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, test.statementCount)
		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected=%v, got=%v", test.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	ifExpression, ok := expressionStatement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expressionStatement.Expression is not *ast.IfExpression. got=%T", expressionStatement.Expression)
	}
	if !testInfixExpression(t, ifExpression.Condition, "x", "<", "y") {
		return
	}
	if len(ifExpression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not '%v'. got=%v", 1, len(ifExpression.Consequence.Statements))
	}
	consequence, ok := ifExpression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", ifExpression.Consequence.Statements[0])
	}
	if !testIdentifierExpression(t, consequence.Expression, "x") {
		return
	}
	if ifExpression.Alternative != nil {
		t.Errorf("ifExpression.Alternative.Statements was not nil. got=%+v", ifExpression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	ifExpression, ok := expressionStatement.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expressionStatement.Expression is not *ast.IfExpression. got=%T", expressionStatement.Expression)
	}
	if !testInfixExpression(t, ifExpression.Condition, "x", "<", "y") {
		return
	}
	if len(ifExpression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not '%v'. got=%v", 1, len(ifExpression.Consequence.Statements))
	}
	consequence, ok := ifExpression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpression.Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", ifExpression.Consequence.Statements[0])
	}
	if !testIdentifierExpression(t, consequence.Expression, "x") {
		return
	}
	alternative, ok := ifExpression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifExpression.Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", ifExpression.Alternative.Statements[0])
	}
	if !testIdentifierExpression(t, alternative.Expression, "y") {
		return
	}
}

func testParseProgram(t *testing.T, input string, statementsCount int) *ast.Program {
	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("parser.ParseProgram() returned nil")
	}
	if len(program.Statements) != statementsCount {
		t.Fatalf("program.Statements does not contain %v statements. got=%v", statementsCount, len(program.Statements))
	}
	return program
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %v errors", len(errors))
	for _, message := range errors {
		t.Errorf("parser error: %v", message)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral not 'let'. got=%v", statement.TokenLiteral())
		return false
	}
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement not *ast.LetStatement. got=%T", statement)
		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not '%v'. got=%v", name, letStatement.Name.Value)
		return false
	}
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral() not '%v'. got=%v", name, letStatement.Name.TokenLiteral())
		return false
	}
	return true
}

func testReturnStatement(t *testing.T, statement ast.Statement) bool {
	returnStatement, ok := statement.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("statement not *ast.ReturnStatement. got=%T", statement)
	}
	if returnStatement.TokenLiteral() != "return" {
		t.Errorf("returnStatement.TokenLiteral not 'return'. got=%v", returnStatement.TokenLiteral())
	}
	return true
}

func testExpressionStatement(t *testing.T, statement ast.Statement) *ast.ExpressionStatement {
	expressionStatement, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not ast.ExpressionStatement. got=%T", statement)
	}
	return expressionStatement
}

func testIdentifierExpression(t *testing.T, expression ast.Expression, value string) bool {
	identifierExpression, ok := expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression is not *ast.Identifier. got=%T", expression)
		return false
	}
	if identifierExpression.Value != value {
		t.Errorf("identifierExpression.Value not %v. got=%v", value, identifierExpression.Value)
		return false
	}
	if identifierExpression.TokenLiteral() != value {
		t.Errorf("identifierExpression.TokenLiteral not %v. got=%v", value, identifierExpression.TokenLiteral())
		return false
	}
	return true
}

func testIntegerExpression(t *testing.T, expression ast.Expression, value int64) bool {
	integerExpression, ok := expression.(*ast.Integer)
	if !ok {
		t.Errorf("expression not *ast.Integer. got=%T", expression)
		return false
	}
	if integerExpression.Value != value {
		t.Errorf("integerExpression.Value not %v. got=%v", value, integerExpression.Value)
		return false
	}
	if integerExpression.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("integerExpression.TokenLiteral not %v. got=%v", value, integerExpression.TokenLiteral())
		return false
	}
	return true
}

func testBooleanExpression(t *testing.T, expression ast.Expression, value bool) bool {
	booleanExpression, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expression not *ast.Boolean. got=%T", expression)
		return false
	}
	if booleanExpression.Value != value {
		t.Errorf("booleanExpression.Value not %v. got=%v", value, booleanExpression.Value)
		return false
	}
	if booleanExpression.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("booleanExpression.TokenLiteral not %v. got=%v", value, booleanExpression.TokenLiteral())
		return false
	}
	return true
}

func testPrefixExpression(t *testing.T, expression ast.Expression, operator string, right any) bool {
	prefixExpression, ok := expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("expression is not *ast.PrefixExpression. got=%T", expression)
		return false
	}
	if prefixExpression.Operator != operator {
		t.Fatalf("prefixExpression.Operator is not '%v'. got=%v", operator, prefixExpression.Operator)
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
		t.Fatalf("expression is not *ast.InfixExpression. got=%T", expression)
		return false
	}
	if !testLiteralExpression(t, infixExpression.Left, left) {
		return false
	}
	if infixExpression.Operator != operator {
		t.Fatalf("infixExpression.Operator is not '%v'. got=%v", operator, infixExpression.Operator)
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
	t.Errorf("type of expression not handled. got=%T", expression)
	return false
}
