package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

// TestParseFullProgram walks the AST of a source exercising every construct the parser supports, checking it was assembled correctly.
func TestParseFullProgram(t *testing.T) {
	input := `
let five = 5;
let ten = 10;
let add = fn(x, y) { return x + y; };
let result = add(five, ten);

!-five;
five < ten > five;
five == five;
five != ten;
true;
false;
!true;

if (five < ten) {
	return five;
} else {
	return ten;
}

let max = fn(a, b) {
	if (a > b) {
		return a;
	}
	return b;
};

max(five * 2, (ten + five) / 3);
`

	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	const expectedStatements = 14
	if len(program.Statements) != expectedStatements {
		t.Fatalf("program.Statements does not contain %v statements. got=%v\n", expectedStatements, len(program.Statements))
	}
	stmts := program.Statements

	// let five = 5;
	if !testLetStatement(t, stmts[0], "five") {
		return
	}
	testLiteralExpression(t, stmts[0].(*ast.LetStatement).ValueExpression, 5)

	// let ten = 10;
	if !testLetStatement(t, stmts[1], "ten") {
		return
	}
	testLiteralExpression(t, stmts[1].(*ast.LetStatement).ValueExpression, 10)

	// let add = fn(x, y) { return x + y; };
	if !testLetStatement(t, stmts[2], "add") {
		return
	}
	addFunction, ok := stmts[2].(*ast.LetStatement).ValueExpression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("add's value is not *ast.FunctionExpression. got=%T\n", stmts[2].(*ast.LetStatement).ValueExpression)
	}
	if len(addFunction.ParameterExpressions) != 2 {
		t.Fatalf("addFunction.Parameters does not contain 2 parameters. got=%v\n", len(addFunction.ParameterExpressions))
	}
	testLiteralExpression(t, addFunction.ParameterExpressions[0], "x")
	testLiteralExpression(t, addFunction.ParameterExpressions[1], "y")
	if len(addFunction.BodyStatement.Statements) != 1 {
		t.Fatalf("addFunction.Body.Statements does not contain 1 statement. got=%v\n", len(addFunction.BodyStatement.Statements))
	}
	if !testReturnStatement(t, addFunction.BodyStatement.Statements[0]) {
		return
	}
	testBinaryExpression(t, addFunction.BodyStatement.Statements[0].(*ast.ReturnStatement).ValueExpression, "x", "+", "y")

	// let result = add(five, ten);
	if !testLetStatement(t, stmts[3], "result") {
		return
	}
	resultCall, ok := stmts[3].(*ast.LetStatement).ValueExpression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("result's value is not *ast.CallExpression. got=%T\n", stmts[3].(*ast.LetStatement).ValueExpression)
	}
	if !testIdentifierExpression(t, resultCall.FunctionExpression, "add") {
		return
	}
	if len(resultCall.ArgumentExpressions) != 2 {
		t.Fatalf("resultCall.Arguments does not contain 2 arguments. got=%v\n", len(resultCall.ArgumentExpressions))
	}
	testIdentifierExpression(t, resultCall.ArgumentExpressions[0], "five")
	testIdentifierExpression(t, resultCall.ArgumentExpressions[1], "ten")

	// !-five;
	bangMinusFive := testExpressionStatement(t, stmts[4])
	outerBang, ok := bangMinusFive.Expression.(*ast.UnaryExpression)
	if !ok {
		t.Fatalf("expression is not *ast.UnaryExpression. got=%T\n", bangMinusFive.Expression)
	}
	if outerBang.Operator != "!" {
		t.Fatalf("outerBang.Operator is not '!'. got=%v\n", outerBang.Operator)
	}
	innerMinus, ok := outerBang.RightExpression.(*ast.UnaryExpression)
	if !ok {
		t.Fatalf("outerBang.Right is not *ast.UnaryExpression. got=%T\n", outerBang.RightExpression)
	}
	if innerMinus.Operator != "-" {
		t.Fatalf("innerMinus.Operator is not '-'. got=%v\n", innerMinus.Operator)
	}
	testIdentifierExpression(t, innerMinus.RightExpression, "five")

	// five < ten > five;
	chainedComparison := testExpressionStatement(t, stmts[5])
	outerComparison, ok := chainedComparison.Expression.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("expression is not *ast.BinaryExpression. got=%T\n", chainedComparison.Expression)
	}
	if outerComparison.Operator != ">" {
		t.Fatalf("outerComparison.Operator is not '>'. got=%v\n", outerComparison.Operator)
	}
	testBinaryExpression(t, outerComparison.LeftExpression, "five", "<", "ten")
	testIdentifierExpression(t, outerComparison.RightExpression, "five")

	// five == five;
	testBinaryExpression(t, testExpressionStatement(t, stmts[6]).Expression, "five", "==", "five")

	// five != ten;
	testBinaryExpression(t, testExpressionStatement(t, stmts[7]).Expression, "five", "!=", "ten")

	// true;
	testBooleanExpression(t, testExpressionStatement(t, stmts[8]).Expression, true)

	// false;
	testBooleanExpression(t, testExpressionStatement(t, stmts[9]).Expression, false)

	// !true;
	testUnaryExpression(t, testExpressionStatement(t, stmts[10]).Expression, "!", true)

	// if (five < ten) { return five; } else { return ten; }
	ifElseExpression, ok := testExpressionStatement(t, stmts[11]).Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expression is not *ast.IfExpression. got=%T\n", testExpressionStatement(t, stmts[11]).Expression)
	}
	testBinaryExpression(t, ifElseExpression.ConditionExpression, "five", "<", "ten")
	if len(ifElseExpression.ConsequenceStatement.Statements) != 1 {
		t.Fatalf("ifElseExpression.Consequence.Statements does not contain 1 statement. got=%v\n", len(ifElseExpression.ConsequenceStatement.Statements))
	}
	if !testReturnStatement(t, ifElseExpression.ConsequenceStatement.Statements[0]) {
		return
	}
	testIdentifierExpression(t, ifElseExpression.ConsequenceStatement.Statements[0].(*ast.ReturnStatement).ValueExpression, "five")
	if ifElseExpression.AlternativeStatement == nil {
		t.Fatalf("ifElseExpression.Alternative is nil\n")
	}
	if len(ifElseExpression.AlternativeStatement.Statements) != 1 {
		t.Fatalf("ifElseExpression.Alternative.Statements does not contain 1 statement. got=%v\n", len(ifElseExpression.AlternativeStatement.Statements))
	}
	if !testReturnStatement(t, ifElseExpression.AlternativeStatement.Statements[0]) {
		return
	}
	testIdentifierExpression(t, ifElseExpression.AlternativeStatement.Statements[0].(*ast.ReturnStatement).ValueExpression, "ten")

	// let max = fn(a, b) { if (a > b) { return a; } return b; };
	if !testLetStatement(t, stmts[12], "max") {
		return
	}
	maxFunction, ok := stmts[12].(*ast.LetStatement).ValueExpression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("max's value is not *ast.FunctionExpression. got=%T\n", stmts[12].(*ast.LetStatement).ValueExpression)
	}
	if len(maxFunction.ParameterExpressions) != 2 {
		t.Fatalf("maxFunction.Parameters does not contain 2 parameters. got=%v\n", len(maxFunction.ParameterExpressions))
	}
	testLiteralExpression(t, maxFunction.ParameterExpressions[0], "a")
	testLiteralExpression(t, maxFunction.ParameterExpressions[1], "b")
	if len(maxFunction.BodyStatement.Statements) != 2 {
		t.Fatalf("maxFunction.Body.Statements does not contain 2 statements. got=%v\n", len(maxFunction.BodyStatement.Statements))
	}
	nestedIf, ok := testExpressionStatement(t, maxFunction.BodyStatement.Statements[0]).Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expression is not *ast.IfExpression. got=%T\n", testExpressionStatement(t, maxFunction.BodyStatement.Statements[0]).Expression)
	}
	testBinaryExpression(t, nestedIf.ConditionExpression, "a", ">", "b")
	if len(nestedIf.ConsequenceStatement.Statements) != 1 {
		t.Fatalf("nestedIf.Consequence.Statements does not contain 1 statement. got=%v\n", len(nestedIf.ConsequenceStatement.Statements))
	}
	if !testReturnStatement(t, nestedIf.ConsequenceStatement.Statements[0]) {
		return
	}
	testIdentifierExpression(t, nestedIf.ConsequenceStatement.Statements[0].(*ast.ReturnStatement).ValueExpression, "a")
	if nestedIf.AlternativeStatement != nil {
		t.Fatalf("nestedIf.Alternative was not nil. got=%v\n", nestedIf.AlternativeStatement)
	}
	if !testReturnStatement(t, maxFunction.BodyStatement.Statements[1]) {
		return
	}
	testIdentifierExpression(t, maxFunction.BodyStatement.Statements[1].(*ast.ReturnStatement).ValueExpression, "b")

	// max(five * 2, (ten + five) / 3);
	finalCall, ok := testExpressionStatement(t, stmts[13]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, stmts[13]).Expression)
	}
	if !testIdentifierExpression(t, finalCall.FunctionExpression, "max") {
		return
	}
	if len(finalCall.ArgumentExpressions) != 2 {
		t.Fatalf("finalCall.Arguments does not contain 2 arguments. got=%v\n", len(finalCall.ArgumentExpressions))
	}
	testBinaryExpression(t, finalCall.ArgumentExpressions[0], "five", "*", 2)
	groupedDivision, ok := finalCall.ArgumentExpressions[1].(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("finalCall.Arguments[1] is not *ast.BinaryExpression. got=%T\n", finalCall.ArgumentExpressions[1])
	}
	if groupedDivision.Operator != "/" {
		t.Fatalf("groupedDivision.Operator is not '/'. got=%v\n", groupedDivision.Operator)
	}
	testBinaryExpression(t, groupedDivision.LeftExpression, "ten", "+", "five")
	testIntegerExpression(t, groupedDivision.RightExpression, 3)
}

// TestParseEmptyProgram checks that an empty or whitespace-only source produces zero statements and zero errors.
func TestParseEmptyProgram(t *testing.T) {
	for _, input := range []string{"", "   ", "\n\n\t\n"} {
		parser := New(lexer.New(input))
		program := parser.ParseProgram()
		checkParserErrors(t, parser)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil for input %q\n", input)
		}
		if len(program.Statements) != 0 {
			t.Fatalf("program.Statements is not empty for input %q. got=%v\n", input, len(program.Statements))
		}
	}
}

// TestLetStatements checks that "let <identifier> = <value>;" statements parse with the right identifier and value.
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
		value := letStatement.(*ast.LetStatement).ValueExpression
		if !testLiteralExpression(t, value, test.value) {
			return
		}
	}
}

// TestLetStatementWithoutTrailingSemicolon checks that the trailing ";" is optional in a let statement.
func TestLetStatementWithoutTrailingSemicolon(t *testing.T) {
	program := testParseProgram(t, "let x = 5", 1)
	if actual := program.Code(); actual != "let x = 5;\n" {
		t.Errorf("expected=%v, got=%v\n", "let x = 5;\n", actual)
	}
}

// TestParseMalformedLetStatementDoesNotPanic checks that a malformed "let" statement records a parser error and never leaves behind a nil-panicking statement.
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
			_ = statement.Code() // must not panic
		}
	}
}

// TestReturnStatements checks that "return <value>;" statements parse with the right return value.
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input       string
		returnValue any
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		returnStatement := program.Statements[0]
		if !testReturnStatement(t, returnStatement) {
			return
		}
		value := returnStatement.(*ast.ReturnStatement).ValueExpression
		if !testLiteralExpression(t, value, test.returnValue) {
			return
		}
	}
}

// TestReturnStatementWithoutTrailingSemicolon checks that the trailing ";" is optional in a return statement.
func TestReturnStatementWithoutTrailingSemicolon(t *testing.T) {
	program := testParseProgram(t, "return 5", 1)
	if actual := program.Code(); actual != "return 5;\n" {
		t.Errorf("expected=%v, got=%v\n", "return 5;\n", actual)
	}
}

// TestOperatorPrecedenceParsing checks that expressions reparse (via String()) with parentheses reflecting the correct operator precedence.
func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input          string
		expected       string
		statementCount int
	}{
		{"-a * b", "((-a) * b);\n", 1},
		{"!-a", "(!(-a));\n", 1},
		{"a + b + c", "((a + b) + c);\n", 1},
		{"a + b - c", "((a + b) - c);\n", 1},
		{"a * b * c", "((a * b) * c);\n", 1},
		{"a * b / c", "((a * b) / c);\n", 1},
		{"a + b / c", "(a + (b / c));\n", 1},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f);\n", 1},
		{"3 + 4; -5 * 5", "(3 + 4);\n((-5) * 5);\n", 2},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4));\n", 1},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4));\n", 1},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));\n", 1},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));\n", 1},
		{"true", "true;\n", 1},
		{"false", "false;\n", 1},
		{"3 > 5 == false", "((3 > 5) == false);\n", 1},
		{"3 < 5 == true", "((3 < 5) == true);\n", 1},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4);\n", 1},
		{"(5 + 5) * 2", "((5 + 5) * 2);\n", 1},
		{"2 / (5 + 5)", "(2 / (5 + 5));\n", 1},
		{"-(5 + 5)", "(-(5 + 5));\n", 1},
		{"!(true == true)", "(!(true == true));\n", 1},
		{"a + add(b * c) + d", "((a + add((b * c))) + d);\n", 1},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));\n", 1},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g));\n", 1},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)", 1},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))", 1},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, test.statementCount)
		actual := program.Code()
		if actual != test.expected {
			t.Errorf("expected=%v, got=%v\n", test.expected, actual)
		}
	}
}

// TestExpressionStatementWithoutTrailingSemicolon checks that the trailing ";" is optional after a bare expression.
func TestExpressionStatementWithoutTrailingSemicolon(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5", "5;\n"},
		{"x + y", "(x + y);\n"},
	}
	for _, test := range tests {
		program := testParseProgram(t, test.input, 1)
		if actual := program.Code(); actual != test.expected {
			t.Errorf("input=%q expected=%v, got=%v\n", test.input, test.expected, actual)
		}
	}
}

// TestParseIllegalTokenRecordsError checks that a token the lexer can't classify records a parser error instead of panicking.
func TestParseIllegalTokenRecordsError(t *testing.T) {
	parser := New(lexer.New("@"))
	parser.ParseProgram()
	if len(parser.Errors()) == 0 {
		t.Fatalf("expected at least 1 error for an illegal token, got 0\n")
	}
}
