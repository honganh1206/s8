package evaluator

import (
	"fmt"
	"math"
	"s8/src/ast"
	"s8/src/object"
	"s8/src/token"
)

// To NOT create new instances of object.Boolean or object.Null and use reference instead
// This improves performance too (pointer comparison is faster than value comparison)
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

const epsilon = 0.000001
const precision = 6

// Traverse the AST recursively
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: toFixed(node.Value, precision)}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node, right, env)
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
	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalPostfixExpression(node, left, env)
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
		if node.Function.TokenLiteral() == "quote" {
			return quote(node.Arguments[0], env)
		}
		// Get the function we want to call
		// Could be either *ast.Identifier (named functions) or *ast.FunctionLiteral (anonymous function)
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		// Stop the evaluation immediately and return the error
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.TernaryExpression:
		con := Eval(node.Condition, env)
		if isError(con) {
			return con
		}
		if isTruthy(con) {
			return Eval(node.Consequence, env)
		}
		return Eval(node.Alternative, env)
	case *ast.ArrayLiteral:
		elems := evalExpressions(node.Elements, env)
		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
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

func evalPrefixExpression(node *ast.PrefixExpression, right object.Object, env *object.Environment) object.Object {
	switch node.Operator {
	case token.BANG:
		return evalBangOperatorExpression(right)
	case token.MINUS:
		return evalMinusPrefixOperatorExpression(right)
	case token.TILDE:
		return evalBitwisePrefixNotOperatorExpression(right)
	case token.INCREMENT, token.DECREMENT:
		return evalIncreDecrePrefixOperatorExpression(node, right, env)
	default:
		return newError("unknown operator:%s%s", node.Operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, left object.Object, env *object.Environment) object.Object {
	if node.Operator == token.INCREMENT || node.Operator == token.DECREMENT {
		return evalIncreDecrePostfixOperatorExpression(node, left, env)
	}

	return newError("unknown operator:%s%s", node.Operator, left.Type())
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
	switch right.Type() {
	case object.INTERGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}

}

func evalBitwisePrefixNotOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTERGER_OBJ {
		return newError("bitwise operators not supported for type: %s", right.Type())
	}
	value := right.(*object.Integer).Value

	return &object.Integer{Value: ^value}
}

func evalIncreDecrePrefixOperatorExpression(node *ast.PrefixExpression, right object.Object, env *object.Environment) object.Object {
	ident, ok := node.Right.(*ast.Identifier)

	if !ok {
		return newError("cannot increment non-identifier: %s", right.Type())
	}

	val, ok := right.(*object.Integer)

	if !ok {
		return newError("cannot increment non-integer: %s", right.Type())
	}

	var newVal int64
	if node.Operator == token.INCREMENT {
		newVal = val.Value + 1
	} else if node.Operator == token.DECREMENT {
		newVal = val.Value - 1
	}

	returnVal := &object.Integer{Value: newVal}
	env.Set(ident.Value, returnVal)
	return returnVal
}

func evalIncreDecrePostfixOperatorExpression(node *ast.PostfixExpression, left object.Object, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)

	if !ok {
		return newError("cannot increment non-identifier: %s", left.Type())
	}

	val, ok := left.(*object.Integer)

	if !ok {
		return newError("cannot increment non-integer: %s", left.Type())
	}

	var newVal int64
	if node.Operator == token.INCREMENT {
		newVal = val.Value + 1
	} else if node.Operator == token.DECREMENT {
		newVal = val.Value - 1
	}

	originalVal := val.Value

	returnVal := &object.Integer{Value: originalVal}
	env.Set(ident.Value, &object.Integer{Value: newVal})
	return returnVal
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// We cannot use pointer comparison here
	// Since we are always allocating NEW instances of object.Integer
	case left.Type() == object.INTERGER_OBJ && right.Type() == object.INTERGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
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
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case ">>":
		return &object.Integer{Value: leftVal >> rightVal}
	case "<<":
		return &object.Integer{Value: leftVal << rightVal}
	case "|":
		return &object.Integer{Value: leftVal | rightVal}
	case "&":
		return &object.Integer{Value: leftVal & rightVal}
	case "^":
		return &object.Integer{Value: leftVal ^ rightVal}
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

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	// Round up to 6th decimal place
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: toFixed(leftVal+rightVal, precision)}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Float{Value: toFixed(leftVal/rightVal, precision)}
	case "<":
		return nativeBoolToBooleanObj(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObj(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObj(math.Abs(leftVal-rightVal) < epsilon)
	case "!=":
		return nativeBoolToBooleanObj(math.Abs(leftVal-rightVal) >= epsilon)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// Round a floating-point number to the nearest integer
//
//	Examples:
//
// - round(3.7) → 3.7 + 0.5 = 4.2 → 4
// - round(3.2) → 3.2 + 0.5 = 3.7 → 3
// - round(-3.7) → -3.7 - 0.5 = -4.2 → -4
func round(value float64) int {
	return int(value + math.Copysign(0.5, value))
}

// Round a floating-point number to a specified number of decimal places (precision)
func toFixed(value float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(value*output)) / output
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
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	// Check if it's a builtin function
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
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

// 	Function's Enclosed Env (Outer)
// ┌────────────────────────┐
// │ outer variables        │
// │                        │
// │   Extended Env (Inner) │
// │   ┌─────────────────┐  │
// │   │ function args   │  │
// │   │ local variables │  │
// │   └─────────────────┘  │
// └────────────────────────┘

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:

		extendedEnv := extendedFunctionEnv(fn, args)
		// We pass the extended env which EXTENDS (not replaces) the function's enclosed environment
		// This means the inner function can access values from its outer/enclosing environment
		evaluated := Eval(fn.Body, extendedEnv)

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:

		return newError("not a function: %s", fn.Type())
	}
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

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTERGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrObj.Elements) - 1)

	// Handle index out of bound
	if idx < 0 || idx > max {
		return NULL
	}

	return arrObj.Elements[idx]
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		// Populate the map
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}

	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as a hash key: %s", index.Type())
	}

	// Hash out the key and use it to retrieve the value from the hash map of the Go host language
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}
