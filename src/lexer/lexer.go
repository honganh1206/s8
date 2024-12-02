package lexer

type Lexer struct {
	input        string
	position     int  // current char
	readPosition int  // after current char
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	return l
}
