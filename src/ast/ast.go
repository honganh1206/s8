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
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token // the identifier token
	Value string      // the identifier name as a literal value
}

type LetStatement struct {
	Token token.Token // the 'let' token
	Name  *Identifier // the identifier holding the variable name
	Value Expression  // the expression producing the value
}
