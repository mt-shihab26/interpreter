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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program := testParseProgram(t, input, 1)
	expressionStatement := testExpressionStatement(t, program.Statements[0])
	testIntegerLiteralExpression(t, expressionStatement.Expression, 5)
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		right    int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
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
		left     int64
		operator string
		right    int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
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
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, test.statementCount)
		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected=%q, got=%q", test.expected, actual)
		}
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
		t.Fatalf("program.Statements does not contain %d statements. got=%d", statementsCount, len(program.Statements))
	}
	return program
}

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement not *ast.LetStatement. got=%T", statement)
		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral() not '%s'. got=%s", name, letStatement.Name.TokenLiteral())
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
		t.Errorf("returnStmt.TokenLiteral not 'return'. got=%q", returnStatement.TokenLiteral())
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
		t.Errorf("identifierExpression.Value not %s. got=%s", value, identifierExpression.Value)
		return false
	}
	if identifierExpression.TokenLiteral() != value {
		t.Errorf("identifierExpression.TokenLiteral not %s. got=%s", value, identifierExpression.TokenLiteral())
		return false
	}
	return true
}

func testIntegerLiteralExpression(t *testing.T, il ast.Expression, value int64) bool {
	integerLiteralExpression, ok := il.(*ast.IntegralLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegralLiteral. got=%T", il)
		return false
	}
	if integerLiteralExpression.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integerLiteralExpression.Value)
		return false
	}
	if integerLiteralExpression.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integerLiteralExpression.TokenLiteral())
		return false
	}
	return true
}

func testPrefixExpression(t *testing.T, expression ast.Expression, operator string, right any) bool {
	prefixExpression, ok := expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("stmt is not *ast.PrefixExpression. got=%T", expression)
		return false
	}
	if prefixExpression.Operator != operator {
		t.Fatalf("exp.Operator is not '%s'. got=%s", operator, prefixExpression.Operator)
		return false
	}
	if !testLiteralExpression(t, prefixExpression.Right, right) {
		return false
	}
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not *ast.InfixExpression. got=%T", exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Fatalf("exp.Operator is not '%s'. got=%s", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteralExpression(t, exp, int64(v))
	case int64:
		return testIntegerLiteralExpression(t, exp, v)
	case string:
		return testIdentifierExpression(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
