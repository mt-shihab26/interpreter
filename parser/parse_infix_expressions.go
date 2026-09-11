package parser

import "monkey/ast"

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExpression := &ast.CallExpression{Token: p.curToken, Function: function}
	callExpression.Arguments = p.parseCallArguments()
	return callExpression
}
