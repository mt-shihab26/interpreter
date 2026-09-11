package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := "=+(){},;"
	tests := []struct {
		Type    token.TokenType
		Literal string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LEFT_PAREN, "("},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}
	l := New(input)
	for i, test := range tests {
		tok := l.NextToken()
		if tok.Type != test.Type {
			t.Fatalf("tests[%v] - tokentype wrong, expected=%v, got=%v", i, test.Type, tok.Type)
		}
		if tok.Literal != test.Literal {
			t.Fatalf("tests[%v] - literal wrong, expected=%v, got=%v", i, test.Literal, tok.Literal)
		}
	}
}

func TestNextTokenTwo(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
`
	tests := []struct {
		Type    token.TokenType
		Literal string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LEFT_PAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.INTEGER, "5"},
		{token.LESS_THAN, "<"},
		{token.INTEGER, "10"},
		{token.GREATER_THAN, ">"},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LEFT_PAREN, "("},
		{token.INTEGER, "5"},
		{token.LESS_THAN, "<"},
		{token.INTEGER, "10"},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.ELSE, "else"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.INTEGER, "10"},
		{token.EQUAL, "=="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.INTEGER, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INTEGER, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, test := range tests {
		tok := l.NextToken()
		if tok.Type != test.Type {
			t.Fatalf("tests[%v] - tokentype wrong, expected=%v, got=%v", i, test.Type, tok.Type)
		}
		if tok.Literal != test.Literal {
			t.Fatalf("tests[%v] - literal wrong, expected=%v, got=%v", i, test.Literal, tok.Literal)
		}
	}
}
