package lexer

import "monkey/token"

type Lexer struct {
	input        string
	curPosition  int  // current position in input (points to current char)
	peekPosition int  // current reading position in input (after current char)
	chracter     byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.peekPosition >= len(l.input) {
		l.chracter = 0
	} else {
		l.chracter = l.input[l.peekPosition]
	}
	l.curPosition = l.peekPosition
	l.peekPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.chracter {
	case '=':
		if l.peekChar() == '=' {
			ch := l.chracter
			l.readChar()
			tok = token.Token{Type: token.EQUALUAL, Literal: string(ch) + string(l.chracter)}
		} else {
			tok = newToken(token.ASSIGN, l.chracter)
		}
	case '+':
		tok = newToken(token.PLUS, l.chracter)
	case '-':
		tok = newToken(token.MINUS, l.chracter)
	case '!':
		if l.peekChar() == '=' {
			ch := l.chracter
			l.readChar()
			tok = token.Token{Type: token.NOT_EQUAL, Literal: string(ch) + string(l.chracter)}
		} else {
			tok = newToken(token.BANG, l.chracter)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.chracter)
	case '/':
		tok = newToken(token.SLASH, l.chracter)
	case '<':
		tok = newToken(token.LESS_THAN, l.chracter)
	case '>':
		tok = newToken(token.GREATER_THAN, l.chracter)
	case ',':
		tok = newToken(token.COMMA, l.chracter)
	case ';':
		tok = newToken(token.SEMICOLON, l.chracter)
	case '(':
		tok = newToken(token.LEFT_PAREN, l.chracter)
	case ')':
		tok = newToken(token.RIGHT_PAREN, l.chracter)
	case '{':
		tok = newToken(token.LEFT_BRACE, l.chracter)
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.chracter)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.chracter) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupKeyword(tok.Literal)
			return tok
		} else if isDigit(l.chracter) {
			tok.Type = token.INTEGER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.chracter)
		}
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.curPosition
	for isLetter(l.chracter) {
		l.readChar()
	}
	return l.input[position:l.curPosition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhitespace() {
	for l.chracter == ' ' || l.chracter == '\t' || l.chracter == '\n' || l.chracter == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.curPosition
	for isDigit(l.chracter) {
		l.readChar()
	}
	return l.input[position:l.curPosition]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.peekPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.peekPosition]
	}
}
