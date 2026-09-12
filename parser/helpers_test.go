package parser

import (
	"fmt"
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

// testParseProgram parses input, fails the test on parser errors or a wrong statement count, and returns the resulting program.
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

// checkParserErrors fails the test immediately, logging every accumulated parser error, if parser has any.
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

// testLetStatement checks that statement is a *ast.LetStatement binding the given identifier name.
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
	if letStatement.IdentifierExpression.Value != name {
		t.Errorf("letStatement.Name.Value not %v. got=%v\n", name, letStatement.IdentifierExpression.Value)
		return false
	}
	if letStatement.IdentifierExpression.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral() not %v. got=%v\n", name, letStatement.IdentifierExpression.TokenLiteral())
		return false
	}
	return true
}

// testReturnStatement checks that statement is a *ast.ReturnStatement.
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

// testExpressionStatement asserts that statement is a *ast.ExpressionStatement and returns it.
func testExpressionStatement(t *testing.T, statement ast.Statement) *ast.ExpressionStatement {
	expressionStatement, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not *ast.ExpressionStatement. got=%T\n", statement)
	}
	return expressionStatement
}

// testIdentifierExpression checks that expression is a *ast.IdentifierExpression with the given value.
func testIdentifierExpression(t *testing.T, expression ast.Expression, value string) bool {
	identifierExpression, ok := expression.(*ast.IdentifierExpression)
	if !ok {
		t.Errorf("expression is not *ast.IdentifierExpression. got=%T\n", expression)
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

// testIntegerExpression checks that expression is a *ast.IntegerExpression with the given value.
func testIntegerExpression(t *testing.T, expression ast.Expression, value int64) bool {
	integerExpression, ok := expression.(*ast.IntegerExpression)
	if !ok {
		t.Errorf("expression is not *ast.IntegerExpression. got=%T\n", expression)
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

// testBooleanExpression checks that expression is a *ast.BooleanExpression with the given value.
func testBooleanExpression(t *testing.T, expression ast.Expression, value bool) bool {
	booleanExpression, ok := expression.(*ast.BooleanExpression)
	if !ok {
		t.Errorf("expression is not *ast.BooleanExpression. got=%T\n", expression)
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

// testUnaryExpression checks that expression is a *ast.UnaryExpression with the given operator and operand.
func testUnaryExpression(t *testing.T, expression ast.Expression, operator string, right any) bool {
	unaryExpression, ok := expression.(*ast.UnaryExpression)
	if !ok {
		t.Fatalf("expression is not *ast.UnaryExpression. got=%T\n", expression)
		return false
	}
	if unaryExpression.Operator != operator {
		t.Fatalf("unaryExpression.Operator is not %v. got=%v\n", operator, unaryExpression.Operator)
		return false
	}
	if !testLiteralExpression(t, unaryExpression.RightExpression, right) {
		return false
	}
	return true
}

// testBinaryExpression checks that expression is a *ast.BinaryExpression with the given left operand, operator, and right operand.
func testBinaryExpression(t *testing.T, expression ast.Expression, left any, operator string, right any) bool {
	binaryExpression, ok := expression.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("expression is not *ast.BinaryExpression. got=%T\n", expression)
		return false
	}
	if !testLiteralExpression(t, binaryExpression.LeftExpression, left) {
		return false
	}
	if binaryExpression.Operator != operator {
		t.Fatalf("binaryExpression.Operator is not %v. got=%v\n", operator, binaryExpression.Operator)
		return false
	}
	if !testLiteralExpression(t, binaryExpression.RightExpression, right) {
		return false
	}
	return true
}

// testLiteralExpression dispatches to the matching testXExpression helper based on expected's Go type.
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
