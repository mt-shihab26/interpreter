package parser

import (
	"monkey/ast"
	"monkey/token"
)

// registerLeds wires up the led (left denotation) parse functions.
// Per Pratt's "Top Down Operator Precedence" paper, a led is a token
// parsed given an already-parsed left-hand expression -- our p.leds map.
func (p *Parser) registerLeds() {
	p.leds[token.PLUS] = p.parseBinaryExpression
	p.leds[token.MINUS] = p.parseBinaryExpression
	p.leds[token.ASTERISK] = p.parseBinaryExpression
	p.leds[token.SLASH] = p.parseBinaryExpression
	p.leds[token.EQUAL] = p.parseBinaryExpression
	p.leds[token.NOT_EQUAL] = p.parseBinaryExpression
	p.leds[token.GREATER_THAN] = p.parseBinaryExpression
	p.leds[token.LESS_THAN] = p.parseBinaryExpression
	p.leds[token.LEFT_PAREN] = p.parseCallExpression
}

// parseBinaryExpression parses a "<left expression> <operator> <right expression>" binary expression.
//
// It expects tokens on entry: a + b  (curToken must be the operator, e.g. "+").
//
// It leaves curToken on the last token of the right-hand expression (e.g. "b").
func (p *Parser) parseBinaryExpression(leftExpression ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:          p.curToken,
		Operator:       p.curToken.Literal,
		LeftExpression: leftExpression,
	}
	precedence := p.curPrecedence()
	p.advanceToken()
	expression.RightExpression = p.parseExpression(precedence)
	return expression
}

// parseCallExpression parses a "<callee expression>(<arguments>)" call expression.
//
// It expects tokens on entry: (a, b)  curToken must be "(".
//
// It leaves curToken on the closing ")" -- it does not consume the ")".
func (p *Parser) parseCallExpression(calleeExpression ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Token: p.curToken, FunctionExpression: calleeExpression}
	p.advanceToken()
	callExpression.ArgumentExpressions = []ast.Expression{}
	for !p.curTokenIs(token.RIGHT_PAREN) && !p.curTokenIs(token.EOF) {
		callExpression.ArgumentExpressions = append(callExpression.ArgumentExpressions, p.parseExpression(LOWEST))
		p.advanceToken()
		if p.curTokenIs(token.RIGHT_PAREN) {
			break
		}
		if !p.curTokenIs(token.COMMA) {
			return nil
		}
		p.advanceToken()
	}
	return callExpression
}
