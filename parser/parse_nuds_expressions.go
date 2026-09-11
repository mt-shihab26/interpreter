package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

func (p *Parser) registerNuds() {
	p.registerNud(token.MINUS, p.parseUnaryExpression)
	p.registerNud(token.BANG, p.parseUnaryExpression)
	p.registerNud(token.IDENT, p.parseIdentifierExpression)
	p.registerNud(token.INT, p.parseIntegerExpression)
	p.registerNud(token.TRUE, p.parseBooleanExpression)
	p.registerNud(token.FALSE, p.parseBooleanExpression)
	p.registerNud(token.IF, p.parseIfExpression)
	p.registerNud(token.FUNCTION, p.parseFunctionExpression)
	p.registerNud(token.LPAREN, p.parseGroupedExpression)
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	unaryExpression := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	unaryExpression.Right = p.parseExpression(PREFIX)
	return unaryExpression
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	identifierExpression := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	return identifierExpression
}

func (p *Parser) parseIntegerExpression() ast.Expression {
	integerExpression := &ast.IntegerExpression{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		message := fmt.Sprintf("cloud not parse %v as integer", p.curToken.Literal)
		p.errors = append(p.errors, message)
	}
	integerExpression.Value = value
	return integerExpression
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	booleanExpression := &ast.BooleanExpression{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	return booleanExpression
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	ifExpression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	ifExpression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		ifExpression.Alternative = p.parseBlockStatement()
	}
	return ifExpression
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	// fn (x, y) { x + y ;}
	functionExpression := &ast.FunctionExpression{Token: p.curToken}
	p.nextToken()
	// (x, y) { x + y ;}
	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	p.nextToken()
	// x, y) { x + y ;}
	functionExpression.Parameters = []*ast.IdentifierExpression{}
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		parameter := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
		functionExpression.Parameters = append(functionExpression.Parameters, parameter)
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
	functionExpression.Body = p.parseBlockStatement()
	return functionExpression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	insideGroupExpression := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return insideGroupExpression
}

// parseBlockStatement parses a "{ ... }" block statement.
//
// It expects tokens on entry: { x + y ;}  (curToken must be "{").
//
// It leaves curToken on the closing "}" -- it does not consume the "}".
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{Token: p.curToken}
	blockStatement.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			blockStatement.Statements = append(blockStatement.Statements, statement)
		}
		p.nextToken()
	}
	return blockStatement
}
