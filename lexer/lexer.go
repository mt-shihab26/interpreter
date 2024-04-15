package lexer

import (
	"monkey/token"
	"strings"
)

type Lexer struct {
	cursor int
	tokens []token.Token
}

func (l *Lexer) NextToken() token.Token {
	tok := l.tokens[l.cursor]
	l.cursor++
	return tok
}

func New(sourceCode string) Lexer {
	return Lexer{
		tokens: parseTokens(sourceCode),
	}
}

func parseTokens(sourceCode string) []token.Token {
	var tokens []token.Token
	characters := strings.Split(sourceCode, "")
	for _, character := range characters {
		tokens = append(tokens, token.Token{Type: tokenType(character), Literal: character})
	}
	return tokens
}

func tokenType(character string) token.TokenType {
	if character == token.ASSIGN {
		return token.ASSIGN
	}
	return token.ILLEGAL
}
