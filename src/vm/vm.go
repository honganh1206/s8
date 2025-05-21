package vm

import (
	"s8/src/code"
	"s8/src/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	sp           int // point to the next value. top of stack is stack[sp-1]
}
