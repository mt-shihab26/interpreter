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
	l.readChar()
	return l
}

// readChar advances the lexer by one character: it moves peekPosition's
// character into chracter, then shifts curPosition/peekPosition forward.
// It sets chracter to 0 (NUL) once the input is exhausted.
func (l *Lexer) readChar() {
	if l.peekPosition >= len(l.input) {
		l.chracter = 0
	} else {
		l.chracter = l.input[l.peekPosition]
	}
	l.curPosition = l.peekPosition
	l.peekPosition += 1
}

// NextToken consumes and returns the next token from the input, skipping
// leading whitespace first. Two-character operators (e.g. "==", "!=") are
// recognized by peeking one character ahead before falling back to the
// single-character token. It leaves chracter on the character right after
// the returned token.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.chracter {
	case '=':
		if l.peekChar() == '=' {
			ch := l.chracter
			l.readChar()
			tok = token.Token{Type: token.EQUAL, Literal: string(ch) + string(l.chracter)}
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

// newToken builds a single-character token.Token of the given type from ch.
func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// readIdentifier consumes consecutive letters starting at curPosition and
// returns them as a string, leaving chracter on the first non-letter after it.
func (l *Lexer) readIdentifier() string {
	position := l.curPosition
	for isLetter(l.chracter) {
		l.readChar()
	}
	return l.input[position:l.curPosition]
}

// isLetter reports whether ch is an ASCII letter or underscore, i.e. valid in an identifier.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// skipWhitespace advances past spaces, tabs, newlines, and carriage returns.
func (l *Lexer) skipWhitespace() {
	for l.chracter == ' ' || l.chracter == '\t' || l.chracter == '\n' || l.chracter == '\r' {
		l.readChar()
	}
}

// readNumber consumes consecutive digits starting at curPosition and
// returns them as a string, leaving chracter on the first non-digit after it.
func (l *Lexer) readNumber() string {
	position := l.curPosition
	for isDigit(l.chracter) {
		l.readChar()
	}
	return l.input[position:l.curPosition]
}

// isDigit reports whether ch is an ASCII digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// peekChar returns the character at peekPosition without advancing the
// lexer, or 0 (NUL) if that position is past the end of the input.
func (l *Lexer) peekChar() byte {
	if l.peekPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.peekPosition]
	}
}
