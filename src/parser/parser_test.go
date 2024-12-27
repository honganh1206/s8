package parser

import (
	"fmt"
	"s8/src/ast"
	"s8/src/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	// no mocks here for readability
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nill")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got: %d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

// Check as many fields of the AST nodes as possible
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() is not 'let'. got: %q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement) // Type assertion to convert ast.Statement into a pointer

	if !ok {
		t.Errorf("s not *ast.LetStatement. got: %T", s)
		return false
	}

	if letStmt.Name.Value != name {

		t.Errorf("letStmt.Name.Value not '%s'. got: %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got: %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got: %d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got: %T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.Name.TokenLiteral() not 'return'. got: %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram() // construct the root node
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got: %d", len(program.Statements))
	}

	// Type assertion to convert from the Statement interface to concrete type ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier. got: %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got: %s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s. got: %s", "foobar", ident.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram() // construct the root node
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got: %d", len(program.Statements))
	}

	// Type assertion to convert from the Statement interface to concrete type ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("exp not ast.IntegerLiteral. got: %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got: %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() not %s. got: %s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	// Table-driven test approach
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram() // construct the root node
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got: %d", 1, len(program.Statements))
		}

		// Type assertion to convert from the Statement interface to concrete type ExpressionStatement
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got: %T", stmt.Expression)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got: '%s'", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) bool {
	integ, ok := expr.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("expr not *ast.IntegerLiteral. got: %T", expr)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got: %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() not %d. got: %s", value, integ.TokenLiteral())
		return false
	}

	return true

}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram() // construct the root node
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got: %d", 1, len(program.Statements))
		}

		// Type assertion to convert from the Statement interface to concrete type ExpressionStatement
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("stmt is not ast.InfixExpression. got: %T", stmt.Expression)
		}

		if !testIntegerLiteral(t, expr.Right, tt.leftValue) {
			return
		}
		if expr.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got: '%s'", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.rightValue) {
			return
		}
	}

}
