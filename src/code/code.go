package code

import (
	"encoding/binary"
	"fmt"
)

// No Instruction byte singular
// Since working with []byte is easier and treat it implicitly as an Instruction
// Instead of encoding the Instruction type to Go's type system
type Instructions []byte

type Opcode byte // Exactly 1 byte wide

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
)

type Definition struct {
	Name          string
	OperandWidths []int // A slice telling how many bytes each operand takes up
}

// Similar to the precedence table, we will store operations like ADD, JUMP, etc. here
var definitions = map[Opcode]*Definition{
	// Two-byte wide operand maximum of 65536
	// That's more than enough. We won't be having more than 65536 references
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)] // Type casting
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{} // Risk of using empty byte slices using unknown opcode
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
		switch width {
		case 2:
			// Take the matching element in the operands slice and put it into the instruction
			// The 1st operand is put behind the opcode
			// Then the 2nd one is put behind the 1st one and so on
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		// Mark where to put the next operand
		offset += width
	}

	return instruction
}
