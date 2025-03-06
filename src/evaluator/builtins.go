package evaluator

import (
	"fmt"
	"math"
	"s8/src/object"
)

// Separate environment of builtin functions
var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got: %d, want: 1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported. got: %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got: %d, want: 1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be of ARRAY type, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got: %d, want: 1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be of ARRAY type, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got: %d, want: 1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be of ARRAY type, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElems := make([]object.Object, length-1, length-1)
				// Exclude the 1st elem
				copy(newElems, arr.Elements[1:length])
				return &object.Array{Elements: newElems}
			}

			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got: %d, want: 2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be of ARRAY type, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElems := make([]object.Object, length+1, length+1)
			copy(newElems, arr.Elements)
			newElems[length] = args[1]
			return &object.Array{Elements: newElems}
		},
	},
	// Print given args to STDOUT
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"power": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for power. got: %d, want: 2", len(args))
			}
			if args[0].Type() != object.INTERGER_OBJ || args[1].Type() != object.INTERGER_OBJ {
				return newError("argument to power must of of INTEGER type, got: %s (1st argument) | %s (2nd argument)", args[0].Type(), args[1].Type())
			}
			base := args[0].(*object.Integer)
			exponent := args[1].(*object.Integer)

			result := math.Pow(float64(base.Value), float64(exponent.Value))
			// Check if the result has no decimal part
			if result == float64(int64(result)) {
				return &object.Integer{Value: int64(result)}
			}
			return &object.Float{Value: result}
		},
		// TODO: Add round and format
	},
}
