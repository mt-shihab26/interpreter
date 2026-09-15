package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// Parser turns a token stream from the lexer into an *ast.Program via Pratt parsing.
type Parser struct {
	lexer     *lexer.Lexer
	errors    []string
	nuds      map[token.Type]nudFuncType
	leds      map[token.Type]ledFuncType
	curToken  token.Token
	peekToken token.Token
}

// New creates a Parser for l, priming curToken/peekToken and registering all nuds and leds.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
		nuds:   make(map[token.Type]nudFuncType),
		leds:   make(map[token.Type]ledFuncType),
	}
	p.advanceToken()
	p.advanceToken()
	p.registerNuds()
	p.registerLeds()
	return p
}

// Errors returns the parser errors accumulated while parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram parses the whole token stream into an *ast.Program of top-level statements.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.advanceToken()
	}
	return program
}

// parseStatement dispatches on curToken's type to parse one top-level or block statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		if statement := p.parseLetStatement(); statement != nil {
			return statement
		}
	case token.RETURN:
		if statement := p.parseReturnStatement(); statement != nil {
			return statement
		}
	default:
		if statement := p.parseExpressionStatement(); statement != nil {
			return statement
		}
	}
	return nil
}

// parseLetStatement parses a "let <identifier expression> = <expression>;" let statement.
//
// It expects tokens on entry: let x = 5;  (curToken must be "let").
//
// It leaves curToken on the trailing ";" if present, otherwise on the value expression's last token (e.g. "5").
func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{Token: p.curToken}
	if !p.expectAdvancePeek(token.IDENTIFIER) {
		return nil
	}
	letStatement.IdentifierExpression = &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectAdvancePeek(token.ASSIGN) {
		return nil
	}
	p.advanceToken()
	letStatement.ValueExpression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.advanceToken()
	}
	return letStatement
}

// parseReturnStatement parses a "return <expression>;" return statement.
//
// It expects tokens on entry: return 5;  (curToken must be "return").
//
// It leaves curToken on the trailing ";" if present, otherwise on the value expression's last token (e.g. "5").
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{Token: p.curToken}
	p.advanceToken()
	returnStatement.ValueExpression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.advanceToken()
	}
	return returnStatement
}

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
		message := fmt.Sprintf("no nud parse function for %v found", p.curToken.Type)
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
