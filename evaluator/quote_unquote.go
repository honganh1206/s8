package evaluator

import (
	"fmt"
	"s8/ast"
	"s8/object"
	"s8/token"
)

// Return an *object.Quote with an un-evaluated ast.Node
func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCall(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCall(quoted ast.Node, env *object.Environment) ast.Node {
	// Traverse every ast.Node inside the quoted argument
	// Punch holes into quote
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		// We can do environment-aware evaluation inside unquote calls
		unquoted := Eval(call.Arguments[0], env)
		return convertObjToASTNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpr, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return callExpr.Function.TokenLiteral() == "unquote"
}

// Create ast.Nodes that represent the passed in obj
func convertObjToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}
	case *object.Quote:
		// Use quote inside unquote
		// Preserve the quoted object
		return obj.Node
	default:
		// TODO: Ignore possible errors and just return nil here
		return nil
	}
}
