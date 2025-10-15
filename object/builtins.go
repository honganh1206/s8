package object

import (
	"fmt"
	"math"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got: %d, want: 1", len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported. got: %s", args[0].Type())
			}
		},
		},
	},
	{
		// Print given args to STDOUT
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return nil
		},
		},
	},
	{

		"first",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got: %d, want: 1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `first` must be of ARRAY type, got %s", args[0].Type())
			}
			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return nil
		},
		},
	},
	{
		"last",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got: %d, want: 1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to `last` must be of ARRAY type, got %s", args[0].Type())
			}
			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return nil
		},
		},
	},
	{
		// Exclude the 1st elem
		"rest",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got: %d, want: 1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `rest` must be of ARRAY type, got %s", args[0].Type())
				}
				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					newElems := make([]Object, length-1)
					copy(newElems, arr.Elements[1:length])
					return &Array{Elements: newElems}
				}

				return nil
			},
		},
	},
	{
		"push",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got: %d, want: 2", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `push` must be of ARRAY type, got %s", args[0].Type())
				}
				arr := args[0].(*Array)
				length := len(arr.Elements)

				// Leave the old array untouched and allocate a new one
				newElems := make([]Object, length+1)
				copy(newElems, arr.Elements)
				newElems[length] = args[1]
				return &Array{Elements: newElems}
			},
		},
	},
	{
		"power",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments for power. got: %d, want: 2", len(args))
				}
				if args[0].Type() != INTERGER_OBJ || args[1].Type() != INTERGER_OBJ {
					return newError("argument to power must of of INTEGER type, got: %s (1st argument) | %s (2nd argument)", args[0].Type(), args[1].Type())
				}
				base := args[0].(*Integer)
				exponent := args[1].(*Integer)

				result := math.Pow(float64(base.Value), float64(exponent.Value))
				// Check if the result has no decimal part
				if result == float64(int64(result)) {
					return &Integer{Value: int64(result)}
				}
				return &Float{Value: result}
			},
			// TODO: Add round and format
		},
	},
}

func newError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
