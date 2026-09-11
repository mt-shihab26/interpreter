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
		p.advanceToken()
	}
	return expressionStatement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	nudFunc := p.nuds[p.curToken.Type]
	if nudFunc == nil {
		message := fmt.Sprintf("no prefix parse function for %v found", p.curToken.Type)
		p.errors = append(p.errors, message)
		return nil
	}
	leftExpression := nudFunc()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		ledFunc := p.leds[p.peekToken.Type]
		if ledFunc == nil {
			return leftExpression
		}
		p.advanceToken()
		leftExpression = ledFunc(leftExpression)
	}
	return leftExpression
}
