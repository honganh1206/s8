package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// A slice of bytes including pointers to opcode and operands (stored in constant pool).
// Why no Instruction byte singular?
// Working with []byte is easier and we can treat it implicitly as an Instruction,
// instead of encoding the Instruction type to Go's type system
type Instructions []byte

// Exactly 1 byte wide
type Opcode byte

const (
	// OpConstant instructs the VM to:
	// 1. Retrieve a constant from the constant pool using the operand as an index.
	// 2. Push the retrieved constant onto the stack.
	// Supposed we have a value 123 and we compile it
	// The compiler then generates the byte code saying:
	// "Use OpConstant and attach operand 7, because 123 is stored at constant pool index 7."
	// The bytecode might look like this:  [OpConstant, 7]
	OpConstant Opcode = iota
	// Each definition later on will have `Op` prefix with the value it refers to determined by iota
	OpTrue
	OpFalse

	// Binary operators
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpAmpersand
	OpPipe
	OpExponent
	OpRShift
	OpLShift

	// Unary operators
	OpMinus
	OpBang
	OpTilde
	OpPreInc
	OpPreDec
	OpPostInc
	OpPostDec

	// Conditional
	OpJumpNotTruthy
	// OpJumpNotTruthy must have an offset
	// pointing to instructions right after BlockStatement of Consequence.
	OpJump
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpGetLocal
	OpSetLocal
	OpArray
	OpHash
	OpIndex

	// Functions
	OpCall        // Tell the VM to start executing *object.CompiledFunction
	OpReturnValue // Return value must be on top of the stack
	OpReturn      // Return no value, resume to parent execution
	OpGetBuiltin
	OpClosure
	OpGetFree        // Get free variables
	OpCurrentClosure // Load the closure it's executing on to the stack (to execute recursive function)
)

// How an instruction looks like
type Definition struct {
	// Opcode
	Name string
	// Widths of its operands
	OperandWidths []int
}

// Similar to the precedence table, we will store operations like ADD, JUMP, etc. here
var definitions = map[Opcode]*Definition{
	// Push constant to top of the stack
	// Two-byte wide operand maximum of 65536 is more than enough.
	// We won't be having more than 65536 references aka values that exceed 65536.
	OpConstant: {"OpConstant", []int{2}},
	// No operand
	OpAdd:         {"OpAdd", []int{}},
	OpPop:         {"OpPop", []int{}},
	OpSub:         {"OpSub", []int{}},
	OpMul:         {"OpMul", []int{}},
	OpDiv:         {"OpDiv", []int{}},
	OpTrue:        {"OpTrue", []int{}},
	OpFalse:       {"OpFalse", []int{}},
	OpEqual:       {"OpEqual", []int{}},
	OpNotEqual:    {"OpNotEqual", []int{}},
	OpGreaterThan: {"OpGreaterThan", []int{}},
	OpAmpersand:   {"OpAmpersand", []int{}},
	OpMinus:       {"OpMinus", []int{}},
	OpBang:        {"OpBang", []int{}},
	OpTilde:       {"OpTilde", []int{}},
	OpPipe:        {"OpPipe", []int{}},
	OpRShift:      {"OpRShift", []int{}},
	OpLShift:      {"OpLShift", []int{}},
	OpExponent:    {"OpExponent", []int{}},
	OpPreInc:      {"OpPreInc", []int{}},
	OpPreDec:      {"OpPreDec", []int{}},
	OpPostInc:     {"OpPostInc", []int{}},
	OpPostDec:     {"OpPostDec", []int{}},
	// We currently do absolute jump here
	// The operand is the index of the instruction
	OpJump:          {"OpJump", []int{2}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpSetLocal:      {"OpSetLocal", []int{1}},
	// One-byte operands
	OpGetLocal: {"OpGetLocal", []int{1}},
	// Operand is number of values in an array
	OpArray: {"OpArray", []int{2}},
	// Operand is number of values x2 in a hash
	OpHash:        {"OpHash", []int{2}},
	OpIndex:       {"OpIndex", []int{}},
	OpCall:        {"OpCall", []int{1}},
	OpReturnValue: {"OpReturnValue", []int{}},
	OpReturn:      {"OpReturn", []int{}},
	// Define up to 256 builtin functions
	OpGetBuiltin: {"OpGetBuiltin", []int{1}},
	// Two operands, what do we have here?
	// The 1st one is the constant index, specifying where in the constant pool we can find *object.CompiledFunction.
	// The two bytes are to ensure we can find the compiled function without worrying the index is too high.
	// The 2nd operand specifies how many free variables sit on the stack and to-be-transferred to the closure. 1 byte (256 free variables) should be enough?
	OpClosure:        {"OpClosure", []int{2, 1}},
	OpGetFree:        {"OpGetFree", []int{1}},
	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)] // Type casting
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// Encode instructions to slice of bytes
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		// Risk of using empty byte slices if using unknown opcode
		return []byte{}
	}

	instructionLen := 1 // 1 byte for the opcode

	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1

	for i, o := range operands {
		width := def.OperandWidths[i]
		// Encode the operands into the instruction as bytes
		switch width {
		case 2:
			// Take the width-matching element in the operands slice and put it into the instruction
			// The 1st operand is put behind the opcode
			// Then the 2nd one is put behind the 1st one and so on
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		// Mark where to put the next operand
		offset += width
	}

	return instruction
}

// Decode the operands of an instruction
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		// Retrieve and decode the operand at that position/offset
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}

// String-tify the instructions
// Implementing the Stringer interface here
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR:%s\n", err)
			continue
		}

		// Skip the opcode and read the operands
		operands, offset := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + offset

	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	// Mismatching number of operands?
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1: // Case for OpConstant
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
