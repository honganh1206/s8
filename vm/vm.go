package vm

import (
	"fmt"
	"s8/code"
	"s8/compiler"
	"s8/object"
)

const StackSize = 2048

// Since each operand is 16 bit-wide,
// we have an upper limit on the number of global bindings
// our VM can support
const GlobalSize = 65536
const MaxFrames = 1024

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	// constant pool
	constants []object.Object
	// The stack elements are numbered/accessed from the bottom up, with index 0 being the bottom. The topmost element is at index sp-1, the second from top at sp-2, etc.
	stack []object.Object
	// stackpointer points to the next free slot in the stack. top of stack is stack[sp-1]
	sp          int
	globals     []object.Object
	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	// Pre-allocate the frames slice with the main frame,
	// now we don't need to initialize the instructions in the VM struct.
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants:   bytecode.Constants,
		stack:       make([]object.Object, StackSize),
		sp:          0,
		globals:     make([]object.Object, GlobalSize),
		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
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
	var ip int
	var ins code.Instructions
	var op code.Opcode
	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		// Helper variables to handle instructions from the current frame.
		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])

		switch op {
		case code.OpConstant:
			// Decode the pointer to the operand right after the opcode
			// We do not use ReadOperands() here for number-of-param reason (and performance?)
			constIndex := code.ReadUint16(ins[ip+1:])
			// Pointing to the NEXT opcode, not an operand
			vm.currentFrame().ip += 2

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
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			// Again pointing to the next opcode
			// just like how OpConstant does that
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				// Set the instruction pointer right before the target instruction
				// then let the for-loop do its work
				vm.currentFrame().ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			// Again, move the pointer to the operand after the opcode
			globalIndex := code.ReadUint16(ins[ip+1:])

			// Increment two bytes for the next instruction
			vm.currentFrame().ip += 2

			vm.globals[globalIndex] = vm.pop()
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			// Operand for local bindings is only 1-byte wide.
			vm.currentFrame().ip += 1

			// Get the frame of the callee function
			frame := vm.currentFrame()
			// Save the binding to the stack frame
			// using the base pointer and the index of the binding as an offset
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()
			err := vm.push(vm.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}
		case code.OpArray:
			numElems := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			// Elements might be scattered around the stack,
			// not necessarily at the top of the stack going down
			// N == endIndex - startIndex
			array := vm.buildArray(vm.sp-numElems, vm.sp)
			// Move the cursor to the start index
			vm.sp = vm.sp - numElems

			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHash:
			numElems := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			hash, err := vm.buildHash(vm.sp-numElems, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElems

			err = vm.push(hash)
			if err != nil {
				return err
			}
		case code.OpIndex:
			// Tip: Look at the tests in compiler_test.go and read from the bottom up.
			// This is equal to from the stack top down
			index := vm.pop()
			left := vm.pop() // Probably an array or hash literal

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}
		case code.OpCall:
			// Arguments now sit on top of function object on the stack
			numArgs := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			err := vm.callFunction(int(numArgs))
			if err != nil {
				return err
			}
		case code.OpReturnValue:
			returnValue := vm.pop()

			// Take the frame for the function call off the stack
			frame := vm.popFrame()
			// At this point the base pointer is pointing to the just-executed function,
			// so when we pop the frame of the function off the stack, we reset the stack pointer as well.
			// It's also an optimization: When we get rid off the local bindings, we leave the just-executed function on the stack.
			// So instead of another vm.pop() to pop the function object off the stack, we decrement vm.sp even further.
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				return err
			}
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

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

	switch {
	case rightType == object.INTERGER_OBJ && leftType == object.INTERGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, left, right)
	case rightType == object.STRING_OBJ && leftType == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s",
			leftType, rightType)
	}

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

func (vm *VM) executeBinaryStringOperation(
	op code.Opcode,
	left, right object.Object,
) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	return vm.push(&object.String{Value: leftValue + rightValue})
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

func (vm *VM) executePostfixIncrementDecrementOperator(_ code.Opcode, val int64) error {
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

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elems := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		// Start with the index 0 of the extracted stack chunk
		elems[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elems}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)
	// K-V so we increment by 2
	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		// The value is right after the key :)
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}
		hashedPairs[hashKey.HashKey()] = pair
	}
	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTERGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(array, index object.Object) error {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	// Start from 0
	max := int64(len(arrayObject.Elements) - 1)
	// Return null if out of bounds
	if i < 0 || i > max {
		return vm.push(Null)
	}
	return vm.push(arrayObject.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index object.Object) error {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}
	return vm.push(pair.Value)
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

// Push a frame to the stack frame
func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

// Pop a frame off a stack frame
func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func (vm *VM) callFunction(numArgs int) error {
	// At the moment the stack pointer is pointing to the next free slot (arrays start at 0 remember?)
	// so we need to access the function by -1 to avoid accessing empty/undefined stack location
	fn, ok := vm.stack[vm.sp-1-int(numArgs)].(*object.CompiledFunction)
	if !ok {
		return fmt.Errorf("calling non-function")
	}

	if numArgs != fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			fn.NumParameters, numArgs)
	}
	// Store the current stack pointer as the base/frame pointer
	// so we know somewhere to resume when we are done with the function call.
	// We also need to subtract the argument indexes so the base pointer does not point to empty stack slots.
	frame := NewFrame(fn, vm.sp-numArgs)
	vm.pushFrame(frame)
	// Create a "hole" - memory region of the stack for the local bindings of the OpCall being executed
	vm.sp = frame.basePointer + fn.NumLocals

	return nil
}
