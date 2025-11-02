package parser

import (
	"fmt"
	"strconv"

	"s8/ast"
	"s8/lexer"
	"s8/token"
)

// Operator precedences
const (
	_ int = iota // Give the following constants incrementing numbers as value
	LOWEST
	ASSIGN
	EQUALS      // ==
	CONDITIONAL // ? and :
	LESSGREATER // > or <
	BITWISE     // &, |, ^, ~, <<, >>
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	POSTFIX     // x++ or x--
	CALL        // myFunction(X)
	INDEX       // arr[index]
)

// Precedence table - Map token types with their precedence
var precedences = map[token.TokenType]int{
	token.EQ:        EQUALS,
	token.NOT_EQ:    EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTERISK:  PRODUCT,
	token.LPAREN:    CALL,
	token.QUESTION:  CONDITIONAL,
	token.LBRACKET:  INDEX,
	token.TILDE:     BITWISE,
	token.EXPONENT:  BITWISE,
	token.PIPE:      BITWISE,
	token.AMPERSAND: BITWISE,
	token.RSHIFT:    BITWISE,
	token.LSHIFT:    BITWISE,
	token.ASSIGN:    ASSIGN,
}

type (
	prefixParseFn func() ast.Expression
	// Infix parsing will ALWAYS has an expression before the operator
	// The Expression-typed argument denotes the "left side" of the infix operator that is being parsed
	// We call it infix but actually it is just ANYTHING BUT PREFIX
	infixParseFn   func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l            *lexer.Lexer // One and only instance of the lexer
	currentToken token.Token  // work as the position field
	peekToken    token.Token  // work as the readPosition field
	prevToken    token.Token
	errors       []string

	// Check if either map has a parsing function associated with currentToken.Type
	// Having separated tables for prefix and infix expressions is important
	// As sometimes we use the same token for different expressions e.g., "(" for grouped expression (prefix) and for call expression (infix)
	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	// Read TWO tokens so currentToken and peekToken are bot set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TILDE, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.TILDE, p.parsePrefixExpression)
	p.registerPrefix(token.MACRO, p.parseMacroLiteral)

	// The actual determination of whether it's prefix or postfix
	// should be done during parsing, not during registration
	precedences[token.INCREMENT] = PREFIX
	precedences[token.DECREMENT] = PREFIX
	p.registerPrefix(token.INCREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.DECREMENT, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.QUESTION, p.parseTernaryExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.EXPONENT, p.parseInfixExpression)
	p.registerInfix(token.PIPE, p.parseInfixExpression)
	p.registerInfix(token.RSHIFT, p.parseInfixExpression)
	p.registerInfix(token.LSHIFT, p.parseInfixExpression)
	p.registerInfix(token.AMPERSAND, p.parseInfixExpression)
	p.registerInfix(token.EXPONENT, p.parseInfixExpression)
	// Assign binds two expressions e.g., a = b + c
	// so it makes sense we make it an infix
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	precedences[token.INCREMENT] = POSTFIX
	precedences[token.DECREMENT] = POSTFIX
	p.registerPostfix(token.INCREMENT, p.parsePostfixExpression)
	p.registerPostfix(token.DECREMENT, p.parsePostfixExpression)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

/*
	HELPER METHODS
*/

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Advance both of our p.currentToken and p.peekToken
func (p *Parser) nextToken() {
	p.prevToken = p.currentToken
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Enforce the correctness of the order of tokens by checking the type of the next token
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

func (p *Parser) previousTokenIs(t token.TokenType) bool {
	return p.prevToken.Type == t
}

// Strengthen left-binding power
// This method makes our operators more right-associative
// By decrementing their precedence
func (p *Parser) peekPrecedence() int {
	// Current token is an operand
	// Check if the left-binding power of the next operator is higher than the current right-binding power
	if pre, ok := precedences[p.peekToken.Type]; ok {
		return pre
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if pre, ok := precedences[p.currentToken.Type]; ok {
		return pre
	}

	return LOWEST
}

/*
	PARSING STATEMENTS
*/

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Construct an *ast.LetStatement node
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	// Bind the variable name to the function
	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.currentToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{Token: p.currentToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

/*
	PARSING EXPRESSIONS
*/

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	// We pass the lowest precedence since we have yet to parse anything and cannot compare precedences
	stmt.Expression = p.parseExpression(LOWEST)

	// This check is optional
	// As our expression statements should have optional semicolons like 5 + 5 (like JS)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Call the corresponding parsing function for the current token type
// Returning an interface type value, so we do not use pointer here
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	// Prefix operators have higher precedence than infix operator, but lower than postfix operator
	// So prefix expressions stay as a separate unit
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	// The default precedence is LOWEST
	// Then we might increase the precedence with each call to parseExpression
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// Current position is the left expression
		if !p.peekTokenIs(token.INCREMENT) && !p.peekTokenIs(token.DECREMENT) {
			infix := p.infixParseFns[p.peekToken.Type]

			if infix == nil {
				return leftExp
			}
			p.nextToken() // Advance to the infix operator

			// The peek token is now the left expression
			// While the right expression is passed in the infix parsing function
			leftExp = infix(leftExp)
			continue // Continue to handle other cases like using both pre/postfix with infix
		} else {
			if !p.currentTokenIs(token.IDENT) && !p.currentTokenIs(token.INT) && !p.currentTokenIs(token.FLOAT) {
				return leftExp
			}

			p.nextToken() // Move to the postfix operator

			postfix := p.postfixParseFns[p.currentToken.Type]
			if postfix == nil {
				p.noPostfixParseFnError(p.currentToken.Type)
				return nil
			}
			leftExp = postfix(leftExp)
		}
	}

	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPostfixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no postfix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	literal := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	literal := &ast.FloatLiteral{Token: p.currentToken}

	// Automatically round to 6 decimal places
	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

// Parse prefix expressions (e.g., "-5") by first constructing the AST node for the prefix operator,
// then parsing the right-side expression. Results in a PrefixExpression node containing the operator
// and its right-side expression as children.
func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))
	expr := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expr.Right = p.parseExpression(PREFIX)

	return expr
}

// When this method is invoked, the current token is the operator
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixStatement"))

	expr := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken() // Current token is now the right expression
	// Only to demonstrate how we increase the right-associative-ness of our operators
	// if expr.Operator == "+" {
	// 	expr.Right = p.parseExpression(precedence - 1)
	// } else {
	expr.Right = p.parseExpression(precedence)
	// }

	return expr
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	expr := &ast.Assignment{
		Token: p.currentToken,
		Name:  left,
	}

	p.nextToken()
	// When we are done parsing the identifier, the current precedence is the lowest,
	// but we still need to pass in LOWEST as an argument
	// so we can recursively parse the right-side value
	expr.Value = p.parseExpression(LOWEST)

	return expr
}

// At this point we already parsed the left expression
// We thus need to consider parsing the precedence
func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	// Current token is the postfix operator
	return &ast.PostfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
}

// Maintain proper precedence by treating the grouped expression as a single unit
func (p *Parser) parseGroupedExpression() ast.Expression {
	// When this function is called, currentToken is at '('
	p.nextToken()

	// Parse everything until ')'
	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Start parsing block statement
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	// Optional ELSE
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expr.Alternative = p.parseBlockStatement()

	}

	return expr
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	expr := &ast.WhileStatement{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	expr := &ast.ForStatement{Token: p.currentToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	// TODO: Init could be a let statement or an assign expression
	expr.Init = p.parseLetStatement()
	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	p.nextToken()

	expr.Update = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	// From here we start parsing multiple statements
	p.nextToken()

	// Quite similar to parseProgram() isnt it?
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	idents := []*ast.Identifier{}

	// Case no param specified
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return idents
	}

	p.nextToken() // At this point our current token is the 1st param

	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	idents = append(idents, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // to the comma
		p.nextToken() // to the next func param

		// Parsing the remaining params
		ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		idents = append(idents, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return idents
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	ce := &ast.CallExpression{Token: p.currentToken, Function: fn}

	ce.Arguments = p.parseExpressionList(token.RPAREN)

	return ce
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.currentToken}

	arr.Elements = p.parseExpressionList(token.RBRACKET)

	return arr
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()

	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // Advance to the comma
		p.nextToken() // Advance to the next argument

		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// A kind of mixfix expression, even though we treat this as infix
func (p *Parser) parseTernaryExpression(condition ast.Expression) ast.Expression {
	expr := &ast.TernaryExpression{
		Token:     p.currentToken,
		Condition: condition,
	}

	pre := p.currentPrecedence()

	p.nextToken() // Move past the "?"
	expr.Consequence = p.parseExpression(pre)

	if !p.expectPeek(token.COLON) {
		msg := fmt.Sprintf("expected next token to be %s", p.peekToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	p.nextToken() // Move past the ":"
	expr.Alternative = p.parseExpression(pre)

	return expr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.nextToken() // To the index value

	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return expr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.currentToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()                    // Move past the "{"
		key := p.parseExpression(LOWEST) // Parsing the key expression

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken() // Move past the ":"
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value
		// Invalid syntax case
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}

	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) parseMacroLiteral() ast.Expression {
	lit := &ast.MacroLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}
