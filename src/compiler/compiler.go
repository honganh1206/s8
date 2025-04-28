package compiler

import (
	"s8/src/code"
	"s8/src/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}
