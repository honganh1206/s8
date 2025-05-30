package ast

import (
	"bytes"
	"s8/token"
	"strings"
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

type FloatLiteral struct {
	token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FloatLiteral) String() string { return fl.Token.Literal }

// +-----------------------+
// |                      |
// |    *ast.Program      |
// +-----------------------+
// |                      |
// |     Statements       |
// +----------+-----------+
//            |
//            |
// +----------v-----------+
// |                      |
// |*ast.ExpressionStmt   |
// +-----------+----------+
//            |
//            |
// +----------v-----------+
// |                      |
// |*ast.PrefixExpression |
// +-----------------------+
// | Token: MINUS         |
// | Operator: "-"        |
// | Right               -+------+
// +-----------------------+     |
//                               |
//                    +----------v----------+
//                    |                     |
//                    | *ast.IntegerLiteral |
//                    +---------------------+
//                    | Token: INT          |
//                    | Value: 5            |
//                    +---------------------+

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// Wrap the expression inside parentheses to make operator precedence explicit
// Example: -5 + 3 -> (-5) + 3 makes it clear that the minus applies to only 5

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// Wraide parentheses to make operator precedence explicit
// Exa + 3 makes it clear that the minus applies to only 5

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type PostfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expressionNode() {}

func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

func (b *Boolean) String() string { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()

}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {}

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()

}

type FunctionLiteral struct {
	Token      token.Token // The 'funk' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("funk")
	out.WriteString("(")

	for i, param := range fl.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}

		out.WriteString(param.String())
	}

	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()

}

type CallExpression struct {
	Token     token.Token // The L.PAREN token
	Function  Expression  // Function identifiers are expressions too. Plus, this could either be an identifier or a function literal
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()

}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

func (sl *StringLiteral) String() string { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, e := range al.Elements {
		elements = append(elements, e.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type TernaryExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode() {}

func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }

func (te *TernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(te.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(te.Consequence.String())
	out.WriteString(" : ")
	out.WriteString(te.Alternative.String())
	out.WriteString(")")

	return out.String()

}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}

	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type MacroLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ml *MacroLiteral) expressionNode() {}

func (ml *MacroLiteral) TokenLiteral() string { return ml.Token.Literal }

func (ml *MacroLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range ml.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(ml.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(ml.Body.String())

	return out.String()
}
