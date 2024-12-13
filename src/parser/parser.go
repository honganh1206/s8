package parser

import (
	"s8/src/ast"
	"s8/src/lexer"
	"s8/src/token"
)

type Parser struct {
	l            *lexer.Lexer // One and only instance of the lexer
	currentToken token.Token  // work as the position field
	peekToken    token.Token  // work as the readPosition field
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read TWO tokens so currentToken and peekToken are bot set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
