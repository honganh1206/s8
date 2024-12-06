package lexer

import "s8/src/token"

// Include pointers to peek further into the input
type Lexer struct {
	input        string
	position     int  // current char
	readPosition int  // after current char
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Get the next char and advance our position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	// Kinda neat: We advance our pointer here
	// So the next time we call `NextToken()` the l.ch field is already updated
	l.readChar()
	return tok
}

// Examine the current character and return a token depending on which character it is
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Changing this function will have a large impact on how our interpreter will parse
// Like the check for '_' means we can use snake case e.g., foo_bar as identifier
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '!' || ch == '?'
}

// Return the literal value of the identifer
func (l *Lexer) readIdentifier() string {
	position := l.position

	// Keep reading characters as long as they are valid identifier characters
	for isLetter(l.ch) {
		l.readChar()
	}

	// Return a slice of the input/the complete identifer as a string
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position

	// Keep reading characters as long as they are numerical values
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
