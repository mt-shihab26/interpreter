package token

import "testing"

func TestLookupKeyword(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		expected Type
	}{
		{"let keyword", "let", LET},
		{"return keyword", "return", RETURN},
		{"if keyword", "if", IF},
		{"else keyword", "else", ELSE},
		{"function keyword", "fn", FUNCTION},
		{"true keyword", "true", TRUE},
		{"false keyword", "false", FALSE},
		{"identifier", "foobar", IDENTIFIER},
		{"numeric string", "1234", IDENTIFIER},
		{"single character identifier", "x", IDENTIFIER},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenType := LookupKeyword(test.word)
			if tokenType != test.expected {
				t.Errorf("LookupKeyword(%q) = %v, expected %v", test.word, tokenType, test.expected)
			}
		})
	}
}
