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
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let foobar = 838383;", "foobar", 838383},
	}
	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nill")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got: %d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
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

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return 838383;", 838383},
	}
	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got: %d",
				len(program.Statements))
		}

		stmt := program.Statements[0]

		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got: %T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.Name.TokenLiteral() not 'return'. got: %q", returnStmt.TokenLiteral())
		}

		val := stmt.(*ast.ReturnStatement).ReturnValue

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
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

	if !testLiteralExpression(t, stmt.Expression, "foobar") {
		return
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
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
		{"~5;", "~", -6}, // TODO: Add parsing for tilde
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

		if !testLiteralExpression(t, expr.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.rightValue, tt.operator) {
			return
		}
	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},

		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},

		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a ? b : c",
			"(a ? b : c)",
		},
		{
			"a + b ? c + d : e + f",
			"((a + b) ? (c + d) : (e + f))",
		},
		{
			"5 > 3 ? 1 : 2",
			"((5 > 3) ? 1 : 2)",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram() // construct the root node
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram() // construct the root node
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	bool, ok := stmt.Expression.(*ast.Boolean)

	if !ok {
		t.Fatalf("exp not ast.Boolean. got: %T", stmt.Expression)
	}

	if bool.Value != true {
		t.Errorf("literal.Value not %t. got: %t", true, bool.Value)
	}

	if bool.TokenLiteral() != "true" {
		t.Errorf("literal.TokenLiteral() not %s. got: %s", "true", bool.TokenLiteral())
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	// no mocks here for readability
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("expr not ast.IfExpression. got: %T", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "y", "<") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got: %d\n", len(expr.Consequence.Statements))
	}

	if expr.Alternative != nil {
		t.Errorf("expr.Alternative.Statements was not nil. got %+v", expr.Alternative)
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] not ast.ExpressionStatement. got: %T", expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression not ast.IfExpression. got: %T", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "y", "<") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got: %d\n", len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] not ast.ExpressionStatement. got: %T", expr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `funk(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got: %T", stmt.Expression)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got: %d\n", stmt.Expression)
	}

	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	if len(fn.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statement. got: %d\n", len(fn.Body.Statements))
	}

	bodyStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function.Body.Statements[0] is not ast.ExpressionStatement. got: %T", fn.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "y", "+")

}

func TestFunctionParametersParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "funk() {};", expectedParams: []string{}},
		{input: "funk(x) {};", expectedParams: []string{"x"}},
		{input: "funk(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fn := stmt.Expression.(*ast.FunctionLiteral)

		if len(fn.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want: %d, got: %d\n", len(tt.expectedParams), len(fn.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, fn.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has not enough statements. got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got: %T", stmt.Expression)
	}

	if !testIdentifier(t, expr.Function, "add") {
		return
	}

	if len(expr.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got: %d", len(expr.Arguments))
	}

	// FIXME: comment this out makes the test pass
	testLiteralExpression(t, expr.Arguments[0], 1)
	testInfixExpression(t, expr.Arguments[1], 2, 3, "*")
	testInfixExpression(t, expr.Arguments[2], 4, 5, "+")

}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got: %T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want: %d, got: %d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want: %q, got: %q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func TestStringLiteralParsing(t *testing.T) {
	input := `"hello world"`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not *ast.StringLiteral. got: %T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got: %q", "hello world", literal.Value)
	}

}

func TestArrayLiteralParsing(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.ArrayLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not *ast.ArrayLiteral. got: %T", stmt.Expression)
	}

	if len(literal.Elements) != 3 {
		t.Fatalf("len(literal.Elements) not 3. got: %d", len(literal.Elements))
	}

	testIntegerLiteral(t, literal.Elements[0], 1)
	testInfixExpression(t, literal.Elements[1], 2, 2, "*")
	testInfixExpression(t, literal.Elements[2], 3, 3, "+")

}

func TestTernaryExpressionParsing(t *testing.T) {
	input := "x > 5 ? 10 : 5"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got: %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got: %T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.TernaryExpression)
	if !ok {
		t.Fatalf("expr not ast.TernaryExpression. got: %T", stmt.Expression)
	}

	testInfixExpression(t, expr.Condition, "x", 5, ">")
	testLiteralExpression(t, expr.TrueValue, 10)
	testLiteralExpression(t, expr.FalseValue, 5)
}

/*
	TEST HELPERS
*/

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}

	t.Errorf("type of expr not handled. got: %T", expr)

	return false
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

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)

	if !ok {
		t.Errorf("expr not *ast.Identifier. got: %T", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got: %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got: %s", value, ident.TokenLiteral())
		return false
	}

	return true

}

func testInfixExpression(t *testing.T, expr ast.Expression, left, right interface{}, operator string) bool {
	opExpr, ok := expr.(*ast.InfixExpression)

	if !ok {
		t.Errorf("expr not ast.InfixExpression. got: %T(%s)", expr, expr)
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != operator {
		t.Errorf("expr.Operator not %s. got: %q", operator, opExpr.Operator)
	}

	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	bo, ok := expr.(*ast.Boolean)

	if !ok {
		t.Errorf("expr not *ast.Boolean. got: %T", expr)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got: %t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("ident.TokenLiteral() not %t. got: %s", value, bo.TokenLiteral())
		return false
	}

	return true
}
