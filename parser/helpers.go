package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
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

type (
	nudFunc func() ast.Expression
	ledFunc func(ast.Expression) ast.Expression
)

func (p *Parser) advance() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		p.advance()
		return true
	} else {
		message := fmt.Sprintf("expected next token to be %v, got %v instead", tokenType, p.peekToken.Type)
		p.errors = append(p.errors, message)
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerNud(tokenType token.TokenType, fn nudFunc) {
	p.nuds[tokenType] = fn
}

func (p *Parser) registerLed(tokenType token.TokenType, fn ledFunc) {
	p.leds[tokenType] = fn
}
