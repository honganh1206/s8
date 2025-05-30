package ast

// Functions that change one node to another
type ModifierFunc func(Node) Node

// TODO: Replace underscores
// This does not update the Token fields of the parent nodes
func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		// Walk down the tree just like how Eval does it
		// Program has children, and we call modify with each child
		// And again we have calls to ast.Modify to the children of children
		// This implements a depth-first traversal of the AST
		for i, stmt := range node.Statements {
			node.Statements[i], _ = Modify(stmt, modifier).(Statement)
		}
	case *ExpressionStatement:
		// Modify the expression within the statement
		// The expression itself might be complex and contain other nodes
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PostfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		for i := range node.Statements {
			node.Statements[i], _ = Modify(node.Statements[i], modifier).(Statement)
		}
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i := range node.Elements {
			node.Elements[i], _ = Modify(node.Elements[i], modifier).(Expression)
		}
	case *HashLiteral:
		// We need to iterate over the map and modify both the keys and the values
		newPairs := make(map[Expression]Expression)
		for k, v := range node.Pairs {
			newKey, _ := Modify(k, modifier).(Expression)
			newVal, _ := Modify(v, modifier).(Expression)
			newPairs[newKey] = newVal
		}
		node.Pairs = newPairs
	}
	// We REPLACE the node passed in as the argument with the node returned by the call
	// Important that we return instead of just modifying the given node so we can actually replace them
	// Also stop the recursion if we reach here (base case)
	return modifier(node)
}
