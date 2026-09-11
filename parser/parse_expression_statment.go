package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionStatement := &ast.ExpressionStatement{Token: p.curToken}
	expressionStatement.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.advance()
	}
	return expressionStatement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParseFn := p.nuds[p.curToken.Type]
	if prefixParseFn == nil {
		p.noPrefixParseError(p.curToken.Type)
		return nil
	}
	leftExpression := prefixParseFn()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixParseFn := p.leds[p.peekToken.Type]
		if infixParseFn == nil {
			return leftExpression
		}
		p.advance()
		leftExpression = infixParseFn(leftExpression)
	}
	return leftExpression
}

func (p *Parser) noPrefixParseError(tokenType token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for %v found", tokenType)
	p.errors = append(p.errors, message)
}
