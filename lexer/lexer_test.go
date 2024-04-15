package lexer_test

import (
	"fmt"
	"monkey/lexer"
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tt := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := lexer.New(input)

	for i, tc := range tt {
		t.Run(fmt.Sprintf("tests[%d]", i), func(t *testing.T) {
			tok := l.NextToken()
			if tok.Type != tc.expectedType {
				t.Errorf(
					"tests[%d] - tokentype wrong. expected=%q, got=%q",
					i,
					tc.expectedType,
					tok.Type,
				)
				return
			}
			if tok.Literal != tc.expectedLiteral {
				t.Fatalf(
					"tests[%d] - literal wrong. expected=%q, got=%q",
					i,
					tc.expectedLiteral,
					tok.Literal,
				)
			}
		})
	}
}
