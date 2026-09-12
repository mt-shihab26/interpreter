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
	debug.PrintProgram(&out, program, evaluated)
	t.Error(out.String())
}
