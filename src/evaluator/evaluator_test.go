package evaluator

import (
	"s8/src/lexer"
	"s8/src/object"
	"s8/src/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"~5", -6}, // Bitwise NOT
		{"5 >> 1", 2},
		{"5 << 1", 10},
		{"5 | 5", 5},
		{"5 & 5", 5},
		{"5 ^ 5", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestPrefixIncrementExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = 5; x++", 5},    // Returns original value
		{"let x = 5; x--", 5},    // Returns original value
		{"let x = 5; x++; x", 6}, // Check value after increment
		{"let x = 5; x--; x", 4}, // Check value after decrement

		{"let x = 5; x++; x++", 6}, // Second ++ should return 6
		{"let x = 5; x++; x--", 6}, // Returns 6 then decrements

		// Arithmetic with postfix
		{"let x = 5; let y = 3; x++ + y", 8},   // 5 + 3, then x becomes 6
		{"let x = 5; let y = 3; x++ + y++", 8}, // 5 + 3, then x and y increment

		// Complex expressions
		{"let x = 5; (x++) + 2", 7},
		{"let x = 5; let y = x++; y", 5}, // Assignment captures original value
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.0", 5.000000},
		{"10.5", 10.500000},
		{"-5.25", -5.250000},
		{"-10.75", -10.750000},
		{"5.5 + 5.25 + 5.125 + 5.625 - 10.5", 11.000000},
		{"2.5 * 2.5", 6.250000},
		{"-50.5 + 100.25 + -50.25", -0.500000},
		{"5.5 * 2.5 + 10.25", 24.000000},
		{"5.25 + 2.5 * 10.125", 30.562500},
		{"20.5 + 2.25 * -10.5", -3.125000},
		{"50.5 / 2.5 * 2.25 + 10.125", 55.575},
		{"2.5 * (5.25 + 10.125)", 38.437500},
		{"3.5 * 3.5 * 3.5 + 10.125", 53.000000},
		{"3.25 * (3.25 * 3.25) + 10.5", 44.828125},
		{"(5.5 + 10.25 * 2.5 + 15.75 / 3.25) * 2.5 + -10.25", 79.677885},
		{"1.23456789", 1.234568},    // Testing rounding
		{"0.333333333", 0.333333},   // Testing precision limit
		{"1.0 / 3.0", 0.333333},     // Division resulting in repeating decimal
		{"0.1 + 0.2", 0.300000},     // Common floating-point precision test
		{"3.14159265359", 3.141593}, // Pi rounded to 6 decimal places
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},

		{"!5", false},
		{"!!true", true},

		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}

	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
if (10 > 1) {
if (10 > 1) {
return 10;
}
return 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
if (10 > 1) {
return true + false;
}
return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		// Only support for + and nothing more
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[funk(x) { x }];`,
			"unusable as a hash key: FUNCTION",
		},
		{
			"5++",
			"cannot apply postfix operator to literal",
		},
		{
			"x++ ++",
			"cannot apply multiple postfix operators",
		},
		{
			"(x++)++",
			"cannot apply postfix operator to postfix expression",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("no object error returned. got: %T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message, expected: %q, got: %q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "funk(x) { x + 2 };"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function. got: %T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got: %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = funk(x) { x; }; identity(5);", 5},
		{"let identity = funk(x) { return x; }; identity(5);", 5},
		{"let double = funk(x) { x * 2; }; double(5);", 10},
		{"let add = funk(x, y) { x + y; }; add(5, 5);", 10},
		// Two forms of *ast.CallExpression
		{"let add = funk(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20}, // Function as an identifier evaluating to a function obj
		{"funk(x) { x; }(5)", 5}, // Function is a function literal aka Anonymous function
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = funk(x) {
funk(y) { x + y };
};
let addTwo = newAdder(2);
addTwo(2);`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`let arr = [1, 2 ,3]; len(arr)`, 3},
		{`let arr = [1, 2 ,3]; first(arr)`, 1},
		{`let arr = [1, 2 ,3]; last(arr)`, 3},
		{`let arr = [1, 2 ,3]; rest(arr)`, [2]int{2, 3}},
		{`let arr = [1, 2 ,3]; push(arr, 4)`, [4]int{1, 2, 3, 4}},
		{`power(2, 3)`, 8},
		{`len(1)`, "argument to `len` not supported. got: INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got: 2, want: 1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got: %T(%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected: %q, got: %q", expected, errObj.Message)
			}
		}

	}
}

func TestTernaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Basic ternary operations
		{"true ? 10 : 20", 10},
		{"false ? 10 : 20", 20},

		// Nested ternary operations
		{"true ? (false ? 10 : 20) : 30", 20},
		{"false ? 10 : (true ? 20 : 30)", 20},

		// With expressions as condition
		{"5 > 3 ? 10 : 20", 10},
		{"5 < 3 ? 10 : 20", 20},

		// With expressions as consequences and alternatives
		{"true ? 5 + 5 : 20", 10},
		{"false ? 10 : 15 + 5", 20},

		// With string operations
		{`true ? "yes" : "no"`, "yes"},
		{`false ? "yes" : "no"`, "no"},

		// With identifiers
		{"let a = 5; let b = 10; a > b ? a : b", 10},
		{"let a = 15; let b = 10; a > b ? a : b", 15},

		// With function calls
		{"let f = funk(x) { x * 2 }; true ? f(5) : 20", 10},
		{"let f = funk(x) { x * 2 }; false ? 20 : f(5)", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("String has wrong value. got=%q, want=%q", str.Value, expected)
			}
		}
	}
}

// Also add error handling tests for ternary expressions
func TestTernaryErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + (true ? true + false : 1)",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true ? 1 + true : 2",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"false ? 1 : true + 2",
			"type mismatch: BOOLEAN + INTEGER",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
{
"one": 10 - 9,
two: 1 + 1,
"thr" + "ee": 6 / 2,
4: 4,
true: 5,
false: 6
}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval did not return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	// HashKey is an interface here
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("hash has wrong num of pairs. got=%d",
			len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		// Ensure the HashKey method implemented by different data types are called correctly
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

/*
HELPER FUNCTIONS
*/

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		// %+v is the enhanced version of %v which prints the default representation
		t.Errorf("object is not Integer. got: %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got: %d, want: %d", result.Value, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		// %+v is the enhanced version of %v which prints the default representation
		t.Errorf("object is not Float. got: %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got: %f, want: %f", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		// %+v is the enhanced version of %v which prints the default representation
		t.Errorf("object is not Boolean. got: %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got: %t, want: %t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got: %T (%+v)", obj, obj)
		return false
	}
	return true
}
