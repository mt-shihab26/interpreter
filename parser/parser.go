package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// Parser turns a token stream from the lexer into an *ast.Program via Pratt parsing.
type Parser struct {
	l         *lexer.Lexer
	errors    []string
	nuds      map[token.TokenType]nudFuncType
	leds      map[token.TokenType]ledFuncType
	curToken  token.Token
	peekToken token.Token
}

// New creates a Parser for l, priming curToken/peekToken and registering all nuds and leds.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
		nuds:   make(map[token.TokenType]nudFuncType),
		leds:   make(map[token.TokenType]ledFuncType),
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
