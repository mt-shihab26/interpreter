package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +, -
	PRODUCT     // *, /
	PREFIX      // -x or !x
	CALL        // myFunction(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %v, got %v instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	p.registerPrefix(token.MINUS, p.parseUnaryExpression)
	p.registerPrefix(token.BANG, p.parseUnaryExpression)
	p.registerPrefix(token.IDENT, p.parseIdentifierExpression)
	p.registerPrefix(token.INT, p.parseIntegerLiteralExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanExpression)
	p.registerPrefix(token.FALSE, p.parseBooleanExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	p.registerInfix(token.PLUS, p.parseBinaryExpression)
	p.registerInfix(token.MINUS, p.parseBinaryExpression)
	p.registerInfix(token.ASTERISK, p.parseBinaryExpression)
	p.registerInfix(token.SLASH, p.parseBinaryExpression)
	p.registerInfix(token.EQ, p.parseBinaryExpression)
	p.registerInfix(token.NOT_EQ, p.parseBinaryExpression)
	p.registerInfix(token.GT, p.parseBinaryExpression)
	p.registerInfix(token.LT, p.parseBinaryExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
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

func (p *Parser) parseCallArguments() []ast.Expression {
	arguments := []ast.Expression{}
	p.nextToken()
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		arguments = append(arguments, p.parseExpression(LOWEST))
		p.nextToken()
		if p.curTokenIs(token.RPAREN) {
			break
		}
		if !p.curTokenIs(token.COMMA) {
			return nil
		}
		p.nextToken()
	}
	return arguments
}
