package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{Token: p.curToken}
	if !p.expectAdvancePeek(token.IDENT) {
		return nil
	}
	letStatement.Name = &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectAdvancePeek(token.ASSIGN) {
		return nil
	}
	p.advanceToken()
	letStatement.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.advanceToken()
	}
	return letStatement
}
