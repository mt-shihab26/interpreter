package output

import (
	"bytes"
	"testing"

	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

func TestNew(t *testing.T) {
	writer := &bytes.Buffer{}
	program := parseProgram(t, "let x = 5;")
	evaluated := &object.Integer{Value: 5}
	output := New(writer, program, evaluated, true)
	if output.Writer != writer {
		t.Errorf("New() Writer = %v, expected %v", output.Writer, writer)
	}
	if output.Program != program {
		t.Errorf("New() Program = %v, expected %v", output.Program, program)
	}
	if output.Evaluated != evaluated {
		t.Errorf("New() Evaluated = %v, expected %v", output.Evaluated, evaluated)
	}
	if !output.Verbose {
		t.Error("New() Verbose = false, expected true")
	}
}

func TestPrint(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		evaluated object.Object
		verbose   bool
		expected  string
	}{
		{
			name:      "normal output",
			code:      "5",
			evaluated: &object.Integer{Value: 5},
			expected:  "5",
		},
		{
			name:      "normal output with nil evaluated",
			code:      "let x = 5;",
			evaluated: nil,
			expected:  "",
		},
		{
			name:      "verbose output",
			code:      "5",
			evaluated: &object.Integer{Value: 5},
			verbose:   true,
			expected: `---CODE---
5
---AST---
Program
└─ Statement[0]: ExpressionStatement
   └─ Expression: IntegerExpression 5

---OUT---
5
`,
		},
		{
			name:      "verbose output with nil evaluated",
			code:      "let x = 5;",
			evaluated: nil,
			verbose:   true,
			expected: `---CODE---
let x = 5;
---AST---
Program
└─ Statement[0]: LetStatement
   ├─ IdentifierExpression: IdentifierExpression "x"
   └─ ValueExpression: IntegerExpression 5
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer bytes.Buffer
			program := parseProgram(t, tt.code)
			output := New(&writer, program, tt.evaluated, tt.verbose)
			output.Print()
			if got := writer.String(); got != tt.expected {
				t.Errorf("Print() = %q, expected %q", got, tt.expected)
			}
		})
	}
}

func TestPrintNormal(t *testing.T) {
	tests := []struct {
		name      string
		evaluated object.Object
		expected  string
	}{
		{
			name:      "integer",
			evaluated: &object.Integer{Value: 5},
			expected:  "5",
		},
		{
			name:      "boolean",
			evaluated: &object.Boolean{Value: true},
			expected:  "true",
		},
		{
			name:      "string",
			evaluated: &object.String{Value: "hello"},
			expected:  "hello",
		},
		{
			name:      "nil evaluated",
			evaluated: nil,
			expected:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer bytes.Buffer
			output := New(&writer, nil, tt.evaluated, false)
			output.printNormal()
			if got := writer.String(); got != tt.expected {
				t.Errorf("printNormal() = %q, expected %q", got, tt.expected)
			}
		})
	}
}

func TestPrintVerbose(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		evaluated object.Object
		expected  string
	}{
		{
			name:      "integer output",
			code:      "5",
			evaluated: &object.Integer{Value: 5},
			expected: `---CODE---
5
---AST---
Program
└─ Statement[0]: ExpressionStatement
   └─ Expression: IntegerExpression 5

---OUT---
5
`,
		},
		{
			name:      "nil evaluated",
			code:      "let x = 5;",
			evaluated: nil,
			expected: `---CODE---
let x = 5;
---AST---
Program
└─ Statement[0]: LetStatement
   ├─ IdentifierExpression: IdentifierExpression "x"
   └─ ValueExpression: IntegerExpression 5
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var writer bytes.Buffer
			program := parseProgram(t, tt.code)
			output := New(&writer, program, tt.evaluated, true)
			output.printVerbose()
			if got := writer.String(); got != tt.expected {
				t.Errorf("printVerbose() = %q, expected %q", got, tt.expected)
			}
		})
	}
}

func parseProgram(t *testing.T, input string) *ast.Program {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	return program
}
