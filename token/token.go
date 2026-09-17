package token

// Type identifies the kind of a Token.
type Type string

// Token is a single lexical token produced by the lexer: its kind and the source text it came from.
type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Indentifiers
	IDENTIFIER = "IDENTIFIER" // add, foobar, x, y, ...

	// Literals
	INTEGER = "INTEGER" // 123456
	STRING  = "STRING"  // 123456

	// Operators
	ASSIGN   = "ASSIGN"
	PLUS     = "PLUS"
	MINUS    = "MINUS"
	BANG     = "BANG"
	ASTERISK = "ASTERISK"
	SLASH    = "SLASH"

	// Comparison operators
	LESS_THAN    = "LESS_THAN"
	GREATER_THAN = "GREATER_THAN"
	EQUAL        = "EQUAL"
	NOT_EQUAL    = "NOT_EQUAL"

	// Separators
	COMMA     = "COMMA"
	COLON     = "COLON"
	SEMICOLON = "SEMICOLON"

	// Brackets
	LEFT_PAREN    = "LEFT_PAREN"
	RIGHT_PAREN   = "RIGHT_PAREN"
	LEFT_BRACE    = "LEFT_BRACE"
	RIGHT_BRACE   = "RIGHT_BRACE"
	LEFT_BRACKET  = "LEFT_BRACKET"
	RIGHT_BRACKET = "RIGHT_BRACKET"

	// Keywords
	LET      = "LET"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

// keywords maps each reserved word to its TokenType.
var keywords = map[string]Type{
	"let":    LET,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
}

// LookupKeyword returns the keyword TokenType for word, or IDENTIFIER if word is not a reserved word.
func LookupKeyword(word string) Type {
	tokenType, ok := keywords[word]
	if !ok {
		return IDENTIFIER
	}
	return tokenType
}
