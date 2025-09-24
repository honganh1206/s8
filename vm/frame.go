package vm

import (
	"s8/code"
	"s8/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

// Return the instructions of the compiled function in a frame
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
