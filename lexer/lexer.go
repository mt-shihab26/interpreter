package lexer

import "monkey/token"

// Lexer turns Monkey source code into a stream of tokens, one NextToken() call at a time.
type Lexer struct {
	input        string
	curPosition  int  // current position in input (points to current char)
	peekPosition int  // current reading position in input (after current char)
	chracter     byte // current char under examination
}

// New creates a Lexer over input, priming chracter with the first character.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.advanceCharacter()
	return l
}

// NextToken skips whitespace and returns the next token, peeking ahead to recognize two-character operators.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.chracter {
	case '=':
		if l.peekCharacter() == '=' {
			ch := l.chracter
			l.advanceCharacter()
			tok = token.Token{Type: token.EQUAL, Literal: string(ch) + string(l.chracter)}
		} else {
			tok = newChToken(token.ASSIGN, l.chracter)
		}
	case '+':
		tok = newChToken(token.PLUS, l.chracter)
	case '-':
		tok = newChToken(token.MINUS, l.chracter)
	case '!':
		if l.peekCharacter() == '=' {
			ch := l.chracter
			l.advanceCharacter()
			tok = token.Token{Type: token.NOT_EQUAL, Literal: string(ch) + string(l.chracter)}
		} else {
			tok = newChToken(token.BANG, l.chracter)
		}
	case '*':
		tok = newChToken(token.ASTERISK, l.chracter)
	case '/':
		tok = newChToken(token.SLASH, l.chracter)
	case '<':
		tok = newChToken(token.LESS_THAN, l.chracter)
	case '>':
		tok = newChToken(token.GREATER_THAN, l.chracter)
	case ',':
		tok = newChToken(token.COMMA, l.chracter)
	case ';':
		tok = newChToken(token.SEMICOLON, l.chracter)
	case '(':
		tok = newChToken(token.LEFT_PAREN, l.chracter)
	case ')':
		tok = newChToken(token.RIGHT_PAREN, l.chracter)
	case '{':
		tok = newChToken(token.LEFT_BRACE, l.chracter)
	case '}':
		tok = newChToken(token.RIGHT_BRACE, l.chracter)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.consumeString()
	default:
		if isLetter(l.chracter) {
			tok.Literal = l.consumeIdentifier()
			tok.Type = token.LookupKeyword(tok.Literal)
			return tok
		} else if isDigit(l.chracter) {
			tok.Type = token.INTEGER
			tok.Literal = l.consumeNumber()
			return tok
		} else {
			tok = newChToken(token.ILLEGAL, l.chracter)
		}
	}
	l.advanceCharacter()
	return tok
}

// skipWhitespace advances past spaces, tabs, newlines, and carriage returns.
func (l *Lexer) skipWhitespace() {
	for l.chracter == ' ' || l.chracter == '\t' || l.chracter == '\n' || l.chracter == '\r' {
		l.advanceCharacter()
	}
}

// advanceCharacter advances the lexer by one character, setting chracter to 0 (NUL) once input is exhausted.
func (l *Lexer) advanceCharacter() {
	if l.peekPosition >= len(l.input) {
		l.chracter = 0
	} else {
		l.chracter = l.input[l.peekPosition]
	}
	l.curPosition = l.peekPosition
	l.peekPosition += 1
}

// peekCharacter returns the character at peekPosition without advancing, or 0 (NUL) at end of input.
func (l *Lexer) peekCharacter() byte {
	if l.peekPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.peekPosition]
	}
}

// consumeIdentifier consumes and returns consecutive letters starting at curPosition.
func (l *Lexer) consumeIdentifier() string {
	position := l.curPosition
	for isLetter(l.chracter) {
		l.advanceCharacter()
	}
	return l.input[position:l.curPosition]
}

// consumeNumber consumes and returns consecutive digits starting at curPosition.
func (l *Lexer) consumeNumber() string {
	position := l.curPosition
	for isDigit(l.chracter) {
		l.advanceCharacter()
	}
	return l.input[position:l.curPosition]
}

// readNumber consumes and returns consecutive digits starting at curPosition.
func (l *Lexer) consumeString() string {
	str := ""
	l.advanceCharacter()
	for l.chracter != '"' {
		str = str + string(l.chracter)
		l.advanceCharacter()
	}
	return str
}

// isLetter reports whether ch is an ASCII letter or underscore, i.e. valid in an identifier.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit reports whether ch is an ASCII digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// newChToken builds a single-character token.Token of the given type from ch.
func newChToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
