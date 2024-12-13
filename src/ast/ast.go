package ast

import (
	"s8/src/token"
)

type Node interface {
	TokenLiteral() string // Only for debugging and testing
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

type Identifier struct {
	Token token.Token // the token.IDENT
	Value string      // the identifier as a literal value
}

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

// This method implementation makes Identifier satisfy the Expression interface
func (i *Identifier) expressionNode() {}

// As Identifier aims to satisfy the Expression interface, it also needs to satisfy the Node interface
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
