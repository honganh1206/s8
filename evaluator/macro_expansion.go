package evaluator

import (
	"s8/ast"
	"s8/object"
	"slices"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
			definitions = append(definitions, i) // Save the indexes to the macro defs
		}
	}

	// Remove the macro definitions from the AST
	// Loop from the end back
	// Only top-level macro definitions are allowed here
	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		definitionIndex := definitions[i]
		// Remove the macro definitions from the program.Statements
		program.Statements = slices.Delete(program.Statements, definitionIndex, definitionIndex+1)
	}
}

func isMacroDefinition(node ast.Statement) bool {
	letStmt, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStmt.Value.(*ast.MacroLiteral)
	if !ok {
		return false
	}

	return true
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStmt, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStmt.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}

	env.Set(letStmt.Name.Value, macro)
}

// Replace the macro calls with the result of their evaluation as generated code (AST nodes)
func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	// Transform the nodes to their quoted versions
	return ast.Modify(program, func(node ast.Node) ast.Node {
		// Call the identifier that uses macro
		callExpr, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		macro, ok := isMacroCall(callExpr, env)
		if !ok {
			return node
		}

		args := quoteArgs(callExpr)
		evalEnv := extendMacroEnv(macro, args)

		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}

		return quote.Node
	})
}

func isMacroCall(expr *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	ident, ok := expr.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(ident.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

func quoteArgs(expr *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}

	for _, a := range expr.Arguments {
		args = append(args, &object.Quote{Node: a})
	}

	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.NewEnclosedEnvironment(macro.Env)

	for paramIdx, param := range macro.Parameters {
		extended.Set(param.Value, args[paramIdx])
	}

	return extended
}
