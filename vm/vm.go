package vm

import (
	"fmt"
	"s8/code"
	"s8/compiler"
	"s8/object"
)

const StackSize = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

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

// Get the topmost element in the VM's stack right before we pop it off
func (vm *VM) LastPoppedStackElement() object.Object {
	return vm.stack[vm.sp]
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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpPipe, code.OpRShift, code.OpLShift, code.OpAmpersand, code.OpExponent:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpMinus, code.OpTilde, code.OpPreInc, code.OpPreDec, code.OpPostInc, code.OpPostDec:
			err := vm.executeUnaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpJump:
			// Set the instruction pointer (ip)
			// to right before the instruction we want to execute
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			// Again pointing to the next opcode
			// just like how OpConstant does that
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				// Set the instruction pointer right before the target instruction
				// then let the for-loop do its work
				ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	if rightType == object.INTERGER_OBJ && leftType == object.INTERGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s",
		leftType, rightType)

}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	case code.OpPipe:
		result = leftValue | rightValue
	case code.OpRShift:
		result = leftValue >> rightValue
	case code.OpLShift:
		result = leftValue << rightValue
	case code.OpAmpersand:
		result = leftValue & rightValue
	case code.OpExponent:
		result = leftValue ^ rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	if leftType == object.INTERGER_OBJ && rightType == object.INTERGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	// Comparing boolean objects like true == false
	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)",
			op, left.Type(), right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeUnaryOperation(op code.Opcode) error {
	operand := vm.pop()
	if operand.Type() != object.INTERGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	val := operand.(*object.Integer).Value
	switch op {
	case code.OpMinus:
		return vm.push(&object.Integer{Value: -val})
	case code.OpTilde:
		return vm.push(&object.Integer{Value: ^val})
	case code.OpPreInc, code.OpPreDec:
		return vm.executePrefixIncrementDecrementOperator(op, val)
	case code.OpPostInc, code.OpPostDec:
		return vm.executePostfixIncrementDecrementOperator(op, val)
	}

	return nil

}

// Push the object from the constant pool to the stack
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

func (vm *VM) executePrefixIncrementDecrementOperator(op code.Opcode, val int64) error {
	var newVal int64
	switch op {
	case code.OpPreInc:
		newVal = val + 1
	case code.OpPreDec:
		newVal = val - 1
	}
	return vm.push(&object.Integer{Value: newVal})
}

func (vm *VM) executePostfixIncrementDecrementOperator(op code.Opcode, val int64) error {
	// TODO: Set to environment later
	return vm.push(&object.Integer{Value: val})
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
