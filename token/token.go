package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string // hold the literal value
}

var keywords = map[string]TokenType{
	"funk":     FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"macro":    MACRO,
	"for":      FOR,
	"while":    WHILE,
	"break":    BREAK,
	"continue": CONTINUE,
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"
	FLOAT = "FLOAT"

	// Operators
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	LT        = "<"
	GT        = ">"
	EQ        = "=="
	NOT_EQ    = "!="
	INCREMENT = "++"
	DECREMENT = "--"
	TILDE     = "~"
	QUESTION  = "?"
	EXPONENT  = "^"
	PIPE      = "|"
	RSHIFT    = ">>" // divided by 2 e.g., n >> x means n divided by 2, x times
	LSHIFT    = "<<" // times 2 e.g., n << x means n times 2, x times
	AMPERSAND = "&"
	// TODO: Implement GTE (>=) and LTE (<=)

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	WHILE    = "WHILE"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"

	// Data types
	STRING = "STRING"

	// Arrays
	LBRACKET = "["
	RBRACKET = "]"

	MACRO = "MACRO"
)

// Check if the given identifier is a keyword or a user-defined identifier
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
