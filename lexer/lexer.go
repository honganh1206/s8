package lexer

import (
	"s8/token"
	"strings"
)

// Include pointers to peek further into the input
type Lexer struct {
	input        string
	position     int  // current char
	readPosition int  // after current char
	ch           rune // current char under examination
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
		l.ch = rune(l.input[l.readPosition])
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// Handle cases like != and ==
func (l *Lexer) makeTwoCharToken(tokenType token.TokenType) token.Token {
	ch := l.ch
	l.readChar()
	literal := string(ch) + string(l.ch)
	tok := token.Token{Type: tokenType, Literal: literal}
	return tok
}

// Get the next char but NOT advance our position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		// Append the 2nd assign token to the 1st one to form the equal token
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.EQ)
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		if l.peekChar() == '+' {
			tok = l.makeTwoCharToken(token.INCREMENT)
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '-':
		if l.peekChar() == '-' {
			tok = l.makeTwoCharToken(token.DECREMENT)
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.NOT_EQ)
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		if l.peekChar() == '<' {
			tok = l.makeTwoCharToken(token.LSHIFT)
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '>' {
			tok = l.makeTwoCharToken(token.RSHIFT)
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '~':
		tok = newToken(token.TILDE, l.ch)
	case '|':
		tok = newToken(token.PIPE, l.ch)
	case '&':
		tok = newToken(token.AMPERSAND, l.ch)
	case '?':
		tok = newToken(token.QUESTION, l.ch)
	case '^':
		tok = newToken(token.EXPONENT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isLetter(l.ch) {
			// Pretty interesting: This is in reverse compared to when we deal with digits
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			literal := l.readNumber()
			if strings.Contains(literal, ".") {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Literal = literal
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
func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Changing this function will have a large impact on how our interpreter will parse
// Like the check for '_' means we can use snake case e.g., foo_bar as identifier
func isLetter(ch rune) bool {
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

// At this point we have yet to support floats or hex notations and things alike
func (l *Lexer) readNumber() string {
	position := l.position
	seenDot := false

	// Keep reading characters as long as they are numerical values
	for isDigit(l.ch) {
		if l.ch == '.' {
			if seenDot {
				break // Do NOT allow multiple decimal points
			}
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func (l *Lexer) readString() string {
	// Move to the chars inside the ""
	pos := l.position + 1
	for {
		l.readChar()
		// Reaching the end of the string
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[pos:l.position]
}
