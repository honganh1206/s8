package evaluator

import (
	"fmt"
	"s8/src/ast"
	"s8/src/object"
)

// To NOT create new instances of object.Boolean or object.Null and use reference instead
// This improves performance too (pointer comparison is faster than value comparison)
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Traverse the AST recursively
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		// Bind the value to the identifier
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		// Get the function we want to call
		// Could be either *ast.Identifier (named functions) or *ast.FunctionLiteral (anonymous function)
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

func evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range p.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			// Unwrap the value here
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

// This function will bubble up to the upper evalBlockStatement (in case of nested block statements)
// Or to evalProgram, then the ReturnValue will get wrapped
/*

 */
func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range bs.Statements {
		result = Eval(stmt, env)
		// The check is necessary for there might be statement types that are not handled
		// Plus result.Type() would cause a panic if result is nil
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				// Still wrap the value inside *object.ReturnValue or *object.Error
				return result
			}
		}
	}

	return result
}

// Take Go's native true and return singleton TRUE pointer
func nativeBoolToBooleanObj(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator:%s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	// If something has the SAME MEMORY ADDRESS value as boolean constants here, then it works
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTERGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// We cannot use pointer comparison here
	// Since we are always allocating NEW instances of object.Integer
	case left.Type() == object.INTERGER_OBJ && right.Type() == object.INTERGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
		// For cases like TRUE == TRUE
		// Here we use POINTER COMPARISON by comparing the memory addresses of two *object.Object pointers
		// The result is a native Go boolean
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case operator == "==":
		return nativeBoolToBooleanObj(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObj(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	// Group 1: Produce values of other types than booleans
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	// Group 2: Produce booleans as their results
	case "<":
		return nativeBoolToBooleanObj(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObj(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObj(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}

	return val
}

// Recursively evaluate each expression from left to right
func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	// If we make assertions about the order of the argument evaluation
	// We are on the conservative and safe side of programming design!
	for _, e := range exprs {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)

	if !ok {
		return newError("not a function: %s", function.Type())
	}

	extendedEnv := extendedFunctionEnv(function, args)
	// We pass the extending env instead of the current env of the function
	// So that the function can access values which are out of its scope
	// And thus using closures
	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
	// return evaluated // For testing
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	// Bind parameters with values inside the enclosed/inner environment
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

// This is critical, since we need to stop the evaluation of the LAST-CALLED function's body (Early return)
// Without this, evalBlockStatement will continue evaluating statements in outer functions
/*

// Example for the necessity of this function
let outer = fn() {
    let inner = fn() {
        return 5;  // This return should only exit the inner function
    };
    inner();
    return 10;  // This should be the final return value
};

*/
func unwrapReturnValue(obj object.Object) object.Object {
	// If there is a return statement inside an inner function, we unwrap it right away
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}
