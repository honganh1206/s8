package evaluator

import (
	"s8/object"
)

// Separate environment of builtin functions
var builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"rest":  object.GetBuiltinByName("rest"),
	"push":  object.GetBuiltinByName("push"),
	"puts":  object.GetBuiltinByName("puts"),
	"power": object.GetBuiltinByName("power"),
}
