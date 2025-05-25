package vm

import (
	"fmt"
	"s8/src/code"
	"s8/src/compiler"
	"s8/src/object"
)

const StackSize = 2048

type VM struct {
	// constant pool
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	// stackpointer points to the next free slot in the stack. top of stack is stack[sp-1]
	sp int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

// Get the element at the top of the VM's stack
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	// Increase the instruction pointer and fetch the current instruction
	// Why not use code.Lookup()? Because then we have to move the byte to here and there
	// then look up the opcode definition, return it and take it apart
	// That's a lot more work!
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			// Decode the pointer to the operand right after the opcode
			// We do not use ReadOperands() here for number-of-param reason (and performance?)
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			// Pointing to the NEXT opcode, not an operand
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftVal := left.(*object.Integer).Value
			rightVal := right.(*object.Integer).Value

			result := leftVal + rightVal
			vm.push(&object.Integer{Value: result})
		}

	}
	return nil
}

// Push the object to the stack's constant pool
func (vm *VM) push(o object.Object) error {
	if vm.sp > StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

// Remove element at the top of the VM's stack
func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp-- // Move the pointer to the next element in line
	return o
}
