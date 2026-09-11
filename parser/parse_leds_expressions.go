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
	p.leds[token.EQ] = p.parseBinaryExpression
	p.leds[token.NOT_EQ] = p.parseBinaryExpression
	p.leds[token.GT] = p.parseBinaryExpression
	p.leds[token.LT] = p.parseBinaryExpression
	p.leds[token.LPAREN] = p.parseCallExpression
}

// parseBinaryExpression parses a "a + b" binary expression.
//
// It expects tokens on entry: a + b  (curToken must be the operator, e.g. "+").
//
// It leaves curToken on the last token of the right-hand expression (e.g. "b").
func (p *Parser) parseBinaryExpression(leftExpression ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     leftExpression,
	}
	precedence := p.curPrecedence()
	p.advanceToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

// parseCallExpression parses a "foo(a, b)" call expression.
//
// It expects tokens on entry: (a, b)  curToken must be "(".
//
// It leaves curToken on the closing ")" -- it does not consume the ")".
func (p *Parser) parseCallExpression(functionExpression ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Token: p.curToken, Function: functionExpression}
	p.advanceToken()
	callExpression.Arguments = []ast.Expression{}
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		callExpression.Arguments = append(callExpression.Arguments, p.parseExpression(LOWEST))
		p.advanceToken()
		if p.curTokenIs(token.RPAREN) {
			break
		}
		if !p.curTokenIs(token.COMMA) {
			return nil
		}
		p.advanceToken()
	}
	return callExpression
}
