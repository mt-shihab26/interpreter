package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
)

type (
	nudFuncType func() ast.Expression
	ledFuncType func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS         // ==
	COMPARISON     // < or >
	ADDITIVE       // +, -
	MULTIPLICATIVE // *, /
	UNARY          // -x or !x
	CALL           // myFunction(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       COMPARISON,
	token.GT:       COMPARISON,
	token.PLUS:     ADDITIVE,
	token.MINUS:    ADDITIVE,
	token.ASTERISK: MULTIPLICATIVE,
	token.SLASH:    MULTIPLICATIVE,
	token.LPAREN:   CALL,
}

// advanceToken shifts peekToken into curToken and reads a new peekToken from the lexer.
func (p *Parser) advanceToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// curTokenIs reports whether curToken is of the given type.
func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

// peekTokenIs reports whether peekToken is of the given type.
func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

// expectAdvancePeek advances past peekToken if it matches tokenType, else records a parser error.
func (p *Parser) expectAdvancePeek(tokenType token.TokenType) bool {
	if !p.peekTokenIs(tokenType) {
		message := fmt.Sprintf("expected token to be %v, got %v instead", tokenType, p.peekToken.Type)
		p.errors = append(p.errors, message)
		return false
	}
	p.advanceToken()
	return true
}

// peekPrecedence returns the binding precedence of peekToken, or LOWEST if it has none.
func (p *Parser) peekPrecedence() int {
	precedence, ok := precedences[p.peekToken.Type]
	if !ok {
		return LOWEST
	}
	return precedence
}

// curPrecedence returns the binding precedence of curToken, or LOWEST if it has none.
func (p *Parser) curPrecedence() int {
	precedence, ok := precedences[p.curToken.Type]
	if !ok {
		return LOWEST
	}
	return precedence
}
