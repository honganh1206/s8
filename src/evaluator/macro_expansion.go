package evaluator

import (
	"s8/src/ast"
	"s8/src/object"
	"slices"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, stmt := range program.Statements {
		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
			definitions = append(definitions, i) // Save the indexes to the macro defs
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
