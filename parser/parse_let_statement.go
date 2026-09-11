package parser

import (
	"monkey/ast"
	"monkey/token"
)

// parseLetStatement parses a "let <identifier expression> = <expression>;" let statement.
//
// It expects tokens on entry: let x = 5;  (curToken must be "let").
//
// It leaves curToken on the trailing ";" if present, otherwise on the value expression's last token (e.g. "5").
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
