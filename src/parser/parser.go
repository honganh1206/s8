package parser

import (
	"fmt"
	"s8/src/ast"
	"s8/src/lexer"
	"s8/src/token"
)

type Parser struct {
	l            *lexer.Lexer // One and only instance of the lexer
	currentToken token.Token  // work as the position field
	peekToken    token.Token  // work as the readPosition field
	errors       []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Read TWO tokens so currentToken and peekToken are bot set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Advance both of our p.currentToken and p.peekToken
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Enfore the correctness of the order of tokens by checking the type of the next token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Construct an *ast.LetStatement node
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Construct an *ast.Identifier node
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Skip expr until a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// TODO: Skip expr until a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // Construct the root node
	program.Statements = make([]ast.Statement, 0)

	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}
