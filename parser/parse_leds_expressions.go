package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) registerLeds() {
	p.registerLed(token.PLUS, p.parseBinaryExpression)
	p.registerLed(token.MINUS, p.parseBinaryExpression)
	p.registerLed(token.ASTERISK, p.parseBinaryExpression)
	p.registerLed(token.SLASH, p.parseBinaryExpression)
	p.registerLed(token.EQ, p.parseBinaryExpression)
	p.registerLed(token.NOT_EQ, p.parseBinaryExpression)
	p.registerLed(token.GT, p.parseBinaryExpression)
	p.registerLed(token.LT, p.parseBinaryExpression)
	p.registerLed(token.LPAREN, p.parseCallExpression)
}

func (p *Parser) parseBinaryExpression(leftExpression ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     leftExpression,
	}
	precedence := p.curPrecedence()
	p.advanceToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseCallExpression(functionExpression ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Token: p.curToken, Function: functionExpression}
	callExpression.Arguments = p.parseCallArguments()
	return callExpression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	callArguments := []ast.Expression{}
	p.advanceToken()
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		callArguments = append(callArguments, p.parseExpression(LOWEST))
		p.advanceToken()
		if p.curTokenIs(token.RPAREN) {
			break
		}
		if !p.curTokenIs(token.COMMA) {
			return nil
		}
		p.advanceToken()
	}
	return callArguments
}
