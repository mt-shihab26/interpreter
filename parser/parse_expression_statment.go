package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
)

// parseExpressionStatement parses a bare expression followed by an optional ";".
//
// It expects tokens on entry: x + y;  (curToken must be the first token of the expression, e.g. "x").
//
// It leaves curToken on the trailing ";" if present, otherwise on the expression's last token (e.g. "y").
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionStatement := &ast.ExpressionStatement{Token: p.curToken}
	expressionStatement.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.advanceToken()
	}
	return expressionStatement
}

// parseExpression runs the core Pratt parsing loop: it parses curToken's nud,
// then keeps folding in leds from peekToken as long as they bind tighter than precedence.
//
// It expects tokens on entry: a + b * c  (curToken must be the first token of the expression, e.g. "a").
//
// It leaves curToken on the last token of the parsed expression (e.g. "c") -- it does not consume a trailing ";".
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
