package ast

// Functions that change one node to another
type ModifierFunc func(Node) Node

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
	}
	// We REPLACE the node passed in as the argument with the node returned by the call
	// Important that we return instead of just modifying the given node so we can actually replace them
	// Also stop the recursion if we reach here (base case)
	return modifier(node)
}
