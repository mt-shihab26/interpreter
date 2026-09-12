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
	SEMICOLON = "SEMICOLON"

	// Brackets
	LEFT_PAREN  = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	LEFT_BRACE  = "LEFT_BRACE"
	RIGHT_BRACE = "RIGHT_BRACE"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// keywords maps each reserved word to its TokenType.
var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupKeyword returns ident's keyword TokenType, or IDENTIFIER if ident is not a reserved word.
func LookupKeyword(ident string) Type {
	tok, ok := keywords[ident]
	if !ok {
		return IDENTIFIER
	}
	return tok
}
