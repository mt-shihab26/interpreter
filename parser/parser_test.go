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
