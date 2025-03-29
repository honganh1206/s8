package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	// Helper functions so we do not have to construct integer literals again and again
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	// A function that modifies a Node
	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	// Outcome is a new Node type
	tests := []struct {
		input    Node
		expected Node
	}{
		{
			// Simply check if an input node is modified and returned
			one(),
			two(),
		},
		{
			// Walk the existing AST, pass each child node into the modifying function for modifications
			// This is how we find calls to unquote and replace it with a new AST node
			&Program{
				Statements: []Statement{
					&ExpressionStatement{
						Expression: one(),
					},
				},
			},
			&Program{
				Statements: []Statement{
					&ExpressionStatement{
						Expression: two(),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)
		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got: %#v, want: %#v", modified, tt.expected)
		}
	}

}
