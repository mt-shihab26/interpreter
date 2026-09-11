package parser

import (
	"testing"

	"monkey/ast"
	"monkey/lexer"
)

// TestParseFullProgram feeds ParseProgram a single source that exercises every
// construct the parser supports (let/return statements, identifiers, integers,
// booleans, prefix and infix expressions, grouped expressions, if/else
// expressions, function literals and call expressions) and walks the
// resulting AST to check it was assembled correctly.
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
	testLiteralExpression(t, stmts[0].(*ast.LetStatement).Value, 5)

	// let ten = 10;
	if !testLetStatement(t, stmts[1], "ten") {
		return
	}
	testLiteralExpression(t, stmts[1].(*ast.LetStatement).Value, 10)

	// let add = fn(x, y) { return x + y; };
	if !testLetStatement(t, stmts[2], "add") {
		return
	}
	addFunction, ok := stmts[2].(*ast.LetStatement).Value.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("add's value is not *ast.FunctionExpression. got=%T\n", stmts[2].(*ast.LetStatement).Value)
	}
	if len(addFunction.Parameters) != 2 {
		t.Fatalf("addFunction.Parameters does not contain 2 parameters. got=%v\n", len(addFunction.Parameters))
	}
	testLiteralExpression(t, addFunction.Parameters[0], "x")
	testLiteralExpression(t, addFunction.Parameters[1], "y")
	if len(addFunction.Body.Statements) != 1 {
		t.Fatalf("addFunction.Body.Statements does not contain 1 statement. got=%v\n", len(addFunction.Body.Statements))
	}
	if !testReturnStatement(t, addFunction.Body.Statements[0]) {
		return
	}
	testBinaryExpression(t, addFunction.Body.Statements[0].(*ast.ReturnStatement).ReturnValue, "x", "+", "y")

	// let result = add(five, ten);
	if !testLetStatement(t, stmts[3], "result") {
		return
	}
	resultCall, ok := stmts[3].(*ast.LetStatement).Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf("result's value is not *ast.CallExpression. got=%T\n", stmts[3].(*ast.LetStatement).Value)
	}
	if !testIdentifierExpression(t, resultCall.Function, "add") {
		return
	}
	if len(resultCall.Arguments) != 2 {
		t.Fatalf("resultCall.Arguments does not contain 2 arguments. got=%v\n", len(resultCall.Arguments))
	}
	testIdentifierExpression(t, resultCall.Arguments[0], "five")
	testIdentifierExpression(t, resultCall.Arguments[1], "ten")

	// !-five;
	bangMinusFive := testExpressionStatement(t, stmts[4])
	outerBang, ok := bangMinusFive.Expression.(*ast.UnaryExpression)
	if !ok {
		t.Fatalf("expression is not *ast.UnaryExpression. got=%T\n", bangMinusFive.Expression)
	}
	if outerBang.Operator != "!" {
		t.Fatalf("outerBang.Operator is not '!'. got=%v\n", outerBang.Operator)
	}
	innerMinus, ok := outerBang.Right.(*ast.UnaryExpression)
	if !ok {
		t.Fatalf("outerBang.Right is not *ast.UnaryExpression. got=%T\n", outerBang.Right)
	}
	if innerMinus.Operator != "-" {
		t.Fatalf("innerMinus.Operator is not '-'. got=%v\n", innerMinus.Operator)
	}
	testIdentifierExpression(t, innerMinus.Right, "five")

	// five < ten > five;
	chainedComparison := testExpressionStatement(t, stmts[5])
	outerComparison, ok := chainedComparison.Expression.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("expression is not *ast.BinaryExpression. got=%T\n", chainedComparison.Expression)
	}
	if outerComparison.Operator != ">" {
		t.Fatalf("outerComparison.Operator is not '>'. got=%v\n", outerComparison.Operator)
	}
	testBinaryExpression(t, outerComparison.Left, "five", "<", "ten")
	testIdentifierExpression(t, outerComparison.Right, "five")

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
	testBinaryExpression(t, ifElseExpression.Condition, "five", "<", "ten")
	if len(ifElseExpression.Consequence.Statements) != 1 {
		t.Fatalf("ifElseExpression.Consequence.Statements does not contain 1 statement. got=%v\n", len(ifElseExpression.Consequence.Statements))
	}
	if !testReturnStatement(t, ifElseExpression.Consequence.Statements[0]) {
		return
	}
	testIdentifierExpression(t, ifElseExpression.Consequence.Statements[0].(*ast.ReturnStatement).ReturnValue, "five")
	if ifElseExpression.Alternative == nil {
		t.Fatalf("ifElseExpression.Alternative is nil\n")
	}
	if len(ifElseExpression.Alternative.Statements) != 1 {
		t.Fatalf("ifElseExpression.Alternative.Statements does not contain 1 statement. got=%v\n", len(ifElseExpression.Alternative.Statements))
	}
	if !testReturnStatement(t, ifElseExpression.Alternative.Statements[0]) {
		return
	}
	testIdentifierExpression(t, ifElseExpression.Alternative.Statements[0].(*ast.ReturnStatement).ReturnValue, "ten")

	// let max = fn(a, b) { if (a > b) { return a; } return b; };
	if !testLetStatement(t, stmts[12], "max") {
		return
	}
	maxFunction, ok := stmts[12].(*ast.LetStatement).Value.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("max's value is not *ast.FunctionExpression. got=%T\n", stmts[12].(*ast.LetStatement).Value)
	}
	if len(maxFunction.Parameters) != 2 {
		t.Fatalf("maxFunction.Parameters does not contain 2 parameters. got=%v\n", len(maxFunction.Parameters))
	}
	testLiteralExpression(t, maxFunction.Parameters[0], "a")
	testLiteralExpression(t, maxFunction.Parameters[1], "b")
	if len(maxFunction.Body.Statements) != 2 {
		t.Fatalf("maxFunction.Body.Statements does not contain 2 statements. got=%v\n", len(maxFunction.Body.Statements))
	}
	nestedIf, ok := testExpressionStatement(t, maxFunction.Body.Statements[0]).Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expression is not *ast.IfExpression. got=%T\n", testExpressionStatement(t, maxFunction.Body.Statements[0]).Expression)
	}
	testBinaryExpression(t, nestedIf.Condition, "a", ">", "b")
	if len(nestedIf.Consequence.Statements) != 1 {
		t.Fatalf("nestedIf.Consequence.Statements does not contain 1 statement. got=%v\n", len(nestedIf.Consequence.Statements))
	}
	if !testReturnStatement(t, nestedIf.Consequence.Statements[0]) {
		return
	}
	testIdentifierExpression(t, nestedIf.Consequence.Statements[0].(*ast.ReturnStatement).ReturnValue, "a")
	if nestedIf.Alternative != nil {
		t.Fatalf("nestedIf.Alternative was not nil. got=%v\n", nestedIf.Alternative)
	}
	if !testReturnStatement(t, maxFunction.Body.Statements[1]) {
		return
	}
	testIdentifierExpression(t, maxFunction.Body.Statements[1].(*ast.ReturnStatement).ReturnValue, "b")

	// max(five * 2, (ten + five) / 3);
	finalCall, ok := testExpressionStatement(t, stmts[13]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, stmts[13]).Expression)
	}
	if !testIdentifierExpression(t, finalCall.Function, "max") {
		return
	}
	if len(finalCall.Arguments) != 2 {
		t.Fatalf("finalCall.Arguments does not contain 2 arguments. got=%v\n", len(finalCall.Arguments))
	}
	testBinaryExpression(t, finalCall.Arguments[0], "five", "*", 2)
	groupedDivision, ok := finalCall.Arguments[1].(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("finalCall.Arguments[1] is not *ast.BinaryExpression. got=%T\n", finalCall.Arguments[1])
	}
	if groupedDivision.Operator != "/" {
		t.Fatalf("groupedDivision.Operator is not '/'. got=%v\n", groupedDivision.Operator)
	}
	testBinaryExpression(t, groupedDivision.Left, "ten", "+", "five")
	testIntegerExpression(t, groupedDivision.Right, 3)
}

// TestParseEmptyProgram checks that an empty (or whitespace-only) source
// produces zero statements and zero errors, rather than nil-panicking.
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

// TestParseStatementsWithoutTrailingSemicolon checks that the trailing ";" is
// optional on the last statement of the input, since every statement parser
// only consumes it when present.
func TestParseStatementsWithoutTrailingSemicolon(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"let x = 5", "let x = 5;"},
		{"return 5", "return 5;"},
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

// TestParseEmptyFunctionLiteral checks a function literal with no parameters
// and an empty body.
func TestParseEmptyFunctionLiteral(t *testing.T) {
	program := testParseProgram(t, "fn() {};", 1)
	functionExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("expression is not *ast.FunctionExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if len(functionExpression.Parameters) != 0 {
		t.Fatalf("functionExpression.Parameters is not empty. got=%v\n", len(functionExpression.Parameters))
	}
	if len(functionExpression.Body.Statements) != 0 {
		t.Fatalf("functionExpression.Body.Statements is not empty. got=%v\n", len(functionExpression.Body.Statements))
	}
}

// TestParseEmptyCallArguments checks a call expression with no arguments.
func TestParseEmptyCallArguments(t *testing.T) {
	program := testParseProgram(t, "foo();", 1)
	callExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if !testIdentifierExpression(t, callExpression.Function, "foo") {
		return
	}
	if len(callExpression.Arguments) != 0 {
		t.Fatalf("callExpression.Arguments is not empty. got=%v\n", len(callExpression.Arguments))
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
	if !testIdentifierExpression(t, ifExpression.Condition, "x") {
		return
	}
	if len(ifExpression.Consequence.Statements) != 0 {
		t.Fatalf("ifExpression.Consequence.Statements is not empty. got=%v\n", len(ifExpression.Consequence.Statements))
	}
	if ifExpression.Alternative != nil {
		t.Fatalf("ifExpression.Alternative was not nil. got=%v\n", ifExpression.Alternative)
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
	testUnaryExpression(t, outer.Right, "-", 5)
}

// TestParseImmediatelyInvokedFunctionExpression checks that a call
// expression's callee can be a function literal, not just an identifier.
func TestParseImmediatelyInvokedFunctionExpression(t *testing.T) {
	program := testParseProgram(t, "fn(x) { x; }(5);", 1)
	callExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	functionExpression, ok := callExpression.Function.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("callExpression.Function is not *ast.FunctionExpression. got=%T\n", callExpression.Function)
	}
	testLiteralExpression(t, functionExpression.Parameters[0], "x")
	if len(callExpression.Arguments) != 1 {
		t.Fatalf("callExpression.Arguments does not contain 1 argument. got=%v\n", len(callExpression.Arguments))
	}
	testIntegerExpression(t, callExpression.Arguments[0], 5)
}

// TestParseChainedCallExpressions checks that a call's callee can itself be
// a call expression, e.g. curried invocation "add(1)(2)".
func TestParseChainedCallExpressions(t *testing.T) {
	program := testParseProgram(t, "add(1)(2);", 1)
	outerCall, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if len(outerCall.Arguments) != 1 {
		t.Fatalf("outerCall.Arguments does not contain 1 argument. got=%v\n", len(outerCall.Arguments))
	}
	testIntegerExpression(t, outerCall.Arguments[0], 2)
	innerCall, ok := outerCall.Function.(*ast.CallExpression)
	if !ok {
		t.Fatalf("outerCall.Function is not *ast.CallExpression. got=%T\n", outerCall.Function)
	}
	if !testIdentifierExpression(t, innerCall.Function, "add") {
		return
	}
	if len(innerCall.Arguments) != 1 {
		t.Fatalf("innerCall.Arguments does not contain 1 argument. got=%v\n", len(innerCall.Arguments))
	}
	testIntegerExpression(t, innerCall.Arguments[0], 1)
}

// TestParseTrailingCommaIsTolerated documents that the parameter/argument
// loops accept (and silently ignore) a trailing comma before the closing
// delimiter, since they only require a COMMA between two identifiers/
// expressions rather than rejecting one right before RPAREN.
func TestParseTrailingCommaIsTolerated(t *testing.T) {
	program := testParseProgram(t, "fn(x, y,) {};", 1)
	functionExpression, ok := testExpressionStatement(t, program.Statements[0]).Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("expression is not *ast.FunctionExpression. got=%T\n", testExpressionStatement(t, program.Statements[0]).Expression)
	}
	if len(functionExpression.Parameters) != 2 {
		t.Fatalf("functionExpression.Parameters does not contain 2 parameters. got=%v\n", len(functionExpression.Parameters))
	}
	testLiteralExpression(t, functionExpression.Parameters[0], "x")
	testLiteralExpression(t, functionExpression.Parameters[1], "y")

	callProgram := testParseProgram(t, "foo(1, 2,);", 1)
	callExpression, ok := testExpressionStatement(t, callProgram.Statements[0]).Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expression is not *ast.CallExpression. got=%T\n", testExpressionStatement(t, callProgram.Statements[0]).Expression)
	}
	if len(callExpression.Arguments) != 2 {
		t.Fatalf("callExpression.Arguments does not contain 2 arguments. got=%v\n", len(callExpression.Arguments))
	}
	testIntegerExpression(t, callExpression.Arguments[0], 1)
	testIntegerExpression(t, callExpression.Arguments[1], 2)
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

// TestParseMalformedLetStatementDoesNotPanic exercises every way a "let"
// statement can be malformed (missing identifier, missing "=", missing
// value). Each case must record a parser error and must not panic: a
// dropped/broken let statement is a concrete *ast.LetStatement returned as
// nil from parseLetStatement, and parseStatement must convert that into a
// true nil ast.Statement -- not a non-nil interface wrapping a typed nil --
// or later calls like .String() on it will nil-dereference and panic.
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
			_ = statement.String() // must not panic
		}
	}
}

// TestParseUnterminatedFunctionParametersIsSilentlyDropped documents a known
// parser limitation: unlike grouped expressions and if conditions,
// parseFunctionExpression and parseCallExpression return nil (as the
// ast.Expression interface, so no typed-nil trap here) without recording an
// error when their parameter/argument list runs into EOF instead of a
// closing ")". Parsing still completes without panicking, but the caller
// gets no diagnostic for genuinely malformed input.
func TestParseUnterminatedFunctionParametersIsSilentlyDropped(t *testing.T) {
	tests := []string{
		"fn(x, y",
		"foo(1, 2",
	}
	for _, input := range tests {
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
}
