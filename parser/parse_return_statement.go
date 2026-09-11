package parser

import (
	"monkey/ast"
	"monkey/token"
)

// parseReturnStatement parses a "return <expression>;" return statement.
//
// It expects tokens on entry: return 5;  (curToken must be "return").
//
// It leaves curToken on the trailing ";" if present, otherwise on the value expression's last token (e.g. "5").
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{Token: p.curToken}
	p.advanceToken()
	returnStatement.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.advanceToken()
	}
	return returnStatement
}
