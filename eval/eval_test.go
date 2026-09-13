package eval

import (
	"bytes"
	"monkey/ast"
	"monkey/debug"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, test := range tests {
		program, evaluated := testEval(test.input)
		if !testIntegerObject(t, evaluated, test.expected) {
			printDebugInfo(t, program, evaluated)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"false == true", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, test := range tests {
		program, evaluated := testEval(test.input)
		if !testBooleanObject(t, evaluated, test.expected) {
			printDebugInfo(t, program, evaluated)
		}
	}
}

func TestUnaryExpressionBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, test := range tests {
		program, evaluated := testEval(test.input)
		if !testBooleanObject(t, evaluated, test.expected) {
			printDebugInfo(t, program, evaluated)
		}
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, test := range tests {
		program, evaluated := testEval(test.input)
		integer, ok := test.expected.(int)
		if !ok {
			if !testNullObject(t, evaluated) {
				printDebugInfo(t, program, evaluated)
			}
		} else {
			if !testIntegerObject(t, evaluated, int64(integer)) {
				printDebugInfo(t, program, evaluated)
			}
		}
	}
}

func TestReturnExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2; 9;", 2},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
	}
	for _, test := range tests {
		program, evaluated := testEval(test.input)
		if !testIntegerObject(t, evaluated, test.expected) {
			printDebugInfo(t, program, evaluated)
		}
	}
}

func TestErrorObject(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5+true", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
				return true + false;
				}
				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"foobar", "identifier not found: foobar"},
		{"let foobar = 5; foobar()", "identifier is not function: foobar"},
		{
			`
			let i = 5;
			let printNum = fn(i) {
				let j = 10;
				i;
			}
			printNum(10);
			j;
			`,
			"identifier not found: j",
		},
	}
	for _, test := range tests {
		program, evaluated := testEval(test.input)
		if !testErrorObject(t, evaluated, test.expected) {
			printDebugInfo(t, program, evaluated)
		}
	}
}

func TestLetStatments(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, test := range tests {
		program, evaluated := testEval(test.input)
		if !testIntegerObject(t, evaluated, test.expected) {
			printDebugInfo(t, program, evaluated)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	_, evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "{ (x + 2) }"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionCalls(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identify = fn(x) { x; }; identify(5);", 5},
		{"let identify = fn(x) { return x; }; identify(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
		{
			`
			let i = 5;
			let printNum = fn(i) {
				i;
			}
			printNum(10);
			i;
			`,
			5,
		},
	}
	for _, test := range tests {
		program, evaluated := testEval(test.input)
		if !testIntegerObject(t, evaluated, test.expected) {
			printDebugInfo(t, program, evaluated)
		}
	}
}

func testErrorObject(t *testing.T, objectValue object.Object, expected string) bool {
	error, ok := objectValue.(*object.Error)
	if !ok {
		t.Errorf("object is not error. got=%T\n", objectValue)
		return false
	}
	if error.Message != expected {
		t.Errorf("wrong error message. expected=%q, got=%q\n", expected, error.Message)
	}
	return true
}

func testNullObject(t *testing.T, objectValue object.Object) bool {
	_, ok := objectValue.(*object.Null)
	if !ok {
		t.Errorf("object is not null. got=%T\n", objectValue)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, objectValue object.Object, expected int64) bool {
	result, ok := objectValue.(*object.Integer)
	if !ok {
		t.Errorf("object is not integer object. got=%T\n", objectValue)
		return false
	}
	if result.Value != expected {
		t.Errorf("integer object has wrong value. got=%v, want=%v\n", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, objectValue object.Object, expected bool) bool {
	result, ok := objectValue.(*object.Boolean)
	if !ok {
		t.Errorf("object is not boolean object. got=%T\n", objectValue)
		return false
	}
	if result.Value != expected {
		t.Errorf("boolean object has wrong value. got=%v, want=%v\n", result.Value, expected)
		return false
	}
	return true
}

func testEval(input string) (*ast.Program, object.Object) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return program, Eval(program, env)
}

func printDebugInfo(t *testing.T, program *ast.Program, evaluated object.Object) {
	var out bytes.Buffer
	out.WriteString("\n")
	out.WriteString("----------------------START DEBUG-----------------------\n")
	debug.PrintProgram(&out, program, evaluated)
	out.WriteString("----------------------END DEBUG-------------------------\n")
	t.Error(out.String())
}
