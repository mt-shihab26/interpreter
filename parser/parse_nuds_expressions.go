package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

// registerNuds wires up the nud (null denotation) parse functions.
// Per Pratt's "Top Down Operator Precedence" paper, a nud is a token
// parsed with no left-hand expression -- our p.nuds map.
func (p *Parser) registerNuds() {
	p.nuds[token.MINUS] = p.parseUnaryExpression
	p.nuds[token.BANG] = p.parseUnaryExpression
	p.nuds[token.IDENT] = p.parseIdentifierExpression
	p.nuds[token.INT] = p.parseIntegerExpression
	p.nuds[token.TRUE] = p.parseBooleanExpression
	p.nuds[token.FALSE] = p.parseBooleanExpression
	p.nuds[token.IF] = p.parseIfExpression
	p.nuds[token.FUNCTION] = p.parseFunctionExpression
	p.nuds[token.LPAREN] = p.parseGroupedExpression
}

// parseUnaryExpression parses a "-<expression>" or "!<expression>" unary expression.
//
// It expects tokens on entry: -x  (curToken must be the operator, e.g. "-").
//
// It leaves curToken on the last token of the operand expression (e.g. "x").
func (p *Parser) parseUnaryExpression() ast.Expression {
	unaryExpression := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.advanceToken()
	unaryExpression.Right = p.parseExpression(UNARY)
	return unaryExpression
}

// parseIdentifierExpression parses a bare identifier, e.g. "foobar".
//
// It expects tokens on entry: foobar  (curToken must be the identifier).
//
// It does not advance -- curToken is left unchanged on the identifier.
func (p *Parser) parseIdentifierExpression() ast.Expression {
	identifierExpression := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	return identifierExpression
}

// parseIntegerExpression parses an integer literal, e.g. "5".
//
// It expects tokens on entry: 5  (curToken must be the INT token).
//
// It does not advance -- curToken is left unchanged on the integer literal.
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

// parseBooleanExpression parses a "true" or "false" literal.
//
// It expects tokens on entry: true  (curToken must be TRUE or FALSE).
//
// It does not advance -- curToken is left unchanged on the boolean literal.
func (p *Parser) parseBooleanExpression() ast.Expression {
	booleanExpression := &ast.BooleanExpression{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	return booleanExpression
}

// parseIfExpression parses an "if (<condition>) { <consequence> } else { <alternative> }" expression.
//
// It expects tokens on entry: if (x < y) { x } else { y }  (curToken must be "if").
//
// It leaves curToken on the closing "}" of whichever block was parsed last (the consequence if there's no "else", otherwise the alternative).
func (p *Parser) parseIfExpression() ast.Expression {
	ifExpression := &ast.IfExpression{Token: p.curToken}
	if !p.expectAdvancePeek(token.LPAREN) {
		return nil
	}
	p.advanceToken()
	ifExpression.Condition = p.parseExpression(LOWEST)
	if !p.expectAdvancePeek(token.RPAREN) {
		return nil
	}
	if !p.expectAdvancePeek(token.LBRACE) {
		return nil
	}
	ifExpression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.advanceToken()
		if !p.expectAdvancePeek(token.LBRACE) {
			return nil
		}
		ifExpression.Alternative = p.parseBlockStatement()
	}
	return ifExpression
}

// parseFunctionExpression parses a "fn(<name>, <name>, ...) { <body> }" function literal.
//
// It expects tokens on entry: fn(x, y) { x + y; }  (curToken must be "fn").
//
// It leaves curToken on the closing "}" of the function body.
func (p *Parser) parseFunctionExpression() ast.Expression {
	// fn (x, y) { x + y ;}
	functionExpression := &ast.FunctionExpression{Token: p.curToken}
	p.advanceToken()
	// (x, y) { x + y ;}
	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	p.advanceToken()
	// x, y) { x + y ;}
	functionExpression.Parameters = []*ast.IdentifierExpression{}
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		parameter := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
		functionExpression.Parameters = append(functionExpression.Parameters, parameter)
		p.advanceToken()
		if p.curTokenIs(token.RPAREN) {
			break
		}
		if !p.curTokenIs(token.COMMA) {
			return nil
		}
		p.advanceToken()
	}
	// ) { x + y ;}
	if !p.curTokenIs(token.RPAREN) {
		return nil
	}
	p.advanceToken()
	// { x + y ;}
	if !p.curTokenIs(token.LBRACE) {
		return nil
	}
	functionExpression.Body = p.parseBlockStatement()
	return functionExpression
}

// parseGroupedExpression parses a parenthesized "(<expression>)" expression.
//
// It expects tokens on entry: (x + y)  (curToken must be "(").
//
// It leaves curToken on the closing ")".
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceToken()
	insideGroupExpression := p.parseExpression(LOWEST)
	if !p.expectAdvancePeek(token.RPAREN) {
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
	p.advanceToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			blockStatement.Statements = append(blockStatement.Statements, statement)
		}
		p.advanceToken()
	}
	return blockStatement
}
