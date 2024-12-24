package ast

import (
	"bytes"
	"s8/src/token"
)

type Node interface {
	TokenLiteral() string // Only for debugging and testing
	String() string       // Print AST nodes for debugging
}

// Nodes implementing the Node interface have to provide the TokenLiteral() method
// that returns the literal value
type Statement interface {
	Node
	statementNode() // Dummy method to guide the Go compiler if we mess up
}

type Expression interface {
	Node
	expressionNode() // Dummy method to guide the Go compiler if we mess up
}

// Root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) <= 0 {
		return ""
	} else {
		return p.Statements[0].TokenLiteral()
	}
}

// Write the return value of each statement's String() method
// Then return the buffer as a string
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT
	Value string      // the identifier as a literal value
}

// This method implementation makes Identifier satisfy the Expression interface
func (i *Identifier) expressionNode() {}

// As Identifier aims to satisfy the Expression interface, it also needs to satisfy the Node interface
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }

///               +-----------------------+
///               |                       |
///               |     *ast.Program      |
///               +-----------------------+
///               |                       |
///               |     Statements        |
///               +----------+------------+
///                          |
///                          |
///                          |
///                          |
///               +----------v------------+
///               |  *ast.LetStatement    |
///               +-----------------------+
///               |                       |
///         +-----+       Name            +-----+
///         |     +-----------------------+     |
///         |     |                       |     |
///         |     |       Value           |     |
///         |     +-----------------------+     |
///         |                                   |
/// +-------v---------+              +----------v------+
/// |                 |              |                 |
/// |                 |              |                 |
/// | *ast.Identifier |              | *ast.Expression |
/// |                 |              |                 |
/// |                 |              |                 |
/// +-----------------+              +-----------------+
///

type LetStatement struct {
	Token token.Token // the token.LET
	Name  *Identifier // the identifier holding the variable name
	Value Expression  // the expression producing the value
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // 1st token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {} // fulfill the ast.Statement interface so we can add this to the Statements[] slice

func (es ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64 // When we build this AST node we have to convert the string to int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *IntegerLiteral) String() string { return il.Token.Literal }
