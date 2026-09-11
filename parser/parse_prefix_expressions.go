package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

func (p *Parser) parseUnaryExpression() ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteralExpression() ast.Expression {
	lit := &ast.Integer{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("cloud not parse %v as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	lit := &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	return lit
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	// fn (x, y) { x + y ;}
	expression := &ast.FunctionLiteral{Token: p.curToken}
	p.nextToken()
	// (x, y) { x + y ;}
	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	p.nextToken()
	// x, y) { x + y ;}
	expression.Parameters = []*ast.Identifier{}
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		expression.Parameters = append(expression.Parameters, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		p.nextToken()
		if p.curTokenIs(token.RPAREN) {
			break
		}
		if !p.curTokenIs(token.COMMA) {
			return nil
		}
		p.nextToken()
	}
	// ) { x + y ;}
	if !p.curTokenIs(token.RPAREN) {
		return nil
	}
	p.nextToken()
	// { x + y ;}
	if !p.curTokenIs(token.LBRACE) {
		return nil
	}
	expression.Body = p.parseBlockStatement()
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	expression := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expression
}
