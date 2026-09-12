package evaluator

import (
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
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
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
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
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

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}
