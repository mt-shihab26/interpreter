package lexer_test

import (
	"fmt"
	"monkey/lexer"
	"monkey/token"
	"testing"
)

type nextTokeTestTable struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tt := []nextTokeTestTable{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}
	runNextTokenTest(t, input, tt)
}

func TestNextToken2(t *testing.T) {
	input := `let five = 5;
let ten = 10;
   let add = fn(x, y) {
     x + y;
};
   let result = add(five, ten);
   `
	tt := []nextTokeTestTable{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	runNextTokenTest(t, input, tt)
}

func runNextTokenTest(t *testing.T, input string, tt []nextTokeTestTable) {
	l := lexer.New(input)

	for i, tc := range tt {
		t.Run(fmt.Sprintf("tests[%d]:%s", i, tc.expectedType), func(t *testing.T) {
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
