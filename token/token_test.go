package token

import "testing"

func TestLookupKeyword(t *testing.T) {
	tests := []struct {
		word      string
		tokenType Type
	}{
		{"let", LET},
		{"return", RETURN},
		{"if", IF},
		{"else", ELSE},
		{"fn", FUNCTION},
		{"true", TRUE},
		{"false", FALSE},
		{"foobar", IDENTIFIER},
		{"1234", IDENTIFIER},
		{"x", IDENTIFIER},
	}

	for _, test := range tests {
		tokenType := LookupKeyword(test.word)
		if tokenType != test.tokenType {
			t.Errorf("LookupKeyword should return %v, got=%v\n", test.tokenType, tokenType)
		}
	}
}
