package vm

import (
	"s8/code"
	"s8/object"
)

type Frame struct {
	cl *object.Closure
	ip int
	// A temp storage to keep track of the stack pointer value
	// before we execute a function call.
	// That way we can use this value to restore the stack pointer's value
	// when we are done executing the function call.
	// aka the frame pointer (Yes it's thing. Look up if you forget.)
	basePointer int
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{cl: cl, ip: -1, basePointer: basePointer}
}

// Return the instructions of the compiled function in a frame
func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
