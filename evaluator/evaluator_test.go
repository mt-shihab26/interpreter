package evaluator

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
	return program, Eval(program)
}

func printDebugInfo(t *testing.T, program *ast.Program, evaluated object.Object) {
	var out bytes.Buffer
	out.WriteString("\n")
	out.WriteString("----------------------START DEBUG-----------------------\n")
	debug.PrintProgram(&out, program, evaluated)
	out.WriteString("----------------------END DEBUG-------------------------\n")
	t.Error(out.String())
}
